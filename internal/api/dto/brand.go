package dto

import "techno-store/internal/domain/bo"

type Brand struct {
	ID       int64  `json:"id,omitempty"`
	Name     string `json:"name" binding:"required"`
	StatusID int64  `json:"status_id" binding:"required"`
}

func (b Brand) Model() bo.Brand {
	return bo.Brand{
		ID:       b.ID,
		Name:     b.Name,
		StatusID: b.StatusID,
	}
}

type BrandUpdate struct {
	ID       int64   `json:"id"`
	Name     *string `json:"name"`
	StatusID *int64  `json:"status_id"`
}

func (b BrandUpdate) Model() bo.BrandUpdate {
	return bo.BrandUpdate{
		ID:       b.ID,
		Name:     b.Name,
		StatusID: b.StatusID,
	}
}

// BrandCollection array
type BrandCollection []Brand

// PaginatedBrandCollection model array with total record
type PaginatedBrandCollection struct {
	// This will always return the total of all records
	Total int64           `json:"total"`
	Data  BrandCollection `json:"data"`
}

func ToPaginatedBrand(bo bo.PaginatedBrandCollection) PaginatedBrandCollection {
	brands := []Brand{}
	for _, brand := range bo.Data {
		brands = append(brands, ToBrandDTO(brand))
	}

	return PaginatedBrandCollection{
		Total: bo.Total,
		Data:  brands,
	}
}

// BrandQuery represent Category model query parameter
type BrandQuery struct {
	Limit  int `form:"limit,default=20" json:"limit,omitempty" binding:"min=1"`
	Offset int `form:"offset" json:"offset,omitempty" binding:"omitempty,min=0"`
}

func (q BrandQuery) Model() bo.BrandQuery {
	// Setup some default behavior
	if q.Limit <= 0 {
		q.Limit = 20
	}
	if q.Offset <= 0 {
		q.Offset = 0
	}

	return bo.BrandQuery{
		Limit:  q.Limit,
		Offset: q.Offset,
	}
}

// Convert BO to DTO
func ToBrandDTO(bo bo.Brand) Brand {
	return Brand{
		ID:       bo.ID,
		Name:     bo.Name,
		StatusID: bo.StatusID,
	}
}
