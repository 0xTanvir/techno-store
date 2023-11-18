package dto

import "techno-store/internal/domain/bo"

type Supplier struct {
	ID                 int64  `json:"id,omitempty"`
	Name               string `json:"name"`
	Email              string `json:"email"`
	Phone              string `json:"phone"`
	StatusID           int64  `json:"status_id"`
	IsVerifiedSupplier bool   `json:"is_verified_supplier"`
}

func ToSupplierDTO(bo bo.Supplier) Supplier {
	return Supplier{
		ID:                 bo.ID,
		Name:               bo.Name,
		Email:              bo.Email,
		Phone:              bo.Phone,
		StatusID:           bo.StatusID,
		IsVerifiedSupplier: bo.IsVerifiedSupplier,
	}
}

func (b Supplier) Model() bo.Supplier {
	return bo.Supplier{
		ID:                 b.ID,
		Name:               b.Name,
		Email:              b.Email,
		Phone:              b.Phone,
		StatusID:           b.StatusID,
		IsVerifiedSupplier: b.IsVerifiedSupplier,
	}
}

type SupplierUpdate struct {
	ID                 int64   `json:"id"`
	Name               *string `json:"name"`
	Email              *string `json:"email"`
	Phone              *string `json:"phone"`
	StatusID           *int64  `json:"status_id"`
	IsVerifiedSupplier *bool   `json:"is_verified_supplier"`
}

func (b SupplierUpdate) Model() bo.SupplierUpdate {
	return bo.SupplierUpdate{
		ID:                 b.ID,
		Name:               b.Name,
		Email:              b.Email,
		Phone:              b.Phone,
		StatusID:           b.StatusID,
		IsVerifiedSupplier: b.IsVerifiedSupplier,
	}
}

// SupplierCollection array
type SupplierCollection []Supplier

// PaginatedSupplierCollection model array with total record
type PaginatedSupplierCollection struct {
	// This will always return the total of all records
	Total int64              `json:"total"`
	Data  SupplierCollection `json:"data"`
}

func ToPaginatedSupplier(bo bo.PaginatedSupplierCollection) PaginatedSupplierCollection {
	suppliers := []Supplier{}
	for _, supplier := range bo.Data {
		suppliers = append(suppliers, ToSupplierDTO(supplier))
	}

	return PaginatedSupplierCollection{
		Total: bo.Total,
		Data:  suppliers,
	}
}

// SupplierQuery represent Category model query parameter
type SupplierQuery struct {
	Limit  int `form:"limit,default=20" json:"limit,omitempty" binding:"min=1"`
	Offset int `form:"offset" json:"offset,omitempty" binding:"omitempty,min=0"`
}

func (q SupplierQuery) Model() bo.SupplierQuery {
	// Setup some default behavior
	if q.Limit <= 0 {
		q.Limit = 20
	}
	if q.Offset <= 0 {
		q.Offset = 0
	}

	return bo.SupplierQuery{
		Limit:  q.Limit,
		Offset: q.Offset,
	}
}
