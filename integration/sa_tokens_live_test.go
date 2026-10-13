//go:build live

package integration_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apierror"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/internal/testsupport"
)

func saTokenDenied(t *testing.T, err error, statuses ...int) {
	t.Helper()
	var api *apierror.Error
	if errors.As(err, &api) {
		for _, status := range statuses {
			if api.StatusCode == status {
				return
			}
		}
	}
	t.Fatalf("expected access refusal: %s", testsupport.SafeError(err))
}

func saTokenMetadata(t *testing.T, token *forward.ServiceAccountToken, subject, identityID string) {
	t.Helper()
	if token == nil || token.AccessToken == "" || token.Type != "service_account_token" || token.TokenType != "Bearer" ||
		token.AuthTokenID == "" || token.CredentialID == "" || token.ServiceAccountID == "" ||
		token.SubjectType != subject || token.IdentityID != identityID || token.ExpiresIn <= 0 || token.ExpiresIn > 600 ||
		!token.ExpiresAt.After(time.Now()) {
		t.Fatal("issued SAT metadata did not match the requested subject and lifetime")
	}
	t.Logf("issued subject=%s auth_token_id=%s credential_id=%s identity_id=%s expires_at=%s", subject, token.AuthTokenID, token.CredentialID, token.IdentityID, token.ExpiresAt.UTC().Format(time.RFC3339))
}

func TestForwardServiceAccountTokensLive(t *testing.T) {
	saRequireConfiguration(t, false)
	if os.Getenv("QODER_SA_LIVE_ALLOW_WRITE") != "true" {
		t.Fatal("QODER_SA_LIVE_ALLOW_WRITE=true is required")
	}
	for _, region := range []string{"cn", "global"} {
		t.Run(region, func(t *testing.T) {
			variable, host := "QODER_SA_GLOBAL_KEY", "https://api.qoder.com"
			if region == "cn" {
				variable, host = "QODER_SA_CN_KEY", "https://api.qoder.com.cn"
			}
			key := os.Getenv(variable)
			opts := []option.RequestOption{option.WithBaseURL(host + "/api/v1/forward/"), option.WithAuthMode("sa"), option.WithMaxRetries(0), option.WithRequestTimeout(30 * time.Second), testsupport.RequestLog(t)}
			issuer := forward.NewClient(opts...)
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Minute)
			defer cancel()
			params := forward.ServiceAccountTokenNewParams{SAKey: key, TTLSeconds: forward.Int(600), Metadata: map[string]any{"purpose": "sdk verification", "sdk": "go"}}
			adminToken, err := issuer.ServiceAccountTokens.New(ctx, params)
			saCheck(t, err)
			admin := forward.NewClient(append(opts, option.WithAuthToken(adminToken.AccessToken))...)
			revoke := func(ctx context.Context, token *forward.ServiceAccountToken) error {
				result, err := issuer.ServiceAccountTokens.Revoke(ctx, token.AuthTokenID, forward.ServiceAccountTokenRevokeParams{SAKey: key, TargetCredentialID: forward.String(token.CredentialID), TargetExpiresAt: param.NewOpt(token.ExpiresAt), Reason: forward.String("sdk verification complete")})
				if err == nil && (result == nil || !result.Revoked) {
					return errors.New("revoke did not acknowledge success")
				}
				return err
			}
			saCleanup(t, "admin SAT "+adminToken.AuthTokenID, func(ctx context.Context) error {
				if err := revoke(ctx, adminToken); err != nil {
					return err
				}
				_, err := admin.Templates.List(ctx, forward.TemplateListParams{Limit: forward.Int(1)})
				saTokenDenied(t, err, 401, 403)
				return nil
			})
			saTokenMetadata(t, adminToken, "admin", "")
			_, err = admin.Templates.List(ctx, forward.TemplateListParams{Limit: forward.Int(1)})
			saCheck(t, err)
			identities := make([]string, 0, 2)
			for i := 0; i < 2; i++ {
				identity, err := admin.Identities.New(ctx, forward.IdentityNewParams{ExternalID: "sdk-go-sat-" + testsupport.Marker(t)})
				saCheck(t, err)
				id := identity.ID
				saCleanup(t, "identity "+id, func(ctx context.Context) error { _, err := admin.Identities.Delete(ctx, id); return err })
				if id == "" {
					t.Fatal("created Identity must have an ID")
				}
				identities = append(identities, id)
			}
			identityToken, err := issuer.Identities.ServiceAccountTokens.New(ctx, identities[0], params)
			saCheck(t, err)
			runtime := forward.NewClient(append(opts, option.WithAuthToken(identityToken.AccessToken))...)
			saCleanup(t, "identity SAT "+identityToken.AuthTokenID, func(ctx context.Context) error {
				if err := revoke(ctx, identityToken); err != nil {
					return err
				}
				_, err := runtime.Identities.Configs.List(ctx, identities[0], forward.IdentityConfigListParams{Limit: forward.Int(1)})
				saTokenDenied(t, err, 401, 403)
				return nil
			})
			saTokenMetadata(t, identityToken, "identity", identities[0])
			_, err = runtime.Identities.Configs.List(ctx, identities[0], forward.IdentityConfigListParams{Limit: forward.Int(1)})
			saCheck(t, err)
			_, err = runtime.Identities.Configs.List(ctx, identities[1], forward.IdentityConfigListParams{Limit: forward.Int(1)})
			saTokenDenied(t, err, 403)
			t.Log("verified admin business access, Identity business access and cross-Identity refusal")
		})
	}
}
