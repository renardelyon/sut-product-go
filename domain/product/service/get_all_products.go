package service

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	productpb "sut-product-go/pb/product"
)

func (s *Service) GetAllProducts(ctx context.Context, reqProduct *productpb.GetAllProductsRequest) (*productpb.GetAllProductsResponse, error) {
	var products []*productpb.Product
	var query string

	if reqProduct.Role.String() == productpb.Role_ADMIN.String() {
		query = fmt.Sprintf(`
			select up.admin_id, up.fullname, up.username, up.user_id, up.request_date, up.status, p.name, p.url 
			from user_products as up 
			inner join products as p 
			on up.product_id = p.id
			where up.admin_id = %s
			and up.status = %s
		`, reqProduct.Id, strings.ToLower(reqProduct.Status.String()))
	}

	if reqProduct.Role.String() == productpb.Role_USER.String() {
		query = fmt.Sprintf(`
			select up.admin_id, up.fullname, up.username, up.user_id, up.request_date, up.status, p.name, p.url 
			from user_products as up 
			inner join products as p 
			on up.product_id = p.id
			where up.user_id = %s
			and up.status = %s
		`, reqProduct.Id, strings.ToLower(reqProduct.Status.String()))
	}

	res := s.H.DB.Raw(query).Scan(&products)
	if res.Error != nil {
		return &productpb.GetAllProductsResponse{
			Status: http.StatusInternalServerError,
			Error:  res.Error.Error(),
		}, nil
	}

	return &productpb.GetAllProductsResponse{
		Status:   http.StatusOK,
		Products: products,
	}, nil
}
