package web

import (
	"log/slog"

	"techno-store/config"
	"techno-store/internal/domain/definition"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS(router *gin.Engine) {
	slog.Warn("Enabling CORS")
	router.Use(cors.New(cors.Config{
		// Only allow your specific domains instead of all
		AllowOrigins: []string{"http://localhost:3000"},

		// Specify the actual methods your frontend uses
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD", "OPTIONS"},

		// Specify the headers your frontend sends
		AllowHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Requested-With",
		},

		AllowCredentials: true,

		// The rest of the properties should be disabled for security reasons unless you specifically need them
		AllowWildcard:          false,
		AllowBrowserExtensions: false,
		AllowWebSockets:        false,
		AllowFiles:             false,
	}))
}

type repos struct {
	config config.ServerConfig
	ds     definition.DataStore
}

func NewAPIService(cfg config.ServerConfig, ds definition.DataStore) *repos {
	return &repos{
		config: cfg,
		ds:     ds,
	}
}

func (r *repos) InstallRoutes(router *gin.Engine) {
	CORS(router)
	router.GET("", health)
	router.GET("/health", health)

	v1 := router.Group("/v1")

	// Brand group
	brandsGroup := v1.Group("/brands")
	brandGroup := v1.Group("/brand")
	{
		brandsGroup.GET("", r.getBrands)
		brandGroup.GET("/:id", r.getBrand)
		brandGroup.POST("", r.addBrand)
		brandGroup.PATCH("/:id", r.updateBrand)
		brandGroup.DELETE("/:id", r.deleteBrand)
	}

	// Category group
	categoriesGroup := v1.Group("/categories")
	categoryGroup := v1.Group("/category")
	{
		categoriesGroup.GET("", r.getCategories)
		categoryGroup.GET("/:id", r.getCategory)
		categoryGroup.POST("", r.addCategory)
		categoryGroup.PATCH("/:id", r.updateCategory)
		categoryGroup.DELETE("/:id", r.deleteCategory)
	}

	// Product group
	productsGroup := v1.Group("/products")
	productGroup := v1.Group("/product")
	{
		productsGroup.GET("", r.getProducts)
		productGroup.GET("/:id", r.getProduct)
		productGroup.POST("", r.addProduct)
		productGroup.PATCH("/:id", r.updateProduct)
		productGroup.DELETE("/:id", r.deleteProduct)
	}

	// Supplier group
	suppliersGroup := v1.Group("/suppliers")
	supplierGroup := v1.Group("/supplier")
	{
		suppliersGroup.GET("", r.getSuppliers)
		supplierGroup.GET("/:id", r.getSupplier)
		supplierGroup.POST("", r.addSupplier)
		supplierGroup.PATCH("/:id", r.updateSupplier)
		supplierGroup.DELETE("/:id", r.deleteSupplier)
	}

	// ProductStock group
	productStocksGroup := v1.Group("/product-stocks")
	productStockGroup := v1.Group("/product-stock")
	{
		productStocksGroup.GET("", r.getProductStocks)
		productStockGroup.GET("/:id", r.getProductStock)
		productStockGroup.POST("", r.addProductStock)
		productStockGroup.PATCH("/:id", r.updateProductStock)
		productStockGroup.DELETE("/:id", r.deleteProductStock)
	}
}
