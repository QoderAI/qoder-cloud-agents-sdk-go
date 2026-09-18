package forward_test

import (
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/forward"
)

func TestForwardResumableStreamPreservesRequestAndCursor(t *testing.T) {
	calls := 0
	client := testClient(func(r *http.Request) (*http.Response, error) {
		calls++
		if r.Header.Get("X-Test") != "preserved" {
			t.Fatalf("option header = %q", r.Header.Get("X-Test"))
		}
		if got := r.URL.Query()["event_deltas[]"]; !reflect.DeepEqual(got, []string{"agent.message", "agent.thinking"}) {
			t.Fatalf("event_deltas = %v", got)
		}
		if r.URL.Query().Get("include_tool_calls") != "true" || r.URL.Query().Get("include_thinking") != "true" {
			t.Fatalf("stream params lost: %s", r.URL.RawQuery)
		}
		wantCursor := "option"
		if calls > 1 {
			wantCursor = "checkpoint"
		}
		if got := r.Header.Get("Last-Event-ID"); got != wantCursor {
			t.Fatalf("call %d cursor = %q, want %q", calls, got, wantCursor)
		}

		var body string
		if calls == 1 {
			body = "id: shared\nevent: event_start\ndata: {\"id\":\"shared\",\"type\":\"event_start\"}\n\n" +
				"id: shared\nevent: event_delta\ndata: {\"id\":\"shared\",\"type\":\"event_delta\"}\n\n" +
				"id: checkpoint\nevent: agent.message\ndata: {\"id\":\"checkpoint\",\"type\":\"agent.message\"}\n\n" +
				"id: partial\nevent: event_delta\ndata: {\"id\":\"partial\""
		} else {
			body = "id: resumed\nevent: session.status_idle\ndata: {\"id\":\"resumed\",\"type\":\"session.status_idle\"}\n\n"
		}
		response := reply(r, http.StatusOK, body)
		response.Header.Set("Content-Type", "text/event-stream")
		return response, nil
	})

	stream := client.Sessions.Events.NewResumableStream(context.Background(), "sess", forward.SessionEventStreamParams{
		EventDeltas:      []string{"agent.message", "agent.thinking"},
		IncludeToolCalls: forward.Bool(true),
		IncludeThinking:  forward.Bool(true),
		LastEventID:      forward.String("initial"),
	}, option.WithHeader("X-Test", "preserved"), option.WithHeader("Last-Event-ID", "option"))
	defer stream.Close()

	var types []string
	for len(types) < 4 && stream.Next() {
		types = append(types, stream.Current().Type)
	}
	if err := stream.Err(); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(types, []string{"event_start", "event_delta", "agent.message", "session.status_idle"}) {
		t.Fatalf("types = %v", types)
	}
	if calls != 2 || stream.LastEventID() != "resumed" {
		t.Fatalf("calls=%d cursor=%q", calls, stream.LastEventID())
	}
}

func TestForwardResumableStreamReconnectsAfterTruncatedHTTPResponse(t *testing.T) {
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
			_, _ = io.WriteString(w, firstFrame)
			w.(http.Flusher).Flush()
			return
		}
		_, _ = io.WriteString(w, terminalFrame)
	}))
	httpClient := server.Client()
	t.Cleanup(func() {
		httpClient.CloseIdleConnections()
		server.CloseClientConnections()
		server.Close()
	})
	client := forward.NewClient(
		option.WithPAT("secret-pat"),
		option.WithBaseURL(server.URL),
		option.WithHTTPClient(httpClient),
		option.WithMaxRetries(0),
	)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	stream := client.Sessions.Events.NewResumableStream(ctx, "sess", forward.SessionEventStreamParams{})

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

func TestForwardResumableStreamHonorsInitialHeaderDeletionAndTerminalEvent(t *testing.T) {
	calls := 0
	client := testClient(func(r *http.Request) (*http.Response, error) {
		calls++
		if values := r.Header.Values("Last-Event-ID"); len(values) != 0 {
			t.Fatalf("Last-Event-ID values = %v", values)
		}
		response := reply(r, http.StatusOK,
			"id: terminal\nevent: session.deleted\ndata: {\"id\":\"terminal\",\"type\":\"session.deleted\"}\n\n"+
				"id: late\nevent: agent.message\ndata: {\"id\":\"late\",\"type\":\"agent.message\"}\n\n")
		response.Header.Set("Content-Type", "text/event-stream")
		return response, nil
	})

	stream := client.Sessions.Events.NewResumableStream(context.Background(), "sess", forward.SessionEventStreamParams{
		LastEventID: forward.String("typed"),
	}, option.WithHeaderDel("Last-Event-ID"))
	defer stream.Close()

	if !stream.Next() || stream.Current().Type != "session.deleted" || stream.LastEventID() != "terminal" {
		t.Fatalf("terminal event not yielded: event=%q cursor=%q err=%v", stream.Current().Type, stream.LastEventID(), stream.Err())
	}
	if stream.Next() || stream.Err() != nil || calls != 1 {
		t.Fatalf("stream continued: err=%v calls=%d", stream.Err(), calls)
	}
}

func TestForwardResumableStreamDoesNotRetryConflict(t *testing.T) {
	calls := 0
	client := testClient(func(r *http.Request) (*http.Response, error) {
		calls++
		response := reply(r, http.StatusConflict, `{"error":{"type":"conflict_error","message":"running"}}`)
		response.Header.Set("x-should-retry", "true")
		return response, nil
	})
	stream := client.Sessions.Events.NewResumableStream(context.Background(), "sess", forward.SessionEventStreamParams{})
	defer stream.Close()
	if stream.Next() {
		t.Fatal("conflict yielded an event")
	}
	var apiErr *convention.Error
	if !errors.As(stream.Err(), &apiErr) || apiErr.StatusCode != http.StatusConflict || calls != 1 {
		t.Fatalf("err=%v calls=%d", stream.Err(), calls)
	}
}
