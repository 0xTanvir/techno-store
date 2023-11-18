package services

import (
	"context"
	"log/slog"
	"sync"

	"techno-store/internal/domain/bo"
	"techno-store/internal/domain/definition"
)

var onceInitProductService sync.Once
var productServiceInstance *productService

type productService struct {
	repo definition.ProductRepository
}

func Product(productRepo definition.ProductRepository) *productService {
	onceInitProductService.Do(func() {
		productServiceInstance = &productService{
			repo: productRepo,
		}
	})

	return productServiceInstance
}

func (s *productService) List(ctx context.Context, query bo.ProductSearchQuery) (bo.PaginatedProductCollection, error) {
	return s.repo.ListProducts(ctx, query)
}

func (s *productService) GetProductByID(ctx context.Context, productID int64) (bo.Product, error) {
	return s.repo.GetProductByID(ctx, productID)
}

func (s *productService) CreateProduct(ctx context.Context, product bo.Product) (int64, error) {
	if err := s.repo.CreateProduct(ctx, &product); err != nil {
		return -1, err
	}

	if product.ID < 1 {
		slog.Warn("inserted product has invalid id", slog.String("name", product.Name))
	}
	return product.ID, nil
}

func (s *productService) UpdateProduct(ctx context.Context, updateProduct bo.ProductUpdate) error {
	return s.repo.UpdateProduct(ctx, updateProduct)
}

func (s *productService) DeleteProduct(ctx context.Context, productID int64) error {
	return s.repo.DeleteProduct(ctx, productID)
}
