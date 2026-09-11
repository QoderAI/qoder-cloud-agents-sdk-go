// Qoder managed API definitions.
package managed

import (
	"reflect"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/apijson"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/thirdparty/gjson"
)

func registerStringPromotion[SliceT ~[]E, E any](wrap func(string) E) {
	apijson.RegisterCustomDecoder[SliceT](func(node gjson.Result, value reflect.Value, defaultDecoder func(gjson.Result, reflect.Value) error) error {
		if node.Type == gjson.String {
			arrayValue := reflect.MakeSlice(value.Type(), 1, 1)
			arrayValue.Index(0).Set(reflect.ValueOf(wrap(node.String())))
			value.Set(arrayValue)
			return nil
		}
		return defaultDecoder(node, value)
	})
}

func init() {
	registerStringPromotion[[]TextBlockParam](func(s string) TextBlockParam {
		return TextBlockParam{Text: s}
	})
	registerStringPromotion[[]TextBlockParam](func(s string) TextBlockParam {
		return TextBlockParam{Text: s}
	})
	registerStringPromotion[[]ContentBlockParamUnion](func(s string) ContentBlockParamUnion {
		return ContentBlockParamUnion{OfText: &TextBlockParam{Text: s}}
	})
	registerStringPromotion[[]ToolResultBlockParamContentUnion](func(s string) ToolResultBlockParamContentUnion {
		return ToolResultBlockParamContentUnion{OfText: &TextBlockParam{Text: s}}
	})
}
