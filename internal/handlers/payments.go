package handler

import (
	"fmt"
	"net/http"
	"restaurant-platform/internal/ctxutil"
	"restaurant-platform/internal/logger"
	"restaurant-platform/internal/models"
	pusherclient "restaurant-platform/internal/pusher"

	"github.com/gin-gonic/gin"
)

func (h *WebHandler) ProcessPayment(c *gin.Context) {
	ctx := c.Request.Context()
	referenceID := c.Param("referenceID")
	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("failed to get tenant from context: ", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	var req models.ProcessPaymentReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("failed to bind the request body", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	// Get order by referenceID
	order, err := h.module.GetOrderByReferenceID(ctx, referenceID, tenantCtx.TenantID)
	if err != nil {
		logger.Error("failed to get order by reference id: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	// Validate order status is "served" (can't pay unserved order)
	if order.Status != string(models.OrderStatusServed) {
		logger.Error("order must be served before paymen: ", err)
		models.ERROR(c, http.StatusBadRequest, fmt.Errorf("order must be served before payment, current status: %s", order.Status))
		return
	}

	// Validate payment_method is valid
	if !models.PaymentMethod.IsValid(models.PaymentMethod(req.PaymentMethod)) {
		logger.Error("Invalid Payment Method: ", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	// Insert payment record with status "completed"
	payment, err := h.module.ProcessPayment(ctx, tenantCtx.TenantID, order.ID, order.TotalAmount, models.ProcessPaymentReq{
		TableID:       order.TableID,
		PaymentMethod: req.PaymentMethod,
		Reference:     req.Reference,
		Notes:         req.Notes,
	})
	if err != nil {
		logger.Error("failed to process payment", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	// Update order status to "paid"
	if err := h.module.UpdateOrderStatus(ctx, order.ID, tenantCtx.TenantID, "paid"); err != nil {
		logger.Error("Ifailed to update order status: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	// Send Pusher event: order.paid on manager channel
	pusherclient.Publish(
		pusherclient.ManageChannel(tenantCtx.TenantID),
		"order.paid",
		map[string]interface{}{
			"reference_id": referenceID,
			"order_id":     order.ID,
			"tenant_id":    order.TenantID,
			"message":      "Order paid",
		},
	)

	_, err = h.module.UpdateTableStatus(ctx, order.TableID, order.TenantID, "available")
	if err != nil {
		logger.Error("failed to update table status: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data: gin.H{
			"payment_id": payment.ID,
			"order_id":   order.ID,
			"amount":     order.TotalAmount,
			"status":     "completed",
			"message":    "payment processed successfully",
		},
	})
}

func (h *WebHandler) GetBill(c *gin.Context) {
	ctx := c.Request.Context()
	referenceID := c.Param("referenceID")

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("failed to get tenant from context", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	bill, err := h.module.GetBill(ctx, referenceID, tenantCtx.TenantID)
	if err != nil {
		logger.Error("Failed to get bill: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data: bill,
	})
}
