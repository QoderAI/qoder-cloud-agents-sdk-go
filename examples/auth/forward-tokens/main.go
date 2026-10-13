// Issue a Forward admin/Identity SAT, make one read-only business call, then revoke it.
// Set QODER_SA_KEY and QODER_SA_REGION=cn|global. Set QODER_IDENTITY_ID for Identity mode.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/live"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func main() {
	flag.Parse()
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Forward token example failed: "+live.Safe(err))
		os.Exit(1)
	}
}

func run() (err error) {
	key := os.Getenv("QODER_SA_KEY")
	if key == "" {
		return fmt.Errorf("configure QODER_SA_KEY")
	}
	region := os.Getenv("QODER_SA_REGION")
	if region == "" {
		region = "cn"
	}
	if region != "cn" && region != "global" {
		return fmt.Errorf("QODER_SA_REGION must be cn or global")
	}
	host := "https://api.qoder.com"
	if region == "cn" {
		host = "https://api.qoder.com.cn"
	}
	baseURL := host + "/api/v1/forward"
	issuer := forward.NewClient(option.WithBaseURL(baseURL), option.WithAuthMode("sa"))
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	identityID := os.Getenv("QODER_IDENTITY_ID")
	params := forward.ServiceAccountTokenNewParams{SAKey: key, TTLSeconds: forward.Int(600), Metadata: map[string]any{"purpose": "sdk example"}}
	var token *forward.ServiceAccountToken
	if identityID == "" {
		token, err = issuer.ServiceAccountTokens.New(ctx, params)
	} else {
		token, err = issuer.Identities.ServiceAccountTokens.New(ctx, identityID, params)
	}
	if err != nil {
		return err
	}
	defer func() {
		cleanup, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		_, revokeErr := issuer.ServiceAccountTokens.Revoke(cleanup, token.AuthTokenID, forward.ServiceAccountTokenRevokeParams{
			SAKey: key, TargetCredentialID: forward.String(token.CredentialID), TargetExpiresAt: param.NewOpt(token.ExpiresAt), Reason: forward.String("sdk example complete"),
		})
		err = errors.Join(err, revokeErr)
	}()
	runtime := forward.NewClient(option.WithBaseURL(baseURL), option.WithAuthMode("sa"), option.WithAuthToken(token.AccessToken), option.WithMaxRetries(0))
	if identityID == "" {
		page, callErr := runtime.Templates.List(ctx, forward.TemplateListParams{Limit: forward.Int(1)})
		if callErr != nil {
			return callErr
		}
		fmt.Println("Forward templates:", len(page.Data))
	} else {
		page, callErr := runtime.Sessions.List(ctx, forward.SessionListParams{IdentityIDs: []string{identityID}, Limit: forward.Int(1)})
		if callErr != nil {
			return callErr
		}
		fmt.Println("Identity sessions:", len(page.Data))
	}
	return nil
}
