package ssestream_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/ssestream"
	"io"
	"net/http"
	"testing"
)

func TestSSEErrorsAndClose(t *testing.T) {
	s := ssestream.NewStream[json.RawMessage](ssestream.NewDecoder(&http.Response{StatusCode: 200, Header: http.Header{}, Body: io.NopCloser(bytes.NewBufferString("event: error\ndata: {\"type\":\"error\",\"request_id\":\"req\",\"error\":{\"type\":\"api_error\",\"message\":\"failed\"}}\n\n"))}), nil)
	if s.Next() {
		t.Fatal("error became data")
	}
	var e *convention.Error
	if !errors.As(s.Err(), &e) || e.RequestID != "req" {
		t.Fatal(s.Err())
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	if err := s.Close(); err != nil {
		t.Fatal(err)
	}
	if s.Next() {
		t.Fatal("closed stream advanced")
	}
}

func TestSSEFraming(t *testing.T) {
	input := "\ufeff:heartbeat\r\nid: first\r\nevent: delta\r\nretry: 1200\r\ndata: {\r\ndata: \"text\":\"hello\"}\r\n\r\nid: bad\x00id\ndata: {}\n\nid:\ndata: {}\n\n"
	d := ssestream.NewDecoder(&http.Response{Body: io.NopCloser(bytes.NewBufferString(input))})
	defer d.Close()
	for i, want := range []struct {
		id, typ, data string
		retry         int64
	}{{"first", "delta", "{\n\"text\":\"hello\"}", 1200}, {"first", "message", "{}", 0}, {"", "message", "{}", 0}} {
		if !d.Next() {
			t.Fatalf("event %d: %v", i, d.Err())
		}
		e := d.Event()
		if e.ID != want.id || e.Type != want.typ || string(e.Data) != want.data || e.Retry != want.retry {
			t.Fatalf("event %d: %#v", i, e)
		}
	}
	if d.Next() || d.Err() != nil {
		t.Fatal("unexpected final event", d.Err())
	}
}

func TestSSEDecoderDiscardsIncompleteEOFFrame(t *testing.T) {
	input := "id: complete\ndata: {\"id\":\"complete\"}\n\nid: incomplete\ndata: {\"id\":\"incomplete\"}"
	d := ssestream.NewDecoder(&http.Response{Body: io.NopCloser(bytes.NewBufferString(input))})
	defer d.Close()

	if !d.Next() || d.Event().ID != "complete" {
		t.Fatalf("complete event missing: event=%#v err=%v", d.Event(), d.Err())
	}
	if d.Next() || d.Err() != nil {
		t.Fatalf("incomplete EOF frame was dispatched: event=%#v err=%v", d.Event(), d.Err())
	}
}

func TestStreamTerminalState(t *testing.T) {
	for _, terminal := range []string{"data: [DONE]\n\n", "data: {broken\n\n", "event: error\ndata: {\"error\":{\"message\":\"failed\"}}\n\n"} {
		t.Run(terminal, func(t *testing.T) {
			body := &countingBody{Reader: bytes.NewBufferString("event: ping\ndata: {}\n\nid: one\ndata: {\"text\":\"hello\"}\n\n" + terminal + "data: {\"text\":\"must not resume\"}\n\n")}
			s := ssestream.NewStream[map[string]any](ssestream.NewDecoder(&http.Response{Header: http.Header{}, Body: body}), nil)
			if !s.Next() || s.Current()["text"] != "hello" || s.LastEventID() != "one" {
				t.Fatal("first event lost")
			}
			if s.Next() || s.Next() {
				t.Fatal("terminal stream resumed")
			}
			_ = s.Close()
			_ = s.Close()
			if body.closes != 1 {
				t.Fatalf("closes=%d", body.closes)
			}
		})
	}
}

type countingBody struct {
	io.Reader
	closes int
}

func (b *countingBody) Close() error { b.closes++; return nil }
func TestStreamReadFailure(t *testing.T) {
	cause := errors.New("connection interrupted")
	s := ssestream.NewStream[json.RawMessage](ssestream.NewDecoder(&http.Response{Body: io.NopCloser(failingReader{cause})}), nil)
	defer s.Close()
	if s.Next() || !errors.Is(s.Err(), cause) {
		t.Fatal(s.Err())
	}
}

type failingReader struct{ err error }

func (r failingReader) Read([]byte) (int, error) { return 0, r.err }
