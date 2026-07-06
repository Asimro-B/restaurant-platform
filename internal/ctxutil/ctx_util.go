package ctxutil

import (
	"fmt"
	"restaurant-platform/internal/models"

	"github.com/gin-gonic/gin"
)

// GetTenantFromContext retrieves the tenant information from the context
func GetTenantFromContext(c *gin.Context) (models.TenantUserContext, error) {
	tenant, ok := c.Request.Context().Value(models.TenantContextKey).(models.TenantUserContext)
	if ok {
		return tenant, nil
	}
	return models.TenantUserContext{}, fmt.Errorf("tenant not found on context")
}
