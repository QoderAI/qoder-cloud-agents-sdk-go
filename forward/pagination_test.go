package forward_test

import (
	"context"
	"encoding/json"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/pagination"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestForwardIDPagination(t *testing.T) {
	for _, direction := range []string{"after_id", "before_id"} {
		t.Run(direction, func(t *testing.T) {
			calls := 0
			client := testClient(func(r *http.Request) (*http.Response, error) {
				calls++
				q := r.URL.Query()
				if q.Get("limit") != "1" || q.Get("status") != "active" {
					t.Fatal(q)
				}
				if calls == 1 {
					return reply(r, 200, `{"data":[{"id":"tmpl_one"}],"first_id":"tmpl_one","last_id":"tmpl_one","has_more":true}`), nil
				}
				if q.Get(direction) != "tmpl_one" {
					t.Fatal(q)
				}
				return reply(r, 200, `{"data":[{"id":"tmpl_two"}],"has_more":false}`), nil
			})
			params := forward.TemplateListParams{Limit: forward.Int(1), Status: forward.String("active")}
			if direction == "before_id" {
				params.BeforeID = forward.String("tmpl_start")
			}
			pager := client.Templates.ListAutoPaging(context.Background(), params)
			var ids []string
			for pager.Next() {
				ids = append(ids, pager.Current().ID)
				if len(ids) > 2 {
					t.Fatal("pagination did not stop")
				}
			}
			if pager.Err() != nil || calls != 2 || !reflect.DeepEqual(ids, []string{"tmpl_one", "tmpl_two"}) {
				t.Fatalf("ids=%v calls=%d err=%v", ids, calls, pager.Err())
			}
		})
	}
}

func TestForwardPageTokenPagination(t *testing.T) {
	calls := 0
	client := testClient(func(r *http.Request) (*http.Response, error) {
		calls++
		q := r.URL.Query()
		if calls == 1 {
			return reply(r, 200, `{"data":[],"next_page":"opaque+/=","has_more":true}`), nil
		}
		if q.Get("page") != "opaque+/=" || q.Has("after_id") || q.Has("before_id") || q.Get("limit") != "1" {
			t.Fatal(q)
		}
		return reply(r, 200, `{"data":[{"id":"env_one"}],"has_more":false}`), nil
	})
	pager := client.Environments.ListAutoPaging(context.Background(), forward.EnvironmentListParams{Limit: forward.Int(1), AfterID: forward.String("legacy")})
	if !pager.Next() || pager.Current().ID != "env_one" || pager.Next() || pager.Err() != nil || calls != 2 {
		t.Fatalf("calls=%d err=%v", calls, pager.Err())
	}
}

func TestForwardPaginationErrors(t *testing.T) {
	var detached pagination.Page[forward.Template]
	if err := json.Unmarshal([]byte(`{"data":[{"id":"a"}],"last_id":"a","has_more":true}`), &detached); err != nil {
		t.Fatal(err)
	}
	if _, err := detached.GetNextPage(); err == nil {
		t.Fatal("detached page must return an error")
	}
	for _, status := range []int{200, 503} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			calls := 0
			client := testClient(func(r *http.Request) (*http.Response, error) {
				calls++
				if calls > 2 {
					t.Fatal("pagination retried a stagnant cursor")
				}
				if calls == 2 && status == 503 {
					return reply(r, status, `{"error":{"code":"unavailable","message":"retry later"}}`), nil
				}
				return reply(r, 200, `{"data":[{"id":"a"}],"first_id":"a","last_id":"a","has_more":true}`), nil
			})
			pager := client.Templates.ListAutoPaging(context.Background(), forward.TemplateListParams{})
			for pager.Next() {
			}
			if pager.Err() == nil || calls != 2 {
				t.Fatalf("calls=%d err=%v", calls, pager.Err())
			}
			if status == 200 && !strings.Contains(pager.Err().Error(), "did not advance") {
				t.Fatal(pager.Err())
			}
		})
	}
}
