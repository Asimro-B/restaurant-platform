package route

import (
	handler "restaurant-platform/internal/handlers"
	"restaurant-platform/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *handler.WebHandler) {
	// Public no token needed
	r.POST("/auth/login", h.HandleLogin)
	r.POST("/tenant", h.CreateTenant)

	// protected token required
	protected := r.Group("/")
	protected.Use(middleware.RequireAuth())
	{
		// Tenant management
		protected.GET("/tenant/slug/:slug", h.GetTenantBySlug)
		protected.GET("/tenants", h.ListTenants)
		protected.GET("/tenants/:tenantID", h.GetTenantByID)
		protected.PUT("/tenants/:tenantID", h.UpdateTenant)
		protected.DELETE("/tenants/:tenantID", h.DeleteTenant)
		protected.PATCH("/tenants/:tenantID/restore", h.RestoreTenant)

		// User management
		protected.GET("/tenants/:tenantID/users", h.ListUsers)
		protected.GET("/tenants/:tenantID/users/:userID", h.GetUserByID)
		protected.PUT("/tenants/:tenantID/users/:userID", h.UpdateUser)
		protected.DELETE("/tenants/:tenantID/users/:userID", h.DeleteUser)
		protected.GET("/tenants/:tenantID/user/get-by-email", h.GetUserByEmail)
		protected.POST("/tenants/:tenantID/user", h.CreateUser)
	}
}
