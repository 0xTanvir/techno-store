package services

import (
	"context"
	"log/slog"
	"sync"

	"techno-store/internal/domain/bo"
	"techno-store/internal/domain/definition"
)

var onceInitSupplierService sync.Once
var supplierServiceInstance *supplierService

type supplierService struct {
	repo definition.SupplierRepository
}

func Supplier(supplierRepo definition.SupplierRepository) *supplierService {
	onceInitSupplierService.Do(func() {
		supplierServiceInstance = &supplierService{
			repo: supplierRepo,
		}
	})

	return supplierServiceInstance
}

func (s *supplierService) List(ctx context.Context, query bo.SupplierQuery) (bo.PaginatedSupplierCollection, error) {
	return s.repo.ListSuppliers(ctx, query)
}

func (s *supplierService) GetSupplierByID(ctx context.Context, supplierID int64) (bo.Supplier, error) {
	return s.repo.GetSupplierByID(ctx, supplierID)
}

func (s *supplierService) CreateSupplier(ctx context.Context, supplier bo.Supplier) (int64, error) {
	if err := s.repo.CreateSupplier(ctx, &supplier); err != nil {
		return -1, err
	}

	if supplier.ID < 1 {
		slog.Warn("inserted supplier has invalid id", slog.String("name", supplier.Name))
	}
	return supplier.ID, nil
}

func (s *supplierService) UpdateSupplier(ctx context.Context, updateSupplier bo.SupplierUpdate) error {
	return s.repo.UpdateSupplier(ctx, updateSupplier)
}

func (s *supplierService) DeleteSupplier(ctx context.Context, supplierID int64) error {
	return s.repo.DeleteSupplier(ctx, supplierID)
}
