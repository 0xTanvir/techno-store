package mockdb

import (
	"techno-store/internal/domain/definition"

	"go.uber.org/mock/gomock"
)

func GetInstance(ctrl *gomock.Controller) definition.DataStore {
	return definition.DataStore{
		Brand:        NewMockBrandRepository(ctrl),
		Category:     NewMockCategoryRepository(ctrl),
		Supplier:     NewMockSupplierRepository(ctrl),
		Product:      NewMockProductRepository(ctrl),
		ProductStock: NewMockProductStockRepository(ctrl),
	}
}
