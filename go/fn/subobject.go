// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fn

import (
	"fmt"
	"reflect"

	"github.com/GoogleContainerTools/kpt-functions-sdk/go/fn/internal"
)

type SubObject struct {
	m *internal.MapVariant
}

// SetString is a helper for setting the specified field to a string value.
func (o *SubObject) SetString(key string, value string) {
	o.m.SetString(key, value)
}

// GetString returns the specified field; if not found or not a string it returns false.
func (o *SubObject) GetString(key string) (string, bool) {
	return o.m.GetString(key)
}

func (o *SubObject) StringStringEntries() (map[string]string, error) {
	entries, err := o.m.Entries()
	if err != nil {
		return nil, fmt.Errorf("error parsing configmap data: %w", err)
	}
	data := make(map[string]string)
	for key, v := range entries {
		vString, ok := internal.AsString(v.Node())
		if ok {
			data[key] = vString
		}
	}
	return data, nil
}

// Set sets a nested field located by fields to the value provided as val. val
// should not be a yaml.RNode. If you want to deal with yaml.RNode, you should
// use Get method and modify the underlying yaml.Node.
func (o *SubObject) Set(val interface{}, fields ...string) error {
	err := func() error {
		if o == nil {
			return fmt.Errorf("the object doesn't exist")
		}
		if val == nil {
			return fmt.Errorf("the passed-in object must not be nil")
		}
		kind := reflect.ValueOf(val).Kind()
		if kind == reflect.Ptr {
			kind = reflect.TypeOf(val).Elem().Kind()
		}

		switch kind {
		case reflect.Struct, reflect.Map:
			m, err := internal.TypedObjectToMapVariant(val)
			if err != nil {
				return err
			}
			return o.m.SetNestedMap(m, fields...)
		case reflect.Slice:
			s, err := internal.TypedObjectToSliceVariant(val)
			if err != nil {
				return err
			}
			return o.m.SetNestedSlice(s, fields...)
		case reflect.String:
			var s string
			switch val := val.(type) {
			case string:
				s = val
			case *string:
				s = *val
			}
			return o.m.SetNestedString(s, fields...)
		case reflect.Int, reflect.Int64:
			var i int
			switch val := val.(type) {
			case int:
				i = val
			case *int:
				i = *val
			case int64:
				i = int(val)
			case *int64:
				i = int(*val)
			}
			return o.m.SetNestedInt(i, fields...)
		case reflect.Float64:
			var f float64
			switch val := val.(type) {
			case float64:
				f = val
			case *float64:
				f = *val
			}
			return o.m.SetNestedFloat(f, fields...)
		case reflect.Bool:
			var b bool
			switch val := val.(type) {
			case bool:
				b = val
			case *bool:
				b = *val
			}
			return o.m.SetNestedBool(b, fields...)
		default:
			return fmt.Errorf("unhandled kind %s", kind)
		}
	}()
	if err != nil {
		return fmt.Errorf("unable to set %v at fields %v with error: %w", val, fields, err)
	}
	return nil
}

func (o *SubObject) UpsertMap(k string) *SubObject {
	m := o.m.UpsertMap(k)
	return &SubObject{m: m}
}
