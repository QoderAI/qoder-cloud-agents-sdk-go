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

// BatchTaskService provides Forward BatchTask operations.
type BatchTaskService struct {
	Options []option.RequestOption
}

func NewBatchTaskService(opts ...option.RequestOption) BatchTaskService {
	return BatchTaskService{Options: slices.Clone(opts)}
}

// 查询 Batch 子任务.
func (r *BatchTaskService) List(ctx context.Context, batchID string, params BatchTaskListParams, opts ...option.RequestOption) (res *pagination.Page[BatchTask], err error) {
	if batchID == "" {
		return nil, fmt.Errorf("missing required batch_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("batches/%s/tasks", url.PathEscape(batchID))
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
func (r *BatchTaskService) ListAutoPaging(ctx context.Context, batchID string, params BatchTaskListParams, opts ...option.RequestOption) *pagination.PageAutoPager[BatchTask] {
	return pagination.NewPageAutoPager(r.List(ctx, batchID, params, opts...))
}

type BatchTaskListParams struct {
	// 按任务状态过滤：`pending`、`running`、`completed`、`failed`、`cancelled`、`expired`。
	Status param.Opt[string] `query:"status,omitzero" json:"-"`
	// 按调用方任务标识精确过滤，仅支持单值；未命中返回空列表。
	CustomID param.Opt[string] `query:"custom_id,omitzero" json:"-"`
	// 分页大小，最大 100。
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// 向后翻页游标，传上一页响应的 `last_id`；游标必须属于当前 Batch。
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	paramObj
}

func (r BatchTaskListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

type BatchTask struct {
	CustomID      string                   `json:"custom_id"`
	Status        string                   `json:"status"`
	StartedAt     time.Time                `json:"started_at" format:"date-time"`
	CompletedAt   time.Time                `json:"completed_at" format:"date-time"`
	OutputSummary string                   `json:"output_summary"`
	Usage         BatchTaskUsage           `json:"usage"`
	Artifacts     []BatchTaskArtifactsItem `json:"artifacts"`
	Error         BatchTaskError           `json:"error"`
	JSON          struct {
		CustomID      respjson.Field
		Status        respjson.Field
		StartedAt     respjson.Field
		CompletedAt   respjson.Field
		OutputSummary respjson.Field
		Usage         respjson.Field
		Artifacts     respjson.Field
		Error         respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

func (r BatchTask) RawJSON() string                  { return r.JSON.raw }
func (r *BatchTask) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type BatchTaskError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	JSON    struct {
		Code        respjson.Field
		Message     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r BatchTaskError) RawJSON() string                  { return r.JSON.raw }
func (r *BatchTaskError) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type BatchTaskArtifactsItem struct {
	FileID      string `json:"file_id"`
	Name        string `json:"name"`
	Size        int64  `json:"size"`
	ContentType string `json:"content_type"`
	JSON        struct {
		FileID      respjson.Field
		Name        respjson.Field
		Size        respjson.Field
		ContentType respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r BatchTaskArtifactsItem) RawJSON() string { return r.JSON.raw }
func (r *BatchTaskArtifactsItem) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type BatchTaskUsage struct {
	TotalCredits float64 `json:"total_credits"`
	JSON         struct {
		TotalCredits respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

func (r BatchTaskUsage) RawJSON() string                  { return r.JSON.raw }
func (r *BatchTaskUsage) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
