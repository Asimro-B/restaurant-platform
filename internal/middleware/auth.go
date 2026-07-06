package middleware

import (
	"context"
	"fmt"
	"net/http"
	"restaurant-platform/internal/models"
	"restaurant-platform/utils"
	"strings"

	"github.com/gin-gonic/gin"
)

func RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			models.ERROR(c, http.StatusUnauthorized, fmt.Errorf("missing token"))
			c.Abort()
			return
		}
		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := utils.ValidateToken(tokenStr)
		if err != nil {
			models.ERROR(c, http.StatusUnauthorized, fmt.Errorf("Invalid token"))
			c.Abort()
			return
		}

		// Set full context object on request context
		tenantCtx := models.TenantUserContext{
			UserID:   claims.UserID,
			TenantID: claims.TenantID,
			Role:     claims.Role,
		}
		ctx := context.WithValue(c.Request.Context(), models.TenantContextKey, tenantCtx)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
