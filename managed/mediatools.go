package managed

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/constant"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/thirdparty/gjson"
)

type ManagedAgentsImageGenToolConfigParams struct {
	Name             constant.ImageGen                                                 `json:"name" default:"ImageGen" api:"required"`
	Type             param.Opt[string]                                                 `json:"type,omitzero"`
	Enabled          param.Opt[bool]                                                   `json:"enabled,omitzero"`
	PermissionPolicy ManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion `json:"permission_policy,omitzero"`
	paramObj
}

func (r ManagedAgentsImageGenToolConfigParams) MarshalJSON() ([]byte, error) {
	type shadow ManagedAgentsImageGenToolConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsImageGenToolConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsImageSearchToolConfigParams struct {
	Name             constant.ImageSearch                                              `json:"name" default:"ImageSearch" api:"required"`
	Type             param.Opt[string]                                                 `json:"type,omitzero"`
	Enabled          param.Opt[bool]                                                   `json:"enabled,omitzero"`
	PermissionPolicy ManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion `json:"permission_policy,omitzero"`
	paramObj
}

func (r ManagedAgentsImageSearchToolConfigParams) MarshalJSON() ([]byte, error) {
	type shadow ManagedAgentsImageSearchToolConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsImageSearchToolConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsVideoGenToolConfigParams struct {
	Name             constant.VideoGen                                                 `json:"name" default:"VideoGen" api:"required"`
	Type             param.Opt[string]                                                 `json:"type,omitzero"`
	Enabled          param.Opt[bool]                                                   `json:"enabled,omitzero"`
	PermissionPolicy ManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion `json:"permission_policy,omitzero"`
	paramObj
}

func (r ManagedAgentsVideoGenToolConfigParams) MarshalJSON() ([]byte, error) {
	type shadow ManagedAgentsVideoGenToolConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsVideoGenToolConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsDeliverArtifactsToolConfigParams struct {
	Name             constant.DeliverArtifacts                                         `json:"name" default:"DeliverArtifacts" api:"required"`
	Type             param.Opt[string]                                                 `json:"type,omitzero"`
	Enabled          param.Opt[bool]                                                   `json:"enabled,omitzero"`
	PermissionPolicy ManagedAgentsAgentToolsetDefaultConfigParamsPermissionPolicyUnion `json:"permission_policy,omitzero"`
	paramObj
}

func (r ManagedAgentsDeliverArtifactsToolConfigParams) MarshalJSON() ([]byte, error) {
	type shadow ManagedAgentsDeliverArtifactsToolConfigParams
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ManagedAgentsDeliverArtifactsToolConfigParams) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

type ManagedAgentsImageGenToolConfig struct {
	Name             string                                                      `json:"name" api:"required"`
	Type             string                                                      `json:"type"`
	Enabled          bool                                                        `json:"enabled"`
	PermissionPolicy ManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion `json:"permission_policy"`
	JSON             struct {
		Name             respjson.Field
		Type             respjson.Field
		Enabled          respjson.Field
		PermissionPolicy respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

type ManagedAgentsImageSearchToolConfig struct {
	Name             string                                                      `json:"name" api:"required"`
	Type             string                                                      `json:"type"`
	Enabled          bool                                                        `json:"enabled"`
	PermissionPolicy ManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion `json:"permission_policy"`
	JSON             struct {
		Name             respjson.Field
		Type             respjson.Field
		Enabled          respjson.Field
		PermissionPolicy respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

type ManagedAgentsVideoGenToolConfig struct {
	Name             string                                                      `json:"name" api:"required"`
	Type             string                                                      `json:"type"`
	Enabled          bool                                                        `json:"enabled"`
	PermissionPolicy ManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion `json:"permission_policy"`
	JSON             struct {
		Name             respjson.Field
		Type             respjson.Field
		Enabled          respjson.Field
		PermissionPolicy respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

type ManagedAgentsDeliverArtifactsToolConfig struct {
	Name             string                                                      `json:"name" api:"required"`
	Type             string                                                      `json:"type"`
	Enabled          bool                                                        `json:"enabled"`
	PermissionPolicy ManagedAgentsAgentToolsetDefaultConfigPermissionPolicyUnion `json:"permission_policy"`
	JSON             struct {
		Name             respjson.Field
		Type             respjson.Field
		Enabled          respjson.Field
		PermissionPolicy respjson.Field
		ExtraFields      map[string]respjson.Field
		raw              string
	} `json:"-"`
}

func (r ManagedAgentsImageGenToolConfig) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsImageGenToolConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
func (r ManagedAgentsImageGenToolConfig) ToParam() ManagedAgentsImageGenToolConfigParams {
	return param.Override[ManagedAgentsImageGenToolConfigParams](json.RawMessage(r.RawJSON()))
}
func (r ManagedAgentsImageSearchToolConfig) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsImageSearchToolConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
func (r ManagedAgentsImageSearchToolConfig) ToParam() ManagedAgentsImageSearchToolConfigParams {
	return param.Override[ManagedAgentsImageSearchToolConfigParams](json.RawMessage(r.RawJSON()))
}
func (r ManagedAgentsVideoGenToolConfig) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsVideoGenToolConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
func (r ManagedAgentsVideoGenToolConfig) ToParam() ManagedAgentsVideoGenToolConfigParams {
	return param.Override[ManagedAgentsVideoGenToolConfigParams](json.RawMessage(r.RawJSON()))
}
func (r ManagedAgentsDeliverArtifactsToolConfig) RawJSON() string { return r.JSON.raw }
func (r *ManagedAgentsDeliverArtifactsToolConfig) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}
func (r ManagedAgentsDeliverArtifactsToolConfig) ToParam() ManagedAgentsDeliverArtifactsToolConfigParams {
	return param.Override[ManagedAgentsDeliverArtifactsToolConfigParams](json.RawMessage(r.RawJSON()))
}

func (ManagedAgentsImageGenToolConfig) implManagedAgentsAgentToolConfigUnion()         {}
func (ManagedAgentsImageSearchToolConfig) implManagedAgentsAgentToolConfigUnion()      {}
func (ManagedAgentsVideoGenToolConfig) implManagedAgentsAgentToolConfigUnion()         {}
func (ManagedAgentsDeliverArtifactsToolConfig) implManagedAgentsAgentToolConfigUnion() {}

func (u ManagedAgentsAgentToolConfigUnion) AsImageGen() (v ManagedAgentsImageGenToolConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}
func (u ManagedAgentsAgentToolConfigUnion) AsImageSearch() (v ManagedAgentsImageSearchToolConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}
func (u ManagedAgentsAgentToolConfigUnion) AsVideoGen() (v ManagedAgentsVideoGenToolConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}
func (u ManagedAgentsAgentToolConfigUnion) AsDeliverArtifacts() (v ManagedAgentsDeliverArtifactsToolConfig) {
	apijson.UnmarshalRoot(json.RawMessage(u.JSON.raw), &v)
	return
}

func init() {
	apijson.RegisterCustomDecoder[ManagedAgentsAgentToolConfigParamsUnion](func(node gjson.Result, value reflect.Value, defaultDecoder func(gjson.Result, reflect.Value) error) error {
		u := ManagedAgentsAgentToolConfigParamsUnion{}
		switch strings.ToLower(node.Get("name").String()) {
		case "bash":
			u.OfBash = &ManagedAgentsBashToolConfigParams{}
			if err := apijson.UnmarshalRoot([]byte(node.Raw), u.OfBash); err != nil {
				return err
			}
		case "edit":
			u.OfEdit = &ManagedAgentsEditToolConfigParams{}
			if err := apijson.UnmarshalRoot([]byte(node.Raw), u.OfEdit); err != nil {
				return err
			}
		case "read":
			u.OfRead = &ManagedAgentsReadToolConfigParams{}
			if err := apijson.UnmarshalRoot([]byte(node.Raw), u.OfRead); err != nil {
				return err
			}
		case "write":
			u.OfWrite = &ManagedAgentsWriteToolConfigParams{}
			if err := apijson.UnmarshalRoot([]byte(node.Raw), u.OfWrite); err != nil {
				return err
			}
		case "glob":
			u.OfGlob = &ManagedAgentsGlobToolConfigParams{}
			if err := apijson.UnmarshalRoot([]byte(node.Raw), u.OfGlob); err != nil {
				return err
			}
		case "grep":
			u.OfGrep = &ManagedAgentsGrepToolConfigParams{}
			if err := apijson.UnmarshalRoot([]byte(node.Raw), u.OfGrep); err != nil {
				return err
			}
		case "webfetch", "web_fetch":
			u.OfWebFetch = &ManagedAgentsWebFetchToolConfigParams{}
			if err := apijson.UnmarshalRoot([]byte(node.Raw), u.OfWebFetch); err != nil {
				return err
			}
		case "websearch", "web_search":
			u.OfWebSearch = &ManagedAgentsWebSearchToolConfigParams{}
			if err := apijson.UnmarshalRoot([]byte(node.Raw), u.OfWebSearch); err != nil {
				return err
			}
		case "imagegen", "image_gen":
			u.OfImageGen = &ManagedAgentsImageGenToolConfigParams{}
			if err := apijson.UnmarshalRoot([]byte(node.Raw), u.OfImageGen); err != nil {
				return err
			}
		case "imagesearch", "image_search":
			u.OfImageSearch = &ManagedAgentsImageSearchToolConfigParams{}
			if err := apijson.UnmarshalRoot([]byte(node.Raw), u.OfImageSearch); err != nil {
				return err
			}
		case "videogen", "video_gen":
			u.OfVideoGen = &ManagedAgentsVideoGenToolConfigParams{}
			if err := apijson.UnmarshalRoot([]byte(node.Raw), u.OfVideoGen); err != nil {
				return err
			}
		case "deliverartifacts", "deliver_artifacts":
			u.OfDeliverArtifacts = &ManagedAgentsDeliverArtifactsToolConfigParams{}
			if err := apijson.UnmarshalRoot([]byte(node.Raw), u.OfDeliverArtifacts); err != nil {
				return err
			}
		default:
			return defaultDecoder(node, value)
		}
		value.Set(reflect.ValueOf(u))
		return nil
	})
}
