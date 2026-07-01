package initiator

import (
	"os"
	handler "restaurant-platform/internal/handlers"
	"restaurant-platform/internal/logger"
	module "restaurant-platform/internal/modules"
	route "restaurant-platform/internal/routers"

	"github.com/gin-gonic/gin"
)

func Initiate() {
	// Initialize logger first
	logger.Init()
	logger.Log.Info("Restorant Platform Started")

	// database
	db, err := InitiatePersistenceDB()
	if err != nil {
		logger.Log.Error("failed to initialize database", "error", err)
		panic(err)
	}

	// service
	mod := module.NewModule(db)

	// handler
	webHandler := handler.NewWebHandler(mod)

	// route
	r := gin.Default()
	api := r.Group("/api/v1")
	route.RegisterRoutes(api, webHandler)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8000"
	}

	if err := r.Run(":" + port); err != nil {
		logger.Log.Error("failed to start server", "error", err)
	}

}
