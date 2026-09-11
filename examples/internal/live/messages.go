package live

import (
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// HistoryMessage extracts only user and assistant text, excluding tool inputs,
// results and other event payloads from the conversation transcript.
func HistoryMessage(raw string) (role, text string, err error) {
	var event struct {
		Type    string
		Content json.RawMessage
	}
	if err = json.Unmarshal([]byte(raw), &event); err != nil {
		return "", "", err
	}
	switch event.Type {
	case "user.message":
		role = "用户"
	case "agent.message":
		role = "助手"
	default:
		return "", "", nil
	}
	text, err = messageText(event.Content)
	return
}

func messageText(raw json.RawMessage) (string, error) {
	var text string
	if json.Unmarshal(raw, &text) == nil {
		return text, nil
	}
	var blocks []struct{ Type, Text string }
	if err := json.Unmarshal(raw, &blocks); err != nil {
		return "", err
	}
	var parts []string
	for _, block := range blocks {
		if block.Type == "text" {
			parts = append(parts, block.Text)
		}
	}
	return strings.Join(parts, "\n"), nil
}

// MessagePreview tracks text by event ID and content-block index. A final
// agent.message replaces its preview, including when some deltas were missed.
type MessagePreview struct {
	blocks map[string]map[int]string
}

type PreviewUpdate struct {
	ID, Text string
	Final    bool
}

func (p *MessagePreview) Observe(raw string) (*PreviewUpdate, error) {
	var event struct {
		ID, Type string
		EventID  string `json:"event_id"`
		Event    struct{ ID, Type string }
		Content  json.RawMessage
		Delta    struct {
			Type    string
			Index   int
			Content struct{ Type, Text string }
		}
	}
	if err := json.Unmarshal([]byte(raw), &event); err != nil {
		return nil, err
	}
	if p.blocks == nil {
		p.blocks = map[string]map[int]string{}
	}
	switch event.Type {
	case "event_start":
		if event.Event.Type == "agent.message" && p.blocks[event.Event.ID] == nil {
			p.blocks[event.Event.ID] = map[int]string{}
		}
	case "event_delta":
		if event.Delta.Type != "content_delta" || event.Delta.Content.Type != "text" {
			return nil, nil
		}
		if event.EventID == "" || event.Delta.Index < 0 {
			return nil, fmt.Errorf("invalid text delta: missing event_id or negative block index")
		}
		blocks := p.blocks[event.EventID]
		if blocks == nil {
			blocks = map[int]string{}
			p.blocks[event.EventID] = blocks
		}
		// Several fragments share an event ID; deduplicating by ID loses text.
		blocks[event.Delta.Index] += event.Delta.Content.Text
		indices := make([]int, 0, len(blocks))
		for index := range blocks {
			indices = append(indices, index)
		}
		sort.Ints(indices)
		parts := make([]string, 0, len(indices))
		for _, index := range indices {
			parts = append(parts, blocks[index])
		}
		return &PreviewUpdate{ID: event.EventID, Text: strings.Join(parts, "\n")}, nil
	case "agent.message":
		text, err := messageText(event.Content)
		if err != nil {
			return nil, err
		}
		delete(p.blocks, event.ID)
		return &PreviewUpdate{ID: event.ID, Text: text, Final: true}, nil
	}
	return nil, nil
}
