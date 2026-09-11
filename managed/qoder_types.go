package managed

import (
	"encoding/json"
	"reflect"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/param"
	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/thirdparty/gjson"
)

// Qoder accepts compact strings as well as objects. Keep the public
// object fields usable for both shapes, including their original RawJSON.
func promoteStringObject[T any](key string, setRaw func(*T, string)) {
	apijson.RegisterCustomDecoder[T](func(node gjson.Result, value reflect.Value, decode func(gjson.Result, reflect.Value) error) error {
		if node.Type != gjson.String {
			return decode(node, value)
		}
		b, err := json.Marshal(map[string]string{key: node.String()})
		if err != nil {
			return err
		}
		if err = decode(gjson.ParseBytes(b), value); err != nil {
			return err
		}
		setRaw(value.Addr().Interface().(*T), node.Raw)
		return nil
	})
}
func init() {
	apijson.RegisterCustomDecoder[ManagedAgentsModelConfigParams](func(node gjson.Result, value reflect.Value, decode func(gjson.Result, reflect.Value) error) error {
		if node.Type != gjson.String {
			return decode(node, value)
		}
		model := param.Override[ManagedAgentsModelConfigParams](node.String())
		model.ID = node.String()
		value.Set(reflect.ValueOf(model))
		return nil
	})
	promoteStringObject[ManagedAgentsEffortLow]("type", func(v *ManagedAgentsEffortLow, raw string) { v.JSON.raw = raw })
	promoteStringObject[ManagedAgentsEffortMedium]("type", func(v *ManagedAgentsEffortMedium, raw string) { v.JSON.raw = raw })
	promoteStringObject[ManagedAgentsEffortHigh]("type", func(v *ManagedAgentsEffortHigh, raw string) { v.JSON.raw = raw })
	promoteStringObject[ManagedAgentsEffortXhigh]("type", func(v *ManagedAgentsEffortXhigh, raw string) { v.JSON.raw = raw })
	promoteStringObject[ManagedAgentsEffortMax]("type", func(v *ManagedAgentsEffortMax, raw string) { v.JSON.raw = raw })

	promoteStringObject[SkillSource]("type", func(v *SkillSource, raw string) { v.JSON.raw = raw })
	promoteStringObject[ManagedAgentsModelConfig]("id", func(v *ManagedAgentsModelConfig, raw string) { v.JSON.raw = raw })
	promoteStringObject[ManagedAgentsModelConfigEffortUnion]("type", func(v *ManagedAgentsModelConfigEffortUnion, raw string) { v.JSON.raw = raw })
}
