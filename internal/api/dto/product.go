package dto

import "techno-store/internal/domain/bo"

type IDWrapper struct {
	ID int64 `uri:"id" json:"id,omitempty" binding:"required,min=1"`
}

type Product struct {
	ID             int64   `json:"id,omitempty"`
	Name           string  `json:"name"`
	Description    string  `json:"description,omitempty"`
	Specifications string  `json:"specifications,omitempty"`
	BrandID        int64   `json:"brand_id"`
	CategoryID     int64   `json:"category_id"`
	SupplierID     int64   `json:"supplier_id"`
	UnitPrice      float64 `json:"unit_price"`
	DiscountPrice  float64 `json:"discount_price,omitempty"`
	Tags           string  `json:"tags,omitempty"`
	StatusID       int64   `json:"status_id"`
}

func ToProductDTO(bo bo.Product) Product {
	return Product{
		ID:             bo.ID,
		Name:           bo.Name,
		Description:    bo.Description,
		Specifications: bo.Specifications,
		BrandID:        bo.BrandID,
		CategoryID:     bo.CategoryID,
		SupplierID:     bo.SupplierID,
		UnitPrice:      bo.UnitPrice,
		DiscountPrice:  bo.DiscountPrice,
		Tags:           bo.Tags,
		StatusID:       bo.StatusID,
	}
}

func (p Product) Model() bo.Product {
	return bo.Product{
		ID:             p.ID,
		Name:           p.Name,
		Description:    p.Description,
		Specifications: p.Specifications,
		BrandID:        p.BrandID,
		CategoryID:     p.CategoryID,
		SupplierID:     p.SupplierID,
		UnitPrice:      p.UnitPrice,
		DiscountPrice:  p.DiscountPrice,
		Tags:           p.Tags,
		StatusID:       p.StatusID,
	}
}

type ProductCollection []Product

// PaginatedProductCollection model array with total record
type PaginatedProduct struct {
	// This will always return the total of all records
	Total int64             `json:"total"`
	Data  ProductCollection `json:"data"`
}

func ToPaginatedProduct(bo bo.PaginatedProductCollection) PaginatedProduct {
	products := []Product{}
	for _, prod := range bo.Data {
		products = append(products, ToProductDTO(prod))
	}
	return PaginatedProduct{
		Data:  products,
		Total: bo.Total,
	}
}

// ProductQuery represent Product model query parameter
type ProductQuery struct {
	Limit            int     `form:"limit,default=20" json:"limit,omitempty" binding:"min=1"`
	Offset           int     `form:"offset" json:"offset,omitempty" binding:"omitempty,min=0"`
	Sort             string  `form:"sort" json:"sort,omitempty"`
	Order            string  `form:"order" json:"order,omitempty"`
	MinPrice         float64 `form:"min_price" json:"min_price,omitempty"`
	MaxPrice         float64 `form:"max_price" json:"max_price,omitempty"`
	VerifiedSupplier bool    `form:"verified_supplier" json:"verified_supplier,omitempty"`
	Supplier         int64   `form:"supplier" json:"supplier,omitempty"`
	Brand            []int64 `form:"brand" json:"brand,omitempty"`
	Category         int64   `form:"category" json:"category,omitempty"`
	Q                string  `form:"q" json:"q,omitempty"`
}

func (p ProductQuery) Model() bo.ProductSearchQuery {
	// Setup some default behavior
	if p.Limit <= 0 {
		p.Limit = 20
	}
	if p.Offset <= 0 {
		p.Offset = 0
	}
	if p.Sort == "" {
		p.Sort = "unit_price"
	}
	if p.Order == "" {
		p.Order = "ASC"
	}
	offset := p.Offset * p.Limit

	return bo.ProductSearchQuery{
		Filter: bo.ProductFilter{
			Query: p.Q,
			PriceRangeFilter: bo.PriceRangeFilter{
				Min: p.MinPrice,
				Max: p.MaxPrice,
			},
			BrandFilter:            p.Brand,
			CategoryFilter:         p.Category,
			SupplierFilter:         p.Supplier,
			VerifiedSupplierFilter: p.VerifiedSupplier,
		},
		Paging: bo.ProductPaging{
			Limit:  p.Limit,
			Offset: offset,
		},
		Sort: bo.ProductSort{
			Field: p.Sort,
			Order: p.Order,
		},
	}
}

type ProductUpdate struct {
	ID             int64    `json:"id"`
	Name           *string  `json:"name,omitempty"`
	Description    *string  `json:"description,omitempty"`
	Specifications *string  `json:"specifications,omitempty"`
	BrandID        *int64   `json:"brand_id,omitempty"`
	CategoryID     *int64   `json:"category_id,omitempty"`
	SupplierID     *int64   `json:"supplier_id,omitempty"`
	UnitPrice      *float64 `json:"unit_price,omitempty"`
	DiscountPrice  *float64 `json:"discount_price,omitempty"`
	Tags           *string  `json:"tags,omitempty"`
	StatusID       *int64   `json:"status_id,omitempty"`
}

func (p ProductUpdate) Model() bo.ProductUpdate {
	return bo.ProductUpdate{
		ID:             p.ID,
		Name:           p.Name,
		Description:    p.Description,
		Specifications: p.Specifications,
		BrandID:        p.BrandID,
		CategoryID:     p.CategoryID,
		SupplierID:     p.SupplierID,
		UnitPrice:      p.UnitPrice,
		DiscountPrice:  p.DiscountPrice,
		Tags:           p.Tags,
		StatusID:       p.StatusID,
	}
}
