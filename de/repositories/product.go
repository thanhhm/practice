package repositories

import (
	"database/sql"
)

const (
	Available   = 1
	UnAvailable = 0
)

type ProductRepo interface {
	CreateProduct(pdInfo Product) error
	GetProducts(sortCondition string) ([]Product, error)
}

type pdImpl struct {
	db *sql.DB
}

func NewProductRepo(db *sql.DB) ProductRepo {
	return &pdImpl{
		db: db,
	}
}

type Product struct {
	ProductID   int64
	ProductName string
	ShopID      int64
	OrderCount  int64
	Price       float64
	Status      int64
	Stock       int64
	TotalRate   int64
	CreateAt    int64
	UpdateAt    int64
}

func (pd *pdImpl) CreateProduct(pdInfo Product) error {
	stmt := `INSERT INTO PRODUCT
				(PRODUCT_ID, PRODUCT_NAME, SHOP_ID, PRICE, STATUS, STOCK)
			VALUES
				(:1, :2, :3, :4, :5, :6)`
	_, err := pd.db.Exec(stmt,
		pdInfo.ProductID,
		pdInfo.ProductName,
		pdInfo.ShopID,
		pdInfo.Price,
		pdInfo.Status,
		pdInfo.Stock,
	)

	return err
}

func (pd *pdImpl) GetProducts(sortCondition string) ([]Product, error) {
	query := `SELECT product_id, product_name, shop_id, order_count, price, status, stock, total_rate 
			  FROM product
			  WHERE status = :1
			  ORDER BY :2`
	if sortCondition == "" {
		sortCondition = "order_count DESC"
	}

	rows, err := pd.db.Query(query, Available, sortCondition)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []Product
	for rows.Next() {
		var product Product
		if err := rows.Scan(
			&product.ProductID,
			&product.ProductName,
			&product.ShopID,
			&product.OrderCount,
			&product.Price,
			&product.Status,
			&product.Stock,
			&product.TotalRate,
		); err != nil {
			return nil, err
		}
		result = append(result, product)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return result, nil
}
