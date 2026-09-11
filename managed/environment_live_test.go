//go:build live

package managed_test

import (
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"testing"
)

func TestEnvironmentListLive(t *testing.T) {
	s := newManagedScenarioSuite(t)
	ctx, cancel := s.context()
	defer cancel()
	liveResult(s.client.Environments.List(ctx, managed.EnvironmentListParams{})).require(t)
}

func TestEnvironmentLifecycleLive(t *testing.T) {
	s := newManagedScenarioSuite(t)
	ctx, cancel := s.context()
	defer cancel()

	environment := s.createEnvironment(t)
	got := liveResult(s.client.Environments.Get(ctx, environment.ID, managed.EnvironmentGetParams{})).require(t)
	if got.ID != environment.ID {
		t.Fatal(got.ID)
	}
	liveResult(s.client.Environments.Update(ctx, environment.ID, managed.EnvironmentUpdateParams{Description: managed.String("updated through Managed Go SDK")})).require(t)
	liveResult(s.client.Environments.Work.List(ctx, environment.ID, managed.EnvironmentWorkListParams{})).require(t)

}
