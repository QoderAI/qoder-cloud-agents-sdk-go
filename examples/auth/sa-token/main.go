// Set QODER_ACCESS_TOKEN and QODER_SA_REGION=cn|global, then run:
// go run ./examples/auth/sa-token
// The caller is responsible for replacing expired SATs.
package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/live"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "SAT example failed: "+live.Safe(err))
		os.Exit(1)
	}
}

func run() error {
	if os.Getenv("QODER_SA_KEY") != "" {
		return fmt.Errorf("unset QODER_SA_KEY before running the SAT example")
	}
	token := os.Getenv("QODER_ACCESS_TOKEN")
	if strings.TrimSpace(token) == "" {
		return fmt.Errorf("configure QODER_ACCESS_TOKEN")
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
	forwardClient := forward.NewClient(option.WithAuthMode("sa"), option.WithAuthToken(token), option.WithBaseURL(host+"/api/v1/forward"), option.WithMaxRetries(0))
	managedClient := managed.NewClient(option.WithAuthMode("sa"), option.WithAuthToken(token), option.WithBaseURL(host+"/api/v1/cloud"), option.WithMaxRetries(0))
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
