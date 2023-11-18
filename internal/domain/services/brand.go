package services

import (
	"context"
	"log/slog"
	"sync"

	"techno-store/internal/domain/bo"
	"techno-store/internal/domain/definition"
)

var onceInitBrandService sync.Once
var brandServiceInstance *brandService

type brandService struct {
	repo definition.BrandRepository
}

func Brand(brandRepo definition.BrandRepository) *brandService {
	onceInitBrandService.Do(func() {
		brandServiceInstance = &brandService{
			repo: brandRepo,
		}
	})

	return brandServiceInstance
}

func (s *brandService) List(ctx context.Context, query bo.BrandQuery) (bo.PaginatedBrandCollection, error) {
	return s.repo.ListBrands(ctx, query)
}

func (s *brandService) GetBrandByID(ctx context.Context, brandID int64) (bo.Brand, error) {
	return s.repo.GetBrandByID(ctx, brandID)
}

func (s *brandService) CreateBrand(ctx context.Context, brand bo.Brand) (int64, error) {
	if err := s.repo.CreateBrand(ctx, &brand); err != nil {
		return -1, err
	}

	if brand.ID < 1 {
		slog.Warn("inserted brand has invalid id", slog.String("name", brand.Name))
	}
	return brand.ID, nil
}

func (s *brandService) UpdateBrand(ctx context.Context, updateBrand bo.BrandUpdate) error {
	return s.repo.UpdateBrand(ctx, updateBrand)
}

func (s *brandService) DeleteBrand(ctx context.Context, brandID int64) error {
	return s.repo.DeleteBrand(ctx, brandID)
}
