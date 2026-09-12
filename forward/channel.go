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

// ChannelService provides Forward Channel operations.
type ChannelService struct {
	Options    []option.RequestOption
	QRSessions ChannelQRSessionService
}

func NewChannelService(opts ...option.RequestOption) ChannelService {
	return ChannelService{Options: slices.Clone(opts), QRSessions: NewChannelQRSessionService(opts...)}
}

// List Channels
func (r *ChannelService) List(ctx context.Context, params ChannelListParams, opts ...option.RequestOption) (res *pagination.Page[Channel], err error) {

	opts = slices.Concat(r.Options, opts)
	path := "channels"
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
func (r *ChannelService) ListAutoPaging(ctx context.Context, params ChannelListParams, opts ...option.RequestOption) *pagination.PageAutoPager[Channel] {
	return pagination.NewPageAutoPager(r.List(ctx, params, opts...))
}

type ChannelListParams struct {
	// Filter by `wechat`, `wecom`, `feishu`, `dingtalk` or `teams` (Global).
	ChannelType param.Opt[string] `query:"channel_type,omitzero" json:"-"`
	// Filter by the manual enable/disable state.
	Enabled param.Opt[bool] `query:"enabled,omitzero" json:"-"`
	// Filter by `unbound`, `bound` or `expired`.
	BindingStatus param.Opt[string] `query:"binding_status,omitzero" json:"-"`
	// Filter by Forward Identity ID.
	IdentityID param.Opt[string] `query:"identity_id,omitzero" json:"-"`
	// Filter by Forward Template ID.
	TemplateID param.Opt[string] `query:"template_id,omitzero" json:"-"`
	// Page size, up to 100.
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// Cursor for the next page.
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// Cursor for the previous page.
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	paramObj
}

func (r ChannelListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// Create Channel
func (r *ChannelService) New(ctx context.Context, params ChannelNewParams, opts ...option.RequestOption) (res *Channel, err error) {
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := "channels"
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type ChannelNewParams struct {
	// Required in `fixed` mode; omit in `pairing` mode.
	IdentityID         param.Opt[string] `json:"identity_id,omitzero"`
	IdentityResolution map[string]any    `json:"identity_resolution,omitzero"`
	// Required in `fixed` mode; omit in `pairing` mode.
	TemplateID param.Opt[string] `json:"template_id,omitzero"`
	// Channel type; currently `wechat`, `wecom`, `feishu`, `dingtalk` and `teams` (Global).
	ChannelType string `json:"channel_type" api:"required"`
	// Channel display name.
	Name param.Opt[string] `json:"name,omitzero"`
	// Manual enable/disable switch, defaults to `true`. Pass `false` to create the Channel
	// without handling inbound messages yet.
	Enabled       param.Opt[bool] `json:"enabled,omitzero"`
	ChannelConfig map[string]any  `json:"channel_config,omitzero"`
	// Optional idempotency key for requests with side effects.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r ChannelNewParams) MarshalJSON() ([]byte, error) {
	type shadow ChannelNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ChannelNewParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

// Get Channel
func (r *ChannelService) Get(ctx context.Context, channelID string, opts ...option.RequestOption) (res *Channel, err error) {
	if channelID == "" {
		return nil, fmt.Errorf("missing required channel_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("channels/%s", url.PathEscape(channelID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update Channel
func (r *ChannelService) Update(ctx context.Context, channelID string, params ChannelUpdateParams, opts ...option.RequestOption) (res *Channel, err error) {
	if channelID == "" {
		return nil, fmt.Errorf("missing required channel_id parameter")
	}
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("channels/%s", url.PathEscape(channelID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type ChannelUpdateParams struct {
	// Channel display name.
	Name param.Opt[string] `json:"name,omitzero"`
	// New Forward Identity ID in `fixed` mode.
	IdentityID param.Opt[string] `json:"identity_id,omitzero"`
	// New Forward Template ID in `fixed` mode.
	TemplateID param.Opt[string] `json:"template_id,omitzero"`
	// Manual enable/disable switch.
	Enabled       param.Opt[bool] `json:"enabled,omitzero"`
	ChannelConfig map[string]any  `json:"channel_config,omitzero"`
	// Optional idempotency key for requests with side effects.
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r ChannelUpdateParams) MarshalJSON() ([]byte, error) {
	type shadow ChannelUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ChannelUpdateParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

// Delete Channel
func (r *ChannelService) Delete(ctx context.Context, channelID string, opts ...option.RequestOption) (res *DeletedChannel, err error) {
	if channelID == "" {
		return nil, fmt.Errorf("missing required channel_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("channels/%s", url.PathEscape(channelID))
	err = convention.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

type DeletedChannel struct {
	// ID of the deleted Channel.
	ID string `json:"id"`
	// Whether the deletion succeeded; always `true` on success.
	Deleted bool `json:"deleted"`
	JSON    struct {
		ID          respjson.Field
		Deleted     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r DeletedChannel) RawJSON() string                  { return r.JSON.raw }
func (r *DeletedChannel) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type Channel struct {
	// Channel ID, prefixed with `channel_` in examples.
	ID string `json:"id"`
	// Always `channel`.
	Type string `json:"type"`
	// Bound Forward Identity ID in `fixed` mode; `null` in `pairing` mode.
	IdentityID         string                    `json:"identity_id" api:"nullable"`
	IdentityResolution ChannelIdentityResolution `json:"identity_resolution"`
	// Bound Forward Template ID in `fixed` mode; `null` in `pairing` mode.
	TemplateID string `json:"template_id" api:"nullable"`
	// External channel type.
	ChannelType string `json:"channel_type"`
	Name        string `json:"name"`
	// Manual enable/disable switch.
	Enabled bool `json:"enabled"`
	// `unbound`, `bound` or `expired`.
	BindingStatus string               `json:"binding_status"`
	ChannelConfig ChannelChannelConfig `json:"channel_config"`
	CreatedAt     time.Time            `json:"created_at" format:"date-time"`
	UpdatedAt     time.Time            `json:"updated_at" format:"date-time"`
	JSON          struct {
		ID                 respjson.Field
		Type               respjson.Field
		IdentityID         respjson.Field
		IdentityResolution respjson.Field
		TemplateID         respjson.Field
		ChannelType        respjson.Field
		Name               respjson.Field
		Enabled            respjson.Field
		BindingStatus      respjson.Field
		ChannelConfig      respjson.Field
		CreatedAt          respjson.Field
		UpdatedAt          respjson.Field
		ExtraFields        map[string]respjson.Field
		raw                string
	} `json:"-"`
}

func (r Channel) RawJSON() string                  { return r.JSON.raw }
func (r *Channel) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type ChannelChannelConfig struct {
	ResponseOptions ChannelChannelConfigResponseOptions `json:"response_options"`
	JSON            struct {
		ResponseOptions respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

func (r ChannelChannelConfig) RawJSON() string { return r.JSON.raw }
func (r *ChannelChannelConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ChannelChannelConfigResponseOptions struct {
	IncludeToolCalls bool `json:"include_tool_calls"`
	IncludeThinking  bool `json:"include_thinking"`
	JSON             struct {
		IncludeToolCalls respjson.Field
		IncludeThinking  respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

func (r ChannelChannelConfigResponseOptions) RawJSON() string { return r.JSON.raw }
func (r *ChannelChannelConfigResponseOptions) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ChannelIdentityResolution struct {
	Mode string `json:"mode"`
	JSON struct {
		Mode        respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r ChannelIdentityResolution) RawJSON() string { return r.JSON.raw }
func (r *ChannelIdentityResolution) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
