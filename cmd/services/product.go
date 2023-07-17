package services

import (
	"context"
	"go-gRPC-server-products/cmd/helpers"
	paginationpb "go-gRPC-server-products/pb/pagination"
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

func (p *ProductService) GetProducts(ctx context.Context, pageParam *productpb.Page) (*productpb.Products, error) {
	var page int64 = 1
	if pageParam.GetPage() != 0 {
		page = pageParam.GetPage()
	}

	var pagination paginationpb.Pagination
	var products []*productpb.Product

	sql := p.DB.Table(
		`products as p`,
	).Joins(
		`left join categories as c on c.id = p.category_id`,
	).Select(
		`p.id`, `p.name`, `p.price`, `p.stock`, `c.id as category_id`, `c.name as category_name`,
	)

	offset, limit := helpers.Pagination(sql, page, &pagination)

	rows, err := sql.Offset(int(offset)).Limit(int(limit)).Rows()

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
			logger.Printf("Failed to Query Rows Data %v", err.Error())
		}
		product.Category = &category
		products = append(products, &product)
	}
	response := &productpb.Products{
		Pagination: &pagination,
		// random input pagination temporary
		Data: products,
	}
	return response, nil
}

func (p *ProductService) GetProduct(ctx context.Context, id *productpb.Id) (*productpb.Product, error) {
	row := p.DB.Table(
		`products as p`,
	).Joins(
		`left join categories as c on c.id = p.category_id`,
	).Select(
		`p.id`, `p.name`, `p.price`, `p.stock`, `c.id as category_id`, `c.name as category_name`,
	).Where(
		`p.id = ?`, id.GetId(),
	).Row()
	var product productpb.Product
	var category productpb.Category

	if err := row.Scan(
		&product.Id,
		&product.Name,
		&product.Price,
		&product.Stock,
		&category.Id,
		&category.Name,
	); err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	product.Category = &category

	return &product, nil
}

func (p *ProductService) CreateProduct(ctx context.Context, productData *productpb.Product) (*productpb.Id, error) {
	var Response productpb.Id

	err := p.DB.Transaction(func(tx *gorm.DB) error {
		category := productpb.Category{
			Id:   0,
			Name: productData.GetCategory().GetName(),
		}
		// check data, if data exist or not exist in db, then will be insert to category
		if err := tx.Table(`categories`).Where(`LCASE(name) = ?`, category.GetName()).FirstOrCreate(&category).Error; err != nil {
			return err
		}

		// anonymous struct
		product := struct {
			Id          uint64
			Name        string
			Price       float64
			Stock       uint32
			Category_id uint32
		}{
			Id:          productData.GetId(),
			Name:        productData.GetName(),
			Price:       productData.GetPrice(),
			Stock:       productData.GetStock(),
			Category_id: category.GetId(),
		}
		// insert data to products table
		if err := tx.Table(`products`).Create(&product).Error; err != nil {
			return err
		}
		Response.Id = product.Id
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &Response, nil
}

func (p *ProductService) UpdateProduct(ctx context.Context, productData *productpb.Product) (*productpb.Status, error) {
	var Response productpb.Status

	err := p.DB.Transaction(func(tx *gorm.DB) error {
		category := productpb.Category{
			Id:   0,
			Name: productData.GetCategory().GetName(),
		}
		// check data, if data exist or not exist in db, then will be insert to category
		if err := tx.Table(`categories`).Where(`LCASE(name) = ?`, category.GetName()).FirstOrCreate(&category).Error; err != nil {
			return err
		}

		// anonymous struct
		product := struct {
			Id          uint64
			Name        string
			Price       float64
			Stock       uint32
			Category_id uint32
		}{
			Id:          productData.GetId(),
			Name:        productData.GetName(),
			Price:       productData.GetPrice(),
			Stock:       productData.GetStock(),
			Category_id: category.GetId(),
		}
		// insert data to products table
		if err := tx.Table(`products`).Where(`id = ?`, product.Id).Updates(&product).Error; err != nil {
			return err
		}
		Response.Status = 1
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &Response, nil
}

func (p *ProductService) DeleteProduct(ctx context.Context, id *productpb.Id) (*productpb.Status, error) {
	var response productpb.Status

	if err := p.DB.Table(`products`).Where(`id = ?`, id.GetId()).Delete(nil).Error; err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	response.Status = 1
	return &response, nil
}
