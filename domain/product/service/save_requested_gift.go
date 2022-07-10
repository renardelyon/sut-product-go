package service

import (
	"context"
	"fmt"
	"log"
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

	var products []model.Product

	query := fmt.Sprintf(`
		select *
		from user_products up
		where product_id  in
		(select p.id
		from user_products up
		inner join products p
		on up.product_id = p.id
		where p."name" = '%s')
		and user_id = '%s'
	`, reqSave.Productname, reqSave.UserId)

	if res := s.H.DB.Raw(query).Scan(&products); res.Error != nil {
		log.Println(res.Error)
		return &productpb.SaveRequestedGiftResponse{
			Status: http.StatusInternalServerError,
			Error:  res.Error.Error(),
		}, nil
	}

	if len(products) < 1 {
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
	}

	return &productpb.SaveRequestedGiftResponse{
		Status: http.StatusOK,
	}, nil
}
