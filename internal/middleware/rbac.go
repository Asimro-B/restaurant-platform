package middleware

import (
	"fmt"
	"net/http"
	"restaurant-platform/internal/ctxutil"
	"restaurant-platform/internal/models"

	"github.com/gin-gonic/gin"
)

// // role hierarchy: higher index have more permissions
// var roleHierarchy = map[models.UserRole]int{
// 	models.RoleCasher:  1,
// 	models.RoleKitchen: 2,
// 	models.RoleWaiter:  3,
// 	models.RoleManager: 4,
// 	models.RoleOwner:   5,
// }

// // require role returns middleware that checks if the user have at least the required role level

// RequireRole checks if the user has ONE OF the allowed roles
func RequireRole(allowed ...models.UserRole) gin.HandlerFunc {
	return func(c *gin.Context) {
		tenantCtx, err := ctxutil.GetTenantFromContext(c)
		if err != nil {
			models.ERROR(c, http.StatusUnauthorized, err)
			c.Abort()
			return
		}

		userRole := models.UserRole(tenantCtx.Role)
		for _, role := range allowed {
			if userRole == role {
				c.Next()
				return
			}
		}

		models.ERROR(c, http.StatusForbidden, fmt.Errorf("insufficient permission"))
		c.Abort()
	}
}
