package handler

import (
	"net/http"
	"restaurant-platform/errkit"
	"restaurant-platform/internal/logger"
	"restaurant-platform/internal/models"
	"strings"

	utils "restaurant-platform/utils"

	"github.com/gin-gonic/gin"
)

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	TenantID int64  `json:"tenant_id,omitempty"`
}

// HandleLogin godoc
// @Summary      Login
// @Description  Authenticate a staff member and return a JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login credentials"
// @Success      200  {object}  models.Response{data=string}
// @Failure      400  {object}  models.ErrorResponse
// @Failure      401  {object}  models.ErrorResponse
// @Router       /auth/login [post]
func (h *WebHandler) HandleLogin(c *gin.Context) {
	ctx := c.Request.Context()

	// Parse request body
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		logger.Error("Failed to decode request body", err)
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	if req.Email == "" || req.Password == "" {
		models.ERROR(c, http.StatusBadRequest, errkit.ErrInvalidData)
		return
	}
	email := strings.ToLower(req.Email)

	// Fetch ALL users with this email across all tenants
	user, err := h.module.GetUserByEmail(ctx, email, req.TenantID)
	if err != nil {
		logger.Error("user not found", err)
		models.ERROR(c, http.StatusUnauthorized, errkit.ErrUserNotFound)
		return
	}

	// verify password
	if err := utils.CheckPassword(req.Password, user.PasswordHash); err != nil {
		models.ERROR(c, http.StatusUnauthorized, errkit.ErrInvalidCredentials)
		return
	}

	// Generate JWT
	token, err := utils.GenerateToken(int64(user.ID), int64(user.TenantID), user.Role)
	if err != nil {
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusOK, models.Response{Data: token, Error: nil})
}
