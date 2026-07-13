package route

import (
	handler "restaurant-platform/internal/handlers"
	"restaurant-platform/internal/middleware"
	"restaurant-platform/internal/models"

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
		// owner and manager only
		admin := protected.Group("/")
		admin.Use(middleware.RequireRole(models.RoleManager))
		{
			// Tenant management
			admin.GET("/tenants", h.ListTenants)
			admin.GET("/tenants/:tenantID", h.GetTenantByID)
			admin.PATCH("/tenants/:tenantID/restore", h.RestoreTenant)
			admin.GET("/tenant/slug/:slug", h.GetTenantBySlug)

			// User management
			admin.GET("/users", h.ListUsers)
			admin.GET("/users/:userID", h.GetUserByID)
			admin.PUT("/users/:userID", h.UpdateUser)
			admin.GET("/user/get-by-email", h.GetUserByEmail)

			// Menu management
			admin.POST("/menus", h.CreateMenu)
			admin.PUT("/menus/:menuID", h.UpdateMenu)
			admin.DELETE("/menus/:menuID", h.DeleteMenu)
			admin.POST("/menus/:menuID/categories", h.CreateMenuCategory)
			admin.PUT("/menus/:menuID/categories/:categoryID", h.UpdateMenuCategory)
			admin.DELETE("/menus/:menuID/categories/:categoryID", h.DeleteMenuCategory)
			admin.POST("/menus/:menuID/categories/:categoryID/items", h.CreateMenuItem)
			admin.PUT("/menus/:menuID/categories/:categoryID/items/:ID", h.UpdateMenuItem)
			admin.DELETE("/menus/:menuID/categories/:categoryID/items/:ID", h.DeleteMenuItem)

			// Reservations management
			admin.GET("/reservations", h.ListReservations)
			admin.PATCH("/reservations/:reservationID/confirm", h.ConfirmReservation)
			admin.PATCH("/reservations/:reservationID/cancel", h.CancelReservation)
			admin.PATCH("/reservations/:reservationID/complete", h.CompleteReservation)
		}

		// waiter and above
		waiter := protected.Group("/")
		waiter.Use(middleware.RequireRole(models.RoleWaiter))
		{
			// Read menus (everyone needs to see the menu)
			waiter.GET("/menus", h.ListMenus)
			waiter.GET("/menus/:menuID", h.GetMenuByID)
			waiter.GET("/menus/:menuID/categories", h.ListMenuCategories)
			waiter.GET("/menus/:menuID/categories/:categoryID", h.GetMenuCategoryByID)
			waiter.GET("/menus/:menuID/categories/:categoryID/items", h.ListMenuItems)
			waiter.GET("/menus/:menuID/categories/:categoryID/items/:ID", h.GetMenuItemByID)

			// Tables
			waiter.GET("/tables", h.ListTables)
			waiter.GET("/tables/:tableID", h.GetTableByID)
			waiter.PATCH("/tables/:tableID/status", h.UpdateTableStatus)

			// Orders
			waiter.POST("/tables/:tableID/orders", h.CreateOrder)
			waiter.GET("/orders/:referenceID/bill", h.GetBill)
			waiter.PATCH("/orders/:referenceID/served", h.OrderServed)

			// Reservations
			waiter.POST("/reservations", h.CreateReservation)
			waiter.GET("/reservations/:reservationID", h.GetReservationByID)
		}

		// kitchen only
		kitchen := protected.Group("/")
		kitchen.Use(middleware.RequireRole(models.RoleKitchen))
		{
			kitchen.PATCH("/orders/:referenceID/kitchen-start", h.KitchenStart)
			kitchen.PATCH("/orders/:referenceID/kitchen-done", h.KitchenDone)
		}

		// casher and above
		cashier := protected.Group("/")
		cashier.Use(middleware.RequireRole(models.RoleCasher))
		{
			cashier.POST("/orders/:referenceID/pay", h.ProcessPayment)
		}

		// owener only
		owner := protected.Group("/")
		owner.Use(middleware.RequireRole(models.RoleOwner))
		{
			owner.PUT("/tenants/:tenantID", h.UpdateTenant)
			owner.DELETE("/tenants/:tenantID", h.DeleteTenant)
			owner.POST("/users", h.CreateUser)
			owner.DELETE("/users/:userID", h.DeleteUser)
		}

	}
}
