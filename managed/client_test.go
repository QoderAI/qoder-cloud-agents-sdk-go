package managed_test

import (
	"context"
	"errors"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	managed "github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

func TestManagedAPIInventory(t *testing.T) {
	fixtures := contracts(t)
	if len(fixtures) != 95 {
		t.Fatalf("got %d API contracts", len(fixtures))
	}
	expected := map[string]bool{}
	for _, c := range fixtures {
		key := c.Service + "." + c.Name
		if expected[key] {
			t.Fatalf("duplicate API %s", key)
		}
		expected[key] = true
	}
	client := managed.NewClient()
	ss := map[string]reflect.Value{}
	services(reflect.ValueOf(&client).Elem(), ss)
	for name, srv := range ss {
		tp := srv.Addr().Type()
		for i := 0; i < tp.NumMethod(); i++ {
			m := tp.Method(i)
			key := name + "." + m.Name
			if strings.HasSuffix(m.Name, "AutoPaging") || key == "SessionEventService.NewResumableStream" {
				continue
			}
			if !expected[key] {
				t.Errorf("API outside scope: %s", key)
			}
			delete(expected, key)
		}
	}
	if len(expected) != 0 {
		t.Fatalf("missing APIs: %v", expected)
	}
}

func TestRetriesAndAPIError(t *testing.T) {
	for _, method := range []string{"read", "write", "idempotent", "conflict"} {
		t.Run(method, func(t *testing.T) {
			calls := 0
			c := testClient(func(r *http.Request) (*http.Response, error) {
				calls++
				status := 503
				if method == "conflict" {
					status = 409
				}
				if calls > 1 {
					return reply(r, 200, `{"id":"ok"}`), nil
				}
				res := reply(r, status, `{"type":"error","request_id":"request-from-body","error":{"type":"conflict_error","code":"VERSION_CONFLICT","message":"stale version"}}`)
				res.Header.Set("Retry-After-Ms", "0")
				return res, nil
			}, option.WithMaxRetries(2))
			var e error
			switch method {
			case "read":
				_, e = c.Agents.Get(context.Background(), "a", managed.AgentGetParams{})
			case "idempotent":
				_, e = c.Agents.New(context.Background(), managed.AgentNewParams{}, option.WithHeader("Idempotency-Key", "operation"))
			default:
				_, e = c.Agents.New(context.Background(), managed.AgentNewParams{})
			}
			if method == "read" || method == "idempotent" {
				if e != nil || calls != 2 {
					t.Fatal(e, calls)
				}
			} else {
				var api *convention.Error
				if !errors.As(e, &api) || calls != 1 || api.RequestID != "request-from-body" || api.Code != "VERSION_CONFLICT" || api.Type() != "conflict_error" {
					t.Fatal(e, calls)
				}
			}
		})
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	c := testClient(func(r *http.Request) (*http.Response, error) { return nil, r.Context().Err() }, option.WithMaxRetries(2))
	_, e := c.Agents.Get(ctx, "a", managed.AgentGetParams{})
	if !errors.Is(e, context.Canceled) {
		t.Fatal(e)
	}
}
