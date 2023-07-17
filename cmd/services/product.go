package services

import (
	"context"
	paginationpb "go-gRPC-server-products/pb/pagination"
	productpb "go-gRPC-server-products/pb/product"

	"gorm.io/gorm"
)

type ProductService struct {
	productpb.UnimplementedProductServiceServer
	DB *gorm.DB
}

func (p *ProductService) GetProducts(context.Context, *productpb.Empty) (*productpb.Products, error) {
	products := &productpb.Products{
		Pagination: &paginationpb.Pagination{
			Total:       10,
			PerPage:     5,
			CurrentPage: 1,
			LastPage:    2,
		},
		Data: []*productpb.Product{
			{
				Id:    1,
				Name:  "Fortuner",
				Price: 500000000,
				Stock: 50,
				Category: &productpb.Category{
					Id:   1,
					Name: "Car",
				},
			},
			{
				Id:    2,
				Name:  "Nmax",
				Price: 20000000,
				Stock: 100,
				Category: &productpb.Category{
					Id:   2,
					Name: "Motor Cycle",
				},
			},
		},
	}
	return products, nil
}
