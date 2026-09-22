package convention_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
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
)

// Fixed upstream provenance: anthropics/anthropic-sdk-go tag v1.74.0,
// commit 3cb26e4450dc5294618ab9cc40377e5dd79ac413.

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
