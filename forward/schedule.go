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

// ScheduleService provides Forward Schedule operations.
type ScheduleService struct {
	Options []option.RequestOption
}

func NewScheduleService(opts ...option.RequestOption) ScheduleService {
	return ScheduleService{Options: slices.Clone(opts)}
}

// 列出 Schedules.
func (r *ScheduleService) List(ctx context.Context, params ScheduleListParams, opts ...option.RequestOption) (res *pagination.Page[Schedule], err error) {

	opts = slices.Concat(r.Options, opts)
	path := "schedules"
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
func (r *ScheduleService) ListAutoPaging(ctx context.Context, params ScheduleListParams, opts ...option.RequestOption) *pagination.PageAutoPager[Schedule] {
	return pagination.NewPageAutoPager(r.List(ctx, params, opts...))
}

type ScheduleListParams struct {
	// PAT 或管理员 SAT 可省略，省略时查询当前 owner 全部 Identity；Identity-bound SAT 省略时自动绑定自身，显式传其他 Identity 返回 403。
	IdentityID param.Opt[string] `query:"identity_id,omitzero" json:"-"`
	// 按 Forward Template ID 过滤。
	TemplateID param.Opt[string] `query:"template_id,omitzero" json:"-"`
	// 按 `active` 或 `paused` 过滤。
	Status param.Opt[string] `query:"status,omitzero" json:"-"`
	// 是否包含已归档 Schedule。
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// 分页大小，最大 100。
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// 向后翻页游标。
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// 向前翻页游标。
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	// 排序字段：`created_at` 或 `upcoming_runs_at`。
	SortBy param.Opt[string] `query:"sort_by,omitzero" json:"-"`
	// 排序方向：`asc` 或 `desc`。
	Order param.Opt[string] `query:"order,omitzero" json:"-"`
	paramObj
}

func (r ScheduleListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// 创建 Schedule.
func (r *ScheduleService) New(ctx context.Context, params ScheduleNewParams, opts ...option.RequestOption) (res *Schedule, err error) {
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := "schedules"
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type ScheduleNewParams struct {
	// Schedule 所属 Forward Identity ID。
	IdentityID string `json:"identity_id" api:"required"`
	// 要执行的 Forward Template ID。
	TemplateID string `json:"template_id" api:"required"`
	// Schedule 名称。
	Name string `json:"name" api:"required"`
	// Schedule 描述。
	Description param.Opt[string] `json:"description,omitzero"`
	// 每次执行注入的初始事件，当前支持 `user.message`。
	InitialEvents []map[string]any `json:"initial_events" api:"required"`
	// 执行策略；省略时使用服务端默认值。
	Execution map[string]any `json:"execution,omitzero"`
	// 触发策略；省略或 `null` 时按 `manual` 处理。
	TriggerPolicy map[string]any `json:"trigger_policy,omitzero" api:"nullable"`
	// 执行环境。
	EnvironmentID string `json:"environment_id" api:"required"`
	// 执行结果推送目标；为兼容性保留数组形式，当前最多允许一个元素。
	Sinks []map[string]any `json:"sinks,omitzero" api:"nullable"`
	// 业务元数据，仅用于标签或透传。
	Metadata map[string]any `json:"metadata,omitzero"`
	// 有副作用请求可选的幂等键。
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r ScheduleNewParams) MarshalJSON() ([]byte, error) {
	type shadow ScheduleNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ScheduleNewParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

// 批量归档 Schedules.
func (r *ScheduleService) ArchiveMany(ctx context.Context, params ScheduleArchiveManyParams, opts ...option.RequestOption) (res *ScheduleArchiveManyResponse, err error) {
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := "schedules/archive"
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type ScheduleArchiveManyParams struct {
	// 归档范围，当前仅支持 by_schedule_ids。
	Scope string `json:"scope" api:"required"`
	// 去重后必须包含 1～50 个非空 Schedule ID。
	ScheduleIDs []string `json:"schedule_ids" api:"required"`
	// 有副作用请求可选的幂等键；相同 owner、路径和请求体可安全重放。
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r ScheduleArchiveManyParams) MarshalJSON() ([]byte, error) {
	type shadow ScheduleArchiveManyParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ScheduleArchiveManyParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// 获取 Schedule.
func (r *ScheduleService) Get(ctx context.Context, scheduleID string, opts ...option.RequestOption) (res *Schedule, err error) {
	if scheduleID == "" {
		return nil, fmt.Errorf("missing required schedule_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("schedules/%s", url.PathEscape(scheduleID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// 更新 Schedule.
func (r *ScheduleService) Update(ctx context.Context, scheduleID string, params ScheduleUpdateParams, opts ...option.RequestOption) (res *Schedule, err error) {
	if scheduleID == "" {
		return nil, fmt.Errorf("missing required schedule_id parameter")
	}
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("schedules/%s", url.PathEscape(scheduleID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type ScheduleUpdateParams struct {
	// 新的 Schedule 名称。
	Name param.Opt[string] `json:"name,omitzero"`
	// 新的 Schedule 描述。
	Description param.Opt[string] `json:"description,omitzero"`
	// 新的 Forward Template ID。
	TemplateID param.Opt[string] `json:"template_id,omitzero"`
	// 替换初始事件列表。
	InitialEvents []map[string]any `json:"initial_events,omitzero"`
	// 合并更新执行策略。
	Execution map[string]any `json:"execution,omitzero"`
	// 更新触发策略；`null` 表示改为 manual。
	TriggerPolicy map[string]any `json:"trigger_policy,omitzero" api:"nullable"`
	// 新的执行环境。
	EnvironmentID param.Opt[string] `json:"environment_id,omitzero"`
	// 执行结果推送目标；为兼容性保留数组形式，当前最多允许一个元素。
	Sinks []map[string]any `json:"sinks,omitzero" api:"nullable"`
	// 合并更新 metadata；value 为 `null` 删除 key。
	Metadata map[string]any `json:"metadata,omitzero"`
	// 有副作用请求可选的幂等键。
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r ScheduleUpdateParams) MarshalJSON() ([]byte, error) {
	type shadow ScheduleUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ScheduleUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// 归档 Schedule.
func (r *ScheduleService) Archive(ctx context.Context, scheduleID string, params ScheduleArchiveParams, opts ...option.RequestOption) (res *Schedule, err error) {
	if scheduleID == "" {
		return nil, fmt.Errorf("missing required schedule_id parameter")
	}
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("schedules/%s/archive", url.PathEscape(scheduleID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type ScheduleArchiveParams struct {
	// 有副作用请求可选的幂等键。
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r ScheduleArchiveParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// 暂停 Schedule.
func (r *ScheduleService) Pause(ctx context.Context, scheduleID string, params SchedulePauseParams, opts ...option.RequestOption) (res *Schedule, err error) {
	if scheduleID == "" {
		return nil, fmt.Errorf("missing required schedule_id parameter")
	}
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("schedules/%s/pause", url.PathEscape(scheduleID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type SchedulePauseParams struct {
	// 有副作用请求可选的幂等键。
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r SchedulePauseParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// 运行 Schedule.
func (r *ScheduleService) Run(ctx context.Context, scheduleID string, params ScheduleRunParams, opts ...option.RequestOption) (res *ScheduleRun, err error) {
	if scheduleID == "" {
		return nil, fmt.Errorf("missing required schedule_id parameter")
	}
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("schedules/%s/run", url.PathEscape(scheduleID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type ScheduleRunParams struct {
	// 有副作用请求可选的幂等键。
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r ScheduleRunParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// 恢复 Schedule.
func (r *ScheduleService) Unpause(ctx context.Context, scheduleID string, params ScheduleUnpauseParams, opts ...option.RequestOption) (res *Schedule, err error) {
	if scheduleID == "" {
		return nil, fmt.Errorf("missing required schedule_id parameter")
	}
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("schedules/%s/unpause", url.PathEscape(scheduleID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type ScheduleUnpauseParams struct {
	// 有副作用请求可选的幂等键。
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r ScheduleUnpauseParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

type ScheduleArchiveManyResponse struct {
	// 本次从未归档状态变为已归档的 Schedule 数量；已经归档的目标不重复计数。
	ArchivedCount int64 `json:"archived_count"`
	JSON          struct {
		ArchivedCount respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

func (r ScheduleArchiveManyResponse) RawJSON() string { return r.JSON.raw }
func (r *ScheduleArchiveManyResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type Schedule struct {
	ID            string                      `json:"id"`
	IdentityID    string                      `json:"identity_id"`
	TemplateID    string                      `json:"template_id"`
	Name          string                      `json:"name"`
	Description   string                      `json:"description"`
	Status        string                      `json:"status"`
	InitialEvents []ScheduleInitialEventsItem `json:"initial_events"`
	Execution     ScheduleExecution           `json:"execution"`
	TriggerPolicy ScheduleTriggerPolicy       `json:"trigger_policy"`
	EnvironmentID string                      `json:"environment_id"`
	Sinks         []ScheduleSinksItem         `json:"sinks"`
	Metadata      map[string]any              `json:"metadata"`
	CreatedAt     time.Time                   `json:"created_at" format:"date-time"`
	UpdatedAt     time.Time                   `json:"updated_at" format:"date-time"`
	ArchivedAt    time.Time                   `json:"archived_at" format:"date-time"`
	PausedReason  SchedulePausedReason        `json:"paused_reason"`
	JSON          struct {
		ID            respjson.Field
		IdentityID    respjson.Field
		TemplateID    respjson.Field
		Name          respjson.Field
		Description   respjson.Field
		Status        respjson.Field
		InitialEvents respjson.Field
		Execution     respjson.Field
		TriggerPolicy respjson.Field
		EnvironmentID respjson.Field
		Sinks         respjson.Field
		Metadata      respjson.Field
		CreatedAt     respjson.Field
		UpdatedAt     respjson.Field
		ArchivedAt    respjson.Field
		PausedReason  respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

func (r Schedule) RawJSON() string                  { return r.JSON.raw }
func (r *Schedule) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type SchedulePausedReason struct {
	Type string `json:"type"`
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r SchedulePausedReason) RawJSON() string { return r.JSON.raw }
func (r *SchedulePausedReason) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ScheduleSinksItem struct {
	Type      string                  `json:"type"`
	ChannelID string                  `json:"channel_id"`
	Target    ScheduleSinksItemTarget `json:"target"`
	JSON      struct {
		Type        respjson.Field
		ChannelID   respjson.Field
		Target      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r ScheduleSinksItem) RawJSON() string                  { return r.JSON.raw }
func (r *ScheduleSinksItem) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type ScheduleSinksItemTarget struct {
	Type       string `json:"type"`
	ExternalID string `json:"external_id"`
	JSON       struct {
		Type        respjson.Field
		ExternalID  respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r ScheduleSinksItemTarget) RawJSON() string { return r.JSON.raw }
func (r *ScheduleSinksItemTarget) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ScheduleTriggerPolicy struct {
	Type           string      `json:"type"`
	Expression     string      `json:"expression"`
	Timezone       string      `json:"timezone"`
	UpcomingRunsAt []time.Time `json:"upcoming_runs_at"`
	JSON           struct {
		Type           respjson.Field
		Expression     respjson.Field
		Timezone       respjson.Field
		UpcomingRunsAt respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

func (r ScheduleTriggerPolicy) RawJSON() string { return r.JSON.raw }
func (r *ScheduleTriggerPolicy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ScheduleExecution struct {
	SessionMode       string `json:"session_mode"`
	MaxConcurrentRuns int64  `json:"max_concurrent_runs"`
	MaxAttempts       int64  `json:"max_attempts"`
	TimeoutMs         int64  `json:"timeout_ms"`
	JSON              struct {
		SessionMode       respjson.Field
		MaxConcurrentRuns respjson.Field
		MaxAttempts       respjson.Field
		TimeoutMs         respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

func (r ScheduleExecution) RawJSON() string                  { return r.JSON.raw }
func (r *ScheduleExecution) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type ScheduleInitialEventsItem struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	JSON    struct {
		Type        respjson.Field
		Content     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r ScheduleInitialEventsItem) RawJSON() string { return r.JSON.raw }
func (r *ScheduleInitialEventsItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
