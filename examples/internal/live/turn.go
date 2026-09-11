package live

import (
	"encoding/json"
	"fmt"
	"strings"
)

// TurnResult only accepts assistant output followed by a terminal idle event.
// User messages and tool output cannot satisfy the result assertion.
type TurnResult struct {
	Text     string
	LastID   string
	ToolUsed bool
	Complete bool
}

func (r *TurnResult) Observe(raw string) error {
	var event struct {
		ID, Type   string
		Content    json.RawMessage
		StopReason json.RawMessage `json:"stop_reason"`
	}
	if err := json.Unmarshal([]byte(raw), &event); err != nil {
		return fmt.Errorf("decode execution event: %w", err)
	}
	if event.ID != "" {
		r.LastID = event.ID
	}
	switch event.Type {
	case "session.error", "session.status_terminated":
		return fmt.Errorf("service execution failed: event=%s id=%s", event.Type, event.ID)
	case "agent.tool_use", "agent.mcp_tool_use":
		r.ToolUsed = true
	case "agent.message":
		var blocks []struct{ Type, Text string }
		if err := json.Unmarshal(event.Content, &blocks); err != nil {
			return fmt.Errorf("decode agent message: %w", err)
		}
		r.Text = ""
		for _, b := range blocks {
			if b.Type == "text" {
				r.Text += b.Text + "\n"
			}
		}
	case "session.status_idle":
		var reason string
		if json.Unmarshal(event.StopReason, &reason) != nil {
			var obj struct{ Type string }
			_ = json.Unmarshal(event.StopReason, &obj)
			reason = obj.Type
		}
		if reason != "" && reason != "end_turn" && reason != "stop_sequence" {
			return fmt.Errorf("execution stopped without completing: reason=%s id=%s", reason, event.ID)
		}
		if r.Text != "" {
			r.Complete = true
		}
	}
	return nil
}
func (r TurnResult) Verify(expected []string, requireTool bool) error {
	if !r.Complete {
		return fmt.Errorf("execution did not reach idle after assistant output; last_event_id=%s", r.LastID)
	}
	for _, token := range expected {
		if !strings.Contains(r.Text, token) {
			return fmt.Errorf("assistant output is missing an expected marker; last_event_id=%s", r.LastID)
		}
	}
	if requireTool && !r.ToolUsed {
		return fmt.Errorf("no tool execution observed; last_event_id=%s", r.LastID)
	}
	return nil
}
