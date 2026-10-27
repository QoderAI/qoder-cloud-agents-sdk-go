//go:build live

package managed_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/pagination"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/internal/testsupport"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
)

func checkSearchPageLive[T any](t *testing.T, search func(map[string]string) (*pagination.PageCursor[T], error)) {
	t.Helper()
	page, err := search(nil)
	if err != nil {
		t.Fatal(testsupport.SafeError(err))
	}
	if len(page.Data) > 1 || !page.JSON.HasMore.Valid() || !page.JSON.Data.Valid() {
		t.Fatal("invalid Search page envelope")
	}
	t.Logf("first_page_items=%d has_more=%t", len(page.Data), page.HasMore)
	if page.HasMore {
		next, err := page.GetNextPage()
		if err != nil {
			t.Fatal(testsupport.SafeError(err))
		}
		if next == nil || len(next.Data) > 1 {
			t.Fatal("invalid second Search page")
		}
		t.Logf("second_page_items=%d", len(next.Data))
	}
	empty, err := search(map[string]string{"sdk_search_probe": fmt.Sprint(time.Now().UnixNano())})
	if err != nil {
		t.Fatal(testsupport.SafeError(err))
	}
	if len(empty.Data) != 0 || empty.HasMore || empty.NextPage != "" {
		t.Fatal("metadata filter was not applied")
	}
}

func TestResourceSearchLive(t *testing.T) {
	s := newManagedScenarioSuite(t)
	ctx, cancel := s.context()
	defer cancel()
	t.Run("agents", func(t *testing.T) {
		checkSearchPageLive(t, func(metadata map[string]string) (*pagination.PageCursor[managed.ManagedAgentsAgent], error) {
			return s.client.Agents.Search(ctx, managed.AgentSearchParams{Limit: managed.Int(1), Metadata: metadata})
		})
	})
	t.Run("sessions", func(t *testing.T) {
		checkSearchPageLive(t, func(metadata map[string]string) (*pagination.PageCursor[managed.ManagedAgentsSession], error) {
			return s.client.Sessions.Search(ctx, managed.SessionSearchParams{Limit: managed.Int(1), Metadata: metadata})
		})
	})
	t.Run("deployments", func(t *testing.T) {
		checkSearchPageLive(t, func(metadata map[string]string) (*pagination.PageCursor[managed.ManagedAgentsDeployment], error) {
			return s.client.Deployments.Search(ctx, managed.DeploymentSearchParams{Limit: managed.Int(1), Metadata: metadata})
		})
	})
	t.Run("environments", func(t *testing.T) {
		checkSearchPageLive(t, func(metadata map[string]string) (*pagination.PageCursor[managed.Environment], error) {
			return s.client.Environments.Search(ctx, managed.EnvironmentSearchParams{Limit: managed.Int(1), Metadata: metadata})
		})
	})
	t.Run("skills", func(t *testing.T) {
		checkSearchPageLive(t, func(metadata map[string]string) (*pagination.PageCursor[managed.Skill], error) {
			return s.client.Skills.Search(ctx, managed.SkillSearchParams{Limit: managed.Int(1), Metadata: metadata})
		})
	})
	t.Run("vaults", func(t *testing.T) {
		checkSearchPageLive(t, func(metadata map[string]string) (*pagination.PageCursor[managed.ManagedAgentsVault], error) {
			return s.client.Vaults.Search(ctx, managed.VaultSearchParams{Limit: managed.Int(1), Metadata: metadata})
		})
	})
	t.Run("files", func(t *testing.T) {
		checkSearchPageLive(t, func(metadata map[string]string) (*pagination.PageCursor[managed.FileMetadata], error) {
			return s.client.Files.Search(ctx, managed.FileSearchParams{Limit: managed.Int(1), Metadata: metadata})
		})
	})
	t.Run("memory_stores", func(t *testing.T) {
		checkSearchPageLive(t, func(metadata map[string]string) (*pagination.PageCursor[managed.ManagedAgentsMemoryStore], error) {
			return s.client.MemoryStores.Search(ctx, managed.MemoryStoreSearchParams{Limit: managed.Int(1), Metadata: metadata})
		})
	})
}
