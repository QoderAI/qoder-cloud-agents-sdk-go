// Evaluate sample state, or set QODER_DECISION_SESSION_ID to evaluate a Session.
// Run: go run ./examples/forward/decisions -env .env.live
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/internal/live"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func main() {
	config, err := live.Load("forward")
	if err != nil {
		fmt.Fprintln(os.Stderr, live.Safe(err))
		os.Exit(1)
	}
	client := forward.NewClient(config.Options()...)
	os.Exit(live.Execute(config, live.Scenario{Name: "decisions", Description: "Create and retrieve a Decision, then list one page", Run: func(ctx context.Context, r *live.Run) error {
		params := forward.DecisionNewParams{
			State: map[string]any{"message": "The issue is resolved; everything works now."},
			Questions: []forward.DecisionQuestionParam{{
				Type: "choice", Name: forward.String("resolved"),
				Instructions: "Has the user explicitly confirmed that their issue is resolved?",
				Choices:      []forward.DecisionChoiceParam{{Value: true}, {Value: false}},
			}},
			IdempotencyKey: forward.String(fmt.Sprintf("sdk-decision-%d", time.Now().UnixNano())),
		}
		listParams := forward.DecisionListParams{Limit: forward.Int(5)}
		if id := os.Getenv("QODER_DECISION_SESSION_ID"); id != "" {
			params.State = nil
			params.SessionID = forward.String(id)
			listParams.SessionID = forward.String(id)
		}
		decision, err := client.Decisions.New(ctx, params)
		if err != nil {
			return err
		}
		stored, err := client.Decisions.Get(ctx, decision.ID)
		if err != nil {
			return err
		}
		fmt.Println(stored.RawJSON())
		page, err := client.Decisions.List(ctx, listParams)
		if err != nil {
			return err
		}
		for _, item := range page.Data {
			fmt.Printf("decision=%s status=%s\n", item.ID, item.Status)
		}
		// Decision records expire according to server retention; there is no delete API.
		return nil
	}}))
}
