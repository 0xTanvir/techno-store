package bo

import (
	"errors"
	"time"
)

var (
	ErrCategoryNotFound = errors.New("the category was not found")
)

type Category struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	ParentID  int64    `db:"parent_id"` // Pointer to handle NULL values
	Sequence  int64     `db:"sequence"`
	StatusID  int64     `db:"status_id"`
	CreatedAt time.Time `db:"created_at"`
}

type CategoryCollection []Category

// PaginatedCategoryCollection model array with total record
type PaginatedCategoryCollection struct {
	Data CategoryCollection
}

type CategoryUpdate struct {
	ID       int64
	Name     *string
	StatusID *int64
	Sequence *int64
}
