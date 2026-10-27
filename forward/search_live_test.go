//go:build live

package forward_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/pagination"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/internal/testsupport"
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
	s := newLiveSuite(t)
	ctx := s.context(t)
	t.Run("environments", func(t *testing.T) {
		checkSearchPageLive(t, func(metadata map[string]string) (*pagination.PageCursor[forward.Environment], error) {
			return s.client.Environments.Search(ctx, forward.EnvironmentSearchParams{Limit: forward.Int(1), Metadata: metadata})
		})
	})
	t.Run("skills", func(t *testing.T) {
		checkSearchPageLive(t, func(metadata map[string]string) (*pagination.PageCursor[forward.Skill], error) {
			return s.client.Skills.Search(ctx, forward.SkillSearchParams{Limit: forward.Int(1), Metadata: metadata})
		})
	})
	t.Run("vaults", func(t *testing.T) {
		checkSearchPageLive(t, func(metadata map[string]string) (*pagination.PageCursor[forward.Vault], error) {
			return s.client.Vaults.Search(ctx, forward.VaultSearchParams{Limit: forward.Int(1), Metadata: metadata})
		})
	})
	t.Run("files", func(t *testing.T) {
		checkSearchPageLive(t, func(metadata map[string]string) (*pagination.PageCursor[forward.FileMetadata], error) {
			return s.client.Files.Search(ctx, forward.FileSearchParams{Limit: forward.Int(1), Metadata: metadata})
		})
	})
}
