package managed_test

import (
	"context"
	managed "github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"reflect"
	"testing"
)

func TestMemoryStoreNew(t *testing.T) {
	client := contractClient(t, "MemoryStoreService", "New")
	response, err := client.MemoryStores.New(context.Background(), managed.MemoryStoreNewParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "MemoryStore.New")
}

func TestMemoryStoreGet(t *testing.T) {
	client := contractClient(t, "MemoryStoreService", "Get")
	response, err := client.MemoryStores.Get(context.Background(), "segment /?%#", managed.MemoryStoreGetParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "MemoryStore.Get")
}

func TestMemoryStoreUpdate(t *testing.T) {
	client := contractClient(t, "MemoryStoreService", "Update")
	response, err := client.MemoryStores.Update(context.Background(), "segment /?%#", managed.MemoryStoreUpdateParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "MemoryStore.Update")
}

func TestMemoryStoreList(t *testing.T) {
	client := contractClient(t, "MemoryStoreService", "List")
	response, err := client.MemoryStores.List(context.Background(), managed.MemoryStoreListParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "MemoryStore.List")
}

func TestMemoryStoreDelete(t *testing.T) {
	client := contractClient(t, "MemoryStoreService", "Delete")
	response, err := client.MemoryStores.Delete(context.Background(), "segment /?%#", managed.MemoryStoreDeleteParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "MemoryStore.Delete")
}

func TestMemoryStoreArchive(t *testing.T) {
	client := contractClient(t, "MemoryStoreService", "Archive")
	response, err := client.MemoryStores.Archive(context.Background(), "segment /?%#", managed.MemoryStoreArchiveParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "MemoryStore.Archive")
}
