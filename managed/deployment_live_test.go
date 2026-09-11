//go:build live

package managed_test

import (
	"context"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"testing"
)

func TestDeploymentListLive(t *testing.T) {
	s := newManagedScenarioSuite(t)
	ctx, cancel := s.context()
	defer cancel()
	liveResult(s.client.Deployments.List(ctx, managed.DeploymentListParams{})).require(t)
}

func TestDeploymentLifecycleLive(t *testing.T) {
	s := newManagedScenarioSuite(t)
	ctx, cancel := s.context()
	defer cancel()

	environment := s.createEnvironment(t)
	agent := s.createAgent(t)
	params := liveJSON[managed.DeploymentNewParams](t, map[string]any{"name": managedUnique("deployment"), "agent": agent.ID, "environment_id": environment.ID, "initial_events": []any{map[string]any{"type": "user.message", "content": []any{map[string]any{"type": "text", "text": "Reply with SDK-LIVE."}}}}})
	created := liveResult(s.client.Deployments.New(ctx, params)).require(t)
	s.cleanup(t, "Deployment "+created.ID, func(ctx context.Context) error {
		_, err := s.client.Deployments.Archive(ctx, created.ID, managed.DeploymentArchiveParams{})
		return err
	})
	liveResult(s.client.Deployments.Get(ctx, created.ID, managed.DeploymentGetParams{})).require(t)
	liveResult(s.client.Deployments.Update(ctx, created.ID, managed.DeploymentUpdateParams{Description: managed.String("updated through Managed Go SDK")})).require(t)
	liveResult(s.client.Deployments.Pause(ctx, created.ID, managed.DeploymentPauseParams{})).require(t)
	liveResult(s.client.Deployments.Unpause(ctx, created.ID, managed.DeploymentUnpauseParams{})).require(t)
	liveResult(s.client.DeploymentRuns.List(ctx, managed.DeploymentRunListParams{})).require(t)

}
