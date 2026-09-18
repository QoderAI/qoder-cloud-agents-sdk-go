// Package resumablestream reconnects resource-level SSE streams after retryable failures.
package resumablestream

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apierror"
)

const (
	initialBackoff    = 500 * time.Millisecond
	maximumBackoff    = 10 * time.Second
	healthyResetAfter = 5 * time.Second
)

// Child is one connection in a resumable stream.
type Child[T any] interface {
	Next() bool
	Current() T
	Err() error
	LastEventID() string
	Close() error
}

// Open starts a child stream at the supplied event cursor. resume is true once
// the cursor came from a yielded event rather than the initial request.
type Open[T any] func(context.Context, string, bool) Child[T]

type settings struct {
	initial      time.Duration
	maximum      time.Duration
	healthyAfter time.Duration
	jitter       func(time.Duration) time.Duration
	wait         func(context.Context, time.Duration) error
	now          func() time.Time
}

func defaultSettings() settings {
	return settings{
		initial:      initialBackoff,
		maximum:      maximumBackoff,
		healthyAfter: healthyResetAfter,
		jitter: func(backoff time.Duration) time.Duration {
			half := backoff / 2
			return half + time.Duration(rand.Int63n(int64(backoff-half)))
		},
		wait: func(ctx context.Context, delay time.Duration) error {
			timer := time.NewTimer(delay)
			defer timer.Stop()
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-timer.C:
				return nil
			}
		},
		now: time.Now,
	}
}

// Stream reconnects child streams while retaining the latest yielded event cursor.
type Stream[T any] struct {
	parent   context.Context
	ctx      context.Context
	cancel   context.CancelFunc
	open     Open[T]
	terminal func(T) bool
	config   settings

	mu            sync.Mutex
	child         Child[T]
	current       T
	lastEventID   string
	hasCheckpoint bool
	connectedAt   time.Time
	err           error
	closed        bool
	done          bool
	backoff       time.Duration
}

// New creates a resumable stream. Connections are opened lazily by Next.
func New[T any](ctx context.Context, lastEventID string, open Open[T], terminal func(T) bool) *Stream[T] {
	return newStream(ctx, lastEventID, open, terminal, defaultSettings())
}

// NewError creates a stream that has already failed.
func NewError[T any](ctx context.Context, err error) *Stream[T] {
	return newStream[T](ctx, "", nil, nil, defaultSettings(), err)
}

func newStream[T any](ctx context.Context, lastEventID string, open Open[T], terminal func(T) bool, config settings, initialErr ...error) *Stream[T] {
	if ctx == nil {
		ctx = context.Background()
	}
	streamCtx, cancel := context.WithCancel(ctx)
	s := &Stream[T]{
		parent:      ctx,
		ctx:         streamCtx,
		cancel:      cancel,
		open:        open,
		terminal:    terminal,
		config:      config,
		lastEventID: lastEventID,
		backoff:     config.initial,
	}
	if len(initialErr) != 0 {
		s.err = initialErr[0]
		cancel()
	} else {
		go s.closeOnContext()
	}
	return s
}

// Next advances to the next decoded event, reconnecting as needed.
func (s *Stream[T]) Next() bool {
	for {
		if s.stopped() {
			return false
		}

		child := s.activeChild()
		if child == nil {
			if s.open == nil {
				s.fail(errors.New("resumable stream opener is nil"))
				return false
			}
			cursor, resume := s.cursor()
			child = s.open(s.ctx, cursor, resume)
			if child == nil {
				s.fail(errors.New("resumable stream opener returned nil"))
				return false
			}
			if !s.installChild(child) {
				_ = child.Close()
				return false
			}
		}

		if child.Next() {
			accepted, terminal := s.checkpoint(child.Current(), child.LastEventID())
			if terminal {
				s.closeChild(child)
				s.cancel()
			}
			return accepted
		}

		childErr := child.Err()
		s.closeChild(child)
		if s.stopped() {
			return false
		}
		if childErr != nil && !retryable(s.ctx, childErr) {
			s.fail(childErr)
			return false
		}
		if err := s.config.wait(s.ctx, s.config.jitter(s.nextBackoff())); err != nil {
			if !s.isClosed() {
				s.fail(s.contextError(err))
			}
			return false
		}
	}
}

// Current returns the most recently yielded event.
func (s *Stream[T]) Current() T {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.current
}

// Err returns the terminal non-retryable or context error, if any.
func (s *Stream[T]) Err() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.err
}

// LastEventID returns the latest non-empty cursor from a yielded child event.
func (s *Stream[T]) LastEventID() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lastEventID
}

// Close cancels pending work and closes the active child stream.
func (s *Stream[T]) Close() error {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return nil
	}
	s.closed = true
	child := s.child
	s.child = nil
	s.mu.Unlock()

	s.cancel()
	if child != nil {
		return child.Close()
	}
	return nil
}

func (s *Stream[T]) closeOnContext() {
	<-s.ctx.Done()
	s.mu.Lock()
	if err := s.parent.Err(); err != nil && s.err == nil && !s.closed {
		s.err = err
	}
	child := s.child
	s.child = nil
	s.mu.Unlock()
	if child != nil {
		_ = child.Close()
	}
}

func (s *Stream[T]) stopped() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed || s.done || s.err != nil {
		return true
	}
	if err := s.parent.Err(); err != nil {
		s.err = err
		return true
	}
	return false
}

func (s *Stream[T]) cursor() (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lastEventID, s.hasCheckpoint
}

func (s *Stream[T]) activeChild() Child[T] {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.child
}

func (s *Stream[T]) installChild(child Child[T]) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed || s.done || s.err != nil || s.parent.Err() != nil {
		if s.err == nil && !s.closed && !s.done {
			s.err = s.parent.Err()
		}
		return false
	}
	s.child = child
	s.connectedAt = s.config.now()
	return true
}

func (s *Stream[T]) checkpoint(current T, lastEventID string) (bool, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed || s.done || s.err != nil {
		return false, false
	}
	if err := s.parent.Err(); err != nil {
		s.err = err
		return false, false
	}
	s.current = current
	if lastEventID != "" {
		s.lastEventID = lastEventID
		s.hasCheckpoint = true
	}
	if s.config.now().Sub(s.connectedAt) > s.config.healthyAfter {
		s.backoff = s.config.initial
	}
	s.done = s.terminal != nil && s.terminal(current)
	return true, s.done
}

func (s *Stream[T]) closeChild(child Child[T]) {
	s.mu.Lock()
	if s.child != child {
		s.mu.Unlock()
		return
	}
	s.child = nil
	s.mu.Unlock()
	_ = child.Close()
}

func (s *Stream[T]) nextBackoff() time.Duration {
	s.mu.Lock()
	defer s.mu.Unlock()
	backoff := s.backoff
	s.backoff = min(s.backoff*2, s.config.maximum)
	return backoff
}

func (s *Stream[T]) fail(err error) {
	s.mu.Lock()
	failed := s.err == nil && !s.closed
	if failed {
		s.err = err
	}
	s.mu.Unlock()
	if failed {
		s.cancel()
	}
}

func (s *Stream[T]) isClosed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.closed
}

func (s *Stream[T]) contextError(fallback error) error {
	if err := s.parent.Err(); err != nil {
		return err
	}
	return fallback
}

func retryable(_ context.Context, err error) bool {
	if err == nil {
		return true
	}
	if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
		return false
	}

	var apiErr *apierror.Error
	if errors.As(err, &apiErr) {
		if apiErr.StatusCode == http.StatusConflict {
			return false
		}
		if apiErr.Response != nil {
			switch apiErr.Response.Header.Get("x-should-retry") {
			case "true":
				return true
			case "false":
				return false
			}
		}
		return apiErr.StatusCode == http.StatusRequestTimeout ||
			apiErr.StatusCode == http.StatusTooManyRequests ||
			apiErr.StatusCode >= http.StatusInternalServerError
	}

	var classified interface{ Retryable() bool }
	if errors.As(err, &classified) {
		return classified.Retryable()
	}
	var noCredentials *convention.NoCredentialsError
	if errors.As(err, &noCredentials) {
		return false
	}
	var syntaxErr *json.SyntaxError
	if errors.As(err, &syntaxErr) {
		return false
	}
	var typeErr *json.UnmarshalTypeError
	if errors.As(err, &typeErr) {
		return false
	}
	var urlErr *url.Error
	if errors.As(err, &urlErr) {
		return true
	}
	var netErr net.Error
	return errors.As(err, &netErr) || errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.ErrClosedPipe)
}
