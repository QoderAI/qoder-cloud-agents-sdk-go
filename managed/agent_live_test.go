//go:build live

package managed_test

import (
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"testing"
)

func TestAgentListLive(t *testing.T) {
	s := newManagedScenarioSuite(t)
	ctx, cancel := s.context()
	defer cancel()
	liveResult(s.client.Agents.List(ctx, managed.AgentListParams{})).require(t)
}

func TestAgentLifecycleLive(t *testing.T) {
	s := newManagedScenarioSuite(t)
	ctx, cancel := s.context()
	defer cancel()

	agent := s.createAgent(t)
	liveResult(s.client.Agents.Get(ctx, agent.ID, managed.AgentGetParams{})).require(t)
	liveResult(s.client.Agents.Update(ctx, agent.ID, managed.AgentUpdateParams{Version: managed.Int(agent.Version), Description: managed.String("updated through Managed Go SDK")})).require(t)
	versions := liveResult(s.client.Agents.Versions.List(ctx, agent.ID, managed.AgentVersionListParams{})).require(t)
	if len(versions.Data) < 2 {
		t.Fatal("expected at least two Agent versions")
	}

}
