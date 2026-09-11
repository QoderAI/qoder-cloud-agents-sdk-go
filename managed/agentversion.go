// Qoder managed API definitions.
package managed

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	requestconfig "github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apiquery"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/pagination"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
)

// AgentVersionService contains methods and other services that help with
// interacting with the Qoder Cloud Agents API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAgentVersionService] method instead.
type AgentVersionService struct {
	Options []option.RequestOption
}

// NewAgentVersionService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewAgentVersionService(opts ...option.RequestOption) (r AgentVersionService) {
	r = AgentVersionService{}
	r.Options = opts
	return
}

// List Agent Versions
func (r *AgentVersionService) List(ctx context.Context, agentID string, params AgentVersionListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ManagedAgentsAgent], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if agentID == "" {
		err = errors.New("missing required agent_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("agents/%s/versions", url.PathEscape(agentID))
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

// List Agent Versions
func (r *AgentVersionService) ListAutoPaging(ctx context.Context, agentID string, params AgentVersionListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ManagedAgentsAgent] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, agentID, params, opts...))
}

type AgentVersionListParams struct {
	// Maximum results per page. Default 20, maximum 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Opaque pagination cursor.
	Page param.Opt[string] `query:"page,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [AgentVersionListParams]'s query parameters as
// `url.Values`.
func (r AgentVersionListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
