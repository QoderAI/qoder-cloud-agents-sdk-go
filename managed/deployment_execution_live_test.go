//go:build live

package managed_test

import (
	"context"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/examples/testutil"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"testing"
)

func TestManagedDeploymentE2ELive(t *testing.T) {
	testutil.RequireE2E(t, "MANAGED")
	s := newManagedScenarioSuite(t)
	s.requireExecution(t)
	ctx, cancel := context.WithTimeout(context.Background(), testutil.ExecutionTimeout(t))
	defer cancel()
	env := s.createEnvironment(t)
	agent := s.createAgent(t)
	marker := testutil.Marker(t)
	params := liveJSON[managed.DeploymentNewParams](t, map[string]any{"name": managedUnique("deployment-e2e"), "agent": agent.ID, "environment_id": env.ID, "initial_events": []any{map[string]any{"type": "user.message", "content": []any{map[string]any{"type": "text", "text": "Reply with exactly " + marker}}}}})
	deployment := liveResult(s.client.Deployments.New(ctx, params)).require(t)
	s.cleanup(t, "deployment "+deployment.ID, func(ctx context.Context) error {
		_, err := s.client.Deployments.Archive(ctx, deployment.ID, managed.DeploymentArchiveParams{})
		return err
	})
	run := liveResult(s.client.Deployments.Run(ctx, deployment.ID, managed.DeploymentRunParams{})).require(t)
	t.Logf("deployment=%s run=%s session=%s", deployment.ID, run.ID, run.SessionID)
	if run.SessionID == "" {
		t.Fatalf("deployment run=%s returned no session", run.ID)
	}
	s.cleanupSession(t, run.SessionID)
	got := liveResult(s.client.DeploymentRuns.Get(ctx, run.ID, managed.DeploymentRunGetParams{})).require(t)
	if got.SessionID != run.SessionID {
		t.Fatal("run session changed")
	}
	s.waitTurn(t, run.SessionID, "", []string{marker}, false, false)
}
