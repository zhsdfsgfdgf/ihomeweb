package main

import (
	"sss/PostAuthUser/handler"
	"sss/PostAuthUser/subscriber"

	"github.com/micro/go-grpc"
	"github.com/micro/go-log"
	"github.com/micro/go-micro"

	example "sss/PostAuthUser/proto/example"
)

func main() {
	// New Service
	service := grpc.NewService(
		micro.Name("go.micro.srv.PostAuthUser"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	example.RegisterExampleHandler(service.Server(), new(handler.Example))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.PostAuthUser", service.Server(), new(subscriber.Example))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.PostAuthUser", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
