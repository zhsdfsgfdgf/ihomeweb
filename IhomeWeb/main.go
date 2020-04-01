package main

import (
	"net/http"

	"github.com/micro/go-log"

	"sss/IhomeWeb/handler"

	"github.com/julienschmidt/httprouter"

	"github.com/micro/go-web"

	_ "sss/IhomeWeb/models"
)

func main() {
	//创建web服务
	service := web.NewService(
		web.Name("go.micro.web.IhomeWeb"),
		web.Version("latest"),
		web.Address(":8990"),
	)

	//初始化服务
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}
	//映射静态页面
	rou := httprouter.New()
	rou.NotFound = http.FileServer(http.Dir("html"))

	//获取地区信息
	//httprouter中GET方法有两个参数,一是匹配路径,二是处理函数
	//其中Handle与HandlerFun类似,但是多了第三个参数用来传参
	rou.GET("/api/v1.0/areas", handler.GetArea)
	//下面两个目前并不实现服务
	//获取session
	rou.GET("/api/v1.0/session", handler.GetSession)
	//获取index
	rou.GET("/api/v1.0/house/index", handler.GetIndex)
	//获取图片验证码
	rou.GET("/api/v1.0/imagecode/:uuid", handler.GetImageCd)
	//获取短信验证码
	rou.GET("/api/v1.0/smscode/:mobile", handler.GetSmscd)
	//注册
	rou.POST("/api/v1.0/users", handler.Postreg)
	//登录
	rou.POST("/api/v1.0/sessions", handler.PostLogin)
	//退出登陆
	rou.DELETE("/api/v1.0/session", handler.DeleteSession)
	//请求用户基本信息
	rou.GET("/api/v1.0/user", handler.GetUserInfo)
	//上传头像 POST
	rou.POST("/api/v1.0/user/avatar", handler.PostAvatar)
	//检查用户实名认证,与请求用户基本信息请求的数据相同
	rou.GET("/api/v1.0/user/auth", handler.GetUserAuth)
	// 更新实名认证信息
	rou.POST("/api/v1.0/user/auth", handler.PostAuthUser)
	// 获取用户房源信息
	rou.GET("/api/v1.0/user/houses", handler.GetUserHouses)
	// 用户发布房源信息
	rou.POST("/api/v1.0/houses", handler.PostHouses)
	// 用户上传房源图片信息
	rou.POST("/api/v1.0/houses/:id/images", handler.PostHousesImage)
	// 获取房源具体信息
	rou.GET("/api/v1.0/houses/:id", handler.GetHouseInfo)
	// register html handler
	//注册服务
	service.Handle("/", rou) //没写这句，用的是下面的，一直404
	//后续陆续添加服务所以这个文件的这个地方会一直添加内容

	// register call handler
	//注册服务
	// service.HandleFunc("/example/call", handler.ExampleCall)

	// run service
	//运行服务
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
