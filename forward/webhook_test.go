package forward_test

import (
	"context"
	"encoding/json"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
	"io"
	"net/http"
	"reflect"
	"testing"
)

func TestWebhookAPIContracts(t *testing.T) {
	for _, status := range []int{200, 400} {
		t.Run(statusName(status), func(t *testing.T) {
			client0 := contractClient(t, "createWebhookEndpoint", status)
			res0, err0 := client0.Webhooks.Endpoints.New(context.Background(), contractParams[forward.WebhookEndpointNewParams](t, "createWebhookEndpoint"))
			checkError(t, status, err0)
			if status == 200 {
				checkDecoded(t, res0)
			}
			client1 := contractClient(t, "listWebhookEndpoints", status)
			res1, err1 := client1.Webhooks.Endpoints.List(context.Background(), contractParams[forward.WebhookEndpointListParams](t, "listWebhookEndpoints"))
			checkError(t, status, err1)
			if status == 200 {
				checkDecoded(t, res1)
			}
			client2 := contractClient(t, "getWebhookEndpoint", status)
			res2, err2 := client2.Webhooks.Endpoints.Get(context.Background(), pathSegment)
			checkError(t, status, err2)
			if status == 200 {
				checkDecoded(t, res2)
			}
			client3 := contractClient(t, "updateWebhookEndpoint", status)
			res3, err3 := client3.Webhooks.Endpoints.Update(context.Background(), pathSegment, contractParams[forward.WebhookEndpointUpdateParams](t, "updateWebhookEndpoint"))
			checkError(t, status, err3)
			if status == 200 {
				checkDecoded(t, res3)
			}
			client4 := contractClient(t, "deleteWebhookEndpoint", status)
			err4 := client4.Webhooks.Endpoints.Delete(context.Background(), pathSegment)
			checkError(t, status, err4)
			client5 := contractClient(t, "enableWebhookEndpoint", status)
			res5, err5 := client5.Webhooks.Endpoints.Enable(context.Background(), pathSegment)
			checkError(t, status, err5)
			if status == 200 {
				checkDecoded(t, res5)
			}
			client6 := contractClient(t, "disableWebhookEndpoint", status)
			res6, err6 := client6.Webhooks.Endpoints.Disable(context.Background(), pathSegment)
			checkError(t, status, err6)
			if status == 200 {
				checkDecoded(t, res6)
			}
			client7 := contractClient(t, "testWebhookEndpoint", status)
			res7, err7 := client7.Webhooks.Endpoints.Test(context.Background(), pathSegment)
			checkError(t, status, err7)
			if status == 200 {
				checkDecoded(t, res7)
			}
		})
	}
}
func TestWebhookPaginationAndBodies(t *testing.T) {
	calls := 0
	client := testClient(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.Header.Get("Authorization") != "Bearer test-sat" {
			t.Fatal("auth lost")
		}
		if r.Method == "GET" {
			if r.URL.Query().Get("limit") != "2" || r.URL.Query().Has("page") {
				t.Fatal(r.URL)
			}
			if calls == 1 {
				return reply(r, 200, `{"data":[{"id":"z"},{"id":"a"}],"last_id":"opaque/+?","has_more":true}`), nil
			}
			if r.URL.Query().Get("after_id") != "opaque/+?" {
				t.Fatal(r.URL)
			}
			return reply(r, 200, `{"data":[],"first_id":null,"last_id":null,"has_more":false}`), nil
		}
		var raw []byte
		if r.Body != nil {
			raw, _ = io.ReadAll(r.Body)
		}
		if r.Method == "PUT" {
			var body map[string]any
			json.Unmarshal(raw, &body)
			if !reflect.DeepEqual(body, map[string]any{"active": false, "metadata": map[string]any{"remove": nil, "keep": "value"}}) {
				t.Fatalf("body: %s", raw)
			}
			return reply(r, 200, `{"active":false}`), nil
		}
		if len(raw) != 0 {
			t.Fatalf("action body: %s", raw)
		}
		if r.Method == "DELETE" {
			return reply(r, 204, ""), nil
		}
		return reply(r, 202, `{"event_id":1,"delivery_rows":0}`), nil
	}, option.WithAuthMode("sa"), option.WithAuthToken("test-sat"))
	pager := client.Webhooks.Endpoints.ListAutoPaging(context.Background(), forward.WebhookEndpointListParams{Limit: forward.Int(2)})
	var ids []string
	for pager.Next() {
		ids = append(ids, pager.Current().ID)
	}
	if pager.Err() != nil || !reflect.DeepEqual(ids, []string{"z", "a"}) || calls != 2 {
		t.Fatalf("pagination: %v %v %d", ids, pager.Err(), calls)
	}
	_, err := client.Webhooks.Endpoints.Update(context.Background(), "id", forward.WebhookEndpointUpdateParams{Active: forward.Bool(false), Metadata: map[string]param.Opt[string]{"remove": param.Null[string](), "keep": forward.String("value")}})
	if err != nil {
		t.Fatal(err)
	}
	result, err := client.Webhooks.Endpoints.Test(context.Background(), "id")
	if err != nil || result.DeliveryRows != 0 {
		t.Fatal(result, err)
	}
	if err = client.Webhooks.Endpoints.Delete(context.Background(), "id"); err != nil {
		t.Fatal(err)
	}
}
func TestWebhookTestIsNotReplayedWithoutIdempotency(t *testing.T) {
	calls := 0
	client := testClient(func(r *http.Request) (*http.Response, error) {
		calls++
		return reply(r, 503, `{"error":{"message":"unavailable"}}`), nil
	}, option.WithMaxRetries(2))
	_, err := client.Webhooks.Endpoints.Test(context.Background(), "id")
	if err == nil || calls != 1 {
		t.Fatalf("err %v calls %d", err, calls)
	}
}

func TestWebhookCreatedAndAcceptedStatuses(t *testing.T) {
	client := testClient(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path == "/api/v1/forward/webhook/endpoints/id/test" {
			return reply(r, 202, `{"event_id":7,"delivery_rows":0}`), nil
		}
		if r.Header.Get("Idempotency-Key") != "one-create" {
			t.Fatal("missing idempotency key")
		}
		return reply(r, 201, `{"id":"id","signing_secret":"synthetic"}`), nil
	})
	var raw *http.Response
	created, err := client.Webhooks.Endpoints.New(context.Background(), forward.WebhookEndpointNewParams{URL: "https://hooks.example.com/qca", Events: []string{"*"}, IdempotencyKey: forward.String("one-create")}, option.WithResponseInto(&raw))
	if err != nil || raw.StatusCode != 201 || created.SigningSecret != "synthetic" {
		t.Fatal(created, err)
	}
	accepted, err := client.Webhooks.Endpoints.Test(context.Background(), created.ID, option.WithResponseInto(&raw))
	if err != nil || raw.StatusCode != 202 || accepted.DeliveryRows != 0 {
		t.Fatal(accepted, err)
	}
}
