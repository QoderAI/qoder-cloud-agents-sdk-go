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

// MemoryStoreService provides Forward MemoryStore operations.
type MemoryStoreService struct {
	Options        []option.RequestOption
	Memories       MemoryStoreMemoryService
	MemoryVersions MemoryStoreMemoryVersionService
}

func NewMemoryStoreService(opts ...option.RequestOption) MemoryStoreService {
	return MemoryStoreService{Options: slices.Clone(opts), Memories: NewMemoryStoreMemoryService(opts...), MemoryVersions: NewMemoryStoreMemoryVersionService(opts...)}
}

// 列出 Memory Store.
func (r *MemoryStoreService) List(ctx context.Context, params MemoryStoreListParams, opts ...option.RequestOption) (res *pagination.Page[MemoryStore], err error) {

	opts = slices.Concat(r.Options, opts)
	path := "memory_stores"
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
func (r *MemoryStoreService) ListAutoPaging(ctx context.Context, params MemoryStoreListParams, opts ...option.RequestOption) *pagination.PageAutoPager[MemoryStore] {
	return pagination.NewPageAutoPager(r.List(ctx, params, opts...))
}

type MemoryStoreListParams struct {
	// 每页返回数量上限，1..100，默认 20。
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// 向前翻页游标，与 `after_id` 互斥。
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	// 向后翻页游标，与 `before_id` 互斥。
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// 三态过滤：`true` 只返回系统默认库；`false` 只返回用户创建的库；不传不过滤。
	SystemManaged param.Opt[bool] `query:"system_managed,omitzero" json:"-"`
	paramObj
}

func (r MemoryStoreListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// 创建 Memory Store.
func (r *MemoryStoreService) New(ctx context.Context, params MemoryStoreNewParams, opts ...option.RequestOption) (res *MemoryStore, err error) {
	opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey))}, opts...)
	opts = slices.Concat(r.Options, opts)
	path := "memory_stores"
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type MemoryStoreNewParams struct {
	// Store 展示名，非空。不允许非打印控制字符（`U+0000`–`U+001F`、`U+007F`），换行 `\n`、回车 `\r`、制表 `\t` 除外。
	Name string `json:"name" api:"required"`
	// 自由文本描述。不允许非打印控制字符。
	Description param.Opt[string] `json:"description,omitzero"`
	// 键值元数据，值必须为字符串。最多 **15** 个键；键 1..64 字符；值 ≤512 字符。`created_by` 是 Forward 保留键，服务端自动写入 `"forward"`；调用方传入 `created_by` 会返回 `400 invalid_request_error`。详见 [Store metadata 约束](./MemoryStore数据结构.md#store-metadata-约束)。
	Metadata map[string]any `json:"metadata,omitzero"`
	// 创建请求幂等键。相同 key 只能用于相同请求体；不传返回 `400`。
	IdempotencyKey string `header:"Idempotency-Key,omitzero" json:"-" api:"required"`
	paramObj
}

func (r MemoryStoreNewParams) MarshalJSON() ([]byte, error) {
	type shadow MemoryStoreNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *MemoryStoreNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// 查询 Memory Store.
func (r *MemoryStoreService) Get(ctx context.Context, memoryStoreID string, opts ...option.RequestOption) (res *MemoryStore, err error) {
	if memoryStoreID == "" {
		return nil, fmt.Errorf("missing required memory_store_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("memory_stores/%s", url.PathEscape(memoryStoreID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// 更新 Memory Store.
func (r *MemoryStoreService) Update(ctx context.Context, memoryStoreID string, params MemoryStoreUpdateParams, opts ...option.RequestOption) (res *MemoryStore, err error) {
	if memoryStoreID == "" {
		return nil, fmt.Errorf("missing required memory_store_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("memory_stores/%s", url.PathEscape(memoryStoreID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type MemoryStoreUpdateParams struct {
	// 新名称。传入时非空且不含非打印控制字符。
	Name param.Opt[string] `json:"name,omitzero"`
	// 新描述。传入时不含非打印控制字符。
	Description param.Opt[string] `json:"description,omitzero"`
	// 新元数据，**整体替换**当前 metadata（非合并）。约束详见 [Store metadata 约束](./MemoryStore数据结构.md#store-metadata-约束)。
	Metadata map[string]any `json:"metadata,omitzero"`
	paramObj
}

func (r MemoryStoreUpdateParams) MarshalJSON() ([]byte, error) {
	type shadow MemoryStoreUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *MemoryStoreUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// 删除 Memory Store.
func (r *MemoryStoreService) Delete(ctx context.Context, memoryStoreID string, opts ...option.RequestOption) (res *DeletedMemoryStore, err error) {
	if memoryStoreID == "" {
		return nil, fmt.Errorf("missing required memory_store_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("memory_stores/%s", url.PathEscape(memoryStoreID))
	err = convention.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// 归档 Memory Store.
func (r *MemoryStoreService) Archive(ctx context.Context, memoryStoreID string, opts ...option.RequestOption) (res *MemoryStore, err error) {
	if memoryStoreID == "" {
		return nil, fmt.Errorf("missing required memory_store_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("memory_stores/%s/archive", url.PathEscape(memoryStoreID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type DeletedMemoryStore struct {
	// 被删除的 Memory Store ID。
	ID string `json:"id"`
	// 固定为 `memory_store_deleted`。
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

func (r DeletedMemoryStore) RawJSON() string                  { return r.JSON.raw }
func (r *DeletedMemoryStore) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type MemoryStore struct {
	ID            string           `json:"id"`
	Type          string           `json:"type"`
	Name          string           `json:"name"`
	Description   string           `json:"description"`
	Status        string           `json:"status"`
	EntryCount    int64            `json:"entry_count"`
	TotalSize     int64            `json:"total_size"`
	Metadata      map[string]any   `json:"metadata"`
	SystemManaged bool             `json:"system_managed"`
	IdentityID    string           `json:"identity_id" api:"nullable"`
	CreatedAt     time.Time        `json:"created_at" format:"date-time"`
	UpdatedAt     time.Time        `json:"updated_at" format:"date-time"`
	ArchivedAt    time.Time        `json:"archived_at" api:"nullable" format:"date-time"`
	BindingInfo   map[string]int64 `json:"binding_info"`
	JSON          struct {
		ID            respjson.Field
		Type          respjson.Field
		Name          respjson.Field
		Description   respjson.Field
		Status        respjson.Field
		EntryCount    respjson.Field
		TotalSize     respjson.Field
		Metadata      respjson.Field
		SystemManaged respjson.Field
		IdentityID    respjson.Field
		CreatedAt     respjson.Field
		UpdatedAt     respjson.Field
		ArchivedAt    respjson.Field
		BindingInfo   respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

func (r MemoryStore) RawJSON() string                  { return r.JSON.raw }
func (r *MemoryStore) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
