package repositories

import (
	"database/sql"
	"time"
)

type CartRepo interface {
	AddCart(c Cart) error
}

type crtImpl struct {
	db *sql.DB
}

func NewCartRepo(db *sql.DB) CartRepo {
	return &crtImpl{
		db: db,
	}
}

type Cart struct {
	CartID     int64
	TotalPrice float64
	CreateAt   time.Time
	UpdateAt   time.Time
}

func (crt *crtImpl) AddCart(c Cart) error {
	stmt := `INSERT INTO cart
				(cart_id, total_price, create_at, update_at)
			VALUES
				(:1, :2, :3, :4)`
	_, err := crt.db.Exec(stmt,
		c.CartID,
		c.TotalPrice,
		c.CreateAt,
		c.UpdateAt,
	)

	return err
}
