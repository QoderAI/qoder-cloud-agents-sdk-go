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
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/paramutil"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
)

// AgentService contains methods and other services that help with interacting
// with the Qoder Cloud Agents API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAgentService] method instead.
type AgentService struct {
	Options  []option.RequestOption
	Versions AgentVersionService
}

// NewAgentService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewAgentService(opts ...option.RequestOption) (r AgentService) {
	r = AgentService{}
	r.Options = opts
	r.Versions = NewAgentVersionService(opts...)
	return
}

// Create Agent
func (r *AgentService) New(ctx context.Context, params AgentNewParams, opts ...option.RequestOption) (res *ManagedAgentsAgent, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	path := "agents"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Get Agent
func (r *AgentService) Get(ctx context.Context, agentID string, params AgentGetParams, opts ...option.RequestOption) (res *ManagedAgentsAgent, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if agentID == "" {
		err = errors.New("missing required agent_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("agents/%s", url.PathEscape(agentID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

// Update Agent
func (r *AgentService) Update(ctx context.Context, agentID string, params AgentUpdateParams, opts ...option.RequestOption) (res *ManagedAgentsAgent, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if agentID == "" {
		err = errors.New("missing required agent_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("agents/%s", url.PathEscape(agentID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List Agents
func (r *AgentService) List(ctx context.Context, params AgentListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ManagedAgentsAgent], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "agents"
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

// List Agents
func (r *AgentService) ListAutoPaging(ctx context.Context, params AgentListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ManagedAgentsAgent] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// Archive Agent
func (r *AgentService) Archive(ctx context.Context, agentID string, body AgentArchiveParams, opts ...option.RequestOption) (res *ManagedAgentsAgent, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if agentID == "" {
		err = errors.New("missing required agent_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("agents/%s/archive", url.PathEscape(agentID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Platform advisor roster entry: a model the session's primary thread may consult
// mid-turn.
type ManagedAgentsAdvisor struct {
	// The advisor model id.
	Model string `json:"model" api:"required"`
	// Any of "advisor".
	Type ManagedAgentsAdvisorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Model       respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsAdvisor) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsAdvisor) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsAdvisorType string

const (
	ManagedAgentsAdvisorTypeAdvisor ManagedAgentsAdvisorType = "advisor"
)

// A Managed Agents `agent`.
type ManagedAgentsAgent struct {
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ArchivedAt time.Time `json:"archived_at" api:"required" format:"date-time"`
	// A timestamp in RFC 3339 format
	CreatedAt   time.Time                             `json:"created_at" api:"required" format:"date-time"`
	Description string                                `json:"description" api:"required"`
	MCPServers  []ManagedAgentsMCPServerURLDefinition `json:"mcp_servers" api:"required"`
	Metadata    map[string]string                     `json:"metadata" api:"required"`
	// Model identifier and configuration.
	Model ManagedAgentsModelConfig `json:"model" api:"required"`
	// Resolved coordinator topology with a concrete agent roster.
	Multiagent ManagedAgentsMultiagent        `json:"multiagent" api:"required"`
	Name       string                         `json:"name" api:"required"`
	Skills     []ManagedAgentsAgentSkillUnion `json:"skills" api:"required"`
	System     string                         `json:"system" api:"required"`
	Tools      []ManagedAgentsAgentToolUnion  `json:"tools" api:"required"`
	// Any of "agent".
	Type ManagedAgentsAgentType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// The agent's current version. Starts at 1 and increments when the agent is
	// modified.
	Version int64 `json:"version" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ArchivedAt  respjson.Field
		CreatedAt   respjson.Field
		Description respjson.Field
		MCPServers  respjson.Field
		Metadata    respjson.Field
		Model       respjson.Field
		Multiagent  respjson.Field
		Name        respjson.Field
		Skills      respjson.Field
		System      respjson.Field
		Tools       respjson.Field
		Type        respjson.Field
		UpdatedAt   respjson.Field
		Version     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsAgent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsAgent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsAgentSkillUnion contains all possible properties and values
// from [ManagedAgentsQoderSkill], [ManagedAgentsCustomSkill].
//
// Use the [ManagedAgentsAgentSkillUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsAgentSkillUnion struct {
	SkillID string `json:"skill_id"`
	// Any of "qoder", "custom".
	Type    string `json:"type"`
	Version string `json:"version"`
	JSON    struct {
		SkillID respjson.Field
		Type    respjson.Field
		Version respjson.Field
		raw     string
	} `json:"-"`
}

// anyManagedAgentsAgentSkill is implemented by each variant of
// [ManagedAgentsAgentSkillUnion] to add type safety for the return type of
// [ManagedAgentsAgentSkillUnion.AsAny]
type anyManagedAgentsAgentSkill interface {
	implManagedAgentsAgentSkillUnion()
}

func (ManagedAgentsQoderSkill) implManagedAgentsAgentSkillUnion()  {}
func (ManagedAgentsCustomSkill) implManagedAgentsAgentSkillUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsAgentSkillUnion.AsAny().(type) {
//	case qoder.ManagedAgentsQoderSkill:
//	case qoder.ManagedAgentsCustomSkill:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsAgentSkillUnion) AsAny() anyManagedAgentsAgentSkill {
	switch u.Type {
	case "qoder":
		return u.AsQoder()
	case "custom":
		return u.AsCustom()
	}
	return nil
}

func (u ManagedAgentsAgentSkillUnion) AsQoder() (v ManagedAgentsQoderSkill) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsAgentSkillUnion) AsCustom() (v ManagedAgentsCustomSkill) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsAgentSkillUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsAgentSkillUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsAgentToolUnion contains all possible properties and values from
// [ManagedAgentsAgentToolset20260401], [ManagedAgentsMCPToolset],
// [ManagedAgentsCustomTool].
//
// Use the [ManagedAgentsAgentToolUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsAgentToolUnion struct {
	EnabledTools    []string `json:"enabled_tools"`
	DisallowedTools []string `json:"disallowed_tools"`
	// This field is a union of [[]ManagedAgentsAgentToolConfigUnion],
	// [[]ManagedAgentsMCPToolConfig]
	Configs ManagedAgentsAgentToolUnionConfigs `json:"configs"`
	// This field is a union of [ManagedAgentsAgentToolsetDefaultConfig],
	// [ManagedAgentsMCPToolsetDefaultConfig]
	DefaultConfig ManagedAgentsAgentToolUnionDefaultConfig `json:"default_config"`
	// Any of "agent_toolset_20260401", "mcp_toolset", "custom".
	Type string `json:"type"`
	// This field is from variant [ManagedAgentsMCPToolset].
	MCPServerName string `json:"mcp_server_name"`
	// This field is from variant [ManagedAgentsCustomTool].
	Description string `json:"description"`
	// This field is from variant [ManagedAgentsCustomTool].
	InputSchema ManagedAgentsCustomToolInputSchema `json:"input_schema"`
	// This field is from variant [ManagedAgentsCustomTool].
	Name string `json:"name"`
	JSON struct {
		EnabledTools    respjson.Field
		DisallowedTools respjson.Field
		Configs         respjson.Field
		DefaultConfig   respjson.Field
		Type            respjson.Field
		MCPServerName   respjson.Field
		Description     respjson.Field
		InputSchema     respjson.Field
		Name            respjson.Field
		raw             string
	} `json:"-"`
}

// anyManagedAgentsAgentTool is implemented by each variant of
// [ManagedAgentsAgentToolUnion] to add type safety for the return type of
// [ManagedAgentsAgentToolUnion.AsAny]
type anyManagedAgentsAgentTool interface {
	implManagedAgentsAgentToolUnion()
}

func (ManagedAgentsAgentToolset20260401) implManagedAgentsAgentToolUnion() {}
func (ManagedAgentsMCPToolset) implManagedAgentsAgentToolUnion()           {}
func (ManagedAgentsCustomTool) implManagedAgentsAgentToolUnion()           {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsAgentToolUnion.AsAny().(type) {
//	case qoder.ManagedAgentsAgentToolset20260401:
//	case qoder.ManagedAgentsMCPToolset:
//	case qoder.ManagedAgentsCustomTool:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsAgentToolUnion) AsAny() anyManagedAgentsAgentTool {
	switch u.Type {
	case "agent_toolset_20260401":
		return u.AsAgentToolset20260401()
	case "mcp_toolset":
		return u.AsMCPToolset()
	case "custom":
		return u.AsCustom()
	}
	return nil
}

func (u ManagedAgentsAgentToolUnion) AsAgentToolset20260401() (v ManagedAgentsAgentToolset20260401) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsAgentToolUnion) AsMCPToolset() (v ManagedAgentsMCPToolset) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsAgentToolUnion) AsCustom() (v ManagedAgentsCustomTool) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsAgentToolUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsAgentToolUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsAgentToolUnionConfigs is an implicit subunion of
// [ManagedAgentsAgentToolUnion]. ManagedAgentsAgentToolUnionConfigs
// provides convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsAgentToolUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfManagedAgentsAgentToolConfigArray
// OfManagedAgentsMCPToolConfigArray]
type ManagedAgentsAgentToolUnionConfigs struct {
	// This field will be present if the value is a
	// [[]ManagedAgentsAgentToolConfigUnion] instead of an object.
	OfManagedAgentsAgentToolConfigArray []ManagedAgentsAgentToolConfigUnion `json:",inline"`
	// This field will be present if the value is a [[]ManagedAgentsMCPToolConfig]
	// instead of an object.
	OfManagedAgentsMCPToolConfigArray []ManagedAgentsMCPToolConfig `json:",inline"`
	JSON                              struct {
		OfManagedAgentsAgentToolConfigArray respjson.Field
		OfManagedAgentsMCPToolConfigArray   respjson.Field
		raw                                 string
	} `json:"-"`
}

func (r *ManagedAgentsAgentToolUnionConfigs) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsAgentToolUnionDefaultConfig is an implicit subunion of
// [ManagedAgentsAgentToolUnion]. ManagedAgentsAgentToolUnionDefaultConfig
// provides convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsAgentToolUnion].
type ManagedAgentsAgentToolUnionDefaultConfig struct {
	Enabled bool `json:"enabled"`
	// This field is a union of
	// [ManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion],
	// [ManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion]
	PermissionPolicy ManagedAgentsAgentToolUnionDefaultConfigPermissionPolicy `json:"permission_policy"`
	JSON             struct {
		Enabled          respjson.Field
		PermissionPolicy respjson.Field
		raw              string
	} `json:"-"`
}

func (r *ManagedAgentsAgentToolUnionDefaultConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsAgentToolUnionDefaultConfigPermissionPolicy is an implicit
// subunion of [ManagedAgentsAgentToolUnion].
// ManagedAgentsAgentToolUnionDefaultConfigPermissionPolicy provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsAgentToolUnion].
type ManagedAgentsAgentToolUnionDefaultConfigPermissionPolicy struct {
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

func (r *ManagedAgentsAgentToolUnionDefaultConfigPermissionPolicy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsAgentType string

const (
	ManagedAgentsAgentTypeAgent ManagedAgentsAgentType = "agent"
)

// A resolved agent reference with a concrete version.
type ManagedAgentsAgentReference struct {
	ID string `json:"id" api:"required"`
	// Any of "agent".
	Type    ManagedAgentsAgentReferenceType `json:"type" api:"required"`
	Version int64                           `json:"version" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		Version     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsAgentReference) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsAgentReference) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsAgentReferenceType string

const (
	ManagedAgentsAgentReferenceTypeAgent ManagedAgentsAgentReferenceType = "agent"
)

// ManagedAgentsAgentToolConfigUnion contains all possible properties and
// values from [ManagedAgentsBashToolConfig],
// [ManagedAgentsEditToolConfig], [ManagedAgentsReadToolConfig],
// [ManagedAgentsWriteToolConfig], [ManagedAgentsGlobToolConfig],
// [ManagedAgentsGrepToolConfig], [ManagedAgentsWebFetchToolConfig],
// [ManagedAgentsWebSearchToolConfig].
//
// Use the [ManagedAgentsAgentToolConfigUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsAgentToolConfigUnion struct {
	Enabled bool   `json:"enabled"`
	Name    string `json:"name"`
	// This field is a union of [ManagedAgentsBashToolConfigPermissionPolicyUnion],
	// [ManagedAgentsEditToolConfigPermissionPolicyUnion],
	// [ManagedAgentsReadToolConfigPermissionPolicyUnion],
	// [ManagedAgentsWriteToolConfigPermissionPolicyUnion],
	// [ManagedAgentsGlobToolConfigPermissionPolicyUnion],
	// [ManagedAgentsGrepToolConfigPermissionPolicyUnion],
	// [ManagedAgentsWebFetchToolConfigPermissionPolicyUnion],
	// [ManagedAgentsWebSearchToolConfigPermissionPolicyUnion]
	PermissionPolicy ManagedAgentsAgentToolConfigUnionPermissionPolicy `json:"permission_policy"`
	// Any of "bash", "edit", "read", "write", "glob", "grep", "web_fetch",
	// "web_search".
	Type           string   `json:"type"`
	AllowedDomains []string `json:"allowed_domains"`
	BlockedDomains []string `json:"blocked_domains"`
	// This field is from variant [ManagedAgentsWebFetchToolConfig].
	MaxContentTokens int64 `json:"max_content_tokens"`
	// This field is from variant [ManagedAgentsWebSearchToolConfig].
	UserLocation ManagedAgentsUserLocation `json:"user_location"`
	JSON         struct {
		Enabled          respjson.Field
		Name             respjson.Field
		PermissionPolicy respjson.Field
		Type             respjson.Field
		AllowedDomains   respjson.Field
		BlockedDomains   respjson.Field
		MaxContentTokens respjson.Field
		UserLocation     respjson.Field
		raw              string
	} `json:"-"`
}

// anyManagedAgentsAgentToolConfig is implemented by each variant of
// [ManagedAgentsAgentToolConfigUnion] to add type safety for the return type
// of [ManagedAgentsAgentToolConfigUnion.AsAny]
type anyManagedAgentsAgentToolConfig interface {
	implManagedAgentsAgentToolConfigUnion()
}

func (ManagedAgentsBashToolConfig) implManagedAgentsAgentToolConfigUnion()      {}
func (ManagedAgentsEditToolConfig) implManagedAgentsAgentToolConfigUnion()      {}
func (ManagedAgentsReadToolConfig) implManagedAgentsAgentToolConfigUnion()      {}
func (ManagedAgentsWriteToolConfig) implManagedAgentsAgentToolConfigUnion()     {}
func (ManagedAgentsGlobToolConfig) implManagedAgentsAgentToolConfigUnion()      {}
func (ManagedAgentsGrepToolConfig) implManagedAgentsAgentToolConfigUnion()      {}
func (ManagedAgentsWebFetchToolConfig) implManagedAgentsAgentToolConfigUnion()  {}
func (ManagedAgentsWebSearchToolConfig) implManagedAgentsAgentToolConfigUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsAgentToolConfigUnion.AsAny().(type) {
//	case qoder.ManagedAgentsBashToolConfig:
//	case qoder.ManagedAgentsEditToolConfig:
//	case qoder.ManagedAgentsReadToolConfig:
//	case qoder.ManagedAgentsWriteToolConfig:
//	case qoder.ManagedAgentsGlobToolConfig:
//	case qoder.ManagedAgentsGrepToolConfig:
//	case qoder.ManagedAgentsWebFetchToolConfig:
//	case qoder.ManagedAgentsWebSearchToolConfig:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsAgentToolConfigUnion) AsAny() anyManagedAgentsAgentToolConfig {
	switch u.Type {
	case "bash":
		return u.AsBash()
	case "edit":
		return u.AsEdit()
	case "read":
		return u.AsRead()
	case "write":
		return u.AsWrite()
	case "glob":
		return u.AsGlob()
	case "grep":
		return u.AsGrep()
	case "web_fetch":
		return u.AsWebFetch()
	case "web_search":
		return u.AsWebSearch()
	}
	return nil
}

func (u ManagedAgentsAgentToolConfigUnion) AsBash() (v ManagedAgentsBashToolConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsAgentToolConfigUnion) AsEdit() (v ManagedAgentsEditToolConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsAgentToolConfigUnion) AsRead() (v ManagedAgentsReadToolConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsAgentToolConfigUnion) AsWrite() (v ManagedAgentsWriteToolConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsAgentToolConfigUnion) AsGlob() (v ManagedAgentsGlobToolConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsAgentToolConfigUnion) AsGrep() (v ManagedAgentsGrepToolConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsAgentToolConfigUnion) AsWebFetch() (v ManagedAgentsWebFetchToolConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsAgentToolConfigUnion) AsWebSearch() (v ManagedAgentsWebSearchToolConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsAgentToolConfigUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsAgentToolConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsAgentToolConfigUnionPermissionPolicy is an implicit subunion of
// [ManagedAgentsAgentToolConfigUnion].
// ManagedAgentsAgentToolConfigUnionPermissionPolicy provides convenient access
// to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsAgentToolConfigUnion].
type ManagedAgentsAgentToolConfigUnionPermissionPolicy struct {
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

func (r *ManagedAgentsAgentToolConfigUnionPermissionPolicy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsAgentToolConfigParamsUnion struct {
	OfBash      *ManagedAgentsBashToolConfigParams      `json:",omitzero,inline"`
	OfEdit      *ManagedAgentsEditToolConfigParams      `json:",omitzero,inline"`
	OfRead      *ManagedAgentsReadToolConfigParams      `json:",omitzero,inline"`
	OfWrite     *ManagedAgentsWriteToolConfigParams     `json:",omitzero,inline"`
	OfGlob      *ManagedAgentsGlobToolConfigParams      `json:",omitzero,inline"`
	OfGrep      *ManagedAgentsGrepToolConfigParams      `json:",omitzero,inline"`
	OfWebFetch  *ManagedAgentsWebFetchToolConfigParams  `json:",omitzero,inline"`
	OfWebSearch *ManagedAgentsWebSearchToolConfigParams `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsAgentToolConfigParamsUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfBash,
		u.OfEdit,
		u.OfRead,
		u.OfWrite,
		u.OfGlob,
		u.OfGrep,
		u.OfWebFetch,
		u.OfWebSearch)
}
func (u *ManagedAgentsAgentToolConfigParamsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsAgentToolConfigParamsUnion) asAny() any {
	if !param.IsOmitted(u.OfBash) {
		return u.OfBash
	} else if !param.IsOmitted(u.OfEdit) {
		return u.OfEdit
	} else if !param.IsOmitted(u.OfRead) {
		return u.OfRead
	} else if !param.IsOmitted(u.OfWrite) {
		return u.OfWrite
	} else if !param.IsOmitted(u.OfGlob) {
		return u.OfGlob
	} else if !param.IsOmitted(u.OfGrep) {
		return u.OfGrep
	} else if !param.IsOmitted(u.OfWebFetch) {
		return u.OfWebFetch
	} else if !param.IsOmitted(u.OfWebSearch) {
		return u.OfWebSearch
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsAgentToolConfigParamsUnion) GetMaxContentTokens() *int64 {
	if vt := u.OfWebFetch; vt != nil && vt.MaxContentTokens.Valid() {
		return &vt.MaxContentTokens.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsAgentToolConfigParamsUnion) GetUserLocation() *ManagedAgentsUserLocationParam {
	if vt := u.OfWebSearch; vt != nil {
		return &vt.UserLocation
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsAgentToolConfigParamsUnion) GetName() *string {
	if vt := u.OfBash; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfEdit; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfRead; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfWrite; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfGlob; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfGrep; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfWebFetch; vt != nil {
		return (*string)(&vt.Name)
	} else if vt := u.OfWebSearch; vt != nil {
		return (*string)(&vt.Name)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsAgentToolConfigParamsUnion) GetEnabled() *bool {
	if vt := u.OfBash; vt != nil && vt.Enabled.Valid() {
		return &vt.Enabled.Value
	} else if vt := u.OfEdit; vt != nil && vt.Enabled.Valid() {
		return &vt.Enabled.Value
	} else if vt := u.OfRead; vt != nil && vt.Enabled.Valid() {
		return &vt.Enabled.Value
	} else if vt := u.OfWrite; vt != nil && vt.Enabled.Valid() {
		return &vt.Enabled.Value
	} else if vt := u.OfGlob; vt != nil && vt.Enabled.Valid() {
		return &vt.Enabled.Value
	} else if vt := u.OfGrep; vt != nil && vt.Enabled.Valid() {
		return &vt.Enabled.Value
	} else if vt := u.OfWebFetch; vt != nil && vt.Enabled.Valid() {
		return &vt.Enabled.Value
	} else if vt := u.OfWebSearch; vt != nil && vt.Enabled.Valid() {
		return &vt.Enabled.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsAgentToolConfigParamsUnion) GetType() *string {
	if vt := u.OfBash; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfEdit; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfRead; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfWrite; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfGlob; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfGrep; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfWebFetch; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfWebSearch; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u ManagedAgentsAgentToolConfigParamsUnion) GetPermissionPolicy() (res managedAgentsAgentToolConfigParamsUnionPermissionPolicy) {
	if vt := u.OfBash; vt != nil {
		res.any = vt.PermissionPolicy.asAny()
	} else if vt := u.OfEdit; vt != nil {
		res.any = vt.PermissionPolicy.asAny()
	} else if vt := u.OfRead; vt != nil {
		res.any = vt.PermissionPolicy.asAny()
	} else if vt := u.OfWrite; vt != nil {
		res.any = vt.PermissionPolicy.asAny()
	} else if vt := u.OfGlob; vt != nil {
		res.any = vt.PermissionPolicy.asAny()
	} else if vt := u.OfGrep; vt != nil {
		res.any = vt.PermissionPolicy.asAny()
	} else if vt := u.OfWebFetch; vt != nil {
		res.any = vt.PermissionPolicy.asAny()
	} else if vt := u.OfWebSearch; vt != nil {
		res.any = vt.PermissionPolicy.asAny()
	}
	return
}

// Can have the runtime types [*ManagedAgentsAlwaysAllowPolicyParam],
// [*ManagedAgentsAlwaysAskPolicyParam]
type managedAgentsAgentToolConfigParamsUnionPermissionPolicy struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *qoder.ManagedAgentsAlwaysAllowPolicyParam:
//	case *qoder.ManagedAgentsAlwaysAskPolicyParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u managedAgentsAgentToolConfigParamsUnionPermissionPolicy) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u managedAgentsAgentToolConfigParamsUnionPermissionPolicy) GetType() *string {
	switch vt := u.any.(type) {
	case *ManagedAgentsBashToolConfigParamsPermissionPolicyUnion:
		return vt.GetType()
	case *ManagedAgentsEditToolConfigParamsPermissionPolicyUnion:
		return vt.GetType()
	case *ManagedAgentsReadToolConfigParamsPermissionPolicyUnion:
		return vt.GetType()
	case *ManagedAgentsWriteToolConfigParamsPermissionPolicyUnion:
		return vt.GetType()
	case *ManagedAgentsGlobToolConfigParamsPermissionPolicyUnion:
		return vt.GetType()
	case *ManagedAgentsGrepToolConfigParamsPermissionPolicyUnion:
		return vt.GetType()
	case *ManagedAgentsWebFetchToolConfigParamsPermissionPolicyUnion:
		return vt.GetType()
	case *ManagedAgentsWebSearchToolConfigParamsPermissionPolicyUnion:
		return vt.GetType()
	}
	return nil
}

// Returns a pointer to the underlying variant's AllowedDomains property, if
// present.
func (u ManagedAgentsAgentToolConfigParamsUnion) GetAllowedDomains() []string {
	if vt := u.OfWebFetch; vt != nil {
		return vt.AllowedDomains
	} else if vt := u.OfWebSearch; vt != nil {
		return vt.AllowedDomains
	}
	return nil
}

// Returns a pointer to the underlying variant's BlockedDomains property, if
// present.
func (u ManagedAgentsAgentToolConfigParamsUnion) GetBlockedDomains() []string {
	if vt := u.OfWebFetch; vt != nil {
		return vt.BlockedDomains
	} else if vt := u.OfWebSearch; vt != nil {
		return vt.BlockedDomains
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ManagedAgentsAgentToolConfigParamsUnion](
		"type",
		apijson.Discriminator[ManagedAgentsBashToolConfigParams]("bash"),
		apijson.Discriminator[ManagedAgentsEditToolConfigParams]("edit"),
		apijson.Discriminator[ManagedAgentsReadToolConfigParams]("read"),
		apijson.Discriminator[ManagedAgentsWriteToolConfigParams]("write"),
		apijson.Discriminator[ManagedAgentsGlobToolConfigParams]("glob"),
		apijson.Discriminator[ManagedAgentsGrepToolConfigParams]("grep"),
		apijson.Discriminator[ManagedAgentsWebFetchToolConfigParams]("web_fetch"),
		apijson.Discriminator[ManagedAgentsWebSearchToolConfigParams]("web_search"),
	)
}

// Resolved default configuration for agent tools.
type ManagedAgentsAgentToolsetDefaultConfig struct {
	Enabled bool `json:"enabled" api:"required"`
	// Permission policy for tool execution.
	PermissionPolicy ManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion `json:"permission_policy" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled          respjson.Field
		PermissionPolicy respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsAgentToolsetDefaultConfig) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsAgentToolsetDefaultConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion contains all
// possible properties and values from [ManagedAgentsAlwaysAllowPolicy],
// [ManagedAgentsAlwaysAskPolicy].
//
// Use the [ManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion.AsAny]
// method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion struct {
	// Any of "always_allow", "always_ask".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyManagedAgentsAgentToolsetDefaultConfigPermissionPolicy is implemented by
// each variant of
// [ManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion] to add type
// safety for the return type of
// [ManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion.AsAny]
type anyManagedAgentsAgentToolsetDefaultConfigPermissionPolicy interface {
	implManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion()
}

func (ManagedAgentsAlwaysAllowPolicy) implManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion() {
}
func (ManagedAgentsAlwaysAskPolicy) implManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion.AsAny().(type) {
//	case qoder.ManagedAgentsAlwaysAllowPolicy:
//	case qoder.ManagedAgentsAlwaysAskPolicy:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion) AsAny() anyManagedAgentsAgentToolsetDefaultConfigPermissionPolicy {
	switch u.Type {
	case "always_allow":
		return u.AsAlwaysAllow()
	case "always_ask":
		return u.AsAlwaysAsk()
	}
	return nil
}

func (u ManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion) AsAlwaysAllow() (v ManagedAgentsAlwaysAllowPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion) AsAlwaysAsk() (v ManagedAgentsAlwaysAskPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *ManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Default configuration for all tools in a toolset.
type ManagedAgentsAgentToolsetDefaultConfigParams struct {
	// Whether tools are enabled and available to the model by default. Defaults to
	// true if not specified.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Permission policy for tool execution.
	PermissionPolicy ManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion `json:"permission_policy,omitzero"`
	paramObj
}

func (r ManagedAgentsAgentToolsetDefaultConfigParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsAgentToolsetDefaultConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsAgentToolsetDefaultConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion struct {
	OfAlwaysAllow *ManagedAgentsAlwaysAllowPolicyParam `json:",omitzero,inline"`
	OfAlwaysAsk   *ManagedAgentsAlwaysAskPolicyParam   `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAlwaysAllow, u.OfAlwaysAsk)
}
func (u *ManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion) asAny() any {
	if !param.IsOmitted(u.OfAlwaysAllow) {
		return u.OfAlwaysAllow
	} else if !param.IsOmitted(u.OfAlwaysAsk) {
		return u.OfAlwaysAsk
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion) GetType() *string {
	if vt := u.OfAlwaysAllow; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAlwaysAsk; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion](
		"type",
		apijson.Discriminator[ManagedAgentsAlwaysAllowPolicyParam]("always_allow"),
		apijson.Discriminator[ManagedAgentsAlwaysAskPolicyParam]("always_ask"),
	)
}

type ManagedAgentsAgentToolset20260401 struct {
	EnabledTools    []string                            `json:"enabled_tools"`
	DisallowedTools []string                            `json:"disallowed_tools"`
	Configs         []ManagedAgentsAgentToolConfigUnion `json:"configs" api:"required"`
	// Resolved default configuration for agent tools.
	DefaultConfig ManagedAgentsAgentToolsetDefaultConfig `json:"default_config" api:"required"`
	// Any of "agent_toolset_20260401".
	Type ManagedAgentsAgentToolset20260401Type `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EnabledTools    respjson.Field
		DisallowedTools respjson.Field
		Configs         respjson.Field
		DefaultConfig   respjson.Field
		Type            respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsAgentToolset20260401) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsAgentToolset20260401) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsAgentToolset20260401Type string

const (
	ManagedAgentsAgentToolset20260401TypeAgentToolset20260401 ManagedAgentsAgentToolset20260401Type = "agent_toolset_20260401"
)

// Input payload for the `bash` tool of the `agent_toolset_20260401` toolset. All
// fields are optional; a normal invocation supplies `command`, while
// `restart=true` (with no `command`) reboots the runner-side bash session.
type ManagedAgentsAgentToolset20260401BashInput struct {
	// Shell command to execute. Omit only when `restart` is true.
	Command string `json:"command"`
	// When true, restart the persistent bash session instead of running a command.
	// Subsequent calls without `restart` will run against the fresh session.
	Restart bool `json:"restart"`
	// Per-call timeout in milliseconds. Defaults to the runner-wide tool timeout when
	// omitted or zero.
	TimeoutMs int64 `json:"timeout_ms"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Command     respjson.Field
		Restart     respjson.Field
		TimeoutMs   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsAgentToolset20260401BashInput) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsAgentToolset20260401BashInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Input payload for the `edit` tool. Performs a string replacement in the named
// file; by default `old_string` must occur exactly once.
type ManagedAgentsAgentToolset20260401EditInput struct {
	// Path of the file to edit.
	FilePath string `json:"file_path" api:"required"`
	// Replacement text.
	NewString string `json:"new_string" api:"required"`
	// Substring to find and replace.
	OldString string `json:"old_string" api:"required"`
	// When true, replace every occurrence of `old_string` instead of requiring a
	// unique match.
	ReplaceAll bool `json:"replace_all"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		FilePath    respjson.Field
		NewString   respjson.Field
		OldString   respjson.Field
		ReplaceAll  respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsAgentToolset20260401EditInput) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsAgentToolset20260401EditInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Input payload for the `glob` tool. Returns paths matching a doublestar glob
// pattern, newest first.
type ManagedAgentsAgentToolset20260401GlobInput struct {
	// Doublestar glob pattern (e.g. `**/*.go`). Absolute patterns are only permitted
	// when the runner is configured to allow them.
	Pattern string `json:"pattern" api:"required"`
	// Optional directory root to search under. Defaults to the runner's working
	// directory.
	Path string `json:"path"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Pattern     respjson.Field
		Path        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsAgentToolset20260401GlobInput) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsAgentToolset20260401GlobInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Input payload for the `grep` tool. Searches file contents for a regular
// expression, returning matching lines.
type ManagedAgentsAgentToolset20260401GrepInput struct {
	// Regular expression to search for.
	Pattern string `json:"pattern" api:"required"`
	// Optional directory root to search under. Defaults to the runner's working
	// directory.
	Path string `json:"path"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Pattern     respjson.Field
		Path        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsAgentToolset20260401GrepInput) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsAgentToolset20260401GrepInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for built-in agent tools. Use this to enable or disable groups of
// tools available to the agent.
//
// The property Type is required.
type ManagedAgentsAgentToolset20260401Params struct {
	DisallowedTools []string `json:"disallowed_tools,omitzero"`
	EnabledTools    []string `json:"enabled_tools,omitzero"`
	// Any of "agent_toolset_20260401".
	Type ManagedAgentsAgentToolset20260401ParamsType `json:"type,omitzero" api:"required"`
	// Per-tool configuration overrides.
	Configs []ManagedAgentsAgentToolConfigParamsUnion `json:"configs,omitzero"`
	// Default configuration for all tools in a toolset.
	DefaultConfig ManagedAgentsAgentToolsetDefaultConfigParams `json:"default_config,omitzero"`
	paramObj
}

func (r ManagedAgentsAgentToolset20260401Params) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsAgentToolset20260401Params
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsAgentToolset20260401Params) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsAgentToolset20260401ParamsType string

const (
	ManagedAgentsAgentToolset20260401ParamsTypeAgentToolset20260401 ManagedAgentsAgentToolset20260401ParamsType = "agent_toolset_20260401"
)

// Input payload for the `read` tool. Reads file contents relative to the runner's
// working directory (or absolute when the runner permits).
type ManagedAgentsAgentToolset20260401ReadInput struct {
	// Path of the file to read.
	FilePath string `json:"file_path" api:"required"`
	// Optional `[start_line, end_line]` 1-indexed inclusive range. When omitted the
	// entire file is returned. `end_line` of 0 or negative means "to end of file".
	ViewRange []int64 `json:"view_range"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		FilePath    respjson.Field
		ViewRange   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsAgentToolset20260401ReadInput) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsAgentToolset20260401ReadInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Input payload for the `write` tool. Writes (overwriting) the entire file
// contents.
type ManagedAgentsAgentToolset20260401WriteInput struct {
	// Full file contents to write.
	Content string `json:"content" api:"required"`
	// Path of the file to write.
	FilePath string `json:"file_path" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Content     respjson.Field
		FilePath    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsAgentToolset20260401WriteInput) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsAgentToolset20260401WriteInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Tool calls are automatically approved without user confirmation.
type ManagedAgentsAlwaysAllowPolicy struct {
	// Any of "always_allow".
	Type ManagedAgentsAlwaysAllowPolicyType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsAlwaysAllowPolicy) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsAlwaysAllowPolicy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ManagedAgentsAlwaysAllowPolicy to a
// ManagedAgentsAlwaysAllowPolicyParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ManagedAgentsAlwaysAllowPolicyParam.Overrides()
func (r ManagedAgentsAlwaysAllowPolicy) ToParam() ManagedAgentsAlwaysAllowPolicyParam {
	return param.Override[ManagedAgentsAlwaysAllowPolicyParam](json.RawMessage(r.RawJSON()))
}

type ManagedAgentsAlwaysAllowPolicyType string

const (
	ManagedAgentsAlwaysAllowPolicyTypeAlwaysAllow ManagedAgentsAlwaysAllowPolicyType = "always_allow"
)

// Tool calls are automatically approved without user confirmation.
//
// The property Type is required.
type ManagedAgentsAlwaysAllowPolicyParam struct {
	// Any of "always_allow".
	Type ManagedAgentsAlwaysAllowPolicyType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsAlwaysAllowPolicyParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsAlwaysAllowPolicyParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsAlwaysAllowPolicyParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Tool calls require user confirmation before execution.
type ManagedAgentsAlwaysAskPolicy struct {
	// Any of "always_ask".
	Type ManagedAgentsAlwaysAskPolicyType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsAlwaysAskPolicy) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsAlwaysAskPolicy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ManagedAgentsAlwaysAskPolicy to a
// ManagedAgentsAlwaysAskPolicyParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ManagedAgentsAlwaysAskPolicyParam.Overrides()
func (r ManagedAgentsAlwaysAskPolicy) ToParam() ManagedAgentsAlwaysAskPolicyParam {
	return param.Override[ManagedAgentsAlwaysAskPolicyParam](json.RawMessage(r.RawJSON()))
}

type ManagedAgentsAlwaysAskPolicyType string

const (
	ManagedAgentsAlwaysAskPolicyTypeAlwaysAsk ManagedAgentsAlwaysAskPolicyType = "always_ask"
)

// Tool calls require user confirmation before execution.
//
// The property Type is required.
type ManagedAgentsAlwaysAskPolicyParam struct {
	// Any of "always_ask".
	Type ManagedAgentsAlwaysAskPolicyType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsAlwaysAskPolicyParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsAlwaysAskPolicyParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsAlwaysAskPolicyParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A resolved Qoder-managed skill.
type ManagedAgentsQoderSkill struct {
	SkillID string `json:"skill_id" api:"required"`
	// Any of "qoder".
	Type    ManagedAgentsQoderSkillType `json:"type" api:"required"`
	Version string                      `json:"version" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		SkillID     respjson.Field
		Type        respjson.Field
		Version     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsQoderSkill) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsQoderSkill) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsQoderSkillType string

const (
	ManagedAgentsQoderSkillTypeQoder ManagedAgentsQoderSkillType = "qoder"
)

// An Qoder-managed skill.
//
// The properties SkillID, Type are required.
type ManagedAgentsQoderSkillParams struct {
	// Identifier of the Qoder skill (e.g., "xlsx").
	SkillID string `json:"skill_id" api:"required"`
	// Any of "qoder".
	Type ManagedAgentsQoderSkillParamsType `json:"type,omitzero" api:"required"`
	// Version to pin. Defaults to latest if omitted.
	Version param.Opt[string] `json:"version,omitzero"`
	paramObj
}

func (r ManagedAgentsQoderSkillParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsQoderSkillParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsQoderSkillParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsQoderSkillParamsType string

const (
	ManagedAgentsQoderSkillParamsTypeQoder ManagedAgentsQoderSkillParamsType = "qoder"
)

// Configuration for the bash tool.
type ManagedAgentsBashToolConfig struct {
	Enabled bool          `json:"enabled" api:"required"`
	Name    constant.Bash `json:"name" default:"bash"`
	// Permission policy for tool execution.
	PermissionPolicy ManagedAgentsBashToolConfigPermissionPolicyUnion `json:"permission_policy" api:"required"`
	Type             constant.Bash                                    `json:"type" default:"bash"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled          respjson.Field
		Name             respjson.Field
		PermissionPolicy respjson.Field
		Type             respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsBashToolConfig) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsBashToolConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsBashToolConfigPermissionPolicyUnion contains all possible
// properties and values from [ManagedAgentsAlwaysAllowPolicy],
// [ManagedAgentsAlwaysAskPolicy].
//
// Use the [ManagedAgentsBashToolConfigPermissionPolicyUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsBashToolConfigPermissionPolicyUnion struct {
	// Any of "always_allow", "always_ask".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyManagedAgentsBashToolConfigPermissionPolicy is implemented by each
// variant of [ManagedAgentsBashToolConfigPermissionPolicyUnion] to add type
// safety for the return type of
// [ManagedAgentsBashToolConfigPermissionPolicyUnion.AsAny]
type anyManagedAgentsBashToolConfigPermissionPolicy interface {
	implManagedAgentsBashToolConfigPermissionPolicyUnion()
}

func (ManagedAgentsAlwaysAllowPolicy) implManagedAgentsBashToolConfigPermissionPolicyUnion() {
}
func (ManagedAgentsAlwaysAskPolicy) implManagedAgentsBashToolConfigPermissionPolicyUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsBashToolConfigPermissionPolicyUnion.AsAny().(type) {
//	case qoder.ManagedAgentsAlwaysAllowPolicy:
//	case qoder.ManagedAgentsAlwaysAskPolicy:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsBashToolConfigPermissionPolicyUnion) AsAny() anyManagedAgentsBashToolConfigPermissionPolicy {
	switch u.Type {
	case "always_allow":
		return u.AsAlwaysAllow()
	case "always_ask":
		return u.AsAlwaysAsk()
	}
	return nil
}

func (u ManagedAgentsBashToolConfigPermissionPolicyUnion) AsAlwaysAllow() (v ManagedAgentsAlwaysAllowPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsBashToolConfigPermissionPolicyUnion) AsAlwaysAsk() (v ManagedAgentsAlwaysAskPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsBashToolConfigPermissionPolicyUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsBashToolConfigPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration override for the bash tool.
//
// The property Name is required.
type ManagedAgentsBashToolConfigParams struct {
	// Whether this tool is enabled and available to the model. Overrides the
	// default_config setting.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Permission policy for tool execution.
	PermissionPolicy ManagedAgentsBashToolConfigParamsPermissionPolicyUnion `json:"permission_policy,omitzero"`
	// Any of "bash".
	Type ManagedAgentsBashToolConfigParamsType `json:"type,omitzero"`
	// Must be "bash".
	//
	// This field can be elided, and will marshal its zero value as "bash".
	Name constant.Bash `json:"name" default:"bash"`
	paramObj
}

func (r ManagedAgentsBashToolConfigParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsBashToolConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsBashToolConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsBashToolConfigParamsPermissionPolicyUnion struct {
	OfAlwaysAllow *ManagedAgentsAlwaysAllowPolicyParam `json:",omitzero,inline"`
	OfAlwaysAsk   *ManagedAgentsAlwaysAskPolicyParam   `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsBashToolConfigParamsPermissionPolicyUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAlwaysAllow, u.OfAlwaysAsk)
}
func (u *ManagedAgentsBashToolConfigParamsPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsBashToolConfigParamsPermissionPolicyUnion) asAny() any {
	if !param.IsOmitted(u.OfAlwaysAllow) {
		return u.OfAlwaysAllow
	} else if !param.IsOmitted(u.OfAlwaysAsk) {
		return u.OfAlwaysAsk
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsBashToolConfigParamsPermissionPolicyUnion) GetType() *string {
	if vt := u.OfAlwaysAllow; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAlwaysAsk; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ManagedAgentsBashToolConfigParamsPermissionPolicyUnion](
		"type",
		apijson.Discriminator[ManagedAgentsAlwaysAllowPolicyParam]("always_allow"),
		apijson.Discriminator[ManagedAgentsAlwaysAskPolicyParam]("always_ask"),
	)
}

type ManagedAgentsBashToolConfigParamsType string

const (
	ManagedAgentsBashToolConfigParamsTypeBash ManagedAgentsBashToolConfigParamsType = "bash"
)

// A resolved user-created custom skill.
type ManagedAgentsCustomSkill struct {
	SkillID string `json:"skill_id" api:"required"`
	// Any of "custom".
	Type    ManagedAgentsCustomSkillType `json:"type" api:"required"`
	Version string                       `json:"version" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		SkillID     respjson.Field
		Type        respjson.Field
		Version     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsCustomSkill) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsCustomSkill) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsCustomSkillType string

const (
	ManagedAgentsCustomSkillTypeCustom ManagedAgentsCustomSkillType = "custom"
)

// A user-created custom skill.
//
// The properties SkillID, Type are required.
type ManagedAgentsCustomSkillParams struct {
	// Tagged ID of the custom skill (e.g., "skill_01XJ5...").
	SkillID string `json:"skill_id" api:"required"`
	// Any of "custom".
	Type ManagedAgentsCustomSkillParamsType `json:"type,omitzero" api:"required"`
	// Version to pin. Defaults to latest if omitted.
	Version param.Opt[string] `json:"version,omitzero"`
	paramObj
}

func (r ManagedAgentsCustomSkillParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsCustomSkillParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsCustomSkillParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsCustomSkillParamsType string

const (
	ManagedAgentsCustomSkillParamsTypeCustom ManagedAgentsCustomSkillParamsType = "custom"
)

// A custom tool as returned in API responses.
type ManagedAgentsCustomTool struct {
	Description string `json:"description" api:"required"`
	// JSON Schema for custom tool input parameters.
	InputSchema ManagedAgentsCustomToolInputSchema `json:"input_schema" api:"required"`
	Name        string                             `json:"name" api:"required"`
	// Any of "custom".
	Type ManagedAgentsCustomToolType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Description respjson.Field
		InputSchema respjson.Field
		Name        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsCustomTool) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsCustomTool) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsCustomToolType string

const (
	ManagedAgentsCustomToolTypeCustom ManagedAgentsCustomToolType = "custom"
)

// JSON Schema for custom tool input parameters.
type ManagedAgentsCustomToolInputSchema struct {
	Type        constant.Object `json:"type" default:"object"`
	Properties  map[string]any  `json:"properties" api:"nullable"`
	Required    []string        `json:"required" api:"nullable"`
	ExtraFields map[string]any  `json:"" api:"extrafields"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		Properties  respjson.Field
		Required    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsCustomToolInputSchema) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsCustomToolInputSchema) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ManagedAgentsCustomToolInputSchema to a
// ManagedAgentsCustomToolInputSchemaParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ManagedAgentsCustomToolInputSchemaParam.Overrides()
func (r ManagedAgentsCustomToolInputSchema) ToParam() ManagedAgentsCustomToolInputSchemaParam {
	return param.Override[ManagedAgentsCustomToolInputSchemaParam](json.RawMessage(r.RawJSON()))
}

// JSON Schema for custom tool input parameters.
//
// The property Type is required.
type ManagedAgentsCustomToolInputSchemaParam struct {
	Properties map[string]any `json:"properties,omitzero"`
	Required   []string       `json:"required,omitzero"`
	// This field can be elided, and will marshal its zero value as "object".
	Type        constant.Object `json:"type" default:"object"`
	ExtraFields map[string]any  `json:"-"`
	paramObj
}

func (r ManagedAgentsCustomToolInputSchemaParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsCustomToolInputSchemaParam
	return param.MarshalWithExtras(r, (*shadow)(&r), r.ExtraFields)
}
func (r *ManagedAgentsCustomToolInputSchemaParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A custom tool that is executed by the API client rather than the agent. When the
// agent calls this tool, an `agent.custom_tool_use` event is emitted and the
// session goes idle, waiting for the client to provide the result via a
// `user.custom_tool_result` event.
//
// The properties Description, InputSchema, Name, Type are required.
type ManagedAgentsCustomToolParams struct {
	// Description of what the tool does, shown to the agent to help it decide when to
	// use the tool.
	Description string `json:"description" api:"required"`
	// JSON Schema for custom tool input parameters.
	InputSchema ManagedAgentsCustomToolInputSchemaParam `json:"input_schema,omitzero" api:"required"`
	// Unique name for the tool. 1-128 characters; letters, digits, underscores, and
	// hyphens.
	Name string `json:"name" api:"required"`
	// Any of "custom".
	Type ManagedAgentsCustomToolParamsType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsCustomToolParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsCustomToolParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsCustomToolParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsCustomToolParamsType string

const (
	ManagedAgentsCustomToolParamsTypeCustom ManagedAgentsCustomToolParamsType = "custom"
)

// Configuration for the edit tool.
type ManagedAgentsEditToolConfig struct {
	Enabled bool          `json:"enabled" api:"required"`
	Name    constant.Edit `json:"name" default:"edit"`
	// Permission policy for tool execution.
	PermissionPolicy ManagedAgentsEditToolConfigPermissionPolicyUnion `json:"permission_policy" api:"required"`
	Type             constant.Edit                                    `json:"type" default:"edit"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled          respjson.Field
		Name             respjson.Field
		PermissionPolicy respjson.Field
		Type             respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsEditToolConfig) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsEditToolConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsEditToolConfigPermissionPolicyUnion contains all possible
// properties and values from [ManagedAgentsAlwaysAllowPolicy],
// [ManagedAgentsAlwaysAskPolicy].
//
// Use the [ManagedAgentsEditToolConfigPermissionPolicyUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsEditToolConfigPermissionPolicyUnion struct {
	// Any of "always_allow", "always_ask".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyManagedAgentsEditToolConfigPermissionPolicy is implemented by each
// variant of [ManagedAgentsEditToolConfigPermissionPolicyUnion] to add type
// safety for the return type of
// [ManagedAgentsEditToolConfigPermissionPolicyUnion.AsAny]
type anyManagedAgentsEditToolConfigPermissionPolicy interface {
	implManagedAgentsEditToolConfigPermissionPolicyUnion()
}

func (ManagedAgentsAlwaysAllowPolicy) implManagedAgentsEditToolConfigPermissionPolicyUnion() {
}
func (ManagedAgentsAlwaysAskPolicy) implManagedAgentsEditToolConfigPermissionPolicyUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsEditToolConfigPermissionPolicyUnion.AsAny().(type) {
//	case qoder.ManagedAgentsAlwaysAllowPolicy:
//	case qoder.ManagedAgentsAlwaysAskPolicy:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsEditToolConfigPermissionPolicyUnion) AsAny() anyManagedAgentsEditToolConfigPermissionPolicy {
	switch u.Type {
	case "always_allow":
		return u.AsAlwaysAllow()
	case "always_ask":
		return u.AsAlwaysAsk()
	}
	return nil
}

func (u ManagedAgentsEditToolConfigPermissionPolicyUnion) AsAlwaysAllow() (v ManagedAgentsAlwaysAllowPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsEditToolConfigPermissionPolicyUnion) AsAlwaysAsk() (v ManagedAgentsAlwaysAskPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsEditToolConfigPermissionPolicyUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsEditToolConfigPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration override for the edit tool.
//
// The property Name is required.
type ManagedAgentsEditToolConfigParams struct {
	// Whether this tool is enabled and available to the model. Overrides the
	// default_config setting.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Permission policy for tool execution.
	PermissionPolicy ManagedAgentsEditToolConfigParamsPermissionPolicyUnion `json:"permission_policy,omitzero"`
	// Any of "edit".
	Type ManagedAgentsEditToolConfigParamsType `json:"type,omitzero"`
	// Must be "edit".
	//
	// This field can be elided, and will marshal its zero value as "edit".
	Name constant.Edit `json:"name" default:"edit"`
	paramObj
}

func (r ManagedAgentsEditToolConfigParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsEditToolConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsEditToolConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsEditToolConfigParamsPermissionPolicyUnion struct {
	OfAlwaysAllow *ManagedAgentsAlwaysAllowPolicyParam `json:",omitzero,inline"`
	OfAlwaysAsk   *ManagedAgentsAlwaysAskPolicyParam   `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsEditToolConfigParamsPermissionPolicyUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAlwaysAllow, u.OfAlwaysAsk)
}
func (u *ManagedAgentsEditToolConfigParamsPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsEditToolConfigParamsPermissionPolicyUnion) asAny() any {
	if !param.IsOmitted(u.OfAlwaysAllow) {
		return u.OfAlwaysAllow
	} else if !param.IsOmitted(u.OfAlwaysAsk) {
		return u.OfAlwaysAsk
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsEditToolConfigParamsPermissionPolicyUnion) GetType() *string {
	if vt := u.OfAlwaysAllow; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAlwaysAsk; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ManagedAgentsEditToolConfigParamsPermissionPolicyUnion](
		"type",
		apijson.Discriminator[ManagedAgentsAlwaysAllowPolicyParam]("always_allow"),
		apijson.Discriminator[ManagedAgentsAlwaysAskPolicyParam]("always_ask"),
	)
}

type ManagedAgentsEditToolConfigParamsType string

const (
	ManagedAgentsEditToolConfigParamsTypeEdit ManagedAgentsEditToolConfigParamsType = "edit"
)

// High effort. Favors reasoning depth.
type ManagedAgentsEffortHigh struct {
	// Any of "high".
	Type ManagedAgentsEffortHighType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsEffortHigh) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsEffortHigh) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ManagedAgentsEffortHigh to a
// ManagedAgentsEffortHighParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ManagedAgentsEffortHighParam.Overrides()
func (r ManagedAgentsEffortHigh) ToParam() ManagedAgentsEffortHighParam {
	return param.Override[ManagedAgentsEffortHighParam](json.RawMessage(r.RawJSON()))
}

type ManagedAgentsEffortHighType string

const (
	ManagedAgentsEffortHighTypeHigh ManagedAgentsEffortHighType = "high"
)

// High effort. Favors reasoning depth.
//
// The property Type is required.
type ManagedAgentsEffortHighParam struct {
	// Any of "high".
	Type ManagedAgentsEffortHighType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsEffortHighParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsEffortHighParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsEffortHighParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Low effort. Favors latency over reasoning depth.
type ManagedAgentsEffortLow struct {
	// Any of "low".
	Type ManagedAgentsEffortLowType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsEffortLow) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsEffortLow) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ManagedAgentsEffortLow to a
// ManagedAgentsEffortLowParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ManagedAgentsEffortLowParam.Overrides()
func (r ManagedAgentsEffortLow) ToParam() ManagedAgentsEffortLowParam {
	return param.Override[ManagedAgentsEffortLowParam](json.RawMessage(r.RawJSON()))
}

type ManagedAgentsEffortLowType string

const (
	ManagedAgentsEffortLowTypeLow ManagedAgentsEffortLowType = "low"
)

// Low effort. Favors latency over reasoning depth.
//
// The property Type is required.
type ManagedAgentsEffortLowParam struct {
	// Any of "low".
	Type ManagedAgentsEffortLowType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsEffortLowParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsEffortLowParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsEffortLowParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Maximum effort. Favors reasoning depth over latency.
type ManagedAgentsEffortMax struct {
	// Any of "max".
	Type ManagedAgentsEffortMaxType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsEffortMax) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsEffortMax) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ManagedAgentsEffortMax to a
// ManagedAgentsEffortMaxParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ManagedAgentsEffortMaxParam.Overrides()
func (r ManagedAgentsEffortMax) ToParam() ManagedAgentsEffortMaxParam {
	return param.Override[ManagedAgentsEffortMaxParam](json.RawMessage(r.RawJSON()))
}

type ManagedAgentsEffortMaxType string

const (
	ManagedAgentsEffortMaxTypeMax ManagedAgentsEffortMaxType = "max"
)

// Maximum effort. Favors reasoning depth over latency.
//
// The property Type is required.
type ManagedAgentsEffortMaxParam struct {
	// Any of "max".
	Type ManagedAgentsEffortMaxType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsEffortMaxParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsEffortMaxParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsEffortMaxParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Medium effort. Balances latency and reasoning depth.
type ManagedAgentsEffortMedium struct {
	// Any of "medium".
	Type ManagedAgentsEffortMediumType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsEffortMedium) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsEffortMedium) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ManagedAgentsEffortMedium to a
// ManagedAgentsEffortMediumParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ManagedAgentsEffortMediumParam.Overrides()
func (r ManagedAgentsEffortMedium) ToParam() ManagedAgentsEffortMediumParam {
	return param.Override[ManagedAgentsEffortMediumParam](json.RawMessage(r.RawJSON()))
}

type ManagedAgentsEffortMediumType string

const (
	ManagedAgentsEffortMediumTypeMedium ManagedAgentsEffortMediumType = "medium"
)

// Medium effort. Balances latency and reasoning depth.
//
// The property Type is required.
type ManagedAgentsEffortMediumParam struct {
	// Any of "medium".
	Type ManagedAgentsEffortMediumType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsEffortMediumParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsEffortMediumParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsEffortMediumParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Extra-high effort. Not all models accept this level.
type ManagedAgentsEffortXhigh struct {
	// Any of "xhigh".
	Type ManagedAgentsEffortXhighType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsEffortXhigh) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsEffortXhigh) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ManagedAgentsEffortXhigh to a
// ManagedAgentsEffortXhighParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ManagedAgentsEffortXhighParam.Overrides()
func (r ManagedAgentsEffortXhigh) ToParam() ManagedAgentsEffortXhighParam {
	return param.Override[ManagedAgentsEffortXhighParam](json.RawMessage(r.RawJSON()))
}

type ManagedAgentsEffortXhighType string

const (
	ManagedAgentsEffortXhighTypeXhigh ManagedAgentsEffortXhighType = "xhigh"
)

// Extra-high effort. Not all models accept this level.
//
// The property Type is required.
type ManagedAgentsEffortXhighParam struct {
	// Any of "xhigh".
	Type ManagedAgentsEffortXhighType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsEffortXhighParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsEffortXhighParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsEffortXhighParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for the glob tool.
type ManagedAgentsGlobToolConfig struct {
	Enabled bool          `json:"enabled" api:"required"`
	Name    constant.Glob `json:"name" default:"glob"`
	// Permission policy for tool execution.
	PermissionPolicy ManagedAgentsGlobToolConfigPermissionPolicyUnion `json:"permission_policy" api:"required"`
	Type             constant.Glob                                    `json:"type" default:"glob"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled          respjson.Field
		Name             respjson.Field
		PermissionPolicy respjson.Field
		Type             respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsGlobToolConfig) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsGlobToolConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsGlobToolConfigPermissionPolicyUnion contains all possible
// properties and values from [ManagedAgentsAlwaysAllowPolicy],
// [ManagedAgentsAlwaysAskPolicy].
//
// Use the [ManagedAgentsGlobToolConfigPermissionPolicyUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsGlobToolConfigPermissionPolicyUnion struct {
	// Any of "always_allow", "always_ask".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyManagedAgentsGlobToolConfigPermissionPolicy is implemented by each
// variant of [ManagedAgentsGlobToolConfigPermissionPolicyUnion] to add type
// safety for the return type of
// [ManagedAgentsGlobToolConfigPermissionPolicyUnion.AsAny]
type anyManagedAgentsGlobToolConfigPermissionPolicy interface {
	implManagedAgentsGlobToolConfigPermissionPolicyUnion()
}

func (ManagedAgentsAlwaysAllowPolicy) implManagedAgentsGlobToolConfigPermissionPolicyUnion() {
}
func (ManagedAgentsAlwaysAskPolicy) implManagedAgentsGlobToolConfigPermissionPolicyUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsGlobToolConfigPermissionPolicyUnion.AsAny().(type) {
//	case qoder.ManagedAgentsAlwaysAllowPolicy:
//	case qoder.ManagedAgentsAlwaysAskPolicy:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsGlobToolConfigPermissionPolicyUnion) AsAny() anyManagedAgentsGlobToolConfigPermissionPolicy {
	switch u.Type {
	case "always_allow":
		return u.AsAlwaysAllow()
	case "always_ask":
		return u.AsAlwaysAsk()
	}
	return nil
}

func (u ManagedAgentsGlobToolConfigPermissionPolicyUnion) AsAlwaysAllow() (v ManagedAgentsAlwaysAllowPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsGlobToolConfigPermissionPolicyUnion) AsAlwaysAsk() (v ManagedAgentsAlwaysAskPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsGlobToolConfigPermissionPolicyUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsGlobToolConfigPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration override for the glob tool.
//
// The property Name is required.
type ManagedAgentsGlobToolConfigParams struct {
	// Whether this tool is enabled and available to the model. Overrides the
	// default_config setting.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Permission policy for tool execution.
	PermissionPolicy ManagedAgentsGlobToolConfigParamsPermissionPolicyUnion `json:"permission_policy,omitzero"`
	// Any of "glob".
	Type ManagedAgentsGlobToolConfigParamsType `json:"type,omitzero"`
	// Must be "glob".
	//
	// This field can be elided, and will marshal its zero value as "glob".
	Name constant.Glob `json:"name" default:"glob"`
	paramObj
}

func (r ManagedAgentsGlobToolConfigParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsGlobToolConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsGlobToolConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsGlobToolConfigParamsPermissionPolicyUnion struct {
	OfAlwaysAllow *ManagedAgentsAlwaysAllowPolicyParam `json:",omitzero,inline"`
	OfAlwaysAsk   *ManagedAgentsAlwaysAskPolicyParam   `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsGlobToolConfigParamsPermissionPolicyUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAlwaysAllow, u.OfAlwaysAsk)
}
func (u *ManagedAgentsGlobToolConfigParamsPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsGlobToolConfigParamsPermissionPolicyUnion) asAny() any {
	if !param.IsOmitted(u.OfAlwaysAllow) {
		return u.OfAlwaysAllow
	} else if !param.IsOmitted(u.OfAlwaysAsk) {
		return u.OfAlwaysAsk
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsGlobToolConfigParamsPermissionPolicyUnion) GetType() *string {
	if vt := u.OfAlwaysAllow; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAlwaysAsk; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ManagedAgentsGlobToolConfigParamsPermissionPolicyUnion](
		"type",
		apijson.Discriminator[ManagedAgentsAlwaysAllowPolicyParam]("always_allow"),
		apijson.Discriminator[ManagedAgentsAlwaysAskPolicyParam]("always_ask"),
	)
}

type ManagedAgentsGlobToolConfigParamsType string

const (
	ManagedAgentsGlobToolConfigParamsTypeGlob ManagedAgentsGlobToolConfigParamsType = "glob"
)

// Configuration for the grep tool.
type ManagedAgentsGrepToolConfig struct {
	Enabled bool          `json:"enabled" api:"required"`
	Name    constant.Grep `json:"name" default:"grep"`
	// Permission policy for tool execution.
	PermissionPolicy ManagedAgentsGrepToolConfigPermissionPolicyUnion `json:"permission_policy" api:"required"`
	Type             constant.Grep                                    `json:"type" default:"grep"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled          respjson.Field
		Name             respjson.Field
		PermissionPolicy respjson.Field
		Type             respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsGrepToolConfig) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsGrepToolConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsGrepToolConfigPermissionPolicyUnion contains all possible
// properties and values from [ManagedAgentsAlwaysAllowPolicy],
// [ManagedAgentsAlwaysAskPolicy].
//
// Use the [ManagedAgentsGrepToolConfigPermissionPolicyUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsGrepToolConfigPermissionPolicyUnion struct {
	// Any of "always_allow", "always_ask".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyManagedAgentsGrepToolConfigPermissionPolicy is implemented by each
// variant of [ManagedAgentsGrepToolConfigPermissionPolicyUnion] to add type
// safety for the return type of
// [ManagedAgentsGrepToolConfigPermissionPolicyUnion.AsAny]
type anyManagedAgentsGrepToolConfigPermissionPolicy interface {
	implManagedAgentsGrepToolConfigPermissionPolicyUnion()
}

func (ManagedAgentsAlwaysAllowPolicy) implManagedAgentsGrepToolConfigPermissionPolicyUnion() {
}
func (ManagedAgentsAlwaysAskPolicy) implManagedAgentsGrepToolConfigPermissionPolicyUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsGrepToolConfigPermissionPolicyUnion.AsAny().(type) {
//	case qoder.ManagedAgentsAlwaysAllowPolicy:
//	case qoder.ManagedAgentsAlwaysAskPolicy:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsGrepToolConfigPermissionPolicyUnion) AsAny() anyManagedAgentsGrepToolConfigPermissionPolicy {
	switch u.Type {
	case "always_allow":
		return u.AsAlwaysAllow()
	case "always_ask":
		return u.AsAlwaysAsk()
	}
	return nil
}

func (u ManagedAgentsGrepToolConfigPermissionPolicyUnion) AsAlwaysAllow() (v ManagedAgentsAlwaysAllowPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsGrepToolConfigPermissionPolicyUnion) AsAlwaysAsk() (v ManagedAgentsAlwaysAskPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsGrepToolConfigPermissionPolicyUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsGrepToolConfigPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration override for the grep tool.
//
// The property Name is required.
type ManagedAgentsGrepToolConfigParams struct {
	// Whether this tool is enabled and available to the model. Overrides the
	// default_config setting.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Permission policy for tool execution.
	PermissionPolicy ManagedAgentsGrepToolConfigParamsPermissionPolicyUnion `json:"permission_policy,omitzero"`
	// Any of "grep".
	Type ManagedAgentsGrepToolConfigParamsType `json:"type,omitzero"`
	// Must be "grep".
	//
	// This field can be elided, and will marshal its zero value as "grep".
	Name constant.Grep `json:"name" default:"grep"`
	paramObj
}

func (r ManagedAgentsGrepToolConfigParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsGrepToolConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsGrepToolConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsGrepToolConfigParamsPermissionPolicyUnion struct {
	OfAlwaysAllow *ManagedAgentsAlwaysAllowPolicyParam `json:",omitzero,inline"`
	OfAlwaysAsk   *ManagedAgentsAlwaysAskPolicyParam   `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsGrepToolConfigParamsPermissionPolicyUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAlwaysAllow, u.OfAlwaysAsk)
}
func (u *ManagedAgentsGrepToolConfigParamsPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsGrepToolConfigParamsPermissionPolicyUnion) asAny() any {
	if !param.IsOmitted(u.OfAlwaysAllow) {
		return u.OfAlwaysAllow
	} else if !param.IsOmitted(u.OfAlwaysAsk) {
		return u.OfAlwaysAsk
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsGrepToolConfigParamsPermissionPolicyUnion) GetType() *string {
	if vt := u.OfAlwaysAllow; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAlwaysAsk; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ManagedAgentsGrepToolConfigParamsPermissionPolicyUnion](
		"type",
		apijson.Discriminator[ManagedAgentsAlwaysAllowPolicyParam]("always_allow"),
		apijson.Discriminator[ManagedAgentsAlwaysAskPolicyParam]("always_ask"),
	)
}

type ManagedAgentsGrepToolConfigParamsType string

const (
	ManagedAgentsGrepToolConfigParamsTypeGrep ManagedAgentsGrepToolConfigParamsType = "grep"
)

// URL-based MCP server connection as returned in API responses.
type ManagedAgentsMCPServerURLDefinition struct {
	Name string `json:"name" api:"required"`
	// Any of "url".
	Type ManagedAgentsMCPServerURLDefinitionType `json:"type" api:"required"`
	URL  string                                  `json:"url" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		Type        respjson.Field
		URL         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsMCPServerURLDefinition) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsMCPServerURLDefinition) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsMCPServerURLDefinitionType string

const (
	ManagedAgentsMCPServerURLDefinitionTypeURL ManagedAgentsMCPServerURLDefinitionType = "url"
)

// Resolved configuration for a specific MCP tool.
type ManagedAgentsMCPToolConfig struct {
	Enabled bool   `json:"enabled" api:"required"`
	Name    string `json:"name" api:"required"`
	// Permission policy for tool execution.
	PermissionPolicy ManagedAgentsMCPToolConfigPermissionPolicyUnion `json:"permission_policy" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled          respjson.Field
		Name             respjson.Field
		PermissionPolicy respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsMCPToolConfig) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsMCPToolConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsMCPToolConfigPermissionPolicyUnion contains all possible
// properties and values from [ManagedAgentsAlwaysAllowPolicy],
// [ManagedAgentsAlwaysAskPolicy].
//
// Use the [ManagedAgentsMCPToolConfigPermissionPolicyUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsMCPToolConfigPermissionPolicyUnion struct {
	// Any of "always_allow", "always_ask".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyManagedAgentsMCPToolConfigPermissionPolicy is implemented by each variant
// of [ManagedAgentsMCPToolConfigPermissionPolicyUnion] to add type safety for
// the return type of [ManagedAgentsMCPToolConfigPermissionPolicyUnion.AsAny]
type anyManagedAgentsMCPToolConfigPermissionPolicy interface {
	implManagedAgentsMCPToolConfigPermissionPolicyUnion()
}

func (ManagedAgentsAlwaysAllowPolicy) implManagedAgentsMCPToolConfigPermissionPolicyUnion() {}
func (ManagedAgentsAlwaysAskPolicy) implManagedAgentsMCPToolConfigPermissionPolicyUnion()   {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsMCPToolConfigPermissionPolicyUnion.AsAny().(type) {
//	case qoder.ManagedAgentsAlwaysAllowPolicy:
//	case qoder.ManagedAgentsAlwaysAskPolicy:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsMCPToolConfigPermissionPolicyUnion) AsAny() anyManagedAgentsMCPToolConfigPermissionPolicy {
	switch u.Type {
	case "always_allow":
		return u.AsAlwaysAllow()
	case "always_ask":
		return u.AsAlwaysAsk()
	}
	return nil
}

func (u ManagedAgentsMCPToolConfigPermissionPolicyUnion) AsAlwaysAllow() (v ManagedAgentsAlwaysAllowPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsMCPToolConfigPermissionPolicyUnion) AsAlwaysAsk() (v ManagedAgentsAlwaysAskPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsMCPToolConfigPermissionPolicyUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsMCPToolConfigPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration override for a specific MCP tool.
//
// The property Name is required.
type ManagedAgentsMCPToolConfigParams struct {
	// Name of the MCP tool to configure. 1-128 characters.
	Name string `json:"name" api:"required"`
	// Whether this tool is enabled. Overrides the `default_config` setting.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Permission policy for tool execution.
	PermissionPolicy ManagedAgentsMCPToolConfigParamsPermissionPolicyUnion `json:"permission_policy,omitzero"`
	paramObj
}

func (r ManagedAgentsMCPToolConfigParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsMCPToolConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsMCPToolConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsMCPToolConfigParamsPermissionPolicyUnion struct {
	OfAlwaysAllow *ManagedAgentsAlwaysAllowPolicyParam `json:",omitzero,inline"`
	OfAlwaysAsk   *ManagedAgentsAlwaysAskPolicyParam   `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsMCPToolConfigParamsPermissionPolicyUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAlwaysAllow, u.OfAlwaysAsk)
}
func (u *ManagedAgentsMCPToolConfigParamsPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsMCPToolConfigParamsPermissionPolicyUnion) asAny() any {
	if !param.IsOmitted(u.OfAlwaysAllow) {
		return u.OfAlwaysAllow
	} else if !param.IsOmitted(u.OfAlwaysAsk) {
		return u.OfAlwaysAsk
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsMCPToolConfigParamsPermissionPolicyUnion) GetType() *string {
	if vt := u.OfAlwaysAllow; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAlwaysAsk; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ManagedAgentsMCPToolConfigParamsPermissionPolicyUnion](
		"type",
		apijson.Discriminator[ManagedAgentsAlwaysAllowPolicyParam]("always_allow"),
		apijson.Discriminator[ManagedAgentsAlwaysAskPolicyParam]("always_ask"),
	)
}

type ManagedAgentsMCPToolset struct {
	Configs []ManagedAgentsMCPToolConfig `json:"configs" api:"required"`
	// Resolved default configuration for all tools from an MCP server.
	DefaultConfig ManagedAgentsMCPToolsetDefaultConfig `json:"default_config" api:"required"`
	MCPServerName string                               `json:"mcp_server_name" api:"required"`
	// Any of "mcp_toolset".
	Type ManagedAgentsMCPToolsetType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Configs       respjson.Field
		DefaultConfig respjson.Field
		MCPServerName respjson.Field
		Type          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsMCPToolset) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsMCPToolset) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsMCPToolsetType string

const (
	ManagedAgentsMCPToolsetTypeMCPToolset ManagedAgentsMCPToolsetType = "mcp_toolset"
)

// Resolved default configuration for all tools from an MCP server.
type ManagedAgentsMCPToolsetDefaultConfig struct {
	Enabled bool `json:"enabled" api:"required"`
	// Permission policy for tool execution.
	PermissionPolicy ManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion `json:"permission_policy" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled          respjson.Field
		PermissionPolicy respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsMCPToolsetDefaultConfig) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsMCPToolsetDefaultConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion contains all
// possible properties and values from [ManagedAgentsAlwaysAllowPolicy],
// [ManagedAgentsAlwaysAskPolicy].
//
// Use the [ManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion.AsAny]
// method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion struct {
	// Any of "always_allow", "always_ask".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyManagedAgentsMCPToolsetDefaultConfigPermissionPolicy is implemented by
// each variant of [ManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion]
// to add type safety for the return type of
// [ManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion.AsAny]
type anyManagedAgentsMCPToolsetDefaultConfigPermissionPolicy interface {
	implManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion()
}

func (ManagedAgentsAlwaysAllowPolicy) implManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion() {
}
func (ManagedAgentsAlwaysAskPolicy) implManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion.AsAny().(type) {
//	case qoder.ManagedAgentsAlwaysAllowPolicy:
//	case qoder.ManagedAgentsAlwaysAskPolicy:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion) AsAny() anyManagedAgentsMCPToolsetDefaultConfigPermissionPolicy {
	switch u.Type {
	case "always_allow":
		return u.AsAlwaysAllow()
	case "always_ask":
		return u.AsAlwaysAsk()
	}
	return nil
}

func (u ManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion) AsAlwaysAllow() (v ManagedAgentsAlwaysAllowPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion) AsAlwaysAsk() (v ManagedAgentsAlwaysAskPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *ManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Default configuration for all tools from an MCP server.
type ManagedAgentsMCPToolsetDefaultConfigParams struct {
	// Whether tools are enabled by default. Defaults to true if not specified.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Permission policy for tool execution.
	PermissionPolicy ManagedAgentsMCPToolsetDefaultConfigParamsPermissionPolicyUnion `json:"permission_policy,omitzero"`
	paramObj
}

func (r ManagedAgentsMCPToolsetDefaultConfigParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsMCPToolsetDefaultConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsMCPToolsetDefaultConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsMCPToolsetDefaultConfigParamsPermissionPolicyUnion struct {
	OfAlwaysAllow *ManagedAgentsAlwaysAllowPolicyParam `json:",omitzero,inline"`
	OfAlwaysAsk   *ManagedAgentsAlwaysAskPolicyParam   `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsMCPToolsetDefaultConfigParamsPermissionPolicyUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAlwaysAllow, u.OfAlwaysAsk)
}
func (u *ManagedAgentsMCPToolsetDefaultConfigParamsPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsMCPToolsetDefaultConfigParamsPermissionPolicyUnion) asAny() any {
	if !param.IsOmitted(u.OfAlwaysAllow) {
		return u.OfAlwaysAllow
	} else if !param.IsOmitted(u.OfAlwaysAsk) {
		return u.OfAlwaysAsk
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsMCPToolsetDefaultConfigParamsPermissionPolicyUnion) GetType() *string {
	if vt := u.OfAlwaysAllow; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAlwaysAsk; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ManagedAgentsMCPToolsetDefaultConfigParamsPermissionPolicyUnion](
		"type",
		apijson.Discriminator[ManagedAgentsAlwaysAllowPolicyParam]("always_allow"),
		apijson.Discriminator[ManagedAgentsAlwaysAskPolicyParam]("always_ask"),
	)
}

// Configuration for tools from an MCP server defined in `mcp_servers`.
//
// The properties MCPServerName, Type are required.
type ManagedAgentsMCPToolsetParams struct {
	// Name of the MCP server. Must match a server name from the mcp_servers array.
	// 1-255 characters.
	MCPServerName string `json:"mcp_server_name" api:"required"`
	// Any of "mcp_toolset".
	Type ManagedAgentsMCPToolsetParamsType `json:"type,omitzero" api:"required"`
	// Per-tool configuration overrides.
	Configs []ManagedAgentsMCPToolConfigParams `json:"configs,omitzero"`
	// Default configuration for all tools from an MCP server.
	DefaultConfig ManagedAgentsMCPToolsetDefaultConfigParams `json:"default_config,omitzero"`
	paramObj
}

func (r ManagedAgentsMCPToolsetParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsMCPToolsetParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsMCPToolsetParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsMCPToolsetParamsType string

const (
	ManagedAgentsMCPToolsetParamsTypeMCPToolset ManagedAgentsMCPToolsetParamsType = "mcp_toolset"
)

// The model that will power your agent.
//
// See [models](https://docs.qoder.com/cloud-agents/api/models/list) for additional
// details and options.
type ManagedAgentsModel = string

// Model identifier and configuration.
type ManagedAgentsModelConfig struct {
	ContextWindow int64 `json:"context_window"`

	// The model that will power your agent.
	//
	// See [models](https://docs.qoder.com/cloud-agents/api/models/list) for additional
	// details and options.
	ID ManagedAgentsModel `json:"id" api:"required"`
	// How hard the model works on each turn. Sets `output_config.effort` on every
	// Messages call the session makes.
	Effort ManagedAgentsModelConfigEffortUnion `json:"effort"`
	// Geographic region for model inference. When unset, requests fall through to the
	// workspace's default_inference_geo.
	InferenceGeo string `json:"inference_geo"`
	// Inference speed mode. `fast` provides significantly faster output token
	// generation at premium pricing. Not all models support `fast`; invalid
	// combinations are rejected at create time.
	//
	// Any of "standard", "fast".
	Speed ManagedAgentsModelConfigSpeed `json:"speed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ContextWindow respjson.Field

		ID           respjson.Field
		Effort       respjson.Field
		InferenceGeo respjson.Field
		Speed        respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsModelConfig) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsModelConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsModelConfigEffortUnion contains all possible properties and
// values from [ManagedAgentsEffortLow], [ManagedAgentsEffortMedium],
// [ManagedAgentsEffortHigh], [ManagedAgentsEffortXhigh],
// [ManagedAgentsEffortMax].
//
// Use the [ManagedAgentsModelConfigEffortUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsModelConfigEffortUnion struct {
	// Any of "low", "medium", "high", "xhigh", "max".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyManagedAgentsModelConfigEffort is implemented by each variant of
// [ManagedAgentsModelConfigEffortUnion] to add type safety for the return type
// of [ManagedAgentsModelConfigEffortUnion.AsAny]
type anyManagedAgentsModelConfigEffort interface {
	implManagedAgentsModelConfigEffortUnion()
}

func (ManagedAgentsEffortLow) implManagedAgentsModelConfigEffortUnion()    {}
func (ManagedAgentsEffortMedium) implManagedAgentsModelConfigEffortUnion() {}
func (ManagedAgentsEffortHigh) implManagedAgentsModelConfigEffortUnion()   {}
func (ManagedAgentsEffortXhigh) implManagedAgentsModelConfigEffortUnion()  {}
func (ManagedAgentsEffortMax) implManagedAgentsModelConfigEffortUnion()    {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsModelConfigEffortUnion.AsAny().(type) {
//	case qoder.ManagedAgentsEffortLow:
//	case qoder.ManagedAgentsEffortMedium:
//	case qoder.ManagedAgentsEffortHigh:
//	case qoder.ManagedAgentsEffortXhigh:
//	case qoder.ManagedAgentsEffortMax:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsModelConfigEffortUnion) AsAny() anyManagedAgentsModelConfigEffort {
	switch u.Type {
	case "low":
		return u.AsLow()
	case "medium":
		return u.AsMedium()
	case "high":
		return u.AsHigh()
	case "xhigh":
		return u.AsXhigh()
	case "max":
		return u.AsMax()
	}
	return nil
}

func (u ManagedAgentsModelConfigEffortUnion) AsLow() (v ManagedAgentsEffortLow) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsModelConfigEffortUnion) AsMedium() (v ManagedAgentsEffortMedium) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsModelConfigEffortUnion) AsHigh() (v ManagedAgentsEffortHigh) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsModelConfigEffortUnion) AsXhigh() (v ManagedAgentsEffortXhigh) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsModelConfigEffortUnion) AsMax() (v ManagedAgentsEffortMax) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsModelConfigEffortUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsModelConfigEffortUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Inference speed mode. `fast` provides significantly faster output token
// generation at premium pricing. Not all models support `fast`; invalid
// combinations are rejected at create time.
type ManagedAgentsModelConfigSpeed string

const (
	ManagedAgentsModelConfigSpeedStandard ManagedAgentsModelConfigSpeed = "standard"
	ManagedAgentsModelConfigSpeedFast     ManagedAgentsModelConfigSpeed = "fast"
)

// An object that defines additional configuration control over model use
//
// The property ID is required.
type ManagedAgentsModelConfigParams struct {
	ContextWindow param.Opt[int64] `json:"context_window,omitzero"`

	// The model that will power your agent.
	//
	// See [models](https://docs.qoder.com/cloud-agents/api/models/list) for additional
	// details and options.
	ID ManagedAgentsModel `json:"id,omitzero" api:"required"`
	// Geographic region for model inference. When unset, requests fall through to the
	// workspace's default_inference_geo. On update, `model` is whole-object
	// replacement — omitting inference_geo clears it.
	InferenceGeo param.Opt[string] `json:"inference_geo,omitzero"`
	// How hard the model works on each inference call. Accepts a bare level string
	// (`"high"`) or `{"type": "high"}`. On create, omitting it resolves the per-model
	// default; on update, omitting it leaves the stored value unchanged.
	Effort ManagedAgentsModelConfigParamsEffortUnion `json:"effort,omitzero"`
	// Inference speed mode. `fast` provides significantly faster output token
	// generation at premium pricing. Not all models support `fast`; invalid
	// combinations are rejected at create time.
	//
	// Any of "standard", "fast".
	Speed ManagedAgentsModelConfigParamsSpeed `json:"speed,omitzero"`
	paramObj
}

func (r ManagedAgentsModelConfigParams) MarshalJSON() (data []byte, err error) {
	if r.Speed != "" {
		return nil, fmt.Errorf("Qoder does not support model.speed")
	}
	type shadow ManagedAgentsModelConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsModelConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsModelConfigParamsEffortUnion struct {
	// Check if union is this variant with
	// !param.IsOmitted(union.OfManagedAgentsModelConfigsEffortManagedAgentsEffortLevel)
	OfManagedAgentsModelConfigsEffortManagedAgentsEffortLevel param.Opt[string]               `json:",omitzero,inline"`
	OfManagedAgentsEffortLow                                  *ManagedAgentsEffortLowParam    `json:",omitzero,inline"`
	OfManagedAgentsEffortMedium                               *ManagedAgentsEffortMediumParam `json:",omitzero,inline"`
	OfManagedAgentsEffortHigh                                 *ManagedAgentsEffortHighParam   `json:",omitzero,inline"`
	OfManagedAgentsEffortXhigh                                *ManagedAgentsEffortXhighParam  `json:",omitzero,inline"`
	OfManagedAgentsEffortMax                                  *ManagedAgentsEffortMaxParam    `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsModelConfigParamsEffortUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfManagedAgentsModelConfigsEffortManagedAgentsEffortLevel,
		u.OfManagedAgentsEffortLow,
		u.OfManagedAgentsEffortMedium,
		u.OfManagedAgentsEffortHigh,
		u.OfManagedAgentsEffortXhigh,
		u.OfManagedAgentsEffortMax)
}
func (u *ManagedAgentsModelConfigParamsEffortUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsModelConfigParamsEffortUnion) asAny() any {
	if !param.IsOmitted(u.OfManagedAgentsModelConfigsEffortManagedAgentsEffortLevel) {
		return &u.OfManagedAgentsModelConfigsEffortManagedAgentsEffortLevel
	} else if !param.IsOmitted(u.OfManagedAgentsEffortLow) {
		return u.OfManagedAgentsEffortLow
	} else if !param.IsOmitted(u.OfManagedAgentsEffortMedium) {
		return u.OfManagedAgentsEffortMedium
	} else if !param.IsOmitted(u.OfManagedAgentsEffortHigh) {
		return u.OfManagedAgentsEffortHigh
	} else if !param.IsOmitted(u.OfManagedAgentsEffortXhigh) {
		return u.OfManagedAgentsEffortXhigh
	} else if !param.IsOmitted(u.OfManagedAgentsEffortMax) {
		return u.OfManagedAgentsEffortMax
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsModelConfigParamsEffortUnion) GetType() *string {
	if vt := u.OfManagedAgentsEffortLow; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfManagedAgentsEffortMedium; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfManagedAgentsEffortHigh; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfManagedAgentsEffortXhigh; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfManagedAgentsEffortMax; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// How hard the model works on each turn. Higher levels favor reasoning depth
// over latency. Not all models accept every level; invalid combinations are
// rejected at create time.
type ManagedAgentsModelConfigParamsEffortManagedAgentsEffortLevel string

const (
	ManagedAgentsModelConfigParamsEffortManagedAgentsEffortLevelLow    ManagedAgentsModelConfigParamsEffortManagedAgentsEffortLevel = "low"
	ManagedAgentsModelConfigParamsEffortManagedAgentsEffortLevelMedium ManagedAgentsModelConfigParamsEffortManagedAgentsEffortLevel = "medium"
	ManagedAgentsModelConfigParamsEffortManagedAgentsEffortLevelHigh   ManagedAgentsModelConfigParamsEffortManagedAgentsEffortLevel = "high"
	ManagedAgentsModelConfigParamsEffortManagedAgentsEffortLevelXhigh  ManagedAgentsModelConfigParamsEffortManagedAgentsEffortLevel = "xhigh"
	ManagedAgentsModelConfigParamsEffortManagedAgentsEffortLevelMax    ManagedAgentsModelConfigParamsEffortManagedAgentsEffortLevel = "max"
)

// Inference speed mode. `fast` provides significantly faster output token
// generation at premium pricing. Not all models support `fast`; invalid
// combinations are rejected at create time.
type ManagedAgentsModelConfigParamsSpeed string

const (
	ManagedAgentsModelConfigParamsSpeedStandard ManagedAgentsModelConfigParamsSpeed = "standard"
	ManagedAgentsModelConfigParamsSpeedFast     ManagedAgentsModelConfigParamsSpeed = "fast"
)

// Sentinel roster entry meaning "the agent that owns this configuration". Resolved
// server-side to a concrete agent reference.
//
// The property Type is required.
type ManagedAgentsMultiagentSelfParams struct {
	// Any of "self".
	Type ManagedAgentsMultiagentSelfParamsType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsMultiagentSelfParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsMultiagentSelfParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsMultiagentSelfParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsMultiagentSelfParamsType string

const (
	ManagedAgentsMultiagentSelfParamsTypeSelf ManagedAgentsMultiagentSelfParamsType = "self"
)

// Configuration for the read tool.
type ManagedAgentsReadToolConfig struct {
	Enabled bool          `json:"enabled" api:"required"`
	Name    constant.Read `json:"name" default:"read"`
	// Permission policy for tool execution.
	PermissionPolicy ManagedAgentsReadToolConfigPermissionPolicyUnion `json:"permission_policy" api:"required"`
	Type             constant.Read                                    `json:"type" default:"read"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled          respjson.Field
		Name             respjson.Field
		PermissionPolicy respjson.Field
		Type             respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsReadToolConfig) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsReadToolConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsReadToolConfigPermissionPolicyUnion contains all possible
// properties and values from [ManagedAgentsAlwaysAllowPolicy],
// [ManagedAgentsAlwaysAskPolicy].
//
// Use the [ManagedAgentsReadToolConfigPermissionPolicyUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsReadToolConfigPermissionPolicyUnion struct {
	// Any of "always_allow", "always_ask".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyManagedAgentsReadToolConfigPermissionPolicy is implemented by each
// variant of [ManagedAgentsReadToolConfigPermissionPolicyUnion] to add type
// safety for the return type of
// [ManagedAgentsReadToolConfigPermissionPolicyUnion.AsAny]
type anyManagedAgentsReadToolConfigPermissionPolicy interface {
	implManagedAgentsReadToolConfigPermissionPolicyUnion()
}

func (ManagedAgentsAlwaysAllowPolicy) implManagedAgentsReadToolConfigPermissionPolicyUnion() {
}
func (ManagedAgentsAlwaysAskPolicy) implManagedAgentsReadToolConfigPermissionPolicyUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsReadToolConfigPermissionPolicyUnion.AsAny().(type) {
//	case qoder.ManagedAgentsAlwaysAllowPolicy:
//	case qoder.ManagedAgentsAlwaysAskPolicy:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsReadToolConfigPermissionPolicyUnion) AsAny() anyManagedAgentsReadToolConfigPermissionPolicy {
	switch u.Type {
	case "always_allow":
		return u.AsAlwaysAllow()
	case "always_ask":
		return u.AsAlwaysAsk()
	}
	return nil
}

func (u ManagedAgentsReadToolConfigPermissionPolicyUnion) AsAlwaysAllow() (v ManagedAgentsAlwaysAllowPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsReadToolConfigPermissionPolicyUnion) AsAlwaysAsk() (v ManagedAgentsAlwaysAskPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsReadToolConfigPermissionPolicyUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsReadToolConfigPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration override for the read tool.
//
// The property Name is required.
type ManagedAgentsReadToolConfigParams struct {
	// Whether this tool is enabled and available to the model. Overrides the
	// default_config setting.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Permission policy for tool execution.
	PermissionPolicy ManagedAgentsReadToolConfigParamsPermissionPolicyUnion `json:"permission_policy,omitzero"`
	// Any of "read".
	Type ManagedAgentsReadToolConfigParamsType `json:"type,omitzero"`
	// Must be "read".
	//
	// This field can be elided, and will marshal its zero value as "read".
	Name constant.Read `json:"name" default:"read"`
	paramObj
}

func (r ManagedAgentsReadToolConfigParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsReadToolConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsReadToolConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsReadToolConfigParamsPermissionPolicyUnion struct {
	OfAlwaysAllow *ManagedAgentsAlwaysAllowPolicyParam `json:",omitzero,inline"`
	OfAlwaysAsk   *ManagedAgentsAlwaysAskPolicyParam   `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsReadToolConfigParamsPermissionPolicyUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAlwaysAllow, u.OfAlwaysAsk)
}
func (u *ManagedAgentsReadToolConfigParamsPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsReadToolConfigParamsPermissionPolicyUnion) asAny() any {
	if !param.IsOmitted(u.OfAlwaysAllow) {
		return u.OfAlwaysAllow
	} else if !param.IsOmitted(u.OfAlwaysAsk) {
		return u.OfAlwaysAsk
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsReadToolConfigParamsPermissionPolicyUnion) GetType() *string {
	if vt := u.OfAlwaysAllow; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAlwaysAsk; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ManagedAgentsReadToolConfigParamsPermissionPolicyUnion](
		"type",
		apijson.Discriminator[ManagedAgentsAlwaysAllowPolicyParam]("always_allow"),
		apijson.Discriminator[ManagedAgentsAlwaysAskPolicyParam]("always_ask"),
	)
}

type ManagedAgentsReadToolConfigParamsType string

const (
	ManagedAgentsReadToolConfigParamsTypeRead ManagedAgentsReadToolConfigParamsType = "read"
)

// Resolved `agent` definition for a single `session_thread`. Snapshot of the agent
// at thread creation time. The multiagent roster is not repeated here; read it
// from `Session.agent`.
type ManagedAgentsSessionThreadAgent struct {
	ID          string                                `json:"id" api:"required"`
	Description string                                `json:"description" api:"required"`
	MCPServers  []ManagedAgentsMCPServerURLDefinition `json:"mcp_servers" api:"required"`
	// Model identifier and configuration.
	Model  ManagedAgentsModelConfig                    `json:"model" api:"required"`
	Name   string                                      `json:"name" api:"required"`
	Skills []ManagedAgentsSessionThreadAgentSkillUnion `json:"skills" api:"required"`
	System string                                      `json:"system" api:"required"`
	Tools  []ManagedAgentsSessionThreadAgentToolUnion  `json:"tools" api:"required"`
	// Any of "agent".
	Type    ManagedAgentsSessionThreadAgentType `json:"type" api:"required"`
	Version int64                               `json:"version" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Description respjson.Field
		MCPServers  respjson.Field
		Model       respjson.Field
		Name        respjson.Field
		Skills      respjson.Field
		System      respjson.Field
		Tools       respjson.Field
		Type        respjson.Field
		Version     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionThreadAgent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionThreadAgent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsSessionThreadAgentSkillUnion contains all possible properties
// and values from [ManagedAgentsQoderSkill],
// [ManagedAgentsCustomSkill].
//
// Use the [ManagedAgentsSessionThreadAgentSkillUnion.AsAny] method to switch
// on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsSessionThreadAgentSkillUnion struct {
	SkillID string `json:"skill_id"`
	// Any of "qoder", "custom".
	Type    string `json:"type"`
	Version string `json:"version"`
	JSON    struct {
		SkillID respjson.Field
		Type    respjson.Field
		Version respjson.Field
		raw     string
	} `json:"-"`
}

// anyManagedAgentsSessionThreadAgentSkill is implemented by each variant of
// [ManagedAgentsSessionThreadAgentSkillUnion] to add type safety for the
// return type of [ManagedAgentsSessionThreadAgentSkillUnion.AsAny]
type anyManagedAgentsSessionThreadAgentSkill interface {
	implManagedAgentsSessionThreadAgentSkillUnion()
}

func (ManagedAgentsQoderSkill) implManagedAgentsSessionThreadAgentSkillUnion()  {}
func (ManagedAgentsCustomSkill) implManagedAgentsSessionThreadAgentSkillUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsSessionThreadAgentSkillUnion.AsAny().(type) {
//	case qoder.ManagedAgentsQoderSkill:
//	case qoder.ManagedAgentsCustomSkill:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsSessionThreadAgentSkillUnion) AsAny() anyManagedAgentsSessionThreadAgentSkill {
	switch u.Type {
	case "qoder":
		return u.AsQoder()
	case "custom":
		return u.AsCustom()
	}
	return nil
}

func (u ManagedAgentsSessionThreadAgentSkillUnion) AsQoder() (v ManagedAgentsQoderSkill) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionThreadAgentSkillUnion) AsCustom() (v ManagedAgentsCustomSkill) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsSessionThreadAgentSkillUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsSessionThreadAgentSkillUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsSessionThreadAgentToolUnion contains all possible properties
// and values from [ManagedAgentsAgentToolset20260401],
// [ManagedAgentsMCPToolset], [ManagedAgentsCustomTool].
//
// Use the [ManagedAgentsSessionThreadAgentToolUnion.AsAny] method to switch on
// the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsSessionThreadAgentToolUnion struct {
	EnabledTools    []string `json:"enabled_tools"`
	DisallowedTools []string `json:"disallowed_tools"`
	// This field is a union of [[]ManagedAgentsAgentToolConfigUnion],
	// [[]ManagedAgentsMCPToolConfig]
	Configs ManagedAgentsSessionThreadAgentToolUnionConfigs `json:"configs"`
	// This field is a union of [ManagedAgentsAgentToolsetDefaultConfig],
	// [ManagedAgentsMCPToolsetDefaultConfig]
	DefaultConfig ManagedAgentsSessionThreadAgentToolUnionDefaultConfig `json:"default_config"`
	// Any of "agent_toolset_20260401", "mcp_toolset", "custom".
	Type string `json:"type"`
	// This field is from variant [ManagedAgentsMCPToolset].
	MCPServerName string `json:"mcp_server_name"`
	// This field is from variant [ManagedAgentsCustomTool].
	Description string `json:"description"`
	// This field is from variant [ManagedAgentsCustomTool].
	InputSchema ManagedAgentsCustomToolInputSchema `json:"input_schema"`
	// This field is from variant [ManagedAgentsCustomTool].
	Name string `json:"name"`
	JSON struct {
		EnabledTools    respjson.Field
		DisallowedTools respjson.Field
		Configs         respjson.Field
		DefaultConfig   respjson.Field
		Type            respjson.Field
		MCPServerName   respjson.Field
		Description     respjson.Field
		InputSchema     respjson.Field
		Name            respjson.Field
		raw             string
	} `json:"-"`
}

// anyManagedAgentsSessionThreadAgentTool is implemented by each variant of
// [ManagedAgentsSessionThreadAgentToolUnion] to add type safety for the return
// type of [ManagedAgentsSessionThreadAgentToolUnion.AsAny]
type anyManagedAgentsSessionThreadAgentTool interface {
	implManagedAgentsSessionThreadAgentToolUnion()
}

func (ManagedAgentsAgentToolset20260401) implManagedAgentsSessionThreadAgentToolUnion() {}
func (ManagedAgentsMCPToolset) implManagedAgentsSessionThreadAgentToolUnion()           {}
func (ManagedAgentsCustomTool) implManagedAgentsSessionThreadAgentToolUnion()           {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsSessionThreadAgentToolUnion.AsAny().(type) {
//	case qoder.ManagedAgentsAgentToolset20260401:
//	case qoder.ManagedAgentsMCPToolset:
//	case qoder.ManagedAgentsCustomTool:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsSessionThreadAgentToolUnion) AsAny() anyManagedAgentsSessionThreadAgentTool {
	switch u.Type {
	case "agent_toolset_20260401":
		return u.AsAgentToolset20260401()
	case "mcp_toolset":
		return u.AsMCPToolset()
	case "custom":
		return u.AsCustom()
	}
	return nil
}

func (u ManagedAgentsSessionThreadAgentToolUnion) AsAgentToolset20260401() (v ManagedAgentsAgentToolset20260401) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionThreadAgentToolUnion) AsMCPToolset() (v ManagedAgentsMCPToolset) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionThreadAgentToolUnion) AsCustom() (v ManagedAgentsCustomTool) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsSessionThreadAgentToolUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsSessionThreadAgentToolUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsSessionThreadAgentToolUnionConfigs is an implicit subunion of
// [ManagedAgentsSessionThreadAgentToolUnion].
// ManagedAgentsSessionThreadAgentToolUnionConfigs provides convenient access
// to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsSessionThreadAgentToolUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfManagedAgentsAgentToolConfigArray
// OfManagedAgentsMCPToolConfigArray]
type ManagedAgentsSessionThreadAgentToolUnionConfigs struct {
	// This field will be present if the value is a
	// [[]ManagedAgentsAgentToolConfigUnion] instead of an object.
	OfManagedAgentsAgentToolConfigArray []ManagedAgentsAgentToolConfigUnion `json:",inline"`
	// This field will be present if the value is a [[]ManagedAgentsMCPToolConfig]
	// instead of an object.
	OfManagedAgentsMCPToolConfigArray []ManagedAgentsMCPToolConfig `json:",inline"`
	JSON                              struct {
		OfManagedAgentsAgentToolConfigArray respjson.Field
		OfManagedAgentsMCPToolConfigArray   respjson.Field
		raw                                 string
	} `json:"-"`
}

func (r *ManagedAgentsSessionThreadAgentToolUnionConfigs) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsSessionThreadAgentToolUnionDefaultConfig is an implicit
// subunion of [ManagedAgentsSessionThreadAgentToolUnion].
// ManagedAgentsSessionThreadAgentToolUnionDefaultConfig provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsSessionThreadAgentToolUnion].
type ManagedAgentsSessionThreadAgentToolUnionDefaultConfig struct {
	Enabled bool `json:"enabled"`
	// This field is a union of
	// [ManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion],
	// [ManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion]
	PermissionPolicy ManagedAgentsSessionThreadAgentToolUnionDefaultConfigPermissionPolicy `json:"permission_policy"`
	JSON             struct {
		Enabled          respjson.Field
		PermissionPolicy respjson.Field
		raw              string
	} `json:"-"`
}

func (r *ManagedAgentsSessionThreadAgentToolUnionDefaultConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsSessionThreadAgentToolUnionDefaultConfigPermissionPolicy is an
// implicit subunion of [ManagedAgentsSessionThreadAgentToolUnion].
// ManagedAgentsSessionThreadAgentToolUnionDefaultConfigPermissionPolicy
// provides convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsSessionThreadAgentToolUnion].
type ManagedAgentsSessionThreadAgentToolUnionDefaultConfigPermissionPolicy struct {
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

func (r *ManagedAgentsSessionThreadAgentToolUnionDefaultConfigPermissionPolicy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSessionThreadAgentType string

const (
	ManagedAgentsSessionThreadAgentTypeAgent ManagedAgentsSessionThreadAgentType = "agent"
)

func ManagedAgentsSkillParamsOfQoder(skillID string) ManagedAgentsSkillParamsUnion {
	var qoder ManagedAgentsQoderSkillParams
	qoder.SkillID = skillID
	return ManagedAgentsSkillParamsUnion{OfQoder: &qoder}
}

func ManagedAgentsSkillParamsOfCustom(skillID string) ManagedAgentsSkillParamsUnion {
	var custom ManagedAgentsCustomSkillParams
	custom.SkillID = skillID
	return ManagedAgentsSkillParamsUnion{OfCustom: &custom}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsSkillParamsUnion struct {
	OfQoder  *ManagedAgentsQoderSkillParams  `json:",omitzero,inline"`
	OfCustom *ManagedAgentsCustomSkillParams `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsSkillParamsUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfQoder, u.OfCustom)
}
func (u *ManagedAgentsSkillParamsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsSkillParamsUnion) asAny() any {
	if !param.IsOmitted(u.OfQoder) {
		return u.OfQoder
	} else if !param.IsOmitted(u.OfCustom) {
		return u.OfCustom
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsSkillParamsUnion) GetSkillID() *string {
	if vt := u.OfQoder; vt != nil {
		return (*string)(&vt.SkillID)
	} else if vt := u.OfCustom; vt != nil {
		return (*string)(&vt.SkillID)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsSkillParamsUnion) GetType() *string {
	if vt := u.OfQoder; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCustom; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsSkillParamsUnion) GetVersion() *string {
	if vt := u.OfQoder; vt != nil && vt.Version.Valid() {
		return &vt.Version.Value
	} else if vt := u.OfCustom; vt != nil && vt.Version.Valid() {
		return &vt.Version.Value
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ManagedAgentsSkillParamsUnion](
		"type",
		apijson.Discriminator[ManagedAgentsQoderSkillParams]("qoder"),
		apijson.Discriminator[ManagedAgentsCustomSkillParams]("custom"),
	)
}

// URL-based MCP server connection.
//
// The properties Name, Type, URL are required.
type ManagedAgentsURLMCPServerParams struct {
	// Unique name for this server, referenced by mcp_toolset configurations. 1-255
	// characters.
	Name string `json:"name" api:"required"`
	// Any of "url".
	Type ManagedAgentsURLMCPServerParamsType `json:"type,omitzero" api:"required"`
	// Endpoint URL for the MCP server.
	URL string `json:"url" api:"required"`
	paramObj
}

func (r ManagedAgentsURLMCPServerParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsURLMCPServerParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsURLMCPServerParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsURLMCPServerParamsType string

const (
	ManagedAgentsURLMCPServerParamsTypeURL ManagedAgentsURLMCPServerParamsType = "url"
)

// Approximate user location for search result localization.
type ManagedAgentsUserLocation struct {
	// Location precision. Only "approximate" is supported.
	Type constant.Approximate `json:"type" default:"approximate"`
	// City name.
	City string `json:"city" api:"nullable"`
	// Two-letter ISO 3166-1 country code, uppercase.
	Country string `json:"country" api:"nullable"`
	// Region or state name.
	Region string `json:"region" api:"nullable"`
	// IANA timezone identifier, e.g. "America/Los_Angeles".
	Timezone string `json:"timezone" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		City        respjson.Field
		Country     respjson.Field
		Region      respjson.Field
		Timezone    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsUserLocation) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsUserLocation) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ManagedAgentsUserLocation to a
// ManagedAgentsUserLocationParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ManagedAgentsUserLocationParam.Overrides()
func (r ManagedAgentsUserLocation) ToParam() ManagedAgentsUserLocationParam {
	return param.Override[ManagedAgentsUserLocationParam](json.RawMessage(r.RawJSON()))
}

// Approximate user location for search result localization.
//
// The property Type is required.
type ManagedAgentsUserLocationParam struct {
	// City name.
	City param.Opt[string] `json:"city,omitzero"`
	// Two-letter ISO 3166-1 country code, uppercase.
	Country param.Opt[string] `json:"country,omitzero"`
	// Region or state name.
	Region param.Opt[string] `json:"region,omitzero"`
	// IANA timezone identifier, e.g. "America/Los_Angeles".
	Timezone param.Opt[string] `json:"timezone,omitzero"`
	// Location precision. Only "approximate" is supported.
	//
	// This field can be elided, and will marshal its zero value as "approximate".
	Type constant.Approximate `json:"type" default:"approximate"`
	paramObj
}

func (r ManagedAgentsUserLocationParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsUserLocationParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsUserLocationParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration for the web_fetch tool.
type ManagedAgentsWebFetchToolConfig struct {
	Enabled bool              `json:"enabled" api:"required"`
	Name    constant.WebFetch `json:"name" default:"web_fetch"`
	// Permission policy for tool execution.
	PermissionPolicy ManagedAgentsWebFetchToolConfigPermissionPolicyUnion `json:"permission_policy" api:"required"`
	Type             constant.WebFetch                                    `json:"type" default:"web_fetch"`
	AllowedDomains   []string                                             `json:"allowed_domains"`
	BlockedDomains   []string                                             `json:"blocked_domains"`
	MaxContentTokens int64                                                `json:"max_content_tokens" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled          respjson.Field
		Name             respjson.Field
		PermissionPolicy respjson.Field
		Type             respjson.Field
		AllowedDomains   respjson.Field
		BlockedDomains   respjson.Field
		MaxContentTokens respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsWebFetchToolConfig) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsWebFetchToolConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsWebFetchToolConfigPermissionPolicyUnion contains all possible
// properties and values from [ManagedAgentsAlwaysAllowPolicy],
// [ManagedAgentsAlwaysAskPolicy].
//
// Use the [ManagedAgentsWebFetchToolConfigPermissionPolicyUnion.AsAny] method
// to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsWebFetchToolConfigPermissionPolicyUnion struct {
	// Any of "always_allow", "always_ask".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyManagedAgentsWebFetchToolConfigPermissionPolicy is implemented by each
// variant of [ManagedAgentsWebFetchToolConfigPermissionPolicyUnion] to add
// type safety for the return type of
// [ManagedAgentsWebFetchToolConfigPermissionPolicyUnion.AsAny]
type anyManagedAgentsWebFetchToolConfigPermissionPolicy interface {
	implManagedAgentsWebFetchToolConfigPermissionPolicyUnion()
}

func (ManagedAgentsAlwaysAllowPolicy) implManagedAgentsWebFetchToolConfigPermissionPolicyUnion() {
}
func (ManagedAgentsAlwaysAskPolicy) implManagedAgentsWebFetchToolConfigPermissionPolicyUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsWebFetchToolConfigPermissionPolicyUnion.AsAny().(type) {
//	case qoder.ManagedAgentsAlwaysAllowPolicy:
//	case qoder.ManagedAgentsAlwaysAskPolicy:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsWebFetchToolConfigPermissionPolicyUnion) AsAny() anyManagedAgentsWebFetchToolConfigPermissionPolicy {
	switch u.Type {
	case "always_allow":
		return u.AsAlwaysAllow()
	case "always_ask":
		return u.AsAlwaysAsk()
	}
	return nil
}

func (u ManagedAgentsWebFetchToolConfigPermissionPolicyUnion) AsAlwaysAllow() (v ManagedAgentsAlwaysAllowPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsWebFetchToolConfigPermissionPolicyUnion) AsAlwaysAsk() (v ManagedAgentsAlwaysAskPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsWebFetchToolConfigPermissionPolicyUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsWebFetchToolConfigPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration override for the web_fetch tool.
//
// The property Name is required.
type ManagedAgentsWebFetchToolConfigParams struct {
	// Whether this tool is enabled and available to the model. Overrides the
	// default_config setting.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Maximum number of tokens of fetched text content to include in context per call.
	// Does not apply to binary content such as PDFs.
	MaxContentTokens param.Opt[int64] `json:"max_content_tokens,omitzero"`
	// Permission policy for tool execution.
	PermissionPolicy ManagedAgentsWebFetchToolConfigParamsPermissionPolicyUnion `json:"permission_policy,omitzero"`
	// Only fetch URLs whose host is one of these domains or a subdomain of one. Each
	// entry is a plain hostname like "docs.example.com" (no scheme, port, or path). At
	// most 64 entries; an empty list is rejected (omit the field instead). Cannot be
	// combined with blocked_domains.
	AllowedDomains []string `json:"allowed_domains,omitzero"`
	// Never fetch URLs whose host is one of these domains or a subdomain of one. Each
	// entry is a plain hostname like "ads.example.com" (no scheme, port, or path). At
	// most 64 entries; an empty list is rejected (omit the field instead). Cannot be
	// combined with allowed_domains.
	BlockedDomains []string `json:"blocked_domains,omitzero"`
	// Any of "web_fetch".
	Type ManagedAgentsWebFetchToolConfigParamsType `json:"type,omitzero"`
	// Must be "web_fetch".
	//
	// This field can be elided, and will marshal its zero value as "web_fetch".
	Name constant.WebFetch `json:"name" default:"web_fetch"`
	paramObj
}

func (r ManagedAgentsWebFetchToolConfigParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsWebFetchToolConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsWebFetchToolConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsWebFetchToolConfigParamsPermissionPolicyUnion struct {
	OfAlwaysAllow *ManagedAgentsAlwaysAllowPolicyParam `json:",omitzero,inline"`
	OfAlwaysAsk   *ManagedAgentsAlwaysAskPolicyParam   `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsWebFetchToolConfigParamsPermissionPolicyUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAlwaysAllow, u.OfAlwaysAsk)
}
func (u *ManagedAgentsWebFetchToolConfigParamsPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsWebFetchToolConfigParamsPermissionPolicyUnion) asAny() any {
	if !param.IsOmitted(u.OfAlwaysAllow) {
		return u.OfAlwaysAllow
	} else if !param.IsOmitted(u.OfAlwaysAsk) {
		return u.OfAlwaysAsk
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsWebFetchToolConfigParamsPermissionPolicyUnion) GetType() *string {
	if vt := u.OfAlwaysAllow; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAlwaysAsk; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ManagedAgentsWebFetchToolConfigParamsPermissionPolicyUnion](
		"type",
		apijson.Discriminator[ManagedAgentsAlwaysAllowPolicyParam]("always_allow"),
		apijson.Discriminator[ManagedAgentsAlwaysAskPolicyParam]("always_ask"),
	)
}

type ManagedAgentsWebFetchToolConfigParamsType string

const (
	ManagedAgentsWebFetchToolConfigParamsTypeWebFetch ManagedAgentsWebFetchToolConfigParamsType = "web_fetch"
)

// Configuration for the web_search tool.
type ManagedAgentsWebSearchToolConfig struct {
	Enabled bool               `json:"enabled" api:"required"`
	Name    constant.WebSearch `json:"name" default:"web_search"`
	// Permission policy for tool execution.
	PermissionPolicy ManagedAgentsWebSearchToolConfigPermissionPolicyUnion `json:"permission_policy" api:"required"`
	Type             constant.WebSearch                                    `json:"type" default:"web_search"`
	AllowedDomains   []string                                              `json:"allowed_domains"`
	BlockedDomains   []string                                              `json:"blocked_domains"`
	// Approximate user location for search result localization.
	UserLocation ManagedAgentsUserLocation `json:"user_location" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled          respjson.Field
		Name             respjson.Field
		PermissionPolicy respjson.Field
		Type             respjson.Field
		AllowedDomains   respjson.Field
		BlockedDomains   respjson.Field
		UserLocation     respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsWebSearchToolConfig) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsWebSearchToolConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsWebSearchToolConfigPermissionPolicyUnion contains all possible
// properties and values from [ManagedAgentsAlwaysAllowPolicy],
// [ManagedAgentsAlwaysAskPolicy].
//
// Use the [ManagedAgentsWebSearchToolConfigPermissionPolicyUnion.AsAny] method
// to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsWebSearchToolConfigPermissionPolicyUnion struct {
	// Any of "always_allow", "always_ask".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyManagedAgentsWebSearchToolConfigPermissionPolicy is implemented by each
// variant of [ManagedAgentsWebSearchToolConfigPermissionPolicyUnion] to add
// type safety for the return type of
// [ManagedAgentsWebSearchToolConfigPermissionPolicyUnion.AsAny]
type anyManagedAgentsWebSearchToolConfigPermissionPolicy interface {
	implManagedAgentsWebSearchToolConfigPermissionPolicyUnion()
}

func (ManagedAgentsAlwaysAllowPolicy) implManagedAgentsWebSearchToolConfigPermissionPolicyUnion() {
}
func (ManagedAgentsAlwaysAskPolicy) implManagedAgentsWebSearchToolConfigPermissionPolicyUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsWebSearchToolConfigPermissionPolicyUnion.AsAny().(type) {
//	case qoder.ManagedAgentsAlwaysAllowPolicy:
//	case qoder.ManagedAgentsAlwaysAskPolicy:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsWebSearchToolConfigPermissionPolicyUnion) AsAny() anyManagedAgentsWebSearchToolConfigPermissionPolicy {
	switch u.Type {
	case "always_allow":
		return u.AsAlwaysAllow()
	case "always_ask":
		return u.AsAlwaysAsk()
	}
	return nil
}

func (u ManagedAgentsWebSearchToolConfigPermissionPolicyUnion) AsAlwaysAllow() (v ManagedAgentsAlwaysAllowPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsWebSearchToolConfigPermissionPolicyUnion) AsAlwaysAsk() (v ManagedAgentsAlwaysAskPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsWebSearchToolConfigPermissionPolicyUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *ManagedAgentsWebSearchToolConfigPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration override for the web_search tool.
//
// The property Name is required.
type ManagedAgentsWebSearchToolConfigParams struct {
	// Whether this tool is enabled and available to the model. Overrides the
	// default_config setting.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Permission policy for tool execution.
	PermissionPolicy ManagedAgentsWebSearchToolConfigParamsPermissionPolicyUnion `json:"permission_policy,omitzero"`
	// Only return search results whose host is one of these domains or a subdomain of
	// one. Each entry is a plain hostname like "docs.example.com" (no scheme or port;
	// an optional path suffix is accepted). At most 64 entries; an empty list is
	// rejected (omit the field instead). Cannot be combined with blocked_domains.
	AllowedDomains []string `json:"allowed_domains,omitzero"`
	// Never return search results whose host is one of these domains or a subdomain of
	// one. Each entry is a plain hostname like "ads.example.com" (no scheme or port;
	// an optional path suffix is accepted). At most 64 entries; an empty list is
	// rejected (omit the field instead). Cannot be combined with allowed_domains.
	BlockedDomains []string `json:"blocked_domains,omitzero"`
	// Any of "web_search".
	Type ManagedAgentsWebSearchToolConfigParamsType `json:"type,omitzero"`
	// Approximate user location for search result localization.
	UserLocation ManagedAgentsUserLocationParam `json:"user_location,omitzero"`
	// Must be "web_search".
	//
	// This field can be elided, and will marshal its zero value as "web_search".
	Name constant.WebSearch `json:"name" default:"web_search"`
	paramObj
}

func (r ManagedAgentsWebSearchToolConfigParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsWebSearchToolConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsWebSearchToolConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsWebSearchToolConfigParamsPermissionPolicyUnion struct {
	OfAlwaysAllow *ManagedAgentsAlwaysAllowPolicyParam `json:",omitzero,inline"`
	OfAlwaysAsk   *ManagedAgentsAlwaysAskPolicyParam   `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsWebSearchToolConfigParamsPermissionPolicyUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAlwaysAllow, u.OfAlwaysAsk)
}
func (u *ManagedAgentsWebSearchToolConfigParamsPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsWebSearchToolConfigParamsPermissionPolicyUnion) asAny() any {
	if !param.IsOmitted(u.OfAlwaysAllow) {
		return u.OfAlwaysAllow
	} else if !param.IsOmitted(u.OfAlwaysAsk) {
		return u.OfAlwaysAsk
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsWebSearchToolConfigParamsPermissionPolicyUnion) GetType() *string {
	if vt := u.OfAlwaysAllow; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAlwaysAsk; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ManagedAgentsWebSearchToolConfigParamsPermissionPolicyUnion](
		"type",
		apijson.Discriminator[ManagedAgentsAlwaysAllowPolicyParam]("always_allow"),
		apijson.Discriminator[ManagedAgentsAlwaysAskPolicyParam]("always_ask"),
	)
}

type ManagedAgentsWebSearchToolConfigParamsType string

const (
	ManagedAgentsWebSearchToolConfigParamsTypeWebSearch ManagedAgentsWebSearchToolConfigParamsType = "web_search"
)

// Configuration for the write tool.
type ManagedAgentsWriteToolConfig struct {
	Enabled bool           `json:"enabled" api:"required"`
	Name    constant.Write `json:"name" default:"write"`
	// Permission policy for tool execution.
	PermissionPolicy ManagedAgentsWriteToolConfigPermissionPolicyUnion `json:"permission_policy" api:"required"`
	Type             constant.Write                                    `json:"type" default:"write"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Enabled          respjson.Field
		Name             respjson.Field
		PermissionPolicy respjson.Field
		Type             respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsWriteToolConfig) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsWriteToolConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsWriteToolConfigPermissionPolicyUnion contains all possible
// properties and values from [ManagedAgentsAlwaysAllowPolicy],
// [ManagedAgentsAlwaysAskPolicy].
//
// Use the [ManagedAgentsWriteToolConfigPermissionPolicyUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsWriteToolConfigPermissionPolicyUnion struct {
	// Any of "always_allow", "always_ask".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyManagedAgentsWriteToolConfigPermissionPolicy is implemented by each
// variant of [ManagedAgentsWriteToolConfigPermissionPolicyUnion] to add type
// safety for the return type of
// [ManagedAgentsWriteToolConfigPermissionPolicyUnion.AsAny]
type anyManagedAgentsWriteToolConfigPermissionPolicy interface {
	implManagedAgentsWriteToolConfigPermissionPolicyUnion()
}

func (ManagedAgentsAlwaysAllowPolicy) implManagedAgentsWriteToolConfigPermissionPolicyUnion() {
}
func (ManagedAgentsAlwaysAskPolicy) implManagedAgentsWriteToolConfigPermissionPolicyUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsWriteToolConfigPermissionPolicyUnion.AsAny().(type) {
//	case qoder.ManagedAgentsAlwaysAllowPolicy:
//	case qoder.ManagedAgentsAlwaysAskPolicy:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsWriteToolConfigPermissionPolicyUnion) AsAny() anyManagedAgentsWriteToolConfigPermissionPolicy {
	switch u.Type {
	case "always_allow":
		return u.AsAlwaysAllow()
	case "always_ask":
		return u.AsAlwaysAsk()
	}
	return nil
}

func (u ManagedAgentsWriteToolConfigPermissionPolicyUnion) AsAlwaysAllow() (v ManagedAgentsAlwaysAllowPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsWriteToolConfigPermissionPolicyUnion) AsAlwaysAsk() (v ManagedAgentsAlwaysAskPolicy) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsWriteToolConfigPermissionPolicyUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsWriteToolConfigPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Configuration override for the write tool.
//
// The property Name is required.
type ManagedAgentsWriteToolConfigParams struct {
	// Whether this tool is enabled and available to the model. Overrides the
	// default_config setting.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Permission policy for tool execution.
	PermissionPolicy ManagedAgentsWriteToolConfigParamsPermissionPolicyUnion `json:"permission_policy,omitzero"`
	// Any of "write".
	Type ManagedAgentsWriteToolConfigParamsType `json:"type,omitzero"`
	// Must be "write".
	//
	// This field can be elided, and will marshal its zero value as "write".
	Name constant.Write `json:"name" default:"write"`
	paramObj
}

func (r ManagedAgentsWriteToolConfigParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsWriteToolConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsWriteToolConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsWriteToolConfigParamsPermissionPolicyUnion struct {
	OfAlwaysAllow *ManagedAgentsAlwaysAllowPolicyParam `json:",omitzero,inline"`
	OfAlwaysAsk   *ManagedAgentsAlwaysAskPolicyParam   `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsWriteToolConfigParamsPermissionPolicyUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAlwaysAllow, u.OfAlwaysAsk)
}
func (u *ManagedAgentsWriteToolConfigParamsPermissionPolicyUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsWriteToolConfigParamsPermissionPolicyUnion) asAny() any {
	if !param.IsOmitted(u.OfAlwaysAllow) {
		return u.OfAlwaysAllow
	} else if !param.IsOmitted(u.OfAlwaysAsk) {
		return u.OfAlwaysAsk
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsWriteToolConfigParamsPermissionPolicyUnion) GetType() *string {
	if vt := u.OfAlwaysAllow; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfAlwaysAsk; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ManagedAgentsWriteToolConfigParamsPermissionPolicyUnion](
		"type",
		apijson.Discriminator[ManagedAgentsAlwaysAllowPolicyParam]("always_allow"),
		apijson.Discriminator[ManagedAgentsAlwaysAskPolicyParam]("always_ask"),
	)
}

type ManagedAgentsWriteToolConfigParamsType string

const (
	ManagedAgentsWriteToolConfigParamsTypeWrite ManagedAgentsWriteToolConfigParamsType = "write"
)

type AgentNewParams struct {
	// Model identifier. Accepts the
	// [model string](https://docs.qoder.com/cloud-agents/api/models/list) or a
	// `model_config` object for additional configuration control
	Model ManagedAgentsModelConfigParams `json:"model,omitzero" api:"required"`
	// Human-readable name for the agent.
	Name string `json:"name" api:"required"`
	// Description of what the agent does.
	Description param.Opt[string] `json:"description,omitzero"`
	// System prompt for the agent.
	System      param.Opt[string] `json:"system,omitzero"`
	WorkspaceID param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`
	// MCP servers this agent connects to. Maximum 20. Names must be unique within the
	// array. Every server must be referenced by an `mcp_toolset` in `tools`;
	// unreferenced servers are rejected. See the
	// [MCP connector guide](https://docs.qoder.com/cloud-agents/api/agents/schemas).
	MCPServers []ManagedAgentsURLMCPServerParams `json:"mcp_servers,omitzero"`
	// Arbitrary key-value metadata. Maximum 16 pairs, keys up to 64 chars, values up
	// to 512 chars.
	Metadata map[string]string `json:"metadata,omitzero"`
	// A coordinator topology: the session's primary thread orchestrates work by
	// spawning session threads, each running an agent drawn from the `agents` roster.
	Multiagent ManagedAgentsMultiagentParams `json:"multiagent,omitzero"`
	// Skills available to the agent.
	Skills []ManagedAgentsSkillParamsUnion `json:"skills,omitzero"`
	// Tool configurations available to the agent. Maximum of 128 tools across all
	// toolsets allowed.
	Tools []AgentNewParamsToolUnion `json:"tools,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

func (r AgentNewParams) MarshalJSON() (data []byte, err error) {
	type shadow AgentNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AgentNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type AgentNewParamsToolUnion struct {
	OfAgentToolset20260401 *ManagedAgentsAgentToolset20260401Params `json:",omitzero,inline"`
	OfMCPToolset           *ManagedAgentsMCPToolsetParams           `json:",omitzero,inline"`
	OfCustom               *ManagedAgentsCustomToolParams           `json:",omitzero,inline"`
	paramUnion
}

func (u AgentNewParamsToolUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAgentToolset20260401, u.OfMCPToolset, u.OfCustom)
}
func (u *AgentNewParamsToolUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *AgentNewParamsToolUnion) asAny() any {
	if !param.IsOmitted(u.OfAgentToolset20260401) {
		return u.OfAgentToolset20260401
	} else if !param.IsOmitted(u.OfMCPToolset) {
		return u.OfMCPToolset
	} else if !param.IsOmitted(u.OfCustom) {
		return u.OfCustom
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentNewParamsToolUnion) GetMCPServerName() *string {
	if vt := u.OfMCPToolset; vt != nil {
		return &vt.MCPServerName
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentNewParamsToolUnion) GetDescription() *string {
	if vt := u.OfCustom; vt != nil {
		return &vt.Description
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentNewParamsToolUnion) GetInputSchema() *ManagedAgentsCustomToolInputSchemaParam {
	if vt := u.OfCustom; vt != nil {
		return &vt.InputSchema
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentNewParamsToolUnion) GetName() *string {
	if vt := u.OfCustom; vt != nil {
		return &vt.Name
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentNewParamsToolUnion) GetType() *string {
	if vt := u.OfAgentToolset20260401; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfMCPToolset; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCustom; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u AgentNewParamsToolUnion) GetConfigs() (res agentNewParamsToolUnionConfigs) {
	if vt := u.OfAgentToolset20260401; vt != nil {
		res.any = &vt.Configs
	} else if vt := u.OfMCPToolset; vt != nil {
		res.any = &vt.Configs
	}
	return
}

// Can have the runtime types [_[]ManagedAgentsAgentToolConfigParamsUnion],
// [_[]ManagedAgentsMCPToolConfigParams]
type agentNewParamsToolUnionConfigs struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *[]qoder.ManagedAgentsAgentToolConfigParamsUnion:
//	case *[]qoder.ManagedAgentsMCPToolConfigParams:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u agentNewParamsToolUnionConfigs) AsAny() any { return u.any }

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u AgentNewParamsToolUnion) GetDefaultConfig() (res agentNewParamsToolUnionDefaultConfig) {
	if vt := u.OfAgentToolset20260401; vt != nil {
		res.any = &vt.DefaultConfig
	} else if vt := u.OfMCPToolset; vt != nil {
		res.any = &vt.DefaultConfig
	}
	return
}

// Can have the runtime types [*ManagedAgentsAgentToolsetDefaultConfigParams],
// [*ManagedAgentsMCPToolsetDefaultConfigParams]
type agentNewParamsToolUnionDefaultConfig struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *qoder.ManagedAgentsAgentToolsetDefaultConfigParams:
//	case *qoder.ManagedAgentsMCPToolsetDefaultConfigParams:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u agentNewParamsToolUnionDefaultConfig) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u agentNewParamsToolUnionDefaultConfig) GetEnabled() *bool {
	switch vt := u.any.(type) {
	case *ManagedAgentsAgentToolsetDefaultConfigParams:
		return paramutil.AddrIfPresent(vt.Enabled)
	case *ManagedAgentsMCPToolsetDefaultConfigParams:
		return paramutil.AddrIfPresent(vt.Enabled)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u agentNewParamsToolUnionDefaultConfig) GetPermissionPolicy() (res agentNewParamsToolUnionDefaultConfigPermissionPolicy) {
	switch vt := u.any.(type) {
	case *ManagedAgentsAgentToolsetDefaultConfigParams:
		res.any = vt.PermissionPolicy
	case *ManagedAgentsMCPToolsetDefaultConfigParams:
		res.any = vt.PermissionPolicy
	}
	return res
}

// Can have the runtime types [*ManagedAgentsAlwaysAllowPolicyParam],
// [*ManagedAgentsAlwaysAskPolicyParam]
type agentNewParamsToolUnionDefaultConfigPermissionPolicy struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *qoder.ManagedAgentsAlwaysAllowPolicyParam:
//	case *qoder.ManagedAgentsAlwaysAskPolicyParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u agentNewParamsToolUnionDefaultConfigPermissionPolicy) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u agentNewParamsToolUnionDefaultConfigPermissionPolicy) GetType() *string {
	switch vt := u.any.(type) {
	case *ManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion:
		return vt.GetType()
	case *ManagedAgentsMCPToolsetDefaultConfigParamsPermissionPolicyUnion:
		return vt.GetType()
	}
	return nil
}

func init() {
	apijson.RegisterUnion[AgentNewParamsToolUnion](
		"type",
		apijson.Discriminator[ManagedAgentsAgentToolset20260401Params]("agent_toolset_20260401"),
		apijson.Discriminator[ManagedAgentsMCPToolsetParams]("mcp_toolset"),
		apijson.Discriminator[ManagedAgentsCustomToolParams]("custom"),
	)
}

type AgentGetParams struct {
	// Agent version. Omit for the most recent version. Must be at least 1 if
	// specified.
	Version     param.Opt[int64]  `query:"version,omitzero" json:"-"`
	WorkspaceID param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AgentGetParams]'s query parameters as `url.Values`.
func (r AgentGetParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type AgentUpdateParams struct {
	// Description. Omit to preserve; send empty string or null to clear.
	Description param.Opt[string] `json:"description,omitzero"`
	// System prompt. Omit to preserve; send empty string or null to clear.
	System param.Opt[string] `json:"system,omitzero"`
	// Human-readable name. Must be non-empty. Omit to preserve. Cannot be cleared.
	Name param.Opt[string] `json:"name,omitzero"`
	// The agent's current version, used to prevent concurrent overwrites. Obtain this
	// value from a create or retrieve response. Must be at least 1 if specified. When
	// supplied, the request fails if it does not match the server's current version;
	// omit to apply the update unconditionally.
	Version     param.Opt[int64]  `json:"version,omitzero"`
	WorkspaceID param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`
	// MCP servers. Full replacement. Omit to preserve; send empty array or `null` to
	// clear. Names must be unique. Maximum 20. Every server must be referenced by an
	// `mcp_toolset` in the agent's resulting `tools`; unreferenced servers are
	// rejected. See the
	// [MCP connector guide](https://docs.qoder.com/cloud-agents/api/agents/schemas).
	MCPServers []ManagedAgentsURLMCPServerParams `json:"mcp_servers,omitzero"`
	// Metadata patch. Set a key to a string to upsert it, or to null to delete it.
	// Omit the field to preserve. The stored bag is limited to 16 keys (up to 64 chars
	// each) with values up to 512 chars.
	Metadata map[string]any `json:"metadata,omitzero"`
	// Skills. Full replacement. Omit to preserve; send empty array or null to clear.
	Skills []ManagedAgentsSkillParamsUnion `json:"skills,omitzero"`
	// Tool configurations available to the agent. Full replacement. Omit to preserve;
	// send empty array or null to clear. Maximum of 128 tools across all toolsets
	// allowed.
	Tools []AgentUpdateParamsToolUnion `json:"tools,omitzero"`
	// Model identifier. Accepts the
	// [model string](https://docs.qoder.com/cloud-agents/api/models/list) or a
	// `model_config` object for additional configuration control. Omit to preserve.
	// Cannot be cleared.
	Model ManagedAgentsModelConfigParams `json:"model,omitzero"`
	// A coordinator topology: the session's primary thread orchestrates work by
	// spawning session threads, each running an agent drawn from the `agents` roster.
	Multiagent ManagedAgentsMultiagentParams `json:"multiagent,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

func (r AgentUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow AgentUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *AgentUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type AgentUpdateParamsToolUnion struct {
	OfAgentToolset20260401 *ManagedAgentsAgentToolset20260401Params `json:",omitzero,inline"`
	OfMCPToolset           *ManagedAgentsMCPToolsetParams           `json:",omitzero,inline"`
	OfCustom               *ManagedAgentsCustomToolParams           `json:",omitzero,inline"`
	paramUnion
}

func (u AgentUpdateParamsToolUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAgentToolset20260401, u.OfMCPToolset, u.OfCustom)
}
func (u *AgentUpdateParamsToolUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *AgentUpdateParamsToolUnion) asAny() any {
	if !param.IsOmitted(u.OfAgentToolset20260401) {
		return u.OfAgentToolset20260401
	} else if !param.IsOmitted(u.OfMCPToolset) {
		return u.OfMCPToolset
	} else if !param.IsOmitted(u.OfCustom) {
		return u.OfCustom
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentUpdateParamsToolUnion) GetMCPServerName() *string {
	if vt := u.OfMCPToolset; vt != nil {
		return &vt.MCPServerName
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentUpdateParamsToolUnion) GetDescription() *string {
	if vt := u.OfCustom; vt != nil {
		return &vt.Description
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentUpdateParamsToolUnion) GetInputSchema() *ManagedAgentsCustomToolInputSchemaParam {
	if vt := u.OfCustom; vt != nil {
		return &vt.InputSchema
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentUpdateParamsToolUnion) GetName() *string {
	if vt := u.OfCustom; vt != nil {
		return &vt.Name
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u AgentUpdateParamsToolUnion) GetType() *string {
	if vt := u.OfAgentToolset20260401; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfMCPToolset; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCustom; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u AgentUpdateParamsToolUnion) GetConfigs() (res agentUpdateParamsToolUnionConfigs) {
	if vt := u.OfAgentToolset20260401; vt != nil {
		res.any = &vt.Configs
	} else if vt := u.OfMCPToolset; vt != nil {
		res.any = &vt.Configs
	}
	return
}

// Can have the runtime types [_[]ManagedAgentsAgentToolConfigParamsUnion],
// [_[]ManagedAgentsMCPToolConfigParams]
type agentUpdateParamsToolUnionConfigs struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *[]qoder.ManagedAgentsAgentToolConfigParamsUnion:
//	case *[]qoder.ManagedAgentsMCPToolConfigParams:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u agentUpdateParamsToolUnionConfigs) AsAny() any { return u.any }

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u AgentUpdateParamsToolUnion) GetDefaultConfig() (res agentUpdateParamsToolUnionDefaultConfig) {
	if vt := u.OfAgentToolset20260401; vt != nil {
		res.any = &vt.DefaultConfig
	} else if vt := u.OfMCPToolset; vt != nil {
		res.any = &vt.DefaultConfig
	}
	return
}

// Can have the runtime types [*ManagedAgentsAgentToolsetDefaultConfigParams],
// [*ManagedAgentsMCPToolsetDefaultConfigParams]
type agentUpdateParamsToolUnionDefaultConfig struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *qoder.ManagedAgentsAgentToolsetDefaultConfigParams:
//	case *qoder.ManagedAgentsMCPToolsetDefaultConfigParams:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u agentUpdateParamsToolUnionDefaultConfig) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u agentUpdateParamsToolUnionDefaultConfig) GetEnabled() *bool {
	switch vt := u.any.(type) {
	case *ManagedAgentsAgentToolsetDefaultConfigParams:
		return paramutil.AddrIfPresent(vt.Enabled)
	case *ManagedAgentsMCPToolsetDefaultConfigParams:
		return paramutil.AddrIfPresent(vt.Enabled)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u agentUpdateParamsToolUnionDefaultConfig) GetPermissionPolicy() (res agentUpdateParamsToolUnionDefaultConfigPermissionPolicy) {
	switch vt := u.any.(type) {
	case *ManagedAgentsAgentToolsetDefaultConfigParams:
		res.any = vt.PermissionPolicy
	case *ManagedAgentsMCPToolsetDefaultConfigParams:
		res.any = vt.PermissionPolicy
	}
	return res
}

// Can have the runtime types [*ManagedAgentsAlwaysAllowPolicyParam],
// [*ManagedAgentsAlwaysAskPolicyParam]
type agentUpdateParamsToolUnionDefaultConfigPermissionPolicy struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *qoder.ManagedAgentsAlwaysAllowPolicyParam:
//	case *qoder.ManagedAgentsAlwaysAskPolicyParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u agentUpdateParamsToolUnionDefaultConfigPermissionPolicy) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u agentUpdateParamsToolUnionDefaultConfigPermissionPolicy) GetType() *string {
	switch vt := u.any.(type) {
	case *ManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion:
		return vt.GetType()
	case *ManagedAgentsMCPToolsetDefaultConfigParamsPermissionPolicyUnion:
		return vt.GetType()
	}
	return nil
}

func init() {
	apijson.RegisterUnion[AgentUpdateParamsToolUnion](
		"type",
		apijson.Discriminator[ManagedAgentsAgentToolset20260401Params]("agent_toolset_20260401"),
		apijson.Discriminator[ManagedAgentsMCPToolsetParams]("mcp_toolset"),
		apijson.Discriminator[ManagedAgentsCustomToolParams]("custom"),
	)
}

type AgentListParams struct {
	// Return agents created at or after this time (inclusive).
	CreatedAtGte param.Opt[time.Time] `query:"created_at[gte],omitzero" format:"date-time" json:"-"`
	// Return agents created at or before this time (inclusive).
	CreatedAtLte param.Opt[time.Time] `query:"created_at[lte],omitzero" format:"date-time" json:"-"`
	// Include archived agents in results. Defaults to false.
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// Maximum results per page. Default 20, maximum 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque pagination cursor from a previous response.
	Page        param.Opt[string] `query:"page,omitzero" json:"-"`
	WorkspaceID param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AgentListParams]'s query parameters as `url.Values`.
func (r AgentListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type AgentArchiveParams struct {
	WorkspaceID param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}
