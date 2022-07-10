package service

import (
	"context"
	"net/http"
	"sut-product-go/domain/product/model"
	productpb "sut-product-go/pb/product"
)

func (s *Service) SaveRequestedGiftBulk(ctx context.Context, reqSave *productpb.SaveRequestedGiftBulkRequest) (*productpb.SaveRequestedGiftBulkResponse, error) {
	products := make([]model.Product, len(reqSave.ProductInfo))

	chanIndex := s.generateProductIndexes(reqSave.ProductInfo)

	chanProductId := s.insertProductInfoIntoProductTable(chanIndex, products, 50)

	s.insertProductIdIntoUserProductTable(chanProductId, reqSave)

	return &productpb.SaveRequestedGiftBulkResponse{
		Status: http.StatusOK,
	}, nil
}
