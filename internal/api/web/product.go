package web

import (
	"context"
	"log/slog"
	"net/http"

	"techno-store/internal/api/dto"
	"techno-store/internal/domain/bo"
	"techno-store/internal/domain/services"

	"github.com/gin-gonic/gin"
)

// Get Products godoc
// @Summary      Get Products by query
// @Description  Get Products by query
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        q       query  string  false  "q is query string of product name"
// @Param        brand   query   []int  false  "brand"
// @Param        category  query   int  false  "category"
// @Param        supplier  query   int  false  "supplier"
// @Param        verified_supplier  query   bool  false  "verified_supplier"
// @Param        min_price  query   float64  false  "min_price"
// @Param        max_price  query   float64  false  "max_price"
// @Param        sort  query   string  false  "sort"
// @Param        order  query   string  false  "order"
// @Param        limit   query   int  false  "limit"
// @Param        offset  query   int  false  "offset"
// @Success      200  {object}  dto.PaginatedProduct
// @Failure      400  {string} string  "Invalid request body"
// @Failure      500  {string}  string  "Error"
// @Router       /v1/products [get]
func (r *repos) getProducts(ctx *gin.Context) {
	var productQueryDto dto.ProductQuery
	if err := ctx.ShouldBindQuery(&productQueryDto); err != nil {
		slog.Error("unable to parse query url", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid query value"))
		return
	}

	getProductCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queryModel := productQueryDto.Model()
	products, err := services.Product(r.ds.Product).List(getProductCtx, queryModel)
	if err != nil {
		slog.Error("unable to get products", "cause", err)
		ctx.JSON(http.StatusInternalServerError, dto.Builder().SetMessage("Internal server error"))
		return
	}

	ctx.JSON(http.StatusOK, dto.ToPaginatedProduct(products))
}

// Get Product godoc
// @Summary      Get a Product by id
// @Description  Get a Product by id
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      200  {object}  dto.Product
// @Failure      400  {string} string  "Invalid request body"
// @Failure      404  {object}  dto.Error
// @Failure      500  {string}  string  "Error"
// @Router       /v1/product/{id} [get]
func (r *repos) getProduct(ctx *gin.Context) {
	var wrappedID dto.IDWrapper
	if err := ctx.ShouldBindUri(&wrappedID); err != nil {
		slog.Error("unable to parse product id", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid query value"))
		return
	}

	getProductCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	product, err := services.Product(r.ds.Product).GetProductByID(getProductCtx, wrappedID.ID)
	if err != nil {
		if err == bo.ErrProductNotFound {
			ctx.JSON(http.StatusNotFound, dto.Builder().SetMessage("product not found"))
			return
		}
		slog.Error("unable to get product from database: ", "cause", err)
		ctx.JSON(http.StatusInternalServerError, dto.Builder().SetMessage("Error"))
		return
	}

	ctx.JSON(http.StatusOK, dto.ToProductDTO(product))
}

// Add Product godoc
// @Summary      Add a new product
// @Description  Create a new product in the system
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        request body dto.Product  true  "Product params"
// @Success      201  {object}  dto.IDWrapper
// @Failure      400  {string} string  "Invalid request body"
// @Failure      404  {object}  dto.Error
// @Failure      500  {string}  string  "Error"
// @Router       /v1/product [post]
func (r *repos) addProduct(ctx *gin.Context) {
	productDto := dto.Product{}
	if err := ctx.ShouldBindJSON(&productDto); err != nil {
		slog.Error("unable to parse product from request body", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid request body"))
		return
	}

	model := productDto.Model()

	addProductCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	id, err := services.Product(r.ds.Product).CreateProduct(addProductCtx, model)
	if err != nil {
		slog.Error("unable to create product", "cause", err)
		ctx.JSON(http.StatusInternalServerError, dto.Builder().SetMessage("Internal server error"))
		return
	}

	ctx.JSON(http.StatusCreated, dto.IDWrapper{ID: id})
}

// UpdateProduct godoc
// @Summary      Update a product by id
// @Description  Update a product by id
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        request body dto.ProductUpdate  true  "product params"
// @Param        id   path      int  true  "Product ID"
// @Success      204  {string}  "productDto updated"
// @Failure      400  {string} string  "Invalid request body"
// @Failure      404  {object}  dto.Error
// @Failure      500  {string}  string  "Error"
// @Router       /v1/product/{id} [patch]
func (r *repos) updateProduct(ctx *gin.Context) {
	var wrappedID dto.IDWrapper
	if err := ctx.ShouldBindUri(&wrappedID); err != nil {
		slog.Error("unable to parse product id", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid query value"))
		return
	}

	var productDto dto.ProductUpdate
	if err := ctx.ShouldBindJSON(&productDto); err != nil {
		slog.Error("unable to parse product from request body", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid request body"))
		return
	}

	updateProductCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	productDto.ID = wrappedID.ID
	if err := services.Product(r.ds.Product).UpdateProduct(updateProductCtx, productDto.Model()); err != nil {
		if err == bo.ErrProductNotFound {
			ctx.JSON(http.StatusNotFound, dto.Builder().SetMessage("product not found"))
			return
		}
		slog.Error("unable to update product", "cause", err)
		ctx.JSON(http.StatusInternalServerError, dto.Builder().SetMessage("Internal server error"))
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"message": "product updated"})
}

// DeleteProduct godoc
// @Summary      Delete a Product by id
// @Description  Delete a Product by id
// @Tags         Product
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Product ID"
// @Success      204  {string}  "Product delete processed"
// @Failure      400  {string} 	string  "Invalid request body"
// @Failure      404  {object}  string  "Product not found"
// @Failure      500  {string}  string  "Error"
// @Router       /v1/product/{id} [delete]
func (r *repos) deleteProduct(ctx *gin.Context) {
	var wrappedID dto.IDWrapper
	if err := ctx.ShouldBindUri(&wrappedID); err != nil {
		slog.Error("unable to parse product id", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid query value"))
		return
	}

	deleteProductCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := services.Product(r.ds.Product).DeleteProduct(deleteProductCtx, wrappedID.ID); err != nil {
		if err == bo.ErrProductNotFound {
			ctx.JSON(http.StatusNotFound, dto.Builder().SetMessage("product not found"))
			return
		}
		slog.Error("unable to delete product", "cause", err)
		ctx.JSON(http.StatusInternalServerError, dto.Builder().SetMessage("Internal server error"))
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"message": "product deleted"})
}
