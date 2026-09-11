// Qoder managed API definitions.
package managed

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"time"

	requestconfig "github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apiquery"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/pagination"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
)

// MemoryStoreService contains methods and other services that help with
// interacting with the Qoder Cloud Agents API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewMemoryStoreService] method instead.
type MemoryStoreService struct {
	Options        []option.RequestOption
	Memories       MemoryStoreMemoryService
	MemoryVersions MemoryStoreMemoryVersionService
}

// NewMemoryStoreService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewMemoryStoreService(opts ...option.RequestOption) (r MemoryStoreService) {
	r = MemoryStoreService{}
	r.Options = opts
	r.Memories = NewMemoryStoreMemoryService(opts...)
	r.MemoryVersions = NewMemoryStoreMemoryVersionService(opts...)
	return
}

// Create a memory store
func (r *MemoryStoreService) New(ctx context.Context, params MemoryStoreNewParams, opts ...option.RequestOption) (res *ManagedAgentsMemoryStore, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	path := "memory_stores"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Retrieve a memory store
func (r *MemoryStoreService) Get(ctx context.Context, memoryStoreID string, query MemoryStoreGetParams, opts ...option.RequestOption) (res *ManagedAgentsMemoryStore, err error) {
	for _, v := range query.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if memoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("memory_stores/%s", url.PathEscape(memoryStoreID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update a memory store
func (r *MemoryStoreService) Update(ctx context.Context, memoryStoreID string, params MemoryStoreUpdateParams, opts ...option.RequestOption) (res *ManagedAgentsMemoryStore, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if memoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("memory_stores/%s", url.PathEscape(memoryStoreID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List memory stores
func (r *MemoryStoreService) List(ctx context.Context, params MemoryStoreListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ManagedAgentsMemoryStore], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "memory_stores"
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

// List memory stores
func (r *MemoryStoreService) ListAutoPaging(ctx context.Context, params MemoryStoreListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ManagedAgentsMemoryStore] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// Delete a memory store
func (r *MemoryStoreService) Delete(ctx context.Context, memoryStoreID string, body MemoryStoreDeleteParams, opts ...option.RequestOption) (res *ManagedAgentsDeletedMemoryStore, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if memoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("memory_stores/%s", url.PathEscape(memoryStoreID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Archive a memory store
func (r *MemoryStoreService) Archive(ctx context.Context, memoryStoreID string, body MemoryStoreArchiveParams, opts ...option.RequestOption) (res *ManagedAgentsMemoryStore, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if memoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("memory_stores/%s/archive", url.PathEscape(memoryStoreID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Confirmation that a `memory_store` was deleted.
type ManagedAgentsDeletedMemoryStore struct {
	// ID of the deleted memory store (a `memstore_...` identifier). The store and all
	// its memories and versions are no longer retrievable.
	ID string `json:"id" api:"required"`
	// Any of "memory_store_deleted".
	Type ManagedAgentsDeletedMemoryStoreType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsDeletedMemoryStore) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsDeletedMemoryStore) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsDeletedMemoryStoreType string

const (
	ManagedAgentsDeletedMemoryStoreTypeMemoryStoreDeleted ManagedAgentsDeletedMemoryStoreType = "memory_store_deleted"
)

// A `memory_store`: a named container for agent memories, scoped to a workspace.
// Attach a store to a session via `resources[]` to mount it as a directory the
// agent can read and write.
type ManagedAgentsMemoryStore struct {
	// Unique identifier for the memory store (a `memstore_...` tagged ID). Use this
	// when attaching the store to a session, or in the `{memory_store_id}` path
	// parameter of subsequent calls.
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Human-readable name for the store. 1–255 characters. The store's mount-path slug
	// under `/mnt/memory/` is derived from this name.
	Name string `json:"name" api:"required"`
	// Any of "memory_store".
	Type ManagedAgentsMemoryStoreType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// A timestamp in RFC 3339 format
	ArchivedAt time.Time `json:"archived_at" api:"nullable" format:"date-time"`
	// Free-text description of what the store contains, up to 1024 characters.
	// Included in the agent's system prompt when the store is attached, so word it to
	// be useful to the agent. Empty string when unset.
	Description string `json:"description"`
	// Arbitrary key-value tags for your own bookkeeping (such as the end user a store
	// belongs to). Up to 16 pairs; keys 1–64 characters; values up to 512 characters.
	// Returned on retrieve/list but not filterable.
	Metadata map[string]string `json:"metadata"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		Name        respjson.Field
		Type        respjson.Field
		UpdatedAt   respjson.Field
		ArchivedAt  respjson.Field
		Description respjson.Field
		Metadata    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsMemoryStore) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsMemoryStore) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsMemoryStoreType string

const (
	ManagedAgentsMemoryStoreTypeMemoryStore ManagedAgentsMemoryStoreType = "memory_store"
)

type MemoryStoreNewParams struct {
	// Human-readable name for the store. Required; 1–255 characters; no control
	// characters. The mount-path slug under `/mnt/memory/` is derived from this name
	// (lowercased, non-alphanumeric runs collapsed to a hyphen). Names need not be
	// unique within a workspace.
	Name string `json:"name" api:"required"`
	// Free-text description of what the store contains, up to 1024 characters.
	// Included in the agent's system prompt when the store is attached, so word it to
	// be useful to the agent.
	Description param.Opt[string] `json:"description,omitzero"`
	WorkspaceID param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`
	// Arbitrary key-value tags for your own bookkeeping (such as the end user a store
	// belongs to). Up to 16 pairs; keys 1–64 characters; values up to 512 characters.
	// Not visible to the agent.
	Metadata map[string]string `json:"metadata,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

func (r MemoryStoreNewParams) MarshalJSON() (data []byte, err error) {
	type shadow MemoryStoreNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *MemoryStoreNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type MemoryStoreGetParams struct {
	WorkspaceID param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

type MemoryStoreUpdateParams struct {
	// New description for the store, up to 1024 characters. Pass an empty string to
	// clear it.
	Description param.Opt[string] `json:"description,omitzero"`
	// New human-readable name for the store. 1–255 characters; no control characters.
	// Renaming changes the slug used for the store's `mount_path` in sessions created
	// after the update.
	Name        param.Opt[string] `json:"name,omitzero"`
	WorkspaceID param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`
	// Metadata patch. Set a key to a string to upsert it, or to null to delete it.
	// Omit the field to preserve. The stored bag is limited to 16 keys (up to 64 chars
	// each) with values up to 512 chars.
	Metadata map[string]any `json:"metadata,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

func (r MemoryStoreUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow MemoryStoreUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *MemoryStoreUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type MemoryStoreListParams struct {
	Name param.Opt[string] `query:"name,omitzero" json:"-"`

	// Return only stores whose `created_at` is at or after this time (inclusive). Sent
	// on the wire as `created_at[gte]`.
	CreatedAtGte param.Opt[time.Time] `query:"created_at[gte],omitzero" format:"date-time" json:"-"`
	// Return only stores whose `created_at` is at or before this time (inclusive).
	// Sent on the wire as `created_at[lte]`.
	CreatedAtLte param.Opt[time.Time] `query:"created_at[lte],omitzero" format:"date-time" json:"-"`
	// When `true`, archived stores are included in the results. Defaults to `false`
	// (archived stores are excluded).
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// Maximum number of stores to return per page. Must be between 1 and 100. Defaults
	// to 20 when omitted.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque pagination cursor (a `page_...` value). Pass the `next_page` value from a
	// previous response to fetch the next page; omit for the first page.
	Page        param.Opt[string] `query:"page,omitzero" json:"-"`
	WorkspaceID param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [MemoryStoreListParams]'s query parameters as
// `url.Values`.
func (r MemoryStoreListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type MemoryStoreDeleteParams struct {
	WorkspaceID param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

type MemoryStoreArchiveParams struct {
	WorkspaceID param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}
