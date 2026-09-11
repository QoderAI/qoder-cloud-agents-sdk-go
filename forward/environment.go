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

// EnvironmentService provides Forward Environment operations.
type EnvironmentService struct {
	Options []option.RequestOption
}

func NewEnvironmentService(opts ...option.RequestOption) EnvironmentService {
	return EnvironmentService{Options: slices.Clone(opts)}
}

// 列出 Environment.
func (r *EnvironmentService) List(ctx context.Context, params EnvironmentListParams, opts ...option.RequestOption) (res *pagination.PageCursor[Environment], err error) {

	opts = slices.Concat(r.Options, opts)
	path := "environments"
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
func (r *EnvironmentService) ListAutoPaging(ctx context.Context, params EnvironmentListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[Environment] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

type EnvironmentListParams struct {
	// 分页大小，最大 100。
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// 分页游标（推荐使用），取值来自上一页响应的 `next_page`；与 `after_id`、`before_id` 互斥。
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// 向后翻页游标；与 `page`、`before_id` 互斥。
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// 向前翻页游标；与 `page`、`after_id` 互斥。
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	paramObj
}

func (r EnvironmentListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// 创建 Environment.
func (r *EnvironmentService) New(ctx context.Context, params EnvironmentNewParams, opts ...option.RequestOption) (res *Environment, err error) {
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := "environments"
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type EnvironmentNewParams struct {
	// Environment 名称；去除首尾空白后不能为空。
	Name string `json:"name" api:"required"`
	// 描述。
	Description param.Opt[string] `json:"description,omitzero"`
	// Environment 运行时配置对象；省略时默认使用 `{"type":"cloud"}`。显式传入时不能为 `null` 或空对象。字段详见 schemas。
	Config map[string]any `json:"config,omitzero"`
	// [Environment metadata](./schemas.md#environment-metadata)；省略时为 `{}`，显式传入时不能为 `null`。
	Metadata map[string]any `json:"metadata,omitzero"`
	// 建议创建请求携带。相同 key 和相同请求可安全重试。
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r EnvironmentNewParams) MarshalJSON() ([]byte, error) {
	type shadow EnvironmentNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *EnvironmentNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// 查询 Environment.
func (r *EnvironmentService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *Environment, err error) {
	if id == "" {
		return nil, fmt.Errorf("missing required id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("environments/%s", url.PathEscape(id))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// 修改 Environment.
func (r *EnvironmentService) Update(ctx context.Context, id string, params EnvironmentUpdateParams, opts ...option.RequestOption) (res *Environment, err error) {
	if id == "" {
		return nil, fmt.Errorf("missing required id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("environments/%s", url.PathEscape(id))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type EnvironmentUpdateParams struct {
	// 新名称。
	Name param.Opt[string] `json:"name,omitzero"`
	// 新描述。
	Description param.Opt[string] `json:"description,omitzero"`
	// 新配置；传入时不能为 `null`，显式 `null` 返回 400。字段详见 schemas。
	Config map[string]any `json:"config,omitzero"`
	// 要合并的 [Environment metadata](./schemas.md#environment-metadata)；传入时不能为 `null`，显式 `null` 返回 400。
	Metadata map[string]any `json:"metadata,omitzero"`
	paramObj
}

func (r EnvironmentUpdateParams) MarshalJSON() ([]byte, error) {
	type shadow EnvironmentUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *EnvironmentUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Archive an environment retained by historical sessions or tool calls.
func (r *EnvironmentService) Archive(ctx context.Context, id string, opts ...option.RequestOption) (res *Environment, err error) {
	if id == "" {
		return nil, fmt.Errorf("missing required id parameter")
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("environments/%s/archive", url.PathEscape(id))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// 删除 Environment.
func (r *EnvironmentService) Delete(ctx context.Context, id string, opts ...option.RequestOption) (err error) {
	if id == "" {
		return fmt.Errorf("missing required id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("environments/%s", url.PathEscape(id))
	err = convention.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

type Environment struct {
	ID          string            `json:"id"`
	Type        string            `json:"type"`
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Config      EnvironmentConfig `json:"config"`
	Metadata    map[string]any    `json:"metadata"`
	CreatedAt   time.Time         `json:"created_at" format:"date-time"`
	UpdatedAt   time.Time         `json:"updated_at" format:"date-time"`
	IdentityID  string            `json:"identity_id" api:"nullable"`
	ArchivedAt  time.Time         `json:"archived_at" api:"nullable" format:"date-time"`
	JSON        struct {
		ID          respjson.Field
		Type        respjson.Field
		Name        respjson.Field
		Description respjson.Field
		Config      respjson.Field
		Metadata    respjson.Field
		CreatedAt   respjson.Field
		UpdatedAt   respjson.Field
		IdentityID  respjson.Field
		ArchivedAt  respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r Environment) RawJSON() string                  { return r.JSON.raw }
func (r *Environment) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type EnvironmentConfig struct {
	Type     string                    `json:"type"`
	Packages EnvironmentConfigPackages `json:"packages"`
	JSON     struct {
		Type        respjson.Field
		Packages    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r EnvironmentConfig) RawJSON() string                  { return r.JSON.raw }
func (r *EnvironmentConfig) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type EnvironmentConfigPackages struct {
	Type  string   `json:"type"`
	Apt   []string `json:"apt"`
	Cargo []string `json:"cargo"`
	Gem   []string `json:"gem"`
	Go    []string `json:"go"`
	Npm   []string `json:"npm"`
	Pip   []string `json:"pip"`
	JSON  struct {
		Type        respjson.Field
		Apt         respjson.Field
		Cargo       respjson.Field
		Gem         respjson.Field
		Go          respjson.Field
		Npm         respjson.Field
		Pip         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r EnvironmentConfigPackages) RawJSON() string { return r.JSON.raw }
func (r *EnvironmentConfigPackages) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
