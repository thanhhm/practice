package services

import (
	"database/sql"
	"time"

	"github.com/pkg/errors"
	"hcmut.vn/de/repositories"
)

type CartService interface {
	AddToCart(req AddToCartRequest) error
}

type crtSvcImpl struct {
	cartRepo        repositories.CartRepo
	cartProductRepo repositories.CartProductRepo
}

func NewCartService(db *sql.DB) CartService {
	return &crtSvcImpl{
		cartRepo:        repositories.NewCartRepo(db),
		cartProductRepo: repositories.NewCartProductRepo(db),
	}
}

type AddToCartRequest struct {
	CartID   int64                `json:"cart_id"`
	Products []CartProductRequest `json:"products"`
}

type CartProductRequest struct {
	ProductID   int64   `json:"product_id"`
	ProductName string  `json:"product_name"`
	ShopID      int64   `json:"shop_id"`
	Quantity    int64   `json:"quantity"`
	Price       float64 `json:"price"`
}

func (cr *crtSvcImpl) AddToCart(req AddToCartRequest) error {
	var totalPrice float64
	for _, v := range req.Products {
		cartProduct := repositories.CartProduct{
			CartID:      req.CartID,
			ProductID:   v.ProductID,
			ProductName: v.ProductName,
			ShopID:      v.ShopID,
			Quantity:    v.Quantity,
			Price:       v.Price,
			CreateAt:    time.Now(),
			UpdateAt:    time.Now(),
		}

		// Add cart product
		err := cr.cartProductRepo.AddCartProduct(cartProduct)
		if err != nil {
			return errors.Wrap(err, "AddCartProduct error")
		}

		totalPrice = totalPrice + v.Price*float64(v.Quantity)
	}

	// Add cart
	cart := repositories.Cart{
		CartID:     req.CartID,
		TotalPrice: totalPrice,
		CreateAt:   time.Now(),
		UpdateAt:   time.Now(),
	}
	err := cr.cartRepo.AddCart(cart)
	return err
}
