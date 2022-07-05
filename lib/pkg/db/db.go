package db

import (
	"log"
	"sut-product-go/domain/product/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Handler struct {
	DB *gorm.DB
}

func Init(url string) Handler {
	db, err := gorm.Open(postgres.Open(url), &gorm.Config{})
	if err != nil {
		log.Fatalln(err.Error())
	}

	models := []interface{}{&model.Product{}, &model.UserProduct{}}
	err = db.AutoMigrate(models...)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return Handler{db}
}
