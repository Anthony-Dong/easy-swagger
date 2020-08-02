package html

import (
	"fmt"
	"github.com/anthony-dong/easy-swagger/util"
	"net/http"
)

var (
	swaggerHtml = make([]SwaggerHtml, 0)
)

func SwaggerHtmls() []SwaggerHtml {
	return swaggerHtml
}

func init() {
	AddSwagger()
}

const (
	js       = "js"
	css      = "css"
	html     = "html"
	byteType = "byte"
	jsonType = "json"
)

func addSwagger(path string, _type string) func(writer http.ResponseWriter, request *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		body, _ := Asset(path)
		switch _type {
		case js:
			writer.Header().Set("content-type", "application/javascript")
		case css:
			writer.Header().Set("content-type", "text/css")
		case jsonType:
			writer.Header().Set("content-type", "application/json")
		}
		if _type == byteType {
			fmt.Fprint(writer, body)
			return
		}
		fmt.Fprint(writer, string(body))
	}
}

func addSwaggerRouter(path string, _type string) {
	sg := SwaggerHtml{
		Method:  http.MethodGet,
		Path:    util.CleanPath("/" + path),
		Handler: addSwagger(path, _type),
	}
	swaggerHtml = append(swaggerHtml, sg)
}

type SwaggerHtml struct {
	Method  string
	Path    string
	Handler func(http.ResponseWriter, *http.Request)
}

func AddSwagger() {
	addSwaggerRouter("swagger-ui/absolute-path.js", js)
	addSwaggerRouter("swagger-ui/favicon-16x16.png", byteType)
	addSwaggerRouter("swagger-ui/favicon-32x32.png", byteType)
	addSwaggerRouter("swagger-ui/index.html", html)
	addSwaggerRouter("swagger-ui/index.js", js)
	addSwaggerRouter("swagger-ui/oauth2-redirect.html", html)
	addSwaggerRouter("swagger-ui/package.json", jsonType)
	addSwaggerRouter("swagger-ui/swagger-ui-bundle.js", js)
	addSwaggerRouter("swagger-ui/swagger-ui-bundle.js.map", html)
	addSwaggerRouter("swagger-ui/swagger-ui-standalone-preset.js", js)
	addSwaggerRouter("swagger-ui/swagger-ui-standalone-preset.js.map", html)
	addSwaggerRouter("swagger-ui/swagger-ui.css", css)
	addSwaggerRouter("swagger-ui/swagger-ui.css.map", html)
	addSwaggerRouter("swagger-ui/swagger-ui.js", js)
	addSwaggerRouter("swagger-ui/swagger-ui.js.map", html)
}
