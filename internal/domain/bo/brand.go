package bo

import (
	"errors"
	"time"
)

var (
	ErrBrandNotFound = errors.New("the brand was not found")
)

// BrandQuery represent brand model query parameter
type BrandQuery struct {
	Limit  int
	Offset int
}

type Brand struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	StatusID  int64     `db:"status_id"`
	CreatedAt time.Time `db:"created_at"`
}

// BrandCollection array
type BrandCollection []Brand

// PaginatedBrandCollection model array with total record
type PaginatedBrandCollection struct {
	// This will always return the total of all records
	Total int64
	Data  BrandCollection
}

type BrandUpdate struct {
	ID       int64
	Name     *string
	StatusID *int64
}
