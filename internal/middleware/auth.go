package middleware

import (
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

		// Inject claims into context for handlers to use
		c.Set("user_id", claims.UserID)
		c.Set("tenant_id", claims.TenantID)
		c.Set("role", claims.Role)

		c.Next()
	}
}
