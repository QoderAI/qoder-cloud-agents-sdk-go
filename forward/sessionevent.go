package forward

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apiquery"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/pagination"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/ssestream"
)

// SessionEventService provides Forward SessionEvent operations.
type SessionEventService struct {
	Options []option.RequestOption
}

func NewSessionEventService(opts ...option.RequestOption) SessionEventService {
	return SessionEventService{Options: slices.Clone(opts)}
}

// 列出 Session Events.
func (r *SessionEventService) List(ctx context.Context, sessionID string, params SessionEventListParams, opts ...option.RequestOption) (res *pagination.Page[SessionEvent], err error) {
	if sessionID == "" {
		return nil, fmt.Errorf("missing required session_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("sessions/%s/events", url.PathEscape(sessionID))
	var raw *http.Response
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	cfg, err := convention.NewRequestConfig(ctx, http.MethodGet, path, params, &res, opts...)
	if err != nil {
		return nil, err
	}
	if err = cfg.Execute(); err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}
func (r *SessionEventService) ListAutoPaging(ctx context.Context, sessionID string, params SessionEventListParams, opts ...option.RequestOption) *pagination.PageAutoPager[SessionEvent] {
	return pagination.NewPageAutoPager(r.List(ctx, sessionID, params, opts...))
}

type SessionEventListParams struct {
	// 分页大小，最大 100。
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// 返回该 Event ID 之后的事件。
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// 返回该 Event ID 之前的事件。
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	// 排序方向：`asc` 或 `desc`。
	Order param.Opt[string] `query:"order,omitzero" json:"-"`
	// 按 Event 类型过滤，支持逗号分隔。
	Type param.Opt[string] `query:"type,omitzero" json:"-"`
	// 数组形式的 Event 类型过滤。
	Types []string `query:"types[],omitzero" json:"-"`
	// 是否包含工具调用类事件。
	IncludeToolCalls param.Opt[bool] `query:"include_tool_calls,omitzero" json:"-"`
	// 是否包含思考过程事件。
	IncludeThinking param.Opt[bool] `query:"include_thinking,omitzero" json:"-"`
	paramObj
}

func (r SessionEventListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// 发送 Session Events.
func (r *SessionEventService) Send(ctx context.Context, sessionID string, params SessionEventSendParams, opts ...option.RequestOption) (res *SessionEventSendResponse, err error) {
	if sessionID == "" {
		return nil, fmt.Errorf("missing required session_id parameter")
	}
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("sessions/%s/events", url.PathEscape(sessionID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type SessionEventSendParams struct {
	Events []SessionEventParam `json:"events" api:"required"`
	// 有副作用请求可选的幂等键。
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r SessionEventSendParams) MarshalJSON() ([]byte, error) {
	type shadow SessionEventSendParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionEventSendParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// 订阅 Session Event Stream.
func (r *SessionEventService) StreamEvents(ctx context.Context, sessionID string, params SessionEventStreamParams, opts ...option.RequestOption) *ssestream.Stream[SessionEvent] {
	if sessionID == "" {
		return ssestream.NewStream[SessionEvent](nil, fmt.Errorf("missing required session_id parameter"))
	}
	if params.LastEventID.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Last-Event-ID", fmt.Sprint(params.LastEventID.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, []option.RequestOption{option.WithHeader("Accept", "text/event-stream")}, opts)
	path := fmt.Sprintf("sessions/%s/events/stream", url.PathEscape(sessionID))
	var raw *http.Response
	err := convention.ExecuteNewRequest(ctx, http.MethodGet, path, params, &raw, opts...)
	return ssestream.NewStream[SessionEvent](ssestream.NewDecoder(raw), err)
}

type SessionEventStreamParams struct {
	// 订阅指定公开事件类型的流式增量事件。支持重复传参，取值见 [流式增量事件](../Session&Event数据结构.md#流式增量事件)。
	EventDeltas []string `query:"event_deltas[],omitzero" json:"-"`
	// 是否包含工具调用类事件。
	IncludeToolCalls param.Opt[bool] `query:"include_tool_calls,omitzero" json:"-"`
	// 是否包含思考过程事件。
	IncludeThinking param.Opt[bool] `query:"include_thinking,omitzero" json:"-"`
	// 从该 Event ID 之后恢复订阅。
	LastEventID param.Opt[string] `header:"Last-Event-ID,omitzero" json:"-"`
	paramObj
}

func (r SessionEventStreamParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

type SessionEvent struct {
	ID                  string          `json:"id"`
	Type                string          `json:"type"`
	SessionID           string          `json:"session_id"`
	Content             json.RawMessage `json:"content"`
	ProcessedAt         time.Time       `json:"processed_at" format:"date-time"`
	SessionThreadID     string          `json:"session_thread_id"`
	Thinking            string          `json:"thinking"`
	Text                string          `json:"text"`
	ToolUseID           string          `json:"tool_use_id"`
	CustomToolUseID     string          `json:"custom_tool_use_id"`
	MCPToolUseID        string          `json:"mcp_tool_use_id"`
	Name                string          `json:"name"`
	MCPServerName       string          `json:"mcp_server_name"`
	Input               map[string]any  `json:"input"`
	IsError             bool            `json:"is_error"`
	Result              string          `json:"result"`
	DenyMessage         string          `json:"deny_message"`
	Description         string          `json:"description"`
	Rubric              string          `json:"rubric"`
	OutcomeID           string          `json:"outcome_id"`
	MaxIterations       int64           `json:"max_iterations"`
	MessageID           string          `json:"message_id"`
	Message             json.RawMessage `json:"message"`
	Index               int64           `json:"index"`
	ContentBlock        json.RawMessage `json:"content_block"`
	Delta               json.RawMessage `json:"delta"`
	Event               json.RawMessage `json:"event"`
	EventID             string          `json:"event_id"`
	Usage               map[string]any  `json:"usage"`
	StopReason          map[string]any  `json:"stop_reason"`
	Error               map[string]any  `json:"error"`
	EvaluatedPermission map[string]any  `json:"evaluated_permission"`
	FileID              string          `json:"file_id"`
	OriginalFilename    string          `json:"original_filename"`
	Size                int64           `json:"size"`
	ContentType         string          `json:"content_type"`
	Agent               map[string]any  `json:"agent"`
	Metadata            map[string]any  `json:"metadata"`
	Title               string          `json:"title"`
	ModelRequestStartID string          `json:"model_request_start_id"`
	ModelUsage          map[string]any  `json:"model_usage"`
	JSON                struct {
		ID                  respjson.Field
		Type                respjson.Field
		SessionID           respjson.Field
		Content             respjson.Field
		ProcessedAt         respjson.Field
		SessionThreadID     respjson.Field
		Thinking            respjson.Field
		Text                respjson.Field
		ToolUseID           respjson.Field
		CustomToolUseID     respjson.Field
		MCPToolUseID        respjson.Field
		Name                respjson.Field
		MCPServerName       respjson.Field
		Input               respjson.Field
		IsError             respjson.Field
		Result              respjson.Field
		DenyMessage         respjson.Field
		Description         respjson.Field
		Rubric              respjson.Field
		OutcomeID           respjson.Field
		MaxIterations       respjson.Field
		MessageID           respjson.Field
		Message             respjson.Field
		Index               respjson.Field
		ContentBlock        respjson.Field
		Delta               respjson.Field
		Event               respjson.Field
		EventID             respjson.Field
		Usage               respjson.Field
		StopReason          respjson.Field
		Error               respjson.Field
		EvaluatedPermission respjson.Field
		FileID              respjson.Field
		OriginalFilename    respjson.Field
		Size                respjson.Field
		ContentType         respjson.Field
		Agent               respjson.Field
		Metadata            respjson.Field
		Title               respjson.Field
		ModelRequestStartID respjson.Field
		ModelUsage          respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

func (r SessionEvent) RawJSON() string                  { return r.JSON.raw }
func (r *SessionEvent) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type SessionEventSendResponse struct {
	Data []SessionEvent `json:"data"`
	JSON struct {
		Data        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r SessionEventSendResponse) RawJSON() string { return r.JSON.raw }
func (r *SessionEventSendResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
