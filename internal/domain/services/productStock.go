package services

import (
	"context"
	"log/slog"
	"sync"

	"techno-store/internal/domain/bo"
	"techno-store/internal/domain/definition"
)

var onceInitProductStockService sync.Once
var productStockServiceInstance *productStockService

type productStockService struct {
	repo definition.ProductStockRepository
}

func ProductStock(productStockRepo definition.ProductStockRepository) *productStockService {
	onceInitProductStockService.Do(func() {
		productStockServiceInstance = &productStockService{
			repo: productStockRepo,
		}
	})

	return productStockServiceInstance
}

func (s *productStockService) List(ctx context.Context, query bo.ProductStockQuery) (bo.PaginatedProductStockCollection, error) {
	return s.repo.ListProductStocks(ctx, query)
}

func (s *productStockService) GetProductStockByID(ctx context.Context, productStockID int64) (bo.ProductStock, error) {
	return s.repo.GetProductStockByID(ctx, productStockID)
}

func (s *productStockService) CreateProductStock(ctx context.Context, productStock bo.ProductStock) (int64, error) {
	if err := s.repo.CreateProductStock(ctx, &productStock); err != nil {
		return -1, err
	}

	if productStock.ID < 1 {
		slog.Warn("inserted productStock has invalid id", slog.Int64("product id", productStock.ProductID))
	}
	return productStock.ID, nil
}

func (s *productStockService) UpdateProductStock(ctx context.Context, updateProductStock bo.ProductStockUpdate) error {
	return s.repo.UpdateProductStock(ctx, updateProductStock)
}

func (s *productStockService) DeleteProductStock(ctx context.Context, productStockID int64) error {
	return s.repo.DeleteProductStock(ctx, productStockID)
}
