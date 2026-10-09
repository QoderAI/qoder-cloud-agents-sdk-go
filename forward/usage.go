package forward

import (
	"context"
	"net/http"
	"net/url"
	"slices"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apiquery"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/pagination"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
)

type UsageService struct{ Options []option.RequestOption }

func NewUsageService(opts ...option.RequestOption) UsageService {
	return UsageService{Options: slices.Clone(opts)}
}

// ListIdentities aggregates an hourly Asia/Shanghai window. PAT or Admin SAT required.
func (r *UsageService) ListIdentities(ctx context.Context, params UsageListParams, opts ...option.RequestOption) (res *pagination.Page[IdentityUsage], err error) {
	opts = slices.Concat(r.Options, opts)
	path := "usage/identities"
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
func (r *UsageService) ListIdentitiesAutoPaging(ctx context.Context, params UsageListParams, opts ...option.RequestOption) *pagination.PageAutoPager[IdentityUsage] {
	return pagination.NewPageAutoPager(r.ListIdentities(ctx, params, opts...))
}

// ListTemplates aggregates an hourly Asia/Shanghai window. PAT or Admin SAT required.
func (r *UsageService) ListTemplates(ctx context.Context, params UsageListParams, opts ...option.RequestOption) (res *pagination.Page[TemplateUsage], err error) {
	opts = slices.Concat(r.Options, opts)
	path := "usage/templates"
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
func (r *UsageService) ListTemplatesAutoPaging(ctx context.Context, params UsageListParams, opts ...option.RequestOption) *pagination.PageAutoPager[TemplateUsage] {
	return pagination.NewPageAutoPager(r.ListTemplates(ctx, params, opts...))
}

// UsageListParams accepts only hourly parameters. StartAt is inclusive; EndAt is exclusive, with a maximum span of 744 hours. Both use YYYY-MM-DDTHH:00:00 in Asia/Shanghai for CN and Global.
type UsageListParams struct {
	StartAt    string            `query:"start_at" json:"-" api:"required"`
	EndAt      string            `query:"end_at" json:"-" api:"required"`
	Limit      param.Opt[int64]  `query:"limit,omitzero" json:"-"`
	AfterID    param.Opt[string] `query:"after_id,omitzero" json:"-"`
	BeforeID   param.Opt[string] `query:"before_id,omitzero" json:"-"`
	IdentityID param.Opt[string] `query:"identity_id,omitzero" json:"-"`
	// Repeated query values; an entry may also contain comma-separated IDs.
	IdentityIDs []string          `query:"identity_ids,omitzero" json:"-"`
	TemplateID  param.Opt[string] `query:"template_id,omitzero" json:"-"`
	// Repeated query values; an entry may also contain comma-separated IDs.
	TemplateIDs []string `query:"template_ids,omitzero" json:"-"`
	paramObj
}

func (r *UsageListParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
func (r UsageListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

type IdentityUsage struct {
	Type          string  `json:"type"`
	IdentityID    string  `json:"identity_id"`
	SessionCount  int64   `json:"session_count"`
	ActiveSeconds float64 `json:"active_seconds"`
	Credits       float64 `json:"credits"`
	JSON          struct {
		Type          respjson.Field
		IdentityID    respjson.Field
		SessionCount  respjson.Field
		ActiveSeconds respjson.Field
		Credits       respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

func (r IdentityUsage) RawJSON() string                  { return r.JSON.raw }
func (r *IdentityUsage) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type TemplateUsage struct {
	Type             string  `json:"type"`
	TemplateID       string  `json:"template_id"`
	ActiveIdentities int64   `json:"active_identities"`
	SessionCount     int64   `json:"session_count"`
	ActiveSeconds    float64 `json:"active_seconds"`
	Credits          float64 `json:"credits"`
	JSON             struct {
		Type             respjson.Field
		TemplateID       respjson.Field
		ActiveIdentities respjson.Field
		SessionCount     respjson.Field
		ActiveSeconds    respjson.Field
		Credits          respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

func (r TemplateUsage) RawJSON() string                  { return r.JSON.raw }
func (r *TemplateUsage) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
