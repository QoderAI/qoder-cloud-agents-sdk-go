//go:build live

package forward_test

import (
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestEnvironmentLifecycleLive(t *testing.T) {
	s := newLiveSuite(t, "WRITE")
	ctx := s.context(t)
	env := s.environment(t)
	got, err := s.client.Environments.Get(ctx, env.ID)
	liveCheck(t, err)
	if got.ID != env.ID || got.Config.Type != "cloud" {
		t.Fatal("environment did not round trip")
	}
	updated, err := s.client.Environments.Update(ctx, env.ID, forward.EnvironmentUpdateParams{Description: forward.String("SDK updated environment")})
	liveCheck(t, err)
	if updated.Description != "SDK updated environment" {
		t.Fatal("description not updated")
	}
}
