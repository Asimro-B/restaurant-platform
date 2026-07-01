package route

import (
	handler "restaurant-platform/internal/handlers"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, h *handler.WebHandler) {
	r.POST("/tenant", h.CreateTenant)
	r.GET("/tenants/:tenantID/user/get-by-email", h.GetUserByEmail)
	r.POST("/tenants/:tenatID/user", h.CreateUser)
}
