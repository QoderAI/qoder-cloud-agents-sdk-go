package forward

import (
	"context"
	"net/http"
	"slices"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/option"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
)

// ModelService provides Forward Model operations.
type ModelService struct {
	Options []option.RequestOption
}

func NewModelService(opts ...option.RequestOption) ModelService {
	return ModelService{Options: slices.Clone(opts)}
}

// 列出模型.
func (r *ModelService) List(ctx context.Context, opts ...option.RequestOption) (res *ModelListResponse, err error) {

	opts = slices.Concat(r.Options, opts)
	path := "models"
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

type Model struct {
	ID                      string   `json:"id"`
	DisplayName             string   `json:"display_name"`
	IsEnabled               bool     `json:"is_enabled"`
	IsNew                   bool     `json:"is_new"`
	IsVl                    bool     `json:"is_vl"`
	SupportDisableReasoning bool     `json:"support_disable_reasoning"`
	PriceFactor             float64  `json:"price_factor"`
	Efforts                 []string `json:"efforts"`
	DefaultEffort           string   `json:"default_effort"`
	Speed                   []string `json:"speed"`
	MaxInputTokens          int64    `json:"max_input_tokens"`
	DefaultContextWindow    int64    `json:"default_context_window"`
	AvailableContextWindows []int64  `json:"available_context_windows"`
	JSON                    struct {
		ID                      respjson.Field
		DisplayName             respjson.Field
		IsEnabled               respjson.Field
		IsNew                   respjson.Field
		IsVl                    respjson.Field
		SupportDisableReasoning respjson.Field
		PriceFactor             respjson.Field
		Efforts                 respjson.Field
		DefaultEffort           respjson.Field
		Speed                   respjson.Field
		MaxInputTokens          respjson.Field
		DefaultContextWindow    respjson.Field
		AvailableContextWindows respjson.Field
		ExtraFields             map[string]respjson.Field
		raw                     string
	} `json:"-"`
}

func (r Model) RawJSON() string                  { return r.JSON.raw }
func (r *Model) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type ModelListResponse struct {
	Data    []Model `json:"data"`
	HasMore bool    `json:"has_more"`
	JSON    struct {
		Data        respjson.Field
		HasMore     respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r ModelListResponse) RawJSON() string                  { return r.JSON.raw }
func (r *ModelListResponse) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }
