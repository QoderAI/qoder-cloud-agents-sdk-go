package managed_test

import (
	"context"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/ssestream"
	managed "github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
	"io"
	"net/http"
	"reflect"
	"testing"
	"time"
)

func TestSessionEventList(t *testing.T) {
	client := contractClient(t, "SessionEventService", "List")
	response, err := client.Sessions.Events.List(context.Background(), "segment /?%#", managed.SessionEventListParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "SessionEvent.List")
}

func TestSessionEventSend(t *testing.T) {
	client := contractClient(t, "SessionEventService", "Send")
	response, err := client.Sessions.Events.Send(context.Background(), "segment /?%#", managed.SessionEventSendParams{})
	if err != nil {
		t.Fatal(err)
	}
	checkFields(t, reflect.ValueOf(response), "SessionEvent.Send")
}

func TestSessionEventStreamEvents(t *testing.T) {
	client := contractClient(t, "SessionEventService", "StreamEvents")
	stream := client.Sessions.Events.StreamEvents(context.Background(), "segment /?%#", managed.SessionEventStreamParams{})
	defer stream.Close()
	if !stream.Next() {
		t.Fatalf("SSE produced no event: %v", stream.Err())
	}
}

func TestStreamingEventsAndResume(t *testing.T) {
	c := testClient(func(r *http.Request) (*http.Response, error) {
		if r.Header.Get("Accept") != "text/event-stream" || r.Header.Get("Last-Event-ID") != "resume" || len(r.URL.Query()["event_deltas[]"]) != 2 {
			t.Fatal("stream request conventions", r.Header, r.URL)
		}
		res := reply(r, 200, ": heartbeat\n\nevent: ping\ndata: {}\n\nid: shared\nevent: event_start\ndata: {\"type\":\"event_start\",\"id\":\"shared\"}\n\nid: shared\nevent: event_delta\ndata: {\"type\":\"event_delta\",\"id\":\"shared\"}\n\nid: final\nevent: agent.message\ndata: {\"type\":\"agent.message\",\"id\":\"final\",\"content\":[{\"type\":\"text\",\"text\":\"hello\"}]}\n\nevent: future.event\ndata: {\"type\":\"future.event\",\"id\":\"future\"}\n\n")
		res.Header.Set("Content-Type", "text/event-stream")
		return res, nil
	})
	stream := c.Sessions.Events.StreamEvents(context.Background(), "s", managed.SessionEventStreamParams{EventDeltas: []managed.ManagedAgentsDeltaType{"agent.message", "agent.thinking"}}, option.WithHeader("Last-Event-ID", "resume"))
	defer stream.Close()
	var types []string
	for stream.Next() {
		types = append(types, string(stream.Current().Type))
		if stream.Current().RawJSON() == "" {
			t.Fatal("raw stream event lost")
		}
	}
	if stream.Err() != nil || len(types) != 4 {
		t.Fatal(stream.Err(), types)
	}
	if types[0] != "event_start" || types[1] != "event_delta" || types[3] != "future.event" || stream.LastEventID() != "final" {
		t.Fatal(types, stream.LastEventID())
	}
	// Exercise a truly incremental stream: Next must return before the server closes.
	reader, writer := io.Pipe()
	res := &http.Response{Header: http.Header{}, Body: reader}
	s := ssestream.NewStream[map[string]any](ssestream.NewDecoder(res), nil)
	defer s.Close()
	done := make(chan bool, 1)
	go func() { done <- s.Next() }()
	go func() { _, _ = io.WriteString(writer, "data: {\"type\":\"message\"}\n\n") }()
	select {
	case ok := <-done:
		if !ok {
			t.Fatal(s.Err())
		}
	case <-time.After(time.Second):
		writer.Close()
		t.Fatal("stream buffered until EOF")
	}
	writer.Close()
}
