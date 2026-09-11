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

// 列出 Batches.
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
	// 按状态过滤。
	Status param.Opt[string] `query:"status,omitzero" json:"-"`
	// 分页大小，最大 100。
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// 向后翻页游标。
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// 向前翻页游标。
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	paramObj
}

func (r BatchListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// 创建 Batch.
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
	// 通过 Files API 上传的 JSONL 文件 ID。
	InputFileID string `json:"input_file_id" api:"required"`
	// 完成窗口：`24h`、`48h`、`72h`。超时后 Batch 自动进入 `expired` 状态。
	CompletionWindow string `json:"completion_window" api:"required"`
	// 调用方业务元数据，最多 16 个 key；value 可为任意 JSON 类型；整体序列化后 ≤ 2KB，key ≤ 64 字符，且不得包含 NUL（U+0000）。
	Metadata map[string]any `json:"metadata,omitzero"`
	// 有副作用请求可选的幂等键。
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r BatchNewParams) MarshalJSON() ([]byte, error) {
	type shadow BatchNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *BatchNewParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

// 查询 Batch 详情.
func (r *BatchService) Get(ctx context.Context, batchID string, opts ...option.RequestOption) (res *Batch, err error) {
	if batchID == "" {
		return nil, fmt.Errorf("missing required batch_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("batches/%s", url.PathEscape(batchID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// 取消 Batch.
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
	// 有副作用请求可选的幂等键。
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r BatchCancelParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// 获取错误文件.
func (r *BatchService) GetError(ctx context.Context, batchID string, opts ...option.RequestOption) (res *BatchFile, err error) {
	if batchID == "" {
		return nil, fmt.Errorf("missing required batch_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("batches/%s/error", url.PathEscape(batchID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// 获取输出文件.
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
	// OSS 预签名下载链接，含 `Expires` / `OSSAccessKeyId` / `Signature` 及 `response-content-disposition`，下载文件名为 `batch-<batch_id>-output.jsonl`。
	URL string `json:"url"`
	// 链接过期时间，RFC 3339，需在此之前完成下载。
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
	// Batch ID，前缀 `batch_`。
	ID string `json:"id"`
	// 固定为 `batch`。
	Object string `json:"object"`
	// Batch 状态，见状态说明。
	Status string `json:"status"`
	// 输入 JSONL 文件 ID。
	InputFileID string `json:"input_file_id"`
	// 成功结果文件 ID；未完成或未生成时省略。
	OutputFileID string `json:"output_file_id"`
	// 完成窗口：`24h`、`48h`、`72h`。
	CompletionWindow string `json:"completion_window"`
	// 创建时间，RFC 3339。
	CreatedAt time.Time `json:"created_at" format:"date-time"`
	// 过期时间，`created_at` + `completion_window`。
	ExpiresAt time.Time `json:"expires_at" format:"date-time"`
	// 任务计数聚合。
	RequestCounts BatchRequestCounts `json:"request_counts"`
	// 创建响应为 `null`；后续 Batch 详情、列表和取消响应中，至少一个子任务已有合法 CAS Session 用量时返回 Credit 汇总。
	Usage BatchUsage `json:"usage" api:"nullable"`
	// 调用方业务元数据。
	Metadata map[string]any `json:"metadata"`
	// 总行数（含校验失败行）。
	Total int64 `json:"total"`
	// 等待执行的行数。
	Pending int64 `json:"pending"`
	// 正在执行的行数。
	Running int64 `json:"running"`
	// 执行成功的行数。
	Completed int64 `json:"completed"`
	// 永久失败的行数（含校验失败）。
	Failed int64 `json:"failed"`
	// 因取消而终止的行数。
	Cancelled int64 `json:"cancelled"`
	// 因过期而终止的行数。
	Expired int64 `json:"expired"`
	// 失败行结果文件 ID；无失败行时省略。
	ErrorFileID string `json:"error_file_id"`
	// Batch 级错误描述；仅 `failed` 状态出现。
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
