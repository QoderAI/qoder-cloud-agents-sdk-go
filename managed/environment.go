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
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/constant"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/pagination"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
)

// EnvironmentService contains methods and other services that help with
// interacting with the Qoder Cloud Agents API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewEnvironmentService] method instead.
type EnvironmentService struct {
	Options []option.RequestOption
	Work    EnvironmentWorkService
}

// NewEnvironmentService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewEnvironmentService(opts ...option.RequestOption) (r EnvironmentService) {
	r = EnvironmentService{}
	r.Options = opts
	r.Work = NewEnvironmentWorkService(opts...)
	return
}

// Create a new environment with the specified configuration.
func (r *EnvironmentService) New(ctx context.Context, params EnvironmentNewParams, opts ...option.RequestOption) (res *Environment, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	path := "environments"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Retrieve a specific environment by ID.
func (r *EnvironmentService) Get(ctx context.Context, environmentID string, query EnvironmentGetParams, opts ...option.RequestOption) (res *Environment, err error) {
	for _, v := range query.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if environmentID == "" {
		err = errors.New("missing required environment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("environments/%s", url.PathEscape(environmentID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update an existing environment's configuration.
func (r *EnvironmentService) Update(ctx context.Context, environmentID string, params EnvironmentUpdateParams, opts ...option.RequestOption) (res *Environment, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if environmentID == "" {
		err = errors.New("missing required environment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("environments/%s", url.PathEscape(environmentID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List environments with pagination support.
func (r *EnvironmentService) List(ctx context.Context, params EnvironmentListParams, opts ...option.RequestOption) (res *pagination.PageCursor[Environment], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "environments"
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

// List environments with pagination support.
func (r *EnvironmentService) ListAutoPaging(ctx context.Context, params EnvironmentListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[Environment] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// Delete an environment by ID. Returns a confirmation of the deletion.
func (r *EnvironmentService) Delete(ctx context.Context, environmentID string, body EnvironmentDeleteParams, opts ...option.RequestOption) (res *EnvironmentDeleteResponse, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if environmentID == "" {
		err = errors.New("missing required environment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("environments/%s", url.PathEscape(environmentID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Archive an environment by ID. Archived environments cannot be used to create new
// sessions.
func (r *EnvironmentService) Archive(ctx context.Context, environmentID string, body EnvironmentArchiveParams, opts ...option.RequestOption) (res *Environment, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if environmentID == "" {
		err = errors.New("missing required environment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("environments/%s/archive", url.PathEscape(environmentID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// `cloud` environment configuration.
type CloudConfig struct {
	SetupScript string `json:"setup_script"`
	// Network configuration policy.
	Networking CloudConfigNetworkingUnion `json:"networking" api:"required"`
	// Package manager configuration.
	Packages Packages `json:"packages" api:"required"`
	// Environment type
	Type constant.Cloud `json:"type" default:"cloud"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		SetupScript respjson.Field
		Networking  respjson.Field
		Packages    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r CloudConfig) RawJSON() string { return r.JSON.raw }
func (r *CloudConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// CloudConfigNetworkingUnion contains all possible properties and values from
// [UnrestrictedNetwork], [LimitedNetwork].
//
// Use the [CloudConfigNetworkingUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type CloudConfigNetworkingUnion struct {
	// Any of "unrestricted", "limited".
	Type string `json:"type"`
	// This field is from variant [LimitedNetwork].
	AllowMCPServers bool `json:"allow_mcp_servers"`
	// This field is from variant [LimitedNetwork].
	AllowPackageManagers bool `json:"allow_package_managers"`
	// This field is from variant [LimitedNetwork].
	AllowedHosts []string `json:"allowed_hosts"`
	JSON         struct {
		Type                 respjson.Field
		AllowMCPServers      respjson.Field
		AllowPackageManagers respjson.Field
		AllowedHosts         respjson.Field
		raw                  string
	} `json:"-"`
}

// anyCloudConfigNetworking is implemented by each variant of
// [CloudConfigNetworkingUnion] to add type safety for the return type of
// [CloudConfigNetworkingUnion.AsAny]
type anyCloudConfigNetworking interface {
	implCloudConfigNetworkingUnion()
}

func (UnrestrictedNetwork) implCloudConfigNetworkingUnion() {}
func (LimitedNetwork) implCloudConfigNetworkingUnion()      {}

// Use the following switch statement to find the correct variant
//
//	switch variant := CloudConfigNetworkingUnion.AsAny().(type) {
//	case qoder.UnrestrictedNetwork:
//	case qoder.LimitedNetwork:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u CloudConfigNetworkingUnion) AsAny() anyCloudConfigNetworking {
	switch u.Type {
	case "unrestricted":
		return u.AsUnrestricted()
	case "limited":
		return u.AsLimited()
	}
	return nil
}

func (u CloudConfigNetworkingUnion) AsUnrestricted() (v UnrestrictedNetwork) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u CloudConfigNetworkingUnion) AsLimited() (v LimitedNetwork) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u CloudConfigNetworkingUnion) RawJSON() string { return u.JSON.raw }

func (r *CloudConfigNetworkingUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Request params for `cloud` environment configuration.
//
// Fields default to null; on update, omitted fields preserve the existing value.
//
// The property Type is required.
type CloudConfigParams struct {
	SetupScript param.Opt[string] `json:"setup_script,omitzero"`
	// Network configuration policy. Omit on update to preserve the existing value.
	Networking CloudConfigParamsNetworkingUnion `json:"networking,omitzero"`
	// Specify packages (and optionally their versions) available in this environment.
	//
	// When versioning, use the version semantics relevant for the package manager,
	// e.g. for `pip` use `package==1.0.0`. You are responsible for validating the
	// package and version exist. Unversioned installs the latest.
	//
	// Under `limited` networking, requires `networking.allow_package_managers` to be
	// `true`.
	Packages PackagesParams `json:"packages,omitzero"`
	// Environment type
	//
	// This field can be elided, and will marshal its zero value as "cloud".
	Type constant.Cloud `json:"type" default:"cloud"`
	paramObj
}

func (r CloudConfigParams) MarshalJSON() (data []byte, err error) {
	type shadow CloudConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *CloudConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type CloudConfigParamsNetworkingUnion struct {
	OfUnrestricted *UnrestrictedNetworkParam `json:",omitzero,inline"`
	OfLimited      *LimitedNetworkParams     `json:",omitzero,inline"`
	paramUnion
}

func (u CloudConfigParamsNetworkingUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfUnrestricted, u.OfLimited)
}
func (u *CloudConfigParamsNetworkingUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *CloudConfigParamsNetworkingUnion) asAny() any {
	if !param.IsOmitted(u.OfUnrestricted) {
		return u.OfUnrestricted
	} else if !param.IsOmitted(u.OfLimited) {
		return u.OfLimited
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CloudConfigParamsNetworkingUnion) GetAllowMCPServers() *bool {
	if vt := u.OfLimited; vt != nil && vt.AllowMCPServers.Valid() {
		return &vt.AllowMCPServers.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CloudConfigParamsNetworkingUnion) GetAllowPackageManagers() *bool {
	if vt := u.OfLimited; vt != nil && vt.AllowPackageManagers.Valid() {
		return &vt.AllowPackageManagers.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CloudConfigParamsNetworkingUnion) GetAllowedHosts() []string {
	if vt := u.OfLimited; vt != nil {
		return vt.AllowedHosts
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u CloudConfigParamsNetworkingUnion) GetType() *string {
	if vt := u.OfUnrestricted; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfLimited; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[CloudConfigParamsNetworkingUnion](
		"type",
		apijson.Discriminator[UnrestrictedNetworkParam]("unrestricted"),
		apijson.Discriminator[LimitedNetworkParams]("limited"),
	)
}

// Unified Environment resource for both cloud and self-hosted environments.
type Environment struct {
	// Environment identifier (e.g., 'env\_...')
	ID string `json:"id" api:"required"`
	// RFC 3339 timestamp when environment was archived, or null if not archived
	ArchivedAt string `json:"archived_at" api:"required"`
	// Environment configuration (either Qoder Cloud or self-hosted)
	Config EnvironmentConfigUnion `json:"config" api:"required"`
	// RFC 3339 timestamp when environment was created
	CreatedAt string `json:"created_at" api:"required"`
	// User-provided description for the environment; null when unset
	Description string `json:"description" api:"required"`
	// User-provided metadata key-value pairs
	Metadata map[string]string `json:"metadata" api:"required"`
	// Human-readable name for the environment
	Name string `json:"name" api:"required"`
	// The type of object (always 'environment')
	Type constant.Environment `json:"type" default:"environment"`
	// RFC 3339 timestamp when environment was last updated
	UpdatedAt string `json:"updated_at" api:"required"`
	// The visibility scope for this environment. 'organization' means visible to all
	// accounts. 'account' means visible only to the owning account.
	//
	// Any of "organization", "account".
	Scope EnvironmentScope `json:"scope"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ArchivedAt  respjson.Field
		Config      respjson.Field
		CreatedAt   respjson.Field
		Description respjson.Field
		Metadata    respjson.Field
		Name        respjson.Field
		Type        respjson.Field
		UpdatedAt   respjson.Field
		Scope       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Environment) RawJSON() string { return r.JSON.raw }
func (r *Environment) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EnvironmentConfigUnion contains all possible properties and values from
// [CloudConfig], [SelfHostedConfig].
//
// Use the [EnvironmentConfigUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type EnvironmentConfigUnion struct {
	// This field is from variant [CloudConfig].
	Networking CloudConfigNetworkingUnion `json:"networking"`
	// This field is from variant [CloudConfig].
	Packages Packages `json:"packages"`
	// Any of "cloud", "self_hosted".
	Type string `json:"type"`
	JSON struct {
		Networking respjson.Field
		Packages   respjson.Field
		Type       respjson.Field
		raw        string
	} `json:"-"`
}

// anyEnvironmentConfig is implemented by each variant of
// [EnvironmentConfigUnion] to add type safety for the return type of
// [EnvironmentConfigUnion.AsAny]
type anyEnvironmentConfig interface {
	implEnvironmentConfigUnion()
}

func (CloudConfig) implEnvironmentConfigUnion()      {}
func (SelfHostedConfig) implEnvironmentConfigUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := EnvironmentConfigUnion.AsAny().(type) {
//	case qoder.CloudConfig:
//	case qoder.SelfHostedConfig:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u EnvironmentConfigUnion) AsAny() anyEnvironmentConfig {
	switch u.Type {
	case "cloud":
		return u.AsCloud()
	case "self_hosted":
		return u.AsSelfHosted()
	}
	return nil
}

func (u EnvironmentConfigUnion) AsCloud() (v CloudConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u EnvironmentConfigUnion) AsSelfHosted() (v SelfHostedConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u EnvironmentConfigUnion) RawJSON() string { return u.JSON.raw }

func (r *EnvironmentConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The visibility scope for this environment. 'organization' means visible to all
// accounts. 'account' means visible only to the owning account.
type EnvironmentScope string

const (
	EnvironmentScopeOrganization EnvironmentScope = "organization"
	EnvironmentScopeAccount      EnvironmentScope = "account"
)

// Response after deleting an environment.
type EnvironmentDeleteResponse struct {
	// Environment identifier
	ID string `json:"id" api:"required"`
	// The type of response
	//
	// Any of "environment_deleted".
	Type EnvironmentDeleteResponseType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r EnvironmentDeleteResponse) RawJSON() string { return r.JSON.raw }
func (r *EnvironmentDeleteResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The type of response
type EnvironmentDeleteResponseType string

const (
	EnvironmentDeleteResponseTypeEnvironmentDeleted EnvironmentDeleteResponseType = "environment_deleted"
)

// Limited network access.
type LimitedNetwork struct {
	// Permits outbound access to MCP server endpoints configured on the agent, beyond
	// those listed in the `allowed_hosts` array.
	AllowMCPServers bool `json:"allow_mcp_servers" api:"required"`
	// Permits outbound access to public package registries (PyPI, npm, etc.) beyond
	// those listed in the `allowed_hosts` array.
	AllowPackageManagers bool `json:"allow_package_managers" api:"required"`
	// Specifies domains the container can reach.
	AllowedHosts []string `json:"allowed_hosts" api:"required"`
	// Network policy type
	Type constant.Limited `json:"type" default:"limited"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		AllowMCPServers      respjson.Field
		AllowPackageManagers respjson.Field
		AllowedHosts         respjson.Field
		Type                 respjson.Field
		ExtraFields          map[string]respjson.Field
		raw                  string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r LimitedNetwork) RawJSON() string { return r.JSON.raw }
func (r *LimitedNetwork) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Limited network request params.
//
// Fields default to null; on update, omitted fields preserve the existing value.
//
// The property Type is required.
type LimitedNetworkParams struct {
	// Permits outbound access to MCP server endpoints configured on the agent, beyond
	// those listed in the `allowed_hosts` array. Defaults to `false`.
	AllowMCPServers param.Opt[bool] `json:"allow_mcp_servers,omitzero"`
	// Permits outbound access to public package registries (PyPI, npm, etc.) beyond
	// those listed in the `allowed_hosts` array. Defaults to `false` on creation. Must
	// be `true` when `packages` are specified.
	AllowPackageManagers param.Opt[bool] `json:"allow_package_managers,omitzero"`
	// Specifies domains the container can reach.
	AllowedHosts []string `json:"allowed_hosts,omitzero"`
	// Network policy type
	//
	// This field can be elided, and will marshal its zero value as "limited".
	Type constant.Limited `json:"type" default:"limited"`
	paramObj
}

func (r LimitedNetworkParams) MarshalJSON() (data []byte, err error) {
	type shadow LimitedNetworkParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *LimitedNetworkParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Packages (and their versions) available in this environment.
type Packages struct {
	// Ubuntu/Debian packages to install
	Apt []string `json:"apt" api:"required"`
	// Rust packages to install
	Cargo []string `json:"cargo" api:"required"`
	// Ruby packages to install
	Gem []string `json:"gem" api:"required"`
	// Go packages to install
	Go []string `json:"go" api:"required"`
	// Node.js packages to install
	Npm []string `json:"npm" api:"required"`
	// Python packages to install
	Pip []string `json:"pip" api:"required"`
	// Package configuration type
	//
	// Any of "packages".
	Type PackagesType `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Apt         respjson.Field
		Cargo       respjson.Field
		Gem         respjson.Field
		Go          respjson.Field
		Npm         respjson.Field
		Pip         respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Packages) RawJSON() string { return r.JSON.raw }
func (r *Packages) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Package configuration type
type PackagesType string

const (
	PackagesTypePackages PackagesType = "packages"
)

// Specify packages (and optionally their versions) available in this environment.
//
// When versioning, use the version semantics relevant for the package manager,
// e.g. for `pip` use `package==1.0.0`. You are responsible for validating the
// package and version exist. Unversioned installs the latest.
//
// Under `limited` networking, requires `networking.allow_package_managers` to be
// `true`.
type PackagesParams struct {
	// Ubuntu/Debian packages to install
	Apt []string `json:"apt,omitzero"`
	// Rust packages to install
	Cargo []string `json:"cargo,omitzero"`
	// Ruby packages to install
	Gem []string `json:"gem,omitzero"`
	// Go packages to install
	Go []string `json:"go,omitzero"`
	// Node.js packages to install
	Npm []string `json:"npm,omitzero"`
	// Python packages to install
	Pip []string `json:"pip,omitzero"`
	// Package configuration type
	//
	// Any of "packages".
	Type PackagesParamsType `json:"type,omitzero"`
	paramObj
}

func (r PackagesParams) MarshalJSON() (data []byte, err error) {
	type shadow PackagesParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *PackagesParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Package configuration type
type PackagesParamsType string

const (
	PackagesParamsTypePackages PackagesParamsType = "packages"
)

// Configuration for self-hosted environments.
type SelfHostedConfig struct {
	// Environment type
	Type constant.SelfHosted `json:"type" default:"self_hosted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SelfHostedConfig) RawJSON() string { return r.JSON.raw }
func (r *SelfHostedConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func NewSelfHostedConfigParams() SelfHostedConfigParams {
	return SelfHostedConfigParams{
		Type: "self_hosted",
	}
}

// Request params for `self_hosted` environment configuration.
//
// This struct has a constant value, construct it with
// [NewSelfHostedConfigParams].
type SelfHostedConfigParams struct {
	// Environment type
	Type constant.SelfHosted `json:"type" default:"self_hosted"`
	paramObj
}

func (r SelfHostedConfigParams) MarshalJSON() (data []byte, err error) {
	type shadow SelfHostedConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SelfHostedConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Unrestricted network access.
type UnrestrictedNetwork struct {
	// Network policy type
	Type constant.Unrestricted `json:"type" default:"unrestricted"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r UnrestrictedNetwork) RawJSON() string { return r.JSON.raw }
func (r *UnrestrictedNetwork) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this UnrestrictedNetwork to a UnrestrictedNetworkParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// UnrestrictedNetworkParam.Overrides()
func (r UnrestrictedNetwork) ToParam() UnrestrictedNetworkParam {
	return param.Override[UnrestrictedNetworkParam](json.RawMessage(r.RawJSON()))
}

func NewUnrestrictedNetworkParam() UnrestrictedNetworkParam {
	return UnrestrictedNetworkParam{
		Type: "unrestricted",
	}
}

// Unrestricted network access.
//
// This struct has a constant value, construct it with
// [NewUnrestrictedNetworkParam].
type UnrestrictedNetworkParam struct {
	// Network policy type
	Type constant.Unrestricted `json:"type" default:"unrestricted"`
	paramObj
}

func (r UnrestrictedNetworkParam) MarshalJSON() (data []byte, err error) {
	type shadow UnrestrictedNetworkParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *UnrestrictedNetworkParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type EnvironmentNewParams struct {
	// Human-readable name for the environment
	Name string `json:"name" api:"required"`
	// Optional description of the environment
	Description param.Opt[string] `json:"description,omitzero"`
	WorkspaceID param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`
	// Environment configuration
	Config EnvironmentNewParamsConfigUnion `json:"config,omitzero"`
	// The visibility scope for this environment. 'organization' makes the environment
	// visible to all accounts. 'account' restricts visibility to the owning account
	// only. Only applicable for self-hosted environments. If not specified, defaults
	// based on organization type.
	//
	// Any of "organization", "account".
	Scope EnvironmentNewParamsScope `json:"scope,omitzero"`
	// User-provided metadata key-value pairs
	Metadata map[string]string `json:"metadata,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

func (r EnvironmentNewParams) MarshalJSON() (data []byte, err error) {
	type shadow EnvironmentNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *EnvironmentNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type EnvironmentNewParamsConfigUnion struct {
	OfCloud      *CloudConfigParams      `json:",omitzero,inline"`
	OfSelfHosted *SelfHostedConfigParams `json:",omitzero,inline"`
	paramUnion
}

func (u EnvironmentNewParamsConfigUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfCloud, u.OfSelfHosted)
}
func (u *EnvironmentNewParamsConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *EnvironmentNewParamsConfigUnion) asAny() any {
	if !param.IsOmitted(u.OfCloud) {
		return u.OfCloud
	} else if !param.IsOmitted(u.OfSelfHosted) {
		return u.OfSelfHosted
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EnvironmentNewParamsConfigUnion) GetNetworking() *CloudConfigParamsNetworkingUnion {
	if vt := u.OfCloud; vt != nil {
		return &vt.Networking
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EnvironmentNewParamsConfigUnion) GetPackages() *PackagesParams {
	if vt := u.OfCloud; vt != nil {
		return &vt.Packages
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EnvironmentNewParamsConfigUnion) GetType() *string {
	if vt := u.OfCloud; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSelfHosted; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[EnvironmentNewParamsConfigUnion](
		"type",
		apijson.Discriminator[CloudConfigParams]("cloud"),
		apijson.Discriminator[SelfHostedConfigParams]("self_hosted"),
	)
}

// The visibility scope for this environment. 'organization' makes the environment
// visible to all accounts. 'account' restricts visibility to the owning account
// only. Only applicable for self-hosted environments. If not specified, defaults
// based on organization type.
type EnvironmentNewParamsScope string

const (
	EnvironmentNewParamsScopeOrganization EnvironmentNewParamsScope = "organization"
	EnvironmentNewParamsScopeAccount      EnvironmentNewParamsScope = "account"
)

type EnvironmentGetParams struct {
	WorkspaceID param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

type EnvironmentUpdateParams struct {
	// Updated description of the environment. Omit to preserve; null clears to null;
	// an empty string is stored as an empty string.
	Description param.Opt[string] `json:"description,omitzero"`
	// Updated name for the environment
	Name        param.Opt[string] `json:"name,omitzero"`
	WorkspaceID param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`
	// Updated environment configuration
	Config EnvironmentUpdateParamsConfigUnion `json:"config,omitzero"`
	// The visibility scope for this environment. 'organization' makes the environment
	// visible to all accounts. 'account' restricts visibility to the owning account
	// only.
	//
	// Any of "organization", "account".
	Scope EnvironmentUpdateParamsScope `json:"scope,omitzero"`
	// User-provided metadata key-value pairs. Set a value to null or empty string to
	// delete the key.
	Metadata map[string]any `json:"metadata,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

func (r EnvironmentUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow EnvironmentUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *EnvironmentUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type EnvironmentUpdateParamsConfigUnion struct {
	OfCloud      *CloudConfigParams      `json:",omitzero,inline"`
	OfSelfHosted *SelfHostedConfigParams `json:",omitzero,inline"`
	paramUnion
}

func (u EnvironmentUpdateParamsConfigUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfCloud, u.OfSelfHosted)
}
func (u *EnvironmentUpdateParamsConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *EnvironmentUpdateParamsConfigUnion) asAny() any {
	if !param.IsOmitted(u.OfCloud) {
		return u.OfCloud
	} else if !param.IsOmitted(u.OfSelfHosted) {
		return u.OfSelfHosted
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EnvironmentUpdateParamsConfigUnion) GetNetworking() *CloudConfigParamsNetworkingUnion {
	if vt := u.OfCloud; vt != nil {
		return &vt.Networking
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EnvironmentUpdateParamsConfigUnion) GetPackages() *PackagesParams {
	if vt := u.OfCloud; vt != nil {
		return &vt.Packages
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u EnvironmentUpdateParamsConfigUnion) GetType() *string {
	if vt := u.OfCloud; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSelfHosted; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[EnvironmentUpdateParamsConfigUnion](
		"type",
		apijson.Discriminator[CloudConfigParams]("cloud"),
		apijson.Discriminator[SelfHostedConfigParams]("self_hosted"),
	)
}

// The visibility scope for this environment. 'organization' makes the environment
// visible to all accounts. 'account' restricts visibility to the owning account
// only.
type EnvironmentUpdateParamsScope string

const (
	EnvironmentUpdateParamsScopeOrganization EnvironmentUpdateParamsScope = "organization"
	EnvironmentUpdateParamsScopeAccount      EnvironmentUpdateParamsScope = "account"
)

type EnvironmentListParams struct {
	CreatedAtGte param.Opt[time.Time] `query:"created_at[gte],omitzero" format:"date-time" json:"-"`
	CreatedAtLte param.Opt[time.Time] `query:"created_at[lte],omitzero" format:"date-time" json:"-"`

	// Opaque cursor from previous response for pagination. Pass the `next_page` value
	// from the previous response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Include archived environments in the response
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// Maximum number of environments to return
	Limit       param.Opt[int64]  `query:"limit,omitzero" json:"-"`
	WorkspaceID param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [EnvironmentListParams]'s query parameters as
// `url.Values`.
func (r EnvironmentListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type EnvironmentDeleteParams struct {
	WorkspaceID param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

type EnvironmentArchiveParams struct {
	WorkspaceID param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}
