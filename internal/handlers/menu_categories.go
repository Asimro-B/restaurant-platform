package handler

import (
	"net/http"
	"restaurant-platform/internal/ctxutil"
	"restaurant-platform/internal/logger"
	"restaurant-platform/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *WebHandler) CreateMenuCategory(c *gin.Context) {
	ctx := c.Request.Context()

	menuIDStr := c.Param("menuID")

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	menuID, err := strconv.ParseInt(menuIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid menu id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	var req models.CreateMenuCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind request body", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	category, err := h.module.CreateMenuCategory(ctx, models.CreateMenuCategoryReq{
		TenantID:    tenantCtx.TenantID,
		MenuID:      menuID,
		Name:        req.Name,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		IsActive:    req.IsActive,
	})
	if err != nil {
		logger.Error("Failed to create menu category", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusCreated, models.Response{
		Data:  category,
		Error: nil,
	})
}

func (h *WebHandler) ListMenuCategories(c *gin.Context) {
	ctx := c.Request.Context()

	menuIDStr := c.Param("menuID")

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	menuID, err := strconv.ParseInt(menuIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid menu id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	page, err := parsePositiveQueryInt(c, "page", 1)
	if err != nil {
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	limit, err := parsePositiveQueryInt(c, "limit", 10)
	if err != nil {
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	offset := (page - 1) * limit

	categories, total, err := h.module.ListMenuCategories(ctx, models.ListMenuCategoriesReq{
		TenantID: tenantCtx.TenantID,
		MenuID:   menuID,
		Page:     page,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		logger.Error("Failed to list menu categories", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data: categories,
		Pagination: &models.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			PageSize:   len(categories),
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
		Error: nil,
	})
}

func (h *WebHandler) GetMenuCategoryByID(c *gin.Context) {
	ctx := c.Request.Context()

	menuIDStr := c.Param("menuID")
	idStr := c.Param("categoryID")

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}
	menuID, err := strconv.ParseInt(menuIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid menu id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("Invalid category id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	category, err := h.module.GetMenuCategoryByID(ctx, tenantCtx.TenantID, menuID, id)
	if err != nil {
		logger.Error("Failed to get menu category", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data:  category,
		Error: nil,
	})
}

func (h *WebHandler) UpdateMenuCategory(c *gin.Context) {
	ctx := c.Request.Context()

	menuIDStr := c.Param("menuID")
	idStr := c.Param("categoryID")

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	menuID, err := strconv.ParseInt(menuIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid menu id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("Invalid category id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	var req models.UpdateMenuCategoryReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind request body", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	category, err := h.module.UpdateMenuCategory(ctx, models.UpdateMenuCategoryReq{
		ID:          id,
		MenuID:      menuID,
		TenantID:    tenantCtx.TenantID,
		Name:        req.Name,
		Description: req.Description,
		SortOrder:   req.SortOrder,
		IsActive:    req.IsActive,
	})
	if err != nil {
		logger.Error("Failed to update menu category", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data:  category,
		Error: nil,
	})
}

func (h *WebHandler) DeleteMenuCategory(c *gin.Context) {
	ctx := c.Request.Context()

	menuIDStr := c.Param("menuID")
	idStr := c.Param("categoryID")

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	menuID, err := strconv.ParseInt(menuIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid menu id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("Invalid category id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	if err := h.module.DeleteMenuCategory(ctx, id, menuID, tenantCtx.TenantID); err != nil {
		logger.Error("Failed to delete menu category", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data:  "Menu Category Deleted Successfully",
		Error: nil,
	})
}
