package service

import (
	"context"
	"log"
	"net/http"
	"strings"
	"sut-product-go/domain/product/model"
	notifpb "sut-product-go/pb/notification"
	productpb "sut-product-go/pb/product"
	"sync"

	"gorm.io/gorm"
)

type statusQtyByUserId map[string]*notifpb.StatusQty

func (s *Service) UpdateProductStatus(ctx context.Context, reqUpdate *productpb.UpdateProductStatusRequest) (*productpb.UpdateProductStatusResponse, error) {
	if reqUpdate.Role.String() != productpb.Role_ADMIN.String() {
		return &productpb.UpdateProductStatusResponse{
			Status: http.StatusUnauthorized,
			Error:  "Unauthorized user",
		}, nil
	}

	chanTx := make(chan *gorm.DB)
	dBCollection := make([]*gorm.DB, len(reqUpdate.UserProducts))

	var errMessage = ""
	var wg sync.WaitGroup
	var mtx sync.Mutex

	wg.Add(len(reqUpdate.UserProducts))

	statusQtyMap := make(statusQtyByUserId)

	for _, userProduct := range reqUpdate.UserProducts {
		go func(up *productpb.UserAndProduct) {
			mtx.Lock()
			res := s.H.DB.Model(&model.UserProduct{}).
				Where(&model.UserProduct{ProductId: up.ProductId, UserId: up.UserId}).
				Update("status", strings.ToLower(reqUpdate.Status.String()))

			if _, ok := statusQtyMap[up.UserId]; !ok {
				statusQtyMap[up.UserId] = &notifpb.StatusQty{
					Status:   reqUpdate.Status.String(),
					Quantity: 1,
				}
			} else {
				statusQtyMap[up.UserId].Quantity = statusQtyMap[up.UserId].Quantity + 1
			}
			mtx.Unlock()

			chanTx <- res

			wg.Done()
		}(userProduct)
	}

	go func() {
		wg.Wait()
		close(chanTx)
	}()

	for tx := range chanTx {
		dBCollection = append(dBCollection, tx)
		if tx.Error != nil {
			errMessage = tx.Error.Error()
		}
	}

	if errMessage != "" {
		for _, tx := range dBCollection {
			tx.Rollback()
		}

		return &productpb.UpdateProductStatusResponse{
			Status: http.StatusInternalServerError,
			Error:  errMessage,
		}, nil
	}

	res, _ := s.NotifRepo.UpdateNotificationByUserId(statusQtyMap)
	if res.Error != "" {
		log.Println(res.Error)
		return &productpb.UpdateProductStatusResponse{
			Status: http.StatusInternalServerError,
			Error:  res.Error,
		}, nil
	}

	return &productpb.UpdateProductStatusResponse{
		Status: http.StatusOK,
	}, nil
}
