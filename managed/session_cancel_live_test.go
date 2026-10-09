//go:build live

package managed_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
)

func TestManagedSessionCancelActiveLive(t *testing.T) {
	s := newManagedScenarioSuite(t)
	s.requireExecution(t)
	ctx, cancel := s.context()
	defer cancel()
	environment := s.createEnvironment(t)
	agentParams := liveJSON[managed.AgentNewParams](t, map[string]any{"name": managedUnique("cancel-agent"), "model": map[string]any{"id": s.model(t)}, "tools": []any{map[string]any{"type": "agent_toolset_20260401"}}})
	agent := liveResult(s.client.Agents.New(ctx, agentParams)).require(t)
	s.cleanup(t, "agent", func(ctx context.Context) error {
		_, err := s.client.Agents.Archive(ctx, agent.ID, managed.AgentArchiveParams{})
		return err
	})
	params := liveJSON[managed.SessionNewParams](t, map[string]any{"agent": agent.ID, "environment_id": environment.ID})
	session := liveResult(s.client.Sessions.New(ctx, params)).require(t)
	s.cleanupSession(t, session.ID)
	events := liveJSON[managed.SessionEventSendParams](t, map[string]any{"events": []any{map[string]any{"type": "user.message", "content": []any{map[string]any{"type": "text", "text": "Use a shell command to sleep 30 seconds, then reply with SDK-LIVE."}}}}})
	liveResult(s.client.Sessions.Events.Send(ctx, session.ID, events)).require(t)
	var raw *http.Response
	ack := liveResult(s.client.Sessions.Cancel(ctx, session.ID, managed.SessionCancelParams{}, option.WithResponseInto(&raw))).require(t)
	// The turn can finish between send and cancel; 200 is the documented idle no-op.
	if (raw.StatusCode != 200 && raw.StatusCode != 202) || ack.ID != session.ID || ack.Type != "session" || ack.Status != "canceling" {
		t.Fatal("cancellation did not return the lightweight acknowledgement")
	}
	t.Logf("session=%s cancellation_http_status=%d", session.ID, raw.StatusCode)
	for {
		current := liveResult(s.client.Sessions.Get(ctx, session.ID, managed.SessionGetParams{})).require(t)
		if current.Status == "idle" || current.Status == "terminated" {
			break
		}
		select {
		case <-ctx.Done():
			t.Fatal("session did not become idle after cancellation")
		case <-time.After(time.Second):
		}
	}
}
