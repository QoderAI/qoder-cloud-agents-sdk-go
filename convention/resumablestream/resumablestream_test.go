package resumablestream

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"sync"
	"testing"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apierror"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/ssestream"
)

type fakeEvent[T any] struct {
	value T
	id    string
}

type fakeChild[T any] struct {
	events     []fakeEvent[T]
	err        error
	index      int
	current    fakeEvent[T]
	closed     int
	beforeNext func()
}

func (s *fakeChild[T]) Next() bool {
	if s.index >= len(s.events) {
		return false
	}
	if s.beforeNext != nil {
		s.beforeNext()
	}
	s.current = s.events[s.index]
	s.index++
	return true
}
func (s *fakeChild[T]) Current() T          { return s.current.value }
func (s *fakeChild[T]) Err() error          { return s.err }
func (s *fakeChild[T]) LastEventID() string { return s.current.id }
func (s *fakeChild[T]) Close() error        { s.closed++; return nil }

type failingBody struct{ err error }

func (b failingBody) Read([]byte) (int, error) { return 0, b.err }
func (failingBody) Close() error               { return nil }

func immediateSettings(delays *[]time.Duration) settings {
	return settings{
		initial:      500 * time.Millisecond,
		maximum:      10 * time.Second,
		healthyAfter: 5 * time.Second,
		jitter:       func(delay time.Duration) time.Duration { return delay },
		wait: func(_ context.Context, delay time.Duration) error {
			*delays = append(*delays, delay)
			return nil
		},
		now: time.Now,
	}
}

func TestStreamAdvancesCursorAndBacksOff(t *testing.T) {
	children := []*fakeChild[string]{
		{err: &url.Error{Op: "Get", URL: "https://sdk.test", Err: errors.New("dial failed")}},
		{},
		{events: []fakeEvent[string]{{value: "first", id: "evt-1"}}, err: io.ErrUnexpectedEOF},
		{events: []fakeEvent[string]{{value: "second", id: "evt-2"}}},
	}
	var cursors []string
	var resumes []bool
	var delays []time.Duration
	stream := newStream(context.Background(), "initial", func(_ context.Context, cursor string, resume bool) Child[string] {
		cursors = append(cursors, cursor)
		resumes = append(resumes, resume)
		child := children[0]
		children = children[1:]
		return child
	}, nil, immediateSettings(&delays))
	defer stream.Close()

	if !stream.Next() || stream.Current() != "first" || stream.LastEventID() != "evt-1" {
		t.Fatalf("first event: current=%q cursor=%q err=%v", stream.Current(), stream.LastEventID(), stream.Err())
	}
	if !stream.Next() || stream.Current() != "second" || stream.LastEventID() != "evt-2" {
		t.Fatalf("second event: current=%q cursor=%q err=%v", stream.Current(), stream.LastEventID(), stream.Err())
	}
	if !reflect.DeepEqual(cursors, []string{"initial", "initial", "initial", "evt-1"}) {
		t.Fatalf("cursors = %v", cursors)
	}
	if !reflect.DeepEqual(resumes, []bool{false, false, false, true}) {
		t.Fatalf("resumes = %v", resumes)
	}
	if !reflect.DeepEqual(delays, []time.Duration{500 * time.Millisecond, time.Second, 2 * time.Second}) {
		t.Fatalf("delays = %v", delays)
	}
}

func TestStreamResetsBackoffOnlyAfterFiveHealthySeconds(t *testing.T) {
	for _, tc := range []struct {
		name       string
		healthyFor time.Duration
		wantDelay  time.Duration
	}{
		{name: "exactly five seconds", healthyFor: 5 * time.Second, wantDelay: 4 * time.Second},
		{name: "over five seconds", healthyFor: 5*time.Second + time.Nanosecond, wantDelay: 500 * time.Millisecond},
	} {
		t.Run(tc.name, func(t *testing.T) {
			now := time.Unix(0, 0)
			first := &fakeChild[string]{
				events: []fakeEvent[string]{{value: "first", id: "evt-1"}},
				err:    io.ErrUnexpectedEOF,
				beforeNext: func() {
					now = now.Add(tc.healthyFor)
				},
			}
			second := &fakeChild[string]{events: []fakeEvent[string]{{value: "second", id: "evt-2"}}}
			children := []*fakeChild[string]{first, second}
			var delays []time.Duration
			config := immediateSettings(&delays)
			config.now = func() time.Time { return now }
			stream := newStream(context.Background(), "", func(context.Context, string, bool) Child[string] {
				child := children[0]
				children = children[1:]
				return child
			}, nil, config)
			stream.backoff = 4 * time.Second
			defer stream.Close()

			if !stream.Next() || !stream.Next() {
				t.Fatalf("stream ended: %v", stream.Err())
			}
			if !reflect.DeepEqual(delays, []time.Duration{tc.wantDelay}) {
				t.Fatalf("delays = %v, want [%v]", delays, tc.wantDelay)
			}
		})
	}
}

func TestStreamYieldsTerminalEventThenEndsNormally(t *testing.T) {
	child := &fakeChild[string]{
		events: []fakeEvent[string]{
			{value: "terminal", id: "evt-terminal"},
			{value: "must-not-yield", id: "evt-late"},
		},
		err: errors.New("must not retry"),
	}
	opens := 0
	var delays []time.Duration
	stream := newStream(context.Background(), "", func(context.Context, string, bool) Child[string] {
		opens++
		return child
	}, func(event string) bool { return event == "terminal" }, immediateSettings(&delays))
	defer stream.Close()

	if !stream.Next() || stream.Current() != "terminal" || stream.LastEventID() != "evt-terminal" {
		t.Fatalf("terminal event not yielded: current=%q cursor=%q err=%v", stream.Current(), stream.LastEventID(), stream.Err())
	}
	next := stream.Next()
	if next || stream.Err() != nil || opens != 1 || len(delays) != 0 || child.closed != 1 {
		t.Fatalf("next=%v err=%v opens=%d delays=%v closes=%d", next, stream.Err(), opens, delays, child.closed)
	}
}

func TestStreamRetriesPartialFrameAndCleanEOF(t *testing.T) {
	for _, tc := range []struct {
		name string
		body string
	}{
		{name: "partial frame", body: "id: partial\ndata: {\"id\":\"partial\""},
		{name: "clean EOF", body: ""},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			var cursors []string
			var delays []time.Duration
			stream := newStream(context.Background(), "start", func(_ context.Context, cursor string, _ bool) Child[map[string]string] {
				calls++
				cursors = append(cursors, cursor)
				body := tc.body
				if calls > 1 {
					body = "id: complete\ndata: {\"id\":\"complete\",\"type\":\"message\"}\n\n"
				}
				response := &http.Response{StatusCode: http.StatusOK, Header: http.Header{}, Body: io.NopCloser(bytes.NewBufferString(body))}
				return ssestream.NewStream[map[string]string](ssestream.NewDecoder(response), nil)
			}, nil, immediateSettings(&delays))
			defer stream.Close()

			if !stream.Next() || stream.Current()["id"] != "complete" {
				t.Fatalf("event=%v err=%v", stream.Current(), stream.Err())
			}
			if stream.LastEventID() != "complete" || !reflect.DeepEqual(cursors, []string{"start", "start"}) || !reflect.DeepEqual(delays, []time.Duration{500 * time.Millisecond}) {
				t.Fatalf("cursor=%q opens=%v delays=%v", stream.LastEventID(), cursors, delays)
			}
		})
	}
}

func TestStreamResumesFromLastCompleteFrame(t *testing.T) {
	calls := 0
	var cursors []string
	var delays []time.Duration
	stream := newStream(context.Background(), "start", func(_ context.Context, cursor string, _ bool) Child[map[string]string] {
		calls++
		cursors = append(cursors, cursor)
		body := "id: checkpoint\ndata: {\"id\":\"checkpoint\"}\n\nid: incomplete\ndata: {\"id\":\"incomplete\"}"
		if calls > 1 {
			body = "id: resumed\ndata: {\"id\":\"resumed\"}\n\n"
		}
		response := &http.Response{StatusCode: http.StatusOK, Header: http.Header{}, Body: io.NopCloser(bytes.NewBufferString(body))}
		return ssestream.NewStream[map[string]string](ssestream.NewDecoder(response), nil)
	}, nil, immediateSettings(&delays))
	defer stream.Close()

	if !stream.Next() || stream.Current()["id"] != "checkpoint" {
		t.Fatalf("checkpoint event=%v err=%v", stream.Current(), stream.Err())
	}
	if !stream.Next() || stream.Current()["id"] != "resumed" {
		t.Fatalf("resumed event=%v err=%v", stream.Current(), stream.Err())
	}
	if !reflect.DeepEqual(cursors, []string{"start", "checkpoint"}) || !reflect.DeepEqual(delays, []time.Duration{500 * time.Millisecond}) {
		t.Fatalf("cursors=%v delays=%v", cursors, delays)
	}
}

func TestStreamRetriesBodyReadFailure(t *testing.T) {
	cause := errors.New("connection interrupted")
	calls := 0
	var delays []time.Duration
	stream := newStream(context.Background(), "checkpoint", func(context.Context, string, bool) Child[map[string]string] {
		calls++
		var body io.ReadCloser = failingBody{err: cause}
		if calls > 1 {
			body = io.NopCloser(bytes.NewBufferString("id: resumed\ndata: {\"id\":\"resumed\"}\n\n"))
		}
		response := &http.Response{StatusCode: http.StatusOK, Header: http.Header{}, Body: body}
		return ssestream.NewStream[map[string]string](ssestream.NewDecoder(response), nil)
	}, nil, immediateSettings(&delays))
	defer stream.Close()

	if !stream.Next() || stream.Current()["id"] != "resumed" {
		t.Fatalf("event=%v err=%v", stream.Current(), stream.Err())
	}
	if calls != 2 || !reflect.DeepEqual(delays, []time.Duration{500 * time.Millisecond}) {
		t.Fatalf("calls=%d delays=%v", calls, delays)
	}
}

func TestStreamKeepsCursorOnNonRetryableError(t *testing.T) {
	calls := 0
	stream := newStream(context.Background(), "invalid", func(_ context.Context, cursor string, _ bool) Child[string] {
		calls++
		if cursor != "invalid" {
			t.Fatalf("cursor = %q", cursor)
		}
		return &fakeChild[string]{err: &apierror.Error{StatusCode: http.StatusNotFound, Response: &http.Response{StatusCode: http.StatusNotFound, Header: http.Header{}}}}
	}, nil, immediateSettings(new([]time.Duration)))
	defer stream.Close()

	if stream.Next() || stream.LastEventID() != "invalid" || calls != 1 {
		t.Fatalf("next succeeded or state changed: cursor=%q calls=%d", stream.LastEventID(), calls)
	}
	var apiErr *apierror.Error
	if !errors.As(stream.Err(), &apiErr) || apiErr.StatusCode != http.StatusNotFound {
		t.Fatalf("err = %v", stream.Err())
	}
}

func TestStreamDoesNotRetryContextOrCredentialErrors(t *testing.T) {
	for _, tc := range []struct {
		name string
		err  error
	}{
		{name: "context canceled", err: context.Canceled},
		{name: "deadline exceeded", err: context.DeadlineExceeded},
		{name: "no credentials", err: &convention.NoCredentialsError{}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			calls := 0
			var delays []time.Duration
			stream := newStream(context.Background(), "", func(context.Context, string, bool) Child[string] {
				calls++
				return &fakeChild[string]{err: tc.err}
			}, nil, immediateSettings(&delays))
			defer stream.Close()

			if stream.Next() || !errors.Is(stream.Err(), tc.err) || calls != 1 || len(delays) != 0 {
				t.Fatalf("err=%v calls=%d delays=%v", stream.Err(), calls, delays)
			}
		})
	}
}

func TestRetryClassificationMatchesRequestPolicy(t *testing.T) {
	api := func(status int, retryHeader string) error {
		response := &http.Response{StatusCode: status, Header: http.Header{}}
		if retryHeader != "" {
			response.Header.Set("x-should-retry", retryHeader)
		}
		return &apierror.Error{StatusCode: status, Response: response}
	}
	var value any
	partialJSON := json.Unmarshal([]byte(`{"id"`), &value)
	malformedJSON := json.Unmarshal([]byte(`{broken}`), &value)
	var typed struct {
		ID int `json:"id"`
	}
	typeJSON := json.Unmarshal([]byte(`{"id":"wrong"}`), &typed)
	for _, tc := range []struct {
		name string
		err  error
		want bool
	}{
		{name: "transport", err: io.ErrUnexpectedEOF, want: true},
		{name: "request timeout", err: api(http.StatusRequestTimeout, ""), want: true},
		{name: "rate limit", err: api(http.StatusTooManyRequests, ""), want: true},
		{name: "server error", err: api(http.StatusServiceUnavailable, ""), want: true},
		{name: "forced retry", err: api(http.StatusBadRequest, "true"), want: true},
		{name: "forced stop", err: api(http.StatusServiceUnavailable, "false"), want: false},
		{name: "conflict cannot be forced", err: api(http.StatusConflict, "true"), want: false},
		{name: "bad request", err: api(http.StatusBadRequest, ""), want: false},
		{name: "cancellation", err: context.Canceled, want: false},
		{name: "deadline", err: context.DeadlineExceeded, want: false},
		{name: "no credentials", err: &convention.NoCredentialsError{}, want: false},
		{name: "partial JSON", err: partialJSON, want: false},
		{name: "malformed JSON", err: malformedJSON, want: false},
		{name: "JSON type error", err: typeJSON, want: false},
		{name: "arbitrary validation error", err: errors.New("missing required parameter"), want: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := retryable(context.Background(), tc.err); got != tc.want {
				t.Fatalf("retryable(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}

type blockingChild struct {
	done      chan struct{}
	closeOnce sync.Once
}

func (s *blockingChild) Next() bool {
	<-s.done
	return false
}
func (s *blockingChild) Current() string     { return "" }
func (s *blockingChild) Err() error          { return nil }
func (s *blockingChild) LastEventID() string { return "" }
func (s *blockingChild) Close() error        { s.closeOnce.Do(func() { close(s.done) }); return nil }

func TestStreamCancellationClosesActiveChild(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	child := &blockingChild{done: make(chan struct{})}
	stream := New(ctx, "", func(context.Context, string, bool) Child[string] { return child }, nil)
	done := make(chan bool, 1)
	go func() { done <- stream.Next() }()

	cancel()
	select {
	case next := <-done:
		if next || !errors.Is(stream.Err(), context.Canceled) {
			t.Fatalf("next=%v err=%v", next, stream.Err())
		}
	case <-time.After(time.Second):
		t.Fatal("cancellation did not stop active child")
	}
}

func TestStreamCloseInterruptsBackoff(t *testing.T) {
	waiting := make(chan struct{})
	stream := newStream(context.Background(), "", func(context.Context, string, bool) Child[string] {
		return &fakeChild[string]{}
	}, nil, settings{
		initial:      500 * time.Millisecond,
		maximum:      10 * time.Second,
		healthyAfter: 5 * time.Second,
		jitter:       func(delay time.Duration) time.Duration { return delay },
		wait: func(ctx context.Context, _ time.Duration) error {
			close(waiting)
			<-ctx.Done()
			return ctx.Err()
		},
		now: time.Now,
	})
	done := make(chan bool, 1)
	go func() { done <- stream.Next() }()
	<-waiting
	if err := stream.Close(); err != nil {
		t.Fatal(err)
	}
	select {
	case next := <-done:
		if next || stream.Err() != nil {
			t.Fatalf("next=%v err=%v", next, stream.Err())
		}
	case <-time.After(time.Second):
		t.Fatal("close did not interrupt backoff")
	}
}
