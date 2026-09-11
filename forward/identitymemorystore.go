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
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
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

// 列出 Identity 上的 Memory Store 挂载.
func (r *IdentityMemoryStoreService) List(ctx context.Context, identityID string, templateID string, opts ...option.RequestOption) (res *IdentityMemoryStoreListResponse, err error) {
	if identityID == "" {
		return nil, fmt.Errorf("missing required identity_id parameter")
	}
	if templateID == "" {
		return nil, fmt.Errorf("missing required template_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("identities/%s/templates/%s/memory_stores", url.PathEscape(identityID), url.PathEscape(templateID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// 挂载 Memory Store 到 Identity.
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
	// 要挂载的 Memory Store ID（`memstore_...`）。必须是当前调用方可见的 active Store。
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

// 解绑 Identity 上的 Memory Store.
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
	// 被解绑的 Memory Store ID。
	ID string `json:"id"`
	// 固定为 `memory_store_binding_deleted`。与删除 Store 的 `memory_store_deleted` 区分：本接口只解除挂载关系，Store 本体仍然存在。
	Type string `json:"type"`
	// 挂载关系是否已解除。
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
	Data    []MemoryStoreMount `json:"data"`
	HasMore bool               `json:"has_more"`
	JSON    struct {
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
