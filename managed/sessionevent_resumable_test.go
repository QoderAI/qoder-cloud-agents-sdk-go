package managed_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/managed"
)

func TestManagedResumableStreamPreservesRequestAcrossEOF(t *testing.T) {
	calls := 0
	client := testClient(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.Header.Get("Last-Event-ID") != "option" || r.Header.Get("X-Test") != "preserved" {
			t.Fatalf("call %d headers = %v", calls, r.Header)
		}
		if got := r.Header.Values("x-qoder-beta"); !reflect.DeepEqual(got, []string{"beta-one", "beta-two"}) {
			t.Fatalf("betas = %v", got)
		}
		if got := r.URL.Query()["event_deltas[]"]; !reflect.DeepEqual(got, []string{"agent.message", "agent.thinking"}) {
			t.Fatalf("event_deltas = %v", got)
		}
		body := ""
		if calls > 1 {
			body = "id: shared\nevent: event_start\ndata: {\"id\":\"shared\",\"type\":\"event_start\"}\n\n" +
				"id: shared\nevent: event_delta\ndata: {\"id\":\"shared\",\"type\":\"event_delta\"}\n\n"
		}
		response := reply(r, http.StatusOK, body)
		response.Header.Set("Content-Type", "text/event-stream")
		return response, nil
	})

	stream := client.Sessions.Events.NewResumableStream(context.Background(), "sess", managed.SessionEventStreamParams{
		EventDeltas: []managed.ManagedAgentsDeltaType{"agent.message", "agent.thinking"},
		Betas:       []managed.QoderBeta{"beta-one", "beta-two"},
		LastEventID: managed.String("initial"),
	}, option.WithHeader("X-Test", "preserved"), option.WithHeader("Last-Event-ID", "option"))
	defer stream.Close()

	if !stream.Next() || !stream.Next() {
		t.Fatalf("events ended: %v", stream.Err())
	}
	if stream.Current().Type != "event_delta" || stream.LastEventID() != "shared" || calls != 2 {
		t.Fatalf("event=%q cursor=%q calls=%d", stream.Current().Type, stream.LastEventID(), calls)
	}
}

func TestManagedResumableStreamReconnectsAfterTruncatedHTTPResponse(t *testing.T) {
	firstFrame := "id: evt-1\nevent: agent.message\ndata: {\"id\":\"evt-1\",\"type\":\"agent.message\"}\n\n"
	terminalFrame := "id: evt-2\nevent: session.deleted\ndata: {\"id\":\"evt-2\",\"type\":\"session.deleted\"}\n\n"
	var mu sync.Mutex
	var cursors []string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		cursors = append(cursors, r.Header.Get("Last-Event-ID"))
		request := len(cursors)
		mu.Unlock()

		w.Header().Set("Content-Type", "text/event-stream")
		if request == 1 {
			w.Header().Set("Content-Length", strconv.Itoa(len(firstFrame)+1))
			_, _ = fmt.Fprint(w, firstFrame)
			w.(http.Flusher).Flush()
			return
		}
		_, _ = fmt.Fprint(w, terminalFrame)
	}))
	httpClient := server.Client()
	t.Cleanup(func() {
		httpClient.CloseIdleConnections()
		server.CloseClientConnections()
		server.Close()
	})
	client := managed.NewClient(
		option.WithPAT("secret-pat"),
		option.WithBaseURL(server.URL),
		option.WithHTTPClient(httpClient),
		option.WithMaxRetries(0),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stream := client.Sessions.Events.NewResumableStream(ctx, "sess", managed.SessionEventStreamParams{})

	var ids []string
	var types []string
	for stream.Next() {
		event := stream.Current()
		ids = append(ids, event.ID)
		types = append(types, event.Type)
	}
	if err := stream.Close(); err != nil {
		t.Fatal(err)
	}
	if err := stream.Err(); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(ids, []string{"evt-1", "evt-2"}) || !reflect.DeepEqual(types, []string{"agent.message", "session.deleted"}) {
		t.Fatalf("events: ids=%v types=%v", ids, types)
	}
	if stream.LastEventID() != "evt-2" {
		t.Fatalf("cursor = %q, want evt-2", stream.LastEventID())
	}
	mu.Lock()
	gotCursors := append([]string(nil), cursors...)
	mu.Unlock()
	if !reflect.DeepEqual(gotCursors, []string{"", "evt-1"}) {
		t.Fatalf("request cursors = %v", gotCursors)
	}
}

func TestManagedResumableStreamHonorsInitialHeaderDeletionAndTerminalEvent(t *testing.T) {
	calls := 0
	client := testClient(func(r *http.Request) (*http.Response, error) {
		calls++
		if values := r.Header.Values("Last-Event-ID"); len(values) != 0 {
			t.Fatalf("Last-Event-ID values = %v", values)
		}
		response := reply(r, http.StatusOK,
			"id: terminal\nevent: session.status_terminated\ndata: {\"id\":\"terminal\",\"type\":\"session.status_terminated\"}\n\n"+
				"id: late\nevent: agent.message\ndata: {\"id\":\"late\",\"type\":\"agent.message\"}\n\n")
		response.Header.Set("Content-Type", "text/event-stream")
		return response, nil
	})

	stream := client.Sessions.Events.NewResumableStream(context.Background(), "sess", managed.SessionEventStreamParams{
		LastEventID: managed.String("typed"),
	}, option.WithHeaderDel("Last-Event-ID"))
	defer stream.Close()

	if !stream.Next() || stream.Current().Type != "session.status_terminated" || stream.LastEventID() != "terminal" {
		t.Fatalf("terminal event not yielded: event=%q cursor=%q err=%v", stream.Current().Type, stream.LastEventID(), stream.Err())
	}
	if stream.Next() || stream.Err() != nil || calls != 1 {
		t.Fatalf("stream continued: err=%v calls=%d", stream.Err(), calls)
	}
}

func TestManagedResumableStreamCloseStopsActiveRead(t *testing.T) {
	requestDone := make(chan struct{})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, "id: first\nevent: event_start\ndata: {\"id\":\"first\",\"type\":\"event_start\"}\n\n")
		w.(http.Flusher).Flush()
		<-r.Context().Done()
		close(requestDone)
	}))
	defer server.Close()
	client := managed.NewClient(
		option.WithPAT("secret-pat"),
		option.WithBaseURL(server.URL),
		option.WithHTTPClient(server.Client()),
		option.WithMaxRetries(0),
	)
	stream := client.Sessions.Events.NewResumableStream(context.Background(), "sess", managed.SessionEventStreamParams{})
	if !stream.Next() {
		t.Fatal(stream.Err())
	}
	nextDone := make(chan bool, 1)
	go func() { nextDone <- stream.Next() }()
	if err := stream.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case next := <-nextDone:
		if next || stream.Err() != nil {
			t.Fatalf("next=%v err=%v", next, stream.Err())
		}
	case <-time.After(time.Second):
		t.Fatal("close did not stop active read")
	}
	select {
	case <-requestDone:
	case <-time.After(time.Second):
		t.Fatal("close did not cancel request context")
	}
}
