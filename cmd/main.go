package main

import (
	"log"
	"sut-product-go/application"
	"sut-product-go/config"
	"sut-product-go/domain/product/service"
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

	s := service.NewService(app.DbClients)

	productpb.RegisterProductServiceServer(app.GrpcServer, s)

	err = app.Run(&c)
	if err != nil {
		log.Fatalln(err.Error())
	}
}
