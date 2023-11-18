package dto

import "techno-store/internal/domain/bo"

type ProductStock struct {
	ID            int64 `json:"id,omitempty"`
	ProductID     int64 `json:"product_id"`
	StockQuantity int64 `json:"stock_quantity"`
}

func ToProductStockDTO(bo bo.ProductStock) ProductStock {
	return ProductStock{
		ID:            bo.ID,
		ProductID:     bo.ProductID,
		StockQuantity: bo.StockQuantity,
	}
}

func (b ProductStock) Model() bo.ProductStock {
	return bo.ProductStock{
		ID:            b.ID,
		ProductID:     b.ProductID,
		StockQuantity: b.StockQuantity,
	}
}

type ProductStockUpdate struct {
	ProductID     int64  `json:"product_id"`
	StockQuantity *int64 `json:"stock_quantity"`
}

func (b ProductStockUpdate) Model() bo.ProductStockUpdate {
	return bo.ProductStockUpdate{
		ProductID:     b.ProductID,
		StockQuantity: b.StockQuantity,
	}
}

// ProductStockCollection array
type ProductStockCollection []ProductStock

// PaginatedProductStockCollection model array with total record
type PaginatedProductStockCollection struct {
	// This will always return the total of all records
	Total int64                  `json:"total"`
	Data  ProductStockCollection `json:"data"`
}

func ToPaginatedProductStock(bo bo.PaginatedProductStockCollection) PaginatedProductStockCollection {
	var dto ProductStockCollection
	for _, v := range bo.Data {
		dto = append(dto, ToProductStockDTO(v))
	}

	return PaginatedProductStockCollection{
		Total: bo.Total,
		Data:  dto,
	}
}

type ProductStockQuery struct {
	Limit  int `form:"limit,default=20" json:"limit,omitempty" binding:"min=1"`
	Offset int `form:"offset" json:"offset,omitempty" binding:"omitempty,min=0"`
}

func (q ProductStockQuery) Model() bo.ProductStockQuery {
	// Setup some default behavior
	if q.Limit <= 0 {
		q.Limit = 20
	}
	if q.Offset <= 0 {
		q.Offset = 0
	}

	return bo.ProductStockQuery{
		Limit:  q.Limit,
		Offset: q.Offset,
	}
}
