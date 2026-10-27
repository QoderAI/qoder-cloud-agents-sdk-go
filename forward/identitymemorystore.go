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

// IdentityMemoryStoreService provides Forward IdentityMemoryStore operations.
type IdentityMemoryStoreService struct {
	Options []option.RequestOption
}

func NewIdentityMemoryStoreService(opts ...option.RequestOption) IdentityMemoryStoreService {
	return IdentityMemoryStoreService{Options: slices.Clone(opts)}
}

// List Memory Store mounts on an Identity
func (r *IdentityMemoryStoreService) List(ctx context.Context, identityID string, templateID string, params IdentityMemoryStoreListParams, opts ...option.RequestOption) (res *pagination.Page[MemoryStoreMount], err error) {
	if identityID == "" {
		return nil, fmt.Errorf("missing required identity_id parameter")
	}
	if templateID == "" {
		return nil, fmt.Errorf("missing required template_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("identities/%s/templates/%s/memory_stores", url.PathEscape(identityID), url.PathEscape(templateID))
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

// Mount Memory Store on an Identity
func (r *IdentityMemoryStoreService) Mount(ctx context.Context, identityID string, templateID string, params IdentityMemoryStoreMountParams, opts ...option.RequestOption) (res *MemoryStoreMount, err error) {
	if identityID == "" {
		return nil, fmt.Errorf("missing required identity_id parameter")
	}
	if templateID == "" {
		return nil, fmt.Errorf("missing required template_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("identities/%s/templates/%s/memory_stores", url.PathEscape(identityID), url.PathEscape(templateID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type IdentityMemoryStoreMountParams struct {
	Access param.Opt[string] `json:"access,omitzero"`
	// ID of the Memory Store to mount (`memstore_...`). It must be an active Store
	// visible to the caller.
	MemoryStoreID string `json:"memory_store_id" api:"required"`
	paramObj
}

func (r IdentityMemoryStoreMountParams) MarshalJSON() ([]byte, error) {
	type shadow IdentityMemoryStoreMountParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *IdentityMemoryStoreMountParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Unmount Memory Store from an Identity
func (r *IdentityMemoryStoreService) Detach(ctx context.Context, identityID string, templateID string, memoryStoreID string, opts ...option.RequestOption) (res *DeletedMemoryStoreMount, err error) {
	if identityID == "" {
		return nil, fmt.Errorf("missing required identity_id parameter")
	}
	if templateID == "" {
		return nil, fmt.Errorf("missing required template_id parameter")
	}
	if memoryStoreID == "" {
		return nil, fmt.Errorf("missing required memory_store_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("identities/%s/templates/%s/memory_stores/%s", url.PathEscape(identityID), url.PathEscape(templateID), url.PathEscape(memoryStoreID))
	err = convention.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

type DeletedMemoryStoreMount struct {
	// ID of the unmounted Memory Store.
	ID string `json:"id"`
	// Always `memory_store_binding_deleted`, as distinct from the `memory_store_deleted`
	// returned when deleting a Store: this endpoint only removes the mount, and the Store
	// itself continues to exist.
	Type string `json:"type"`
	// Whether the mount was removed.
	Deleted bool `json:"deleted"`
	JSON    struct {
		ID          respjson.Field
		Type        respjson.Field
		Deleted     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r DeletedMemoryStoreMount) RawJSON() string { return r.JSON.raw }
func (r *DeletedMemoryStoreMount) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type MemoryStoreMount struct {
	MemoryStoreID string    `json:"memory_store_id"`
	IdentityID    string    `json:"identity_id"`
	TemplateID    string    `json:"template_id"`
	Access        string    `json:"access"`
	SystemManaged bool      `json:"system_managed"`
	Name          string    `json:"name"`
	Status        string    `json:"status"`
	EntryCount    int64     `json:"entry_count"`
	CreatedAt     time.Time `json:"created_at" format:"date-time"`
	JSON          struct {
		MemoryStoreID respjson.Field
		IdentityID    respjson.Field
		TemplateID    respjson.Field
		Access        respjson.Field
		SystemManaged respjson.Field
		Name          respjson.Field
		Status        respjson.Field
		EntryCount    respjson.Field
		CreatedAt     respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

func (r MemoryStoreMount) RawJSON() string                  { return r.JSON.raw }
func (r *MemoryStoreMount) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type IdentityMemoryStoreListResponse struct {
	FirstID string             `json:"first_id"`
	LastID  string             `json:"last_id"`
	Data    []MemoryStoreMount `json:"data"`
	HasMore bool               `json:"has_more"`
	JSON    struct {
		FirstID     respjson.Field
		LastID      respjson.Field
		Data        respjson.Field
		HasMore     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r IdentityMemoryStoreListResponse) RawJSON() string { return r.JSON.raw }
func (r *IdentityMemoryStoreListResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type IdentityMemoryStoreListParams struct {
	Limit    param.Opt[int64]  `query:"limit,omitzero" json:"-"`
	AfterID  param.Opt[string] `query:"after_id,omitzero" json:"-"`
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	paramObj
}

func (r IdentityMemoryStoreListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{})
}
func (r *IdentityMemoryStoreService) ListAutoPaging(ctx context.Context, identityID, templateID string, params IdentityMemoryStoreListParams, opts ...option.RequestOption) *pagination.PageAutoPager[MemoryStoreMount] {
	return pagination.NewPageAutoPager(r.List(ctx, identityID, templateID, params, opts...))
}
