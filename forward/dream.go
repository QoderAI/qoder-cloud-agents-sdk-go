package forward

import (
	"context"
	"fmt"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apiquery"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/pagination"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
	"net/http"
	"net/url"
	"slices"
	"time"
)

// DreamService consolidates Forward memory; requires PAT or an admin SAT.
type DreamService struct{ Options []option.RequestOption }

func NewDreamService(opts ...option.RequestOption) DreamService {
	return DreamService{Options: slices.Clone(opts)}
}
func (r *DreamService) New(ctx context.Context, params DreamNewParams, opts ...option.RequestOption) (res *Dream, err error) {
	opts = slices.Concat(r.Options, opts)
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, "dreams", params, &res, opts...)
	return res, err
}
func (r *DreamService) List(ctx context.Context, params DreamListParams, opts ...option.RequestOption) (res *pagination.Page[Dream], err error) {
	opts = slices.Concat(r.Options, opts)
	var raw *http.Response
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	cfg, err := convention.NewRequestConfig(ctx, http.MethodGet, "dreams", params, &res, opts...)
	if err != nil {
		return nil, err
	}
	if err = cfg.Execute(); err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}
func (r *DreamService) ListAutoPaging(ctx context.Context, params DreamListParams, opts ...option.RequestOption) *pagination.PageAutoPager[Dream] {
	return pagination.NewPageAutoPager(r.List(ctx, params, opts...))
}
func (r *DreamService) Get(ctx context.Context, dreamID string, opts ...option.RequestOption) (res *Dream, err error) {
	if dreamID == "" {
		return nil, fmt.Errorf("missing required dream_id parameter")
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("dreams/%s", url.PathEscape(dreamID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}
func (r *DreamService) Cancel(ctx context.Context, dreamID string, opts ...option.RequestOption) (res *Dream, err error) {
	if dreamID == "" {
		return nil, fmt.Errorf("missing required dream_id parameter")
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("dreams/%s/cancel", url.PathEscape(dreamID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}
func (r *DreamService) Archive(ctx context.Context, dreamID string, opts ...option.RequestOption) (res *Dream, err error) {
	if dreamID == "" {
		return nil, fmt.Errorf("missing required dream_id parameter")
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("dreams/%s/archive", url.PathEscape(dreamID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type DreamListParams struct {
	Limit           param.Opt[int64]  `query:"limit,omitzero" json:"-"`
	AfterID         param.Opt[string] `query:"after_id,omitzero" json:"-"`
	BeforeID        param.Opt[string] `query:"before_id,omitzero" json:"-"`
	IncludeArchived param.Opt[bool]   `query:"include_archived,omitzero" json:"-"`
	paramObj
}

func (r DreamListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{})
}

type DreamNewParams struct {
	Inputs       []DreamInputParam        `json:"inputs,omitzero"`
	Model        DreamNewParamsModelUnion `json:"model,omitzero"`
	Instructions param.Opt[string]        `json:"instructions,omitzero"`
	paramObj
}

func (r DreamNewParams) MarshalJSON() ([]byte, error) {
	type shadow DreamNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DreamNewParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type DreamInputParam struct {
	Type          string            `json:"type"`
	MemoryStoreID param.Opt[string] `json:"memory_store_id,omitzero"`
	IdentityID    param.Opt[string] `json:"identity_id,omitzero"`
	TemplateID    param.Opt[string] `json:"template_id,omitzero"`
	SessionIDs    []string          `json:"session_ids,omitzero"`
	paramObj
}

func (r DreamInputParam) MarshalJSON() ([]byte, error) {
	type shadow DreamInputParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DreamInputParam) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type DreamModelConfigParam struct {
	ID            string            `json:"id"`
	Effort        param.Opt[string] `json:"effort,omitzero"`
	ContextWindow param.Opt[int64]  `json:"context_window,omitzero"`
	paramObj
}

func (r DreamModelConfigParam) MarshalJSON() ([]byte, error) {
	type shadow DreamModelConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DreamModelConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DreamNewParamsModelUnion accepts a model ID or tuning object; omit to use auto.
type DreamNewParamsModelUnion struct {
	OfString param.Opt[string]      `json:",omitzero,inline"`
	OfConfig *DreamModelConfigParam `json:",omitzero,inline"`
	paramUnion
}

func (r DreamNewParamsModelUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(r, r.OfString, r.OfConfig)
}
func (r *DreamNewParamsModelUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type Dream struct {
	ID           string           `json:"id"`
	Type         string           `json:"type"`
	Status       string           `json:"status"`
	Inputs       []DreamInput     `json:"inputs"`
	Outputs      []DreamOutput    `json:"outputs"`
	Model        DreamModelConfig `json:"model"`
	Instructions string           `json:"instructions"`
	Usage        DreamUsage       `json:"usage"`
	Error        DreamError       `json:"error"`
	CreatedAt    time.Time        `json:"created_at"`
	EndedAt      time.Time        `json:"ended_at"`
	ArchivedAt   time.Time        `json:"archived_at"`
	JSON         struct {
		ID           respjson.Field
		Type         respjson.Field
		Status       respjson.Field
		Inputs       respjson.Field
		Outputs      respjson.Field
		Model        respjson.Field
		Instructions respjson.Field
		Usage        respjson.Field
		Error        respjson.Field
		CreatedAt    respjson.Field
		EndedAt      respjson.Field
		ArchivedAt   respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

func (r Dream) RawJSON() string                  { return r.JSON.raw }
func (r *Dream) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type DreamInput struct {
	Type          string   `json:"type"`
	MemoryStoreID string   `json:"memory_store_id"`
	IdentityID    string   `json:"identity_id"`
	TemplateID    string   `json:"template_id"`
	SessionIDs    []string `json:"session_ids"`
	JSON          struct {
		Type          respjson.Field
		MemoryStoreID respjson.Field
		IdentityID    respjson.Field
		TemplateID    respjson.Field
		SessionIDs    respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

func (r DreamInput) RawJSON() string                  { return r.JSON.raw }
func (r *DreamInput) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type DreamOutput struct {
	Type          string   `json:"type"`
	MemoryStoreID string   `json:"memory_store_id"`
	FilesTouched  []string `json:"files_touched"`
	JSON          struct {
		Type          respjson.Field
		MemoryStoreID respjson.Field
		FilesTouched  respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

func (r DreamOutput) RawJSON() string                  { return r.JSON.raw }
func (r *DreamOutput) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type DreamModelConfig struct {
	ID            string `json:"id"`
	Effort        string `json:"effort"`
	ContextWindow int64  `json:"context_window"`
	JSON          struct {
		ID            respjson.Field
		Effort        respjson.Field
		ContextWindow respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

func (r DreamModelConfig) RawJSON() string                  { return r.JSON.raw }
func (r *DreamModelConfig) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type DreamUsage struct {
	InputTokens              int64 `json:"input_tokens"`
	OutputTokens             int64 `json:"output_tokens"`
	CacheCreationInputTokens int64 `json:"cache_creation_input_tokens"`
	CacheReadInputTokens     int64 `json:"cache_read_input_tokens"`
	JSON                     struct {
		InputTokens              respjson.Field
		OutputTokens             respjson.Field
		CacheCreationInputTokens respjson.Field
		CacheReadInputTokens     respjson.Field
		ExtraFields              map[string]respjson.Field
		raw                      string
	} `json:"-"`
}

func (r DreamUsage) RawJSON() string                  { return r.JSON.raw }
func (r *DreamUsage) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type DreamError struct {
	Type    string `json:"type"`
	Code    string `json:"code"`
	Message string `json:"message"`
	JSON    struct {
		Type        respjson.Field
		Code        respjson.Field
		Message     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r DreamError) RawJSON() string                  { return r.JSON.raw }
func (r *DreamError) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
