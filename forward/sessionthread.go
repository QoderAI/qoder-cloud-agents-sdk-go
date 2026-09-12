package forward

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apiquery"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/pagination"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
)

// SessionThreadService provides Forward SessionThread operations.
type SessionThreadService struct {
	Options []option.RequestOption
	Events  SessionThreadEventService
}

func NewSessionThreadService(opts ...option.RequestOption) SessionThreadService {
	return SessionThreadService{Options: slices.Clone(opts), Events: NewSessionThreadEventService(opts...)}
}

// List Session Threads
func (r *SessionThreadService) List(ctx context.Context, sessionID string, params SessionThreadListParams, opts ...option.RequestOption) (res *pagination.Page[SessionThread], err error) {
	if sessionID == "" {
		return nil, fmt.Errorf("missing required session_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("sessions/%s/threads", url.PathEscape(sessionID))
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
func (r *SessionThreadService) ListAutoPaging(ctx context.Context, sessionID string, params SessionThreadListParams, opts ...option.RequestOption) *pagination.PageAutoPager[SessionThread] {
	return pagination.NewPageAutoPager(r.List(ctx, sessionID, params, opts...))
}

type SessionThreadListParams struct {
	// Page size, 1–100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Return records after this Thread ID.
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// Return records before this Thread ID.
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	paramObj
}

func (r SessionThreadListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// Get Session Thread
func (r *SessionThreadService) Get(ctx context.Context, sessionID string, threadID string, opts ...option.RequestOption) (res *SessionThread, err error) {
	if sessionID == "" {
		return nil, fmt.Errorf("missing required session_id parameter")
	}
	if threadID == "" {
		return nil, fmt.Errorf("missing required thread_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("sessions/%s/threads/%s", url.PathEscape(sessionID), url.PathEscape(threadID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Archive Session Thread
func (r *SessionThreadService) Archive(ctx context.Context, sessionID string, threadID string, params SessionThreadArchiveParams, opts ...option.RequestOption) (res *SessionThread, err error) {
	if sessionID == "" {
		return nil, fmt.Errorf("missing required session_id parameter")
	}
	if threadID == "" {
		return nil, fmt.Errorf("missing required thread_id parameter")
	}
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("sessions/%s/threads/%s/archive", url.PathEscape(sessionID), url.PathEscape(threadID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type SessionThreadArchiveParams struct {
	// Identifies one logical archive attempt; generate a unique value for each new
	// attempt.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r SessionThreadArchiveParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

type SessionThread struct {
	ID                 string                  `json:"id"`
	Type               string                  `json:"type"`
	SessionID          string                  `json:"session_id"`
	TemplateID         string                  `json:"template_id"`
	Role               string                  `json:"role"`
	Status             string                  `json:"status"`
	StopReason         SessionThreadStopReason `json:"stop_reason"`
	CreatedAt          time.Time               `json:"created_at" format:"date-time"`
	UpdatedAt          time.Time               `json:"updated_at" format:"date-time"`
	ParentThreadID     string                  `json:"parent_thread_id"`
	Name               string                  `json:"name"`
	CreatedByToolUseID string                  `json:"created_by_tool_use_id"`
	ArchivedAt         time.Time               `json:"archived_at" format:"date-time"`
	JSON               struct {
		ID                 respjson.Field
		Type               respjson.Field
		SessionID          respjson.Field
		TemplateID         respjson.Field
		Role               respjson.Field
		Status             respjson.Field
		StopReason         respjson.Field
		CreatedAt          respjson.Field
		UpdatedAt          respjson.Field
		ParentThreadID     respjson.Field
		Name               respjson.Field
		CreatedByToolUseID respjson.Field
		ArchivedAt         respjson.Field
		ExtraFields        map[string]respjson.Field
		raw                string
	} `json:"-"`
}

func (r SessionThread) RawJSON() string                  { return r.JSON.raw }
func (r *SessionThread) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type SessionThreadStopReason struct {
	Type string `json:"type"`
	JSON struct {
		Type        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r SessionThreadStopReason) RawJSON() string { return r.JSON.raw }
func (r *SessionThreadStopReason) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
