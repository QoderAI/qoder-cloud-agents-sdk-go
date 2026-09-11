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

// IdentityService provides Forward Identity operations.
type IdentityService struct {
	Options      []option.RequestOption
	Configs      IdentityConfigService
	MemoryStores IdentityMemoryStoreService
}

func NewIdentityService(opts ...option.RequestOption) IdentityService {
	return IdentityService{Options: slices.Clone(opts), Configs: NewIdentityConfigService(opts...), MemoryStores: NewIdentityMemoryStoreService(opts...)}
}

// 列出 Identities.
func (r *IdentityService) List(ctx context.Context, params IdentityListParams, opts ...option.RequestOption) (res *pagination.Page[Identity], err error) {

	opts = slices.Concat(r.Options, opts)
	path := "identities"
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
func (r *IdentityService) ListAutoPaging(ctx context.Context, params IdentityListParams, opts ...option.RequestOption) *pagination.PageAutoPager[Identity] {
	return pagination.NewPageAutoPager(r.List(ctx, params, opts...))
}

type IdentityListParams struct {
	// 按集成方终端用户 ID 过滤。
	ExternalID param.Opt[string] `query:"external_id,omitzero" json:"-"`
	// 按多个 Identity ID 过滤；支持逗号分隔或重复 query 参数，去重后最多 100 个。
	IdentityIDs []string `query:"identity_ids,omitzero" json:"-"`
	// 匹配 Identity ID、名称或外部 ID。
	Search param.Opt[string] `query:"search,omitzero" json:"-"`
	// 按是否启用过滤；非布尔值返回 400。
	Enabled param.Opt[bool] `query:"enabled,omitzero" json:"-"`
	// 分页大小，最大 100；超过上限时按最大值处理。
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// 向后翻页游标，不能与 `before_id` 同用。
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// 向前翻页游标，不能与 `after_id` 同用。
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	paramObj
}

func (r IdentityListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// 创建 Identity.
func (r *IdentityService) New(ctx context.Context, params IdentityNewParams, opts ...option.RequestOption) (res *Identity, err error) {
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := "identities"
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type IdentityNewParams struct {
	// 集成方系统中的终端用户 ID，不能是空串或纯空白。
	ExternalID string `json:"external_id" api:"required"`
	// 展示名，传入时不能是空串或纯空白。
	Name param.Opt[string] `json:"name,omitzero"`
	// 是否启用该 Identity，默认 `true`。
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// 业务元数据，建议最多 16 个 key。
	Metadata map[string]any `json:"metadata,omitzero"`
	// 有副作用请求可选的幂等键。
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r IdentityNewParams) MarshalJSON() ([]byte, error) {
	type shadow IdentityNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *IdentityNewParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

// 确保管理员 Identity.
func (r *IdentityService) EnsureAdmin(ctx context.Context, opts ...option.RequestOption) (res *Identity, err error) {

	opts = slices.Concat(r.Options, opts)
	path := "identities/admin/ensure"
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// 获取 Identity 统计.
func (r *IdentityService) Stats(ctx context.Context, opts ...option.RequestOption) (res *IdentityStats, err error) {

	opts = slices.Concat(r.Options, opts)
	path := "identities/stats"
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// 获取 Identity.
func (r *IdentityService) Get(ctx context.Context, identityID string, opts ...option.RequestOption) (res *Identity, err error) {
	if identityID == "" {
		return nil, fmt.Errorf("missing required identity_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("identities/%s", url.PathEscape(identityID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// 更新 Identity.
func (r *IdentityService) Update(ctx context.Context, identityID string, params IdentityUpdateParams, opts ...option.RequestOption) (res *Identity, err error) {
	if identityID == "" {
		return nil, fmt.Errorf("missing required identity_id parameter")
	}
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("identities/%s", url.PathEscape(identityID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type IdentityUpdateParams struct {
	// 替换原有终端用户 ID。
	ExternalID param.Opt[string] `json:"external_id,omitzero"`
	// 替换展示名。
	Name param.Opt[string] `json:"name,omitzero"`
	// 更新 Identity 是否可用。
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// 合并更新业务元数据；空字符串 value 删除对应 key。
	Metadata map[string]any `json:"metadata,omitzero"`
	// 有副作用请求可选的幂等键。
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r IdentityUpdateParams) MarshalJSON() ([]byte, error) {
	type shadow IdentityUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *IdentityUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// 删除 Identity.
func (r *IdentityService) Delete(ctx context.Context, identityID string, opts ...option.RequestOption) (res *DeletedIdentity, err error) {
	if identityID == "" {
		return nil, fmt.Errorf("missing required identity_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("identities/%s", url.PathEscape(identityID))
	err = convention.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// 列出 Identity 使用的 Template.
func (r *IdentityService) ListTemplates(ctx context.Context, identityID string, opts ...option.RequestOption) (res *IdentityListTemplatesResponse, err error) {
	if identityID == "" {
		return nil, fmt.Errorf("missing required identity_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("identities/%s/agents", url.PathEscape(identityID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// 清理 Identity.
func (r *IdentityService) Clear(ctx context.Context, identityID string, params IdentityClearParams, opts ...option.RequestOption) (res *IdentityClearResponse, err error) {
	if identityID == "" {
		return nil, fmt.Errorf("missing required identity_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("identities/%s/clear", url.PathEscape(identityID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type IdentityClearParams struct {
	// 清理原因，仅用于记录调用意图。
	Reason param.Opt[string] `json:"reason,omitzero"`
	paramObj
}

func (r IdentityClearParams) MarshalJSON() ([]byte, error) {
	type shadow IdentityClearParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *IdentityClearParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

// 停用 Identity.
func (r *IdentityService) Disable(ctx context.Context, identityID string, opts ...option.RequestOption) (res *Identity, err error) {
	if identityID == "" {
		return nil, fmt.Errorf("missing required identity_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("identities/%s/disable", url.PathEscape(identityID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// 启用 Identity.
func (r *IdentityService) Enable(ctx context.Context, identityID string, opts ...option.RequestOption) (res *Identity, err error) {
	if identityID == "" {
		return nil, fmt.Errorf("missing required identity_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("identities/%s/enable", url.PathEscape(identityID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type IdentityClearResponse struct {
	// 被清理的 Identity ID。
	IdentityID string `json:"identity_id"`
	// 清理状态，成功时为 `completed`。
	Status string `json:"status"`
	// Forward 侧清理完成时间，使用 RFC 3339 格式。
	CompletedAt time.Time                    `json:"completed_at" format:"date-time"`
	Summary     IdentityClearResponseSummary `json:"summary"`
	JSON        struct {
		IdentityID  respjson.Field
		Status      respjson.Field
		CompletedAt respjson.Field
		Summary     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r IdentityClearResponse) RawJSON() string { return r.JSON.raw }
func (r *IdentityClearResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type IdentityClearResponseSummary struct {
	IdentityConfigsArchived        int64 `json:"identity_configs_archived"`
	ResourceBindingsArchived       int64 `json:"resource_bindings_archived"`
	IdentityOwnedResourcesArchived int64 `json:"identity_owned_resources_archived"`
	SchedulesArchived              int64 `json:"schedules_archived"`
	ScheduleRunsSkipped            int64 `json:"schedule_runs_skipped"`
	SessionsArchived               int64 `json:"sessions_archived"`
	JSON                           struct {
		IdentityConfigsArchived        respjson.Field
		ResourceBindingsArchived       respjson.Field
		IdentityOwnedResourcesArchived respjson.Field
		SchedulesArchived              respjson.Field
		ScheduleRunsSkipped            respjson.Field
		SessionsArchived               respjson.Field
		ExtraFields                    map[string]respjson.Field
		raw                            string
	} `json:"-"`
}

func (r IdentityClearResponseSummary) RawJSON() string { return r.JSON.raw }
func (r *IdentityClearResponseSummary) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type IdentityTemplate struct {
	TemplateID   string    `json:"template_id"`
	TemplateName string    `json:"template_name"`
	SessionCount int64     `json:"session_count"`
	LastActiveAt time.Time `json:"last_active_at" format:"date-time"`
	JSON         struct {
		TemplateID   respjson.Field
		TemplateName respjson.Field
		SessionCount respjson.Field
		LastActiveAt respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

func (r IdentityTemplate) RawJSON() string                  { return r.JSON.raw }
func (r *IdentityTemplate) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type DeletedIdentity struct {
	// 被删除的 Identity ID。
	ID string `json:"id"`
	// 是否已完成删除。成功响应为 `true`。
	Deleted bool `json:"deleted"`
	JSON    struct {
		ID          respjson.Field
		Deleted     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r DeletedIdentity) RawJSON() string                  { return r.JSON.raw }
func (r *DeletedIdentity) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type IdentityStats struct {
	// 当前账号的 Identity 总数。
	TotalIdentities int64 `json:"total_identities"`
	// 最近活跃的 Identity 数量。
	ActiveIdentities int64 `json:"active_identities"`
	// 当前账号使用过的 Template 总数。
	TotalAgents int64 `json:"total_agents"`
	// 当前账号的 Session 总数。
	TotalSessions int64 `json:"total_sessions"`
	JSON          struct {
		TotalIdentities  respjson.Field
		ActiveIdentities respjson.Field
		TotalAgents      respjson.Field
		TotalSessions    respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

func (r IdentityStats) RawJSON() string                  { return r.JSON.raw }
func (r *IdentityStats) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type Identity struct {
	// Forward Identity ID，建议前缀 `idn_`。
	ID string `json:"id"`
	// 集成方系统中的终端用户 ID。
	ExternalID string `json:"external_id"`
	// Identity 展示名。
	Name string `json:"name"`
	// Identity 类型。普通 Identity 返回 `normal`。
	IdentityType string `json:"identity_type"`
	// 是否允许继续使用该 Identity。
	Enabled bool `json:"enabled"`
	// 业务元数据。
	Metadata map[string]any `json:"metadata"`
	// 创建时间，RFC 3339 格式。
	CreatedAt time.Time `json:"created_at" format:"date-time"`
	// 最近更新时间，RFC 3339 格式。
	UpdatedAt time.Time `json:"updated_at" format:"date-time"`
	JSON      struct {
		ID           respjson.Field
		ExternalID   respjson.Field
		Name         respjson.Field
		IdentityType respjson.Field
		Enabled      respjson.Field
		Metadata     respjson.Field
		CreatedAt    respjson.Field
		UpdatedAt    respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

func (r Identity) RawJSON() string                  { return r.JSON.raw }
func (r *Identity) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type IdentityListTemplatesResponse struct {
	Data []IdentityTemplate `json:"data"`
	JSON struct {
		Data        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r IdentityListTemplatesResponse) RawJSON() string { return r.JSON.raw }
func (r *IdentityListTemplatesResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
