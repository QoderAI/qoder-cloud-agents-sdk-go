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

// SessionService provides Forward Session operations.
type SessionService struct {
	Options   []option.RequestOption
	Events    SessionEventService
	Resources SessionResourceService
	Threads   SessionThreadService
}

func NewSessionService(opts ...option.RequestOption) SessionService {
	return SessionService{Options: slices.Clone(opts), Events: NewSessionEventService(opts...), Resources: NewSessionResourceService(opts...), Threads: NewSessionThreadService(opts...)}
}

// 列出 Sessions.
func (r *SessionService) List(ctx context.Context, params SessionListParams, opts ...option.RequestOption) (res *pagination.Page[Session], err error) {

	opts = slices.Concat(r.Options, opts)
	path := "sessions"
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
func (r *SessionService) ListAutoPaging(ctx context.Context, params SessionListParams, opts ...option.RequestOption) *pagination.PageAutoPager[Session] {
	return pagination.NewPageAutoPager(r.List(ctx, params, opts...))
}

type SessionListParams struct {
	// 按一个或多个 Identity ID 过滤，支持逗号分隔。
	IdentityIDs []string `query:"identity_ids,omitzero" json:"-"`
	// 按 Forward Template ID 过滤。
	TemplateID param.Opt[string] `query:"template_id,omitzero" json:"-"`
	// 按 `api`、`im`、`schedule` 或 `batch` 过滤。
	SourceType param.Opt[string] `query:"source_type,omitzero" json:"-"`
	// 创建时间严格大于该 RFC 3339 时间。
	CreatedAtGt param.Opt[time.Time] `query:"created_at[gt],omitzero" json:"-" format:"date-time"`
	// 创建时间大于等于该 RFC 3339 时间。
	CreatedAtGte param.Opt[time.Time] `query:"created_at[gte],omitzero" json:"-" format:"date-time"`
	// 创建时间严格小于该 RFC 3339 时间。
	CreatedAtLt param.Opt[time.Time] `query:"created_at[lt],omitzero" json:"-" format:"date-time"`
	// 创建时间小于等于该 RFC 3339 时间。
	CreatedAtLte param.Opt[time.Time] `query:"created_at[lte],omitzero" json:"-" format:"date-time"`
	// 更新时间严格大于该 RFC 3339 时间。
	UpdatedAtGt param.Opt[time.Time] `query:"updated_at[gt],omitzero" json:"-" format:"date-time"`
	// 更新时间大于等于该 RFC 3339 时间。
	UpdatedAtGte param.Opt[time.Time] `query:"updated_at[gte],omitzero" json:"-" format:"date-time"`
	// 更新时间严格小于该 RFC 3339 时间。
	UpdatedAtLt param.Opt[time.Time] `query:"updated_at[lt],omitzero" json:"-" format:"date-time"`
	// 更新时间小于等于该 RFC 3339 时间。
	UpdatedAtLte param.Opt[time.Time] `query:"updated_at[lte],omitzero" json:"-" format:"date-time"`
	// 分页大小，最大 100。
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// 向后翻页游标，传入上一页响应的 `last_id`。
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// 向前翻页游标，传入当前页响应的 `first_id`。
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	// 创建时间排序方向：`desc` 或 `asc`。
	Order param.Opt[string] `query:"order,omitzero" json:"-"`
	// 是否包含已归档 Session。
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	paramObj
}

func (r SessionListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// 创建 Session.
func (r *SessionService) New(ctx context.Context, params SessionNewParams, opts ...option.RequestOption) (res *Session, err error) {
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := "sessions"
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type SessionNewParams struct {
	// Forward Identity ID。
	IdentityID string `json:"identity_id" api:"required"`
	// Forward Template ID。
	TemplateID string `json:"template_id" api:"required"`
	// Session 标题。
	Title param.Opt[string] `json:"title,omitzero"`
	// 业务元数据。
	Metadata  map[string]any              `json:"metadata,omitzero"`
	Config    SessionNewParamsConfigParam `json:"config,omitzero"`
	Resources []SessionResourceSpecParam  `json:"resources,omitzero"`
	// 有副作用请求可选的幂等键。
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r SessionNewParams) MarshalJSON() ([]byte, error) {
	type shadow SessionNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionNewParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type SessionNewParamsConfigParam struct {
	EnvironmentVariables map[string]any `json:"environment_variables,omitzero"`
	paramObj
}

func (r SessionNewParamsConfigParam) MarshalJSON() ([]byte, error) {
	type shadow SessionNewParamsConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionNewParamsConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// 获取 Session.
func (r *SessionService) Get(ctx context.Context, sessionID string, opts ...option.RequestOption) (res *Session, err error) {
	if sessionID == "" {
		return nil, fmt.Errorf("missing required session_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("sessions/%s", url.PathEscape(sessionID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// 更新 Session.
func (r *SessionService) Update(ctx context.Context, sessionID string, params SessionUpdateParams, opts ...option.RequestOption) (res *Session, err error) {
	if sessionID == "" {
		return nil, fmt.Errorf("missing required session_id parameter")
	}
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("sessions/%s", url.PathEscape(sessionID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type SessionUpdateParams struct {
	// 新的 Session 标题。
	Title param.Opt[string] `json:"title,omitzero"`
	// metadata merge patch；传入的 key 覆盖已有 key，未出现的 key 保留。
	Metadata map[string]any                 `json:"metadata,omitzero"`
	Config   SessionUpdateParamsConfigParam `json:"config,omitzero"`
	// 有副作用请求可选的幂等键。
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r SessionUpdateParams) MarshalJSON() ([]byte, error) {
	type shadow SessionUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionUpdateParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type SessionUpdateParamsConfigParam struct {
	EnvironmentVariables map[string]any `json:"environment_variables,omitzero"`
	paramObj
}

func (r SessionUpdateParamsConfigParam) MarshalJSON() ([]byte, error) {
	type shadow SessionUpdateParamsConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionUpdateParamsConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// 归档 Session.
func (r *SessionService) Archive(ctx context.Context, sessionID string, params SessionArchiveParams, opts ...option.RequestOption) (res *Session, err error) {
	if sessionID == "" {
		return nil, fmt.Errorf("missing required session_id parameter")
	}
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("sessions/%s/archive", url.PathEscape(sessionID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type SessionArchiveParams struct {
	// 有副作用请求可选的幂等键。
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r SessionArchiveParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// 取消当前 Turn.
func (r *SessionService) Cancel(ctx context.Context, sessionID string, params SessionCancelParams, opts ...option.RequestOption) (res *Session, err error) {
	if sessionID == "" {
		return nil, fmt.Errorf("missing required session_id parameter")
	}
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("sessions/%s/cancel", url.PathEscape(sessionID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type SessionCancelParams struct {
	// 有副作用请求可选的幂等键。
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r SessionCancelParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

type Session struct {
	// Session ID。
	ID string `json:"id"`
	// 固定为 `session`。
	Type string `json:"type"`
	// Forward Identity ID。
	IdentityID string `json:"identity_id"`
	// Template 摘要。
	Template SessionTemplate `json:"template"`
	// Session 来源，直接 API 创建为 `api`。
	SourceType string `json:"source_type"`
	// `idle`、`running`、`rescheduling`、`canceling` 或 `terminated`。
	Status     string            `json:"status"`
	Title      string            `json:"title"`
	Metadata   map[string]any    `json:"metadata"`
	Config     SessionConfig     `json:"config"`
	Resources  []SessionResource `json:"resources"`
	Stats      SessionStats      `json:"stats"`
	Usage      SessionUsage      `json:"usage"`
	ArchivedAt time.Time         `json:"archived_at" api:"nullable" format:"date-time"`
	CreatedAt  time.Time         `json:"created_at" format:"date-time"`
	UpdatedAt  time.Time         `json:"updated_at" format:"date-time"`
	JSON       struct {
		ID          respjson.Field
		Type        respjson.Field
		IdentityID  respjson.Field
		Template    respjson.Field
		SourceType  respjson.Field
		Status      respjson.Field
		Title       respjson.Field
		Metadata    respjson.Field
		Config      respjson.Field
		Resources   respjson.Field
		Stats       respjson.Field
		Usage       respjson.Field
		ArchivedAt  respjson.Field
		CreatedAt   respjson.Field
		UpdatedAt   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r Session) RawJSON() string                  { return r.JSON.raw }
func (r *Session) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type SessionUsage struct {
	TotalCredits float64 `json:"total_credits"`
	JSON         struct {
		TotalCredits respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

func (r SessionUsage) RawJSON() string                  { return r.JSON.raw }
func (r *SessionUsage) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type SessionStats struct {
	ActiveSeconds   int64 `json:"active_seconds"`
	DurationSeconds int64 `json:"duration_seconds"`
	JSON            struct {
		ActiveSeconds   respjson.Field
		DurationSeconds respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

func (r SessionStats) RawJSON() string                  { return r.JSON.raw }
func (r *SessionStats) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type SessionConfig struct {
	EnvironmentVariables map[string]string `json:"environment_variables"`
	JSON                 struct {
		EnvironmentVariables respjson.Field
		ExtraFields          map[string]respjson.Field
		raw                  string
	} `json:"-"`
}

func (r SessionConfig) RawJSON() string                  { return r.JSON.raw }
func (r *SessionConfig) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type SessionTemplate struct {
	ID      string      `json:"id"`
	Type    string      `json:"type"`
	Name    string      `json:"name"`
	Model   ModelConfig `json:"model"`
	Version int64       `json:"version"`
	JSON    struct {
		ID          respjson.Field
		Type        respjson.Field
		Name        respjson.Field
		Model       respjson.Field
		Version     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r SessionTemplate) RawJSON() string                  { return r.JSON.raw }
func (r *SessionTemplate) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
