package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	productpb "sut-product-go/pb/product"
)

type tempProductInfo struct {
	AdminId     string
	Fullname    string
	Username    string
	UserId      string
	ProductId   string
	RequestDate string
	Status      string
	Name        string
	Url         string
}

func (s *Service) GetAllProducts(ctx context.Context, reqProduct *productpb.GetAllProductsRequest) (*productpb.GetAllProductsResponse, error) {
	var productsTemp []tempProductInfo
	var query string

	if reqProduct.Role.String() == productpb.Role_ADMIN.String() {
		query = fmt.Sprintf(`
			select up.admin_id, up.fullname, up.username, up.user_id, up.request_date, up.status, p.name, p.url 
			from user_products as up 
			inner join products as p 
			on up.product_id = p.id
			where up.admin_id = '%s'
			and up.status = '%s'
		`, reqProduct.Id, strings.ToLower(reqProduct.Status.String()))
	}

	if reqProduct.Role.String() == productpb.Role_USER.String() {
		query = fmt.Sprintf(`
			select up.admin_id, up.fullname, up.username, up.user_id, up.request_date, up.status, p.name, p.url 
			from user_products as up 
			inner join products as p 
			on up.product_id = p.id
			where up.user_id = '%s'
			and up.status = '%s'
		`, reqProduct.Id, strings.ToLower(reqProduct.Status.String()))
	}

	res := s.H.DB.Raw(query).Scan(&productsTemp)
	if res.Error != nil {
		return &productpb.GetAllProductsResponse{
			Status: http.StatusInternalServerError,
			Error:  res.Error.Error(),
		}, nil
	}

	products := make([]*productpb.Product, 0)
	for _, productTmp := range productsTemp {
		var procStatus int
		switch strings.ToUpper(productTmp.Status) {
		case productpb.Status_ACCEPTED.String():
			procStatus = int(productpb.Status_ACCEPTED)
		case productpb.Status_PENDING.String():
			procStatus = int(productpb.Status_PENDING)
		case productpb.Status_REJECTED.String():
			procStatus = int(productpb.Status_REJECTED)
		}

		products = append(products, &productpb.Product{
			AdminId:     productTmp.AdminId,
			Fullname:    productTmp.Fullname,
			Username:    productTmp.Username,
			UserId:      productTmp.UserId,
			RequestDate: productTmp.RequestDate,
			Status:      productpb.Status(procStatus),
			ProductName: productTmp.Name,
			Url:         productTmp.Url,
		})
	}

	return &productpb.GetAllProductsResponse{
		Status:   http.StatusOK,
		Products: products,
	}, nil
}
