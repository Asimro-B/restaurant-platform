package handler

import (
	"errors"
	"net/http"
	"restaurant-platform/internal/logger"
	"restaurant-platform/internal/models"
	module "restaurant-platform/internal/modules"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type WebHandler struct {
	module *module.WebModule
}

func NewWebHandler(webModule *module.WebModule) *WebHandler {
	return &WebHandler{
		module: webModule,
	}
}

func (h *WebHandler) CreateTenant(c *gin.Context) {
	ctx := c.Request.Context()

	var req models.CreateTenantReq
	if err := c.ShouldBindJSON(&req); err != nil {
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	tenant, err := h.module.CreateTenant(ctx, req)
	if err != nil {
		logger.Error("failed to create tenant")
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusCreated, models.Response{
		Data:  tenant,
		Error: nil,
	})
}

func (h *WebHandler) ListTenants(c *gin.Context) {
	ctx := c.Request.Context()

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
	tenants, total, err := h.module.ListTenants(ctx, models.ListTenantsReq{
		Page:   page,
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		logger.Error("failed to list tenants")
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data: tenants,
		Pagination: &models.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			PageSize:   len(tenants),
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
		Error: nil,
	})
}

func (h *WebHandler) GetTenantByID(c *gin.Context) {
	ctx := c.Request.Context()

	tenantID, err := parseTenantID(c)
	if err != nil {
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	tenant, err := h.module.GetTenantByID(ctx, tenantID)
	if err != nil {
		handleTenantError(c, err, "failed to get tenant by id")
		return
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data:  tenant,
		Error: nil,
	})
}

func (h *WebHandler) GetTenantBySlug(c *gin.Context) {
	ctx := c.Request.Context()
	slug := c.Param("slug")

	tenant, err := h.module.GetTenantBySlug(ctx, slug)
	if err != nil {
		handleTenantError(c, err, "failed to get tenant by slug")
		return
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data:  tenant,
		Error: nil,
	})
}

func (h *WebHandler) UpdateTenant(c *gin.Context) {
	ctx := c.Request.Context()

	tenantID, err := parseTenantID(c)
	if err != nil {
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	var req models.UpdateTenantReq
	if err := c.ShouldBindJSON(&req); err != nil {
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}
	req.ID = tenantID

	tenant, err := h.module.UpdateTenant(ctx, req)
	if err != nil {
		handleTenantError(c, err, "failed to update tenant")
		return
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data:  tenant,
		Error: nil,
	})
}

func (h *WebHandler) DeleteTenant(c *gin.Context) {
	ctx := c.Request.Context()

	tenantID, err := parseTenantID(c)
	if err != nil {
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	tenant, err := h.module.DeleteTenant(ctx, tenantID)
	if err != nil {
		handleTenantError(c, err, "failed to delete tenant")
		return
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data:  tenant,
		Error: nil,
	})
}

func (h *WebHandler) RestoreTenant(c *gin.Context) {
	ctx := c.Request.Context()

	tenantID, err := parseTenantID(c)
	if err != nil {
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	tenant, err := h.module.RestoreTenant(ctx, tenantID)
	if err != nil {
		handleTenantError(c, err, "failed to restore tenant")
		return
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data:  tenant,
		Error: nil,
	})
}

func parseTenantID(c *gin.Context) (int64, error) {
	return strconv.ParseInt(c.Param("tenantID"), 10, 64)
}

func parsePositiveQueryInt(c *gin.Context, key string, defaultValue int) (int, error) {
	value := c.Query(key)
	if value == "" {
		return defaultValue, nil
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	if parsed < 1 {
		return 0, errors.New(key + " must be greater than 0")
	}

	return parsed, nil
}

func handleTenantError(c *gin.Context, err error, message string) {
	if errors.Is(err, pgx.ErrNoRows) {
		models.ERROR(c, http.StatusNotFound, err)
		return
	}

	logger.Error(message)
	models.ERROR(c, http.StatusInternalServerError, err)
}
