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
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/constant"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/pagination"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
)

// MemoryStoreMemoryVersionService contains methods and other services that
// help with interacting with the Qoder Cloud Agents API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewMemoryStoreMemoryVersionService] method instead.
type MemoryStoreMemoryVersionService struct {
	Options []option.RequestOption
}

// NewMemoryStoreMemoryVersionService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewMemoryStoreMemoryVersionService(opts ...option.RequestOption) (r MemoryStoreMemoryVersionService) {
	r = MemoryStoreMemoryVersionService{}
	r.Options = opts
	return
}

// Retrieve a memory version
func (r *MemoryStoreMemoryVersionService) Get(ctx context.Context, memoryVersionID string, params MemoryStoreMemoryVersionGetParams, opts ...option.RequestOption) (res *ManagedAgentsMemoryVersion, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if params.MemoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	if memoryVersionID == "" {
		err = errors.New("missing required memory_version_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("memory_stores/%s/memory_versions/%s", url.PathEscape(params.MemoryStoreID), url.PathEscape(memoryVersionID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

// List memory versions
func (r *MemoryStoreMemoryVersionService) List(ctx context.Context, memoryStoreID string, params MemoryStoreMemoryVersionListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ManagedAgentsMemoryVersion], err error) {
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
	path := fmt.Sprintf("memory_stores/%s/memory_versions", url.PathEscape(memoryStoreID))
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

// List memory versions
func (r *MemoryStoreMemoryVersionService) ListAutoPaging(ctx context.Context, memoryStoreID string, params MemoryStoreMemoryVersionListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ManagedAgentsMemoryVersion] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, memoryStoreID, params, opts...))
}

// Redact a memory version
func (r *MemoryStoreMemoryVersionService) Redact(ctx context.Context, memoryVersionID string, params MemoryStoreMemoryVersionRedactParams, opts ...option.RequestOption) (res *ManagedAgentsMemoryVersion, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if params.MemoryStoreID == "" {
		err = errors.New("missing required memory_store_id parameter")
		return nil, err
	}
	if memoryVersionID == "" {
		err = errors.New("missing required memory_version_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("memory_stores/%s/memory_versions/%s/redact", url.PathEscape(params.MemoryStoreID), url.PathEscape(memoryVersionID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// ManagedAgentsActorUnion contains all possible properties and values from
// [ManagedAgentsSessionActor], [ManagedAgentsAPIActor],
// [ManagedAgentsUserActor], [ManagedAgentsServiceAccountActor].
//
// Use the [ManagedAgentsActorUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsActorUnion struct {
	// This field is from variant [ManagedAgentsSessionActor].
	SessionID string `json:"session_id"`
	// Any of "session_actor", "api_actor", "user_actor", "service_account_actor".
	Type string `json:"type"`
	// This field is from variant [ManagedAgentsAPIActor].
	APIKeyID string `json:"api_key_id"`
	// This field is from variant [ManagedAgentsUserActor].
	UserID string `json:"user_id"`
	// This field is from variant [ManagedAgentsServiceAccountActor].
	ServiceAccountID string `json:"service_account_id"`
	JSON             struct {
		SessionID        respjson.Field
		Type             respjson.Field
		APIKeyID         respjson.Field
		UserID           respjson.Field
		ServiceAccountID respjson.Field
		raw              string
	} `json:"-"`
}

// anyManagedAgentsActor is implemented by each variant of
// [ManagedAgentsActorUnion] to add type safety for the return type of
// [ManagedAgentsActorUnion.AsAny]
type anyManagedAgentsActor interface {
	implManagedAgentsActorUnion()
}

func (ManagedAgentsSessionActor) implManagedAgentsActorUnion()        {}
func (ManagedAgentsAPIActor) implManagedAgentsActorUnion()            {}
func (ManagedAgentsUserActor) implManagedAgentsActorUnion()           {}
func (ManagedAgentsServiceAccountActor) implManagedAgentsActorUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsActorUnion.AsAny().(type) {
//	case qoder.ManagedAgentsSessionActor:
//	case qoder.ManagedAgentsAPIActor:
//	case qoder.ManagedAgentsUserActor:
//	case qoder.ManagedAgentsServiceAccountActor:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsActorUnion) AsAny() anyManagedAgentsActor {
	switch u.Type {
	case "session_actor":
		return u.AsSessionActor()
	case "api_actor":
		return u.AsAPIActor()
	case "user_actor":
		return u.AsUserActor()
	case "service_account_actor":
		return u.AsServiceAccountActor()
	}
	return nil
}

func (u ManagedAgentsActorUnion) AsSessionActor() (v ManagedAgentsSessionActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsActorUnion) AsAPIActor() (v ManagedAgentsAPIActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsActorUnion) AsUserActor() (v ManagedAgentsUserActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsActorUnion) AsServiceAccountActor() (v ManagedAgentsServiceAccountActor) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsActorUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsActorUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Attribution for a write made directly via the public API (outside of any
// session).
type ManagedAgentsAPIActor struct {
	// ID of the API key that performed the write. This identifies the key, not the
	// secret.
	APIKeyID string `json:"api_key_id" api:"required"`
	// Any of "api_actor".
	Type ManagedAgentsAPIActorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		APIKeyID    respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsAPIActor) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsAPIActor) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsAPIActorType string

const (
	ManagedAgentsAPIActorTypeAPIActor ManagedAgentsAPIActorType = "api_actor"
)

// A `memory_version` object: one immutable, attributed row in a memory's
// append-only history. Every non-no-op mutation to a memory produces a new
// version. Versions belong to the store (not the individual memory) and are not
// deleted with the memory; each version is retained for at least the version
// retention period after it was written, unless the store itself is deleted.
// Retrieving a redacted version returns 200 with `content`, `path`,
// `content_size_bytes`, and `content_sha256` set to `null`; branch on
// `redacted_at`, not HTTP status.
type ManagedAgentsMemoryVersion struct {
	// Unique identifier for this version (a `memver_...` value).
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// ID of the memory this version snapshots (a `mem_...` value). Remains valid after
	// the memory is deleted; pass it as `memory_id` to
	// [List memory versions](/en/api/beta/memory_stores/memory_versions/list) to
	// retrieve the memory's retained versions, including the `deleted` row while the
	// lineage is retained.
	MemoryID string `json:"memory_id" api:"required"`
	// ID of the memory store this version belongs to (a `memstore_...` value).
	MemoryStoreID string `json:"memory_store_id" api:"required"`
	// The kind of mutation a `memory_version` records. Every non-no-op mutation to a
	// memory appends exactly one version row with one of these values.
	//
	// Any of "created", "modified", "deleted".
	Operation ManagedAgentsMemoryVersionOperation `json:"operation" api:"required"`
	// Any of "memory_version".
	Type ManagedAgentsMemoryVersionType `json:"type" api:"required"`
	// The memory's UTF-8 text content as of this version. `null` when `view=basic`,
	// when `operation` is `deleted`, or when `redacted_at` is set.
	Content string `json:"content" api:"nullable"`
	// Lowercase hex SHA-256 digest of `content` as of this version (64 characters).
	// `null` when `redacted_at` is set or `operation` is `deleted`. Populated
	// regardless of `view` otherwise.
	ContentSha256 string `json:"content_sha256" api:"nullable"`
	// Size of `content` in bytes as of this version. `null` when `redacted_at` is set
	// or `operation` is `deleted`. Populated regardless of `view` otherwise.
	ContentSizeBytes int64 `json:"content_size_bytes" api:"nullable"`
	// Identifies who performed a write or redact operation. Captured at write time on
	// the `memory_version` row. The API key that created a session is not recorded on
	// agent writes; attribution answers who made the write, not who is ultimately
	// responsible. Look up session provenance separately via the
	// [Sessions API](/en/api/beta/sessions/retrieve).
	CreatedBy ManagedAgentsActorUnion `json:"created_by"`
	// The memory's path at the time of this write. `null` if and only if `redacted_at`
	// is set.
	Path string `json:"path" api:"nullable"`
	// A timestamp in RFC 3339 format
	RedactedAt time.Time `json:"redacted_at" api:"nullable" format:"date-time"`
	// Identifies who performed a write or redact operation. Captured at write time on
	// the `memory_version` row. The API key that created a session is not recorded on
	// agent writes; attribution answers who made the write, not who is ultimately
	// responsible. Look up session provenance separately via the
	// [Sessions API](/en/api/beta/sessions/retrieve).
	RedactedBy ManagedAgentsActorUnion `json:"redacted_by"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID               respjson.Field
		CreatedAt        respjson.Field
		MemoryID         respjson.Field
		MemoryStoreID    respjson.Field
		Operation        respjson.Field
		Type             respjson.Field
		Content          respjson.Field
		ContentSha256    respjson.Field
		ContentSizeBytes respjson.Field
		CreatedBy        respjson.Field
		Path             respjson.Field
		RedactedAt       respjson.Field
		RedactedBy       respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsMemoryVersion) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsMemoryVersion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsMemoryVersionType string

const (
	ManagedAgentsMemoryVersionTypeMemoryVersion ManagedAgentsMemoryVersionType = "memory_version"
)

// The kind of mutation a `memory_version` records. Every non-no-op mutation to a
// memory appends exactly one version row with one of these values.
type ManagedAgentsMemoryVersionOperation string

const (
	ManagedAgentsMemoryVersionOperationCreated  ManagedAgentsMemoryVersionOperation = "created"
	ManagedAgentsMemoryVersionOperationModified ManagedAgentsMemoryVersionOperation = "modified"
	ManagedAgentsMemoryVersionOperationDeleted  ManagedAgentsMemoryVersionOperation = "deleted"
)

// Attribution for a write made by a workload authenticated as a service account,
// for example via Workload Identity Federation.
type ManagedAgentsServiceAccountActor struct {
	// ID of the service account that performed the write (a `svac_...` value).
	ServiceAccountID string                       `json:"service_account_id" api:"required"`
	Type             constant.ServiceAccountActor `json:"type" default:"service_account_actor"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ServiceAccountID respjson.Field
		Type             respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsServiceAccountActor) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsServiceAccountActor) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Attribution for a write made by an agent during a session, through the mounted
// filesystem at `/mnt/memory/`.
type ManagedAgentsSessionActor struct {
	// ID of the session that performed the write (a `sesn_...` value). Look up the
	// session via [Retrieve a session](/en/api/beta/sessions/retrieve) for further
	// provenance.
	SessionID string `json:"session_id" api:"required"`
	// Any of "session_actor".
	Type ManagedAgentsSessionActorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		SessionID   respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionActor) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionActor) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSessionActorType string

const (
	ManagedAgentsSessionActorTypeSessionActor ManagedAgentsSessionActorType = "session_actor"
)

// Attribution for a write made by a human user through the Qoder Console.
type ManagedAgentsUserActor struct {
	// Any of "user_actor".
	Type ManagedAgentsUserActorType `json:"type" api:"required"`
	// ID of the user who performed the write (a `user_...` value).
	UserID string `json:"user_id" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		UserID      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsUserActor) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsUserActor) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsUserActorType string

const (
	ManagedAgentsUserActorTypeUserActor ManagedAgentsUserActorType = "user_actor"
)

type MemoryStoreMemoryVersionGetParams struct {
	MemoryStoreID string            `path:"memory_store_id" api:"required" json:"-"`
	WorkspaceID   param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`
	// Query parameter for view
	//
	// Any of "basic", "full".
	View ManagedAgentsMemoryView `query:"view,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [MemoryStoreMemoryVersionGetParams]'s query parameters
// as `url.Values`.
func (r MemoryStoreMemoryVersionGetParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type MemoryStoreMemoryVersionListParams struct {
	// Query parameter for api_key_id
	APIKeyID param.Opt[string] `query:"api_key_id,omitzero" json:"-"`
	// Return versions created at or after this time (inclusive).
	CreatedAtGte param.Opt[time.Time] `query:"created_at[gte],omitzero" format:"date-time" json:"-"`
	// Return versions created at or before this time (inclusive).
	CreatedAtLte param.Opt[time.Time] `query:"created_at[lte],omitzero" format:"date-time" json:"-"`
	// Query parameter for limit
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Query parameter for memory_id
	MemoryID param.Opt[string] `query:"memory_id,omitzero" json:"-"`
	// Query parameter for page
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Query parameter for service_account_id
	ServiceAccountID param.Opt[string] `query:"service_account_id,omitzero" json:"-"`
	// Query parameter for session_id
	SessionID   param.Opt[string] `query:"session_id,omitzero" json:"-"`
	WorkspaceID param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`
	// Query parameter for operation
	//
	// Any of "created", "modified", "deleted".
	Operation ManagedAgentsMemoryVersionOperation `query:"operation,omitzero" json:"-"`
	// Query parameter for view
	//
	// Any of "basic", "full".
	View ManagedAgentsMemoryView `query:"view,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [MemoryStoreMemoryVersionListParams]'s query parameters
// as `url.Values`.
func (r MemoryStoreMemoryVersionListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type MemoryStoreMemoryVersionRedactParams struct {
	MemoryStoreID string            `path:"memory_store_id" api:"required" json:"-"`
	WorkspaceID   param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}
