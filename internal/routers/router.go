package route

import (
	handler "restaurant-platform/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *handler.WebHandler) {
	r.POST("/tenant", h.CreateTenant)
	r.GET("/tenant/slug/:slug", h.GetTenantBySlug)
	r.GET("/tenants", h.ListTenants)
	r.GET("/tenants/:tenantID", h.GetTenantByID)
	r.PUT("/tenants/:tenantID", h.UpdateTenant)
	r.DELETE("/tenants/:tenantID", h.DeleteTenant)
	r.PATCH("/tenants/:tenantID/restore", h.RestoreTenant)
	r.GET("/tenants/:tenantID/user/get-by-email", h.GetUserByEmail)
	r.POST("/tenants/:tenatID/user", h.CreateUser)
}
