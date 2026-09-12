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

// List Memory Stores
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
	// Maximum number of items per page, 1..100, defaults to 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Cursor for the previous page, mutually exclusive with `after_id`.
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	// Cursor for the next page, mutually exclusive with `before_id`.
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// Three-state filter: `true` returns only system default stores; `false` returns only
	// user-created stores; omit it to apply no filter.
	SystemManaged param.Opt[bool] `query:"system_managed,omitzero" json:"-"`
	paramObj
}

func (r MemoryStoreListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// Create Memory Store
func (r *MemoryStoreService) New(ctx context.Context, params MemoryStoreNewParams, opts ...option.RequestOption) (res *MemoryStore, err error) {
	opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey))}, opts...)
	opts = slices.Concat(r.Options, opts)
	path := "memory_stores"
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type MemoryStoreNewParams struct {
	// Display name for the store, non-empty. Non-printable control characters
	// (`U+0000`–`U+001F`, `U+007F`) are not allowed, except newline `\n`, carriage return
	// `\r` and tab `\t`.
	Name string `json:"name" api:"required"`
	// Free-text description. Non-printable control characters are not allowed.
	Description param.Opt[string] `json:"description,omitzero"`
	// Key-value metadata whose values must be strings. At most **15** keys; keys 1..64
	// characters; values ≤512 characters. `created_by` is a Forward reserved key that the
	// server fills in with `"forward"`; sending `created_by` returns
	// `400 invalid_request_error`. See the Store metadata constraints in the MemoryStore
	// data structure reference for details.
	Metadata map[string]any `json:"metadata,omitzero"`
	// Idempotency key for the create request. The same key may only be reused with an
	// identical request body; omitting it returns `400`.
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

// Get Memory Store
func (r *MemoryStoreService) Get(ctx context.Context, memoryStoreID string, opts ...option.RequestOption) (res *MemoryStore, err error) {
	if memoryStoreID == "" {
		return nil, fmt.Errorf("missing required memory_store_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("memory_stores/%s", url.PathEscape(memoryStoreID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update Memory Store
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
	// New name. Must be non-empty and free of non-printable control characters when sent.
	Name param.Opt[string] `json:"name,omitzero"`
	// New description. Must be free of non-printable control characters when sent.
	Description param.Opt[string] `json:"description,omitzero"`
	// New metadata that **replaces** the current metadata entirely (not a merge). See the
	// Store metadata constraints in the MemoryStore data structure reference for details.
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

// Delete Memory Store
func (r *MemoryStoreService) Delete(ctx context.Context, memoryStoreID string, opts ...option.RequestOption) (res *DeletedMemoryStore, err error) {
	if memoryStoreID == "" {
		return nil, fmt.Errorf("missing required memory_store_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("memory_stores/%s", url.PathEscape(memoryStoreID))
	err = convention.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Archive Memory Store
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
	// ID of the deleted Memory Store.
	ID string `json:"id"`
	// Always `memory_store_deleted`.
	Type string `json:"type"`
	// Whether the Memory Store was deleted.
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
