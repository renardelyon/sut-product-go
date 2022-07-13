package service

import (
	"sut-product-go/domain/notification"
	"sut-product-go/lib/pkg/db"
)

type Service struct {
	H         db.Handler
	NotifRepo notification.NotificationClientInterface
}

func NewService(H db.Handler, notifRepo notification.NotificationClientInterface) *Service {
	return &Service{
		H,
		notifRepo,
	}
}
