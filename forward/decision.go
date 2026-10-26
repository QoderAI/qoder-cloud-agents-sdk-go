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

// DecisionService evaluates explicit state or an accessible Session's recent history.
type DecisionService struct {
	Options []option.RequestOption
}

func NewDecisionService(opts ...option.RequestOption) DecisionService {
	return DecisionService{Options: slices.Clone(opts)}
}

// New synchronously evaluates questions and stores the resulting Decision.
// Provide exactly one of SessionID and State. Session context is selected and
// truncated by Forward; explicit State and Questions are never truncated.
func (r *DecisionService) New(ctx context.Context, params DecisionNewParams, opts ...option.RequestOption) (res *Decision, err error) {
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", params.IdempotencyKey.Value)}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, "decisions", params, &res, opts...)
	return res, err
}

// Get retrieves a Decision accessible to the current caller.
func (r *DecisionService) Get(ctx context.Context, decisionID string, opts ...option.RequestOption) (res *Decision, err error) {
	if decisionID == "" {
		return nil, fmt.Errorf("missing required decision_id parameter")
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("decisions/%s", url.PathEscape(decisionID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List retrieves Decisions, optionally restricted to an accessible Session.
func (r *DecisionService) List(ctx context.Context, params DecisionListParams, opts ...option.RequestOption) (res *pagination.Page[Decision], err error) {
	opts = slices.Concat(r.Options, opts)
	var raw *http.Response
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	cfg, err := convention.NewRequestConfig(ctx, http.MethodGet, "decisions", params, &res, opts...)
	if err != nil {
		return nil, err
	}
	if err = cfg.Execute(); err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

func (r *DecisionService) ListAutoPaging(ctx context.Context, params DecisionListParams, opts ...option.RequestOption) *pagination.PageAutoPager[Decision] {
	return pagination.NewPageAutoPager(r.List(ctx, params, opts...))
}

type DecisionNewParams struct {
	// An accessible Session ID. Mutually exclusive with State.
	SessionID param.Opt[string] `json:"session_id,omitzero"`
	// A nonempty JSON object to evaluate. Mutually exclusive with SessionID.
	// The top-level state.session_context key is reserved for server-built history.
	State map[string]any `json:"state,omitzero"`
	// One to 32 choice or score questions.
	Questions []DecisionQuestionParam `json:"questions" api:"required"`
	// Reusing a key with identical request bytes replays the stored result.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r DecisionNewParams) MarshalJSON() ([]byte, error) {
	type shadow DecisionNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionNewParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type DecisionQuestionParam struct {
	// Either "choice" or "score".
	Type string `json:"type" api:"required"`
	// Optional unique, nonempty question name, up to 128 bytes.
	Name param.Opt[string] `json:"name,omitzero"`
	// Nonempty instructions for evaluating this question.
	Instructions string `json:"instructions" api:"required"`
	// Two to 32 options for a choice question; omit for score questions.
	Choices []DecisionChoiceParam `json:"choices,omitzero"`
	// Two to 32 discrete levels for a score question; omit for choice questions.
	Levels []DecisionLevelParam `json:"levels,omitzero"`
	paramObj
}

func (r DecisionQuestionParam) MarshalJSON() ([]byte, error) {
	type shadow DecisionQuestionParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionQuestionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type DecisionChoiceParam struct {
	// A nonempty string or boolean. Values must be distinct within a question.
	Value       any               `json:"value" api:"required"`
	Description param.Opt[string] `json:"description,omitzero"`
	paramObj
}

func (r DecisionChoiceParam) MarshalJSON() ([]byte, error) {
	type shadow DecisionChoiceParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionChoiceParam) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type DecisionLevelParam struct {
	Label string `json:"label" api:"required"`
	// A finite numeric score, including zero. The result selects a configured level.
	Score       float64           `json:"score" api:"required"`
	Description param.Opt[string] `json:"description,omitzero"`
	paramObj
}

func (r DecisionLevelParam) MarshalJSON() ([]byte, error) {
	type shadow DecisionLevelParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *DecisionLevelParam) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type DecisionListParams struct {
	SessionID param.Opt[string] `query:"session_id,omitzero" json:"-"`
	// Page size from 1 to 100; defaults to 20.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Cursor for the next page. Mutually exclusive with BeforeID.
	AfterID  param.Opt[string] `query:"after_id,omitzero" json:"-"`
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	paramObj
}

func (r DecisionListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// Decision is a stored evaluation. GET and List may also return running or failed records.
type Decision struct {
	ID string `json:"id"`
	// Null for explicit-state input.
	SessionID string `json:"session_id" api:"nullable"`
	Source    string `json:"source"`
	// "running", "completed", or "failed".
	Status    string             `json:"status"`
	Questions []DecisionQuestion `json:"questions"`
	Answers   []DecisionAnswer   `json:"answers"`
	Context   DecisionContext    `json:"context"`
	// Reported model usage, if available. Context estimates are not usage.
	Usage       map[string]any `json:"usage" api:"nullable"`
	Error       DecisionError  `json:"error" api:"nullable"`
	CreatedAt   time.Time      `json:"created_at" format:"date-time"`
	CompletedAt time.Time      `json:"completed_at" format:"date-time" api:"nullable"`
	JSON        struct {
		ID, SessionID, Source, Status, Questions, Answers, Context, Usage, Error, CreatedAt, CompletedAt respjson.Field
		ExtraFields                                                                                      map[string]respjson.Field
		raw                                                                                              string
	} `json:"-"`
}

func (r Decision) RawJSON() string                  { return r.JSON.raw }
func (r *Decision) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type DecisionQuestion struct {
	Type         string           `json:"type"`
	Name         string           `json:"name"`
	Instructions string           `json:"instructions"`
	Choices      []DecisionChoice `json:"choices"`
	Levels       []DecisionLevel  `json:"levels"`
	JSON         struct {
		Type, Name, Instructions, Choices, Levels respjson.Field
		ExtraFields                               map[string]respjson.Field
		raw                                       string
	} `json:"-"`
}

func (r DecisionQuestion) RawJSON() string                  { return r.JSON.raw }
func (r *DecisionQuestion) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type DecisionChoice struct {
	// The original string or boolean value, including false.
	Value       any    `json:"value"`
	Description string `json:"description"`
	JSON        struct {
		Value, Description respjson.Field
		ExtraFields        map[string]respjson.Field
		raw                string
	} `json:"-"`
}

func (r DecisionChoice) RawJSON() string                  { return r.JSON.raw }
func (r *DecisionChoice) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type DecisionLevel struct {
	Label       string  `json:"label"`
	Score       float64 `json:"score"`
	Description string  `json:"description"`
	JSON        struct {
		Label, Score, Description respjson.Field
		ExtraFields               map[string]respjson.Field
		raw                       string
	} `json:"-"`
}

func (r DecisionLevel) RawJSON() string                  { return r.JSON.raw }
func (r *DecisionLevel) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type DecisionAnswer struct {
	// "choice", "score", or "refusal".
	Type string `json:"type"`
	Name string `json:"name"`
	// A string or boolean for choice answers; absent for scores or refusals.
	Value any `json:"value"`
	// Selected numeric level for score answers. Inspect JSON.Score for presence.
	Score   float64 `json:"score"`
	Refusal string  `json:"refusal"`
	JSON    struct {
		Type, Name, Value, Score, Refusal respjson.Field
		ExtraFields                       map[string]respjson.Field
		raw                               string
	} `json:"-"`
}

func (r DecisionAnswer) RawJSON() string                  { return r.JSON.raw }
func (r *DecisionAnswer) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

// DecisionContext describes the server-selected input. Session-only fields may
// be absent for explicit state or when preparation fails; inspect JSON for presence.
type DecisionContext struct {
	Mode           string `json:"mode"`
	PolicyVersion  string `json:"policy_version"`
	TokenEstimator string `json:"token_estimator"`
	// Estimate of the full serialized model request, not billed model usage.
	EstimatedInputTokens int64     `json:"estimated_input_tokens"`
	MaxInputTokens       int64     `json:"max_input_tokens"`
	SnapshotAt           time.Time `json:"snapshot_at" format:"date-time"`
	AnchorEventID        string    `json:"anchor_event_id"`
	FirstEventID         string    `json:"first_event_id"`
	LastEventID          string    `json:"last_event_id"`
	IncludedSegmentCount int64     `json:"included_segment_count"`
	IncludedEventCount   int64     `json:"included_event_count"`
	ScanLimited          bool      `json:"scan_limited"`
	Truncated            bool      `json:"truncated"`
	TruncationReasons    []string  `json:"truncation_reasons"`
	JSON                 struct {
		Mode, PolicyVersion, TokenEstimator, EstimatedInputTokens, MaxInputTokens, SnapshotAt respjson.Field
		AnchorEventID, FirstEventID, LastEventID, IncludedSegmentCount, IncludedEventCount    respjson.Field
		ScanLimited, Truncated, TruncationReasons                                             respjson.Field
		ExtraFields                                                                           map[string]respjson.Field
		raw                                                                                   string
	} `json:"-"`
}

func (r DecisionContext) RawJSON() string                  { return r.JSON.raw }
func (r *DecisionContext) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type DecisionError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	JSON    struct {
		Code, Message respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

func (r DecisionError) RawJSON() string                  { return r.JSON.raw }
func (r *DecisionError) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
