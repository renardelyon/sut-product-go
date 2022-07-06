package service

import (
	"context"
	"net/http"
	"sut-product-go/domain/product/model"
	productpb "sut-product-go/pb/product"
	"time"

	"github.com/google/uuid"
)

func (s *Service) SaveRequestedGift(ctx context.Context, reqSave *productpb.SaveRequestedGiftRequest) (*productpb.SaveRequestedGiftResponse, error) {
	var product model.Product

	productId := uuid.NewString()
	timeNow := time.Now().Local().UTC().String()

	if result := s.H.DB.Where(&model.Product{Name: reqSave.Productname}).First(&product); result.Error != nil {
		res := s.H.DB.Create(&model.Product{
			Id:   productId,
			Name: reqSave.Productname,
			Url:  reqSave.Url,
		})
		if res.Error != nil {
			return &productpb.SaveRequestedGiftResponse{
				Status: http.StatusInternalServerError,
				Error:  res.Error.Error(),
			}, nil
		}
	}

	res := s.H.DB.Create(&model.UserProduct{
		AdminId:     reqSave.AdminId,
		Fullname:    reqSave.Fullname,
		Username:    reqSave.Username,
		UserId:      reqSave.UserId,
		ProductId:   productId,
		RequestDate: timeNow,
		Status:      model.Pending,
	})

	if res.Error != nil {
		return &productpb.SaveRequestedGiftResponse{
			Status: http.StatusInternalServerError,
			Error:  res.Error.Error(),
		}, nil
	}

	return &productpb.SaveRequestedGiftResponse{
		Status: http.StatusOK,
	}, nil
}
