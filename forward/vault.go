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

// VaultService provides Forward Vault operations.
type VaultService struct {
	Options     []option.RequestOption
	Credentials VaultCredentialService
}

func NewVaultService(opts ...option.RequestOption) VaultService {
	return VaultService{Options: slices.Clone(opts), Credentials: NewVaultCredentialService(opts...)}
}

// 列出 Vault.
func (r *VaultService) List(ctx context.Context, params VaultListParams, opts ...option.RequestOption) (res *pagination.PageCursor[Vault], err error) {

	opts = slices.Concat(r.Options, opts)
	path := "vaults"
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
func (r *VaultService) ListAutoPaging(ctx context.Context, params VaultListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[Vault] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

type VaultListParams struct {
	// 分页大小，最大 100。
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// 分页游标（推荐使用），取值来自上一页响应的 `next_page`；与 `after_id`、`before_id` 互斥。
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// 向后翻页游标；与 `page`、`before_id` 互斥。
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// 向前翻页游标；与 `page`、`after_id` 互斥。
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	// 按 `display_name` 搜索。
	Name param.Opt[string] `query:"name,omitzero" json:"-"`
	paramObj
}

func (r VaultListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// 创建 Vault.
func (r *VaultService) New(ctx context.Context, params VaultNewParams, opts ...option.RequestOption) (res *Vault, err error) {
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := "vaults"
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type VaultNewParams struct {
	// Vault 展示名。
	DisplayName string `json:"display_name" api:"required"`
	// 元数据对象；`created_by` 为保留字段，不可传入（传入返回 400）。
	Metadata map[string]any `json:"metadata,omitzero"`
	// 可选创建请求幂等键。传入时相同 key 只能用于相同请求；不传时不提供本地幂等重放保护。
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r VaultNewParams) MarshalJSON() ([]byte, error) {
	type shadow VaultNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultNewParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

// 查询 Vault.
func (r *VaultService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *Vault, err error) {
	if id == "" {
		return nil, fmt.Errorf("missing required id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("vaults/%s", url.PathEscape(id))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// 删除 Vault.
func (r *VaultService) Delete(ctx context.Context, id string, opts ...option.RequestOption) (err error) {
	if id == "" {
		return fmt.Errorf("missing required id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("vaults/%s", url.PathEscape(id))
	err = convention.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

type Vault struct {
	// Vault ID。
	ID string `json:"id"`
	// 固定为 `vault`。
	Type string `json:"type"`
	// Vault 展示名。
	DisplayName string `json:"display_name"`
	// Vault 元数据。
	Metadata map[string]any `json:"metadata"`
	// 创建时间，RFC 3339 格式。
	CreatedAt time.Time `json:"created_at" format:"date-time"`
	// 最后更新时间，RFC 3339 格式。
	UpdatedAt time.Time `json:"updated_at" format:"date-time"`
	// Forward 归属身份。
	IdentityID string `json:"identity_id" api:"nullable"`
	JSON       struct {
		ID          respjson.Field
		Type        respjson.Field
		DisplayName respjson.Field
		Metadata    respjson.Field
		CreatedAt   respjson.Field
		UpdatedAt   respjson.Field
		IdentityID  respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r Vault) RawJSON() string                  { return r.JSON.raw }
func (r *Vault) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
