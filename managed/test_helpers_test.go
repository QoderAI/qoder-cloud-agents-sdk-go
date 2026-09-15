package managed_test

import (
	"encoding/json"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/testutil"
	managed "github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"io"
	"net/http"
	"net/url"
	"os"
	"reflect"
	"strings"
	"testing"
)

type roundTripFunc func(*http.Request) (*http.Response, error)

func (f roundTripFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func reply(r *http.Request, status int, body string) *http.Response {
	return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": []string{"application/json"}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}
}

type apiContract struct {
	Service, Name, Method, Route, Doc string
	Response                          json.RawMessage
}

func contracts(t *testing.T) []apiContract {
	t.Helper()
	b, e := os.ReadFile("testdata/api-contracts.json")
	if e != nil {
		t.Fatal(e)
	}
	var c []apiContract
	if e = json.Unmarshal(b, &c); e != nil {
		t.Fatal(e)
	}
	return c
}

func services(v reflect.Value, into map[string]reflect.Value) {
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Kind() == reflect.Struct && strings.HasSuffix(f.Type().Name(), "Service") {
			into[f.Type().Name()] = f
			services(f, into)
		}
	}
}

// Contract fixtures come from official documentation, independently of the adapted services.
// The inventory also detects extra HTTP methods outside the authorized 95 endpoints.

func checkFields(t *testing.T, v reflect.Value, path string) {
	t.Helper()
	testutil.CheckFields(t, v, path)
}

func contractClient(t *testing.T, service, method string) managed.Client {
	t.Helper()
	var c apiContract
	for _, fixture := range contracts(t) {
		if fixture.Service == service && fixture.Name == method {
			c = fixture
			break
		}
	}
	if c.Service == "" {
		t.Fatalf("missing API contract: %s.%s", service, method)
	}
	count := 0
	t.Cleanup(func() {
		want := 1
		if c.Service == "FileService" && c.Name == "Download" {
			want = 2
		}
		if count != want {
			t.Errorf("HTTP calls: got %d want %d", count, want)
		}
	})
	client := managed.NewClient(option.WithPAT("test-token"), option.WithBaseURL("https://qoder.test/api/v1/cloud"), option.WithMaxRetries(0), option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
		count++
		if r.URL.Host == "storage.test" {
			return reply(r, 200, "file bytes"), nil
		}
		if r.Method != c.Method {
			t.Errorf("method: got %s want %s", r.Method, c.Method)
		}
		path := c.Route
		for strings.Contains(path, "{") {
			start := strings.Index(path, "{")
			end := strings.Index(path[start:], "}") + start
			path = path[:start] + url.PathEscape("segment /?%#") + path[end+1:]
		}
		if r.URL.EscapedPath() != "/api/v1/cloud"+path {
			t.Errorf("path: %s want %s", r.URL.EscapedPath(), path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Error("missing Qoder authentication")
		}
		if r.URL.Query().Has("beta") {
			t.Error("unexpected beta query parameter")
		}
		for header := range r.Header {
			if header == "Authorization" || header == "Content-Type" || header == "Accept" || header == "User-Agent" || strings.HasPrefix(header, "X-Qoder-") {
				continue
			}
			t.Errorf("unexpected API header: %s", header)
		}

		body := string(c.Response)
		if c.Service == "FileService" && c.Name == "Download" {
			body = `{"url":"https://storage.test/file","expires_at":"2026-09-08T00:00:00Z"}`
		}
		if c.Name == "StreamEvents" {
			res := reply(r, 200, "event: session.status_idle\ndata: {\"id\":\"evt_1\",\"type\":\"session.status_idle\"}\n\n")
			res.Header.Set("Content-Type", "text/event-stream")
			return res, nil
		}
		return reply(r, 200, body), nil
	})}))
	return client
}

func testClient(fn roundTripFunc, opts ...option.RequestOption) managed.Client {
	return managed.NewClient(append([]option.RequestOption{option.WithPAT("secret-pat"), option.WithBaseURL("https://qoder.test/api/v1/cloud"), option.WithMaxRetries(0), option.WithHTTPClient(&http.Client{Transport: fn})}, opts...)...)
}

func jsonObject(t *testing.T, v any) map[string]any {
	t.Helper()
	b, e := json.Marshal(v)
	if e != nil {
		t.Fatal(e)
	}
	var m map[string]any
	if e = json.Unmarshal(b, &m); e != nil {
		t.Fatal(e)
	}
	return m
}

type doerFunc func(*http.Request) (*http.Response, error)

func (f doerFunc) Do(r *http.Request) (*http.Response, error) { return f(r) }
