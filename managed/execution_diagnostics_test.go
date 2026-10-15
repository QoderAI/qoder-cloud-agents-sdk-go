package managed_test

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"testing"
)

const managedDiagnosticMaxEvents = 128
const managedDiagnosticMaxBytes = 256 * 1024

var managedDiagnosticURL = regexp.MustCompile(`https?://[^\s"<>]+`)

// Keep public turn evidence before cleanup deletes the session. Limit the tail
// so noisy streams cannot exhaust test memory or hide the final failure in logs.
type managedTurnDiagnostics struct {
	secret  string
	events  []string
	bytes   int
	omitted int
}

func (d *managedTurnDiagnostics) text(s string) string {
	if d.secret != "" {
		s = strings.ReplaceAll(s, d.secret, "[REDACTED]")
	}
	s = managedDiagnosticURL.ReplaceAllString(s, "[URL OMITTED]")
	if len(s) > 4096 {
		s = s[:4096] + "[TRUNCATED]"
	}
	return s
}

func (d *managedTurnDiagnostics) observe(raw string) {
	var event map[string]any
	safe := map[string]any{}
	if json.Unmarshal([]byte(raw), &event) != nil || event == nil {
		safe["diagnostic"] = "invalid event JSON; payload omitted"
	} else {
		// Explicit fields exclude reasoning, binary content, auth headers and
		// unknown payload extensions. Text is from this test account only.
		for _, key := range []string{"id", "type", "processed_at", "name", "tool_use_id", "mcp_tool_use_id"} {
			if value, ok := event[key].(string); ok {
				safe[key] = d.text(value)
			}
		}
		if value, ok := event["is_error"].(bool); ok {
			safe["is_error"] = value
		}
		if reason, ok := event["stop_reason"].(map[string]any); ok {
			if value, ok := reason["type"].(string); ok {
				safe["stop_reason"] = d.text(value)
			}
		} else if value, ok := event["stop_reason"].(string); ok {
			safe["stop_reason"] = d.text(value)
		}
		if input, ok := event["input"].(map[string]any); ok {
			fields := map[string]string{}
			for _, key := range []string{"file_path", "path", "command"} {
				if value, ok := input[key].(string); ok {
					fields[key] = d.text(value)
				}
			}
			safe["input"] = fields
		}
		if content, ok := event["content"].([]any); ok {
			var texts []string
			for _, block := range content {
				if block, ok := block.(map[string]any); ok && block["type"] == "text" {
					if value, ok := block["text"].(string); ok {
						texts = append(texts, value)
					}
				}
			}
			if len(texts) > 0 {
				safe["text"] = d.text(strings.Join(texts, "\n"))
			}
		}
	}
	encoded, _ := json.Marshal(safe)
	d.events = append(d.events, string(encoded))
	d.bytes += len(encoded)
	for len(d.events) > managedDiagnosticMaxEvents || d.bytes > managedDiagnosticMaxBytes {
		d.bytes -= len(d.events[0])
		d.events[0] = ""
		d.events = d.events[1:]
		d.omitted++
	}
}

func TestManagedTurnDiagnostics(t *testing.T) {
	d := managedTurnDiagnostics{secret: "test-pat-secret"}
	for _, raw := range []string{
		`{"id":"tool-1","type":"agent.tool_use","name":"Read","input":{"file_path":"/data/.agent/awareness/sdk-e2e/proof.md","authorization":"secret-header"}}`,
		`{"id":"result-1","type":"agent.tool_result","tool_use_id":"tool-1","is_error":false,"content":[{"type":"text","text":"SDK_E2E_MEMORY_TOKEN=random-marker"}]}`,
		`{"id":"message-1","type":"agent.message","authorization":"secret-header","content":[{"type":"thinking","thinking":"private-reasoning"},{"type":"image","source":{"data":"private-image"}},{"type":"text","text":"Missing skill. test-pat-secret https://storage.example/file?signature=private-signature"}]}`,
		`{"id":"idle-1","type":"session.status_idle","stop_reason":{"type":"end_turn"}}`,
		`invalid payload test-pat-secret`,
	} {
		d.observe(raw)
	}
	output := strings.Join(d.events, "\n")
	for _, want := range []string{"tool-1", "result-1", "random-marker", "Missing skill.", "end_turn", "/data/.agent/awareness/sdk-e2e/proof.md", "invalid event JSON", "[REDACTED]", "[URL OMITTED]"} {
		if !strings.Contains(output, want) {
			t.Errorf("diagnostics lost %q", want)
		}
	}
	for _, secret := range []string{"test-pat-secret", "private-reasoning", "private-image", "private-signature", "secret-header"} {
		if strings.Contains(output, secret) {
			t.Errorf("diagnostics exposed %q", secret)
		}
	}
	for _, event := range d.events {
		if !json.Valid([]byte(event)) {
			t.Fatal("diagnostics must retain valid JSON")
		}
	}
}

func TestManagedTurnDiagnosticsBoundedTail(t *testing.T) {
	for _, size := range []int{1, 10000} {
		t.Run(fmt.Sprint(size), func(t *testing.T) {
			d := managedTurnDiagnostics{}
			for i := 0; i < 200; i++ {
				raw, _ := json.Marshal(map[string]any{"id": fmt.Sprint(i), "type": "agent.message", "content": []any{map[string]any{"type": "text", "text": strings.Repeat("x", size)}}})
				d.observe(string(raw))
			}
			d.observe(`{"id":"last","type":"session.status_idle","stop_reason":"end_turn"}`)
			if d.omitted == 0 || d.bytes > managedDiagnosticMaxBytes || len(d.events) > managedDiagnosticMaxEvents {
				t.Fatal("diagnostic tail exceeded its budget or did not report omitted events")
			}
			if !strings.Contains(d.events[len(d.events)-1], `"id":"last"`) {
				t.Fatal("diagnostic tail lost the terminal event")
			}
			if size > 4096 && !strings.Contains(strings.Join(d.events, "\n"), "[TRUNCATED]") {
				t.Fatal("oversized text was not marked as truncated")
			}
		})
	}
}
