package managed_test

import (
	"context"
	"encoding/json"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	managed "github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"net/http"
	"reflect"
	"testing"
)

func TestAgentNew(t *testing.T) {
	client := contractClient(t, "AgentService", "New")
	response, err := client.Agents.New(context.Background(), managed.AgentNewParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Agent.New")
}

func TestAgentGet(t *testing.T) {
	client := contractClient(t, "AgentService", "Get")
	response, err := client.Agents.Get(context.Background(), "segment /?%#", managed.AgentGetParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Agent.Get")
}

func TestAgentUpdate(t *testing.T) {
	client := contractClient(t, "AgentService", "Update")
	response, err := client.Agents.Update(context.Background(), "segment /?%#", managed.AgentUpdateParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Agent.Update")
}

func TestAgentList(t *testing.T) {
	client := contractClient(t, "AgentService", "List")
	response, err := client.Agents.List(context.Background(), managed.AgentListParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Agent.List")
}

func TestAgentArchive(t *testing.T) {
	client := contractClient(t, "AgentService", "Archive")
	response, err := client.Agents.Archive(context.Background(), "segment /?%#", managed.AgentArchiveParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Agent.Archive")
}

func TestRequestFieldSemantics(t *testing.T) {
	p := managed.AgentUpdateParams{Version: managed.Int(3), Multiagent: param.NullStruct[managed.ManagedAgentsMultiagentParams](), Metadata: map[string]any{"team": "sdk", "project": "cloud"}}
	m := jsonObject(t, p)
	if v, ok := m["multiagent"]; !ok || v != nil {
		t.Fatalf("explicit object null lost: %v", m)
	}
	if _, ok := m["description"]; ok {
		t.Fatal("omitted description was sent")
	}
	if m["version"] != float64(3) {
		t.Fatal(m)
	}
	p.SetExtraFields(map[string]any{"metadata": map[string]any{"team": "sdk", "obsolete": nil}})
	m = jsonObject(t, p)
	if v, ok := m["metadata"].(map[string]any)["obsolete"]; !ok || v != nil {
		t.Fatal("metadata deletion lost")
	}
	d := managed.DeploymentUpdateParams{Metadata: param.NullMap[map[string]any]()}
	m = jsonObject(t, d)
	if v, ok := m["metadata"]; !ok || v != nil {
		t.Fatal("metadata null lost")
	}
	s := managed.SessionUpdateParams{Title: param.Null[string](), EnvironmentVariables: map[string]string{"LANG": "zh_CN"}}
	m = jsonObject(t, s)
	if m["title"] != nil || m["environment_variables"].(map[string]any)["LANG"] != "zh_CN" {
		t.Fatal(m)
	}
	r := managed.SessionResourceUpdateParams{SessionID: "session", Password: managed.String("rotated")}
	m = jsonObject(t, r)
	if m["password"] != "rotated" {
		t.Fatal(m)
	}
	if _, ok := m["authorization_token"]; ok {
		t.Fatal("empty legacy authorization_token sent")
	}
	model := managed.ManagedAgentsModelConfigParams{ID: "qoder-model", ContextWindow: managed.Int(128000), Effort: managed.ManagedAgentsModelConfigParamsEffortUnion{OfManagedAgentsModelConfigsEffortManagedAgentsEffortLevel: managed.String("high")}}
	m = jsonObject(t, model)
	if m["effort"] != "high" || m["context_window"] != float64(128000) {
		t.Fatal(m)
	}
	model.Speed = "fast"
	if _, err := json.Marshal(model); err == nil {
		t.Fatal("unsupported model.speed should be rejected")
	}
}

func TestQoderResponseFields(t *testing.T) {
	var a managed.ManagedAgentsAgent
	for _, body := range []string{`{"id":"a","model":"qoder-model","metadata":{"team":"sdk","project":"cloud"},"future":{"nested":true}}`, `{"id":"a","model":{"id":"qoder-model","effort":"high","context_window":128000},"metadata":{"team":"sdk","project":"cloud"},"future":{"nested":true}}`} {
		if e := json.Unmarshal([]byte(body), &a); e != nil {
			t.Fatal(e)
		}
		if a.Model.ID != "qoder-model" || a.Metadata["project"] != "cloud" {
			t.Fatalf("response lost fields: %+v", a)
		}
		if a.JSON.ExtraFields["future"].Raw() != `{"nested":true}` || a.RawJSON() != body {
			t.Fatal("raw/unknown JSON lost")
		}
	}
	if a.Model.Effort.Type != "high" || a.Model.Effort.AsHigh().Type != "high" || a.Model.ContextWindow != 128000 {
		t.Fatalf("model fields: %+v", a.Model)
	}
	var skill managed.Skill
	if e := json.Unmarshal([]byte(`{"display_title":"demo","latest_version":"123","source":"custom","metadata":{"project":"cloud"}}`), &skill); e != nil {
		t.Fatal(e)
	}
	if skill.DisplayName != "demo" || skill.LatestVersionID != "123" || skill.Source.Type != "custom" || skill.Metadata["project"] != "cloud" {
		t.Fatalf("skill mapping: %+v", skill)
	}
	var model managed.ModelInfo
	if e := json.Unmarshal([]byte(`{"id":"qoder-model","is_enabled":true,"efforts":["high","max"],"available_context_windows":[128000,256000],"max_input_tokens":256000}`), &model); e != nil {
		t.Fatal(e)
	}
	if !model.IsEnabled || len(model.AvailableContextWindows) != 2 || model.MaxInputTokens != 256000 {
		t.Fatalf("catalog fields: %+v", model)
	}
}

func TestPaginationAndOptionPrecedence(t *testing.T) {
	calls := 0
	c := testClient(func(r *http.Request) (*http.Response, error) {
		calls++
		q := r.URL.Query()
		if r.Header.Get("X-Test") != "request" || q.Get("limit") != "2" {
			t.Fatalf("options not preserved: %s", r.URL)
		}
		switch calls {
		case 1:
			if q.Get("after_id") != "legacy" {
				t.Fatal(q)
			}
			return reply(r, 200, `{"data":[{"id":"a"}],"has_more":true,"next_page":"opaque+/=?"}`), nil
		case 2:
			if q.Get("page") != "opaque+/=?" || q.Has("after_id") {
				t.Fatal(q)
			}
			return reply(r, 200, `{"data":[],"has_more":true,"next_page":"last"}`), nil
		case 3:
			return reply(r, 200, `{"data":[{"id":"b"}],"has_more":false,"next_page":"must-not-follow"}`), nil
		default:
			t.Fatal("pager ignored has_more=false")
			return nil, nil
		}
	}, option.WithHeader("X-Test", "client"))
	pager := c.Agents.ListAutoPaging(context.Background(), managed.AgentListParams{Limit: managed.Int(2)}, option.WithQuery("after_id", "legacy"), option.WithHeader("X-Test", "request"))
	var ids []string
	for pager.Next() {
		ids = append(ids, pager.Current().ID)
	}
	if pager.Err() != nil {
		t.Fatal(pager.Err())
	}
	if !reflect.DeepEqual(ids, []string{"a", "b"}) || calls != 3 {
		t.Fatal(ids, calls)
	}
}

func TestPaginationPreservesCustomHTTPDoer(t *testing.T) {
	calls := 0
	c := managed.NewClient(option.WithBaseURL("https://no-network.test/"), option.WithHTTPClient(doerFunc(func(r *http.Request) (*http.Response, error) {
		calls++
		if calls == 1 {
			return reply(r, 200, `{"data":[{"id":"a"}],"next_page":"next"}`), nil
		}
		return reply(r, 200, `{"data":[{"id":"b"}]}`), nil
	})))
	pager := c.Agents.ListAutoPaging(context.Background(), managed.AgentListParams{})
	for pager.Next() {
	}
	if pager.Err() != nil || calls != 2 {
		t.Fatal("custom HTTP client lost between pages", pager.Err(), calls)
	}
}
