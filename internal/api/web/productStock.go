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

// Get ProductStocks godoc
// @Summary      Get ProductStocks
// @Description  Get ProductStocks
// @Tags         ProductStock
// @Accept       json
// @Produce      json
// @Param        limit   query   int  false  "limit"
// @Param        offset  query   int  false  "offset"
// @Success      200  {object}  dto.PaginatedProductStockCollection
// @Failure      400  {string} string  "Invalid request body"
// @Failure      500  {string}  string  "Error"
// @Router       /v1/product-stocks [get]
func (r *repos) getProductStocks(ctx *gin.Context) {
	var productStockQueryDto dto.ProductStockQuery
	if err := ctx.ShouldBindQuery(&productStockQueryDto); err != nil {
		slog.Error("unable to parse query url", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid query value"))
		return
	}

	getProductStockCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pbc, err := services.ProductStock(r.ds.ProductStock).List(getProductStockCtx, productStockQueryDto.Model())
	if err != nil {
		slog.Error("unable to get productStocks", "cause", err)
		ctx.JSON(http.StatusInternalServerError, dto.Builder().SetMessage("Internal server error"))
		return
	}

	ctx.JSON(http.StatusOK, dto.ToPaginatedProductStock(pbc))
}

// Get ProductStock godoc
// @Summary      Get a ProductStock by id
// @Description  Get a ProductStock by id
// @Tags         ProductStock
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ProductStock ID"
// @Success      200  {object}  dto.ProductStock
// @Failure      400  {string} string  "Invalid request body"
// @Failure      404  {object}  dto.Error
// @Failure      500  {string}  string  "Error"
// @Router       /v1/product-stock/{id} [get]
func (r *repos) getProductStock(ctx *gin.Context) {
	var wrappedID dto.IDWrapper
	if err := ctx.ShouldBindUri(&wrappedID); err != nil {
		slog.Error("unable to parse productStock id", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid query value"))
		return
	}

	getProductStockCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	productStock, err := services.ProductStock(r.ds.ProductStock).GetProductStockByID(getProductStockCtx, wrappedID.ID)
	if err != nil {
		if err == bo.ErrProductStockNotFound {
			ctx.JSON(http.StatusNotFound, dto.Builder().SetMessage("productStock not found"))
			return
		}
		slog.Error("unable to get productStock from database: ", "cause", err)
		ctx.JSON(http.StatusInternalServerError, dto.Builder().SetMessage("Error"))
		return
	}

	ctx.JSON(http.StatusOK, dto.ToProductStockDTO(productStock))
}

// Add ProductStock godoc
// @Summary      Add a new ProductStock
// @Description  Create a new ProductStock in the system
// @Tags         ProductStock
// @Accept       json
// @Produce      json
// @Param        request body dto.ProductStock  true  "ProductStock params"
// @Success      201  {object}  dto.IDWrapper
// @Failure      400  {string} string  "Invalid request body"
// @Failure      404  {object}  dto.Error
// @Failure      500  {string}  string  "Error"
// @Router       /v1/product-stock [post]
func (r *repos) addProductStock(ctx *gin.Context) {
	productStockDto := dto.ProductStock{}
	if err := ctx.ShouldBindJSON(&productStockDto); err != nil {
		slog.Error("unable to parse productStock from request body", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid request body"))
		return
	}

	model := productStockDto.Model()

	addProductStockCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	id, err := services.ProductStock(r.ds.ProductStock).CreateProductStock(addProductStockCtx, model)
	if err != nil {
		slog.Error("unable to create productStock", "cause", err)
		ctx.JSON(http.StatusInternalServerError, dto.Builder().SetMessage("Internal server error"))
		return
	}

	ctx.JSON(http.StatusCreated, dto.IDWrapper{ID: id})
}

// UpdateProductStock godoc
// @Summary      Update a ProductStock by id
// @Description  Update a ProductStock by id
// @Tags         ProductStock
// @Accept       json
// @Produce      json
// @Param        request body dto.ProductStockUpdate  true  "ProductStock params"
// @Param        id   path      int  true  "ProductStock ID"
// @Success      204  {string}  "ProductStockDto updated"
// @Failure      400  {string} string  "Invalid request body"
// @Failure      404  {object}  dto.Error
// @Failure      500  {string}  string  "Error"
// @Router       /v1/product-stock/{id} [patch]
func (r *repos) updateProductStock(ctx *gin.Context) {
	var wrappedID dto.IDWrapper
	if err := ctx.ShouldBindUri(&wrappedID); err != nil {
		slog.Error("unable to parse productStock id", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid query value"))
		return
	}

	var productStockDto dto.ProductStockUpdate
	if err := ctx.ShouldBindJSON(&productStockDto); err != nil {
		slog.Error("unable to parse productStock from request body", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid request body"))
		return
	}

	updateProductStockCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	productStockDto.ProductID = wrappedID.ID
	if err := services.ProductStock(r.ds.ProductStock).UpdateProductStock(updateProductStockCtx, productStockDto.Model()); err != nil {
		if err == bo.ErrProductStockNotFound {
			ctx.JSON(http.StatusNotFound, dto.Builder().SetMessage("productStock not found"))
			return
		}
		slog.Error("unable to update productStock", "cause", err)
		ctx.JSON(http.StatusInternalServerError, dto.Builder().SetMessage("Internal server error"))
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"message": "productStock updated"})
}

// DeleteProductStock godoc
// @Summary      Delete a ProductStock by id
// @Description  Delete a ProductStock by id
// @Tags         ProductStock
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "ProductStock ID"
// @Success      204  {string}  "ProductStock delete processed"
// @Failure      400  {string} 	string  "Invalid request body"
// @Failure      404  {object}  string  "ProductStock not found"
// @Failure      500  {string}  string  "Error"
// @Router       /v1/product-stock/{id} [delete]
func (r *repos) deleteProductStock(ctx *gin.Context) {
	var wrappedID dto.IDWrapper
	if err := ctx.ShouldBindUri(&wrappedID); err != nil {
		slog.Error("unable to parse productStock id", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid query value"))
		return
	}

	deleteProductStockCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := services.ProductStock(r.ds.ProductStock).DeleteProductStock(deleteProductStockCtx, wrappedID.ID); err != nil {
		if err == bo.ErrProductStockNotFound {
			ctx.JSON(http.StatusNotFound, dto.Builder().SetMessage("productStock not found"))
			return
		}
		slog.Error("unable to delete productStock", "cause", err)
		ctx.JSON(http.StatusInternalServerError, dto.Builder().SetMessage("Internal server error"))
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"message": "productStock deleted"})
}
