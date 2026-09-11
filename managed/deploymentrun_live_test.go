//go:build live

package managed_test

import (
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"testing"
)

func TestDeploymentRunListLive(t *testing.T) {
	s := newManagedScenarioSuite(t)
	ctx, cancel := s.context()
	defer cancel()
	liveResult(s.client.DeploymentRuns.List(ctx, managed.DeploymentRunListParams{})).require(t)
}
