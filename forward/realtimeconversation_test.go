package forward_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestRealtimeConversationAPIContracts(t *testing.T) {
	for _, status := range []int{200, 400} {
		t.Run("create_"+statusName(status), func(t *testing.T) {
			client := contractClient(t, "createRealtimeConversation", status)
			res, err := client.Realtime.Conversations.New(context.Background(), contractParams[forward.RealtimeConversationNewParams](t, "createRealtimeConversation"))
			checkError(t, status, err)
			if status == 200 && (res.ID != "conv_contract" || res.CreatedAt != "2026-10-13T01:02:03.123456789Z" || res.Config.Audio.Output.Voice != "longanqian") {
				t.Fatalf("conversation = %#v", res)
			}
		})
		t.Run("history_"+statusName(status), func(t *testing.T) {
			client := contractClient(t, "getRealtimeConversationHistory", status)
			res, err := client.Realtime.Conversations.History(context.Background(), pathSegment, contractParams[forward.RealtimeConversationHistoryParams](t, "getRealtimeConversationHistory"))
			checkError(t, status, err)
			if status == 200 && (res.Conversation.ID != "conv_contract" || len(res.Events) != 1 || res.Events[0].ID != "vevt_contract" || !strings.Contains(res.Events[0].RawJSON(), `"future": true`) || res.HasNextPage()) {
				t.Fatalf("history = %#v", res)
			}
		})
	}
}

func TestRealtimeCreationIdempotencyAndRetry(t *testing.T) {
	for _, key := range []string{"", strings.Repeat("a", 257), " key", "key ", "key\n", strings.Repeat("你", 86)} {
		t.Run("invalid_"+statusName(len(key)), func(t *testing.T) {
			client := testClient(func(*http.Request) (*http.Response, error) { t.Fatal("invalid key reached transport"); return nil, nil })
			_, err := client.Realtime.Conversations.New(context.Background(), forward.RealtimeConversationNewParams{IdentityID: "idn", TemplateID: "tmpl", IdempotencyKey: forward.String(key)})
			if err == nil {
				t.Fatal("expected key validation error")
			}
		})
	}
	calls := 0
	client := testClient(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.Header.Get("Idempotency-Key") != "request-key" {
			t.Fatal("retry changed operation key")
		}
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		if len(body) != 2 || body["identity_id"] != "idn" || body["template_id"] != "tmpl" {
			t.Fatalf("unexpected body %#v", body)
		}
		if calls == 1 {
			response := reply(r, 503, `{"error":{"code":"idempotency_unavailable"}}`)
			response.Header.Set("Retry-After", "0.001")
			return response, nil
		}
		return reply(r, 201, `{"id":"conv","type":"voice.conversation"}`), nil
	}, option.WithMaxRetries(1))
	_, err := client.Realtime.Conversations.New(context.Background(), forward.RealtimeConversationNewParams{IdentityID: "idn", TemplateID: "tmpl", IdempotencyKey: forward.String("params-key")}, option.WithHeader("Idempotency-Key", "request-key"))
	if err != nil || calls != 2 {
		t.Fatalf("calls=%d err=%v", calls, err)
	}
	client = testClient(func(r *http.Request) (*http.Response, error) {
		response := reply(r, 409, `{"error":{"code":"idempotency_key_in_progress","message":"busy"}}`)
		response.Header.Set("Retry-After", "1")
		return response, nil
	}, option.WithMaxRetries(1))
	_, err = client.Realtime.Conversations.New(context.Background(), forward.RealtimeConversationNewParams{IdempotencyKey: forward.String("key")})
	var api *forward.Error
	if !errors.As(err, &api) || api.Code != "idempotency_key_in_progress" || api.Response.Header.Get("Retry-After") != "1" {
		t.Fatalf("lost conflict: %v", err)
	}
}

func TestRealtimeHistoryPaginationAndBounds(t *testing.T) {
	calls := 0
	client := testClient(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.URL.Query().Get("types") != "work" || r.URL.Query().Get("limit") != "1" || r.Header.Get("X-Test") != "page-options" {
			t.Fatal("pagination lost options")
		}
		switch calls {
		case 1:
			if r.URL.Query().Get("before") != "" {
				t.Fatal("unexpected first cursor")
			}
			return reply(r, 200, `{"conversation":{"id":"conv"},"events":[{"id":"vevt_new","type":"voice.work.running"}],"page":{"has_more":true,"next_before":"older /?"}}`), nil
		case 2:
			if r.URL.Query().Get("before") != "older /?" {
				t.Fatal("cursor was not encoded")
			}
			return reply(r, 200, `{"conversation":{"id":"conv"},"events":[],"page":{"has_more":true,"next_before":"oldest"}}`), nil
		default:
			if r.URL.Query().Get("before") != "oldest" {
				t.Fatal("empty page ended pagination")
			}
			return reply(r, 200, `{"conversation":{"id":"conv"},"events":[{"id":"vevt_old","type":"voice.work.accepted"}],"page":{"has_more":false,"next_before":null}}`), nil
		}
	})
	pager := client.Realtime.Conversations.HistoryAutoPaging(context.Background(), "conv", forward.RealtimeConversationHistoryParams{Types: forward.String("work"), Limit: forward.Int(1)}, option.WithHeader("X-Test", "page-options"))
	var ids []string
	for pager.Next() {
		ids = append(ids, pager.Current().ID)
	}
	if pager.Err() != nil || strings.Join(ids, ",") != "vevt_new,vevt_old" || calls != 3 {
		t.Fatalf("ids=%v calls=%d err=%v", ids, calls, pager.Err())
	}
	for _, params := range []forward.RealtimeConversationHistoryParams{{Limit: forward.Int(0)}, {Limit: forward.Int(101)}, {Before: forward.String("")}, {Types: forward.String("all")}} {
		_, err := client.Realtime.Conversations.History(context.Background(), "conv", params)
		if err == nil {
			t.Fatal("invalid history filter accepted")
		}
	}
	page := &forward.RealtimeHistoryPage{}
	if next, err := page.GetNextPage(context.Background()); err != nil || next != nil {
		t.Fatal("exhausted page should terminate")
	}
}
