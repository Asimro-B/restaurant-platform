package handler

import (
	"net/http"
	"restaurant-platform/internal/ctxutil"
	"restaurant-platform/internal/logger"
	"restaurant-platform/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *WebHandler) CreateMenu(c *gin.Context) {
	ctx := c.Request.Context()

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	var req models.CreateMenuReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("failed to bind the request body", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	menu, err := h.module.CreateMenu(ctx, models.CreateMenuReq{
		TenantID:    tenantCtx.TenantID,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
	})
	if err != nil {
		logger.Error("failed to create menu", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusCreated, models.Response{Data: menu, Error: nil})
}

func (h *WebHandler) ListMenus(c *gin.Context) {
	ctx := c.Request.Context()

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
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

	menus, total, err := h.module.ListMenus(ctx, models.ListMenusReq{
		TenantID: tenantCtx.TenantID,
		Page:     page,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		logger.Error("failed to list munus", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}
	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data: menus,
		Pagination: &models.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			PageSize:   len(menus),
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
		Error: nil,
	})
}

func (h *WebHandler) GetMenuByID(c *gin.Context) {
	ctx := c.Request.Context()

	IDStr := c.Param("menuID")

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}
	id, err := strconv.ParseInt(IDStr, 10, 64)
	if err != nil {
		logger.Error("Inalid menu id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	menu, err := h.module.GetMenuByID(ctx, id, tenantCtx.TenantID)
	if err != nil {
		logger.Error("Failed to get the menu", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data:  menu,
		Error: nil,
	})
}

func (h *WebHandler) UpdateMenu(c *gin.Context) {
	ctx := c.Request.Context()

	IDStr := c.Param("menuID")

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}
	id, err := strconv.ParseInt(IDStr, 10, 64)
	if err != nil {
		logger.Error("Inalid menu id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	var req models.UpdateMenuReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind json", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	menu, err := h.module.UpdateMenu(ctx, models.UpdateMenuReq{
		ID:          id,
		TenantID:    tenantCtx.TenantID,
		Name:        req.Name,
		Description: req.Description,
		IsActive:    req.IsActive,
	})
	if err != nil {
		logger.Error("Failed to update menu", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusOK, models.Response{Data: menu, Error: nil})
}

func (h *WebHandler) DeleteMenu(c *gin.Context) {
	ctx := c.Request.Context()

	IDStr := c.Param("menuID")

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}
	id, err := strconv.ParseInt(IDStr, 10, 64)
	if err != nil {
		logger.Error("Inalid menu id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	if err := h.module.DeleteMenu(ctx, id, tenantCtx.TenantID); err != nil {
		logger.Error("Failed to delete menu", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusOK, models.Response{Data: "Menu Deleted Successfully", Error: nil})
}

func (h *WebHandler) RestoreMenu(c *gin.Context) {
	ctx := c.Request.Context()

	IDStr := c.Param("menuID")

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}
	id, err := strconv.ParseInt(IDStr, 10, 64)
	if err != nil {
		logger.Error("Inalid menu id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	if err := h.module.RestoreMenu(ctx, id, tenantCtx.TenantID); err != nil {
		logger.Error("Failed to restore menu", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusOK, models.Response{Data: "Menu Restored Successfully", Error: nil})
}
