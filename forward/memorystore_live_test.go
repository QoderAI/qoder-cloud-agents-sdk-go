//go:build live

package forward_test

import (
	"context"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestMemoryStoreAndMemoryLifecycleLive(t *testing.T) {
	s := newLiveSuite(t, "WRITE")
	ctx := s.context(t)
	store, err := s.client.MemoryStores.New(ctx, forward.MemoryStoreNewParams{Name: liveName("memory"), IdempotencyKey: liveName("key")})
	liveCheck(t, err)
	s.cleanup(t, "memory store", func(ctx context.Context) error { _, err := s.client.MemoryStores.Delete(ctx, store.ID); return err })
	memory, err := s.client.MemoryStores.Memories.New(ctx, store.ID, forward.MemoryStoreMemoryNewParams{Path: "sdk/live.md", Content: "SDK initial content"})
	liveCheck(t, err)
	s.cleanup(t, "memory", func(ctx context.Context) error {
		_, err := s.client.MemoryStores.Memories.Delete(ctx, store.ID, memory.ID)
		return err
	})
	updated, err := s.client.MemoryStores.Memories.Update(ctx, store.ID, memory.ID, forward.MemoryStoreMemoryUpdateParams{Content: "SDK updated content", ContentSHA256: forward.String(memory.ContentSHA256)})
	liveCheck(t, err)
	if updated.Content != "SDK updated content" || updated.ContentSHA256 == memory.ContentSHA256 {
		t.Fatal("memory update did not round trip")
	}
	_, err = s.client.MemoryStores.MemoryVersions.List(ctx, store.ID, forward.MemoryStoreMemoryVersionListParams{})
	liveCheck(t, err)
	identity := s.identity(t)
	template := s.template(t, s.environment(t).ID)
	_, err = s.client.Identities.MemoryStores.Mount(ctx, identity.ID, template.ID, forward.IdentityMemoryStoreMountParams{MemoryStoreID: store.ID})
	liveCheck(t, err)
	s.cleanup(t, "memory mount", func(ctx context.Context) error {
		_, err := s.client.Identities.MemoryStores.Detach(ctx, identity.ID, template.ID, store.ID)
		return err
	})
	_, err = s.client.Identities.MemoryStores.List(ctx, identity.ID, template.ID)
	liveCheck(t, err)
}
