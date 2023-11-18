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

// Get Suppliers godoc
// @Summary      Get Suppliers
// @Description  Get Suppliers
// @Tags         Supplier
// @Accept       json
// @Produce      json
// @Param        limit   query   int  false  "limit"
// @Param        offset  query   int  false  "offset"
// @Success      200  {object}  dto.PaginatedSupplierCollection
// @Failure      400  {string} string  "Invalid request body"
// @Failure      500  {string}  string  "Error"
// @Router       /v1/suppliers [get]
func (r *repos) getSuppliers(ctx *gin.Context) {
	var supplierQueryDto dto.SupplierQuery
	if err := ctx.ShouldBindQuery(&supplierQueryDto); err != nil {
		slog.Error("unable to parse query url", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid query value"))
		return
	}

	getSupplierCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pbc, err := services.Supplier(r.ds.Supplier).List(getSupplierCtx, supplierQueryDto.Model())
	if err != nil {
		slog.Error("unable to get suppliers", "cause", err)
		ctx.JSON(http.StatusInternalServerError, dto.Builder().SetMessage("Internal server error"))
		return
	}

	ctx.JSON(http.StatusOK, dto.ToPaginatedSupplier(pbc))
}

// Get Supplier godoc
// @Summary      Get a Supplier by id
// @Description  Get a Supplier by id
// @Tags         Supplier
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Supplier ID"
// @Success      200  {object}  dto.Supplier
// @Failure      400  {string} string  "Invalid request body"
// @Failure      404  {object}  dto.Error
// @Failure      500  {string}  string  "Error"
// @Router       /v1/supplier/{id} [get]
func (r *repos) getSupplier(ctx *gin.Context) {
	var wrappedID dto.IDWrapper
	if err := ctx.ShouldBindUri(&wrappedID); err != nil {
		slog.Error("unable to parse supplier id", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid query value"))
		return
	}

	getSupplierCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	supplier, err := services.Supplier(r.ds.Supplier).GetSupplierByID(getSupplierCtx, wrappedID.ID)
	if err != nil {
		if err == bo.ErrSupplierNotFound {
			ctx.JSON(http.StatusNotFound, dto.Builder().SetMessage("supplier not found"))
			return
		}
		slog.Error("unable to get supplier from database: ", "cause", err)
		ctx.JSON(http.StatusInternalServerError, dto.Builder().SetMessage("Error"))
		return
	}

	ctx.JSON(http.StatusOK, dto.ToSupplierDTO(supplier))
}

// Add Supplier godoc
// @Summary      Add a new Supplier
// @Description  Create a new Supplier in the system
// @Tags         Supplier
// @Accept       json
// @Produce      json
// @Param        request body dto.Supplier  true  "Supplier params"
// @Success      201  {object}  dto.IDWrapper
// @Failure      400  {string} string  "Invalid request body"
// @Failure      404  {object}  dto.Error
// @Failure      500  {string}  string  "Error"
// @Router       /v1/supplier [post]
func (r *repos) addSupplier(ctx *gin.Context) {
	supplierDto := dto.Supplier{}
	if err := ctx.ShouldBindJSON(&supplierDto); err != nil {
		slog.Error("unable to parse supplier from request body", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid request body"))
		return
	}

	model := supplierDto.Model()

	addSupplierCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	id, err := services.Supplier(r.ds.Supplier).CreateSupplier(addSupplierCtx, model)
	if err != nil {
		slog.Error("unable to create supplier", "cause", err)
		ctx.JSON(http.StatusInternalServerError, dto.Builder().SetMessage("Internal server error"))
		return
	}

	ctx.JSON(http.StatusCreated, dto.IDWrapper{ID: id})
}

// UpdateSupplier godoc
// @Summary      Update a Supplier by id
// @Description  Update a Supplier by id
// @Tags         Supplier
// @Accept       json
// @Produce      json
// @Param        request body dto.SupplierUpdate  true  "Supplier params"
// @Param        id   path      int  true  "Supplier ID"
// @Success      204  {string}  "SupplierDto updated"
// @Failure      400  {string} string  "Invalid request body"
// @Failure      404  {object}  dto.Error
// @Failure      500  {string}  string  "Error"
// @Router       /v1/supplier/{id} [patch]
func (r *repos) updateSupplier(ctx *gin.Context) {
	var wrappedID dto.IDWrapper
	if err := ctx.ShouldBindUri(&wrappedID); err != nil {
		slog.Error("unable to parse supplier id", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid query value"))
		return
	}

	var supplierDto dto.SupplierUpdate
	if err := ctx.ShouldBindJSON(&supplierDto); err != nil {
		slog.Error("unable to parse supplier from request body", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid request body"))
		return
	}

	updateSupplierCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	supplierDto.ID = wrappedID.ID
	if err := services.Supplier(r.ds.Supplier).UpdateSupplier(updateSupplierCtx, supplierDto.Model()); err != nil {
		if err == bo.ErrSupplierNotFound {
			ctx.JSON(http.StatusNotFound, dto.Builder().SetMessage("supplier not found"))
			return
		}
		slog.Error("unable to update supplier", "cause", err)
		ctx.JSON(http.StatusInternalServerError, dto.Builder().SetMessage("Internal server error"))
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"message": "supplier updated"})
}

// DeleteSupplier godoc
// @Summary      Delete a Supplier by id
// @Description  Delete a Supplier by id
// @Tags         Supplier
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Supplier ID"
// @Success      204  {string}  "Supplier delete processed"
// @Failure      400  {string} 	string  "Invalid request body"
// @Failure      404  {object}  string  "Supplier not found"
// @Failure      500  {string}  string  "Error"
// @Router       /v1/supplier/{id} [delete]
func (r *repos) deleteSupplier(ctx *gin.Context) {
	var wrappedID dto.IDWrapper
	if err := ctx.ShouldBindUri(&wrappedID); err != nil {
		slog.Error("unable to parse supplier id", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid query value"))
		return
	}

	deleteSupplierCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := services.Supplier(r.ds.Supplier).DeleteSupplier(deleteSupplierCtx, wrappedID.ID); err != nil {
		if err == bo.ErrSupplierNotFound {
			ctx.JSON(http.StatusNotFound, dto.Builder().SetMessage("supplier not found"))
			return
		}
		slog.Error("unable to delete supplier", "cause", err)
		ctx.JSON(http.StatusInternalServerError, dto.Builder().SetMessage("Internal server error"))
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"message": "supplier deleted"})
}
