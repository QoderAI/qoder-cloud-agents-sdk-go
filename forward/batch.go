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

// BatchService provides Forward Batch operations.
type BatchService struct {
	Options []option.RequestOption
	Tasks   BatchTaskService
}

func NewBatchService(opts ...option.RequestOption) BatchService {
	return BatchService{Options: slices.Clone(opts), Tasks: NewBatchTaskService(opts...)}
}

// List Batches
func (r *BatchService) List(ctx context.Context, params BatchListParams, opts ...option.RequestOption) (res *pagination.Page[Batch], err error) {

	opts = slices.Concat(r.Options, opts)
	path := "batches"
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
func (r *BatchService) ListAutoPaging(ctx context.Context, params BatchListParams, opts ...option.RequestOption) *pagination.PageAutoPager[Batch] {
	return pagination.NewPageAutoPager(r.List(ctx, params, opts...))
}

type BatchListParams struct {
	// Filter by status.
	Status param.Opt[string] `query:"status,omitzero" json:"-"`
	// Page size, up to 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Cursor for the next page.
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// Cursor for the previous page.
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	paramObj
}

func (r BatchListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// Create Batch
func (r *BatchService) New(ctx context.Context, params BatchNewParams, opts ...option.RequestOption) (res *Batch, err error) {
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := "batches"
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type BatchNewParams struct {
	// ID of a JSONL file uploaded through the Files API.
	InputFileID string `json:"input_file_id" api:"required"`
	// Completion window: `24h`, `48h` or `72h`. The Batch moves to `expired` when it elapses.
	CompletionWindow string `json:"completion_window" api:"required"`
	// Caller business metadata, up to 16 keys; values may be any JSON type; the serialized
	// object must be ≤ 2KB, keys ≤ 64 characters, and must not contain NUL (U+0000).
	Metadata map[string]any `json:"metadata,omitzero"`
	// Optional idempotency key for requests with side effects.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r BatchNewParams) MarshalJSON() ([]byte, error) {
	type shadow BatchNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BatchNewParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

// Get Batch
func (r *BatchService) Get(ctx context.Context, batchID string, opts ...option.RequestOption) (res *Batch, err error) {
	if batchID == "" {
		return nil, fmt.Errorf("missing required batch_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("batches/%s", url.PathEscape(batchID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Cancel Batch
func (r *BatchService) Cancel(ctx context.Context, batchID string, params BatchCancelParams, opts ...option.RequestOption) (res *Batch, err error) {
	if batchID == "" {
		return nil, fmt.Errorf("missing required batch_id parameter")
	}
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("batches/%s/cancel", url.PathEscape(batchID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type BatchCancelParams struct {
	// Optional idempotency key for requests with side effects.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r BatchCancelParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// Get error file
func (r *BatchService) GetError(ctx context.Context, batchID string, opts ...option.RequestOption) (res *BatchFile, err error) {
	if batchID == "" {
		return nil, fmt.Errorf("missing required batch_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("batches/%s/error", url.PathEscape(batchID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Get output file
func (r *BatchService) GetOutput(ctx context.Context, batchID string, opts ...option.RequestOption) (res *BatchFile, err error) {
	if batchID == "" {
		return nil, fmt.Errorf("missing required batch_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("batches/%s/output", url.PathEscape(batchID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

type BatchFile struct {
	// Pre-signed OSS download URL with `Expires` / `OSSAccessKeyId` / `Signature` and
	// `response-content-disposition`; the file downloads as `batch-<batch_id>-output.jsonl`.
	URL string `json:"url"`
	// Link expiration time in RFC 3339; the download must finish before it.
	ExpiresAt time.Time `json:"expires_at" format:"date-time"`
	JSON      struct {
		URL         respjson.Field
		ExpiresAt   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r BatchFile) RawJSON() string                  { return r.JSON.raw }
func (r *BatchFile) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type Batch struct {
	// Batch ID, prefixed with `batch_`.
	ID string `json:"id"`
	// Always `batch`.
	Object string `json:"object"`
	// Batch status; see the status descriptions.
	Status string `json:"status"`
	// Input JSONL file ID.
	InputFileID string `json:"input_file_id"`
	// File ID of successful results; omitted while incomplete or not yet generated.
	OutputFileID string `json:"output_file_id"`
	// Completion window: `24h`, `48h` or `72h`.
	CompletionWindow string `json:"completion_window"`
	// Creation time in RFC 3339.
	CreatedAt time.Time `json:"created_at" format:"date-time"`
	// Expiration time, `created_at` + `completion_window`.
	ExpiresAt time.Time `json:"expires_at" format:"date-time"`
	// Aggregated task counts.
	RequestCounts BatchRequestCounts `json:"request_counts"`
	// `null` on the create response; on later Batch get, list and cancel responses it returns
	// the Credit total once at least one subtask has valid CAS Session usage.
	Usage BatchUsage `json:"usage" api:"nullable"`
	// Caller business metadata.
	Metadata map[string]any `json:"metadata"`
	// Total number of lines, including lines that failed validation.
	Total int64 `json:"total"`
	// Number of lines waiting to run.
	Pending int64 `json:"pending"`
	// Number of lines currently running.
	Running int64 `json:"running"`
	// Number of lines that completed successfully.
	Completed int64 `json:"completed"`
	// Number of permanently failed lines, including validation failures.
	Failed int64 `json:"failed"`
	// Number of lines terminated by cancellation.
	Cancelled int64 `json:"cancelled"`
	// Number of lines terminated by expiration.
	Expired int64 `json:"expired"`
	// File ID of failed-line results; omitted when there are no failed lines.
	ErrorFileID string `json:"error_file_id"`
	// Batch-level error description; present only in the `failed` status.
	ErrorMessage string `json:"error_message"`
	JSON         struct {
		ID               respjson.Field
		Object           respjson.Field
		Status           respjson.Field
		InputFileID      respjson.Field
		OutputFileID     respjson.Field
		CompletionWindow respjson.Field
		CreatedAt        respjson.Field
		ExpiresAt        respjson.Field
		RequestCounts    respjson.Field
		Usage            respjson.Field
		Metadata         respjson.Field
		Total            respjson.Field
		Pending          respjson.Field
		Running          respjson.Field
		Completed        respjson.Field
		Failed           respjson.Field
		Cancelled        respjson.Field
		Expired          respjson.Field
		ErrorFileID      respjson.Field
		ErrorMessage     respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

func (r Batch) RawJSON() string                  { return r.JSON.raw }
func (r *Batch) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type BatchUsage struct {
	TotalCredits float64 `json:"total_credits"`
	JSON         struct {
		TotalCredits respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

func (r BatchUsage) RawJSON() string                  { return r.JSON.raw }
func (r *BatchUsage) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type BatchRequestCounts struct {
	Total     int64 `json:"total"`
	Pending   int64 `json:"pending"`
	Running   int64 `json:"running"`
	Completed int64 `json:"completed"`
	Failed    int64 `json:"failed"`
	Cancelled int64 `json:"cancelled"`
	Expired   int64 `json:"expired"`
	JSON      struct {
		Total       respjson.Field
		Pending     respjson.Field
		Running     respjson.Field
		Completed   respjson.Field
		Failed      respjson.Field
		Cancelled   respjson.Field
		Expired     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r BatchRequestCounts) RawJSON() string                  { return r.JSON.raw }
func (r *BatchRequestCounts) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
