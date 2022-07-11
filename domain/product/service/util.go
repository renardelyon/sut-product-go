package service

import (
	"fmt"
	"log"
	"sut-product-go/domain/product/model"
	productpb "sut-product-go/pb/product"
	"sync"
	"time"

	"github.com/google/uuid"
)

type ProductInfo struct {
	index int
	name  string
	url   string
}

type productIdAndName map[string]string

func (s *Service) generateProductIndexes(productInfos []*productpb.ProductInfo) <-chan ProductInfo {
	chanOut := make(chan ProductInfo)
	go func() {
		for index, productInfo := range productInfos {
			chanOut <- ProductInfo{
				index: index,
				name:  productInfo.Productname,
				url:   productInfo.Url,
			}
		}

		close(chanOut)
	}()

	return chanOut
}

func (s *Service) insertProductInfoIntoProductTable(chanIn <-chan ProductInfo, products []model.Product, numOfWorkers int) <-chan productIdAndName {
	chanOut := make(chan productIdAndName)
	chanIsDataInDB := make(chan bool)
	var product model.Product
	var wg sync.WaitGroup

	wg.Add(numOfWorkers)

	go func() {
		for workerIndex := 0; workerIndex < numOfWorkers; workerIndex++ {
			go func() {
				for productInfo := range chanIn {
					productId := uuid.NewString()
					if result := s.H.DB.Where(&model.Product{Name: productInfo.name}).First(&product); result.Error != nil {
						products[productInfo.index] = model.Product{
							Id:   productId,
							Name: productInfo.name,
							Url:  productInfo.url,
						}
						chanOut <- productIdAndName{
							"productId":   productId,
							"productName": productInfo.name,
						}
						chanIsDataInDB <- false
					} else {
						chanOut <- productIdAndName{
							"productId":   product.Id,
							"productName": product.Name,
						}
						chanIsDataInDB <- true
					}
				}

				wg.Done()
			}()
		}
	}()

	go func() {
		var dataInDBQty int

		for isDataInDb := range chanIsDataInDB {
			if isDataInDb {
				dataInDBQty++
			}
		}

		if dataInDBQty <= 0 {
			res := s.H.DB.Create(&products)
			if res.Error != nil {
				log.Println(res.Error)
			}
		}
	}()

	go func() {
		wg.Wait()
		close(chanOut)
		close(chanIsDataInDB)
	}()

	return chanOut
}

func (s *Service) insertProductIdIntoUserProductTable(chanIn <-chan productIdAndName, tempProductInfo Temp) {
	var userProducts = make([]model.UserProduct, 0)
	timeNow := time.Now().Local().UTC().String()

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		for productIdAndName := range chanIn {
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
			`, productIdAndName["productName"], tempProductInfo.UserId)

			if res := s.H.DB.Raw(query).Scan(&products); res.Error != nil {
				log.Println(res.Error)
			}

			if len(products) < 1 {
				userProducts = append(userProducts, model.UserProduct{
					AdminId:     tempProductInfo.AdminId,
					Fullname:    tempProductInfo.Fullname,
					Username:    tempProductInfo.Username,
					UserId:      tempProductInfo.UserId,
					ProductId:   productIdAndName["productId"],
					RequestDate: timeNow,
					Status:      model.Pending,
				})
			}
		}

		wg.Done()
	}()

	wg.Wait()
	if len(userProducts) > 0 {
		res := s.H.DB.Create(&userProducts)
		if res.Error != nil {
			log.Fatalln(res.Error)
		}
	}
}
