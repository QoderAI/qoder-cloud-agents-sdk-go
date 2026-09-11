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

// TemplateService provides Forward Template operations.
type TemplateService struct {
	Options []option.RequestOption
}

func NewTemplateService(opts ...option.RequestOption) TemplateService {
	return TemplateService{Options: slices.Clone(opts)}
}

// 列出 Templates.
func (r *TemplateService) List(ctx context.Context, params TemplateListParams, opts ...option.RequestOption) (res *pagination.Page[Template], err error) {

	opts = slices.Concat(r.Options, opts)
	path := "templates"
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
func (r *TemplateService) ListAutoPaging(ctx context.Context, params TemplateListParams, opts ...option.RequestOption) *pagination.PageAutoPager[Template] {
	return pagination.NewPageAutoPager(r.List(ctx, params, opts...))
}

type TemplateListParams struct {
	// 按 `active` 或 `archived` 过滤。
	Status param.Opt[string] `query:"status,omitzero" json:"-"`
	// 分页大小，最大 100。
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// 向后翻页游标，不能与 `before_id` 同用。
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// 向前翻页游标，不能与 `after_id` 同用。
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	paramObj
}

func (r TemplateListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// 创建 Template.
func (r *TemplateService) New(ctx context.Context, params TemplateNewParams, opts ...option.RequestOption) (res *Template, err error) {
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	if params.QoderBeta.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("X-Qoder-Beta", fmt.Sprint(params.QoderBeta.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := "templates"
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type TemplateNewParams struct {
	// Template 名称，1-256 个字符，租户内唯一。
	Name string `json:"name" api:"required"`
	// 模型标识。可传 string（如 `"ultimate"`），或传 Agent model 对象以同时配置 `effort` 或 `context_window`。可通过列出模型接口查询可用值。
	Model ModelConfigUnionParam `json:"model" api:"required"`
	// 创建 Session 时默认使用的 Environment ID。
	EnvironmentID string `json:"environment_id" api:"required"`
	// Template 描述，最多 2048 个字符。
	Description param.Opt[string] `json:"description,omitzero"`
	// System Prompt，最多 100,000 个字符。
	System param.Opt[string] `json:"system,omitzero"`
	// 工具配置列表，最多 128 项。
	Tools []ToolParam `json:"tools,omitzero"`
	// MCP Server 配置列表，最多 20 项。
	MCPServers []MCPServerParam `json:"mcp_servers,omitzero"`
	// Skill 绑定列表，最多 20 项。
	Skills []SkillBindingParam `json:"skills,omitzero"`
	// Multi-agent 协作配置。`type` 必须为 `coordinator`；省略或传 `null` 表示不启用。
	Multiagent MultiagentConfigParam `json:"multiagent,omitzero"`
	// 默认 Vault 配置，按 Vault ID 组织。
	Vaults map[string]ResourceBindingParam `json:"vaults,omitzero"`
	// 默认文件资源配置，按 file ID 组织。
	Files map[string]ResourceBindingParam `json:"files,omitzero"`
	// 默认 GitHub 仓库配置，按调用方指定的 binding key 组织，最多 20 项。
	GitHubRepositories map[string]GitHubRepositoryParam `json:"github_repositories,omitzero"`
	// 默认 Session 环境变量。
	EnvironmentVariables EnvironmentVariablesUnionParam `json:"environment_variables,omitzero"`
	// 自定义元数据。
	Metadata map[string]any `json:"metadata,omitzero"`
	// 有副作用请求可选的幂等键。
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	// 启用 Browser Use 时必须设置为 `browser-use-2026-07-14`。
	QoderBeta param.Opt[string] `header:"X-Qoder-Beta,omitzero" json:"-"`
	paramObj
}

func (r TemplateNewParams) MarshalJSON() ([]byte, error) {
	type shadow TemplateNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *TemplateNewParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

// 获取 Template.
func (r *TemplateService) Get(ctx context.Context, templateID string, opts ...option.RequestOption) (res *Template, err error) {
	if templateID == "" {
		return nil, fmt.Errorf("missing required template_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("templates/%s", url.PathEscape(templateID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// 更新 Template.
func (r *TemplateService) Update(ctx context.Context, templateID string, params TemplateUpdateParams, opts ...option.RequestOption) (res *Template, err error) {
	if templateID == "" {
		return nil, fmt.Errorf("missing required template_id parameter")
	}
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	if params.QoderBeta.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("X-Qoder-Beta", fmt.Sprint(params.QoderBeta.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("templates/%s", url.PathEscape(templateID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type TemplateUpdateParams struct {
	// 新的 Template 名称。
	Name param.Opt[string] `json:"name,omitzero"`
	// 新的 Template 描述。
	Description param.Opt[string] `json:"description,omitzero"`
	// 新的模型标识。可传 string，或传 Agent model 对象以同时配置 `effort` 或 `context_window`。可通过列出模型接口查询可用值。
	Model ModelConfigUnionParam `json:"model,omitzero"`
	// 新的 System Prompt。
	System param.Opt[string] `json:"system,omitzero"`
	// 整体替换工具配置列表。
	Tools []ToolParam `json:"tools,omitzero"`
	// 整体替换 MCP Server 列表。
	MCPServers []MCPServerParam `json:"mcp_servers,omitzero"`
	// 整体替换 Skill 绑定列表。
	Skills []SkillBindingParam `json:"skills,omitzero"`
	// 整体替换 Multi-agent 协作配置；传 `null` 表示清空，省略则保留当前配置。
	Multiagent MultiagentConfigParam `json:"multiagent,omitzero"`
	// 替换默认 Environment ID；`null` 或空字符串表示清空。
	EnvironmentID param.Opt[string] `json:"environment_id,omitzero" api:"nullable"`
	// 整体替换默认 Vault 配置；按 Vault ID 组织，`null` 表示清空。
	Vaults map[string]ResourceBindingParam `json:"vaults,omitzero"`
	// 整体替换默认文件资源配置；`null` 表示清空。
	Files map[string]ResourceBindingParam `json:"files,omitzero"`
	// 整体替换默认 GitHub 仓库配置；按 binding key 组织，`null` 或空 object 表示清空。
	GitHubRepositories map[string]GitHubRepositoryParam `json:"github_repositories,omitzero"`
	// 整体替换默认环境变量；`null` 表示清空。
	EnvironmentVariables EnvironmentVariablesUnionParam `json:"environment_variables,omitzero"`
	// 合并更新自定义元数据。
	Metadata map[string]any `json:"metadata,omitzero"`
	// 有副作用请求可选的幂等键。
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	// 更新后的 `tools` 中包含 Browser Use 工具集时，必须设置为 `browser-use-2026-07-14`。
	QoderBeta param.Opt[string] `header:"X-Qoder-Beta,omitzero" json:"-"`
	paramObj
}

func (r TemplateUpdateParams) MarshalJSON() ([]byte, error) {
	type shadow TemplateUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *TemplateUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// 归档 Template.
func (r *TemplateService) Archive(ctx context.Context, templateID string, params TemplateArchiveParams, opts ...option.RequestOption) (res *Template, err error) {
	if templateID == "" {
		return nil, fmt.Errorf("missing required template_id parameter")
	}
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("templates/%s/archive", url.PathEscape(templateID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type TemplateArchiveParams struct {
	// 有副作用请求可选的幂等键。
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r TemplateArchiveParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// 克隆 Template.
func (r *TemplateService) Clone(ctx context.Context, templateID string, params TemplateCloneParams, opts ...option.RequestOption) (res *Template, err error) {
	if templateID == "" {
		return nil, fmt.Errorf("missing required template_id parameter")
	}
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("templates/%s/clone", url.PathEscape(templateID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type TemplateCloneParams struct {
	// 新 Template 名称；不传时使用 `<源名称> Copy <随机短 ID>`。
	Name param.Opt[string] `json:"name,omitzero"`
	// 新 Template 描述；不传时沿用源描述。
	Description param.Opt[string] `json:"description,omitzero"`
	// 有副作用请求可选的幂等键。
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r TemplateCloneParams) MarshalJSON() ([]byte, error) {
	type shadow TemplateCloneParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *TemplateCloneParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type Template struct {
	// 固定为 `template`。
	Type string `json:"type"`
	// Template ID。
	ID string `json:"id"`
	// Template 名称。
	Name        string `json:"name"`
	Description string `json:"description"`
	// `active` 或 `archived`。
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at" format:"date-time"`
	UpdatedAt time.Time `json:"updated_at" format:"date-time"`
	// 与请求提交的形态一致；对象形态保留 `id`、`effort` 和 `context_window`。
	Model      ModelConfig      `json:"model"`
	Multiagent MultiagentConfig `json:"multiagent"`
	// 默认 Environment ID。
	EnvironmentID        string                      `json:"environment_id"`
	Vaults               map[string]ResourceBinding  `json:"vaults"`
	Files                map[string]ResourceBinding  `json:"files"`
	GitHubRepositories   map[string]GitHubRepository `json:"github_repositories"`
	System               string                      `json:"system"`
	Tools                []Tool                      `json:"tools"`
	MCPServers           []MCPServer                 `json:"mcp_servers"`
	Skills               []SkillBinding              `json:"skills"`
	EnvironmentVariables map[string]string           `json:"environment_variables"`
	// 自定义元数据。
	Metadata map[string]any `json:"metadata"`
	JSON     struct {
		Type                 respjson.Field
		ID                   respjson.Field
		Name                 respjson.Field
		Description          respjson.Field
		Status               respjson.Field
		CreatedAt            respjson.Field
		UpdatedAt            respjson.Field
		Model                respjson.Field
		Multiagent           respjson.Field
		EnvironmentID        respjson.Field
		Vaults               respjson.Field
		Files                respjson.Field
		GitHubRepositories   respjson.Field
		System               respjson.Field
		Tools                respjson.Field
		MCPServers           respjson.Field
		Skills               respjson.Field
		EnvironmentVariables respjson.Field
		Metadata             respjson.Field
		ExtraFields          map[string]respjson.Field
		raw                  string
	} `json:"-"`
}

func (r Template) RawJSON() string                  { return r.JSON.raw }
func (r *Template) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
