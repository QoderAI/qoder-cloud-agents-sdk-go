package forward

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"reflect"
	"strings"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apiform"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/thirdparty/gjson"
)

type paramObj = param.APIObject
type paramUnion = param.APIUnion
type Error = convention.Error

func String(v string) param.Opt[string]  { return param.NewOpt(v) }
func Int(v int64) param.Opt[int64]       { return param.NewOpt(v) }
func Bool(v bool) param.Opt[bool]        { return param.NewOpt(v) }
func Float(v float64) param.Opt[float64] { return param.NewOpt(v) }

// ModelConfig preserves either a model ID string or an object from the API.
// ID is populated for both shapes; RawJSON retains the original representation.
type ModelConfig struct {
	ID            string `json:"id"`
	Effort        string `json:"effort"`
	ContextWindow int64  `json:"context_window"`
	JSON          struct {
		ID            respjson.Field
		Effort        respjson.Field
		ContextWindow respjson.Field
		ExtraFields   map[string]respjson.Field
		raw           string
	} `json:"-"`
}

func (r ModelConfig) RawJSON() string                  { return r.JSON.raw }
func (r *ModelConfig) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

type ModelConfigParam struct {
	ID            string            `json:"id" api:"required"`
	Effort        param.Opt[string] `json:"effort,omitzero"`
	ContextWindow param.Opt[int64]  `json:"context_window,omitzero"`
	paramObj
}

func (r ModelConfigParam) MarshalJSON() ([]byte, error) {
	type shadow ModelConfigParam
	return param.MarshalObject(r, (*shadow)(&r))
}
func (r *ModelConfigParam) UnmarshalJSON(data []byte) error { return apijson.UnmarshalRoot(data, r) }

// ModelConfigUnionParam accepts one of a model ID or model configuration.
type ModelConfigUnionParam struct {
	OfString param.Opt[string] `json:",omitzero,inline"`
	OfConfig *ModelConfigParam `json:",omitzero,inline"`
	paramUnion
}

func (r ModelConfigUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(r, r.OfString, r.OfConfig)
}
func (r *ModelConfigUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EnvironmentVariablesUnionParam accepts an environment map or the legacy string form.
// Map values may be nil to express an explicit null during updates.
type EnvironmentVariablesUnionParam struct {
	OfMap    map[string]any    `json:",omitzero,inline"`
	OfString param.Opt[string] `json:",omitzero,inline"`
	paramUnion
}

func (r EnvironmentVariablesUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(r, r.OfMap, r.OfString)
}
func (r *EnvironmentVariablesUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// EventContentUnionParam represents message blocks or a tool result payload.
type EventContentUnionParam struct {
	OfBlocks []ContentBlockParam `json:",omitzero,inline"`
	OfString param.Opt[string]   `json:",omitzero,inline"`
	OfObject map[string]any      `json:",omitzero,inline"`
	paramUnion
}

func (r EventContentUnionParam) MarshalJSON() ([]byte, error) {
	return param.MarshalUnion(r, r.OfBlocks, r.OfString, r.OfObject)
}
func (r *EventContentUnionParam) UnmarshalJSON(data []byte) error {
	return apijson.UnmarshalRoot(data, r)
}

// ContentBlocks decodes block content; non-block event payloads remain in Content and RawJSON.
func (r SessionEvent) ContentBlocks() ([]ContentBlock, error) {
	var blocks []ContentBlock
	err := json.Unmarshal(r.Content, &blocks)
	return blocks, err
}

func init() {
	registerParamExtras[ModelConfigParam]()
	registerParamExtras[ToolParam]()
	registerParamExtras[ToolConfigParam]()
	registerParamExtras[MCPServerParam]()
	registerParamExtras[SkillBindingParam]()
	registerParamExtras[MultiagentConfigParam]()
	registerParamExtras[MultiagentEntryParam]()
	registerParamExtras[ResourceBindingParam]()
	registerParamExtras[GitHubRepositoryParam]()
	registerParamExtras[IdentityConfigSpecParam]()
	registerParamExtras[ToolOverrideParam]()
	registerParamExtras[MCPServerOverrideParam]()
	registerParamExtras[SkillOverrideParam]()
	registerParamExtras[ContentBlockParam]()
	registerParamExtras[ImageSourceParam]()
	registerParamExtras[SessionEventParam]()
	registerParamExtras[SessionResourceSpecParam]()
	apijson.RegisterCustomDecoder[ModelConfig](func(node gjson.Result, value reflect.Value, decode func(gjson.Result, reflect.Value) error) error {
		if node.Type != gjson.String {
			return decode(node, value)
		}
		data, err := json.Marshal(map[string]string{"id": node.String()})
		if err != nil {
			return err
		}
		if err = decode(gjson.ParseBytes(data), value); err != nil {
			return err
		}
		value.Addr().Interface().(*ModelConfig).JSON.raw = node.Raw
		return nil
	})
}

// These configuration objects were open maps in the wire contract. Keep unknown
// fields when callers decode an existing JSON configuration into typed params.
func registerParamExtras[T any]() {
	known := map[string]bool{}
	typ := reflect.TypeFor[T]()
	for i := 0; i < typ.NumField(); i++ {
		name := strings.Split(typ.Field(i).Tag.Get("json"), ",")[0]
		if name != "" && name != "-" {
			known[name] = true
		}
	}
	apijson.RegisterCustomDecoder[T](func(node gjson.Result, value reflect.Value, decode func(gjson.Result, reflect.Value) error) error {
		if err := decode(node, value); err != nil {
			return err
		}
		if !node.IsObject() {
			return nil
		}
		extra := map[string]any{}
		node.ForEach(func(key, item gjson.Result) bool {
			if !known[key.String()] {
				extra[key.String()] = json.RawMessage(item.Raw)
			}
			return true
		})
		if len(extra) != 0 {
			value.Addr().Interface().(interface{ SetExtraFields(map[string]any) }).SetExtraFields(extra)
		}
		return nil
	})
}

func marshalMultipart(value any, extras map[string]any) ([]byte, string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	if err := apiform.MarshalRoot(value, writer); err != nil {
		_ = writer.Close()
		return nil, "", err
	}
	if err := apiform.WriteExtras(writer, extras); err != nil {
		_ = writer.Close()
		return nil, "", err
	}
	if err := writer.Close(); err != nil {
		return nil, "", err
	}
	return buf.Bytes(), writer.FormDataContentType(), nil
}
