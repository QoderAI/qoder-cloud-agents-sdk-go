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

// List Templates
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
	// Filter by `active` or `archived`.
	Status param.Opt[string] `query:"status,omitzero" json:"-"`
	// Page size, maximum 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Cursor for the next page; cannot be combined with `before_id`.
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// Cursor for the previous page; cannot be combined with `after_id`.
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	paramObj
}

func (r TemplateListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// Create Template
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
	// Template name, 1-256 characters, unique per user.
	Name string `json:"name" api:"required"`
	// Model identifier. Accepts a string (such as `"ultimate"`), or an Agent model
	// object to configure `effort` or `context_window` at the same time. Use the
	// list models endpoint to discover the available values.
	Model ModelConfigUnionParam `json:"model" api:"required"`
	// Environment ID used by default when creating a Session.
	EnvironmentID string `json:"environment_id" api:"required"`
	// Template description, up to 2048 characters.
	Description param.Opt[string] `json:"description,omitzero"`
	// System Prompt, up to 100,000 characters.
	System param.Opt[string] `json:"system,omitzero"`
	// Tool configuration list, up to 128 entries.
	Tools []ToolParam `json:"tools,omitzero"`
	// MCP Server configuration list, up to 20 entries.
	MCPServers []MCPServerParam `json:"mcp_servers,omitzero"`
	// Skill binding list, up to 20 entries.
	Skills []SkillBindingParam `json:"skills,omitzero"`
	// Multi-agent collaboration configuration. `type` must be `coordinator`; omit
	// it or pass `null` to leave it disabled.
	Multiagent MultiagentConfigParam `json:"multiagent,omitzero"`
	// Default Vault configuration, keyed by Vault ID.
	Vaults map[string]ResourceBindingParam `json:"vaults,omitzero"`
	// Default file resource configuration, keyed by file ID.
	Files map[string]ResourceBindingParam `json:"files,omitzero"`
	// Default GitHub repository configuration, keyed by a caller-chosen binding
	// key, up to 20 entries.
	GitHubRepositories map[string]GitHubRepositoryParam `json:"github_repositories,omitzero"`
	// Default Session environment variables.
	EnvironmentVariables EnvironmentVariablesUnionParam `json:"environment_variables,omitzero"`
	// Custom metadata.
	Metadata map[string]any `json:"metadata,omitzero"`
	// Optional idempotency key for requests with side effects.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	// Must be set to `browser-use-2026-07-14` when Browser Use is enabled.
	QoderBeta param.Opt[string] `header:"X-Qoder-Beta,omitzero" json:"-"`
	paramObj
}

func (r TemplateNewParams) MarshalJSON() ([]byte, error) {
	type shadow TemplateNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *TemplateNewParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

// Get Template
func (r *TemplateService) Get(ctx context.Context, templateID string, opts ...option.RequestOption) (res *Template, err error) {
	if templateID == "" {
		return nil, fmt.Errorf("missing required template_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("templates/%s", url.PathEscape(templateID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update Template
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
	// New Template name.
	Name param.Opt[string] `json:"name,omitzero"`
	// New Template description.
	Description param.Opt[string] `json:"description,omitzero"`
	// New model identifier. Accepts a string, or an Agent model object to configure
	// `effort` or `context_window` at the same time. Use the list models endpoint
	// to discover the available values.
	Model ModelConfigUnionParam `json:"model,omitzero"`
	// New System Prompt.
	System param.Opt[string] `json:"system,omitzero"`
	// Replaces the tool configuration list wholesale.
	Tools []ToolParam `json:"tools,omitzero"`
	// Replaces the MCP Server list wholesale.
	MCPServers []MCPServerParam `json:"mcp_servers,omitzero"`
	// Replaces the Skill binding list wholesale.
	Skills []SkillBindingParam `json:"skills,omitzero"`
	// Replaces the Multi-agent collaboration configuration wholesale; `null` clears
	// it, omitting it keeps the current configuration.
	Multiagent MultiagentConfigParam `json:"multiagent,omitzero"`
	// Replaces the default Environment ID; `null` or an empty string clears it.
	EnvironmentID param.Opt[string] `json:"environment_id,omitzero" api:"nullable"`
	// Replaces the default Vault configuration wholesale, keyed by Vault ID;
	// `null` clears it.
	Vaults map[string]ResourceBindingParam `json:"vaults,omitzero"`
	// Replaces the default file resource configuration wholesale; `null` clears it.
	Files map[string]ResourceBindingParam `json:"files,omitzero"`
	// Replaces the default GitHub repository configuration wholesale, keyed by
	// binding key; `null` or an empty object clears it.
	GitHubRepositories map[string]GitHubRepositoryParam `json:"github_repositories,omitzero"`
	// Replaces the default environment variables wholesale; `null` clears it.
	EnvironmentVariables EnvironmentVariablesUnionParam `json:"environment_variables,omitzero"`
	// Merges updates into the custom metadata.
	Metadata map[string]any `json:"metadata,omitzero"`
	// Optional idempotency key for requests with side effects.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	// Must be set to `browser-use-2026-07-14` when the updated `tools` include the
	// Browser Use toolset.
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

// Archive Template
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
	// Optional idempotency key for requests with side effects.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r TemplateArchiveParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// Clone Template
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
	// Name of the new Template; defaults to `<source name> Copy <random short ID>`.
	Name param.Opt[string] `json:"name,omitzero"`
	// Description of the new Template; defaults to the source description.
	Description param.Opt[string] `json:"description,omitzero"`
	// Optional idempotency key for requests with side effects.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r TemplateCloneParams) MarshalJSON() ([]byte, error) {
	type shadow TemplateCloneParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *TemplateCloneParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type Template struct {
	// Always `template`.
	Type string `json:"type"`
	// Template ID.
	ID string `json:"id"`
	// Template name.
	Name        string `json:"name"`
	Description string `json:"description"`
	// Either `active` or `archived`.
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at" format:"date-time"`
	UpdatedAt time.Time `json:"updated_at" format:"date-time"`
	// Mirrors the shape submitted in the request; the object form keeps `id`,
	// `effort` and `context_window`.
	Model      ModelConfig      `json:"model"`
	Multiagent MultiagentConfig `json:"multiagent"`
	// Default Environment ID.
	EnvironmentID        string                      `json:"environment_id"`
	Vaults               map[string]ResourceBinding  `json:"vaults"`
	Files                map[string]ResourceBinding  `json:"files"`
	GitHubRepositories   map[string]GitHubRepository `json:"github_repositories"`
	System               string                      `json:"system"`
	Tools                []Tool                      `json:"tools"`
	MCPServers           []MCPServer                 `json:"mcp_servers"`
	Skills               []SkillBinding              `json:"skills"`
	EnvironmentVariables map[string]string           `json:"environment_variables"`
	// Custom metadata.
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
