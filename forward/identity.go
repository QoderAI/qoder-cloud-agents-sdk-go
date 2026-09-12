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

// IdentityService provides Forward Identity operations.
type IdentityService struct {
	Options      []option.RequestOption
	Configs      IdentityConfigService
	MemoryStores IdentityMemoryStoreService
}

func NewIdentityService(opts ...option.RequestOption) IdentityService {
	return IdentityService{Options: slices.Clone(opts), Configs: NewIdentityConfigService(opts...), MemoryStores: NewIdentityMemoryStoreService(opts...)}
}

// List Identities
func (r *IdentityService) List(ctx context.Context, params IdentityListParams, opts ...option.RequestOption) (res *pagination.Page[Identity], err error) {

	opts = slices.Concat(r.Options, opts)
	path := "identities"
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
func (r *IdentityService) ListAutoPaging(ctx context.Context, params IdentityListParams, opts ...option.RequestOption) *pagination.PageAutoPager[Identity] {
	return pagination.NewPageAutoPager(r.List(ctx, params, opts...))
}

type IdentityListParams struct {
	// Filter by the integrator's end-user ID.
	ExternalID param.Opt[string] `query:"external_id,omitzero" json:"-"`
	// Filter by multiple Identity IDs; accepts a comma-separated list or repeated
	// query parameters, up to 100 after deduplication.
	IdentityIDs []string `query:"identity_ids,omitzero" json:"-"`
	// Matches the Identity ID, name or external ID.
	Search param.Opt[string] `query:"search,omitzero" json:"-"`
	// Filter by enabled state; a non-boolean value returns 400.
	Enabled param.Opt[bool] `query:"enabled,omitzero" json:"-"`
	// Page size, maximum 100; larger values are clamped to the maximum.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Cursor for the next page; cannot be combined with `before_id`.
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// Cursor for the previous page; cannot be combined with `after_id`.
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	paramObj
}

func (r IdentityListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// Create Identity
func (r *IdentityService) New(ctx context.Context, params IdentityNewParams, opts ...option.RequestOption) (res *Identity, err error) {
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := "identities"
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type IdentityNewParams struct {
	// End-user ID in the integrator's system; cannot be empty or whitespace only.
	ExternalID string `json:"external_id" api:"required"`
	// Display name; when provided it cannot be empty or whitespace only.
	Name param.Opt[string] `json:"name,omitzero"`
	// Whether the Identity is enabled; defaults to `true`.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Business metadata; at most 16 keys is recommended.
	Metadata map[string]any `json:"metadata,omitzero"`
	// Optional idempotency key for requests with side effects.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r IdentityNewParams) MarshalJSON() ([]byte, error) {
	type shadow IdentityNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *IdentityNewParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

// Ensure the admin Identity
func (r *IdentityService) EnsureAdmin(ctx context.Context, opts ...option.RequestOption) (res *Identity, err error) {

	opts = slices.Concat(r.Options, opts)
	path := "identities/admin/ensure"
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Get Identity stats
func (r *IdentityService) Stats(ctx context.Context, opts ...option.RequestOption) (res *IdentityStats, err error) {

	opts = slices.Concat(r.Options, opts)
	path := "identities/stats"
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Get Identity
func (r *IdentityService) Get(ctx context.Context, identityID string, opts ...option.RequestOption) (res *Identity, err error) {
	if identityID == "" {
		return nil, fmt.Errorf("missing required identity_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("identities/%s", url.PathEscape(identityID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update Identity
func (r *IdentityService) Update(ctx context.Context, identityID string, params IdentityUpdateParams, opts ...option.RequestOption) (res *Identity, err error) {
	if identityID == "" {
		return nil, fmt.Errorf("missing required identity_id parameter")
	}
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("identities/%s", url.PathEscape(identityID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type IdentityUpdateParams struct {
	// Replaces the existing end-user ID.
	ExternalID param.Opt[string] `json:"external_id,omitzero"`
	// Replaces the display name.
	Name param.Opt[string] `json:"name,omitzero"`
	// Updates whether the Identity is usable.
	Enabled param.Opt[bool] `json:"enabled,omitzero"`
	// Merges updates into the business metadata; an empty string value deletes the
	// corresponding key.
	Metadata map[string]any `json:"metadata,omitzero"`
	// Optional idempotency key for requests with side effects.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r IdentityUpdateParams) MarshalJSON() ([]byte, error) {
	type shadow IdentityUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *IdentityUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Delete Identity
func (r *IdentityService) Delete(ctx context.Context, identityID string, opts ...option.RequestOption) (res *DeletedIdentity, err error) {
	if identityID == "" {
		return nil, fmt.Errorf("missing required identity_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("identities/%s", url.PathEscape(identityID))
	err = convention.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// List Templates used by an Identity
func (r *IdentityService) ListTemplates(ctx context.Context, identityID string, opts ...option.RequestOption) (res *IdentityListTemplatesResponse, err error) {
	if identityID == "" {
		return nil, fmt.Errorf("missing required identity_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("identities/%s/agents", url.PathEscape(identityID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Clear Identity
func (r *IdentityService) Clear(ctx context.Context, identityID string, params IdentityClearParams, opts ...option.RequestOption) (res *IdentityClearResponse, err error) {
	if identityID == "" {
		return nil, fmt.Errorf("missing required identity_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("identities/%s/clear", url.PathEscape(identityID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type IdentityClearParams struct {
	// Reason for clearing; recorded only to capture the caller's intent.
	Reason param.Opt[string] `json:"reason,omitzero"`
	paramObj
}

func (r IdentityClearParams) MarshalJSON() ([]byte, error) {
	type shadow IdentityClearParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *IdentityClearParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

// Disable Identity
func (r *IdentityService) Disable(ctx context.Context, identityID string, opts ...option.RequestOption) (res *Identity, err error) {
	if identityID == "" {
		return nil, fmt.Errorf("missing required identity_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("identities/%s/disable", url.PathEscape(identityID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Enable Identity
func (r *IdentityService) Enable(ctx context.Context, identityID string, opts ...option.RequestOption) (res *Identity, err error) {
	if identityID == "" {
		return nil, fmt.Errorf("missing required identity_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("identities/%s/enable", url.PathEscape(identityID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type IdentityClearResponse struct {
	// ID of the cleared Identity.
	IdentityID string `json:"identity_id"`
	// Clear status; `completed` on success.
	Status string `json:"status"`
	// Time the Forward-side cleanup completed, in RFC 3339 format.
	CompletedAt time.Time                    `json:"completed_at" format:"date-time"`
	Summary     IdentityClearResponseSummary `json:"summary"`
	JSON        struct {
		IdentityID  respjson.Field
		Status      respjson.Field
		CompletedAt respjson.Field
		Summary     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r IdentityClearResponse) RawJSON() string { return r.JSON.raw }
func (r *IdentityClearResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type IdentityClearResponseSummary struct {
	IdentityConfigsArchived        int64 `json:"identity_configs_archived"`
	ResourceBindingsArchived       int64 `json:"resource_bindings_archived"`
	IdentityOwnedResourcesArchived int64 `json:"identity_owned_resources_archived"`
	SchedulesArchived              int64 `json:"schedules_archived"`
	ScheduleRunsSkipped            int64 `json:"schedule_runs_skipped"`
	SessionsArchived               int64 `json:"sessions_archived"`
	JSON                           struct {
		IdentityConfigsArchived        respjson.Field
		ResourceBindingsArchived       respjson.Field
		IdentityOwnedResourcesArchived respjson.Field
		SchedulesArchived              respjson.Field
		ScheduleRunsSkipped            respjson.Field
		SessionsArchived               respjson.Field
		ExtraFields                    map[string]respjson.Field
		raw                            string
	} `json:"-"`
}

func (r IdentityClearResponseSummary) RawJSON() string { return r.JSON.raw }
func (r *IdentityClearResponseSummary) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type IdentityTemplate struct {
	TemplateID   string    `json:"template_id"`
	TemplateName string    `json:"template_name"`
	SessionCount int64     `json:"session_count"`
	LastActiveAt time.Time `json:"last_active_at" format:"date-time"`
	JSON         struct {
		TemplateID   respjson.Field
		TemplateName respjson.Field
		SessionCount respjson.Field
		LastActiveAt respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

func (r IdentityTemplate) RawJSON() string                  { return r.JSON.raw }
func (r *IdentityTemplate) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type DeletedIdentity struct {
	// ID of the deleted Identity.
	ID string `json:"id"`
	// Whether the deletion completed. `true` in a successful response.
	Deleted bool `json:"deleted"`
	JSON    struct {
		ID          respjson.Field
		Deleted     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r DeletedIdentity) RawJSON() string                  { return r.JSON.raw }
func (r *DeletedIdentity) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type IdentityStats struct {
	// Total number of Identities for the current account.
	TotalIdentities int64 `json:"total_identities"`
	// Number of recently active Identities.
	ActiveIdentities int64 `json:"active_identities"`
	// Total number of Templates the current account has used.
	TotalAgents int64 `json:"total_agents"`
	// Total number of Sessions for the current account.
	TotalSessions int64 `json:"total_sessions"`
	JSON          struct {
		TotalIdentities  respjson.Field
		ActiveIdentities respjson.Field
		TotalAgents      respjson.Field
		TotalSessions    respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

func (r IdentityStats) RawJSON() string                  { return r.JSON.raw }
func (r *IdentityStats) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type Identity struct {
	// Forward Identity ID; the `idn_` prefix is recommended.
	ID string `json:"id"`
	// End-user ID in the integrator's system.
	ExternalID string `json:"external_id"`
	// Identity display name.
	Name string `json:"name"`
	// Identity type. A regular Identity returns `normal`.
	IdentityType string `json:"identity_type"`
	// Whether the Identity may continue to be used.
	Enabled bool `json:"enabled"`
	// Business metadata.
	Metadata map[string]any `json:"metadata"`
	// Creation time, in RFC 3339 format.
	CreatedAt time.Time `json:"created_at" format:"date-time"`
	// Last update time, in RFC 3339 format.
	UpdatedAt time.Time `json:"updated_at" format:"date-time"`
	JSON      struct {
		ID           respjson.Field
		ExternalID   respjson.Field
		Name         respjson.Field
		IdentityType respjson.Field
		Enabled      respjson.Field
		Metadata     respjson.Field
		CreatedAt    respjson.Field
		UpdatedAt    respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

func (r Identity) RawJSON() string                  { return r.JSON.raw }
func (r *Identity) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type IdentityListTemplatesResponse struct {
	Data []IdentityTemplate `json:"data"`
	JSON struct {
		Data        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r IdentityListTemplatesResponse) RawJSON() string { return r.JSON.raw }
func (r *IdentityListTemplatesResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
