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

// SessionService provides Forward Session operations.
type SessionService struct {
	Options   []option.RequestOption
	Events    SessionEventService
	Resources SessionResourceService
	Threads   SessionThreadService
}

func NewSessionService(opts ...option.RequestOption) SessionService {
	return SessionService{Options: slices.Clone(opts), Events: NewSessionEventService(opts...), Resources: NewSessionResourceService(opts...), Threads: NewSessionThreadService(opts...)}
}

// List Sessions
func (r *SessionService) List(ctx context.Context, params SessionListParams, opts ...option.RequestOption) (res *pagination.Page[Session], err error) {

	opts = slices.Concat(r.Options, opts)
	path := "sessions"
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
func (r *SessionService) ListAutoPaging(ctx context.Context, params SessionListParams, opts ...option.RequestOption) *pagination.PageAutoPager[Session] {
	return pagination.NewPageAutoPager(r.List(ctx, params, opts...))
}

type SessionListParams struct {
	// Filter by one or more Identity IDs, comma-separated.
	IdentityIDs []string `query:"identity_ids,omitzero" json:"-"`
	// Filter by Forward Template ID.
	TemplateID param.Opt[string] `query:"template_id,omitzero" json:"-"`
	// Filter by `api`, `im`, `schedule` or `batch`.
	SourceType param.Opt[string] `query:"source_type,omitzero" json:"-"`
	// Creation time strictly after this RFC 3339 timestamp.
	CreatedAtGt param.Opt[time.Time] `query:"created_at[gt],omitzero" json:"-" format:"date-time"`
	// Creation time at or after this RFC 3339 timestamp.
	CreatedAtGte param.Opt[time.Time] `query:"created_at[gte],omitzero" json:"-" format:"date-time"`
	// Creation time strictly before this RFC 3339 timestamp.
	CreatedAtLt param.Opt[time.Time] `query:"created_at[lt],omitzero" json:"-" format:"date-time"`
	// Creation time at or before this RFC 3339 timestamp.
	CreatedAtLte param.Opt[time.Time] `query:"created_at[lte],omitzero" json:"-" format:"date-time"`
	// Update time strictly after this RFC 3339 timestamp.
	UpdatedAtGt param.Opt[time.Time] `query:"updated_at[gt],omitzero" json:"-" format:"date-time"`
	// Update time at or after this RFC 3339 timestamp.
	UpdatedAtGte param.Opt[time.Time] `query:"updated_at[gte],omitzero" json:"-" format:"date-time"`
	// Update time strictly before this RFC 3339 timestamp.
	UpdatedAtLt param.Opt[time.Time] `query:"updated_at[lt],omitzero" json:"-" format:"date-time"`
	// Update time at or before this RFC 3339 timestamp.
	UpdatedAtLte param.Opt[time.Time] `query:"updated_at[lte],omitzero" json:"-" format:"date-time"`
	// Page size, up to 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Cursor for the next page; pass the `last_id` from the previous response.
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// Cursor for the previous page; pass the `first_id` from the current response.
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	// Sort direction by creation time: `desc` or `asc`.
	Order param.Opt[string] `query:"order,omitzero" json:"-"`
	// Whether to include archived Sessions.
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	paramObj
}

func (r SessionListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// Create Session
func (r *SessionService) New(ctx context.Context, params SessionNewParams, opts ...option.RequestOption) (res *Session, err error) {
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := "sessions"
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type SessionNewParams struct {
	// Forward Identity ID.
	IdentityID string `json:"identity_id" api:"required"`
	// Forward Template ID.
	TemplateID string `json:"template_id" api:"required"`
	// Session title.
	Title param.Opt[string] `json:"title,omitzero"`
	// Business metadata.
	Metadata  map[string]any              `json:"metadata,omitzero"`
	Config    SessionNewParamsConfigParam `json:"config,omitzero"`
	Resources []SessionResourceSpecParam  `json:"resources,omitzero"`
	// Optional idempotency key for requests with side effects.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r SessionNewParams) MarshalJSON() ([]byte, error) {
	type shadow SessionNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionNewParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type SessionNewParamsConfigParam struct {
	EnvironmentVariables map[string]any `json:"environment_variables,omitzero"`
	paramObj
}

func (r SessionNewParamsConfigParam) MarshalJSON() ([]byte, error) {
	type shadow SessionNewParamsConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionNewParamsConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Get Session
func (r *SessionService) Get(ctx context.Context, sessionID string, opts ...option.RequestOption) (res *Session, err error) {
	if sessionID == "" {
		return nil, fmt.Errorf("missing required session_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("sessions/%s", url.PathEscape(sessionID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update Session
func (r *SessionService) Update(ctx context.Context, sessionID string, params SessionUpdateParams, opts ...option.RequestOption) (res *Session, err error) {
	if sessionID == "" {
		return nil, fmt.Errorf("missing required session_id parameter")
	}
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("sessions/%s", url.PathEscape(sessionID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type SessionUpdateParams struct {
	// New Session title.
	Title param.Opt[string] `json:"title,omitzero"`
	// Metadata merge patch; supplied keys overwrite existing keys, omitted keys are kept.
	Metadata map[string]any                 `json:"metadata,omitzero"`
	Config   SessionUpdateParamsConfigParam `json:"config,omitzero"`
	// Optional idempotency key for requests with side effects.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r SessionUpdateParams) MarshalJSON() ([]byte, error) {
	type shadow SessionUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionUpdateParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type SessionUpdateParamsConfigParam struct {
	EnvironmentVariables map[string]any `json:"environment_variables,omitzero"`
	paramObj
}

func (r SessionUpdateParamsConfigParam) MarshalJSON() ([]byte, error) {
	type shadow SessionUpdateParamsConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionUpdateParamsConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Archive Session
func (r *SessionService) Archive(ctx context.Context, sessionID string, params SessionArchiveParams, opts ...option.RequestOption) (res *Session, err error) {
	if sessionID == "" {
		return nil, fmt.Errorf("missing required session_id parameter")
	}
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("sessions/%s/archive", url.PathEscape(sessionID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type SessionArchiveParams struct {
	// Optional idempotency key for requests with side effects.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r SessionArchiveParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// Cancel the current Turn
func (r *SessionService) Cancel(ctx context.Context, sessionID string, params SessionCancelParams, opts ...option.RequestOption) (res *Session, err error) {
	if sessionID == "" {
		return nil, fmt.Errorf("missing required session_id parameter")
	}
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("sessions/%s/cancel", url.PathEscape(sessionID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type SessionCancelParams struct {
	// Optional idempotency key for requests with side effects.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r SessionCancelParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

type Session struct {
	// Session ID.
	ID string `json:"id"`
	// Always `session`.
	Type string `json:"type"`
	// Forward Identity ID.
	IdentityID string `json:"identity_id"`
	// Template summary.
	Template SessionTemplate `json:"template"`
	// Session source; `api` when created directly through the API.
	SourceType string `json:"source_type"`
	// `idle`, `running`, `rescheduling`, `canceling` or `terminated`.
	Status     string            `json:"status"`
	Title      string            `json:"title"`
	Metadata   map[string]any    `json:"metadata"`
	Config     SessionConfig     `json:"config"`
	Resources  []SessionResource `json:"resources"`
	Stats      SessionStats      `json:"stats"`
	Usage      SessionUsage      `json:"usage"`
	ArchivedAt time.Time         `json:"archived_at" api:"nullable" format:"date-time"`
	CreatedAt  time.Time         `json:"created_at" format:"date-time"`
	UpdatedAt  time.Time         `json:"updated_at" format:"date-time"`
	JSON       struct {
		ID          respjson.Field
		Type        respjson.Field
		IdentityID  respjson.Field
		Template    respjson.Field
		SourceType  respjson.Field
		Status      respjson.Field
		Title       respjson.Field
		Metadata    respjson.Field
		Config      respjson.Field
		Resources   respjson.Field
		Stats       respjson.Field
		Usage       respjson.Field
		ArchivedAt  respjson.Field
		CreatedAt   respjson.Field
		UpdatedAt   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r Session) RawJSON() string                  { return r.JSON.raw }
func (r *Session) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type SessionUsage struct {
	TotalCredits float64 `json:"total_credits"`
	JSON         struct {
		TotalCredits respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

func (r SessionUsage) RawJSON() string                  { return r.JSON.raw }
func (r *SessionUsage) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type SessionStats struct {
	ActiveSeconds   int64 `json:"active_seconds"`
	DurationSeconds int64 `json:"duration_seconds"`
	JSON            struct {
		ActiveSeconds   respjson.Field
		DurationSeconds respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

func (r SessionStats) RawJSON() string                  { return r.JSON.raw }
func (r *SessionStats) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type SessionConfig struct {
	EnvironmentVariables map[string]string `json:"environment_variables"`
	JSON                 struct {
		EnvironmentVariables respjson.Field
		ExtraFields          map[string]respjson.Field
		raw                  string
	} `json:"-"`
}

func (r SessionConfig) RawJSON() string                  { return r.JSON.raw }
func (r *SessionConfig) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type SessionTemplate struct {
	ID      string      `json:"id"`
	Type    string      `json:"type"`
	Name    string      `json:"name"`
	Model   ModelConfig `json:"model"`
	Version int64       `json:"version"`
	JSON    struct {
		ID          respjson.Field
		Type        respjson.Field
		Name        respjson.Field
		Model       respjson.Field
		Version     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r SessionTemplate) RawJSON() string                  { return r.JSON.raw }
func (r *SessionTemplate) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
