package swagger

import (
	"encoding/json"
	"fmt"
	"github.com/anthony-dong/easy-swagger/html"
	"github.com/anthony-dong/easy-swagger/util"
	"net/http"
	"strings"
)

//---------------------------- init swagger -----------------------------------
type SwInfo struct {
	Version    string // api版本
	Host       string // host
	BasePath   string // 基本路径，默认是 "/"
	Title      string // 描述
	Desc       string // 描述
	Email      string // 联系方式
	ServerName string // server address
}

func ApiHost(host string) InfoOp {
	return func(swagger *SwInfo) {
		host, port := util.GetHost(host)
		swagger.Host = fmt.Sprintf("%s:%s", host, port)
	}
}
func ApiTitle(title string) InfoOp {
	return func(swagger *SwInfo) {
		swagger.Title = title
	}
}
func ApiDesc(desc string) InfoOp {
	return func(swagger *SwInfo) {
		swagger.Desc = desc
	}
}
func ApiContact(email string) InfoOp {
	return func(swagger *SwInfo) {
		swagger.Email = email
	}
}
func ApiServerAddress(serverName string) InfoOp {
	return func(swagger *SwInfo) {
		swagger.ServerName = serverName
	}
}

type InfoOp func(swagger *SwInfo)

func New(op ...InfoOp) *Swagger {
	swagger := &Swagger{
		Version:     Version,
		BasePath:    "/",
		Tags:        []*Tag{},
		Schemes:     []string{"http"},
		Paths:       map[string]map[string]*Opera{},
		Definitions: map[string]*Definition{},
	}
	swInfo := new(SwInfo)
	for _, elem := range op {
		elem(swInfo)
	}
	if swInfo.Host != "" {
		swagger.Host = swInfo.Host
	}
	if swInfo.BasePath != "" {
		swagger.BasePath = swInfo.BasePath
	}
	swagger.Host = swInfo.Host
	swagger.Info = &Info{
		Title:       swInfo.Title,
		ApiVersion:  swInfo.Version,
		Description: swInfo.Desc,
		Contact: &Contact{
			Email: swInfo.Email,
		},
		License: &License{
			Name: "The Apache License",
			Url:  "https://opensource.org/licenses/MIT",
		},
		TermsOfService: swInfo.ServerName,
	}
	return swagger
}

// ---------------------------------- http --------------------------------------------
type HttpFunc func(http.ResponseWriter, *http.Request)

func (this *Swagger) ExportRestApi() (path string, fun HttpFunc) {
	return "/swagger.json", func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("content-type", "application/json")
		json.NewEncoder(writer).Encode(this)
	}
}

func (this *Swagger) ExportDefaultHttpHandler() {
	path, fun := this.ExportRestApi()
	http.HandleFunc(path, fun)
	htmls := html.SwaggerHtmls()
	for _, elem := range htmls {
		http.HandleFunc(elem.Path, elem.Handler)
	}
}

// ----------------------------------- init api  --------------------------

type API struct {
	sw   *Swagger
	tags []Tag
	path string
}

type NewApiFunc func(api *API)

func (this *Swagger) NewApi(ops ...NewApiFunc) *API {
	api := &API{
		sw:   this,
		tags: []Tag{},
	}
	for _, elem := range ops {
		if elem != nil {
			elem(api)
		}
	}
	for _, elem := range api.tags {
		this.Tags = append(this.Tags, &elem)
	}
	return api
}
func ApiTag(tag string, desc string) NewApiFunc {
	return func(api *API) {
		api.tags = append(api.tags, Tag{
			Name:        tag,
			Description: desc,
		})
	}
}
func ApiPath(path string) NewApiFunc {
	return func(api *API) {
		api.path = path
	}
}

func (this *API) ApiOperation(method, path string, pathFunc ...ApiOperation) {
	for _, tag := range this.tags {
		pathFunc = append(pathFunc, func(opera *SwOpera) {
			opera.op.Tags = append(opera.op.Tags, tag.Name)
		})
	}
	path = util.CleanPath(path)
	path = strings.Join([]string{this.path, path}, "/")
	this.sw.ApiOperation(method, util.CleanPath(path), pathFunc ...)
}

