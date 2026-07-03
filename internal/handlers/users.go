package handler

import (
	"errors"
	"net/http"
	"restaurant-platform/internal/logger"
	"restaurant-platform/internal/models"
	"restaurant-platform/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func (h *WebHandler) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()

	tenantID, err := parseTenantID(c)
	if err != nil {
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	var req models.CreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	req.TenantID = tenantID
	hashedPassword, err := utils.HashPassword(req.PasswordHash)
	if err != nil {
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	req.PasswordHash = hashedPassword

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

func (h *WebHandler) ListUsers(c *gin.Context) {
	ctx := c.Request.Context()

	tenantID, err := parseTenantID(c)
	if err != nil {
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
	users, err := h.module.ListUsers(ctx, models.ListUsersReq{
		TenantID:    tenantID,
		Search:      c.Query("search"),
		Role:        c.Query("role"),
		SortBy:      queryDefault(c, "sort_by", "created_at"),
		SortOrder:   queryDefault(c, "sort_order", "desc"),
		LimitCount:  int32(limit),
		OffsetCount: int32(offset),
	})
	if err != nil {
		logger.Error("failed to list users")
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data: users,
		Pagination: &models.PaginationMeta{
			Page:     page,
			Limit:    limit,
			PageSize: len(users),
			HasNext:  len(users) == limit,
			HasPrev:  page > 1,
		},
		Error: nil,
	})
}

func (h *WebHandler) GetUserByEmail(c *gin.Context) {
	ctx := c.Request.Context()

	tenantID, err := parseTenantID(c)
	if err != nil {
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	var req models.GetUserByEmailReq
	if err := c.ShouldBindJSON(&req); err != nil {
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	req.TenantID = int(tenantID)

	user, err := h.module.GetUserByEmail(ctx, req.Email, int64(req.TenantID))
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

func (h *WebHandler) GetUserByID(c *gin.Context) {
	ctx := c.Request.Context()

	tenantID, userID, err := parseUserRouteIDs(c)
	if err != nil {
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	user, err := h.module.GetUserByID(ctx, userID, tenantID)
	if err != nil {
		handleUserError(c, err, "failed to get user by id")
		return
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data:  user,
		Error: nil,
	})
}

func (h *WebHandler) UpdateUser(c *gin.Context) {
	ctx := c.Request.Context()

	tenantID, userID, err := parseUserRouteIDs(c)
	if err != nil {
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	var req models.UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}
	req.ID = userID
	req.TenantID = tenantID

	hashedPassword, err := utils.HashPassword(req.PasswordHash)
	if err != nil {
		models.ERROR(c, http.StatusInternalServerError, err)
		return
	}
	req.PasswordHash = hashedPassword

	user, err := h.module.UpdateUser(ctx, req)
	if err != nil {
		handleUserError(c, err, "failed to update user")
		return
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data:  user,
		Error: nil,
	})
}

func (h *WebHandler) DeleteUser(c *gin.Context) {
	ctx := c.Request.Context()

	tenantID, userID, err := parseUserRouteIDs(c)
	if err != nil {
		models.ERROR(c, http.StatusBadRequest, err)
		return
	}

	user, err := h.module.DeleteUser(ctx, userID, tenantID)
	if err != nil {
		handleUserError(c, err, "failed to delete user")
		return
	}

	models.JSON(c, http.StatusOK, models.Response{
		Data:  user,
		Error: nil,
	})
}

func parseUserRouteIDs(c *gin.Context) (int64, int64, error) {
	tenantID, err := parseTenantID(c)
	if err != nil {
		return 0, 0, err
	}

	userID, err := strconv.ParseInt(c.Param("userID"), 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return tenantID, userID, nil
}

func queryDefault(c *gin.Context, key string, defaultValue string) string {
	value := c.Query(key)
	if value == "" {
		return defaultValue
	}

	return value
}

func handleUserError(c *gin.Context, err error, message string) {
	if errors.Is(err, pgx.ErrNoRows) {
		models.ERROR(c, http.StatusNotFound, err)
		return
	}

	logger.Error(message)
	models.ERROR(c, http.StatusInternalServerError, err)
}
