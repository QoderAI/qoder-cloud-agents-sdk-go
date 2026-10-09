package managed

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apiquery"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/pagination"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
)

type DeploymentScopedRunService struct{ Options []option.RequestOption }

func NewDeploymentScopedRunService(opts ...option.RequestOption) DeploymentScopedRunService {
	return DeploymentScopedRunService{Options: slices.Clone(opts)}
}

// List DeploymentScopedRun.
func (r *DeploymentScopedRunService) List(ctx context.Context, deploymentID string, params DeploymentScopedRunListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ManagedAgentsDeploymentRun], err error) {
	if deploymentID == "" {
		return nil, fmt.Errorf("missing required deploymentID parameter")
	}
	if params.WorkspaceID.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("qoder-workspace-id", params.WorkspaceID.Value)}, opts...)
	}
	for _, beta := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", string(beta)))
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("deployments/%s/runs", url.PathEscape(deploymentID))
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
func (r *DeploymentScopedRunService) ListAutoPaging(ctx context.Context, deploymentID string, params DeploymentScopedRunListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ManagedAgentsDeploymentRun] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, deploymentID, params, opts...))
}

// Get DeploymentScopedRun.
func (r *DeploymentScopedRunService) Get(ctx context.Context, deploymentID string, runID string, params DeploymentScopedRunGetParams, opts ...option.RequestOption) (res *ManagedAgentsDeploymentRun, err error) {
	if deploymentID == "" {
		return nil, fmt.Errorf("missing required deploymentID parameter")
	}
	if runID == "" {
		return nil, fmt.Errorf("missing required runID parameter")
	}
	if params.WorkspaceID.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("qoder-workspace-id", params.WorkspaceID.Value)}, opts...)
	}
	for _, beta := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", string(beta)))
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("deployments/%s/runs/%s", url.PathEscape(deploymentID), url.PathEscape(runID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

type DeploymentScopedRunListParams struct {
	Limit           param.Opt[int64]  `query:"limit,omitzero" json:"-"`
	Page            param.Opt[string] `query:"page,omitzero" json:"-"`
	AfterID         param.Opt[string] `query:"after_id,omitzero" json:"-"`
	BeforeID        param.Opt[string] `query:"before_id,omitzero" json:"-"`
	TriggeredAfter  param.Opt[string] `query:"triggered_after,omitzero" json:"-"`
	TriggeredBefore param.Opt[string] `query:"triggered_before,omitzero" json:"-"`
	WorkspaceID     param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`
	Betas           []QoderBeta       `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

func (r *DeploymentScopedRunListParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
func (r DeploymentScopedRunListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

type DeploymentScopedRunGetParams struct {
	WorkspaceID param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`
	Betas       []QoderBeta       `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

func (r *DeploymentScopedRunGetParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
