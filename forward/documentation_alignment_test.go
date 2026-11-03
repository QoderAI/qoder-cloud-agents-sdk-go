package forward_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/ssestream"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

// Official Forward Session schemas and Environment schemas, retrieved 2026-11-03.
func TestForwardDocumentationBudgetAndModelWire(t *testing.T) {
	budget := forward.SessionBudgetParam{Type: "limit", MaxCreditCost: "100.00"}
	model := forward.ModelConfigUnionParam{OfConfig: &forward.ModelConfigParam{ID: "ultimate", Speed: forward.String("high"), ContextWindow: forward.Int(200000)}}
	for _, tc := range []struct {
		name   string
		params forward.SessionUpdateParams
		want   string
	}{
		{"omit", forward.SessionUpdateParams{}, `{}`},
		{"clear", forward.SessionUpdateParams{Budget: param.NullStruct[forward.SessionBudgetParam](), Config: forward.SessionUpdateParamsConfigParam{Model: param.NullStruct[forward.ModelConfigUnionParam]()}}, `{"budget":null,"config":{"model":null}}`},
		{"replace", forward.SessionUpdateParams{Budget: budget, Config: forward.SessionUpdateParamsConfigParam{Model: model}}, `{"budget":{"type":"limit","max_credit_cost":"100.00"},"config":{"model":{"id":"ultimate","speed":"high","context_window":200000}}}`},
		{"string", forward.SessionUpdateParams{Config: forward.SessionUpdateParamsConfigParam{Model: forward.ModelConfigUnionParam{OfString: forward.String("ultimate")}}}, `{"config":{"model":"ultimate"}}`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			client := testClient(func(r *http.Request) (*http.Response, error) {
				raw, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatal(err)
				}
				var got, want any
				if err = json.Unmarshal(raw, &got); err != nil {
					t.Fatal(err)
				}
				if err = json.Unmarshal([]byte(tc.want), &want); err != nil {
					t.Fatal(err)
				}
				if !reflect.DeepEqual(got, want) {
					t.Fatalf("got %s want %s", raw, tc.want)
				}
				return reply(r, 200, tc.want), nil
			})
			result, err := client.Sessions.Update(context.Background(), "s", tc.params)
			if err != nil {
				t.Fatal(err)
			}
			if tc.name == "clear" && result.JSON.Budget.Raw() != "null" {
				t.Fatal("lost explicit null")
			}
			if tc.name == "omit" && result.JSON.Budget.Raw() != "" {
				t.Fatal("lost omission")
			}
			if tc.name == "replace" && (result.Budget.MaxCreditCost != "100.00" || result.Config.Model.Speed != "high") {
				t.Fatalf("result=%+v", result)
			}
			if tc.name == "string" && result.Config.Model.ID != "ultimate" {
				t.Fatal("string model lost")
			}
		})
	}
	metadata := map[string]any{"business": "kept"}
	client := testClient(func(r *http.Request) (*http.Response, error) {
		var body map[string]any
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		if body["budget"].(map[string]any)["max_credit_cost"] != "100.00" {
			t.Fatal(body)
		}
		if body["metadata"].(map[string]any)["__qca_sdk_created"] != "true" {
			t.Fatal(body)
		}
		if body["config"].(map[string]any)["model"].(map[string]any)["speed"] != "high" {
			t.Fatal(body)
		}
		return reply(r, 200, `{"id":"s"}`), nil
	})
	_, err := client.Sessions.New(context.Background(), forward.SessionNewParams{IdentityID: "idn", TemplateID: "tmpl", Budget: budget, Config: forward.SessionNewParamsConfigParam{Model: model}, Metadata: metadata})
	if err != nil {
		t.Fatal(err)
	}
	if len(metadata) != 1 {
		t.Fatal("mutated caller metadata")
	}
}

func TestForwardDocumentationEventTypes(t *testing.T) {
	for _, tc := range []struct{ decision, policy, reason string }{{"allow", "always_allow", ""}, {"ask", "always_ask", ""}, {"allow", "auto", ""}, {"ask", "auto", "indeterminate"}, {"deny", "auto", "high_risk"}} {
		t.Run(tc.decision+tc.policy, func(t *testing.T) {
			detail := map[string]any{"type": tc.policy}
			if tc.policy == "auto" {
				detail["evaluated_permission"] = map[string]any{"type": tc.decision, "reason_code": tc.reason}
			}
			raw, _ := json.Marshal(map[string]any{"evaluated_permission": tc.decision, "evaluation": detail, "future": 17})
			var e forward.SessionEvent
			if err := json.Unmarshal(raw, &e); err != nil {
				t.Fatal(err)
			}
			if e.EvaluatedPermission != tc.decision || e.Evaluation.Type != tc.policy {
				t.Fatal(e)
			}
			if tc.policy == "auto" && (e.Evaluation.EvaluatedPermission.Type != tc.decision || e.Evaluation.EvaluatedPermission.ReasonCode != tc.reason) {
				t.Fatal(e.Evaluation)
			}
			if e.JSON.ExtraFields["future"].Raw() != "17" {
				t.Fatal("unknown field lost")
			}
		})
	}
	for _, raw := range []string{`{}`, `{"budget":null}`, `{"budget":{"type":"limit","max_credit_cost":"0.01"}}`} {
		var e forward.SessionEvent
		if err := json.Unmarshal([]byte(raw), &e); err != nil {
			t.Fatal(err)
		}
		if strings.Contains(raw, "null") && e.JSON.Budget.Raw() != "null" {
			t.Fatal("null budget lost")
		}
		if raw == `{}` && e.JSON.Budget.Raw() != "" {
			t.Fatal("missing budget lost")
		}
		if strings.Contains(raw, "0.01") && e.Budget.MaxCreditCost != "0.01" {
			t.Fatal(e.Budget)
		}
	}
	payloads := []string{`{"id":"1","type":"session.usage","model_credits":0,"sandbox_runtime_credits":0,"total_credits":0}`, `{"id":"2","type":"session.usage","model_credits":0.68,"sandbox_runtime_credits":0.04,"total_credits":0.72}`, `{"id":"3","type":"session.status_idle","reason":"outcome_evaluation"}`, `{"id":"4","type":"span.outcome_evaluation_end","iteration":2,"explanation":"done","outcome_evaluation_start_id":"start"}`}
	client := testClient(func(r *http.Request) (*http.Response, error) {
		if strings.HasSuffix(r.URL.Path, "/stream") {
			var b strings.Builder
			for i, p := range payloads {
				fmt.Fprintf(&b, "id: %d\ndata: %s\n\n", i, p)
			}
			res := reply(r, 200, b.String())
			res.Header.Set("Content-Type", "text/event-stream")
			return res, nil
		}
		return reply(r, 200, `{"data":[`+strings.Join(payloads, ",")+`],"has_more":false}`), nil
	})
	for _, thread := range []bool{false, true} {
		var stream *ssestream.Stream[forward.SessionEvent]
		if thread {
			stream = client.Sessions.Threads.Events.StreamEvents(context.Background(), "s", "t", forward.SessionThreadEventStreamParams{})
		} else {
			stream = client.Sessions.Events.StreamEvents(context.Background(), "s", forward.SessionEventStreamParams{})
		}
		var got []forward.SessionEvent
		for stream.Next() {
			got = append(got, stream.Current())
		}
		err := stream.Err()
		stream.Close()
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 4 || !got[0].JSON.TotalCredits.Valid() || got[1].TotalCredits != 0.72 || got[1].ModelCredits != 0.68 || got[1].SandboxRuntimeCredits != 0.04 || got[2].Reason != "outcome_evaluation" || got[3].Iteration != 2 || got[3].Explanation != "done" || got[3].OutcomeEvaluationStartID != "start" {
			t.Fatalf("events=%+v", got)
		}
	}
	page, err := client.Sessions.Events.List(context.Background(), "s", forward.SessionEventListParams{})
	if err != nil {
		t.Fatal(err)
	}
	if page.Data[1].TotalCredits != 0.72 {
		t.Fatal("history type differs")
	}
	var env forward.Environment
	if err = json.Unmarshal([]byte(`{"config":{"type":"cloud","startup_script":"echo ready"}}`), &env); err != nil {
		t.Fatal(err)
	}
	if env.Config.StartupScript != "echo ready" {
		t.Fatal("missing startup script")
	}
}
