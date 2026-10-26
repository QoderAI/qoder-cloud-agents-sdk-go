package forward_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestDecisionAPIContracts(t *testing.T) {
	for _, status := range []int{200, 400} {
		t.Run("new_"+statusName(status), func(t *testing.T) {
			client := contractClient(t, "createDecision", status)
			res, err := client.Decisions.New(context.Background(), contractParams[forward.DecisionNewParams](t, "createDecision"))
			checkError(t, status, err)
			if status == 200 {
				checkDecoded(t, res)
				if res.Answers[0].Value != false || res.Answers[1].Score != 0 || !res.Answers[1].JSON.Score.Valid() || res.JSON.Usage.Raw() != "null" || res.JSON.SessionID.Raw() != "null" {
					t.Fatalf("lost false, zero, or null: %s", res.RawJSON())
				}
			}
		})
		t.Run("get_"+statusName(status), func(t *testing.T) {
			client := contractClient(t, "getDecision", status)
			res, err := client.Decisions.Get(context.Background(), pathSegment)
			checkError(t, status, err)
			if status == 200 {
				checkDecoded(t, res)
			}
		})
		t.Run("list_"+statusName(status), func(t *testing.T) {
			client := contractClient(t, "listDecisions", status)
			res, err := client.Decisions.List(context.Background(), contractParams[forward.DecisionListParams](t, "listDecisions"))
			checkError(t, status, err)
			if status == 200 {
				checkDecoded(t, res)
			}
		})
	}
}

func TestDecisionSessionInput(t *testing.T) {
	client := testClient(func(r *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if len(body) != 2 || body["session_id"] != "sess_own" || body["questions"] == nil {
			t.Fatalf("unexpected session input: %#v", body)
		}
		if r.Header.Get("X-Qoder-Beta") != "" || r.Header.Get("X-Request-Only") != "present" {
			t.Fatalf("unexpected headers: %v", r.Header)
		}
		return reply(r, 200, `{"id":"dec_session","session_id":"sess_own","status":"completed","context":{"mode":"session","policy_version":"decision-context-v3","token_estimator":"caw-heuristic-v1","estimated_input_tokens":28464,"max_input_tokens":30000,"included_segment_count":13,"included_event_count":14,"snapshot_at":"2026-10-23T03:50:37.005193664Z","scan_limited":false,"truncated":true,"truncation_reasons":["segment_budget"],"anchor_event_id":"evt_last","first_event_id":"evt_first","last_event_id":"evt_last","future_field":true},"answers":[{"type":"refusal","name":"private","refusal":"Cannot assess"}],"future_top":1}`), nil
	})
	params := contractParams[forward.DecisionNewParams](t, "createDecision")
	params.State = nil
	params.SessionID = forward.String("sess_own")
	res, err := client.Decisions.New(context.Background(), params, option.WithHeader("X-Request-Only", "present"))
	if err != nil {
		t.Fatal(err)
	}
	if res.Context.IncludedSegmentCount != 13 || res.Context.EstimatedInputTokens != 28464 || !res.Context.Truncated || res.Context.ScanLimited || res.Context.SnapshotAt.Nanosecond() != 5193664 || res.Answers[0].Refusal != "Cannot assess" || res.Answers[0].JSON.Value.Valid() {
		t.Fatalf("lost session metadata or refusal: %s", res.RawJSON())
	}
	if res.Context.JSON.ExtraFields["future_field"].Raw() != "true" || res.JSON.ExtraFields["future_top"].Raw() != "1" {
		t.Fatal("unknown response fields were lost")
	}
}

func TestDecisionPagination(t *testing.T) {
	for _, reverse := range []bool{false, true} {
		name := "forward"
		if reverse {
			name = "backward"
		}
		t.Run(name, func(t *testing.T) {
			calls := 0
			client := testClient(func(r *http.Request) (*http.Response, error) {
				calls++
				if r.URL.Query().Get("session_id") != "sess_own" || r.URL.Query().Get("limit") != "1" || r.Header.Get("X-Page") != "kept" {
					t.Fatal("pagination lost filter/options")
				}
				if calls == 1 {
					return reply(r, 200, `{"data":[{"id":"dec_one"}],"first_id":"dec_one","last_id":"dec_one","has_more":true}`), nil
				}
				cursor := "after_id"
				if reverse {
					cursor = "before_id"
				}
				if r.URL.Query().Get(cursor) != "dec_one" {
					t.Fatal(r.URL.Query())
				}
				return reply(r, 200, `{"data":[{"id":"dec_two","status":"failed","error":{"code":"decision_upstream_timeout","message":"Timed out"},"completed_at":"2026-10-23T05:00:00Z"}],"has_more":false,"first_id":"dec_two","last_id":"dec_two"}`), nil
			})
			params := forward.DecisionListParams{SessionID: forward.String("sess_own"), Limit: forward.Int(1)}
			if reverse {
				params.BeforeID = forward.String("dec_start")
			}
			pager := client.Decisions.ListAutoPaging(context.Background(), params, option.WithHeader("X-Page", "kept"))
			var ids []string
			for pager.Next() {
				d := pager.Current()
				ids = append(ids, d.ID)
				if len(ids) > 2 {
					t.Fatal("pagination did not stop")
				}
				if d.ID == "dec_two" && (d.Error.Code != "decision_upstream_timeout" || d.CompletedAt.IsZero()) {
					t.Fatalf("failed record: %s", d.RawJSON())
				}
			}
			if pager.Err() != nil || calls != 2 || !reflect.DeepEqual(ids, []string{"dec_one", "dec_two"}) {
				t.Fatalf("ids=%v calls=%d err=%v", ids, calls, pager.Err())
			}
		})
	}
}

func TestDecisionIdempotencyAndErrors(t *testing.T) {
	for _, keyed := range []bool{false, true} {
		calls := 0
		var bodies []string
		client := testClient(func(r *http.Request) (*http.Response, error) {
			calls++
			body, _ := io.ReadAll(r.Body)
			bodies = append(bodies, string(body))
			wantKey := ""
			if keyed {
				wantKey = "request-key"
			}
			if r.Header.Get("Idempotency-Key") != wantKey || strings.Contains(string(body), "Idempotency") {
				t.Fatal("key not confined to header")
			}
			if calls == 1 {
				res := reply(r, 503, `{"error":{"code":"unavailable","message":"Retry later"}}`)
				res.Header.Set("Retry-After", "0.001")
				return res, nil
			}
			return reply(r, 200, `{"id":"dec_once","status":"completed"}`), nil
		}, option.WithMaxRetries(1))
		params := forward.DecisionNewParams{State: map[string]any{"ok": false}, Questions: contractParams[forward.DecisionNewParams](t, "createDecision").Questions}
		var opts []option.RequestOption
		if keyed {
			params.IdempotencyKey = forward.String("params-key")
			opts = append(opts, option.WithHeader("Idempotency-Key", "request-key"))
		}
		_, err := client.Decisions.New(context.Background(), params, opts...)
		if keyed {
			if err != nil || calls != 2 || bodies[0] != bodies[1] {
				t.Fatalf("keyed retry changed body: calls=%d err=%v", calls, err)
			}
		} else if err == nil || calls != 1 {
			t.Fatalf("unkeyed request retried: calls=%d err=%v", calls, err)
		}
	}
	for _, status := range []int{400, 404, 409, 422, 504} {
		client := testClient(func(r *http.Request) (*http.Response, error) {
			return reply(r, status, `{"type":"error","request_id":"req_decision","decision_id":"dec_failed","error":{"type":"api_error","code":"decision_test_error","message":"Decision failed"}}`), nil
		})
		_, err := client.Decisions.New(context.Background(), forward.DecisionNewParams{})
		var api *forward.Error
		if !errors.As(err, &api) || api.StatusCode != status || api.Code != "decision_test_error" || !strings.Contains(api.RawJSON(), "dec_failed") {
			t.Fatalf("lost error correlation: %v", err)
		}
	}
	client := testClient(func(*http.Request) (*http.Response, error) { t.Fatal("empty ID reached HTTP"); return nil, nil })
	if _, err := client.Decisions.Get(context.Background(), ""); err == nil {
		t.Fatal("empty ID accepted")
	}
}

func TestDecisionAuthenticationAndConflict(t *testing.T) {
	for _, mode := range []string{"pat", "sa"} {
		t.Run(mode, func(t *testing.T) {
			calls := 0
			client := testClient(func(r *http.Request) (*http.Response, error) {
				calls++
				if r.Header.Get("Authorization") != "Bearer decision-token" || r.Header.Get("Idempotency-Key") != "decision-key" {
					t.Fatal("Decision did not use the configured authentication/idempotency key")
				}
				return reply(r, 409, `{"decision_id":"dec_running","error":{"code":"decision_in_progress","message":"In progress"}}`), nil
			}, option.WithAuthMode(mode), option.WithAuthToken("decision-token"), option.WithMaxRetries(2))
			params := contractParams[forward.DecisionNewParams](t, "createDecision")
			params.IdempotencyKey = forward.String("decision-key")
			_, err := client.Decisions.New(context.Background(), params)
			var api *forward.Error
			if !errors.As(err, &api) || api.StatusCode != 409 || calls != 1 {
				t.Fatalf("conflict must not retry: calls=%d err=%v", calls, err)
			}
		})
	}
}
