package forward

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apiquery"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/pagination"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/ssestream"
)

// SessionThreadEventService provides Forward SessionThreadEvent operations.
type SessionThreadEventService struct {
	Options []option.RequestOption
}

func NewSessionThreadEventService(opts ...option.RequestOption) SessionThreadEventService {
	return SessionThreadEventService{Options: slices.Clone(opts)}
}

// 列出 Session Thread Events.
func (r *SessionThreadEventService) List(ctx context.Context, sessionID string, threadID string, params SessionThreadEventListParams, opts ...option.RequestOption) (res *pagination.Page[SessionEvent], err error) {
	if sessionID == "" {
		return nil, fmt.Errorf("missing required session_id parameter")
	}
	if threadID == "" {
		return nil, fmt.Errorf("missing required thread_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("sessions/%s/threads/%s/events", url.PathEscape(sessionID), url.PathEscape(threadID))
	var raw *http.Response
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	cfg, err := convention.NewRequestConfig(ctx, http.MethodGet, path, params, &res, opts...)
	if err != nil {
		return nil, err
	}
	if err = cfg.Execute(); err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}
func (r *SessionThreadEventService) ListAutoPaging(ctx context.Context, sessionID string, threadID string, params SessionThreadEventListParams, opts ...option.RequestOption) *pagination.PageAutoPager[SessionEvent] {
	return pagination.NewPageAutoPager(r.List(ctx, sessionID, threadID, params, opts...))
}

type SessionThreadEventListParams struct {
	// 分页大小，范围为 1–100。
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// 返回该 Event ID 之后的记录。
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// 返回该 Event ID 之前的记录。
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	paramObj
}

func (r SessionThreadEventListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// 订阅 Session Thread Event Stream.
func (r *SessionThreadEventService) StreamEvents(ctx context.Context, sessionID string, threadID string, params SessionThreadEventStreamParams, opts ...option.RequestOption) *ssestream.Stream[SessionEvent] {
	if sessionID == "" {
		return ssestream.NewStream[SessionEvent](nil, fmt.Errorf("missing required session_id parameter"))
	}
	if threadID == "" {
		return ssestream.NewStream[SessionEvent](nil, fmt.Errorf("missing required thread_id parameter"))
	}
	if params.LastEventID.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Last-Event-ID", fmt.Sprint(params.LastEventID.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, []option.RequestOption{option.WithHeader("Accept", "text/event-stream")}, opts)
	path := fmt.Sprintf("sessions/%s/threads/%s/stream", url.PathEscape(sessionID), url.PathEscape(threadID))
	var raw *http.Response
	err := convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &raw, opts...)
	return ssestream.NewStream[SessionEvent](ssestream.NewDecoder(raw), err)
}

type SessionThreadEventStreamParams struct {
	// 从该 Thread Event 之后继续订阅。
	LastEventID param.Opt[string] `header:"Last-Event-ID,omitzero" json:"-"`
	paramObj
}

func (r SessionThreadEventStreamParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}
