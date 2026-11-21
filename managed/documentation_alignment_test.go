package managed_test

import (
	"encoding/json"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"reflect"
	"testing"
)

// Official Memory Store, Vault and Environment schemas, retrieved 2026-11-03.
func TestManagedDocumentationResponses(t *testing.T) {
	var memory managed.ManagedAgentsMemory
	raw := []byte(`{"type":"memory","version":2,"metadata":{"source":"test"}}`)
	if err := json.Unmarshal(raw, &memory); err != nil {
		t.Fatal(err)
	}
	if memory.Version != 2 || !memory.JSON.Version.Valid() {
		t.Fatal(memory)
	}
	var item managed.ManagedAgentsMemoryListItemUnion
	if err := json.Unmarshal(raw, &item); err != nil {
		t.Fatal(err)
	}
	if item.Version != 2 || item.Metadata["source"] != "test" || item.AsMemory().Version != 2 {
		t.Fatal(item)
	}
	for _, redacted := range []string{"false", "true"} {
		var version managed.ManagedAgentsMemoryVersion
		if err := json.Unmarshal([]byte(`{"operation":"updated","redacted":`+redacted+`}`), &version); err != nil {
			t.Fatal(err)
		}
		if version.Operation != managed.ManagedAgentsMemoryVersionOperationUpdated || !version.JSON.Redacted.Valid() || version.Redacted != (redacted == "true") {
			t.Fatal(version)
		}
	}
	var vault managed.ManagedAgentsVault
	if err := json.Unmarshal([]byte(`{"credentials":[]}`), &vault); err != nil {
		t.Fatal(err)
	}
	if !vault.JSON.Credentials.Valid() || vault.Credentials == nil || len(vault.Credentials) != 0 {
		t.Fatal(vault)
	}
	for _, id := range []string{`"ultimate"`, `{"id":"ultimate","effort":"none","context_window":200000,"effective_context_window":180000}`} {
		var thread managed.ManagedAgentsSessionThreadAgentUnionModel
		var coordinator managed.ManagedAgentsSessionMultiagentCoordinatorAgentUnionModel
		if err := json.Unmarshal([]byte(id), &thread); err != nil {
			t.Fatal(err)
		}
		if err := json.Unmarshal([]byte(id), &coordinator); err != nil {
			t.Fatal(err)
		}
		if id[0] == '{' && (thread.Effort != "none" || thread.ContextWindow != 200000 || thread.EffectiveContextWindow != 180000 || coordinator.Effort != "none" || coordinator.ContextWindow != 200000 || coordinator.EffectiveContextWindow != 180000) {
			t.Fatal(thread, coordinator)
		}
		if id[0] == '"' && (thread.OfString != "ultimate" || coordinator.OfString != "ultimate") {
			t.Fatal("string variant lost")
		}
	}
	for _, kind := range []string{"cloud", "self_hosted"} {
		var env managed.Environment
		if err := json.Unmarshal([]byte(`{"config":{"type":"`+kind+`","startup_script":"echo ready"}}`), &env); err != nil {
			t.Fatal(err)
		}
		if env.Config.StartupScript != "echo ready" {
			t.Fatal(env)
		}
	}
	for _, config := range []any{managed.CloudConfigParams{StartupScript: managed.String("echo ready")}, managed.SelfHostedConfigParams{StartupScript: managed.String("echo ready")}} {
		raw, err := json.Marshal(config)
		if err != nil {
			t.Fatal(err)
		}
		var fields map[string]any
		if err = json.Unmarshal(raw, &fields); err != nil {
			t.Fatal(err)
		}
		if fields["startup_script"] != "echo ready" {
			t.Fatal(string(raw))
		}
	}
	var thread managed.ManagedAgentsSessionThread
	if err := json.Unmarshal([]byte(`{"id":"t"}`), &thread); err != nil {
		t.Fatal(err)
	}
	if thread.JSON.Usage.Valid() {
		t.Fatal("usage must be optional")
	}
}

// Official Agent schema, retrieved 2026-11-20: roster child versions accept "latest".
func TestManagedAgentRosterLatestVersion(t *testing.T) {
	agent := func(id string, version managed.ManagedAgentsMultiagentAgentVersionUnionParam) managed.ManagedAgentsMultiagentRosterEntryParamsUnion {
		return managed.ManagedAgentsMultiagentRosterEntryParamsUnion{OfManagedAgentsAgents: &managed.ManagedAgentsMultiagentAgentParams{ID: id, Type: managed.ManagedAgentsAgentParamsTypeAgent, Version: version}}
	}
	roster := managed.ManagedAgentsMultiagentParams{Type: managed.ManagedAgentsMultiagentParamsTypeCoordinator, Agents: []managed.ManagedAgentsMultiagentRosterEntryParamsUnion{
		agent("agent_latest", managed.ManagedAgentsMultiagentAgentVersionUnionParam{OfString: managed.String("latest")}),
		agent("agent_pinned", managed.ManagedAgentsMultiagentAgentVersionUnionParam{OfInt: managed.Int(3)}),
		agent("agent_saved", managed.ManagedAgentsMultiagentAgentVersionUnionParam{}),
	}}
	body, err := json.Marshal(roster)
	if err != nil {
		t.Fatal(err)
	}
	var got, want any
	_ = json.Unmarshal(body, &got)
	_ = json.Unmarshal([]byte(`{"type":"coordinator","agents":[{"id":"agent_latest","type":"agent","version":"latest"},{"id":"agent_pinned","type":"agent","version":3},{"id":"agent_saved","type":"agent"}]}`), &want)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("roster request = %s", body)
	}
	if v := roster.Agents[1].GetVersion(); v == nil || *v != 3 || roster.Agents[0].GetVersion() != nil {
		t.Fatal("GetVersion must report only numeric versions")
	}

	var saved managed.ManagedAgentsMultiagent
	if err := json.Unmarshal([]byte(`{"type":"coordinator","agents":[{"type":"agent","id":"agent_latest","version":"latest"},{"type":"agent","id":"agent_pinned","version":3},{"type":"advisor","model":"ultimate"}]}`), &saved); err != nil {
		t.Fatal(err)
	}
	latest, pinned := saved.Agents[0].Version, saved.Agents[1].AsAgent().Version
	if latest.OfString != "latest" || !latest.JSON.OfString.Valid() || latest.JSON.OfInt.Valid() {
		t.Fatalf("latest version = %+v", latest)
	}
	if pinned.OfInt != 3 || !pinned.JSON.OfInt.Valid() || pinned.JSON.OfString.Valid() {
		t.Fatalf("pinned version = %+v", pinned)
	}
	if _, ok := saved.Agents[2].AsAny().(managed.ManagedAgentsAdvisor); !ok {
		t.Fatal("advisor variant lost")
	}
}
