// Qoder managed API definitions.
package managed

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	requestconfig "github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apiquery"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/constant"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/pagination"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
)

// EnvironmentWorkService contains methods and other services that help with
// interacting with the Qoder Cloud Agents API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewEnvironmentWorkService] method instead.
type EnvironmentWorkService struct {
	Options []option.RequestOption
}

// NewEnvironmentWorkService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewEnvironmentWorkService(opts ...option.RequestOption) (r EnvironmentWorkService) {
	r = EnvironmentWorkService{}
	r.Options = opts
	return
}

// Note: these endpoints are called automatically by the pre-built environment
// worker provided in the SDKs and CLI, for orchestrating sessions with self-hosted
// sandbox environments. They are included here as a reference; you do not need to
// invoke them directly.
//
// Retrieve detailed information about a specific work item.
func (r *EnvironmentWorkService) Get(ctx context.Context, workID string, params EnvironmentWorkGetParams, opts ...option.RequestOption) (res *SelfHostedWork, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if params.EnvironmentID == "" {
		err = errors.New("missing required environment_id parameter")
		return nil, err
	}
	if workID == "" {
		err = errors.New("missing required work_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("environments/%s/work/%s", url.PathEscape(params.EnvironmentID), url.PathEscape(workID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Note: these endpoints are called automatically by the pre-built environment
// worker provided in the SDKs and CLI, for orchestrating sessions with self-hosted
// sandbox environments. They are included here as a reference; you do not need to
// invoke them directly.
//
// Update work item metadata with merge semantics.
func (r *EnvironmentWorkService) Update(ctx context.Context, workID string, params EnvironmentWorkUpdateParams, opts ...option.RequestOption) (res *SelfHostedWork, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if params.EnvironmentID == "" {
		err = errors.New("missing required environment_id parameter")
		return nil, err
	}
	if workID == "" {
		err = errors.New("missing required work_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("environments/%s/work/%s", url.PathEscape(params.EnvironmentID), url.PathEscape(workID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Note: these endpoints are called automatically by the pre-built environment
// worker provided in the SDKs and CLI, for orchestrating sessions with self-hosted
// sandbox environments. They are included here as a reference; you do not need to
// invoke them directly.
//
// List work items in an environment.
func (r *EnvironmentWorkService) List(ctx context.Context, environmentID string, params EnvironmentWorkListParams, opts ...option.RequestOption) (res *pagination.PageCursor[SelfHostedWork], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if environmentID == "" {
		err = errors.New("missing required environment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("environments/%s/work", url.PathEscape(environmentID))
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, params, &res, opts...)
	if err != nil {
		return nil, err
	}
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

// Note: these endpoints are called automatically by the pre-built environment
// worker provided in the SDKs and CLI, for orchestrating sessions with self-hosted
// sandbox environments. They are included here as a reference; you do not need to
// invoke them directly.
//
// List work items in an environment.
func (r *EnvironmentWorkService) ListAutoPaging(ctx context.Context, environmentID string, params EnvironmentWorkListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[SelfHostedWork] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, environmentID, params, opts...))
}

// Note: these endpoints are called automatically by the pre-built environment
// worker provided in the SDKs and CLI, for orchestrating sessions with self-hosted
// sandbox environments. They are included here as a reference; you do not need to
// invoke them directly.
//
// Acknowledge receipt of a work item, transitioning it from 'queued' to 'starting'
// and removing it from the queue.
func (r *EnvironmentWorkService) Ack(ctx context.Context, workID string, params EnvironmentWorkAckParams, opts ...option.RequestOption) (res *SelfHostedWork, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if params.EnvironmentID == "" {
		err = errors.New("missing required environment_id parameter")
		return nil, err
	}
	if workID == "" {
		err = errors.New("missing required work_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("environments/%s/work/%s/ack", url.PathEscape(params.EnvironmentID), url.PathEscape(workID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Note: these endpoints are called automatically by the pre-built environment
// worker provided in the SDKs and CLI, for orchestrating sessions with self-hosted
// sandbox environments. They are included here as a reference; you do not need to
// invoke them directly.
//
// Record a heartbeat for a work item to maintain the lease.
func (r *EnvironmentWorkService) Heartbeat(ctx context.Context, workID string, params EnvironmentWorkHeartbeatParams, opts ...option.RequestOption) (res *SelfHostedWorkHeartbeatResponse, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if params.EnvironmentID == "" {
		err = errors.New("missing required environment_id parameter")
		return nil, err
	}
	if workID == "" {
		err = errors.New("missing required work_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("environments/%s/work/%s/heartbeat", url.PathEscape(params.EnvironmentID), url.PathEscape(workID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Note: these endpoints are called automatically by the pre-built environment
// worker provided in the SDKs and CLI, for orchestrating sessions with self-hosted
// sandbox environments. They are included here as a reference; you do not need to
// invoke them directly.
//
// Long poll for work items in the queue.
func (r *EnvironmentWorkService) Poll(ctx context.Context, environmentID string, params EnvironmentWorkPollParams, opts ...option.RequestOption) (res *SelfHostedWork, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	if !param.IsOmitted(params.QoderWorkerID) {
		opts = append(opts, option.WithHeader("Worker-ID", fmt.Sprintf("%v", params.QoderWorkerID.Value)))
	}
	opts = slices.Concat(r.Options, opts)

	if environmentID == "" {
		err = errors.New("missing required environment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("environments/%s/work/poll", url.PathEscape(environmentID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

// Get statistics about the work queue for an environment.
func (r *EnvironmentWorkService) Stats(ctx context.Context, environmentID string, query EnvironmentWorkStatsParams, opts ...option.RequestOption) (res *SelfHostedWorkQueueStats, err error) {
	for _, v := range query.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if environmentID == "" {
		err = errors.New("missing required environment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("environments/%s/work/stats", url.PathEscape(environmentID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Note: these endpoints are called automatically by the pre-built environment
// worker provided in the SDKs and CLI, for orchestrating sessions with self-hosted
// sandbox environments. They are included here as a reference; you do not need to
// invoke them directly.
//
// Stop a work item, initiating graceful or forced shutdown.
func (r *EnvironmentWorkService) Stop(ctx context.Context, workID string, params EnvironmentWorkStopParams, opts ...option.RequestOption) (res *SelfHostedWork, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if params.EnvironmentID == "" {
		err = errors.New("missing required environment_id parameter")
		return nil, err
	}
	if workID == "" {
		err = errors.New("missing required work_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("environments/%s/work/%s/stop", url.PathEscape(params.EnvironmentID), url.PathEscape(workID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Work data for environment health checks.
//
// This resource type is used for assessing the health of containers where work
// occurs. The data is opaque to users; the runner handles the health check by
// probing connectivity to required services.
type HealthCheckWorkData struct {
	// Health check identifier
	ID string `json:"id" api:"required"`
	// Type of work data
	//
	// Any of "healthcheck".
	Type HealthCheckWorkDataType `json:"type"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r HealthCheckWorkData) RawJSON() string { return r.JSON.raw }
func (r *HealthCheckWorkData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Type of work data
type HealthCheckWorkDataType string

const (
	HealthCheckWorkDataTypeHealthcheck HealthCheckWorkDataType = "healthcheck"
)

// Work resource representing a unit of work in a self-hosted environment.
//
// Work items are queued when sessions are created or when long-dormant sessions
// receive new messages. The Environment Manager polls for work items and executes
// them on customer-hosted infrastructure.
type SelfHostedWork struct {
	// Work identifier (e.g., 'work\_...')
	ID string `json:"id" api:"required"`
	// RFC 3339 timestamp when work was acknowledged by Environment Manager
	AcknowledgedAt string `json:"acknowledged_at" api:"required"`
	// RFC 3339 timestamp when work was created
	CreatedAt string `json:"created_at" api:"required"`
	// The actual work to be performed
	Data SessionWorkData `json:"data" api:"required"`
	// Environment identifier this work belongs to (e.g., `env_...`)
	EnvironmentID string `json:"environment_id" api:"required"`
	// RFC 3339 timestamp of the most recent heartbeat
	LatestHeartbeatAt string `json:"latest_heartbeat_at" api:"required"`
	// User-provided metadata key-value pairs associated with this work item
	Metadata map[string]string `json:"metadata" api:"required"`
	// Credential payload used by the environment worker to execute this work item. May
	// be populated when polling for work; null on all other retrieval paths.
	Secret string `json:"secret" api:"required"`
	// RFC 3339 timestamp when work execution started
	StartedAt string `json:"started_at" api:"required"`
	// Current state of the work item
	//
	// Any of "queued", "starting", "active", "stopping", "stopped".
	State SelfHostedWorkState `json:"state" api:"required"`
	// RFC 3339 timestamp when stop was requested
	StopRequestedAt string `json:"stop_requested_at" api:"required"`
	// RFC 3339 timestamp when work execution stopped
	StoppedAt string `json:"stopped_at" api:"required"`
	// The type of object (always 'work')
	Type constant.Work `json:"type" default:"work"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID                respjson.Field
		AcknowledgedAt    respjson.Field
		CreatedAt         respjson.Field
		Data              respjson.Field
		EnvironmentID     respjson.Field
		LatestHeartbeatAt respjson.Field
		Metadata          respjson.Field
		Secret            respjson.Field
		StartedAt         respjson.Field
		State             respjson.Field
		StopRequestedAt   respjson.Field
		StoppedAt         respjson.Field
		Type              respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SelfHostedWork) RawJSON() string { return r.JSON.raw }
func (r *SelfHostedWork) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Current state of the work item
type SelfHostedWorkState string

const (
	SelfHostedWorkStateQueued   SelfHostedWorkState = "queued"
	SelfHostedWorkStateStarting SelfHostedWorkState = "starting"
	SelfHostedWorkStateActive   SelfHostedWorkState = "active"
	SelfHostedWorkStateStopping SelfHostedWorkState = "stopping"
	SelfHostedWorkStateStopped  SelfHostedWorkState = "stopped"
)

// Response after recording a heartbeat for a work item.
type SelfHostedWorkHeartbeatResponse struct {
	// RFC 3339 timestamp of the actual heartbeat from DB
	LastHeartbeat string `json:"last_heartbeat" api:"required"`
	// Whether the heartbeat succeeded in extending the lease
	LeaseExtended bool `json:"lease_extended" api:"required"`
	// Current state of the work item (active/stopping/stopped)
	//
	// Any of "queued", "starting", "active", "stopping", "stopped".
	State SelfHostedWorkHeartbeatResponseState `json:"state" api:"required"`
	// Effective TTL applied to the lease
	TTLSeconds int64 `json:"ttl_seconds" api:"required"`
	// The type of response
	Type constant.WorkHeartbeat `json:"type" default:"work_heartbeat"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		LastHeartbeat respjson.Field
		LeaseExtended respjson.Field
		State         respjson.Field
		TTLSeconds    respjson.Field
		Type          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SelfHostedWorkHeartbeatResponse) RawJSON() string { return r.JSON.raw }
func (r *SelfHostedWorkHeartbeatResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Current state of the work item (active/stopping/stopped)
type SelfHostedWorkHeartbeatResponseState string

const (
	SelfHostedWorkHeartbeatResponseStateQueued   SelfHostedWorkHeartbeatResponseState = "queued"
	SelfHostedWorkHeartbeatResponseStateStarting SelfHostedWorkHeartbeatResponseState = "starting"
	SelfHostedWorkHeartbeatResponseStateActive   SelfHostedWorkHeartbeatResponseState = "active"
	SelfHostedWorkHeartbeatResponseStateStopping SelfHostedWorkHeartbeatResponseState = "stopping"
	SelfHostedWorkHeartbeatResponseStateStopped  SelfHostedWorkHeartbeatResponseState = "stopped"
)

// Response when listing work items with cursor-based pagination.
type SelfHostedWorkListResponse struct {
	// List of work items
	Data []SelfHostedWork `json:"data" api:"required"`
	// Opaque cursor for fetching the next page of results
	NextPage string `json:"next_page" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Data        respjson.Field
		NextPage    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SelfHostedWorkListResponse) RawJSON() string { return r.JSON.raw }
func (r *SelfHostedWorkListResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Statistics about the work queue for an environment.
//
// Uses Redis Stream consumer group metrics for O(1) queries.
type SelfHostedWorkQueueStats struct {
	// Number of work items waiting to be picked up (lag from consumer group)
	Depth int64 `json:"depth" api:"required"`
	// RFC 3339 timestamp of oldest item in the work stream (includes both queued and
	// pending items), null if stream empty
	OldestQueuedAt string `json:"oldest_queued_at" api:"required"`
	// Number of work items being processed (polled but not acknowledged)
	Pending int64 `json:"pending" api:"required"`
	// The type of object
	Type constant.WorkQueueStats `json:"type" default:"work_queue_stats"`
	// Number of workers that have polled for work in the last 30 seconds. Requires
	// worker_id to be sent with poll requests.
	WorkersPolling int64 `json:"workers_polling" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Depth          respjson.Field
		OldestQueuedAt respjson.Field
		Pending        respjson.Field
		Type           respjson.Field
		WorkersPolling respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SelfHostedWorkQueueStats) RawJSON() string { return r.JSON.raw }
func (r *SelfHostedWorkQueueStats) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Request to stop a work item.
type SelfHostedWorkStopRequestParam struct {
	// If true, immediately stop work without graceful shutdown
	Force param.Opt[bool] `json:"force,omitzero"`
	paramObj
}

func (r SelfHostedWorkStopRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow SelfHostedWorkStopRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SelfHostedWorkStopRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Request to update work item metadata.
//
// The property Metadata is required.
type SelfHostedWorkUpdateRequestParam struct {
	// Metadata patch. Set a key to a string to upsert it, or to null to delete it.
	// Omit the field to preserve existing metadata.
	Metadata map[string]any `json:"metadata,omitzero" api:"required"`
	paramObj
}

func (r SelfHostedWorkUpdateRequestParam) MarshalJSON() (data []byte, err error) {
	type shadow SelfHostedWorkUpdateRequestParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SelfHostedWorkUpdateRequestParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Work data for session work items.
//
// This resource type is used when work represents a session that needs to be
// executed in a self-hosted environment.
type SessionWorkData struct {
	// Session identifier (e.g., 'session\_...')
	ID string `json:"id" api:"required"`
	// Type of work data
	Type constant.Session `json:"type" default:"session"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r SessionWorkData) RawJSON() string { return r.JSON.raw }
func (r *SessionWorkData) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type EnvironmentWorkGetParams struct {
	EnvironmentID string `path:"environment_id" api:"required" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

type EnvironmentWorkUpdateParams struct {
	EnvironmentID string `path:"environment_id" api:"required" json:"-"`
	// Request to update work item metadata.
	SelfHostedWorkUpdateRequest SelfHostedWorkUpdateRequestParam
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

func (r EnvironmentWorkUpdateParams) MarshalJSON() (data []byte, err error) {
	return param.MarshalObject(r, r.SelfHostedWorkUpdateRequest)
}
func (r *EnvironmentWorkUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, &r.SelfHostedWorkUpdateRequest)
}

type EnvironmentWorkListParams struct {
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	AfterID  param.Opt[string] `query:"after_id,omitzero" json:"-"`

	// Opaque cursor from previous response for pagination
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Maximum number of work items to return
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [EnvironmentWorkListParams]'s query parameters as
// `url.Values`.
func (r EnvironmentWorkListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type EnvironmentWorkAckParams struct {
	EnvironmentID string `path:"environment_id" api:"required" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

type EnvironmentWorkHeartbeatParams struct {
	EnvironmentID string `path:"environment_id" api:"required" json:"-"`
	// Desired TTL in seconds
	DesiredTTLSeconds param.Opt[int64] `query:"desired_ttl_seconds,omitzero" json:"-"`
	// Expected last_heartbeat for conditional update (optimistic concurrency). Use
	// literal 'NO_HEARTBEAT' to claim an unclaimed lease (first heartbeat). For
	// subsequent heartbeats, echo the server's previous last_heartbeat value exactly.
	// Returns 412 Precondition Failed if the actual value doesn't match.
	ExpectedLastHeartbeat param.Opt[string] `query:"expected_last_heartbeat,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [EnvironmentWorkHeartbeatParams]'s query parameters as
// `url.Values`.
func (r EnvironmentWorkHeartbeatParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type EnvironmentWorkPollParams struct {
	// How long to wait for work to arrive before returning. Must be 1-999 in
	// milliseconds. Defaults to non-blocking (returns immediately if no work is
	// available).
	BlockMs param.Opt[int64] `query:"block_ms,omitzero" json:"-"`
	// Reclaim unacknowledged work items older than this many milliseconds. If omitted,
	// uses the default (5000ms).
	ReclaimOlderThanMs param.Opt[int64] `query:"reclaim_older_than_ms,omitzero" json:"-"`
	// Unique identifier for the specific worker polling, used to track aggregated
	// environment-level work metrics in Console
	QoderWorkerID param.Opt[string] `header:"Worker-ID,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [EnvironmentWorkPollParams]'s query parameters as
// `url.Values`.
func (r EnvironmentWorkPollParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type EnvironmentWorkStatsParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

type EnvironmentWorkStopParams struct {
	EnvironmentID string `path:"environment_id" api:"required" json:"-"`
	// Request to stop a work item.
	SelfHostedWorkStopRequest SelfHostedWorkStopRequestParam
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

func (r EnvironmentWorkStopParams) MarshalJSON() (data []byte, err error) {
	return param.MarshalObject(r, r.SelfHostedWorkStopRequest)
}
func (r *EnvironmentWorkStopParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, &r.SelfHostedWorkStopRequest)
}
