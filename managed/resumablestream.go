package managed

import (
	"context"
	"errors"
	"slices"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/resumablestream"
)

// NewResumableStream subscribes to Session Events and reconnects retryable failures
// with the latest yielded Last-Event-ID cursor.
func (r *SessionEventService) NewResumableStream(ctx context.Context, sessionID string, params SessionEventStreamParams, opts ...option.RequestOption) *ResumableSessionEventStream {
	if sessionID == "" {
		return &ResumableSessionEventStream{stream: resumablestream.NewError[ManagedAgentsStreamSessionEventsUnion](ctx, errors.New("missing required session_id parameter"))}
	}
	params.EventDeltas = slices.Clone(params.EventDeltas)
	params.Betas = slices.Clone(params.Betas)
	opts = slices.Clone(opts)
	service := NewSessionEventService(slices.Clone(r.Options)...)
	lastEventID := ""
	if params.LastEventID.Valid() {
		lastEventID = params.LastEventID.Value
	}
	stream := resumablestream.New(ctx, lastEventID, func(ctx context.Context, cursor string, resume bool) resumablestream.Child[ManagedAgentsStreamSessionEventsUnion] {
		requestOpts := slices.Clone(opts)
		if resume {
			requestOpts = append(requestOpts, option.WithHeader("Last-Event-ID", cursor))
		}
		return service.StreamEvents(ctx, sessionID, params, requestOpts...)
	}, func(event ManagedAgentsStreamSessionEventsUnion) bool {
		return event.Type == "session.status_terminated" || event.Type == "session.deleted"
	})
	return &ResumableSessionEventStream{stream: stream}
}

// ResumableSessionEventStream is a Session Event stream that reconnects after
// retryable HTTP, transport, read, and unexpected EOF failures.
type ResumableSessionEventStream struct {
	stream *resumablestream.Stream[ManagedAgentsStreamSessionEventsUnion]
}

func (s *ResumableSessionEventStream) Next() bool { return s.stream.Next() }
func (s *ResumableSessionEventStream) Current() ManagedAgentsStreamSessionEventsUnion {
	return s.stream.Current()
}
func (s *ResumableSessionEventStream) Err() error          { return s.stream.Err() }
func (s *ResumableSessionEventStream) LastEventID() string { return s.stream.LastEventID() }
func (s *ResumableSessionEventStream) Close() error        { return s.stream.Close() }
