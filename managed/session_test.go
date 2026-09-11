package managed_test

import (
	"context"
	managed "github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"reflect"
	"testing"
)

func TestSessionNew(t *testing.T) {
	client := contractClient(t, "SessionService", "New")
	response, err := client.Sessions.New(context.Background(), managed.SessionNewParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Session.New")
}

func TestSessionGet(t *testing.T) {
	client := contractClient(t, "SessionService", "Get")
	response, err := client.Sessions.Get(context.Background(), "segment /?%#", managed.SessionGetParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Session.Get")
}

func TestSessionUpdate(t *testing.T) {
	client := contractClient(t, "SessionService", "Update")
	response, err := client.Sessions.Update(context.Background(), "segment /?%#", managed.SessionUpdateParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Session.Update")
}

func TestSessionList(t *testing.T) {
	client := contractClient(t, "SessionService", "List")
	response, err := client.Sessions.List(context.Background(), managed.SessionListParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Session.List")
}

func TestSessionDelete(t *testing.T) {
	client := contractClient(t, "SessionService", "Delete")
	response, err := client.Sessions.Delete(context.Background(), "segment /?%#", managed.SessionDeleteParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Session.Delete")
}

func TestSessionArchive(t *testing.T) {
	client := contractClient(t, "SessionService", "Archive")
	response, err := client.Sessions.Archive(context.Background(), "segment /?%#", managed.SessionArchiveParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "Session.Archive")
}
