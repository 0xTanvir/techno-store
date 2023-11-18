package bo

import (
	"errors"
	"time"
)

var (
	ErrSupplierNotFound = errors.New("the supplier was not found")
)

// SupplierQuery represent Supplier model query parameter
type SupplierQuery struct {
	Limit  int
	Offset int
}

type Supplier struct {
	ID                 int64     `db:"id"`
	Name               string    `db:"name"`
	Email              string    `db:"email"`
	Phone              string    `db:"phone"`
	StatusID           int64     `db:"status_id"`
	IsVerifiedSupplier bool      `db:"is_verified_supplier"`
	CreatedAt          time.Time `db:"created_at"`
}

type SupplierCollection []Supplier

// PaginatedSupplierCollection model array with total record
type PaginatedSupplierCollection struct {
	Data SupplierCollection

	// This will always return the total of all records
	Total int64
}

type SupplierUpdate struct {
	ID                 int64
	Name               *string
	Email              *string
	Phone              *string
	StatusID           *int64
	IsVerifiedSupplier *bool
}
