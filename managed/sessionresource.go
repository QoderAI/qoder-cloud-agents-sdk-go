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

// SessionResourceService contains methods and other services that help with
// interacting with the Qoder Cloud Agents API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSessionResourceService] method instead.
type SessionResourceService struct {
	Options []option.RequestOption
}

// NewSessionResourceService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewSessionResourceService(opts ...option.RequestOption) (r SessionResourceService) {
	r = SessionResourceService{}
	r.Options = opts
	return
}

// Get Session Resource
func (r *SessionResourceService) Get(ctx context.Context, resourceID string, params SessionResourceGetParams, opts ...option.RequestOption) (res *SessionResourceGetResponseUnion, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if params.SessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	if resourceID == "" {
		err = errors.New("missing required resource_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("sessions/%s/resources/%s", url.PathEscape(params.SessionID), url.PathEscape(resourceID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update Session Resource
func (r *SessionResourceService) Update(ctx context.Context, resourceID string, params SessionResourceUpdateParams, opts ...option.RequestOption) (res *SessionResourceUpdateResponseUnion, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if params.SessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	if resourceID == "" {
		err = errors.New("missing required resource_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("sessions/%s/resources/%s", url.PathEscape(params.SessionID), url.PathEscape(resourceID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List Session Resources
func (r *SessionResourceService) List(ctx context.Context, sessionID string, params SessionResourceListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ManagedAgentsSessionResourceUnion], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("sessions/%s/resources", url.PathEscape(sessionID))
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

// List Session Resources
func (r *SessionResourceService) ListAutoPaging(ctx context.Context, sessionID string, params SessionResourceListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ManagedAgentsSessionResourceUnion] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, sessionID, params, opts...))
}

// Delete Session Resource
func (r *SessionResourceService) Delete(ctx context.Context, resourceID string, params SessionResourceDeleteParams, opts ...option.RequestOption) (res *ManagedAgentsDeleteSessionResource, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if params.SessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	if resourceID == "" {
		err = errors.New("missing required resource_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("sessions/%s/resources/%s", url.PathEscape(params.SessionID), url.PathEscape(resourceID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Add Session Resource
func (r *SessionResourceService) Add(ctx context.Context, sessionID string, params SessionResourceAddParams, opts ...option.RequestOption) (res *ManagedAgentsFileResource, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("sessions/%s/resources", url.PathEscape(sessionID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Confirmation of resource deletion.
type ManagedAgentsDeleteSessionResource struct {
	ID string `json:"id" api:"required"`
	// Any of "session_resource_deleted".
	Type ManagedAgentsDeleteSessionResourceType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsDeleteSessionResource) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsDeleteSessionResource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsDeleteSessionResourceType string

const (
	ManagedAgentsDeleteSessionResourceTypeSessionResourceDeleted ManagedAgentsDeleteSessionResourceType = "session_resource_deleted"
)

type ManagedAgentsFileResource struct {
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	FileID    string    `json:"file_id" api:"required"`
	MountPath string    `json:"mount_path" api:"required"`
	// Any of "file".
	Type ManagedAgentsFileResourceType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		FileID      respjson.Field
		MountPath   respjson.Field
		Type        respjson.Field
		UpdatedAt   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsFileResource) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsFileResource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsFileResourceType string

const (
	ManagedAgentsFileResourceTypeFile ManagedAgentsFileResourceType = "file"
)

type ManagedAgentsGitHubRepositoryResource struct {
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	MountPath string    `json:"mount_path" api:"required"`
	// Any of "github_repository".
	Type ManagedAgentsGitHubRepositoryResourceType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	UpdatedAt time.Time                                          `json:"updated_at" api:"required" format:"date-time"`
	URL       string                                             `json:"url" api:"required"`
	Checkout  ManagedAgentsGitHubRepositoryResourceCheckoutUnion `json:"checkout" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		CreatedAt   respjson.Field
		MountPath   respjson.Field
		Type        respjson.Field
		UpdatedAt   respjson.Field
		URL         respjson.Field
		Checkout    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsGitHubRepositoryResource) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsGitHubRepositoryResource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsGitHubRepositoryResourceType string

const (
	ManagedAgentsGitHubRepositoryResourceTypeGitHubRepository ManagedAgentsGitHubRepositoryResourceType = "github_repository"
)

// ManagedAgentsGitHubRepositoryResourceCheckoutUnion contains all possible
// properties and values from [ManagedAgentsBranchCheckout],
// [ManagedAgentsCommitCheckout].
//
// Use the [ManagedAgentsGitHubRepositoryResourceCheckoutUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsGitHubRepositoryResourceCheckoutUnion struct {
	// This field is from variant [ManagedAgentsBranchCheckout].
	Name string `json:"name"`
	// Any of "branch", "commit".
	Type string `json:"type"`
	// This field is from variant [ManagedAgentsCommitCheckout].
	Sha  string `json:"sha"`
	JSON struct {
		Name respjson.Field
		Type respjson.Field
		Sha  respjson.Field
		raw  string
	} `json:"-"`
}

// anyManagedAgentsGitHubRepositoryResourceCheckout is implemented by each
// variant of [ManagedAgentsGitHubRepositoryResourceCheckoutUnion] to add type
// safety for the return type of
// [ManagedAgentsGitHubRepositoryResourceCheckoutUnion.AsAny]
type anyManagedAgentsGitHubRepositoryResourceCheckout interface {
	implManagedAgentsGitHubRepositoryResourceCheckoutUnion()
}

func (ManagedAgentsBranchCheckout) implManagedAgentsGitHubRepositoryResourceCheckoutUnion() {}
func (ManagedAgentsCommitCheckout) implManagedAgentsGitHubRepositoryResourceCheckoutUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsGitHubRepositoryResourceCheckoutUnion.AsAny().(type) {
//	case qoder.ManagedAgentsBranchCheckout:
//	case qoder.ManagedAgentsCommitCheckout:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsGitHubRepositoryResourceCheckoutUnion) AsAny() anyManagedAgentsGitHubRepositoryResourceCheckout {
	switch u.Type {
	case "branch":
		return u.AsBranch()
	case "commit":
		return u.AsCommit()
	}
	return nil
}

func (u ManagedAgentsGitHubRepositoryResourceCheckoutUnion) AsBranch() (v ManagedAgentsBranchCheckout) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsGitHubRepositoryResourceCheckoutUnion) AsCommit() (v ManagedAgentsCommitCheckout) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsGitHubRepositoryResourceCheckoutUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsGitHubRepositoryResourceCheckoutUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A memory store attached to an agent session.
type ManagedAgentsMemoryStoreResource struct {
	// The memory store ID (memstore\_...). Must belong to the caller's organization
	// and workspace.
	MemoryStoreID string `json:"memory_store_id" api:"required"`
	// Any of "memory_store".
	Type ManagedAgentsMemoryStoreResourceType `json:"type" api:"required"`
	// Access mode for an attached memory store.
	//
	// Any of "read_write", "read_only".
	Access ManagedAgentsMemoryStoreResourceAccess `json:"access" api:"nullable"`
	// Description of the memory store, snapshotted at attach time. Rendered into the
	// agent's system prompt. Empty string when the store has no description.
	Description string `json:"description"`
	// Per-attachment guidance for the agent on how to use this store. Rendered into
	// the memory section of the system prompt. Max 4096 chars.
	Instructions string `json:"instructions" api:"nullable"`
	// Filesystem path where the store is mounted in the session container, e.g.
	// /mnt/memory/user-preferences. Derived from the store's name. Output-only.
	MountPath string `json:"mount_path" api:"nullable"`
	// Display name of the memory store, snapshotted at attach time. Later edits to the
	// store's name do not propagate to this resource.
	Name string `json:"name" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		MemoryStoreID respjson.Field
		Type          respjson.Field
		Access        respjson.Field
		Description   respjson.Field
		Instructions  respjson.Field
		MountPath     respjson.Field
		Name          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsMemoryStoreResource) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsMemoryStoreResource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsMemoryStoreResourceType string

const (
	ManagedAgentsMemoryStoreResourceTypeMemoryStore ManagedAgentsMemoryStoreResourceType = "memory_store"
)

// Access mode for an attached memory store.
type ManagedAgentsMemoryStoreResourceAccess string

const (
	ManagedAgentsMemoryStoreResourceAccessReadWrite ManagedAgentsMemoryStoreResourceAccess = "read_write"
	ManagedAgentsMemoryStoreResourceAccessReadOnly  ManagedAgentsMemoryStoreResourceAccess = "read_only"
)

// ManagedAgentsSessionResourceUnion contains all possible properties and
// values from [ManagedAgentsGitHubRepositoryResource],
// [ManagedAgentsFileResource], [ManagedAgentsMemoryStoreResource].
//
// Use the [ManagedAgentsSessionResourceUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsSessionResourceUnion struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	MountPath string    `json:"mount_path"`
	// Any of "github_repository", "file", "memory_store".
	Type      string    `json:"type"`
	UpdatedAt time.Time `json:"updated_at"`
	// This field is from variant [ManagedAgentsGitHubRepositoryResource].
	URL string `json:"url"`
	// This field is from variant [ManagedAgentsGitHubRepositoryResource].
	Checkout ManagedAgentsGitHubRepositoryResourceCheckoutUnion `json:"checkout"`
	// This field is from variant [ManagedAgentsFileResource].
	FileID string `json:"file_id"`
	// This field is from variant [ManagedAgentsMemoryStoreResource].
	MemoryStoreID string `json:"memory_store_id"`
	// This field is from variant [ManagedAgentsMemoryStoreResource].
	Access ManagedAgentsMemoryStoreResourceAccess `json:"access"`
	// This field is from variant [ManagedAgentsMemoryStoreResource].
	Description string `json:"description"`
	// This field is from variant [ManagedAgentsMemoryStoreResource].
	Instructions string `json:"instructions"`
	// This field is from variant [ManagedAgentsMemoryStoreResource].
	Name string `json:"name"`
	JSON struct {
		ID            respjson.Field
		CreatedAt     respjson.Field
		MountPath     respjson.Field
		Type          respjson.Field
		UpdatedAt     respjson.Field
		URL           respjson.Field
		Checkout      respjson.Field
		FileID        respjson.Field
		MemoryStoreID respjson.Field
		Access        respjson.Field
		Description   respjson.Field
		Instructions  respjson.Field
		Name          respjson.Field
		raw           string
	} `json:"-"`
}

// anyManagedAgentsSessionResource is implemented by each variant of
// [ManagedAgentsSessionResourceUnion] to add type safety for the return type
// of [ManagedAgentsSessionResourceUnion.AsAny]
type anyManagedAgentsSessionResource interface {
	implManagedAgentsSessionResourceUnion()
}

func (ManagedAgentsGitHubRepositoryResource) implManagedAgentsSessionResourceUnion() {}
func (ManagedAgentsFileResource) implManagedAgentsSessionResourceUnion()             {}
func (ManagedAgentsMemoryStoreResource) implManagedAgentsSessionResourceUnion()      {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsSessionResourceUnion.AsAny().(type) {
//	case qoder.ManagedAgentsGitHubRepositoryResource:
//	case qoder.ManagedAgentsFileResource:
//	case qoder.ManagedAgentsMemoryStoreResource:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsSessionResourceUnion) AsAny() anyManagedAgentsSessionResource {
	switch u.Type {
	case "github_repository":
		return u.AsGitHubRepository()
	case "file":
		return u.AsFile()
	case "memory_store":
		return u.AsMemoryStore()
	}
	return nil
}

func (u ManagedAgentsSessionResourceUnion) AsGitHubRepository() (v ManagedAgentsGitHubRepositoryResource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionResourceUnion) AsFile() (v ManagedAgentsFileResource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionResourceUnion) AsMemoryStore() (v ManagedAgentsMemoryStoreResource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsSessionResourceUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsSessionResourceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// SessionResourceGetResponseUnion contains all possible properties and values
// from [ManagedAgentsGitHubRepositoryResource],
// [ManagedAgentsFileResource], [ManagedAgentsMemoryStoreResource].
//
// Use the [SessionResourceGetResponseUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type SessionResourceGetResponseUnion struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	MountPath string    `json:"mount_path"`
	// Any of "github_repository", "file", "memory_store".
	Type      string    `json:"type"`
	UpdatedAt time.Time `json:"updated_at"`
	// This field is from variant [ManagedAgentsGitHubRepositoryResource].
	URL string `json:"url"`
	// This field is from variant [ManagedAgentsGitHubRepositoryResource].
	Checkout ManagedAgentsGitHubRepositoryResourceCheckoutUnion `json:"checkout"`
	// This field is from variant [ManagedAgentsFileResource].
	FileID string `json:"file_id"`
	// This field is from variant [ManagedAgentsMemoryStoreResource].
	MemoryStoreID string `json:"memory_store_id"`
	// This field is from variant [ManagedAgentsMemoryStoreResource].
	Access ManagedAgentsMemoryStoreResourceAccess `json:"access"`
	// This field is from variant [ManagedAgentsMemoryStoreResource].
	Description string `json:"description"`
	// This field is from variant [ManagedAgentsMemoryStoreResource].
	Instructions string `json:"instructions"`
	// This field is from variant [ManagedAgentsMemoryStoreResource].
	Name string `json:"name"`
	JSON struct {
		ID            respjson.Field
		CreatedAt     respjson.Field
		MountPath     respjson.Field
		Type          respjson.Field
		UpdatedAt     respjson.Field
		URL           respjson.Field
		Checkout      respjson.Field
		FileID        respjson.Field
		MemoryStoreID respjson.Field
		Access        respjson.Field
		Description   respjson.Field
		Instructions  respjson.Field
		Name          respjson.Field
		raw           string
	} `json:"-"`
}

// anySessionResourceGetResponse is implemented by each variant of
// [SessionResourceGetResponseUnion] to add type safety for the return type of
// [SessionResourceGetResponseUnion.AsAny]
type anySessionResourceGetResponse interface {
	implSessionResourceGetResponseUnion()
}

func (ManagedAgentsGitHubRepositoryResource) implSessionResourceGetResponseUnion() {}
func (ManagedAgentsFileResource) implSessionResourceGetResponseUnion()             {}
func (ManagedAgentsMemoryStoreResource) implSessionResourceGetResponseUnion()      {}

// Use the following switch statement to find the correct variant
//
//	switch variant := SessionResourceGetResponseUnion.AsAny().(type) {
//	case qoder.ManagedAgentsGitHubRepositoryResource:
//	case qoder.ManagedAgentsFileResource:
//	case qoder.ManagedAgentsMemoryStoreResource:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u SessionResourceGetResponseUnion) AsAny() anySessionResourceGetResponse {
	switch u.Type {
	case "github_repository":
		return u.AsGitHubRepository()
	case "file":
		return u.AsFile()
	case "memory_store":
		return u.AsMemoryStore()
	}
	return nil
}

func (u SessionResourceGetResponseUnion) AsGitHubRepository() (v ManagedAgentsGitHubRepositoryResource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SessionResourceGetResponseUnion) AsFile() (v ManagedAgentsFileResource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SessionResourceGetResponseUnion) AsMemoryStore() (v ManagedAgentsMemoryStoreResource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u SessionResourceGetResponseUnion) RawJSON() string { return u.JSON.raw }

func (r *SessionResourceGetResponseUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// SessionResourceUpdateResponseUnion contains all possible properties and
// values from [ManagedAgentsGitHubRepositoryResource],
// [ManagedAgentsFileResource], [ManagedAgentsMemoryStoreResource].
//
// Use the [SessionResourceUpdateResponseUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type SessionResourceUpdateResponseUnion struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	MountPath string    `json:"mount_path"`
	// Any of "github_repository", "file", "memory_store".
	Type      string    `json:"type"`
	UpdatedAt time.Time `json:"updated_at"`
	// This field is from variant [ManagedAgentsGitHubRepositoryResource].
	URL string `json:"url"`
	// This field is from variant [ManagedAgentsGitHubRepositoryResource].
	Checkout ManagedAgentsGitHubRepositoryResourceCheckoutUnion `json:"checkout"`
	// This field is from variant [ManagedAgentsFileResource].
	FileID string `json:"file_id"`
	// This field is from variant [ManagedAgentsMemoryStoreResource].
	MemoryStoreID string `json:"memory_store_id"`
	// This field is from variant [ManagedAgentsMemoryStoreResource].
	Access ManagedAgentsMemoryStoreResourceAccess `json:"access"`
	// This field is from variant [ManagedAgentsMemoryStoreResource].
	Description string `json:"description"`
	// This field is from variant [ManagedAgentsMemoryStoreResource].
	Instructions string `json:"instructions"`
	// This field is from variant [ManagedAgentsMemoryStoreResource].
	Name string `json:"name"`
	JSON struct {
		ID            respjson.Field
		CreatedAt     respjson.Field
		MountPath     respjson.Field
		Type          respjson.Field
		UpdatedAt     respjson.Field
		URL           respjson.Field
		Checkout      respjson.Field
		FileID        respjson.Field
		MemoryStoreID respjson.Field
		Access        respjson.Field
		Description   respjson.Field
		Instructions  respjson.Field
		Name          respjson.Field
		raw           string
	} `json:"-"`
}

// anySessionResourceUpdateResponse is implemented by each variant of
// [SessionResourceUpdateResponseUnion] to add type safety for the return type
// of [SessionResourceUpdateResponseUnion.AsAny]
type anySessionResourceUpdateResponse interface {
	implSessionResourceUpdateResponseUnion()
}

func (ManagedAgentsGitHubRepositoryResource) implSessionResourceUpdateResponseUnion() {}
func (ManagedAgentsFileResource) implSessionResourceUpdateResponseUnion()             {}
func (ManagedAgentsMemoryStoreResource) implSessionResourceUpdateResponseUnion()      {}

// Use the following switch statement to find the correct variant
//
//	switch variant := SessionResourceUpdateResponseUnion.AsAny().(type) {
//	case qoder.ManagedAgentsGitHubRepositoryResource:
//	case qoder.ManagedAgentsFileResource:
//	case qoder.ManagedAgentsMemoryStoreResource:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u SessionResourceUpdateResponseUnion) AsAny() anySessionResourceUpdateResponse {
	switch u.Type {
	case "github_repository":
		return u.AsGitHubRepository()
	case "file":
		return u.AsFile()
	case "memory_store":
		return u.AsMemoryStore()
	}
	return nil
}

func (u SessionResourceUpdateResponseUnion) AsGitHubRepository() (v ManagedAgentsGitHubRepositoryResource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SessionResourceUpdateResponseUnion) AsFile() (v ManagedAgentsFileResource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u SessionResourceUpdateResponseUnion) AsMemoryStore() (v ManagedAgentsMemoryStoreResource) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u SessionResourceUpdateResponseUnion) RawJSON() string { return u.JSON.raw }

func (r *SessionResourceUpdateResponseUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SessionResourceGetParams struct {
	SessionID string `path:"session_id" api:"required" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

type SessionResourceUpdateParams struct {
	Password param.Opt[string] `json:"password,omitzero"`

	SessionID string `path:"session_id" api:"required" json:"-"`
	// New authorization token for the resource. Currently only `github_repository`
	// resources support token rotation.
	AuthorizationToken string `json:"authorization_token,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

func (r SessionResourceUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow SessionResourceUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *SessionResourceUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type SessionResourceListParams struct {
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	AfterID  param.Opt[string] `query:"after_id,omitzero" json:"-"`

	// Maximum number of resources to return per page (max 1000). If omitted, returns
	// all resources.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque cursor from a previous response's `next_page` field.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [SessionResourceListParams]'s query parameters as
// `url.Values`.
func (r SessionResourceListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SessionResourceDeleteParams struct {
	SessionID string `path:"session_id" api:"required" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

type SessionResourceAddParams struct {
	// Mount a file uploaded via the Files API into the session.
	ManagedAgentsFileResourceParams ManagedAgentsFileResourceParams
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

func (r SessionResourceAddParams) MarshalJSON() (data []byte, err error) {
	return param.MarshalObject(r, r.ManagedAgentsFileResourceParams)
}
func (r *SessionResourceAddParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, &r.ManagedAgentsFileResourceParams)
}
