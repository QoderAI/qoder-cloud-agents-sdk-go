// Qoder managed API definitions.
package managed

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"time"

	requestconfig "github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apiquery"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/pagination"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/ssestream"
)

// SessionEventService contains methods and other services that help with
// interacting with the Qoder Cloud Agents API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSessionEventService] method instead.
type SessionEventService struct {
	Options []option.RequestOption
}

// NewSessionEventService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewSessionEventService(opts ...option.RequestOption) (r SessionEventService) {
	r = SessionEventService{}
	r.Options = opts
	return
}

// List Events
func (r *SessionEventService) List(ctx context.Context, sessionID string, params SessionEventListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ManagedAgentsSessionEventUnion], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("sessions/%s/events", url.PathEscape(sessionID))
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, params, &res, opts...)
	if err != nil {
		return nil, err
	}
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

// List Events
func (r *SessionEventService) ListAutoPaging(ctx context.Context, sessionID string, params SessionEventListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ManagedAgentsSessionEventUnion] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, sessionID, params, opts...))
}

// Send Events
func (r *SessionEventService) Send(ctx context.Context, sessionID string, params SessionEventSendParams, opts ...option.RequestOption) (res *ManagedAgentsSendSessionEvents, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("sessions/%s/events", url.PathEscape(sessionID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Stream Events
func (r *SessionEventService) StreamEvents(ctx context.Context, sessionID string, params SessionEventStreamParams, opts ...option.RequestOption) (stream *ssestream.Stream[ManagedAgentsStreamSessionEventsUnion]) {
	var (
		raw *http.Response
		err error
	)
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, []option.RequestOption{option.WithHeader("Accept", "text/event-stream")}, opts)

	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return ssestream.NewStream[ManagedAgentsStreamSessionEventsUnion](nil, err)
	}
	path := fmt.Sprintf("sessions/%s/events/stream", url.PathEscape(sessionID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &raw, opts...)
	return ssestream.NewStream[ManagedAgentsStreamSessionEventsUnion](ssestream.NewDecoder(raw), err)
}

// Event emitted when the agent calls a custom tool. The session goes idle until
// the client sends a `user.custom_tool_result` event with the result.
type ManagedAgentsAgentCustomToolUseEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// Input parameters for the tool call.
	Input map[string]any `json:"input" api:"required"`
	// Name of the custom tool being called.
	Name string `json:"name" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "agent.custom_tool_use".
	Type ManagedAgentsAgentCustomToolUseEventType `json:"type" api:"required"`
	// When set, this event was cross-posted from a subagent's thread to surface its
	// custom tool use on the primary thread's stream. Empty on the thread's own
	// events. Echo this on a `user.custom_tool_result` event to route the result back.
	SessionThreadID string `json:"session_thread_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID              respjson.Field
		Input           respjson.Field
		Name            respjson.Field
		ProcessedAt     respjson.Field
		Type            respjson.Field
		SessionThreadID respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsAgentCustomToolUseEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsAgentCustomToolUseEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsAgentCustomToolUseEventType string

const (
	ManagedAgentsAgentCustomToolUseEventTypeAgentCustomToolUse ManagedAgentsAgentCustomToolUseEventType = "agent.custom_tool_use"
)

// Event representing the result of an MCP tool execution.
type ManagedAgentsAgentMCPToolResultEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// The id of the `agent.mcp_tool_use` event this result corresponds to.
	MCPToolUseID string `json:"mcp_tool_use_id" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "agent.mcp_tool_result".
	Type ManagedAgentsAgentMCPToolResultEventType `json:"type" api:"required"`
	// The result content returned by the tool.
	Content []ManagedAgentsAgentMCPToolResultEventContentUnion `json:"content"`
	// Whether the tool execution resulted in an error.
	IsError bool `json:"is_error" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID           respjson.Field
		MCPToolUseID respjson.Field
		ProcessedAt  respjson.Field
		Type         respjson.Field
		Content      respjson.Field
		IsError      respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsAgentMCPToolResultEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsAgentMCPToolResultEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsAgentMCPToolResultEventType string

const (
	ManagedAgentsAgentMCPToolResultEventTypeAgentMCPToolResult ManagedAgentsAgentMCPToolResultEventType = "agent.mcp_tool_result"
)

// ManagedAgentsAgentMCPToolResultEventContentUnion contains all possible
// properties and values from [ManagedAgentsTextBlock],
// [ManagedAgentsImageBlock], [ManagedAgentsDocumentBlock],
// [ManagedAgentsSearchResultBlock].
//
// Use the [ManagedAgentsAgentMCPToolResultEventContentUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsAgentMCPToolResultEventContentUnion struct {
	// This field is from variant [ManagedAgentsTextBlock].
	Text string `json:"text"`
	// Any of "text", "image", "document", "search_result".
	Type string `json:"type"`
	// This field is a union of [ManagedAgentsImageBlockSourceUnion],
	// [ManagedAgentsDocumentBlockSourceUnion], [string]
	Source ManagedAgentsAgentMCPToolResultEventContentUnionSource `json:"source"`
	// This field is from variant [ManagedAgentsDocumentBlock].
	Context string `json:"context"`
	Title   string `json:"title"`
	// This field is from variant [ManagedAgentsSearchResultBlock].
	Citations ManagedAgentsSearchResultCitations `json:"citations"`
	// This field is from variant [ManagedAgentsSearchResultBlock].
	Content []ManagedAgentsSearchResultContent `json:"content"`
	JSON    struct {
		Text      respjson.Field
		Type      respjson.Field
		Source    respjson.Field
		Context   respjson.Field
		Title     respjson.Field
		Citations respjson.Field
		Content   respjson.Field
		raw       string
	} `json:"-"`
}

// anyManagedAgentsAgentMCPToolResultEventContent is implemented by each
// variant of [ManagedAgentsAgentMCPToolResultEventContentUnion] to add type
// safety for the return type of
// [ManagedAgentsAgentMCPToolResultEventContentUnion.AsAny]
type anyManagedAgentsAgentMCPToolResultEventContent interface {
	implManagedAgentsAgentMCPToolResultEventContentUnion()
}

func (ManagedAgentsTextBlock) implManagedAgentsAgentMCPToolResultEventContentUnion()     {}
func (ManagedAgentsImageBlock) implManagedAgentsAgentMCPToolResultEventContentUnion()    {}
func (ManagedAgentsDocumentBlock) implManagedAgentsAgentMCPToolResultEventContentUnion() {}
func (ManagedAgentsSearchResultBlock) implManagedAgentsAgentMCPToolResultEventContentUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsAgentMCPToolResultEventContentUnion.AsAny().(type) {
//	case qoder.ManagedAgentsTextBlock:
//	case qoder.ManagedAgentsImageBlock:
//	case qoder.ManagedAgentsDocumentBlock:
//	case qoder.ManagedAgentsSearchResultBlock:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsAgentMCPToolResultEventContentUnion) AsAny() anyManagedAgentsAgentMCPToolResultEventContent {
	switch u.Type {
	case "text":
		return u.AsText()
	case "image":
		return u.AsImage()
	case "document":
		return u.AsDocument()
	case "search_result":
		return u.AsSearchResult()
	}
	return nil
}

func (u ManagedAgentsAgentMCPToolResultEventContentUnion) AsText() (v ManagedAgentsTextBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsAgentMCPToolResultEventContentUnion) AsImage() (v ManagedAgentsImageBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsAgentMCPToolResultEventContentUnion) AsDocument() (v ManagedAgentsDocumentBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsAgentMCPToolResultEventContentUnion) AsSearchResult() (v ManagedAgentsSearchResultBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsAgentMCPToolResultEventContentUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsAgentMCPToolResultEventContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsAgentMCPToolResultEventContentUnionSource is an implicit
// subunion of [ManagedAgentsAgentMCPToolResultEventContentUnion].
// ManagedAgentsAgentMCPToolResultEventContentUnionSource provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsAgentMCPToolResultEventContentUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString]
type ManagedAgentsAgentMCPToolResultEventContentUnionSource struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString  string `json:",inline"`
	Data      string `json:"data"`
	MediaType string `json:"media_type"`
	Type      string `json:"type"`
	URL       string `json:"url"`
	FileID    string `json:"file_id"`
	JSON      struct {
		OfString  respjson.Field
		Data      respjson.Field
		MediaType respjson.Field
		Type      respjson.Field
		URL       respjson.Field
		FileID    respjson.Field
		raw       string
	} `json:"-"`
}

func (r *ManagedAgentsAgentMCPToolResultEventContentUnionSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Event emitted when the agent invokes a tool provided by an MCP server.
type ManagedAgentsAgentMCPToolUseEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// Input parameters for the tool call.
	Input map[string]any `json:"input" api:"required"`
	// Name of the MCP server providing the tool.
	MCPServerName string `json:"mcp_server_name" api:"required"`
	// Name of the MCP tool being used.
	Name string `json:"name" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "agent.mcp_tool_use".
	Type ManagedAgentsAgentMCPToolUseEventType `json:"type" api:"required"`
	// AgentEvaluatedPermission enum
	//
	// Any of "allow", "ask", "deny".
	EvaluatedPermission ManagedAgentsAgentMCPToolUseEventEvaluatedPermission `json:"evaluated_permission"`
	// When set, this event was cross-posted from a subagent's thread to surface its
	// permission request on the primary thread's stream. Empty on the thread's own
	// events. Echo this on a `user.tool_confirmation` event to route the approval
	// back.
	SessionThreadID string `json:"session_thread_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                  respjson.Field
		Input               respjson.Field
		MCPServerName       respjson.Field
		Name                respjson.Field
		ProcessedAt         respjson.Field
		Type                respjson.Field
		EvaluatedPermission respjson.Field
		SessionThreadID     respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsAgentMCPToolUseEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsAgentMCPToolUseEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsAgentMCPToolUseEventType string

const (
	ManagedAgentsAgentMCPToolUseEventTypeAgentMCPToolUse ManagedAgentsAgentMCPToolUseEventType = "agent.mcp_tool_use"
)

// AgentEvaluatedPermission enum
type ManagedAgentsAgentMCPToolUseEventEvaluatedPermission string

const (
	ManagedAgentsAgentMCPToolUseEventEvaluatedPermissionAllow ManagedAgentsAgentMCPToolUseEventEvaluatedPermission = "allow"
	ManagedAgentsAgentMCPToolUseEventEvaluatedPermissionAsk   ManagedAgentsAgentMCPToolUseEventEvaluatedPermission = "ask"
	ManagedAgentsAgentMCPToolUseEventEvaluatedPermissionDeny  ManagedAgentsAgentMCPToolUseEventEvaluatedPermission = "deny"
)

// An agent response event in the session conversation.
type ManagedAgentsAgentMessageEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// Array of text blocks comprising the agent response.
	Content []ManagedAgentsAgentMessageEventContentUnion `json:"content" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "agent.message".
	Type ManagedAgentsAgentMessageEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Content     respjson.Field
		ProcessedAt respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsAgentMessageEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsAgentMessageEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsAgentMessageEventContentUnion contains all possible properties
// and values from [ManagedAgentsTextBlock], [ManagedAgentsRedactedBlock].
//
// Use the [ManagedAgentsAgentMessageEventContentUnion.AsAny] method to switch
// on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsAgentMessageEventContentUnion struct {
	// This field is from variant [ManagedAgentsTextBlock].
	Text string `json:"text"`
	// Any of "text", "redacted".
	Type string `json:"type"`
	JSON struct {
		Text respjson.Field
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyManagedAgentsAgentMessageEventContent is implemented by each variant of
// [ManagedAgentsAgentMessageEventContentUnion] to add type safety for the
// return type of [ManagedAgentsAgentMessageEventContentUnion.AsAny]
type anyManagedAgentsAgentMessageEventContent interface {
	implManagedAgentsAgentMessageEventContentUnion()
}

func (ManagedAgentsTextBlock) implManagedAgentsAgentMessageEventContentUnion()     {}
func (ManagedAgentsRedactedBlock) implManagedAgentsAgentMessageEventContentUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsAgentMessageEventContentUnion.AsAny().(type) {
//	case qoder.ManagedAgentsTextBlock:
//	case qoder.ManagedAgentsRedactedBlock:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsAgentMessageEventContentUnion) AsAny() anyManagedAgentsAgentMessageEventContent {
	switch u.Type {
	case "text":
		return u.AsText()
	case "redacted":
		return u.AsRedacted()
	}
	return nil
}

func (u ManagedAgentsAgentMessageEventContentUnion) AsText() (v ManagedAgentsTextBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsAgentMessageEventContentUnion) AsRedacted() (v ManagedAgentsRedactedBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsAgentMessageEventContentUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsAgentMessageEventContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsAgentMessageEventType string

const (
	ManagedAgentsAgentMessageEventTypeAgentMessage ManagedAgentsAgentMessageEventType = "agent.message"
)

// Indicates the agent is making forward progress via extended thinking. A progress
// signal, not a content carrier.
type ManagedAgentsAgentThinkingEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "agent.thinking".
	Type ManagedAgentsAgentThinkingEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ProcessedAt respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsAgentThinkingEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsAgentThinkingEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsAgentThinkingEventType string

const (
	ManagedAgentsAgentThinkingEventTypeAgentThinking ManagedAgentsAgentThinkingEventType = "agent.thinking"
)

// Indicates that context compaction (summarization) occurred during the session.
type ManagedAgentsAgentThreadContextCompactedEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "agent.thread_context_compacted".
	Type ManagedAgentsAgentThreadContextCompactedEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ProcessedAt respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsAgentThreadContextCompactedEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsAgentThreadContextCompactedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsAgentThreadContextCompactedEventType string

const (
	ManagedAgentsAgentThreadContextCompactedEventTypeAgentThreadContextCompacted ManagedAgentsAgentThreadContextCompactedEventType = "agent.thread_context_compacted"
)

// Delivery event written to the target thread's input stream when an
// agent-to-agent message arrives.
type ManagedAgentsAgentThreadMessageReceivedEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// Message content blocks.
	Content []ManagedAgentsAgentThreadMessageReceivedEventContentUnion `json:"content" api:"required"`
	// Public `sthr_` ID of the thread that sent the message.
	FromSessionThreadID string `json:"from_session_thread_id" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "agent.thread_message_received".
	Type ManagedAgentsAgentThreadMessageReceivedEventType `json:"type" api:"required"`
	// Name of the callable agent this message came from. Absent when received from the
	// primary agent.
	FromAgentName string `json:"from_agent_name" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                  respjson.Field
		Content             respjson.Field
		FromSessionThreadID respjson.Field
		ProcessedAt         respjson.Field
		Type                respjson.Field
		FromAgentName       respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsAgentThreadMessageReceivedEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsAgentThreadMessageReceivedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsAgentThreadMessageReceivedEventContentUnion contains all
// possible properties and values from [ManagedAgentsTextBlock],
// [ManagedAgentsImageBlock], [ManagedAgentsDocumentBlock],
// [ManagedAgentsRedactedBlock].
//
// Use the [ManagedAgentsAgentThreadMessageReceivedEventContentUnion.AsAny]
// method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsAgentThreadMessageReceivedEventContentUnion struct {
	// This field is from variant [ManagedAgentsTextBlock].
	Text string `json:"text"`
	// Any of "text", "image", "document", "redacted".
	Type string `json:"type"`
	// This field is a union of [ManagedAgentsImageBlockSourceUnion],
	// [ManagedAgentsDocumentBlockSourceUnion]
	Source ManagedAgentsAgentThreadMessageReceivedEventContentUnionSource `json:"source"`
	// This field is from variant [ManagedAgentsDocumentBlock].
	Context string `json:"context"`
	// This field is from variant [ManagedAgentsDocumentBlock].
	Title string `json:"title"`
	JSON  struct {
		Text    respjson.Field
		Type    respjson.Field
		Source  respjson.Field
		Context respjson.Field
		Title   respjson.Field
		raw     string
	} `json:"-"`
}

// anyManagedAgentsAgentThreadMessageReceivedEventContent is implemented by
// each variant of [ManagedAgentsAgentThreadMessageReceivedEventContentUnion]
// to add type safety for the return type of
// [ManagedAgentsAgentThreadMessageReceivedEventContentUnion.AsAny]
type anyManagedAgentsAgentThreadMessageReceivedEventContent interface {
	implManagedAgentsAgentThreadMessageReceivedEventContentUnion()
}

func (ManagedAgentsTextBlock) implManagedAgentsAgentThreadMessageReceivedEventContentUnion() {
}
func (ManagedAgentsImageBlock) implManagedAgentsAgentThreadMessageReceivedEventContentUnion() {
}
func (ManagedAgentsDocumentBlock) implManagedAgentsAgentThreadMessageReceivedEventContentUnion() {
}
func (ManagedAgentsRedactedBlock) implManagedAgentsAgentThreadMessageReceivedEventContentUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsAgentThreadMessageReceivedEventContentUnion.AsAny().(type) {
//	case qoder.ManagedAgentsTextBlock:
//	case qoder.ManagedAgentsImageBlock:
//	case qoder.ManagedAgentsDocumentBlock:
//	case qoder.ManagedAgentsRedactedBlock:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsAgentThreadMessageReceivedEventContentUnion) AsAny() anyManagedAgentsAgentThreadMessageReceivedEventContent {
	switch u.Type {
	case "text":
		return u.AsText()
	case "image":
		return u.AsImage()
	case "document":
		return u.AsDocument()
	case "redacted":
		return u.AsRedacted()
	}
	return nil
}

func (u ManagedAgentsAgentThreadMessageReceivedEventContentUnion) AsText() (v ManagedAgentsTextBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsAgentThreadMessageReceivedEventContentUnion) AsImage() (v ManagedAgentsImageBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsAgentThreadMessageReceivedEventContentUnion) AsDocument() (v ManagedAgentsDocumentBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsAgentThreadMessageReceivedEventContentUnion) AsRedacted() (v ManagedAgentsRedactedBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsAgentThreadMessageReceivedEventContentUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *ManagedAgentsAgentThreadMessageReceivedEventContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsAgentThreadMessageReceivedEventContentUnionSource is an
// implicit subunion of
// [ManagedAgentsAgentThreadMessageReceivedEventContentUnion].
// ManagedAgentsAgentThreadMessageReceivedEventContentUnionSource provides
// convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsAgentThreadMessageReceivedEventContentUnion].
type ManagedAgentsAgentThreadMessageReceivedEventContentUnionSource struct {
	Data      string `json:"data"`
	MediaType string `json:"media_type"`
	Type      string `json:"type"`
	URL       string `json:"url"`
	FileID    string `json:"file_id"`
	JSON      struct {
		Data      respjson.Field
		MediaType respjson.Field
		Type      respjson.Field
		URL       respjson.Field
		FileID    respjson.Field
		raw       string
	} `json:"-"`
}

func (r *ManagedAgentsAgentThreadMessageReceivedEventContentUnionSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsAgentThreadMessageReceivedEventType string

const (
	ManagedAgentsAgentThreadMessageReceivedEventTypeAgentThreadMessageReceived ManagedAgentsAgentThreadMessageReceivedEventType = "agent.thread_message_received"
)

// Observability event emitted to the sender's output stream when an agent-to-agent
// message is sent.
type ManagedAgentsAgentThreadMessageSentEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// Message content blocks.
	Content []ManagedAgentsAgentThreadMessageSentEventContentUnion `json:"content" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Public `sthr_` ID of the thread the message was sent to.
	ToSessionThreadID string `json:"to_session_thread_id" api:"required"`
	// Any of "agent.thread_message_sent".
	Type ManagedAgentsAgentThreadMessageSentEventType `json:"type" api:"required"`
	// Name of the callable agent this message was sent to. Absent when sent to the
	// primary agent.
	ToAgentName string `json:"to_agent_name" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                respjson.Field
		Content           respjson.Field
		ProcessedAt       respjson.Field
		ToSessionThreadID respjson.Field
		Type              respjson.Field
		ToAgentName       respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsAgentThreadMessageSentEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsAgentThreadMessageSentEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsAgentThreadMessageSentEventContentUnion contains all possible
// properties and values from [ManagedAgentsTextBlock],
// [ManagedAgentsImageBlock], [ManagedAgentsDocumentBlock],
// [ManagedAgentsRedactedBlock].
//
// Use the [ManagedAgentsAgentThreadMessageSentEventContentUnion.AsAny] method
// to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsAgentThreadMessageSentEventContentUnion struct {
	// This field is from variant [ManagedAgentsTextBlock].
	Text string `json:"text"`
	// Any of "text", "image", "document", "redacted".
	Type string `json:"type"`
	// This field is a union of [ManagedAgentsImageBlockSourceUnion],
	// [ManagedAgentsDocumentBlockSourceUnion]
	Source ManagedAgentsAgentThreadMessageSentEventContentUnionSource `json:"source"`
	// This field is from variant [ManagedAgentsDocumentBlock].
	Context string `json:"context"`
	// This field is from variant [ManagedAgentsDocumentBlock].
	Title string `json:"title"`
	JSON  struct {
		Text    respjson.Field
		Type    respjson.Field
		Source  respjson.Field
		Context respjson.Field
		Title   respjson.Field
		raw     string
	} `json:"-"`
}

// anyManagedAgentsAgentThreadMessageSentEventContent is implemented by each
// variant of [ManagedAgentsAgentThreadMessageSentEventContentUnion] to add
// type safety for the return type of
// [ManagedAgentsAgentThreadMessageSentEventContentUnion.AsAny]
type anyManagedAgentsAgentThreadMessageSentEventContent interface {
	implManagedAgentsAgentThreadMessageSentEventContentUnion()
}

func (ManagedAgentsTextBlock) implManagedAgentsAgentThreadMessageSentEventContentUnion()  {}
func (ManagedAgentsImageBlock) implManagedAgentsAgentThreadMessageSentEventContentUnion() {}
func (ManagedAgentsDocumentBlock) implManagedAgentsAgentThreadMessageSentEventContentUnion() {
}
func (ManagedAgentsRedactedBlock) implManagedAgentsAgentThreadMessageSentEventContentUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsAgentThreadMessageSentEventContentUnion.AsAny().(type) {
//	case qoder.ManagedAgentsTextBlock:
//	case qoder.ManagedAgentsImageBlock:
//	case qoder.ManagedAgentsDocumentBlock:
//	case qoder.ManagedAgentsRedactedBlock:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsAgentThreadMessageSentEventContentUnion) AsAny() anyManagedAgentsAgentThreadMessageSentEventContent {
	switch u.Type {
	case "text":
		return u.AsText()
	case "image":
		return u.AsImage()
	case "document":
		return u.AsDocument()
	case "redacted":
		return u.AsRedacted()
	}
	return nil
}

func (u ManagedAgentsAgentThreadMessageSentEventContentUnion) AsText() (v ManagedAgentsTextBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsAgentThreadMessageSentEventContentUnion) AsImage() (v ManagedAgentsImageBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsAgentThreadMessageSentEventContentUnion) AsDocument() (v ManagedAgentsDocumentBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsAgentThreadMessageSentEventContentUnion) AsRedacted() (v ManagedAgentsRedactedBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsAgentThreadMessageSentEventContentUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsAgentThreadMessageSentEventContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsAgentThreadMessageSentEventContentUnionSource is an implicit
// subunion of [ManagedAgentsAgentThreadMessageSentEventContentUnion].
// ManagedAgentsAgentThreadMessageSentEventContentUnionSource provides
// convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsAgentThreadMessageSentEventContentUnion].
type ManagedAgentsAgentThreadMessageSentEventContentUnionSource struct {
	Data      string `json:"data"`
	MediaType string `json:"media_type"`
	Type      string `json:"type"`
	URL       string `json:"url"`
	FileID    string `json:"file_id"`
	JSON      struct {
		Data      respjson.Field
		MediaType respjson.Field
		Type      respjson.Field
		URL       respjson.Field
		FileID    respjson.Field
		raw       string
	} `json:"-"`
}

func (r *ManagedAgentsAgentThreadMessageSentEventContentUnionSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsAgentThreadMessageSentEventType string

const (
	ManagedAgentsAgentThreadMessageSentEventTypeAgentThreadMessageSent ManagedAgentsAgentThreadMessageSentEventType = "agent.thread_message_sent"
)

// Event representing the result of an agent tool execution.
type ManagedAgentsAgentToolResultEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// The id of the `agent.tool_use` event this result corresponds to.
	ToolUseID string `json:"tool_use_id" api:"required"`
	// Any of "agent.tool_result".
	Type ManagedAgentsAgentToolResultEventType `json:"type" api:"required"`
	// The result content returned by the tool.
	Content []ManagedAgentsAgentToolResultEventContentUnion `json:"content"`
	// Whether the tool execution resulted in an error.
	IsError bool `json:"is_error" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ProcessedAt respjson.Field
		ToolUseID   respjson.Field
		Type        respjson.Field
		Content     respjson.Field
		IsError     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsAgentToolResultEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsAgentToolResultEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsAgentToolResultEventType string

const (
	ManagedAgentsAgentToolResultEventTypeAgentToolResult ManagedAgentsAgentToolResultEventType = "agent.tool_result"
)

// ManagedAgentsAgentToolResultEventContentUnion contains all possible
// properties and values from [ManagedAgentsTextBlock],
// [ManagedAgentsImageBlock], [ManagedAgentsDocumentBlock],
// [ManagedAgentsSearchResultBlock].
//
// Use the [ManagedAgentsAgentToolResultEventContentUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsAgentToolResultEventContentUnion struct {
	// This field is from variant [ManagedAgentsTextBlock].
	Text string `json:"text"`
	// Any of "text", "image", "document", "search_result".
	Type string `json:"type"`
	// This field is a union of [ManagedAgentsImageBlockSourceUnion],
	// [ManagedAgentsDocumentBlockSourceUnion], [string]
	Source ManagedAgentsAgentToolResultEventContentUnionSource `json:"source"`
	// This field is from variant [ManagedAgentsDocumentBlock].
	Context string `json:"context"`
	Title   string `json:"title"`
	// This field is from variant [ManagedAgentsSearchResultBlock].
	Citations ManagedAgentsSearchResultCitations `json:"citations"`
	// This field is from variant [ManagedAgentsSearchResultBlock].
	Content []ManagedAgentsSearchResultContent `json:"content"`
	JSON    struct {
		Text      respjson.Field
		Type      respjson.Field
		Source    respjson.Field
		Context   respjson.Field
		Title     respjson.Field
		Citations respjson.Field
		Content   respjson.Field
		raw       string
	} `json:"-"`
}

// anyManagedAgentsAgentToolResultEventContent is implemented by each variant
// of [ManagedAgentsAgentToolResultEventContentUnion] to add type safety for
// the return type of [ManagedAgentsAgentToolResultEventContentUnion.AsAny]
type anyManagedAgentsAgentToolResultEventContent interface {
	implManagedAgentsAgentToolResultEventContentUnion()
}

func (ManagedAgentsTextBlock) implManagedAgentsAgentToolResultEventContentUnion()         {}
func (ManagedAgentsImageBlock) implManagedAgentsAgentToolResultEventContentUnion()        {}
func (ManagedAgentsDocumentBlock) implManagedAgentsAgentToolResultEventContentUnion()     {}
func (ManagedAgentsSearchResultBlock) implManagedAgentsAgentToolResultEventContentUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsAgentToolResultEventContentUnion.AsAny().(type) {
//	case qoder.ManagedAgentsTextBlock:
//	case qoder.ManagedAgentsImageBlock:
//	case qoder.ManagedAgentsDocumentBlock:
//	case qoder.ManagedAgentsSearchResultBlock:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsAgentToolResultEventContentUnion) AsAny() anyManagedAgentsAgentToolResultEventContent {
	switch u.Type {
	case "text":
		return u.AsText()
	case "image":
		return u.AsImage()
	case "document":
		return u.AsDocument()
	case "search_result":
		return u.AsSearchResult()
	}
	return nil
}

func (u ManagedAgentsAgentToolResultEventContentUnion) AsText() (v ManagedAgentsTextBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsAgentToolResultEventContentUnion) AsImage() (v ManagedAgentsImageBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsAgentToolResultEventContentUnion) AsDocument() (v ManagedAgentsDocumentBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsAgentToolResultEventContentUnion) AsSearchResult() (v ManagedAgentsSearchResultBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsAgentToolResultEventContentUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsAgentToolResultEventContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsAgentToolResultEventContentUnionSource is an implicit subunion
// of [ManagedAgentsAgentToolResultEventContentUnion].
// ManagedAgentsAgentToolResultEventContentUnionSource provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsAgentToolResultEventContentUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString]
type ManagedAgentsAgentToolResultEventContentUnionSource struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString  string `json:",inline"`
	Data      string `json:"data"`
	MediaType string `json:"media_type"`
	Type      string `json:"type"`
	URL       string `json:"url"`
	FileID    string `json:"file_id"`
	JSON      struct {
		OfString  respjson.Field
		Data      respjson.Field
		MediaType respjson.Field
		Type      respjson.Field
		URL       respjson.Field
		FileID    respjson.Field
		raw       string
	} `json:"-"`
}

func (r *ManagedAgentsAgentToolResultEventContentUnionSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Event emitted when the agent invokes a built-in agent tool.
type ManagedAgentsAgentToolUseEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// Input parameters for the tool call.
	Input map[string]any `json:"input" api:"required"`
	// Name of the agent tool being used.
	Name string `json:"name" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "agent.tool_use".
	Type ManagedAgentsAgentToolUseEventType `json:"type" api:"required"`
	// AgentEvaluatedPermission enum
	//
	// Any of "allow", "ask", "deny".
	EvaluatedPermission ManagedAgentsAgentToolUseEventEvaluatedPermission `json:"evaluated_permission"`
	// When set, this event was cross-posted from a subagent's thread to surface its
	// permission request on the primary thread's stream. Empty on the thread's own
	// events. Echo this on a `user.tool_confirmation` event to route the approval
	// back.
	SessionThreadID string `json:"session_thread_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                  respjson.Field
		Input               respjson.Field
		Name                respjson.Field
		ProcessedAt         respjson.Field
		Type                respjson.Field
		EvaluatedPermission respjson.Field
		SessionThreadID     respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsAgentToolUseEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsAgentToolUseEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsAgentToolUseEventType string

const (
	ManagedAgentsAgentToolUseEventTypeAgentToolUse ManagedAgentsAgentToolUseEventType = "agent.tool_use"
)

// AgentEvaluatedPermission enum
type ManagedAgentsAgentToolUseEventEvaluatedPermission string

const (
	ManagedAgentsAgentToolUseEventEvaluatedPermissionAllow ManagedAgentsAgentToolUseEventEvaluatedPermission = "allow"
	ManagedAgentsAgentToolUseEventEvaluatedPermissionAsk   ManagedAgentsAgentToolUseEventEvaluatedPermission = "ask"
	ManagedAgentsAgentToolUseEventEvaluatedPermissionDeny  ManagedAgentsAgentToolUseEventEvaluatedPermission = "deny"
)

// Base64-encoded document data.
type ManagedAgentsBase64DocumentSource struct {
	// Base64-encoded document data.
	Data string `json:"data" api:"required"`
	// MIME type of the document (e.g., "application/pdf").
	MediaType string `json:"media_type" api:"required"`
	// Any of "base64".
	Type ManagedAgentsBase64DocumentSourceType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Data        respjson.Field
		MediaType   respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsBase64DocumentSource) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsBase64DocumentSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ManagedAgentsBase64DocumentSource to a
// ManagedAgentsBase64DocumentSourceParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ManagedAgentsBase64DocumentSourceParam.Overrides()
func (r ManagedAgentsBase64DocumentSource) ToParam() ManagedAgentsBase64DocumentSourceParam {
	return param.Override[ManagedAgentsBase64DocumentSourceParam](json.RawMessage(r.RawJSON()))
}

type ManagedAgentsBase64DocumentSourceType string

const (
	ManagedAgentsBase64DocumentSourceTypeBase64 ManagedAgentsBase64DocumentSourceType = "base64"
)

// Base64-encoded document data.
//
// The properties Data, MediaType, Type are required.
type ManagedAgentsBase64DocumentSourceParam struct {
	// Base64-encoded document data.
	Data string `json:"data" api:"required"`
	// MIME type of the document (e.g., "application/pdf").
	MediaType string `json:"media_type" api:"required"`
	// Any of "base64".
	Type ManagedAgentsBase64DocumentSourceType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsBase64DocumentSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsBase64DocumentSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsBase64DocumentSourceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Base64-encoded image data.
type ManagedAgentsBase64ImageSource struct {
	// Base64-encoded image data.
	Data string `json:"data" api:"required"`
	// MIME type of the image (e.g., "image/png", "image/jpeg", "image/gif",
	// "image/webp").
	MediaType string `json:"media_type" api:"required"`
	// Any of "base64".
	Type ManagedAgentsBase64ImageSourceType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Data        respjson.Field
		MediaType   respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsBase64ImageSource) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsBase64ImageSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ManagedAgentsBase64ImageSource to a
// ManagedAgentsBase64ImageSourceParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ManagedAgentsBase64ImageSourceParam.Overrides()
func (r ManagedAgentsBase64ImageSource) ToParam() ManagedAgentsBase64ImageSourceParam {
	return param.Override[ManagedAgentsBase64ImageSourceParam](json.RawMessage(r.RawJSON()))
}

type ManagedAgentsBase64ImageSourceType string

const (
	ManagedAgentsBase64ImageSourceTypeBase64 ManagedAgentsBase64ImageSourceType = "base64"
)

// Base64-encoded image data.
//
// The properties Data, MediaType, Type are required.
type ManagedAgentsBase64ImageSourceParam struct {
	// Base64-encoded image data.
	Data string `json:"data" api:"required"`
	// MIME type of the image (e.g., "image/png", "image/jpeg", "image/gif",
	// "image/webp").
	MediaType string `json:"media_type" api:"required"`
	// Any of "base64".
	Type ManagedAgentsBase64ImageSourceType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsBase64ImageSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsBase64ImageSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsBase64ImageSourceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The caller's organization or workspace cannot make model requests — out of
// credits or spend limit reached. Retrying with the same credentials will not
// succeed; the caller must resolve the billing state.
type ManagedAgentsBillingError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// What the client should do next in response to this error.
	RetryStatus ManagedAgentsBillingErrorRetryStatusUnion `json:"retry_status" api:"required"`
	// Any of "billing_error".
	Type ManagedAgentsBillingErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		RetryStatus respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsBillingError) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsBillingError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsBillingErrorRetryStatusUnion contains all possible properties
// and values from [ManagedAgentsRetryStatusRetrying],
// [ManagedAgentsRetryStatusExhausted], [ManagedAgentsRetryStatusTerminal].
//
// Use the [ManagedAgentsBillingErrorRetryStatusUnion.AsAny] method to switch
// on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsBillingErrorRetryStatusUnion struct {
	// Any of "retrying", "exhausted", "terminal".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyManagedAgentsBillingErrorRetryStatus is implemented by each variant of
// [ManagedAgentsBillingErrorRetryStatusUnion] to add type safety for the
// return type of [ManagedAgentsBillingErrorRetryStatusUnion.AsAny]
type anyManagedAgentsBillingErrorRetryStatus interface {
	implManagedAgentsBillingErrorRetryStatusUnion()
}

func (ManagedAgentsRetryStatusRetrying) implManagedAgentsBillingErrorRetryStatusUnion()  {}
func (ManagedAgentsRetryStatusExhausted) implManagedAgentsBillingErrorRetryStatusUnion() {}
func (ManagedAgentsRetryStatusTerminal) implManagedAgentsBillingErrorRetryStatusUnion()  {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsBillingErrorRetryStatusUnion.AsAny().(type) {
//	case qoder.ManagedAgentsRetryStatusRetrying:
//	case qoder.ManagedAgentsRetryStatusExhausted:
//	case qoder.ManagedAgentsRetryStatusTerminal:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsBillingErrorRetryStatusUnion) AsAny() anyManagedAgentsBillingErrorRetryStatus {
	switch u.Type {
	case "retrying":
		return u.AsRetrying()
	case "exhausted":
		return u.AsExhausted()
	case "terminal":
		return u.AsTerminal()
	}
	return nil
}

func (u ManagedAgentsBillingErrorRetryStatusUnion) AsRetrying() (v ManagedAgentsRetryStatusRetrying) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsBillingErrorRetryStatusUnion) AsExhausted() (v ManagedAgentsRetryStatusExhausted) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsBillingErrorRetryStatusUnion) AsTerminal() (v ManagedAgentsRetryStatusTerminal) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsBillingErrorRetryStatusUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsBillingErrorRetryStatusUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsBillingErrorType string

const (
	ManagedAgentsBillingErrorTypeBillingError ManagedAgentsBillingErrorType = "billing_error"
)

// An `environment_variable` credential's `auth.networking.allowed_hosts` includes
// a host the environment's network policy does not permit.
type ManagedAgentsCredentialHostUnreachableError struct {
	// ID of the affected credential.
	CredentialID string `json:"credential_id" api:"required"`
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// What the client should do next in response to this error.
	RetryStatus ManagedAgentsCredentialHostUnreachableErrorRetryStatusUnion `json:"retry_status" api:"required"`
	// Any of "credential_host_unreachable_error".
	Type ManagedAgentsCredentialHostUnreachableErrorType `json:"type" api:"required"`
	// ID of the vault containing the affected credential.
	VaultID string `json:"vault_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CredentialID respjson.Field
		Message      respjson.Field
		RetryStatus  respjson.Field
		Type         respjson.Field
		VaultID      respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsCredentialHostUnreachableError) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsCredentialHostUnreachableError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsCredentialHostUnreachableErrorRetryStatusUnion contains all
// possible properties and values from [ManagedAgentsRetryStatusRetrying],
// [ManagedAgentsRetryStatusExhausted], [ManagedAgentsRetryStatusTerminal].
//
// Use the [ManagedAgentsCredentialHostUnreachableErrorRetryStatusUnion.AsAny]
// method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsCredentialHostUnreachableErrorRetryStatusUnion struct {
	// Any of "retrying", "exhausted", "terminal".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyManagedAgentsCredentialHostUnreachableErrorRetryStatus is implemented by
// each variant of
// [ManagedAgentsCredentialHostUnreachableErrorRetryStatusUnion] to add type
// safety for the return type of
// [ManagedAgentsCredentialHostUnreachableErrorRetryStatusUnion.AsAny]
type anyManagedAgentsCredentialHostUnreachableErrorRetryStatus interface {
	implManagedAgentsCredentialHostUnreachableErrorRetryStatusUnion()
}

func (ManagedAgentsRetryStatusRetrying) implManagedAgentsCredentialHostUnreachableErrorRetryStatusUnion() {
}
func (ManagedAgentsRetryStatusExhausted) implManagedAgentsCredentialHostUnreachableErrorRetryStatusUnion() {
}
func (ManagedAgentsRetryStatusTerminal) implManagedAgentsCredentialHostUnreachableErrorRetryStatusUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsCredentialHostUnreachableErrorRetryStatusUnion.AsAny().(type) {
//	case qoder.ManagedAgentsRetryStatusRetrying:
//	case qoder.ManagedAgentsRetryStatusExhausted:
//	case qoder.ManagedAgentsRetryStatusTerminal:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsCredentialHostUnreachableErrorRetryStatusUnion) AsAny() anyManagedAgentsCredentialHostUnreachableErrorRetryStatus {
	switch u.Type {
	case "retrying":
		return u.AsRetrying()
	case "exhausted":
		return u.AsExhausted()
	case "terminal":
		return u.AsTerminal()
	}
	return nil
}

func (u ManagedAgentsCredentialHostUnreachableErrorRetryStatusUnion) AsRetrying() (v ManagedAgentsRetryStatusRetrying) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsCredentialHostUnreachableErrorRetryStatusUnion) AsExhausted() (v ManagedAgentsRetryStatusExhausted) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsCredentialHostUnreachableErrorRetryStatusUnion) AsTerminal() (v ManagedAgentsRetryStatusTerminal) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsCredentialHostUnreachableErrorRetryStatusUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *ManagedAgentsCredentialHostUnreachableErrorRetryStatusUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsCredentialHostUnreachableErrorType string

const (
	ManagedAgentsCredentialHostUnreachableErrorTypeCredentialHostUnreachableError ManagedAgentsCredentialHostUnreachableErrorType = "credential_host_unreachable_error"
)

// Document content, either specified directly as base64 data, as text, or as a
// reference via a URL.
type ManagedAgentsDocumentBlock struct {
	// Union type for document source variants.
	Source ManagedAgentsDocumentBlockSourceUnion `json:"source" api:"required"`
	// Any of "document".
	Type ManagedAgentsDocumentBlockType `json:"type" api:"required"`
	// Additional context about the document for the model.
	Context string `json:"context" api:"nullable"`
	// The title of the document.
	Title string `json:"title" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Source      respjson.Field
		Type        respjson.Field
		Context     respjson.Field
		Title       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsDocumentBlock) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsDocumentBlock) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ManagedAgentsDocumentBlock to a
// ManagedAgentsDocumentBlockParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ManagedAgentsDocumentBlockParam.Overrides()
func (r ManagedAgentsDocumentBlock) ToParam() ManagedAgentsDocumentBlockParam {
	return param.Override[ManagedAgentsDocumentBlockParam](json.RawMessage(r.RawJSON()))
}

// ManagedAgentsDocumentBlockSourceUnion contains all possible properties and
// values from [ManagedAgentsBase64DocumentSource],
// [ManagedAgentsPlainTextDocumentSource],
// [ManagedAgentsURLDocumentSource], [ManagedAgentsFileDocumentSource].
//
// Use the [ManagedAgentsDocumentBlockSourceUnion.AsAny] method to switch on
// the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsDocumentBlockSourceUnion struct {
	Data      string `json:"data"`
	MediaType string `json:"media_type"`
	// Any of "base64", "text", "url", "file".
	Type string `json:"type"`
	// This field is from variant [ManagedAgentsURLDocumentSource].
	URL string `json:"url"`
	// This field is from variant [ManagedAgentsFileDocumentSource].
	FileID string `json:"file_id"`
	JSON   struct {
		Data      respjson.Field
		MediaType respjson.Field
		Type      respjson.Field
		URL       respjson.Field
		FileID    respjson.Field
		raw       string
	} `json:"-"`
}

// anyManagedAgentsDocumentBlockSource is implemented by each variant of
// [ManagedAgentsDocumentBlockSourceUnion] to add type safety for the return
// type of [ManagedAgentsDocumentBlockSourceUnion.AsAny]
type anyManagedAgentsDocumentBlockSource interface {
	implManagedAgentsDocumentBlockSourceUnion()
}

func (ManagedAgentsBase64DocumentSource) implManagedAgentsDocumentBlockSourceUnion()    {}
func (ManagedAgentsPlainTextDocumentSource) implManagedAgentsDocumentBlockSourceUnion() {}
func (ManagedAgentsURLDocumentSource) implManagedAgentsDocumentBlockSourceUnion()       {}
func (ManagedAgentsFileDocumentSource) implManagedAgentsDocumentBlockSourceUnion()      {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsDocumentBlockSourceUnion.AsAny().(type) {
//	case qoder.ManagedAgentsBase64DocumentSource:
//	case qoder.ManagedAgentsPlainTextDocumentSource:
//	case qoder.ManagedAgentsURLDocumentSource:
//	case qoder.ManagedAgentsFileDocumentSource:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsDocumentBlockSourceUnion) AsAny() anyManagedAgentsDocumentBlockSource {
	switch u.Type {
	case "base64":
		return u.AsBase64()
	case "text":
		return u.AsText()
	case "url":
		return u.AsURL()
	case "file":
		return u.AsFile()
	}
	return nil
}

func (u ManagedAgentsDocumentBlockSourceUnion) AsBase64() (v ManagedAgentsBase64DocumentSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDocumentBlockSourceUnion) AsText() (v ManagedAgentsPlainTextDocumentSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDocumentBlockSourceUnion) AsURL() (v ManagedAgentsURLDocumentSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDocumentBlockSourceUnion) AsFile() (v ManagedAgentsFileDocumentSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsDocumentBlockSourceUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsDocumentBlockSourceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsDocumentBlockType string

const (
	ManagedAgentsDocumentBlockTypeDocument ManagedAgentsDocumentBlockType = "document"
)

// Document content, either specified directly as base64 data, as text, or as a
// reference via a URL.
//
// The properties Source, Type are required.
type ManagedAgentsDocumentBlockParam struct {
	// Union type for document source variants.
	Source ManagedAgentsDocumentBlockSourceUnionParam `json:"source,omitzero" api:"required"`
	// Any of "document".
	Type ManagedAgentsDocumentBlockType `json:"type,omitzero" api:"required"`
	// Additional context about the document for the model.
	Context param.Opt[string] `json:"context,omitzero"`
	// The title of the document.
	Title param.Opt[string] `json:"title,omitzero"`
	paramObj
}

func (r ManagedAgentsDocumentBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsDocumentBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsDocumentBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsDocumentBlockSourceUnionParam struct {
	OfBase64 *ManagedAgentsBase64DocumentSourceParam    `json:",omitzero,inline"`
	OfText   *ManagedAgentsPlainTextDocumentSourceParam `json:",omitzero,inline"`
	OfURL    *ManagedAgentsURLDocumentSourceParam       `json:",omitzero,inline"`
	OfFile   *ManagedAgentsFileDocumentSourceParam      `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsDocumentBlockSourceUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfBase64, u.OfText, u.OfURL, u.OfFile)
}
func (u *ManagedAgentsDocumentBlockSourceUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsDocumentBlockSourceUnionParam) asAny() any {
	if !param.IsOmitted(u.OfBase64) {
		return u.OfBase64
	} else if !param.IsOmitted(u.OfText) {
		return u.OfText
	} else if !param.IsOmitted(u.OfURL) {
		return u.OfURL
	} else if !param.IsOmitted(u.OfFile) {
		return u.OfFile
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsDocumentBlockSourceUnionParam) GetURL() *string {
	if vt := u.OfURL; vt != nil {
		return &vt.URL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsDocumentBlockSourceUnionParam) GetFileID() *string {
	if vt := u.OfFile; vt != nil {
		return &vt.FileID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsDocumentBlockSourceUnionParam) GetData() *string {
	if vt := u.OfBase64; vt != nil {
		return (*string)(&vt.Data)
	} else if vt := u.OfText; vt != nil {
		return (*string)(&vt.Data)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsDocumentBlockSourceUnionParam) GetMediaType() *string {
	if vt := u.OfBase64; vt != nil {
		return (*string)(&vt.MediaType)
	} else if vt := u.OfText; vt != nil {
		return (*string)(&vt.MediaType)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsDocumentBlockSourceUnionParam) GetType() *string {
	if vt := u.OfBase64; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfURL; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFile; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ManagedAgentsDocumentBlockSourceUnionParam](
		"type",
		apijson.Discriminator[ManagedAgentsBase64DocumentSourceParam]("base64"),
		apijson.Discriminator[ManagedAgentsPlainTextDocumentSourceParam]("text"),
		apijson.Discriminator[ManagedAgentsURLDocumentSourceParam]("url"),
		apijson.Discriminator[ManagedAgentsFileDocumentSourceParam]("file"),
	)
}

func ManagedAgentsEventParamsOfUserMessage(content []ManagedAgentsUserMessageEventParamsContentUnion) ManagedAgentsEventParamsUnion {
	var userMessage ManagedAgentsUserMessageEventParams
	userMessage.Content = content
	return ManagedAgentsEventParamsUnion{OfUserMessage: &userMessage}
}

func ManagedAgentsEventParamsOfUserInterrupt(type_ ManagedAgentsUserInterruptEventParamsType) ManagedAgentsEventParamsUnion {
	var userInterrupt ManagedAgentsUserInterruptEventParams
	userInterrupt.Type = type_
	return ManagedAgentsEventParamsUnion{OfUserInterrupt: &userInterrupt}
}

func ManagedAgentsEventParamsOfUserToolConfirmation(result ManagedAgentsUserToolConfirmationEventParamsResult, toolUseID string, type_ ManagedAgentsUserToolConfirmationEventParamsType) ManagedAgentsEventParamsUnion {
	var userToolConfirmation ManagedAgentsUserToolConfirmationEventParams
	userToolConfirmation.Result = result
	userToolConfirmation.ToolUseID = toolUseID
	userToolConfirmation.Type = type_
	return ManagedAgentsEventParamsUnion{OfUserToolConfirmation: &userToolConfirmation}
}

func ManagedAgentsEventParamsOfUserCustomToolResult(customToolUseID string) ManagedAgentsEventParamsUnion {
	var userCustomToolResult ManagedAgentsUserCustomToolResultEventParams
	userCustomToolResult.CustomToolUseID = customToolUseID
	return ManagedAgentsEventParamsUnion{OfUserCustomToolResult: &userCustomToolResult}
}

func ManagedAgentsEventParamsOfUserDefineOutcome[
	T ManagedAgentsFileRubricParams | ManagedAgentsTextRubricParams,
](description string, rubric T, type_ ManagedAgentsUserDefineOutcomeEventParamsType) ManagedAgentsEventParamsUnion {
	var userDefineOutcome ManagedAgentsUserDefineOutcomeEventParams
	userDefineOutcome.Description = description
	switch v := any(rubric).(type) {
	case ManagedAgentsFileRubricParams:
		userDefineOutcome.Rubric.OfFile = &v
	case ManagedAgentsTextRubricParams:
		userDefineOutcome.Rubric.OfText = &v
	}
	userDefineOutcome.Type = type_
	return ManagedAgentsEventParamsUnion{OfUserDefineOutcome: &userDefineOutcome}
}

func ManagedAgentsEventParamsOfUserToolResult(toolUseID string) ManagedAgentsEventParamsUnion {
	var userToolResult ManagedAgentsUserToolResultEventParams
	userToolResult.ToolUseID = toolUseID
	return ManagedAgentsEventParamsUnion{OfUserToolResult: &userToolResult}
}

func ManagedAgentsEventParamsOfSystemMessage(content []ManagedAgentsSystemContentBlockParam) ManagedAgentsEventParamsUnion {
	var systemMessage ManagedAgentsSystemMessageEventParams
	systemMessage.Content = content
	return ManagedAgentsEventParamsUnion{OfSystemMessage: &systemMessage}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsEventParamsUnion struct {
	OfUserMessage          *ManagedAgentsUserMessageEventParams          `json:",omitzero,inline"`
	OfUserInterrupt        *ManagedAgentsUserInterruptEventParams        `json:",omitzero,inline"`
	OfUserToolConfirmation *ManagedAgentsUserToolConfirmationEventParams `json:",omitzero,inline"`
	OfUserCustomToolResult *ManagedAgentsUserCustomToolResultEventParams `json:",omitzero,inline"`
	OfUserDefineOutcome    *ManagedAgentsUserDefineOutcomeEventParams    `json:",omitzero,inline"`
	OfUserToolResult       *ManagedAgentsUserToolResultEventParams       `json:",omitzero,inline"`
	OfSystemMessage        *ManagedAgentsSystemMessageEventParams        `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsEventParamsUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfUserMessage,
		u.OfUserInterrupt,
		u.OfUserToolConfirmation,
		u.OfUserCustomToolResult,
		u.OfUserDefineOutcome,
		u.OfUserToolResult,
		u.OfSystemMessage)
}
func (u *ManagedAgentsEventParamsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsEventParamsUnion) asAny() any {
	if !param.IsOmitted(u.OfUserMessage) {
		return u.OfUserMessage
	} else if !param.IsOmitted(u.OfUserInterrupt) {
		return u.OfUserInterrupt
	} else if !param.IsOmitted(u.OfUserToolConfirmation) {
		return u.OfUserToolConfirmation
	} else if !param.IsOmitted(u.OfUserCustomToolResult) {
		return u.OfUserCustomToolResult
	} else if !param.IsOmitted(u.OfUserDefineOutcome) {
		return u.OfUserDefineOutcome
	} else if !param.IsOmitted(u.OfUserToolResult) {
		return u.OfUserToolResult
	} else if !param.IsOmitted(u.OfSystemMessage) {
		return u.OfSystemMessage
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsEventParamsUnion) GetSessionThreadID() *string {
	if vt := u.OfUserInterrupt; vt != nil && vt.SessionThreadID.Valid() {
		return &vt.SessionThreadID.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsEventParamsUnion) GetResult() *string {
	if vt := u.OfUserToolConfirmation; vt != nil {
		return (*string)(&vt.Result)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsEventParamsUnion) GetDenyMessage() *string {
	if vt := u.OfUserToolConfirmation; vt != nil && vt.DenyMessage.Valid() {
		return &vt.DenyMessage.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsEventParamsUnion) GetCustomToolUseID() *string {
	if vt := u.OfUserCustomToolResult; vt != nil {
		return &vt.CustomToolUseID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsEventParamsUnion) GetDescription() *string {
	if vt := u.OfUserDefineOutcome; vt != nil {
		return &vt.Description
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsEventParamsUnion) GetRubric() *ManagedAgentsUserDefineOutcomeEventParamsRubricUnion {
	if vt := u.OfUserDefineOutcome; vt != nil {
		return &vt.Rubric
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsEventParamsUnion) GetMaxIterations() *int64 {
	if vt := u.OfUserDefineOutcome; vt != nil && vt.MaxIterations.Valid() {
		return &vt.MaxIterations.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsEventParamsUnion) GetType() *string {
	if vt := u.OfUserMessage; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfUserInterrupt; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfUserToolConfirmation; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfUserCustomToolResult; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfUserDefineOutcome; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfUserToolResult; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSystemMessage; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsEventParamsUnion) GetToolUseID() *string {
	if vt := u.OfUserToolConfirmation; vt != nil {
		return (*string)(&vt.ToolUseID)
	} else if vt := u.OfUserToolResult; vt != nil {
		return (*string)(&vt.ToolUseID)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsEventParamsUnion) GetIsError() *bool {
	if vt := u.OfUserCustomToolResult; vt != nil && vt.IsError.Valid() {
		return &vt.IsError.Value
	} else if vt := u.OfUserToolResult; vt != nil && vt.IsError.Valid() {
		return &vt.IsError.Value
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u ManagedAgentsEventParamsUnion) GetContent() (res managedAgentsEventParamsUnionContent) {
	if vt := u.OfUserMessage; vt != nil {
		res.any = &vt.Content
	} else if vt := u.OfUserCustomToolResult; vt != nil {
		res.any = &vt.Content
	} else if vt := u.OfUserToolResult; vt != nil {
		res.any = &vt.Content
	} else if vt := u.OfSystemMessage; vt != nil {
		res.any = &vt.Content
	}
	return
}

// Can have the runtime types
// [_[]ManagedAgentsUserMessageEventParamsContentUnion],
// [_[]ManagedAgentsUserCustomToolResultEventParamsContentUnion],
// [_[]ManagedAgentsUserToolResultEventParamsContentUnion],
// [_[]ManagedAgentsSystemContentBlockParam]
type managedAgentsEventParamsUnionContent struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *[]qoder.ManagedAgentsUserMessageEventParamsContentUnion:
//	case *[]qoder.ManagedAgentsUserCustomToolResultEventParamsContentUnion:
//	case *[]qoder.ManagedAgentsUserToolResultEventParamsContentUnion:
//	case *[]qoder.ManagedAgentsSystemContentBlockParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u managedAgentsEventParamsUnionContent) AsAny() any { return u.any }

func init() {
	apijson.RegisterUnion[ManagedAgentsEventParamsUnion](
		"type",
		apijson.Discriminator[ManagedAgentsUserMessageEventParams]("user.message"),
		apijson.Discriminator[ManagedAgentsUserInterruptEventParams]("user.interrupt"),
		apijson.Discriminator[ManagedAgentsUserToolConfirmationEventParams]("user.tool_confirmation"),
		apijson.Discriminator[ManagedAgentsUserCustomToolResultEventParams]("user.custom_tool_result"),
		apijson.Discriminator[ManagedAgentsUserDefineOutcomeEventParams]("user.define_outcome"),
		apijson.Discriminator[ManagedAgentsUserToolResultEventParams]("user.tool_result"),
		apijson.Discriminator[ManagedAgentsSystemMessageEventParams]("system.message"),
	)
}

// Document referenced by file ID.
type ManagedAgentsFileDocumentSource struct {
	// ID of a previously uploaded file.
	FileID string `json:"file_id" api:"required"`
	// Any of "file".
	Type ManagedAgentsFileDocumentSourceType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		FileID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsFileDocumentSource) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsFileDocumentSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ManagedAgentsFileDocumentSource to a
// ManagedAgentsFileDocumentSourceParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ManagedAgentsFileDocumentSourceParam.Overrides()
func (r ManagedAgentsFileDocumentSource) ToParam() ManagedAgentsFileDocumentSourceParam {
	return param.Override[ManagedAgentsFileDocumentSourceParam](json.RawMessage(r.RawJSON()))
}

type ManagedAgentsFileDocumentSourceType string

const (
	ManagedAgentsFileDocumentSourceTypeFile ManagedAgentsFileDocumentSourceType = "file"
)

// Document referenced by file ID.
//
// The properties FileID, Type are required.
type ManagedAgentsFileDocumentSourceParam struct {
	// ID of a previously uploaded file.
	FileID string `json:"file_id" api:"required"`
	// Any of "file".
	Type ManagedAgentsFileDocumentSourceType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsFileDocumentSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsFileDocumentSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsFileDocumentSourceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Image referenced by file ID.
type ManagedAgentsFileImageSource struct {
	// ID of a previously uploaded file.
	FileID string `json:"file_id" api:"required"`
	// Any of "file".
	Type ManagedAgentsFileImageSourceType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		FileID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsFileImageSource) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsFileImageSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ManagedAgentsFileImageSource to a
// ManagedAgentsFileImageSourceParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ManagedAgentsFileImageSourceParam.Overrides()
func (r ManagedAgentsFileImageSource) ToParam() ManagedAgentsFileImageSourceParam {
	return param.Override[ManagedAgentsFileImageSourceParam](json.RawMessage(r.RawJSON()))
}

type ManagedAgentsFileImageSourceType string

const (
	ManagedAgentsFileImageSourceTypeFile ManagedAgentsFileImageSourceType = "file"
)

// Image referenced by file ID.
//
// The properties FileID, Type are required.
type ManagedAgentsFileImageSourceParam struct {
	// ID of a previously uploaded file.
	FileID string `json:"file_id" api:"required"`
	// Any of "file".
	Type ManagedAgentsFileImageSourceType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsFileImageSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsFileImageSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsFileImageSourceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Rubric referenced by a file uploaded via the Files API.
type ManagedAgentsFileRubric struct {
	// ID of the rubric file.
	FileID string `json:"file_id" api:"required"`
	// Any of "file".
	Type ManagedAgentsFileRubricType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		FileID      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsFileRubric) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsFileRubric) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsFileRubricType string

const (
	ManagedAgentsFileRubricTypeFile ManagedAgentsFileRubricType = "file"
)

// Rubric referenced by a file uploaded via the Files API.
//
// The properties FileID, Type are required.
type ManagedAgentsFileRubricParams struct {
	// ID of the rubric file.
	FileID string `json:"file_id" api:"required"`
	// Any of "file".
	Type ManagedAgentsFileRubricParamsType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsFileRubricParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsFileRubricParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsFileRubricParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsFileRubricParamsType string

const (
	ManagedAgentsFileRubricParamsTypeFile ManagedAgentsFileRubricParamsType = "file"
)

// Image content specified directly as base64 data or as a reference via a URL.
type ManagedAgentsImageBlock struct {
	// Union type for image source variants.
	Source ManagedAgentsImageBlockSourceUnion `json:"source" api:"required"`
	// Any of "image".
	Type ManagedAgentsImageBlockType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Source      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsImageBlock) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsImageBlock) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ManagedAgentsImageBlock to a
// ManagedAgentsImageBlockParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ManagedAgentsImageBlockParam.Overrides()
func (r ManagedAgentsImageBlock) ToParam() ManagedAgentsImageBlockParam {
	return param.Override[ManagedAgentsImageBlockParam](json.RawMessage(r.RawJSON()))
}

// ManagedAgentsImageBlockSourceUnion contains all possible properties and
// values from [ManagedAgentsBase64ImageSource],
// [ManagedAgentsURLImageSource], [ManagedAgentsFileImageSource].
//
// Use the [ManagedAgentsImageBlockSourceUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsImageBlockSourceUnion struct {
	// This field is from variant [ManagedAgentsBase64ImageSource].
	Data string `json:"data"`
	// This field is from variant [ManagedAgentsBase64ImageSource].
	MediaType string `json:"media_type"`
	// Any of "base64", "url", "file".
	Type string `json:"type"`
	// This field is from variant [ManagedAgentsURLImageSource].
	URL string `json:"url"`
	// This field is from variant [ManagedAgentsFileImageSource].
	FileID string `json:"file_id"`
	JSON   struct {
		Data      respjson.Field
		MediaType respjson.Field
		Type      respjson.Field
		URL       respjson.Field
		FileID    respjson.Field
		raw       string
	} `json:"-"`
}

// anyManagedAgentsImageBlockSource is implemented by each variant of
// [ManagedAgentsImageBlockSourceUnion] to add type safety for the return type
// of [ManagedAgentsImageBlockSourceUnion.AsAny]
type anyManagedAgentsImageBlockSource interface {
	implManagedAgentsImageBlockSourceUnion()
}

func (ManagedAgentsBase64ImageSource) implManagedAgentsImageBlockSourceUnion() {}
func (ManagedAgentsURLImageSource) implManagedAgentsImageBlockSourceUnion()    {}
func (ManagedAgentsFileImageSource) implManagedAgentsImageBlockSourceUnion()   {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsImageBlockSourceUnion.AsAny().(type) {
//	case qoder.ManagedAgentsBase64ImageSource:
//	case qoder.ManagedAgentsURLImageSource:
//	case qoder.ManagedAgentsFileImageSource:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsImageBlockSourceUnion) AsAny() anyManagedAgentsImageBlockSource {
	switch u.Type {
	case "base64":
		return u.AsBase64()
	case "url":
		return u.AsURL()
	case "file":
		return u.AsFile()
	}
	return nil
}

func (u ManagedAgentsImageBlockSourceUnion) AsBase64() (v ManagedAgentsBase64ImageSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsImageBlockSourceUnion) AsURL() (v ManagedAgentsURLImageSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsImageBlockSourceUnion) AsFile() (v ManagedAgentsFileImageSource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsImageBlockSourceUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsImageBlockSourceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsImageBlockType string

const (
	ManagedAgentsImageBlockTypeImage ManagedAgentsImageBlockType = "image"
)

// Image content specified directly as base64 data or as a reference via a URL.
//
// The properties Source, Type are required.
type ManagedAgentsImageBlockParam struct {
	// Union type for image source variants.
	Source ManagedAgentsImageBlockSourceUnionParam `json:"source,omitzero" api:"required"`
	// Any of "image".
	Type ManagedAgentsImageBlockType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsImageBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsImageBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsImageBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsImageBlockSourceUnionParam struct {
	OfBase64 *ManagedAgentsBase64ImageSourceParam `json:",omitzero,inline"`
	OfURL    *ManagedAgentsURLImageSourceParam    `json:",omitzero,inline"`
	OfFile   *ManagedAgentsFileImageSourceParam   `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsImageBlockSourceUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfBase64, u.OfURL, u.OfFile)
}
func (u *ManagedAgentsImageBlockSourceUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsImageBlockSourceUnionParam) asAny() any {
	if !param.IsOmitted(u.OfBase64) {
		return u.OfBase64
	} else if !param.IsOmitted(u.OfURL) {
		return u.OfURL
	} else if !param.IsOmitted(u.OfFile) {
		return u.OfFile
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsImageBlockSourceUnionParam) GetData() *string {
	if vt := u.OfBase64; vt != nil {
		return &vt.Data
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsImageBlockSourceUnionParam) GetMediaType() *string {
	if vt := u.OfBase64; vt != nil {
		return &vt.MediaType
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsImageBlockSourceUnionParam) GetURL() *string {
	if vt := u.OfURL; vt != nil {
		return &vt.URL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsImageBlockSourceUnionParam) GetFileID() *string {
	if vt := u.OfFile; vt != nil {
		return &vt.FileID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsImageBlockSourceUnionParam) GetType() *string {
	if vt := u.OfBase64; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfURL; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFile; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ManagedAgentsImageBlockSourceUnionParam](
		"type",
		apijson.Discriminator[ManagedAgentsBase64ImageSourceParam]("base64"),
		apijson.Discriminator[ManagedAgentsURLImageSourceParam]("url"),
		apijson.Discriminator[ManagedAgentsFileImageSourceParam]("file"),
	)
}

// Authentication to an MCP server failed.
type ManagedAgentsMCPAuthenticationFailedError struct {
	// Name of the MCP server that failed authentication.
	MCPServerName string `json:"mcp_server_name" api:"required"`
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// What the client should do next in response to this error.
	RetryStatus ManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion `json:"retry_status" api:"required"`
	// Any of "mcp_authentication_failed_error".
	Type ManagedAgentsMCPAuthenticationFailedErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		MCPServerName respjson.Field
		Message       respjson.Field
		RetryStatus   respjson.Field
		Type          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsMCPAuthenticationFailedError) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsMCPAuthenticationFailedError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion contains all
// possible properties and values from [ManagedAgentsRetryStatusRetrying],
// [ManagedAgentsRetryStatusExhausted], [ManagedAgentsRetryStatusTerminal].
//
// Use the [ManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion.AsAny]
// method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion struct {
	// Any of "retrying", "exhausted", "terminal".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyManagedAgentsMCPAuthenticationFailedErrorRetryStatus is implemented by
// each variant of [ManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion]
// to add type safety for the return type of
// [ManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion.AsAny]
type anyManagedAgentsMCPAuthenticationFailedErrorRetryStatus interface {
	implManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion()
}

func (ManagedAgentsRetryStatusRetrying) implManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion() {
}
func (ManagedAgentsRetryStatusExhausted) implManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion() {
}
func (ManagedAgentsRetryStatusTerminal) implManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion.AsAny().(type) {
//	case qoder.ManagedAgentsRetryStatusRetrying:
//	case qoder.ManagedAgentsRetryStatusExhausted:
//	case qoder.ManagedAgentsRetryStatusTerminal:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion) AsAny() anyManagedAgentsMCPAuthenticationFailedErrorRetryStatus {
	switch u.Type {
	case "retrying":
		return u.AsRetrying()
	case "exhausted":
		return u.AsExhausted()
	case "terminal":
		return u.AsTerminal()
	}
	return nil
}

func (u ManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion) AsRetrying() (v ManagedAgentsRetryStatusRetrying) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion) AsExhausted() (v ManagedAgentsRetryStatusExhausted) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion) AsTerminal() (v ManagedAgentsRetryStatusTerminal) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *ManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsMCPAuthenticationFailedErrorType string

const (
	ManagedAgentsMCPAuthenticationFailedErrorTypeMCPAuthenticationFailedError ManagedAgentsMCPAuthenticationFailedErrorType = "mcp_authentication_failed_error"
)

// Failed to connect to an MCP server.
type ManagedAgentsMCPConnectionFailedError struct {
	// Name of the MCP server that failed to connect.
	MCPServerName string `json:"mcp_server_name" api:"required"`
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// What the client should do next in response to this error.
	RetryStatus ManagedAgentsMCPConnectionFailedErrorRetryStatusUnion `json:"retry_status" api:"required"`
	// Any of "mcp_connection_failed_error".
	Type ManagedAgentsMCPConnectionFailedErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		MCPServerName respjson.Field
		Message       respjson.Field
		RetryStatus   respjson.Field
		Type          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsMCPConnectionFailedError) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsMCPConnectionFailedError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsMCPConnectionFailedErrorRetryStatusUnion contains all possible
// properties and values from [ManagedAgentsRetryStatusRetrying],
// [ManagedAgentsRetryStatusExhausted], [ManagedAgentsRetryStatusTerminal].
//
// Use the [ManagedAgentsMCPConnectionFailedErrorRetryStatusUnion.AsAny] method
// to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsMCPConnectionFailedErrorRetryStatusUnion struct {
	// Any of "retrying", "exhausted", "terminal".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyManagedAgentsMCPConnectionFailedErrorRetryStatus is implemented by each
// variant of [ManagedAgentsMCPConnectionFailedErrorRetryStatusUnion] to add
// type safety for the return type of
// [ManagedAgentsMCPConnectionFailedErrorRetryStatusUnion.AsAny]
type anyManagedAgentsMCPConnectionFailedErrorRetryStatus interface {
	implManagedAgentsMCPConnectionFailedErrorRetryStatusUnion()
}

func (ManagedAgentsRetryStatusRetrying) implManagedAgentsMCPConnectionFailedErrorRetryStatusUnion() {
}
func (ManagedAgentsRetryStatusExhausted) implManagedAgentsMCPConnectionFailedErrorRetryStatusUnion() {
}
func (ManagedAgentsRetryStatusTerminal) implManagedAgentsMCPConnectionFailedErrorRetryStatusUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsMCPConnectionFailedErrorRetryStatusUnion.AsAny().(type) {
//	case qoder.ManagedAgentsRetryStatusRetrying:
//	case qoder.ManagedAgentsRetryStatusExhausted:
//	case qoder.ManagedAgentsRetryStatusTerminal:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsMCPConnectionFailedErrorRetryStatusUnion) AsAny() anyManagedAgentsMCPConnectionFailedErrorRetryStatus {
	switch u.Type {
	case "retrying":
		return u.AsRetrying()
	case "exhausted":
		return u.AsExhausted()
	case "terminal":
		return u.AsTerminal()
	}
	return nil
}

func (u ManagedAgentsMCPConnectionFailedErrorRetryStatusUnion) AsRetrying() (v ManagedAgentsRetryStatusRetrying) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsMCPConnectionFailedErrorRetryStatusUnion) AsExhausted() (v ManagedAgentsRetryStatusExhausted) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsMCPConnectionFailedErrorRetryStatusUnion) AsTerminal() (v ManagedAgentsRetryStatusTerminal) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsMCPConnectionFailedErrorRetryStatusUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *ManagedAgentsMCPConnectionFailedErrorRetryStatusUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsMCPConnectionFailedErrorType string

const (
	ManagedAgentsMCPConnectionFailedErrorTypeMCPConnectionFailedError ManagedAgentsMCPConnectionFailedErrorType = "mcp_connection_failed_error"
)

// The model is currently overloaded. Emitted after automatic retries are
// exhausted.
type ManagedAgentsModelOverloadedError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// What the client should do next in response to this error.
	RetryStatus ManagedAgentsModelOverloadedErrorRetryStatusUnion `json:"retry_status" api:"required"`
	// Any of "model_overloaded_error".
	Type ManagedAgentsModelOverloadedErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		RetryStatus respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsModelOverloadedError) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsModelOverloadedError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsModelOverloadedErrorRetryStatusUnion contains all possible
// properties and values from [ManagedAgentsRetryStatusRetrying],
// [ManagedAgentsRetryStatusExhausted], [ManagedAgentsRetryStatusTerminal].
//
// Use the [ManagedAgentsModelOverloadedErrorRetryStatusUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsModelOverloadedErrorRetryStatusUnion struct {
	// Any of "retrying", "exhausted", "terminal".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyManagedAgentsModelOverloadedErrorRetryStatus is implemented by each
// variant of [ManagedAgentsModelOverloadedErrorRetryStatusUnion] to add type
// safety for the return type of
// [ManagedAgentsModelOverloadedErrorRetryStatusUnion.AsAny]
type anyManagedAgentsModelOverloadedErrorRetryStatus interface {
	implManagedAgentsModelOverloadedErrorRetryStatusUnion()
}

func (ManagedAgentsRetryStatusRetrying) implManagedAgentsModelOverloadedErrorRetryStatusUnion() {
}
func (ManagedAgentsRetryStatusExhausted) implManagedAgentsModelOverloadedErrorRetryStatusUnion() {
}
func (ManagedAgentsRetryStatusTerminal) implManagedAgentsModelOverloadedErrorRetryStatusUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsModelOverloadedErrorRetryStatusUnion.AsAny().(type) {
//	case qoder.ManagedAgentsRetryStatusRetrying:
//	case qoder.ManagedAgentsRetryStatusExhausted:
//	case qoder.ManagedAgentsRetryStatusTerminal:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsModelOverloadedErrorRetryStatusUnion) AsAny() anyManagedAgentsModelOverloadedErrorRetryStatus {
	switch u.Type {
	case "retrying":
		return u.AsRetrying()
	case "exhausted":
		return u.AsExhausted()
	case "terminal":
		return u.AsTerminal()
	}
	return nil
}

func (u ManagedAgentsModelOverloadedErrorRetryStatusUnion) AsRetrying() (v ManagedAgentsRetryStatusRetrying) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsModelOverloadedErrorRetryStatusUnion) AsExhausted() (v ManagedAgentsRetryStatusExhausted) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsModelOverloadedErrorRetryStatusUnion) AsTerminal() (v ManagedAgentsRetryStatusTerminal) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsModelOverloadedErrorRetryStatusUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsModelOverloadedErrorRetryStatusUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsModelOverloadedErrorType string

const (
	ManagedAgentsModelOverloadedErrorTypeModelOverloadedError ManagedAgentsModelOverloadedErrorType = "model_overloaded_error"
)

// The model request was rate-limited.
type ManagedAgentsModelRateLimitedError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// What the client should do next in response to this error.
	RetryStatus ManagedAgentsModelRateLimitedErrorRetryStatusUnion `json:"retry_status" api:"required"`
	// Any of "model_rate_limited_error".
	Type ManagedAgentsModelRateLimitedErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		RetryStatus respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsModelRateLimitedError) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsModelRateLimitedError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsModelRateLimitedErrorRetryStatusUnion contains all possible
// properties and values from [ManagedAgentsRetryStatusRetrying],
// [ManagedAgentsRetryStatusExhausted], [ManagedAgentsRetryStatusTerminal].
//
// Use the [ManagedAgentsModelRateLimitedErrorRetryStatusUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsModelRateLimitedErrorRetryStatusUnion struct {
	// Any of "retrying", "exhausted", "terminal".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyManagedAgentsModelRateLimitedErrorRetryStatus is implemented by each
// variant of [ManagedAgentsModelRateLimitedErrorRetryStatusUnion] to add type
// safety for the return type of
// [ManagedAgentsModelRateLimitedErrorRetryStatusUnion.AsAny]
type anyManagedAgentsModelRateLimitedErrorRetryStatus interface {
	implManagedAgentsModelRateLimitedErrorRetryStatusUnion()
}

func (ManagedAgentsRetryStatusRetrying) implManagedAgentsModelRateLimitedErrorRetryStatusUnion() {
}
func (ManagedAgentsRetryStatusExhausted) implManagedAgentsModelRateLimitedErrorRetryStatusUnion() {
}
func (ManagedAgentsRetryStatusTerminal) implManagedAgentsModelRateLimitedErrorRetryStatusUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsModelRateLimitedErrorRetryStatusUnion.AsAny().(type) {
//	case qoder.ManagedAgentsRetryStatusRetrying:
//	case qoder.ManagedAgentsRetryStatusExhausted:
//	case qoder.ManagedAgentsRetryStatusTerminal:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsModelRateLimitedErrorRetryStatusUnion) AsAny() anyManagedAgentsModelRateLimitedErrorRetryStatus {
	switch u.Type {
	case "retrying":
		return u.AsRetrying()
	case "exhausted":
		return u.AsExhausted()
	case "terminal":
		return u.AsTerminal()
	}
	return nil
}

func (u ManagedAgentsModelRateLimitedErrorRetryStatusUnion) AsRetrying() (v ManagedAgentsRetryStatusRetrying) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsModelRateLimitedErrorRetryStatusUnion) AsExhausted() (v ManagedAgentsRetryStatusExhausted) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsModelRateLimitedErrorRetryStatusUnion) AsTerminal() (v ManagedAgentsRetryStatusTerminal) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsModelRateLimitedErrorRetryStatusUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsModelRateLimitedErrorRetryStatusUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsModelRateLimitedErrorType string

const (
	ManagedAgentsModelRateLimitedErrorTypeModelRateLimitedError ManagedAgentsModelRateLimitedErrorType = "model_rate_limited_error"
)

// A model request failed for a reason other than overload or rate-limiting.
type ManagedAgentsModelRequestFailedError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// What the client should do next in response to this error.
	RetryStatus ManagedAgentsModelRequestFailedErrorRetryStatusUnion `json:"retry_status" api:"required"`
	// Any of "model_request_failed_error".
	Type ManagedAgentsModelRequestFailedErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		RetryStatus respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsModelRequestFailedError) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsModelRequestFailedError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsModelRequestFailedErrorRetryStatusUnion contains all possible
// properties and values from [ManagedAgentsRetryStatusRetrying],
// [ManagedAgentsRetryStatusExhausted], [ManagedAgentsRetryStatusTerminal].
//
// Use the [ManagedAgentsModelRequestFailedErrorRetryStatusUnion.AsAny] method
// to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsModelRequestFailedErrorRetryStatusUnion struct {
	// Any of "retrying", "exhausted", "terminal".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyManagedAgentsModelRequestFailedErrorRetryStatus is implemented by each
// variant of [ManagedAgentsModelRequestFailedErrorRetryStatusUnion] to add
// type safety for the return type of
// [ManagedAgentsModelRequestFailedErrorRetryStatusUnion.AsAny]
type anyManagedAgentsModelRequestFailedErrorRetryStatus interface {
	implManagedAgentsModelRequestFailedErrorRetryStatusUnion()
}

func (ManagedAgentsRetryStatusRetrying) implManagedAgentsModelRequestFailedErrorRetryStatusUnion() {
}
func (ManagedAgentsRetryStatusExhausted) implManagedAgentsModelRequestFailedErrorRetryStatusUnion() {
}
func (ManagedAgentsRetryStatusTerminal) implManagedAgentsModelRequestFailedErrorRetryStatusUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsModelRequestFailedErrorRetryStatusUnion.AsAny().(type) {
//	case qoder.ManagedAgentsRetryStatusRetrying:
//	case qoder.ManagedAgentsRetryStatusExhausted:
//	case qoder.ManagedAgentsRetryStatusTerminal:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsModelRequestFailedErrorRetryStatusUnion) AsAny() anyManagedAgentsModelRequestFailedErrorRetryStatus {
	switch u.Type {
	case "retrying":
		return u.AsRetrying()
	case "exhausted":
		return u.AsExhausted()
	case "terminal":
		return u.AsTerminal()
	}
	return nil
}

func (u ManagedAgentsModelRequestFailedErrorRetryStatusUnion) AsRetrying() (v ManagedAgentsRetryStatusRetrying) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsModelRequestFailedErrorRetryStatusUnion) AsExhausted() (v ManagedAgentsRetryStatusExhausted) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsModelRequestFailedErrorRetryStatusUnion) AsTerminal() (v ManagedAgentsRetryStatusTerminal) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsModelRequestFailedErrorRetryStatusUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsModelRequestFailedErrorRetryStatusUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsModelRequestFailedErrorType string

const (
	ManagedAgentsModelRequestFailedErrorTypeModelRequestFailedError ManagedAgentsModelRequestFailedErrorType = "model_request_failed_error"
)

// Plain text document content.
type ManagedAgentsPlainTextDocumentSource struct {
	// The plain text content.
	Data string `json:"data" api:"required"`
	// MIME type of the text content. Must be "text/plain".
	//
	// Any of "text/plain".
	MediaType ManagedAgentsPlainTextDocumentSourceMediaType `json:"media_type" api:"required"`
	// Any of "text".
	Type ManagedAgentsPlainTextDocumentSourceType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Data        respjson.Field
		MediaType   respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsPlainTextDocumentSource) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsPlainTextDocumentSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ManagedAgentsPlainTextDocumentSource to a
// ManagedAgentsPlainTextDocumentSourceParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ManagedAgentsPlainTextDocumentSourceParam.Overrides()
func (r ManagedAgentsPlainTextDocumentSource) ToParam() ManagedAgentsPlainTextDocumentSourceParam {
	return param.Override[ManagedAgentsPlainTextDocumentSourceParam](json.RawMessage(r.RawJSON()))
}

// MIME type of the text content. Must be "text/plain".
type ManagedAgentsPlainTextDocumentSourceMediaType string

const (
	ManagedAgentsPlainTextDocumentSourceMediaTypeTextPlain ManagedAgentsPlainTextDocumentSourceMediaType = "text/plain"
)

type ManagedAgentsPlainTextDocumentSourceType string

const (
	ManagedAgentsPlainTextDocumentSourceTypeText ManagedAgentsPlainTextDocumentSourceType = "text"
)

// Plain text document content.
//
// The properties Data, MediaType, Type are required.
type ManagedAgentsPlainTextDocumentSourceParam struct {
	// The plain text content.
	Data string `json:"data" api:"required"`
	// MIME type of the text content. Must be "text/plain".
	//
	// Any of "text/plain".
	MediaType ManagedAgentsPlainTextDocumentSourceMediaType `json:"media_type,omitzero" api:"required"`
	// Any of "text".
	Type ManagedAgentsPlainTextDocumentSourceType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsPlainTextDocumentSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsPlainTextDocumentSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsPlainTextDocumentSourceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Placeholder for content withheld by Qoder model policy.
type ManagedAgentsRedactedBlock struct {
	// Any of "redacted".
	Type ManagedAgentsRedactedBlockType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsRedactedBlock) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsRedactedBlock) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ManagedAgentsRedactedBlock to a
// ManagedAgentsRedactedBlockParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ManagedAgentsRedactedBlockParam.Overrides()
func (r ManagedAgentsRedactedBlock) ToParam() ManagedAgentsRedactedBlockParam {
	return param.Override[ManagedAgentsRedactedBlockParam](json.RawMessage(r.RawJSON()))
}

type ManagedAgentsRedactedBlockType string

const (
	ManagedAgentsRedactedBlockTypeRedacted ManagedAgentsRedactedBlockType = "redacted"
)

// Placeholder for content withheld by Qoder model policy.
//
// The property Type is required.
type ManagedAgentsRedactedBlockParam struct {
	// Any of "redacted".
	Type ManagedAgentsRedactedBlockType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsRedactedBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsRedactedBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsRedactedBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// This turn is dead; queued inputs are flushed and the session returns to idle.
// Client may send a new prompt.
type ManagedAgentsRetryStatusExhausted struct {
	// Any of "exhausted".
	Type ManagedAgentsRetryStatusExhaustedType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsRetryStatusExhausted) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsRetryStatusExhausted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsRetryStatusExhaustedType string

const (
	ManagedAgentsRetryStatusExhaustedTypeExhausted ManagedAgentsRetryStatusExhaustedType = "exhausted"
)

// The server is retrying automatically. Client should wait; the same error type
// may fire again as retrying, then once as exhausted when the retry budget runs
// out.
type ManagedAgentsRetryStatusRetrying struct {
	// Any of "retrying".
	Type ManagedAgentsRetryStatusRetryingType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsRetryStatusRetrying) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsRetryStatusRetrying) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsRetryStatusRetryingType string

const (
	ManagedAgentsRetryStatusRetryingTypeRetrying ManagedAgentsRetryStatusRetryingType = "retrying"
)

// The session encountered a terminal error and will transition to `terminated`
// state.
type ManagedAgentsRetryStatusTerminal struct {
	// Any of "terminal".
	Type ManagedAgentsRetryStatusTerminalType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsRetryStatusTerminal) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsRetryStatusTerminal) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsRetryStatusTerminalType string

const (
	ManagedAgentsRetryStatusTerminalTypeTerminal ManagedAgentsRetryStatusTerminalType = "terminal"
)

// A block containing a web search result.
type ManagedAgentsSearchResultBlock struct {
	// Citation settings for a search result.
	Citations ManagedAgentsSearchResultCitations `json:"citations" api:"required"`
	// Array of text content blocks from the search result.
	Content []ManagedAgentsSearchResultContent `json:"content" api:"required"`
	// The URL source of the search result.
	Source string `json:"source" api:"required"`
	// The title of the search result.
	Title string `json:"title" api:"required"`
	// Any of "search_result".
	Type ManagedAgentsSearchResultBlockType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Citations   respjson.Field
		Content     respjson.Field
		Source      respjson.Field
		Title       respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSearchResultBlock) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSearchResultBlock) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ManagedAgentsSearchResultBlock to a
// ManagedAgentsSearchResultBlockParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ManagedAgentsSearchResultBlockParam.Overrides()
func (r ManagedAgentsSearchResultBlock) ToParam() ManagedAgentsSearchResultBlockParam {
	return param.Override[ManagedAgentsSearchResultBlockParam](json.RawMessage(r.RawJSON()))
}

type ManagedAgentsSearchResultBlockType string

const (
	ManagedAgentsSearchResultBlockTypeSearchResult ManagedAgentsSearchResultBlockType = "search_result"
)

// A block containing a web search result.
//
// The properties Citations, Content, Source, Title, Type are required.
type ManagedAgentsSearchResultBlockParam struct {
	// Citation settings for a search result.
	Citations ManagedAgentsSearchResultCitationsParam `json:"citations,omitzero" api:"required"`
	// Array of text content blocks from the search result.
	Content []ManagedAgentsSearchResultContentParam `json:"content,omitzero" api:"required"`
	// The URL source of the search result.
	Source string `json:"source" api:"required"`
	// The title of the search result.
	Title string `json:"title" api:"required"`
	// Any of "search_result".
	Type ManagedAgentsSearchResultBlockType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsSearchResultBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsSearchResultBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsSearchResultBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Citation settings for a search result.
type ManagedAgentsSearchResultCitations struct {
	// Whether citations are enabled for this search result.
	Enabled bool `json:"enabled" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSearchResultCitations) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSearchResultCitations) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ManagedAgentsSearchResultCitations to a
// ManagedAgentsSearchResultCitationsParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ManagedAgentsSearchResultCitationsParam.Overrides()
func (r ManagedAgentsSearchResultCitations) ToParam() ManagedAgentsSearchResultCitationsParam {
	return param.Override[ManagedAgentsSearchResultCitationsParam](json.RawMessage(r.RawJSON()))
}

// Citation settings for a search result.
//
// The property Enabled is required.
type ManagedAgentsSearchResultCitationsParam struct {
	// Whether citations are enabled for this search result.
	Enabled bool `json:"enabled" api:"required"`
	paramObj
}

func (r ManagedAgentsSearchResultCitationsParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsSearchResultCitationsParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsSearchResultCitationsParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Text content within a search result.
type ManagedAgentsSearchResultContent struct {
	// The text content.
	Text string `json:"text" api:"required"`
	// Any of "text".
	Type ManagedAgentsSearchResultContentType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Text        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSearchResultContent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSearchResultContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ManagedAgentsSearchResultContent to a
// ManagedAgentsSearchResultContentParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ManagedAgentsSearchResultContentParam.Overrides()
func (r ManagedAgentsSearchResultContent) ToParam() ManagedAgentsSearchResultContentParam {
	return param.Override[ManagedAgentsSearchResultContentParam](json.RawMessage(r.RawJSON()))
}

type ManagedAgentsSearchResultContentType string

const (
	ManagedAgentsSearchResultContentTypeText ManagedAgentsSearchResultContentType = "text"
)

// Text content within a search result.
//
// The properties Text, Type are required.
type ManagedAgentsSearchResultContentParam struct {
	// The text content.
	Text string `json:"text" api:"required"`
	// Any of "text".
	Type ManagedAgentsSearchResultContentType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsSearchResultContentParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsSearchResultContentParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsSearchResultContentParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Events that were successfully sent to the session.
type ManagedAgentsSendSessionEvents struct {
	// Sent events
	Data []ManagedAgentsSendSessionEventsDataUnion `json:"data"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Data        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSendSessionEvents) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSendSessionEvents) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsSendSessionEventsDataUnion contains all possible properties and
// values from [ManagedAgentsUserMessageEvent],
// [ManagedAgentsUserInterruptEvent],
// [ManagedAgentsUserToolConfirmationEvent],
// [ManagedAgentsUserCustomToolResultEvent],
// [ManagedAgentsUserDefineOutcomeEvent],
// [ManagedAgentsUserToolResultEvent], [ManagedAgentsSystemMessageEvent].
//
// Use the [ManagedAgentsSendSessionEventsDataUnion.AsAny] method to switch on
// the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsSendSessionEventsDataUnion struct {
	ID string `json:"id"`
	// This field is a union of [[]ManagedAgentsUserMessageEventContentUnion],
	// [[]ManagedAgentsUserCustomToolResultEventContentUnion],
	// [[]ManagedAgentsUserToolResultEventContentUnion],
	// [[]ManagedAgentsSystemContentBlock]
	Content ManagedAgentsSendSessionEventsDataUnionContent `json:"content"`
	// Any of "user.message", "user.interrupt", "user.tool_confirmation",
	// "user.custom_tool_result", "user.define_outcome", "user.tool_result",
	// "system.message".
	Type            string    `json:"type"`
	ProcessedAt     time.Time `json:"processed_at"`
	SessionThreadID string    `json:"session_thread_id"`
	// This field is from variant [ManagedAgentsUserToolConfirmationEvent].
	Result    ManagedAgentsUserToolConfirmationEventResult `json:"result"`
	ToolUseID string                                       `json:"tool_use_id"`
	// This field is from variant [ManagedAgentsUserToolConfirmationEvent].
	DenyMessage string `json:"deny_message"`
	// This field is from variant [ManagedAgentsUserCustomToolResultEvent].
	CustomToolUseID string `json:"custom_tool_use_id"`
	IsError         bool   `json:"is_error"`
	// This field is from variant [ManagedAgentsUserDefineOutcomeEvent].
	Description string `json:"description"`
	// This field is from variant [ManagedAgentsUserDefineOutcomeEvent].
	MaxIterations int64 `json:"max_iterations"`
	// This field is from variant [ManagedAgentsUserDefineOutcomeEvent].
	OutcomeID string `json:"outcome_id"`
	// This field is from variant [ManagedAgentsUserDefineOutcomeEvent].
	Rubric ManagedAgentsUserDefineOutcomeEventRubricUnion `json:"rubric"`
	JSON   struct {
		ID              respjson.Field
		Content         respjson.Field
		Type            respjson.Field
		ProcessedAt     respjson.Field
		SessionThreadID respjson.Field
		Result          respjson.Field
		ToolUseID       respjson.Field
		DenyMessage     respjson.Field
		CustomToolUseID respjson.Field
		IsError         respjson.Field
		Description     respjson.Field
		MaxIterations   respjson.Field
		OutcomeID       respjson.Field
		Rubric          respjson.Field
		raw             string
	} `json:"-"`
}

// anyManagedAgentsSendSessionEventsData is implemented by each variant of
// [ManagedAgentsSendSessionEventsDataUnion] to add type safety for the return
// type of [ManagedAgentsSendSessionEventsDataUnion.AsAny]
type anyManagedAgentsSendSessionEventsData interface {
	implManagedAgentsSendSessionEventsDataUnion()
}

func (ManagedAgentsUserMessageEvent) implManagedAgentsSendSessionEventsDataUnion()          {}
func (ManagedAgentsUserInterruptEvent) implManagedAgentsSendSessionEventsDataUnion()        {}
func (ManagedAgentsUserToolConfirmationEvent) implManagedAgentsSendSessionEventsDataUnion() {}
func (ManagedAgentsUserCustomToolResultEvent) implManagedAgentsSendSessionEventsDataUnion() {}
func (ManagedAgentsUserDefineOutcomeEvent) implManagedAgentsSendSessionEventsDataUnion()    {}
func (ManagedAgentsUserToolResultEvent) implManagedAgentsSendSessionEventsDataUnion()       {}
func (ManagedAgentsSystemMessageEvent) implManagedAgentsSendSessionEventsDataUnion()        {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsSendSessionEventsDataUnion.AsAny().(type) {
//	case qoder.ManagedAgentsUserMessageEvent:
//	case qoder.ManagedAgentsUserInterruptEvent:
//	case qoder.ManagedAgentsUserToolConfirmationEvent:
//	case qoder.ManagedAgentsUserCustomToolResultEvent:
//	case qoder.ManagedAgentsUserDefineOutcomeEvent:
//	case qoder.ManagedAgentsUserToolResultEvent:
//	case qoder.ManagedAgentsSystemMessageEvent:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsSendSessionEventsDataUnion) AsAny() anyManagedAgentsSendSessionEventsData {
	switch u.Type {
	case "user.message":
		return u.AsUserMessage()
	case "user.interrupt":
		return u.AsUserInterrupt()
	case "user.tool_confirmation":
		return u.AsUserToolConfirmation()
	case "user.custom_tool_result":
		return u.AsUserCustomToolResult()
	case "user.define_outcome":
		return u.AsUserDefineOutcome()
	case "user.tool_result":
		return u.AsUserToolResult()
	case "system.message":
		return u.AsSystemMessage()
	}
	return nil
}

func (u ManagedAgentsSendSessionEventsDataUnion) AsUserMessage() (v ManagedAgentsUserMessageEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSendSessionEventsDataUnion) AsUserInterrupt() (v ManagedAgentsUserInterruptEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSendSessionEventsDataUnion) AsUserToolConfirmation() (v ManagedAgentsUserToolConfirmationEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSendSessionEventsDataUnion) AsUserCustomToolResult() (v ManagedAgentsUserCustomToolResultEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSendSessionEventsDataUnion) AsUserDefineOutcome() (v ManagedAgentsUserDefineOutcomeEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSendSessionEventsDataUnion) AsUserToolResult() (v ManagedAgentsUserToolResultEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSendSessionEventsDataUnion) AsSystemMessage() (v ManagedAgentsSystemMessageEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsSendSessionEventsDataUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsSendSessionEventsDataUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsSendSessionEventsDataUnionContent is an implicit subunion of
// [ManagedAgentsSendSessionEventsDataUnion].
// ManagedAgentsSendSessionEventsDataUnionContent provides convenient access to
// the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsSendSessionEventsDataUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfManagedAgentsUserMessageEventContentArray
// OfManagedAgentsUserCustomToolResultEventContentArray
// OfManagedAgentsUserToolResultEventContentArray
// OfManagedAgentsSystemContentBlockArray]
type ManagedAgentsSendSessionEventsDataUnionContent struct {
	// This field will be present if the value is a
	// [[]ManagedAgentsUserMessageEventContentUnion] instead of an object.
	OfManagedAgentsUserMessageEventContentArray []ManagedAgentsUserMessageEventContentUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]ManagedAgentsUserCustomToolResultEventContentUnion] instead of an object.
	OfManagedAgentsUserCustomToolResultEventContentArray []ManagedAgentsUserCustomToolResultEventContentUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]ManagedAgentsUserToolResultEventContentUnion] instead of an object.
	OfManagedAgentsUserToolResultEventContentArray []ManagedAgentsUserToolResultEventContentUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]ManagedAgentsSystemContentBlock] instead of an object.
	OfManagedAgentsSystemContentBlockArray []ManagedAgentsSystemContentBlock `json:",inline"`
	JSON                                   struct {
		OfManagedAgentsUserMessageEventContentArray          respjson.Field
		OfManagedAgentsUserCustomToolResultEventContentArray respjson.Field
		OfManagedAgentsUserToolResultEventContentArray       respjson.Field
		OfManagedAgentsSystemContentBlockArray               respjson.Field
		raw                                                  string
	} `json:"-"`
}

func (r *ManagedAgentsSendSessionEventsDataUnionContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The agent stopped because the session's tracked list cost reached its budget, or
// because its usage includes a model with no list price (which the budget cannot
// measure). Raise the budget to continue — or, if raising is rejected because a
// model has no list price, remove the budget.
type ManagedAgentsSessionBudgetReached struct {
	// Any of "budget_reached".
	Type ManagedAgentsSessionBudgetReachedType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionBudgetReached) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionBudgetReached) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSessionBudgetReachedType string

const (
	ManagedAgentsSessionBudgetReachedTypeBudgetReached ManagedAgentsSessionBudgetReachedType = "budget_reached"
)

// Emitted when a session has been deleted. Terminates any active event stream — no
// further events will be emitted for this session.
type ManagedAgentsSessionDeletedEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "session.deleted".
	Type ManagedAgentsSessionDeletedEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ProcessedAt respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionDeletedEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionDeletedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSessionDeletedEventType string

const (
	ManagedAgentsSessionDeletedEventTypeSessionDeleted ManagedAgentsSessionDeletedEventType = "session.deleted"
)

// The agent completed its turn naturally and is ready for the next user message.
type ManagedAgentsSessionEndTurn struct {
	// Any of "end_turn".
	Type ManagedAgentsSessionEndTurnType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionEndTurn) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionEndTurn) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSessionEndTurnType string

const (
	ManagedAgentsSessionEndTurnTypeEndTurn ManagedAgentsSessionEndTurnType = "end_turn"
)

// An error event indicating a problem occurred during session execution.
type ManagedAgentsSessionErrorEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// An unknown or unexpected error occurred during session execution. A fallback
	// variant; clients that don't recognize a new error code can match on
	// `retry_status` and `message` alone.
	Error ManagedAgentsSessionErrorEventErrorUnion `json:"error" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "session.error".
	Type ManagedAgentsSessionErrorEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Error       respjson.Field
		ProcessedAt respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionErrorEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionErrorEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsSessionErrorEventErrorUnion contains all possible properties
// and values from [ManagedAgentsUnknownError],
// [ManagedAgentsModelOverloadedError],
// [ManagedAgentsModelRateLimitedError],
// [ManagedAgentsModelRequestFailedError],
// [ManagedAgentsMCPConnectionFailedError],
// [ManagedAgentsMCPAuthenticationFailedError],
// [ManagedAgentsBillingError],
// [ManagedAgentsCredentialHostUnreachableError].
//
// Use the [ManagedAgentsSessionErrorEventErrorUnion.AsAny] method to switch on
// the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsSessionErrorEventErrorUnion struct {
	Message string `json:"message"`
	// This field is a union of [ManagedAgentsUnknownErrorRetryStatusUnion],
	// [ManagedAgentsModelOverloadedErrorRetryStatusUnion],
	// [ManagedAgentsModelRateLimitedErrorRetryStatusUnion],
	// [ManagedAgentsModelRequestFailedErrorRetryStatusUnion],
	// [ManagedAgentsMCPConnectionFailedErrorRetryStatusUnion],
	// [ManagedAgentsMCPAuthenticationFailedErrorRetryStatusUnion],
	// [ManagedAgentsBillingErrorRetryStatusUnion],
	// [ManagedAgentsCredentialHostUnreachableErrorRetryStatusUnion]
	RetryStatus ManagedAgentsSessionErrorEventErrorUnionRetryStatus `json:"retry_status"`
	// Any of "unknown_error", "model_overloaded_error", "model_rate_limited_error",
	// "model_request_failed_error", "mcp_connection_failed_error",
	// "mcp_authentication_failed_error", "billing_error",
	// "credential_host_unreachable_error".
	Type          string `json:"type"`
	MCPServerName string `json:"mcp_server_name"`
	// This field is from variant [ManagedAgentsCredentialHostUnreachableError].
	CredentialID string `json:"credential_id"`
	// This field is from variant [ManagedAgentsCredentialHostUnreachableError].
	VaultID string `json:"vault_id"`
	JSON    struct {
		Message       respjson.Field
		RetryStatus   respjson.Field
		Type          respjson.Field
		MCPServerName respjson.Field
		CredentialID  respjson.Field
		VaultID       respjson.Field
		raw           string
	} `json:"-"`
}

// anyManagedAgentsSessionErrorEventError is implemented by each variant of
// [ManagedAgentsSessionErrorEventErrorUnion] to add type safety for the return
// type of [ManagedAgentsSessionErrorEventErrorUnion.AsAny]
type anyManagedAgentsSessionErrorEventError interface {
	implManagedAgentsSessionErrorEventErrorUnion()
}

func (ManagedAgentsUnknownError) implManagedAgentsSessionErrorEventErrorUnion()             {}
func (ManagedAgentsModelOverloadedError) implManagedAgentsSessionErrorEventErrorUnion()     {}
func (ManagedAgentsModelRateLimitedError) implManagedAgentsSessionErrorEventErrorUnion()    {}
func (ManagedAgentsModelRequestFailedError) implManagedAgentsSessionErrorEventErrorUnion()  {}
func (ManagedAgentsMCPConnectionFailedError) implManagedAgentsSessionErrorEventErrorUnion() {}
func (ManagedAgentsMCPAuthenticationFailedError) implManagedAgentsSessionErrorEventErrorUnion() {
}
func (ManagedAgentsBillingError) implManagedAgentsSessionErrorEventErrorUnion() {}
func (ManagedAgentsCredentialHostUnreachableError) implManagedAgentsSessionErrorEventErrorUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsSessionErrorEventErrorUnion.AsAny().(type) {
//	case qoder.ManagedAgentsUnknownError:
//	case qoder.ManagedAgentsModelOverloadedError:
//	case qoder.ManagedAgentsModelRateLimitedError:
//	case qoder.ManagedAgentsModelRequestFailedError:
//	case qoder.ManagedAgentsMCPConnectionFailedError:
//	case qoder.ManagedAgentsMCPAuthenticationFailedError:
//	case qoder.ManagedAgentsBillingError:
//	case qoder.ManagedAgentsCredentialHostUnreachableError:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsSessionErrorEventErrorUnion) AsAny() anyManagedAgentsSessionErrorEventError {
	switch u.Type {
	case "unknown_error":
		return u.AsUnknownError()
	case "model_overloaded_error":
		return u.AsModelOverloadedError()
	case "model_rate_limited_error":
		return u.AsModelRateLimitedError()
	case "model_request_failed_error":
		return u.AsModelRequestFailedError()
	case "mcp_connection_failed_error":
		return u.AsMCPConnectionFailedError()
	case "mcp_authentication_failed_error":
		return u.AsMCPAuthenticationFailedError()
	case "billing_error":
		return u.AsBillingError()
	case "credential_host_unreachable_error":
		return u.AsCredentialHostUnreachableError()
	}
	return nil
}

func (u ManagedAgentsSessionErrorEventErrorUnion) AsUnknownError() (v ManagedAgentsUnknownError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionErrorEventErrorUnion) AsModelOverloadedError() (v ManagedAgentsModelOverloadedError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionErrorEventErrorUnion) AsModelRateLimitedError() (v ManagedAgentsModelRateLimitedError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionErrorEventErrorUnion) AsModelRequestFailedError() (v ManagedAgentsModelRequestFailedError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionErrorEventErrorUnion) AsMCPConnectionFailedError() (v ManagedAgentsMCPConnectionFailedError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionErrorEventErrorUnion) AsMCPAuthenticationFailedError() (v ManagedAgentsMCPAuthenticationFailedError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionErrorEventErrorUnion) AsBillingError() (v ManagedAgentsBillingError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionErrorEventErrorUnion) AsCredentialHostUnreachableError() (v ManagedAgentsCredentialHostUnreachableError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsSessionErrorEventErrorUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsSessionErrorEventErrorUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsSessionErrorEventErrorUnionRetryStatus is an implicit subunion
// of [ManagedAgentsSessionErrorEventErrorUnion].
// ManagedAgentsSessionErrorEventErrorUnionRetryStatus provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsSessionErrorEventErrorUnion].
type ManagedAgentsSessionErrorEventErrorUnionRetryStatus struct {
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

func (r *ManagedAgentsSessionErrorEventErrorUnionRetryStatus) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSessionErrorEventType string

const (
	ManagedAgentsSessionErrorEventTypeSessionError ManagedAgentsSessionErrorEventType = "session.error"
)

// ManagedAgentsSessionEventUnion contains all possible properties and values
// from [ManagedAgentsUserMessageEvent], [ManagedAgentsUserInterruptEvent],
// [ManagedAgentsUserToolConfirmationEvent],
// [ManagedAgentsUserCustomToolResultEvent],
// [ManagedAgentsAgentCustomToolUseEvent],
// [ManagedAgentsAgentMessageEvent], [ManagedAgentsAgentThinkingEvent],
// [ManagedAgentsAgentMCPToolUseEvent],
// [ManagedAgentsAgentMCPToolResultEvent],
// [ManagedAgentsAgentToolUseEvent], [ManagedAgentsAgentToolResultEvent],
// [ManagedAgentsAgentThreadMessageReceivedEvent],
// [ManagedAgentsAgentThreadMessageSentEvent],
// [ManagedAgentsAgentThreadContextCompactedEvent],
// [ManagedAgentsSessionErrorEvent],
// [ManagedAgentsSessionStatusRescheduledEvent],
// [ManagedAgentsSessionStatusRunningEvent],
// [ManagedAgentsSessionStatusIdleEvent],
// [ManagedAgentsSessionStatusTerminatedEvent],
// [ManagedAgentsSessionThreadCreatedEvent],
// [ManagedAgentsSpanOutcomeEvaluationStartEvent],
// [ManagedAgentsSpanOutcomeEvaluationEndEvent],
// [ManagedAgentsSpanModelRequestStartEvent],
// [ManagedAgentsSpanModelRequestEndEvent],
// [ManagedAgentsSpanOutcomeEvaluationOngoingEvent],
// [ManagedAgentsUserDefineOutcomeEvent],
// [ManagedAgentsSessionDeletedEvent],
// [ManagedAgentsSessionThreadStatusRunningEvent],
// [ManagedAgentsSessionThreadStatusIdleEvent],
// [ManagedAgentsSessionThreadStatusTerminatedEvent],
// [ManagedAgentsUserToolResultEvent],
// [ManagedAgentsSessionThreadStatusRescheduledEvent],
// [ManagedAgentsSessionUpdatedEvent], [ManagedAgentsSystemMessageEvent],
// [ManagedAgentsSessionUsageEvent].
//
// Use the [ManagedAgentsSessionEventUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsSessionEventUnion struct {
	ID string `json:"id"`
	// This field is a union of [[]ManagedAgentsUserMessageEventContentUnion],
	// [[]ManagedAgentsUserCustomToolResultEventContentUnion],
	// [[]ManagedAgentsAgentMessageEventContentUnion],
	// [[]ManagedAgentsAgentMCPToolResultEventContentUnion],
	// [[]ManagedAgentsAgentToolResultEventContentUnion],
	// [[]ManagedAgentsAgentThreadMessageReceivedEventContentUnion],
	// [[]ManagedAgentsAgentThreadMessageSentEventContentUnion],
	// [[]ManagedAgentsUserToolResultEventContentUnion],
	// [[]ManagedAgentsSystemContentBlock]
	Content ManagedAgentsSessionEventUnionContent `json:"content"`
	// Any of "user.message", "user.interrupt", "user.tool_confirmation",
	// "user.custom_tool_result", "agent.custom_tool_use", "agent.message",
	// "agent.thinking", "agent.mcp_tool_use", "agent.mcp_tool_result",
	// "agent.tool_use", "agent.tool_result", "agent.thread_message_received",
	// "agent.thread_message_sent", "agent.thread_context_compacted", "session.error",
	// "session.status_rescheduled", "session.status_running", "session.status_idle",
	// "session.status_terminated", "session.thread_created",
	// "span.outcome_evaluation_start", "span.outcome_evaluation_end",
	// "span.model_request_start", "span.model_request_end",
	// "span.outcome_evaluation_ongoing", "user.define_outcome", "session.deleted",
	// "session.thread_status_running", "session.thread_status_idle",
	// "session.thread_status_terminated", "user.tool_result",
	// "session.thread_status_rescheduled", "session.updated", "system.message",
	// "session.usage".
	Type            string    `json:"type"`
	ProcessedAt     time.Time `json:"processed_at"`
	SessionThreadID string    `json:"session_thread_id"`
	Result          string    `json:"result"`
	ToolUseID       string    `json:"tool_use_id"`
	// This field is from variant [ManagedAgentsUserToolConfirmationEvent].
	DenyMessage string `json:"deny_message"`
	// This field is from variant [ManagedAgentsUserCustomToolResultEvent].
	CustomToolUseID string `json:"custom_tool_use_id"`
	IsError         bool   `json:"is_error"`
	Input           any    `json:"input"`
	Name            string `json:"name"`
	// This field is from variant [ManagedAgentsAgentMCPToolUseEvent].
	MCPServerName       string `json:"mcp_server_name"`
	EvaluatedPermission string `json:"evaluated_permission"`
	// This field is from variant [ManagedAgentsAgentMCPToolResultEvent].
	MCPToolUseID string `json:"mcp_tool_use_id"`
	// This field is from variant [ManagedAgentsAgentThreadMessageReceivedEvent].
	FromSessionThreadID string `json:"from_session_thread_id"`
	// This field is from variant [ManagedAgentsAgentThreadMessageReceivedEvent].
	FromAgentName string `json:"from_agent_name"`
	// This field is from variant [ManagedAgentsAgentThreadMessageSentEvent].
	ToSessionThreadID string `json:"to_session_thread_id"`
	// This field is from variant [ManagedAgentsAgentThreadMessageSentEvent].
	ToAgentName string `json:"to_agent_name"`
	// This field is from variant [ManagedAgentsSessionErrorEvent].
	Error ManagedAgentsSessionErrorEventErrorUnion `json:"error"`
	// This field is a union of
	// [ManagedAgentsSessionStatusIdleEventStopReasonUnion],
	// [ManagedAgentsSessionThreadStatusIdleEventStopReasonUnion]
	StopReason ManagedAgentsSessionEventUnionStopReason `json:"stop_reason"`
	AgentName  string                                   `json:"agent_name"`
	Iteration  int64                                    `json:"iteration"`
	OutcomeID  string                                   `json:"outcome_id"`
	// This field is from variant [ManagedAgentsSpanOutcomeEvaluationEndEvent].
	Explanation string `json:"explanation"`
	// This field is from variant [ManagedAgentsSpanOutcomeEvaluationEndEvent].
	OutcomeEvaluationStartID string `json:"outcome_evaluation_start_id"`
	// This field is a union of [ManagedAgentsSpanModelUsage],
	// [ManagedAgentsSessionUsageSnapshot]
	Usage ManagedAgentsSessionEventUnionUsage `json:"usage"`
	// This field is from variant [ManagedAgentsSpanModelRequestEndEvent].
	ModelRequestStartID string `json:"model_request_start_id"`
	// This field is from variant [ManagedAgentsSpanModelRequestEndEvent].
	ModelUsage ManagedAgentsSpanModelUsage `json:"model_usage"`
	// This field is from variant [ManagedAgentsUserDefineOutcomeEvent].
	Description string `json:"description"`
	// This field is from variant [ManagedAgentsUserDefineOutcomeEvent].
	MaxIterations int64 `json:"max_iterations"`
	// This field is from variant [ManagedAgentsUserDefineOutcomeEvent].
	Rubric ManagedAgentsUserDefineOutcomeEventRubricUnion `json:"rubric"`
	// This field is from variant [ManagedAgentsSessionUpdatedEvent].
	Agent ManagedAgentsSessionAgent `json:"agent"`
	// This field is from variant [ManagedAgentsSessionUpdatedEvent].
	Budget ManagedAgentsBudgetLimit `json:"budget"`
	// This field is from variant [ManagedAgentsSessionUpdatedEvent].
	Metadata map[string]string `json:"metadata"`
	// This field is from variant [ManagedAgentsSessionUpdatedEvent].
	Title string `json:"title"`
	JSON  struct {
		ID                       respjson.Field
		Content                  respjson.Field
		Type                     respjson.Field
		ProcessedAt              respjson.Field
		SessionThreadID          respjson.Field
		Result                   respjson.Field
		ToolUseID                respjson.Field
		DenyMessage              respjson.Field
		CustomToolUseID          respjson.Field
		IsError                  respjson.Field
		Input                    respjson.Field
		Name                     respjson.Field
		MCPServerName            respjson.Field
		EvaluatedPermission      respjson.Field
		MCPToolUseID             respjson.Field
		FromSessionThreadID      respjson.Field
		FromAgentName            respjson.Field
		ToSessionThreadID        respjson.Field
		ToAgentName              respjson.Field
		Error                    respjson.Field
		StopReason               respjson.Field
		AgentName                respjson.Field
		Iteration                respjson.Field
		OutcomeID                respjson.Field
		Explanation              respjson.Field
		OutcomeEvaluationStartID respjson.Field
		Usage                    respjson.Field
		ModelRequestStartID      respjson.Field
		ModelUsage               respjson.Field
		Description              respjson.Field
		MaxIterations            respjson.Field
		Rubric                   respjson.Field
		Agent                    respjson.Field
		Budget                   respjson.Field
		Metadata                 respjson.Field
		Title                    respjson.Field
		raw                      string
	} `json:"-"`
}

// anyManagedAgentsSessionEvent is implemented by each variant of
// [ManagedAgentsSessionEventUnion] to add type safety for the return type of
// [ManagedAgentsSessionEventUnion.AsAny]
type anyManagedAgentsSessionEvent interface {
	implManagedAgentsSessionEventUnion()
}

func (ManagedAgentsUserMessageEvent) implManagedAgentsSessionEventUnion()                   {}
func (ManagedAgentsUserInterruptEvent) implManagedAgentsSessionEventUnion()                 {}
func (ManagedAgentsUserToolConfirmationEvent) implManagedAgentsSessionEventUnion()          {}
func (ManagedAgentsUserCustomToolResultEvent) implManagedAgentsSessionEventUnion()          {}
func (ManagedAgentsAgentCustomToolUseEvent) implManagedAgentsSessionEventUnion()            {}
func (ManagedAgentsAgentMessageEvent) implManagedAgentsSessionEventUnion()                  {}
func (ManagedAgentsAgentThinkingEvent) implManagedAgentsSessionEventUnion()                 {}
func (ManagedAgentsAgentMCPToolUseEvent) implManagedAgentsSessionEventUnion()               {}
func (ManagedAgentsAgentMCPToolResultEvent) implManagedAgentsSessionEventUnion()            {}
func (ManagedAgentsAgentToolUseEvent) implManagedAgentsSessionEventUnion()                  {}
func (ManagedAgentsAgentToolResultEvent) implManagedAgentsSessionEventUnion()               {}
func (ManagedAgentsAgentThreadMessageReceivedEvent) implManagedAgentsSessionEventUnion()    {}
func (ManagedAgentsAgentThreadMessageSentEvent) implManagedAgentsSessionEventUnion()        {}
func (ManagedAgentsAgentThreadContextCompactedEvent) implManagedAgentsSessionEventUnion()   {}
func (ManagedAgentsSessionErrorEvent) implManagedAgentsSessionEventUnion()                  {}
func (ManagedAgentsSessionStatusRescheduledEvent) implManagedAgentsSessionEventUnion()      {}
func (ManagedAgentsSessionStatusRunningEvent) implManagedAgentsSessionEventUnion()          {}
func (ManagedAgentsSessionStatusIdleEvent) implManagedAgentsSessionEventUnion()             {}
func (ManagedAgentsSessionStatusTerminatedEvent) implManagedAgentsSessionEventUnion()       {}
func (ManagedAgentsSessionThreadCreatedEvent) implManagedAgentsSessionEventUnion()          {}
func (ManagedAgentsSpanOutcomeEvaluationStartEvent) implManagedAgentsSessionEventUnion()    {}
func (ManagedAgentsSpanOutcomeEvaluationEndEvent) implManagedAgentsSessionEventUnion()      {}
func (ManagedAgentsSpanModelRequestStartEvent) implManagedAgentsSessionEventUnion()         {}
func (ManagedAgentsSpanModelRequestEndEvent) implManagedAgentsSessionEventUnion()           {}
func (ManagedAgentsSpanOutcomeEvaluationOngoingEvent) implManagedAgentsSessionEventUnion()  {}
func (ManagedAgentsUserDefineOutcomeEvent) implManagedAgentsSessionEventUnion()             {}
func (ManagedAgentsSessionDeletedEvent) implManagedAgentsSessionEventUnion()                {}
func (ManagedAgentsSessionThreadStatusRunningEvent) implManagedAgentsSessionEventUnion()    {}
func (ManagedAgentsSessionThreadStatusIdleEvent) implManagedAgentsSessionEventUnion()       {}
func (ManagedAgentsSessionThreadStatusTerminatedEvent) implManagedAgentsSessionEventUnion() {}
func (ManagedAgentsUserToolResultEvent) implManagedAgentsSessionEventUnion()                {}
func (ManagedAgentsSessionThreadStatusRescheduledEvent) implManagedAgentsSessionEventUnion() {
}
func (ManagedAgentsSessionUpdatedEvent) implManagedAgentsSessionEventUnion() {}
func (ManagedAgentsSystemMessageEvent) implManagedAgentsSessionEventUnion()  {}
func (ManagedAgentsSessionUsageEvent) implManagedAgentsSessionEventUnion()   {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsSessionEventUnion.AsAny().(type) {
//	case qoder.ManagedAgentsUserMessageEvent:
//	case qoder.ManagedAgentsUserInterruptEvent:
//	case qoder.ManagedAgentsUserToolConfirmationEvent:
//	case qoder.ManagedAgentsUserCustomToolResultEvent:
//	case qoder.ManagedAgentsAgentCustomToolUseEvent:
//	case qoder.ManagedAgentsAgentMessageEvent:
//	case qoder.ManagedAgentsAgentThinkingEvent:
//	case qoder.ManagedAgentsAgentMCPToolUseEvent:
//	case qoder.ManagedAgentsAgentMCPToolResultEvent:
//	case qoder.ManagedAgentsAgentToolUseEvent:
//	case qoder.ManagedAgentsAgentToolResultEvent:
//	case qoder.ManagedAgentsAgentThreadMessageReceivedEvent:
//	case qoder.ManagedAgentsAgentThreadMessageSentEvent:
//	case qoder.ManagedAgentsAgentThreadContextCompactedEvent:
//	case qoder.ManagedAgentsSessionErrorEvent:
//	case qoder.ManagedAgentsSessionStatusRescheduledEvent:
//	case qoder.ManagedAgentsSessionStatusRunningEvent:
//	case qoder.ManagedAgentsSessionStatusIdleEvent:
//	case qoder.ManagedAgentsSessionStatusTerminatedEvent:
//	case qoder.ManagedAgentsSessionThreadCreatedEvent:
//	case qoder.ManagedAgentsSpanOutcomeEvaluationStartEvent:
//	case qoder.ManagedAgentsSpanOutcomeEvaluationEndEvent:
//	case qoder.ManagedAgentsSpanModelRequestStartEvent:
//	case qoder.ManagedAgentsSpanModelRequestEndEvent:
//	case qoder.ManagedAgentsSpanOutcomeEvaluationOngoingEvent:
//	case qoder.ManagedAgentsUserDefineOutcomeEvent:
//	case qoder.ManagedAgentsSessionDeletedEvent:
//	case qoder.ManagedAgentsSessionThreadStatusRunningEvent:
//	case qoder.ManagedAgentsSessionThreadStatusIdleEvent:
//	case qoder.ManagedAgentsSessionThreadStatusTerminatedEvent:
//	case qoder.ManagedAgentsUserToolResultEvent:
//	case qoder.ManagedAgentsSessionThreadStatusRescheduledEvent:
//	case qoder.ManagedAgentsSessionUpdatedEvent:
//	case qoder.ManagedAgentsSystemMessageEvent:
//	case qoder.ManagedAgentsSessionUsageEvent:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsSessionEventUnion) AsAny() anyManagedAgentsSessionEvent {
	switch u.Type {
	case "user.message":
		return u.AsUserMessage()
	case "user.interrupt":
		return u.AsUserInterrupt()
	case "user.tool_confirmation":
		return u.AsUserToolConfirmation()
	case "user.custom_tool_result":
		return u.AsUserCustomToolResult()
	case "agent.custom_tool_use":
		return u.AsAgentCustomToolUse()
	case "agent.message":
		return u.AsAgentMessage()
	case "agent.thinking":
		return u.AsAgentThinking()
	case "agent.mcp_tool_use":
		return u.AsAgentMCPToolUse()
	case "agent.mcp_tool_result":
		return u.AsAgentMCPToolResult()
	case "agent.tool_use":
		return u.AsAgentToolUse()
	case "agent.tool_result":
		return u.AsAgentToolResult()
	case "agent.thread_message_received":
		return u.AsAgentThreadMessageReceived()
	case "agent.thread_message_sent":
		return u.AsAgentThreadMessageSent()
	case "agent.thread_context_compacted":
		return u.AsAgentThreadContextCompacted()
	case "session.error":
		return u.AsSessionError()
	case "session.status_rescheduled":
		return u.AsSessionStatusRescheduled()
	case "session.status_running":
		return u.AsSessionStatusRunning()
	case "session.status_idle":
		return u.AsSessionStatusIdle()
	case "session.status_terminated":
		return u.AsSessionStatusTerminated()
	case "session.thread_created":
		return u.AsSessionThreadCreated()
	case "span.outcome_evaluation_start":
		return u.AsSpanOutcomeEvaluationStart()
	case "span.outcome_evaluation_end":
		return u.AsSpanOutcomeEvaluationEnd()
	case "span.model_request_start":
		return u.AsSpanModelRequestStart()
	case "span.model_request_end":
		return u.AsSpanModelRequestEnd()
	case "span.outcome_evaluation_ongoing":
		return u.AsSpanOutcomeEvaluationOngoing()
	case "user.define_outcome":
		return u.AsUserDefineOutcome()
	case "session.deleted":
		return u.AsSessionDeleted()
	case "session.thread_status_running":
		return u.AsSessionThreadStatusRunning()
	case "session.thread_status_idle":
		return u.AsSessionThreadStatusIdle()
	case "session.thread_status_terminated":
		return u.AsSessionThreadStatusTerminated()
	case "user.tool_result":
		return u.AsUserToolResult()
	case "session.thread_status_rescheduled":
		return u.AsSessionThreadStatusRescheduled()
	case "session.updated":
		return u.AsSessionUpdated()
	case "system.message":
		return u.AsSystemMessage()
	case "session.usage":
		return u.AsSessionUsage()
	}
	return nil
}

func (u ManagedAgentsSessionEventUnion) AsUserMessage() (v ManagedAgentsUserMessageEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsUserInterrupt() (v ManagedAgentsUserInterruptEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsUserToolConfirmation() (v ManagedAgentsUserToolConfirmationEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsUserCustomToolResult() (v ManagedAgentsUserCustomToolResultEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsAgentCustomToolUse() (v ManagedAgentsAgentCustomToolUseEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsAgentMessage() (v ManagedAgentsAgentMessageEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsAgentThinking() (v ManagedAgentsAgentThinkingEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsAgentMCPToolUse() (v ManagedAgentsAgentMCPToolUseEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsAgentMCPToolResult() (v ManagedAgentsAgentMCPToolResultEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsAgentToolUse() (v ManagedAgentsAgentToolUseEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsAgentToolResult() (v ManagedAgentsAgentToolResultEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsAgentThreadMessageReceived() (v ManagedAgentsAgentThreadMessageReceivedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsAgentThreadMessageSent() (v ManagedAgentsAgentThreadMessageSentEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsAgentThreadContextCompacted() (v ManagedAgentsAgentThreadContextCompactedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsSessionError() (v ManagedAgentsSessionErrorEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsSessionStatusRescheduled() (v ManagedAgentsSessionStatusRescheduledEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsSessionStatusRunning() (v ManagedAgentsSessionStatusRunningEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsSessionStatusIdle() (v ManagedAgentsSessionStatusIdleEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsSessionStatusTerminated() (v ManagedAgentsSessionStatusTerminatedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsSessionThreadCreated() (v ManagedAgentsSessionThreadCreatedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsSpanOutcomeEvaluationStart() (v ManagedAgentsSpanOutcomeEvaluationStartEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsSpanOutcomeEvaluationEnd() (v ManagedAgentsSpanOutcomeEvaluationEndEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsSpanModelRequestStart() (v ManagedAgentsSpanModelRequestStartEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsSpanModelRequestEnd() (v ManagedAgentsSpanModelRequestEndEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsSpanOutcomeEvaluationOngoing() (v ManagedAgentsSpanOutcomeEvaluationOngoingEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsUserDefineOutcome() (v ManagedAgentsUserDefineOutcomeEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsSessionDeleted() (v ManagedAgentsSessionDeletedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsSessionThreadStatusRunning() (v ManagedAgentsSessionThreadStatusRunningEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsSessionThreadStatusIdle() (v ManagedAgentsSessionThreadStatusIdleEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsSessionThreadStatusTerminated() (v ManagedAgentsSessionThreadStatusTerminatedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsUserToolResult() (v ManagedAgentsUserToolResultEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsSessionThreadStatusRescheduled() (v ManagedAgentsSessionThreadStatusRescheduledEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsSessionUpdated() (v ManagedAgentsSessionUpdatedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsSystemMessage() (v ManagedAgentsSystemMessageEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionEventUnion) AsSessionUsage() (v ManagedAgentsSessionUsageEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsSessionEventUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsSessionEventUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsSessionEventUnionContent is an implicit subunion of
// [ManagedAgentsSessionEventUnion]. ManagedAgentsSessionEventUnionContent
// provides convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsSessionEventUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfManagedAgentsUserMessageEventContentArray
// OfManagedAgentsUserCustomToolResultEventContentArray
// OfManagedAgentsAgentMessageEventContentArray
// OfManagedAgentsAgentMCPToolResultEventContentArray
// OfManagedAgentsAgentToolResultEventContentArray
// OfManagedAgentsAgentThreadMessageReceivedEventContentArray
// OfManagedAgentsAgentThreadMessageSentEventContentArray
// OfManagedAgentsUserToolResultEventContentArray
// OfManagedAgentsSystemContentBlockArray]
type ManagedAgentsSessionEventUnionContent struct {
	// This field will be present if the value is a
	// [[]ManagedAgentsUserMessageEventContentUnion] instead of an object.
	OfManagedAgentsUserMessageEventContentArray []ManagedAgentsUserMessageEventContentUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]ManagedAgentsUserCustomToolResultEventContentUnion] instead of an object.
	OfManagedAgentsUserCustomToolResultEventContentArray []ManagedAgentsUserCustomToolResultEventContentUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]ManagedAgentsAgentMessageEventContentUnion] instead of an object.
	OfManagedAgentsAgentMessageEventContentArray []ManagedAgentsAgentMessageEventContentUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]ManagedAgentsAgentMCPToolResultEventContentUnion] instead of an object.
	OfManagedAgentsAgentMCPToolResultEventContentArray []ManagedAgentsAgentMCPToolResultEventContentUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]ManagedAgentsAgentToolResultEventContentUnion] instead of an object.
	OfManagedAgentsAgentToolResultEventContentArray []ManagedAgentsAgentToolResultEventContentUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]ManagedAgentsAgentThreadMessageReceivedEventContentUnion] instead of an
	// object.
	OfManagedAgentsAgentThreadMessageReceivedEventContentArray []ManagedAgentsAgentThreadMessageReceivedEventContentUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]ManagedAgentsAgentThreadMessageSentEventContentUnion] instead of an
	// object.
	OfManagedAgentsAgentThreadMessageSentEventContentArray []ManagedAgentsAgentThreadMessageSentEventContentUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]ManagedAgentsUserToolResultEventContentUnion] instead of an object.
	OfManagedAgentsUserToolResultEventContentArray []ManagedAgentsUserToolResultEventContentUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]ManagedAgentsSystemContentBlock] instead of an object.
	OfManagedAgentsSystemContentBlockArray []ManagedAgentsSystemContentBlock `json:",inline"`
	JSON                                   struct {
		OfManagedAgentsUserMessageEventContentArray                respjson.Field
		OfManagedAgentsUserCustomToolResultEventContentArray       respjson.Field
		OfManagedAgentsAgentMessageEventContentArray               respjson.Field
		OfManagedAgentsAgentMCPToolResultEventContentArray         respjson.Field
		OfManagedAgentsAgentToolResultEventContentArray            respjson.Field
		OfManagedAgentsAgentThreadMessageReceivedEventContentArray respjson.Field
		OfManagedAgentsAgentThreadMessageSentEventContentArray     respjson.Field
		OfManagedAgentsUserToolResultEventContentArray             respjson.Field
		OfManagedAgentsSystemContentBlockArray                     respjson.Field
		raw                                                        string
	} `json:"-"`
}

func (r *ManagedAgentsSessionEventUnionContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsSessionEventUnionStopReason is an implicit subunion of
// [ManagedAgentsSessionEventUnion].
// ManagedAgentsSessionEventUnionStopReason provides convenient access to the
// sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsSessionEventUnion].
type ManagedAgentsSessionEventUnionStopReason struct {
	Type string `json:"type"`
	// This field is from variant
	// [ManagedAgentsSessionStatusIdleEventStopReasonUnion],
	// [ManagedAgentsSessionThreadStatusIdleEventStopReasonUnion].
	EventIDs []string `json:"event_ids"`
	JSON     struct {
		Type     respjson.Field
		EventIDs respjson.Field
		raw      string
	} `json:"-"`
}

func (r *ManagedAgentsSessionEventUnionStopReason) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsSessionEventUnionUsage is an implicit subunion of
// [ManagedAgentsSessionEventUnion]. ManagedAgentsSessionEventUnionUsage
// provides convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsSessionEventUnion].
type ManagedAgentsSessionEventUnionUsage struct {
	// This field is from variant [ManagedAgentsSpanModelUsage].
	CacheCreationInputTokens int64 `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int64 `json:"cache_read_input_tokens"`
	InputTokens              int64 `json:"input_tokens"`
	OutputTokens             int64 `json:"output_tokens"`
	// This field is from variant [ManagedAgentsSpanModelUsage].
	Speed ManagedAgentsSpanModelUsageSpeed `json:"speed"`
	// This field is from variant [ManagedAgentsSessionUsageSnapshot].
	ActiveSeconds float64 `json:"active_seconds"`
	// This field is from variant [ManagedAgentsSessionUsageSnapshot].
	CacheCreation ManagedAgentsCacheCreationUsage `json:"cache_creation"`
	// This field is from variant [ManagedAgentsSessionUsageSnapshot].
	ListCost MonetaryAmount `json:"list_cost"`
	// This field is from variant [ManagedAgentsSessionUsageSnapshot].
	ServerToolUse ManagedAgentsServerToolUsage `json:"server_tool_use"`
	JSON          struct {
		CacheCreationInputTokens respjson.Field
		CacheReadInputTokens     respjson.Field
		InputTokens              respjson.Field
		OutputTokens             respjson.Field
		Speed                    respjson.Field
		ActiveSeconds            respjson.Field
		CacheCreation            respjson.Field
		ListCost                 respjson.Field
		ServerToolUse            respjson.Field
		raw                      string
	} `json:"-"`
}

func (r *ManagedAgentsSessionEventUnionUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The agent is idle waiting on one or more blocking user-input events (tool
// confirmation, custom tool result, etc.). Resolving all of them transitions the
// session back to running.
type ManagedAgentsSessionRequiresAction struct {
	// The ids of events the agent is blocked on. Resolving fewer than all re-emits
	// `session.status_idle` with the remainder.
	EventIDs []string `json:"event_ids" api:"required"`
	// Any of "requires_action".
	Type ManagedAgentsSessionRequiresActionType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EventIDs    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionRequiresAction) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionRequiresAction) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSessionRequiresActionType string

const (
	ManagedAgentsSessionRequiresActionTypeRequiresAction ManagedAgentsSessionRequiresActionType = "requires_action"
)

// The turn ended because repeated errors exhausted the retry budget or an error
// escalated to `retry_status: 'exhausted'`.
type ManagedAgentsSessionRetriesExhausted struct {
	// Any of "retries_exhausted".
	Type ManagedAgentsSessionRetriesExhaustedType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionRetriesExhausted) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionRetriesExhausted) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSessionRetriesExhaustedType string

const (
	ManagedAgentsSessionRetriesExhaustedTypeRetriesExhausted ManagedAgentsSessionRetriesExhaustedType = "retries_exhausted"
)

// Indicates the agent has paused and is awaiting user input.
type ManagedAgentsSessionStatusIdleEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// The agent completed its turn naturally and is ready for the next user message.
	StopReason ManagedAgentsSessionStatusIdleEventStopReasonUnion `json:"stop_reason" api:"required"`
	// Any of "session.status_idle".
	Type ManagedAgentsSessionStatusIdleEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ProcessedAt respjson.Field
		StopReason  respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionStatusIdleEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionStatusIdleEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsSessionStatusIdleEventStopReasonUnion contains all possible
// properties and values from [ManagedAgentsSessionEndTurn],
// [ManagedAgentsSessionRequiresAction],
// [ManagedAgentsSessionRetriesExhausted],
// [ManagedAgentsSessionBudgetReached].
//
// Use the [ManagedAgentsSessionStatusIdleEventStopReasonUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsSessionStatusIdleEventStopReasonUnion struct {
	// Any of "end_turn", "requires_action", "retries_exhausted", "budget_reached".
	Type string `json:"type"`
	// This field is from variant [ManagedAgentsSessionRequiresAction].
	EventIDs []string `json:"event_ids"`
	JSON     struct {
		Type     respjson.Field
		EventIDs respjson.Field
		raw      string
	} `json:"-"`
}

// anyManagedAgentsSessionStatusIdleEventStopReason is implemented by each
// variant of [ManagedAgentsSessionStatusIdleEventStopReasonUnion] to add type
// safety for the return type of
// [ManagedAgentsSessionStatusIdleEventStopReasonUnion.AsAny]
type anyManagedAgentsSessionStatusIdleEventStopReason interface {
	implManagedAgentsSessionStatusIdleEventStopReasonUnion()
}

func (ManagedAgentsSessionEndTurn) implManagedAgentsSessionStatusIdleEventStopReasonUnion() {}
func (ManagedAgentsSessionRequiresAction) implManagedAgentsSessionStatusIdleEventStopReasonUnion() {
}
func (ManagedAgentsSessionRetriesExhausted) implManagedAgentsSessionStatusIdleEventStopReasonUnion() {
}
func (ManagedAgentsSessionBudgetReached) implManagedAgentsSessionStatusIdleEventStopReasonUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsSessionStatusIdleEventStopReasonUnion.AsAny().(type) {
//	case qoder.ManagedAgentsSessionEndTurn:
//	case qoder.ManagedAgentsSessionRequiresAction:
//	case qoder.ManagedAgentsSessionRetriesExhausted:
//	case qoder.ManagedAgentsSessionBudgetReached:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsSessionStatusIdleEventStopReasonUnion) AsAny() anyManagedAgentsSessionStatusIdleEventStopReason {
	switch u.Type {
	case "end_turn":
		return u.AsEndTurn()
	case "requires_action":
		return u.AsRequiresAction()
	case "retries_exhausted":
		return u.AsRetriesExhausted()
	case "budget_reached":
		return u.AsBudgetReached()
	}
	return nil
}

func (u ManagedAgentsSessionStatusIdleEventStopReasonUnion) AsEndTurn() (v ManagedAgentsSessionEndTurn) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionStatusIdleEventStopReasonUnion) AsRequiresAction() (v ManagedAgentsSessionRequiresAction) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionStatusIdleEventStopReasonUnion) AsRetriesExhausted() (v ManagedAgentsSessionRetriesExhausted) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionStatusIdleEventStopReasonUnion) AsBudgetReached() (v ManagedAgentsSessionBudgetReached) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsSessionStatusIdleEventStopReasonUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsSessionStatusIdleEventStopReasonUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSessionStatusIdleEventType string

const (
	ManagedAgentsSessionStatusIdleEventTypeSessionStatusIdle ManagedAgentsSessionStatusIdleEventType = "session.status_idle"
)

// Indicates the session is recovering from an error state and is rescheduled for
// execution.
type ManagedAgentsSessionStatusRescheduledEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "session.status_rescheduled".
	Type ManagedAgentsSessionStatusRescheduledEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ProcessedAt respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionStatusRescheduledEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionStatusRescheduledEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSessionStatusRescheduledEventType string

const (
	ManagedAgentsSessionStatusRescheduledEventTypeSessionStatusRescheduled ManagedAgentsSessionStatusRescheduledEventType = "session.status_rescheduled"
)

// Indicates the session is actively running and the agent is working.
type ManagedAgentsSessionStatusRunningEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "session.status_running".
	Type ManagedAgentsSessionStatusRunningEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ProcessedAt respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionStatusRunningEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionStatusRunningEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSessionStatusRunningEventType string

const (
	ManagedAgentsSessionStatusRunningEventTypeSessionStatusRunning ManagedAgentsSessionStatusRunningEventType = "session.status_running"
)

// Indicates the session has terminated, either due to an error or completion.
type ManagedAgentsSessionStatusTerminatedEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "session.status_terminated".
	Type ManagedAgentsSessionStatusTerminatedEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ProcessedAt respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionStatusTerminatedEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionStatusTerminatedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSessionStatusTerminatedEventType string

const (
	ManagedAgentsSessionStatusTerminatedEventTypeSessionStatusTerminated ManagedAgentsSessionStatusTerminatedEventType = "session.status_terminated"
)

// Emitted when a subagent is spawned as a new thread. Written to the parent
// thread's output stream so clients observing the session see child creation.
type ManagedAgentsSessionThreadCreatedEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// Name of the callable agent the thread runs.
	AgentName string `json:"agent_name" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Public `sthr_` ID of the newly created thread.
	SessionThreadID string `json:"session_thread_id" api:"required"`
	// Any of "session.thread_created".
	Type ManagedAgentsSessionThreadCreatedEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID              respjson.Field
		AgentName       respjson.Field
		ProcessedAt     respjson.Field
		SessionThreadID respjson.Field
		Type            respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionThreadCreatedEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionThreadCreatedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSessionThreadCreatedEventType string

const (
	ManagedAgentsSessionThreadCreatedEventTypeSessionThreadCreated ManagedAgentsSessionThreadCreatedEventType = "session.thread_created"
)

// A session thread has yielded and is awaiting input. Emitted on the thread's own
// stream and cross-posted to the primary stream for child threads.
type ManagedAgentsSessionThreadStatusIdleEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// Name of the agent the thread runs.
	AgentName string `json:"agent_name" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Public sthr\_ ID of the thread that went idle.
	SessionThreadID string `json:"session_thread_id" api:"required"`
	// The agent completed its turn naturally and is ready for the next user message.
	StopReason ManagedAgentsSessionThreadStatusIdleEventStopReasonUnion `json:"stop_reason" api:"required"`
	// Any of "session.thread_status_idle".
	Type ManagedAgentsSessionThreadStatusIdleEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID              respjson.Field
		AgentName       respjson.Field
		ProcessedAt     respjson.Field
		SessionThreadID respjson.Field
		StopReason      respjson.Field
		Type            respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionThreadStatusIdleEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionThreadStatusIdleEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsSessionThreadStatusIdleEventStopReasonUnion contains all
// possible properties and values from [ManagedAgentsSessionEndTurn],
// [ManagedAgentsSessionRequiresAction],
// [ManagedAgentsSessionRetriesExhausted],
// [ManagedAgentsSessionBudgetReached].
//
// Use the [ManagedAgentsSessionThreadStatusIdleEventStopReasonUnion.AsAny]
// method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsSessionThreadStatusIdleEventStopReasonUnion struct {
	// Any of "end_turn", "requires_action", "retries_exhausted", "budget_reached".
	Type string `json:"type"`
	// This field is from variant [ManagedAgentsSessionRequiresAction].
	EventIDs []string `json:"event_ids"`
	JSON     struct {
		Type     respjson.Field
		EventIDs respjson.Field
		raw      string
	} `json:"-"`
}

// anyManagedAgentsSessionThreadStatusIdleEventStopReason is implemented by
// each variant of [ManagedAgentsSessionThreadStatusIdleEventStopReasonUnion]
// to add type safety for the return type of
// [ManagedAgentsSessionThreadStatusIdleEventStopReasonUnion.AsAny]
type anyManagedAgentsSessionThreadStatusIdleEventStopReason interface {
	implManagedAgentsSessionThreadStatusIdleEventStopReasonUnion()
}

func (ManagedAgentsSessionEndTurn) implManagedAgentsSessionThreadStatusIdleEventStopReasonUnion() {
}
func (ManagedAgentsSessionRequiresAction) implManagedAgentsSessionThreadStatusIdleEventStopReasonUnion() {
}
func (ManagedAgentsSessionRetriesExhausted) implManagedAgentsSessionThreadStatusIdleEventStopReasonUnion() {
}
func (ManagedAgentsSessionBudgetReached) implManagedAgentsSessionThreadStatusIdleEventStopReasonUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsSessionThreadStatusIdleEventStopReasonUnion.AsAny().(type) {
//	case qoder.ManagedAgentsSessionEndTurn:
//	case qoder.ManagedAgentsSessionRequiresAction:
//	case qoder.ManagedAgentsSessionRetriesExhausted:
//	case qoder.ManagedAgentsSessionBudgetReached:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsSessionThreadStatusIdleEventStopReasonUnion) AsAny() anyManagedAgentsSessionThreadStatusIdleEventStopReason {
	switch u.Type {
	case "end_turn":
		return u.AsEndTurn()
	case "requires_action":
		return u.AsRequiresAction()
	case "retries_exhausted":
		return u.AsRetriesExhausted()
	case "budget_reached":
		return u.AsBudgetReached()
	}
	return nil
}

func (u ManagedAgentsSessionThreadStatusIdleEventStopReasonUnion) AsEndTurn() (v ManagedAgentsSessionEndTurn) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionThreadStatusIdleEventStopReasonUnion) AsRequiresAction() (v ManagedAgentsSessionRequiresAction) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionThreadStatusIdleEventStopReasonUnion) AsRetriesExhausted() (v ManagedAgentsSessionRetriesExhausted) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionThreadStatusIdleEventStopReasonUnion) AsBudgetReached() (v ManagedAgentsSessionBudgetReached) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsSessionThreadStatusIdleEventStopReasonUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *ManagedAgentsSessionThreadStatusIdleEventStopReasonUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSessionThreadStatusIdleEventType string

const (
	ManagedAgentsSessionThreadStatusIdleEventTypeSessionThreadStatusIdle ManagedAgentsSessionThreadStatusIdleEventType = "session.thread_status_idle"
)

// A session thread hit a transient error and is retrying automatically. Emitted on
// the thread's own stream and cross-posted to the primary stream for child
// threads.
type ManagedAgentsSessionThreadStatusRescheduledEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// Name of the agent the thread runs.
	AgentName string `json:"agent_name" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Public sthr\_ ID of the thread that is retrying.
	SessionThreadID string `json:"session_thread_id" api:"required"`
	// Any of "session.thread_status_rescheduled".
	Type ManagedAgentsSessionThreadStatusRescheduledEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID              respjson.Field
		AgentName       respjson.Field
		ProcessedAt     respjson.Field
		SessionThreadID respjson.Field
		Type            respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionThreadStatusRescheduledEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionThreadStatusRescheduledEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSessionThreadStatusRescheduledEventType string

const (
	ManagedAgentsSessionThreadStatusRescheduledEventTypeSessionThreadStatusRescheduled ManagedAgentsSessionThreadStatusRescheduledEventType = "session.thread_status_rescheduled"
)

// A session thread has begun executing. Emitted on the thread's own stream and
// cross-posted to the primary stream for child threads.
type ManagedAgentsSessionThreadStatusRunningEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// Name of the agent the thread runs.
	AgentName string `json:"agent_name" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Public sthr\_ ID of the thread that started running.
	SessionThreadID string `json:"session_thread_id" api:"required"`
	// Any of "session.thread_status_running".
	Type ManagedAgentsSessionThreadStatusRunningEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID              respjson.Field
		AgentName       respjson.Field
		ProcessedAt     respjson.Field
		SessionThreadID respjson.Field
		Type            respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionThreadStatusRunningEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionThreadStatusRunningEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSessionThreadStatusRunningEventType string

const (
	ManagedAgentsSessionThreadStatusRunningEventTypeSessionThreadStatusRunning ManagedAgentsSessionThreadStatusRunningEventType = "session.thread_status_running"
)

// A session thread has terminated and will accept no further input. Emitted on the
// thread's own stream and cross-posted to the primary stream for child threads.
type ManagedAgentsSessionThreadStatusTerminatedEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// Name of the agent the thread runs.
	AgentName string `json:"agent_name" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Public sthr\_ ID of the thread that terminated.
	SessionThreadID string `json:"session_thread_id" api:"required"`
	// Any of "session.thread_status_terminated".
	Type ManagedAgentsSessionThreadStatusTerminatedEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID              respjson.Field
		AgentName       respjson.Field
		ProcessedAt     respjson.Field
		SessionThreadID respjson.Field
		Type            respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionThreadStatusTerminatedEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionThreadStatusTerminatedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSessionThreadStatusTerminatedEventType string

const (
	ManagedAgentsSessionThreadStatusTerminatedEventTypeSessionThreadStatusTerminated ManagedAgentsSessionThreadStatusTerminatedEventType = "session.thread_status_terminated"
)

// Point-in-time snapshot of a session's cumulative usage.
type ManagedAgentsSessionUsageSnapshot struct {
	// Cumulative time in seconds during which the session had at least one thread in
	// running status. Overlapping activity from concurrent threads is counted once.
	// This is the duration the session's runtime cost is priced on.
	ActiveSeconds float64 `json:"active_seconds"`
	// Prompt-cache creation token usage broken down by cache lifetime.
	CacheCreation ManagedAgentsCacheCreationUsage `json:"cache_creation"`
	// Total tokens read from prompt cache.
	CacheReadInputTokens int64 `json:"cache_read_input_tokens"`
	// Total input tokens consumed across all turns.
	InputTokens int64 `json:"input_tokens"`
	// A monetary amount in a specific currency.
	ListCost MonetaryAmount `json:"list_cost"`
	// Total output tokens generated across all turns.
	OutputTokens int64 `json:"output_tokens"`
	// Cumulative count of server-executed tool invocations, broken down by tool.
	ServerToolUse ManagedAgentsServerToolUsage `json:"server_tool_use"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ActiveSeconds        respjson.Field
		CacheCreation        respjson.Field
		CacheReadInputTokens respjson.Field
		InputTokens          respjson.Field
		ListCost             respjson.Field
		OutputTokens         respjson.Field
		ServerToolUse        respjson.Field
		ExtraFields          map[string]respjson.Field
		raw                  string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionUsageSnapshot) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionUsageSnapshot) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when a model request completes.
type ManagedAgentsSpanModelRequestEndEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// Whether the model request resulted in an error.
	IsError bool `json:"is_error" api:"required"`
	// The id of the corresponding `span.model_request_start` event.
	ModelRequestStartID string `json:"model_request_start_id" api:"required"`
	// Token usage for a single model request.
	ModelUsage ManagedAgentsSpanModelUsage `json:"model_usage" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "span.model_request_end".
	Type ManagedAgentsSpanModelRequestEndEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                  respjson.Field
		IsError             respjson.Field
		ModelRequestStartID respjson.Field
		ModelUsage          respjson.Field
		ProcessedAt         respjson.Field
		Type                respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSpanModelRequestEndEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSpanModelRequestEndEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSpanModelRequestEndEventType string

const (
	ManagedAgentsSpanModelRequestEndEventTypeSpanModelRequestEnd ManagedAgentsSpanModelRequestEndEventType = "span.model_request_end"
)

// Emitted when a model request is initiated by the agent.
type ManagedAgentsSpanModelRequestStartEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "span.model_request_start".
	Type ManagedAgentsSpanModelRequestStartEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ProcessedAt respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSpanModelRequestStartEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSpanModelRequestStartEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSpanModelRequestStartEventType string

const (
	ManagedAgentsSpanModelRequestStartEventTypeSpanModelRequestStart ManagedAgentsSpanModelRequestStartEventType = "span.model_request_start"
)

// Token usage for a single model request.
type ManagedAgentsSpanModelUsage struct {
	// Tokens used to create prompt cache in this request.
	CacheCreationInputTokens int64 `json:"cache_creation_input_tokens" api:"required"`
	// Tokens read from prompt cache in this request.
	CacheReadInputTokens int64 `json:"cache_read_input_tokens" api:"required"`
	// Input tokens consumed by this request.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// Output tokens generated by this request.
	OutputTokens int64 `json:"output_tokens" api:"required"`
	// Inference speed mode. `fast` provides significantly faster output token
	// generation at premium pricing. Not all models support `fast`; invalid
	// combinations are rejected at create time.
	//
	// Any of "standard", "fast".
	Speed ManagedAgentsSpanModelUsageSpeed `json:"speed" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CacheCreationInputTokens respjson.Field
		CacheReadInputTokens     respjson.Field
		InputTokens              respjson.Field
		OutputTokens             respjson.Field
		Speed                    respjson.Field
		ExtraFields              map[string]respjson.Field
		raw                      string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSpanModelUsage) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSpanModelUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Inference speed mode. `fast` provides significantly faster output token
// generation at premium pricing. Not all models support `fast`; invalid
// combinations are rejected at create time.
type ManagedAgentsSpanModelUsageSpeed string

const (
	ManagedAgentsSpanModelUsageSpeedStandard ManagedAgentsSpanModelUsageSpeed = "standard"
	ManagedAgentsSpanModelUsageSpeedFast     ManagedAgentsSpanModelUsageSpeed = "fast"
)

// Emitted when an outcome evaluation cycle completes. Carries the verdict and
// aggregate token usage. A verdict of `needs_revision` means another evaluation
// cycle follows; `satisfied`, `max_iterations_reached`, `failed`, or `interrupted`
// are terminal — no further evaluation cycles follow.
type ManagedAgentsSpanOutcomeEvaluationEndEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// Human-readable explanation of the verdict. For `needs_revision`, describes which
	// criteria failed and why.
	Explanation string `json:"explanation" api:"required"`
	// 0-indexed revision cycle, matching the corresponding
	// `span.outcome_evaluation_start`.
	Iteration int64 `json:"iteration" api:"required"`
	// The id of the corresponding `span.outcome_evaluation_start` event.
	OutcomeEvaluationStartID string `json:"outcome_evaluation_start_id" api:"required"`
	// The `outc_` ID of the outcome being evaluated.
	OutcomeID string `json:"outcome_id" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Evaluation verdict. 'satisfied': criteria met, session goes idle.
	// 'needs_revision': criteria not met, another revision cycle follows.
	// 'max_iterations_reached': evaluation budget exhausted with criteria still unmet
	// — one final acknowledgment turn follows before the session goes idle, but no
	// further evaluation runs. 'failed': grader determined the rubric does not apply
	// to the deliverables. 'interrupted': user sent an interrupt while evaluation was
	// in progress.
	Result string `json:"result" api:"required"`
	// Any of "span.outcome_evaluation_end".
	Type ManagedAgentsSpanOutcomeEvaluationEndEventType `json:"type" api:"required"`
	// Token usage for a single model request.
	Usage ManagedAgentsSpanModelUsage `json:"usage" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                       respjson.Field
		Explanation              respjson.Field
		Iteration                respjson.Field
		OutcomeEvaluationStartID respjson.Field
		OutcomeID                respjson.Field
		ProcessedAt              respjson.Field
		Result                   respjson.Field
		Type                     respjson.Field
		Usage                    respjson.Field
		ExtraFields              map[string]respjson.Field
		raw                      string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSpanOutcomeEvaluationEndEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSpanOutcomeEvaluationEndEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSpanOutcomeEvaluationEndEventType string

const (
	ManagedAgentsSpanOutcomeEvaluationEndEventTypeSpanOutcomeEvaluationEnd ManagedAgentsSpanOutcomeEvaluationEndEventType = "span.outcome_evaluation_end"
)

// Periodic heartbeat emitted while an outcome evaluation cycle is in progress.
// Distinguishes 'evaluation is actively running' from 'evaluation is stuck'
// between the corresponding `span.outcome_evaluation_start` and
// `span.outcome_evaluation_end` events.
type ManagedAgentsSpanOutcomeEvaluationOngoingEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// 0-indexed revision cycle, matching the corresponding
	// `span.outcome_evaluation_start`.
	Iteration int64 `json:"iteration" api:"required"`
	// The `outc_` ID of the outcome being evaluated.
	OutcomeID string `json:"outcome_id" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "span.outcome_evaluation_ongoing".
	Type ManagedAgentsSpanOutcomeEvaluationOngoingEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Iteration   respjson.Field
		OutcomeID   respjson.Field
		ProcessedAt respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSpanOutcomeEvaluationOngoingEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSpanOutcomeEvaluationOngoingEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSpanOutcomeEvaluationOngoingEventType string

const (
	ManagedAgentsSpanOutcomeEvaluationOngoingEventTypeSpanOutcomeEvaluationOngoing ManagedAgentsSpanOutcomeEvaluationOngoingEventType = "span.outcome_evaluation_ongoing"
)

// Emitted when an outcome evaluation cycle begins.
type ManagedAgentsSpanOutcomeEvaluationStartEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// 0-indexed revision cycle. 0 is the first evaluation; 1 is the re-evaluation
	// after the first revision; etc.
	Iteration int64 `json:"iteration" api:"required"`
	// The `outc_` ID of the outcome being evaluated.
	OutcomeID string `json:"outcome_id" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "span.outcome_evaluation_start".
	Type ManagedAgentsSpanOutcomeEvaluationStartEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Iteration   respjson.Field
		OutcomeID   respjson.Field
		ProcessedAt respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSpanOutcomeEvaluationStartEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSpanOutcomeEvaluationStartEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSpanOutcomeEvaluationStartEventType string

const (
	ManagedAgentsSpanOutcomeEvaluationStartEventTypeSpanOutcomeEvaluationStart ManagedAgentsSpanOutcomeEvaluationStartEventType = "span.outcome_evaluation_start"
)

// ManagedAgentsStreamSessionEventsUnion contains all possible properties and
// values from [ManagedAgentsUserMessageEvent],
// [ManagedAgentsUserInterruptEvent],
// [ManagedAgentsUserToolConfirmationEvent],
// [ManagedAgentsUserCustomToolResultEvent],
// [ManagedAgentsAgentCustomToolUseEvent],
// [ManagedAgentsAgentMessageEvent], [ManagedAgentsAgentThinkingEvent],
// [ManagedAgentsAgentMCPToolUseEvent],
// [ManagedAgentsAgentMCPToolResultEvent],
// [ManagedAgentsAgentToolUseEvent], [ManagedAgentsAgentToolResultEvent],
// [ManagedAgentsAgentThreadMessageReceivedEvent],
// [ManagedAgentsAgentThreadMessageSentEvent],
// [ManagedAgentsAgentThreadContextCompactedEvent],
// [ManagedAgentsSessionErrorEvent],
// [ManagedAgentsSessionStatusRescheduledEvent],
// [ManagedAgentsSessionStatusRunningEvent],
// [ManagedAgentsSessionStatusIdleEvent],
// [ManagedAgentsSessionStatusTerminatedEvent],
// [ManagedAgentsSessionThreadCreatedEvent],
// [ManagedAgentsSpanOutcomeEvaluationStartEvent],
// [ManagedAgentsSpanOutcomeEvaluationEndEvent],
// [ManagedAgentsSpanModelRequestStartEvent],
// [ManagedAgentsSpanModelRequestEndEvent],
// [ManagedAgentsSpanOutcomeEvaluationOngoingEvent],
// [ManagedAgentsUserDefineOutcomeEvent],
// [ManagedAgentsSessionDeletedEvent],
// [ManagedAgentsSessionThreadStatusRunningEvent],
// [ManagedAgentsSessionThreadStatusIdleEvent],
// [ManagedAgentsSessionThreadStatusTerminatedEvent],
// [ManagedAgentsUserToolResultEvent],
// [ManagedAgentsSessionThreadStatusRescheduledEvent],
// [ManagedAgentsSessionUpdatedEvent], [ManagedAgentsStartEvent],
// [ManagedAgentsDeltaEvent], [ManagedAgentsSystemMessageEvent],
// [ManagedAgentsSessionUsageEvent].
//
// Use the [ManagedAgentsStreamSessionEventsUnion.AsAny] method to switch on
// the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsStreamSessionEventsUnion struct {
	ID string `json:"id"`
	// This field is a union of [[]ManagedAgentsUserMessageEventContentUnion],
	// [[]ManagedAgentsUserCustomToolResultEventContentUnion],
	// [[]ManagedAgentsAgentMessageEventContentUnion],
	// [[]ManagedAgentsAgentMCPToolResultEventContentUnion],
	// [[]ManagedAgentsAgentToolResultEventContentUnion],
	// [[]ManagedAgentsAgentThreadMessageReceivedEventContentUnion],
	// [[]ManagedAgentsAgentThreadMessageSentEventContentUnion],
	// [[]ManagedAgentsUserToolResultEventContentUnion],
	// [[]ManagedAgentsSystemContentBlock]
	Content ManagedAgentsStreamSessionEventsUnionContent `json:"content"`
	// Any of "user.message", "user.interrupt", "user.tool_confirmation",
	// "user.custom_tool_result", "agent.custom_tool_use", "agent.message",
	// "agent.thinking", "agent.mcp_tool_use", "agent.mcp_tool_result",
	// "agent.tool_use", "agent.tool_result", "agent.thread_message_received",
	// "agent.thread_message_sent", "agent.thread_context_compacted", "session.error",
	// "session.status_rescheduled", "session.status_running", "session.status_idle",
	// "session.status_terminated", "session.thread_created",
	// "span.outcome_evaluation_start", "span.outcome_evaluation_end",
	// "span.model_request_start", "span.model_request_end",
	// "span.outcome_evaluation_ongoing", "user.define_outcome", "session.deleted",
	// "session.thread_status_running", "session.thread_status_idle",
	// "session.thread_status_terminated", "user.tool_result",
	// "session.thread_status_rescheduled", "session.updated", "event_start",
	// "event_delta", "system.message", "session.usage".
	Type            string    `json:"type"`
	ProcessedAt     time.Time `json:"processed_at"`
	SessionThreadID string    `json:"session_thread_id"`
	Result          string    `json:"result"`
	ToolUseID       string    `json:"tool_use_id"`
	// This field is from variant [ManagedAgentsUserToolConfirmationEvent].
	DenyMessage string `json:"deny_message"`
	// This field is from variant [ManagedAgentsUserCustomToolResultEvent].
	CustomToolUseID string `json:"custom_tool_use_id"`
	IsError         bool   `json:"is_error"`
	Input           any    `json:"input"`
	Name            string `json:"name"`
	// This field is from variant [ManagedAgentsAgentMCPToolUseEvent].
	MCPServerName       string `json:"mcp_server_name"`
	EvaluatedPermission string `json:"evaluated_permission"`
	// This field is from variant [ManagedAgentsAgentMCPToolResultEvent].
	MCPToolUseID string `json:"mcp_tool_use_id"`
	// This field is from variant [ManagedAgentsAgentThreadMessageReceivedEvent].
	FromSessionThreadID string `json:"from_session_thread_id"`
	// This field is from variant [ManagedAgentsAgentThreadMessageReceivedEvent].
	FromAgentName string `json:"from_agent_name"`
	// This field is from variant [ManagedAgentsAgentThreadMessageSentEvent].
	ToSessionThreadID string `json:"to_session_thread_id"`
	// This field is from variant [ManagedAgentsAgentThreadMessageSentEvent].
	ToAgentName string `json:"to_agent_name"`
	// This field is from variant [ManagedAgentsSessionErrorEvent].
	Error ManagedAgentsSessionErrorEventErrorUnion `json:"error"`
	// This field is a union of
	// [ManagedAgentsSessionStatusIdleEventStopReasonUnion],
	// [ManagedAgentsSessionThreadStatusIdleEventStopReasonUnion]
	StopReason ManagedAgentsStreamSessionEventsUnionStopReason `json:"stop_reason"`
	AgentName  string                                          `json:"agent_name"`
	Iteration  int64                                           `json:"iteration"`
	OutcomeID  string                                          `json:"outcome_id"`
	// This field is from variant [ManagedAgentsSpanOutcomeEvaluationEndEvent].
	Explanation string `json:"explanation"`
	// This field is from variant [ManagedAgentsSpanOutcomeEvaluationEndEvent].
	OutcomeEvaluationStartID string `json:"outcome_evaluation_start_id"`
	// This field is a union of [ManagedAgentsSpanModelUsage],
	// [ManagedAgentsSessionUsageSnapshot]
	Usage ManagedAgentsStreamSessionEventsUnionUsage `json:"usage"`
	// This field is from variant [ManagedAgentsSpanModelRequestEndEvent].
	ModelRequestStartID string `json:"model_request_start_id"`
	// This field is from variant [ManagedAgentsSpanModelRequestEndEvent].
	ModelUsage ManagedAgentsSpanModelUsage `json:"model_usage"`
	// This field is from variant [ManagedAgentsUserDefineOutcomeEvent].
	Description string `json:"description"`
	// This field is from variant [ManagedAgentsUserDefineOutcomeEvent].
	MaxIterations int64 `json:"max_iterations"`
	// This field is from variant [ManagedAgentsUserDefineOutcomeEvent].
	Rubric ManagedAgentsUserDefineOutcomeEventRubricUnion `json:"rubric"`
	// This field is from variant [ManagedAgentsSessionUpdatedEvent].
	Agent ManagedAgentsSessionAgent `json:"agent"`
	// This field is from variant [ManagedAgentsSessionUpdatedEvent].
	Budget ManagedAgentsBudgetLimit `json:"budget"`
	// This field is from variant [ManagedAgentsSessionUpdatedEvent].
	Metadata map[string]string `json:"metadata"`
	// This field is from variant [ManagedAgentsSessionUpdatedEvent].
	Title string `json:"title"`
	// This field is from variant [ManagedAgentsStartEvent].
	Event ManagedAgentsStartEventPreviewUnion `json:"event"`
	// This field is from variant [ManagedAgentsDeltaEvent].
	Delta ManagedAgentsDeltaContent `json:"delta"`
	// This field is from variant [ManagedAgentsDeltaEvent].
	EventID string `json:"event_id"`
	JSON    struct {
		ID                       respjson.Field
		Content                  respjson.Field
		Type                     respjson.Field
		ProcessedAt              respjson.Field
		SessionThreadID          respjson.Field
		Result                   respjson.Field
		ToolUseID                respjson.Field
		DenyMessage              respjson.Field
		CustomToolUseID          respjson.Field
		IsError                  respjson.Field
		Input                    respjson.Field
		Name                     respjson.Field
		MCPServerName            respjson.Field
		EvaluatedPermission      respjson.Field
		MCPToolUseID             respjson.Field
		FromSessionThreadID      respjson.Field
		FromAgentName            respjson.Field
		ToSessionThreadID        respjson.Field
		ToAgentName              respjson.Field
		Error                    respjson.Field
		StopReason               respjson.Field
		AgentName                respjson.Field
		Iteration                respjson.Field
		OutcomeID                respjson.Field
		Explanation              respjson.Field
		OutcomeEvaluationStartID respjson.Field
		Usage                    respjson.Field
		ModelRequestStartID      respjson.Field
		ModelUsage               respjson.Field
		Description              respjson.Field
		MaxIterations            respjson.Field
		Rubric                   respjson.Field
		Agent                    respjson.Field
		Budget                   respjson.Field
		Metadata                 respjson.Field
		Title                    respjson.Field
		Event                    respjson.Field
		Delta                    respjson.Field
		EventID                  respjson.Field
		raw                      string
	} `json:"-"`
}

// anyManagedAgentsStreamSessionEvents is implemented by each variant of
// [ManagedAgentsStreamSessionEventsUnion] to add type safety for the return
// type of [ManagedAgentsStreamSessionEventsUnion.AsAny]
type anyManagedAgentsStreamSessionEvents interface {
	implManagedAgentsStreamSessionEventsUnion()
}

func (ManagedAgentsUserMessageEvent) implManagedAgentsStreamSessionEventsUnion()          {}
func (ManagedAgentsUserInterruptEvent) implManagedAgentsStreamSessionEventsUnion()        {}
func (ManagedAgentsUserToolConfirmationEvent) implManagedAgentsStreamSessionEventsUnion() {}
func (ManagedAgentsUserCustomToolResultEvent) implManagedAgentsStreamSessionEventsUnion() {}
func (ManagedAgentsAgentCustomToolUseEvent) implManagedAgentsStreamSessionEventsUnion()   {}
func (ManagedAgentsAgentMessageEvent) implManagedAgentsStreamSessionEventsUnion()         {}
func (ManagedAgentsAgentThinkingEvent) implManagedAgentsStreamSessionEventsUnion()        {}
func (ManagedAgentsAgentMCPToolUseEvent) implManagedAgentsStreamSessionEventsUnion()      {}
func (ManagedAgentsAgentMCPToolResultEvent) implManagedAgentsStreamSessionEventsUnion()   {}
func (ManagedAgentsAgentToolUseEvent) implManagedAgentsStreamSessionEventsUnion()         {}
func (ManagedAgentsAgentToolResultEvent) implManagedAgentsStreamSessionEventsUnion()      {}
func (ManagedAgentsAgentThreadMessageReceivedEvent) implManagedAgentsStreamSessionEventsUnion() {
}
func (ManagedAgentsAgentThreadMessageSentEvent) implManagedAgentsStreamSessionEventsUnion() {}
func (ManagedAgentsAgentThreadContextCompactedEvent) implManagedAgentsStreamSessionEventsUnion() {
}
func (ManagedAgentsSessionErrorEvent) implManagedAgentsStreamSessionEventsUnion() {}
func (ManagedAgentsSessionStatusRescheduledEvent) implManagedAgentsStreamSessionEventsUnion() {
}
func (ManagedAgentsSessionStatusRunningEvent) implManagedAgentsStreamSessionEventsUnion() {}
func (ManagedAgentsSessionStatusIdleEvent) implManagedAgentsStreamSessionEventsUnion()    {}
func (ManagedAgentsSessionStatusTerminatedEvent) implManagedAgentsStreamSessionEventsUnion() {
}
func (ManagedAgentsSessionThreadCreatedEvent) implManagedAgentsStreamSessionEventsUnion() {}
func (ManagedAgentsSpanOutcomeEvaluationStartEvent) implManagedAgentsStreamSessionEventsUnion() {
}
func (ManagedAgentsSpanOutcomeEvaluationEndEvent) implManagedAgentsStreamSessionEventsUnion() {
}
func (ManagedAgentsSpanModelRequestStartEvent) implManagedAgentsStreamSessionEventsUnion() {}
func (ManagedAgentsSpanModelRequestEndEvent) implManagedAgentsStreamSessionEventsUnion()   {}
func (ManagedAgentsSpanOutcomeEvaluationOngoingEvent) implManagedAgentsStreamSessionEventsUnion() {
}
func (ManagedAgentsUserDefineOutcomeEvent) implManagedAgentsStreamSessionEventsUnion() {}
func (ManagedAgentsSessionDeletedEvent) implManagedAgentsStreamSessionEventsUnion()    {}
func (ManagedAgentsSessionThreadStatusRunningEvent) implManagedAgentsStreamSessionEventsUnion() {
}
func (ManagedAgentsSessionThreadStatusIdleEvent) implManagedAgentsStreamSessionEventsUnion() {
}
func (ManagedAgentsSessionThreadStatusTerminatedEvent) implManagedAgentsStreamSessionEventsUnion() {
}
func (ManagedAgentsUserToolResultEvent) implManagedAgentsStreamSessionEventsUnion() {}
func (ManagedAgentsSessionThreadStatusRescheduledEvent) implManagedAgentsStreamSessionEventsUnion() {
}
func (ManagedAgentsSessionUpdatedEvent) implManagedAgentsStreamSessionEventsUnion() {}
func (ManagedAgentsStartEvent) implManagedAgentsStreamSessionEventsUnion()          {}
func (ManagedAgentsDeltaEvent) implManagedAgentsStreamSessionEventsUnion()          {}
func (ManagedAgentsSystemMessageEvent) implManagedAgentsStreamSessionEventsUnion()  {}
func (ManagedAgentsSessionUsageEvent) implManagedAgentsStreamSessionEventsUnion()   {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsStreamSessionEventsUnion.AsAny().(type) {
//	case qoder.ManagedAgentsUserMessageEvent:
//	case qoder.ManagedAgentsUserInterruptEvent:
//	case qoder.ManagedAgentsUserToolConfirmationEvent:
//	case qoder.ManagedAgentsUserCustomToolResultEvent:
//	case qoder.ManagedAgentsAgentCustomToolUseEvent:
//	case qoder.ManagedAgentsAgentMessageEvent:
//	case qoder.ManagedAgentsAgentThinkingEvent:
//	case qoder.ManagedAgentsAgentMCPToolUseEvent:
//	case qoder.ManagedAgentsAgentMCPToolResultEvent:
//	case qoder.ManagedAgentsAgentToolUseEvent:
//	case qoder.ManagedAgentsAgentToolResultEvent:
//	case qoder.ManagedAgentsAgentThreadMessageReceivedEvent:
//	case qoder.ManagedAgentsAgentThreadMessageSentEvent:
//	case qoder.ManagedAgentsAgentThreadContextCompactedEvent:
//	case qoder.ManagedAgentsSessionErrorEvent:
//	case qoder.ManagedAgentsSessionStatusRescheduledEvent:
//	case qoder.ManagedAgentsSessionStatusRunningEvent:
//	case qoder.ManagedAgentsSessionStatusIdleEvent:
//	case qoder.ManagedAgentsSessionStatusTerminatedEvent:
//	case qoder.ManagedAgentsSessionThreadCreatedEvent:
//	case qoder.ManagedAgentsSpanOutcomeEvaluationStartEvent:
//	case qoder.ManagedAgentsSpanOutcomeEvaluationEndEvent:
//	case qoder.ManagedAgentsSpanModelRequestStartEvent:
//	case qoder.ManagedAgentsSpanModelRequestEndEvent:
//	case qoder.ManagedAgentsSpanOutcomeEvaluationOngoingEvent:
//	case qoder.ManagedAgentsUserDefineOutcomeEvent:
//	case qoder.ManagedAgentsSessionDeletedEvent:
//	case qoder.ManagedAgentsSessionThreadStatusRunningEvent:
//	case qoder.ManagedAgentsSessionThreadStatusIdleEvent:
//	case qoder.ManagedAgentsSessionThreadStatusTerminatedEvent:
//	case qoder.ManagedAgentsUserToolResultEvent:
//	case qoder.ManagedAgentsSessionThreadStatusRescheduledEvent:
//	case qoder.ManagedAgentsSessionUpdatedEvent:
//	case qoder.ManagedAgentsStartEvent:
//	case qoder.ManagedAgentsDeltaEvent:
//	case qoder.ManagedAgentsSystemMessageEvent:
//	case qoder.ManagedAgentsSessionUsageEvent:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsStreamSessionEventsUnion) AsAny() anyManagedAgentsStreamSessionEvents {
	switch u.Type {
	case "user.message":
		return u.AsUserMessage()
	case "user.interrupt":
		return u.AsUserInterrupt()
	case "user.tool_confirmation":
		return u.AsUserToolConfirmation()
	case "user.custom_tool_result":
		return u.AsUserCustomToolResult()
	case "agent.custom_tool_use":
		return u.AsAgentCustomToolUse()
	case "agent.message":
		return u.AsAgentMessage()
	case "agent.thinking":
		return u.AsAgentThinking()
	case "agent.mcp_tool_use":
		return u.AsAgentMCPToolUse()
	case "agent.mcp_tool_result":
		return u.AsAgentMCPToolResult()
	case "agent.tool_use":
		return u.AsAgentToolUse()
	case "agent.tool_result":
		return u.AsAgentToolResult()
	case "agent.thread_message_received":
		return u.AsAgentThreadMessageReceived()
	case "agent.thread_message_sent":
		return u.AsAgentThreadMessageSent()
	case "agent.thread_context_compacted":
		return u.AsAgentThreadContextCompacted()
	case "session.error":
		return u.AsSessionError()
	case "session.status_rescheduled":
		return u.AsSessionStatusRescheduled()
	case "session.status_running":
		return u.AsSessionStatusRunning()
	case "session.status_idle":
		return u.AsSessionStatusIdle()
	case "session.status_terminated":
		return u.AsSessionStatusTerminated()
	case "session.thread_created":
		return u.AsSessionThreadCreated()
	case "span.outcome_evaluation_start":
		return u.AsSpanOutcomeEvaluationStart()
	case "span.outcome_evaluation_end":
		return u.AsSpanOutcomeEvaluationEnd()
	case "span.model_request_start":
		return u.AsSpanModelRequestStart()
	case "span.model_request_end":
		return u.AsSpanModelRequestEnd()
	case "span.outcome_evaluation_ongoing":
		return u.AsSpanOutcomeEvaluationOngoing()
	case "user.define_outcome":
		return u.AsUserDefineOutcome()
	case "session.deleted":
		return u.AsSessionDeleted()
	case "session.thread_status_running":
		return u.AsSessionThreadStatusRunning()
	case "session.thread_status_idle":
		return u.AsSessionThreadStatusIdle()
	case "session.thread_status_terminated":
		return u.AsSessionThreadStatusTerminated()
	case "user.tool_result":
		return u.AsUserToolResult()
	case "session.thread_status_rescheduled":
		return u.AsSessionThreadStatusRescheduled()
	case "session.updated":
		return u.AsSessionUpdated()
	case "event_start":
		return u.AsEventStart()
	case "event_delta":
		return u.AsEventDelta()
	case "system.message":
		return u.AsSystemMessage()
	case "session.usage":
		return u.AsSessionUsage()
	}
	return nil
}

func (u ManagedAgentsStreamSessionEventsUnion) AsUserMessage() (v ManagedAgentsUserMessageEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsUserInterrupt() (v ManagedAgentsUserInterruptEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsUserToolConfirmation() (v ManagedAgentsUserToolConfirmationEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsUserCustomToolResult() (v ManagedAgentsUserCustomToolResultEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsAgentCustomToolUse() (v ManagedAgentsAgentCustomToolUseEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsAgentMessage() (v ManagedAgentsAgentMessageEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsAgentThinking() (v ManagedAgentsAgentThinkingEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsAgentMCPToolUse() (v ManagedAgentsAgentMCPToolUseEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsAgentMCPToolResult() (v ManagedAgentsAgentMCPToolResultEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsAgentToolUse() (v ManagedAgentsAgentToolUseEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsAgentToolResult() (v ManagedAgentsAgentToolResultEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsAgentThreadMessageReceived() (v ManagedAgentsAgentThreadMessageReceivedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsAgentThreadMessageSent() (v ManagedAgentsAgentThreadMessageSentEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsAgentThreadContextCompacted() (v ManagedAgentsAgentThreadContextCompactedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsSessionError() (v ManagedAgentsSessionErrorEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsSessionStatusRescheduled() (v ManagedAgentsSessionStatusRescheduledEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsSessionStatusRunning() (v ManagedAgentsSessionStatusRunningEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsSessionStatusIdle() (v ManagedAgentsSessionStatusIdleEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsSessionStatusTerminated() (v ManagedAgentsSessionStatusTerminatedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsSessionThreadCreated() (v ManagedAgentsSessionThreadCreatedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsSpanOutcomeEvaluationStart() (v ManagedAgentsSpanOutcomeEvaluationStartEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsSpanOutcomeEvaluationEnd() (v ManagedAgentsSpanOutcomeEvaluationEndEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsSpanModelRequestStart() (v ManagedAgentsSpanModelRequestStartEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsSpanModelRequestEnd() (v ManagedAgentsSpanModelRequestEndEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsSpanOutcomeEvaluationOngoing() (v ManagedAgentsSpanOutcomeEvaluationOngoingEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsUserDefineOutcome() (v ManagedAgentsUserDefineOutcomeEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsSessionDeleted() (v ManagedAgentsSessionDeletedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsSessionThreadStatusRunning() (v ManagedAgentsSessionThreadStatusRunningEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsSessionThreadStatusIdle() (v ManagedAgentsSessionThreadStatusIdleEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsSessionThreadStatusTerminated() (v ManagedAgentsSessionThreadStatusTerminatedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsUserToolResult() (v ManagedAgentsUserToolResultEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsSessionThreadStatusRescheduled() (v ManagedAgentsSessionThreadStatusRescheduledEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsSessionUpdated() (v ManagedAgentsSessionUpdatedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsEventStart() (v ManagedAgentsStartEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsEventDelta() (v ManagedAgentsDeltaEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsSystemMessage() (v ManagedAgentsSystemMessageEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionEventsUnion) AsSessionUsage() (v ManagedAgentsSessionUsageEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsStreamSessionEventsUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsStreamSessionEventsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsStreamSessionEventsUnionContent is an implicit subunion of
// [ManagedAgentsStreamSessionEventsUnion].
// ManagedAgentsStreamSessionEventsUnionContent provides convenient access to
// the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsStreamSessionEventsUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfManagedAgentsUserMessageEventContentArray
// OfManagedAgentsUserCustomToolResultEventContentArray
// OfManagedAgentsAgentMessageEventContentArray
// OfManagedAgentsAgentMCPToolResultEventContentArray
// OfManagedAgentsAgentToolResultEventContentArray
// OfManagedAgentsAgentThreadMessageReceivedEventContentArray
// OfManagedAgentsAgentThreadMessageSentEventContentArray
// OfManagedAgentsUserToolResultEventContentArray
// OfManagedAgentsSystemContentBlockArray]
type ManagedAgentsStreamSessionEventsUnionContent struct {
	// This field will be present if the value is a
	// [[]ManagedAgentsUserMessageEventContentUnion] instead of an object.
	OfManagedAgentsUserMessageEventContentArray []ManagedAgentsUserMessageEventContentUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]ManagedAgentsUserCustomToolResultEventContentUnion] instead of an object.
	OfManagedAgentsUserCustomToolResultEventContentArray []ManagedAgentsUserCustomToolResultEventContentUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]ManagedAgentsAgentMessageEventContentUnion] instead of an object.
	OfManagedAgentsAgentMessageEventContentArray []ManagedAgentsAgentMessageEventContentUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]ManagedAgentsAgentMCPToolResultEventContentUnion] instead of an object.
	OfManagedAgentsAgentMCPToolResultEventContentArray []ManagedAgentsAgentMCPToolResultEventContentUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]ManagedAgentsAgentToolResultEventContentUnion] instead of an object.
	OfManagedAgentsAgentToolResultEventContentArray []ManagedAgentsAgentToolResultEventContentUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]ManagedAgentsAgentThreadMessageReceivedEventContentUnion] instead of an
	// object.
	OfManagedAgentsAgentThreadMessageReceivedEventContentArray []ManagedAgentsAgentThreadMessageReceivedEventContentUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]ManagedAgentsAgentThreadMessageSentEventContentUnion] instead of an
	// object.
	OfManagedAgentsAgentThreadMessageSentEventContentArray []ManagedAgentsAgentThreadMessageSentEventContentUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]ManagedAgentsUserToolResultEventContentUnion] instead of an object.
	OfManagedAgentsUserToolResultEventContentArray []ManagedAgentsUserToolResultEventContentUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]ManagedAgentsSystemContentBlock] instead of an object.
	OfManagedAgentsSystemContentBlockArray []ManagedAgentsSystemContentBlock `json:",inline"`
	JSON                                   struct {
		OfManagedAgentsUserMessageEventContentArray                respjson.Field
		OfManagedAgentsUserCustomToolResultEventContentArray       respjson.Field
		OfManagedAgentsAgentMessageEventContentArray               respjson.Field
		OfManagedAgentsAgentMCPToolResultEventContentArray         respjson.Field
		OfManagedAgentsAgentToolResultEventContentArray            respjson.Field
		OfManagedAgentsAgentThreadMessageReceivedEventContentArray respjson.Field
		OfManagedAgentsAgentThreadMessageSentEventContentArray     respjson.Field
		OfManagedAgentsUserToolResultEventContentArray             respjson.Field
		OfManagedAgentsSystemContentBlockArray                     respjson.Field
		raw                                                        string
	} `json:"-"`
}

func (r *ManagedAgentsStreamSessionEventsUnionContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsStreamSessionEventsUnionStopReason is an implicit subunion of
// [ManagedAgentsStreamSessionEventsUnion].
// ManagedAgentsStreamSessionEventsUnionStopReason provides convenient access
// to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsStreamSessionEventsUnion].
type ManagedAgentsStreamSessionEventsUnionStopReason struct {
	Type string `json:"type"`
	// This field is from variant
	// [ManagedAgentsSessionStatusIdleEventStopReasonUnion],
	// [ManagedAgentsSessionThreadStatusIdleEventStopReasonUnion].
	EventIDs []string `json:"event_ids"`
	JSON     struct {
		Type     respjson.Field
		EventIDs respjson.Field
		raw      string
	} `json:"-"`
}

func (r *ManagedAgentsStreamSessionEventsUnionStopReason) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsStreamSessionEventsUnionUsage is an implicit subunion of
// [ManagedAgentsStreamSessionEventsUnion].
// ManagedAgentsStreamSessionEventsUnionUsage provides convenient access to the
// sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsStreamSessionEventsUnion].
type ManagedAgentsStreamSessionEventsUnionUsage struct {
	// This field is from variant [ManagedAgentsSpanModelUsage].
	CacheCreationInputTokens int64 `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int64 `json:"cache_read_input_tokens"`
	InputTokens              int64 `json:"input_tokens"`
	OutputTokens             int64 `json:"output_tokens"`
	// This field is from variant [ManagedAgentsSpanModelUsage].
	Speed ManagedAgentsSpanModelUsageSpeed `json:"speed"`
	// This field is from variant [ManagedAgentsSessionUsageSnapshot].
	ActiveSeconds float64 `json:"active_seconds"`
	// This field is from variant [ManagedAgentsSessionUsageSnapshot].
	CacheCreation ManagedAgentsCacheCreationUsage `json:"cache_creation"`
	// This field is from variant [ManagedAgentsSessionUsageSnapshot].
	ListCost MonetaryAmount `json:"list_cost"`
	// This field is from variant [ManagedAgentsSessionUsageSnapshot].
	ServerToolUse ManagedAgentsServerToolUsage `json:"server_tool_use"`
	JSON          struct {
		CacheCreationInputTokens respjson.Field
		CacheReadInputTokens     respjson.Field
		InputTokens              respjson.Field
		OutputTokens             respjson.Field
		Speed                    respjson.Field
		ActiveSeconds            respjson.Field
		CacheCreation            respjson.Field
		ListCost                 respjson.Field
		ServerToolUse            respjson.Field
		raw                      string
	} `json:"-"`
}

func (r *ManagedAgentsStreamSessionEventsUnionUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Privileged context for the accompanying turn and all subsequent turns, appended
// to the session's system context as a `role: "system"` turn rather than replacing
// the top-level system prompt. At most one per request: it must be the final event
// and immediately follow the `user.message`, `user.tool_result`, or
// `user.custom_tool_result` it accompanies. Only supported on models that accept
// mid-conversation system messages.
//
// The properties Content, Type are required.
type ManagedAgentsSystemMessageEventParams struct {
	// System content blocks to append. Text-only.
	Content []ManagedAgentsSystemContentBlockParam `json:"content,omitzero" api:"required"`
	// Any of "system.message".
	Type ManagedAgentsSystemMessageEventParamsType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsSystemMessageEventParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsSystemMessageEventParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsSystemMessageEventParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSystemMessageEventParamsType string

const (
	ManagedAgentsSystemMessageEventParamsTypeSystemMessage ManagedAgentsSystemMessageEventParamsType = "system.message"
)

// Regular text content.
type ManagedAgentsTextBlock struct {
	// The text content.
	Text string `json:"text" api:"required"`
	// Any of "text".
	Type ManagedAgentsTextBlockType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Text        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsTextBlock) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsTextBlock) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ManagedAgentsTextBlock to a
// ManagedAgentsTextBlockParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ManagedAgentsTextBlockParam.Overrides()
func (r ManagedAgentsTextBlock) ToParam() ManagedAgentsTextBlockParam {
	return param.Override[ManagedAgentsTextBlockParam](json.RawMessage(r.RawJSON()))
}

type ManagedAgentsTextBlockType string

const (
	ManagedAgentsTextBlockTypeText ManagedAgentsTextBlockType = "text"
)

// Regular text content.
//
// The properties Text, Type are required.
type ManagedAgentsTextBlockParam struct {
	// The text content.
	Text string `json:"text" api:"required"`
	// Any of "text".
	Type ManagedAgentsTextBlockType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsTextBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsTextBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsTextBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Rubric content provided inline as text.
type ManagedAgentsTextRubric struct {
	// Rubric content. Plain text or markdown — the grader treats it as freeform text.
	Content string `json:"content" api:"required"`
	// Any of "text".
	Type ManagedAgentsTextRubricType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Content     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsTextRubric) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsTextRubric) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsTextRubricType string

const (
	ManagedAgentsTextRubricTypeText ManagedAgentsTextRubricType = "text"
)

// Rubric content provided inline as text.
//
// The properties Content, Type are required.
type ManagedAgentsTextRubricParams struct {
	// Rubric content. Plain text or markdown — the grader treats it as freeform text.
	// Maximum 262144 characters.
	Content string `json:"content" api:"required"`
	// Any of "text".
	Type ManagedAgentsTextRubricParamsType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsTextRubricParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsTextRubricParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsTextRubricParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsTextRubricParamsType string

const (
	ManagedAgentsTextRubricParamsTypeText ManagedAgentsTextRubricParamsType = "text"
)

// An unknown or unexpected error occurred during session execution. A fallback
// variant; clients that don't recognize a new error code can match on
// `retry_status` and `message` alone.
type ManagedAgentsUnknownError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// What the client should do next in response to this error.
	RetryStatus ManagedAgentsUnknownErrorRetryStatusUnion `json:"retry_status" api:"required"`
	// Any of "unknown_error".
	Type ManagedAgentsUnknownErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		RetryStatus respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsUnknownError) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsUnknownError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsUnknownErrorRetryStatusUnion contains all possible properties
// and values from [ManagedAgentsRetryStatusRetrying],
// [ManagedAgentsRetryStatusExhausted], [ManagedAgentsRetryStatusTerminal].
//
// Use the [ManagedAgentsUnknownErrorRetryStatusUnion.AsAny] method to switch
// on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsUnknownErrorRetryStatusUnion struct {
	// Any of "retrying", "exhausted", "terminal".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyManagedAgentsUnknownErrorRetryStatus is implemented by each variant of
// [ManagedAgentsUnknownErrorRetryStatusUnion] to add type safety for the
// return type of [ManagedAgentsUnknownErrorRetryStatusUnion.AsAny]
type anyManagedAgentsUnknownErrorRetryStatus interface {
	implManagedAgentsUnknownErrorRetryStatusUnion()
}

func (ManagedAgentsRetryStatusRetrying) implManagedAgentsUnknownErrorRetryStatusUnion()  {}
func (ManagedAgentsRetryStatusExhausted) implManagedAgentsUnknownErrorRetryStatusUnion() {}
func (ManagedAgentsRetryStatusTerminal) implManagedAgentsUnknownErrorRetryStatusUnion()  {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsUnknownErrorRetryStatusUnion.AsAny().(type) {
//	case qoder.ManagedAgentsRetryStatusRetrying:
//	case qoder.ManagedAgentsRetryStatusExhausted:
//	case qoder.ManagedAgentsRetryStatusTerminal:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsUnknownErrorRetryStatusUnion) AsAny() anyManagedAgentsUnknownErrorRetryStatus {
	switch u.Type {
	case "retrying":
		return u.AsRetrying()
	case "exhausted":
		return u.AsExhausted()
	case "terminal":
		return u.AsTerminal()
	}
	return nil
}

func (u ManagedAgentsUnknownErrorRetryStatusUnion) AsRetrying() (v ManagedAgentsRetryStatusRetrying) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsUnknownErrorRetryStatusUnion) AsExhausted() (v ManagedAgentsRetryStatusExhausted) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsUnknownErrorRetryStatusUnion) AsTerminal() (v ManagedAgentsRetryStatusTerminal) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsUnknownErrorRetryStatusUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsUnknownErrorRetryStatusUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsUnknownErrorType string

const (
	ManagedAgentsUnknownErrorTypeUnknownError ManagedAgentsUnknownErrorType = "unknown_error"
)

// Document referenced by URL.
type ManagedAgentsURLDocumentSource struct {
	// Any of "url".
	Type ManagedAgentsURLDocumentSourceType `json:"type" api:"required"`
	// URL of the document to fetch.
	URL string `json:"url" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		URL         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsURLDocumentSource) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsURLDocumentSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ManagedAgentsURLDocumentSource to a
// ManagedAgentsURLDocumentSourceParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ManagedAgentsURLDocumentSourceParam.Overrides()
func (r ManagedAgentsURLDocumentSource) ToParam() ManagedAgentsURLDocumentSourceParam {
	return param.Override[ManagedAgentsURLDocumentSourceParam](json.RawMessage(r.RawJSON()))
}

type ManagedAgentsURLDocumentSourceType string

const (
	ManagedAgentsURLDocumentSourceTypeURL ManagedAgentsURLDocumentSourceType = "url"
)

// Document referenced by URL.
//
// The properties Type, URL are required.
type ManagedAgentsURLDocumentSourceParam struct {
	// Any of "url".
	Type ManagedAgentsURLDocumentSourceType `json:"type,omitzero" api:"required"`
	// URL of the document to fetch.
	URL string `json:"url" api:"required"`
	paramObj
}

func (r ManagedAgentsURLDocumentSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsURLDocumentSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsURLDocumentSourceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Image referenced by URL.
type ManagedAgentsURLImageSource struct {
	// Any of "url".
	Type ManagedAgentsURLImageSourceType `json:"type" api:"required"`
	// URL of the image to fetch.
	URL string `json:"url" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		URL         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsURLImageSource) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsURLImageSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ManagedAgentsURLImageSource to a
// ManagedAgentsURLImageSourceParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ManagedAgentsURLImageSourceParam.Overrides()
func (r ManagedAgentsURLImageSource) ToParam() ManagedAgentsURLImageSourceParam {
	return param.Override[ManagedAgentsURLImageSourceParam](json.RawMessage(r.RawJSON()))
}

type ManagedAgentsURLImageSourceType string

const (
	ManagedAgentsURLImageSourceTypeURL ManagedAgentsURLImageSourceType = "url"
)

// Image referenced by URL.
//
// The properties Type, URL are required.
type ManagedAgentsURLImageSourceParam struct {
	// Any of "url".
	Type ManagedAgentsURLImageSourceType `json:"type,omitzero" api:"required"`
	// URL of the image to fetch.
	URL string `json:"url" api:"required"`
	paramObj
}

func (r ManagedAgentsURLImageSourceParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsURLImageSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsURLImageSourceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Event sent by the client providing the result of a custom tool execution.
type ManagedAgentsUserCustomToolResultEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// The id of the `agent.custom_tool_use` event this result corresponds to, which
	// can be found in the last `session.status_idle`
	// [event's](https://docs.qoder.com/cloud-agents/api/sessions/schemas)
	// `stop_reason.event_ids` field.
	CustomToolUseID string `json:"custom_tool_use_id" api:"required"`
	// Any of "user.custom_tool_result".
	Type ManagedAgentsUserCustomToolResultEventType `json:"type" api:"required"`
	// The result content returned by the tool.
	Content []ManagedAgentsUserCustomToolResultEventContentUnion `json:"content"`
	// Whether the tool execution resulted in an error.
	IsError bool `json:"is_error" api:"nullable"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"nullable" format:"date-time"`
	// Routes this result to a subagent thread. Copy from the `agent.custom_tool_use`
	// event's `session_thread_id`.
	SessionThreadID string `json:"session_thread_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID              respjson.Field
		CustomToolUseID respjson.Field
		Type            respjson.Field
		Content         respjson.Field
		IsError         respjson.Field
		ProcessedAt     respjson.Field
		SessionThreadID respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsUserCustomToolResultEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsUserCustomToolResultEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsUserCustomToolResultEventType string

const (
	ManagedAgentsUserCustomToolResultEventTypeUserCustomToolResult ManagedAgentsUserCustomToolResultEventType = "user.custom_tool_result"
)

// ManagedAgentsUserCustomToolResultEventContentUnion contains all possible
// properties and values from [ManagedAgentsTextBlock],
// [ManagedAgentsImageBlock], [ManagedAgentsDocumentBlock],
// [ManagedAgentsSearchResultBlock].
//
// Use the [ManagedAgentsUserCustomToolResultEventContentUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsUserCustomToolResultEventContentUnion struct {
	// This field is from variant [ManagedAgentsTextBlock].
	Text string `json:"text"`
	// Any of "text", "image", "document", "search_result".
	Type string `json:"type"`
	// This field is a union of [ManagedAgentsImageBlockSourceUnion],
	// [ManagedAgentsDocumentBlockSourceUnion], [string]
	Source ManagedAgentsUserCustomToolResultEventContentUnionSource `json:"source"`
	// This field is from variant [ManagedAgentsDocumentBlock].
	Context string `json:"context"`
	Title   string `json:"title"`
	// This field is from variant [ManagedAgentsSearchResultBlock].
	Citations ManagedAgentsSearchResultCitations `json:"citations"`
	// This field is from variant [ManagedAgentsSearchResultBlock].
	Content []ManagedAgentsSearchResultContent `json:"content"`
	JSON    struct {
		Text      respjson.Field
		Type      respjson.Field
		Source    respjson.Field
		Context   respjson.Field
		Title     respjson.Field
		Citations respjson.Field
		Content   respjson.Field
		raw       string
	} `json:"-"`
}

// anyManagedAgentsUserCustomToolResultEventContent is implemented by each
// variant of [ManagedAgentsUserCustomToolResultEventContentUnion] to add type
// safety for the return type of
// [ManagedAgentsUserCustomToolResultEventContentUnion.AsAny]
type anyManagedAgentsUserCustomToolResultEventContent interface {
	implManagedAgentsUserCustomToolResultEventContentUnion()
}

func (ManagedAgentsTextBlock) implManagedAgentsUserCustomToolResultEventContentUnion()     {}
func (ManagedAgentsImageBlock) implManagedAgentsUserCustomToolResultEventContentUnion()    {}
func (ManagedAgentsDocumentBlock) implManagedAgentsUserCustomToolResultEventContentUnion() {}
func (ManagedAgentsSearchResultBlock) implManagedAgentsUserCustomToolResultEventContentUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsUserCustomToolResultEventContentUnion.AsAny().(type) {
//	case qoder.ManagedAgentsTextBlock:
//	case qoder.ManagedAgentsImageBlock:
//	case qoder.ManagedAgentsDocumentBlock:
//	case qoder.ManagedAgentsSearchResultBlock:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsUserCustomToolResultEventContentUnion) AsAny() anyManagedAgentsUserCustomToolResultEventContent {
	switch u.Type {
	case "text":
		return u.AsText()
	case "image":
		return u.AsImage()
	case "document":
		return u.AsDocument()
	case "search_result":
		return u.AsSearchResult()
	}
	return nil
}

func (u ManagedAgentsUserCustomToolResultEventContentUnion) AsText() (v ManagedAgentsTextBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsUserCustomToolResultEventContentUnion) AsImage() (v ManagedAgentsImageBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsUserCustomToolResultEventContentUnion) AsDocument() (v ManagedAgentsDocumentBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsUserCustomToolResultEventContentUnion) AsSearchResult() (v ManagedAgentsSearchResultBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsUserCustomToolResultEventContentUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsUserCustomToolResultEventContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsUserCustomToolResultEventContentUnionSource is an implicit
// subunion of [ManagedAgentsUserCustomToolResultEventContentUnion].
// ManagedAgentsUserCustomToolResultEventContentUnionSource provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsUserCustomToolResultEventContentUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString]
type ManagedAgentsUserCustomToolResultEventContentUnionSource struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString  string `json:",inline"`
	Data      string `json:"data"`
	MediaType string `json:"media_type"`
	Type      string `json:"type"`
	URL       string `json:"url"`
	FileID    string `json:"file_id"`
	JSON      struct {
		OfString  respjson.Field
		Data      respjson.Field
		MediaType respjson.Field
		Type      respjson.Field
		URL       respjson.Field
		FileID    respjson.Field
		raw       string
	} `json:"-"`
}

func (r *ManagedAgentsUserCustomToolResultEventContentUnionSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Parameters for providing the result of a custom tool execution.
//
// The properties CustomToolUseID, Type are required.
type ManagedAgentsUserCustomToolResultEventParams struct {
	// The id of the `agent.custom_tool_use` event this result corresponds to, which
	// can be found in the last `session.status_idle`
	// [event's](https://docs.qoder.com/cloud-agents/api/sessions/schemas)
	// `stop_reason.event_ids` field.
	CustomToolUseID string `json:"custom_tool_use_id" api:"required"`
	// Any of "user.custom_tool_result".
	Type ManagedAgentsUserCustomToolResultEventParamsType `json:"type,omitzero" api:"required"`
	// Whether the tool execution resulted in an error.
	IsError param.Opt[bool] `json:"is_error,omitzero"`
	// The result content returned by the tool.
	Content []ManagedAgentsUserCustomToolResultEventParamsContentUnion `json:"content,omitzero"`
	paramObj
}

func (r ManagedAgentsUserCustomToolResultEventParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsUserCustomToolResultEventParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsUserCustomToolResultEventParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsUserCustomToolResultEventParamsType string

const (
	ManagedAgentsUserCustomToolResultEventParamsTypeUserCustomToolResult ManagedAgentsUserCustomToolResultEventParamsType = "user.custom_tool_result"
)

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsUserCustomToolResultEventParamsContentUnion struct {
	OfText         *ManagedAgentsTextBlockParam         `json:",omitzero,inline"`
	OfImage        *ManagedAgentsImageBlockParam        `json:",omitzero,inline"`
	OfDocument     *ManagedAgentsDocumentBlockParam     `json:",omitzero,inline"`
	OfSearchResult *ManagedAgentsSearchResultBlockParam `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsUserCustomToolResultEventParamsContentUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfText, u.OfImage, u.OfDocument, u.OfSearchResult)
}
func (u *ManagedAgentsUserCustomToolResultEventParamsContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsUserCustomToolResultEventParamsContentUnion) asAny() any {
	if !param.IsOmitted(u.OfText) {
		return u.OfText
	} else if !param.IsOmitted(u.OfImage) {
		return u.OfImage
	} else if !param.IsOmitted(u.OfDocument) {
		return u.OfDocument
	} else if !param.IsOmitted(u.OfSearchResult) {
		return u.OfSearchResult
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsUserCustomToolResultEventParamsContentUnion) GetText() *string {
	if vt := u.OfText; vt != nil {
		return &vt.Text
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsUserCustomToolResultEventParamsContentUnion) GetContext() *string {
	if vt := u.OfDocument; vt != nil && vt.Context.Valid() {
		return &vt.Context.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsUserCustomToolResultEventParamsContentUnion) GetCitations() *ManagedAgentsSearchResultCitationsParam {
	if vt := u.OfSearchResult; vt != nil {
		return &vt.Citations
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsUserCustomToolResultEventParamsContentUnion) GetContent() []ManagedAgentsSearchResultContentParam {
	if vt := u.OfSearchResult; vt != nil {
		return vt.Content
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsUserCustomToolResultEventParamsContentUnion) GetType() *string {
	if vt := u.OfText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfImage; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDocument; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSearchResult; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsUserCustomToolResultEventParamsContentUnion) GetTitle() *string {
	if vt := u.OfDocument; vt != nil && vt.Title.Valid() {
		return &vt.Title.Value
	} else if vt := u.OfSearchResult; vt != nil {
		return (*string)(&vt.Title)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u ManagedAgentsUserCustomToolResultEventParamsContentUnion) GetSource() (res managedAgentsUserCustomToolResultEventParamsContentUnionSource) {
	if vt := u.OfImage; vt != nil {
		res.any = vt.Source.asAny()
	} else if vt := u.OfDocument; vt != nil {
		res.any = vt.Source.asAny()
	} else if vt := u.OfSearchResult; vt != nil {
		res.any = &vt.Source
	}
	return
}

// Can have the runtime types [*ManagedAgentsBase64ImageSourceParam],
// [*ManagedAgentsURLImageSourceParam],
// [*ManagedAgentsFileImageSourceParam],
// [*ManagedAgentsBase64DocumentSourceParam],
// [*ManagedAgentsPlainTextDocumentSourceParam],
// [*ManagedAgentsURLDocumentSourceParam],
// [*ManagedAgentsFileDocumentSourceParam], [*string]
type managedAgentsUserCustomToolResultEventParamsContentUnionSource struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *qoder.ManagedAgentsBase64ImageSourceParam:
//	case *qoder.ManagedAgentsURLImageSourceParam:
//	case *qoder.ManagedAgentsFileImageSourceParam:
//	case *qoder.ManagedAgentsBase64DocumentSourceParam:
//	case *qoder.ManagedAgentsPlainTextDocumentSourceParam:
//	case *qoder.ManagedAgentsURLDocumentSourceParam:
//	case *qoder.ManagedAgentsFileDocumentSourceParam:
//	case *string:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u managedAgentsUserCustomToolResultEventParamsContentUnionSource) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u managedAgentsUserCustomToolResultEventParamsContentUnionSource) GetData() *string {
	switch vt := u.any.(type) {
	case *ManagedAgentsImageBlockSourceUnionParam:
		return vt.GetData()
	case *ManagedAgentsDocumentBlockSourceUnionParam:
		return vt.GetData()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u managedAgentsUserCustomToolResultEventParamsContentUnionSource) GetMediaType() *string {
	switch vt := u.any.(type) {
	case *ManagedAgentsImageBlockSourceUnionParam:
		return vt.GetMediaType()
	case *ManagedAgentsDocumentBlockSourceUnionParam:
		return vt.GetMediaType()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u managedAgentsUserCustomToolResultEventParamsContentUnionSource) GetType() *string {
	switch vt := u.any.(type) {
	case *ManagedAgentsImageBlockSourceUnionParam:
		return vt.GetType()
	case *ManagedAgentsDocumentBlockSourceUnionParam:
		return vt.GetType()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u managedAgentsUserCustomToolResultEventParamsContentUnionSource) GetURL() *string {
	switch vt := u.any.(type) {
	case *ManagedAgentsImageBlockSourceUnionParam:
		return vt.GetURL()
	case *ManagedAgentsDocumentBlockSourceUnionParam:
		return vt.GetURL()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u managedAgentsUserCustomToolResultEventParamsContentUnionSource) GetFileID() *string {
	switch vt := u.any.(type) {
	case *ManagedAgentsImageBlockSourceUnionParam:
		return vt.GetFileID()
	case *ManagedAgentsDocumentBlockSourceUnionParam:
		return vt.GetFileID()
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ManagedAgentsUserCustomToolResultEventParamsContentUnion](
		"type",
		apijson.Discriminator[ManagedAgentsTextBlockParam]("text"),
		apijson.Discriminator[ManagedAgentsImageBlockParam]("image"),
		apijson.Discriminator[ManagedAgentsDocumentBlockParam]("document"),
		apijson.Discriminator[ManagedAgentsSearchResultBlockParam]("search_result"),
	)
}

// Echo of a `user.define_outcome` input event. Carries the server-generated
// `outcome_id` that subsequent `span.outcome_evaluation_*` events reference.
type ManagedAgentsUserDefineOutcomeEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// What the agent should produce. Copied from the input event.
	Description string `json:"description" api:"required"`
	// Evaluate-then-revise cycles before giving up. Default 3, max 20.
	MaxIterations int64 `json:"max_iterations" api:"required"`
	// Server-generated `outc_` ID for this outcome. Referenced by
	// `span.outcome_evaluation_*` events and the session's `outcome_evaluations` list.
	OutcomeID string `json:"outcome_id" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Rubric for grading the quality of an outcome.
	Rubric ManagedAgentsUserDefineOutcomeEventRubricUnion `json:"rubric" api:"required"`
	// Any of "user.define_outcome".
	Type ManagedAgentsUserDefineOutcomeEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID            respjson.Field
		Description   respjson.Field
		MaxIterations respjson.Field
		OutcomeID     respjson.Field
		ProcessedAt   respjson.Field
		Rubric        respjson.Field
		Type          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsUserDefineOutcomeEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsUserDefineOutcomeEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsUserDefineOutcomeEventRubricUnion contains all possible
// properties and values from [ManagedAgentsFileRubric],
// [ManagedAgentsTextRubric].
//
// Use the [ManagedAgentsUserDefineOutcomeEventRubricUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsUserDefineOutcomeEventRubricUnion struct {
	// This field is from variant [ManagedAgentsFileRubric].
	FileID string `json:"file_id"`
	// Any of "file", "text".
	Type string `json:"type"`
	// This field is from variant [ManagedAgentsTextRubric].
	Content string `json:"content"`
	JSON    struct {
		FileID  respjson.Field
		Type    respjson.Field
		Content respjson.Field
		raw     string
	} `json:"-"`
}

// anyManagedAgentsUserDefineOutcomeEventRubric is implemented by each variant
// of [ManagedAgentsUserDefineOutcomeEventRubricUnion] to add type safety for
// the return type of [ManagedAgentsUserDefineOutcomeEventRubricUnion.AsAny]
type anyManagedAgentsUserDefineOutcomeEventRubric interface {
	implManagedAgentsUserDefineOutcomeEventRubricUnion()
}

func (ManagedAgentsFileRubric) implManagedAgentsUserDefineOutcomeEventRubricUnion() {}
func (ManagedAgentsTextRubric) implManagedAgentsUserDefineOutcomeEventRubricUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsUserDefineOutcomeEventRubricUnion.AsAny().(type) {
//	case qoder.ManagedAgentsFileRubric:
//	case qoder.ManagedAgentsTextRubric:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsUserDefineOutcomeEventRubricUnion) AsAny() anyManagedAgentsUserDefineOutcomeEventRubric {
	switch u.Type {
	case "file":
		return u.AsFile()
	case "text":
		return u.AsText()
	}
	return nil
}

func (u ManagedAgentsUserDefineOutcomeEventRubricUnion) AsFile() (v ManagedAgentsFileRubric) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsUserDefineOutcomeEventRubricUnion) AsText() (v ManagedAgentsTextRubric) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsUserDefineOutcomeEventRubricUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsUserDefineOutcomeEventRubricUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsUserDefineOutcomeEventType string

const (
	ManagedAgentsUserDefineOutcomeEventTypeUserDefineOutcome ManagedAgentsUserDefineOutcomeEventType = "user.define_outcome"
)

// Parameters for defining an outcome the agent should work toward. The agent
// begins work on receipt.
//
// The properties Description, Rubric, Type are required.
type ManagedAgentsUserDefineOutcomeEventParams struct {
	// What the agent should produce. This is the task specification.
	Description string `json:"description" api:"required"`
	// Rubric for grading the quality of an outcome.
	Rubric ManagedAgentsUserDefineOutcomeEventParamsRubricUnion `json:"rubric,omitzero" api:"required"`
	// Any of "user.define_outcome".
	Type ManagedAgentsUserDefineOutcomeEventParamsType `json:"type,omitzero" api:"required"`
	// Eval→revision cycles before giving up. Default 3, max 20.
	MaxIterations param.Opt[int64] `json:"max_iterations,omitzero"`
	paramObj
}

func (r ManagedAgentsUserDefineOutcomeEventParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsUserDefineOutcomeEventParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsUserDefineOutcomeEventParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsUserDefineOutcomeEventParamsRubricUnion struct {
	OfFile *ManagedAgentsFileRubricParams `json:",omitzero,inline"`
	OfText *ManagedAgentsTextRubricParams `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsUserDefineOutcomeEventParamsRubricUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfFile, u.OfText)
}
func (u *ManagedAgentsUserDefineOutcomeEventParamsRubricUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsUserDefineOutcomeEventParamsRubricUnion) asAny() any {
	if !param.IsOmitted(u.OfFile) {
		return u.OfFile
	} else if !param.IsOmitted(u.OfText) {
		return u.OfText
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsUserDefineOutcomeEventParamsRubricUnion) GetFileID() *string {
	if vt := u.OfFile; vt != nil {
		return &vt.FileID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsUserDefineOutcomeEventParamsRubricUnion) GetContent() *string {
	if vt := u.OfText; vt != nil {
		return &vt.Content
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsUserDefineOutcomeEventParamsRubricUnion) GetType() *string {
	if vt := u.OfFile; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfText; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ManagedAgentsUserDefineOutcomeEventParamsRubricUnion](
		"type",
		apijson.Discriminator[ManagedAgentsFileRubricParams]("file"),
		apijson.Discriminator[ManagedAgentsTextRubricParams]("text"),
	)
}

type ManagedAgentsUserDefineOutcomeEventParamsType string

const (
	ManagedAgentsUserDefineOutcomeEventParamsTypeUserDefineOutcome ManagedAgentsUserDefineOutcomeEventParamsType = "user.define_outcome"
)

// An interrupt event that pauses agent execution and returns control to the user.
type ManagedAgentsUserInterruptEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// Any of "user.interrupt".
	Type ManagedAgentsUserInterruptEventType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"nullable" format:"date-time"`
	// If absent, interrupts every non-archived thread in a multiagent session (or the
	// primary alone in a single-agent session). If present, interrupts only the named
	// thread.
	SessionThreadID string `json:"session_thread_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID              respjson.Field
		Type            respjson.Field
		ProcessedAt     respjson.Field
		SessionThreadID respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsUserInterruptEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsUserInterruptEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsUserInterruptEventType string

const (
	ManagedAgentsUserInterruptEventTypeUserInterrupt ManagedAgentsUserInterruptEventType = "user.interrupt"
)

// Parameters for sending an interrupt to pause the agent.
//
// The property Type is required.
type ManagedAgentsUserInterruptEventParams struct {
	// Any of "user.interrupt".
	Type ManagedAgentsUserInterruptEventParamsType `json:"type,omitzero" api:"required"`
	// If absent, interrupts every non-archived thread in a multiagent session (or the
	// primary alone in a single-agent session). If present, interrupts only the named
	// thread.
	SessionThreadID param.Opt[string] `json:"session_thread_id,omitzero"`
	paramObj
}

func (r ManagedAgentsUserInterruptEventParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsUserInterruptEventParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsUserInterruptEventParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsUserInterruptEventParamsType string

const (
	ManagedAgentsUserInterruptEventParamsTypeUserInterrupt ManagedAgentsUserInterruptEventParamsType = "user.interrupt"
)

// A user message event in the session conversation.
type ManagedAgentsUserMessageEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// Array of content blocks comprising the user message.
	Content []ManagedAgentsUserMessageEventContentUnion `json:"content" api:"required"`
	// Any of "user.message".
	Type ManagedAgentsUserMessageEventType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"nullable" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Content     respjson.Field
		Type        respjson.Field
		ProcessedAt respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsUserMessageEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsUserMessageEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsUserMessageEventContentUnion contains all possible properties
// and values from [ManagedAgentsTextBlock], [ManagedAgentsImageBlock],
// [ManagedAgentsDocumentBlock], [ManagedAgentsRedactedBlock].
//
// Use the [ManagedAgentsUserMessageEventContentUnion.AsAny] method to switch
// on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsUserMessageEventContentUnion struct {
	// This field is from variant [ManagedAgentsTextBlock].
	Text string `json:"text"`
	// Any of "text", "image", "document", "redacted".
	Type string `json:"type"`
	// This field is a union of [ManagedAgentsImageBlockSourceUnion],
	// [ManagedAgentsDocumentBlockSourceUnion]
	Source ManagedAgentsUserMessageEventContentUnionSource `json:"source"`
	// This field is from variant [ManagedAgentsDocumentBlock].
	Context string `json:"context"`
	// This field is from variant [ManagedAgentsDocumentBlock].
	Title string `json:"title"`
	JSON  struct {
		Text    respjson.Field
		Type    respjson.Field
		Source  respjson.Field
		Context respjson.Field
		Title   respjson.Field
		raw     string
	} `json:"-"`
}

// anyManagedAgentsUserMessageEventContent is implemented by each variant of
// [ManagedAgentsUserMessageEventContentUnion] to add type safety for the
// return type of [ManagedAgentsUserMessageEventContentUnion.AsAny]
type anyManagedAgentsUserMessageEventContent interface {
	implManagedAgentsUserMessageEventContentUnion()
}

func (ManagedAgentsTextBlock) implManagedAgentsUserMessageEventContentUnion()     {}
func (ManagedAgentsImageBlock) implManagedAgentsUserMessageEventContentUnion()    {}
func (ManagedAgentsDocumentBlock) implManagedAgentsUserMessageEventContentUnion() {}
func (ManagedAgentsRedactedBlock) implManagedAgentsUserMessageEventContentUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsUserMessageEventContentUnion.AsAny().(type) {
//	case qoder.ManagedAgentsTextBlock:
//	case qoder.ManagedAgentsImageBlock:
//	case qoder.ManagedAgentsDocumentBlock:
//	case qoder.ManagedAgentsRedactedBlock:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsUserMessageEventContentUnion) AsAny() anyManagedAgentsUserMessageEventContent {
	switch u.Type {
	case "text":
		return u.AsText()
	case "image":
		return u.AsImage()
	case "document":
		return u.AsDocument()
	case "redacted":
		return u.AsRedacted()
	}
	return nil
}

func (u ManagedAgentsUserMessageEventContentUnion) AsText() (v ManagedAgentsTextBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsUserMessageEventContentUnion) AsImage() (v ManagedAgentsImageBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsUserMessageEventContentUnion) AsDocument() (v ManagedAgentsDocumentBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsUserMessageEventContentUnion) AsRedacted() (v ManagedAgentsRedactedBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsUserMessageEventContentUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsUserMessageEventContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsUserMessageEventContentUnionSource is an implicit subunion of
// [ManagedAgentsUserMessageEventContentUnion].
// ManagedAgentsUserMessageEventContentUnionSource provides convenient access
// to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsUserMessageEventContentUnion].
type ManagedAgentsUserMessageEventContentUnionSource struct {
	Data      string `json:"data"`
	MediaType string `json:"media_type"`
	Type      string `json:"type"`
	URL       string `json:"url"`
	FileID    string `json:"file_id"`
	JSON      struct {
		Data      respjson.Field
		MediaType respjson.Field
		Type      respjson.Field
		URL       respjson.Field
		FileID    respjson.Field
		raw       string
	} `json:"-"`
}

func (r *ManagedAgentsUserMessageEventContentUnionSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsUserMessageEventType string

const (
	ManagedAgentsUserMessageEventTypeUserMessage ManagedAgentsUserMessageEventType = "user.message"
)

// Parameters for sending a user message to the session.
//
// The properties Content, Type are required.
type ManagedAgentsUserMessageEventParams struct {
	// Array of content blocks for the user message.
	Content []ManagedAgentsUserMessageEventParamsContentUnion `json:"content,omitzero" api:"required"`
	// Any of "user.message".
	Type ManagedAgentsUserMessageEventParamsType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsUserMessageEventParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsUserMessageEventParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsUserMessageEventParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsUserMessageEventParamsContentUnion struct {
	OfText     *ManagedAgentsTextBlockParam     `json:",omitzero,inline"`
	OfImage    *ManagedAgentsImageBlockParam    `json:",omitzero,inline"`
	OfDocument *ManagedAgentsDocumentBlockParam `json:",omitzero,inline"`
	OfRedacted *ManagedAgentsRedactedBlockParam `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsUserMessageEventParamsContentUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfText, u.OfImage, u.OfDocument, u.OfRedacted)
}
func (u *ManagedAgentsUserMessageEventParamsContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsUserMessageEventParamsContentUnion) asAny() any {
	if !param.IsOmitted(u.OfText) {
		return u.OfText
	} else if !param.IsOmitted(u.OfImage) {
		return u.OfImage
	} else if !param.IsOmitted(u.OfDocument) {
		return u.OfDocument
	} else if !param.IsOmitted(u.OfRedacted) {
		return u.OfRedacted
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsUserMessageEventParamsContentUnion) GetText() *string {
	if vt := u.OfText; vt != nil {
		return &vt.Text
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsUserMessageEventParamsContentUnion) GetContext() *string {
	if vt := u.OfDocument; vt != nil && vt.Context.Valid() {
		return &vt.Context.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsUserMessageEventParamsContentUnion) GetTitle() *string {
	if vt := u.OfDocument; vt != nil && vt.Title.Valid() {
		return &vt.Title.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsUserMessageEventParamsContentUnion) GetType() *string {
	if vt := u.OfText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfImage; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDocument; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRedacted; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u ManagedAgentsUserMessageEventParamsContentUnion) GetSource() (res managedAgentsUserMessageEventParamsContentUnionSource) {
	if vt := u.OfImage; vt != nil {
		res.any = vt.Source.asAny()
	} else if vt := u.OfDocument; vt != nil {
		res.any = vt.Source.asAny()
	}
	return
}

// Can have the runtime types [*ManagedAgentsBase64ImageSourceParam],
// [*ManagedAgentsURLImageSourceParam],
// [*ManagedAgentsFileImageSourceParam],
// [*ManagedAgentsBase64DocumentSourceParam],
// [*ManagedAgentsPlainTextDocumentSourceParam],
// [*ManagedAgentsURLDocumentSourceParam],
// [*ManagedAgentsFileDocumentSourceParam]
type managedAgentsUserMessageEventParamsContentUnionSource struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *qoder.ManagedAgentsBase64ImageSourceParam:
//	case *qoder.ManagedAgentsURLImageSourceParam:
//	case *qoder.ManagedAgentsFileImageSourceParam:
//	case *qoder.ManagedAgentsBase64DocumentSourceParam:
//	case *qoder.ManagedAgentsPlainTextDocumentSourceParam:
//	case *qoder.ManagedAgentsURLDocumentSourceParam:
//	case *qoder.ManagedAgentsFileDocumentSourceParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u managedAgentsUserMessageEventParamsContentUnionSource) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u managedAgentsUserMessageEventParamsContentUnionSource) GetData() *string {
	switch vt := u.any.(type) {
	case *ManagedAgentsImageBlockSourceUnionParam:
		return vt.GetData()
	case *ManagedAgentsDocumentBlockSourceUnionParam:
		return vt.GetData()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u managedAgentsUserMessageEventParamsContentUnionSource) GetMediaType() *string {
	switch vt := u.any.(type) {
	case *ManagedAgentsImageBlockSourceUnionParam:
		return vt.GetMediaType()
	case *ManagedAgentsDocumentBlockSourceUnionParam:
		return vt.GetMediaType()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u managedAgentsUserMessageEventParamsContentUnionSource) GetType() *string {
	switch vt := u.any.(type) {
	case *ManagedAgentsImageBlockSourceUnionParam:
		return vt.GetType()
	case *ManagedAgentsDocumentBlockSourceUnionParam:
		return vt.GetType()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u managedAgentsUserMessageEventParamsContentUnionSource) GetURL() *string {
	switch vt := u.any.(type) {
	case *ManagedAgentsImageBlockSourceUnionParam:
		return vt.GetURL()
	case *ManagedAgentsDocumentBlockSourceUnionParam:
		return vt.GetURL()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u managedAgentsUserMessageEventParamsContentUnionSource) GetFileID() *string {
	switch vt := u.any.(type) {
	case *ManagedAgentsImageBlockSourceUnionParam:
		return vt.GetFileID()
	case *ManagedAgentsDocumentBlockSourceUnionParam:
		return vt.GetFileID()
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ManagedAgentsUserMessageEventParamsContentUnion](
		"type",
		apijson.Discriminator[ManagedAgentsTextBlockParam]("text"),
		apijson.Discriminator[ManagedAgentsImageBlockParam]("image"),
		apijson.Discriminator[ManagedAgentsDocumentBlockParam]("document"),
		apijson.Discriminator[ManagedAgentsRedactedBlockParam]("redacted"),
	)
}

type ManagedAgentsUserMessageEventParamsType string

const (
	ManagedAgentsUserMessageEventParamsTypeUserMessage ManagedAgentsUserMessageEventParamsType = "user.message"
)

// A tool confirmation event that approves or denies a pending tool execution.
type ManagedAgentsUserToolConfirmationEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// UserToolConfirmationResult enum
	//
	// Any of "allow", "deny".
	Result ManagedAgentsUserToolConfirmationEventResult `json:"result" api:"required"`
	// The id of the `agent.tool_use` or `agent.mcp_tool_use` event this result
	// corresponds to, which can be found in the last `session.status_idle`
	// [event's](https://docs.qoder.com/cloud-agents/api/sessions/schemas)
	// `stop_reason.event_ids` field.
	ToolUseID string `json:"tool_use_id" api:"required"`
	// Any of "user.tool_confirmation".
	Type ManagedAgentsUserToolConfirmationEventType `json:"type" api:"required"`
	// Optional message providing context for a 'deny' decision. Only allowed when
	// result is 'deny'.
	DenyMessage string `json:"deny_message" api:"nullable"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"nullable" format:"date-time"`
	// When set, the confirmation routes to this subagent's thread rather than the
	// primary. Echo this from the `session_thread_id` on the `agent.tool_use` or
	// `agent.mcp_tool_use` event that prompted the approval.
	SessionThreadID string `json:"session_thread_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID              respjson.Field
		Result          respjson.Field
		ToolUseID       respjson.Field
		Type            respjson.Field
		DenyMessage     respjson.Field
		ProcessedAt     respjson.Field
		SessionThreadID respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsUserToolConfirmationEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsUserToolConfirmationEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// UserToolConfirmationResult enum
type ManagedAgentsUserToolConfirmationEventResult string

const (
	ManagedAgentsUserToolConfirmationEventResultAllow ManagedAgentsUserToolConfirmationEventResult = "allow"
	ManagedAgentsUserToolConfirmationEventResultDeny  ManagedAgentsUserToolConfirmationEventResult = "deny"
)

type ManagedAgentsUserToolConfirmationEventType string

const (
	ManagedAgentsUserToolConfirmationEventTypeUserToolConfirmation ManagedAgentsUserToolConfirmationEventType = "user.tool_confirmation"
)

// Parameters for confirming or denying a tool execution request.
//
// The properties Result, ToolUseID, Type are required.
type ManagedAgentsUserToolConfirmationEventParams struct {
	// UserToolConfirmationResult enum
	//
	// Any of "allow", "deny".
	Result ManagedAgentsUserToolConfirmationEventParamsResult `json:"result,omitzero" api:"required"`
	// The id of the `agent.tool_use` or `agent.mcp_tool_use` event this result
	// corresponds to, which can be found in the last `session.status_idle`
	// [event's](https://docs.qoder.com/cloud-agents/api/sessions/schemas)
	// `stop_reason.event_ids` field.
	ToolUseID string `json:"tool_use_id" api:"required"`
	// Any of "user.tool_confirmation".
	Type ManagedAgentsUserToolConfirmationEventParamsType `json:"type,omitzero" api:"required"`
	// Optional message providing context for a 'deny' decision. Only allowed when
	// result is 'deny'.
	DenyMessage param.Opt[string] `json:"deny_message,omitzero"`
	paramObj
}

func (r ManagedAgentsUserToolConfirmationEventParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsUserToolConfirmationEventParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsUserToolConfirmationEventParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// UserToolConfirmationResult enum
type ManagedAgentsUserToolConfirmationEventParamsResult string

const (
	ManagedAgentsUserToolConfirmationEventParamsResultAllow ManagedAgentsUserToolConfirmationEventParamsResult = "allow"
	ManagedAgentsUserToolConfirmationEventParamsResultDeny  ManagedAgentsUserToolConfirmationEventParamsResult = "deny"
)

type ManagedAgentsUserToolConfirmationEventParamsType string

const (
	ManagedAgentsUserToolConfirmationEventParamsTypeUserToolConfirmation ManagedAgentsUserToolConfirmationEventParamsType = "user.tool_confirmation"
)

// Parameters for providing the result of an agent-toolset tool execution. Only
// valid on `self_hosted` environments, where sandbox-routed tools are executed by
// the client rather than the server.
//
// The properties ToolUseID, Type are required.
type ManagedAgentsUserToolResultEventParams struct {
	// The id of the `agent.tool_use` event this result corresponds to, which can be
	// found in the last `session.status_idle`
	// [event's](https://docs.qoder.com/cloud-agents/api/sessions/schemas)
	// `stop_reason.event_ids` field.
	ToolUseID string `json:"tool_use_id" api:"required"`
	// Any of "user.tool_result".
	Type ManagedAgentsUserToolResultEventParamsType `json:"type,omitzero" api:"required"`
	// Whether the tool execution resulted in an error.
	IsError param.Opt[bool] `json:"is_error,omitzero"`
	// The result content returned by the tool.
	Content []ManagedAgentsUserToolResultEventParamsContentUnion `json:"content,omitzero"`
	paramObj
}

func (r ManagedAgentsUserToolResultEventParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsUserToolResultEventParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsUserToolResultEventParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsUserToolResultEventParamsType string

const (
	ManagedAgentsUserToolResultEventParamsTypeUserToolResult ManagedAgentsUserToolResultEventParamsType = "user.tool_result"
)

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsUserToolResultEventParamsContentUnion struct {
	OfText         *ManagedAgentsTextBlockParam         `json:",omitzero,inline"`
	OfImage        *ManagedAgentsImageBlockParam        `json:",omitzero,inline"`
	OfDocument     *ManagedAgentsDocumentBlockParam     `json:",omitzero,inline"`
	OfSearchResult *ManagedAgentsSearchResultBlockParam `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsUserToolResultEventParamsContentUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfText, u.OfImage, u.OfDocument, u.OfSearchResult)
}
func (u *ManagedAgentsUserToolResultEventParamsContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsUserToolResultEventParamsContentUnion) asAny() any {
	if !param.IsOmitted(u.OfText) {
		return u.OfText
	} else if !param.IsOmitted(u.OfImage) {
		return u.OfImage
	} else if !param.IsOmitted(u.OfDocument) {
		return u.OfDocument
	} else if !param.IsOmitted(u.OfSearchResult) {
		return u.OfSearchResult
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsUserToolResultEventParamsContentUnion) GetText() *string {
	if vt := u.OfText; vt != nil {
		return &vt.Text
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsUserToolResultEventParamsContentUnion) GetContext() *string {
	if vt := u.OfDocument; vt != nil && vt.Context.Valid() {
		return &vt.Context.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsUserToolResultEventParamsContentUnion) GetCitations() *ManagedAgentsSearchResultCitationsParam {
	if vt := u.OfSearchResult; vt != nil {
		return &vt.Citations
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsUserToolResultEventParamsContentUnion) GetContent() []ManagedAgentsSearchResultContentParam {
	if vt := u.OfSearchResult; vt != nil {
		return vt.Content
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsUserToolResultEventParamsContentUnion) GetType() *string {
	if vt := u.OfText; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfImage; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfDocument; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSearchResult; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsUserToolResultEventParamsContentUnion) GetTitle() *string {
	if vt := u.OfDocument; vt != nil && vt.Title.Valid() {
		return &vt.Title.Value
	} else if vt := u.OfSearchResult; vt != nil {
		return (*string)(&vt.Title)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u ManagedAgentsUserToolResultEventParamsContentUnion) GetSource() (res managedAgentsUserToolResultEventParamsContentUnionSource) {
	if vt := u.OfImage; vt != nil {
		res.any = vt.Source.asAny()
	} else if vt := u.OfDocument; vt != nil {
		res.any = vt.Source.asAny()
	} else if vt := u.OfSearchResult; vt != nil {
		res.any = &vt.Source
	}
	return
}

// Can have the runtime types [*ManagedAgentsBase64ImageSourceParam],
// [*ManagedAgentsURLImageSourceParam],
// [*ManagedAgentsFileImageSourceParam],
// [*ManagedAgentsBase64DocumentSourceParam],
// [*ManagedAgentsPlainTextDocumentSourceParam],
// [*ManagedAgentsURLDocumentSourceParam],
// [*ManagedAgentsFileDocumentSourceParam], [*string]
type managedAgentsUserToolResultEventParamsContentUnionSource struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *qoder.ManagedAgentsBase64ImageSourceParam:
//	case *qoder.ManagedAgentsURLImageSourceParam:
//	case *qoder.ManagedAgentsFileImageSourceParam:
//	case *qoder.ManagedAgentsBase64DocumentSourceParam:
//	case *qoder.ManagedAgentsPlainTextDocumentSourceParam:
//	case *qoder.ManagedAgentsURLDocumentSourceParam:
//	case *qoder.ManagedAgentsFileDocumentSourceParam:
//	case *string:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u managedAgentsUserToolResultEventParamsContentUnionSource) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u managedAgentsUserToolResultEventParamsContentUnionSource) GetData() *string {
	switch vt := u.any.(type) {
	case *ManagedAgentsImageBlockSourceUnionParam:
		return vt.GetData()
	case *ManagedAgentsDocumentBlockSourceUnionParam:
		return vt.GetData()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u managedAgentsUserToolResultEventParamsContentUnionSource) GetMediaType() *string {
	switch vt := u.any.(type) {
	case *ManagedAgentsImageBlockSourceUnionParam:
		return vt.GetMediaType()
	case *ManagedAgentsDocumentBlockSourceUnionParam:
		return vt.GetMediaType()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u managedAgentsUserToolResultEventParamsContentUnionSource) GetType() *string {
	switch vt := u.any.(type) {
	case *ManagedAgentsImageBlockSourceUnionParam:
		return vt.GetType()
	case *ManagedAgentsDocumentBlockSourceUnionParam:
		return vt.GetType()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u managedAgentsUserToolResultEventParamsContentUnionSource) GetURL() *string {
	switch vt := u.any.(type) {
	case *ManagedAgentsImageBlockSourceUnionParam:
		return vt.GetURL()
	case *ManagedAgentsDocumentBlockSourceUnionParam:
		return vt.GetURL()
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u managedAgentsUserToolResultEventParamsContentUnionSource) GetFileID() *string {
	switch vt := u.any.(type) {
	case *ManagedAgentsImageBlockSourceUnionParam:
		return vt.GetFileID()
	case *ManagedAgentsDocumentBlockSourceUnionParam:
		return vt.GetFileID()
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ManagedAgentsUserToolResultEventParamsContentUnion](
		"type",
		apijson.Discriminator[ManagedAgentsTextBlockParam]("text"),
		apijson.Discriminator[ManagedAgentsImageBlockParam]("image"),
		apijson.Discriminator[ManagedAgentsDocumentBlockParam]("document"),
		apijson.Discriminator[ManagedAgentsSearchResultBlockParam]("search_result"),
	)
}

type SessionEventListParams struct {
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	AfterID  param.Opt[string] `query:"after_id,omitzero" json:"-"`

	// Return events created after this time (exclusive). Compared against the event's
	// `processed_at` value.
	CreatedAtGt param.Opt[time.Time] `query:"created_at[gt],omitzero" format:"date-time" json:"-"`
	// Return events created at or after this time (inclusive). Compared against the
	// event's `processed_at` value.
	CreatedAtGte param.Opt[time.Time] `query:"created_at[gte],omitzero" format:"date-time" json:"-"`
	// Return events created before this time (exclusive). Compared against the event's
	// `processed_at` value.
	CreatedAtLt param.Opt[time.Time] `query:"created_at[lt],omitzero" format:"date-time" json:"-"`
	// Return events created at or before this time (inclusive). Compared against the
	// event's `processed_at` value.
	CreatedAtLte param.Opt[time.Time] `query:"created_at[lte],omitzero" format:"date-time" json:"-"`
	// Query parameter for limit
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque pagination cursor from a previous response's `next_page`.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Sort direction for results, ordered by the event's `processed_at`. Defaults to
	// `asc` (chronological).
	//
	// Any of "asc", "desc".
	Order SessionEventListParamsOrder `query:"order,omitzero" json:"-"`
	// Filter by event type. Values match the `type` field on returned events (for
	// example, `user.message` or `agent.tool_use`). Omit to return all event types.
	Types []string `query:"types,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [SessionEventListParams]'s query parameters as
// `url.Values`.
func (r SessionEventListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort direction for results, ordered by the event's `processed_at`. Defaults to
// `asc` (chronological).
type SessionEventListParamsOrder string

const (
	SessionEventListParamsOrderAsc  SessionEventListParamsOrder = "asc"
	SessionEventListParamsOrderDesc SessionEventListParamsOrder = "desc"
)

type SessionEventSendParams struct {
	// Events to send to the `session`.
	Events []ManagedAgentsEventParamsUnion `json:"events,omitzero" api:"required"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

func (r SessionEventSendParams) MarshalJSON() (data []byte, err error) {
	type shadow SessionEventSendParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionEventSendParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SessionEventStreamParams struct {
	// When set, this connection also receives streaming deltas (`event_start`,
	// `event_delta`) while an event is being produced, before the event itself
	// arrives. Deltas are best-effort; when the final event is produced it carries the
	// complete content. A model request that ends early (an error or interrupt)
	// produces no final event — its terminal `span.model_request_end` closes the
	// preview. Accepts one or more event types to preview and may be repeated:
	// `agent.message` streams `content_delta` fragments; `agent.thinking` is
	// start-only — a signal that the agent has begun extended thinking, concluded by
	// the `agent.thinking` event itself. Only previews of the requested event types
	// are sent.
	EventDeltas []ManagedAgentsDeltaType `query:"event_deltas,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [SessionEventStreamParams]'s query parameters as
// `url.Values`.
func (r SessionEventStreamParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
