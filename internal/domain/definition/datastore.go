package definition

import (
	"context"

	"techno-store/internal/domain/bo"
)

type DataStore struct {
	Brand        BrandRepository
	Category     CategoryRepository
	Supplier     SupplierRepository
	Product      ProductRepository
	ProductStock ProductStockRepository
}

// BrandRepository is the interface that wraps the basic CRUD operations
// defines the rules around what a Brand repository has to be able to perform
// For datastore implementations, see internal/infrastructure/datastores
type BrandRepository interface {
	GetBrandByID(ctx context.Context, brandID int64) (bo.Brand, error)
	CreateBrand(ctx context.Context, brand *bo.Brand) error
	UpdateBrand(ctx context.Context, updateBrand bo.BrandUpdate) error
	DeleteBrand(ctx context.Context, brandID int64) error
	ListBrands(ctx context.Context, brandQuery bo.BrandQuery) (bo.PaginatedBrandCollection, error)
}

// CategoryRepository is the interface that wraps the basic CRUD operations
// defines the rules around what a Category repository has to be able to perform
// For datastore implementations, see internal/infrastructure/datastores
type CategoryRepository interface {
	GetCategoryByID(ctx context.Context, categoryID int64) (bo.Category, error)
	CreateCategory(ctx context.Context, category *bo.Category) error
	UpdateCategory(ctx context.Context, updateCategory bo.CategoryUpdate) error
	DeleteCategory(ctx context.Context, categoryID int64) error
	ListCategories(ctx context.Context) (bo.PaginatedCategoryCollection, error)
}

// SupplierRepository is the interface that wraps the basic CRUD operations
// defines the rules around what a Supplier repository has to be able to perform
// For datastore implementations, see internal/infrastructure/datastores
type SupplierRepository interface {
	GetSupplierByID(ctx context.Context, supplierID int64) (bo.Supplier, error)
	CreateSupplier(ctx context.Context, supplier *bo.Supplier) error
	UpdateSupplier(ctx context.Context, updateSupplier bo.SupplierUpdate) error
	DeleteSupplier(ctx context.Context, supplierID int64) error
	ListSuppliers(ctx context.Context, supplierQuery bo.SupplierQuery) (bo.PaginatedSupplierCollection, error)
}

// ProductRepository is the interface that wraps the basic CRUD operations
// defines the rules around what a Product repository has to be able to perform
// For datastore implementations, see internal/infrastructure/datastores
type ProductRepository interface {
	GetProductByID(ctx context.Context, productID int64) (bo.Product, error)
	CreateProduct(ctx context.Context, product *bo.Product) error
	UpdateProduct(ctx context.Context, updateProduct bo.ProductUpdate) error
	DeleteProduct(ctx context.Context, productID int64) error
	ListProducts(ctx context.Context, productQuery bo.ProductSearchQuery) (bo.PaginatedProductCollection, error)
}

// StockRepository is the interface that wraps the basic CRUD operations
// defines the rules around what a Stock repository has to be able to perform
// For datastore implementations, see internal/infrastructure/datastores
type ProductStockRepository interface {
	GetProductStockByID(ctx context.Context, productStockID int64) (bo.ProductStock, error)
	CreateProductStock(ctx context.Context, productStock *bo.ProductStock) error
	UpdateProductStock(ctx context.Context, updateProductStock bo.ProductStockUpdate) error
	DeleteProductStock(ctx context.Context, productStockID int64) error
	ListProductStocks(ctx context.Context, productStockQuery bo.ProductStockQuery) (bo.PaginatedProductStockCollection, error)
}
