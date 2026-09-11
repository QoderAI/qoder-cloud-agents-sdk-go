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

// VaultCredentialService provides Forward VaultCredential operations.
type VaultCredentialService struct {
	Options []option.RequestOption
}

func NewVaultCredentialService(opts ...option.RequestOption) VaultCredentialService {
	return VaultCredentialService{Options: slices.Clone(opts)}
}

// 列出 Credential.
func (r *VaultCredentialService) List(ctx context.Context, id string, params VaultCredentialListParams, opts ...option.RequestOption) (res *pagination.PageCursor[VaultCredential], err error) {
	if id == "" {
		return nil, fmt.Errorf("missing required id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("vaults/%s/credentials", url.PathEscape(id))
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
func (r *VaultCredentialService) ListAutoPaging(ctx context.Context, id string, params VaultCredentialListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[VaultCredential] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, id, params, opts...))
}

type VaultCredentialListParams struct {
	// 分页大小，最大 100。
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// 分页游标（推荐使用），取值来自上一页响应的 `next_page`；与 `after_id`、`before_id` 互斥。
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// 向后翻页游标；与 `page`、`before_id` 互斥。
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// 向前翻页游标；与 `page`、`after_id` 互斥。
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	// 按 `mcp_server_url` 搜索。
	Name param.Opt[string] `query:"name,omitzero" json:"-"`
	paramObj
}

func (r VaultCredentialListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// 创建 Credential.
func (r *VaultCredentialService) New(ctx context.Context, id string, params VaultCredentialNewParams, opts ...option.RequestOption) (res *VaultCredential, err error) {
	if id == "" {
		return nil, fmt.Errorf("missing required id parameter")
	}
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("vaults/%s/credentials", url.PathEscape(id))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type VaultCredentialNewParams struct {
	// Credential 认证信息，支持 `static_bearer`、`mcp_oauth`；响应只返回脱敏后的非密文字段。
	Auth map[string]any `json:"auth" api:"required"`
	// 兼容字段；当前不持久化，Forward 响应固定为空字符串。
	DisplayName param.Opt[string] `json:"display_name,omitzero"`
	// 元数据对象；`created_by` 为保留字段，不可传入（传入返回 400）。
	Metadata map[string]any `json:"metadata,omitzero"`
	// 可选创建请求幂等键。相同 key 只能用于相同请求。
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r VaultCredentialNewParams) MarshalJSON() ([]byte, error) {
	type shadow VaultCredentialNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultCredentialNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// 查询 Credential.
func (r *VaultCredentialService) Get(ctx context.Context, id string, credID string, opts ...option.RequestOption) (res *VaultCredential, err error) {
	if id == "" {
		return nil, fmt.Errorf("missing required id parameter")
	}
	if credID == "" {
		return nil, fmt.Errorf("missing required cred_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("vaults/%s/credentials/%s", url.PathEscape(id), url.PathEscape(credID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// 删除 Credential.
func (r *VaultCredentialService) Delete(ctx context.Context, id string, credID string, opts ...option.RequestOption) (err error) {
	if id == "" {
		return fmt.Errorf("missing required id parameter")
	}
	if credID == "" {
		return fmt.Errorf("missing required cred_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("vaults/%s/credentials/%s", url.PathEscape(id), url.PathEscape(credID))
	err = convention.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

type VaultCredential struct {
	// Credential ID。
	ID string `json:"id"`
	// 固定为 `vault_credential`。
	Type string `json:"type"`
	// 所属 Vault ID。
	VaultID string `json:"vault_id"`
	// 脱敏后的认证信息。
	Auth VaultCredentialAuth `json:"auth"`
	// 当前固定为空字符串。
	DisplayName string `json:"display_name"`
	// Credential 元数据。
	Metadata map[string]any `json:"metadata"`
	// 创建时间，RFC 3339 格式。
	CreatedAt time.Time `json:"created_at" format:"date-time"`
	// 最后更新时间，RFC 3339 格式。
	UpdatedAt time.Time `json:"updated_at" format:"date-time"`
	JSON      struct {
		ID          respjson.Field
		Type        respjson.Field
		VaultID     respjson.Field
		Auth        respjson.Field
		DisplayName respjson.Field
		Metadata    respjson.Field
		CreatedAt   respjson.Field
		UpdatedAt   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r VaultCredential) RawJSON() string                  { return r.JSON.raw }
func (r *VaultCredential) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type VaultCredentialAuth struct {
	Type         string `json:"type"`
	MCPServerURL string `json:"mcp_server_url"`
	JSON         struct {
		Type         respjson.Field
		MCPServerURL respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

func (r VaultCredentialAuth) RawJSON() string                  { return r.JSON.raw }
func (r *VaultCredentialAuth) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
