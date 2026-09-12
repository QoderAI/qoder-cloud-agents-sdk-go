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

// ScheduleService provides Forward Schedule operations.
type ScheduleService struct {
	Options []option.RequestOption
}

func NewScheduleService(opts ...option.RequestOption) ScheduleService {
	return ScheduleService{Options: slices.Clone(opts)}
}

// List Schedules
func (r *ScheduleService) List(ctx context.Context, params ScheduleListParams, opts ...option.RequestOption) (res *pagination.Page[Schedule], err error) {

	opts = slices.Concat(r.Options, opts)
	path := "schedules"
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
func (r *ScheduleService) ListAutoPaging(ctx context.Context, params ScheduleListParams, opts ...option.RequestOption) *pagination.PageAutoPager[Schedule] {
	return pagination.NewPageAutoPager(r.List(ctx, params, opts...))
}

type ScheduleListParams struct {
	// Optional for a PAT or an admin SAT; when omitted, all Identities of the
	// current owner are queried. An Identity-bound SAT binds to itself when this is
	// omitted, and returns 403 if another Identity is passed explicitly.
	IdentityID param.Opt[string] `query:"identity_id,omitzero" json:"-"`
	// Filter by Forward Template ID.
	TemplateID param.Opt[string] `query:"template_id,omitzero" json:"-"`
	// Filter by `active` or `paused`.
	Status param.Opt[string] `query:"status,omitzero" json:"-"`
	// Whether to include archived Schedules.
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// Page size, maximum 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Cursor for the next page.
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// Cursor for the previous page.
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	// Sort field: `created_at` or `upcoming_runs_at`.
	SortBy param.Opt[string] `query:"sort_by,omitzero" json:"-"`
	// Sort direction: `asc` or `desc`.
	Order param.Opt[string] `query:"order,omitzero" json:"-"`
	paramObj
}

func (r ScheduleListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// Create Schedule
func (r *ScheduleService) New(ctx context.Context, params ScheduleNewParams, opts ...option.RequestOption) (res *Schedule, err error) {
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := "schedules"
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type ScheduleNewParams struct {
	// Forward Identity ID the Schedule belongs to.
	IdentityID string `json:"identity_id" api:"required"`
	// Forward Template ID to execute.
	TemplateID string `json:"template_id" api:"required"`
	// Schedule name.
	Name string `json:"name" api:"required"`
	// Schedule description.
	Description param.Opt[string] `json:"description,omitzero"`
	// Initial events injected on every run; currently `user.message` is supported.
	InitialEvents []map[string]any `json:"initial_events" api:"required"`
	// Execution policy; server-side defaults apply when omitted.
	Execution map[string]any `json:"execution,omitzero"`
	// Trigger policy; treated as `manual` when omitted or `null`.
	TriggerPolicy map[string]any `json:"trigger_policy,omitzero" api:"nullable"`
	// Environment used for execution.
	EnvironmentID string `json:"environment_id" api:"required"`
	// Destinations the run result is pushed to; kept as an array for compatibility,
	// currently at most one element is allowed.
	Sinks []map[string]any `json:"sinks,omitzero" api:"nullable"`
	// Business metadata, used only for labeling or pass-through.
	Metadata map[string]any `json:"metadata,omitzero"`
	// Optional idempotency key for requests with side effects.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r ScheduleNewParams) MarshalJSON() ([]byte, error) {
	type shadow ScheduleNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ScheduleNewParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

// Archive Schedules in bulk
func (r *ScheduleService) ArchiveMany(ctx context.Context, params ScheduleArchiveManyParams, opts ...option.RequestOption) (res *ScheduleArchiveManyResponse, err error) {
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := "schedules/archive"
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type ScheduleArchiveManyParams struct {
	// Archive scope; currently only by_schedule_ids is supported.
	Scope string `json:"scope" api:"required"`
	// Must contain 1-50 non-empty Schedule IDs after deduplication.
	ScheduleIDs []string `json:"schedule_ids" api:"required"`
	// Optional idempotency key for requests with side effects; the same owner, path
	// and body can be replayed safely.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r ScheduleArchiveManyParams) MarshalJSON() ([]byte, error) {
	type shadow ScheduleArchiveManyParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ScheduleArchiveManyParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Get Schedule
func (r *ScheduleService) Get(ctx context.Context, scheduleID string, opts ...option.RequestOption) (res *Schedule, err error) {
	if scheduleID == "" {
		return nil, fmt.Errorf("missing required schedule_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("schedules/%s", url.PathEscape(scheduleID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update Schedule
func (r *ScheduleService) Update(ctx context.Context, scheduleID string, params ScheduleUpdateParams, opts ...option.RequestOption) (res *Schedule, err error) {
	if scheduleID == "" {
		return nil, fmt.Errorf("missing required schedule_id parameter")
	}
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("schedules/%s", url.PathEscape(scheduleID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type ScheduleUpdateParams struct {
	// New Schedule name.
	Name param.Opt[string] `json:"name,omitzero"`
	// New Schedule description.
	Description param.Opt[string] `json:"description,omitzero"`
	// New Forward Template ID.
	TemplateID param.Opt[string] `json:"template_id,omitzero"`
	// Replaces the initial event list.
	InitialEvents []map[string]any `json:"initial_events,omitzero"`
	// Merges updates into the execution policy.
	Execution map[string]any `json:"execution,omitzero"`
	// Updates the trigger policy; `null` switches it to manual.
	TriggerPolicy map[string]any `json:"trigger_policy,omitzero" api:"nullable"`
	// New Environment used for execution.
	EnvironmentID param.Opt[string] `json:"environment_id,omitzero"`
	// Destinations the run result is pushed to; kept as an array for compatibility,
	// currently at most one element is allowed.
	Sinks []map[string]any `json:"sinks,omitzero" api:"nullable"`
	// Merges updates into metadata; a `null` value deletes the key.
	Metadata map[string]any `json:"metadata,omitzero"`
	// Optional idempotency key for requests with side effects.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r ScheduleUpdateParams) MarshalJSON() ([]byte, error) {
	type shadow ScheduleUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ScheduleUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Archive Schedule
func (r *ScheduleService) Archive(ctx context.Context, scheduleID string, params ScheduleArchiveParams, opts ...option.RequestOption) (res *Schedule, err error) {
	if scheduleID == "" {
		return nil, fmt.Errorf("missing required schedule_id parameter")
	}
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("schedules/%s/archive", url.PathEscape(scheduleID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type ScheduleArchiveParams struct {
	// Optional idempotency key for requests with side effects.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r ScheduleArchiveParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// Pause Schedule
func (r *ScheduleService) Pause(ctx context.Context, scheduleID string, params SchedulePauseParams, opts ...option.RequestOption) (res *Schedule, err error) {
	if scheduleID == "" {
		return nil, fmt.Errorf("missing required schedule_id parameter")
	}
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("schedules/%s/pause", url.PathEscape(scheduleID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type SchedulePauseParams struct {
	// Optional idempotency key for requests with side effects.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r SchedulePauseParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// Run Schedule
func (r *ScheduleService) Run(ctx context.Context, scheduleID string, params ScheduleRunParams, opts ...option.RequestOption) (res *ScheduleRun, err error) {
	if scheduleID == "" {
		return nil, fmt.Errorf("missing required schedule_id parameter")
	}
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("schedules/%s/run", url.PathEscape(scheduleID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type ScheduleRunParams struct {
	// Optional idempotency key for requests with side effects.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r ScheduleRunParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// Unpause Schedule
func (r *ScheduleService) Unpause(ctx context.Context, scheduleID string, params ScheduleUnpauseParams, opts ...option.RequestOption) (res *Schedule, err error) {
	if scheduleID == "" {
		return nil, fmt.Errorf("missing required schedule_id parameter")
	}
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("schedules/%s/unpause", url.PathEscape(scheduleID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type ScheduleUnpauseParams struct {
	// Optional idempotency key for requests with side effects.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r ScheduleUnpauseParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

type ScheduleArchiveManyResponse struct {
	// Number of Schedules that moved from unarchived to archived in this call;
	// targets that were already archived are not counted again.
	ArchivedCount int64 `json:"archived_count"`
	JSON          struct {
		ArchivedCount respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

func (r ScheduleArchiveManyResponse) RawJSON() string { return r.JSON.raw }
func (r *ScheduleArchiveManyResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type Schedule struct {
	ID            string                      `json:"id"`
	IdentityID    string                      `json:"identity_id"`
	TemplateID    string                      `json:"template_id"`
	Name          string                      `json:"name"`
	Description   string                      `json:"description"`
	Status        string                      `json:"status"`
	InitialEvents []ScheduleInitialEventsItem `json:"initial_events"`
	Execution     ScheduleExecution           `json:"execution"`
	TriggerPolicy ScheduleTriggerPolicy       `json:"trigger_policy"`
	EnvironmentID string                      `json:"environment_id"`
	Sinks         []ScheduleSinksItem         `json:"sinks"`
	Metadata      map[string]any              `json:"metadata"`
	CreatedAt     time.Time                   `json:"created_at" format:"date-time"`
	UpdatedAt     time.Time                   `json:"updated_at" format:"date-time"`
	ArchivedAt    time.Time                   `json:"archived_at" format:"date-time"`
	PausedReason  SchedulePausedReason        `json:"paused_reason"`
	JSON          struct {
		ID            respjson.Field
		IdentityID    respjson.Field
		TemplateID    respjson.Field
		Name          respjson.Field
		Description   respjson.Field
		Status        respjson.Field
		InitialEvents respjson.Field
		Execution     respjson.Field
		TriggerPolicy respjson.Field
		EnvironmentID respjson.Field
		Sinks         respjson.Field
		Metadata      respjson.Field
		CreatedAt     respjson.Field
		UpdatedAt     respjson.Field
		ArchivedAt    respjson.Field
		PausedReason  respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

func (r Schedule) RawJSON() string                  { return r.JSON.raw }
func (r *Schedule) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type SchedulePausedReason struct {
	Type string `json:"type"`
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r SchedulePausedReason) RawJSON() string { return r.JSON.raw }
func (r *SchedulePausedReason) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ScheduleSinksItem struct {
	Type      string                  `json:"type"`
	ChannelID string                  `json:"channel_id"`
	Target    ScheduleSinksItemTarget `json:"target"`
	JSON      struct {
		Type        respjson.Field
		ChannelID   respjson.Field
		Target      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r ScheduleSinksItem) RawJSON() string                  { return r.JSON.raw }
func (r *ScheduleSinksItem) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type ScheduleSinksItemTarget struct {
	Type       string `json:"type"`
	ExternalID string `json:"external_id"`
	JSON       struct {
		Type        respjson.Field
		ExternalID  respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r ScheduleSinksItemTarget) RawJSON() string { return r.JSON.raw }
func (r *ScheduleSinksItemTarget) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ScheduleTriggerPolicy struct {
	Type           string      `json:"type"`
	Expression     string      `json:"expression"`
	Timezone       string      `json:"timezone"`
	UpcomingRunsAt []time.Time `json:"upcoming_runs_at"`
	JSON           struct {
		Type           respjson.Field
		Expression     respjson.Field
		Timezone       respjson.Field
		UpcomingRunsAt respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

func (r ScheduleTriggerPolicy) RawJSON() string { return r.JSON.raw }
func (r *ScheduleTriggerPolicy) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ScheduleExecution struct {
	SessionMode       string `json:"session_mode"`
	MaxConcurrentRuns int64  `json:"max_concurrent_runs"`
	MaxAttempts       int64  `json:"max_attempts"`
	TimeoutMs         int64  `json:"timeout_ms"`
	JSON              struct {
		SessionMode       respjson.Field
		MaxConcurrentRuns respjson.Field
		MaxAttempts       respjson.Field
		TimeoutMs         respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

func (r ScheduleExecution) RawJSON() string                  { return r.JSON.raw }
func (r *ScheduleExecution) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type ScheduleInitialEventsItem struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	JSON    struct {
		Type        respjson.Field
		Content     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r ScheduleInitialEventsItem) RawJSON() string { return r.JSON.raw }
func (r *ScheduleInitialEventsItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
