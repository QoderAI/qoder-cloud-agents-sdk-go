// Package ssestream incrementally decodes Qoder server-sent events.
package ssestream

import (
	"bufio"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apierror"
)

type Event struct {
	Type  string
	Data  []byte
	ID    string
	Retry int64
}
type Decoder interface {
	Event() Event
	Next() bool
	Close() error
	Err() error
}

type eventDecoder struct {
	body    io.ReadCloser
	scanner *bufio.Scanner
	event   Event
	err     error
	first   bool
	lastID  string
}

func NewDecoder(res *http.Response) Decoder {
	if res == nil || res.Body == nil {
		return nil
	}
	s := bufio.NewScanner(res.Body)
	s.Buffer(make([]byte, 4096), 32<<20)
	return &responseDecoder{eventDecoder: &eventDecoder{body: res.Body, scanner: s, first: true}, response: res}
}

type responseDecoder struct {
	*eventDecoder
	response *http.Response
}

type readError struct{ error }

func (e *readError) Unwrap() error { return e.error }
func (*readError) Retryable() bool { return true }

type decodeError struct{ error }

func (e *decodeError) Unwrap() error { return e.error }
func (*decodeError) Retryable() bool { return false }

func (d *eventDecoder) Event() Event { return d.event }
func (d *eventDecoder) Err() error   { return d.err }
func (d *eventDecoder) Close() error { return d.body.Close() }
func (d *eventDecoder) Next() bool {
	if d.err != nil {
		return false
	}
	ev := Event{Type: "message", ID: d.lastID}
	var data []string
	dispatch := func() bool {
		if len(data) == 0 {
			return false
		}
		ev.Data = []byte(strings.Join(data, "\n"))
		d.event = ev
		return true
	}
	for d.scanner.Scan() {
		line := d.scanner.Text()
		if d.first {
			line = strings.TrimPrefix(line, "\ufeff")
			d.first = false
		}
		if line == "" {
			if dispatch() {
				return true
			}
			ev = Event{Type: "message", ID: d.lastID}
			continue
		}
		if strings.HasPrefix(line, ":") {
			continue
		}
		key, value, _ := strings.Cut(line, ":")
		value = strings.TrimPrefix(value, " ")
		switch key {
		case "event":
			ev.Type = value
		case "data":
			data = append(data, value)
		case "id":
			if !strings.ContainsRune(value, '\x00') {
				d.lastID = value
				ev.ID = value
			}
		case "retry":
			if n, e := strconv.ParseInt(value, 10, 64); e == nil && n >= 0 {
				ev.Retry = n
			}
		}
	}
	d.err = d.scanner.Err()
	return false
}

type Stream[T any] struct {
	decoder Decoder
	cur     T
	err     error
	lastID  string
	closed  atomic.Bool
	done    bool
}

func NewStream[T any](d Decoder, err error) *Stream[T] { return &Stream[T]{decoder: d, err: err} }
func (s *Stream[T]) Next() bool {
	if s.closed.Load() || s.done || s.err != nil || s.decoder == nil {
		return false
	}
	for s.decoder.Next() {
		ev := s.decoder.Event()
		s.lastID = ev.ID
		if ev.Type == "ping" || len(ev.Data) == 0 {
			continue
		}
		if string(ev.Data) == "[DONE]" {
			s.done = true
			return false
		}
		if ev.Type == "error" {
			e := &apierror.Error{}
			if d, ok := s.decoder.(*responseDecoder); ok {
				e.Request = d.response.Request
				e.Response = d.response
				e.StatusCode = d.response.StatusCode
				e.RequestID = d.response.Header.Get("x-request-id")
			}
			_ = e.UnmarshalJSON(ev.Data)
			s.err = e
			return false
		}
		var value T
		if err := json.Unmarshal(ev.Data, &value); err != nil {
			s.err = &decodeError{error: err}
			return false
		}
		s.cur = value
		return true
	}
	if err := s.decoder.Err(); err != nil {
		s.err = &readError{error: err}
	}
	s.done = true
	return false
}
func (s *Stream[T]) Current() T { return s.cur }
func (s *Stream[T]) Err() error { return s.err }

// LastEventID can be passed as a Last-Event-ID header when the caller reconnects.
func (s *Stream[T]) LastEventID() string { return s.lastID }
func (s *Stream[T]) Close() error {
	if s.closed.Swap(true) || s.decoder == nil {
		return nil
	}
	return s.decoder.Close()
}
