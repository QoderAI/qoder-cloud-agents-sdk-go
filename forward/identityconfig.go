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

// IdentityConfigService provides Forward IdentityConfig operations.
type IdentityConfigService struct {
	Options []option.RequestOption
}

func NewIdentityConfigService(opts ...option.RequestOption) IdentityConfigService {
	return IdentityConfigService{Options: slices.Clone(opts)}
}

// List Identity Configs
func (r *IdentityConfigService) List(ctx context.Context, identityID string, params IdentityConfigListParams, opts ...option.RequestOption) (res *pagination.Page[IdentityConfig], err error) {
	if identityID == "" {
		return nil, fmt.Errorf("missing required identity_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("identities/%s/templates", url.PathEscape(identityID))
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
func (r *IdentityConfigService) ListAutoPaging(ctx context.Context, identityID string, params IdentityConfigListParams, opts ...option.RequestOption) *pagination.PageAutoPager[IdentityConfig] {
	return pagination.NewPageAutoPager(r.List(ctx, identityID, params, opts...))
}

type IdentityConfigListParams struct {
	// Filter by Forward Template ID.
	TemplateID param.Opt[string] `query:"template_id,omitzero" json:"-"`
	// Filter by `active` or `archived`.
	Status param.Opt[string] `query:"status,omitzero" json:"-"`
	// Page size, maximum 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Cursor for the next page, taken from `last_id` in the previous response.
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// Cursor for the previous page, taken from `first_id` in the previous response.
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	paramObj
}

func (r IdentityConfigListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// Get Identity Config
func (r *IdentityConfigService) Get(ctx context.Context, identityID string, templateID string, opts ...option.RequestOption) (res *IdentityConfig, err error) {
	if identityID == "" {
		return nil, fmt.Errorf("missing required identity_id parameter")
	}
	if templateID == "" {
		return nil, fmt.Errorf("missing required template_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("identities/%s/templates/%s/config", url.PathEscape(identityID), url.PathEscape(templateID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Create or update Identity Config
func (r *IdentityConfigService) Upsert(ctx context.Context, identityID string, templateID string, params IdentityConfigUpsertParams, opts ...option.RequestOption) (res *IdentityConfig, err error) {
	if identityID == "" {
		return nil, fmt.Errorf("missing required identity_id parameter")
	}
	if templateID == "" {
		return nil, fmt.Errorf("missing required template_id parameter")
	}
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("identities/%s/templates/%s/config", url.PathEscape(identityID), url.PathEscape(templateID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type IdentityConfigUpsertParams struct {
	// Display name for the Config.
	Name param.Opt[string] `json:"name,omitzero"`
	// User-level configuration overrides.
	IdentityConfig IdentityConfigSpecParam `json:"identity_config" api:"required"`
	// Business metadata; when sent it replaces the existing metadata entirely.
	Metadata map[string]any `json:"metadata,omitzero"`
	// Optional idempotency key for requests with side effects.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r IdentityConfigUpsertParams) MarshalJSON() ([]byte, error) {
	type shadow IdentityConfigUpsertParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *IdentityConfigUpsertParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Get Effective Config
func (r *IdentityConfigService) GetEffective(ctx context.Context, identityID string, templateID string, opts ...option.RequestOption) (res *EffectiveConfig, err error) {
	if identityID == "" {
		return nil, fmt.Errorf("missing required identity_id parameter")
	}
	if templateID == "" {
		return nil, fmt.Errorf("missing required template_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("identities/%s/templates/%s/effective", url.PathEscape(identityID), url.PathEscape(templateID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

type EffectiveConfig struct {
	// Always `effective_spec`.
	Type string `json:"type"`
	// Hash of the compiled Agent section.
	AgentEffectiveHash string `json:"agent_effective_hash"`
	// Hash of the compiled Session section.
	SessionEffectiveHash string `json:"session_effective_hash"`
	// Hash of the full effective configuration.
	EffectiveHash string `json:"effective_hash"`
	// Compiled Agent configuration.
	Agent EffectiveConfigAgent `json:"agent"`
	// Compiled Session defaults.
	Session    EffectiveConfigSession `json:"session"`
	ID         string                 `json:"id"`
	IdentityID string                 `json:"identity_id"`
	TemplateID string                 `json:"template_id"`
	JSON       struct {
		Type                 respjson.Field
		AgentEffectiveHash   respjson.Field
		SessionEffectiveHash respjson.Field
		EffectiveHash        respjson.Field
		Agent                respjson.Field
		Session              respjson.Field
		ID                   respjson.Field
		IdentityID           respjson.Field
		TemplateID           respjson.Field
		ExtraFields          map[string]respjson.Field
		raw                  string
	} `json:"-"`
}

func (r EffectiveConfig) RawJSON() string                  { return r.JSON.raw }
func (r *EffectiveConfig) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type EffectiveConfigSession struct {
	EnvironmentID        string                                `json:"environment_id"`
	EnvironmentVariables map[string]string                     `json:"environment_variables"`
	VaultIDs             []string                              `json:"vault_ids"`
	Resources            []EffectiveConfigSessionResourcesItem `json:"resources"`
	JSON                 struct {
		EnvironmentID        respjson.Field
		EnvironmentVariables respjson.Field
		VaultIDs             respjson.Field
		Resources            respjson.Field
		ExtraFields          map[string]respjson.Field
		raw                  string
	} `json:"-"`
}

func (r EffectiveConfigSession) RawJSON() string { return r.JSON.raw }
func (r *EffectiveConfigSession) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type EffectiveConfigSessionResourcesItem struct {
	Type      string `json:"type"`
	FileID    string `json:"file_id"`
	URL       string `json:"url"`
	MountPath string `json:"mount_path"`
	JSON      struct {
		Type        respjson.Field
		FileID      respjson.Field
		URL         respjson.Field
		MountPath   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r EffectiveConfigSessionResourcesItem) RawJSON() string { return r.JSON.raw }
func (r *EffectiveConfigSessionResourcesItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type EffectiveConfigAgent struct {
	Model      ModelConfig                     `json:"model"`
	System     string                          `json:"system"`
	Tools      []EffectiveConfigAgentToolsItem `json:"tools"`
	MCPServers []map[string]any                `json:"mcp_servers"`
	Skills     []map[string]any                `json:"skills"`
	JSON       struct {
		Model       respjson.Field
		System      respjson.Field
		Tools       respjson.Field
		MCPServers  respjson.Field
		Skills      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r EffectiveConfigAgent) RawJSON() string { return r.JSON.raw }
func (r *EffectiveConfigAgent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type EffectiveConfigAgentToolsItem struct {
	Type    string                                     `json:"type"`
	Configs []EffectiveConfigAgentToolsItemConfigsItem `json:"configs"`
	JSON    struct {
		Type        respjson.Field
		Configs     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r EffectiveConfigAgentToolsItem) RawJSON() string { return r.JSON.raw }
func (r *EffectiveConfigAgentToolsItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type EffectiveConfigAgentToolsItemConfigsItem struct {
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
	JSON    struct {
		Name        respjson.Field
		Enabled     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r EffectiveConfigAgentToolsItemConfigsItem) RawJSON() string { return r.JSON.raw }
func (r *EffectiveConfigAgentToolsItemConfigsItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type IdentityConfig struct {
	// Always `config`.
	Type string `json:"type"`
	ID   string `json:"id"`
	// Forward Identity ID.
	IdentityID string `json:"identity_id"`
	// Forward Template ID.
	TemplateID string `json:"template_id"`
	Name       string `json:"name"`
	Status     string `json:"status"`
	// Hash of the compiled Effective Config.
	EffectiveHash  string             `json:"effective_hash"`
	CreatedAt      time.Time          `json:"created_at" format:"date-time"`
	UpdatedAt      time.Time          `json:"updated_at" format:"date-time"`
	IdentityConfig IdentityConfigSpec `json:"identity_config"`
	// Business metadata.
	Metadata map[string]any `json:"metadata"`
	JSON     struct {
		Type           respjson.Field
		ID             respjson.Field
		IdentityID     respjson.Field
		TemplateID     respjson.Field
		Name           respjson.Field
		Status         respjson.Field
		EffectiveHash  respjson.Field
		CreatedAt      respjson.Field
		UpdatedAt      respjson.Field
		IdentityConfig respjson.Field
		Metadata       respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

func (r IdentityConfig) RawJSON() string                  { return r.JSON.raw }
func (r *IdentityConfig) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
