package repositories

import (
	"database/sql"
	"time"
)

type CartProductRepo interface {
	AddCartProduct(c CartProduct) error
}

type cpdImpl struct {
	db *sql.DB
}

func NewCartProductRepo(db *sql.DB) CartProductRepo {
	return &cpdImpl{
		db: db,
	}
}

type CartProduct struct {
	CartID      int64
	ProductID   int64
	ProductName string
	ShopID      int64
	Quantity    int64
	Price       float64
	CreateAt    time.Time
	UpdateAt    time.Time
}

func (crt *cpdImpl) AddCartProduct(c CartProduct) error {
	stmt := `INSERT INTO cart_product
				(cart_id, product_id, product_name, shop_id, quantity, price, create_at, update_at)
			VALUES
				(:1, :2, :3, :4, :5, :6, :7, :8)`
	_, err := crt.db.Exec(stmt,
		c.CartID,
		c.ProductID,
		c.ProductName,
		c.ShopID,
		c.Quantity,
		c.Price,
		c.CreateAt,
		c.UpdateAt,
	)

	return err
}
