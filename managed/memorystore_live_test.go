//go:build live

package managed_test

import (
	"context"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"testing"
)

func TestMemoryStoreListLive(t *testing.T) {
	s := newManagedScenarioSuite(t)
	ctx, cancel := s.context()
	defer cancel()
	liveResult(s.client.MemoryStores.List(ctx, managed.MemoryStoreListParams{})).require(t)
}

func TestMemorystoreLifecycleLive(t *testing.T) {
	s := newManagedScenarioSuite(t)
	ctx, cancel := s.context()
	defer cancel()

	store := s.createMemoryStore(t)
	liveResult(s.client.MemoryStores.Get(ctx, store.ID, managed.MemoryStoreGetParams{})).require(t)
	liveResult(s.client.MemoryStores.Update(ctx, store.ID, managed.MemoryStoreUpdateParams{Description: managed.String("updated through Managed Go SDK")})).require(t)
	memory := liveResult(s.client.MemoryStores.Memories.New(ctx, store.ID, managed.MemoryStoreMemoryNewParams{Path: "notes/sdk-live.md", Content: managed.String("Managed SDK memory v1")})).require(t)
	s.cleanup(t, "Memory "+memory.ID, func(ctx context.Context) error {
		_, err := s.client.MemoryStores.Memories.Delete(ctx, memory.ID, managed.MemoryStoreMemoryDeleteParams{MemoryStoreID: store.ID})
		return err
	})
	liveResult(s.client.MemoryStores.Memories.Get(ctx, memory.ID, managed.MemoryStoreMemoryGetParams{MemoryStoreID: store.ID})).require(t)
	liveResult(s.client.MemoryStores.Memories.Update(ctx, memory.ID, managed.MemoryStoreMemoryUpdateParams{MemoryStoreID: store.ID, Content: managed.String("Managed SDK memory v2")})).require(t)
	versions := liveResult(s.client.MemoryStores.MemoryVersions.List(ctx, store.ID, managed.MemoryStoreMemoryVersionListParams{})).require(t)
	if len(versions.Data) == 0 {
		t.Fatal("expected at least one memory version")
	}
	liveResult(s.client.MemoryStores.MemoryVersions.Get(ctx, versions.Data[0].ID, managed.MemoryStoreMemoryVersionGetParams{MemoryStoreID: store.ID})).require(t)

}
