// Qoder managed API definitions.
package managed

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"time"

	requestconfig "github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apiquery"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/pagination"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
)

// VaultService contains methods and other services that help with interacting
// with the Qoder Cloud Agents API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewVaultService] method instead.
type VaultService struct {
	Options     []option.RequestOption
	Credentials VaultCredentialService
}

// NewVaultService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewVaultService(opts ...option.RequestOption) (r VaultService) {
	r = VaultService{}
	r.Options = opts
	r.Credentials = NewVaultCredentialService(opts...)
	return
}

// Create Vault
func (r *VaultService) New(ctx context.Context, params VaultNewParams, opts ...option.RequestOption) (res *ManagedAgentsVault, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	path := "vaults"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Get Vault
func (r *VaultService) Get(ctx context.Context, vaultID string, query VaultGetParams, opts ...option.RequestOption) (res *ManagedAgentsVault, err error) {
	for _, v := range query.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if vaultID == "" {
		err = errors.New("missing required vault_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("vaults/%s", url.PathEscape(vaultID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List Vaults
func (r *VaultService) List(ctx context.Context, params VaultListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ManagedAgentsVault], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "vaults"
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, params, &res, opts...)
	if err != nil {
		return nil, err
	}
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

// List Vaults
func (r *VaultService) ListAutoPaging(ctx context.Context, params VaultListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ManagedAgentsVault] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// Delete Vault
func (r *VaultService) Delete(ctx context.Context, vaultID string, body VaultDeleteParams, opts ...option.RequestOption) (res *ManagedAgentsDeletedVault, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if vaultID == "" {
		err = errors.New("missing required vault_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("vaults/%s", url.PathEscape(vaultID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Archive Vault
func (r *VaultService) Archive(ctx context.Context, vaultID string, body VaultArchiveParams, opts ...option.RequestOption) (res *ManagedAgentsVault, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if vaultID == "" {
		err = errors.New("missing required vault_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("vaults/%s/archive", url.PathEscape(vaultID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Confirmation of a deleted vault.
type ManagedAgentsDeletedVault struct {
	// Unique identifier of the deleted vault.
	ID string `json:"id" api:"required"`
	// Any of "vault_deleted".
	Type ManagedAgentsDeletedVaultType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsDeletedVault) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsDeletedVault) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsDeletedVaultType string

const (
	ManagedAgentsDeletedVaultTypeVaultDeleted ManagedAgentsDeletedVaultType = "vault_deleted"
)

// A vault that stores credentials for use by agents during sessions.
type ManagedAgentsVault struct {
	// Unique identifier for the vault.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ArchivedAt time.Time `json:"archived_at" api:"required" format:"date-time"`
	// A timestamp in RFC 3339 format
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Human-readable name for the vault.
	DisplayName string `json:"display_name" api:"required"`
	// Arbitrary key-value metadata attached to the vault.
	Metadata map[string]string `json:"metadata" api:"required"`
	// Any of "vault".
	Type ManagedAgentsVaultType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ArchivedAt  respjson.Field
		CreatedAt   respjson.Field
		DisplayName respjson.Field
		Metadata    respjson.Field
		Type        respjson.Field
		UpdatedAt   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsVault) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsVault) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsVaultType string

const (
	ManagedAgentsVaultTypeVault ManagedAgentsVaultType = "vault"
)

type VaultNewParams struct {
	// Human-readable name for the vault. 1-255 characters.
	DisplayName string `json:"display_name" api:"required"`
	// Arbitrary key-value metadata to attach to the vault. Maximum 16 pairs, keys up
	// to 64 chars, values up to 512 chars.
	Metadata map[string]string `json:"metadata,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

func (r VaultNewParams) MarshalJSON() (data []byte, err error) {
	type shadow VaultNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type VaultGetParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

type VaultListParams struct {
	Name param.Opt[string] `query:"name,omitzero" json:"-"`

	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	AfterID  param.Opt[string] `query:"after_id,omitzero" json:"-"`

	// Whether to include archived vaults in the results.
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// Maximum number of vaults to return per page. Defaults to 20, maximum 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque pagination token from a previous `list_vaults` response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [VaultListParams]'s query parameters as `url.Values`.
func (r VaultListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type VaultDeleteParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

type VaultArchiveParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}
