//go:build live

package forward_test

import (
	"context"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestTemplateLifecycleLive(t *testing.T) {
	s := newLiveSuite(t, "WRITE")
	ctx := s.context(t)
	template := s.template(t, s.environment(t).ID)
	got, err := s.client.Templates.Get(ctx, template.ID)
	liveCheck(t, err)
	if got.ID != template.ID || got.Model.ID == "" {
		t.Fatal("template/model did not round trip")
	}
	updated, err := s.client.Templates.Update(ctx, template.ID, forward.TemplateUpdateParams{Name: forward.String(template.Name + "-updated")})
	liveCheck(t, err)
	if updated.Name != template.Name+"-updated" {
		t.Fatal("name not updated")
	}
	clone, err := s.client.Templates.Clone(ctx, template.ID, forward.TemplateCloneParams{})
	liveCheck(t, err)
	s.cleanup(t, "cloned template", func(ctx context.Context) error {
		_, err := s.client.Templates.Archive(ctx, clone.ID, forward.TemplateArchiveParams{})
		return err
	})
	if clone.ID == template.ID {
		t.Fatal("clone reused original ID")
	}
}
