//go:build live

package managed_test

import (
	"context"
	"fmt"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/internal/testsupport"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"strings"
	"testing"
)

func TestManagedDreamE2ELive(t *testing.T) {
	testsupport.RequireE2E(t, "MANAGED")
	s := newManagedScenarioSuite(t)
	s.requireExecution(t)
	ctx, cancel := context.WithTimeout(context.Background(), testsupport.ExecutionTimeout(t))
	defer cancel()
	store := s.createMemoryStore(t)
	marker := testsupport.Marker(t)
	liveResult(s.client.MemoryStores.Memories.New(ctx, store.ID, managed.MemoryStoreMemoryNewParams{Path: "sdk-e2e/source.md", Content: managed.String("Permanent project verification code: " + marker + ". Preserve this exact code during consolidation.")})).require(t)
	params := liveJSON[managed.DreamNewParams](t, map[string]any{"inputs": []any{map[string]any{"type": "memory_store", "memory_store_id": store.ID}}, "model": s.model(t), "instructions": "Consolidate the supplied memory into sdk-e2e/consolidated.md. Preserve the exact project verification code. Keep the original source."})
	dream := liveResult(s.client.Dreams.New(ctx, params)).require(t)
	t.Logf("dream=%s input_memory_store=%s", dream.ID, store.ID)
	s.cleanup(t, "dream "+dream.ID, func(ctx context.Context) error {
		current, err := s.client.Dreams.Get(ctx, dream.ID, managed.DreamGetParams{})
		if err != nil {
			return err
		}
		if current.Status == "pending" || current.Status == "running" {
			if _, err = s.client.Dreams.Cancel(ctx, dream.ID, managed.DreamCancelParams{}); err != nil {
				return err
			}
			for {
				current, err = s.client.Dreams.Get(ctx, dream.ID, managed.DreamGetParams{})
				if err != nil {
					return err
				}
				if current.Status != "pending" && current.Status != "running" {
					break
				}
				if err = testsupport.PollPause(ctx); err != nil {
					return fmt.Errorf("dream=%s remains %s: %w", dream.ID, current.Status, err)
				}
			}
		}
		// Dream-created resources belong to this test even if execution failed.
		for _, out := range current.Outputs {
			if out.MemoryStoreID != "" && out.MemoryStoreID != store.ID {
				outputID := out.MemoryStoreID
				s.cleanup(t, "dream output "+outputID, func(ctx context.Context) error {
					_, err := s.client.MemoryStores.Delete(ctx, outputID, managed.MemoryStoreDeleteParams{})
					return err
				})
			}
		}
		if current.SessionID != "" {
			s.cleanupSession(t, current.SessionID)
		}
		_, err = s.client.Dreams.Archive(ctx, dream.ID, managed.DreamArchiveParams{})
		return err
	})
	for dream.Status == "pending" || dream.Status == "running" {
		if err := testsupport.PollPause(ctx); err != nil {
			t.Fatalf("dream=%s status=%s: %v", dream.ID, dream.Status, err)
		}
		dream = liveResult(s.client.Dreams.Get(ctx, dream.ID, managed.DreamGetParams{})).require(t)
	}
	if dream.Status != "completed" || len(dream.Outputs) == 0 {
		t.Fatalf("dream=%s status=%s outputs=%d", dream.ID, dream.Status, len(dream.Outputs))
	}
	found := false
	for _, out := range dream.Outputs {
		page := s.client.MemoryStores.Memories.ListAutoPaging(ctx, out.MemoryStoreID, managed.MemoryStoreMemoryListParams{})
		for page.Next() {
			memory := page.Current()
			if memory.Path != "sdk-e2e/consolidated.md" {
				continue
			}
			got := liveResult(s.client.MemoryStores.Memories.Get(ctx, memory.ID, managed.MemoryStoreMemoryGetParams{MemoryStoreID: out.MemoryStoreID})).require(t)
			if strings.Contains(got.Content, marker) {
				found = true
			}
		}
		if err := page.Err(); err != nil {
			t.Fatal(err)
		}
	}
	if !found {
		t.Fatal("Dream did not persist the requested consolidated memory with the original verification code")
	}
}
