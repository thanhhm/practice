package services

import (
	"database/sql"

	"hcmut.vn/de/repositories"
)

type ProductService interface {
	CreateProduct(req CreateProductRequest) error
	GetProduct(req GetProductsRequest) ([]repositories.Product, error)
}

type pdSvcImpl struct {
	productRepo repositories.ProductRepo
}

func NewProductService(db *sql.DB) ProductService {
	return &pdSvcImpl{
		productRepo: repositories.NewProductRepo(db),
	}
}

type CreateProductRequest struct {
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"product_name"`
	ShopID      int64   `json:"shop_id"`
	Price       float64 `json:"price"`
	Status      int64   `json:"status"`
	Stock       int64   `json:"stock"`
}

func (pds *pdSvcImpl) CreateProduct(req CreateProductRequest) error {
	tsData := repositories.Product{
		ProductID:   req.ProductID,
		ProductName: req.ProductName,
		ShopID:      req.ShopID,
		Price:       req.Price,
		Status:      req.Status,
		Stock:       req.Stock,
	}

	err := pds.productRepo.CreateProduct(tsData)
	if err != nil {
		return err
	}

	return nil
}

type GetProductsRequest struct {
	SortCondition string
}

func (pds *pdSvcImpl) GetProduct(req GetProductsRequest) ([]repositories.Product, error) {
	result, err := pds.productRepo.GetProducts(req.SortCondition)
	if err != nil {
		return nil, err
	}

	return result, nil
}
