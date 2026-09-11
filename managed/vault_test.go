package managed_test

import (
	"context"
	managed "github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"reflect"
	"testing"
)

func TestVaultNew(t *testing.T) {
	client := contractClient(t, "VaultService", "New")
	response, err := client.Vaults.New(context.Background(), managed.VaultNewParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Vault.New")
}

func TestVaultGet(t *testing.T) {
	client := contractClient(t, "VaultService", "Get")
	response, err := client.Vaults.Get(context.Background(), "segment /?%#", managed.VaultGetParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Vault.Get")
}

func TestVaultList(t *testing.T) {
	client := contractClient(t, "VaultService", "List")
	response, err := client.Vaults.List(context.Background(), managed.VaultListParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Vault.List")
}

func TestVaultDelete(t *testing.T) {
	client := contractClient(t, "VaultService", "Delete")
	response, err := client.Vaults.Delete(context.Background(), "segment /?%#", managed.VaultDeleteParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Vault.Delete")
}

func TestVaultArchive(t *testing.T) {
	client := contractClient(t, "VaultService", "Archive")
	response, err := client.Vaults.Archive(context.Background(), "segment /?%#", managed.VaultArchiveParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Vault.Archive")
}
