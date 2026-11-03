package managed_test

import (
	"encoding/json"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
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
