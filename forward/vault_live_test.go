//go:build live

package forward_test

import (
	"context"
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
	credential, err := s.client.Vaults.Credentials.New(ctx, vault.ID, forward.VaultCredentialNewParams{Auth: map[string]any{"type": "static_bearer", "mcp_server_url": "https://example.com/" + liveName("mcp"), "token": secret}})
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
}
