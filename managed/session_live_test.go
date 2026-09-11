//go:build live

package managed_test

import (
	"context"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"testing"
)

func TestSessionListLive(t *testing.T) {
	s := newManagedScenarioSuite(t)
	ctx, cancel := s.context()
	defer cancel()
	liveResult(s.client.Sessions.List(ctx, managed.SessionListParams{})).require(t)
}

func TestSessionLifecycleLive(t *testing.T) {
	s := newManagedScenarioSuite(t)
	ctx, cancel := s.context()
	defer cancel()

	environment := s.createEnvironment(t)
	agent := s.createAgent(t)
	params := liveJSON[managed.SessionNewParams](t, map[string]any{"agent": agent.ID, "environment_id": environment.ID, "title": "Managed SDK live session"})
	session := liveResult(s.client.Sessions.New(ctx, params)).require(t)
	s.cleanup(t, "Session "+session.ID, func(ctx context.Context) error {
		_, err := s.client.Sessions.Delete(ctx, session.ID, managed.SessionDeleteParams{})
		return err
	})
	liveResult(s.client.Sessions.Get(ctx, session.ID, managed.SessionGetParams{})).require(t)
	liveResult(s.client.Sessions.Update(ctx, session.ID, managed.SessionUpdateParams{Title: managed.String("Managed SDK renamed session")})).require(t)
	liveResult(s.client.Sessions.Events.List(ctx, session.ID, managed.SessionEventListParams{})).require(t)
	liveResult(s.client.Sessions.Resources.List(ctx, session.ID, managed.SessionResourceListParams{})).require(t)
	liveResult(s.client.Sessions.Threads.List(ctx, session.ID, managed.SessionThreadListParams{})).require(t)

}
