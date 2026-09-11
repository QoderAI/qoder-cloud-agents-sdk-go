package managed_test

import (
	"context"
	managed "github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"reflect"
	"testing"
)

func TestMemoryStoreMemoryVersionGet(t *testing.T) {
	client := contractClient(t, "MemoryStoreMemoryVersionService", "Get")
	response, err := client.MemoryStores.MemoryVersions.Get(context.Background(), "segment /?%#", managed.MemoryStoreMemoryVersionGetParams{MemoryStoreID: "segment /?%#"})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "MemoryStoreMemoryVersion.Get")
}

func TestMemoryStoreMemoryVersionList(t *testing.T) {
	client := contractClient(t, "MemoryStoreMemoryVersionService", "List")
	response, err := client.MemoryStores.MemoryVersions.List(context.Background(), "segment /?%#", managed.MemoryStoreMemoryVersionListParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "MemoryStoreMemoryVersion.List")
}

func TestMemoryStoreMemoryVersionRedact(t *testing.T) {
	client := contractClient(t, "MemoryStoreMemoryVersionService", "Redact")
	response, err := client.MemoryStores.MemoryVersions.Redact(context.Background(), "segment /?%#", managed.MemoryStoreMemoryVersionRedactParams{MemoryStoreID: "segment /?%#"})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "MemoryStoreMemoryVersion.Redact")
}
