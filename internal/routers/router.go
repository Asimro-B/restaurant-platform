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
		protected.GET("/users", h.ListUsers)
		protected.GET("/users/:userID", h.GetUserByID)
		protected.PUT("/users/:userID", h.UpdateUser)
		protected.DELETE("/users/:userID", h.DeleteUser)
		protected.GET("/user/get-by-email", h.GetUserByEmail)
		protected.POST("/user", h.CreateUser)

		// Menu Management
		protected.POST("/menus", h.CreateMenu)
		protected.GET("/menus", h.ListMenus)
		protected.GET("/menus/:ID", h.GetMenuByID)
		protected.PUT("/menus/:ID", h.UpdateMenu)
		protected.DELETE("/menus/:ID", h.DeleteMenu)

		// =========================
		// Menu Categories
		// =========================
		protected.POST("/menus/:menuID/categories", h.CreateMenuCategory)
		protected.GET("/menus/:menuID/categories", h.ListMenuCategories)
		protected.GET("/menus/:menuID/categories/:ID", h.GetMenuCategoryByID)
		protected.PUT("/menus/:menuID/categories/:ID", h.UpdateMenuCategory)
		protected.DELETE("/menus/:menuID/categories/:ID", h.DeleteMenuCategory)

		// =========================
		// Menu Items
		// =========================
		protected.POST("/menus/:menuID/categories/:categoryID/items", h.CreateMenuItem)
		protected.GET("/menus/:menuID/categories/:categoryID/items", h.ListMenuItems)
		protected.GET("/menus/:menuID/categories/:categoryID/items/:ID", h.GetMenuItemByID)
		protected.PUT("/menus/:menuID/categories/:categoryID/items/:ID", h.UpdateMenuItem)
		protected.DELETE("/menus/:menuID/categories/:categoryID/items/:ID", h.DeleteMenuItem)

		// Table Management
		protected.POST("/tables", h.CreateTable)
		protected.GET("/tables", h.ListTables)
		protected.GET("/tables/:tableID", h.GetTableByID)
		protected.PUT("/tables/:tableID", h.UpdateTable)
		protected.DELETE("/tables/:tableID", h.DeleteTable)
		protected.PATCH("/tables/:tableID/status", h.UpdateTableStatus)

		// =========================
		// Order Management
		// =========================
		protected.POST("/tables/:tableID/users/:userID/orders", h.CreateOrder)
		protected.GET("/tables/:tableID/users/:userID/orders", h.ListOrders)
		protected.GET("/tables/:tableID/users/:userID/orders/:ID", h.GetOrderByID)
		protected.PATCH("/tables/:tableID/users/:userID/orders/:ID/status", h.UpdateOrderStatus)
		protected.DELETE("/tables/:tableID/users/:userID/orders/:ID", h.DeleteOrder)
	}
}
