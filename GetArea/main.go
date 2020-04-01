package main

import (
	"sss/GetArea/handler"
	"sss/GetArea/subscriber"

	"github.com/micro/go-grpc"

	"github.com/micro/go-log"
	"github.com/micro/go-micro"

	example "sss/GetArea/proto/example"
)

func main() {
	// New Service
	//换成grpc
	service := grpc.NewService(
		micro.Name("go.micro.srv.GetArea"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	example.RegisterExampleHandler(service.Server(), new(handler.Example))

	// Register Struct as Subscriber
	micro.RegisterSubscriber("go.micro.srv.GetArea", service.Server(), new(subscriber.Example))

	// Register Function as Subscriber
	micro.RegisterSubscriber("go.micro.srv.GetArea", service.Server(), subscriber.Handler)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
