package managed_test

import (
	"context"
	"errors"
	"fmt"
	"io"
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

func TestManagedThreadResumableStreamReconnects(t *testing.T) {
	for _, ending := range []string{"eof", "truncated"} {
		t.Run(ending, func(t *testing.T) {
			first := "id: checkpoint\nevent: session.thread_status_idle\ndata: {\"id\":\"payload-id\",\"type\":\"session.thread_status_idle\"}\n\n" +
				"id: partial\ndata: {\"id\":\"not-delivered\"}"
			last := "id: terminal\nevent: session.thread_status_terminated\ndata: {\"type\":\"session.thread_status_terminated\"}\n\n"
			var mu sync.Mutex
			var cursors []string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				mu.Lock()
				defer mu.Unlock()
				cursors = append(cursors, r.Header.Get("Last-Event-ID"))
				if r.Method != http.MethodGet || r.URL.Path != "/sessions/sess/threads/thread/stream" || r.URL.RawQuery != "" {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL)
				}
				if r.Header.Get("X-Test") != "kept" || r.Header.Get("Accept") != "text/event-stream" {
					t.Errorf("headers lost: %v", r.Header)
				}
				if r.Header.Get("x-qoder-beta") != "beta-one" {
					t.Errorf("beta lost")
				}
				w.Header().Set("Content-Type", "text/event-stream")
				if len(cursors) == 1 {
					if ending == "truncated" {
						w.Header().Set("Content-Length", strconv.Itoa(len(first)+128))
					}
					_, _ = io.WriteString(w, first)
				} else {
					_, _ = io.WriteString(w, last)
				}
			}))
			defer server.Close()
			client := managed.NewClient(option.WithPAT("test"), option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()), option.WithMaxRetries(0))
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			stream := client.Sessions.Threads.Events.NewResumableStream(ctx, "thread", managed.SessionThreadEventStreamParams{SessionID: "sess", LastEventID: managed.String("initial"), Betas: []managed.QoderBeta{"beta-one"}}, option.WithHeader("X-Test", "kept"), option.WithHeader("last-event-id", "override"))
			defer stream.Close()
			var types []string
			for stream.Next() {
				types = append(types, stream.Current().Type)
			}
			if stream.Err() != nil {
				t.Fatal(stream.Err())
			}
			if !reflect.DeepEqual(types, []string{"session.thread_status_idle", "session.thread_status_terminated"}) || stream.LastEventID() != "terminal" {
				t.Fatalf("types=%v cursor=%s", types, stream.LastEventID())
			}
			mu.Lock()
			defer mu.Unlock()
			if !reflect.DeepEqual(cursors, []string{"override", "checkpoint"}) {
				t.Fatalf("cursors=%v", cursors)
			}
		})
	}
}

func TestManagedThreadResumableStreamHTTPPolicy(t *testing.T) {
	for _, status := range []int{400, 401, 403, 404, 409, 429, 503} {
		t.Run(strconv.Itoa(status), func(t *testing.T) {
			calls := 0
			client := testClient(func(r *http.Request) (*http.Response, error) {
				calls++
				if r.Header.Get("Last-Event-ID") != "initial" {
					t.Errorf("cursor lost: %v", r.Header)
				}
				if calls == 1 {
					return reply(r, status, `{"error":{"type":"invalid_request_error","message":"test failure"}}`), nil
				}
				return reply(r, 200, "event: session.thread_status_terminated\ndata: {\"type\":\"session.thread_status_terminated\"}\n\n"), nil
			}, option.WithMaxRetries(0))
			ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			stream := client.Sessions.Threads.Events.NewResumableStream(ctx, "thread", managed.SessionThreadEventStreamParams{SessionID: "sess", LastEventID: managed.String("initial"), Betas: []managed.QoderBeta{"beta-one"}})
			defer stream.Close()
			got := stream.Next()
			if status == 429 || status == 503 {
				if !got || stream.Next() || stream.Err() != nil || calls != 2 {
					t.Fatalf("retry failed: calls=%d err=%v", calls, stream.Err())
				}
			} else if got || stream.Err() == nil || calls != 1 {
				t.Fatalf("non-retryable request repeated: calls=%d err=%v", calls, stream.Err())
			}
		})
	}
}

func TestManagedThreadResumableStreamCancellation(t *testing.T) {
	for _, closeStream := range []bool{false, true} {
		t.Run(fmt.Sprint(closeStream), func(t *testing.T) {
			started := make(chan struct{})
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "text/event-stream")
				w.(http.Flusher).Flush()
				close(started)
				<-r.Context().Done()
			}))
			defer server.Close()
			client := managed.NewClient(option.WithPAT("test"), option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()))
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			stream := client.Sessions.Threads.Events.NewResumableStream(ctx, "thread", managed.SessionThreadEventStreamParams{SessionID: "sess", LastEventID: managed.String("initial"), Betas: []managed.QoderBeta{"beta-one"}})
			defer stream.Close()
			done := make(chan bool, 1)
			go func() { done <- stream.Next() }()
			select {
			case <-started:
			case <-time.After(3 * time.Second):
				t.Fatal("request did not start")
			}
			if closeStream {
				_ = stream.Close()
			} else {
				cancel()
			}
			select {
			case got := <-done:
				if got {
					t.Fatal("yielded after cancellation")
				}
			case <-time.After(3 * time.Second):
				t.Fatal("read not interrupted")
			}
			if !closeStream && !errors.Is(stream.Err(), context.Canceled) {
				t.Fatalf("error=%v", stream.Err())
			}
		})
	}
}

func TestManagedThreadResumableStreamMissingIDs(t *testing.T) {
	client := managed.NewClient(option.WithPAT("test"))
	for _, ids := range [][2]string{{"", "thread"}, {"sess", ""}} {
		stream := client.Sessions.Threads.Events.NewResumableStream(context.Background(), ids[1], managed.SessionThreadEventStreamParams{SessionID: ids[0]})
		if stream.Next() || stream.Err() == nil {
			t.Fatal("missing IDs accepted")
		}
		_ = stream.Close()
	}
}
