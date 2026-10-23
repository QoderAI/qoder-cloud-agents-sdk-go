package forward

import (
	"context"
	"fmt"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apiquery"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/pagination"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
	"net/http"
	"net/url"
	"slices"
	"time"
)

// WebhookService manages account webhook subscriptions using a PAT or administrator SAT.
type WebhookService struct{ Endpoints WebhookEndpointService }

func NewWebhookService(opts ...option.RequestOption) WebhookService {
	return WebhookService{Endpoints: NewWebhookEndpointService(opts...)}
}

type WebhookEndpointService struct{ Options []option.RequestOption }

func NewWebhookEndpointService(opts ...option.RequestOption) WebhookEndpointService {
	return WebhookEndpointService{Options: slices.Clone(opts)}
}

// New a webhook endpoint.
func (r *WebhookEndpointService) New(ctx context.Context, params WebhookEndpointNewParams, opts ...option.RequestOption) (res *WebhookEndpointCreateResponse, err error) {
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", params.IdempotencyKey.Value)}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := "webhook/endpoints"
	err = convention.ExecuteNewRequest(ctx, "POST", path, params, &res, opts...)
	return res, err
}

// List a webhook endpoint page.
func (r *WebhookEndpointService) List(ctx context.Context, params WebhookEndpointListParams, opts ...option.RequestOption) (res *pagination.Page[WebhookEndpoint], err error) {
	opts = slices.Concat(r.Options, opts)
	path := "webhook/endpoints"
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
func (r *WebhookEndpointService) ListAutoPaging(ctx context.Context, params WebhookEndpointListParams, opts ...option.RequestOption) *pagination.PageAutoPager[WebhookEndpoint] {
	return pagination.NewPageAutoPager(r.List(ctx, params, opts...))
}

// Get a webhook endpoint.
func (r *WebhookEndpointService) Get(ctx context.Context, endpointID string, opts ...option.RequestOption) (res *WebhookEndpoint, err error) {
	if endpointID == "" {
		return nil, fmt.Errorf("missing required endpointID parameter")
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("webhook/endpoints/%s", url.PathEscape(endpointID))
	err = convention.ExecuteNewRequest(ctx, "GET", path, nil, &res, opts...)
	return res, err
}

// Update a webhook endpoint.
func (r *WebhookEndpointService) Update(ctx context.Context, endpointID string, params WebhookEndpointUpdateParams, opts ...option.RequestOption) (res *WebhookEndpoint, err error) {
	if endpointID == "" {
		return nil, fmt.Errorf("missing required endpointID parameter")
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("webhook/endpoints/%s", url.PathEscape(endpointID))
	err = convention.ExecuteNewRequest(ctx, "PUT", path, params, &res, opts...)
	return res, err
}

// Delete a webhook endpoint.
func (r *WebhookEndpointService) Delete(ctx context.Context, endpointID string, opts ...option.RequestOption) (err error) {
	if endpointID == "" {
		return fmt.Errorf("missing required endpointID parameter")
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("webhook/endpoints/%s", url.PathEscape(endpointID))
	err = convention.ExecuteNewRequest(ctx, "DELETE", path, nil, nil, opts...)
	return err
}

// Enable a webhook endpoint.
func (r *WebhookEndpointService) Enable(ctx context.Context, endpointID string, opts ...option.RequestOption) (res *WebhookEndpoint, err error) {
	if endpointID == "" {
		return nil, fmt.Errorf("missing required endpointID parameter")
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("webhook/endpoints/%s/enable", url.PathEscape(endpointID))
	err = convention.ExecuteNewRequest(ctx, "POST", path, http.NoBody, &res, opts...)
	return res, err
}

// Disable a webhook endpoint.
func (r *WebhookEndpointService) Disable(ctx context.Context, endpointID string, opts ...option.RequestOption) (res *WebhookEndpoint, err error) {
	if endpointID == "" {
		return nil, fmt.Errorf("missing required endpointID parameter")
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("webhook/endpoints/%s/disable", url.PathEscape(endpointID))
	err = convention.ExecuteNewRequest(ctx, "POST", path, http.NoBody, &res, opts...)
	return res, err
}

// Test queues a synthetic event; acceptance does not confirm delivery.
func (r *WebhookEndpointService) Test(ctx context.Context, endpointID string, opts ...option.RequestOption) (res *WebhookEndpointTestResponse, err error) {
	if endpointID == "" {
		return nil, fmt.Errorf("missing required endpointID parameter")
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("webhook/endpoints/%s/test", url.PathEscape(endpointID))
	err = convention.ExecuteNewRequest(ctx, "POST", path, http.NoBody, &res, opts...)
	return res, err
}

type WebhookEndpointListParams struct {
	Limit   param.Opt[int64]  `query:"limit,omitzero" json:"-"`
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	paramObj
}

func (r WebhookEndpointListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

type WebhookEndpointNewParams struct {
	URL string `json:"url" api:"required"`
	// Full event names or "*". Unknown event names are passed to the server.
	Events         []string          `json:"events" api:"required"`
	Description    param.Opt[string] `json:"description,omitzero"`
	Metadata       map[string]string `json:"metadata,omitzero"`
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}
type WebhookEndpointUpdateParams struct {
	URL         param.Opt[string] `json:"url,omitzero"`
	Description param.Opt[string] `json:"description,omitzero"`
	Events      []string          `json:"events,omitzero"`
	Active      param.Opt[bool]   `json:"active,omitzero"`
	// Values are merged; param.Null[string]() removes a key.
	Metadata map[string]param.Opt[string] `json:"metadata,omitzero"`
	paramObj
}

func (r WebhookEndpointNewParams) MarshalJSON() ([]byte, error) {
	type shadow WebhookEndpointNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *WebhookEndpointNewParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
func (r WebhookEndpointUpdateParams) MarshalJSON() ([]byte, error) {
	type shadow WebhookEndpointUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *WebhookEndpointUpdateParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// WebhookEndpoint preserves unknown response fields.
type WebhookEndpoint struct {
	ID              string         `json:"id"`
	URL             string         `json:"url"`
	Description     string         `json:"description"`
	Events          []string       `json:"events"`
	Metadata        map[string]any `json:"metadata"`
	Active          bool           `json:"active"`
	CreatedAt       time.Time      `json:"created_at" format:"date-time"`
	UpdatedAt       time.Time      `json:"updated_at" format:"date-time"`
	LastSuccessAt   time.Time      `json:"last_success_at" format:"date-time" api:"nullable"`
	LastFailureAt   time.Time      `json:"last_failure_at" format:"date-time" api:"nullable"`
	DisabledReason  string         `json:"disabled_reason"`
	ConsecutiveFail int64          `json:"consecutive_fail"`
	JSON            struct {
		ID              respjson.Field
		URL             respjson.Field
		Description     respjson.Field
		Events          respjson.Field
		Metadata        respjson.Field
		Active          respjson.Field
		CreatedAt       respjson.Field
		UpdatedAt       respjson.Field
		LastSuccessAt   respjson.Field
		LastFailureAt   respjson.Field
		DisabledReason  respjson.Field
		ConsecutiveFail respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

func (r WebhookEndpoint) RawJSON() string                  { return r.JSON.raw }
func (r *WebhookEndpoint) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

// WebhookEndpointCreateResponse preserves unknown response fields.
type WebhookEndpointCreateResponse struct {
	ID          string         `json:"id"`
	URL         string         `json:"url"`
	Description string         `json:"description"`
	Events      []string       `json:"events"`
	Metadata    map[string]any `json:"metadata"`
	Active      bool           `json:"active"`
	CreatedAt   time.Time      `json:"created_at" format:"date-time"`
	// Returned only at creation. Store securely; do not log.
	SigningSecret string `json:"signing_secret"`
	JSON          struct {
		ID            respjson.Field
		URL           respjson.Field
		Description   respjson.Field
		Events        respjson.Field
		Metadata      respjson.Field
		Active        respjson.Field
		CreatedAt     respjson.Field
		SigningSecret respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

func (r WebhookEndpointCreateResponse) RawJSON() string { return r.JSON.raw }
func (r *WebhookEndpointCreateResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// WebhookEndpointTestResponse preserves unknown response fields.
type WebhookEndpointTestResponse struct {
	EventID      int64 `json:"event_id"`
	DeliveryRows int64 `json:"delivery_rows"`
	JSON         struct {
		EventID      respjson.Field
		DeliveryRows respjson.Field
		ExtraFields  map[string]respjson.Field
		raw          string
	} `json:"-"`
}

func (r WebhookEndpointTestResponse) RawJSON() string { return r.JSON.raw }
func (r *WebhookEndpointTestResponse) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
