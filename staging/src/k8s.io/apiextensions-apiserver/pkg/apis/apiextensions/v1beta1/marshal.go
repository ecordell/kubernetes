/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"errors"
	"strings"

	"k8s.io/apimachinery/pkg/util/json"
)

var jsTrue = []byte("true")
var jsFalse = []byte("false")

func (s JSONSchemaPropsOrBool) MarshalJSON() ([]byte, error) {
	if s.Schema != nil {
		return json.Marshal(s.Schema)
	}

	if s.Schema == nil && !s.Allows {
		return jsFalse, nil
	}
	return jsTrue, nil
}

func (s *JSONSchemaPropsOrBool) UnmarshalJSON(data []byte) error {
	var nw JSONSchemaPropsOrBool
	switch {
	case len(data) == 0:
	case data[0] == '{':
		var sch JSONSchemaProps
		if err := json.Unmarshal(data, &sch); err != nil {
			return err
		}
		nw.Allows = true
		nw.Schema = &sch
	case len(data) == 4 && string(data) == "true":
		nw.Allows = true
	case len(data) == 5 && string(data) == "false":
		nw.Allows = false
	default:
		return errors.New("boolean or JSON schema expected")
	}
	*s = nw
	return nil
}

func (s JSONSchemaPropsOrStringArray) MarshalJSON() ([]byte, error) {
	if len(s.Property) > 0 {
		return json.Marshal(s.Property)
	}
	if s.Schema != nil {
		return json.Marshal(s.Schema)
	}
	return []byte("null"), nil
}

func (s *JSONSchemaPropsOrStringArray) UnmarshalJSON(data []byte) error {
	var first byte
	if len(data) > 1 {
		first = data[0]
	}
	var nw JSONSchemaPropsOrStringArray
	if first == '{' {
		var sch JSONSchemaProps
		if err := json.Unmarshal(data, &sch); err != nil {
			return err
		}
		nw.Schema = &sch
	}
	if first == '[' {
		if err := json.Unmarshal(data, &nw.Property); err != nil {
			return err
		}
	}
	*s = nw
	return nil
}

func (s JSONSchemaPropsOrArray) MarshalJSON() ([]byte, error) {
	if len(s.JSONSchemas) > 0 {
		return json.Marshal(s.JSONSchemas)
	}
	return json.Marshal(s.Schema)
}

func (s *JSONSchemaPropsOrArray) UnmarshalJSON(data []byte) error {
	var nw JSONSchemaPropsOrArray
	var first byte
	if len(data) > 1 {
		first = data[0]
	}
	if first == '{' {
		var sch JSONSchemaProps
		if err := json.Unmarshal(data, &sch); err != nil {
			return err
		}
		nw.Schema = &sch
	}
	if first == '[' {
		if err := json.Unmarshal(data, &nw.JSONSchemas); err != nil {
			return err
		}
	}
	*s = nw
	return nil
}

func (s JSONSchemaProps) MarshalJSON() ([]byte, error) {
	type JSONSchemaPropsAlias JSONSchemaProps

	if len(s.SpecExtensions) == 0 {
		return json.Marshal(JSONSchemaPropsAlias(s))
	}
	all := make(map[string]interface{})
	bytes, err := json.Marshal(JSONSchemaPropsAlias(s))
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bytes, &all); err != nil {
		return nil, err
	}
	for x, val := range s.SpecExtensions {
		all[x] = val
	}
	return json.Marshal(all)
}

func (s *JSONSchemaProps) UnmarshalJSON(data []byte) error {
	type JSONSchemaPropsAddressAlias *JSONSchemaProps

	all := make(map[string]interface{})
	if err := json.Unmarshal(data, &all); err != nil {
		return err
	}
	extensions := make(map[string]interface{})
	for k, v := range all {
		if !strings.HasPrefix(k, "x-") || strings.HasPrefix(k, "x-kubernetes") {
			continue
		}
		extensions[k] = v
	}
	if err := json.Unmarshal(data, JSONSchemaPropsAddressAlias(s)); err != nil {
		return err
	}
	s.SpecExtensions = extensions
	return nil
}

func (s JSON) MarshalJSON() ([]byte, error) {
	if len(s.Raw) > 0 {
		return s.Raw, nil
	}
	return []byte("null"), nil

}

func (s *JSON) UnmarshalJSON(data []byte) error {
	if len(data) > 0 && string(data) != "null" {
		s.Raw = data
	}
	return nil
}
