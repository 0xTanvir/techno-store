package bo

import (
	"errors"
	"time"
)

var (
	ErrProductStockNotFound = errors.New("the product stock was not found")
)

// ProductStockQuery represent ProductStock model query parameter
type ProductStockQuery struct {
	Limit  int
	Offset int
}

type ProductStock struct {
	ID            int64     `db:"id"`
	ProductID     int64     `db:"product_id"`
	StockQuantity int64     `db:"stock_quantity"`
	UpdatedAt     time.Time `db:"updated_at"`
}

type ProductStockCollection []ProductStock

// PaginatedProductStockCollection model array with total record
type PaginatedProductStockCollection struct {
	Data ProductStockCollection

	// This will always return the total of all records
	Total int64
}

type ProductStockUpdate struct {
	ProductID     int64
	StockQuantity *int64
}
