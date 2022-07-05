package service

import (
	"sut-product-go/lib/pkg/db"
)

type Service struct {
	H db.Handler
}

func NewService(H db.Handler) *Service {
	return &Service{
		H,
	}
}
