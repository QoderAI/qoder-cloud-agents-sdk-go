// List endpoints with shared PAT/administrator-SA configuration.
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/live"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
	"net/http"
	"os"
)

// parseDelivery is the receiver integration point; use the original HTTP body bytes.
func parseDelivery(raw []byte, headers http.Header, secret string) (map[string]any, error) {
	if err := forward.VerifyWebhookSignature(raw, headers, secret); err != nil {
		return nil, err
	}
	var event map[string]any
	err := json.Unmarshal(raw, &event)
	return event, err
}
func main() {
	config, err := live.Load("forward")
	if err != nil {
		fmt.Fprintln(os.Stderr, live.Safe(err))
		os.Exit(1)
	}
	client := forward.NewClient(config.Options()...)
	os.Exit(live.Execute(config, live.Scenario{Name: "webhooks", Description: "List webhook endpoints", Run: func(ctx context.Context, r *live.Run) error {
		pages := client.Webhooks.Endpoints.ListAutoPaging(ctx, forward.WebhookEndpointListParams{Limit: forward.Int(20)})
		for pages.Next() {
			fmt.Printf("endpoint=%s active=%t\n", pages.Current().ID, pages.Current().Active)
		}
		return pages.Err()
	}}))
}
