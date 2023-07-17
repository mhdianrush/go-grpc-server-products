package services

import (
	"context"
	"go-gRPC-server-products/pb/pagination"
	productpb "go-gRPC-server-products/pb/product"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

type ProductService struct {
	productpb.UnimplementedProductServiceServer
	DB *gorm.DB
}

var logger = logrus.New()

func (p *ProductService) GetProducts(context.Context, *productpb.Empty) (*productpb.Products, error) {
	var products []*productpb.Product

	rows, err := p.DB.Table(
		`products as p`,
	).Joins(
		`left join categories as c on c.id = p.category_id`,
	).Select(
		`p.id`, `p.name`, `p.price`, `p.stock`, `c.id as category_id`, `c.name as category_name`,
	).Rows()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	defer rows.Close()

	for rows.Next() {
		var product productpb.Product
		var category productpb.Category

		err = rows.Scan(
			&product.Id,
			&product.Name,
			&product.Price,
			&product.Stock,
			&category.Id,
			&category.Name,
		)
		if err != nil {
			logger.Printf("Failed to Query Row Data %v", err.Error())
		}
		product.Category = &category
		products = append(products, &product)
	}
	response := &productpb.Products{
		Pagination: &pagination.Pagination{
			Total:       2,
			PerPage:     1,
			CurrentPage: 1,
			LastPage:    1,
		},
		// random input pagination temporary
		Data: products,
	}
	return response, nil
}
