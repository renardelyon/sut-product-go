package grpc

import (
	"context"
	notifpb "sut-product-go/pb/notification"
)

type repo struct {
	notifClient notifpb.NotificationServiceClient
}

func NewGrpcRepo(notifClient notifpb.NotificationServiceClient) *repo {
	return &repo{
		notifClient: notifClient,
	}
}

func (r *repo) UpdateNotificationByUserId(statusQtyMap map[string]*notifpb.StatusQty) (*notifpb.UpdateNotificationResponse, error) {
	req := &notifpb.UpdateNotificationRequest{
		StatusQtyMap: statusQtyMap,
	}
	return r.notifClient.UpdateNotificationByUserId(context.Background(), req)
}
