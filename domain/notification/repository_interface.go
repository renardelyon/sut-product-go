package notification

import notifpb "sut-product-go/pb/notification"

type NotificationClientInterface interface {
	UpdateNotificationByUserId(statusQtyMap map[string]*notifpb.StatusQty) (*notifpb.UpdateNotificationResponse, error)
}
