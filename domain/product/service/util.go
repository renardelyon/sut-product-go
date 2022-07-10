package service

import (
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

func (s *Service) insertProductInfoIntoProductTable(chanIn <-chan ProductInfo, products []model.Product, numOfWorkers int) <-chan string {
	chanOut := make(chan string)
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
						chanOut <- productId
					}
				}

				wg.Done()
			}()
		}
	}()

	go func() {
		wg.Wait()
		close(chanOut)
		res := s.H.DB.Create(&products)
		if res.Error != nil {
			log.Fatalln(res.Error)
		}
	}()

	return chanOut
}

func (s *Service) insertProductIdIntoUserProductTable(chanIn <-chan string, reqSave *productpb.SaveRequestedGiftBulkRequest) {
	var userProducts = make([]model.UserProduct, 0)
	var chanOut = make(chan []model.UserProduct)
	timeNow := time.Now().Local().UTC().String()

	go func() {
		for productId := range chanIn {
			userProducts = append(userProducts, model.UserProduct{
				AdminId:     reqSave.AdminId,
				Fullname:    reqSave.Fullname,
				Username:    reqSave.Username,
				UserId:      reqSave.UserId,
				ProductId:   productId,
				RequestDate: timeNow,
				Status:      model.Pending,
			})
		}

		chanOut <- userProducts
		close(chanOut)
	}()

	up := <-chanOut
	res := s.H.DB.Create(&up)
	if res.Error != nil {
		log.Fatalln(res.Error)
	}
}
