package managed

import (
	"context"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
	"net/http"
	"slices"
)

type OAuthService struct{ Options []option.RequestOption }

func NewOAuthService(opts ...option.RequestOption) OAuthService {
	return OAuthService{Options: slices.Clone(opts)}
}

// Start returns an authorization URL; the user completes authorization in their browser.
func (r *OAuthService) Start(ctx context.Context, params OAuthStartParams, opts ...option.RequestOption) (res *OAuthStartResponse, err error) {
	if params.WorkspaceID.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("qoder-workspace-id", params.WorkspaceID.Value)}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, "oauth/start", params, &res, opts...)
	return res, err
}

type OAuthStartParams struct {
	WorkspaceID param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`

	VaultID      string            `json:"vault_id"`
	MCPServerURL string            `json:"mcp_server_url"`
	ClientID     param.Opt[string] `json:"client_id,omitzero"`
	ClientSecret param.Opt[string] `json:"client_secret,omitzero"`
	paramObj
}

func (r OAuthStartParams) MarshalJSON() ([]byte, error) {
	type shadow OAuthStartParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *OAuthStartParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type OAuthStartResponse struct {
	AuthorizationURL string `json:"authorization_url"`
	State            string `json:"state"`
	CallbackOrigin   string `json:"callback_origin"`
	JSON             struct {
		AuthorizationURL respjson.Field
		State            respjson.Field
		CallbackOrigin   respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

func (r OAuthStartResponse) RawJSON() string                  { return r.JSON.raw }
func (r *OAuthStartResponse) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
