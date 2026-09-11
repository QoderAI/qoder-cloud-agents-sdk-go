package forward

import (
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
)

type ResourceBinding struct {
	Enabled bool `json:"enabled"`
	JSON    struct {
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r ResourceBinding) RawJSON() string                  { return r.JSON.raw }
func (r *ResourceBinding) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type ResourceBindingParam struct {
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	paramObj
}

func (r ResourceBindingParam) MarshalJSON() ([]byte, error) {
	type shadow ResourceBindingParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ResourceBindingParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type GitHubRepository struct {
	URL       string `json:"url"`
	MountPath string `json:"mount_path"`
	Enabled   bool   `json:"enabled"`
	JSON      struct {
		URL         respjson.Field
		MountPath   respjson.Field
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r GitHubRepository) RawJSON() string                  { return r.JSON.raw }
func (r *GitHubRepository) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type GitHubRepositoryParam struct {
	URL                param.Opt[string] `json:"url,omitzero"`
	MountPath          param.Opt[string] `json:"mount_path,omitzero"`
	Enabled            param.Opt[bool]   `json:"enabled,omitzero"`
	AuthorizationToken param.Opt[string] `json:"authorization_token,omitzero"`
	paramObj
}

func (r GitHubRepositoryParam) MarshalJSON() ([]byte, error) {
	type shadow GitHubRepositoryParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *GitHubRepositoryParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type PermissionPolicy struct {
	Type string `json:"type" api:"required"`
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r PermissionPolicy) RawJSON() string                  { return r.JSON.raw }
func (r *PermissionPolicy) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type PermissionPolicyParam struct {
	Type string `json:"type" api:"required"`
	paramObj
}

func (r PermissionPolicyParam) MarshalJSON() ([]byte, error) {
	type shadow PermissionPolicyParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *PermissionPolicyParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ToolConfig struct {
	Name             string           `json:"name" api:"required"`
	Enabled          bool             `json:"enabled"`
	PermissionPolicy PermissionPolicy `json:"permission_policy"`
	JSON             struct {
		Name             respjson.Field
		Enabled          respjson.Field
		PermissionPolicy respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

func (r ToolConfig) RawJSON() string                  { return r.JSON.raw }
func (r *ToolConfig) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type ToolConfigParam struct {
	Name             string                `json:"name" api:"required"`
	Enabled          param.Opt[bool]       `json:"enabled,omitzero"`
	PermissionPolicy PermissionPolicyParam `json:"permission_policy,omitzero"`
	paramObj
}

func (r ToolConfigParam) MarshalJSON() ([]byte, error) {
	type shadow ToolConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ToolConfigParam) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type Tool struct {
	Type            string         `json:"type" api:"required"`
	EnabledTools    []string       `json:"enabled_tools"`
	DisallowedTools []string       `json:"disallowed_tools"`
	Configs         []ToolConfig   `json:"configs"`
	MCPServerName   string         `json:"mcp_server_name"`
	Name            string         `json:"name"`
	Description     string         `json:"description"`
	InputSchema     map[string]any `json:"input_schema"`
	JSON            struct {
		Type            respjson.Field
		EnabledTools    respjson.Field
		DisallowedTools respjson.Field
		Configs         respjson.Field
		MCPServerName   respjson.Field
		Name            respjson.Field
		Description     respjson.Field
		InputSchema     respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

func (r Tool) RawJSON() string                  { return r.JSON.raw }
func (r *Tool) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type ToolParam struct {
	Type            string            `json:"type" api:"required"`
	EnabledTools    []string          `json:"enabled_tools,omitzero"`
	DisallowedTools []string          `json:"disallowed_tools,omitzero"`
	Configs         []ToolConfigParam `json:"configs,omitzero"`
	MCPServerName   param.Opt[string] `json:"mcp_server_name,omitzero"`
	Name            param.Opt[string] `json:"name,omitzero"`
	Description     param.Opt[string] `json:"description,omitzero"`
	InputSchema     map[string]any    `json:"input_schema,omitzero"`
	paramObj
}

func (r ToolParam) MarshalJSON() ([]byte, error) {
	type shadow ToolParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ToolParam) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type MCPServer struct {
	Type string `json:"type"`
	Name string `json:"name" api:"required"`
	URL  string `json:"url" api:"required"`
	JSON struct {
		Type        respjson.Field
		Name        respjson.Field
		URL         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r MCPServer) RawJSON() string                  { return r.JSON.raw }
func (r *MCPServer) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type MCPServerParam struct {
	Type param.Opt[string] `json:"type,omitzero"`
	Name string            `json:"name" api:"required"`
	URL  string            `json:"url" api:"required"`
	paramObj
}

func (r MCPServerParam) MarshalJSON() ([]byte, error) {
	type shadow MCPServerParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *MCPServerParam) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type SkillBinding struct {
	Type    string `json:"type" api:"required"`
	SkillID string `json:"skill_id" api:"required"`
	Version string `json:"version"`
	Enabled bool   `json:"enabled"`
	JSON    struct {
		Type        respjson.Field
		SkillID     respjson.Field
		Version     respjson.Field
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r SkillBinding) RawJSON() string                  { return r.JSON.raw }
func (r *SkillBinding) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type SkillBindingParam struct {
	Type    string            `json:"type" api:"required"`
	SkillID string            `json:"skill_id" api:"required"`
	Version param.Opt[string] `json:"version,omitzero"`
	Enabled param.Opt[bool]   `json:"enabled,omitzero"`
	paramObj
}

func (r SkillBindingParam) MarshalJSON() ([]byte, error) {
	type shadow SkillBindingParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SkillBindingParam) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type MultiagentEntry struct {
	Type       string `json:"type" api:"required"`
	TemplateID string `json:"template_id"`
	Name       string `json:"name"`
	JSON       struct {
		Type        respjson.Field
		TemplateID  respjson.Field
		Name        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r MultiagentEntry) RawJSON() string                  { return r.JSON.raw }
func (r *MultiagentEntry) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type MultiagentEntryParam struct {
	Type       string            `json:"type" api:"required"`
	TemplateID param.Opt[string] `json:"template_id,omitzero"`
	Name       param.Opt[string] `json:"name,omitzero"`
	paramObj
}

func (r MultiagentEntryParam) MarshalJSON() ([]byte, error) {
	type shadow MultiagentEntryParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *MultiagentEntryParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type MultiagentConfig struct {
	Type   string            `json:"type" api:"required"`
	Agents []MultiagentEntry `json:"agents" api:"required"`
	JSON   struct {
		Type        respjson.Field
		Agents      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r MultiagentConfig) RawJSON() string                  { return r.JSON.raw }
func (r *MultiagentConfig) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type MultiagentConfigParam struct {
	Type   string                 `json:"type" api:"required"`
	Agents []MultiagentEntryParam `json:"agents" api:"required"`
	paramObj
}

func (r MultiagentConfigParam) MarshalJSON() ([]byte, error) {
	type shadow MultiagentConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *MultiagentConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ImageSource struct {
	Type      string `json:"type" api:"required"`
	MediaType string `json:"media_type"`
	Data      string `json:"data"`
	URL       string `json:"url"`
	FileID    string `json:"file_id"`
	JSON      struct {
		Type        respjson.Field
		MediaType   respjson.Field
		Data        respjson.Field
		URL         respjson.Field
		FileID      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r ImageSource) RawJSON() string                  { return r.JSON.raw }
func (r *ImageSource) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type ImageSourceParam struct {
	Type      string            `json:"type" api:"required"`
	MediaType param.Opt[string] `json:"media_type,omitzero"`
	Data      param.Opt[string] `json:"data,omitzero"`
	URL       param.Opt[string] `json:"url,omitzero"`
	FileID    param.Opt[string] `json:"file_id,omitzero"`
	paramObj
}

func (r ImageSourceParam) MarshalJSON() ([]byte, error) {
	type shadow ImageSourceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ImageSourceParam) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type ContentBlock struct {
	Type     string      `json:"type" api:"required"`
	Text     string      `json:"text"`
	Thinking string      `json:"thinking"`
	Source   ImageSource `json:"source"`
	JSON     struct {
		Type        respjson.Field
		Text        respjson.Field
		Thinking    respjson.Field
		Source      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r ContentBlock) RawJSON() string                  { return r.JSON.raw }
func (r *ContentBlock) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type ContentBlockParam struct {
	Type     string            `json:"type" api:"required"`
	Text     param.Opt[string] `json:"text,omitzero"`
	Thinking param.Opt[string] `json:"thinking,omitzero"`
	Source   ImageSourceParam  `json:"source,omitzero"`
	paramObj
}

func (r ContentBlockParam) MarshalJSON() ([]byte, error) {
	type shadow ContentBlockParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ContentBlockParam) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type SessionResourceSpec struct {
	Type      string `json:"type" api:"required"`
	FileID    string `json:"file_id" api:"required"`
	MountPath string `json:"mount_path"`
	JSON      struct {
		Type        respjson.Field
		FileID      respjson.Field
		MountPath   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r SessionResourceSpec) RawJSON() string                  { return r.JSON.raw }
func (r *SessionResourceSpec) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type SessionResourceSpecParam struct {
	Type      string            `json:"type" api:"required"`
	FileID    string            `json:"file_id" api:"required"`
	MountPath param.Opt[string] `json:"mount_path,omitzero"`
	paramObj
}

func (r SessionResourceSpecParam) MarshalJSON() ([]byte, error) {
	type shadow SessionResourceSpecParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionResourceSpecParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type EnvironmentVariableOverride struct {
	Op    string `json:"op" api:"required"`
	Value string `json:"value"`
	JSON  struct {
		Op          respjson.Field
		Value       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r EnvironmentVariableOverride) RawJSON() string { return r.JSON.raw }
func (r *EnvironmentVariableOverride) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type EnvironmentVariableOverrideParam struct {
	Op    string            `json:"op" api:"required"`
	Value param.Opt[string] `json:"value,omitzero"`
	paramObj
}

func (r EnvironmentVariableOverrideParam) MarshalJSON() ([]byte, error) {
	type shadow EnvironmentVariableOverrideParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *EnvironmentVariableOverrideParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SystemOverride struct {
	Mode    string `json:"mode"`
	Content string `json:"content"`
	JSON    struct {
		Mode        respjson.Field
		Content     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r SystemOverride) RawJSON() string                  { return r.JSON.raw }
func (r *SystemOverride) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type SystemOverrideParam struct {
	Mode    param.Opt[string] `json:"mode,omitzero"`
	Content param.Opt[string] `json:"content,omitzero"`
	paramObj
}

func (r SystemOverrideParam) MarshalJSON() ([]byte, error) {
	type shadow SystemOverrideParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SystemOverrideParam) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type ToolOverride struct {
	Enabled          bool             `json:"enabled"`
	PermissionPolicy PermissionPolicy `json:"permission_policy"`
	JSON             struct {
		Enabled          respjson.Field
		PermissionPolicy respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

func (r ToolOverride) RawJSON() string                  { return r.JSON.raw }
func (r *ToolOverride) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type ToolOverrideParam struct {
	Enabled          param.Opt[bool]       `json:"enabled,omitzero"`
	PermissionPolicy PermissionPolicyParam `json:"permission_policy,omitzero"`
	paramObj
}

func (r ToolOverrideParam) MarshalJSON() ([]byte, error) {
	type shadow ToolOverrideParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ToolOverrideParam) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type MCPServerOverride struct {
	Enabled bool   `json:"enabled"`
	Type    string `json:"type"`
	URL     string `json:"url"`
	JSON    struct {
		Enabled     respjson.Field
		Type        respjson.Field
		URL         respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r MCPServerOverride) RawJSON() string                  { return r.JSON.raw }
func (r *MCPServerOverride) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type MCPServerOverrideParam struct {
	Enabled param.Opt[bool]   `json:"enabled,omitzero"`
	Type    param.Opt[string] `json:"type,omitzero"`
	URL     param.Opt[string] `json:"url,omitzero"`
	paramObj
}

func (r MCPServerOverrideParam) MarshalJSON() ([]byte, error) {
	type shadow MCPServerOverrideParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *MCPServerOverrideParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SkillOverride struct {
	Enabled bool   `json:"enabled"`
	Type    string `json:"type"`
	Version string `json:"version"`
	JSON    struct {
		Enabled     respjson.Field
		Type        respjson.Field
		Version     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r SkillOverride) RawJSON() string                  { return r.JSON.raw }
func (r *SkillOverride) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type SkillOverrideParam struct {
	Enabled param.Opt[bool]   `json:"enabled,omitzero"`
	Type    param.Opt[string] `json:"type,omitzero"`
	Version param.Opt[string] `json:"version,omitzero"`
	paramObj
}

func (r SkillOverrideParam) MarshalJSON() ([]byte, error) {
	type shadow SkillOverrideParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SkillOverrideParam) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type IdentityConfigSpec struct {
	System               SystemOverride                         `json:"system"`
	Model                ModelConfig                            `json:"model"`
	Tools                map[string]ToolOverride                `json:"tools"`
	MCPServers           map[string]MCPServerOverride           `json:"mcp_servers"`
	Skills               map[string]SkillOverride               `json:"skills"`
	Toolsets             map[string]any                         `json:"toolsets"`
	AgentMetadata        map[string]any                         `json:"agent_metadata"`
	Vaults               map[string]ResourceBinding             `json:"vaults"`
	Files                map[string]ResourceBinding             `json:"files"`
	GitHubRepositories   map[string]GitHubRepository            `json:"github_repositories"`
	EnvironmentVariables map[string]EnvironmentVariableOverride `json:"environment_variables"`
	JSON                 struct {
		System               respjson.Field
		Model                respjson.Field
		Tools                respjson.Field
		MCPServers           respjson.Field
		Skills               respjson.Field
		Toolsets             respjson.Field
		AgentMetadata        respjson.Field
		Vaults               respjson.Field
		Files                respjson.Field
		GitHubRepositories   respjson.Field
		EnvironmentVariables respjson.Field
		ExtraFields          map[string]respjson.Field
		raw                  string
	} `json:"-"`
}

func (r IdentityConfigSpec) RawJSON() string                  { return r.JSON.raw }
func (r *IdentityConfigSpec) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type IdentityConfigSpecParam struct {
	System               SystemOverrideParam                         `json:"system,omitzero"`
	Model                ModelConfigUnionParam                       `json:"model,omitzero"`
	Tools                map[string]ToolOverrideParam                `json:"tools,omitzero"`
	MCPServers           map[string]MCPServerOverrideParam           `json:"mcp_servers,omitzero"`
	Skills               map[string]SkillOverrideParam               `json:"skills,omitzero"`
	Toolsets             map[string]any                              `json:"toolsets,omitzero"`
	AgentMetadata        map[string]any                              `json:"agent_metadata,omitzero"`
	Vaults               map[string]ResourceBindingParam             `json:"vaults,omitzero"`
	Files                map[string]ResourceBindingParam             `json:"files,omitzero"`
	GitHubRepositories   map[string]GitHubRepositoryParam            `json:"github_repositories,omitzero"`
	EnvironmentVariables map[string]EnvironmentVariableOverrideParam `json:"environment_variables,omitzero"`
	paramObj
}

func (r IdentityConfigSpecParam) MarshalJSON() ([]byte, error) {
	type shadow IdentityConfigSpecParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *IdentityConfigSpecParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SessionEventParam struct {
	Type            string                 `json:"type" api:"required"`
	Content         EventContentUnionParam `json:"content,omitzero"`
	ToolUseID       param.Opt[string]      `json:"tool_use_id,omitzero"`
	CustomToolUseID param.Opt[string]      `json:"custom_tool_use_id,omitzero"`
	Result          param.Opt[string]      `json:"result,omitzero"`
	DenyMessage     param.Opt[string]      `json:"deny_message,omitzero"`
	IsError         param.Opt[bool]        `json:"is_error,omitzero"`
	Description     param.Opt[string]      `json:"description,omitzero"`
	Rubric          param.Opt[string]      `json:"rubric,omitzero"`
	OutcomeID       param.Opt[string]      `json:"outcome_id,omitzero"`
	MaxIterations   param.Opt[int64]       `json:"max_iterations,omitzero"`
	paramObj
}

func (r SessionEventParam) MarshalJSON() ([]byte, error) {
	type shadow SessionEventParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionEventParam) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
