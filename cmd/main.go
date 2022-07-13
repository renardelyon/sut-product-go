package main

import (
	"log"
	"sut-product-go/application"
	"sut-product-go/config"
	notifGrpc "sut-product-go/domain/notification/grpc"
	"sut-product-go/domain/product/service"
	notifpb "sut-product-go/pb/notification"
	productpb "sut-product-go/pb/product"
)

func main() {
	c, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config: ", err.Error())
	}

	app, err := application.Setup(&c)
	if err != nil {
		log.Fatalln("Failed at application setup: ", err.Error())
	}

	notifClient := notifpb.NewNotificationServiceClient(app.GrpcClients["notification-service"])
	notifRepo := notifGrpc.NewGrpcRepo(notifClient)

	s := service.NewService(app.DbClients, notifRepo)

	productpb.RegisterProductServiceServer(app.GrpcServer, s)

	err = app.Run(&c)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
