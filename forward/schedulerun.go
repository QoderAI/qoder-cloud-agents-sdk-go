package forward

import (
	"context"
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
)

// ScheduleRunService provides Forward ScheduleRun operations.
type ScheduleRunService struct {
	Options []option.RequestOption
}

func NewScheduleRunService(opts ...option.RequestOption) ScheduleRunService {
	return ScheduleRunService{Options: slices.Clone(opts)}
}

// 列出 Schedule Runs.
func (r *ScheduleRunService) List(ctx context.Context, params ScheduleRunListParams, opts ...option.RequestOption) (res *pagination.Page[ScheduleRun], err error) {

	opts = slices.Concat(r.Options, opts)
	path := "schedule_runs"
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
func (r *ScheduleRunService) ListAutoPaging(ctx context.Context, params ScheduleRunListParams, opts ...option.RequestOption) *pagination.PageAutoPager[ScheduleRun] {
	return pagination.NewPageAutoPager(r.List(ctx, params, opts...))
}

type ScheduleRunListParams struct {
	// Run 所属 Forward Identity ID。
	IdentityID string `query:"identity_id,omitzero" json:"-" api:"required"`
	// 按 Schedule ID 过滤。
	ScheduleID param.Opt[string] `query:"schedule_id,omitzero" json:"-"`
	// 按 `pending`、`running`、`completed`、`failed` 或 `skipped` 过滤。
	Status param.Opt[string] `query:"status,omitzero" json:"-"`
	// 按 `schedule` 或 `manual` 过滤。
	TriggerType param.Opt[string] `query:"trigger_type,omitzero" json:"-"`
	// 是否只返回有错误或无错误的 Run。
	HasError param.Opt[bool] `query:"has_error,omitzero" json:"-"`
	// 分页大小，最大 100。
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// 向后翻页游标。
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// 向前翻页游标。
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	// 排序字段：`created_at` 或 `triggered_at`。
	SortBy param.Opt[string] `query:"sort_by,omitzero" json:"-"`
	// 排序方向：`asc` 或 `desc`。
	Order param.Opt[string] `query:"order,omitzero" json:"-"`
	paramObj
}

func (r ScheduleRunListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// 获取 Schedule Run.
func (r *ScheduleRunService) Get(ctx context.Context, runID string, params ScheduleRunGetParams, opts ...option.RequestOption) (res *ScheduleRun, err error) {
	if runID == "" {
		return nil, fmt.Errorf("missing required run_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("schedule_runs/%s", url.PathEscape(runID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

type ScheduleRunGetParams struct {
	// 额外归属约束。
	IdentityID param.Opt[string] `query:"identity_id,omitzero" json:"-"`
	paramObj
}

func (r ScheduleRunGetParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

type ScheduleRun struct {
	// Schedule Run ID。
	ID string `json:"id"`
	// 所属 Schedule ID。
	ScheduleID string `json:"schedule_id"`
	// Forward Identity ID。
	IdentityID string `json:"identity_id"`
	// Forward Template ID。
	TemplateID string `json:"template_id"`
	// null|本次执行创建或使用的 Session。
	SessionID string `json:"session_id"`
	// `pending`、`running`、`completed`、`failed` 或 `skipped`。
	Status string `json:"status"`
	// 触发来源。
	TriggerContext ScheduleRunTriggerContext `json:"trigger_context"`
	// null|主流程文本结果。
	ResultPayload string `json:"result_payload"`
	// null|本次 IM 投递使用的 Sink 类型；未配置投递时为 `null`。
	PushSink string `json:"push_sink"`
	// IM 投递状态：`pending`、`succeeded`、`failed` 或 `skipped`。主流程状态与投递状态相互独立。
	PushStatus string `json:"push_status"`
	// null|IM 投递结束时间。
	PushFinishedAt time.Time `json:"push_finished_at" format:"date-time"`
	// 当前或最终实际执行到第几次，从 `1` 开始；当 Schedule 的 `execution.max_attempts=2` 且服务端完成自动重试时，可能返回 `2`。
	Attempt int64 `json:"attempt"`
	// 触发时间。
	TriggeredAt time.Time `json:"triggered_at" format:"date-time"`
	// null|开始执行时间。
	StartedAt time.Time `json:"started_at" format:"date-time"`
	// null|结束时间。
	CompletedAt time.Time `json:"completed_at" format:"date-time"`
	// null|执行耗时，单位毫秒。
	DurationMs int64 `json:"duration_ms"`
	// 记录创建时间。
	CreatedAt time.Time `json:"created_at" format:"date-time"`
	// null|失败或跳过时的结构化错误。
	Error map[string]any `json:"error"`
	// null|便于展示的错误信息；结构化信息保留在 `error`。
	ErrorMessage string `json:"error_message"`
	JSON         struct {
		ID             respjson.Field
		ScheduleID     respjson.Field
		IdentityID     respjson.Field
		TemplateID     respjson.Field
		SessionID      respjson.Field
		Status         respjson.Field
		TriggerContext respjson.Field
		ResultPayload  respjson.Field
		PushSink       respjson.Field
		PushStatus     respjson.Field
		PushFinishedAt respjson.Field
		Attempt        respjson.Field
		TriggeredAt    respjson.Field
		StartedAt      respjson.Field
		CompletedAt    respjson.Field
		DurationMs     respjson.Field
		CreatedAt      respjson.Field
		Error          respjson.Field
		ErrorMessage   respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

func (r ScheduleRun) RawJSON() string                  { return r.JSON.raw }
func (r *ScheduleRun) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type ScheduleRunTriggerContext struct {
	Type        string    `json:"type"`
	ScheduledAt time.Time `json:"scheduled_at" format:"date-time"`
	JSON        struct {
		Type        respjson.Field
		ScheduledAt respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r ScheduleRunTriggerContext) RawJSON() string { return r.JSON.raw }
func (r *ScheduleRunTriggerContext) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
