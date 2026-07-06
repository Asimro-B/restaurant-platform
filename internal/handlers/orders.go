package handler

import (
	"net/http"
	"restaurant-platform/internal/ctxutil"
	"restaurant-platform/internal/logger"
	"restaurant-platform/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *WebHandler) CreateOrder(c *gin.Context) {
	ctx := c.Request.Context()

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	tableIDStr := c.Param("tableID")
	tableID, err := strconv.ParseInt(tableIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid table id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	userIDStr := c.Param("userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid user id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	var req models.CreateOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind request body", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	order, err := h.module.CreateOrder(ctx, models.CreateOrderReq{
		TenantID:    tenantCtx.TenantID,
		TableID:     tableID,
		UserID:      userID,
		Notes:       req.Notes,
		TotalAmount: req.TotalAmount,
		Status:      req.Status,
	})
	if err != nil {
		logger.Error("Failed to create order", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusCreated, models.Response{
		Data:  order,
		Error: nil,
	})
}

func (h *WebHandler) ListOrders(c *gin.Context) {
	ctx := c.Request.Context()

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	tableIDStr := c.Param("tableID")
	tableID, err := strconv.ParseInt(tableIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid table id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	userIDStr := c.Param("userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid user id", err)
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

	orders, total, err := h.module.ListOrders(ctx, models.ListOrdersReq{
		TenantID: tenantCtx.TenantID,
		TableID:  tableID,
		UserID:   userID,
		Page:     page,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		logger.Error("Failed to list orders", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data: orders,
		Pagination: &models.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			PageSize:   len(orders),
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
		Error: nil,
	})
}

func (h *WebHandler) GetOrderByID(c *gin.Context) {
	ctx := c.Request.Context()

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	orderIDStr := c.Param("ID")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid order id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	tableIDStr := c.Param("tableID")
	tableID, err := strconv.ParseInt(tableIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid table id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	userIDStr := c.Param("userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid user id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	order, err := h.module.GetOrderByID(ctx, orderID, tenantCtx.TenantID, tableID, userID)
	if err != nil {
		logger.Error("Failed to get order", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data:  order,
		Error: nil,
	})
}

func (h *WebHandler) UpdateOrderStatus(c *gin.Context) {
	ctx := c.Request.Context()

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	orderIDStr := c.Param("ID")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid order id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	tableIDStr := c.Param("tableID")
	tableID, err := strconv.ParseInt(tableIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid table id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	userIDStr := c.Param("userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid user id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	var req models.UpdateOrderStatusReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to bind request body", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	order, err := h.module.UpdateOrderStatus(ctx, models.UpdateOrderStatusReq{
		ID:       orderID,
		TenantID: tenantCtx.TenantID,
		TableID:  tableID,
		UserID:   userID,
		Status:   req.Status,
	})
	if err != nil {
		logger.Error("Failed to update order status", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data:  order,
		Error: nil,
	})
}

func (h *WebHandler) DeleteOrder(c *gin.Context) {
	ctx := c.Request.Context()

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	orderIDStr := c.Param("ID")
	orderID, err := strconv.ParseInt(orderIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid order id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	tableIDStr := c.Param("tableID")
	tableID, err := strconv.ParseInt(tableIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid table id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	userIDStr := c.Param("userID")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		logger.Error("Invalid user id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	if err := h.module.DeleteOrder(ctx, orderID, tenantCtx.TenantID, tableID, userID); err != nil {
		logger.Error("Failed to delete order", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data:  "Order Deleted Successfully",
		Error: nil,
	})
}
