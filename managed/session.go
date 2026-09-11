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
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/paramutil"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
)

// SessionService contains methods and other services that help with
// interacting with the Qoder Cloud Agents API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSessionService] method instead.
type SessionService struct {
	Options   []option.RequestOption
	Events    SessionEventService
	Resources SessionResourceService
	Threads   SessionThreadService
}

// NewSessionService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewSessionService(opts ...option.RequestOption) (r SessionService) {
	r = SessionService{}
	r.Options = opts
	r.Events = NewSessionEventService(opts...)
	r.Resources = NewSessionResourceService(opts...)
	r.Threads = NewSessionThreadService(opts...)
	return
}

// Create Session
func (r *SessionService) New(ctx context.Context, params SessionNewParams, opts ...option.RequestOption) (res *ManagedAgentsSession, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	path := "sessions"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Get Session
func (r *SessionService) Get(ctx context.Context, sessionID string, query SessionGetParams, opts ...option.RequestOption) (res *ManagedAgentsSession, err error) {
	for _, v := range query.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("sessions/%s", url.PathEscape(sessionID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update Session
func (r *SessionService) Update(ctx context.Context, sessionID string, params SessionUpdateParams, opts ...option.RequestOption) (res *ManagedAgentsSession, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("sessions/%s", url.PathEscape(sessionID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List Sessions
func (r *SessionService) List(ctx context.Context, params SessionListParams, opts ...option.RequestOption) (res *pagination.BidirectionalPageCursor[ManagedAgentsSession], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "sessions"
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

// List Sessions
func (r *SessionService) ListAutoPaging(ctx context.Context, params SessionListParams, opts ...option.RequestOption) *pagination.BidirectionalPageCursorAutoPager[ManagedAgentsSession] {
	return pagination.NewBidirectionalPageCursorAutoPager(r.List(ctx, params, opts...))
}

// Delete Session
func (r *SessionService) Delete(ctx context.Context, sessionID string, body SessionDeleteParams, opts ...option.RequestOption) (res *ManagedAgentsDeletedSession, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("sessions/%s", url.PathEscape(sessionID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Archive Session
func (r *SessionService) Archive(ctx context.Context, sessionID string, body SessionArchiveParams, opts ...option.RequestOption) (res *ManagedAgentsSession, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("sessions/%s/archive", url.PathEscape(sessionID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Platform advisor roster entry: a model the session's primary thread may consult
// mid-turn. At most one per roster; the entry occupies the roster name
// `qoder.advisor`.
//
// The properties Model, Type are required.
type ManagedAgentsAdvisorParams struct {
	// A model id. The model must be permitted as an advisor for this agent's model
	// — see the sessions/threads/advisor spec.
	Model string `json:"model" api:"required"`
	// Any of "advisor".
	Type ManagedAgentsAdvisorParamsType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsAdvisorParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsAdvisorParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsAdvisorParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsAdvisorParamsType string

const (
	ManagedAgentsAdvisorParamsTypeAdvisor ManagedAgentsAdvisorParamsType = "advisor"
)

type ManagedAgentsAgentMessagePreview struct {
	// The id the buffered agent.message will carry if it is emitted. Matches the
	// event_id on this preview's event_delta events.
	ID string `json:"id" api:"required"`
	// Any of "agent.message".
	Type ManagedAgentsAgentMessagePreviewType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsAgentMessagePreview) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsAgentMessagePreview) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsAgentMessagePreviewType string

const (
	ManagedAgentsAgentMessagePreviewTypeAgentMessage ManagedAgentsAgentMessagePreviewType = "agent.message"
)

// Specification for an Agent. Provide a specific `version` or use the short-form
// `agent="agent_id"` for the most recent version
//
// The properties ID, Type are required.
type ManagedAgentsAgentParams struct {
	// The `agent` ID.
	ID string `json:"id" api:"required"`
	// Any of "agent".
	Type ManagedAgentsAgentParamsType `json:"type,omitzero" api:"required"`
	// The specific `agent` version to use. Omit to use the latest version. Must be at
	// least 1 if specified.
	Version param.Opt[int64] `json:"version,omitzero"`
	paramObj
}

func (r ManagedAgentsAgentParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsAgentParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsAgentParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsAgentParamsType string

const (
	ManagedAgentsAgentParamsTypeAgent ManagedAgentsAgentParamsType = "agent"
)

type ManagedAgentsAgentThinkingPreview struct {
	// The id the buffered agent.thinking will carry if it is emitted. Start-only — no
	// event_delta events follow.
	ID string `json:"id" api:"required"`
	// Any of "agent.thinking".
	Type ManagedAgentsAgentThinkingPreviewType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsAgentThinkingPreview) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsAgentThinkingPreview) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsAgentThinkingPreviewType string

const (
	ManagedAgentsAgentThinkingPreviewTypeAgentThinking ManagedAgentsAgentThinkingPreviewType = "agent.thinking"
)

// Reference to an `agent` plus optional configuration overrides. Each provided
// field replaces the agent's value for the caller's use; the agent resource is
// unchanged.
//
// The properties ID, Type are required.
type ManagedAgentsAgentWithOverridesParams struct {
	// The `agent` ID.
	ID string `json:"id" api:"required"`
	// Any of "agent_with_overrides".
	Type ManagedAgentsAgentWithOverridesParamsType `json:"type,omitzero" api:"required"`
	// Replacement system prompt. Up to 100,000 characters. Set to null to clear the
	// agent's system prompt; omit to preserve it.
	System param.Opt[string] `json:"system,omitzero"`
	// The specific `agent` version to use. Omit to use the latest version.
	Version param.Opt[int64] `json:"version,omitzero"`
	// Replacement MCP server list. Full replacement: the provided array becomes the
	// MCP servers. Send an empty array to clear; omit to preserve the agent's servers.
	MCPServers []ManagedAgentsURLMCPServerParams `json:"mcp_servers,omitzero"`
	// Replacement model. Accepts the model string or a `model_config` object. Omit to
	// use the agent's model.
	Model ManagedAgentsModelConfigParams `json:"model,omitzero"`
	// Replacement skill list. Full replacement: the provided array becomes the skills.
	// Send an empty array to clear; omit to preserve the agent's skills.
	Skills []ManagedAgentsSkillParamsUnion `json:"skills,omitzero"`
	// Replacement tool list. Full replacement: the provided array becomes the tool
	// configuration. Send an empty array to clear; omit to preserve the agent's tools.
	Tools []ManagedAgentsAgentWithOverridesParamsToolUnion `json:"tools,omitzero"`
	paramObj
}

func (r ManagedAgentsAgentWithOverridesParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsAgentWithOverridesParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsAgentWithOverridesParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsAgentWithOverridesParamsType string

const (
	ManagedAgentsAgentWithOverridesParamsTypeAgentWithOverrides ManagedAgentsAgentWithOverridesParamsType = "agent_with_overrides"
)

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsAgentWithOverridesParamsToolUnion struct {
	OfAgentToolset20260401 *ManagedAgentsAgentToolset20260401Params `json:",omitzero,inline"`
	OfMCPToolset           *ManagedAgentsMCPToolsetParams           `json:",omitzero,inline"`
	OfCustom               *ManagedAgentsCustomToolParams           `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsAgentWithOverridesParamsToolUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAgentToolset20260401, u.OfMCPToolset, u.OfCustom)
}
func (u *ManagedAgentsAgentWithOverridesParamsToolUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsAgentWithOverridesParamsToolUnion) asAny() any {
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
func (u ManagedAgentsAgentWithOverridesParamsToolUnion) GetMCPServerName() *string {
	if vt := u.OfMCPToolset; vt != nil {
		return &vt.MCPServerName
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsAgentWithOverridesParamsToolUnion) GetDescription() *string {
	if vt := u.OfCustom; vt != nil {
		return &vt.Description
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsAgentWithOverridesParamsToolUnion) GetInputSchema() *ManagedAgentsCustomToolInputSchemaParam {
	if vt := u.OfCustom; vt != nil {
		return &vt.InputSchema
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsAgentWithOverridesParamsToolUnion) GetName() *string {
	if vt := u.OfCustom; vt != nil {
		return &vt.Name
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsAgentWithOverridesParamsToolUnion) GetType() *string {
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
func (u ManagedAgentsAgentWithOverridesParamsToolUnion) GetConfigs() (res managedAgentsAgentWithOverridesParamsToolUnionConfigs) {
	if vt := u.OfAgentToolset20260401; vt != nil {
		res.any = &vt.Configs
	} else if vt := u.OfMCPToolset; vt != nil {
		res.any = &vt.Configs
	}
	return
}

// Can have the runtime types [_[]ManagedAgentsAgentToolConfigParamsUnion],
// [_[]ManagedAgentsMCPToolConfigParams]
type managedAgentsAgentWithOverridesParamsToolUnionConfigs struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *[]qoder.ManagedAgentsAgentToolConfigParamsUnion:
//	case *[]qoder.ManagedAgentsMCPToolConfigParams:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u managedAgentsAgentWithOverridesParamsToolUnionConfigs) AsAny() any { return u.any }

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u ManagedAgentsAgentWithOverridesParamsToolUnion) GetDefaultConfig() (res managedAgentsAgentWithOverridesParamsToolUnionDefaultConfig) {
	if vt := u.OfAgentToolset20260401; vt != nil {
		res.any = &vt.DefaultConfig
	} else if vt := u.OfMCPToolset; vt != nil {
		res.any = &vt.DefaultConfig
	}
	return
}

// Can have the runtime types [*ManagedAgentsAgentToolsetDefaultConfigParams],
// [*ManagedAgentsMCPToolsetDefaultConfigParams]
type managedAgentsAgentWithOverridesParamsToolUnionDefaultConfig struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *qoder.ManagedAgentsAgentToolsetDefaultConfigParams:
//	case *qoder.ManagedAgentsMCPToolsetDefaultConfigParams:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u managedAgentsAgentWithOverridesParamsToolUnionDefaultConfig) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u managedAgentsAgentWithOverridesParamsToolUnionDefaultConfig) GetEnabled() *bool {
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
func (u managedAgentsAgentWithOverridesParamsToolUnionDefaultConfig) GetPermissionPolicy() (res managedAgentsAgentWithOverridesParamsToolUnionDefaultConfigPermissionPolicy) {
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
type managedAgentsAgentWithOverridesParamsToolUnionDefaultConfigPermissionPolicy struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *qoder.ManagedAgentsAlwaysAllowPolicyParam:
//	case *qoder.ManagedAgentsAlwaysAskPolicyParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u managedAgentsAgentWithOverridesParamsToolUnionDefaultConfigPermissionPolicy) AsAny() any {
	return u.any
}

// Returns a pointer to the underlying variant's property, if present.
func (u managedAgentsAgentWithOverridesParamsToolUnionDefaultConfigPermissionPolicy) GetType() *string {
	switch vt := u.any.(type) {
	case *ManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion:
		return vt.GetType()
	case *ManagedAgentsMCPToolsetDefaultConfigParamsPermissionPolicyUnion:
		return vt.GetType()
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ManagedAgentsAgentWithOverridesParamsToolUnion](
		"type",
		apijson.Discriminator[ManagedAgentsAgentToolset20260401Params]("agent_toolset_20260401"),
		apijson.Discriminator[ManagedAgentsMCPToolsetParams]("mcp_toolset"),
		apijson.Discriminator[ManagedAgentsCustomToolParams]("custom"),
	)
}

type ManagedAgentsBranchCheckout struct {
	// Branch name to check out.
	Name string `json:"name" api:"required"`
	// Any of "branch".
	Type ManagedAgentsBranchCheckoutType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Name        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsBranchCheckout) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsBranchCheckout) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ManagedAgentsBranchCheckout to a
// ManagedAgentsBranchCheckoutParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ManagedAgentsBranchCheckoutParam.Overrides()
func (r ManagedAgentsBranchCheckout) ToParam() ManagedAgentsBranchCheckoutParam {
	return param.Override[ManagedAgentsBranchCheckoutParam](json.RawMessage(r.RawJSON()))
}

type ManagedAgentsBranchCheckoutType string

const (
	ManagedAgentsBranchCheckoutTypeBranch ManagedAgentsBranchCheckoutType = "branch"
)

// The properties Name, Type are required.
type ManagedAgentsBranchCheckoutParam struct {
	// Branch name to check out.
	Name string `json:"name" api:"required"`
	// Any of "branch".
	Type ManagedAgentsBranchCheckoutType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsBranchCheckoutParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsBranchCheckoutParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsBranchCheckoutParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A hard spend ceiling. The session stops issuing new model requests once the
// tracked list cost reaches `max_list_cost`.
type ManagedAgentsBudgetLimit struct {
	// A monetary amount in a specific currency.
	MaxListCost MonetaryAmount `json:"max_list_cost" api:"required"`
	// Any of "limit".
	Type ManagedAgentsBudgetLimitType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		MaxListCost respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsBudgetLimit) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsBudgetLimit) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ManagedAgentsBudgetLimit to a
// ManagedAgentsBudgetLimitParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ManagedAgentsBudgetLimitParam.Overrides()
func (r ManagedAgentsBudgetLimit) ToParam() ManagedAgentsBudgetLimitParam {
	return param.Override[ManagedAgentsBudgetLimitParam](json.RawMessage(r.RawJSON()))
}

type ManagedAgentsBudgetLimitType string

const (
	ManagedAgentsBudgetLimitTypeLimit ManagedAgentsBudgetLimitType = "limit"
)

// A hard spend ceiling. The session stops issuing new model requests once the
// tracked list cost reaches `max_list_cost`.
//
// The properties MaxListCost, Type are required.
type ManagedAgentsBudgetLimitParam struct {
	// A monetary amount in a specific currency.
	MaxListCost MonetaryAmountParam `json:"max_list_cost,omitzero" api:"required"`
	// Any of "limit".
	Type ManagedAgentsBudgetLimitType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsBudgetLimitParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsBudgetLimitParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsBudgetLimitParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Prompt-cache creation token usage broken down by cache lifetime.
type ManagedAgentsCacheCreationUsage struct {
	// Tokens used to create 1-hour ephemeral cache entries.
	Ephemeral1hInputTokens int64 `json:"ephemeral_1h_input_tokens"`
	// Tokens used to create 5-minute ephemeral cache entries.
	Ephemeral5mInputTokens int64 `json:"ephemeral_5m_input_tokens"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Ephemeral1hInputTokens respjson.Field
		Ephemeral5mInputTokens respjson.Field
		ExtraFields            map[string]respjson.Field
		raw                    string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsCacheCreationUsage) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsCacheCreationUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsCommitCheckout struct {
	// Full commit SHA to check out.
	Sha string `json:"sha" api:"required"`
	// Any of "commit".
	Type ManagedAgentsCommitCheckoutType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Sha         respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsCommitCheckout) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsCommitCheckout) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ManagedAgentsCommitCheckout to a
// ManagedAgentsCommitCheckoutParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ManagedAgentsCommitCheckoutParam.Overrides()
func (r ManagedAgentsCommitCheckout) ToParam() ManagedAgentsCommitCheckoutParam {
	return param.Override[ManagedAgentsCommitCheckoutParam](json.RawMessage(r.RawJSON()))
}

type ManagedAgentsCommitCheckoutType string

const (
	ManagedAgentsCommitCheckoutTypeCommit ManagedAgentsCommitCheckoutType = "commit"
)

// The properties Sha, Type are required.
type ManagedAgentsCommitCheckoutParam struct {
	// Full commit SHA to check out.
	Sha string `json:"sha" api:"required"`
	// Any of "commit".
	Type ManagedAgentsCommitCheckoutType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsCommitCheckoutParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsCommitCheckoutParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsCommitCheckoutParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Confirmation that a `session` has been permanently deleted.
type ManagedAgentsDeletedSession struct {
	ID string `json:"id" api:"required"`
	// Any of "session_deleted".
	Type ManagedAgentsDeletedSessionType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsDeletedSession) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsDeletedSession) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsDeletedSessionType string

const (
	ManagedAgentsDeletedSessionTypeSessionDeleted ManagedAgentsDeletedSessionType = "session_deleted"
)

type ManagedAgentsDeltaContent struct {
	// Regular text content.
	Content ManagedAgentsTextBlock `json:"content" api:"required"`
	// Any of "content_delta".
	Type ManagedAgentsDeltaContentType `json:"type" api:"required"`
	// Which entry in the previewed event's content array this fragment lands in.
	// Insert content as that entry when the index is new; append to the existing entry
	// otherwise.
	Index int64 `json:"index"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Content     respjson.Field
		Type        respjson.Field
		Index       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsDeltaContent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsDeltaContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsDeltaContentType string

const (
	ManagedAgentsDeltaContentTypeContentDelta ManagedAgentsDeltaContentType = "content_delta"
)

// An incremental update to an event that is still being streamed. Deltas are
// best-effort and may stop early; when the buffered event with id == event_id is
// produced it carries the complete content. A model request that ends early (an
// error or interrupt) produces no buffered event — its terminal
// span.model_request_end closes the preview. Only sent on stream connections that
// opt in via event_deltas; never appears in event history.
type ManagedAgentsDeltaEvent struct {
	// One fragment of the previewed event. The delta type is named for the previewed
	// event's field it streams into: agent.message events stream content_delta
	// fragments, each a partial element of the content array.
	Delta ManagedAgentsDeltaContent `json:"delta" api:"required"`
	// The id of the event being previewed. Matches event.id on the corresponding
	// event_start and the buffered event that reconciles the preview.
	EventID string `json:"event_id" api:"required"`
	// Any of "event_delta".
	Type ManagedAgentsDeltaEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Delta       respjson.Field
		EventID     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsDeltaEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsDeltaEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsDeltaEventType string

const (
	ManagedAgentsDeltaEventTypeEventDelta ManagedAgentsDeltaEventType = "event_delta"
)

// EventDeltaType enum
type ManagedAgentsDeltaType string

const (
	ManagedAgentsDeltaTypeAgentMessage  ManagedAgentsDeltaType = "agent.message"
	ManagedAgentsDeltaTypeAgentThinking ManagedAgentsDeltaType = "agent.thinking"
)

// Mount a file uploaded via the Files API into the session.
//
// The properties FileID, Type are required.
type ManagedAgentsFileResourceParams struct {
	// ID of a previously uploaded file.
	FileID string `json:"file_id" api:"required"`
	// Any of "file".
	Type ManagedAgentsFileResourceParamsType `json:"type,omitzero" api:"required"`
	// Mount path in the container. Defaults to `/mnt/session/uploads/<file_id>`.
	MountPath param.Opt[string] `json:"mount_path,omitzero"`
	paramObj
}

func (r ManagedAgentsFileResourceParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsFileResourceParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsFileResourceParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsFileResourceParamsType string

const (
	ManagedAgentsFileResourceParamsTypeFile ManagedAgentsFileResourceParamsType = "file"
)

// Mount a GitHub repository into the session's container.
//
// The properties AuthorizationToken, Type, URL are required.
type ManagedAgentsGitHubRepositoryResourceParams struct {
	// GitHub authorization token used to clone the repository.
	AuthorizationToken string `json:"authorization_token" api:"required"`
	// Any of "github_repository".
	Type ManagedAgentsGitHubRepositoryResourceParamsType `json:"type,omitzero" api:"required"`
	// Github URL of the repository
	URL string `json:"url" api:"required"`
	// Mount path in the container. Defaults to `/workspace/<repo-name>`.
	MountPath param.Opt[string] `json:"mount_path,omitzero"`
	// Branch or commit to check out. Defaults to the repository's default branch.
	Checkout ManagedAgentsGitHubRepositoryResourceParamsCheckoutUnion `json:"checkout,omitzero"`
	paramObj
}

func (r ManagedAgentsGitHubRepositoryResourceParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsGitHubRepositoryResourceParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsGitHubRepositoryResourceParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsGitHubRepositoryResourceParamsType string

const (
	ManagedAgentsGitHubRepositoryResourceParamsTypeGitHubRepository ManagedAgentsGitHubRepositoryResourceParamsType = "github_repository"
)

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsGitHubRepositoryResourceParamsCheckoutUnion struct {
	OfBranch *ManagedAgentsBranchCheckoutParam `json:",omitzero,inline"`
	OfCommit *ManagedAgentsCommitCheckoutParam `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsGitHubRepositoryResourceParamsCheckoutUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfBranch, u.OfCommit)
}
func (u *ManagedAgentsGitHubRepositoryResourceParamsCheckoutUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsGitHubRepositoryResourceParamsCheckoutUnion) asAny() any {
	if !param.IsOmitted(u.OfBranch) {
		return u.OfBranch
	} else if !param.IsOmitted(u.OfCommit) {
		return u.OfCommit
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsGitHubRepositoryResourceParamsCheckoutUnion) GetName() *string {
	if vt := u.OfBranch; vt != nil {
		return &vt.Name
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsGitHubRepositoryResourceParamsCheckoutUnion) GetSha() *string {
	if vt := u.OfCommit; vt != nil {
		return &vt.Sha
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsGitHubRepositoryResourceParamsCheckoutUnion) GetType() *string {
	if vt := u.OfBranch; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfCommit; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ManagedAgentsGitHubRepositoryResourceParamsCheckoutUnion](
		"type",
		apijson.Discriminator[ManagedAgentsBranchCheckoutParam]("branch"),
		apijson.Discriminator[ManagedAgentsCommitCheckoutParam]("commit"),
	)
}

// Parameters for attaching a memory store to an agent session.
//
// The properties MemoryStoreID, Type are required.
type ManagedAgentsMemoryStoreResourceParam struct {
	// The memory store ID (memstore\_...). Must belong to the caller's organization
	// and workspace.
	MemoryStoreID string `json:"memory_store_id" api:"required"`
	// Any of "memory_store".
	Type ManagedAgentsMemoryStoreResourceParamType `json:"type,omitzero" api:"required"`
	// Per-attachment guidance for the agent on how to use this store. Rendered into
	// the memory section of the system prompt. Max 4096 chars.
	Instructions param.Opt[string] `json:"instructions,omitzero"`
	// Access mode for an attached memory store.
	//
	// Any of "read_write", "read_only".
	Access ManagedAgentsMemoryStoreResourceParamAccess `json:"access,omitzero"`
	paramObj
}

func (r ManagedAgentsMemoryStoreResourceParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsMemoryStoreResourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsMemoryStoreResourceParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsMemoryStoreResourceParamType string

const (
	ManagedAgentsMemoryStoreResourceParamTypeMemoryStore ManagedAgentsMemoryStoreResourceParamType = "memory_store"
)

// Access mode for an attached memory store.
type ManagedAgentsMemoryStoreResourceParamAccess string

const (
	ManagedAgentsMemoryStoreResourceParamAccessReadWrite ManagedAgentsMemoryStoreResourceParamAccess = "read_write"
	ManagedAgentsMemoryStoreResourceParamAccessReadOnly  ManagedAgentsMemoryStoreResourceParamAccess = "read_only"
)

// Resolved coordinator topology with a concrete agent roster.
type ManagedAgentsMultiagent struct {
	// Agents the coordinator may spawn as session threads, each resolved to a specific
	// version.
	Agents []ManagedAgentsMultiagentAgentUnion `json:"agents" api:"required"`
	// Any of "coordinator".
	Type ManagedAgentsMultiagentType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Agents      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsMultiagent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsMultiagent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsMultiagentAgentUnion contains all possible properties and
// values from [ManagedAgentsAgentReference], [ManagedAgentsAdvisor].
//
// Use the [ManagedAgentsMultiagentAgentUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsMultiagentAgentUnion struct {
	// This field is from variant [ManagedAgentsAgentReference].
	ID string `json:"id"`
	// Any of "agent", "advisor".
	Type string `json:"type"`
	// This field is from variant [ManagedAgentsAgentReference].
	Version int64 `json:"version"`
	// This field is from variant [ManagedAgentsAdvisor].
	Model string `json:"model"`
	JSON  struct {
		ID      respjson.Field
		Type    respjson.Field
		Version respjson.Field
		Model   respjson.Field
		raw     string
	} `json:"-"`
}

// anyManagedAgentsMultiagentAgent is implemented by each variant of
// [ManagedAgentsMultiagentAgentUnion] to add type safety for the return type
// of [ManagedAgentsMultiagentAgentUnion.AsAny]
type anyManagedAgentsMultiagentAgent interface {
	implManagedAgentsMultiagentAgentUnion()
}

func (ManagedAgentsAgentReference) implManagedAgentsMultiagentAgentUnion() {}
func (ManagedAgentsAdvisor) implManagedAgentsMultiagentAgentUnion()        {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsMultiagentAgentUnion.AsAny().(type) {
//	case qoder.ManagedAgentsAgentReference:
//	case qoder.ManagedAgentsAdvisor:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsMultiagentAgentUnion) AsAny() anyManagedAgentsMultiagentAgent {
	switch u.Type {
	case "agent":
		return u.AsAgent()
	case "advisor":
		return u.AsAdvisor()
	}
	return nil
}

func (u ManagedAgentsMultiagentAgentUnion) AsAgent() (v ManagedAgentsAgentReference) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsMultiagentAgentUnion) AsAdvisor() (v ManagedAgentsAdvisor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsMultiagentAgentUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsMultiagentAgentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsMultiagentType string

const (
	ManagedAgentsMultiagentTypeCoordinator ManagedAgentsMultiagentType = "coordinator"
)

// A coordinator topology: the session's primary thread orchestrates work by
// spawning session threads, each running an agent drawn from the `agents` roster.
//
// The properties Agents, Type are required.
type ManagedAgentsMultiagentParams struct {
	// Agents the coordinator may spawn as session threads. 1–20 entries. Each entry is
	// an agent ID string, a versioned `{"type":"agent","id","version"}` reference, or
	// `{"type":"self"}` to allow recursive self-invocation. Entries must reference
	// distinct agents (after resolving `self` and string forms); at most one `self`.
	// Referenced agents must exist, must not be archived, and must not themselves have
	// `multiagent` set (depth limit 1).
	Agents []ManagedAgentsMultiagentRosterEntryParamsUnion `json:"agents,omitzero" api:"required"`
	// Any of "coordinator".
	Type ManagedAgentsMultiagentParamsType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsMultiagentParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsMultiagentParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsMultiagentParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsMultiagentParamsType string

const (
	ManagedAgentsMultiagentParamsTypeCoordinator ManagedAgentsMultiagentParamsType = "coordinator"
)

func ManagedAgentsMultiagentRosterEntryParamsOfManagedAgentsAgents(id string, type_ ManagedAgentsAgentParamsType) ManagedAgentsMultiagentRosterEntryParamsUnion {
	var variant ManagedAgentsAgentParams
	variant.ID = id
	variant.Type = type_
	return ManagedAgentsMultiagentRosterEntryParamsUnion{OfManagedAgentsAgents: &variant}
}

func ManagedAgentsMultiagentRosterEntryParamsOfManagedAgentsMultiagentSelfs(type_ ManagedAgentsMultiagentSelfParamsType) ManagedAgentsMultiagentRosterEntryParamsUnion {
	var variant ManagedAgentsMultiagentSelfParams
	variant.Type = type_
	return ManagedAgentsMultiagentRosterEntryParamsUnion{OfManagedAgentsMultiagentSelfs: &variant}
}

func ManagedAgentsMultiagentRosterEntryParamsOfManagedAgentsAdvisors(model string, type_ ManagedAgentsAdvisorParamsType) ManagedAgentsMultiagentRosterEntryParamsUnion {
	var variant ManagedAgentsAdvisorParams
	variant.Model = model
	variant.Type = type_
	return ManagedAgentsMultiagentRosterEntryParamsUnion{OfManagedAgentsAdvisors: &variant}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsMultiagentRosterEntryParamsUnion struct {
	OfString                       param.Opt[string]                  `json:",omitzero,inline"`
	OfManagedAgentsAgents          *ManagedAgentsAgentParams          `json:",omitzero,inline"`
	OfManagedAgentsMultiagentSelfs *ManagedAgentsMultiagentSelfParams `json:",omitzero,inline"`
	OfManagedAgentsAdvisors        *ManagedAgentsAdvisorParams        `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsMultiagentRosterEntryParamsUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfManagedAgentsAgents, u.OfManagedAgentsMultiagentSelfs, u.OfManagedAgentsAdvisors)
}
func (u *ManagedAgentsMultiagentRosterEntryParamsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsMultiagentRosterEntryParamsUnion) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfManagedAgentsAgents) {
		return u.OfManagedAgentsAgents
	} else if !param.IsOmitted(u.OfManagedAgentsMultiagentSelfs) {
		return u.OfManagedAgentsMultiagentSelfs
	} else if !param.IsOmitted(u.OfManagedAgentsAdvisors) {
		return u.OfManagedAgentsAdvisors
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsMultiagentRosterEntryParamsUnion) GetID() *string {
	if vt := u.OfManagedAgentsAgents; vt != nil {
		return &vt.ID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsMultiagentRosterEntryParamsUnion) GetVersion() *int64 {
	if vt := u.OfManagedAgentsAgents; vt != nil && vt.Version.Valid() {
		return &vt.Version.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsMultiagentRosterEntryParamsUnion) GetModel() *string {
	if vt := u.OfManagedAgentsAdvisors; vt != nil {
		return &vt.Model
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsMultiagentRosterEntryParamsUnion) GetType() *string {
	if vt := u.OfManagedAgentsAgents; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfManagedAgentsMultiagentSelfs; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfManagedAgentsAdvisors; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Evaluation state for a single outcome defined via a `define_outcome` event.
type ManagedAgentsOutcomeEvaluationResource struct {
	// A timestamp in RFC 3339 format
	CompletedAt time.Time `json:"completed_at" api:"required" format:"date-time"`
	// What the agent should produce.
	Description string `json:"description" api:"required"`
	// Grader's verdict text from the most recent evaluation. For `satisfied`, explains
	// why criteria are met; for `needs_revision` (intermediate), what's missing; for
	// `failed`, why unrecoverable.
	Explanation string `json:"explanation" api:"required"`
	// 0-indexed revision cycle the outcome is currently on.
	Iteration int64 `json:"iteration" api:"required"`
	// Server-generated outc\_ ID for this outcome.
	OutcomeID string `json:"outcome_id" api:"required"`
	// Current evaluation state. `pending` before the agent begins work; `running`
	// while producing or revising; `evaluating` while the grader scores;
	// `satisfied`/`max_iterations_reached`/`failed`/`interrupted` are terminal.
	Result string `json:"result" api:"required"`
	// Any of "outcome_evaluation".
	Type ManagedAgentsOutcomeEvaluationResourceType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CompletedAt respjson.Field
		Description respjson.Field
		Explanation respjson.Field
		Iteration   respjson.Field
		OutcomeID   respjson.Field
		Result      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsOutcomeEvaluationResource) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsOutcomeEvaluationResource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsOutcomeEvaluationResourceType string

const (
	ManagedAgentsOutcomeEvaluationResourceTypeOutcomeEvaluation ManagedAgentsOutcomeEvaluationResourceType = "outcome_evaluation"
)

// Cumulative count of server-executed tool invocations, broken down by tool.
type ManagedAgentsServerToolUsage struct {
	// Number of server-executed web fetch requests.
	WebFetchRequests int64 `json:"web_fetch_requests"`
	// Number of server-executed web search requests.
	WebSearchRequests int64 `json:"web_search_requests"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		WebFetchRequests  respjson.Field
		WebSearchRequests respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsServerToolUsage) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsServerToolUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A Managed Agents `session`.
type ManagedAgentsSession struct {
	EnvironmentVariables map[string]string `json:"environment_variables"`

	ID string `json:"id" api:"required"`
	// Resolved `agent` definition for a `session`. Snapshot of the `agent` at
	// `session` creation time.
	Agent ManagedAgentsSessionAgent `json:"agent" api:"required"`
	// A timestamp in RFC 3339 format
	ArchivedAt time.Time `json:"archived_at" api:"required" format:"date-time"`
	// A hard spend ceiling. The session stops issuing new model requests once the
	// tracked list cost reaches `max_list_cost`.
	Budget ManagedAgentsBudgetLimit `json:"budget" api:"required"`
	// A timestamp in RFC 3339 format
	CreatedAt     time.Time         `json:"created_at" api:"required" format:"date-time"`
	EnvironmentID string            `json:"environment_id" api:"required"`
	Metadata      map[string]string `json:"metadata" api:"required"`
	// Per-outcome evaluation state. One entry per `define_outcome` event sent to the
	// session.
	OutcomeEvaluations []ManagedAgentsOutcomeEvaluationResource `json:"outcome_evaluations" api:"required"`
	Resources          []ManagedAgentsSessionResourceUnion      `json:"resources" api:"required"`
	// Timing statistics for a session.
	Stats ManagedAgentsSessionStats `json:"stats" api:"required"`
	// SessionStatus enum
	//
	// Any of "rescheduling", "running", "idle", "terminated".
	Status ManagedAgentsSessionStatus `json:"status" api:"required"`
	Title  string                     `json:"title" api:"required"`
	// Any of "session".
	Type ManagedAgentsSessionType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// Cumulative token usage for a session across all turns.
	Usage ManagedAgentsSessionUsage `json:"usage" api:"required"`
	// Vault IDs attached to the session at creation. Empty when no vaults were
	// supplied.
	VaultIDs []string `json:"vault_ids" api:"required"`
	// Deployment ID when the session was created from a deployment reference. Null
	// otherwise.
	DeploymentID string `json:"deployment_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EnvironmentVariables respjson.Field

		ID                 respjson.Field
		Agent              respjson.Field
		ArchivedAt         respjson.Field
		Budget             respjson.Field
		CreatedAt          respjson.Field
		EnvironmentID      respjson.Field
		Metadata           respjson.Field
		OutcomeEvaluations respjson.Field
		Resources          respjson.Field
		Stats              respjson.Field
		Status             respjson.Field
		Title              respjson.Field
		Type               respjson.Field
		UpdatedAt          respjson.Field
		Usage              respjson.Field
		VaultIDs           respjson.Field
		DeploymentID       respjson.Field
		ExtraFields        map[string]respjson.Field
		raw                string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSession) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSession) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// SessionStatus enum
type ManagedAgentsSessionStatus string

const (
	ManagedAgentsSessionStatusRescheduling ManagedAgentsSessionStatus = "rescheduling"
	ManagedAgentsSessionStatusRunning      ManagedAgentsSessionStatus = "running"
	ManagedAgentsSessionStatusIdle         ManagedAgentsSessionStatus = "idle"
	ManagedAgentsSessionStatusTerminated   ManagedAgentsSessionStatus = "terminated"
)

type ManagedAgentsSessionType string

const (
	ManagedAgentsSessionTypeSession ManagedAgentsSessionType = "session"
)

// Resolved `agent` definition for a `session`. Snapshot of the `agent` at
// `session` creation time.
type ManagedAgentsSessionAgent struct {
	ID          string                                `json:"id" api:"required"`
	Description string                                `json:"description" api:"required"`
	MCPServers  []ManagedAgentsMCPServerURLDefinition `json:"mcp_servers" api:"required"`
	// Model identifier and configuration.
	Model ManagedAgentsModelConfig `json:"model" api:"required"`
	// Resolved coordinator topology with full agent definitions for each roster
	// member.
	Multiagent ManagedAgentsSessionMultiagentCoordinator `json:"multiagent" api:"required"`
	Name       string                                    `json:"name" api:"required"`
	Skills     []ManagedAgentsSessionAgentSkillUnion     `json:"skills" api:"required"`
	System     string                                    `json:"system" api:"required"`
	Tools      []ManagedAgentsSessionAgentToolUnion      `json:"tools" api:"required"`
	// Any of "agent".
	Type    ManagedAgentsSessionAgentType `json:"type" api:"required"`
	Version int64                         `json:"version" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Description respjson.Field
		MCPServers  respjson.Field
		Model       respjson.Field
		Multiagent  respjson.Field
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
func (r ManagedAgentsSessionAgent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionAgent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsSessionAgentSkillUnion contains all possible properties and
// values from [ManagedAgentsQoderSkill], [ManagedAgentsCustomSkill].
//
// Use the [ManagedAgentsSessionAgentSkillUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsSessionAgentSkillUnion struct {
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

// anyManagedAgentsSessionAgentSkill is implemented by each variant of
// [ManagedAgentsSessionAgentSkillUnion] to add type safety for the return type
// of [ManagedAgentsSessionAgentSkillUnion.AsAny]
type anyManagedAgentsSessionAgentSkill interface {
	implManagedAgentsSessionAgentSkillUnion()
}

func (ManagedAgentsQoderSkill) implManagedAgentsSessionAgentSkillUnion()  {}
func (ManagedAgentsCustomSkill) implManagedAgentsSessionAgentSkillUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsSessionAgentSkillUnion.AsAny().(type) {
//	case qoder.ManagedAgentsQoderSkill:
//	case qoder.ManagedAgentsCustomSkill:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsSessionAgentSkillUnion) AsAny() anyManagedAgentsSessionAgentSkill {
	switch u.Type {
	case "qoder":
		return u.AsQoder()
	case "custom":
		return u.AsCustom()
	}
	return nil
}

func (u ManagedAgentsSessionAgentSkillUnion) AsQoder() (v ManagedAgentsQoderSkill) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionAgentSkillUnion) AsCustom() (v ManagedAgentsCustomSkill) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsSessionAgentSkillUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsSessionAgentSkillUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsSessionAgentToolUnion contains all possible properties and
// values from [ManagedAgentsAgentToolset20260401],
// [ManagedAgentsMCPToolset], [ManagedAgentsCustomTool].
//
// Use the [ManagedAgentsSessionAgentToolUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsSessionAgentToolUnion struct {
	EnabledTools    []string `json:"enabled_tools"`
	DisallowedTools []string `json:"disallowed_tools"`
	// This field is a union of [[]ManagedAgentsAgentToolConfigUnion],
	// [[]ManagedAgentsMCPToolConfig]
	Configs ManagedAgentsSessionAgentToolUnionConfigs `json:"configs"`
	// This field is a union of [ManagedAgentsAgentToolsetDefaultConfig],
	// [ManagedAgentsMCPToolsetDefaultConfig]
	DefaultConfig ManagedAgentsSessionAgentToolUnionDefaultConfig `json:"default_config"`
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

// anyManagedAgentsSessionAgentTool is implemented by each variant of
// [ManagedAgentsSessionAgentToolUnion] to add type safety for the return type
// of [ManagedAgentsSessionAgentToolUnion.AsAny]
type anyManagedAgentsSessionAgentTool interface {
	implManagedAgentsSessionAgentToolUnion()
}

func (ManagedAgentsAgentToolset20260401) implManagedAgentsSessionAgentToolUnion() {}
func (ManagedAgentsMCPToolset) implManagedAgentsSessionAgentToolUnion()           {}
func (ManagedAgentsCustomTool) implManagedAgentsSessionAgentToolUnion()           {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsSessionAgentToolUnion.AsAny().(type) {
//	case qoder.ManagedAgentsAgentToolset20260401:
//	case qoder.ManagedAgentsMCPToolset:
//	case qoder.ManagedAgentsCustomTool:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsSessionAgentToolUnion) AsAny() anyManagedAgentsSessionAgentTool {
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

func (u ManagedAgentsSessionAgentToolUnion) AsAgentToolset20260401() (v ManagedAgentsAgentToolset20260401) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionAgentToolUnion) AsMCPToolset() (v ManagedAgentsMCPToolset) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionAgentToolUnion) AsCustom() (v ManagedAgentsCustomTool) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsSessionAgentToolUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsSessionAgentToolUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsSessionAgentToolUnionConfigs is an implicit subunion of
// [ManagedAgentsSessionAgentToolUnion].
// ManagedAgentsSessionAgentToolUnionConfigs provides convenient access to the
// sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsSessionAgentToolUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfManagedAgentsAgentToolConfigArray
// OfManagedAgentsMCPToolConfigArray]
type ManagedAgentsSessionAgentToolUnionConfigs struct {
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

func (r *ManagedAgentsSessionAgentToolUnionConfigs) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsSessionAgentToolUnionDefaultConfig is an implicit subunion of
// [ManagedAgentsSessionAgentToolUnion].
// ManagedAgentsSessionAgentToolUnionDefaultConfig provides convenient access
// to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsSessionAgentToolUnion].
type ManagedAgentsSessionAgentToolUnionDefaultConfig struct {
	Enabled bool `json:"enabled"`
	// This field is a union of
	// [ManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion],
	// [ManagedAgentsMCPToolsetDefaultConfigPermissionPolicyUnion]
	PermissionPolicy ManagedAgentsSessionAgentToolUnionDefaultConfigPermissionPolicy `json:"permission_policy"`
	JSON             struct {
		Enabled          respjson.Field
		PermissionPolicy respjson.Field
		raw              string
	} `json:"-"`
}

func (r *ManagedAgentsSessionAgentToolUnionDefaultConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsSessionAgentToolUnionDefaultConfigPermissionPolicy is an
// implicit subunion of [ManagedAgentsSessionAgentToolUnion].
// ManagedAgentsSessionAgentToolUnionDefaultConfigPermissionPolicy provides
// convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsSessionAgentToolUnion].
type ManagedAgentsSessionAgentToolUnionDefaultConfigPermissionPolicy struct {
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

func (r *ManagedAgentsSessionAgentToolUnionDefaultConfigPermissionPolicy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSessionAgentType string

const (
	ManagedAgentsSessionAgentTypeAgent ManagedAgentsSessionAgentType = "agent"
)

// Mid-session agent configuration update. Only `tools` and `mcp_servers` are
// updatable. Full replacement: the provided array becomes the new value. To
// preserve existing entries, GET the session, modify the array, and POST it back.
type ManagedAgentsSessionAgentUpdateParam struct {
	Model  ManagedAgentsModelConfigParams  `json:"model,omitzero"`
	System param.Opt[string]               `json:"system,omitzero"`
	Skills []ManagedAgentsSkillParamsUnion `json:"skills,omitzero"`

	// Replacement MCP server list. Full replacement: the provided array becomes the
	// new value. Send an empty array to clear; omit to preserve.
	MCPServers []ManagedAgentsURLMCPServerParams `json:"mcp_servers,omitzero"`
	// Replacement tool list. Full replacement: the provided array becomes the new
	// value. Send an empty array to clear; omit to preserve.
	Tools []ManagedAgentsSessionAgentUpdateToolUnionParam `json:"tools,omitzero"`
	paramObj
}

func (r ManagedAgentsSessionAgentUpdateParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsSessionAgentUpdateParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsSessionAgentUpdateParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsSessionAgentUpdateToolUnionParam struct {
	OfAgentToolset20260401 *ManagedAgentsAgentToolset20260401Params `json:",omitzero,inline"`
	OfMCPToolset           *ManagedAgentsMCPToolsetParams           `json:",omitzero,inline"`
	OfCustom               *ManagedAgentsCustomToolParams           `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsSessionAgentUpdateToolUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfAgentToolset20260401, u.OfMCPToolset, u.OfCustom)
}
func (u *ManagedAgentsSessionAgentUpdateToolUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsSessionAgentUpdateToolUnionParam) asAny() any {
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
func (u ManagedAgentsSessionAgentUpdateToolUnionParam) GetMCPServerName() *string {
	if vt := u.OfMCPToolset; vt != nil {
		return &vt.MCPServerName
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsSessionAgentUpdateToolUnionParam) GetDescription() *string {
	if vt := u.OfCustom; vt != nil {
		return &vt.Description
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsSessionAgentUpdateToolUnionParam) GetInputSchema() *ManagedAgentsCustomToolInputSchemaParam {
	if vt := u.OfCustom; vt != nil {
		return &vt.InputSchema
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsSessionAgentUpdateToolUnionParam) GetName() *string {
	if vt := u.OfCustom; vt != nil {
		return &vt.Name
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsSessionAgentUpdateToolUnionParam) GetType() *string {
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
func (u ManagedAgentsSessionAgentUpdateToolUnionParam) GetConfigs() (res managedAgentsSessionAgentUpdateToolUnionParamConfigs) {
	if vt := u.OfAgentToolset20260401; vt != nil {
		res.any = &vt.Configs
	} else if vt := u.OfMCPToolset; vt != nil {
		res.any = &vt.Configs
	}
	return
}

// Can have the runtime types [_[]ManagedAgentsAgentToolConfigParamsUnion],
// [_[]ManagedAgentsMCPToolConfigParams]
type managedAgentsSessionAgentUpdateToolUnionParamConfigs struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *[]qoder.ManagedAgentsAgentToolConfigParamsUnion:
//	case *[]qoder.ManagedAgentsMCPToolConfigParams:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u managedAgentsSessionAgentUpdateToolUnionParamConfigs) AsAny() any { return u.any }

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u ManagedAgentsSessionAgentUpdateToolUnionParam) GetDefaultConfig() (res managedAgentsSessionAgentUpdateToolUnionParamDefaultConfig) {
	if vt := u.OfAgentToolset20260401; vt != nil {
		res.any = &vt.DefaultConfig
	} else if vt := u.OfMCPToolset; vt != nil {
		res.any = &vt.DefaultConfig
	}
	return
}

// Can have the runtime types [*ManagedAgentsAgentToolsetDefaultConfigParams],
// [*ManagedAgentsMCPToolsetDefaultConfigParams]
type managedAgentsSessionAgentUpdateToolUnionParamDefaultConfig struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *qoder.ManagedAgentsAgentToolsetDefaultConfigParams:
//	case *qoder.ManagedAgentsMCPToolsetDefaultConfigParams:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u managedAgentsSessionAgentUpdateToolUnionParamDefaultConfig) AsAny() any { return u.any }

// Returns a pointer to the underlying variant's property, if present.
func (u managedAgentsSessionAgentUpdateToolUnionParamDefaultConfig) GetEnabled() *bool {
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
func (u managedAgentsSessionAgentUpdateToolUnionParamDefaultConfig) GetPermissionPolicy() (res managedAgentsSessionAgentUpdateToolUnionParamDefaultConfigPermissionPolicy) {
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
type managedAgentsSessionAgentUpdateToolUnionParamDefaultConfigPermissionPolicy struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *qoder.ManagedAgentsAlwaysAllowPolicyParam:
//	case *qoder.ManagedAgentsAlwaysAskPolicyParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u managedAgentsSessionAgentUpdateToolUnionParamDefaultConfigPermissionPolicy) AsAny() any {
	return u.any
}

// Returns a pointer to the underlying variant's property, if present.
func (u managedAgentsSessionAgentUpdateToolUnionParamDefaultConfigPermissionPolicy) GetType() *string {
	switch vt := u.any.(type) {
	case *ManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion:
		return vt.GetType()
	case *ManagedAgentsMCPToolsetDefaultConfigParamsPermissionPolicyUnion:
		return vt.GetType()
	}
	return nil
}

func init() {
	apijson.RegisterUnion[ManagedAgentsSessionAgentUpdateToolUnionParam](
		"type",
		apijson.Discriminator[ManagedAgentsAgentToolset20260401Params]("agent_toolset_20260401"),
		apijson.Discriminator[ManagedAgentsMCPToolsetParams]("mcp_toolset"),
		apijson.Discriminator[ManagedAgentsCustomToolParams]("custom"),
	)
}

// Resolved coordinator topology with full agent definitions for each roster
// member.
type ManagedAgentsSessionMultiagentCoordinator struct {
	// Full `agent` definitions the coordinator may spawn as session threads.
	Agents []ManagedAgentsSessionMultiagentCoordinatorAgentUnion `json:"agents" api:"required"`
	// Any of "coordinator".
	Type ManagedAgentsSessionMultiagentCoordinatorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Agents      respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionMultiagentCoordinator) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionMultiagentCoordinator) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsSessionMultiagentCoordinatorAgentUnion contains all possible
// properties and values from [ManagedAgentsSessionThreadAgent],
// [ManagedAgentsAdvisor].
//
// Use the [ManagedAgentsSessionMultiagentCoordinatorAgentUnion.AsAny] method
// to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsSessionMultiagentCoordinatorAgentUnion struct {
	// This field is from variant [ManagedAgentsSessionThreadAgent].
	ID string `json:"id"`
	// This field is from variant [ManagedAgentsSessionThreadAgent].
	Description string `json:"description"`
	// This field is from variant [ManagedAgentsSessionThreadAgent].
	MCPServers []ManagedAgentsMCPServerURLDefinition `json:"mcp_servers"`
	// This field is a union of [ManagedAgentsModelConfig], [string]
	Model ManagedAgentsSessionMultiagentCoordinatorAgentUnionModel `json:"model"`
	// This field is from variant [ManagedAgentsSessionThreadAgent].
	Name string `json:"name"`
	// This field is from variant [ManagedAgentsSessionThreadAgent].
	Skills []ManagedAgentsSessionThreadAgentSkillUnion `json:"skills"`
	// This field is from variant [ManagedAgentsSessionThreadAgent].
	System string `json:"system"`
	// This field is from variant [ManagedAgentsSessionThreadAgent].
	Tools []ManagedAgentsSessionThreadAgentToolUnion `json:"tools"`
	// Any of "agent", "advisor".
	Type string `json:"type"`
	// This field is from variant [ManagedAgentsSessionThreadAgent].
	Version int64 `json:"version"`
	JSON    struct {
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
		raw         string
	} `json:"-"`
}

// anyManagedAgentsSessionMultiagentCoordinatorAgent is implemented by each
// variant of [ManagedAgentsSessionMultiagentCoordinatorAgentUnion] to add type
// safety for the return type of
// [ManagedAgentsSessionMultiagentCoordinatorAgentUnion.AsAny]
type anyManagedAgentsSessionMultiagentCoordinatorAgent interface {
	implManagedAgentsSessionMultiagentCoordinatorAgentUnion()
}

func (ManagedAgentsSessionThreadAgent) implManagedAgentsSessionMultiagentCoordinatorAgentUnion() {
}
func (ManagedAgentsAdvisor) implManagedAgentsSessionMultiagentCoordinatorAgentUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsSessionMultiagentCoordinatorAgentUnion.AsAny().(type) {
//	case qoder.ManagedAgentsSessionThreadAgent:
//	case qoder.ManagedAgentsAdvisor:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsSessionMultiagentCoordinatorAgentUnion) AsAny() anyManagedAgentsSessionMultiagentCoordinatorAgent {
	switch u.Type {
	case "agent":
		return u.AsAgent()
	case "advisor":
		return u.AsAdvisor()
	}
	return nil
}

func (u ManagedAgentsSessionMultiagentCoordinatorAgentUnion) AsAgent() (v ManagedAgentsSessionThreadAgent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionMultiagentCoordinatorAgentUnion) AsAdvisor() (v ManagedAgentsAdvisor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsSessionMultiagentCoordinatorAgentUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsSessionMultiagentCoordinatorAgentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsSessionMultiagentCoordinatorAgentUnionModel is an implicit
// subunion of [ManagedAgentsSessionMultiagentCoordinatorAgentUnion].
// ManagedAgentsSessionMultiagentCoordinatorAgentUnionModel provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsSessionMultiagentCoordinatorAgentUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString]
type ManagedAgentsSessionMultiagentCoordinatorAgentUnionModel struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString string `json:",inline"`
	// This field is from variant [ManagedAgentsModelConfig].
	ID ManagedAgentsModel `json:"id"`
	// This field is from variant [ManagedAgentsModelConfig].
	Effort ManagedAgentsModelConfigEffortUnion `json:"effort"`
	// This field is from variant [ManagedAgentsModelConfig].
	InferenceGeo string `json:"inference_geo"`
	// This field is from variant [ManagedAgentsModelConfig].
	Speed ManagedAgentsModelConfigSpeed `json:"speed"`
	JSON  struct {
		OfString     respjson.Field
		ID           respjson.Field
		Effort       respjson.Field
		InferenceGeo respjson.Field
		Speed        respjson.Field
		raw          string
	} `json:"-"`
}

func (r *ManagedAgentsSessionMultiagentCoordinatorAgentUnionModel) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSessionMultiagentCoordinatorType string

const (
	ManagedAgentsSessionMultiagentCoordinatorTypeCoordinator ManagedAgentsSessionMultiagentCoordinatorType = "coordinator"
)

// Timing statistics for a session.
type ManagedAgentsSessionStats struct {
	// Cumulative time in seconds the session spent in `running` status. Excludes idle
	// time.
	ActiveSeconds float64 `json:"active_seconds"`
	// Elapsed time since session creation in seconds. For terminated sessions, frozen
	// at the final update.
	DurationSeconds float64 `json:"duration_seconds"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ActiveSeconds   respjson.Field
		DurationSeconds respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionStats) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionStats) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Emitted when an UpdateSession request changed at least one field. Carries only
// the fields that changed; absent fields were not part of the update. The new
// configuration applies from the next turn.
type ManagedAgentsSessionUpdatedEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "session.updated".
	Type ManagedAgentsSessionUpdatedEventType `json:"type" api:"required"`
	// Resolved `agent` definition for a `session`. Snapshot of the `agent` at
	// `session` creation time.
	Agent ManagedAgentsSessionAgent `json:"agent" api:"nullable"`
	// A hard spend ceiling. The session stops issuing new model requests once the
	// tracked list cost reaches `max_list_cost`.
	Budget ManagedAgentsBudgetLimit `json:"budget" api:"nullable"`
	// The session's full metadata bag after the update. Present when the update set
	// non-empty metadata; absent when metadata was unchanged or cleared to empty.
	Metadata map[string]string `json:"metadata"`
	// The session's new title. Present only when the update changed it.
	Title string `json:"title" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ProcessedAt respjson.Field
		Type        respjson.Field
		Agent       respjson.Field
		Budget      respjson.Field
		Metadata    respjson.Field
		Title       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionUpdatedEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionUpdatedEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSessionUpdatedEventType string

const (
	ManagedAgentsSessionUpdatedEventTypeSessionUpdated ManagedAgentsSessionUpdatedEventType = "session.updated"
)

// Cumulative token usage for a session across all turns.
type ManagedAgentsSessionUsage struct {
	// Cumulative time in seconds during which the session had at least one thread in
	// running status. Overlapping activity from concurrent threads is counted once,
	// unlike `stats.active_seconds`, which sums each thread's own active time. This is
	// the duration the session's runtime cost is priced on.
	ActiveSeconds float64 `json:"active_seconds"`
	// Prompt-cache creation token usage broken down by cache lifetime.
	CacheCreation ManagedAgentsCacheCreationUsage `json:"cache_creation"`
	// Total tokens read from prompt cache.
	CacheReadInputTokens int64 `json:"cache_read_input_tokens"`
	// Total input tokens consumed across all turns.
	InputTokens int64 `json:"input_tokens"`
	// A monetary amount in a specific currency.
	ListCost MonetaryAmount `json:"list_cost" api:"nullable"`
	// Total output tokens generated across all turns.
	OutputTokens int64 `json:"output_tokens"`
	// Cumulative count of server-executed tool invocations, broken down by tool.
	ServerToolUse ManagedAgentsServerToolUsage `json:"server_tool_use" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ActiveSeconds        respjson.Field
		CacheCreation        respjson.Field
		CacheReadInputTokens respjson.Field
		InputTokens          respjson.Field
		ListCost             respjson.Field
		OutputTokens         respjson.Field
		ServerToolUse        respjson.Field
		ExtraFields          map[string]respjson.Field
		raw                  string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionUsage) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Periodic snapshot of the session's cumulative usage and tracked list cost.
type ManagedAgentsSessionUsageEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"required" format:"date-time"`
	// Any of "session.usage".
	Type ManagedAgentsSessionUsageEventType `json:"type" api:"required"`
	// Point-in-time snapshot of a session's cumulative usage.
	Usage ManagedAgentsSessionUsageSnapshot `json:"usage" api:"required"`
	// A hard spend ceiling. The session stops issuing new model requests once the
	// tracked list cost reaches `max_list_cost`.
	Budget ManagedAgentsBudgetLimit `json:"budget" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		ProcessedAt respjson.Field
		Type        respjson.Field
		Usage       respjson.Field
		Budget      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionUsageEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionUsageEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSessionUsageEventType string

const (
	ManagedAgentsSessionUsageEventTypeSessionUsage ManagedAgentsSessionUsageEventType = "session.usage"
)

// Opens a preview of a buffered event. Carries the previewed event's type and id
// only. Followed by zero or more event_delta events with the same event id,
// normally concluded by the buffered event carrying that id. If the producing
// model request ends without that event (an error or interrupt mid-stream), its
// terminal span.model_request_end closes the preview. Only sent on stream
// connections that opt in via event_deltas; never appears in event history.
type ManagedAgentsStartEvent struct {
	// The previewed event's type and id. The event type determines which delta types
	// the preview's event_delta events carry: agent.message events stream
	// content_delta fragments; agent.thinking previews are start-only — no deltas
	// follow, and the buffered agent.thinking with the same id concludes them.
	Event ManagedAgentsStartEventPreviewUnion `json:"event" api:"required"`
	// Any of "event_start".
	Type ManagedAgentsStartEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Event       respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsStartEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsStartEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsStartEventType string

const (
	ManagedAgentsStartEventTypeEventStart ManagedAgentsStartEventType = "event_start"
)

// ManagedAgentsStartEventPreviewUnion contains all possible properties and
// values from [ManagedAgentsAgentMessagePreview],
// [ManagedAgentsAgentThinkingPreview].
//
// Use the [ManagedAgentsStartEventPreviewUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsStartEventPreviewUnion struct {
	ID string `json:"id"`
	// Any of "agent.message", "agent.thinking".
	Type string `json:"type"`
	JSON struct {
		ID   respjson.Field
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyManagedAgentsStartEventPreview is implemented by each variant of
// [ManagedAgentsStartEventPreviewUnion] to add type safety for the return type
// of [ManagedAgentsStartEventPreviewUnion.AsAny]
type anyManagedAgentsStartEventPreview interface {
	implManagedAgentsStartEventPreviewUnion()
}

func (ManagedAgentsAgentMessagePreview) implManagedAgentsStartEventPreviewUnion()  {}
func (ManagedAgentsAgentThinkingPreview) implManagedAgentsStartEventPreviewUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsStartEventPreviewUnion.AsAny().(type) {
//	case qoder.ManagedAgentsAgentMessagePreview:
//	case qoder.ManagedAgentsAgentThinkingPreview:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsStartEventPreviewUnion) AsAny() anyManagedAgentsStartEventPreview {
	switch u.Type {
	case "agent.message":
		return u.AsAgentMessage()
	case "agent.thinking":
		return u.AsAgentThinking()
	}
	return nil
}

func (u ManagedAgentsStartEventPreviewUnion) AsAgentMessage() (v ManagedAgentsAgentMessagePreview) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsStartEventPreviewUnion) AsAgentThinking() (v ManagedAgentsAgentThinkingPreview) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsStartEventPreviewUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsStartEventPreviewUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Regular text content.
type ManagedAgentsSystemContentBlock struct {
	// The text content.
	Text string `json:"text" api:"required"`
	// Any of "text".
	Type ManagedAgentsSystemContentBlockType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Text        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSystemContentBlock) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSystemContentBlock) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this ManagedAgentsSystemContentBlock to a
// ManagedAgentsSystemContentBlockParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// ManagedAgentsSystemContentBlockParam.Overrides()
func (r ManagedAgentsSystemContentBlock) ToParam() ManagedAgentsSystemContentBlockParam {
	return param.Override[ManagedAgentsSystemContentBlockParam](json.RawMessage(r.RawJSON()))
}

type ManagedAgentsSystemContentBlockType string

const (
	ManagedAgentsSystemContentBlockTypeText ManagedAgentsSystemContentBlockType = "text"
)

// Regular text content.
//
// The properties Text, Type are required.
type ManagedAgentsSystemContentBlockParam struct {
	// The text content.
	Text string `json:"text" api:"required"`
	// Any of "text".
	Type ManagedAgentsSystemContentBlockType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsSystemContentBlockParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsSystemContentBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsSystemContentBlockParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A mid-conversation system message event. Carries system-role content that is
// appended to the session as a `role: "system"` turn.
type ManagedAgentsSystemMessageEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// System content blocks. Text-only.
	Content []ManagedAgentsSystemContentBlock `json:"content" api:"required"`
	// Any of "system.message".
	Type ManagedAgentsSystemMessageEventType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"nullable" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Content     respjson.Field
		Type        respjson.Field
		ProcessedAt respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSystemMessageEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSystemMessageEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSystemMessageEventType string

const (
	ManagedAgentsSystemMessageEventTypeSystemMessage ManagedAgentsSystemMessageEventType = "system.message"
)

// Event sent by the client providing the result of an agent-toolset tool
// execution. Only valid on `self_hosted` environments, where sandbox-routed tools
// are executed by the client rather than the server.
type ManagedAgentsUserToolResultEvent struct {
	// Unique identifier for this event.
	ID string `json:"id" api:"required"`
	// The id of the `agent.tool_use` event this result corresponds to, which can be
	// found in the last `session.status_idle`
	// [event's](https://docs.qoder.com/cloud-agents/api/sessions/schemas)
	// `stop_reason.event_ids` field.
	ToolUseID string `json:"tool_use_id" api:"required"`
	// Any of "user.tool_result".
	Type ManagedAgentsUserToolResultEventType `json:"type" api:"required"`
	// The result content returned by the tool.
	Content []ManagedAgentsUserToolResultEventContentUnion `json:"content"`
	// Whether the tool execution resulted in an error.
	IsError bool `json:"is_error" api:"nullable"`
	// A timestamp in RFC 3339 format
	ProcessedAt time.Time `json:"processed_at" api:"nullable" format:"date-time"`
	// Routes this result to a subagent thread. Copy from the `agent.tool_use` event's
	// `session_thread_id`.
	SessionThreadID string `json:"session_thread_id" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID              respjson.Field
		ToolUseID       respjson.Field
		Type            respjson.Field
		Content         respjson.Field
		IsError         respjson.Field
		ProcessedAt     respjson.Field
		SessionThreadID respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsUserToolResultEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsUserToolResultEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsUserToolResultEventType string

const (
	ManagedAgentsUserToolResultEventTypeUserToolResult ManagedAgentsUserToolResultEventType = "user.tool_result"
)

// ManagedAgentsUserToolResultEventContentUnion contains all possible
// properties and values from [ManagedAgentsTextBlock],
// [ManagedAgentsImageBlock], [ManagedAgentsDocumentBlock],
// [ManagedAgentsSearchResultBlock].
//
// Use the [ManagedAgentsUserToolResultEventContentUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsUserToolResultEventContentUnion struct {
	// This field is from variant [ManagedAgentsTextBlock].
	Text string `json:"text"`
	// Any of "text", "image", "document", "search_result".
	Type string `json:"type"`
	// This field is a union of [ManagedAgentsImageBlockSourceUnion],
	// [ManagedAgentsDocumentBlockSourceUnion], [string]
	Source ManagedAgentsUserToolResultEventContentUnionSource `json:"source"`
	// This field is from variant [ManagedAgentsDocumentBlock].
	Context string `json:"context"`
	Title   string `json:"title"`
	// This field is from variant [ManagedAgentsSearchResultBlock].
	Citations ManagedAgentsSearchResultCitations `json:"citations"`
	// This field is from variant [ManagedAgentsSearchResultBlock].
	Content []ManagedAgentsSearchResultContent `json:"content"`
	JSON    struct {
		Text      respjson.Field
		Type      respjson.Field
		Source    respjson.Field
		Context   respjson.Field
		Title     respjson.Field
		Citations respjson.Field
		Content   respjson.Field
		raw       string
	} `json:"-"`
}

// anyManagedAgentsUserToolResultEventContent is implemented by each variant of
// [ManagedAgentsUserToolResultEventContentUnion] to add type safety for the
// return type of [ManagedAgentsUserToolResultEventContentUnion.AsAny]
type anyManagedAgentsUserToolResultEventContent interface {
	implManagedAgentsUserToolResultEventContentUnion()
}

func (ManagedAgentsTextBlock) implManagedAgentsUserToolResultEventContentUnion()         {}
func (ManagedAgentsImageBlock) implManagedAgentsUserToolResultEventContentUnion()        {}
func (ManagedAgentsDocumentBlock) implManagedAgentsUserToolResultEventContentUnion()     {}
func (ManagedAgentsSearchResultBlock) implManagedAgentsUserToolResultEventContentUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsUserToolResultEventContentUnion.AsAny().(type) {
//	case qoder.ManagedAgentsTextBlock:
//	case qoder.ManagedAgentsImageBlock:
//	case qoder.ManagedAgentsDocumentBlock:
//	case qoder.ManagedAgentsSearchResultBlock:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsUserToolResultEventContentUnion) AsAny() anyManagedAgentsUserToolResultEventContent {
	switch u.Type {
	case "text":
		return u.AsText()
	case "image":
		return u.AsImage()
	case "document":
		return u.AsDocument()
	case "search_result":
		return u.AsSearchResult()
	}
	return nil
}

func (u ManagedAgentsUserToolResultEventContentUnion) AsText() (v ManagedAgentsTextBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsUserToolResultEventContentUnion) AsImage() (v ManagedAgentsImageBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsUserToolResultEventContentUnion) AsDocument() (v ManagedAgentsDocumentBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsUserToolResultEventContentUnion) AsSearchResult() (v ManagedAgentsSearchResultBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsUserToolResultEventContentUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsUserToolResultEventContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsUserToolResultEventContentUnionSource is an implicit subunion
// of [ManagedAgentsUserToolResultEventContentUnion].
// ManagedAgentsUserToolResultEventContentUnionSource provides convenient
// access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsUserToolResultEventContentUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfString]
type ManagedAgentsUserToolResultEventContentUnionSource struct {
	// This field will be present if the value is a [string] instead of an object.
	OfString  string `json:",inline"`
	Data      string `json:"data"`
	MediaType string `json:"media_type"`
	Type      string `json:"type"`
	URL       string `json:"url"`
	FileID    string `json:"file_id"`
	JSON      struct {
		OfString  respjson.Field
		Data      respjson.Field
		MediaType respjson.Field
		Type      respjson.Field
		URL       respjson.Field
		FileID    respjson.Field
		raw       string
	} `json:"-"`
}

func (r *ManagedAgentsUserToolResultEventContentUnionSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SessionNewParams struct {
	EnvironmentVariables map[string]string `json:"environment_variables,omitzero"`

	// Agent identifier. Accepts the `agent` ID string, which pins the latest version
	// for the session, or an `agent` object with both id and version specified.
	Agent SessionNewParamsAgentUnion `json:"agent,omitzero" api:"required"`
	// ID of the `environment` defining the container configuration for this session.
	EnvironmentID string `json:"environment_id" api:"required"`
	// Human-readable session title.
	Title param.Opt[string] `json:"title,omitzero"`
	// A hard spend ceiling. The session stops issuing new model requests once the
	// tracked list cost reaches `max_list_cost`.
	Budget ManagedAgentsBudgetLimitParam `json:"budget,omitzero"`
	// Initial events to send to the `session` at creation, processed in order.
	// Supports `user.message` and `user.define_outcome` events. Maximum 50 events.
	InitialEvents []SessionNewParamsInitialEventUnion `json:"initial_events,omitzero"`
	// Arbitrary key-value metadata attached to the session. Maximum 16 pairs, keys up
	// to 64 chars, values up to 512 chars.
	Metadata map[string]string `json:"metadata,omitzero"`
	// Resources (e.g. repositories, files) to mount into the session's container.
	Resources []SessionNewParamsResourceUnion `json:"resources,omitzero"`
	// Vault IDs for stored credentials the agent can use during the session.
	VaultIDs []string `json:"vault_ids,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

func (r SessionNewParams) MarshalJSON() (data []byte, err error) {
	type shadow SessionNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type SessionNewParamsAgentUnion struct {
	OfString                           param.Opt[string]                      `json:",omitzero,inline"`
	OfManagedAgentsAgents              *ManagedAgentsAgentParams              `json:",omitzero,inline"`
	OfManagedAgentsAgentWithOverridess *ManagedAgentsAgentWithOverridesParams `json:",omitzero,inline"`
	paramUnion
}

func (u SessionNewParamsAgentUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfManagedAgentsAgents, u.OfManagedAgentsAgentWithOverridess)
}
func (u *SessionNewParamsAgentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *SessionNewParamsAgentUnion) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfManagedAgentsAgents) {
		return u.OfManagedAgentsAgents
	} else if !param.IsOmitted(u.OfManagedAgentsAgentWithOverridess) {
		return u.OfManagedAgentsAgentWithOverridess
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionNewParamsAgentUnion) GetMCPServers() []ManagedAgentsURLMCPServerParams {
	if vt := u.OfManagedAgentsAgentWithOverridess; vt != nil {
		return vt.MCPServers
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionNewParamsAgentUnion) GetModel() *ManagedAgentsModelConfigParams {
	if vt := u.OfManagedAgentsAgentWithOverridess; vt != nil {
		return &vt.Model
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionNewParamsAgentUnion) GetSkills() []ManagedAgentsSkillParamsUnion {
	if vt := u.OfManagedAgentsAgentWithOverridess; vt != nil {
		return vt.Skills
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionNewParamsAgentUnion) GetSystem() *string {
	if vt := u.OfManagedAgentsAgentWithOverridess; vt != nil && vt.System.Valid() {
		return &vt.System.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionNewParamsAgentUnion) GetTools() []ManagedAgentsAgentWithOverridesParamsToolUnion {
	if vt := u.OfManagedAgentsAgentWithOverridess; vt != nil {
		return vt.Tools
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionNewParamsAgentUnion) GetID() *string {
	if vt := u.OfManagedAgentsAgents; vt != nil {
		return (*string)(&vt.ID)
	} else if vt := u.OfManagedAgentsAgentWithOverridess; vt != nil {
		return (*string)(&vt.ID)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionNewParamsAgentUnion) GetType() *string {
	if vt := u.OfManagedAgentsAgents; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfManagedAgentsAgentWithOverridess; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionNewParamsAgentUnion) GetVersion() *int64 {
	if vt := u.OfManagedAgentsAgents; vt != nil && vt.Version.Valid() {
		return &vt.Version.Value
	} else if vt := u.OfManagedAgentsAgentWithOverridess; vt != nil && vt.Version.Valid() {
		return &vt.Version.Value
	}
	return nil
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type SessionNewParamsInitialEventUnion struct {
	OfUserMessage       *ManagedAgentsUserMessageEventParams       `json:",omitzero,inline"`
	OfUserDefineOutcome *ManagedAgentsUserDefineOutcomeEventParams `json:",omitzero,inline"`
	paramUnion
}

func (u SessionNewParamsInitialEventUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfUserMessage, u.OfUserDefineOutcome)
}
func (u *SessionNewParamsInitialEventUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *SessionNewParamsInitialEventUnion) asAny() any {
	if !param.IsOmitted(u.OfUserMessage) {
		return u.OfUserMessage
	} else if !param.IsOmitted(u.OfUserDefineOutcome) {
		return u.OfUserDefineOutcome
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionNewParamsInitialEventUnion) GetContent() []ManagedAgentsUserMessageEventParamsContentUnion {
	if vt := u.OfUserMessage; vt != nil {
		return vt.Content
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionNewParamsInitialEventUnion) GetDescription() *string {
	if vt := u.OfUserDefineOutcome; vt != nil {
		return &vt.Description
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionNewParamsInitialEventUnion) GetRubric() *ManagedAgentsUserDefineOutcomeEventParamsRubricUnion {
	if vt := u.OfUserDefineOutcome; vt != nil {
		return &vt.Rubric
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionNewParamsInitialEventUnion) GetMaxIterations() *int64 {
	if vt := u.OfUserDefineOutcome; vt != nil && vt.MaxIterations.Valid() {
		return &vt.MaxIterations.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionNewParamsInitialEventUnion) GetType() *string {
	if vt := u.OfUserMessage; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfUserDefineOutcome; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[SessionNewParamsInitialEventUnion](
		"type",
		apijson.Discriminator[ManagedAgentsUserMessageEventParams]("user.message"),
		apijson.Discriminator[ManagedAgentsUserDefineOutcomeEventParams]("user.define_outcome"),
	)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type SessionNewParamsResourceUnion struct {
	OfGitHubRepository *ManagedAgentsGitHubRepositoryResourceParams `json:",omitzero,inline"`
	OfFile             *ManagedAgentsFileResourceParams             `json:",omitzero,inline"`
	OfMemoryStore      *ManagedAgentsMemoryStoreResourceParam       `json:",omitzero,inline"`
	paramUnion
}

func (u SessionNewParamsResourceUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfGitHubRepository, u.OfFile, u.OfMemoryStore)
}
func (u *SessionNewParamsResourceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *SessionNewParamsResourceUnion) asAny() any {
	if !param.IsOmitted(u.OfGitHubRepository) {
		return u.OfGitHubRepository
	} else if !param.IsOmitted(u.OfFile) {
		return u.OfFile
	} else if !param.IsOmitted(u.OfMemoryStore) {
		return u.OfMemoryStore
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionNewParamsResourceUnion) GetAuthorizationToken() *string {
	if vt := u.OfGitHubRepository; vt != nil {
		return &vt.AuthorizationToken
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionNewParamsResourceUnion) GetURL() *string {
	if vt := u.OfGitHubRepository; vt != nil {
		return &vt.URL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionNewParamsResourceUnion) GetCheckout() *ManagedAgentsGitHubRepositoryResourceParamsCheckoutUnion {
	if vt := u.OfGitHubRepository; vt != nil {
		return &vt.Checkout
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionNewParamsResourceUnion) GetFileID() *string {
	if vt := u.OfFile; vt != nil {
		return &vt.FileID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionNewParamsResourceUnion) GetMemoryStoreID() *string {
	if vt := u.OfMemoryStore; vt != nil {
		return &vt.MemoryStoreID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionNewParamsResourceUnion) GetAccess() *string {
	if vt := u.OfMemoryStore; vt != nil {
		return (*string)(&vt.Access)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionNewParamsResourceUnion) GetInstructions() *string {
	if vt := u.OfMemoryStore; vt != nil && vt.Instructions.Valid() {
		return &vt.Instructions.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionNewParamsResourceUnion) GetType() *string {
	if vt := u.OfGitHubRepository; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFile; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfMemoryStore; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u SessionNewParamsResourceUnion) GetMountPath() *string {
	if vt := u.OfGitHubRepository; vt != nil && vt.MountPath.Valid() {
		return &vt.MountPath.Value
	} else if vt := u.OfFile; vt != nil && vt.MountPath.Valid() {
		return &vt.MountPath.Value
	}
	return nil
}

func init() {
	apijson.RegisterUnion[SessionNewParamsResourceUnion](
		"type",
		apijson.Discriminator[ManagedAgentsGitHubRepositoryResourceParams]("github_repository"),
		apijson.Discriminator[ManagedAgentsFileResourceParams]("file"),
		apijson.Discriminator[ManagedAgentsMemoryStoreResourceParam]("memory_store"),
	)
}

type SessionGetParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

type SessionUpdateParams struct {
	EnvironmentVariables map[string]string `json:"environment_variables,omitzero"`

	// Human-readable session title.
	Title param.Opt[string] `json:"title,omitzero"`
	// Metadata patch. Set a key to a string to upsert it, or to null to delete it.
	// Omit the field to preserve.
	Metadata map[string]any `json:"metadata,omitzero"`
	// Mid-session agent configuration update. Only `tools` and `mcp_servers` are
	// updatable. Full replacement: the provided array becomes the new value. To
	// preserve existing entries, GET the session, modify the array, and POST it back.
	Agent ManagedAgentsSessionAgentUpdateParam `json:"agent,omitzero"`
	// A hard spend ceiling. The session stops issuing new model requests once the
	// tracked list cost reaches `max_list_cost`.
	Budget ManagedAgentsBudgetLimitParam `json:"budget,omitzero"`
	// Vault IDs (`vlt_*`) to attach to the session. Not yet supported; requests
	// setting this field are rejected. Reserved for future use.
	VaultIDs []string `json:"vault_ids,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

func (r SessionUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow SessionUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SessionListParams struct {
	// Filter sessions created with this agent ID.
	AgentID param.Opt[string] `query:"agent_id,omitzero" json:"-"`
	// Filter by agent version. Only applies when `agent_id` is also set.
	AgentVersion param.Opt[int64] `query:"agent_version,omitzero" json:"-"`
	// Return sessions created after this time (exclusive).
	CreatedAtGt param.Opt[time.Time] `query:"created_at[gt],omitzero" format:"date-time" json:"-"`
	// Return sessions created at or after this time (inclusive).
	CreatedAtGte param.Opt[time.Time] `query:"created_at[gte],omitzero" format:"date-time" json:"-"`
	// Return sessions created before this time (exclusive).
	CreatedAtLt param.Opt[time.Time] `query:"created_at[lt],omitzero" format:"date-time" json:"-"`
	// Return sessions created at or before this time (inclusive).
	CreatedAtLte param.Opt[time.Time] `query:"created_at[lte],omitzero" format:"date-time" json:"-"`
	// Filter sessions created by this deployment ID.
	DeploymentID param.Opt[string] `query:"deployment_id,omitzero" json:"-"`
	// When true, includes archived sessions. Default: false (exclude archived).
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// Maximum number of results to return.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Filter sessions whose resources contain a `memory_store` with this memory store
	// ID.
	MemoryStoreID param.Opt[string] `query:"memory_store_id,omitzero" json:"-"`
	// Opaque pagination cursor from a previous response.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Sort direction for results, ordered by `created_at`. Defaults to `desc` (newest
	// first).
	//
	// Any of "asc", "desc".
	Order SessionListParamsOrder `query:"order,omitzero" json:"-"`
	// Filter by session status. Repeat the parameter to match any of multiple
	// statuses.
	//
	// Any of "rescheduling", "running", "idle", "terminated".
	Statuses []string `query:"statuses,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [SessionListParams]'s query parameters as `url.Values`.
func (r SessionListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Sort direction for results, ordered by `created_at`. Defaults to `desc` (newest
// first).
type SessionListParamsOrder string

const (
	SessionListParamsOrderAsc  SessionListParamsOrder = "asc"
	SessionListParamsOrderDesc SessionListParamsOrder = "desc"
)

type SessionDeleteParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

type SessionArchiveParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}
