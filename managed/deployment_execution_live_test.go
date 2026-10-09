//go:build live

package managed_test

import (
	"context"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/internal/testsupport"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"testing"
)

func TestManagedDeploymentE2ELive(t *testing.T) {
	testsupport.RequireE2E(t, "MANAGED")
	s := newManagedScenarioSuite(t)
	s.requireExecution(t)
	ctx, cancel := context.WithTimeout(context.Background(), testsupport.ExecutionTimeout(t))
	defer cancel()
	env := s.createEnvironment(t)
	agent := s.createAgent(t)
	marker := testsupport.Marker(t)
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
	scoped := liveResult(s.client.Deployments.Runs.Get(ctx, deployment.ID, run.ID, managed.DeploymentScopedRunGetParams{})).require(t)
	if scoped.ID != run.ID || scoped.DeploymentID != deployment.ID || scoped.SessionID != run.SessionID {
		t.Fatal("scoped run changed its identity, deployment, or session")
	}
	runs := liveResult(s.client.Deployments.Runs.List(ctx, deployment.ID, managed.DeploymentScopedRunListParams{Limit: managed.Int(10)})).require(t)
	found := false
	for _, item := range runs.Data {
		if item.DeploymentID != deployment.ID {
			t.Fatal("scoped list returned another deployment")
		}
		found = found || item.ID == run.ID
	}
	if !found {
		t.Fatal("scoped list did not include the created Run")
	}
}
