# Easy-Swagger

## 功能

- Go-Web项目快速介入Swagger
- 与Web项目节藕，不再是重度依赖
- 目前只支持 json+rest 方式，不允许 path模式匹配，类似于Java的`@RequestBody`一样
- 与gin框架完美整合，未完成

## 快速引入

普通模式

```shell
go get -u github.com/anthony-dong/easy-swagger
```

go mod 模式

```go
// go mod 
import "github.com/anthony-dong/easy-swagger"
```

## 快速开始

[示例代码](./example/test_test.go)

```go
func main(){
    // 1、初始化swagger
    sr := swagger.New(
      swagger.ApiTitle("rest-api"),
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
      swagger.ApiJsonParams(new(UserInfoRequest)),
      swagger.ApiJsonResponse(new(UserInfoResponse)),
    )

    // 4、暴露api
    sr.ExportDefaultHttpHandler()

    // 5、编写 /user/info 接口

    //6、启动
    http.ListenAndServe(":8888", nil)
}
```

启动程序，访问`http://localhost:8888/swagger-ui/index.html`

启动如下图展示：

<img src="https://tyut.oss-accelerate.aliyuncs.com/image/2020-80-88/22ca1701-d8dc-4ec0-b608-954ce019e3c7.png" alt="image-20200802194524790" style="zoom:40%;" />

接口详情：

<img src="https://tyut.oss-accelerate.aliyuncs.com/image/2020-80-88/e09af5e5-86fe-438b-b6d2-f4091ede470a.png" alt="image-20200802194923015" style="zoom:50%;" />

model详情:

<img src="https://tyut.oss-accelerate.aliyuncs.com/image/2020-80-88/acfc4a86-de3d-49d4-8046-7700e66bf605.png" alt="image-20200802195057715" style="zoom:50%;" />

