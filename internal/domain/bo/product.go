package bo

import "errors"

var (
	ErrProductNotFound = errors.New("the product was not found")
)

type PriceRangeFilter struct {
	Min float64
	Max float64
}

type ProductFilter struct {
	Query                  string
	PriceRangeFilter       PriceRangeFilter
	BrandFilter            []int64
	CategoryFilter         int64
	SupplierFilter         int64
	VerifiedSupplierFilter bool
}

type ProductPaging struct {
	Limit  int
	Offset int
}

type ProductSort struct {
	Field string
	Order string
}

// ProductQuery represent Product model query parameter
type ProductSearchQuery struct {
	Filter ProductFilter
	Paging ProductPaging
	Sort   ProductSort
}

type Product struct {
	ID             int64   `db:"id"`
	Name           string  `db:"name"`
	Description    string  `db:"description"`
	Specifications string  `db:"specifications"`
	BrandID        int64   `db:"brand_id"`
	CategoryID     int64   `db:"category_id"`
	SupplierID     int64   `db:"supplier_id"`
	UnitPrice      float64 `db:"unit_price"`
	DiscountPrice  float64 `db:"discount_price"`
	Tags           string  `db:"tags"`
	StatusID       int64   `db:"status_id"`
}

type ProductCollection []Product

// PaginatedProductCollection model array with total record
type PaginatedProductCollection struct {
	Data ProductCollection

	// This will always return the total of all records
	Total int64
}

type ProductUpdate struct {
	ID             int64
	Name           *string
	Description    *string
	Specifications *string
	BrandID        *int64
	CategoryID     *int64
	SupplierID     *int64
	UnitPrice      *float64
	DiscountPrice  *float64
	Tags           *string
	StatusID       *int64
}
