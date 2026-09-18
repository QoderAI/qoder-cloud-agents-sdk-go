package forward_test

import (
	"context"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestForwardAPIInventory(t *testing.T) {
	client := forward.NewClient()
	root := reflect.ValueOf(&client).Elem()
	actual := map[string]bool{}
	var walk func(reflect.Value)
	walk = func(value reflect.Value) {
		for i := 0; i < value.NumField(); i++ {
			field := value.Field(i)
			if field.Kind() != reflect.Struct || !strings.HasSuffix(field.Type().Name(), "Service") {
				continue
			}
			methods := field.Addr().Type()
			for j := 0; j < methods.NumMethod(); j++ {
				name := methods.Method(j).Name
				key := field.Type().Name() + "." + name
				if !strings.HasSuffix(name, "AutoPaging") && key != "SessionEventService.NewResumableStream" {
					actual[key] = true
				}
			}
			walk(field)
		}
	}
	walk(root)
	if len(actual) != 110 {
		t.Fatalf("HTTP methods: %d want 110", len(actual))
	}
	baseline := map[string]bool{}
	for _, c := range operations(t) {
		baseline[c.OperationID] = true
	}
	for _, c := range contracts(t) {
		key := c.Service + "." + c.Name
		if !actual[key] || !baseline[c.OperationID] {
			t.Errorf("unmatched contract %s: %s", c.OperationID, key)
		}
		delete(actual, key)
		delete(baseline, c.OperationID)
	}
	if len(actual) != 0 || len(baseline) != 0 {
		t.Fatalf("uncovered methods=%v baseline=%v", actual, baseline)
	}
}

func TestClientConfigurationPrecedence(t *testing.T) {
	t.Setenv("QODER_PAT", "environment-token")
	t.Setenv("QODER_FORWARD_BASE_URL", "https://env.test/api/v1/forward")
	for _, explicit := range []bool{false, true} {
		wantHost, wantToken := "env.test", "environment-token"
		var opts []option.RequestOption
		if explicit {
			wantHost, wantToken = "explicit.test", "explicit-token"
			opts = append(opts, option.WithBaseURL("https://explicit.test/api/v1/forward"), option.WithPAT(wantToken))
		}
		opts = append(opts, option.WithHTTPClient(&http.Client{Transport: roundTripFunc(func(r *http.Request) (*http.Response, error) {
			if r.URL.Host != wantHost || r.URL.Path != "/api/v1/forward/models" || r.Header.Get("Authorization") != "Bearer "+wantToken {
				t.Fatalf("request: %s %v", r.URL, r.Header)
			}
			if r.Header.Get("X-Test") != "request" {
				t.Fatal("request options did not win")
			}
			return reply(r, 200, `{"data":[],"has_more":false}`), nil
		})}), option.WithHeader("X-Test", "client"))
		client := forward.NewClient(opts...)
		client.Models.Options = append(client.Models.Options, option.WithHeader("X-Test", "service"))
		if _, err := client.Models.List(context.Background(), option.WithHeader("X-Test", "request")); err != nil {
			t.Fatal(err)
		}
	}
}

func TestServiceOptionsDoNotShareBackingSlices(t *testing.T) {
	options := make([]option.RequestOption, 1, 8)
	options[0] = option.WithHeader("X-Test", "initial")
	service := forward.NewSessionService(options...)
	service.Options = append(service.Options, option.WithHeader("X-Test", "parent"))
	service.Events.Options = append(service.Events.Options, option.WithHeader("X-Test", "child"))
	if len(options) != 1 || len(service.Threads.Options) != 1 {
		t.Fatal("service options modified siblings")
	}
	client := testClient(func(r *http.Request) (*http.Response, error) {
		if r.Header.Get("X-Test") != "parent" {
			t.Fatalf("parent options overwritten: %v", r.Header)
		}
		return reply(r, 200, `{"data":[],"has_more":false}`), nil
	}, service.Options...)
	if _, err := client.Sessions.List(context.Background(), forward.SessionListParams{}); err != nil {
		t.Fatal(err)
	}
}
