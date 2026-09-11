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
)

// SessionThreadService contains methods and other services that help with
// interacting with the Qoder Cloud Agents API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSessionThreadService] method instead.
type SessionThreadService struct {
	Options []option.RequestOption
	Events  SessionThreadEventService
}

// NewSessionThreadService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewSessionThreadService(opts ...option.RequestOption) (r SessionThreadService) {
	r = SessionThreadService{}
	r.Options = opts
	r.Events = NewSessionThreadEventService(opts...)
	return
}

// Get Session Thread
func (r *SessionThreadService) Get(ctx context.Context, threadID string, params SessionThreadGetParams, opts ...option.RequestOption) (res *ManagedAgentsSessionThread, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if params.SessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("sessions/%s/threads/%s", url.PathEscape(params.SessionID), url.PathEscape(threadID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List Session Threads
func (r *SessionThreadService) List(ctx context.Context, sessionID string, params SessionThreadListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ManagedAgentsSessionThread], err error) {
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
	path := fmt.Sprintf("sessions/%s/threads", url.PathEscape(sessionID))
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

// List Session Threads
func (r *SessionThreadService) ListAutoPaging(ctx context.Context, sessionID string, params SessionThreadListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ManagedAgentsSessionThread] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, sessionID, params, opts...))
}

// Archive Session Thread
func (r *SessionThreadService) Archive(ctx context.Context, threadID string, params SessionThreadArchiveParams, opts ...option.RequestOption) (res *ManagedAgentsSessionThread, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if params.SessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("sessions/%s/threads/%s/archive", url.PathEscape(params.SessionID), url.PathEscape(threadID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// An execution thread within a `session`. Each session has one primary thread plus
// zero or more child threads spawned by the coordinator.
type ManagedAgentsSessionThread struct {
	// Unique identifier for this thread.
	ID string `json:"id" api:"required"`
	// The resolved agent a session thread runs: a saved-agent snapshot, the platform
	// advisor entry, or an inline-defined (ephemeral) agent snapshot.
	Agent ManagedAgentsSessionThreadAgentUnion `json:"agent" api:"required"`
	// A timestamp in RFC 3339 format
	ArchivedAt time.Time `json:"archived_at" api:"required" format:"date-time"`
	// A timestamp in RFC 3339 format
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Parent thread that spawned this thread. Null for the primary thread.
	ParentThreadID string `json:"parent_thread_id" api:"required"`
	// The session this thread belongs to.
	SessionID string `json:"session_id" api:"required"`
	// Timing statistics for a session thread.
	Stats ManagedAgentsSessionThreadStats `json:"stats" api:"required"`
	// SessionThreadStatus enum
	//
	// Any of "running", "idle", "rescheduling", "terminated".
	Status ManagedAgentsSessionThreadStatus `json:"status" api:"required"`
	// Any of "session_thread".
	Type ManagedAgentsSessionThreadType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// Cumulative token usage for a session thread across all turns.
	Usage ManagedAgentsSessionThreadUsage `json:"usage" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		Agent          respjson.Field
		ArchivedAt     respjson.Field
		CreatedAt      respjson.Field
		ParentThreadID respjson.Field
		SessionID      respjson.Field
		Stats          respjson.Field
		Status         respjson.Field
		Type           respjson.Field
		UpdatedAt      respjson.Field
		Usage          respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionThread) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionThread) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsSessionThreadAgentUnion contains all possible properties and
// values from [ManagedAgentsSessionThreadAgent], [ManagedAgentsAdvisor].
//
// Use the [ManagedAgentsSessionThreadAgentUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsSessionThreadAgentUnion struct {
	// This field is from variant [ManagedAgentsSessionThreadAgent].
	ID string `json:"id"`
	// This field is from variant [ManagedAgentsSessionThreadAgent].
	Description string `json:"description"`
	// This field is from variant [ManagedAgentsSessionThreadAgent].
	MCPServers []ManagedAgentsMCPServerURLDefinition `json:"mcp_servers"`
	// This field is a union of [ManagedAgentsModelConfig], [string]
	Model ManagedAgentsSessionThreadAgentUnionModel `json:"model"`
	// This field is from variant [ManagedAgentsSessionThreadAgent].
	Name string `json:"name"`
	// This field is from variant [ManagedAgentsSessionThreadAgent].
	Skills []ManagedAgentsSessionThreadAgentSkillUnion `json:"skills"`
	// This field is from variant [ManagedAgentsSessionThreadAgent].
	System string `json:"system"`
	// This field is from variant [ManagedAgentsSessionThreadAgent].
	Tools []ManagedAgentsSessionThreadAgentToolUnion `json:"tools"`
	// Any of "agent", "advisor".
	Type string `json:"type"`
	// This field is from variant [ManagedAgentsSessionThreadAgent].
	Version int64 `json:"version"`
	JSON    struct {
		ID          respjson.Field
		Description respjson.Field
		MCPServers  respjson.Field
		Model       respjson.Field
		Name        respjson.Field
		Skills      respjson.Field
		System      respjson.Field
		Tools       respjson.Field
		Type        respjson.Field
		Version     respjson.Field
		raw         string
	} `json:"-"`
}

// anyManagedAgentsSessionThreadAgent is implemented by each variant of
// [ManagedAgentsSessionThreadAgentUnion] to add type safety for the return
// type of [ManagedAgentsSessionThreadAgentUnion.AsAny]
type anyManagedAgentsSessionThreadAgent interface {
	implManagedAgentsSessionThreadAgentUnion()
}

func (ManagedAgentsSessionThreadAgent) implManagedAgentsSessionThreadAgentUnion() {}
func (ManagedAgentsAdvisor) implManagedAgentsSessionThreadAgentUnion()            {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsSessionThreadAgentUnion.AsAny().(type) {
//	case qoder.ManagedAgentsSessionThreadAgent:
//	case qoder.ManagedAgentsAdvisor:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsSessionThreadAgentUnion) AsAny() anyManagedAgentsSessionThreadAgent {
	switch u.Type {
	case "agent":
		return u.AsAgent()
	case "advisor":
		return u.AsAdvisor()
	}
	return nil
}

func (u ManagedAgentsSessionThreadAgentUnion) AsAgent() (v ManagedAgentsSessionThreadAgent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionThreadAgentUnion) AsAdvisor() (v ManagedAgentsAdvisor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsSessionThreadAgentUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsSessionThreadAgentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsSessionThreadAgentUnionModel is an implicit subunion of
// [ManagedAgentsSessionThreadAgentUnion].
// ManagedAgentsSessionThreadAgentUnionModel provides convenient access to the
// sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsSessionThreadAgentUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString]
type ManagedAgentsSessionThreadAgentUnionModel struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	// This field is from variant [ManagedAgentsModelConfig].
	ID ManagedAgentsModel `json:"id"`
	// This field is from variant [ManagedAgentsModelConfig].
	Effort ManagedAgentsModelConfigEffortUnion `json:"effort"`
	// This field is from variant [ManagedAgentsModelConfig].
	InferenceGeo string `json:"inference_geo"`
	// This field is from variant [ManagedAgentsModelConfig].
	Speed ManagedAgentsModelConfigSpeed `json:"speed"`
	JSON  struct {
		OfString     respjson.Field
		ID           respjson.Field
		Effort       respjson.Field
		InferenceGeo respjson.Field
		Speed        respjson.Field
		raw          string
	} `json:"-"`
}

func (r *ManagedAgentsSessionThreadAgentUnionModel) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSessionThreadType string

const (
	ManagedAgentsSessionThreadTypeSessionThread ManagedAgentsSessionThreadType = "session_thread"
)

// Timing statistics for a session thread.
type ManagedAgentsSessionThreadStats struct {
	// Cumulative time in seconds the thread spent actively running. Excludes idle
	// time.
	ActiveSeconds float64 `json:"active_seconds"`
	// Elapsed time since thread creation in seconds. For archived threads, frozen at
	// the final update.
	DurationSeconds float64 `json:"duration_seconds"`
	// Time in seconds for the thread to begin running. Zero for child threads, which
	// start immediately.
	StartupSeconds float64 `json:"startup_seconds"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ActiveSeconds   respjson.Field
		DurationSeconds respjson.Field
		StartupSeconds  respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionThreadStats) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionThreadStats) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// SessionThreadStatus enum
type ManagedAgentsSessionThreadStatus string

const (
	ManagedAgentsSessionThreadStatusRunning      ManagedAgentsSessionThreadStatus = "running"
	ManagedAgentsSessionThreadStatusIdle         ManagedAgentsSessionThreadStatus = "idle"
	ManagedAgentsSessionThreadStatusRescheduling ManagedAgentsSessionThreadStatus = "rescheduling"
	ManagedAgentsSessionThreadStatusTerminated   ManagedAgentsSessionThreadStatus = "terminated"
)

// Cumulative token usage for a session thread across all turns.
type ManagedAgentsSessionThreadUsage struct {
	// Cumulative time in seconds this thread spent in running status. Equal to
	// `stats.active_seconds`; surfaced here so a thread's usage carries every quantity
	// its cost is priced on.
	ActiveSeconds float64 `json:"active_seconds"`
	// Prompt-cache creation token usage broken down by cache lifetime.
	CacheCreation ManagedAgentsCacheCreationUsage `json:"cache_creation"`
	// Total tokens read from prompt cache.
	CacheReadInputTokens int64 `json:"cache_read_input_tokens"`
	// Total input tokens consumed across all turns.
	InputTokens int64 `json:"input_tokens"`
	// A monetary amount in a specific currency.
	ListCost MonetaryAmount `json:"list_cost" api:"nullable"`
	// Total output tokens generated across all turns.
	OutputTokens int64 `json:"output_tokens"`
	// Cumulative count of server-executed tool invocations, broken down by tool.
	ServerToolUse ManagedAgentsServerToolUsage `json:"server_tool_use" api:"nullable"`
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
func (r ManagedAgentsSessionThreadUsage) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionThreadUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsStreamSessionThreadEventsUnion contains all possible properties
// and values from [ManagedAgentsUserMessageEvent],
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
// Use the [ManagedAgentsStreamSessionThreadEventsUnion.AsAny] method to switch
// on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsStreamSessionThreadEventsUnion struct {
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
	Content ManagedAgentsStreamSessionThreadEventsUnionContent `json:"content"`
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
	StopReason ManagedAgentsStreamSessionThreadEventsUnionStopReason `json:"stop_reason"`
	AgentName  string                                                `json:"agent_name"`
	Iteration  int64                                                 `json:"iteration"`
	OutcomeID  string                                                `json:"outcome_id"`
	// This field is from variant [ManagedAgentsSpanOutcomeEvaluationEndEvent].
	Explanation string `json:"explanation"`
	// This field is from variant [ManagedAgentsSpanOutcomeEvaluationEndEvent].
	OutcomeEvaluationStartID string `json:"outcome_evaluation_start_id"`
	// This field is a union of [ManagedAgentsSpanModelUsage],
	// [ManagedAgentsSessionUsageSnapshot]
	Usage ManagedAgentsStreamSessionThreadEventsUnionUsage `json:"usage"`
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

// anyManagedAgentsStreamSessionThreadEvents is implemented by each variant of
// [ManagedAgentsStreamSessionThreadEventsUnion] to add type safety for the
// return type of [ManagedAgentsStreamSessionThreadEventsUnion.AsAny]
type anyManagedAgentsStreamSessionThreadEvents interface {
	implManagedAgentsStreamSessionThreadEventsUnion()
}

func (ManagedAgentsUserMessageEvent) implManagedAgentsStreamSessionThreadEventsUnion()   {}
func (ManagedAgentsUserInterruptEvent) implManagedAgentsStreamSessionThreadEventsUnion() {}
func (ManagedAgentsUserToolConfirmationEvent) implManagedAgentsStreamSessionThreadEventsUnion() {
}
func (ManagedAgentsUserCustomToolResultEvent) implManagedAgentsStreamSessionThreadEventsUnion() {
}
func (ManagedAgentsAgentCustomToolUseEvent) implManagedAgentsStreamSessionThreadEventsUnion() {
}
func (ManagedAgentsAgentMessageEvent) implManagedAgentsStreamSessionThreadEventsUnion()    {}
func (ManagedAgentsAgentThinkingEvent) implManagedAgentsStreamSessionThreadEventsUnion()   {}
func (ManagedAgentsAgentMCPToolUseEvent) implManagedAgentsStreamSessionThreadEventsUnion() {}
func (ManagedAgentsAgentMCPToolResultEvent) implManagedAgentsStreamSessionThreadEventsUnion() {
}
func (ManagedAgentsAgentToolUseEvent) implManagedAgentsStreamSessionThreadEventsUnion()    {}
func (ManagedAgentsAgentToolResultEvent) implManagedAgentsStreamSessionThreadEventsUnion() {}
func (ManagedAgentsAgentThreadMessageReceivedEvent) implManagedAgentsStreamSessionThreadEventsUnion() {
}
func (ManagedAgentsAgentThreadMessageSentEvent) implManagedAgentsStreamSessionThreadEventsUnion() {
}
func (ManagedAgentsAgentThreadContextCompactedEvent) implManagedAgentsStreamSessionThreadEventsUnion() {
}
func (ManagedAgentsSessionErrorEvent) implManagedAgentsStreamSessionThreadEventsUnion() {}
func (ManagedAgentsSessionStatusRescheduledEvent) implManagedAgentsStreamSessionThreadEventsUnion() {
}
func (ManagedAgentsSessionStatusRunningEvent) implManagedAgentsStreamSessionThreadEventsUnion() {
}
func (ManagedAgentsSessionStatusIdleEvent) implManagedAgentsStreamSessionThreadEventsUnion() {
}
func (ManagedAgentsSessionStatusTerminatedEvent) implManagedAgentsStreamSessionThreadEventsUnion() {
}
func (ManagedAgentsSessionThreadCreatedEvent) implManagedAgentsStreamSessionThreadEventsUnion() {
}
func (ManagedAgentsSpanOutcomeEvaluationStartEvent) implManagedAgentsStreamSessionThreadEventsUnion() {
}
func (ManagedAgentsSpanOutcomeEvaluationEndEvent) implManagedAgentsStreamSessionThreadEventsUnion() {
}
func (ManagedAgentsSpanModelRequestStartEvent) implManagedAgentsStreamSessionThreadEventsUnion() {
}
func (ManagedAgentsSpanModelRequestEndEvent) implManagedAgentsStreamSessionThreadEventsUnion() {
}
func (ManagedAgentsSpanOutcomeEvaluationOngoingEvent) implManagedAgentsStreamSessionThreadEventsUnion() {
}
func (ManagedAgentsUserDefineOutcomeEvent) implManagedAgentsStreamSessionThreadEventsUnion() {
}
func (ManagedAgentsSessionDeletedEvent) implManagedAgentsStreamSessionThreadEventsUnion() {}
func (ManagedAgentsSessionThreadStatusRunningEvent) implManagedAgentsStreamSessionThreadEventsUnion() {
}
func (ManagedAgentsSessionThreadStatusIdleEvent) implManagedAgentsStreamSessionThreadEventsUnion() {
}
func (ManagedAgentsSessionThreadStatusTerminatedEvent) implManagedAgentsStreamSessionThreadEventsUnion() {
}
func (ManagedAgentsUserToolResultEvent) implManagedAgentsStreamSessionThreadEventsUnion() {}
func (ManagedAgentsSessionThreadStatusRescheduledEvent) implManagedAgentsStreamSessionThreadEventsUnion() {
}
func (ManagedAgentsSessionUpdatedEvent) implManagedAgentsStreamSessionThreadEventsUnion() {}
func (ManagedAgentsStartEvent) implManagedAgentsStreamSessionThreadEventsUnion()          {}
func (ManagedAgentsDeltaEvent) implManagedAgentsStreamSessionThreadEventsUnion()          {}
func (ManagedAgentsSystemMessageEvent) implManagedAgentsStreamSessionThreadEventsUnion()  {}
func (ManagedAgentsSessionUsageEvent) implManagedAgentsStreamSessionThreadEventsUnion()   {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsStreamSessionThreadEventsUnion.AsAny().(type) {
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
func (u ManagedAgentsStreamSessionThreadEventsUnion) AsAny() anyManagedAgentsStreamSessionThreadEvents {
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

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsUserMessage() (v ManagedAgentsUserMessageEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsUserInterrupt() (v ManagedAgentsUserInterruptEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsUserToolConfirmation() (v ManagedAgentsUserToolConfirmationEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsUserCustomToolResult() (v ManagedAgentsUserCustomToolResultEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsAgentCustomToolUse() (v ManagedAgentsAgentCustomToolUseEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsAgentMessage() (v ManagedAgentsAgentMessageEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsAgentThinking() (v ManagedAgentsAgentThinkingEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsAgentMCPToolUse() (v ManagedAgentsAgentMCPToolUseEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsAgentMCPToolResult() (v ManagedAgentsAgentMCPToolResultEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsAgentToolUse() (v ManagedAgentsAgentToolUseEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsAgentToolResult() (v ManagedAgentsAgentToolResultEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsAgentThreadMessageReceived() (v ManagedAgentsAgentThreadMessageReceivedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsAgentThreadMessageSent() (v ManagedAgentsAgentThreadMessageSentEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsAgentThreadContextCompacted() (v ManagedAgentsAgentThreadContextCompactedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsSessionError() (v ManagedAgentsSessionErrorEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsSessionStatusRescheduled() (v ManagedAgentsSessionStatusRescheduledEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsSessionStatusRunning() (v ManagedAgentsSessionStatusRunningEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsSessionStatusIdle() (v ManagedAgentsSessionStatusIdleEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsSessionStatusTerminated() (v ManagedAgentsSessionStatusTerminatedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsSessionThreadCreated() (v ManagedAgentsSessionThreadCreatedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsSpanOutcomeEvaluationStart() (v ManagedAgentsSpanOutcomeEvaluationStartEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsSpanOutcomeEvaluationEnd() (v ManagedAgentsSpanOutcomeEvaluationEndEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsSpanModelRequestStart() (v ManagedAgentsSpanModelRequestStartEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsSpanModelRequestEnd() (v ManagedAgentsSpanModelRequestEndEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsSpanOutcomeEvaluationOngoing() (v ManagedAgentsSpanOutcomeEvaluationOngoingEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsUserDefineOutcome() (v ManagedAgentsUserDefineOutcomeEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsSessionDeleted() (v ManagedAgentsSessionDeletedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsSessionThreadStatusRunning() (v ManagedAgentsSessionThreadStatusRunningEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsSessionThreadStatusIdle() (v ManagedAgentsSessionThreadStatusIdleEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsSessionThreadStatusTerminated() (v ManagedAgentsSessionThreadStatusTerminatedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsUserToolResult() (v ManagedAgentsUserToolResultEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsSessionThreadStatusRescheduled() (v ManagedAgentsSessionThreadStatusRescheduledEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsSessionUpdated() (v ManagedAgentsSessionUpdatedEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsEventStart() (v ManagedAgentsStartEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsEventDelta() (v ManagedAgentsDeltaEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsSystemMessage() (v ManagedAgentsSystemMessageEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStreamSessionThreadEventsUnion) AsSessionUsage() (v ManagedAgentsSessionUsageEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsStreamSessionThreadEventsUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsStreamSessionThreadEventsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsStreamSessionThreadEventsUnionContent is an implicit subunion
// of [ManagedAgentsStreamSessionThreadEventsUnion].
// ManagedAgentsStreamSessionThreadEventsUnionContent provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsStreamSessionThreadEventsUnion].
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
type ManagedAgentsStreamSessionThreadEventsUnionContent struct {
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

func (r *ManagedAgentsStreamSessionThreadEventsUnionContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsStreamSessionThreadEventsUnionStopReason is an implicit
// subunion of [ManagedAgentsStreamSessionThreadEventsUnion].
// ManagedAgentsStreamSessionThreadEventsUnionStopReason provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsStreamSessionThreadEventsUnion].
type ManagedAgentsStreamSessionThreadEventsUnionStopReason struct {
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

func (r *ManagedAgentsStreamSessionThreadEventsUnionStopReason) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsStreamSessionThreadEventsUnionUsage is an implicit subunion of
// [ManagedAgentsStreamSessionThreadEventsUnion].
// ManagedAgentsStreamSessionThreadEventsUnionUsage provides convenient access
// to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsStreamSessionThreadEventsUnion].
type ManagedAgentsStreamSessionThreadEventsUnionUsage struct {
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

func (r *ManagedAgentsStreamSessionThreadEventsUnionUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SessionThreadGetParams struct {
	SessionID string `path:"session_id" api:"required" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

type SessionThreadListParams struct {
	// Maximum results per page. Defaults to 1000.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque pagination cursor from a previous response's `next_page`. Forward-only.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [SessionThreadListParams]'s query parameters as
// `url.Values`.
func (r SessionThreadListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SessionThreadArchiveParams struct {
	SessionID string `path:"session_id" api:"required" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}
