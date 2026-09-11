// Qoder managed API definitions.
package managed

import (
	"context"
	"encoding/json"
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

// MemoryStoreMemoryService contains methods and other services that help with
// interacting with the Qoder Cloud Agents API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewMemoryStoreMemoryService] method instead.
type MemoryStoreMemoryService struct {
	Options []option.RequestOption
}

// NewMemoryStoreMemoryService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewMemoryStoreMemoryService(opts ...option.RequestOption) (r MemoryStoreMemoryService) {
	r = MemoryStoreMemoryService{}
	r.Options = opts
	return
}

// Create a memory
func (r *MemoryStoreMemoryService) New(ctx context.Context, memoryStoreID string, params MemoryStoreMemoryNewParams, opts ...option.RequestOption) (res *ManagedAgentsMemory, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if memoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("memory_stores/%s/memories", url.PathEscape(memoryStoreID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Retrieve a memory
func (r *MemoryStoreMemoryService) Get(ctx context.Context, memoryID string, params MemoryStoreMemoryGetParams, opts ...option.RequestOption) (res *ManagedAgentsMemory, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if params.MemoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	if memoryID == "" {
		err = errors.New("missing required memory_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("memory_stores/%s/memories/%s", url.PathEscape(params.MemoryStoreID), url.PathEscape(memoryID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

// Update a memory
func (r *MemoryStoreMemoryService) Update(ctx context.Context, memoryID string, params MemoryStoreMemoryUpdateParams, opts ...option.RequestOption) (res *ManagedAgentsMemory, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if params.MemoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	if memoryID == "" {
		err = errors.New("missing required memory_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("memory_stores/%s/memories/%s", url.PathEscape(params.MemoryStoreID), url.PathEscape(memoryID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List memories
func (r *MemoryStoreMemoryService) List(ctx context.Context, memoryStoreID string, params MemoryStoreMemoryListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ManagedAgentsMemoryListItemUnion], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if memoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("memory_stores/%s/memories", url.PathEscape(memoryStoreID))
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

// List memories
func (r *MemoryStoreMemoryService) ListAutoPaging(ctx context.Context, memoryStoreID string, params MemoryStoreMemoryListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ManagedAgentsMemoryListItemUnion] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, memoryStoreID, params, opts...))
}

// Delete a memory
func (r *MemoryStoreMemoryService) Delete(ctx context.Context, memoryID string, params MemoryStoreMemoryDeleteParams, opts ...option.RequestOption) (res *ManagedAgentsDeletedMemory, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if params.MemoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	if memoryID == "" {
		err = errors.New("missing required memory_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("memory_stores/%s/memories/%s", url.PathEscape(params.MemoryStoreID), url.PathEscape(memoryID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, params, &res, opts...)
	return res, err
}

// Tombstone returned by
// [Delete a memory](/en/api/beta/memory_stores/memories/delete). Deleting a memory
// does not erase its version history: its versions remain listable via
// [List memory versions](/en/api/beta/memory_stores/memory_versions/list) while
// they are retained (each version is kept for at least the version retention
// period after it was written, unless the store itself is deleted).
type ManagedAgentsDeletedMemory struct {
	// ID of the deleted memory (a `mem_...` value).
	ID string `json:"id" api:"required"`
	// Any of "memory_deleted".
	Type ManagedAgentsDeletedMemoryType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsDeletedMemory) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsDeletedMemory) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsDeletedMemoryType string

const (
	ManagedAgentsDeletedMemoryTypeMemoryDeleted ManagedAgentsDeletedMemoryType = "memory_deleted"
)

// A `memory` object: a single text document at a hierarchical path inside a memory
// store. The `content` field is populated when `view=full` and `null` when
// `view=basic`; the `content_size_bytes` and `content_sha256` fields are always
// populated so sync clients can diff without fetching content. Memories are
// addressed by their `mem_...` ID; the path is the create key and can be changed
// via update.
type ManagedAgentsMemory struct {
	Metadata map[string]string `json:"metadata"`

	// Unique identifier for this memory (a `mem_...` value). Stable across renames;
	// use this ID, not the path, to read, update, or delete the memory.
	ID string `json:"id" api:"required"`
	// Lowercase hex SHA-256 digest of the UTF-8 `content` bytes (64 characters). The
	// server applies no normalization, so clients can compute the same hash locally
	// for staleness checks and as the value for a `content_sha256` precondition on
	// update. Always populated, regardless of `view`.
	ContentSha256 string `json:"content_sha256" api:"required"`
	// Size of `content` in bytes (the UTF-8 plaintext length). Always populated,
	// regardless of `view`.
	ContentSizeBytes int64 `json:"content_size_bytes" api:"required"`
	// A timestamp in RFC 3339 format
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// ID of the memory store this memory belongs to (a `memstore_...` value).
	MemoryStoreID string `json:"memory_store_id" api:"required"`
	// ID of the `memory_version` representing this memory's current content (a
	// `memver_...` value). This is the authoritative head pointer; `memory_version`
	// objects do not carry an `is_latest` flag, so compare against this field instead.
	// Enumerate the history via
	// [List memory versions](/en/api/beta/memory_stores/memory_versions/list).
	MemoryVersionID string `json:"memory_version_id" api:"required"`
	// Hierarchical path of the memory within the store, e.g. `/projects/foo/notes.md`.
	// Always starts with `/`. Paths are case-sensitive and unique within a store.
	// Maximum 1,024 bytes.
	Path string `json:"path" api:"required"`
	// Any of "memory".
	Type ManagedAgentsMemoryType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// The memory's UTF-8 text content. Populated when `view=full`; `null` when
	// `view=basic`. Maximum 100 kB (102,400 bytes).
	Content string `json:"content" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Metadata respjson.Field

		ID               respjson.Field
		ContentSha256    respjson.Field
		ContentSizeBytes respjson.Field
		CreatedAt        respjson.Field
		MemoryStoreID    respjson.Field
		MemoryVersionID  respjson.Field
		Path             respjson.Field
		Type             respjson.Field
		UpdatedAt        respjson.Field
		Content          respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsMemory) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsMemory) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsMemoryType string

const (
	ManagedAgentsMemoryTypeMemory ManagedAgentsMemoryType = "memory"
)

// ManagedAgentsMemoryListItemUnion contains all possible properties and values
// from [ManagedAgentsMemory], [ManagedAgentsMemoryPrefix].
//
// Use the [ManagedAgentsMemoryListItemUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsMemoryListItemUnion struct {
	// This field is from variant [ManagedAgentsMemory].
	ID string `json:"id"`
	// This field is from variant [ManagedAgentsMemory].
	ContentSha256 string `json:"content_sha256"`
	// This field is from variant [ManagedAgentsMemory].
	ContentSizeBytes int64 `json:"content_size_bytes"`
	// This field is from variant [ManagedAgentsMemory].
	CreatedAt time.Time `json:"created_at"`
	// This field is from variant [ManagedAgentsMemory].
	MemoryStoreID string `json:"memory_store_id"`
	// This field is from variant [ManagedAgentsMemory].
	MemoryVersionID string `json:"memory_version_id"`
	Path            string `json:"path"`
	// Any of "memory", "memory_prefix".
	Type string `json:"type"`
	// This field is from variant [ManagedAgentsMemory].
	UpdatedAt time.Time `json:"updated_at"`
	// This field is from variant [ManagedAgentsMemory].
	Content string `json:"content"`
	JSON    struct {
		ID               respjson.Field
		ContentSha256    respjson.Field
		ContentSizeBytes respjson.Field
		CreatedAt        respjson.Field
		MemoryStoreID    respjson.Field
		MemoryVersionID  respjson.Field
		Path             respjson.Field
		Type             respjson.Field
		UpdatedAt        respjson.Field
		Content          respjson.Field
		raw              string
	} `json:"-"`
}

// anyManagedAgentsMemoryListItem is implemented by each variant of
// [ManagedAgentsMemoryListItemUnion] to add type safety for the return type of
// [ManagedAgentsMemoryListItemUnion.AsAny]
type anyManagedAgentsMemoryListItem interface {
	implManagedAgentsMemoryListItemUnion()
}

func (ManagedAgentsMemory) implManagedAgentsMemoryListItemUnion()       {}
func (ManagedAgentsMemoryPrefix) implManagedAgentsMemoryListItemUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsMemoryListItemUnion.AsAny().(type) {
//	case qoder.ManagedAgentsMemory:
//	case qoder.ManagedAgentsMemoryPrefix:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsMemoryListItemUnion) AsAny() anyManagedAgentsMemoryListItem {
	switch u.Type {
	case "memory":
		return u.AsMemory()
	case "memory_prefix":
		return u.AsMemoryPrefix()
	}
	return nil
}

func (u ManagedAgentsMemoryListItemUnion) AsMemory() (v ManagedAgentsMemory) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsMemoryListItemUnion) AsMemoryPrefix() (v ManagedAgentsMemoryPrefix) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsMemoryListItemUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsMemoryListItemUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A rolled-up directory marker returned by
// [List memories](/en/api/beta/memory_stores/memories/list) when `depth` is set.
// Indicates that one or more memories exist deeper than the requested depth under
// this prefix. This is a list-time rollup, not a stored resource; it has no ID and
// no lifecycle. Each prefix counts toward the page `limit` and interleaves with
// `memory` items in path order.
type ManagedAgentsMemoryPrefix struct {
	// The rolled-up path prefix, including a trailing `/` (e.g. `/projects/foo/`).
	// Pass this value as `path_prefix` on a subsequent list call to drill into the
	// directory.
	Path string `json:"path" api:"required"`
	// Any of "memory_prefix".
	Type ManagedAgentsMemoryPrefixType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Path        respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsMemoryPrefix) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsMemoryPrefix) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsMemoryPrefixType string

const (
	ManagedAgentsMemoryPrefixTypeMemoryPrefix ManagedAgentsMemoryPrefixType = "memory_prefix"
)

// Selects which projection of a `memory` or `memory_version` the server returns.
// `basic` returns the object with `content` set to `null`; `full` populates
// `content`. When omitted, the default is endpoint-specific: retrieve operations
// default to `full`; list, create, and update operations default to `basic`.
// Listing with `view=full` caps `limit` at 20.
type ManagedAgentsMemoryView string

const (
	ManagedAgentsMemoryViewBasic ManagedAgentsMemoryView = "basic"
	ManagedAgentsMemoryViewFull  ManagedAgentsMemoryView = "full"
)

// Optimistic-concurrency precondition: the update applies only if the memory's
// stored `content_sha256` equals the supplied value. On mismatch, the request
// returns `memory_precondition_failed_error` (HTTP 409); re-read the memory and
// retry against the fresh state. If the precondition fails but the stored state
// already exactly matches the requested `content` and `path`, the server returns
// 200 instead of 409.
//
// The property Type is required.
type ManagedAgentsPreconditionParam struct {
	// Any of "content_sha256".
	Type ManagedAgentsPreconditionType `json:"type,omitzero" api:"required"`
	// Expected `content_sha256` of the stored memory (64 lowercase hexadecimal
	// characters). Typically the `content_sha256` returned by a prior read or list
	// call. Because the server applies no content normalization, clients can also
	// compute this locally as the SHA-256 of the UTF-8 content bytes.
	ContentSha256 param.Opt[string] `json:"content_sha256,omitzero"`
	paramObj
}

func (r ManagedAgentsPreconditionParam) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsPreconditionParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsPreconditionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsPreconditionType string

const (
	ManagedAgentsPreconditionTypeContentSha256 ManagedAgentsPreconditionType = "content_sha256"
)

type MemoryStoreMemoryNewParams struct {
	Metadata map[string]string `json:"metadata,omitzero"`

	// UTF-8 text content for the new memory. Maximum 100 kB (102,400 bytes). Required;
	// pass `""` explicitly to create an empty memory.
	Content param.Opt[string] `json:"content,omitzero" api:"required"`
	// Hierarchical path for the new memory, e.g. `/projects/foo/notes.md`. Must start
	// with `/`, contain at least one non-empty segment, and be at most 1,024 bytes.
	// Must not contain empty segments, `.` or `..` segments, control or format
	// characters, or the Unicode line and paragraph separators (U+2028, U+2029), and
	// must be NFC-normalized. Paths are case-sensitive.
	Path string `json:"path" api:"required"`
	// Query parameter for view
	//
	// Any of "basic", "full".
	View ManagedAgentsMemoryView `query:"view,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

func (r MemoryStoreMemoryNewParams) MarshalJSON() (data []byte, err error) {
	type shadow MemoryStoreMemoryNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *MemoryStoreMemoryNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// URLQuery serializes [MemoryStoreMemoryNewParams]'s query parameters as
// `url.Values`.
func (r MemoryStoreMemoryNewParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type MemoryStoreMemoryGetParams struct {
	MemoryStoreID string `path:"memory_store_id" api:"required" json:"-"`
	// Query parameter for view
	//
	// Any of "basic", "full".
	View ManagedAgentsMemoryView `query:"view,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [MemoryStoreMemoryGetParams]'s query parameters as
// `url.Values`.
func (r MemoryStoreMemoryGetParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type MemoryStoreMemoryUpdateParams struct {
	ContentSha256 param.Opt[string] `json:"content_sha256,omitzero"`

	Metadata map[string]any `json:"metadata,omitzero"`

	MemoryStoreID string `path:"memory_store_id" api:"required" json:"-"`
	// New UTF-8 text content for the memory. Maximum 100 kB (102,400 bytes). Omit to
	// leave the content unchanged (e.g., for a rename-only update).
	Content param.Opt[string] `json:"content,omitzero"`
	// New path for the memory (a rename). Must start with `/`, contain at least one
	// non-empty segment, and be at most 1,024 bytes. Must not contain empty segments,
	// `.` or `..` segments, control or format characters, or the Unicode line and
	// paragraph separators (U+2028, U+2029), and must be NFC-normalized. Paths are
	// case-sensitive. The memory's `id` is preserved across renames. Omit to leave the
	// path unchanged.
	Path param.Opt[string] `json:"path,omitzero"`
	// Query parameter for view
	//
	// Any of "basic", "full".
	View ManagedAgentsMemoryView `query:"view,omitzero" json:"-"`
	// Optimistic-concurrency precondition: the update applies only if the memory's
	// stored `content_sha256` equals the supplied value. On mismatch, the request
	// returns `memory_precondition_failed_error` (HTTP 409); re-read the memory and
	// retry against the fresh state. If the precondition fails but the stored state
	// already exactly matches the requested `content` and `path`, the server returns
	// 200 instead of 409.
	Precondition ManagedAgentsPreconditionParam `json:"precondition,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

func (r MemoryStoreMemoryUpdateParams) MarshalJSON() (data []byte, err error) {
	if !param.IsOmitted(r.Precondition) {
		if r.Precondition.Type != ManagedAgentsPreconditionTypeContentSha256 {
			return nil, fmt.Errorf("unsupported memory precondition")
		}
		if !param.IsOmitted(r.ContentSha256) && r.ContentSha256 != r.Precondition.ContentSha256 {
			return nil, fmt.Errorf("conflicting memory content hashes")
		}
		r.ContentSha256 = r.Precondition.ContentSha256
		r.Precondition = ManagedAgentsPreconditionParam{}
	}

	type shadow MemoryStoreMemoryUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *MemoryStoreMemoryUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// URLQuery serializes [MemoryStoreMemoryUpdateParams]'s query parameters as
// `url.Values`.
func (r MemoryStoreMemoryUpdateParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type MemoryStoreMemoryListParams struct {
	// `0` (or omitted) returns all descendants below `path_prefix` (recursive). `1`
	// returns immediate children only; deeper entries roll up as `memory_prefix`
	// items. `depth=1` behaves like `ls`; omitting `depth` behaves like `find`.
	Depth param.Opt[int64] `query:"depth,omitzero" json:"-"`
	// Maximum number of items to return per page. Must be between 1 and 100. Defaults
	// to 20 when omitted. Capped at 20 when `view=full`. Both `memory` and
	// `memory_prefix` items count toward the limit.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque pagination cursor (a `page_...` value). Pass the `next_page` value from a
	// previous response to fetch the next page; omit for the first page.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Optional path prefix filter. Must end with `/` (segment-aligned), e.g.,
	// `/notes/`. This value appears in request URLs. Do not include secrets or
	// personally identifiable information.
	PathPrefix param.Opt[string] `query:"path_prefix,omitzero" json:"-"`
	// Which projection of each `memory` to return. Defaults to `basic` (content
	// omitted). `full` populates `content` on each item and caps `limit` at 20; use
	// this as the bulk-read path for export and sync.
	//
	// Any of "basic", "full".
	View ManagedAgentsMemoryView `query:"view,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [MemoryStoreMemoryListParams]'s query parameters as
// `url.Values`.
func (r MemoryStoreMemoryListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type MemoryStoreMemoryDeleteParams struct {
	MemoryStoreID string `path:"memory_store_id" api:"required" json:"-"`
	// Query parameter for expected_content_sha256
	ExpectedContentSha256 param.Opt[string] `query:"expected_content_sha256,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [MemoryStoreMemoryDeleteParams]'s query parameters as
// `url.Values`.
func (r MemoryStoreMemoryDeleteParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
