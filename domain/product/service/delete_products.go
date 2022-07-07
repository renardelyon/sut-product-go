package service

import (
	"context"
	"net/http"
	"sut-product-go/domain/product/model"
	productpb "sut-product-go/pb/product"
	"sync"

	"gorm.io/gorm"
)

func (s *Service) DeleteProducts(ctx context.Context, reqDel *productpb.DeleteProductsRequest) (*productpb.DeleteProductsResponse, error) {
	chanTx := make(chan *gorm.DB)
	dBCollection := make([]*gorm.DB, len(reqDel.ProductIds))

	var errMessage = ""
	var wg sync.WaitGroup
	var mtx sync.Mutex

	wg.Add(len(reqDel.ProductIds))

	for _, productId := range reqDel.ProductIds {
		go func(productId string) {
			mtx.Lock()
			res := s.H.DB.Model(&model.UserProduct{}).
				Where(&model.UserProduct{ProductId: productId, UserId: reqDel.UserId}).
				Delete(&model.UserProduct{})
			mtx.Unlock()

			chanTx <- res

			wg.Done()
		}(productId)
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

		return &productpb.DeleteProductsResponse{
			Status: http.StatusInternalServerError,
			Error:  errMessage,
		}, nil
	}

	return &productpb.DeleteProductsResponse{
		Status: http.StatusOK,
	}, nil
}
