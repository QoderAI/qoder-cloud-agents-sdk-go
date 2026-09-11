package managed_test

import (
	"context"
	managed "github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"reflect"
	"testing"
)

func TestSessionThreadGet(t *testing.T) {
	client := contractClient(t, "SessionThreadService", "Get")
	response, err := client.Sessions.Threads.Get(context.Background(), "segment /?%#", managed.SessionThreadGetParams{SessionID: "segment /?%#"})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "SessionThread.Get")
}

func TestSessionThreadList(t *testing.T) {
	client := contractClient(t, "SessionThreadService", "List")
	response, err := client.Sessions.Threads.List(context.Background(), "segment /?%#", managed.SessionThreadListParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "SessionThread.List")
}

func TestSessionThreadArchive(t *testing.T) {
	client := contractClient(t, "SessionThreadService", "Archive")
	response, err := client.Sessions.Threads.Archive(context.Background(), "segment /?%#", managed.SessionThreadArchiveParams{SessionID: "segment /?%#"})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "SessionThread.Archive")
}
