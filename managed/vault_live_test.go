//go:build live

package managed_test

import (
	"context"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"strings"
	"testing"
)

func TestVaultListLive(t *testing.T) {
	s := newManagedScenarioSuite(t)
	ctx, cancel := s.context()
	defer cancel()
	liveResult(s.client.Vaults.List(ctx, managed.VaultListParams{})).require(t)
}

func TestVaultLifecycleLive(t *testing.T) {
	s := newManagedScenarioSuite(t)
	ctx, cancel := s.context()
	defer cancel()

	s.requireWrite(t)
	created := liveResult(s.client.Vaults.New(ctx, managed.VaultNewParams{DisplayName: managedUnique("vault"), Metadata: map[string]string{"suite": "sdk-live"}})).require(t)
	s.cleanup(t, "Vault "+created.ID, func(ctx context.Context) error {
		_, err := s.client.Vaults.Delete(ctx, created.ID, managed.VaultDeleteParams{})
		return err
	})
	liveResult(s.client.Vaults.Get(ctx, created.ID, managed.VaultGetParams{})).require(t)
	liveResult(s.client.Vaults.Credentials.List(ctx, created.ID, managed.VaultCredentialListParams{})).require(t)

}

func TestVaultCredentialUpdateLifecycleLive(t *testing.T) {
	s := newManagedScenarioSuite(t)
	s.requireWrite(t)
	ctx, cancel := s.context()
	defer cancel()
	vault := liveResult(s.client.Vaults.New(ctx, managed.VaultNewParams{DisplayName: managedUnique("credential-update")})).require(t)
	s.cleanup(t, "Vault "+vault.ID, func(ctx context.Context) error {
		_, err := s.client.Vaults.Delete(ctx, vault.ID, managed.VaultDeleteParams{})
		return err
	})
	credential := liveResult(s.client.Vaults.Credentials.New(ctx, vault.ID, managed.VaultCredentialNewParams{
		Auth: managed.VaultCredentialNewParamsAuthUnion{OfStaticBearer: &managed.ManagedAgentsStaticBearerCreateParams{
			Type: "static_bearer", MCPServerURL: "https://example.com/" + managedUnique("mcp"), Token: managedUnique("synthetic-token"),
		}},
		Metadata: map[string]string{"suite": "sdk-live", "remove": "old"},
	})).require(t)
	s.cleanup(t, "Credential "+credential.ID, func(ctx context.Context) error {
		_, err := s.client.Vaults.Credentials.Delete(ctx, credential.ID, managed.VaultCredentialDeleteParams{VaultID: vault.ID})
		return err
	})
	_, err := s.client.Vaults.Credentials.Update(ctx, credential.ID, managed.VaultCredentialUpdateParams{
		VaultID: vault.ID, DisplayName: managed.String("unsupported update"),
	})
	if err == nil || !strings.Contains(err.Error(), "does not support updating credential display_name") {
		t.Fatal("unsupported credential name must fail locally")
	}
	updated := liveResult(s.client.Vaults.Credentials.Update(ctx, credential.ID, managed.VaultCredentialUpdateParams{
		VaultID: vault.ID,
		Auth: managed.VaultCredentialUpdateParamsAuthUnion{OfStaticBearer: &managed.ManagedAgentsStaticBearerUpdateParams{
			Type: "static_bearer", Token: managed.String(managedUnique("rotated-synthetic-token")),
		}},
		Metadata: map[string]any{"suite": "sdk-credential-update", "remove": nil},
	})).require(t)
	if updated.ID != credential.ID || updated.Metadata["suite"] != "sdk-credential-update" {
		t.Fatal("credential update response did not preserve ID and metadata patch")
	}
	got := liveResult(s.client.Vaults.Credentials.Get(ctx, credential.ID, managed.VaultCredentialGetParams{VaultID: vault.ID})).require(t)
	if got.ID != credential.ID || got.Metadata["suite"] != "sdk-credential-update" {
		t.Fatal("credential metadata update did not persist")
	}
	if _, exists := got.Metadata["remove"]; exists {
		t.Fatal("credential metadata null patch did not remove the key")
	}
}
