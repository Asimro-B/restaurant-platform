package initiator

import (
	"os"
	handler "restaurant-platform/internal/handlers"
	"restaurant-platform/internal/logger"
	module "restaurant-platform/internal/modules"
	pusherclient "restaurant-platform/internal/pusher"
	route "restaurant-platform/internal/routers"
	"restaurant-platform/internal/worker"
	orderworkflow "restaurant-platform/internal/workflows/order"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Initiate() {
	// Initialize logger first
	logger.Init()
	logger.Log.Info("Restorant Platform Started")

	pusherclient.Init()
	logger.Log.Info("Pusher initialized",
		"app_id", os.Getenv("PUSHER_APP_ID"),
		"cluster", os.Getenv("PUSHER_CLUSTER"),
	)

	// database
	db, err := InitiatePersistenceDB()
	if err != nil {
		logger.Log.Error("failed to initialize database", "error", err)
		panic(err)
	}

	temporalClient, err := InitiateTemporalClient()
	if err != nil {
		logger.Log.Error("failed to initialize temporal", "error", err)
		panic(err)
	}
	defer temporalClient.Close()

	// service
	mod := module.NewModule(db)

	activities := orderworkflow.NewOrderActivities(mod)
	go worker.StartWorker(temporalClient, activities)

	// handler
	webHandler := handler.NewWebHandler(mod, temporalClient)

	// route
	r := gin.Default()
	// CORS
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))
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
