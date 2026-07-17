package handler

import (
	"fmt"
	"net/http"
	"restaurant-platform/internal/ctxutil"
	"restaurant-platform/internal/logger"
	"restaurant-platform/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *WebHandler) CreateTable(c *gin.Context) {
	ctx := c.Request.Context()

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	var req models.CreateTableReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("failed to bind the request body", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	table, err := h.module.CreateTable(ctx, models.CreateTableReq{
		TenantID: tenantCtx.TenantID,
		Name:     req.Name,
		Capacity: req.Capacity,
	})
	if err != nil {
		logger.Error("failed to create table", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusCreated, models.Response{Data: table, Error: nil})
}

// ListTables godoc
// @Summary      List tables
// @Description  Get all tables for the authenticated tenant
// @Tags         tables
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        page   query  int  false  "Page number"  default(1)
// @Param        limit  query  int  false  "Items per page"  default(10)
// @Success      200  {object}  models.Response{data=[]models.Table}
// @Failure      401  {object}  models.ErrorResponse
// @Failure      403  {object}  models.ErrorResponse
// @Router       /tables [get]
func (h *WebHandler) ListTables(c *gin.Context) {
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
	statusFilter := c.Query("status") // empty string if not provided — sqlc query handles it

	tables, total, err := h.module.ListTables(ctx, models.ListTablesReq{
		TenantID: tenantCtx.TenantID,
		Column2:  statusFilter,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		logger.Error("failed to list tables", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}
	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data: tables,
		Pagination: &models.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			PageSize:   len(tables),
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
		Error: nil,
	})
}

func (h *WebHandler) GetTableByID(c *gin.Context) {
	ctx := c.Request.Context()

	IDStr := c.Param("ID")

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}
	id, err := strconv.ParseInt(IDStr, 10, 64)
	if err != nil {
		logger.Error("Inalid table id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	table, err := h.module.GetTableByID(ctx, id, tenantCtx.TenantID)
	if err != nil {
		logger.Error("Failed to get the table", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data:  table,
		Error: nil,
	})
}

func (h *WebHandler) UpdateTable(c *gin.Context) {
	ctx := c.Request.Context()

	IDStr := c.Param("ID")

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}
	id, err := strconv.ParseInt(IDStr, 10, 64)
	if err != nil {
		logger.Error("Inalid table id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	var req models.UpdateTableReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind json", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	table, err := h.module.UpdateTable(ctx, models.UpdateTableReq{
		Name:     req.Name,
		Capacity: req.Capacity,
		ID:       id,
		TenantID: tenantCtx.TenantID,
	})
	if err != nil {
		logger.Error("Failed to update table", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusOK, models.Response{Data: table, Error: nil})
}

func (h *WebHandler) DeleteTable(c *gin.Context) {
	ctx := c.Request.Context()

	IDStr := c.Param("ID")

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}
	id, err := strconv.ParseInt(IDStr, 10, 64)
	if err != nil {
		logger.Error("Inalid table id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	if err := h.module.DeleteTable(ctx, id, tenantCtx.TenantID); err != nil {
		logger.Error("Failed to update table", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusOK, models.Response{Data: "table Deleted Successfully", Error: nil})
}

func (h *WebHandler) UpdateTableStatus(c *gin.Context) {
	ctx := c.Request.Context()

	tableIDStr := c.Param("tableID")
	tableID, err := strconv.ParseInt(tableIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid table id: ", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind the request body", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	if !models.TableStatus(req.Status).IsValid() {
		logger.Error("invalid status: must be one of available, occupied, reserved, inactive", err)
		models.ERROR(c, http.StatusBadRequest, fmt.Errorf("invalid status: must be one of available, occupied, reserved, inactive"))
		return
	}

	table, err := h.module.UpdateTableStatus(ctx, tableID, tenantCtx.TenantID, req.Status)
	if err != nil {
		logger.Error("Failed to update the table status: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data:  table,
		Error: nil,
	})

}
