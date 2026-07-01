package handler

import (
	"net/http"
	"restaurant-platform/internal/logger"
	"restaurant-platform/internal/models"
	module "restaurant-platform/internal/modules"

	"github.com/gin-gonic/gin"
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
