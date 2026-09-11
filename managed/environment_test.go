package managed_test

import (
	"context"
	managed "github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"reflect"
	"testing"
)

func TestEnvironmentNew(t *testing.T) {
	client := contractClient(t, "EnvironmentService", "New")
	response, err := client.Environments.New(context.Background(), managed.EnvironmentNewParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Environment.New")
}

func TestEnvironmentGet(t *testing.T) {
	client := contractClient(t, "EnvironmentService", "Get")
	response, err := client.Environments.Get(context.Background(), "segment /?%#", managed.EnvironmentGetParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Environment.Get")
}

func TestEnvironmentUpdate(t *testing.T) {
	client := contractClient(t, "EnvironmentService", "Update")
	response, err := client.Environments.Update(context.Background(), "segment /?%#", managed.EnvironmentUpdateParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Environment.Update")
}

func TestEnvironmentList(t *testing.T) {
	client := contractClient(t, "EnvironmentService", "List")
	response, err := client.Environments.List(context.Background(), managed.EnvironmentListParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Environment.List")
}

func TestEnvironmentDelete(t *testing.T) {
	client := contractClient(t, "EnvironmentService", "Delete")
	response, err := client.Environments.Delete(context.Background(), "segment /?%#", managed.EnvironmentDeleteParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Environment.Delete")
}

func TestEnvironmentArchive(t *testing.T) {
	client := contractClient(t, "EnvironmentService", "Archive")
	response, err := client.Environments.Archive(context.Background(), "segment /?%#", managed.EnvironmentArchiveParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Environment.Archive")
}
