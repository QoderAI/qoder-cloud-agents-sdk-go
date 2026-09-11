// Qoder managed API definitions.
package managed

import (
	"context"
	"encoding/json"
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

// VaultCredentialService contains methods and other services that help with
// interacting with the Qoder Cloud Agents API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewVaultCredentialService] method instead.
type VaultCredentialService struct {
	Options []option.RequestOption
}

// NewVaultCredentialService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewVaultCredentialService(opts ...option.RequestOption) (r VaultCredentialService) {
	r = VaultCredentialService{}
	r.Options = opts
	return
}

// Create Credential
func (r *VaultCredentialService) New(ctx context.Context, vaultID string, params VaultCredentialNewParams, opts ...option.RequestOption) (res *ManagedAgentsCredential, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if vaultID == "" {
		err = errors.New("missing required vault_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("vaults/%s/credentials", url.PathEscape(vaultID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Get Credential
func (r *VaultCredentialService) Get(ctx context.Context, credentialID string, params VaultCredentialGetParams, opts ...option.RequestOption) (res *ManagedAgentsCredential, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if params.VaultID == "" {
		err = errors.New("missing required vault_id parameter")
		return nil, err
	}
	if credentialID == "" {
		err = errors.New("missing required credential_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("vaults/%s/credentials/%s", url.PathEscape(params.VaultID), url.PathEscape(credentialID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update Credential
func (r *VaultCredentialService) Update(ctx context.Context, credentialID string, params VaultCredentialUpdateParams, opts ...option.RequestOption) (res *ManagedAgentsCredential, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if params.VaultID == "" {
		err = errors.New("missing required vault_id parameter")
		return nil, err
	}
	if credentialID == "" {
		err = errors.New("missing required credential_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("vaults/%s/credentials/%s", url.PathEscape(params.VaultID), url.PathEscape(credentialID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List Credentials
func (r *VaultCredentialService) List(ctx context.Context, vaultID string, params VaultCredentialListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ManagedAgentsCredential], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if vaultID == "" {
		err = errors.New("missing required vault_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("vaults/%s/credentials", url.PathEscape(vaultID))
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

// List Credentials
func (r *VaultCredentialService) ListAutoPaging(ctx context.Context, vaultID string, params VaultCredentialListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ManagedAgentsCredential] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, vaultID, params, opts...))
}

// Delete Credential
func (r *VaultCredentialService) Delete(ctx context.Context, credentialID string, params VaultCredentialDeleteParams, opts ...option.RequestOption) (res *ManagedAgentsDeletedCredential, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if params.VaultID == "" {
		err = errors.New("missing required vault_id parameter")
		return nil, err
	}
	if credentialID == "" {
		err = errors.New("missing required credential_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("vaults/%s/credentials/%s", url.PathEscape(params.VaultID), url.PathEscape(credentialID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Archive Credential
func (r *VaultCredentialService) Archive(ctx context.Context, credentialID string, params VaultCredentialArchiveParams, opts ...option.RequestOption) (res *ManagedAgentsCredential, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if params.VaultID == "" {
		err = errors.New("missing required vault_id parameter")
		return nil, err
	}
	if credentialID == "" {
		err = errors.New("missing required credential_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("vaults/%s/credentials/%s/archive", url.PathEscape(params.VaultID), url.PathEscape(credentialID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Validate Credential
func (r *VaultCredentialService) MCPOAuthValidate(ctx context.Context, credentialID string, params VaultCredentialMCPOAuthValidateParams, opts ...option.RequestOption) (res *ManagedAgentsCredentialValidation, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if params.VaultID == "" {
		err = errors.New("missing required vault_id parameter")
		return nil, err
	}
	if credentialID == "" {
		err = errors.New("missing required credential_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("vaults/%s/credentials/%s/mcp_oauth_validate", url.PathEscape(params.VaultID), url.PathEscape(credentialID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// A credential stored in a vault. Sensitive fields are never returned in
// responses.
type ManagedAgentsCredential struct {
	// Unique identifier for the credential.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ArchivedAt time.Time `json:"archived_at" api:"required" format:"date-time"`
	// Authentication details for a credential.
	Auth ManagedAgentsCredentialAuthUnion `json:"auth" api:"required"`
	// A timestamp in RFC 3339 format
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Arbitrary key-value metadata attached to the credential.
	Metadata map[string]string `json:"metadata" api:"required"`
	// Any of "vault_credential".
	Type ManagedAgentsCredentialType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// Identifier of the vault this credential belongs to.
	VaultID string `json:"vault_id" api:"required"`
	// Human-readable name for the credential.
	DisplayName string `json:"display_name" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ArchivedAt  respjson.Field
		Auth        respjson.Field
		CreatedAt   respjson.Field
		Metadata    respjson.Field
		Type        respjson.Field
		UpdatedAt   respjson.Field
		VaultID     respjson.Field
		DisplayName respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsCredential) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsCredential) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsCredentialAuthUnion contains all possible properties and values
// from [ManagedAgentsMCPOAuthAuthResponse],
// [ManagedAgentsStaticBearerAuthResponse],
// [ManagedAgentsEnvironmentVariableAuthResponse].
//
// Use the [ManagedAgentsCredentialAuthUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsCredentialAuthUnion struct {
	MCPServerURL string `json:"mcp_server_url"`
	// Any of "mcp_oauth", "static_bearer", "environment_variable".
	Type string `json:"type"`
	// This field is from variant [ManagedAgentsMCPOAuthAuthResponse].
	ExpiresAt time.Time `json:"expires_at"`
	// This field is from variant [ManagedAgentsMCPOAuthAuthResponse].
	Refresh ManagedAgentsMCPOAuthRefreshResponse `json:"refresh"`
	// This field is from variant [ManagedAgentsEnvironmentVariableAuthResponse].
	InjectionLocation ManagedAgentsInjectionLocationResponse `json:"injection_location"`
	// This field is from variant [ManagedAgentsEnvironmentVariableAuthResponse].
	Networking ManagedAgentsEnvironmentVariableAuthResponseNetworkingUnion `json:"networking"`
	// This field is from variant [ManagedAgentsEnvironmentVariableAuthResponse].
	SecretName string `json:"secret_name"`
	JSON       struct {
		MCPServerURL      respjson.Field
		Type              respjson.Field
		ExpiresAt         respjson.Field
		Refresh           respjson.Field
		InjectionLocation respjson.Field
		Networking        respjson.Field
		SecretName        respjson.Field
		raw               string
	} `json:"-"`
}

// anyManagedAgentsCredentialAuth is implemented by each variant of
// [ManagedAgentsCredentialAuthUnion] to add type safety for the return type of
// [ManagedAgentsCredentialAuthUnion.AsAny]
type anyManagedAgentsCredentialAuth interface {
	implManagedAgentsCredentialAuthUnion()
}

func (ManagedAgentsMCPOAuthAuthResponse) implManagedAgentsCredentialAuthUnion()            {}
func (ManagedAgentsStaticBearerAuthResponse) implManagedAgentsCredentialAuthUnion()        {}
func (ManagedAgentsEnvironmentVariableAuthResponse) implManagedAgentsCredentialAuthUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsCredentialAuthUnion.AsAny().(type) {
//	case qoder.ManagedAgentsMCPOAuthAuthResponse:
//	case qoder.ManagedAgentsStaticBearerAuthResponse:
//	case qoder.ManagedAgentsEnvironmentVariableAuthResponse:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsCredentialAuthUnion) AsAny() anyManagedAgentsCredentialAuth {
	switch u.Type {
	case "mcp_oauth":
		return u.AsMCPOAuth()
	case "static_bearer":
		return u.AsStaticBearer()
	case "environment_variable":
		return u.AsEnvironmentVariable()
	}
	return nil
}

func (u ManagedAgentsCredentialAuthUnion) AsMCPOAuth() (v ManagedAgentsMCPOAuthAuthResponse) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsCredentialAuthUnion) AsStaticBearer() (v ManagedAgentsStaticBearerAuthResponse) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsCredentialAuthUnion) AsEnvironmentVariable() (v ManagedAgentsEnvironmentVariableAuthResponse) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsCredentialAuthUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsCredentialAuthUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsCredentialType string

const (
	ManagedAgentsCredentialTypeVaultCredential ManagedAgentsCredentialType = "vault_credential"
)

func ManagedAgentsCredentialNetworkingParamsOfUnrestricted(type_ ManagedAgentsUnrestrictedCredentialNetworkingParamsType) ManagedAgentsCredentialNetworkingParamsUnion {
	var unrestricted ManagedAgentsUnrestrictedCredentialNetworkingParams
	unrestricted.Type = type_
	return ManagedAgentsCredentialNetworkingParamsUnion{OfUnrestricted: &unrestricted}
}

func ManagedAgentsCredentialNetworkingParamsOfLimited(allowedHosts []string) ManagedAgentsCredentialNetworkingParamsUnion {
	var limited ManagedAgentsLimitedCredentialNetworkingParams
	limited.AllowedHosts = allowedHosts
	return ManagedAgentsCredentialNetworkingParamsUnion{OfLimited: &limited}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsCredentialNetworkingParamsUnion struct {
	OfUnrestricted *ManagedAgentsUnrestrictedCredentialNetworkingParams `json:",omitzero,inline"`
	OfLimited      *ManagedAgentsLimitedCredentialNetworkingParams      `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsCredentialNetworkingParamsUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfUnrestricted, u.OfLimited)
}
func (u *ManagedAgentsCredentialNetworkingParamsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsCredentialNetworkingParamsUnion) asAny() any {
	if !param.IsOmitted(u.OfUnrestricted) {
		return u.OfUnrestricted
	} else if !param.IsOmitted(u.OfLimited) {
		return u.OfLimited
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsCredentialNetworkingParamsUnion) GetAllowedHosts() []string {
	if vt := u.OfLimited; vt != nil {
		return vt.AllowedHosts
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsCredentialNetworkingParamsUnion) GetType() *string {
	if vt := u.OfUnrestricted; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfLimited; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ManagedAgentsCredentialNetworkingParamsUnion](
		"type",
		apijson.Discriminator[ManagedAgentsUnrestrictedCredentialNetworkingParams]("unrestricted"),
		apijson.Discriminator[ManagedAgentsLimitedCredentialNetworkingParams]("limited"),
	)
}

// Result of live-probing a credential against its configured MCP server.
type ManagedAgentsCredentialValidation struct {
	// Unique identifier of the credential that was validated.
	CredentialID string `json:"credential_id" api:"required"`
	// Whether the credential has a refresh token configured.
	HasRefreshToken bool `json:"has_refresh_token" api:"required"`
	// The failing step of an MCP validation probe.
	MCPProbe ManagedAgentsMCPProbe `json:"mcp_probe" api:"required"`
	// Outcome of a refresh-token exchange attempted during credential validation.
	Refresh ManagedAgentsRefreshObject `json:"refresh" api:"required"`
	// Overall verdict of a credential validation probe.
	//
	// Any of "valid", "invalid", "unknown".
	Status ManagedAgentsCredentialValidationStatus `json:"status" api:"required"`
	// Any of "vault_credential_validation".
	Type ManagedAgentsCredentialValidationType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	ValidatedAt time.Time `json:"validated_at" api:"required" format:"date-time"`
	// Identifier of the vault containing the credential.
	VaultID string `json:"vault_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CredentialID    respjson.Field
		HasRefreshToken respjson.Field
		MCPProbe        respjson.Field
		Refresh         respjson.Field
		Status          respjson.Field
		Type            respjson.Field
		ValidatedAt     respjson.Field
		VaultID         respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsCredentialValidation) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsCredentialValidation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsCredentialValidationType string

const (
	ManagedAgentsCredentialValidationTypeVaultCredentialValidation ManagedAgentsCredentialValidationType = "vault_credential_validation"
)

// Overall verdict of a credential validation probe.
type ManagedAgentsCredentialValidationStatus string

const (
	ManagedAgentsCredentialValidationStatusValid   ManagedAgentsCredentialValidationStatus = "valid"
	ManagedAgentsCredentialValidationStatusInvalid ManagedAgentsCredentialValidationStatus = "invalid"
	ManagedAgentsCredentialValidationStatusUnknown ManagedAgentsCredentialValidationStatus = "unknown"
)

// Confirmation of a deleted credential.
type ManagedAgentsDeletedCredential struct {
	// Unique identifier of the deleted credential.
	ID string `json:"id" api:"required"`
	// Any of "vault_credential_deleted".
	Type ManagedAgentsDeletedCredentialType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsDeletedCredential) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsDeletedCredential) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsDeletedCredentialType string

const (
	ManagedAgentsDeletedCredentialTypeVaultCredentialDeleted ManagedAgentsDeletedCredentialType = "vault_credential_deleted"
)

// Environment variable credential details. The secret value is never returned.
type ManagedAgentsEnvironmentVariableAuthResponse struct {
	// Where in the outbound request the secret value is substituted.
	InjectionLocation ManagedAgentsInjectionLocationResponse `json:"injection_location" api:"required"`
	// Outbound hosts the secret value is substituted on.
	Networking ManagedAgentsEnvironmentVariableAuthResponseNetworkingUnion `json:"networking" api:"required"`
	// Name of the environment variable.
	SecretName string `json:"secret_name" api:"required"`
	// Any of "environment_variable".
	Type ManagedAgentsEnvironmentVariableAuthResponseType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		InjectionLocation respjson.Field
		Networking        respjson.Field
		SecretName        respjson.Field
		Type              respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsEnvironmentVariableAuthResponse) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsEnvironmentVariableAuthResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsEnvironmentVariableAuthResponseNetworkingUnion contains all
// possible properties and values from
// [ManagedAgentsUnrestrictedCredentialNetworkingResponse],
// [ManagedAgentsLimitedCredentialNetworkingResponse].
//
// Use the [ManagedAgentsEnvironmentVariableAuthResponseNetworkingUnion.AsAny]
// method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsEnvironmentVariableAuthResponseNetworkingUnion struct {
	// Any of "unrestricted", "limited".
	Type string `json:"type"`
	// This field is from variant
	// [ManagedAgentsLimitedCredentialNetworkingResponse].
	AllowedHosts []string `json:"allowed_hosts"`
	JSON         struct {
		Type         respjson.Field
		AllowedHosts respjson.Field
		raw          string
	} `json:"-"`
}

// anyManagedAgentsEnvironmentVariableAuthResponseNetworking is implemented by
// each variant of
// [ManagedAgentsEnvironmentVariableAuthResponseNetworkingUnion] to add type
// safety for the return type of
// [ManagedAgentsEnvironmentVariableAuthResponseNetworkingUnion.AsAny]
type anyManagedAgentsEnvironmentVariableAuthResponseNetworking interface {
	implManagedAgentsEnvironmentVariableAuthResponseNetworkingUnion()
}

func (ManagedAgentsUnrestrictedCredentialNetworkingResponse) implManagedAgentsEnvironmentVariableAuthResponseNetworkingUnion() {
}
func (ManagedAgentsLimitedCredentialNetworkingResponse) implManagedAgentsEnvironmentVariableAuthResponseNetworkingUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsEnvironmentVariableAuthResponseNetworkingUnion.AsAny().(type) {
//	case qoder.ManagedAgentsUnrestrictedCredentialNetworkingResponse:
//	case qoder.ManagedAgentsLimitedCredentialNetworkingResponse:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsEnvironmentVariableAuthResponseNetworkingUnion) AsAny() anyManagedAgentsEnvironmentVariableAuthResponseNetworking {
	switch u.Type {
	case "unrestricted":
		return u.AsUnrestricted()
	case "limited":
		return u.AsLimited()
	}
	return nil
}

func (u ManagedAgentsEnvironmentVariableAuthResponseNetworkingUnion) AsUnrestricted() (v ManagedAgentsUnrestrictedCredentialNetworkingResponse) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsEnvironmentVariableAuthResponseNetworkingUnion) AsLimited() (v ManagedAgentsLimitedCredentialNetworkingResponse) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsEnvironmentVariableAuthResponseNetworkingUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *ManagedAgentsEnvironmentVariableAuthResponseNetworkingUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsEnvironmentVariableAuthResponseType string

const (
	ManagedAgentsEnvironmentVariableAuthResponseTypeEnvironmentVariable ManagedAgentsEnvironmentVariableAuthResponseType = "environment_variable"
)

// Parameters for creating an environment variable credential.
//
// The properties Networking, SecretName, SecretValue, Type are required.
type ManagedAgentsEnvironmentVariableCreateParams struct {
	// Outbound hosts the secret value is substituted on.
	Networking ManagedAgentsCredentialNetworkingParamsUnion `json:"networking,omitzero" api:"required"`
	// Name of the environment variable. Immutable after create.
	SecretName string `json:"secret_name" api:"required"`
	// Secret value. Write-only; never returned in responses.
	SecretValue string `json:"secret_value" api:"required"`
	// Any of "environment_variable".
	Type ManagedAgentsEnvironmentVariableCreateParamsType `json:"type,omitzero" api:"required"`
	// Where in the outbound request the secret value may be substituted.
	InjectionLocation ManagedAgentsInjectionLocationParams `json:"injection_location,omitzero"`
	paramObj
}

func (r ManagedAgentsEnvironmentVariableCreateParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsEnvironmentVariableCreateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsEnvironmentVariableCreateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsEnvironmentVariableCreateParamsType string

const (
	ManagedAgentsEnvironmentVariableCreateParamsTypeEnvironmentVariable ManagedAgentsEnvironmentVariableCreateParamsType = "environment_variable"
)

// Parameters for updating an environment variable credential. `secret_name` is
// immutable.
//
// The property Type is required.
type ManagedAgentsEnvironmentVariableUpdateParams struct {
	// Any of "environment_variable".
	Type ManagedAgentsEnvironmentVariableUpdateParamsType `json:"type,omitzero" api:"required"`
	// Updated secret value.
	SecretValue param.Opt[string] `json:"secret_value,omitzero"`
	// Updated injection location.
	InjectionLocation ManagedAgentsInjectionLocationUpdateParams `json:"injection_location,omitzero"`
	// Updated networking scope. Full replacement.
	Networking ManagedAgentsCredentialNetworkingParamsUnion `json:"networking,omitzero"`
	paramObj
}

func (r ManagedAgentsEnvironmentVariableUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsEnvironmentVariableUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsEnvironmentVariableUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsEnvironmentVariableUpdateParamsType string

const (
	ManagedAgentsEnvironmentVariableUpdateParamsTypeEnvironmentVariable ManagedAgentsEnvironmentVariableUpdateParamsType = "environment_variable"
)

// Where in the outbound request the secret value may be substituted.
type ManagedAgentsInjectionLocationParams struct {
	// Substitute when the placeholder appears in the request body.
	Body param.Opt[bool] `json:"body,omitzero"`
	// Substitute when the placeholder appears in a request header value.
	Header param.Opt[bool] `json:"header,omitzero"`
	paramObj
}

func (r ManagedAgentsInjectionLocationParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsInjectionLocationParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsInjectionLocationParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Where in the outbound request the secret value is substituted.
type ManagedAgentsInjectionLocationResponse struct {
	// Whether the placeholder is substituted in the request body.
	Body bool `json:"body" api:"required"`
	// Whether the placeholder is substituted in request header values.
	Header bool `json:"header" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Body        respjson.Field
		Header      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsInjectionLocationResponse) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsInjectionLocationResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Updated injection location.
type ManagedAgentsInjectionLocationUpdateParams struct {
	// Substitute when the placeholder appears in the request body.
	Body param.Opt[bool] `json:"body,omitzero"`
	// Substitute when the placeholder appears in a request header value.
	Header param.Opt[bool] `json:"header,omitzero"`
	paramObj
}

func (r ManagedAgentsInjectionLocationUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsInjectionLocationUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsInjectionLocationUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Substitute the secret only on requests to the listed hosts.
//
// The properties AllowedHosts, Type are required.
type ManagedAgentsLimitedCredentialNetworkingParams struct {
	// Hostnames on which the secret will be substituted. Each entry is a bare hostname
	// (`api.example.com`), an IPv4 address (`192.0.2.1`), or a `*.`-prefixed wildcard
	// (`*.example.com`). URLs, ports, paths, and IPv6 addresses are not accepted. At
	// most 16 entries.
	AllowedHosts []string `json:"allowed_hosts,omitzero" api:"required"`
	// Any of "limited".
	Type ManagedAgentsLimitedCredentialNetworkingParamsType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsLimitedCredentialNetworkingParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsLimitedCredentialNetworkingParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsLimitedCredentialNetworkingParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsLimitedCredentialNetworkingParamsType string

const (
	ManagedAgentsLimitedCredentialNetworkingParamsTypeLimited ManagedAgentsLimitedCredentialNetworkingParamsType = "limited"
)

// The secret is substituted only on requests to the listed hosts.
type ManagedAgentsLimitedCredentialNetworkingResponse struct {
	// Hostnames on which the secret will be substituted. An entry matches the request
	// host exactly; a `*.`-prefixed entry matches any subdomain of the named domain
	// but not the domain itself.
	AllowedHosts []string `json:"allowed_hosts" api:"required"`
	// Any of "limited".
	Type ManagedAgentsLimitedCredentialNetworkingResponseType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AllowedHosts respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsLimitedCredentialNetworkingResponse) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsLimitedCredentialNetworkingResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsLimitedCredentialNetworkingResponseType string

const (
	ManagedAgentsLimitedCredentialNetworkingResponseTypeLimited ManagedAgentsLimitedCredentialNetworkingResponseType = "limited"
)

// OAuth credential details for an MCP server.
type ManagedAgentsMCPOAuthAuthResponse struct {
	// URL of the MCP server this credential authenticates against.
	MCPServerURL string `json:"mcp_server_url" api:"required"`
	// Any of "mcp_oauth".
	Type ManagedAgentsMCPOAuthAuthResponseType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	ExpiresAt time.Time `json:"expires_at" api:"nullable" format:"date-time"`
	// OAuth refresh token configuration returned in credential responses.
	Refresh ManagedAgentsMCPOAuthRefreshResponse `json:"refresh" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		MCPServerURL respjson.Field
		Type         respjson.Field
		ExpiresAt    respjson.Field
		Refresh      respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsMCPOAuthAuthResponse) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsMCPOAuthAuthResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsMCPOAuthAuthResponseType string

const (
	ManagedAgentsMCPOAuthAuthResponseTypeMCPOAuth ManagedAgentsMCPOAuthAuthResponseType = "mcp_oauth"
)

// Parameters for creating an MCP OAuth credential.
//
// The properties AccessToken, MCPServerURL, Type are required.
type ManagedAgentsMCPOAuthCreateParams struct {
	// OAuth access token.
	AccessToken string `json:"access_token" api:"required"`
	// URL of the MCP server this credential authenticates against.
	MCPServerURL string `json:"mcp_server_url" api:"required"`
	// Any of "mcp_oauth".
	Type ManagedAgentsMCPOAuthCreateParamsType `json:"type,omitzero" api:"required"`
	// A timestamp in RFC 3339 format
	ExpiresAt param.Opt[time.Time] `json:"expires_at,omitzero" format:"date-time"`
	// OAuth refresh token parameters for creating a credential with refresh support.
	Refresh ManagedAgentsMCPOAuthRefreshParams `json:"refresh,omitzero"`
	paramObj
}

func (r ManagedAgentsMCPOAuthCreateParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsMCPOAuthCreateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsMCPOAuthCreateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsMCPOAuthCreateParamsType string

const (
	ManagedAgentsMCPOAuthCreateParamsTypeMCPOAuth ManagedAgentsMCPOAuthCreateParamsType = "mcp_oauth"
)

// OAuth refresh token parameters for creating a credential with refresh support.
//
// The properties ClientID, RefreshToken, TokenEndpoint, TokenEndpointAuth are
// required.
type ManagedAgentsMCPOAuthRefreshParams struct {
	// OAuth client ID.
	ClientID string `json:"client_id" api:"required"`
	// OAuth refresh token.
	RefreshToken string `json:"refresh_token" api:"required"`
	// Token endpoint URL used to refresh the access token.
	TokenEndpoint string `json:"token_endpoint" api:"required"`
	// Token endpoint requires no client authentication.
	TokenEndpointAuth ManagedAgentsMCPOAuthRefreshParamsTokenEndpointAuthUnion `json:"token_endpoint_auth,omitzero" api:"required"`
	// OAuth resource indicator.
	Resource param.Opt[string] `json:"resource,omitzero"`
	// OAuth scope for the refresh request.
	Scope param.Opt[string] `json:"scope,omitzero"`
	paramObj
}

func (r ManagedAgentsMCPOAuthRefreshParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsMCPOAuthRefreshParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsMCPOAuthRefreshParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsMCPOAuthRefreshParamsTokenEndpointAuthUnion struct {
	OfNone              *ManagedAgentsTokenEndpointAuthNoneParam  `json:",omitzero,inline"`
	OfClientSecretBasic *ManagedAgentsTokenEndpointAuthBasicParam `json:",omitzero,inline"`
	OfClientSecretPost  *ManagedAgentsTokenEndpointAuthPostParam  `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsMCPOAuthRefreshParamsTokenEndpointAuthUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfNone, u.OfClientSecretBasic, u.OfClientSecretPost)
}
func (u *ManagedAgentsMCPOAuthRefreshParamsTokenEndpointAuthUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsMCPOAuthRefreshParamsTokenEndpointAuthUnion) asAny() any {
	if !param.IsOmitted(u.OfNone) {
		return u.OfNone
	} else if !param.IsOmitted(u.OfClientSecretBasic) {
		return u.OfClientSecretBasic
	} else if !param.IsOmitted(u.OfClientSecretPost) {
		return u.OfClientSecretPost
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsMCPOAuthRefreshParamsTokenEndpointAuthUnion) GetType() *string {
	if vt := u.OfNone; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfClientSecretBasic; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfClientSecretPost; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsMCPOAuthRefreshParamsTokenEndpointAuthUnion) GetClientSecret() *string {
	if vt := u.OfClientSecretBasic; vt != nil {
		return (*string)(&vt.ClientSecret)
	} else if vt := u.OfClientSecretPost; vt != nil {
		return (*string)(&vt.ClientSecret)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ManagedAgentsMCPOAuthRefreshParamsTokenEndpointAuthUnion](
		"type",
		apijson.Discriminator[ManagedAgentsTokenEndpointAuthNoneParam]("none"),
		apijson.Discriminator[ManagedAgentsTokenEndpointAuthBasicParam]("client_secret_basic"),
		apijson.Discriminator[ManagedAgentsTokenEndpointAuthPostParam]("client_secret_post"),
	)
}

// OAuth refresh token configuration returned in credential responses.
type ManagedAgentsMCPOAuthRefreshResponse struct {
	// OAuth client ID.
	ClientID string `json:"client_id" api:"required"`
	// Token endpoint URL used to refresh the access token.
	TokenEndpoint string `json:"token_endpoint" api:"required"`
	// Token endpoint requires no client authentication.
	TokenEndpointAuth ManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuthUnion `json:"token_endpoint_auth" api:"required"`
	// OAuth resource indicator.
	Resource string `json:"resource" api:"nullable"`
	// OAuth scope for the refresh request.
	Scope string `json:"scope" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ClientID          respjson.Field
		TokenEndpoint     respjson.Field
		TokenEndpointAuth respjson.Field
		Resource          respjson.Field
		Scope             respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsMCPOAuthRefreshResponse) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsMCPOAuthRefreshResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuthUnion contains all
// possible properties and values from
// [ManagedAgentsTokenEndpointAuthNoneResponse],
// [ManagedAgentsTokenEndpointAuthBasicResponse],
// [ManagedAgentsTokenEndpointAuthPostResponse].
//
// Use the [ManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuthUnion.AsAny]
// method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuthUnion struct {
	// Any of "none", "client_secret_basic", "client_secret_post".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuth is implemented by
// each variant of [ManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuthUnion]
// to add type safety for the return type of
// [ManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuthUnion.AsAny]
type anyManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuth interface {
	implManagedAgentsMcpoAuthRefreshResponseTokenEndpointAuthUnion()
}

func (ManagedAgentsTokenEndpointAuthNoneResponse) implManagedAgentsMcpoAuthRefreshResponseTokenEndpointAuthUnion() {
}
func (ManagedAgentsTokenEndpointAuthBasicResponse) implManagedAgentsMcpoAuthRefreshResponseTokenEndpointAuthUnion() {
}
func (ManagedAgentsTokenEndpointAuthPostResponse) implManagedAgentsMcpoAuthRefreshResponseTokenEndpointAuthUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuthUnion.AsAny().(type) {
//	case qoder.ManagedAgentsTokenEndpointAuthNoneResponse:
//	case qoder.ManagedAgentsTokenEndpointAuthBasicResponse:
//	case qoder.ManagedAgentsTokenEndpointAuthPostResponse:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuthUnion) AsAny() anyManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuth {
	switch u.Type {
	case "none":
		return u.AsNone()
	case "client_secret_basic":
		return u.AsClientSecretBasic()
	case "client_secret_post":
		return u.AsClientSecretPost()
	}
	return nil
}

func (u ManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuthUnion) AsNone() (v ManagedAgentsTokenEndpointAuthNoneResponse) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuthUnion) AsClientSecretBasic() (v ManagedAgentsTokenEndpointAuthBasicResponse) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuthUnion) AsClientSecretPost() (v ManagedAgentsTokenEndpointAuthPostResponse) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuthUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *ManagedAgentsMCPOAuthRefreshResponseTokenEndpointAuthUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Parameters for updating OAuth refresh token configuration.
type ManagedAgentsMCPOAuthRefreshUpdateParams struct {
	// Updated OAuth refresh token.
	RefreshToken param.Opt[string] `json:"refresh_token,omitzero"`
	// Updated OAuth scope for the refresh request.
	Scope param.Opt[string] `json:"scope,omitzero"`
	// Updated HTTP Basic authentication parameters for the token endpoint.
	TokenEndpointAuth ManagedAgentsMCPOAuthRefreshUpdateParamsTokenEndpointAuthUnion `json:"token_endpoint_auth,omitzero"`
	paramObj
}

func (r ManagedAgentsMCPOAuthRefreshUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsMCPOAuthRefreshUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsMCPOAuthRefreshUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsMCPOAuthRefreshUpdateParamsTokenEndpointAuthUnion struct {
	OfClientSecretBasic *ManagedAgentsTokenEndpointAuthBasicUpdateParam `json:",omitzero,inline"`
	OfClientSecretPost  *ManagedAgentsTokenEndpointAuthPostUpdateParam  `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsMCPOAuthRefreshUpdateParamsTokenEndpointAuthUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfClientSecretBasic, u.OfClientSecretPost)
}
func (u *ManagedAgentsMCPOAuthRefreshUpdateParamsTokenEndpointAuthUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsMCPOAuthRefreshUpdateParamsTokenEndpointAuthUnion) asAny() any {
	if !param.IsOmitted(u.OfClientSecretBasic) {
		return u.OfClientSecretBasic
	} else if !param.IsOmitted(u.OfClientSecretPost) {
		return u.OfClientSecretPost
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsMCPOAuthRefreshUpdateParamsTokenEndpointAuthUnion) GetType() *string {
	if vt := u.OfClientSecretBasic; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfClientSecretPost; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsMCPOAuthRefreshUpdateParamsTokenEndpointAuthUnion) GetClientSecret() *string {
	if vt := u.OfClientSecretBasic; vt != nil && vt.ClientSecret.Valid() {
		return &vt.ClientSecret.Value
	} else if vt := u.OfClientSecretPost; vt != nil && vt.ClientSecret.Valid() {
		return &vt.ClientSecret.Value
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ManagedAgentsMCPOAuthRefreshUpdateParamsTokenEndpointAuthUnion](
		"type",
		apijson.Discriminator[ManagedAgentsTokenEndpointAuthBasicUpdateParam]("client_secret_basic"),
		apijson.Discriminator[ManagedAgentsTokenEndpointAuthPostUpdateParam]("client_secret_post"),
	)
}

// Parameters for updating an MCP OAuth credential. The `mcp_server_url` is
// immutable.
//
// The property Type is required.
type ManagedAgentsMCPOAuthUpdateParams struct {
	// Any of "mcp_oauth".
	Type ManagedAgentsMCPOAuthUpdateParamsType `json:"type,omitzero" api:"required"`
	// Updated OAuth access token.
	AccessToken param.Opt[string] `json:"access_token,omitzero"`
	// A timestamp in RFC 3339 format
	ExpiresAt param.Opt[time.Time] `json:"expires_at,omitzero" format:"date-time"`
	// Parameters for updating OAuth refresh token configuration.
	Refresh ManagedAgentsMCPOAuthRefreshUpdateParams `json:"refresh,omitzero"`
	paramObj
}

func (r ManagedAgentsMCPOAuthUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsMCPOAuthUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsMCPOAuthUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsMCPOAuthUpdateParamsType string

const (
	ManagedAgentsMCPOAuthUpdateParamsTypeMCPOAuth ManagedAgentsMCPOAuthUpdateParamsType = "mcp_oauth"
)

// The failing step of an MCP validation probe.
type ManagedAgentsMCPProbe struct {
	// An HTTP response captured during a credential validation probe.
	HTTPResponse ManagedAgentsRefreshHTTPResponse `json:"http_response" api:"required"`
	// The MCP method that failed (for example `initialize` or `tools/list`).
	Method string `json:"method" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		HTTPResponse respjson.Field
		Method       respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsMCPProbe) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsMCPProbe) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// An HTTP response captured during a credential validation probe.
type ManagedAgentsRefreshHTTPResponse struct {
	// Response body. May be truncated and has sensitive values scrubbed.
	Body string `json:"body" api:"required"`
	// Whether `body` was truncated.
	BodyTruncated bool `json:"body_truncated" api:"required"`
	// Value of the `Content-Type` response header.
	ContentType string `json:"content_type" api:"required"`
	// HTTP status code.
	StatusCode int64 `json:"status_code" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Body          respjson.Field
		BodyTruncated respjson.Field
		ContentType   respjson.Field
		StatusCode    respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsRefreshHTTPResponse) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsRefreshHTTPResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Outcome of a refresh-token exchange attempted during credential validation.
type ManagedAgentsRefreshObject struct {
	// An HTTP response captured during a credential validation probe.
	HTTPResponse ManagedAgentsRefreshHTTPResponse `json:"http_response" api:"required"`
	// Outcome of a refresh-token exchange attempted during credential validation.
	//
	// Any of "succeeded", "failed", "connect_error", "no_refresh_token".
	Status ManagedAgentsRefreshObjectStatus `json:"status" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		HTTPResponse respjson.Field
		Status       respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsRefreshObject) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsRefreshObject) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Outcome of a refresh-token exchange attempted during credential validation.
type ManagedAgentsRefreshObjectStatus string

const (
	ManagedAgentsRefreshObjectStatusSucceeded      ManagedAgentsRefreshObjectStatus = "succeeded"
	ManagedAgentsRefreshObjectStatusFailed         ManagedAgentsRefreshObjectStatus = "failed"
	ManagedAgentsRefreshObjectStatusConnectError   ManagedAgentsRefreshObjectStatus = "connect_error"
	ManagedAgentsRefreshObjectStatusNoRefreshToken ManagedAgentsRefreshObjectStatus = "no_refresh_token"
)

// Static bearer token credential details for an MCP server.
type ManagedAgentsStaticBearerAuthResponse struct {
	// URL of the MCP server this credential authenticates against.
	MCPServerURL string `json:"mcp_server_url" api:"required"`
	// Any of "static_bearer".
	Type ManagedAgentsStaticBearerAuthResponseType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		MCPServerURL respjson.Field
		Type         respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsStaticBearerAuthResponse) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsStaticBearerAuthResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsStaticBearerAuthResponseType string

const (
	ManagedAgentsStaticBearerAuthResponseTypeStaticBearer ManagedAgentsStaticBearerAuthResponseType = "static_bearer"
)

// Parameters for creating a static bearer token credential.
//
// The properties Token, MCPServerURL, Type are required.
type ManagedAgentsStaticBearerCreateParams struct {
	// Static bearer token value.
	Token string `json:"token" api:"required"`
	// URL of the MCP server this credential authenticates against.
	MCPServerURL string `json:"mcp_server_url" api:"required"`
	// Any of "static_bearer".
	Type ManagedAgentsStaticBearerCreateParamsType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsStaticBearerCreateParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsStaticBearerCreateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsStaticBearerCreateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsStaticBearerCreateParamsType string

const (
	ManagedAgentsStaticBearerCreateParamsTypeStaticBearer ManagedAgentsStaticBearerCreateParamsType = "static_bearer"
)

// Parameters for updating a static bearer token credential. The `mcp_server_url`
// is immutable.
//
// The property Type is required.
type ManagedAgentsStaticBearerUpdateParams struct {
	// Any of "static_bearer".
	Type ManagedAgentsStaticBearerUpdateParamsType `json:"type,omitzero" api:"required"`
	// Updated static bearer token value.
	Token param.Opt[string] `json:"token,omitzero"`
	paramObj
}

func (r ManagedAgentsStaticBearerUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsStaticBearerUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsStaticBearerUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsStaticBearerUpdateParamsType string

const (
	ManagedAgentsStaticBearerUpdateParamsTypeStaticBearer ManagedAgentsStaticBearerUpdateParamsType = "static_bearer"
)

// Token endpoint uses HTTP Basic authentication with client credentials.
//
// The properties ClientSecret, Type are required.
type ManagedAgentsTokenEndpointAuthBasicParam struct {
	// OAuth client secret.
	ClientSecret string `json:"client_secret" api:"required"`
	// Any of "client_secret_basic".
	Type ManagedAgentsTokenEndpointAuthBasicParamType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsTokenEndpointAuthBasicParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsTokenEndpointAuthBasicParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsTokenEndpointAuthBasicParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsTokenEndpointAuthBasicParamType string

const (
	ManagedAgentsTokenEndpointAuthBasicParamTypeClientSecretBasic ManagedAgentsTokenEndpointAuthBasicParamType = "client_secret_basic"
)

// Token endpoint uses HTTP Basic authentication with client credentials.
type ManagedAgentsTokenEndpointAuthBasicResponse struct {
	// Any of "client_secret_basic".
	Type ManagedAgentsTokenEndpointAuthBasicResponseType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsTokenEndpointAuthBasicResponse) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsTokenEndpointAuthBasicResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsTokenEndpointAuthBasicResponseType string

const (
	ManagedAgentsTokenEndpointAuthBasicResponseTypeClientSecretBasic ManagedAgentsTokenEndpointAuthBasicResponseType = "client_secret_basic"
)

// Updated HTTP Basic authentication parameters for the token endpoint.
//
// The property Type is required.
type ManagedAgentsTokenEndpointAuthBasicUpdateParam struct {
	// Any of "client_secret_basic".
	Type ManagedAgentsTokenEndpointAuthBasicUpdateParamType `json:"type,omitzero" api:"required"`
	// Updated OAuth client secret.
	ClientSecret param.Opt[string] `json:"client_secret,omitzero"`
	paramObj
}

func (r ManagedAgentsTokenEndpointAuthBasicUpdateParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsTokenEndpointAuthBasicUpdateParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsTokenEndpointAuthBasicUpdateParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsTokenEndpointAuthBasicUpdateParamType string

const (
	ManagedAgentsTokenEndpointAuthBasicUpdateParamTypeClientSecretBasic ManagedAgentsTokenEndpointAuthBasicUpdateParamType = "client_secret_basic"
)

// Token endpoint requires no client authentication.
//
// The property Type is required.
type ManagedAgentsTokenEndpointAuthNoneParam struct {
	// Any of "none".
	Type ManagedAgentsTokenEndpointAuthNoneParamType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsTokenEndpointAuthNoneParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsTokenEndpointAuthNoneParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsTokenEndpointAuthNoneParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsTokenEndpointAuthNoneParamType string

const (
	ManagedAgentsTokenEndpointAuthNoneParamTypeNone ManagedAgentsTokenEndpointAuthNoneParamType = "none"
)

// Token endpoint requires no client authentication.
type ManagedAgentsTokenEndpointAuthNoneResponse struct {
	// Any of "none".
	Type ManagedAgentsTokenEndpointAuthNoneResponseType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsTokenEndpointAuthNoneResponse) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsTokenEndpointAuthNoneResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsTokenEndpointAuthNoneResponseType string

const (
	ManagedAgentsTokenEndpointAuthNoneResponseTypeNone ManagedAgentsTokenEndpointAuthNoneResponseType = "none"
)

// Token endpoint uses POST body authentication with client credentials.
//
// The properties ClientSecret, Type are required.
type ManagedAgentsTokenEndpointAuthPostParam struct {
	// OAuth client secret.
	ClientSecret string `json:"client_secret" api:"required"`
	// Any of "client_secret_post".
	Type ManagedAgentsTokenEndpointAuthPostParamType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsTokenEndpointAuthPostParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsTokenEndpointAuthPostParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsTokenEndpointAuthPostParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsTokenEndpointAuthPostParamType string

const (
	ManagedAgentsTokenEndpointAuthPostParamTypeClientSecretPost ManagedAgentsTokenEndpointAuthPostParamType = "client_secret_post"
)

// Token endpoint uses POST body authentication with client credentials.
type ManagedAgentsTokenEndpointAuthPostResponse struct {
	// Any of "client_secret_post".
	Type ManagedAgentsTokenEndpointAuthPostResponseType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsTokenEndpointAuthPostResponse) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsTokenEndpointAuthPostResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsTokenEndpointAuthPostResponseType string

const (
	ManagedAgentsTokenEndpointAuthPostResponseTypeClientSecretPost ManagedAgentsTokenEndpointAuthPostResponseType = "client_secret_post"
)

// Updated POST body authentication parameters for the token endpoint.
//
// The property Type is required.
type ManagedAgentsTokenEndpointAuthPostUpdateParam struct {
	// Any of "client_secret_post".
	Type ManagedAgentsTokenEndpointAuthPostUpdateParamType `json:"type,omitzero" api:"required"`
	// Updated OAuth client secret.
	ClientSecret param.Opt[string] `json:"client_secret,omitzero"`
	paramObj
}

func (r ManagedAgentsTokenEndpointAuthPostUpdateParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsTokenEndpointAuthPostUpdateParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsTokenEndpointAuthPostUpdateParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsTokenEndpointAuthPostUpdateParamType string

const (
	ManagedAgentsTokenEndpointAuthPostUpdateParamTypeClientSecretPost ManagedAgentsTokenEndpointAuthPostUpdateParamType = "client_secret_post"
)

// Substitute the secret on any host the session's Environment network policy
// permits egress to. The Environment's network policy is the only boundary on
// where the secret can reach.
//
// The property Type is required.
type ManagedAgentsUnrestrictedCredentialNetworkingParams struct {
	// Any of "unrestricted".
	Type ManagedAgentsUnrestrictedCredentialNetworkingParamsType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsUnrestrictedCredentialNetworkingParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsUnrestrictedCredentialNetworkingParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsUnrestrictedCredentialNetworkingParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsUnrestrictedCredentialNetworkingParamsType string

const (
	ManagedAgentsUnrestrictedCredentialNetworkingParamsTypeUnrestricted ManagedAgentsUnrestrictedCredentialNetworkingParamsType = "unrestricted"
)

// The secret is substituted on any host the session's Environment network policy
// permits egress to.
type ManagedAgentsUnrestrictedCredentialNetworkingResponse struct {
	// Any of "unrestricted".
	Type ManagedAgentsUnrestrictedCredentialNetworkingResponseType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsUnrestrictedCredentialNetworkingResponse) RawJSON() string {
	return r.JSON.raw
}
func (r *ManagedAgentsUnrestrictedCredentialNetworkingResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsUnrestrictedCredentialNetworkingResponseType string

const (
	ManagedAgentsUnrestrictedCredentialNetworkingResponseTypeUnrestricted ManagedAgentsUnrestrictedCredentialNetworkingResponseType = "unrestricted"
)

type VaultCredentialNewParams struct {
	// Authentication details for creating a credential.
	Auth VaultCredentialNewParamsAuthUnion `json:"auth,omitzero" api:"required"`
	// Human-readable name for the credential. Up to 255 characters.
	DisplayName param.Opt[string] `json:"display_name,omitzero"`
	// Arbitrary key-value metadata to attach to the credential. Maximum 16 pairs, keys
	// up to 64 chars, values up to 512 chars.
	Metadata map[string]string `json:"metadata,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

func (r VaultCredentialNewParams) MarshalJSON() (data []byte, err error) {
	type shadow VaultCredentialNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultCredentialNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type VaultCredentialNewParamsAuthUnion struct {
	OfMCPOAuth            *ManagedAgentsMCPOAuthCreateParams            `json:",omitzero,inline"`
	OfStaticBearer        *ManagedAgentsStaticBearerCreateParams        `json:",omitzero,inline"`
	OfEnvironmentVariable *ManagedAgentsEnvironmentVariableCreateParams `json:",omitzero,inline"`
	paramUnion
}

func (u VaultCredentialNewParamsAuthUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfMCPOAuth, u.OfStaticBearer, u.OfEnvironmentVariable)
}
func (u *VaultCredentialNewParamsAuthUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *VaultCredentialNewParamsAuthUnion) asAny() any {
	if !param.IsOmitted(u.OfMCPOAuth) {
		return u.OfMCPOAuth
	} else if !param.IsOmitted(u.OfStaticBearer) {
		return u.OfStaticBearer
	} else if !param.IsOmitted(u.OfEnvironmentVariable) {
		return u.OfEnvironmentVariable
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u VaultCredentialNewParamsAuthUnion) GetAccessToken() *string {
	if vt := u.OfMCPOAuth; vt != nil {
		return &vt.AccessToken
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u VaultCredentialNewParamsAuthUnion) GetExpiresAt() *time.Time {
	if vt := u.OfMCPOAuth; vt != nil && vt.ExpiresAt.Valid() {
		return &vt.ExpiresAt.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u VaultCredentialNewParamsAuthUnion) GetRefresh() *ManagedAgentsMCPOAuthRefreshParams {
	if vt := u.OfMCPOAuth; vt != nil {
		return &vt.Refresh
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u VaultCredentialNewParamsAuthUnion) GetToken() *string {
	if vt := u.OfStaticBearer; vt != nil {
		return &vt.Token
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u VaultCredentialNewParamsAuthUnion) GetNetworking() *ManagedAgentsCredentialNetworkingParamsUnion {
	if vt := u.OfEnvironmentVariable; vt != nil {
		return &vt.Networking
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u VaultCredentialNewParamsAuthUnion) GetSecretName() *string {
	if vt := u.OfEnvironmentVariable; vt != nil {
		return &vt.SecretName
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u VaultCredentialNewParamsAuthUnion) GetSecretValue() *string {
	if vt := u.OfEnvironmentVariable; vt != nil {
		return &vt.SecretValue
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u VaultCredentialNewParamsAuthUnion) GetInjectionLocation() *ManagedAgentsInjectionLocationParams {
	if vt := u.OfEnvironmentVariable; vt != nil {
		return &vt.InjectionLocation
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u VaultCredentialNewParamsAuthUnion) GetMCPServerURL() *string {
	if vt := u.OfMCPOAuth; vt != nil {
		return (*string)(&vt.MCPServerURL)
	} else if vt := u.OfStaticBearer; vt != nil {
		return (*string)(&vt.MCPServerURL)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u VaultCredentialNewParamsAuthUnion) GetType() *string {
	if vt := u.OfMCPOAuth; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfStaticBearer; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfEnvironmentVariable; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[VaultCredentialNewParamsAuthUnion](
		"type",
		apijson.Discriminator[ManagedAgentsMCPOAuthCreateParams]("mcp_oauth"),
		apijson.Discriminator[ManagedAgentsStaticBearerCreateParams]("static_bearer"),
		apijson.Discriminator[ManagedAgentsEnvironmentVariableCreateParams]("environment_variable"),
	)
}

type VaultCredentialGetParams struct {
	VaultID string `path:"vault_id" api:"required" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

type VaultCredentialUpdateParams struct {
	VaultID string `path:"vault_id" api:"required" json:"-"`
	// Deprecated: Qoder does not support updating credential display names.
	// Leave this field omitted; setting it returns an error before sending a request.
	DisplayName param.Opt[string] `json:"display_name,omitzero"`
	// Metadata patch. Set a key to a string to upsert it, or to null to delete it.
	// Omitted keys are preserved.
	Metadata map[string]any `json:"metadata,omitzero"`
	// Updated authentication details for a credential.
	Auth VaultCredentialUpdateParamsAuthUnion `json:"auth,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

func (r VaultCredentialUpdateParams) MarshalJSON() (data []byte, err error) {
	if !param.IsOmitted(r.DisplayName) {
		return nil, errors.New("Qoder does not support updating credential display_name; omit DisplayName and update auth or metadata instead")
	}
	type shadow VaultCredentialUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *VaultCredentialUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type VaultCredentialUpdateParamsAuthUnion struct {
	OfMCPOAuth            *ManagedAgentsMCPOAuthUpdateParams            `json:",omitzero,inline"`
	OfStaticBearer        *ManagedAgentsStaticBearerUpdateParams        `json:",omitzero,inline"`
	OfEnvironmentVariable *ManagedAgentsEnvironmentVariableUpdateParams `json:",omitzero,inline"`
	paramUnion
}

func (u VaultCredentialUpdateParamsAuthUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfMCPOAuth, u.OfStaticBearer, u.OfEnvironmentVariable)
}
func (u *VaultCredentialUpdateParamsAuthUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *VaultCredentialUpdateParamsAuthUnion) asAny() any {
	if !param.IsOmitted(u.OfMCPOAuth) {
		return u.OfMCPOAuth
	} else if !param.IsOmitted(u.OfStaticBearer) {
		return u.OfStaticBearer
	} else if !param.IsOmitted(u.OfEnvironmentVariable) {
		return u.OfEnvironmentVariable
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u VaultCredentialUpdateParamsAuthUnion) GetAccessToken() *string {
	if vt := u.OfMCPOAuth; vt != nil && vt.AccessToken.Valid() {
		return &vt.AccessToken.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u VaultCredentialUpdateParamsAuthUnion) GetExpiresAt() *time.Time {
	if vt := u.OfMCPOAuth; vt != nil && vt.ExpiresAt.Valid() {
		return &vt.ExpiresAt.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u VaultCredentialUpdateParamsAuthUnion) GetRefresh() *ManagedAgentsMCPOAuthRefreshUpdateParams {
	if vt := u.OfMCPOAuth; vt != nil {
		return &vt.Refresh
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u VaultCredentialUpdateParamsAuthUnion) GetToken() *string {
	if vt := u.OfStaticBearer; vt != nil && vt.Token.Valid() {
		return &vt.Token.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u VaultCredentialUpdateParamsAuthUnion) GetInjectionLocation() *ManagedAgentsInjectionLocationUpdateParams {
	if vt := u.OfEnvironmentVariable; vt != nil {
		return &vt.InjectionLocation
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u VaultCredentialUpdateParamsAuthUnion) GetNetworking() *ManagedAgentsCredentialNetworkingParamsUnion {
	if vt := u.OfEnvironmentVariable; vt != nil {
		return &vt.Networking
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u VaultCredentialUpdateParamsAuthUnion) GetSecretValue() *string {
	if vt := u.OfEnvironmentVariable; vt != nil && vt.SecretValue.Valid() {
		return &vt.SecretValue.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u VaultCredentialUpdateParamsAuthUnion) GetType() *string {
	if vt := u.OfMCPOAuth; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfStaticBearer; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfEnvironmentVariable; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[VaultCredentialUpdateParamsAuthUnion](
		"type",
		apijson.Discriminator[ManagedAgentsMCPOAuthUpdateParams]("mcp_oauth"),
		apijson.Discriminator[ManagedAgentsStaticBearerUpdateParams]("static_bearer"),
		apijson.Discriminator[ManagedAgentsEnvironmentVariableUpdateParams]("environment_variable"),
	)
}

type VaultCredentialListParams struct {
	Name param.Opt[string] `query:"name,omitzero" json:"-"`

	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	AfterID  param.Opt[string] `query:"after_id,omitzero" json:"-"`

	// Whether to include archived credentials in the results.
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// Maximum number of credentials to return per page. Defaults to 20, maximum 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque pagination token from a previous `list_credentials` response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [VaultCredentialListParams]'s query parameters as
// `url.Values`.
func (r VaultCredentialListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type VaultCredentialDeleteParams struct {
	VaultID string `path:"vault_id" api:"required" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

type VaultCredentialArchiveParams struct {
	VaultID string `path:"vault_id" api:"required" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

type VaultCredentialMCPOAuthValidateParams struct {
	VaultID string `path:"vault_id" api:"required" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}
