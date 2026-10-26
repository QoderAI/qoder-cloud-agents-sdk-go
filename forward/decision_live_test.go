//go:build live

package forward_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

// Opt in separately because Decision evaluation consumes model usage.
func TestDecisionE2ELive(t *testing.T) {
	s := newLiveSuite(t, "WRITE", "EXECUTION", "DECISIONS")
	ctx := s.context(t)
	questions := []forward.DecisionQuestionParam{
		{Type: "choice", Name: forward.String("resolved"), Instructions: "Return true only if state or the latest user message contains RESOLVED_SDK_TEST.", Choices: []forward.DecisionChoiceParam{{Value: true, Description: forward.String("The marker is present")}, {Value: false, Description: forward.String("The marker is absent")}}},
		{Type: "score", Name: forward.String("risk"), Instructions: "Select the risk level explicitly indicated by RISK_LEVEL=low or RISK_LEVEL=high in state or the latest user message.", Levels: []forward.DecisionLevelParam{{Label: "low", Score: 0}, {Label: "high", Score: 10}}},
	}
	check := func(d *forward.Decision, mode string) {
		t.Helper()
		if d == nil || d.ID == "" || d.Status != "completed" || len(d.Answers) != 2 {
			t.Fatalf("incomplete Decision: %#v", d)
		}
		if d.Answers[0].Value != true || d.Answers[1].Score != 0 || !d.Answers[1].JSON.Score.Valid() {
			t.Fatalf("unexpected answers: %+v", d.Answers)
		}
		if d.Context.Mode != mode || d.Context.EstimatedInputTokens <= 0 || d.Context.MaxInputTokens <= 0 || d.Context.EstimatedInputTokens > d.Context.MaxInputTokens || d.Context.TokenEstimator == "" {
			t.Fatalf("invalid input metadata: %+v", d.Context)
		}
		t.Logf("verified decision=%s mode=%s tokens=%d budget=%d segments=%d", d.ID, mode, d.Context.EstimatedInputTokens, d.Context.MaxInputTokens, d.Context.IncludedSegmentCount)
	}

	stateParams := forward.DecisionNewParams{State: map[string]any{"message": "RESOLVED_SDK_TEST RISK_LEVEL=low"}, Questions: questions, IdempotencyKey: forward.String(liveName("decision-state"))}
	state, err := s.client.Decisions.New(ctx, stateParams)
	liveCheck(t, err)
	check(state, "state")
	if state.JSON.SessionID.Raw() != "null" {
		t.Fatal("state Decision unexpectedly belongs to a Session")
	}
	replay, err := s.client.Decisions.New(ctx, stateParams)
	liveCheck(t, err)
	if replay.ID != state.ID || !replay.CompletedAt.Equal(state.CompletedAt) {
		t.Fatal("idempotency replay created a new Decision")
	}
	got, err := s.client.Decisions.Get(ctx, state.ID)
	liveCheck(t, err)
	check(got, "state")
	if got.ID != state.ID {
		t.Fatal("Get returned another Decision")
	}
	// Decision has no delete endpoint; the server's record retention handles cleanup.
	t.Logf("Decision records are retained by the server: %s", state.ID)

	environmentID := os.Getenv("QODER_FORWARD_DECISION_ENVIRONMENT_ID")
	if environmentID == "" {
		environmentID = s.environment(t).ID
	}
	identity, err := s.client.Identities.New(ctx, forward.IdentityNewParams{ExternalID: liveName("decision-identity"), Name: forward.String(liveName("decision"))}, option.WithJSONSet("auto_memory_store", false))
	liveCheck(t, err)
	s.cleanup(t, "decision identity", func(ctx context.Context) error { _, err := s.client.Identities.Delete(ctx, identity.ID); return err })
	template := s.template(t, environmentID)
	session := s.session(t, identity.ID, template.ID)
	events := make([]forward.SessionEventParam, 12)
	for i := range events {
		text := fmt.Sprintf("SDK_DECISION_SEGMENT_%02d RESOLVED_SDK_TEST RISK_LEVEL=low. This is inert test data. Please reply briefly without using tools.", i+1)
		events[i] = forward.SessionEventParam{Type: "user.message", Content: forward.EventContentUnionParam{OfBlocks: []forward.ContentBlockParam{{Type: "text", Text: forward.String(text)}}}}
	}
	sent, err := s.client.Sessions.Events.Send(ctx, session.ID, forward.SessionEventSendParams{Events: events, IdempotencyKey: forward.String(liveName("decision-events"))})
	liveCheck(t, err)
	if len(sent.Data) != len(events) {
		t.Fatalf("accepted %d of %d events", len(sent.Data), len(events))
	}
	s.waitTurn(t, session.ID, sent.Data[len(sent.Data)-1].ID, nil, false, false)
	ids := map[string]bool{}
	for i := 0; i < 2; i++ {
		d, err := s.client.Decisions.New(ctx, forward.DecisionNewParams{SessionID: forward.String(session.ID), Questions: questions, IdempotencyKey: forward.String(liveName("decision-session"))})
		liveCheck(t, err)
		check(d, "session")
		if d.SessionID != session.ID || d.Context.IncludedSegmentCount != 12 || d.Context.Truncated || d.Context.ScanLimited {
			t.Fatalf("unexpected session context: %+v", d.Context)
		}
		ids[d.ID] = true
	}
	if len(ids) != 2 {
		t.Fatal("independent keys returned the same Decision")
	}
	pager := s.client.Decisions.ListAutoPaging(ctx, forward.DecisionListParams{SessionID: forward.String(session.ID), Limit: forward.Int(1)})
	count := 0
	for pager.Next() {
		d := pager.Current()
		if !ids[d.ID] || d.SessionID != session.ID {
			t.Fatalf("unexpected Decision in Session-filtered page: %s", d.ID)
		}
		delete(ids, d.ID)
		count++
	}
	liveCheck(t, pager.Err())
	if count != 2 || len(ids) != 0 {
		t.Fatalf("pagination count=%d missing=%d", count, len(ids))
	}

	_, err = s.client.Decisions.New(ctx, forward.DecisionNewParams{State: stateParams.State, SessionID: forward.String(session.ID), Questions: questions})
	var api *forward.Error
	if !errors.As(err, &api) || api.StatusCode != 400 {
		t.Fatalf("combined state/session must return 400: %v", err)
	}
	if state.Context.MaxInputTokens > 100000 {
		t.Log("skipping token-overflow probe: it would exceed the independent HTTP body bound")
		return
	}
	_, err = s.client.Decisions.New(ctx, forward.DecisionNewParams{State: map[string]any{"text": strings.Repeat("测", int(state.Context.MaxInputTokens))}, Questions: questions})
	if !errors.As(err, &api) || api.StatusCode != 422 || api.Code != "decision_context_too_large" {
		t.Fatalf("over-budget explicit state must return 422: %v", err)
	}
}
