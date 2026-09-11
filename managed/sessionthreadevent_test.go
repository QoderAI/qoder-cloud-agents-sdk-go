package managed_test

import (
	"context"
	managed "github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"reflect"
	"testing"
)

func TestSessionThreadEventList(t *testing.T) {
	client := contractClient(t, "SessionThreadEventService", "List")
	response, err := client.Sessions.Threads.Events.List(context.Background(), "segment /?%#", managed.SessionThreadEventListParams{SessionID: "segment /?%#"})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "SessionThreadEvent.List")
}

func TestSessionThreadEventStreamEvents(t *testing.T) {
	client := contractClient(t, "SessionThreadEventService", "StreamEvents")
	stream := client.Sessions.Threads.Events.StreamEvents(context.Background(), "segment /?%#", managed.SessionThreadEventStreamParams{SessionID: "segment /?%#"})
	defer stream.Close()
	if !stream.Next() {
		t.Fatalf("SSE produced no event: %v", stream.Err())
	}
}
