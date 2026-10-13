package forward

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
	"time"
)

const RealtimeProtocolVersion = "voice.realtime.v1"

var realtimeTimestampPattern = regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}(?:\.\d{1,9})?(?:Z|[+-](?:[01]\d|2[0-3]):[0-5]\d)$`)

// RealtimeEvent is a server event. Payload and RawJSON preserve Beta extension
// fields; Timestamp preserves nanosecond precision. Sequence is not a cursor.
type RealtimeEvent struct {
	Version        string          `json:"version"`
	Type           string          `json:"type"`
	EventID        string          `json:"event_id"`
	Sequence       int64           `json:"sequence"`
	ConversationID string          `json:"conversation_id"`
	Timestamp      string          `json:"timestamp"`
	WorkID         string          `json:"work_id,omitempty"`
	AnnouncementID string          `json:"announcement_id,omitempty"`
	Payload        json.RawMessage `json:"payload"`
	raw            string
}

func (e RealtimeEvent) RawJSON() string { return e.raw }

// DecodePayload decodes a known typed payload or an application's extension.
func (e RealtimeEvent) DecodePayload(dst any) error { return json.Unmarshal(e.Payload, dst) }

// Audio decodes audio.delta PCM16. It does not decode WAV or play audio.
func (e RealtimeEvent) Audio() ([]byte, error) {
	if e.Type != "audio.delta" {
		return nil, fmt.Errorf("event is not audio.delta")
	}
	var p RealtimeAudioPayload
	if err := e.DecodePayload(&p); err != nil {
		return nil, err
	}
	pcm, err := base64.StdEncoding.DecodeString(p.Audio)
	if err != nil || len(pcm) == 0 || len(pcm)%2 != 0 {
		return nil, &RealtimeProtocolError{Message: "invalid PCM16 audio.delta"}
	}
	return pcm, nil
}

type RealtimeAudioFormat struct {
	Format     string `json:"format"`
	SampleRate int    `json:"sample_rate"`
	Channels   int    `json:"channels"`
}
type RealtimeCapabilities struct {
	GracefulClose bool `json:"graceful_close"`
}
type RealtimeReadyPayload struct {
	Config       RealtimeConfig       `json:"config"`
	Provider     string               `json:"provider"`
	InputAudio   RealtimeAudioFormat  `json:"input_audio"`
	OutputAudio  RealtimeAudioFormat  `json:"output_audio"`
	Capabilities RealtimeCapabilities `json:"capabilities"`
}
type RealtimeStatePayload struct {
	State string `json:"state"`
}
type RealtimeReasonPayload struct {
	Reason string `json:"reason"`
}
type RealtimeTranscriptPayload struct {
	Role       string `json:"role"`
	Text       string `json:"text"`
	ItemID     string `json:"item_id,omitempty"`
	ResponseID string `json:"response_id,omitempty"`
}
type RealtimeAudioPayload struct {
	Audio string `json:"audio"`
}
type RealtimeWorkAcceptedPayload struct {
	Status    string `json:"status"`
	Objective string `json:"objective"`
}
type RealtimeWorkRunningPayload struct {
	Status string `json:"status"`
}
type RealtimeWorkProgressPayload struct {
	Kind    string `json:"kind"`
	Detail  string `json:"detail"`
	Title   string `json:"title,omitempty"`
	IsError bool   `json:"is_error,omitempty"`
	EventID string `json:"event_id,omitempty"`
}
type RealtimeWorkMilestonePayload struct {
	AnnouncementID string `json:"announcement_id"`
	MilestoneType  string `json:"milestone_type"`
	Summary        string `json:"summary"`
}
type RealtimeWorkCompletedPayload struct {
	Status  string `json:"status"`
	Result  string `json:"result"`
	EventID string `json:"event_id,omitempty"`
}
type RealtimeWorkError struct {
	Code string `json:"code"`
}
type RealtimeWorkFailedPayload struct {
	Status  string            `json:"status"`
	Error   RealtimeWorkError `json:"error"`
	EventID string            `json:"event_id,omitempty"`
}
type RealtimeClosedPayload struct {
	RequestID string `json:"request_id"`
	Outcome   string `json:"outcome"`
}
type RealtimeErrorPayload struct {
	Code         string `json:"code"`
	Message      string `json:"message"`
	Retryable    *bool  `json:"retryable,omitempty"`
	RetryAfterMS *int64 `json:"retry_after_ms,omitempty"`
}

// RealtimeClientEvent is always encoded as a JSON text frame. Use SendEvent for
// protocol extensions; convenience methods provide validated public events.
type RealtimeClientEvent struct {
	Version string         `json:"version"`
	Type    string         `json:"type"`
	Payload map[string]any `json:"payload"`
}
type RealtimePlaybackReceipt struct {
	WorkID         string `json:"work_id,omitempty"`
	AnnouncementID string `json:"announcement_id,omitempty"`
}

type RealtimeProtocolError struct{ Message string }

func (e *RealtimeProtocolError) Error() string { return "qca realtime protocol: " + e.Message }

type RealtimeServerError struct {
	RealtimeErrorPayload
	ConversationID string
}

func (e *RealtimeServerError) Error() string {
	return fmt.Sprintf("qca realtime: %s: %s (conversation_id=%s)", e.Code, e.Message, e.ConversationID)
}

func decodeRealtimeEvent(b []byte, conversationID string) (RealtimeEvent, error) {
	var e RealtimeEvent
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(b, &fields); err != nil {
		return e, &RealtimeProtocolError{Message: "invalid JSON event"}
	}
	if fields == nil {
		return e, &RealtimeProtocolError{Message: "event must be an object"}
	}
	if err := json.Unmarshal(b, &e); err != nil {
		return e, &RealtimeProtocolError{Message: "invalid event envelope"}
	}
	if e.Version != RealtimeProtocolVersion || e.Type == "" || e.EventID == "" || e.ConversationID != conversationID || fields["sequence"] == nil || string(fields["sequence"]) == "null" || e.Sequence < 0 {
		return e, &RealtimeProtocolError{Message: "missing or invalid event envelope fields"}
	}
	if _, err := time.Parse(time.RFC3339Nano, e.Timestamp); err != nil || !realtimeTimestampPattern.MatchString(e.Timestamp) {
		return e, &RealtimeProtocolError{Message: "invalid event timestamp"}
	}
	var payload map[string]json.RawMessage
	if err := json.Unmarshal(e.Payload, &payload); err != nil || payload == nil {
		return e, &RealtimeProtocolError{Message: "event payload must be an object"}
	}
	required := map[string][]string{
		"voice.state": {"state"}, "voice.ready": {"config", "provider", "input_audio", "output_audio", "capabilities"}, "voice.replaced": {"reason"},
		"transcript.delta": {"role", "text"}, "transcript.final": {"role", "text"}, "audio.delta": {"audio"}, "playback.interrupt": {"reason"},
		"work.accepted": {"status", "objective"}, "work.running": {"status"}, "work.progress": {"kind", "detail"}, "work.milestone": {"announcement_id", "milestone_type", "summary"},
		"work.completed": {"status", "result"}, "work.failed": {"status", "error"}, "connection.closed": {"request_id", "outcome"}, "error": {"code", "message"},
	}
	for _, key := range required[e.Type] {
		if payload[key] == nil || string(payload[key]) == "null" {
			return e, &RealtimeProtocolError{Message: "missing " + e.Type + " payload field " + key}
		}
	}
	var typed any
	switch e.Type {
	case "voice.ready":
		typed = &RealtimeReadyPayload{}
	case "voice.state":
		typed = &RealtimeStatePayload{}
	case "voice.replaced", "playback.interrupt":
		typed = &RealtimeReasonPayload{}
	case "transcript.delta", "transcript.final":
		typed = &RealtimeTranscriptPayload{}
	case "audio.delta":
		typed = &RealtimeAudioPayload{}
	case "work.accepted":
		typed = &RealtimeWorkAcceptedPayload{}
	case "work.running":
		typed = &RealtimeWorkRunningPayload{}
	case "work.progress":
		typed = &RealtimeWorkProgressPayload{}
	case "work.milestone":
		typed = &RealtimeWorkMilestonePayload{}
	case "work.completed":
		typed = &RealtimeWorkCompletedPayload{}
	case "work.failed":
		typed = &RealtimeWorkFailedPayload{}
	case "connection.closed":
		typed = &RealtimeClosedPayload{}
	case "error":
		typed = &RealtimeErrorPayload{}
	}
	if typed != nil {
		if err := json.Unmarshal(e.Payload, typed); err != nil {
			return e, &RealtimeProtocolError{Message: "invalid " + e.Type + " payload types"}
		}
	}
	if e.Type == "work.failed" {
		var failure RealtimeWorkFailedPayload
		_ = json.Unmarshal(e.Payload, &failure)
		if failure.Error.Code == "" {
			return e, &RealtimeProtocolError{Message: "missing work.failed error code"}
		}
	}
	if expected, ok := map[string]string{"work.accepted": "accepted", "work.running": "running", "work.completed": "completed", "work.failed": "failed"}[e.Type]; ok {
		var status string
		if json.Unmarshal(payload["status"], &status) != nil || status != expected {
			return e, &RealtimeProtocolError{Message: "invalid " + e.Type + " status"}
		}
	}
	if e.Type == "error" {
		var failure RealtimeErrorPayload
		_ = json.Unmarshal(e.Payload, &failure)
		if failure.Code == "" || (failure.RetryAfterMS != nil && *failure.RetryAfterMS < 0) {
			return e, &RealtimeProtocolError{Message: "invalid error code or retry_after_ms"}
		}
	}
	e.raw = string(b)
	return e, nil
}
