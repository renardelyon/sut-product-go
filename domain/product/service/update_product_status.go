package service

import (
	"context"
	"net/http"
	"strings"
	"sut-product-go/domain/product/model"
	productpb "sut-product-go/pb/product"
	"sync"

	"gorm.io/gorm"
)

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

	wg.Add(50)

	for _, userProduct := range reqUpdate.UserProducts {
		go func(up *productpb.UserAndProduct) {
			mtx.Lock()
			res := s.H.DB.Model(&model.UserProduct{}).
				Where(&model.UserProduct{ProductId: up.ProductId, UserId: up.UserId}).
				Update("status", strings.ToLower(reqUpdate.Status.String()))
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

	//TODO: Notification service

	return &productpb.UpdateProductStatusResponse{
		Status: http.StatusOK,
	}, nil
}
