// Set QODER_SA_KEY and QODER_SA_REGION=cn|global, then run:
// go run ./examples/auth/sa-key
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/live"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "SA Key example failed: "+live.Safe(err))
		os.Exit(1)
	}
}

func run() error {
	if os.Getenv("QODER_ACCESS_TOKEN") != "" {
		return fmt.Errorf("unset QODER_ACCESS_TOKEN before running the SA Key example")
	}
	region := os.Getenv("QODER_SA_REGION")
	if region == "" {
		region = "cn"
	}
	credential, err := convention.NewServiceAccountCredential(convention.ServiceAccountCredentialConfig{
		SAKey: os.Getenv("QODER_SA_KEY"), Region: region, Scopes: []string{"qca.access", "forward.access"},
	})
	if err != nil {
		return err
	}
	host := "https://api.qoder.com"
	if region == "cn" {
		host = "https://api.qoder.com.cn"
	}
	// The clients share one automatically renewed SAT cache.
	forwardClient := forward.NewClient(option.WithCredential(credential), option.WithBaseURL(host+"/api/v1/forward"), option.WithMaxRetries(0))
	managedClient := managed.NewClient(option.WithCredential(credential), option.WithBaseURL(host+"/api/v1/cloud"), option.WithMaxRetries(0))
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	templates, err := forwardClient.Templates.List(ctx, forward.TemplateListParams{})
	if err != nil {
		return err
	}
	fmt.Println("Forward templates:", len(templates.Data))
	agents, err := managedClient.Agents.List(ctx, managed.AgentListParams{})
	if err != nil {
		return err
	}
	fmt.Println("Managed agents:", len(agents.Data))
	return nil
}
