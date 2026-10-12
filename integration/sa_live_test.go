//go:build live

package integration_test

import (
	"context"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"os"
	"testing"
	"time"
)

func TestServiceAccountReadOnlyLive(t *testing.T) {
	saRequireConfiguration(t, false)
	for _, region := range []string{"cn", "global"} {
		t.Run(region, func(t *testing.T) {
			name, host := "QODER_SA_GLOBAL_KEY", "https://api.qoder.com"
			if region == "cn" {
				name, host = "QODER_SA_CN_KEY", "https://api.qoder.com.cn"
			}
			key := os.Getenv(name)
			if key == "" {
				t.Fatalf("%s is required", name)
			}
			credential, err := convention.NewServiceAccountCredential(convention.ServiceAccountCredentialConfig{SAKey: key, Region: region, Scopes: []string{"qca.access", "forward.access"}})
			if err != nil {
				saCheck(t, err)
			}
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
			defer cancel()
			for _, auth := range []string{"automatic", "static"} {
				t.Run(auth, func(t *testing.T) {
					opts := []option.RequestOption{option.WithMaxRetries(0), option.WithRequestTimeout(30 * time.Second)}
					if auth == "automatic" {
						opts = append(opts, option.WithCredential(credential))
					} else {
						token, err := credential.Token(ctx)
						if err != nil {
							saCheck(t, err)
						}
						opts = append(opts, option.WithAuthMode("sa"), option.WithAuthToken(token))
					}
					f := forward.NewClient(append(opts, option.WithBaseURL(host+"/api/v1/forward/"))...)
					if _, err := f.Templates.List(ctx, forward.TemplateListParams{}); err != nil {
						saCheck(t, err)
					}
					m := managed.NewClient(append(opts, option.WithBaseURL(host+"/api/v1/cloud/"))...)
					if _, err := m.Agents.List(ctx, managed.AgentListParams{}); err != nil {
						saCheck(t, err)
					}
				})
			}
		})
	}
}
