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

// List Vaults
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
	// Select identity ownership with a PAT. Omit when using a SAT.
	IdentityID param.Opt[string] `query:"identity_id,omitzero" json:"-"`
	// Page size, maximum 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Pagination cursor (recommended), taken from `next_page` in the previous
	// response; mutually exclusive with `after_id` and `before_id`.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Cursor for the next page; mutually exclusive with `page` and `before_id`.
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// Cursor for the previous page; mutually exclusive with `page` and `after_id`.
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	// Search by `display_name`.
	Name param.Opt[string] `query:"name,omitzero" json:"-"`
	paramObj
}

func (r VaultListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// Create Vault
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
	// Select identity ownership with a PAT. Omit when using a SAT.
	IdentityID param.Opt[string] `query:"identity_id,omitzero" json:"-"`
	// Vault display name.
	DisplayName string `json:"display_name" api:"required"`
	// Metadata object; `created_by` is reserved and cannot be supplied (supplying it
	// returns 400).
	Metadata map[string]any `json:"metadata,omitzero"`
	// Optional idempotency key for the create request. A given key may only be reused
	// for an identical request; when omitted, no local replay protection applies.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r VaultNewParams) MarshalJSON() ([]byte, error) {
	type shadow VaultNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultNewParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

// Get Vault
func (r *VaultService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *Vault, err error) {
	if id == "" {
		return nil, fmt.Errorf("missing required id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("vaults/%s", url.PathEscape(id))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Delete Vault
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
	// Always `vault`.
	Type string `json:"type"`
	// Vault display name.
	DisplayName string `json:"display_name"`
	// Vault metadata.
	Metadata map[string]any `json:"metadata"`
	// Creation time in RFC 3339 format.
	CreatedAt time.Time `json:"created_at" format:"date-time"`
	// Last update time in RFC 3339 format.
	UpdatedAt time.Time `json:"updated_at" format:"date-time"`
	// Owning Forward identity.
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

func (r VaultNewParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// Search vaults by metadata and resource filters. Adds the required search Beta header.
func (r *VaultService) Search(ctx context.Context, params VaultSearchParams, opts ...option.RequestOption) (*pagination.PageCursor[Vault], error) {
	for _, beta := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", string(beta)))
	}
	return pagination.NewSearchPage[Vault](ctx, "vaults/search", params, slices.Concat(r.Options, opts)...)
}

// SearchAutoPaging iterates all matching vaults using opaque body cursors.
func (r *VaultService) SearchAutoPaging(ctx context.Context, params VaultSearchParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[Vault] {
	return pagination.NewPageCursorAutoPager(r.Search(ctx, params, opts...))
}

// VaultSearchParams contains JSON body filters for Search.
type VaultSearchParams struct {
	// Exact string metadata filters combined with AND (at most 16 entries).
	Metadata map[string]string `json:"metadata,omitzero"`
	// Page size; default 20, maximum 100.
	Limit param.Opt[int64] `json:"limit,omitzero"`
	// Opaque next_page cursor returned by the previous search response.
	Page param.Opt[string] `json:"page,omitzero"`
	// Include archived resources; defaults to false.
	IncludeArchived param.Opt[bool] `json:"include_archived,omitzero"`
	// Case-insensitive name substring filter.
	Name param.Opt[string] `json:"name,omitzero"`
	// Select Identity ownership with a PAT; omit with any SAT. Sent in the query only.
	IdentityID param.Opt[string] `query:"identity_id,omitzero" json:"-"`
	// Additional beta features; search-2026-08-31 is always included.
	Betas []string `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

func (r VaultSearchParams) MarshalJSON() ([]byte, error) {
	type shadow VaultSearchParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultSearchParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
func (r VaultSearchParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}
