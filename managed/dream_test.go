package managed_test

import (
	"context"
	managed "github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"reflect"
	"testing"
)

func TestDreamNew(t *testing.T) {
	client := contractClient(t, "DreamService", "New")
	response, err := client.Dreams.New(context.Background(), managed.DreamNewParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Dream.New")
}

func TestDreamGet(t *testing.T) {
	client := contractClient(t, "DreamService", "Get")
	response, err := client.Dreams.Get(context.Background(), "segment /?%#", managed.DreamGetParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Dream.Get")
}

func TestDreamList(t *testing.T) {
	client := contractClient(t, "DreamService", "List")
	response, err := client.Dreams.List(context.Background(), managed.DreamListParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Dream.List")
}

func TestDreamArchive(t *testing.T) {
	client := contractClient(t, "DreamService", "Archive")
	response, err := client.Dreams.Archive(context.Background(), "segment /?%#", managed.DreamArchiveParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Dream.Archive")
}

func TestDreamCancel(t *testing.T) {
	client := contractClient(t, "DreamService", "Cancel")
	response, err := client.Dreams.Cancel(context.Background(), "segment /?%#", managed.DreamCancelParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Dream.Cancel")
}
