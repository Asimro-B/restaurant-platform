package route

import (
	handler "restaurant-platform/internal/handlers"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRegisterRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	api := r.Group("/api/v1")

	RegisterRoutes(api, handler.NewWebHandler(nil, nil))
}
