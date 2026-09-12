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

// List Memories
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
	// Maximum number of entries per page, 1..100, default 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Cursor for the previous page; mutually exclusive with `after_id`.
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	// Cursor for the next page; mutually exclusive with `before_id`.
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// Filter by `path` prefix. **This is a plain string prefix match, not directory
	// semantics** — `path_prefix=a/b` also matches `a/bc.md`.
	PathPrefix param.Opt[string] `query:"path_prefix,omitzero" json:"-"`
	paramObj
}

func (r MemoryStoreMemoryListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// Create Memory
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
	// Path relative to the store root, case-sensitive; it must be relative, with no
	// leading `/`.
	Path string `json:"path" api:"required"`
	// UTF-8 plain text content, not base64; at most 100 KiB of raw bytes.
	Content string `json:"content" api:"required"`
	// Key-value metadata whose values must be strings, up to **16** keys.
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

// Get Memory
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

// Update Memory
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
	// New content as UTF-8 plain text; at most 100 KiB of raw bytes.
	Content string `json:"content" api:"required"`
	// Expected SHA-256 of the current content, for optimistic concurrency control.
	// Returns `409` when it does not match.
	ContentSHA256 param.Opt[string] `json:"content_sha256,omitzero"`
	// New metadata, which **replaces** the current metadata wholesale rather than
	// merging into it. Omit to leave the existing metadata unchanged.
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

// Delete Memory
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
	// ID of the deleted Memory.
	ID string `json:"id"`
	// Always `memory_deleted`.
	Type string `json:"type"`
	// Whether the Memory was deleted.
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
