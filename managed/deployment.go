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

// DeploymentService contains methods and other services that help with
// interacting with the Qoder Cloud Agents API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewDeploymentService] method instead.
type DeploymentService struct {
	Options []option.RequestOption
}

// NewDeploymentService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewDeploymentService(opts ...option.RequestOption) (r DeploymentService) {
	r = DeploymentService{}
	r.Options = opts
	return
}

// Create Deployment
func (r *DeploymentService) New(ctx context.Context, params DeploymentNewParams, opts ...option.RequestOption) (res *ManagedAgentsDeployment, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	path := "deployments"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Get Deployment
func (r *DeploymentService) Get(ctx context.Context, deploymentID string, query DeploymentGetParams, opts ...option.RequestOption) (res *ManagedAgentsDeployment, err error) {
	for _, v := range query.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if deploymentID == "" {
		err = errors.New("missing required deployment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("deployments/%s", url.PathEscape(deploymentID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update Deployment
func (r *DeploymentService) Update(ctx context.Context, deploymentID string, params DeploymentUpdateParams, opts ...option.RequestOption) (res *ManagedAgentsDeployment, err error) {
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if deploymentID == "" {
		err = errors.New("missing required deployment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("deployments/%s", url.PathEscape(deploymentID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// List Deployments
func (r *DeploymentService) List(ctx context.Context, params DeploymentListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ManagedAgentsDeployment], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "deployments"
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

// List Deployments
func (r *DeploymentService) ListAutoPaging(ctx context.Context, params DeploymentListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ManagedAgentsDeployment] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// Archive Deployment
func (r *DeploymentService) Archive(ctx context.Context, deploymentID string, body DeploymentArchiveParams, opts ...option.RequestOption) (res *ManagedAgentsDeployment, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if deploymentID == "" {
		err = errors.New("missing required deployment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("deployments/%s/archive", url.PathEscape(deploymentID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Pause Deployment
func (r *DeploymentService) Pause(ctx context.Context, deploymentID string, body DeploymentPauseParams, opts ...option.RequestOption) (res *ManagedAgentsDeployment, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if deploymentID == "" {
		err = errors.New("missing required deployment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("deployments/%s/pause", url.PathEscape(deploymentID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Run Deployment Now
func (r *DeploymentService) Run(ctx context.Context, deploymentID string, body DeploymentRunParams, opts ...option.RequestOption) (res *ManagedAgentsDeploymentRun, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if deploymentID == "" {
		err = errors.New("missing required deployment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("deployments/%s/run", url.PathEscape(deploymentID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Unpause Deployment
func (r *DeploymentService) Unpause(ctx context.Context, deploymentID string, body DeploymentUnpauseParams, opts ...option.RequestOption) (res *ManagedAgentsDeployment, err error) {
	for _, v := range body.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if deploymentID == "" {
		err = errors.New("missing required deployment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("deployments/%s/unpause", url.PathEscape(deploymentID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// The deployment's agent was archived.
type ManagedAgentsAgentArchivedDeploymentPausedReasonError struct {
	// Any of "agent_archived_error".
	Type ManagedAgentsAgentArchivedDeploymentPausedReasonErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsAgentArchivedDeploymentPausedReasonError) RawJSON() string {
	return r.JSON.raw
}
func (r *ManagedAgentsAgentArchivedDeploymentPausedReasonError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsAgentArchivedDeploymentPausedReasonErrorType string

const (
	ManagedAgentsAgentArchivedDeploymentPausedReasonErrorTypeAgentArchivedError ManagedAgentsAgentArchivedDeploymentPausedReasonErrorType = "agent_archived_error"
)

// A deployment is a configured instance of an agent — it binds the agent to
// everything needed to run it autonomously: an environment, credentials, initial
// events, and an optional schedule.
type ManagedAgentsDeployment struct {
	EnvironmentVariables string `json:"environment_variables"`

	// Unique identifier for this deployment.
	ID string `json:"id" api:"required"`
	// A resolved agent reference with a concrete version.
	Agent ManagedAgentsAgentReference `json:"agent" api:"required"`
	// A timestamp in RFC 3339 format
	ArchivedAt time.Time `json:"archived_at" api:"required" format:"date-time"`
	// A timestamp in RFC 3339 format
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// Description of what the deployment does.
	Description string `json:"description" api:"required"`
	// ID of the `environment` where sessions run.
	EnvironmentID string `json:"environment_id" api:"required"`
	// Events sent to each session immediately after creation.
	InitialEvents []ManagedAgentsDeploymentInitialEventUnion `json:"initial_events" api:"required"`
	// Arbitrary key-value metadata. Maximum 16 pairs.
	Metadata map[string]string `json:"metadata" api:"required"`
	// Human-readable name.
	Name string `json:"name" api:"required"`
	// Why a deployment is paused. Non-null exactly when `status` is `paused`.
	PausedReason ManagedAgentsDeploymentPausedReasonUnion `json:"paused_reason" api:"required"`
	// Resources attached to sessions created from this deployment. Echoes the input
	// minus write-only credentials.
	Resources []ManagedAgentsSessionResourceConfigUnion `json:"resources" api:"required"`
	// 5-field POSIX cron schedule with computed runtime timestamps.
	Schedule ManagedAgentsSchedule `json:"schedule" api:"required"`
	// Lifecycle status of a deployment.
	//
	// Any of "active", "paused".
	Status ManagedAgentsDeploymentStatus `json:"status" api:"required"`
	// Any of "deployment".
	Type ManagedAgentsDeploymentType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// Vault IDs supplying stored credentials for sessions created from this
	// deployment.
	VaultIDs []string `json:"vault_ids" api:"required"`
	// A hard spend ceiling. The session stops issuing new model requests once the
	// tracked list cost reaches `max_list_cost`.
	Budget ManagedAgentsBudgetLimit `json:"budget" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		EnvironmentVariables respjson.Field

		ID            respjson.Field
		Agent         respjson.Field
		ArchivedAt    respjson.Field
		CreatedAt     respjson.Field
		Description   respjson.Field
		EnvironmentID respjson.Field
		InitialEvents respjson.Field
		Metadata      respjson.Field
		Name          respjson.Field
		PausedReason  respjson.Field
		Resources     respjson.Field
		Schedule      respjson.Field
		Status        respjson.Field
		Type          respjson.Field
		UpdatedAt     respjson.Field
		VaultIDs      respjson.Field
		Budget        respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsDeployment) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsDeployment) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsDeploymentType string

const (
	ManagedAgentsDeploymentTypeDeployment ManagedAgentsDeploymentType = "deployment"
)

// ManagedAgentsDeploymentInitialEventUnion contains all possible properties
// and values from [ManagedAgentsDeploymentUserMessageEvent],
// [ManagedAgentsDeploymentUserDefineOutcomeEvent],
// [ManagedAgentsDeploymentSystemMessageEvent].
//
// Use the [ManagedAgentsDeploymentInitialEventUnion.AsAny] method to switch on
// the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsDeploymentInitialEventUnion struct {
	// This field is a union of
	// [[]ManagedAgentsDeploymentUserMessageEventContentUnion],
	// [[]ManagedAgentsSystemContentBlock]
	Content ManagedAgentsDeploymentInitialEventUnionContent `json:"content"`
	// Any of "user.message", "user.define_outcome", "system.message".
	Type string `json:"type"`
	// This field is from variant [ManagedAgentsDeploymentUserDefineOutcomeEvent].
	Description string `json:"description"`
	// This field is from variant [ManagedAgentsDeploymentUserDefineOutcomeEvent].
	Rubric ManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion `json:"rubric"`
	// This field is from variant [ManagedAgentsDeploymentUserDefineOutcomeEvent].
	MaxIterations int64 `json:"max_iterations"`
	JSON          struct {
		Content       respjson.Field
		Type          respjson.Field
		Description   respjson.Field
		Rubric        respjson.Field
		MaxIterations respjson.Field
		raw           string
	} `json:"-"`
}

// anyManagedAgentsDeploymentInitialEvent is implemented by each variant of
// [ManagedAgentsDeploymentInitialEventUnion] to add type safety for the return
// type of [ManagedAgentsDeploymentInitialEventUnion.AsAny]
type anyManagedAgentsDeploymentInitialEvent interface {
	implManagedAgentsDeploymentInitialEventUnion()
}

func (ManagedAgentsDeploymentUserMessageEvent) implManagedAgentsDeploymentInitialEventUnion() {
}
func (ManagedAgentsDeploymentUserDefineOutcomeEvent) implManagedAgentsDeploymentInitialEventUnion() {
}
func (ManagedAgentsDeploymentSystemMessageEvent) implManagedAgentsDeploymentInitialEventUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsDeploymentInitialEventUnion.AsAny().(type) {
//	case qoder.ManagedAgentsDeploymentUserMessageEvent:
//	case qoder.ManagedAgentsDeploymentUserDefineOutcomeEvent:
//	case qoder.ManagedAgentsDeploymentSystemMessageEvent:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsDeploymentInitialEventUnion) AsAny() anyManagedAgentsDeploymentInitialEvent {
	switch u.Type {
	case "user.message":
		return u.AsUserMessage()
	case "user.define_outcome":
		return u.AsUserDefineOutcome()
	case "system.message":
		return u.AsSystemMessage()
	}
	return nil
}

func (u ManagedAgentsDeploymentInitialEventUnion) AsUserMessage() (v ManagedAgentsDeploymentUserMessageEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentInitialEventUnion) AsUserDefineOutcome() (v ManagedAgentsDeploymentUserDefineOutcomeEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentInitialEventUnion) AsSystemMessage() (v ManagedAgentsDeploymentSystemMessageEvent) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsDeploymentInitialEventUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsDeploymentInitialEventUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsDeploymentInitialEventUnionContent is an implicit subunion of
// [ManagedAgentsDeploymentInitialEventUnion].
// ManagedAgentsDeploymentInitialEventUnionContent provides convenient access
// to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsDeploymentInitialEventUnion].
//
// If the underlying value is not a json object, one of the following properties
// will be valid: OfManagedAgentsDeploymentUserMessageEventContentArray
// OfManagedAgentsSystemContentBlockArray]
type ManagedAgentsDeploymentInitialEventUnionContent struct {
	// This field will be present if the value is a
	// [[]ManagedAgentsDeploymentUserMessageEventContentUnion] instead of an
	// object.
	OfManagedAgentsDeploymentUserMessageEventContentArray []ManagedAgentsDeploymentUserMessageEventContentUnion `json:",inline"`
	// This field will be present if the value is a
	// [[]ManagedAgentsSystemContentBlock] instead of an object.
	OfManagedAgentsSystemContentBlockArray []ManagedAgentsSystemContentBlock `json:",inline"`
	JSON                                   struct {
		OfManagedAgentsDeploymentUserMessageEventContentArray respjson.Field
		OfManagedAgentsSystemContentBlockArray                respjson.Field
		raw                                                   string
	} `json:"-"`
}

func (r *ManagedAgentsDeploymentInitialEventUnionContent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

func ManagedAgentsDeploymentInitialEventParamsOfUserMessage(content []ManagedAgentsUserMessageEventParamsContentUnion) ManagedAgentsDeploymentInitialEventParamsUnion {
	var userMessage ManagedAgentsUserMessageEventParams
	userMessage.Content = content
	return ManagedAgentsDeploymentInitialEventParamsUnion{OfUserMessage: &userMessage}
}

func ManagedAgentsDeploymentInitialEventParamsOfUserDefineOutcome[
	T ManagedAgentsFileRubricParams | ManagedAgentsTextRubricParams,
](description string, rubric T, type_ ManagedAgentsUserDefineOutcomeEventParamsType) ManagedAgentsDeploymentInitialEventParamsUnion {
	var userDefineOutcome ManagedAgentsUserDefineOutcomeEventParams
	userDefineOutcome.Description = description
	switch v := any(rubric).(type) {
	case ManagedAgentsFileRubricParams:
		userDefineOutcome.Rubric.OfFile = &v
	case ManagedAgentsTextRubricParams:
		userDefineOutcome.Rubric.OfText = &v
	}
	userDefineOutcome.Type = type_
	return ManagedAgentsDeploymentInitialEventParamsUnion{OfUserDefineOutcome: &userDefineOutcome}
}

func ManagedAgentsDeploymentInitialEventParamsOfSystemMessage(content []ManagedAgentsSystemContentBlockParam) ManagedAgentsDeploymentInitialEventParamsUnion {
	var systemMessage ManagedAgentsSystemMessageEventParams
	systemMessage.Content = content
	return ManagedAgentsDeploymentInitialEventParamsUnion{OfSystemMessage: &systemMessage}
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type ManagedAgentsDeploymentInitialEventParamsUnion struct {
	OfUserMessage       *ManagedAgentsUserMessageEventParams       `json:",omitzero,inline"`
	OfUserDefineOutcome *ManagedAgentsUserDefineOutcomeEventParams `json:",omitzero,inline"`
	OfSystemMessage     *ManagedAgentsSystemMessageEventParams     `json:",omitzero,inline"`
	paramUnion
}

func (u ManagedAgentsDeploymentInitialEventParamsUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfUserMessage, u.OfUserDefineOutcome, u.OfSystemMessage)
}
func (u *ManagedAgentsDeploymentInitialEventParamsUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *ManagedAgentsDeploymentInitialEventParamsUnion) asAny() any {
	if !param.IsOmitted(u.OfUserMessage) {
		return u.OfUserMessage
	} else if !param.IsOmitted(u.OfUserDefineOutcome) {
		return u.OfUserDefineOutcome
	} else if !param.IsOmitted(u.OfSystemMessage) {
		return u.OfSystemMessage
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsDeploymentInitialEventParamsUnion) GetDescription() *string {
	if vt := u.OfUserDefineOutcome; vt != nil {
		return &vt.Description
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsDeploymentInitialEventParamsUnion) GetRubric() *ManagedAgentsUserDefineOutcomeEventParamsRubricUnion {
	if vt := u.OfUserDefineOutcome; vt != nil {
		return &vt.Rubric
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsDeploymentInitialEventParamsUnion) GetMaxIterations() *int64 {
	if vt := u.OfUserDefineOutcome; vt != nil && vt.MaxIterations.Valid() {
		return &vt.MaxIterations.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u ManagedAgentsDeploymentInitialEventParamsUnion) GetType() *string {
	if vt := u.OfUserMessage; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfUserDefineOutcome; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfSystemMessage; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a subunion which exports methods to access subproperties
//
// Or use AsAny() to get the underlying value
func (u ManagedAgentsDeploymentInitialEventParamsUnion) GetContent() (res managedAgentsDeploymentInitialEventParamsUnionContent) {
	if vt := u.OfUserMessage; vt != nil {
		res.any = &vt.Content
	} else if vt := u.OfSystemMessage; vt != nil {
		res.any = &vt.Content
	}
	return
}

// Can have the runtime types
// [_[]ManagedAgentsUserMessageEventParamsContentUnion],
// [_[]ManagedAgentsSystemContentBlockParam]
type managedAgentsDeploymentInitialEventParamsUnionContent struct{ any }

// Use the following switch statement to get the type of the union:
//
//	switch u.AsAny().(type) {
//	case *[]qoder.ManagedAgentsUserMessageEventParamsContentUnion:
//	case *[]qoder.ManagedAgentsSystemContentBlockParam:
//	default:
//	    fmt.Errorf("not present")
//	}
func (u managedAgentsDeploymentInitialEventParamsUnionContent) AsAny() any { return u.any }

func init() {
	apijson.RegisterUnion[ManagedAgentsDeploymentInitialEventParamsUnion](
		"type",
		apijson.Discriminator[ManagedAgentsUserMessageEventParams]("user.message"),
		apijson.Discriminator[ManagedAgentsUserDefineOutcomeEventParams]("user.define_outcome"),
		apijson.Discriminator[ManagedAgentsSystemMessageEventParams]("system.message"),
	)
}

// ManagedAgentsDeploymentPausedReasonUnion contains all possible properties
// and values from [ManagedAgentsManualDeploymentPausedReason],
// [ManagedAgentsErrorDeploymentPausedReason].
//
// Use the [ManagedAgentsDeploymentPausedReasonUnion.AsAny] method to switch on
// the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsDeploymentPausedReasonUnion struct {
	// Any of "manual", "error".
	Type string `json:"type"`
	// This field is from variant [ManagedAgentsErrorDeploymentPausedReason].
	Error ManagedAgentsDeploymentPausedReasonErrorUnion `json:"error"`
	JSON  struct {
		Type  respjson.Field
		Error respjson.Field
		raw   string
	} `json:"-"`
}

// anyManagedAgentsDeploymentPausedReason is implemented by each variant of
// [ManagedAgentsDeploymentPausedReasonUnion] to add type safety for the return
// type of [ManagedAgentsDeploymentPausedReasonUnion.AsAny]
type anyManagedAgentsDeploymentPausedReason interface {
	implManagedAgentsDeploymentPausedReasonUnion()
}

func (ManagedAgentsManualDeploymentPausedReason) implManagedAgentsDeploymentPausedReasonUnion() {
}
func (ManagedAgentsErrorDeploymentPausedReason) implManagedAgentsDeploymentPausedReasonUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsDeploymentPausedReasonUnion.AsAny().(type) {
//	case qoder.ManagedAgentsManualDeploymentPausedReason:
//	case qoder.ManagedAgentsErrorDeploymentPausedReason:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsDeploymentPausedReasonUnion) AsAny() anyManagedAgentsDeploymentPausedReason {
	switch u.Type {
	case "manual":
		return u.AsManual()
	case "error":
		return u.AsError()
	}
	return nil
}

func (u ManagedAgentsDeploymentPausedReasonUnion) AsManual() (v ManagedAgentsManualDeploymentPausedReason) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentPausedReasonUnion) AsError() (v ManagedAgentsErrorDeploymentPausedReason) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsDeploymentPausedReasonUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsDeploymentPausedReasonUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsDeploymentPausedReasonErrorUnion contains all possible
// properties and values from
// [ManagedAgentsEnvironmentArchivedDeploymentPausedReasonError],
// [ManagedAgentsAgentArchivedDeploymentPausedReasonError],
// [ManagedAgentsEnvironmentNotFoundDeploymentPausedReasonError],
// [ManagedAgentsVaultNotFoundDeploymentPausedReasonError],
// [ManagedAgentsFileNotFoundDeploymentPausedReasonError],
// [ManagedAgentsSessionResourceNotFoundDeploymentPausedReasonError],
// [ManagedAgentsWorkspaceArchivedDeploymentPausedReasonError],
// [ManagedAgentsOrganizationDisabledDeploymentPausedReasonError],
// [ManagedAgentsMemoryStoreArchivedDeploymentPausedReasonError],
// [ManagedAgentsSkillNotFoundDeploymentPausedReasonError],
// [ManagedAgentsVaultArchivedDeploymentPausedReasonError],
// [ManagedAgentsUnknownDeploymentPausedReasonError],
// [ManagedAgentsSelfHostedResourcesUnsupportedDeploymentPausedReasonError],
// [ManagedAgentsMCPEgressBlockedDeploymentPausedReasonError].
//
// Use the [ManagedAgentsDeploymentPausedReasonErrorUnion.AsAny] method to
// switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsDeploymentPausedReasonErrorUnion struct {
	// Any of "environment_archived_error", "agent_archived_error",
	// "environment_not_found_error", "vault_not_found_error", "file_not_found_error",
	// "session_resource_not_found_error", "workspace_archived_error",
	// "organization_disabled_error", "memory_store_archived_error",
	// "skill_not_found_error", "vault_archived_error", "unknown_error",
	// "self_hosted_resources_unsupported_error", "mcp_egress_blocked_error".
	Type string `json:"type"`
	JSON struct {
		Type respjson.Field
		raw  string
	} `json:"-"`
}

// anyManagedAgentsDeploymentPausedReasonError is implemented by each variant
// of [ManagedAgentsDeploymentPausedReasonErrorUnion] to add type safety for
// the return type of [ManagedAgentsDeploymentPausedReasonErrorUnion.AsAny]
type anyManagedAgentsDeploymentPausedReasonError interface {
	implManagedAgentsDeploymentPausedReasonErrorUnion()
}

func (ManagedAgentsEnvironmentArchivedDeploymentPausedReasonError) implManagedAgentsDeploymentPausedReasonErrorUnion() {
}
func (ManagedAgentsAgentArchivedDeploymentPausedReasonError) implManagedAgentsDeploymentPausedReasonErrorUnion() {
}
func (ManagedAgentsEnvironmentNotFoundDeploymentPausedReasonError) implManagedAgentsDeploymentPausedReasonErrorUnion() {
}
func (ManagedAgentsVaultNotFoundDeploymentPausedReasonError) implManagedAgentsDeploymentPausedReasonErrorUnion() {
}
func (ManagedAgentsFileNotFoundDeploymentPausedReasonError) implManagedAgentsDeploymentPausedReasonErrorUnion() {
}
func (ManagedAgentsSessionResourceNotFoundDeploymentPausedReasonError) implManagedAgentsDeploymentPausedReasonErrorUnion() {
}
func (ManagedAgentsWorkspaceArchivedDeploymentPausedReasonError) implManagedAgentsDeploymentPausedReasonErrorUnion() {
}
func (ManagedAgentsOrganizationDisabledDeploymentPausedReasonError) implManagedAgentsDeploymentPausedReasonErrorUnion() {
}
func (ManagedAgentsMemoryStoreArchivedDeploymentPausedReasonError) implManagedAgentsDeploymentPausedReasonErrorUnion() {
}
func (ManagedAgentsSkillNotFoundDeploymentPausedReasonError) implManagedAgentsDeploymentPausedReasonErrorUnion() {
}
func (ManagedAgentsVaultArchivedDeploymentPausedReasonError) implManagedAgentsDeploymentPausedReasonErrorUnion() {
}
func (ManagedAgentsUnknownDeploymentPausedReasonError) implManagedAgentsDeploymentPausedReasonErrorUnion() {
}
func (ManagedAgentsSelfHostedResourcesUnsupportedDeploymentPausedReasonError) implManagedAgentsDeploymentPausedReasonErrorUnion() {
}
func (ManagedAgentsMCPEgressBlockedDeploymentPausedReasonError) implManagedAgentsDeploymentPausedReasonErrorUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsDeploymentPausedReasonErrorUnion.AsAny().(type) {
//	case qoder.ManagedAgentsEnvironmentArchivedDeploymentPausedReasonError:
//	case qoder.ManagedAgentsAgentArchivedDeploymentPausedReasonError:
//	case qoder.ManagedAgentsEnvironmentNotFoundDeploymentPausedReasonError:
//	case qoder.ManagedAgentsVaultNotFoundDeploymentPausedReasonError:
//	case qoder.ManagedAgentsFileNotFoundDeploymentPausedReasonError:
//	case qoder.ManagedAgentsSessionResourceNotFoundDeploymentPausedReasonError:
//	case qoder.ManagedAgentsWorkspaceArchivedDeploymentPausedReasonError:
//	case qoder.ManagedAgentsOrganizationDisabledDeploymentPausedReasonError:
//	case qoder.ManagedAgentsMemoryStoreArchivedDeploymentPausedReasonError:
//	case qoder.ManagedAgentsSkillNotFoundDeploymentPausedReasonError:
//	case qoder.ManagedAgentsVaultArchivedDeploymentPausedReasonError:
//	case qoder.ManagedAgentsUnknownDeploymentPausedReasonError:
//	case qoder.ManagedAgentsSelfHostedResourcesUnsupportedDeploymentPausedReasonError:
//	case qoder.ManagedAgentsMCPEgressBlockedDeploymentPausedReasonError:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsDeploymentPausedReasonErrorUnion) AsAny() anyManagedAgentsDeploymentPausedReasonError {
	switch u.Type {
	case "environment_archived_error":
		return u.AsEnvironmentArchivedError()
	case "agent_archived_error":
		return u.AsAgentArchivedError()
	case "environment_not_found_error":
		return u.AsEnvironmentNotFoundError()
	case "vault_not_found_error":
		return u.AsVaultNotFoundError()
	case "file_not_found_error":
		return u.AsFileNotFoundError()
	case "session_resource_not_found_error":
		return u.AsSessionResourceNotFoundError()
	case "workspace_archived_error":
		return u.AsWorkspaceArchivedError()
	case "organization_disabled_error":
		return u.AsOrganizationDisabledError()
	case "memory_store_archived_error":
		return u.AsMemoryStoreArchivedError()
	case "skill_not_found_error":
		return u.AsSkillNotFoundError()
	case "vault_archived_error":
		return u.AsVaultArchivedError()
	case "unknown_error":
		return u.AsUnknownError()
	case "self_hosted_resources_unsupported_error":
		return u.AsSelfHostedResourcesUnsupportedError()
	case "mcp_egress_blocked_error":
		return u.AsMCPEgressBlockedError()
	}
	return nil
}

func (u ManagedAgentsDeploymentPausedReasonErrorUnion) AsEnvironmentArchivedError() (v ManagedAgentsEnvironmentArchivedDeploymentPausedReasonError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentPausedReasonErrorUnion) AsAgentArchivedError() (v ManagedAgentsAgentArchivedDeploymentPausedReasonError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentPausedReasonErrorUnion) AsEnvironmentNotFoundError() (v ManagedAgentsEnvironmentNotFoundDeploymentPausedReasonError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentPausedReasonErrorUnion) AsVaultNotFoundError() (v ManagedAgentsVaultNotFoundDeploymentPausedReasonError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentPausedReasonErrorUnion) AsFileNotFoundError() (v ManagedAgentsFileNotFoundDeploymentPausedReasonError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentPausedReasonErrorUnion) AsSessionResourceNotFoundError() (v ManagedAgentsSessionResourceNotFoundDeploymentPausedReasonError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentPausedReasonErrorUnion) AsWorkspaceArchivedError() (v ManagedAgentsWorkspaceArchivedDeploymentPausedReasonError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentPausedReasonErrorUnion) AsOrganizationDisabledError() (v ManagedAgentsOrganizationDisabledDeploymentPausedReasonError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentPausedReasonErrorUnion) AsMemoryStoreArchivedError() (v ManagedAgentsMemoryStoreArchivedDeploymentPausedReasonError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentPausedReasonErrorUnion) AsSkillNotFoundError() (v ManagedAgentsSkillNotFoundDeploymentPausedReasonError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentPausedReasonErrorUnion) AsVaultArchivedError() (v ManagedAgentsVaultArchivedDeploymentPausedReasonError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentPausedReasonErrorUnion) AsUnknownError() (v ManagedAgentsUnknownDeploymentPausedReasonError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentPausedReasonErrorUnion) AsSelfHostedResourcesUnsupportedError() (v ManagedAgentsSelfHostedResourcesUnsupportedDeploymentPausedReasonError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentPausedReasonErrorUnion) AsMCPEgressBlockedError() (v ManagedAgentsMCPEgressBlockedDeploymentPausedReasonError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsDeploymentPausedReasonErrorUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsDeploymentPausedReasonErrorUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Lifecycle status of a deployment.
type ManagedAgentsDeploymentStatus string

const (
	ManagedAgentsDeploymentStatusActive ManagedAgentsDeploymentStatus = "active"
	ManagedAgentsDeploymentStatusPaused ManagedAgentsDeploymentStatus = "paused"
)

// Privileged context for the accompanying turn and all subsequent turns, appended
// to the session's system context as a `role: "system"` turn rather than replacing
// the top-level system prompt.
type ManagedAgentsDeploymentSystemMessageEvent struct {
	// System content blocks to append. Text-only.
	Content []ManagedAgentsSystemContentBlock `json:"content" api:"required"`
	// Any of "system.message".
	Type ManagedAgentsDeploymentSystemMessageEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Content     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsDeploymentSystemMessageEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsDeploymentSystemMessageEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsDeploymentSystemMessageEventType string

const (
	ManagedAgentsDeploymentSystemMessageEventTypeSystemMessage ManagedAgentsDeploymentSystemMessageEventType = "system.message"
)

// An outcome the agent should work toward. The agent begins work on receipt.
type ManagedAgentsDeploymentUserDefineOutcomeEvent struct {
	// What the agent should produce. This is the task specification.
	Description string `json:"description" api:"required"`
	// Rubric for grading the quality of an outcome.
	Rubric ManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion `json:"rubric" api:"required"`
	// Any of "user.define_outcome".
	Type ManagedAgentsDeploymentUserDefineOutcomeEventType `json:"type" api:"required"`
	// Eval→revision cycles before giving up. Default 3, max 20.
	MaxIterations int64 `json:"max_iterations" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Description   respjson.Field
		Rubric        respjson.Field
		Type          respjson.Field
		MaxIterations respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsDeploymentUserDefineOutcomeEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsDeploymentUserDefineOutcomeEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion contains all
// possible properties and values from [ManagedAgentsFileRubric],
// [ManagedAgentsTextRubric].
//
// Use the [ManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion.AsAny]
// method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion struct {
	// This field is from variant [ManagedAgentsFileRubric].
	FileID string `json:"file_id"`
	// Any of "file", "text".
	Type string `json:"type"`
	// This field is from variant [ManagedAgentsTextRubric].
	Content string `json:"content"`
	JSON    struct {
		FileID  respjson.Field
		Type    respjson.Field
		Content respjson.Field
		raw     string
	} `json:"-"`
}

// anyManagedAgentsDeploymentUserDefineOutcomeEventRubric is implemented by
// each variant of [ManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion]
// to add type safety for the return type of
// [ManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion.AsAny]
type anyManagedAgentsDeploymentUserDefineOutcomeEventRubric interface {
	implManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion()
}

func (ManagedAgentsFileRubric) implManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion() {
}
func (ManagedAgentsTextRubric) implManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion.AsAny().(type) {
//	case qoder.ManagedAgentsFileRubric:
//	case qoder.ManagedAgentsTextRubric:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion) AsAny() anyManagedAgentsDeploymentUserDefineOutcomeEventRubric {
	switch u.Type {
	case "file":
		return u.AsFile()
	case "text":
		return u.AsText()
	}
	return nil
}

func (u ManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion) AsFile() (v ManagedAgentsFileRubric) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion) AsText() (v ManagedAgentsTextRubric) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *ManagedAgentsDeploymentUserDefineOutcomeEventRubricUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsDeploymentUserDefineOutcomeEventType string

const (
	ManagedAgentsDeploymentUserDefineOutcomeEventTypeUserDefineOutcome ManagedAgentsDeploymentUserDefineOutcomeEventType = "user.define_outcome"
)

// A user message sent to the session.
type ManagedAgentsDeploymentUserMessageEvent struct {
	// Array of content blocks for the user message.
	Content []ManagedAgentsDeploymentUserMessageEventContentUnion `json:"content" api:"required"`
	// Any of "user.message".
	Type ManagedAgentsDeploymentUserMessageEventType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Content     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsDeploymentUserMessageEvent) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsDeploymentUserMessageEvent) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsDeploymentUserMessageEventContentUnion contains all possible
// properties and values from [ManagedAgentsTextBlock],
// [ManagedAgentsImageBlock], [ManagedAgentsDocumentBlock],
// [ManagedAgentsRedactedBlock].
//
// Use the [ManagedAgentsDeploymentUserMessageEventContentUnion.AsAny] method
// to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsDeploymentUserMessageEventContentUnion struct {
	// This field is from variant [ManagedAgentsTextBlock].
	Text string `json:"text"`
	// Any of "text", "image", "document", "redacted".
	Type string `json:"type"`
	// This field is a union of [ManagedAgentsImageBlockSourceUnion],
	// [ManagedAgentsDocumentBlockSourceUnion]
	Source ManagedAgentsDeploymentUserMessageEventContentUnionSource `json:"source"`
	// This field is from variant [ManagedAgentsDocumentBlock].
	Context string `json:"context"`
	// This field is from variant [ManagedAgentsDocumentBlock].
	Title string `json:"title"`
	JSON  struct {
		Text    respjson.Field
		Type    respjson.Field
		Source  respjson.Field
		Context respjson.Field
		Title   respjson.Field
		raw     string
	} `json:"-"`
}

// anyManagedAgentsDeploymentUserMessageEventContent is implemented by each
// variant of [ManagedAgentsDeploymentUserMessageEventContentUnion] to add type
// safety for the return type of
// [ManagedAgentsDeploymentUserMessageEventContentUnion.AsAny]
type anyManagedAgentsDeploymentUserMessageEventContent interface {
	implManagedAgentsDeploymentUserMessageEventContentUnion()
}

func (ManagedAgentsTextBlock) implManagedAgentsDeploymentUserMessageEventContentUnion()     {}
func (ManagedAgentsImageBlock) implManagedAgentsDeploymentUserMessageEventContentUnion()    {}
func (ManagedAgentsDocumentBlock) implManagedAgentsDeploymentUserMessageEventContentUnion() {}
func (ManagedAgentsRedactedBlock) implManagedAgentsDeploymentUserMessageEventContentUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsDeploymentUserMessageEventContentUnion.AsAny().(type) {
//	case qoder.ManagedAgentsTextBlock:
//	case qoder.ManagedAgentsImageBlock:
//	case qoder.ManagedAgentsDocumentBlock:
//	case qoder.ManagedAgentsRedactedBlock:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsDeploymentUserMessageEventContentUnion) AsAny() anyManagedAgentsDeploymentUserMessageEventContent {
	switch u.Type {
	case "text":
		return u.AsText()
	case "image":
		return u.AsImage()
	case "document":
		return u.AsDocument()
	case "redacted":
		return u.AsRedacted()
	}
	return nil
}

func (u ManagedAgentsDeploymentUserMessageEventContentUnion) AsText() (v ManagedAgentsTextBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentUserMessageEventContentUnion) AsImage() (v ManagedAgentsImageBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentUserMessageEventContentUnion) AsDocument() (v ManagedAgentsDocumentBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentUserMessageEventContentUnion) AsRedacted() (v ManagedAgentsRedactedBlock) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsDeploymentUserMessageEventContentUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsDeploymentUserMessageEventContentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsDeploymentUserMessageEventContentUnionSource is an implicit
// subunion of [ManagedAgentsDeploymentUserMessageEventContentUnion].
// ManagedAgentsDeploymentUserMessageEventContentUnionSource provides
// convenient access to the sub-properties of the union.
//
// For type safety it is recommended to directly use a variant of the
// [ManagedAgentsDeploymentUserMessageEventContentUnion].
type ManagedAgentsDeploymentUserMessageEventContentUnionSource struct {
	Data      string `json:"data"`
	MediaType string `json:"media_type"`
	Type      string `json:"type"`
	URL       string `json:"url"`
	FileID    string `json:"file_id"`
	JSON      struct {
		Data      respjson.Field
		MediaType respjson.Field
		Type      respjson.Field
		URL       respjson.Field
		FileID    respjson.Field
		raw       string
	} `json:"-"`
}

func (r *ManagedAgentsDeploymentUserMessageEventContentUnionSource) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsDeploymentUserMessageEventType string

const (
	ManagedAgentsDeploymentUserMessageEventTypeUserMessage ManagedAgentsDeploymentUserMessageEventType = "user.message"
)

// The deployment's environment was archived.
type ManagedAgentsEnvironmentArchivedDeploymentPausedReasonError struct {
	// Any of "environment_archived_error".
	Type ManagedAgentsEnvironmentArchivedDeploymentPausedReasonErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsEnvironmentArchivedDeploymentPausedReasonError) RawJSON() string {
	return r.JSON.raw
}
func (r *ManagedAgentsEnvironmentArchivedDeploymentPausedReasonError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsEnvironmentArchivedDeploymentPausedReasonErrorType string

const (
	ManagedAgentsEnvironmentArchivedDeploymentPausedReasonErrorTypeEnvironmentArchivedError ManagedAgentsEnvironmentArchivedDeploymentPausedReasonErrorType = "environment_archived_error"
)

// The deployment's environment no longer exists.
type ManagedAgentsEnvironmentNotFoundDeploymentPausedReasonError struct {
	// Any of "environment_not_found_error".
	Type ManagedAgentsEnvironmentNotFoundDeploymentPausedReasonErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsEnvironmentNotFoundDeploymentPausedReasonError) RawJSON() string {
	return r.JSON.raw
}
func (r *ManagedAgentsEnvironmentNotFoundDeploymentPausedReasonError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsEnvironmentNotFoundDeploymentPausedReasonErrorType string

const (
	ManagedAgentsEnvironmentNotFoundDeploymentPausedReasonErrorTypeEnvironmentNotFoundError ManagedAgentsEnvironmentNotFoundDeploymentPausedReasonErrorType = "environment_not_found_error"
)

// A scheduled fire recorded a failed run whose error auto-pauses the deployment.
type ManagedAgentsErrorDeploymentPausedReason struct {
	// The error that triggered an auto-pause. Matches the failed run's `error.type`.
	Error ManagedAgentsDeploymentPausedReasonErrorUnion `json:"error" api:"required"`
	// Any of "error".
	Type ManagedAgentsErrorDeploymentPausedReasonType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Error       respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsErrorDeploymentPausedReason) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsErrorDeploymentPausedReason) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsErrorDeploymentPausedReasonType string

const (
	ManagedAgentsErrorDeploymentPausedReasonTypeError ManagedAgentsErrorDeploymentPausedReasonType = "error"
)

// A file resource referenced by the deployment no longer exists.
type ManagedAgentsFileNotFoundDeploymentPausedReasonError struct {
	// Any of "file_not_found_error".
	Type ManagedAgentsFileNotFoundDeploymentPausedReasonErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsFileNotFoundDeploymentPausedReasonError) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsFileNotFoundDeploymentPausedReasonError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsFileNotFoundDeploymentPausedReasonErrorType string

const (
	ManagedAgentsFileNotFoundDeploymentPausedReasonErrorTypeFileNotFoundError ManagedAgentsFileNotFoundDeploymentPausedReasonErrorType = "file_not_found_error"
)

// A file mounted into each session's container.
type ManagedAgentsFileResourceConfig struct {
	// ID of a previously uploaded file.
	FileID string `json:"file_id" api:"required"`
	// Any of "file".
	Type ManagedAgentsFileResourceConfigType `json:"type" api:"required"`
	// Mount path in the container. Defaults to `/mnt/session/uploads/<file_id>`.
	MountPath string `json:"mount_path" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		FileID      respjson.Field
		Type        respjson.Field
		MountPath   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsFileResourceConfig) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsFileResourceConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsFileResourceConfigType string

const (
	ManagedAgentsFileResourceConfigTypeFile ManagedAgentsFileResourceConfigType = "file"
)

// A GitHub repository mounted into each session's container. The authorization
// token is write-only and never returned.
type ManagedAgentsGitHubRepositoryResourceConfig struct {
	// Any of "github_repository".
	Type ManagedAgentsGitHubRepositoryResourceConfigType `json:"type" api:"required"`
	// Github URL of the repository
	URL string `json:"url" api:"required"`
	// Branch or commit to check out. Defaults to the repository's default branch.
	Checkout ManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion `json:"checkout" api:"nullable"`
	// Mount path in the container. Defaults to `/workspace/<repo-name>`.
	MountPath string `json:"mount_path" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		URL         respjson.Field
		Checkout    respjson.Field
		MountPath   respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsGitHubRepositoryResourceConfig) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsGitHubRepositoryResourceConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsGitHubRepositoryResourceConfigType string

const (
	ManagedAgentsGitHubRepositoryResourceConfigTypeGitHubRepository ManagedAgentsGitHubRepositoryResourceConfigType = "github_repository"
)

// ManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion contains all
// possible properties and values from [ManagedAgentsBranchCheckout],
// [ManagedAgentsCommitCheckout].
//
// Use the [ManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion.AsAny]
// method to switch on the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion struct {
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

// anyManagedAgentsGitHubRepositoryResourceConfigCheckout is implemented by
// each variant of [ManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion]
// to add type safety for the return type of
// [ManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion.AsAny]
type anyManagedAgentsGitHubRepositoryResourceConfigCheckout interface {
	implManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion()
}

func (ManagedAgentsBranchCheckout) implManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion() {
}
func (ManagedAgentsCommitCheckout) implManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion() {
}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion.AsAny().(type) {
//	case qoder.ManagedAgentsBranchCheckout:
//	case qoder.ManagedAgentsCommitCheckout:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion) AsAny() anyManagedAgentsGitHubRepositoryResourceConfigCheckout {
	switch u.Type {
	case "branch":
		return u.AsBranch()
	case "commit":
		return u.AsCommit()
	}
	return nil
}

func (u ManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion) AsBranch() (v ManagedAgentsBranchCheckout) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion) AsCommit() (v ManagedAgentsCommitCheckout) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion) RawJSON() string {
	return u.JSON.raw
}

func (r *ManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// The caller invoked the pause endpoint on the deployment.
type ManagedAgentsManualDeploymentPausedReason struct {
	// Any of "manual".
	Type ManagedAgentsManualDeploymentPausedReasonType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsManualDeploymentPausedReason) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsManualDeploymentPausedReason) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsManualDeploymentPausedReasonType string

const (
	ManagedAgentsManualDeploymentPausedReasonTypeManual ManagedAgentsManualDeploymentPausedReasonType = "manual"
)

// An MCP server host used by the deployment's agent is blocked by the
// environment's network policy.
type ManagedAgentsMCPEgressBlockedDeploymentPausedReasonError struct {
	// Any of "mcp_egress_blocked_error".
	Type ManagedAgentsMCPEgressBlockedDeploymentPausedReasonErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsMCPEgressBlockedDeploymentPausedReasonError) RawJSON() string {
	return r.JSON.raw
}
func (r *ManagedAgentsMCPEgressBlockedDeploymentPausedReasonError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsMCPEgressBlockedDeploymentPausedReasonErrorType string

const (
	ManagedAgentsMCPEgressBlockedDeploymentPausedReasonErrorTypeMCPEgressBlockedError ManagedAgentsMCPEgressBlockedDeploymentPausedReasonErrorType = "mcp_egress_blocked_error"
)

// A memory store referenced by the deployment is archived.
type ManagedAgentsMemoryStoreArchivedDeploymentPausedReasonError struct {
	// Any of "memory_store_archived_error".
	Type ManagedAgentsMemoryStoreArchivedDeploymentPausedReasonErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsMemoryStoreArchivedDeploymentPausedReasonError) RawJSON() string {
	return r.JSON.raw
}
func (r *ManagedAgentsMemoryStoreArchivedDeploymentPausedReasonError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsMemoryStoreArchivedDeploymentPausedReasonErrorType string

const (
	ManagedAgentsMemoryStoreArchivedDeploymentPausedReasonErrorTypeMemoryStoreArchivedError ManagedAgentsMemoryStoreArchivedDeploymentPausedReasonErrorType = "memory_store_archived_error"
)

// A memory store attached to each session created from this deployment.
type ManagedAgentsMemoryStoreResourceConfig struct {
	// The memory store ID (memstore\_...). Must belong to the caller's organization
	// and workspace.
	MemoryStoreID string `json:"memory_store_id" api:"required"`
	// Any of "memory_store".
	Type ManagedAgentsMemoryStoreResourceConfigType `json:"type" api:"required"`
	// Access mode for an attached memory store.
	//
	// Any of "read_write", "read_only".
	Access ManagedAgentsMemoryStoreResourceConfigAccess `json:"access" api:"nullable"`
	// Per-attachment guidance for the agent on how to use this store. Rendered into
	// the memory section of the system prompt. Max 4096 chars.
	Instructions string `json:"instructions" api:"nullable"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		MemoryStoreID respjson.Field
		Type          respjson.Field
		Access        respjson.Field
		Instructions  respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsMemoryStoreResourceConfig) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsMemoryStoreResourceConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsMemoryStoreResourceConfigType string

const (
	ManagedAgentsMemoryStoreResourceConfigTypeMemoryStore ManagedAgentsMemoryStoreResourceConfigType = "memory_store"
)

// Access mode for an attached memory store.
type ManagedAgentsMemoryStoreResourceConfigAccess string

const (
	ManagedAgentsMemoryStoreResourceConfigAccessReadWrite ManagedAgentsMemoryStoreResourceConfigAccess = "read_write"
	ManagedAgentsMemoryStoreResourceConfigAccessReadOnly  ManagedAgentsMemoryStoreResourceConfigAccess = "read_only"
)

// The deployment's organization is disabled.
type ManagedAgentsOrganizationDisabledDeploymentPausedReasonError struct {
	// Any of "organization_disabled_error".
	Type ManagedAgentsOrganizationDisabledDeploymentPausedReasonErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsOrganizationDisabledDeploymentPausedReasonError) RawJSON() string {
	return r.JSON.raw
}
func (r *ManagedAgentsOrganizationDisabledDeploymentPausedReasonError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsOrganizationDisabledDeploymentPausedReasonErrorType string

const (
	ManagedAgentsOrganizationDisabledDeploymentPausedReasonErrorTypeOrganizationDisabledError ManagedAgentsOrganizationDisabledDeploymentPausedReasonErrorType = "organization_disabled_error"
)

// 5-field POSIX cron schedule with computed runtime timestamps.
type ManagedAgentsSchedule struct {
	// 5-field POSIX cron expression: minute hour day-of-month month day-of-week (e.g.,
	// "0 9 \* \* 1-5" for weekdays at 9am). Day-of-week is 0-7 where 0 and 7 both mean
	// Sunday. Extended cron syntax - seconds or year fields, and the special
	// characters L, W, #, and ? - is not supported, nor are predefined shortcuts
	// (@daily).
	Expression string `json:"expression" api:"required"`
	// IANA timezone identifier (e.g., "America/Los_Angeles", "UTC").
	Timezone string `json:"timezone" api:"required"`
	// Any of "cron".
	Type ManagedAgentsScheduleType `json:"type" api:"required"`
	// A timestamp in RFC 3339 format
	LastRunAt time.Time `json:"last_run_at" api:"nullable" format:"date-time"`
	// Up to 5 timestamps of upcoming cron occurrences. Non-empty for active and paused
	// deployments (reflects what the schedule would do if unpaused); empty once the
	// deployment is archived (`archived_at` set). Each fire is offset by a small
	// per-schedule jitter, so a run will actually start at or shortly after its listed
	// time.
	UpcomingRunsAt []time.Time `json:"upcoming_runs_at" format:"date-time"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Expression     respjson.Field
		Timezone       respjson.Field
		Type           respjson.Field
		LastRunAt      respjson.Field
		UpcomingRunsAt respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSchedule) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSchedule) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsScheduleType string

const (
	ManagedAgentsScheduleTypeCron ManagedAgentsScheduleType = "cron"
)

// 5-field POSIX cron schedule. Literal wall-clock matching in the configured
// timezone.
//
// The properties Expression, Timezone, Type are required.
type ManagedAgentsScheduleParams struct {
	// 5-field POSIX cron expression: minute hour day-of-month month day-of-week (e.g.,
	// "0 9 \* \* 1-5" for weekdays at 9am). Day-of-week is 0-7 where 0 and 7 both mean
	// Sunday. Extended cron syntax - seconds or year fields, and the special
	// characters L, W, #, and ? - is not supported, nor are predefined shortcuts
	// (@daily).
	Expression string `json:"expression" api:"required"`
	// Required. IANA timezone identifier (e.g., "America/Los_Angeles", "UTC").
	// Validated against the IANA timezone database.
	Timezone string `json:"timezone" api:"required"`
	// Any of "cron".
	Type ManagedAgentsScheduleParamsType `json:"type,omitzero" api:"required"`
	paramObj
}

func (r ManagedAgentsScheduleParams) MarshalJSON() (data []byte, err error) {
	type shadow ManagedAgentsScheduleParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsScheduleParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsScheduleParamsType string

const (
	ManagedAgentsScheduleParamsTypeCron ManagedAgentsScheduleParamsType = "cron"
)

// The deployment configures resources, but its environment is self-hosted and
// cannot mount them.
type ManagedAgentsSelfHostedResourcesUnsupportedDeploymentPausedReasonError struct {
	// Any of "self_hosted_resources_unsupported_error".
	Type ManagedAgentsSelfHostedResourcesUnsupportedDeploymentPausedReasonErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSelfHostedResourcesUnsupportedDeploymentPausedReasonError) RawJSON() string {
	return r.JSON.raw
}
func (r *ManagedAgentsSelfHostedResourcesUnsupportedDeploymentPausedReasonError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSelfHostedResourcesUnsupportedDeploymentPausedReasonErrorType string

const (
	ManagedAgentsSelfHostedResourcesUnsupportedDeploymentPausedReasonErrorTypeSelfHostedResourcesUnsupportedError ManagedAgentsSelfHostedResourcesUnsupportedDeploymentPausedReasonErrorType = "self_hosted_resources_unsupported_error"
)

// ManagedAgentsSessionResourceConfigUnion contains all possible properties and
// values from [ManagedAgentsGitHubRepositoryResourceConfig],
// [ManagedAgentsFileResourceConfig],
// [ManagedAgentsMemoryStoreResourceConfig].
//
// Use the [ManagedAgentsSessionResourceConfigUnion.AsAny] method to switch on
// the variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsSessionResourceConfigUnion struct {
	// Any of "github_repository", "file", "memory_store".
	Type string `json:"type"`
	// This field is from variant [ManagedAgentsGitHubRepositoryResourceConfig].
	URL string `json:"url"`
	// This field is from variant [ManagedAgentsGitHubRepositoryResourceConfig].
	Checkout  ManagedAgentsGitHubRepositoryResourceConfigCheckoutUnion `json:"checkout"`
	MountPath string                                                   `json:"mount_path"`
	// This field is from variant [ManagedAgentsFileResourceConfig].
	FileID string `json:"file_id"`
	// This field is from variant [ManagedAgentsMemoryStoreResourceConfig].
	MemoryStoreID string `json:"memory_store_id"`
	// This field is from variant [ManagedAgentsMemoryStoreResourceConfig].
	Access ManagedAgentsMemoryStoreResourceConfigAccess `json:"access"`
	// This field is from variant [ManagedAgentsMemoryStoreResourceConfig].
	Instructions string `json:"instructions"`
	JSON         struct {
		Type          respjson.Field
		URL           respjson.Field
		Checkout      respjson.Field
		MountPath     respjson.Field
		FileID        respjson.Field
		MemoryStoreID respjson.Field
		Access        respjson.Field
		Instructions  respjson.Field
		raw           string
	} `json:"-"`
}

// anyManagedAgentsSessionResourceConfig is implemented by each variant of
// [ManagedAgentsSessionResourceConfigUnion] to add type safety for the return
// type of [ManagedAgentsSessionResourceConfigUnion.AsAny]
type anyManagedAgentsSessionResourceConfig interface {
	implManagedAgentsSessionResourceConfigUnion()
}

func (ManagedAgentsGitHubRepositoryResourceConfig) implManagedAgentsSessionResourceConfigUnion() {
}
func (ManagedAgentsFileResourceConfig) implManagedAgentsSessionResourceConfigUnion()        {}
func (ManagedAgentsMemoryStoreResourceConfig) implManagedAgentsSessionResourceConfigUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsSessionResourceConfigUnion.AsAny().(type) {
//	case qoder.ManagedAgentsGitHubRepositoryResourceConfig:
//	case qoder.ManagedAgentsFileResourceConfig:
//	case qoder.ManagedAgentsMemoryStoreResourceConfig:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsSessionResourceConfigUnion) AsAny() anyManagedAgentsSessionResourceConfig {
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

func (u ManagedAgentsSessionResourceConfigUnion) AsGitHubRepository() (v ManagedAgentsGitHubRepositoryResourceConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionResourceConfigUnion) AsFile() (v ManagedAgentsFileResourceConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsSessionResourceConfigUnion) AsMemoryStore() (v ManagedAgentsMemoryStoreResourceConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsSessionResourceConfigUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsSessionResourceConfigUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// A referenced resource no longer exists and its kind was not reported.
type ManagedAgentsSessionResourceNotFoundDeploymentPausedReasonError struct {
	// Any of "session_resource_not_found_error".
	Type ManagedAgentsSessionResourceNotFoundDeploymentPausedReasonErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionResourceNotFoundDeploymentPausedReasonError) RawJSON() string {
	return r.JSON.raw
}
func (r *ManagedAgentsSessionResourceNotFoundDeploymentPausedReasonError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSessionResourceNotFoundDeploymentPausedReasonErrorType string

const (
	ManagedAgentsSessionResourceNotFoundDeploymentPausedReasonErrorTypeSessionResourceNotFoundError ManagedAgentsSessionResourceNotFoundDeploymentPausedReasonErrorType = "session_resource_not_found_error"
)

// A skill referenced by the deployment's agent no longer exists.
type ManagedAgentsSkillNotFoundDeploymentPausedReasonError struct {
	// Any of "skill_not_found_error".
	Type ManagedAgentsSkillNotFoundDeploymentPausedReasonErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSkillNotFoundDeploymentPausedReasonError) RawJSON() string {
	return r.JSON.raw
}
func (r *ManagedAgentsSkillNotFoundDeploymentPausedReasonError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSkillNotFoundDeploymentPausedReasonErrorType string

const (
	ManagedAgentsSkillNotFoundDeploymentPausedReasonErrorTypeSkillNotFoundError ManagedAgentsSkillNotFoundDeploymentPausedReasonErrorType = "skill_not_found_error"
)

// An unrecognized error auto-paused the deployment. A fallback variant; matches a
// run whose `error.type` is `unknown_error`.
type ManagedAgentsUnknownDeploymentPausedReasonError struct {
	// Any of "unknown_error".
	Type ManagedAgentsUnknownDeploymentPausedReasonErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsUnknownDeploymentPausedReasonError) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsUnknownDeploymentPausedReasonError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsUnknownDeploymentPausedReasonErrorType string

const (
	ManagedAgentsUnknownDeploymentPausedReasonErrorTypeUnknownError ManagedAgentsUnknownDeploymentPausedReasonErrorType = "unknown_error"
)

// A vault referenced by the deployment is archived.
type ManagedAgentsVaultArchivedDeploymentPausedReasonError struct {
	// Any of "vault_archived_error".
	Type ManagedAgentsVaultArchivedDeploymentPausedReasonErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsVaultArchivedDeploymentPausedReasonError) RawJSON() string {
	return r.JSON.raw
}
func (r *ManagedAgentsVaultArchivedDeploymentPausedReasonError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsVaultArchivedDeploymentPausedReasonErrorType string

const (
	ManagedAgentsVaultArchivedDeploymentPausedReasonErrorTypeVaultArchivedError ManagedAgentsVaultArchivedDeploymentPausedReasonErrorType = "vault_archived_error"
)

// A vault referenced by the deployment no longer exists.
type ManagedAgentsVaultNotFoundDeploymentPausedReasonError struct {
	// Any of "vault_not_found_error".
	Type ManagedAgentsVaultNotFoundDeploymentPausedReasonErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsVaultNotFoundDeploymentPausedReasonError) RawJSON() string {
	return r.JSON.raw
}
func (r *ManagedAgentsVaultNotFoundDeploymentPausedReasonError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsVaultNotFoundDeploymentPausedReasonErrorType string

const (
	ManagedAgentsVaultNotFoundDeploymentPausedReasonErrorTypeVaultNotFoundError ManagedAgentsVaultNotFoundDeploymentPausedReasonErrorType = "vault_not_found_error"
)

// The deployment's workspace was archived.
type ManagedAgentsWorkspaceArchivedDeploymentPausedReasonError struct {
	// Any of "workspace_archived_error".
	Type ManagedAgentsWorkspaceArchivedDeploymentPausedReasonErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsWorkspaceArchivedDeploymentPausedReasonError) RawJSON() string {
	return r.JSON.raw
}
func (r *ManagedAgentsWorkspaceArchivedDeploymentPausedReasonError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsWorkspaceArchivedDeploymentPausedReasonErrorType string

const (
	ManagedAgentsWorkspaceArchivedDeploymentPausedReasonErrorTypeWorkspaceArchivedError ManagedAgentsWorkspaceArchivedDeploymentPausedReasonErrorType = "workspace_archived_error"
)

type DeploymentNewParams struct {
	EnvironmentVariables param.Opt[string] `json:"environment_variables,omitzero"`

	// Agent to deploy. Accepts the `agent` ID string, which pins the latest version,
	// or an `agent` object with both id and version specified. The agent must exist
	// and not be archived.
	Agent DeploymentNewParamsAgentUnion `json:"agent,omitzero" api:"required"`
	// ID of the `environment` defining the container configuration for sessions
	// created from this deployment.
	EnvironmentID string `json:"environment_id" api:"required"`
	// Events to send to each session immediately after creation. At least 1,
	// maximum 50.
	InitialEvents []ManagedAgentsDeploymentInitialEventParamsUnion `json:"initial_events,omitzero" api:"required"`
	// Human-readable name for the deployment.
	Name string `json:"name" api:"required"`
	// Description of what the deployment does.
	Description param.Opt[string] `json:"description,omitzero"`
	// A hard spend ceiling. The session stops issuing new model requests once the
	// tracked list cost reaches `max_list_cost`.
	Budget ManagedAgentsBudgetLimitParam `json:"budget,omitzero"`
	// Arbitrary key-value metadata. Maximum 16 pairs, keys up to 64 chars, values up
	// to 512 chars.
	Metadata map[string]string `json:"metadata,omitzero"`
	// Resources (e.g. repositories, files) to mount into each session's container.
	// Maximum 500.
	Resources []DeploymentNewParamsResourceUnion `json:"resources,omitzero"`
	// 5-field POSIX cron schedule. Literal wall-clock matching in the configured
	// timezone.
	Schedule ManagedAgentsScheduleParams `json:"schedule,omitzero"`
	// Vault IDs for stored credentials the agent can use during sessions created from
	// this deployment. Maximum 50.
	VaultIDs []string `json:"vault_ids,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

func (r DeploymentNewParams) MarshalJSON() (data []byte, err error) {
	type shadow DeploymentNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DeploymentNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DeploymentNewParamsAgentUnion struct {
	OfString              param.Opt[string]         `json:",omitzero,inline"`
	OfManagedAgentsAgents *ManagedAgentsAgentParams `json:",omitzero,inline"`
	paramUnion
}

func (u DeploymentNewParamsAgentUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfManagedAgentsAgents)
}
func (u *DeploymentNewParamsAgentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DeploymentNewParamsAgentUnion) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfManagedAgentsAgents) {
		return u.OfManagedAgentsAgents
	}
	return nil
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DeploymentNewParamsResourceUnion struct {
	OfGitHubRepository *ManagedAgentsGitHubRepositoryResourceParams `json:",omitzero,inline"`
	OfFile             *ManagedAgentsFileResourceParams             `json:",omitzero,inline"`
	OfMemoryStore      *ManagedAgentsMemoryStoreResourceParam       `json:",omitzero,inline"`
	paramUnion
}

func (u DeploymentNewParamsResourceUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfGitHubRepository, u.OfFile, u.OfMemoryStore)
}
func (u *DeploymentNewParamsResourceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DeploymentNewParamsResourceUnion) asAny() any {
	if !param.IsOmitted(u.OfGitHubRepository) {
		return u.OfGitHubRepository
	} else if !param.IsOmitted(u.OfFile) {
		return u.OfFile
	} else if !param.IsOmitted(u.OfMemoryStore) {
		return u.OfMemoryStore
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DeploymentNewParamsResourceUnion) GetAuthorizationToken() *string {
	if vt := u.OfGitHubRepository; vt != nil {
		return &vt.AuthorizationToken
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DeploymentNewParamsResourceUnion) GetURL() *string {
	if vt := u.OfGitHubRepository; vt != nil {
		return &vt.URL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DeploymentNewParamsResourceUnion) GetCheckout() *ManagedAgentsGitHubRepositoryResourceParamsCheckoutUnion {
	if vt := u.OfGitHubRepository; vt != nil {
		return &vt.Checkout
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DeploymentNewParamsResourceUnion) GetFileID() *string {
	if vt := u.OfFile; vt != nil {
		return &vt.FileID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DeploymentNewParamsResourceUnion) GetMemoryStoreID() *string {
	if vt := u.OfMemoryStore; vt != nil {
		return &vt.MemoryStoreID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DeploymentNewParamsResourceUnion) GetAccess() *string {
	if vt := u.OfMemoryStore; vt != nil {
		return (*string)(&vt.Access)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DeploymentNewParamsResourceUnion) GetInstructions() *string {
	if vt := u.OfMemoryStore; vt != nil && vt.Instructions.Valid() {
		return &vt.Instructions.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DeploymentNewParamsResourceUnion) GetType() *string {
	if vt := u.OfGitHubRepository; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFile; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfMemoryStore; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DeploymentNewParamsResourceUnion) GetMountPath() *string {
	if vt := u.OfGitHubRepository; vt != nil && vt.MountPath.Valid() {
		return &vt.MountPath.Value
	} else if vt := u.OfFile; vt != nil && vt.MountPath.Valid() {
		return &vt.MountPath.Value
	}
	return nil
}

func init() {
	apijson.RegisterUnion[DeploymentNewParamsResourceUnion](
		"type",
		apijson.Discriminator[ManagedAgentsGitHubRepositoryResourceParams]("github_repository"),
		apijson.Discriminator[ManagedAgentsFileResourceParams]("file"),
		apijson.Discriminator[ManagedAgentsMemoryStoreResourceParam]("memory_store"),
	)
}

type DeploymentGetParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

type DeploymentUpdateParams struct {
	EnvironmentVariables param.Opt[string] `json:"environment_variables,omitzero"`

	// Description. Omit to preserve; send empty string or null to clear.
	Description param.Opt[string] `json:"description,omitzero"`
	// ID of the `environment` where sessions run. Omit to preserve. Cannot be cleared.
	EnvironmentID param.Opt[string] `json:"environment_id,omitzero"`
	// Human-readable name. Must be non-empty. Omit to preserve. Cannot be cleared.
	Name param.Opt[string] `json:"name,omitzero"`
	// Metadata patch. Set a key to a string to upsert it, or to null to delete it.
	// Omit the field to preserve. The stored bag is limited to 16 keys (up to 64 chars
	// each) with values up to 512 chars.
	Metadata map[string]any `json:"metadata,omitzero"`
	// Session resources. Full replacement. Omit to preserve; send empty array or null
	// to clear. Maximum 500.
	Resources []DeploymentUpdateParamsResourceUnion `json:"resources,omitzero"`
	// Vault IDs. Full replacement. Omit to preserve; send empty array or null to
	// clear. Maximum 50.
	VaultIDs []string `json:"vault_ids,omitzero"`
	// Agent to deploy. Accepts the `agent` ID string, which re-pins to the latest
	// version, or an `agent` object with both id and version specified. Omit to
	// preserve. Cannot be cleared.
	Agent DeploymentUpdateParamsAgentUnion `json:"agent,omitzero"`
	// A hard spend ceiling. The session stops issuing new model requests once the
	// tracked list cost reaches `max_list_cost`.
	Budget ManagedAgentsBudgetLimitParam `json:"budget,omitzero"`
	// Initial events. Full replacement. Omit to preserve. Cannot be cleared. At least
	// 1, maximum 50.
	InitialEvents []ManagedAgentsDeploymentInitialEventParamsUnion `json:"initial_events,omitzero"`
	// 5-field POSIX cron schedule. Literal wall-clock matching in the configured
	// timezone.
	Schedule ManagedAgentsScheduleParams `json:"schedule,omitzero"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

func (r DeploymentUpdateParams) MarshalJSON() (data []byte, err error) {
	type shadow DeploymentUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DeploymentUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DeploymentUpdateParamsAgentUnion struct {
	OfString              param.Opt[string]         `json:",omitzero,inline"`
	OfManagedAgentsAgents *ManagedAgentsAgentParams `json:",omitzero,inline"`
	paramUnion
}

func (u DeploymentUpdateParamsAgentUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfString, u.OfManagedAgentsAgents)
}
func (u *DeploymentUpdateParamsAgentUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DeploymentUpdateParamsAgentUnion) asAny() any {
	if !param.IsOmitted(u.OfString) {
		return &u.OfString.Value
	} else if !param.IsOmitted(u.OfManagedAgentsAgents) {
		return u.OfManagedAgentsAgents
	}
	return nil
}

// Only one field can be non-zero.
//
// Use [param.IsOmitted] to confirm if a field is set.
type DeploymentUpdateParamsResourceUnion struct {
	OfGitHubRepository *ManagedAgentsGitHubRepositoryResourceParams `json:",omitzero,inline"`
	OfFile             *ManagedAgentsFileResourceParams             `json:",omitzero,inline"`
	OfMemoryStore      *ManagedAgentsMemoryStoreResourceParam       `json:",omitzero,inline"`
	paramUnion
}

func (u DeploymentUpdateParamsResourceUnion) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(u, u.OfGitHubRepository, u.OfFile, u.OfMemoryStore)
}
func (u *DeploymentUpdateParamsResourceUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, u)
}

func (u *DeploymentUpdateParamsResourceUnion) asAny() any {
	if !param.IsOmitted(u.OfGitHubRepository) {
		return u.OfGitHubRepository
	} else if !param.IsOmitted(u.OfFile) {
		return u.OfFile
	} else if !param.IsOmitted(u.OfMemoryStore) {
		return u.OfMemoryStore
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DeploymentUpdateParamsResourceUnion) GetAuthorizationToken() *string {
	if vt := u.OfGitHubRepository; vt != nil {
		return &vt.AuthorizationToken
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DeploymentUpdateParamsResourceUnion) GetURL() *string {
	if vt := u.OfGitHubRepository; vt != nil {
		return &vt.URL
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DeploymentUpdateParamsResourceUnion) GetCheckout() *ManagedAgentsGitHubRepositoryResourceParamsCheckoutUnion {
	if vt := u.OfGitHubRepository; vt != nil {
		return &vt.Checkout
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DeploymentUpdateParamsResourceUnion) GetFileID() *string {
	if vt := u.OfFile; vt != nil {
		return &vt.FileID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DeploymentUpdateParamsResourceUnion) GetMemoryStoreID() *string {
	if vt := u.OfMemoryStore; vt != nil {
		return &vt.MemoryStoreID
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DeploymentUpdateParamsResourceUnion) GetAccess() *string {
	if vt := u.OfMemoryStore; vt != nil {
		return (*string)(&vt.Access)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DeploymentUpdateParamsResourceUnion) GetInstructions() *string {
	if vt := u.OfMemoryStore; vt != nil && vt.Instructions.Valid() {
		return &vt.Instructions.Value
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DeploymentUpdateParamsResourceUnion) GetType() *string {
	if vt := u.OfGitHubRepository; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfFile; vt != nil {
		return (*string)(&vt.Type)
	} else if vt := u.OfMemoryStore; vt != nil {
		return (*string)(&vt.Type)
	}
	return nil
}

// Returns a pointer to the underlying variant's property, if present.
func (u DeploymentUpdateParamsResourceUnion) GetMountPath() *string {
	if vt := u.OfGitHubRepository; vt != nil && vt.MountPath.Valid() {
		return &vt.MountPath.Value
	} else if vt := u.OfFile; vt != nil && vt.MountPath.Valid() {
		return &vt.MountPath.Value
	}
	return nil
}

func init() {
	apijson.RegisterUnion[DeploymentUpdateParamsResourceUnion](
		"type",
		apijson.Discriminator[ManagedAgentsGitHubRepositoryResourceParams]("github_repository"),
		apijson.Discriminator[ManagedAgentsFileResourceParams]("file"),
		apijson.Discriminator[ManagedAgentsMemoryStoreResourceParam]("memory_store"),
	)
}

type DeploymentListParams struct {
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	AfterID  param.Opt[string] `query:"after_id,omitzero" json:"-"`

	// Filter by agent ID.
	AgentID param.Opt[string] `query:"agent_id,omitzero" json:"-"`
	// Return deployments created at or after this time (inclusive).
	CreatedAtGte param.Opt[time.Time] `query:"created_at[gte],omitzero" format:"date-time" json:"-"`
	// Return deployments created at or before this time (inclusive).
	CreatedAtLte param.Opt[time.Time] `query:"created_at[lte],omitzero" format:"date-time" json:"-"`
	// When true, includes archived deployments. Default: false (exclude archived).
	IncludeArchived param.Opt[bool] `query:"include_archived,omitzero" json:"-"`
	// Maximum results per page. Default 20, maximum 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque pagination cursor.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Filter by status: `active` or `paused`. Omit for both. To include archived
	// deployments, use `include_archived` instead; the two cannot be combined.
	//
	// Any of "active", "paused".
	Status ManagedAgentsDeploymentStatus `query:"status,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [DeploymentListParams]'s query parameters as
// `url.Values`.
func (r DeploymentListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type DeploymentArchiveParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

type DeploymentPauseParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

type DeploymentRunParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

type DeploymentUnpauseParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}
