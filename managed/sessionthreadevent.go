// Qoder managed API definitions.
package managed

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	requestconfig "github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apiquery"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/pagination"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/ssestream"
)

// SessionThreadEventService contains methods and other services that help with
// interacting with the Qoder Cloud Agents API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSessionThreadEventService] method instead.
type SessionThreadEventService struct {
	Options []option.RequestOption
}

// NewSessionThreadEventService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewSessionThreadEventService(opts ...option.RequestOption) (r SessionThreadEventService) {
	r = SessionThreadEventService{}
	r.Options = opts
	return
}

// List Session Thread Events
func (r *SessionThreadEventService) List(ctx context.Context, threadID string, params SessionThreadEventListParams, opts ...option.RequestOption) (res *pagination.PageCursor[ManagedAgentsSessionEventUnion], err error) {
	var raw *http.Response
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if params.SessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("sessions/%s/threads/%s/events", url.PathEscape(params.SessionID), url.PathEscape(threadID))
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, params, &res, opts...)
	if err != nil {
		return nil, err
	}
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

// List Session Thread Events
func (r *SessionThreadEventService) ListAutoPaging(ctx context.Context, threadID string, params SessionThreadEventListParams, opts ...option.RequestOption) *pagination.PageCursorAutoPager[ManagedAgentsSessionEventUnion] {
	return pagination.NewPageCursorAutoPager(r.List(ctx, threadID, params, opts...))
}

// Stream Session Thread Events
func (r *SessionThreadEventService) StreamEvents(ctx context.Context, threadID string, params SessionThreadEventStreamParams, opts ...option.RequestOption) (stream *ssestream.Stream[ManagedAgentsStreamSessionThreadEventsUnion]) {
	var (
		raw *http.Response
		err error
	)
	for _, v := range params.Betas {
		opts = append(opts, option.WithHeaderAdd("x-qoder-beta", fmt.Sprintf("%v", v)))
	}
	opts = slices.Concat(r.Options, []option.RequestOption{option.WithHeader("Accept", "text/event-stream")}, opts)

	if params.SessionID == "" {
		err = errors.New("missing required session_id parameter")
		return ssestream.NewStream[ManagedAgentsStreamSessionThreadEventsUnion](nil, err)
	}
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return ssestream.NewStream[ManagedAgentsStreamSessionThreadEventsUnion](nil, err)
	}
	path := fmt.Sprintf("sessions/%s/threads/%s/stream", url.PathEscape(params.SessionID), url.PathEscape(threadID))
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &raw, opts...)
	return ssestream.NewStream[ManagedAgentsStreamSessionThreadEventsUnion](ssestream.NewDecoder(raw), err)
}

type SessionThreadEventListParams struct {
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	AfterID  param.Opt[string] `query:"after_id,omitzero" json:"-"`

	SessionID string `path:"session_id" api:"required" json:"-"`
	// Query parameter for limit
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Query parameter for page
	Page        param.Opt[string] `query:"page,omitzero" json:"-"`
	WorkspaceID param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [SessionThreadEventListParams]'s query parameters as
// `url.Values`.
func (r SessionThreadEventListParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SessionThreadEventStreamParams struct {
	SessionID   string            `path:"session_id" api:"required" json:"-"`
	WorkspaceID param.Opt[string] `header:"qoder-workspace-id,omitzero" json:"-"`
	// When set, this connection also receives streaming deltas (`event_start`,
	// `event_delta`) while an event is being produced, before the event itself
	// arrives. Deltas are best-effort; when the final event is produced it carries the
	// complete content. A model request that ends early (an error or interrupt)
	// produces no final event — its terminal `span.model_request_end` closes the
	// preview. Accepts one or more event types to preview and may be repeated:
	// `agent.message` streams `content_delta` fragments; `agent.thinking` is
	// start-only — a signal that the agent has begun extended thinking, concluded by
	// the `agent.thinking` event itself. Only previews of the requested event types
	// are sent.
	EventDeltas []ManagedAgentsDeltaType `query:"event_deltas,omitzero" json:"-"`
	// Optional header to specify the beta version(s) you want to use.
	Betas []QoderBeta `header:"x-qoder-beta,omitzero" json:"-"`
	paramObj
}

// URLQuery serializes [SessionThreadEventStreamParams]'s query parameters as
// `url.Values`.
func (r SessionThreadEventStreamParams) URLQuery() (v url.Values, err error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatBrackets,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
