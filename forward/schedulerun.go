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

// ScheduleRunService provides Forward ScheduleRun operations.
type ScheduleRunService struct {
	Options []option.RequestOption
}

func NewScheduleRunService(opts ...option.RequestOption) ScheduleRunService {
	return ScheduleRunService{Options: slices.Clone(opts)}
}

// List Schedule Runs
func (r *ScheduleRunService) List(ctx context.Context, params ScheduleRunListParams, opts ...option.RequestOption) (res *pagination.Page[ScheduleRun], err error) {

	opts = slices.Concat(r.Options, opts)
	path := "schedule_runs"
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
func (r *ScheduleRunService) ListAutoPaging(ctx context.Context, params ScheduleRunListParams, opts ...option.RequestOption) *pagination.PageAutoPager[ScheduleRun] {
	return pagination.NewPageAutoPager(r.List(ctx, params, opts...))
}

type ScheduleRunListParams struct {
	// Forward Identity ID that owns the Run.
	IdentityID string `query:"identity_id,omitzero" json:"-" api:"required"`
	// Filter by Schedule ID.
	ScheduleID param.Opt[string] `query:"schedule_id,omitzero" json:"-"`
	// Filter by `pending`, `running`, `completed`, `failed` or `skipped`.
	Status param.Opt[string] `query:"status,omitzero" json:"-"`
	// Filter by `schedule` or `manual`.
	TriggerType param.Opt[string] `query:"trigger_type,omitzero" json:"-"`
	// Whether to return only Runs with or without an error.
	HasError param.Opt[bool] `query:"has_error,omitzero" json:"-"`
	// Page size, up to 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Cursor for the next page.
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// Cursor for the previous page.
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	// Sort field: `created_at` or `triggered_at`.
	SortBy param.Opt[string] `query:"sort_by,omitzero" json:"-"`
	// Sort direction: `asc` or `desc`.
	Order param.Opt[string] `query:"order,omitzero" json:"-"`
	paramObj
}

func (r ScheduleRunListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// Get Schedule Run
func (r *ScheduleRunService) Get(ctx context.Context, runID string, params ScheduleRunGetParams, opts ...option.RequestOption) (res *ScheduleRun, err error) {
	if runID == "" {
		return nil, fmt.Errorf("missing required run_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("schedule_runs/%s", url.PathEscape(runID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

type ScheduleRunGetParams struct {
	// Additional ownership constraint.
	IdentityID param.Opt[string] `query:"identity_id,omitzero" json:"-"`
	paramObj
}

func (r ScheduleRunGetParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

type ScheduleRun struct {
	// Schedule Run ID.
	ID string `json:"id"`
	// ID of the owning Schedule.
	ScheduleID string `json:"schedule_id"`
	// Forward Identity ID.
	IdentityID string `json:"identity_id"`
	// Forward Template ID.
	TemplateID string `json:"template_id"`
	// null|Session created or reused by this run.
	SessionID string `json:"session_id"`
	// `pending`, `running`, `completed`, `failed` or `skipped`.
	Status string `json:"status"`
	// Trigger source.
	TriggerContext ScheduleRunTriggerContext `json:"trigger_context"`
	// null|Text result of the main flow.
	ResultPayload string `json:"result_payload"`
	// null|Sink type used for this IM delivery; `null` when no delivery is configured.
	PushSink string `json:"push_sink"`
	// IM delivery status: `pending`, `succeeded`, `failed` or `skipped`. The main flow
	// status and the delivery status are independent.
	PushStatus string `json:"push_status"`
	// null|Time the IM delivery finished.
	PushFinishedAt time.Time `json:"push_finished_at" format:"date-time"`
	// Current or final attempt number, starting at `1`; may return `2` when the Schedule sets
	// `execution.max_attempts=2` and the server has performed an automatic retry.
	Attempt int64 `json:"attempt"`
	// Trigger time.
	TriggeredAt time.Time `json:"triggered_at" format:"date-time"`
	// null|Time execution started.
	StartedAt time.Time `json:"started_at" format:"date-time"`
	// null|Completion time.
	CompletedAt time.Time `json:"completed_at" format:"date-time"`
	// null|Execution duration in milliseconds.
	DurationMs int64 `json:"duration_ms"`
	// Record creation time.
	CreatedAt time.Time `json:"created_at" format:"date-time"`
	// null|Structured error when the Run failed or was skipped.
	Error map[string]any `json:"error"`
	// null|Human-readable error message; structured details stay in `error`.
	ErrorMessage string `json:"error_message"`
	JSON         struct {
		ID             respjson.Field
		ScheduleID     respjson.Field
		IdentityID     respjson.Field
		TemplateID     respjson.Field
		SessionID      respjson.Field
		Status         respjson.Field
		TriggerContext respjson.Field
		ResultPayload  respjson.Field
		PushSink       respjson.Field
		PushStatus     respjson.Field
		PushFinishedAt respjson.Field
		Attempt        respjson.Field
		TriggeredAt    respjson.Field
		StartedAt      respjson.Field
		CompletedAt    respjson.Field
		DurationMs     respjson.Field
		CreatedAt      respjson.Field
		Error          respjson.Field
		ErrorMessage   respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

func (r ScheduleRun) RawJSON() string                  { return r.JSON.raw }
func (r *ScheduleRun) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type ScheduleRunTriggerContext struct {
	Type        string    `json:"type"`
	ScheduledAt time.Time `json:"scheduled_at" format:"date-time"`
	JSON        struct {
		Type        respjson.Field
		ScheduledAt respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r ScheduleRunTriggerContext) RawJSON() string { return r.JSON.raw }
func (r *ScheduleRunTriggerContext) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
