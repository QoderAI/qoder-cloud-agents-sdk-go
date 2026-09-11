package forward

import (
	"context"
	"encoding/json"
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

// MemoryStoreMemoryVersionService provides Forward MemoryStoreMemoryVersion operations.
type MemoryStoreMemoryVersionService struct {
	Options []option.RequestOption
}

func NewMemoryStoreMemoryVersionService(opts ...option.RequestOption) MemoryStoreMemoryVersionService {
	return MemoryStoreMemoryVersionService{Options: slices.Clone(opts)}
}

// 列出 Memory 版本.
func (r *MemoryStoreMemoryVersionService) List(ctx context.Context, memoryStoreID string, params MemoryStoreMemoryVersionListParams, opts ...option.RequestOption) (res *pagination.Page[MemoryVersion], err error) {
	if memoryStoreID == "" {
		return nil, fmt.Errorf("missing required memory_store_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("memory_stores/%s/memory_versions", url.PathEscape(memoryStoreID))
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
func (r *MemoryStoreMemoryVersionService) ListAutoPaging(ctx context.Context, memoryStoreID string, params MemoryStoreMemoryVersionListParams, opts ...option.RequestOption) *pagination.PageAutoPager[MemoryVersion] {
	return pagination.NewPageAutoPager(r.List(ctx, memoryStoreID, params, opts...))
}

type MemoryStoreMemoryVersionListParams struct {
	// 每页返回数量上限，1..100，默认 20。
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// 向前翻页游标，与 `after_id` 互斥。
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	// 向后翻页游标，与 `before_id` 互斥。
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// 只返回该 memory（`mem_...`）的版本，用于查看单条记忆的变更历史。
	MemoryID param.Opt[string] `query:"memory_id,omitzero" json:"-"`
	paramObj
}

func (r MemoryStoreMemoryVersionListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// 查询 Memory 版本.
func (r *MemoryStoreMemoryVersionService) Get(ctx context.Context, memoryStoreID string, memoryVersionID string, opts ...option.RequestOption) (res *MemoryVersion, err error) {
	if memoryStoreID == "" {
		return nil, fmt.Errorf("missing required memory_store_id parameter")
	}
	if memoryVersionID == "" {
		return nil, fmt.Errorf("missing required memory_version_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("memory_stores/%s/memory_versions/%s", url.PathEscape(memoryStoreID), url.PathEscape(memoryVersionID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Redact Memory 版本.
func (r *MemoryStoreMemoryVersionService) Redact(ctx context.Context, memoryStoreID string, memoryVersionID string, opts ...option.RequestOption) (res *MemoryVersion, err error) {
	if memoryStoreID == "" {
		return nil, fmt.Errorf("missing required memory_store_id parameter")
	}
	if memoryVersionID == "" {
		return nil, fmt.Errorf("missing required memory_version_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("memory_stores/%s/memory_versions/%s/redact", url.PathEscape(memoryStoreID), url.PathEscape(memoryVersionID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type MemoryVersion struct {
	ID               string          `json:"id"`
	Type             string          `json:"type"`
	MemoryStoreID    string          `json:"memory_store_id"`
	MemoryID         string          `json:"memory_id"`
	Path             string          `json:"path"`
	ContentSizeBytes int64           `json:"content_size_bytes"`
	ContentSHA256    string          `json:"content_sha256"`
	Operation        string          `json:"operation"`
	Redacted         bool            `json:"redacted"`
	RedactedAt       time.Time       `json:"redacted_at" api:"nullable" format:"date-time"`
	CreatedAt        time.Time       `json:"created_at" format:"date-time"`
	Content          json.RawMessage `json:"content"`
	JSON             struct {
		ID               respjson.Field
		Type             respjson.Field
		MemoryStoreID    respjson.Field
		MemoryID         respjson.Field
		Path             respjson.Field
		ContentSizeBytes respjson.Field
		ContentSHA256    respjson.Field
		Operation        respjson.Field
		Redacted         respjson.Field
		RedactedAt       respjson.Field
		CreatedAt        respjson.Field
		Content          respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

func (r MemoryVersion) RawJSON() string                  { return r.JSON.raw }
func (r *MemoryVersion) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
