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

// DeploymentRunService contains methods and other services that help with
// interacting with the Qoder Cloud Agents API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewDeploymentRunService] method instead.
type DeploymentRunService struct {
	Options []option.RequestOption
}

// NewDeploymentRunService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewDeploymentRunService(opts ...option.RequestOption) (r DeploymentRunService) {
	r = DeploymentRunService{}
	r.Options = opts
	return
}

// Get Deployment Run
func (r *DeploymentRunService) Get(ctx context.Context, deploymentRunID string, query DeploymentRunGetParams, opts ...option.RequestOption) (res *ManagedAgentsDeploymentRun, err error) {
	for _, v := range query.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)

	if deploymentRunID == "" {
		err = errors.New("missing required deployment_run_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("deployment_runs/%s", url.PathEscape(deploymentRunID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List Deployment Runs
func (r *DeploymentRunService) List(ctx context.Context, params DeploymentRunListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ManagedAgentsDeploymentRun], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "deployment_runs"
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

// List Deployment Runs
func (r *DeploymentRunService) ListAutoPaging(ctx context.Context, params DeploymentRunListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ManagedAgentsDeploymentRun] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, params, opts...))
}

// The deployment's agent was archived.
type ManagedAgentsAgentArchivedRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "agent_archived_error".
	Type ManagedAgentsAgentArchivedRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsAgentArchivedRunError) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsAgentArchivedRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsAgentArchivedRunErrorType string

const (
	ManagedAgentsAgentArchivedRunErrorTypeAgentArchivedError ManagedAgentsAgentArchivedRunErrorType = "agent_archived_error"
)

// A persistent, append-only record of a single deployment execution. Records
// session creation success or failure — no session lifecycle tracking.
type ManagedAgentsDeploymentRun struct {
	// Unique identifier for this run (`drun_...`).
	ID string `json:"id" api:"required"`
	// A resolved agent reference with a concrete version.
	Agent ManagedAgentsAgentReference `json:"agent" api:"required"`
	// A timestamp in RFC 3339 format
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	// ID of the deployment that produced this run.
	DeploymentID string `json:"deployment_id" api:"required"`
	// Why the run failed to create a session. The type identifies the failure; message
	// is human-readable detail.
	Error ManagedAgentsDeploymentRunErrorUnion `json:"error" api:"required"`
	// Populated on success. Null on creation failure. Exactly one of `session_id` or
	// `error` is non-null.
	SessionID string `json:"session_id" api:"required"`
	// Describes what triggered a deployment run, with trigger-specific metadata.
	TriggerContext ManagedAgentsTriggerContextUnion `json:"trigger_context" api:"required"`
	// Any of "deployment_run".
	Type ManagedAgentsDeploymentRunType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ID             respjson.Field
		Agent          respjson.Field
		CreatedAt      respjson.Field
		DeploymentID   respjson.Field
		Error          respjson.Field
		SessionID      respjson.Field
		TriggerContext respjson.Field
		Type           respjson.Field
		ExtraFields    map[string]respjson.Field
		raw            string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsDeploymentRun) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsDeploymentRun) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ManagedAgentsDeploymentRunErrorUnion contains all possible properties and
// values from [ManagedAgentsEnvironmentArchivedRunError],
// [ManagedAgentsAgentArchivedRunError],
// [ManagedAgentsEnvironmentNotFoundRunError],
// [ManagedAgentsVaultNotFoundRunError],
// [ManagedAgentsVaultArchivedRunError],
// [ManagedAgentsFileNotFoundRunError],
// [ManagedAgentsMemoryStoreArchivedRunError],
// [ManagedAgentsSkillNotFoundRunError],
// [ManagedAgentsSessionResourceNotFoundRunError],
// [ManagedAgentsWorkspaceArchivedRunError],
// [ManagedAgentsOrganizationDisabledRunError],
// [ManagedAgentsSessionRateLimitedRunError],
// [ManagedAgentsSessionCreationRejectedRunError],
// [ManagedAgentsUnknownRunError],
// [ManagedAgentsSelfHostedResourcesUnsupportedRunError],
// [ManagedAgentsMCPEgressBlockedRunError].
//
// Use the [ManagedAgentsDeploymentRunErrorUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsDeploymentRunErrorUnion struct {
	Message string `json:"message"`
	// Any of "environment_archived_error", "agent_archived_error",
	// "environment_not_found_error", "vault_not_found_error", "vault_archived_error",
	// "file_not_found_error", "memory_store_archived_error", "skill_not_found_error",
	// "session_resource_not_found_error", "workspace_archived_error",
	// "organization_disabled_error", "session_rate_limited_error",
	// "session_creation_rejected_error", "unknown_error",
	// "self_hosted_resources_unsupported_error", "mcp_egress_blocked_error".
	Type string `json:"type"`
	JSON struct {
		Message respjson.Field
		Type    respjson.Field
		raw     string
	} `json:"-"`
}

// anyManagedAgentsDeploymentRunError is implemented by each variant of
// [ManagedAgentsDeploymentRunErrorUnion] to add type safety for the return
// type of [ManagedAgentsDeploymentRunErrorUnion.AsAny]
type anyManagedAgentsDeploymentRunError interface {
	implManagedAgentsDeploymentRunErrorUnion()
}

func (ManagedAgentsEnvironmentArchivedRunError) implManagedAgentsDeploymentRunErrorUnion() {}
func (ManagedAgentsAgentArchivedRunError) implManagedAgentsDeploymentRunErrorUnion()       {}
func (ManagedAgentsEnvironmentNotFoundRunError) implManagedAgentsDeploymentRunErrorUnion() {}
func (ManagedAgentsVaultNotFoundRunError) implManagedAgentsDeploymentRunErrorUnion()       {}
func (ManagedAgentsVaultArchivedRunError) implManagedAgentsDeploymentRunErrorUnion()       {}
func (ManagedAgentsFileNotFoundRunError) implManagedAgentsDeploymentRunErrorUnion()        {}
func (ManagedAgentsMemoryStoreArchivedRunError) implManagedAgentsDeploymentRunErrorUnion() {}
func (ManagedAgentsSkillNotFoundRunError) implManagedAgentsDeploymentRunErrorUnion()       {}
func (ManagedAgentsSessionResourceNotFoundRunError) implManagedAgentsDeploymentRunErrorUnion() {
}
func (ManagedAgentsWorkspaceArchivedRunError) implManagedAgentsDeploymentRunErrorUnion()    {}
func (ManagedAgentsOrganizationDisabledRunError) implManagedAgentsDeploymentRunErrorUnion() {}
func (ManagedAgentsSessionRateLimitedRunError) implManagedAgentsDeploymentRunErrorUnion()   {}
func (ManagedAgentsSessionCreationRejectedRunError) implManagedAgentsDeploymentRunErrorUnion() {
}
func (ManagedAgentsUnknownRunError) implManagedAgentsDeploymentRunErrorUnion() {}
func (ManagedAgentsSelfHostedResourcesUnsupportedRunError) implManagedAgentsDeploymentRunErrorUnion() {
}
func (ManagedAgentsMCPEgressBlockedRunError) implManagedAgentsDeploymentRunErrorUnion() {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsDeploymentRunErrorUnion.AsAny().(type) {
//	case qoder.ManagedAgentsEnvironmentArchivedRunError:
//	case qoder.ManagedAgentsAgentArchivedRunError:
//	case qoder.ManagedAgentsEnvironmentNotFoundRunError:
//	case qoder.ManagedAgentsVaultNotFoundRunError:
//	case qoder.ManagedAgentsVaultArchivedRunError:
//	case qoder.ManagedAgentsFileNotFoundRunError:
//	case qoder.ManagedAgentsMemoryStoreArchivedRunError:
//	case qoder.ManagedAgentsSkillNotFoundRunError:
//	case qoder.ManagedAgentsSessionResourceNotFoundRunError:
//	case qoder.ManagedAgentsWorkspaceArchivedRunError:
//	case qoder.ManagedAgentsOrganizationDisabledRunError:
//	case qoder.ManagedAgentsSessionRateLimitedRunError:
//	case qoder.ManagedAgentsSessionCreationRejectedRunError:
//	case qoder.ManagedAgentsUnknownRunError:
//	case qoder.ManagedAgentsSelfHostedResourcesUnsupportedRunError:
//	case qoder.ManagedAgentsMCPEgressBlockedRunError:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsDeploymentRunErrorUnion) AsAny() anyManagedAgentsDeploymentRunError {
	switch u.Type {
	case "environment_archived_error":
		return u.AsEnvironmentArchivedError()
	case "agent_archived_error":
		return u.AsAgentArchivedError()
	case "environment_not_found_error":
		return u.AsEnvironmentNotFoundError()
	case "vault_not_found_error":
		return u.AsVaultNotFoundError()
	case "vault_archived_error":
		return u.AsVaultArchivedError()
	case "file_not_found_error":
		return u.AsFileNotFoundError()
	case "memory_store_archived_error":
		return u.AsMemoryStoreArchivedError()
	case "skill_not_found_error":
		return u.AsSkillNotFoundError()
	case "session_resource_not_found_error":
		return u.AsSessionResourceNotFoundError()
	case "workspace_archived_error":
		return u.AsWorkspaceArchivedError()
	case "organization_disabled_error":
		return u.AsOrganizationDisabledError()
	case "session_rate_limited_error":
		return u.AsSessionRateLimitedError()
	case "session_creation_rejected_error":
		return u.AsSessionCreationRejectedError()
	case "unknown_error":
		return u.AsUnknownError()
	case "self_hosted_resources_unsupported_error":
		return u.AsSelfHostedResourcesUnsupportedError()
	case "mcp_egress_blocked_error":
		return u.AsMCPEgressBlockedError()
	}
	return nil
}

func (u ManagedAgentsDeploymentRunErrorUnion) AsEnvironmentArchivedError() (v ManagedAgentsEnvironmentArchivedRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentRunErrorUnion) AsAgentArchivedError() (v ManagedAgentsAgentArchivedRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentRunErrorUnion) AsEnvironmentNotFoundError() (v ManagedAgentsEnvironmentNotFoundRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentRunErrorUnion) AsVaultNotFoundError() (v ManagedAgentsVaultNotFoundRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentRunErrorUnion) AsVaultArchivedError() (v ManagedAgentsVaultArchivedRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentRunErrorUnion) AsFileNotFoundError() (v ManagedAgentsFileNotFoundRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentRunErrorUnion) AsMemoryStoreArchivedError() (v ManagedAgentsMemoryStoreArchivedRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentRunErrorUnion) AsSkillNotFoundError() (v ManagedAgentsSkillNotFoundRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentRunErrorUnion) AsSessionResourceNotFoundError() (v ManagedAgentsSessionResourceNotFoundRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentRunErrorUnion) AsWorkspaceArchivedError() (v ManagedAgentsWorkspaceArchivedRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentRunErrorUnion) AsOrganizationDisabledError() (v ManagedAgentsOrganizationDisabledRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentRunErrorUnion) AsSessionRateLimitedError() (v ManagedAgentsSessionRateLimitedRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentRunErrorUnion) AsSessionCreationRejectedError() (v ManagedAgentsSessionCreationRejectedRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentRunErrorUnion) AsUnknownError() (v ManagedAgentsUnknownRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentRunErrorUnion) AsSelfHostedResourcesUnsupportedError() (v ManagedAgentsSelfHostedResourcesUnsupportedRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsDeploymentRunErrorUnion) AsMCPEgressBlockedError() (v ManagedAgentsMCPEgressBlockedRunError) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsDeploymentRunErrorUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsDeploymentRunErrorUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsDeploymentRunType string

const (
	ManagedAgentsDeploymentRunTypeDeploymentRun ManagedAgentsDeploymentRunType = "deployment_run"
)

// The deployment's environment was archived.
type ManagedAgentsEnvironmentArchivedRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "environment_archived_error".
	Type ManagedAgentsEnvironmentArchivedRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsEnvironmentArchivedRunError) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsEnvironmentArchivedRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsEnvironmentArchivedRunErrorType string

const (
	ManagedAgentsEnvironmentArchivedRunErrorTypeEnvironmentArchivedError ManagedAgentsEnvironmentArchivedRunErrorType = "environment_archived_error"
)

// The deployment's environment no longer exists.
type ManagedAgentsEnvironmentNotFoundRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "environment_not_found_error".
	Type ManagedAgentsEnvironmentNotFoundRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsEnvironmentNotFoundRunError) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsEnvironmentNotFoundRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsEnvironmentNotFoundRunErrorType string

const (
	ManagedAgentsEnvironmentNotFoundRunErrorTypeEnvironmentNotFoundError ManagedAgentsEnvironmentNotFoundRunErrorType = "environment_not_found_error"
)

// A file resource referenced by the deployment no longer exists.
type ManagedAgentsFileNotFoundRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "file_not_found_error".
	Type ManagedAgentsFileNotFoundRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsFileNotFoundRunError) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsFileNotFoundRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsFileNotFoundRunErrorType string

const (
	ManagedAgentsFileNotFoundRunErrorTypeFileNotFoundError ManagedAgentsFileNotFoundRunErrorType = "file_not_found_error"
)

// The run was started manually by creating a session directly against the
// deployment.
type ManagedAgentsManualTriggerContext struct {
	// Any of "manual".
	Type ManagedAgentsManualTriggerContextType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsManualTriggerContext) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsManualTriggerContext) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsManualTriggerContextType string

const (
	ManagedAgentsManualTriggerContextTypeManual ManagedAgentsManualTriggerContextType = "manual"
)

// An MCP server host used by the deployment's agent is blocked by the
// environment's network policy.
type ManagedAgentsMCPEgressBlockedRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "mcp_egress_blocked_error".
	Type ManagedAgentsMCPEgressBlockedRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsMCPEgressBlockedRunError) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsMCPEgressBlockedRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsMCPEgressBlockedRunErrorType string

const (
	ManagedAgentsMCPEgressBlockedRunErrorTypeMCPEgressBlockedError ManagedAgentsMCPEgressBlockedRunErrorType = "mcp_egress_blocked_error"
)

// A memory store referenced by the deployment is archived.
type ManagedAgentsMemoryStoreArchivedRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "memory_store_archived_error".
	Type ManagedAgentsMemoryStoreArchivedRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsMemoryStoreArchivedRunError) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsMemoryStoreArchivedRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsMemoryStoreArchivedRunErrorType string

const (
	ManagedAgentsMemoryStoreArchivedRunErrorTypeMemoryStoreArchivedError ManagedAgentsMemoryStoreArchivedRunErrorType = "memory_store_archived_error"
)

// The deployment's organization is disabled.
type ManagedAgentsOrganizationDisabledRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "organization_disabled_error".
	Type ManagedAgentsOrganizationDisabledRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsOrganizationDisabledRunError) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsOrganizationDisabledRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsOrganizationDisabledRunErrorType string

const (
	ManagedAgentsOrganizationDisabledRunErrorTypeOrganizationDisabledError ManagedAgentsOrganizationDisabledRunErrorType = "organization_disabled_error"
)

// The run was fired by the deployment's cron schedule.
type ManagedAgentsScheduleTriggerContext struct {
	// A timestamp in RFC 3339 format
	ScheduledAt time.Time `json:"scheduled_at" api:"required" format:"date-time"`
	// Any of "schedule".
	Type ManagedAgentsScheduleTriggerContextType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		ScheduledAt respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsScheduleTriggerContext) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsScheduleTriggerContext) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsScheduleTriggerContextType string

const (
	ManagedAgentsScheduleTriggerContextTypeSchedule ManagedAgentsScheduleTriggerContextType = "schedule"
)

// The deployment configures resources, but its environment is self-hosted and
// cannot mount them.
type ManagedAgentsSelfHostedResourcesUnsupportedRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "self_hosted_resources_unsupported_error".
	Type ManagedAgentsSelfHostedResourcesUnsupportedRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSelfHostedResourcesUnsupportedRunError) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSelfHostedResourcesUnsupportedRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSelfHostedResourcesUnsupportedRunErrorType string

const (
	ManagedAgentsSelfHostedResourcesUnsupportedRunErrorTypeSelfHostedResourcesUnsupportedError ManagedAgentsSelfHostedResourcesUnsupportedRunErrorType = "self_hosted_resources_unsupported_error"
)

// The session create request was rejected with a non-retryable validation error.
type ManagedAgentsSessionCreationRejectedRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "session_creation_rejected_error".
	Type ManagedAgentsSessionCreationRejectedRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionCreationRejectedRunError) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionCreationRejectedRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSessionCreationRejectedRunErrorType string

const (
	ManagedAgentsSessionCreationRejectedRunErrorTypeSessionCreationRejectedError ManagedAgentsSessionCreationRejectedRunErrorType = "session_creation_rejected_error"
)

// Session creation was rejected due to rate limiting. The schedule keeps firing;
// subsequent runs may succeed.
type ManagedAgentsSessionRateLimitedRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "session_rate_limited_error".
	Type ManagedAgentsSessionRateLimitedRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionRateLimitedRunError) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionRateLimitedRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSessionRateLimitedRunErrorType string

const (
	ManagedAgentsSessionRateLimitedRunErrorTypeSessionRateLimitedError ManagedAgentsSessionRateLimitedRunErrorType = "session_rate_limited_error"
)

// A referenced resource no longer exists and its kind was not reported.
type ManagedAgentsSessionResourceNotFoundRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "session_resource_not_found_error".
	Type ManagedAgentsSessionResourceNotFoundRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSessionResourceNotFoundRunError) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSessionResourceNotFoundRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSessionResourceNotFoundRunErrorType string

const (
	ManagedAgentsSessionResourceNotFoundRunErrorTypeSessionResourceNotFoundError ManagedAgentsSessionResourceNotFoundRunErrorType = "session_resource_not_found_error"
)

// A skill referenced by the deployment's agent no longer exists.
type ManagedAgentsSkillNotFoundRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "skill_not_found_error".
	Type ManagedAgentsSkillNotFoundRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsSkillNotFoundRunError) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsSkillNotFoundRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsSkillNotFoundRunErrorType string

const (
	ManagedAgentsSkillNotFoundRunErrorTypeSkillNotFoundError ManagedAgentsSkillNotFoundRunErrorType = "skill_not_found_error"
)

// ManagedAgentsTriggerContextUnion contains all possible properties and values
// from [ManagedAgentsScheduleTriggerContext],
// [ManagedAgentsManualTriggerContext].
//
// Use the [ManagedAgentsTriggerContextUnion.AsAny] method to switch on the
// variant.
//
// Use the methods beginning with 'As' to cast the union to one of its variants.
type ManagedAgentsTriggerContextUnion struct {
	// This field is from variant [ManagedAgentsScheduleTriggerContext].
	ScheduledAt time.Time `json:"scheduled_at"`
	// Any of "schedule", "manual".
	Type string `json:"type"`
	JSON struct {
		ScheduledAt respjson.Field
		Type        respjson.Field
		raw         string
	} `json:"-"`
}

// anyManagedAgentsTriggerContext is implemented by each variant of
// [ManagedAgentsTriggerContextUnion] to add type safety for the return type of
// [ManagedAgentsTriggerContextUnion.AsAny]
type anyManagedAgentsTriggerContext interface {
	implManagedAgentsTriggerContextUnion()
}

func (ManagedAgentsScheduleTriggerContext) implManagedAgentsTriggerContextUnion() {}
func (ManagedAgentsManualTriggerContext) implManagedAgentsTriggerContextUnion()   {}

// Use the following switch statement to find the correct variant
//
//	switch variant := ManagedAgentsTriggerContextUnion.AsAny().(type) {
//	case qoder.ManagedAgentsScheduleTriggerContext:
//	case qoder.ManagedAgentsManualTriggerContext:
//	default:
//	  fmt.Errorf("no variant present")
//	}
func (u ManagedAgentsTriggerContextUnion) AsAny() anyManagedAgentsTriggerContext {
	switch u.Type {
	case "schedule":
		return u.AsSchedule()
	case "manual":
		return u.AsManual()
	}
	return nil
}

func (u ManagedAgentsTriggerContextUnion) AsSchedule() (v ManagedAgentsScheduleTriggerContext) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func (u ManagedAgentsTriggerContextUnion) AsManual() (v ManagedAgentsManualTriggerContext) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

// Returns the unmodified JSON received from the API
func (u ManagedAgentsTriggerContextUnion) RawJSON() string { return u.JSON.raw }

func (r *ManagedAgentsTriggerContextUnion) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// What triggered a deployment run.
type ManagedAgentsTriggerType string

const (
	ManagedAgentsTriggerTypeSchedule ManagedAgentsTriggerType = "schedule"
	ManagedAgentsTriggerTypeManual   ManagedAgentsTriggerType = "manual"
)

// An unknown or unexpected error caused the run to fail. A fallback variant;
// clients that do not recognize a new error type can match on message alone.
type ManagedAgentsUnknownRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "unknown_error".
	Type ManagedAgentsUnknownRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsUnknownRunError) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsUnknownRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsUnknownRunErrorType string

const (
	ManagedAgentsUnknownRunErrorTypeUnknownError ManagedAgentsUnknownRunErrorType = "unknown_error"
)

// A vault referenced by the deployment is archived.
type ManagedAgentsVaultArchivedRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "vault_archived_error".
	Type ManagedAgentsVaultArchivedRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsVaultArchivedRunError) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsVaultArchivedRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsVaultArchivedRunErrorType string

const (
	ManagedAgentsVaultArchivedRunErrorTypeVaultArchivedError ManagedAgentsVaultArchivedRunErrorType = "vault_archived_error"
)

// A vault referenced by the deployment no longer exists.
type ManagedAgentsVaultNotFoundRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "vault_not_found_error".
	Type ManagedAgentsVaultNotFoundRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsVaultNotFoundRunError) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsVaultNotFoundRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsVaultNotFoundRunErrorType string

const (
	ManagedAgentsVaultNotFoundRunErrorTypeVaultNotFoundError ManagedAgentsVaultNotFoundRunErrorType = "vault_not_found_error"
)

// The deployment's workspace was archived.
type ManagedAgentsWorkspaceArchivedRunError struct {
	// Human-readable error description.
	Message string `json:"message" api:"required"`
	// Any of "workspace_archived_error".
	Type ManagedAgentsWorkspaceArchivedRunErrorType `json:"type" api:"required"`
	// JSON contains metadata for fields, check presence with [respjson.Field.Valid].
	JSON struct {
		Message     respjson.Field
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

// Returns the unmodified JSON received from the API
func (r ManagedAgentsWorkspaceArchivedRunError) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsWorkspaceArchivedRunError) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsWorkspaceArchivedRunErrorType string

const (
	ManagedAgentsWorkspaceArchivedRunErrorTypeWorkspaceArchivedError ManagedAgentsWorkspaceArchivedRunErrorType = "workspace_archived_error"
)

type DeploymentRunGetParams struct {
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

type DeploymentRunListParams struct {
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	AfterID  param.Opt[string] `query:"after_id,omitzero" json:"-"`

	// Return runs created strictly after this time (exclusive).
	CreatedAtGt param.Opt[time.Time] `query:"created_at[gt],omitzero" format:"date-time" json:"-"`
	// Return runs created at or after this time (inclusive).
	CreatedAtGte param.Opt[time.Time] `query:"created_at[gte],omitzero" format:"date-time" json:"-"`
	// Return runs created strictly before this time (exclusive).
	CreatedAtLt param.Opt[time.Time] `query:"created_at[lt],omitzero" format:"date-time" json:"-"`
	// Return runs created at or before this time (inclusive).
	CreatedAtLte param.Opt[time.Time] `query:"created_at[lte],omitzero" format:"date-time" json:"-"`
	// Filter to a specific deployment. Omit to list across all deployments in the
	// workspace. Filtering by a non-existent `deployment_id` returns 200 with empty
	// data.
	DeploymentID param.Opt[string] `query:"deployment_id,omitzero" json:"-"`
	// Filter: true for runs with non-null `error`, false for runs with non-null
	// `session_id`. Omit for all.
	HasError param.Opt[bool] `query:"has_error,omitzero" json:"-"`
	// Maximum results per page. Default 20, maximum 1000.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque pagination cursor. Pass `next_page` from the previous response. Invalid
	// or expired cursors return 400.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Filter runs by what triggered them. Omit to return all runs.
	//
	// Any of "schedule", "manual".
	TriggerType ManagedAgentsTriggerType `query:"trigger_type,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [DeploymentRunListParams]'s query parameters as
// `url.Values`.
func (r DeploymentRunListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
