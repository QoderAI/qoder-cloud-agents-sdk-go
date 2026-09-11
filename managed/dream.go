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

// DreamService contains methods and other services that help with interacting
// with the Qoder Cloud Agents API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewDreamService] method instead.
type DreamService struct {
	Options []option.RequestOption
}

// NewDreamService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewDreamService(opts ...option.RequestOption) (r DreamService) {
	r = DreamService{}
	r.Options = opts
	return
}

// Create a Dream
func (r *DreamService) New(ctx context.Context, params DreamNewParams, opts ...option.RequestOption) (res *Dream, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	path := "dreams"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Get a Dream
func (r *DreamService) Get(ctx context.Context, dreamID string, query DreamGetParams, opts ...option.RequestOption) (res *Dream, err error) {
	for _, v := range query.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if dreamID == "" {
		err = errors.New("missing required dream_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("dreams/%s", url.PathEscape(dreamID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List Dreams
func (r *DreamService) List(ctx context.Context, params DreamListParams, opts ...option.RequestOption) (res *pagination.PageCursor[Dream], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "dreams"
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

// List Dreams
func (r *DreamService) ListAutoPaging(ctx context.Context, params DreamListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[Dream] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// Archive a Dream
func (r *DreamService) Archive(ctx context.Context, dreamID string, body DreamArchiveParams, opts ...option.RequestOption) (res *Dream, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if dreamID == "" {
		err = errors.New("missing required dream_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("dreams/%s/archive", url.PathEscape(dreamID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Cancel a Dream
func (r *DreamService) Cancel(ctx context.Context, dreamID string, body DreamCancelParams, opts ...option.RequestOption) (res *Dream, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if dreamID == "" {
		err = errors.New("missing required dream_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("dreams/%s/cancel", url.PathEscape(dreamID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// An asynchronous memory-consolidation job that reads a memory store plus a set of
// session transcripts and writes consolidated memories into an output memory store
// — a new store by default, or an existing store chosen via output_behavior. The
// Dreams API is in research preview: the request and response shapes are volatile
// and may change without the deprecation period that applies to
// generally-available endpoints.
type Dream struct {
	ID string `json:"id" api:"required"`
	// A timestamp in RFC 3339 format
	ArchivedAt time.Time `json:"archived_at" api:"required" format:"date-time"`
	// A timestamp in RFC 3339 format
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// A timestamp in RFC 3339 format
	EndedAt time.Time `json:"ended_at" api:"required" format:"date-time"`
	// Failure detail for a Dream whose `status` is `failed`.
	Error        DreamError        `json:"error" api:"required"`
	Inputs       []DreamInputUnion `json:"inputs" api:"required"`
	Instructions string            `json:"instructions" api:"required"`
	// Model identifier and configuration applied to every pipeline stage. Same wire
	// shape as the Agents API ModelConfig.
	Model DreamModelConfig `json:"model" api:"required"`
	// The default destination: the job creates a new output memory store as a clone of
	// the memory_store input and writes the consolidated memories into it. The input
	// store is never mutated.
	OutputBehavior OutputBehaviorUnion `json:"output_behavior" api:"required"`
	Outputs        []DreamOutput       `json:"outputs" api:"required"`
	SessionID      string              `json:"session_id" api:"required"`
	// Lifecycle status of a Dream.
	//
	// Any of "pending", "running", "completed", "failed", "canceled".
	Status DreamStatus `json:"status" api:"required"`
	// Any of "dream".
	Type DreamType `json:"type" api:"required"`
	// Cumulative token usage for the dream across every pipeline stage.
	Usage DreamUsage `json:"usage" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		ArchivedAt     respjson.Field
		CreatedAt      respjson.Field
		EndedAt        respjson.Field
		Error          respjson.Field
		Inputs         respjson.Field
		Instructions   respjson.Field
		Model          respjson.Field
		OutputBehavior respjson.Field
		Outputs        respjson.Field
		SessionID      respjson.Field
		Status         respjson.Field
		Type           respjson.Field
		Usage          respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r Dream) RawJSON() string { return r.JSON.raw }
func (r *Dream) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DreamType string

const (
	DreamTypeDream DreamType = "dream"
)

// Failure detail for a Dream whose `status` is `failed`.
type DreamError struct {
	Message string `json:"message" api:"required"`
	Type    string `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DreamError) RawJSON() string { return r.JSON.raw }
func (r *DreamError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// DreamInputUnion contains all possible properties and values from
// [DreamMemoryStoreInput], [DreamSessionsInput].
//
// Use the [DreamInputUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type DreamInputUnion struct {
	// This field is from variant [DreamMemoryStoreInput].
	MemoryStoreID string `json:"memory_store_id"`
	// Any of "memory_store", "sessions".
	Type string `json:"type"`
	// This field is from variant [DreamSessionsInput].
	SessionIDs []string `json:"session_ids"`
	JSON       struct {
		MemoryStoreID respjson.Field
		Type          respjson.Field
		SessionIDs    respjson.Field
		raw           string
	} `json:"-"`
}

// anyDreamInput is implemented by each variant of [DreamInputUnion] to add
// type safety for the return type of [DreamInputUnion.AsAny]
type anyDreamInput interface {
	implDreamInputUnion()
}

func (DreamMemoryStoreInput) implDreamInputUnion() {}
func (DreamSessionsInput) implDreamInputUnion()    {}

// Use the following switch statement to find the correct variant
//
//	switch variant := DreamInputUnion.AsAny().(type) {
//	case qoder.DreamMemoryStoreInput:
//	case qoder.DreamSessionsInput:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u DreamInputUnion) AsAny() anyDreamInput {
	switch u.Type {
	case "memory_store":
		return u.AsMemoryStore()
	case "sessions":
		return u.AsSessions()
	}
	return nil
}

func (u DreamInputUnion) AsMemoryStore() (v DreamMemoryStoreInput) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u DreamInputUnion) AsSessions() (v DreamSessionsInput) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u DreamInputUnion) RawJSON() string { return u.JSON.raw }

func (r *DreamInputUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this DreamInputUnion to a DreamInputUnionParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// DreamInputUnionParam.Overrides()
func (r DreamInputUnion) ToParam() DreamInputUnionParam {
	return param.Override[DreamInputUnionParam](json.RawMessage(r.RawJSON()))
}

func DreamInputParamOfMemoryStore(memoryStoreID string) DreamInputUnionParam {
	var memoryStore DreamMemoryStoreInputParam
	memoryStore.MemoryStoreID = memoryStoreID
	return DreamInputUnionParam{OfMemoryStore: &memoryStore}
}

func DreamInputParamOfSessions(sessionIDs []string) DreamInputUnionParam {
	var sessions DreamSessionsInputParam
	sessions.SessionIDs = sessionIDs
	return DreamInputUnionParam{OfSessions: &sessions}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DreamInputUnionParam struct {
	OfMemoryStore *DreamMemoryStoreInputParam `json:",omitzero,inline"`
	OfSessions    *DreamSessionsInputParam    `json:",omitzero,inline"`
	paramUnion
}

func (u DreamInputUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfMemoryStore, u.OfSessions)
}
func (u *DreamInputUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DreamInputUnionParam) asAny() any {
	if !param.IsOmitted(u.OfMemoryStore) {
		return u.OfMemoryStore
	} else if !param.IsOmitted(u.OfSessions) {
		return u.OfSessions
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DreamInputUnionParam) GetMemoryStoreID() *string {
	if vt := u.OfMemoryStore; vt != nil {
		return &vt.MemoryStoreID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DreamInputUnionParam) GetSessionIDs() []string {
	if vt := u.OfSessions; vt != nil {
		return vt.SessionIDs
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DreamInputUnionParam) GetType() *string {
	if vt := u.OfMemoryStore; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSessions; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[DreamInputUnionParam](
		"type",
		apijson.Discriminator[DreamMemoryStoreInputParam]("memory_store"),
		apijson.Discriminator[DreamSessionsInputParam]("sessions"),
	)
}

// An input memory store the dream reads from. The dream never mutates this store
// unless it is also the destination: with output_behavior {type:
// "update_existing"} the job consolidates this store in place.
type DreamMemoryStoreInput struct {
	MemoryStoreID string `json:"memory_store_id" api:"required"`
	// Any of "memory_store".
	Type DreamMemoryStoreInputType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		MemoryStoreID respjson.Field
		Type          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DreamMemoryStoreInput) RawJSON() string { return r.JSON.raw }
func (r *DreamMemoryStoreInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this DreamMemoryStoreInput to a
// DreamMemoryStoreInputParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// DreamMemoryStoreInputParam.Overrides()
func (r DreamMemoryStoreInput) ToParam() DreamMemoryStoreInputParam {
	return param.Override[DreamMemoryStoreInputParam](json.RawMessage(r.RawJSON()))
}

type DreamMemoryStoreInputType string

const (
	DreamMemoryStoreInputTypeMemoryStore DreamMemoryStoreInputType = "memory_store"
)

// An input memory store the dream reads from. The dream never mutates this store
// unless it is also the destination: with output_behavior {type:
// "update_existing"} the job consolidates this store in place.
//
// The properties MemoryStoreID, Type are required.
type DreamMemoryStoreInputParam struct {
	MemoryStoreID string `json:"memory_store_id" api:"required"`
	// Any of "memory_store".
	Type DreamMemoryStoreInputType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r DreamMemoryStoreInputParam) MarshalJSON() (data []byte, err error) {
	type shadow DreamMemoryStoreInputParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DreamMemoryStoreInputParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Model identifier and configuration applied to every pipeline stage. Same wire
// shape as the Agents API ModelConfig.
type DreamModelConfig struct {
	// Model identifier. 1-256 characters.
	ID string `json:"id" api:"required"`
	// Inference speed mode. `fast` provides significantly faster output token
	// generation at premium pricing. Not all models support `fast`; invalid
	// combinations are rejected at create time.
	//
	// Any of "standard", "fast".
	Speed DreamModelConfigSpeed `json:"speed"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID          respjson.Field
		Speed       respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DreamModelConfig) RawJSON() string { return r.JSON.raw }
func (r *DreamModelConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Inference speed mode. `fast` provides significantly faster output token
// generation at premium pricing. Not all models support `fast`; invalid
// combinations are rejected at create time.
type DreamModelConfigSpeed string

const (
	DreamModelConfigSpeedStandard DreamModelConfigSpeed = "standard"
	DreamModelConfigSpeedFast     DreamModelConfigSpeed = "fast"
)

// Model identifier and configuration applied to every pipeline stage.
//
// The property ID is required.
type DreamModelConfigParam struct {
	// Model identifier. 1-256 characters.
	ID string `json:"id" api:"required"`
	// Inference speed mode. `fast` provides significantly faster output token
	// generation at premium pricing. Not all models support `fast`; invalid
	// combinations are rejected at create time.
	//
	// Any of "standard", "fast".
	Speed DreamModelConfigParamSpeed `json:"speed,omitzero"`
	paramObj
}

func (r DreamModelConfigParam) MarshalJSON() (data []byte, err error) {
	type shadow DreamModelConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DreamModelConfigParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Inference speed mode. `fast` provides significantly faster output token
// generation at premium pricing. Not all models support `fast`; invalid
// combinations are rejected at create time.
type DreamModelConfigParamSpeed string

const (
	DreamModelConfigParamSpeedStandard DreamModelConfigParamSpeed = "standard"
	DreamModelConfigParamSpeedFast     DreamModelConfigParamSpeed = "fast"
)

// An output memory store the dream writes consolidated memories into.
type DreamOutput struct {
	MemoryStoreID string `json:"memory_store_id" api:"required"`
	// Any of "memory_store".
	Type DreamOutputType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		MemoryStoreID respjson.Field
		Type          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DreamOutput) RawJSON() string { return r.JSON.raw }
func (r *DreamOutput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DreamOutputType string

const (
	DreamOutputTypeMemoryStore DreamOutputType = "memory_store"
)

// Input session transcripts the dream reads.
type DreamSessionsInput struct {
	SessionIDs []string `json:"session_ids" api:"required"`
	// Any of "sessions".
	Type DreamSessionsInputType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		SessionIDs  respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DreamSessionsInput) RawJSON() string { return r.JSON.raw }
func (r *DreamSessionsInput) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this DreamSessionsInput to a DreamSessionsInputParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// DreamSessionsInputParam.Overrides()
func (r DreamSessionsInput) ToParam() DreamSessionsInputParam {
	return param.Override[DreamSessionsInputParam](json.RawMessage(r.RawJSON()))
}

type DreamSessionsInputType string

const (
	DreamSessionsInputTypeSessions DreamSessionsInputType = "sessions"
)

// Input session transcripts the dream reads.
//
// The properties SessionIDs, Type are required.
type DreamSessionsInputParam struct {
	SessionIDs []string `json:"session_ids,omitzero" api:"required"`
	// Any of "sessions".
	Type DreamSessionsInputType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r DreamSessionsInputParam) MarshalJSON() (data []byte, err error) {
	type shadow DreamSessionsInputParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DreamSessionsInputParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Lifecycle status of a Dream.
type DreamStatus string

const (
	DreamStatusPending   DreamStatus = "pending"
	DreamStatusRunning   DreamStatus = "running"
	DreamStatusCompleted DreamStatus = "completed"
	DreamStatusFailed    DreamStatus = "failed"
	DreamStatusCanceled  DreamStatus = "canceled"
)

// Cumulative token usage for the dream across every pipeline stage.
type DreamUsage struct {
	// Total tokens used to create prompt-cache entries (sum of all TTL tiers).
	CacheCreationInputTokens int64 `json:"cache_creation_input_tokens" api:"required"`
	// Total tokens read from prompt cache.
	CacheReadInputTokens int64 `json:"cache_read_input_tokens" api:"required"`
	// Total uncached input tokens consumed across every pipeline stage.
	InputTokens int64 `json:"input_tokens" api:"required"`
	// Total output tokens generated across every pipeline stage.
	OutputTokens int64 `json:"output_tokens" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		CacheCreationInputTokens respjson.Field
		CacheReadInputTokens     respjson.Field
		InputTokens              respjson.Field
		OutputTokens             respjson.Field
		ExtraFields              map[string]respjson.Field
		raw                      string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r DreamUsage) RawJSON() string { return r.JSON.raw }
func (r *DreamUsage) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// OutputBehaviorUnion contains all possible properties and values from
// [OutputBehaviorCreateNew], [OutputBehaviorUpdateExisting].
//
// Use the [OutputBehaviorUnion.AsAny] method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type OutputBehaviorUnion struct {
	// Any of "create_new", "update_existing".
	Type string `json:"type"`
	// This field is from variant [OutputBehaviorUpdateExisting].
	MemoryStoreID string `json:"memory_store_id"`
	JSON          struct {
		Type          respjson.Field
		MemoryStoreID respjson.Field
		raw           string
	} `json:"-"`
}

// anyOutputBehavior is implemented by each variant of
// [OutputBehaviorUnion] to add type safety for the return type of
// [OutputBehaviorUnion.AsAny]
type anyOutputBehavior interface {
	implOutputBehaviorUnion()
}

func (OutputBehaviorCreateNew) implOutputBehaviorUnion()      {}
func (OutputBehaviorUpdateExisting) implOutputBehaviorUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := OutputBehaviorUnion.AsAny().(type) {
//	case qoder.OutputBehaviorCreateNew:
//	case qoder.OutputBehaviorUpdateExisting:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u OutputBehaviorUnion) AsAny() anyOutputBehavior {
	switch u.Type {
	case "create_new":
		return u.AsCreateNew()
	case "update_existing":
		return u.AsUpdateExisting()
	}
	return nil
}

func (u OutputBehaviorUnion) AsCreateNew() (v OutputBehaviorCreateNew) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u OutputBehaviorUnion) AsUpdateExisting() (v OutputBehaviorUpdateExisting) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u OutputBehaviorUnion) RawJSON() string { return u.JSON.raw }

func (r *OutputBehaviorUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this OutputBehaviorUnion to a OutputBehaviorUnionParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// OutputBehaviorUnionParam.Overrides()
func (r OutputBehaviorUnion) ToParam() OutputBehaviorUnionParam {
	return param.Override[OutputBehaviorUnionParam](json.RawMessage(r.RawJSON()))
}

func OutputBehaviorParamOfCreateNew(type_ OutputBehaviorCreateNewType) OutputBehaviorUnionParam {
	var createNew OutputBehaviorCreateNewParam
	createNew.Type = type_
	return OutputBehaviorUnionParam{OfCreateNew: &createNew}
}

func OutputBehaviorParamOfUpdateExisting(memoryStoreID string) OutputBehaviorUnionParam {
	var updateExisting OutputBehaviorUpdateExistingParam
	updateExisting.MemoryStoreID = memoryStoreID
	return OutputBehaviorUnionParam{OfUpdateExisting: &updateExisting}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type OutputBehaviorUnionParam struct {
	OfCreateNew      *OutputBehaviorCreateNewParam      `json:",omitzero,inline"`
	OfUpdateExisting *OutputBehaviorUpdateExistingParam `json:",omitzero,inline"`
	paramUnion
}

func (u OutputBehaviorUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfCreateNew, u.OfUpdateExisting)
}
func (u *OutputBehaviorUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *OutputBehaviorUnionParam) asAny() any {
	if !param.IsOmitted(u.OfCreateNew) {
		return u.OfCreateNew
	} else if !param.IsOmitted(u.OfUpdateExisting) {
		return u.OfUpdateExisting
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u OutputBehaviorUnionParam) GetMemoryStoreID() *string {
	if vt := u.OfUpdateExisting; vt != nil {
		return &vt.MemoryStoreID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u OutputBehaviorUnionParam) GetType() *string {
	if vt := u.OfCreateNew; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfUpdateExisting; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

func init() {
	apijson.RegisterUnion[OutputBehaviorUnionParam](
		"type",
		apijson.Discriminator[OutputBehaviorCreateNewParam]("create_new"),
		apijson.Discriminator[OutputBehaviorUpdateExistingParam]("update_existing"),
	)
}

// The default destination: the job creates a new output memory store as a clone of
// the memory_store input and writes the consolidated memories into it. The input
// store is never mutated.
type OutputBehaviorCreateNew struct {
	// Any of "create_new".
	Type OutputBehaviorCreateNewType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OutputBehaviorCreateNew) RawJSON() string { return r.JSON.raw }
func (r *OutputBehaviorCreateNew) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this OutputBehaviorCreateNew to a
// OutputBehaviorCreateNewParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// OutputBehaviorCreateNewParam.Overrides()
func (r OutputBehaviorCreateNew) ToParam() OutputBehaviorCreateNewParam {
	return param.Override[OutputBehaviorCreateNewParam](json.RawMessage(r.RawJSON()))
}

type OutputBehaviorCreateNewType string

const (
	OutputBehaviorCreateNewTypeCreateNew OutputBehaviorCreateNewType = "create_new"
)

// The default destination: the job creates a new output memory store as a clone of
// the memory_store input and writes the consolidated memories into it. The input
// store is never mutated.
//
// The property Type is required.
type OutputBehaviorCreateNewParam struct {
	// Any of "create_new".
	Type OutputBehaviorCreateNewType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r OutputBehaviorCreateNewParam) MarshalJSON() (data []byte, err error) {
	type shadow OutputBehaviorCreateNewParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OutputBehaviorCreateNewParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The job writes the consolidated memories into this existing memory store instead
// of creating one. In EAP the store must be the job's own memory_store input, so
// the job consolidates the store in place.
type OutputBehaviorUpdateExisting struct {
	MemoryStoreID string `json:"memory_store_id" api:"required"`
	// Any of "update_existing".
	Type OutputBehaviorUpdateExistingType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		MemoryStoreID respjson.Field
		Type          respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r OutputBehaviorUpdateExisting) RawJSON() string { return r.JSON.raw }
func (r *OutputBehaviorUpdateExisting) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ToParam converts this OutputBehaviorUpdateExisting to a
// OutputBehaviorUpdateExistingParam.
//
// Warning: the fields of the param type will not be present. ToParam should only
// be used at the last possible moment before sending a request. Test for this with
// OutputBehaviorUpdateExistingParam.Overrides()
func (r OutputBehaviorUpdateExisting) ToParam() OutputBehaviorUpdateExistingParam {
	return param.Override[OutputBehaviorUpdateExistingParam](json.RawMessage(r.RawJSON()))
}

type OutputBehaviorUpdateExistingType string

const (
	OutputBehaviorUpdateExistingTypeUpdateExisting OutputBehaviorUpdateExistingType = "update_existing"
)

// The job writes the consolidated memories into this existing memory store instead
// of creating one. In EAP the store must be the job's own memory_store input, so
// the job consolidates the store in place.
//
// The properties MemoryStoreID, Type are required.
type OutputBehaviorUpdateExistingParam struct {
	MemoryStoreID string `json:"memory_store_id" api:"required"`
	// Any of "update_existing".
	Type OutputBehaviorUpdateExistingType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r OutputBehaviorUpdateExistingParam) MarshalJSON() (data []byte, err error) {
	type shadow OutputBehaviorUpdateExistingParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OutputBehaviorUpdateExistingParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DreamNewParams struct {
	Inputs []DreamInputUnionParam `json:"inputs,omitzero" api:"required"`
	// Model identifier and configuration applied to every pipeline stage.
	Model        DreamNewParamsModelUnion `json:"model,omitzero" api:"required"`
	Instructions param.Opt[string]        `json:"instructions,omitzero"`
	WorkspaceID  param.Opt[string]        `header:"qoder-workspace-id,omitzero" json:"-"`
	// The default destination: the job creates a new output memory store as a clone of
	// the memory_store input and writes the consolidated memories into it. The input
	// store is never mutated.
	OutputBehavior OutputBehaviorUnionParam `json:"output_behavior,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

func (r DreamNewParams) MarshalJSON() (data []byte, err error) {
	type shadow DreamNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DreamNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DreamNewParamsModelUnion struct {
	OfString           param.Opt[string]      `json:",omitzero,inline"`
	OfDreamModelConfig *DreamModelConfigParam `json:",omitzero,inline"`
	paramUnion
}

func (u DreamNewParamsModelUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfDreamModelConfig)
}
func (u *DreamNewParamsModelUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DreamNewParamsModelUnion) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfDreamModelConfig) {
		return u.OfDreamModelConfig
	}
	return nil
}

type DreamGetParams struct {
	WorkspaceID param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

type DreamListParams struct {
	// Return dreams with `created_at` strictly after this timestamp (exclusive lower
	// bound, RFC 3339). Unset applies no lower bound.
	CreatedAtGt param.Opt[time.Time] `query:"created_at[gt],omitzero" format:"date-time" json:"-"`
	// Return dreams with `created_at` strictly before this timestamp (exclusive upper
	// bound, RFC 3339). Unset applies no upper bound.
	CreatedAtLt param.Opt[time.Time] `query:"created_at[lt],omitzero" format:"date-time" json:"-"`
	// Query parameter for include_archived
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// Query parameter for limit
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Query parameter for page
	Page        param.Opt[string] `query:"page,omitzero" json:"-"`
	WorkspaceID param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`
	// Filter by lifecycle status. Repeat the parameter to match any of multiple
	// statuses. Empty applies no status filter.
	Statuses []DreamStatus `query:"statuses,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [DreamListParams]'s query parameters as `url.Values`.
func (r DreamListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type DreamArchiveParams struct {
	WorkspaceID param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

type DreamCancelParams struct {
	WorkspaceID param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}
