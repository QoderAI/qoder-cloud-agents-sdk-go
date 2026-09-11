//go:build live

package forward_test

import (
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestIdentityAndConfigLifecycleLive(t *testing.T) {
	s := newLiveSuite(t, "WRITE")
	ctx := s.context(t)
	identity := s.identity(t)
	template := s.template(t, s.environment(t).ID)
	updated, err := s.client.Identities.Update(ctx, identity.ID, forward.IdentityUpdateParams{Name: forward.String("SDK renamed identity")})
	liveCheck(t, err)
	if updated.Name != "SDK renamed identity" {
		t.Fatal("identity update not returned")
	}
	disabled, err := s.client.Identities.Disable(ctx, identity.ID)
	liveCheck(t, err)
	if disabled.Enabled {
		t.Fatal("identity still enabled")
	}
	enabled, err := s.client.Identities.Enable(ctx, identity.ID)
	liveCheck(t, err)
	if !enabled.Enabled {
		t.Fatal("identity still disabled")
	}
	_, err = s.client.Identities.Configs.Upsert(ctx, identity.ID, template.ID, forward.IdentityConfigUpsertParams{
		IdentityConfig: forward.IdentityConfigSpecParam{EnvironmentVariables: map[string]forward.EnvironmentVariableOverrideParam{"SDK_LIVE": {Op: "set", Value: forward.String("one")}}},
	})
	liveCheck(t, err)
	config, err := s.client.Identities.Configs.Get(ctx, identity.ID, template.ID)
	liveCheck(t, err)
	if config.IdentityConfig.EnvironmentVariables["SDK_LIVE"].Value != "one" {
		t.Fatal("identity config did not round trip")
	}
	_, err = s.client.Identities.Configs.GetEffective(ctx, identity.ID, template.ID)
	liveCheck(t, err)
	_, err = s.client.Identities.Configs.List(ctx, identity.ID, forward.IdentityConfigListParams{})
	liveCheck(t, err)
	_, err = s.client.Identities.ListTemplates(ctx, identity.ID)
	liveCheck(t, err)
}
