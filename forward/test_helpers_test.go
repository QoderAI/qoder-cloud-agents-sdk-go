package forward_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/testutil"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

const pathSegment = "segment /?%#"

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func reply(r *http.Request, status int, body string) *http.Response {
	return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": []string{"application/json"}, "X-Request-Id": []string{"req_contract"}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}
}

func testClient(fn roundTripFunc, opts ...option.RequestOption) forward.Client {
	return forward.NewClient(append([]option.RequestOption{option.WithPAT("secret-pat"), option.WithBaseURL("https://qoder.test/api/v1/forward"), option.WithMaxRetries(0), option.WithHTTPClient(&http.Client{Transport: fn})}, opts...)...)
}

type operationCase struct {
	OperationID string          `json:"operation_id"`
	Method      string          `json:"http_method"`
	Path        string          `json:"path"`
	Parameters  []parameterCase `json:"parameters"`
	Body        *bodyCase       `json:"request_body"`
}
type parameterCase struct {
	Name     string `json:"wire_name"`
	Location string `json:"location"`
	Value    any    `json:"value"`
}
type bodyCase struct {
	Kind   string          `json:"kind"`
	Value  map[string]any  `json:"value"`
	Fields []formFieldCase `json:"fields"`
}
type formFieldCase struct {
	Name   string `json:"wire_name"`
	Binary bool   `json:"binary"`
	Value  any    `json:"value"`
}
type contractCase struct {
	OperationID string          `json:"operation_id"`
	Service     string          `json:"service"`
	Entry       string          `json:"entry"`
	Name        string          `json:"name"`
	Response    json.RawMessage `json:"response"`
	Stream      bool            `json:"stream"`
	Download    bool            `json:"download"`
	Raw         bool            `json:"raw"`
	Empty       bool            `json:"empty"`
}

func readJSON(t *testing.T, path string, dst any) {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if err = json.Unmarshal(b, dst); err != nil {
		t.Fatal(err)
	}
}
func operations(t *testing.T) []operationCase {
	t.Helper()
	var fixture struct {
		Count int             `json:"operation_count"`
		Cases []operationCase `json:"cases"`
	}
	readJSON(t, "testdata/api-operation-cases.json", &fixture)
	if fixture.Count != 110 || len(fixture.Cases) != 110 {
		t.Fatal("Forward API scope must contain 110 operations")
	}
	return fixture.Cases
}
func operation(t *testing.T, id string) operationCase {
	t.Helper()
	for _, c := range operations(t) {
		if c.OperationID == id {
			return c
		}
	}
	t.Fatalf("missing baseline operation %s", id)
	return operationCase{}
}
func contracts(t *testing.T) []contractCase {
	t.Helper()
	var cases []contractCase
	readJSON(t, "testdata/api-contracts.json", &cases)
	return cases
}
func contract(t *testing.T, id string) contractCase {
	t.Helper()
	for _, c := range contracts(t) {
		if c.OperationID == id {
			return c
		}
	}
	t.Fatalf("missing contract %s", id)
	return contractCase{}
}

// The original fixture's converter represented identity_ids
// as an empty object, though the published API accepts strings or arrays.
func parameterValue(p parameterCase) any {
	if p.Name == "identity_ids" {
		return []string{"idn_one", "idn_two"}
	}
	if strings.HasSuffix(p.Name, "[]") {
		return []string{fmt.Sprint(p.Value), "second-value"}
	}
	return p.Value
}

func bodyValue(c operationCase) map[string]any {
	if c.OperationID == "sendSessionEvents" {
		// The original fixture used an untyped empty event; use the documented input shape.
		return map[string]any{"events": []any{map[string]any{"type": "user.message", "content": []any{map[string]any{"type": "text", "text": "hello"}}}}}
	}
	return c.Body.Value
}

func setField(t *testing.T, value reflect.Value, tag, key string, input any) {
	t.Helper()
	for i := 0; i < value.NumField(); i++ {
		field := value.Type().Field(i)
		if strings.Split(field.Tag.Get(tag), ",")[0] != key {
			continue
		}
		b, err := json.Marshal(input)
		if err != nil {
			t.Fatal(err)
		}
		if err = json.Unmarshal(b, value.Field(i).Addr().Interface()); err != nil {
			t.Fatalf("%s.%s: %v", value.Type(), field.Name, err)
		}
		return
	}
	t.Fatalf("%s lacks %s field %q", value.Type(), tag, key)
}

func contractParams[T any](t *testing.T, id string) T {
	t.Helper()
	c := operation(t, id)
	var params T
	if c.Body != nil && c.Body.Kind == "json" {
		b, err := json.Marshal(bodyValue(c))
		if err != nil {
			t.Fatal(err)
		}
		if err = json.Unmarshal(b, &params); err != nil {
			t.Fatalf("%s params: %v", id, err)
		}
	}
	v := reflect.ValueOf(&params).Elem()
	for _, p := range c.Parameters {
		if p.Location != "path" {
			setField(t, v, p.Location, p.Name, parameterValue(p))
		}
	}
	if c.Body != nil && c.Body.Kind == "multipart" {
		for _, f := range c.Body.Fields {
			if !f.Binary {
				setField(t, v, "json", f.Name, f.Value)
				continue
			}
			for i := 0; i < v.NumField(); i++ {
				if strings.Split(v.Type().Field(i).Tag.Get("json"), ",")[0] != f.Name {
					continue
				}
				file := convention.UploadFile{Reader: strings.NewReader("sdk-contract"), Name: "skill/SKILL.md", MediaType: "text/markdown"}
				if v.Field(i).Kind() == reflect.Slice {
					v.Field(i).Set(reflect.ValueOf([]io.Reader{file}))
				} else {
					v.Field(i).Set(reflect.ValueOf(file))
				}
			}
		}
	}
	return params
}

func contractClient(t *testing.T, id string, status int) forward.Client {
	t.Helper()
	c, meta := operation(t, id), contract(t, id)
	calls := 0
	t.Cleanup(func() {
		want := 1
		if meta.Download && status == 200 {
			want = 2
		}
		if calls != want {
			t.Errorf("HTTP calls: %d want %d", calls, want)
		}
	})
	return testClient(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.URL.Host == "storage.test" {
			if r.Header.Get("Authorization") != "" {
				t.Error("PAT leaked to storage")
			}
			return reply(r, 200, "file bytes"), nil
		}
		if r.Method != c.Method {
			t.Errorf("method: %s want %s", r.Method, c.Method)
		}
		path := c.Path
		for strings.Contains(path, "{") {
			start := strings.Index(path, "{")
			end := strings.Index(path[start:], "}") + start
			path = path[:start] + url.PathEscape(pathSegment) + path[end+1:]
		}
		if r.URL.EscapedPath() != "/api/v1/forward"+path {
			t.Errorf("path: %s want %s", r.URL.EscapedPath(), path)
		}
		if r.Header.Get("Authorization") != "Bearer secret-pat" {
			t.Error("missing bearer authentication")
		}
		for _, p := range c.Parameters {
			switch p.Location {
			case "header":
				if r.Header.Get(p.Name) != fmt.Sprint(parameterValue(p)) {
					t.Errorf("header %s: %q", p.Name, r.Header.Get(p.Name))
				}
			case "query":
				want := []string{fmt.Sprint(parameterValue(p))}
				if values, ok := parameterValue(p).([]string); ok {
					want = values
				}
				if !reflect.DeepEqual(r.URL.Query()[p.Name], want) {
					t.Errorf("query %s: %v want %v", p.Name, r.URL.Query()[p.Name], want)
				}
			}
		}
		if c.Body != nil && c.Body.Kind == "json" {
			var actual map[string]any
			if err := json.NewDecoder(r.Body).Decode(&actual); err != nil {
				t.Error(err)
			}
			if !reflect.DeepEqual(actual, bodyValue(c)) {
				t.Errorf("body: %#v want %#v", actual, bodyValue(c))
			}
		} else if c.Body != nil && c.Body.Kind == "multipart" {
			if err := r.ParseMultipartForm(1 << 20); err != nil {
				t.Fatal(err)
			}
			defer r.MultipartForm.RemoveAll()
			for _, f := range c.Body.Fields {
				if f.Binary {
					files := r.MultipartForm.File[f.Name]
					if len(files) != 1 {
						t.Errorf("multipart files %s: %d", f.Name, len(files))
						continue
					}
					file, err := files[0].Open()
					if err != nil {
						t.Fatal(err)
					}
					b, err := io.ReadAll(file)
					file.Close()
					if err != nil || string(b) != "sdk-contract" {
						t.Errorf("multipart content: %q %v", b, err)
					}
				} else if r.FormValue(f.Name) != fmt.Sprint(f.Value) {
					t.Errorf("multipart %s: %s", f.Name, r.FormValue(f.Name))
				}
			}
		}
		if status != 200 {
			return reply(r, status, `{"type":"error","error":{"type":"invalid_request_error","code":"contract_probe","message":"contract probe"}}`), nil
		}
		if meta.Download {
			return reply(r, 200, `{"url":"https://storage.test/file?signature=test","expires_at":"2026-09-09T00:00:00Z"}`), nil
		}
		if meta.Raw {
			res := reply(r, 200, "PK-zip-content")
			res.Header.Set("Content-Type", "application/zip")
			return res, nil
		}
		if meta.Stream {
			if r.Header.Get("Accept") != "text/event-stream" {
				t.Error("missing SSE accept header")
			}
			res := reply(r, 200, "id: evt_one\nevent: agent.message\ndata: {\"id\":\"evt_one\",\"type\":\"agent.message\",\"session_id\":\"sess_one\",\"content\":[{\"type\":\"text\",\"text\":\"hello\"}]}\n\n")
			res.Header.Set("Content-Type", "text/event-stream")
			return res, nil
		}
		if meta.Empty {
			return reply(r, http.StatusNoContent, ""), nil
		}
		return reply(r, 200, string(meta.Response)), nil
	})
}

func checkError(t *testing.T, status int, err error) {
	t.Helper()
	if status == 200 {
		if err != nil {
			t.Fatal(err)
		}
		return
	}
	var apiError *forward.Error
	if !errors.As(err, &apiError) || apiError.StatusCode != status || apiError.Code != "contract_probe" || apiError.RequestID != "req_contract" {
		t.Fatalf("unexpected API error: %#v, %v", apiError, err)
	}
}
func statusName(status int) string { return strconv.Itoa(status) }

func checkDecoded(t *testing.T, result any) {
	t.Helper()
	if response, ok := result.(*http.Response); ok {
		defer response.Body.Close()
		b, err := io.ReadAll(response.Body)
		if err != nil || len(b) == 0 {
			t.Fatalf("empty download: %s %v", b, err)
		}
		return
	}
	checkFields(t, reflect.ValueOf(result), "response")
}
func checkFields(t *testing.T, v reflect.Value, path string) {
	t.Helper()
	testutil.CheckFields(t, v, path)
}

func jsonObject(t *testing.T, value any) map[string]any {
	t.Helper()
	b, err := json.Marshal(value)
	if err != nil {
		t.Fatal(err)
	}
	var result map[string]any
	if err = json.NewDecoder(bytes.NewReader(b)).Decode(&result); err != nil {
		t.Fatal(err)
	}
	return result
}
