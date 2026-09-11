package forward

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apiquery"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
)

// ChannelQRSessionService provides Forward ChannelQRSession operations.
type ChannelQRSessionService struct {
	Options []option.RequestOption
}

func NewChannelQRSessionService(opts ...option.RequestOption) ChannelQRSessionService {
	return ChannelQRSessionService{Options: slices.Clone(opts)}
}

// 创建 Channel QR Session.
func (r *ChannelQRSessionService) New(ctx context.Context, channelID string, params ChannelQRSessionNewParams, opts ...option.RequestOption) (res *ChannelQRSession, err error) {
	if channelID == "" {
		return nil, fmt.Errorf("missing required channel_id parameter")
	}
	if params.IdempotencyKey.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("Idempotency-Key", fmt.Sprint(params.IdempotencyKey.Value))}, opts...)
	}
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("channels/%s/qr_sessions", url.PathEscape(channelID))
	err = convention.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type ChannelQRSessionNewParams struct {
	// 有副作用请求可选的幂等键。
	IdempotencyKey param.Opt[string] `header:"Idempotency-Key,omitzero" json:"-"`
	paramObj
}

func (r ChannelQRSessionNewParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

// 获取 Channel QR Session.
func (r *ChannelQRSessionService) Get(ctx context.Context, sessionKey string, opts ...option.RequestOption) (res *ChannelQRSession, err error) {
	if sessionKey == "" {
		return nil, fmt.Errorf("missing required session_key parameter")
	}

	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("qr_sessions/%s", url.PathEscape(sessionKey))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

type ChannelQRSession struct {
	// 用于轮询状态的不透明 QR session key。
	SessionKey string `json:"session_key"`
	// 关联的 Channel ID。
	ChannelID string `json:"channel_id"`
	// `wechat`、`feishu`、`dingtalk` 或 `wecom`。
	ChannelType string `json:"channel_type"`
	// 初始状态，通常为 `waiting`。
	Status string `json:"status"`
	// 二维码原始内容，通常是三方授权 URL。
	QRCodeContent string `json:"qr_code_content"`
	// 服务端生成的二维码图片。
	QRCodeImageBase64 string `json:"qr_code_image_base64"`
	// 过期时间。
	ExpiresAt string `json:"expires_at"`
	// null|失败时的渠道错误码。
	ErrCode string `json:"err_code"`
	// null|失败时的渠道错误信息。
	ErrMsg string `json:"err_msg"`
	JSON   struct {
		SessionKey        respjson.Field
		ChannelID         respjson.Field
		ChannelType       respjson.Field
		Status            respjson.Field
		QRCodeContent     respjson.Field
		QRCodeImageBase64 respjson.Field
		ExpiresAt         respjson.Field
		ErrCode           respjson.Field
		ErrMsg            respjson.Field
		ExtraFields       map[string]respjson.Field
		raw               string
	} `json:"-"`
}

func (r ChannelQRSession) RawJSON() string                  { return r.JSON.raw }
func (r *ChannelQRSession) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
