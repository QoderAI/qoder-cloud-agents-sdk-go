package forward

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
)

// ServiceAccountTokenService issues admin SATs and revokes Forward SAT bindings.
// These methods require a per-call SA Key, rather than the client's PAT or SAT.
type ServiceAccountTokenService struct {
	Options []option.RequestOption
}

func NewServiceAccountTokenService(opts ...option.RequestOption) ServiceAccountTokenService {
	return ServiceAccountTokenService{Options: slices.Clone(opts)}
}

// New issues an admin SAT. Automatic retries are disabled to avoid duplicate issuance.
func (r *ServiceAccountTokenService) New(ctx context.Context, params ServiceAccountTokenNewParams, opts ...option.RequestOption) (res *ServiceAccountToken, err error) {
	opts, err = serviceAccountTokenOptions(params.SAKey, slices.Concat(r.Options, opts))
	if err != nil {
		return nil, err
	}
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, "service_account_tokens", params, &res, opts...)
	return res, err
}

// Revoke revokes a Forward SAT binding belonging to the same Service Account.
// Use the issuance response's CredentialID and ExpiresAt to disambiguate the target.
func (r *ServiceAccountTokenService) Revoke(ctx context.Context, authTokenID string, params ServiceAccountTokenRevokeParams, opts ...option.RequestOption) (res *ServiceAccountTokenRevokeResponse, err error) {
	if authTokenID == "" {
		return nil, fmt.Errorf("missing required auth_token_id parameter")
	}
	opts, err = serviceAccountTokenOptions(params.SAKey, slices.Concat(r.Options, opts))
	if err != nil {
		return nil, err
	}
	path := fmt.Sprintf("service_account_tokens/%s/revoke", url.PathEscape(authTokenID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

func serviceAccountTokenOptions(saKey string, opts []option.RequestOption) ([]option.RequestOption, error) {
	if saKey == "" || strings.ContainsFunc(saKey, func(r rune) bool { return r <= 0x20 || r == 0x7f }) {
		return nil, fmt.Errorf("sa_key must be non-empty and contain no whitespace or control characters")
	}
	// Apply last so client/request auth overrides cannot replace the explicit Key.
	return append(opts, option.WithHeader("Authorization", "Bearer "+saKey), option.WithHeader("Content-Type", "application/json"), option.WithMaxRetries(0)), nil
}

type ServiceAccountTokenNewParams struct {
	// SA Key used only in this request's Authorization header; never serialized as JSON.
	SAKey string `json:"-"`
	// Lifetime in seconds, 1-43200. Omit to use the issuer's default.
	TTLSeconds param.Opt[int64] `json:"ttl_seconds,omitzero"`
	// Business audit metadata; do not include credentials or other secrets.
	Metadata map[string]any `json:"metadata,omitzero"`
	paramObj
}

func (r ServiceAccountTokenNewParams) MarshalJSON() ([]byte, error) {
	type shadow ServiceAccountTokenNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ServiceAccountTokenNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
func (r ServiceAccountTokenNewParams) String() string {
	return "ServiceAccountTokenNewParams{SAKey: [redacted]}"
}
func (r ServiceAccountTokenNewParams) GoString() string { return r.String() }

type ServiceAccountTokenRevokeParams struct {
	SAKey              string               `json:"-"`
	TargetCredentialID param.Opt[string]    `json:"target_credential_id,omitzero"`
	TargetExpiresAt    param.Opt[time.Time] `json:"target_expires_at,omitzero" format:"date-time"`
	Reason             param.Opt[string]    `json:"reason,omitzero"`
	paramObj
}

func (r ServiceAccountTokenRevokeParams) MarshalJSON() ([]byte, error) {
	type shadow ServiceAccountTokenRevokeParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ServiceAccountTokenRevokeParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
func (r ServiceAccountTokenRevokeParams) String() string {
	return "ServiceAccountTokenRevokeParams{SAKey: [redacted]}"
}
func (r ServiceAccountTokenRevokeParams) GoString() string { return r.String() }

// ServiceAccountToken contains a newly issued Forward SAT. Treat AccessToken and RawJSON as secrets.
// Use WithAuthMode("sa") and WithAuthToken(AccessToken) for business calls; renewal is caller-owned.
type ServiceAccountToken struct {
	Type             string    `json:"type"`
	AccessToken      string    `json:"access_token"`
	TokenType        string    `json:"token_type"`
	ExpiresIn        int64     `json:"expires_in"`
	ExpiresAt        time.Time `json:"expires_at" format:"date-time"`
	AuthTokenID      string    `json:"auth_token_id"`
	ServiceAccountID string    `json:"service_account_id"`
	CredentialID     string    `json:"credential_id"`
	// admin or identity. IdentityID is absent for admin tokens.
	SubjectType string `json:"subject_type"`
	IdentityID  string `json:"identity_id"`
	JSON        struct {
		Type             respjson.Field
		AccessToken      respjson.Field
		TokenType        respjson.Field
		ExpiresIn        respjson.Field
		ExpiresAt        respjson.Field
		AuthTokenID      respjson.Field
		ServiceAccountID respjson.Field
		CredentialID     respjson.Field
		SubjectType      respjson.Field
		IdentityID       respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

func (r ServiceAccountToken) RawJSON() string                  { return r.JSON.raw }
func (r *ServiceAccountToken) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
func (r ServiceAccountToken) String() string {
	return fmt.Sprintf("ServiceAccountToken{AuthTokenID: %q, SubjectType: %q, AccessToken: [redacted]}", r.AuthTokenID, r.SubjectType)
}
func (r ServiceAccountToken) GoString() string { return r.String() }

type ServiceAccountTokenRevokeResponse struct {
	Revoked bool `json:"revoked"`
	JSON    struct {
		Revoked     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r ServiceAccountTokenRevokeResponse) RawJSON() string { return r.JSON.raw }
func (r *ServiceAccountTokenRevokeResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
