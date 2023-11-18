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

// Get Categories godoc
// @Summary      Get categories
// @Description  Get categories
// @Tags         Category
// @Accept       json
// @Produce      json
// @Success      200  {object}  dto.CategoriesTree
// @Failure      400  {string} string  "Invalid request body"
// @Failure      500  {string}  string  "Error"
// @Router       /v1/categories [get]
func (r *repos) getCategories(ctx *gin.Context) {
	getCtgCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pcc, err := services.Category(r.ds.Category).List(getCtgCtx)
	if err != nil {
		slog.Error("unable to get categories", "cause", err)
		ctx.JSON(http.StatusInternalServerError, dto.Builder().SetMessage("Internal server error"))
		return
	}

	ctx.JSON(http.StatusOK, dto.CategoriesToTree(pcc))
}

// Get Category godoc
// @Summary      Get a Category by id
// @Description  Get a Category by id
// @Tags         Category
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      200  {object}  dto.Category
// @Failure      400  {string} string  "Invalid request body"
// @Failure      404  {object}  dto.Error
// @Failure      500  {string}  string  "Error"
// @Router       /v1/category/{id} [get]
func (r *repos) getCategory(ctx *gin.Context) {
	var wrappedID dto.IDWrapper
	if err := ctx.ShouldBindUri(&wrappedID); err != nil {
		slog.Error("unable to parse category id", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid query value"))
		return
	}

	getCategoryCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	category, err := services.Category(r.ds.Category).GetCategoryByID(getCategoryCtx, wrappedID.ID)
	if err != nil {
		if err == bo.ErrCategoryNotFound {
			ctx.JSON(http.StatusNotFound, dto.Builder().SetMessage("category not found"))
			return
		}
		slog.Error("unable to get category from database: ", "cause", err)
		ctx.JSON(http.StatusInternalServerError, dto.Builder().SetMessage("Error"))
		return
	}

	ctx.JSON(http.StatusOK, dto.ToCategoryDTO(category))
}

// Add Category godoc
// @Summary      Add a new Category
// @Description  Create a new Category in the system
// @Tags         Category
// @Accept       json
// @Produce      json
// @Param        request body dto.Category  true  "Category params"
// @Success      201  {object}  dto.IDWrapper
// @Failure      400  {string} string  "Invalid request body"
// @Failure      404  {object}  dto.Error
// @Failure      500  {string}  string  "Error"
// @Router       /v1/category [post]
func (r *repos) addCategory(ctx *gin.Context) {
	categoryDto := dto.Category{}
	if err := ctx.ShouldBindJSON(&categoryDto); err != nil {
		slog.Error("unable to parse category from request body", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid request body"))
		return
	}

	model := categoryDto.Model()

	addCategoryCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	id, err := services.Category(r.ds.Category).CreateCategory(addCategoryCtx, model)
	if err != nil {
		slog.Error("unable to create category", "cause", err)
		ctx.JSON(http.StatusInternalServerError, dto.Builder().SetMessage("Internal server error"))
		return
	}

	ctx.JSON(http.StatusCreated, dto.IDWrapper{ID: id})
}

// UpdateCategory godoc
// @Summary      Update a Category by id
// @Description  Update a Category by id
// @Tags         Category
// @Accept       json
// @Produce      json
// @Param        request body dto.CategoryUpdate  true  "Category params"
// @Param        id   path      int  true  "Category ID"
// @Success      204  {string}  "CategoryDto updated"
// @Failure      400  {string} string  "Invalid request body"
// @Failure      404  {object}  dto.Error
// @Failure      500  {string}  string  "Error"
// @Router       /v1/category/{id} [patch]
func (r *repos) updateCategory(ctx *gin.Context) {
	var wrappedID dto.IDWrapper
	if err := ctx.ShouldBindUri(&wrappedID); err != nil {
		slog.Error("unable to parse category id", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid query value"))
		return
	}

	var categoryDto dto.CategoryUpdate
	if err := ctx.ShouldBindJSON(&categoryDto); err != nil {
		slog.Error("unable to parse category from request body", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid request body"))
		return
	}

	updateCategoryCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	categoryDto.ID = wrappedID.ID
	if err := services.Category(r.ds.Category).UpdateCategory(updateCategoryCtx, categoryDto.Model()); err != nil {
		if err == bo.ErrCategoryNotFound {
			ctx.JSON(http.StatusNotFound, dto.Builder().SetMessage("category not found"))
			return
		}
		slog.Error("unable to update category", "cause", err)
		ctx.JSON(http.StatusInternalServerError, dto.Builder().SetMessage("Internal server error"))
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"message": "category updated"})
}

// DeleteCategory godoc
// @Summary      Delete a Category by id
// @Description  Delete a Category by id
// @Tags         Category
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Category ID"
// @Success      204  {string}  "Category delete processed"
// @Failure      400  {string} 	string  "Invalid request body"
// @Failure      404  {object}  string  "Category not found"
// @Failure      500  {string}  string  "Error"
// @Router       /v1/category/{id} [delete]
func (r *repos) deleteCategory(ctx *gin.Context) {
	var wrappedID dto.IDWrapper
	if err := ctx.ShouldBindUri(&wrappedID); err != nil {
		slog.Error("unable to parse category id", "cause", err)
		ctx.JSON(http.StatusBadRequest, dto.Builder().SetMessage("Invalid query value"))
		return
	}

	deleteCategoryCtx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := services.Category(r.ds.Category).DeleteCategory(deleteCategoryCtx, wrappedID.ID); err != nil {
		if err == bo.ErrCategoryNotFound {
			ctx.JSON(http.StatusNotFound, dto.Builder().SetMessage("category not found"))
			return
		}
		slog.Error("unable to delete category", "cause", err)
		ctx.JSON(http.StatusInternalServerError, dto.Builder().SetMessage("Internal server error"))
		return
	}

	ctx.JSON(http.StatusNoContent, gin.H{"message": "category deleted"})
}
