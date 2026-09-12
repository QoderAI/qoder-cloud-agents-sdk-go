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

// List Credentials
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
	// Page size, maximum 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Pagination cursor (recommended), taken from `next_page` in the previous response;
	// mutually exclusive with `after_id` and `before_id`.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Cursor for the next page; mutually exclusive with `page` and `before_id`.
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// Cursor for the previous page; mutually exclusive with `page` and `after_id`.
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	// Search by `mcp_server_url`.
	Name param.Opt[string] `query:"name,omitzero" json:"-"`
	paramObj
}

func (r VaultCredentialListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// Create Credential
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
	// Credential authentication material; supports `static_bearer` and `mcp_oauth`. The
	// response only returns redacted, non-secret fields.
	Auth map[string]any `json:"auth" api:"required"`
	// Compatibility field; currently not persisted, and Forward always returns an empty
	// string.
	DisplayName param.Opt[string] `json:"display_name,omitzero"`
	// Metadata object. `created_by` is reserved and must not be sent (sending it returns
	// 400).
	Metadata map[string]any `json:"metadata,omitzero"`
	// Optional idempotency key for the create request. The same key may only be reused
	// for an identical request.
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

// Get Credential
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

// Delete Credential
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
	// Credential ID.
	ID string `json:"id"`
	// Always `vault_credential`.
	Type string `json:"type"`
	// ID of the owning Vault.
	VaultID string `json:"vault_id"`
	// Redacted authentication material.
	Auth VaultCredentialAuth `json:"auth"`
	// Currently always an empty string.
	DisplayName string `json:"display_name"`
	// Credential metadata.
	Metadata map[string]any `json:"metadata"`
	// Creation time in RFC 3339 format.
	CreatedAt time.Time `json:"created_at" format:"date-time"`
	// Last update time in RFC 3339 format.
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
