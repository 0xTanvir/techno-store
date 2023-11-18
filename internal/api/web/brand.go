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

// Get Brands godoc
// @Summary      Get Brands
// @Description  Get Brands
// @Tags         Brand
// @Accept       json
// @Produce      json
// @Param        limit   query   int  false  "limit"
// @Param        offset  query   int  false  "offset"
// @Success      200  {object}  dto.PaginatedBrandCollection
// @Failure      400  {string} string  "Invalid request body"
// @Failure      500  {string}  string  "Error"
// @Router       /v1/brands [get]
func (r *repos) getBrands(ctx *gin.Context) {
	var brandQueryDto dto.BrandQuery
	if err := ctx.ShouldBindQuery(&brandQueryDto); err != nil {
		slog.Error("unable to parse query url", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid query value"))
		return
	}

	getBrandCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pbc, err := services.Brand(r.ds.Brand).List(getBrandCtx, brandQueryDto.Model())
	if err != nil {
		slog.Error("unable to get brands", "cause", err)
		ctx.JSON(http.StatusInternalServerError, dto.Builder().SetMessage("Internal server error"))
		return
	}

	ctx.JSON(http.StatusOK, dto.ToPaginatedBrand(pbc))
}

// Get Brand godoc
// @Summary      Get a Brand by id
// @Description  Get a Brand by id
// @Tags         Brand
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Brand ID"
// @Success      200  {object}  dto.Brand
// @Failure      400  {string} string  "Invalid request body"
// @Failure      404  {object}  dto.Error
// @Failure      500  {string}  string  "Error"
// @Router       /v1/brand/{id} [get]
func (r *repos) getBrand(ctx *gin.Context) {
	var wrappedID dto.IDWrapper
	if err := ctx.ShouldBindUri(&wrappedID); err != nil {
		slog.Error("unable to parse brand id", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid query value"))
		return
	}

	getBrandCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	brand, err := services.Brand(r.ds.Brand).GetBrandByID(getBrandCtx, wrappedID.ID)
	if err != nil {
		if err == bo.ErrBrandNotFound {
			ctx.JSON(http.StatusNotFound, dto.Builder().SetMessage("brand not found"))
			return
		}
		slog.Error("unable to get brand from database: ", "cause", err)
		ctx.JSON(http.StatusInternalServerError, dto.Builder().SetMessage("Error"))
		return
	}

	ctx.JSON(http.StatusOK, dto.ToBrandDTO(brand))
}

// Add Brand godoc
// @Summary      Add a new Brand
// @Description  Create a new Brand in the system
// @Tags         Brand
// @Accept       json
// @Produce      json
// @Param        request body dto.Brand  true  "Brand params"
// @Success      201  {object}  dto.IDWrapper
// @Failure      400  {string} string  "Invalid request body"
// @Failure      404  {object}  dto.Error
// @Failure      500  {string}  string  "Error"
// @Router       /v1/brand [post]
func (r *repos) addBrand(ctx *gin.Context) {
	brandDto := dto.Brand{}
	if err := ctx.ShouldBindJSON(&brandDto); err != nil {
		slog.Error("unable to parse brand from request body", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid request body"))
		return
	}

	model := brandDto.Model()

	addBrandCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	id, err := services.Brand(r.ds.Brand).CreateBrand(addBrandCtx, model)
	if err != nil {
		slog.Error("unable to create brand", "cause", err)
		ctx.JSON(http.StatusInternalServerError, dto.Builder().SetMessage("Internal server error"))
		return
	}

	ctx.JSON(http.StatusCreated, dto.IDWrapper{ID: id})
}

// UpdateBrand godoc
// @Summary      Update a Brand by id
// @Description  Update a Brand by id
// @Tags         Brand
// @Accept       json
// @Produce      json
// @Param        request body dto.BrandUpdate  true  "Brand params"
// @Param        id   path      int  true  "Brand ID"
// @Success      204  {string}  "BrandDto updated"
// @Failure      400  {string} string  "Invalid request body"
// @Failure      404  {object}  dto.Error
// @Failure      500  {string}  string  "Error"
// @Router       /v1/brand/{id} [patch]
func (r *repos) updateBrand(ctx *gin.Context) {
	var wrappedID dto.IDWrapper
	if err := ctx.ShouldBindUri(&wrappedID); err != nil {
		slog.Error("unable to parse brand id", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid query value"))
		return
	}

	var brandDto dto.BrandUpdate
	if err := ctx.ShouldBindJSON(&brandDto); err != nil {
		slog.Error("unable to parse brand from request body", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid request body"))
		return
	}

	updateBrandCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	brandDto.ID = wrappedID.ID
	if err := services.Brand(r.ds.Brand).UpdateBrand(updateBrandCtx, brandDto.Model()); err != nil {
		if err == bo.ErrBrandNotFound {
			ctx.JSON(http.StatusNotFound, dto.Builder().SetMessage("brand not found"))
			return
		}
		slog.Error("unable to update brand", "cause", err)
		ctx.JSON(http.StatusInternalServerError, dto.Builder().SetMessage("Internal server error"))
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"message": "brand updated"})
}

// DeleteBrand godoc
// @Summary      Delete a Brand by id
// @Description  Delete a Brand by id
// @Tags         Brand
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Brand ID"
// @Success      204  {string}  "Brand delete processed"
// @Failure      400  {string} 	string  "Invalid request body"
// @Failure      404  {object}  string  "Brand not found"
// @Failure      500  {string}  string  "Error"
// @Router       /v1/brand/{id} [delete]
func (r *repos) deleteBrand(ctx *gin.Context) {
	var wrappedID dto.IDWrapper
	if err := ctx.ShouldBindUri(&wrappedID); err != nil {
		slog.Error("unable to parse brand id", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid query value"))
		return
	}

	deleteBrandCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := services.Brand(r.ds.Brand).DeleteBrand(deleteBrandCtx, wrappedID.ID); err != nil {
		if err == bo.ErrBrandNotFound {
			ctx.JSON(http.StatusNotFound, dto.Builder().SetMessage("brand not found"))
			return
		}
		slog.Error("unable to delete brand", "cause", err)
		ctx.JSON(http.StatusInternalServerError, dto.Builder().SetMessage("Internal server error"))
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"message": "brand deleted"})
}
