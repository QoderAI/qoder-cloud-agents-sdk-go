package forward_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"reflect"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

// Existing API contracts from docs.qoder.com/cloud-agents/api, 2026-10-20.
func alignmentFixture(t *testing.T) map[string]json.RawMessage {
	t.Helper()
	data, err := os.ReadFile("testdata/contract-alignment-0.8.json")
	if err != nil {
		t.Fatal(err)
	}
	var fixture map[string]json.RawMessage
	if err := json.Unmarshal(data, &fixture); err != nil {
		t.Fatal(err)
	}
	return fixture
}

func TestForwardContractAlignment080Requests(t *testing.T) {
	ctx := context.Background()
	var body map[string]any
	var requests int
	client := testClient(func(r *http.Request) (*http.Response, error) {
		requests++
		if r.Method == "GET" {
			if r.URL.Query().Get("search") != "report & 中文" {
				t.Fatal(r.URL)
			}
			return reply(r, 200, `{"data":[],"has_more":false}`), nil
		}
		data, err := io.ReadAll(r.Body)
		if err != nil {
			t.Fatal(err)
		}
		body = nil
		if err := json.Unmarshal(data, &body); err != nil {
			t.Fatal(err)
		}
		if r.URL.Path == "/api/v1/forward/memory_stores" && r.Header.Get("Idempotency-Key") != "mount-one" {
			t.Fatal("missing idempotency key")
		}
		return reply(r, 200, `{"id":"contract_one"}`), nil
	})
	for _, value := range []param.Opt[int64]{{}, param.Null[int64](), forward.Int(3)} {
		_, err := client.Templates.New(ctx, forward.TemplateNewParams{Name: "sdk-contract", EnvironmentID: "env_one", Model: forward.ModelConfigUnionParam{OfString: forward.String("model")}, MaxToolRounds: value})
		if err != nil {
			t.Fatal(err)
		}
		got, present := body["max_tool_rounds"]
		if value.Valid() && got != float64(3) {
			t.Fatal(body)
		}
		if param.IsNull(value) && (!present || got != nil) {
			t.Fatal(body)
		}
		if param.IsOmitted(value) && present {
			t.Fatal(body)
		}
		_, err = client.Templates.Update(ctx, "tmpl_one", forward.TemplateUpdateParams{MaxToolRounds: value})
		if err != nil {
			t.Fatal(err)
		}
		got, present = body["max_tool_rounds"]
		if value.Valid() && got != float64(3) || param.IsNull(value) && (!present || got != nil) || param.IsOmitted(value) && present {
			t.Fatal(body)
		}
	}
	for _, value := range []param.Opt[bool]{{}, forward.Bool(false), forward.Bool(true)} {
		_, err := client.Batches.New(ctx, forward.BatchNewParams{InputFileID: "file_one", CompletionWindow: "24h", IgnoreIdleWindow: value})
		if err != nil {
			t.Fatal(err)
		}
		got, present := body["ignore_idle_window"]
		if present != value.Valid() || present && got != value.Value {
			t.Fatal(body)
		}
	}
	if _, err := client.Sessions.List(ctx, forward.SessionListParams{Search: forward.String("report & 中文")}); err != nil {
		t.Fatal(err)
	}
	before := requests
	for i := 0; i < 2; i++ {
		_, err := client.MemoryStores.New(ctx, forward.MemoryStoreNewParams{Name: "sdk-contract", IdempotencyKey: "mount-one", Mount: forward.MemoryStoreNewParamsMount{IdentityID: "idn_one", TemplateID: "tmpl_one", Access: forward.String("read_write")}})
		if err != nil {
			t.Fatal(err)
		}
		var mount any
		if err := json.Unmarshal(alignmentFixture(t)["mount"], &mount); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(body, map[string]any{"name": "sdk-contract", "mount": mount}) {
			t.Fatal(body)
		}
	}
	if requests != before+2 {
		t.Fatal("mount must use a single native operation")
	}
}

func TestForwardContractAlignment080Responses(t *testing.T) {
	f := alignmentFixture(t)
	decode := func(key string, dst any) {
		t.Helper()
		if err := json.Unmarshal(f[key], dst); err != nil {
			t.Fatal(err)
		}
	}
	var session forward.Session
	decode("session", &session)
	if len(session.OutcomeEvaluations) != 1 || !session.OutcomeEvaluations[0].JSON.Iteration.Valid() || session.OutcomeEvaluations[0].JSON.Explanation.Raw() != "null" {
		t.Fatal(session)
	}
	var thread forward.SessionThread
	decode("thread", &thread)
	if !thread.Usage.JSON.TotalCredits.Valid() || thread.Usage.TotalCredits != 0 || !thread.Stats.JSON.StartupSeconds.Valid() {
		t.Fatal(thread)
	}
	var absent forward.SessionThread
	if err := json.Unmarshal([]byte(`{"id":"empty"}`), &absent); err != nil {
		t.Fatal(err)
	}
	if absent.JSON.Usage.Valid() || absent.JSON.Stats.Valid() {
		t.Fatal("absent must remain unknown")
	}
	var schedule forward.Schedule
	decode("schedule", &schedule)
	if schedule.Source != "api" || schedule.JSON.SourceSessionID.Raw() != "null" || schedule.InitialEventsVersion != 2 || schedule.TriggerPolicy.StartAt.Day() != 20 || schedule.TriggerPolicy.StopAt.Day() != 21 || schedule.TriggerPolicy.JSON.LastRunAt.Raw() != "null" {
		t.Fatal(schedule)
	}
	var run forward.ScheduleRun
	decode("run", &run)
	if run.Source != "tool" || run.SourceSessionID != "sess_one" || run.JSON.InitialEventsVersion.Valid() {
		t.Fatal(run)
	}
	var env forward.Environment
	decode("environment", &env)
	if env.Config.SetupScript != "echo sdk-contract" || !env.BindingInfo.JSON.AgentTemplateCount.Valid() {
		t.Fatal(env)
	}
	var batch forward.Batch
	if err := json.Unmarshal([]byte(`{"ignore_idle_window":false,"queue_reason":"capacity"}`), &batch); err != nil {
		t.Fatal(err)
	}
	if batch.IgnoreIdleWindow || !batch.JSON.IgnoreIdleWindow.Valid() || batch.QueueReason != "capacity" {
		t.Fatal(batch)
	}
}
