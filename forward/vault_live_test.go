//go:build live

package forward_test

import (
	"context"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"strings"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestVaultAndCredentialLifecycleLive(t *testing.T) {
	s := newLiveSuite(t, "WRITE")
	ctx := s.context(t)
	vault, err := s.client.Vaults.New(ctx, forward.VaultNewParams{DisplayName: liveName("vault")})
	liveCheck(t, err)
	s.cleanup(t, "vault", func(ctx context.Context) error { return s.client.Vaults.Delete(ctx, vault.ID) })
	got, err := s.client.Vaults.Get(ctx, vault.ID)
	liveCheck(t, err)
	if got.ID != vault.ID {
		t.Fatal("vault did not round trip")
	}
	secret := "sdk-live-placeholder-secret"
	mcpURL := "https://example.com/" + liveName("mcp")
	credential, err := s.client.Vaults.Credentials.New(ctx, vault.ID, forward.VaultCredentialNewParams{Auth: map[string]any{"type": "static_bearer", "mcp_server_url": mcpURL, "token": secret}})
	liveCheck(t, err)
	s.cleanup(t, "credential", func(ctx context.Context) error {
		return s.client.Vaults.Credentials.Delete(ctx, vault.ID, credential.ID)
	})
	gotCredential, err := s.client.Vaults.Credentials.Get(ctx, vault.ID, credential.ID)
	liveCheck(t, err)
	if gotCredential.ID != credential.ID || strings.Contains(gotCredential.RawJSON(), secret) || strings.Contains(credential.RawJSON(), secret) {
		t.Fatal("credential response failed ID or secret-redaction check")
	}
	_, err = s.client.Vaults.Credentials.List(ctx, vault.ID, forward.VaultCredentialListParams{})
	liveCheck(t, err)
	_, err = s.client.Vaults.Credentials.Update(ctx, vault.ID, credential.ID, forward.VaultCredentialUpdateParams{Metadata: map[string]any{"keep": "original", "remove": "old"}})
	liveCheck(t, err)
	rotated := liveName("rotated-secret")
	updated, err := s.client.Vaults.Credentials.Update(ctx, vault.ID, credential.ID, forward.VaultCredentialUpdateParams{
		Auth:     forward.VaultCredentialUpdateAuthUnionParam{OfStaticBearer: &forward.StaticBearerUpdateParam{Type: "static_bearer", Token: forward.String(rotated)}},
		Metadata: map[string]any{"remove": nil, "added": "new"},
	})
	liveCheck(t, err)
	saved, err := s.client.Vaults.Credentials.Get(ctx, vault.ID, credential.ID)
	liveCheck(t, err)
	for _, value := range []*forward.VaultCredential{updated, saved} {
		_, removed := value.Metadata["remove"]
		if value.ID != credential.ID || value.Auth.Type != "static_bearer" || value.Auth.MCPServerURL != mcpURL || value.Metadata["keep"] != "original" || value.Metadata["added"] != "new" || removed || strings.Contains(value.RawJSON(), secret) || strings.Contains(value.RawJSON(), rotated) {
			t.Fatal("credential patch failed preservation or secret-redaction check")
		}
	}
	cleared, err := s.client.Vaults.Credentials.Update(ctx, vault.ID, credential.ID, forward.VaultCredentialUpdateParams{Metadata: param.NullMap[map[string]any]()})
	liveCheck(t, err)
	saved, err = s.client.Vaults.Credentials.Get(ctx, vault.ID, credential.ID)
	liveCheck(t, err)
	if len(cleared.Metadata) != 0 || len(saved.Metadata) != 0 || strings.Contains(cleared.RawJSON()+saved.RawJSON(), rotated) || strings.Contains(cleared.RawJSON()+saved.RawJSON(), secret) {
		t.Fatal("null metadata did not clear metadata or disclosed the secret")
	}
}
