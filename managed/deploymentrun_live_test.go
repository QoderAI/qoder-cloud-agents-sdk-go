//go:build live

package managed_test

import (
	"context"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"testing"
)

func TestDeploymentRunListLive(t *testing.T) {
	s := newManagedScenarioSuite(t)
	ctx, cancel := s.context()
	defer cancel()
	liveResult(s.client.DeploymentRuns.List(ctx, managed.DeploymentRunListParams{})).require(t)
}

func TestManagedDeploymentScopedRunsLive(t *testing.T) {
	s := newManagedScenarioSuite(t)
	s.requireExecution(t)
	ctx, cancel := s.context()
	defer cancel()
	environment, agent := s.createEnvironment(t), s.createAgent(t)
	params := liveJSON[managed.DeploymentNewParams](t, map[string]any{"name": managedUnique("scoped-runs"), "agent": agent.ID, "environment_id": environment.ID, "initial_events": []any{map[string]any{"type": "user.message", "content": []any{map[string]any{"type": "text", "text": "Reply with SDK-LIVE."}}}}})
	deployment := liveResult(s.client.Deployments.New(ctx, params)).require(t)
	s.cleanup(t, "deployment", func(ctx context.Context) error {
		_, err := s.client.Deployments.Archive(ctx, deployment.ID, managed.DeploymentArchiveParams{})
		return err
	})
	run := liveResult(s.client.Deployments.Run(ctx, deployment.ID, managed.DeploymentRunParams{})).require(t)
	if run.SessionID == "" {
		t.Fatal("deployment run returned no session for cleanup")
	}
	s.cleanupSession(t, run.SessionID)
	scoped := liveResult(s.client.Deployments.Runs.Get(ctx, deployment.ID, run.ID, managed.DeploymentScopedRunGetParams{})).require(t)
	if scoped.ID != run.ID || scoped.DeploymentID != deployment.ID || scoped.SessionID != run.SessionID {
		t.Fatal("scoped run changed identity, deployment, or session")
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
		t.Fatal("scoped list omitted the created Run")
	}
}
