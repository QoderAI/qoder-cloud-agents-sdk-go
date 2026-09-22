// Package testutil contains shared assertions for the two SDK contract suites.
package testsupport

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apierror"
)

type Endpoint struct{ Service, Method string }
type Transport func(*http.Request) (*http.Response, error)

func (f Transport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

// FailureContracts exercises every documented endpoint through its public service method.
// Expected HTTP failures are independent of parameter/response fixture generation.
func FailureContracts(t *testing.T, endpoints []Endpoint, client func(Transport) any) {
	t.Helper()
	for _, ep := range endpoints {
		t.Run(ep.Service+"/"+ep.Method, func(t *testing.T) {
			for _, status := range []int{400, 401, 403, 404, 409, 422, 429, 500, 503} {
				t.Run(fmt.Sprint(status), func(t *testing.T) {
					calls := 0
					c := client(func(r *http.Request) (*http.Response, error) {
						calls++
						return &http.Response{StatusCode: status, Header: http.Header{"Content-Type": {"application/json"}, "X-Request-Id": {"header-id"}}, Body: io.NopCloser(strings.NewReader(`{"request_id":"body-id","error":{"type":"test_error","code":"TEST_FAILURE","message":"expected failure"}}`)), Request: r}, nil
					})
					_, err := invoke(t, c, ep, context.Background(), -1)
					var api *apierror.Error
					if !errors.As(err, &api) {
						t.Fatalf("expected API error, got %T: %v", err, err)
					}
					if api.StatusCode != status || api.Code != "TEST_FAILURE" || api.Message != "expected failure" || api.RequestID != "body-id" || api.Type() != "test_error" || api.Request == nil || api.Response == nil {
						t.Fatalf("lost error details: %#v", api)
					}
					if calls != 1 {
						t.Fatalf("HTTP calls=%d, want 1", calls)
					}
				})
			}
			for _, plain := range []bool{false, true} {
				t.Run(fmt.Sprintf("gateway_error_plain_%t", plain), func(t *testing.T) {
					body := "{broken"
					if plain {
						body = "upstream unavailable"
					}
					c := client(func(r *http.Request) (*http.Response, error) {
						return &http.Response{StatusCode: 502, Header: http.Header{"X-Request-Id": {"gateway-id"}}, Body: io.NopCloser(strings.NewReader(body)), Request: r}, nil
					})
					_, err := invoke(t, c, ep, context.Background(), -1)
					var api *apierror.Error
					if !errors.As(err, &api) || api.StatusCode != 502 || api.RequestID != "gateway-id" || api.Message != body {
						t.Fatalf("gateway error lost: %v", err)
					}
				})
			}
			t.Run("transport_error", func(t *testing.T) {
				cause := errors.New("connection failed")
				_, err := invoke(t, client(func(*http.Request) (*http.Response, error) { return nil, cause }), ep, context.Background(), -1)
				if !errors.Is(err, cause) {
					t.Fatalf("lost transport cause: %v", err)
				}
			})
			t.Run("canceled", func(t *testing.T) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				_, err := invoke(t, client(func(r *http.Request) (*http.Response, error) { return nil, r.Context().Err() }), ep, ctx, -1)
				if !errors.Is(err, context.Canceled) {
					t.Fatalf("lost cancellation: %v", err)
				}
			})
			method := serviceMethod(t, client(func(*http.Request) (*http.Response, error) {
				t.Fatal("missing path reached transport")
				return nil, nil
			}), ep)
			for i := 1; i < method.Type().NumIn()-1; i++ {
				typ := method.Type().In(i)
				if typ.Kind() == reflect.String {
					t.Run(fmt.Sprintf("missing_path_%d", i), func(t *testing.T) {
						_, err := call(t, method, context.Background(), i)
						if err == nil || !strings.Contains(err.Error(), "missing required") {
							t.Fatalf("missing path must fail locally: %v", err)
						}
					})
				}
				if typ.Kind() == reflect.Struct {
					for j := 0; j < typ.NumField(); j++ {
						if typ.Field(j).Tag.Get("path") == "" {
							continue
						}
						t.Run("missing_path_"+typ.Field(j).Tag.Get("path"), func(t *testing.T) {
							_, err := call(t, method, context.Background(), 1000*i+j)
							if err == nil || !strings.Contains(err.Error(), "missing required") {
								t.Fatalf("missing path must fail locally: %v", err)
							}
						})
					}
				}
			}
		})
	}
}
func serviceMethod(t *testing.T, client any, ep Endpoint) reflect.Value {
	t.Helper()
	root := reflect.ValueOf(client)
	if root.Kind() != reflect.Pointer {
		copy := reflect.New(root.Type())
		copy.Elem().Set(root)
		root = copy
	}
	var found reflect.Value
	var walk func(reflect.Value)
	walk = func(v reflect.Value) {
		if v.Kind() == reflect.Pointer {
			v = v.Elem()
		}
		for i := 0; i < v.NumField(); i++ {
			f := v.Field(i)
			if f.Kind() == reflect.Struct && strings.HasSuffix(f.Type().Name(), "Service") {
				if f.Type().Name() == ep.Service {
					found = f.Addr().MethodByName(ep.Method)
					return
				}
				walk(f)
			}
		}
	}
	walk(root)
	if !found.IsValid() {
		t.Fatalf("missing %s.%s", ep.Service, ep.Method)
	}
	return found
}
func invoke(t *testing.T, client any, ep Endpoint, ctx context.Context, emptyPath int) (any, error) {
	return call(t, serviceMethod(t, client, ep), ctx, emptyPath)
}
func call(t *testing.T, method reflect.Value, ctx context.Context, emptyPath int, bodies ...json.RawMessage) (any, error) {
	t.Helper()
	args := []reflect.Value{reflect.ValueOf(ctx)}
	for i := 1; i < method.Type().NumIn()-1; i++ {
		typ := method.Type().In(i)
		v := reflect.New(typ).Elem()
		if typ.Kind() == reflect.String && i != emptyPath {
			v.SetString("segment /?%#")
		}
		if typ.Kind() == reflect.Struct {
			if len(bodies) > 0 {
				if err := json.Unmarshal(bodies[0], v.Addr().Interface()); err != nil {
					t.Fatalf("request parameters: %v", err)
				}
			}
			for j := 0; j < typ.NumField(); j++ {
				if typ.Field(j).Tag.Get("path") != "" && v.Field(j).Kind() == reflect.String && emptyPath != 1000*i+j {
					v.Field(j).SetString("segment /?%#")
				}
			}
		}
		args = append(args, v)
	}
	out := method.Call(args)
	if len(out) == 1 {
		if stream, ok := out[0].Interface().(interface {
			Next() bool
			Err() error
			Close() error
		}); ok {
			defer stream.Close()
			if stream.Next() {
				t.Fatal("failed stream returned an event")
			}
			return nil, stream.Err()
		}
		if out[0].IsNil() {
			return nil, nil
		}
		return nil, out[0].Interface().(error)
	}
	if !out[1].IsNil() {
		return nil, out[1].Interface().(error)
	}
	return out[0].Interface(), nil
}
