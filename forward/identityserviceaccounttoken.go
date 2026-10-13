package forward

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
)

// IdentityServiceAccountTokenService issues SATs bound to a Forward Identity.
type IdentityServiceAccountTokenService struct {
	Options []option.RequestOption
}

func NewIdentityServiceAccountTokenService(opts ...option.RequestOption) IdentityServiceAccountTokenService {
	return IdentityServiceAccountTokenService{Options: slices.Clone(opts)}
}

// New issues an Identity SAT using the explicit SA Key. Automatic retries are disabled.
func (r *IdentityServiceAccountTokenService) New(ctx context.Context, identityID string, params ServiceAccountTokenNewParams, opts ...option.RequestOption) (res *ServiceAccountToken, err error) {
	if identityID == "" {
		return nil, fmt.Errorf("missing required identity_id parameter")
	}
	opts, err = serviceAccountTokenOptions(params.SAKey, slices.Concat(r.Options, opts))
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("identities/%s/service_account_tokens", url.PathEscape(identityID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}
