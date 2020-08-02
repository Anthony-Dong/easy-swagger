package main

import (
	"encoding/json"
	"fmt"
	"github.com/anthony-dong/easy-swagger/swagger"
	"log"
	"net/http"
	"time"
)

type UserInfoResponse struct {
	Id           uint64            `json:"id" desc:"用户ID"`
	Name         string            `json:"name" desc:"用户名"`
	Birthday     int64             `json:"birthday" desc:"生日"`
	Hobbies      []Hobby           `json:"hobbies" desc:"喜爱"`
}
type Hobby struct {
	Name string `json:"name" desc:"名称"`
}
type UserInfoRequest struct {
	Id   uint64 `json:"id" desc:"用户ID"`
	Name string `json:"name" desc:"用户名"`
}

func main() {
	// 1、初始化swagger
	sr := swagger.New(swagger.ApiTitle("rest-api"),
		swagger.ApiServerAddress("www.google.com.cn"),
		swagger.ApiDesc("用户支付服务"),
		swagger.ApiHost(":8888"),
		swagger.ApiContact("574986060@qq.com"),
	)
	// 2、初始化-api
	api := sr.NewApi(swagger.ApiTag("user", "用户基本信息"),
		swagger.ApiPath("/user"))

	// 3、添加handler
	api.ApiOperation(http.MethodPost, "/info",
		swagger.ApiName("用户详情"),
		swagger.ApiDetail("用户详情详细描述"),
		swagger.ApiJsonParams(new(UserInfoRequest), swagger.ApiParameterName("请求参数")),
		swagger.ApiJsonResponse(new(UserInfoResponse)),
	)

	// 4、暴露api
	sr.ExportDefaultHttpHandler()

	// 5、业务接口
	http.HandleFunc("/user/info", func(writer http.ResponseWriter, request *http.Request) {
		req := new(UserInfoRequest)
		err := json.NewDecoder(request.Body).Decode(req)
		if err != nil {
			writer.WriteHeader(http.StatusUnsupportedMediaType)
			return
		}
		fmt.Printf("%+v\n", req)
		response := &UserInfoResponse{
			Id:       req.Id,
			Name:     req.Name,
			Birthday: time.Now().Unix(),
			Hobbies:  []Hobby{{"basketball"}},
		}
		writer.Header().Add("content-type", "application/json")
		json.NewEncoder(writer).Encode(response)
	})
	log.Fatal(http.ListenAndServe(":8888", nil))
}
