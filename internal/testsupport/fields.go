package testsupport

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

	"github.com/QoderAI/qoder-cloud-agents-sdk-go/convention/respjson"
)

// CheckFields validates concrete decoded values as well as presence metadata,
// including objects inside dynamic maps and unions.
func CheckFields(t *testing.T, v reflect.Value, path string) {
	t.Helper()
	if !v.IsValid() || (v.Kind() == reflect.Pointer && v.IsNil()) {
		t.Fatalf("%s: nil response", path)
	}
	checkFields(t, v, path)
}
func checkFields(t *testing.T, v reflect.Value, path string) {
	t.Helper()
	for v.IsValid() && (v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface) {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}
	if !v.IsValid() {
		return
	}
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < v.Len(); i++ {
			checkFields(t, v.Index(i), fmt.Sprintf("%s[%d]", path, i))
		}
	case reflect.Map:
		it := v.MapRange()
		for it.Next() {
			checkFields(t, it.Value(), fmt.Sprintf("%s[%v]", path, it.Key()))
		}
	case reflect.Struct:
		if meta := v.FieldByName("JSON"); meta.IsValid() && meta.Kind() == reflect.Struct {
			for i := 0; i < meta.NumField(); i++ {
				if !meta.Field(i).CanInterface() {
					continue
				}
				f, ok := meta.Field(i).Interface().(respjson.Field)
				if !ok || f.Raw() == "" || f.Raw() == "null" {
					continue
				}
				name := meta.Type().Field(i).Name
				if !f.Valid() {
					t.Errorf("%s.%s invalid decoded field: %s", path, name, f.Raw())
					continue
				}
				actual := v.FieldByName(name)
				if !actual.IsValid() || !actual.CanInterface() {
					continue
				}
				switch actual.Kind() {
				case reflect.String, reflect.Bool, reflect.Int, reflect.Int64, reflect.Float64:
					want := reflect.New(actual.Type())
					if err := json.Unmarshal([]byte(f.Raw()), want.Interface()); err != nil {
						t.Errorf("%s.%s: %v", path, name, err)
					} else if !reflect.DeepEqual(actual.Interface(), want.Elem().Interface()) {
						t.Errorf("%s.%s decoded=%v, wire=%s", path, name, actual.Interface(), f.Raw())
					}
				}
			}
		}
		for i := 0; i < v.NumField(); i++ {
			f := v.Type().Field(i)
			if f.IsExported() && f.Name != "JSON" {
				checkFields(t, v.Field(i), path+"."+f.Name)
			}
		}
	}
}

// InvokeJSON calls one service with a documented request example. Path arguments
// are deliberately awkward to exercise escaping independently of payload values.
func InvokeJSON(t *testing.T, client any, ep Endpoint, body json.RawMessage) (any, error) {
	return call(t, serviceMethod(t, client, ep), context.Background(), -1, body)
}
