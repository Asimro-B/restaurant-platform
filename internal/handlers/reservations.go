package handler

import (
	"net/http"
	"restaurant-platform/internal/ctxutil"
	"restaurant-platform/internal/logger"
	"restaurant-platform/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *WebHandler) CreateReservation(c *gin.Context) {
	ctx := c.Request.Context()

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("failed to get tenant from the context, ", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	var req models.CreateReservationReq
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("failed to bind the request body, ", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	reservation, err := h.module.CreateReservation(ctx, tenantCtx.TenantID, models.CreateReservationReq{
		TableID:       req.TableID,
		CustomerName:  req.CustomerName,
		CustomerPhone: req.CustomerPhone,
		PartySize:     req.PartySize,
		ReservedAt:    req.ReservedAt,
		Notes:         req.Notes,
	})
	if err != nil {
		logger.Error("failed to create reservation, ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusCreated, models.Response{Data: reservation})
}

func (h *WebHandler) ListReservations(c *gin.Context) {
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

	statusFilter := c.Query("status")

	reservations, total, err := h.module.ListReservations(ctx, models.ListReservationsReq{
		TenantID: tenantCtx.TenantID,
		Status:   statusFilter,
		Page:     page,
		Limit:    limit,
		Offset:   offset,
	})
	if err != nil {
		logger.Error("failed to list reservations", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}
	totalPages := 0
	if total > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit))
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data: reservations,
		Pagination: &models.PaginationMeta{
			Page:       page,
			Limit:      limit,
			Total:      total,
			TotalPages: totalPages,
			PageSize:   len(reservations),
			HasNext:    page < totalPages,
			HasPrev:    page > 1,
		},
		Error: nil,
	})
}

func (h *WebHandler) GetReservationByID(c *gin.Context) {
	ctx := c.Request.Context()

	IDStr := c.Param("reservationID")

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}
	id, err := strconv.ParseInt(IDStr, 10, 64)
	if err != nil {
		logger.Error("Inalid Reservation id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	reservation, err := h.module.GetReservationByID(ctx, id, tenantCtx.TenantID)
	if err != nil {
		logger.Error("Failed to get the Reservation", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data:  reservation,
		Error: nil,
	})
}

func (h *WebHandler) ConfirmReservation(c *gin.Context) {
	ctx := c.Request.Context()

	IDStr := c.Param("reservationID")

	id, err := strconv.ParseInt(IDStr, 10, 64)
	if err != nil {
		logger.Error("Inalid table id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	reservation, err := h.module.UpdateReservationStatus(ctx, id, tenantCtx.TenantID, models.ReservationStatusConfirmed)
	if err != nil {
		logger.Error("Failed to update Status", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusOK, models.Response{Data: reservation, Error: nil})
}

func (h *WebHandler) CancelReservation(c *gin.Context) {
	ctx := c.Request.Context()

	IDStr := c.Param("reservationID")

	id, err := strconv.ParseInt(IDStr, 10, 64)
	if err != nil {
		logger.Error("Inalid table id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	reservation, err := h.module.UpdateReservationStatus(ctx, id, tenantCtx.TenantID, models.ReservationStatusCancelled)
	if err != nil {
		logger.Error("Failed to update Status", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	_, err = h.module.UpdateTableStatus(ctx, reservation.TableID, reservation.TenantID, string(models.TableStatusAvailable))
	if err != nil {
		logger.Error("failed to update table status, ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusOK, models.Response{Data: reservation, Error: nil})
}

func (h *WebHandler) CompleteReservation(c *gin.Context) {
	ctx := c.Request.Context()

	IDStr := c.Param("reservationID")

	id, err := strconv.ParseInt(IDStr, 10, 64)
	if err != nil {
		logger.Error("Inalid table id", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	tenantCtx, err := ctxutil.GetTenantFromContext(c)
	if err != nil {
		logger.Error("Failed to get the user from the context: ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	reservation, err := h.module.UpdateReservationStatus(ctx, id, tenantCtx.TenantID, models.ReservationStatusCompleted)
	if err != nil {
		logger.Error("Failed to update Status", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	_, err = h.module.UpdateTableStatus(ctx, reservation.TableID, reservation.TenantID, string(models.TableStatusOccupied))
	if err != nil {
		logger.Error("failed to update table status, ", err)
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusOK, models.Response{Data: reservation, Error: nil})
}
