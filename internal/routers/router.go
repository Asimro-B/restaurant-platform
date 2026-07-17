package route

import (
	handler "restaurant-platform/internal/handlers"
	"restaurant-platform/internal/middleware"
	"restaurant-platform/internal/models"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *handler.WebHandler) {
	// Public
	r.POST("/auth/login", h.HandleLogin)
	r.POST("/tenant", h.CreateTenant)

	// Protected
	protected := r.Group("/")
	protected.Use(middleware.RequireAuth())
	{
		// Tenant management
		protected.GET("/tenants", middleware.RequireRole(models.RoleOwner, models.RoleManager), h.ListTenants)
		protected.GET("/tenants/:tenantID", middleware.RequireRole(models.RoleOwner, models.RoleManager), h.GetTenantByID)
		protected.GET("/tenant/slug/:slug", middleware.RequireRole(models.RoleOwner, models.RoleManager), h.GetTenantBySlug)
		protected.PUT("/tenants/:tenantID", middleware.RequireRole(models.RoleOwner), h.UpdateTenant)
		protected.DELETE("/tenants/:tenantID", middleware.RequireRole(models.RoleOwner), h.DeleteTenant)
		protected.PATCH("/tenants/:tenantID/restore", middleware.RequireRole(models.RoleOwner), h.RestoreTenant)

		// User management
		protected.GET("/users", middleware.RequireRole(models.RoleOwner, models.RoleManager), h.ListUsers)
		protected.GET("/users/:userID", middleware.RequireRole(models.RoleOwner, models.RoleManager), h.GetUserByID)
		protected.GET("/user/get-by-email", middleware.RequireRole(models.RoleOwner, models.RoleManager), h.GetUserByEmail)
		protected.PUT("/users/:userID", middleware.RequireRole(models.RoleOwner, models.RoleManager), h.UpdateUser)
		protected.POST("/users", middleware.RequireRole(models.RoleOwner), h.CreateUser)
		protected.DELETE("/users/:userID", middleware.RequireRole(models.RoleOwner), h.DeleteUser)
		protected.PATCH("/users/:userID/restore", middleware.RequireRole(models.RoleOwner), h.RestoreUser)

		// Menu management — write: owner/manager, read: all roles
		protected.POST("/menus", middleware.RequireRole(models.RoleOwner, models.RoleManager), h.CreateMenu)
		protected.PUT("/menus/:menuID", middleware.RequireRole(models.RoleOwner, models.RoleManager), h.UpdateMenu)
		protected.DELETE("/menus/:menuID", middleware.RequireRole(models.RoleOwner, models.RoleManager), h.DeleteMenu)
		protected.PATCH("/menus/:menuID/restore", middleware.RequireRole(models.RoleOwner, models.RoleManager), h.RestoreMenu)
		protected.GET("/menus", middleware.RequireRole(models.RoleOwner, models.RoleManager, models.RoleWaiter, models.RoleKitchen, models.RoleCashier), h.ListMenus)
		protected.GET("/menus/:menuID", middleware.RequireRole(models.RoleOwner, models.RoleManager, models.RoleWaiter, models.RoleKitchen, models.RoleCashier), h.GetMenuByID)

		// Menu categories
		protected.POST("/menus/:menuID/categories", middleware.RequireRole(models.RoleOwner, models.RoleManager), h.CreateMenuCategory)
		protected.PUT("/menus/:menuID/categories/:categoryID", middleware.RequireRole(models.RoleOwner, models.RoleManager), h.UpdateMenuCategory)
		protected.DELETE("/menus/:menuID/categories/:categoryID", middleware.RequireRole(models.RoleOwner, models.RoleManager), h.DeleteMenuCategory)
		protected.PATCH("/menus/:menuID/categories/:categoryID/restore", middleware.RequireRole(models.RoleOwner, models.RoleManager), h.RestoreMenuCategory)
		protected.GET("/menus/:menuID/categories", middleware.RequireRole(models.RoleOwner, models.RoleManager, models.RoleWaiter, models.RoleKitchen, models.RoleCashier), h.ListMenuCategories)
		protected.GET("/menus/:menuID/categories/:categoryID", middleware.RequireRole(models.RoleOwner, models.RoleManager, models.RoleWaiter, models.RoleKitchen, models.RoleCashier), h.GetMenuCategoryByID)

		// Menu items
		protected.POST("/menus/:menuID/categories/:categoryID/items", middleware.RequireRole(models.RoleOwner, models.RoleManager), h.CreateMenuItem)
		protected.PUT("/menus/:menuID/categories/:categoryID/items/:ID", middleware.RequireRole(models.RoleOwner, models.RoleManager), h.UpdateMenuItem)
		protected.DELETE("/menus/:menuID/categories/:categoryID/items/:ID", middleware.RequireRole(models.RoleOwner, models.RoleManager), h.DeleteMenuItem)
		protected.PATCH("/menus/:menuID/categories/:categoryID/items/:ID/restore", middleware.RequireRole(models.RoleOwner, models.RoleManager), h.RestoreMenuItem)
		protected.GET("/menus/:menuID/categories/:categoryID/items", middleware.RequireRole(models.RoleOwner, models.RoleManager, models.RoleWaiter, models.RoleKitchen, models.RoleCashier), h.ListMenuItems)
		protected.GET("/menus/:menuID/categories/:categoryID/items/:ID", middleware.RequireRole(models.RoleOwner, models.RoleManager, models.RoleWaiter, models.RoleKitchen, models.RoleCashier), h.GetMenuItemByID)

		// Tables
		protected.POST("/tables", middleware.RequireRole(models.RoleOwner, models.RoleManager), h.CreateTable)
		protected.PUT("/tables/:tableID", middleware.RequireRole(models.RoleOwner, models.RoleManager), h.UpdateTable)
		protected.DELETE("/tables/:tableID", middleware.RequireRole(models.RoleOwner, models.RoleManager), h.DeleteTable)
		protected.GET("/tables", middleware.RequireRole(models.RoleOwner, models.RoleManager, models.RoleWaiter, models.RoleCashier), h.ListTables)
		protected.GET("/tables/:tableID", middleware.RequireRole(models.RoleOwner, models.RoleManager, models.RoleWaiter, models.RoleCashier), h.GetTableByID)
		protected.PATCH("/tables/:tableID/status", middleware.RequireRole(models.RoleOwner, models.RoleManager, models.RoleWaiter), h.UpdateTableStatus)

		// Orders
		protected.GET("/orders", middleware.RequireRole(models.RoleOwner, models.RoleManager, models.RoleWaiter, models.RoleCashier), h.ListOrders)
		protected.POST("/tables/:tableID/orders", middleware.RequireRole(models.RoleOwner, models.RoleManager, models.RoleWaiter), h.CreateOrder)
		protected.GET("/orders/:referenceID/bill", middleware.RequireRole(models.RoleOwner, models.RoleManager, models.RoleWaiter, models.RoleCashier), h.GetBill)
		protected.PATCH("/orders/:referenceID/served", middleware.RequireRole(models.RoleOwner, models.RoleManager, models.RoleWaiter), h.OrderServed)
		protected.PATCH("/orders/:referenceID/kitchen-start", middleware.RequireRole(models.RoleOwner, models.RoleKitchen), h.KitchenStart)
		protected.PATCH("/orders/:referenceID/kitchen-done", middleware.RequireRole(models.RoleOwner, models.RoleKitchen), h.KitchenDone)
		protected.POST("/orders/:referenceID/pay", middleware.RequireRole(models.RoleOwner, models.RoleCashier), h.ProcessPayment)

		// Reservations
		protected.POST("/reservations", middleware.RequireRole(models.RoleOwner, models.RoleManager, models.RoleWaiter), h.CreateReservation)
		protected.GET("/reservations", middleware.RequireRole(models.RoleOwner, models.RoleManager), h.ListReservations)
		protected.GET("/reservations/:reservationID", middleware.RequireRole(models.RoleOwner, models.RoleManager, models.RoleWaiter), h.GetReservationByID)
		protected.PATCH("/reservations/:reservationID/confirm", middleware.RequireRole(models.RoleOwner, models.RoleManager), h.ConfirmReservation)
		protected.PATCH("/reservations/:reservationID/cancel", middleware.RequireRole(models.RoleOwner, models.RoleManager), h.CancelReservation)
		protected.PATCH("/reservations/:reservationID/complete", middleware.RequireRole(models.RoleOwner, models.RoleManager), h.CompleteReservation)
	}
}
