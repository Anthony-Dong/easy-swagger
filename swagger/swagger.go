// Copyright 2016 HenryLee. All Rights Reserved.
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

// Package swagger struct definition
package swagger

import (
	"reflect"
)

// Version show the current swagger version
const Version = "2.0"

type (
	// Swagger object
	Swagger struct {
		Version  string                       `json:"swagger"`
		Info     *Info                        `json:"info"`
		Host     string                       `json:"host"`
		BasePath string                       `json:"basePath"`
		Tags     []*Tag                       `json:"tags"`
		Schemes  []string                     `json:"schemes"`
		Paths    map[string]map[string]*Opera `json:"paths,omitempty"` // {"prefix":{"method":{...}}}
		//SecurityDefinitions map[string]map[string]interface{} `json:"securityDefinitions,omitempty"`
		Definitions map[string]*Definition `json:"definitions,omitempty"`
		//ExternalDocs map[string]string      `json:"externalDocs,omitempty"`
	}
	// Info object
	Info struct {
		Title          string   `json:"title"`
		ApiVersion     string   `json:"version"`
		Description    string   `json:"description"`
		Contact        *Contact `json:"contact"`
		TermsOfService string   `json:"termsOfService"`
		License        *License `json:"license,omitempty"`
	}
	// Contact object
	Contact struct {
		Email string `json:"email,omitempty"`
	}
	// License object
	License struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	// Tag object
	Tag struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	// Opera object
	Opera struct {
		Tags        []string              `json:"tags"`
		Summary     string                `json:"summary"`
		Description string                `json:"description"`
		OperationId string                `json:"operationId"`
		Consumes    []string              `json:"consumes,omitempty"`
		Produces    []string              `json:"produces,omitempty"`
		Parameters  []*Parameter          `json:"parameters,omitempty"`
		Responses   map[string]*Resp      `json:"responses"` // {"httpcode":resp}
		Security    []map[string][]string `json:"security,omitempty"`
	}
	// Parameter object
	Parameter struct {
		In               string      `json:"in"` // the position of the parameter
		Name             string      `json:"name"`
		Description      string      `json:"description"`
		Required         bool        `json:"required"`
		Type             string      `json:"type,omitempty"` // "array"|"integer"|"object"
		Items            *Items      `json:"items,omitempty"`
		Schema           *Schema     `json:"schema,omitempty"`
		CollectionFormat string      `json:"collectionFormat,omitempty"` // "multi"
		Format           string      `json:"format,omitempty"`           // "int64"
		Default          interface{} `json:"default,omitempty"`
	}
	// Items object
	Items struct {
		Ref   string `json:"$ref,omitempty"`
		Type  string `json:"type,omitempty"`  // "string"
		Items *Items `json:"items,omitempty"` //  子类型
		//Format string `json:"format,omitempty"` // "int64"
		//Enum    interface{} `json:"enum,omitempty"`  // slice
		//Default interface{} `json:"default,omitempty"`
	}
	// Schema object
	Schema struct {
		Ref                  string            `json:"$ref,omitempty"`
		Type                 string            `json:"type,omitempty"` // "array"|"integer"|"object"
		Items                *Items            `json:"items,omitempty"`
		Description          string            `json:"description,omitempty"`
		AdditionalProperties map[string]string `json:"additionalProperties,omitempty"`
	}
	// Resp object
	Resp struct {
		Schema      *Schema `json:"schema"`
		Description string  `json:"description"`
	}
	// Definition object
	Definition struct {
		Type       string               `json:"type,omitempty"` // "object"
		Properties map[string]*Property `json:"properties,omitempty"`
		//Xml        *Xml                 `json:"xml,omitempty"`
	}
	// Property object
	Property struct {
		Ref         string `json:"$ref,omitempty"`        // 可以引用到其他地方
		Type        string `json:"type,omitempty"`        // "array"|"integer"|"object"|"string"
		Items       *Items `json:"items,omitempty"`       // 数组需要子类型
		Format      string `json:"format,omitempty"`      // "int64"，真实类型
		Description string `json:"description,omitempty"` // 描述
		//Enum        []string `json:"enum,omitempty"`        // 枚举，展示效果为数组，但是目前go不支持
		//Example            interface{}       `json:"example,omitempty"`              // 基本类型展示，其他类型不展示
		Default            interface{}       `json:"default,omitempty"`              // 默认值，go不支持
		MapValueProperties *MapValueProperty `json:"additionalProperties,omitempty"` // map 的value
	}

	MapValueProperty struct {
		Type string `json:"type,omitempty"` // 类型 "array"|"integer"|"object"|"string"
		Ref  string `json:"$ref,omitempty"` // 引用
		//Format             string            `json:"format,omitempty"`               // format 就是真正类型，由于map的key必须是string，所以只记录value的属性
		MapValueProperties *MapValueProperty `json:"additionalProperties,omitempty"` // map 的value
	}
	// Xml object
	Xml struct {
		Name    string `json:"name"`
		Wrapped bool   `json:"wrapped,omitempty"`
	}
)

// CommonMIMETypes common MIME types
var CommonMIMETypes = []string{
	"application/json",
	"application/javascript",
	"application/xml",
	"application/x-www-form-urlencoded",
	"application/protobuf",
	"application/msgpack",
	"text/html",
	"text/plain",
	"multipart/form-data",
	"application/octet-stream",
}

// github.com/mcuadros/go-jsonschema-generator
var mapping = map[reflect.Kind]string{
	reflect.Bool:      "bool",
	reflect.Int:       "integer",
	reflect.Int8:      "integer",
	reflect.Int16:     "integer",
	reflect.Int32:     "integer",
	reflect.Int64:     "integer",
	reflect.Uint:      "integer",
	reflect.Uint8:     "integer",
	reflect.Uint16:    "integer",
	reflect.Uint32:    "integer",
	reflect.Uint64:    "integer",
	reflect.Float32:   "number",
	reflect.Float64:   "number",
	reflect.String:    "string",
	reflect.Interface: "object",
}

// ParamType type of the parameter value passed in
func ParamType(value interface{}) string {
	if value == nil {
		return ""
	}
	rv, ok := value.(reflect.Type)
	if !ok {
		rv = reflect.TypeOf(value)
	}
	if rv.Kind() == reflect.Ptr {
		rv = rv.Elem()
	}
	tn := rv.String()
	if tn == "multipart.FileHeader" || tn == "[]*multipart.FileHeader" || tn == "[]multipart.FileHeader" {
		return "file"
	}
	return mapping[rv.Kind()]
}
