package handler

import (
	"fmt"
	"net/http"
	"restaurant-platform/internal/logger"
	"restaurant-platform/internal/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h *WebHandler) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()
	tenantIDStr := c.Param("tenatID")

	tenant_id, err := strconv.ParseInt(tenantIDStr, 10, 64)
	if err != nil {
		fmt.Println("Error during conversion:", err)
		return
	}

	var req models.CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	req.TenantID = tenant_id

	user, err := h.module.CreateUser(ctx, req)
	if err != nil {
		logger.Error("failed to create user")
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusCreated, models.Response{
		Data:  user,
		Error: nil,
	})
}

func (h *WebHandler) GetUserByEmail(c *gin.Context) {
	ctx := c.Request.Context()

	tenantIdStr := c.Param("tenantID")
	tenantId, err := strconv.ParseInt(tenantIdStr, 10, 64)
	if err != nil {
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	var req models.GetUserByEmailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	req.TenantID = int(tenantId)

	user, err := h.module.GetUserByEmail(ctx, req)
	if err != nil {
		logger.Error("failed to get user by email")
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data:  user,
		Error: nil,
	})
}
