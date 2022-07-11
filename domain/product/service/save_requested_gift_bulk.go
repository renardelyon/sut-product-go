package service

import (
	"context"
	"net/http"
	"sut-product-go/domain/product/model"
	productpb "sut-product-go/pb/product"
)

type Temp struct {
	AdminId  string
	Fullname string
	Username string
	UserId   string
}

func (s *Service) SaveRequestedGiftBulk(ctx context.Context, reqSave *productpb.SaveRequestedGiftBulkRequest) (*productpb.SaveRequestedGiftBulkResponse, error) {
	temp := Temp{
		AdminId:  reqSave.AdminId,
		Fullname: reqSave.Fullname,
		Username: reqSave.Username,
		UserId:   reqSave.UserId,
	}

	products := make([]model.Product, len(reqSave.ProductInfo))

	chanIndex := s.generateProductIndexes(reqSave.ProductInfo)

	chanProductId := s.insertProductInfoIntoProductTable(chanIndex, products, 50)

	s.insertProductIdIntoUserProductTable(chanProductId, temp)

	return &productpb.SaveRequestedGiftBulkResponse{
		Status: http.StatusOK,
	}, nil
}
