package convention_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"mime/multipart"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apiform"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apiquery"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/pagination"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
)

// Fixed upstream provenance: anthropics/anthropic-sdk-go tag v1.74.0,
// commit 3cb26e4450dc5294618ab9cc40377e5dd79ac413.

func TestConformanceDefaultResponseHeaderTimeout(t *testing.T) {
	for name, options := range map[string][]option.RequestOption{
		"forward": forward.NewClient(option.WithPAT("test-token")).Options,
		"managed": managed.NewClient(option.WithPAT("test-token")).Options,
	} {
		t.Run(name, func(t *testing.T) {
			cfg, err := convention.NewRequestConfig(context.Background(), http.MethodGet, "models", nil, nil, options...)
			if err != nil {
				t.Fatal(err)
			}
			transport, ok := cfg.HTTPClient.Transport.(*http.Transport)
			if !ok || transport == http.DefaultTransport {
				t.Fatal("client must have its own transport without modifying http.DefaultTransport")
			}
			if transport.ResponseHeaderTimeout != 10*time.Minute {
				t.Fatalf("response header timeout = %v", transport.ResponseHeaderTimeout)
			}
			if cfg.HTTPClient.Timeout != 0 || cfg.RequestTimeout != 0 {
				t.Fatal("default timeout must not limit the response body or total stream duration")
			}
			next, err := convention.NewRequestConfig(context.Background(), http.MethodGet, "models", nil, nil, options...)
			if err != nil || next.HTTPClient != cfg.HTTPClient {
				t.Fatalf("requests must reuse the client's connection pool: %v", err)
			}
			custom := &http.Client{Timeout: time.Second}
			overridden, err := convention.NewRequestConfig(context.Background(), http.MethodGet, "models", nil, nil, append(options, option.WithHTTPClient(custom))...)
			if err != nil || overridden.HTTPClient != custom {
				t.Fatalf("explicit HTTP client must take precedence: %v", err)
			}
		})
	}
}

func TestConformanceHeaderTimeoutAndStreamLifetime(t *testing.T) {
	const timeout = 100 * time.Millisecond
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/headers" {
			<-r.Context().Done()
			return
		}
		w.Header().Set("Content-Type", "text/event-stream")
		w.WriteHeader(http.StatusOK)
		w.(http.Flusher).Flush()
		if r.URL.Path == "/request-timeout" {
			<-r.Context().Done()
			return
		}
		time.Sleep(2 * timeout)
		_, _ = io.WriteString(w, "data: {}\n\n")
	}))
	defer server.Close()
	client := convention.DefaultHTTPClient()
	defer client.CloseIdleConnections()
	client.Transport.(*http.Transport).ResponseHeaderTimeout = timeout
	for _, path := range []string{"headers", "stream", "request-timeout"} {
		t.Run(path, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			opts := []option.RequestOption{option.WithBaseURL(server.URL), option.WithHTTPClient(client), option.WithMaxRetries(0)}
			if path == "request-timeout" {
				opts = append(opts, option.WithRequestTimeout(timeout))
			}
			var response *http.Response
			err := convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &response, opts...)
			var body []byte
			if err == nil {
				defer response.Body.Close()
				body, err = io.ReadAll(response.Body)
			}
			if path == "stream" {
				if err != nil || string(body) != "data: {}\n\n" {
					t.Fatalf("response header timeout interrupted stream body: %q, %v", body, err)
				}
			} else {
				var timeoutError net.Error
				if !errors.As(err, &timeoutError) || !timeoutError.Timeout() || ctx.Err() != nil {
					t.Fatalf("expected transport/request timeout before caller deadline, got %v", err)
				}
			}
		})
	}
}

func TestConformanceWrappedDefaultTransport(t *testing.T) {
	previous := http.DefaultTransport
	t.Cleanup(func() { http.DefaultTransport = previous })
	calls := 0
	http.DefaultTransport = roundTripFunc(func(*http.Request) (*http.Response, error) {
		calls++
		return jsonResponse(http.StatusOK, `{}`), nil
	})
	client := convention.DefaultHTTPClient()
	request, err := http.NewRequest(http.MethodGet, "https://sdk.test/resource", nil)
	if err != nil {
		t.Fatal(err)
	}
	response, err := client.Do(request)
	if err != nil {
		t.Fatal(err)
	}
	defer response.Body.Close()
	if calls != 1 || client.Timeout != 0 {
		t.Fatal("custom default transport must be preserved without a total timeout")
	}
}

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func jsonResponse(status int, body string) *http.Response {
	return &http.Response{
		StatusCode: status,
		Header:     http.Header{"Content-Type": {"application/json"}},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

type bodyParams struct {
	Title    param.Opt[string] `json:"title,omitzero"`
	Enabled  param.Opt[bool]   `json:"enabled,omitzero"`
	Metadata map[string]any    `json:"metadata,omitzero"`
	paramObj
}

type paramObj = param.APIObject

func (p bodyParams) MarshalJSON() ([]byte, error) {
	type shadow bodyParams
	return param.MarshalObject(p, (*shadow)(&p))
}

func TestConformanceBodyParamOmitNullValue(t *testing.T) {
	omitted, err := json.Marshal(bodyParams{})
	if err != nil {
		t.Fatalf("marshal omitted: %v", err)
	}
	if string(omitted) != "{}" {
		t.Fatalf("unset optional fields must be omitted, got %s", omitted)
	}
	explicit, err := json.Marshal(bodyParams{
		Title:    param.Null[string](),
		Enabled:  param.NewOpt(false),
		Metadata: map[string]any{},
	})
	if err != nil {
		t.Fatalf("marshal explicit: %v", err)
	}
	var got map[string]any
	if err := json.Unmarshal(explicit, &got); err != nil {
		t.Fatalf("decode explicit: %v", err)
	}
	if v, ok := got["title"]; !ok || v != nil {
		t.Fatalf("explicit null must serialize as JSON null, got %v (present=%v)", v, ok)
	}
	if v, ok := got["enabled"]; !ok || v != false {
		t.Fatalf("explicit false must be retained, got %v (present=%v)", v, ok)
	}
	if _, ok := got["metadata"]; !ok {
		t.Fatalf("explicit empty map must be retained")
	}
}

type queryParamsT struct {
	Limit          param.Opt[int64]     `query:"limit,omitzero"`
	Status         param.Opt[string]    `query:"status,omitzero"`
	IncludeArchive param.Opt[bool]      `query:"include_archived,omitzero"`
	CreatedAtGt    param.Opt[time.Time] `query:"created_at[gt],omitzero" format:"date-time"`
	Order          param.Opt[string]    `query:"order,omitzero"`
	Tags           []string             `query:"tags,omitzero"`
	paramObj
}

func (q queryParamsT) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(q, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

func TestConformanceQueryEncoding(t *testing.T) {
	values, err := queryParamsT{
		Limit:          param.NewOpt[int64](0),
		Status:         param.NewOpt("active"),
		IncludeArchive: param.NewOpt(false),
		CreatedAtGt:    param.NewOpt(time.Date(2026, 1, 2, 3, 4, 5, 0, time.UTC)),
		Tags:           []string{"a", "b"},
	}.URLQuery()
	if err != nil {
		t.Fatalf("marshal query: %v", err)
	}
	if values.Get("limit") != "0" {
		t.Fatalf("zero scalar must encode, got %q", values.Get("limit"))
	}
	if values.Get("status") != "active" {
		t.Fatalf("scalar string mismatch, got %q", values.Get("status"))
	}
	if values.Get("include_archived") != "false" {
		t.Fatalf("false scalar must encode, got %q", values.Get("include_archived"))
	}
	if values.Get("created_at[gt]") != "2026-01-02T03:04:05Z" {
		t.Fatalf("date-time format mismatch, got %q", values.Get("created_at[gt]"))
	}
	if tags := values["tags"]; len(tags) != 2 || tags[0] != "a" || tags[1] != "b" {
		t.Fatalf("repeated array format mismatch, got %v", tags)
	}
	if _, ok := values["order"]; ok {
		t.Fatalf("unset optional query field must be omitted")
	}
}

type formParams struct {
	Name    string    `json:"name"`
	Content io.Reader `json:"content" format:"binary"`
	paramObj
}

func (p formParams) MarshalMultipart() ([]byte, string, error) {
	buf := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buf)
	if err := apiform.MarshalRoot(p, writer); err != nil {
		return nil, "", err
	}
	if err := writer.Close(); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), writer.FormDataContentType(), nil
}

func TestConformanceMultipartEncoding(t *testing.T) {
	raw, contentType, err := formParams{Name: "asset", Content: strings.NewReader("hello")}.MarshalMultipart()
	if err != nil {
		t.Fatalf("marshal multipart: %v", err)
	}
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil || mediaType != "multipart/form-data" {
		t.Fatalf("unexpected content type %q (err=%v)", contentType, err)
	}
	reader := multipart.NewReader(bytes.NewReader(raw), params["boundary"])
	fields := map[string]string{}
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("read part: %v", err)
		}
		data, _ := io.ReadAll(part)
		fields[part.FormName()] = string(data)
	}
	if fields["name"] != "asset" {
		t.Fatalf("scalar form field mismatch, got %q", fields["name"])
	}
	if fields["content"] != "hello" {
		t.Fatalf("binary form field mismatch, got %q", fields["content"])
	}
}

func TestConformanceRequestConfigContract(t *testing.T) {
	var captured *http.Request
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		captured = r
		return jsonResponse(200, `{"id":"one"}`), nil
	})}
	var dst map[string]any
	err := convention.ExecuteNewRequest(
		context.Background(), http.MethodGet, "resource", nil, &dst,
		option.WithBaseURL("https://sdk.test/api/"),
		option.WithMaxRetries(0),
		option.WithHTTPClient(client),
		option.WithQuery("a", "1"),
		option.WithHeader("X-Caller", "override"),
	)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	if captured.Method != http.MethodGet {
		t.Fatalf("method mismatch: %s", captured.Method)
	}
	if captured.URL.String() != "https://sdk.test/api/resource?a=1" {
		t.Fatalf("baseURL/path/query join mismatch: %s", captured.URL.String())
	}
	if captured.Header.Get("X-Caller") != "override" {
		t.Fatalf("caller header override lost: %q", captured.Header.Get("X-Caller"))
	}
	if dst["id"] != "one" {
		t.Fatalf("body selection mismatch: %v", dst)
	}
}

type listItem struct {
	ID   string `json:"id"`
	JSON struct {
		ID          respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r *listItem) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

func TestConformanceResponseAndPageTerminal(t *testing.T) {
	var page pagination.Page[listItem]
	if err := json.Unmarshal([]byte(`{"data":[{"id":"a"}],"has_more":false,"last_id":"a"}`), &page); err != nil {
		t.Fatalf("unmarshal page: %v", err)
	}
	if len(page.Data) != 1 || page.Data[0].ID != "a" {
		t.Fatalf("shared response shape mismatch: %+v", page.Data)
	}
	next, err := page.GetNextPage()
	if err != nil {
		t.Fatalf("terminal page must not error: %v", err)
	}
	if next != nil {
		t.Fatalf("terminal page (has_more=false) must yield no next page")
	}

	var token pagination.TokenPage[listItem]
	if err := json.Unmarshal([]byte(`{"data":[{"id":"a"}],"has_more":false}`), &token); err != nil {
		t.Fatalf("unmarshal token page: %v", err)
	}
	tnext, err := token.GetNextPage()
	if err != nil {
		t.Fatalf("terminal token page must not error: %v", err)
	}
	if tnext != nil {
		t.Fatalf("terminal token page must yield no next page")
	}
}

func TestConformanceAutoPagerWalksPages(t *testing.T) {
	calls := 0
	client := &http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		if calls == 1 {
			return jsonResponse(200, `{"data":[{"id":"a"}],"has_more":true,"last_id":"a"}`), nil
		}
		return jsonResponse(200, `{"data":[{"id":"b"}],"has_more":false,"last_id":"b"}`), nil
	})}

	var raw *http.Response
	var res *pagination.Page[listItem]
	opts := []option.RequestOption{
		option.WithResponseInto(&raw),
		option.WithBaseURL("https://sdk.test/api/"),
		option.WithMaxRetries(0),
		option.WithHTTPClient(client),
	}
	cfg, err := convention.NewRequestConfig(context.Background(), http.MethodGet, "resource", nil, &res, opts...)
	if err != nil {
		t.Fatalf("new request config: %v", err)
	}
	if err := cfg.Execute(); err != nil {
		t.Fatalf("execute first page: %v", err)
	}
	res.SetPageConfig(cfg, raw)

	var ids []string
	pager := pagination.NewPageAutoPager(res, nil)
	for pager.Next() {
		ids = append(ids, pager.Current().ID)
	}
	if err := pager.Err(); err != nil {
		t.Fatalf("auto-pager error: %v", err)
	}
	if len(ids) != 2 || ids[0] != "a" || ids[1] != "b" {
		t.Fatalf("auto-pager must walk both pages, got %v", ids)
	}
	if calls != 2 {
		t.Fatalf("auto-pager must issue exactly two requests, got %d", calls)
	}
}
