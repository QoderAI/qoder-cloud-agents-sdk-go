package managed_test

import (
	"context"
	managed "github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"reflect"
	"testing"
)

func TestSessionResourceGet(t *testing.T) {
	client := contractClient(t, "SessionResourceService", "Get")
	response, err := client.Sessions.Resources.Get(context.Background(), "segment /?%#", managed.SessionResourceGetParams{SessionID: "segment /?%#"})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "SessionResource.Get")
}

func TestSessionResourceUpdate(t *testing.T) {
	client := contractClient(t, "SessionResourceService", "Update")
	response, err := client.Sessions.Resources.Update(context.Background(), "segment /?%#", managed.SessionResourceUpdateParams{SessionID: "segment /?%#"})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "SessionResource.Update")
}

func TestSessionResourceList(t *testing.T) {
	client := contractClient(t, "SessionResourceService", "List")
	response, err := client.Sessions.Resources.List(context.Background(), "segment /?%#", managed.SessionResourceListParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "SessionResource.List")
}

func TestSessionResourceDelete(t *testing.T) {
	client := contractClient(t, "SessionResourceService", "Delete")
	response, err := client.Sessions.Resources.Delete(context.Background(), "segment /?%#", managed.SessionResourceDeleteParams{SessionID: "segment /?%#"})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "SessionResource.Delete")
}

func TestSessionResourceAdd(t *testing.T) {
	client := contractClient(t, "SessionResourceService", "Add")
	response, err := client.Sessions.Resources.Add(context.Background(), "segment /?%#", managed.SessionResourceAddParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "SessionResource.Add")
}
