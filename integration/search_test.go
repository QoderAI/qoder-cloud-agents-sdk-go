package integration_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
)

type searchTransport func(*http.Request) (*http.Response, error)

func (f searchTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type searchCase struct {
	Mode, Resource, Entry, Service, Path string
	Body                                 json.RawMessage
	Response                             map[string]any
}

func searchCases(t *testing.T) []searchCase {
	t.Helper()
	data, err := os.ReadFile("testdata/search-contracts.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixture struct{ Cases []searchCase }
	if err := json.Unmarshal(data, &fixture); err != nil {
		t.Fatal(err)
	}
	if len(fixture.Cases) != 12 {
		t.Fatalf("got %d search contracts", len(fixture.Cases))
	}
	return fixture.Cases
}
func searchResponse(r *http.Request, body any) *http.Response {
	data, _ := json.Marshal(body)
	return &http.Response{StatusCode: 200, Header: http.Header{"Content-Type": {"application/json"}, "X-Request-Id": {"req_search"}}, Body: io.NopCloser(strings.NewReader(string(data))), Request: r}
}
func searchResource(c searchCase, fn searchTransport) reflect.Value {
	opts := []option.RequestOption{option.WithPAT("test"), option.WithMaxRetries(0), option.WithBaseURL("https://search.test/"), option.WithHeader("x-qoder-beta", "client-beta"), option.WithHTTPClient(&http.Client{Transport: fn})}
	if c.Mode == "forward" {
		client := forward.NewClient(opts...)
		return reflect.ValueOf(&client).Elem().FieldByName(strings.ToUpper(c.Entry[:1]) + c.Entry[1:]).Addr()
	}
	client := managed.NewClient(opts...)
	return reflect.ValueOf(&client).Elem().FieldByName(strings.ToUpper(c.Entry[:1]) + c.Entry[1:]).Addr()
}
func searchArguments(t *testing.T, method reflect.Value, data []byte, identity bool) []reflect.Value {
	t.Helper()
	params := reflect.New(method.Type().In(1))
	if err := json.Unmarshal(data, params.Interface()); err != nil {
		t.Fatal(err)
	}
	if identity {
		params.Elem().FieldByName("IdentityID").Set(reflect.ValueOf(forward.String("idn +/")))
	}
	return []reflect.Value{reflect.ValueOf(context.Background()), params.Elem()}
}
func TestResourceSearchBodyPagination(t *testing.T) {
	for _, c := range searchCases(t) {
		t.Run(c.Mode+"/"+c.Resource, func(t *testing.T) {
			calls := 0
			var original map[string]any
			if err := json.Unmarshal(c.Body, &original); err != nil {
				t.Fatal(err)
			}
			resource := searchResource(c, func(r *http.Request) (*http.Response, error) {
				var got map[string]any
				if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
					t.Fatal(err)
				}
				want := map[string]any{}
				for k, v := range original {
					want[k] = v
				}
				cursors := []string{original["page"].(string), "opaque +/=? one", "opaque +/=? two"}
				if calls >= len(cursors) {
					t.Fatal("unexpected extra request")
				}
				want["page"] = cursors[calls]
				if !reflect.DeepEqual(want, got) {
					t.Fatalf("body got=%v want=%v", got, want)
				}
				if r.Method != "POST" || r.URL.Path != c.Path {
					t.Fatalf("%s %s", r.Method, r.URL)
				}
				query := r.URL.Query()
				if c.Mode == "forward" {
					if query.Get("identity_id") != "idn +/" {
						t.Fatalf("ownership: %v", query)
					}
					query.Del("identity_id")
				}
				if len(query) != 0 {
					t.Fatalf("filters leaked into query: %v", query)
				}
				if r.Header.Get("x-qoder-beta") != "client-beta,search-2026-08-31" {
					t.Fatalf("beta: %v", r.Header.Values("x-qoder-beta"))
				}
				calls++
				body := map[string]any{}
				for k, v := range c.Response {
					body[k] = v
				}
				if calls == 2 {
					body["data"] = []any{}
				}
				body["has_more"] = calls < 3
				body["next_page"] = []any{"opaque +/=? one", "opaque +/=? two", nil}[calls-1]
				return searchResponse(r, body), nil
			})
			method := resource.MethodByName("SearchAutoPaging")
			args := searchArguments(t, method, c.Body, c.Mode == "forward")
			pager := method.Call(args)[0].Interface().(interface {
				Next() bool
				Err() error
			})
			count := 0
			for pager.Next() {
				count++
			}
			if pager.Err() != nil {
				t.Fatal(pager.Err())
			}
			if calls != 3 || count != 2*len(c.Response["data"].([]any)) {
				t.Fatalf("calls=%d items=%d", calls, count)
			}
			raw, err := json.Marshal(args[1].Interface())
			if err != nil {
				t.Fatal(err)
			}
			var unchanged map[string]any
			_ = json.Unmarshal(raw, &unchanged)
			if !reflect.DeepEqual(unchanged, original) {
				t.Fatal("caller filters mutated")
			}
		})
	}
}
func TestResourceSearchEmptyBody(t *testing.T) {
	for _, c := range searchCases(t) {
		t.Run(c.Mode+"/"+c.Resource, func(t *testing.T) {
			calls := 0
			resource := searchResource(c, func(r *http.Request) (*http.Response, error) {
				calls++
				b, _ := io.ReadAll(r.Body)
				if string(b) != "{}" || r.URL.RawQuery != "" {
					t.Fatalf("body=%s query=%s", b, r.URL.RawQuery)
				}
				return searchResponse(r, map[string]any{"data": []any{}, "has_more": false, "next_page": nil, "total": nil}), nil
			})
			method := resource.MethodByName("Search")
			result := method.Call(searchArguments(t, method, []byte("{}"), false))
			if !result[1].IsNil() {
				t.Fatal(result[1].Interface())
			}
			next := result[0].MethodByName("GetNextPage").Call(nil)
			if !next[0].IsNil() || !next[1].IsNil() || calls != 1 {
				t.Fatal("terminal page fetched again")
			}
		})
	}
}
func TestResourceSearchRepeatedCursor(t *testing.T) {
	c := searchCases(t)[0]
	calls := 0
	resource := searchResource(c, func(r *http.Request) (*http.Response, error) {
		calls++
		return searchResponse(r, map[string]any{"data": []any{}, "has_more": true, "next_page": "same"}), nil
	})
	method := resource.MethodByName("SearchAutoPaging")
	pager := method.Call(searchArguments(t, method, []byte(`{"page":"same"}`), false))[0].Interface().(interface {
		Next() bool
		Err() error
	})
	if pager.Next() || pager.Err() == nil || !strings.Contains(pager.Err().Error(), "did not advance") || calls != 1 {
		t.Fatalf("calls=%d err=%v", calls, pager.Err())
	}
}
func TestResourceSearchErrors(t *testing.T) {
	for _, status := range []int{400, 403, 429, 503} {
		t.Run(fmt.Sprint(status), func(t *testing.T) {
			c := searchCases(t)[0]
			resource := searchResource(c, func(r *http.Request) (*http.Response, error) {
				res := searchResponse(r, map[string]any{"error": map[string]any{"message": "search failed"}})
				res.StatusCode = status
				return res, nil
			})
			method := resource.MethodByName("Search")
			results := method.Call(searchArguments(t, method, []byte("{}"), false))
			err, ok := results[1].Interface().(*convention.Error)
			if !ok || err.StatusCode != status || err.RequestID != "req_search" {
				t.Fatalf("error: %v", results[1].Interface())
			}
		})
	}
}

func TestResourceSearchCursorCycle(t *testing.T) {
	c := searchCases(t)[0]
	calls := 0
	resource := searchResource(c, func(r *http.Request) (*http.Response, error) {
		calls++
		if calls > 3 {
			t.Fatal("cursor cycle was not stopped")
		}
		return searchResponse(r, map[string]any{"data": []any{}, "has_more": true, "next_page": []string{"a", "b", "a"}[calls-1]}), nil
	})
	method := resource.MethodByName("SearchAutoPaging")
	pager := method.Call(searchArguments(t, method, []byte(`{}`), false))[0].Interface().(interface {
		Next() bool
		Err() error
	})
	if pager.Next() || pager.Err() == nil || !strings.Contains(pager.Err().Error(), "cycle") || calls != 3 {
		t.Fatalf("calls=%d err=%v", calls, pager.Err())
	}
}

func TestResourceSearchRetryAndCancellation(t *testing.T) {
	calls := 0
	client := forward.NewClient(option.WithPAT("test"), option.WithBaseURL("https://search.test/"), option.WithMaxRetries(1), option.WithHTTPClient(&http.Client{Transport: searchTransport(func(r *http.Request) (*http.Response, error) {
		if err := r.Context().Err(); err != nil {
			return nil, err
		}
		calls++
		data, _ := io.ReadAll(r.Body)
		if string(data) != `{"metadata":{"team":"core"}}` || r.Header.Get("x-qoder-beta") != "search-2026-08-31" {
			t.Fatalf("body=%s beta=%s", data, r.Header.Get("x-qoder-beta"))
		}
		res := searchResponse(r, map[string]any{"data": []any{}, "has_more": false})
		if calls == 1 {
			res.StatusCode = 429
			res.Header.Set("Retry-After-Ms", "0")
		}
		return res, nil
	})}))
	_, err := client.Files.Search(context.Background(), forward.FileSearchParams{Metadata: map[string]string{"team": "core"}})
	if err != nil || calls != 2 {
		t.Fatalf("calls=%d err=%v", calls, err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err = client.Files.Search(ctx, forward.FileSearchParams{})
	if err == nil {
		t.Fatal("canceled request succeeded")
	}
}

func TestResourceSearchRequestOptionsSurvivePagination(t *testing.T) {
	calls := 0
	client := managed.NewClient(option.WithPAT("test"), option.WithBaseURL("https://search.test/"), option.WithHTTPClient(&http.Client{Transport: searchTransport(func(r *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		want := "start"
		if calls > 0 {
			want = "next"
		}
		if body["page"] != want || body["limit"] != float64(1) || r.URL.Query().Get("trace") != "keep" {
			t.Fatalf("body=%v url=%v", body, r.URL)
		}
		if r.Header.Get("x-qoder-beta") != "search-2026-08-31" || r.Header.Get("qoder-workspace-id") != "workspace" {
			t.Fatal("Search headers changed")
		}
		calls++
		if calls > 2 {
			t.Fatal("unexpected request")
		}
		return searchResponse(r, map[string]any{"data": []any{}, "has_more": calls == 1, "next_page": "next", "total": 0}), nil
	})}))
	page, err := client.Agents.Search(context.Background(), managed.AgentSearchParams{Limit: managed.Int(1), WorkspaceID: managed.String("workspace"), Betas: []managed.QoderBeta{"search-2026-08-31"}}, option.WithJSONSet("page", "start"), option.WithQuery("trace", "keep"))
	if err != nil {
		t.Fatal(err)
	}
	if page.Total != 0 || !page.JSON.Total.Valid() {
		t.Fatal("total was not decoded")
	}
	next, err := page.GetNextPage()
	if err != nil || next == nil || calls != 2 {
		t.Fatalf("calls=%d err=%v", calls, err)
	}
}
