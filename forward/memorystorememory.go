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

// MemoryStoreMemoryService provides Forward MemoryStoreMemory operations.
type MemoryStoreMemoryService struct {
	Options []option.RequestOption
}

func NewMemoryStoreMemoryService(opts ...option.RequestOption) MemoryStoreMemoryService {
	return MemoryStoreMemoryService{Options: slices.Clone(opts)}
}

// 列出 Memory.
func (r *MemoryStoreMemoryService) List(ctx context.Context, memoryStoreID string, params MemoryStoreMemoryListParams, opts ...option.RequestOption) (res *pagination.Page[Memory], err error) {
	if memoryStoreID == "" {
		return nil, fmt.Errorf("missing required memory_store_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("memory_stores/%s/memories", url.PathEscape(memoryStoreID))
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
func (r *MemoryStoreMemoryService) ListAutoPaging(ctx context.Context, memoryStoreID string, params MemoryStoreMemoryListParams, opts ...option.RequestOption) *pagination.PageAutoPager[Memory] {
	return pagination.NewPageAutoPager(r.List(ctx, memoryStoreID, params, opts...))
}

type MemoryStoreMemoryListParams struct {
	// 每页返回数量上限，1..100，默认 20。
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// 向前翻页游标，与 `after_id` 互斥。
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	// 向后翻页游标，与 `before_id` 互斥。
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// 按 `path` 前缀过滤。**这是纯字符串前缀匹配，不是目录语义** —— `path_prefix=a/b` 也会命中 `a/bc.md`。
	PathPrefix param.Opt[string] `query:"path_prefix,omitzero" json:"-"`
	paramObj
}

func (r MemoryStoreMemoryListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// 创建 Memory.
func (r *MemoryStoreMemoryService) New(ctx context.Context, memoryStoreID string, params MemoryStoreMemoryNewParams, opts ...option.RequestOption) (res *Memory, err error) {
	if memoryStoreID == "" {
		return nil, fmt.Errorf("missing required memory_store_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("memory_stores/%s/memories", url.PathEscape(memoryStoreID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type MemoryStoreMemoryNewParams struct {
	// 库内相对路径，大小写敏感。约束详见 [path 规则](../MemoryStore数据结构.md#path-规则)。
	Path string `json:"path" api:"required"`
	// UTF-8 明文内容，非 base64；原始字节 ≤100 KiB。约束详见 [content 约束](../MemoryStore数据结构.md#content-约束)。
	Content string `json:"content" api:"required"`
	// 键值元数据，值必须为字符串。最多 **16** 个键。约束详见 [Memory metadata 约束](../MemoryStore数据结构.md#memory-metadata-约束)。
	Metadata map[string]any `json:"metadata,omitzero"`
	paramObj
}

func (r MemoryStoreMemoryNewParams) MarshalJSON() ([]byte, error) {
	type shadow MemoryStoreMemoryNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *MemoryStoreMemoryNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// 查询 Memory.
func (r *MemoryStoreMemoryService) Get(ctx context.Context, memoryStoreID string, memoryID string, opts ...option.RequestOption) (res *Memory, err error) {
	if memoryStoreID == "" {
		return nil, fmt.Errorf("missing required memory_store_id parameter")
	}
	if memoryID == "" {
		return nil, fmt.Errorf("missing required memory_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("memory_stores/%s/memories/%s", url.PathEscape(memoryStoreID), url.PathEscape(memoryID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// 更新 Memory.
func (r *MemoryStoreMemoryService) Update(ctx context.Context, memoryStoreID string, memoryID string, params MemoryStoreMemoryUpdateParams, opts ...option.RequestOption) (res *Memory, err error) {
	if memoryStoreID == "" {
		return nil, fmt.Errorf("missing required memory_store_id parameter")
	}
	if memoryID == "" {
		return nil, fmt.Errorf("missing required memory_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("memory_stores/%s/memories/%s", url.PathEscape(memoryStoreID), url.PathEscape(memoryID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type MemoryStoreMemoryUpdateParams struct {
	// 新内容，UTF-8 明文；原始字节 ≤100 KiB。约束详见 [content 约束](../MemoryStore数据结构.md#content-约束)。
	Content string `json:"content" api:"required"`
	// 期望的当前内容 SHA-256，用于乐观并发控制。不一致时返回 `409`。
	ContentSHA256 param.Opt[string] `json:"content_sha256,omitzero"`
	// 新元数据，**整体替换**当前 metadata（非合并）。未传入时保持原 metadata 不变。约束详见 [Memory metadata 约束](../MemoryStore数据结构.md#memory-metadata-约束)。
	Metadata map[string]any `json:"metadata,omitzero"`
	paramObj
}

func (r MemoryStoreMemoryUpdateParams) MarshalJSON() ([]byte, error) {
	type shadow MemoryStoreMemoryUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *MemoryStoreMemoryUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// 删除 Memory.
func (r *MemoryStoreMemoryService) Delete(ctx context.Context, memoryStoreID string, memoryID string, opts ...option.RequestOption) (res *DeletedMemory, err error) {
	if memoryStoreID == "" {
		return nil, fmt.Errorf("missing required memory_store_id parameter")
	}
	if memoryID == "" {
		return nil, fmt.Errorf("missing required memory_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("memory_stores/%s/memories/%s", url.PathEscape(memoryStoreID), url.PathEscape(memoryID))
	err = convention.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

type DeletedMemory struct {
	// 被删除的 Memory ID。
	ID string `json:"id"`
	// 固定为 `memory_deleted`。
	Type string `json:"type"`
	// 是否已删除。
	Deleted bool `json:"deleted"`
	JSON    struct {
		ID          respjson.Field
		Type        respjson.Field
		Deleted     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r DeletedMemory) RawJSON() string                  { return r.JSON.raw }
func (r *DeletedMemory) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type Memory struct {
	ID               string         `json:"id"`
	Type             string         `json:"type"`
	MemoryStoreID    string         `json:"memory_store_id"`
	Path             string         `json:"path"`
	ContentSizeBytes int64          `json:"content_size_bytes"`
	ContentSHA256    string         `json:"content_sha256"`
	Metadata         map[string]any `json:"metadata"`
	CreatedAt        time.Time      `json:"created_at" format:"date-time"`
	UpdatedAt        time.Time      `json:"updated_at" format:"date-time"`
	Content          string         `json:"content"`
	JSON             struct {
		ID               respjson.Field
		Type             respjson.Field
		MemoryStoreID    respjson.Field
		Path             respjson.Field
		ContentSizeBytes respjson.Field
		ContentSHA256    respjson.Field
		Metadata         respjson.Field
		CreatedAt        respjson.Field
		UpdatedAt        respjson.Field
		Content          respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

func (r Memory) RawJSON() string                  { return r.JSON.raw }
func (r *Memory) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
