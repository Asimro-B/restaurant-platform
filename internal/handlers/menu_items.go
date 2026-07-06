package handler

import (
	"net/http"
	"restaurant-platform/internal/ctxutil"
	"restaurant-platform/internal/logger"
	"restaurant-platform/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *WebHandler) CreateMenuItem(c *gin.Context) {
	ctx := c.Request.Context()

	categoryIDStr := c.Param("categoryID")
	menuIDStr := c.Param("menuID")
	menuID, err := strconv.ParseInt(menuIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid menu id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid category id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	var req models.CreateMenuItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind request body", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	item, err := h.module.CreateMenuItem(ctx, models.CreateMenuItemReq{
		TenantID:    tenantCtx.TenantID,
		MenuID:      menuID,
		CategoryID:  categoryID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		IsAvailable: req.IsAvailable,
	})
	if err != nil {
		logger.Error("Failed to create menu item", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusCreated, models.Response{
		Data:  item,
		Error: nil,
	})
}

func (h *WebHandler) ListMenuItems(c *gin.Context) {
	ctx := c.Request.Context()

	categoryIDStr := c.Param("categoryID")
	menuIDStr := c.Param("menuID")
	menuID, err := strconv.ParseInt(menuIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid menu id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid category id", err)
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

	items, total, err := h.module.ListMenuItems(ctx, models.ListMenuItemsReq{
		TenantID:   tenantCtx.TenantID,
		MenuID:     menuID,
		CategoryID: categoryID,
		Page:       page,
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		logger.Error("Failed to list menu items", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data: items,
		Pagination: &models.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			PageSize:   len(items),
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
		Error: nil,
	})
}

func (h *WebHandler) GetMenuItemByID(c *gin.Context) {
	ctx := c.Request.Context()

	categoryIDStr := c.Param("categoryID")
	idStr := c.Param("ID")
	menuIDStr := c.Param("menuID")
	menuID, err := strconv.ParseInt(menuIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid menu id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid category id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("Invalid menu item id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	item, err := h.module.GetMenuItemByID(ctx, tenantCtx.TenantID, menuID, categoryID, id)
	if err != nil {
		logger.Error("Failed to get menu item", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data:  item,
		Error: nil,
	})
}

func (h *WebHandler) UpdateMenuItem(c *gin.Context) {
	ctx := c.Request.Context()

	categoryIDStr := c.Param("categoryID")
	idStr := c.Param("ID")
	menuIDStr := c.Param("menuID")
	menuID, err := strconv.ParseInt(menuIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid menu id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid category id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("Invalid menu item id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	var req models.UpdateMenuItemReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind request body", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	item, err := h.module.UpdateMenuItem(ctx, models.UpdateMenuItemReq{
		ID:          id,
		CategoryID:  categoryID,
		MenuID:      menuID,
		TenantID:    tenantCtx.TenantID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		IsAvailable: req.IsAvailable,
	})
	if err != nil {
		logger.Error("Failed to update menu item", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data:  item,
		Error: nil,
	})
}

func (h *WebHandler) DeleteMenuItem(c *gin.Context) {
	ctx := c.Request.Context()

	categoryIDStr := c.Param("categoryID")
	idStr := c.Param("ID")
	menuIDStr := c.Param("menuID")
	menuID, err := strconv.ParseInt(menuIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid menu id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	categoryID, err := strconv.ParseInt(categoryIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid category id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		logger.Error("Invalid menu item id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	if err := h.module.DeleteMenuItem(ctx, id, categoryID, menuID, tenantCtx.TenantID); err != nil {
		logger.Error("Failed to delete menu item", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data:  "Menu Item Deleted Successfully",
		Error: nil,
	})
}
