package forward

import (
	"context"
	"encoding/json"
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

type TemplateA2AConfigParam struct {
	Enabled     param.Opt[bool]         `json:"enabled,omitzero"`
	Description param.Opt[string]       `json:"description,omitzero"`
	Skills      []TemplateA2ASkillParam `json:"skills,omitzero"`
	configParamObj
}

func (r TemplateA2AConfigParam) MarshalJSON() ([]byte, error) {
	type shadow TemplateA2AConfigParam
	return marshalConfigParam(r, (*shadow)(&r))
}
func (r *TemplateA2AConfigParam) UnmarshalAPIJSON(data []byte) error { return r.UnmarshalJSON(data) }

func (r *TemplateA2AConfigParam) UnmarshalJSON(data []byte) error {
	return unmarshalObjectParam[TemplateA2AConfigParam](data, r)
}

func unmarshalTemplateA2A(data []byte, dst *TemplateA2AConfigParam) error {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(data, &fields); err != nil {
		return err
	}
	if raw, present := fields["a2a"]; present {
		return dst.UnmarshalJSON(raw)
	}
	return nil
}

type TemplateA2ASkillParam struct {
	ID          string   `json:"id" api:"required"`
	Name        string   `json:"name" api:"required"`
	Description string   `json:"description" api:"required"`
	Tags        []string `json:"tags,omitzero"`
	Examples    []string `json:"examples,omitzero"`
	configParamObj
}

func (r TemplateA2ASkillParam) MarshalJSON() ([]byte, error) {
	type shadow TemplateA2ASkillParam
	return marshalConfigParam(r, (*shadow)(&r))
}
func (r *TemplateA2ASkillParam) UnmarshalAPIJSON(data []byte) error { return r.UnmarshalJSON(data) }
func (r *TemplateA2ASkillParam) UnmarshalJSON(data []byte) error {
	return unmarshalObjectParam[TemplateA2ASkillParam](data, r)
}

type TemplateA2AConfig struct {
	Enabled     bool               `json:"enabled"`
	Description string             `json:"description"`
	Skills      []TemplateA2ASkill `json:"skills"`
	JSON        struct {
		Enabled     respjson.Field
		Description respjson.Field
		Skills      respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r TemplateA2AConfig) RawJSON() string                  { return r.JSON.raw }
func (r *TemplateA2AConfig) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type TemplateA2ASkill struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Examples    []string `json:"examples"`
	JSON        struct {
		ID          respjson.Field
		Name        respjson.Field
		Description respjson.Field
		Tags        respjson.Field
		Examples    respjson.Field
		ExtraFields map[string]respjson.Field
		raw         string
	} `json:"-"`
}

func (r TemplateA2ASkill) RawJSON() string                  { return r.JSON.raw }
func (r *TemplateA2ASkill) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

func (r *TemplateService) GetAgentCard(ctx context.Context, templateID string, params TemplateGetAgentCardParams, opts ...option.RequestOption) (res *TemplateAgentCard, err error) {
	if templateID == "" {
		return nil, fmt.Errorf("missing required template_id parameter")
	}
	if params.A2AVersion.Valid() {
		opts = append([]option.RequestOption{option.WithHeader("A2A-Version", params.A2AVersion.Value)}, opts...)
	}
	path := fmt.Sprintf("templates/%s/a2a/agent-card.json", url.PathEscape(templateID))
	err = convention.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, slices.Concat(r.Options, opts)...)
	return res, err
}

type TemplateGetAgentCardParams struct {
	IdentityID param.Opt[string] `query:"identity_id,omitzero" json:"-"`
	A2AVersion param.Opt[string] `header:"A2A-Version,omitzero" json:"-"`
	paramObj
}

func (r TemplateGetAgentCardParams) URLQuery() (url.Values, error) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{ArrayFormat: apiquery.ArrayQueryFormatRepeat, NestedFormat: apiquery.NestedQueryFormatBrackets})
}

type TemplateAgentCard struct {
	Name                string                       `json:"name"`
	Description         string                       `json:"description"`
	Version             string                       `json:"version"`
	Skills              []TemplateA2ASkill           `json:"skills"`
	DefaultInputModes   []string                     `json:"defaultInputModes"`
	DefaultOutputModes  []string                     `json:"defaultOutputModes"`
	URL                 string                       `json:"url"`
	SupportedInterfaces []TemplateAgentCardInterface `json:"supportedInterfaces"`
	JSON                struct {
		Name                respjson.Field
		Description         respjson.Field
		Version             respjson.Field
		Skills              respjson.Field
		DefaultInputModes   respjson.Field
		DefaultOutputModes  respjson.Field
		URL                 respjson.Field
		SupportedInterfaces respjson.Field
		ExtraFields         map[string]respjson.Field
		raw                 string
	} `json:"-"`
}

func (r TemplateAgentCard) RawJSON() string                  { return r.JSON.raw }
func (r *TemplateAgentCard) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type TemplateAgentCardInterface struct {
	URL             string `json:"url"`
	ProtocolBinding string `json:"protocolBinding"`
	ProtocolVersion string `json:"protocolVersion"`
	JSON            struct {
		URL             respjson.Field
		ProtocolBinding respjson.Field
		ProtocolVersion respjson.Field
		ExtraFields     map[string]respjson.Field
		raw             string
	} `json:"-"`
}

func (r TemplateAgentCardInterface) RawJSON() string { return r.JSON.raw }
func (r *TemplateAgentCardInterface) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
