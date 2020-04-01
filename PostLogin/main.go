package main

import (
	"sss/PostLogin/handler"
	"sss/PostLogin/subscriber"

	"github.com/micro/go-grpc"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"

	example "sss/PostLogin/proto/example"
	//redis缓存操作与支持驱动

	_ "github.com/astaxie/beego/cache/redis"
	_ "github.com/garyburd/redigo/redis"
	_ "github.com/gomodule/redigo/redis"
)

func main() {
	// New Service
	service := grpc.NewService(
		micro.Name("go.micro.srv.PostLogin"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	example.RegisterExampleHandler(service.Server(), new(handler.Example))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.PostLogin", service.Server(), new(subscriber.Example))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.PostLogin", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
