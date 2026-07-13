package middleware

import (
	"fmt"
	"net/http"
	"restaurant-platform/internal/ctxutil"
	"restaurant-platform/internal/models"

	"github.com/gin-gonic/gin"
)

// role hierarchy: higher index have more permissions
var roleHierarchy = map[models.UserRole]int{
	models.RoleCasher:  1,
	models.RoleKitchen: 2,
	models.RoleWaiter:  3,
	models.RoleManager: 4,
	models.RoleOwner:   5,
}

// require role returns middleware that checks if the user have at least the required role level
func RequireRole(required models.UserRole) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tenantCtx, err := ctxutil.GetTenantFromContext(ctx)
		if err != nil {
			models.ERROR(ctx, http.StatusUnauthorized, err)
			ctx.Abort()
			return
		}

		userLevel := roleHierarchy[models.UserRole(tenantCtx.Role)]
		requiredLevel := roleHierarchy[required]

		if userLevel < requiredLevel {
			models.ERROR(ctx, http.StatusForbidden, fmt.Errorf("Insufficient permissions: required %s or above", required))
			ctx.Abort()
			return
		}
		ctx.Next()
	}
}
