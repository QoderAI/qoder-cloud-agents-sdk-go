package managed_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
)

// Contract: https://docs.qoder.com/cloud-agents/api/sessions/schemas (2026-10-20).
func TestSessionCreditBudgetWireContract(t *testing.T) {
	for _, tc := range []struct {
		name   string
		params managed.SessionUpdateParams
		want   string
	}{
		{"decimal", managed.SessionUpdateParams{Budget: managed.ManagedAgentsBudgetLimitParam{Type: "limit", MaxCreditCost: "100.00"}}, `{"budget":{"type":"limit","max_credit_cost":"100.00"}}`},
		{"omit", managed.SessionUpdateParams{}, `{}`},
		{"clear", managed.SessionUpdateParams{Budget: param.NullStruct[managed.ManagedAgentsBudgetLimitParam]()}, `{"budget":null}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := testClient(func(r *http.Request) (*http.Response, error) {
				body, _ := io.ReadAll(r.Body)
				var got, want any
				if err := json.Unmarshal(body, &got); err != nil {
					t.Fatal(err)
				}
				if err := json.Unmarshal([]byte(tc.want), &want); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("got %s want %s", body, tc.want)
				}
				return reply(r, 200, `{"id":"sess_one","budget":{"type":"limit","max_credit_cost":"100.00"}}`), nil
			})
			session, err := client.Sessions.Update(context.Background(), "sess_one", tc.params)
			if err != nil {
				t.Fatal(err)
			}
			if session.Budget.MaxCreditCost != "100.00" {
				t.Fatalf("budget=%+v", session.Budget)
			}
		})
	}
}

func TestSessionCreditsSnapshotsAndEvents(t *testing.T) {
	const usage = `{"model_credits":7.66,"sandbox_runtime_credits":0.34,"total_credits":8}`
	var session managed.ManagedAgentsSession
	if err := json.Unmarshal([]byte(`{"id":"sess_one","usage":`+usage+`}`), &session); err != nil {
		t.Fatal(err)
	}
	if session.Usage.ModelCredits != 7.66 || session.Usage.SandboxRuntimeCredits != 0.34 || session.Usage.TotalCredits != 8 {
		t.Fatalf("usage=%+v", session.Usage)
	}
	var fsession forward.Session
	if err := json.Unmarshal([]byte(`{"id":"sess_one","usage":`+usage+`}`), &fsession); err != nil {
		t.Fatal(err)
	}
	if fsession.Usage.ModelCredits != 7.66 || fsession.Usage.SandboxRuntimeCredits != 0.34 || fsession.Usage.TotalCredits != 8 {
		t.Fatalf("forward usage=%+v", fsession.Usage)
	}
	var empty managed.ManagedAgentsSession
	if err := json.Unmarshal([]byte(`{"id":"sess_empty"}`), &empty); err != nil {
		t.Fatal(err)
	}
	if empty.JSON.Usage.Valid() {
		t.Fatal("missing usage must stay absent")
	}
	const event = `{"id":"evt_one","type":"session.usage","model_credits":7.66,"sandbox_runtime_credits":0.34,"total_credits":8}`
	var history managed.ManagedAgentsSessionEventUnion
	if err := json.Unmarshal([]byte(event), &history); err != nil {
		t.Fatal(err)
	}
	if history.TotalCredits != 8 || history.AsSessionUsage().TotalCredits != 8 {
		t.Fatal("history lost top-level credits")
	}
	var thread managed.ManagedAgentsStreamSessionThreadEventsUnion
	if err := json.Unmarshal([]byte(event), &thread); err != nil {
		t.Fatal(err)
	}
	if thread.TotalCredits != 8 || thread.AsSessionUsage().ModelCredits != 7.66 {
		t.Fatal("thread stream lost credits")
	}
	client := testClient(func(r *http.Request) (*http.Response, error) {
		res := reply(r, 200, "id: evt_one\nevent: session.usage\ndata: "+event+"\n\nid: evt_two\nevent: agent.thread_context_compacted\ndata: {\"id\":\"evt_two\",\"type\":\"agent.thread_context_compacted\",\"usage\":{\"credits\":0}}\n\nid: evt_three\nevent: span.model_request_end\ndata: {\"id\":\"evt_three\",\"type\":\"span.model_request_end\",\"model_usage\":{\"credits\":1.25}}\n\n")
		res.Header.Set("Content-Type", "text/event-stream")
		return res, nil
	})
	stream := client.Sessions.Events.StreamEvents(context.Background(), "sess_one", managed.SessionEventStreamParams{})
	defer stream.Close()
	var seen []string
	for stream.Next() {
		e := stream.Current()
		seen = append(seen, e.Type)
		switch e.Type {
		case "session.usage":
			if e.TotalCredits != 8 || e.AsSessionUsage().TotalCredits != 8 {
				t.Fatal("SSE lost credits")
			}
		case "agent.thread_context_compacted":
			if !e.Usage.JSON.Credits.Valid() || e.AsAgentThreadContextCompacted().Usage.Credits != 0 {
				t.Fatal("zero credits must be preserved")
			}
		case "span.model_request_end":
			if e.ModelUsage.Credits != 1.25 {
				t.Fatal("model credits lost")
			}
		}
	}
	if err := stream.Err(); err != nil {
		t.Fatal(err)
	}
	if strings.Join(seen, ",") != "session.usage,agent.thread_context_compacted,span.model_request_end" {
		t.Fatal(seen)
	}
}
