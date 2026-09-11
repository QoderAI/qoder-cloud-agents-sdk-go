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

// 列出 Channels.
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
	// 按 `wechat`、`wecom`、`feishu`、`dingtalk` 或 `teams`（Global）过滤。
	ChannelType param.Opt[string] `query:"channel_type,omitzero" json:"-"`
	// 按人工启停状态过滤。
	Enabled param.Opt[bool] `query:"enabled,omitzero" json:"-"`
	// 按 `unbound`、`bound` 或 `expired` 过滤。
	BindingStatus param.Opt[string] `query:"binding_status,omitzero" json:"-"`
	// 按 Forward Identity ID 过滤。
	IdentityID param.Opt[string] `query:"identity_id,omitzero" json:"-"`
	// 按 Forward Template ID 过滤。
	TemplateID param.Opt[string] `query:"template_id,omitzero" json:"-"`
	// 分页大小，最大 100。
	Limit param.Opt[int64] `query:"limit,omitzero" json:"-"`
	// 向后翻页游标。
	AfterID param.Opt[string] `query:"after_id,omitzero" json:"-"`
	// 向前翻页游标。
	BeforeID param.Opt[string] `query:"before_id,omitzero" json:"-"`
	paramObj
}

func (r ChannelListParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// 创建 Channel.
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
	// `fixed` 模式必填；`pairing` 模式不传。
	IdentityID         param.Opt[string] `json:"identity_id,omitzero"`
	IdentityResolution map[string]any    `json:"identity_resolution,omitzero"`
	// `fixed` 模式必填；`pairing` 模式不传。
	TemplateID param.Opt[string] `json:"template_id,omitzero"`
	// 渠道类型，当前支持 `wechat`、`wecom`、`feishu`、`dingtalk` 和 `teams`（Global）。
	ChannelType string `json:"channel_type" api:"required"`
	// Channel 展示名。
	Name param.Opt[string] `json:"name,omitzero"`
	// 人工启停开关，默认 `true`。传 `false` 可创建后暂不处理上行消息。
	Enabled       param.Opt[bool] `json:"enabled,omitzero"`
	ChannelConfig map[string]any  `json:"channel_config,omitzero"`
	// 有副作用请求可选的幂等键。
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r ChannelNewParams) MarshalJSON() ([]byte, error) {
	type shadow ChannelNewParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ChannelNewParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

// 获取 Channel.
func (r *ChannelService) Get(ctx context.Context, channelID string, opts ...option.RequestOption) (res *Channel, err error) {
	if channelID == "" {
		return nil, fmt.Errorf("missing required channel_id parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("channels/%s", url.PathEscape(channelID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// 更新 Channel.
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
	// Channel 展示名。
	Name param.Opt[string] `json:"name,omitzero"`
	// `fixed` 模式下新的 Forward Identity ID。
	IdentityID param.Opt[string] `json:"identity_id,omitzero"`
	// `fixed` 模式下新的 Forward Template ID。
	TemplateID param.Opt[string] `json:"template_id,omitzero"`
	// 人工启停开关。
	Enabled       param.Opt[bool] `json:"enabled,omitzero"`
	ChannelConfig map[string]any  `json:"channel_config,omitzero"`
	// 有副作用请求可选的幂等键。
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r ChannelUpdateParams) MarshalJSON() ([]byte, error) {
	type shadow ChannelUpdateParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ChannelUpdateParams) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

// 删除 Channel.
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
	// 被删除的 Channel ID。
	ID string `json:"id"`
	// 是否删除成功，成功时恒为 `true`。
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
	// Channel ID，示例前缀 `channel_`。
	ID string `json:"id"`
	// 固定为 `channel`。
	Type string `json:"type"`
	// `fixed` 模式为绑定的 Forward Identity ID；`pairing` 模式为 `null`。
	IdentityID         string                    `json:"identity_id" api:"nullable"`
	IdentityResolution ChannelIdentityResolution `json:"identity_resolution"`
	// `fixed` 模式为绑定的 Forward Template ID；`pairing` 模式为 `null`。
	TemplateID string `json:"template_id" api:"nullable"`
	// 外部渠道类型。
	ChannelType string `json:"channel_type"`
	Name        string `json:"name"`
	// 人工启停开关。
	Enabled bool `json:"enabled"`
	// `unbound`、`bound` 或 `expired`。
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
