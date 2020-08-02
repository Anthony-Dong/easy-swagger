package swagger

import (
	"github.com/anthony-dong/easy-swagger/logger"
	"reflect"
	"strings"
)

// -----------------init api handler-----------------------------
type ParamsType string

const (
	Query = ParamsType("query")
	Body  = ParamsType("body")
)

type SwOpera struct {
	op *Opera
	sw *Swagger
}
type ApiOperation func(*SwOpera)

type ApiParameter struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}
type ApiParameterFunc func(parameter *ApiParameter)

func ApiParameterName(name string) ApiParameterFunc {
	return func(parameter *ApiParameter) {
		parameter.Name = name
	}
}
func ApiDetail(detail string) ApiOperation {
	return func(opera *SwOpera) {
		opera.op.Description = detail
	}
}
func ApiName(name string) ApiOperation {
	return func(opera *SwOpera) {
		opera.op.Summary = name
	}
}

func ApiJsonParams(param interface{}, ops ...ApiParameterFunc) ApiOperation {
	return func(opera *SwOpera) {
		opera.op.Produces = []string{"application/json"}
		tp := reflect.TypeOf(param)
		if tp.Kind() == reflect.Ptr {
			tp = tp.Elem()
		}
		if tp.Kind() != reflect.Struct {
			panic("")
		}
		parameter := Parameter{
			In:       string(Body),
			Name:     "Request-Body",
			Required: true,
			Schema:   &Schema{Ref: GetRef(tp), Type: "object"},
		}
		AddProperties(reflect.TypeOf(param), opera.sw.Definitions, nil)
		// load custom
		param := new(ApiParameter)
		for _, elem := range ops {
			if elem != nil {
				elem(param)
			}
		}
		parameter.Required = param.Required
		parameter.Description = param.Description
		parameter.Name = param.Name
		if opera.op.Parameters == nil {
			opera.op.Parameters = []*Parameter{}
		}
		opera.op.Parameters = append(opera.op.Parameters, &parameter)
	}
}

func ApiJsonResponse(response interface{}) ApiOperation {
	return func(opera *SwOpera) {
		opera.op.Consumes = []string{"application/json"}
		tp := reflect.TypeOf(response)
		if tp.Kind() == reflect.Ptr {
			tp = tp.Elem()
		}
		if tp.Kind() != reflect.Struct {
			panic("struct must")
		}
		AddProperties(tp, opera.sw.Definitions, nil)
		if opera.op.Responses == nil {
			opera.op.Responses = map[string]*Resp{}
		}
		resp := &Resp{
			Description: "OK",
			Schema:      &Schema{Ref: GetRef(tp)},
		}
		opera.op.Responses["200"] = resp
	}
}
func (this *Swagger) ApiOperation(method, path string, pathFunc ...ApiOperation) {
	method = strings.ToLower(method)
	desc, isExist := this.Paths[path]
	if !isExist {
		desc = map[string]*Opera{}
		this.Paths[path] = desc
	} else {
		logger.FatalF("[HandlerPath] is exist %s path", path)
	}
	op, isExist := desc[method]
	if !isExist {
		op = &Opera{
			Tags:       []string{},
			Consumes:   CommonMIMETypes,
			Produces:   CommonMIMETypes,
			Parameters: []*Parameter{},
			Responses:  map[string]*Resp{},
			Security:   []map[string][]string{},
		}
		desc[method] = op
	} else {
		logger.FatalF("[HandlerPath] is exist %s-%s path-method", path, method)
	}
	op.OperationId = strings.Join([]string{method, path}, "-")
	opera := SwOpera{
		op: op,
		sw: this,
	}
	for _, elem := range pathFunc {
		if elem != nil {
			elem(&opera)
		}
	}
}
