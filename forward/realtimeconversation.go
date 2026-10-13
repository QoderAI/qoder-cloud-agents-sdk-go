package forward

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apiquery"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
)

// RealtimeConversationService creates Realtime Conversations and queries their
// persisted history. Realtime is Beta; a Template must already have it enabled.
type RealtimeConversationService struct{ Options []option.RequestOption }

func NewRealtimeConversationService(opts ...option.RequestOption) RealtimeConversationService {
	return RealtimeConversationService{Options: slices.Clone(opts)}
}

// New creates a Conversation. Supply a stable Idempotency-Key for this operation
// with IdempotencyKey or option.WithHeader; it is required by the service.
func (r *RealtimeConversationService) New(ctx context.Context, params RealtimeConversationNewParams, opts ...option.RequestOption) (res *RealtimeConversation, err error) {
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", params.IdempotencyKey.Value)}, opts...)
	}
	cfg, err := convention.NewRequestConfig(ctx, http.MethodPost, "realtime/conversations", params, &res, slices.Concat(r.Options, opts)...)
	if err != nil {
		return nil, err
	}
	key := cfg.Request.Header.Get("Idempotency-Key")
	if len(key) < 1 || len(key) > 256 || !utf8.ValidString(key) || strings.TrimSpace(key) != key || strings.ContainsFunc(key, unicode.IsControl) {
		return nil, fmt.Errorf("Realtime creation requires Idempotency-Key of 1–256 UTF-8 bytes without surrounding whitespace or controls")
	}
	err = cfg.Execute()
	return res, err
}

type RealtimeConversationNewParams struct {
	IdentityID     string              `json:"identity_id" api:"required"`
	TemplateID     string              `json:"template_id" api:"required"`
	Title          param.Opt[string]   `json:"title,omitzero"`
	Metadata       map[string]any      `json:"metadata,omitzero"`
	Config         RealtimeConfigParam `json:"config,omitzero"`
	IdempotencyKey param.Opt[string]   `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r RealtimeConversationNewParams) MarshalJSON() ([]byte, error) {
	type shadow RealtimeConversationNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}

type RealtimeConfigParam struct {
	Audio RealtimeAudioConfigParam `json:"audio,omitzero"`
	paramObj
}

func (r RealtimeConfigParam) MarshalJSON() ([]byte, error) {
	type shadow RealtimeConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type RealtimeAudioConfigParam struct {
	Output RealtimeAudioOutputConfigParam `json:"output,omitzero"`
	paramObj
}

func (r RealtimeAudioConfigParam) MarshalJSON() ([]byte, error) {
	type shadow RealtimeAudioConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}

type RealtimeAudioOutputConfigParam struct {
	Voice param.Opt[string] `json:"voice,omitzero"`
	paramObj
}

func (r RealtimeAudioOutputConfigParam) MarshalJSON() ([]byte, error) {
	type shadow RealtimeAudioOutputConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}

// RealtimeConversation timestamps remain text to preserve the wire precision.
type RealtimeConversation struct {
	ID        string         `json:"id"`
	Type      string         `json:"type"`
	Status    string         `json:"status"`
	Title     *string        `json:"title"`
	Metadata  map[string]any `json:"metadata"`
	Config    RealtimeConfig `json:"config"`
	CreatedAt string         `json:"created_at"`
	UpdatedAt string         `json:"updated_at"`
}
type RealtimeConfig struct {
	Audio RealtimeAudioConfig `json:"audio"`
}
type RealtimeAudioConfig struct {
	Output RealtimeAudioOutputConfig `json:"output"`
}
type RealtimeAudioOutputConfig struct {
	Voice string `json:"voice"`
}

type RealtimeConversationHistoryParams struct {
	Limit  param.Opt[int64]  `query:"limit,omitzero" json:"-"`
	Before param.Opt[string] `query:"before,omitzero" json:"-"`
	// "message", "work", or "message,work"; the default includes both.
	Types param.Opt[string] `query:"types,omitzero" json:"-"`
	paramObj
}

func (r RealtimeConversationHistoryParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{})
}

// History fetches a page of persisted events. Its cursor points towards older
// history, and is unrelated to WebSocket event IDs or sequence numbers.
func (r *RealtimeConversationService) History(ctx context.Context, conversationID string, params RealtimeConversationHistoryParams, opts ...option.RequestOption) (res *RealtimeHistoryPage, err error) {
	if conversationID == "" {
		return nil, fmt.Errorf("missing required conversation_id parameter")
	}
	if params.Limit.Valid() && (params.Limit.Value < 1 || params.Limit.Value > 100) {
		return nil, fmt.Errorf("history limit must be between 1 and 100")
	}
	if params.Before.Valid() && params.Before.Value == "" {
		return nil, fmt.Errorf("history before cursor must not be empty")
	}
	if params.Types.Valid() && params.Types.Value != "message" && params.Types.Value != "work" && params.Types.Value != "message,work" {
		return nil, fmt.Errorf("history types must be message, work, or message,work")
	}
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, "realtime/conversations/"+url.PathEscape(conversationID)+"/history", params, &res, slices.Concat(r.Options, opts)...)
	if err == nil && res != nil {
		if res.Page.HasMore && (res.Page.NextBefore == nil || *res.Page.NextBefore == "") {
			return nil, &RealtimeProtocolError{Message: "history page has_more without next_before cursor"}
		}
		res.service = r
		res.params = params
		res.options = slices.Clone(opts)
		res.conversationID = conversationID
	}
	return res, err
}

type RealtimeHistoryConversation struct {
	ID                   string         `json:"id"`
	Title                *string        `json:"title"`
	InitializationStatus string         `json:"initialization_status"`
	Metadata             map[string]any `json:"metadata"`
	CreatedAt            string         `json:"created_at"`
	UpdatedAt            string         `json:"updated_at"`
}

// RealtimeHistoryEvent retains unknown fields through RawJSON. ID is a history
// ID (vevt_*), not an ephemeral WebSocket event_id.
type RealtimeHistoryEvent struct {
	ID                 string             `json:"id"`
	Type               string             `json:"type"`
	Status             string             `json:"status"`
	OccurredAt         string             `json:"occurred_at"`
	Role               string             `json:"role,omitempty"`
	Text               string             `json:"text,omitempty"`
	WorkID             string             `json:"work_id,omitempty"`
	Objective          string             `json:"objective,omitempty"`
	Result             string             `json:"result,omitempty"`
	Error              *RealtimeWorkError `json:"error,omitempty"`
	TurnID             string             `json:"turn_id,omitempty"`
	UserMessageEventID string             `json:"user_message_event_id,omitempty"`
	AnnouncementID     string             `json:"announcement_id,omitempty"`
	MilestoneType      string             `json:"milestone_type,omitempty"`
	raw                string
}

func (e *RealtimeHistoryEvent) UnmarshalJSON(b []byte) error {
	type shadow RealtimeHistoryEvent
	if err := json.Unmarshal(b, (*shadow)(e)); err != nil {
		return err
	}
	e.raw = string(b)
	return nil
}
func (e RealtimeHistoryEvent) RawJSON() string { return e.raw }

type RealtimeHistoryPageInfo struct {
	NextBefore *string `json:"next_before"`
	HasMore    bool    `json:"has_more"`
}
type RealtimeHistoryPage struct {
	Conversation   RealtimeHistoryConversation `json:"conversation"`
	Events         []RealtimeHistoryEvent      `json:"events"`
	Page           RealtimeHistoryPageInfo     `json:"page"`
	service        *RealtimeConversationService
	params         RealtimeConversationHistoryParams
	options        []option.RequestOption
	conversationID string
}

func (p *RealtimeHistoryPage) HasNextPage() bool {
	return p.Page.HasMore && p.Page.NextBefore != nil && *p.Page.NextBefore != ""
}

// GetNextPage returns older history, or nil at the end.
func (p *RealtimeHistoryPage) GetNextPage(ctx context.Context) (*RealtimeHistoryPage, error) {
	if !p.HasNextPage() {
		return nil, nil
	}
	if p.service == nil {
		return nil, fmt.Errorf("history page has no request configuration")
	}
	if p.params.Before.Valid() && p.params.Before.Value == *p.Page.NextBefore {
		return nil, fmt.Errorf("history cursor did not advance")
	}
	params := p.params
	params.Before = param.NewOpt(*p.Page.NextBefore)
	return p.service.History(ctx, p.conversationID, params, p.options...)
}

// HistoryAutoPaging iterates persisted history towards older records.
func (r *RealtimeConversationService) HistoryAutoPaging(ctx context.Context, conversationID string, params RealtimeConversationHistoryParams, opts ...option.RequestOption) *RealtimeHistoryAutoPager {
	page, err := r.History(ctx, conversationID, params, opts...)
	return &RealtimeHistoryAutoPager{ctx: ctx, page: page, err: err}
}

type RealtimeHistoryAutoPager struct {
	ctx     context.Context
	page    *RealtimeHistoryPage
	index   int
	current RealtimeHistoryEvent
	err     error
	seen    map[string]bool
}

func (p *RealtimeHistoryAutoPager) Next() bool {
	for p.err == nil && p.page != nil {
		if p.index < len(p.page.Events) {
			p.current = p.page.Events[p.index]
			p.index++
			return true
		}
		if !p.page.HasNextPage() {
			p.page = nil
			return false
		}
		cursor := *p.page.Page.NextBefore
		if p.seen == nil {
			p.seen = map[string]bool{}
		}
		if p.seen[cursor] {
			p.err = fmt.Errorf("history cursor cycle detected")
			return false
		}
		p.seen[cursor] = true
		p.page, p.err = p.page.GetNextPage(p.ctx)
		p.index = 0
	}
	return false
}
func (p *RealtimeHistoryAutoPager) Current() RealtimeHistoryEvent { return p.current }
func (p *RealtimeHistoryAutoPager) Err() error                    { return p.err }
