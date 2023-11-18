package services

import (
	"context"
	"log/slog"
	"sync"

	"techno-store/internal/domain/bo"
	"techno-store/internal/domain/definition"
)

var onceInitCategoryService sync.Once
var categoryServiceInstance *categoryService

type categoryService struct {
	repo definition.CategoryRepository
}

func Category(categoryRepo definition.CategoryRepository) *categoryService {
	onceInitCategoryService.Do(func() {
		categoryServiceInstance = &categoryService{
			repo: categoryRepo,
		}
	})

	return categoryServiceInstance
}

func (s *categoryService) List(ctx context.Context) (bo.PaginatedCategoryCollection, error) {
	return s.repo.ListCategories(ctx)
}

func (s *categoryService) GetCategoryByID(ctx context.Context, categoryID int64) (bo.Category, error) {
	return s.repo.GetCategoryByID(ctx, categoryID)
}

func (s *categoryService) CreateCategory(ctx context.Context, category bo.Category) (int64, error) {
	if err := s.repo.CreateCategory(ctx, &category); err != nil {
		return -1, err
	}

	if category.ID < 1 {
		slog.Warn("inserted category has invalid id", slog.String("name", category.Name))
	}
	return category.ID, nil
}

func (s *categoryService) UpdateCategory(ctx context.Context, updateCategory bo.CategoryUpdate) error {
	return s.repo.UpdateCategory(ctx, updateCategory)
}

func (s *categoryService) DeleteCategory(ctx context.Context, categoryID int64) error {
	return s.repo.DeleteCategory(ctx, categoryID)
}
