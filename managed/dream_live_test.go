//go:build live

package managed_test

import (
	"context"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"testing"
)

func TestDreamListLive(t *testing.T) {
	s := newManagedScenarioSuite(t)
	ctx, cancel := s.context()
	defer cancel()
	liveResult(s.client.Dreams.List(ctx, managed.DreamListParams{})).require(t)
}

func TestDreamLifecycleLive(t *testing.T) {
	s := newManagedScenarioSuite(t)
	ctx, cancel := s.context()
	defer cancel()

	s.requireExecution(t)
	store := s.createMemoryStore(t)
	params := liveJSON[managed.DreamNewParams](t, map[string]any{"inputs": []any{map[string]any{"type": "memory_store", "memory_store_id": store.ID}}, "model": s.model(t), "instructions": "Summarize the test memory store."})
	created := liveResult(s.client.Dreams.New(ctx, params)).require(t)
	s.cleanup(t, "Dream "+created.ID, func(ctx context.Context) error {
		_, err := s.client.Dreams.Archive(ctx, created.ID, managed.DreamArchiveParams{})
		return err
	})
	liveResult(s.client.Dreams.Get(ctx, created.ID, managed.DreamGetParams{})).require(t)
	if status := created.Status; status != "completed" && status != "failed" && status != "canceled" {
		liveResult(s.client.Dreams.Cancel(ctx, created.ID, managed.DreamCancelParams{})).require(t)
	}

}
