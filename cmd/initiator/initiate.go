package initiator

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"restaurant-platform/internal/cache"
	handler "restaurant-platform/internal/handlers"
	"restaurant-platform/internal/logger"
	module "restaurant-platform/internal/modules"
	pusherclient "restaurant-platform/internal/pusher"
	route "restaurant-platform/internal/routers"
	"restaurant-platform/internal/worker"
	orderworkflow "restaurant-platform/internal/workflows/order"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Initiate() {
	// Initialize logger first
	logger.Init()
	logger.Log.Info("Restorant Platform Started")

	cache.Init()
	logger.Log.Info("Redis Initialized")

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

	// Create the HTTP server instance
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}

	// Start the server inside a separate goroutine
	go func() {
		log.Println("Server is starting on port", port)
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("server failed to listen: %v\n", err)
		}
	}()

	// Create a channel to listen for OS signals
	// SIGINT = Ctrl+C, SIGTERM = Sent by Kubernetes or Docker to terminate pods/containers
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	// Block the main goroutine untill a signal is recieved
	<-shutdownSignal
	log.Println("shotdown signal recieved. starting graceful termination...")

	// create a context with a timeout for the shutdown process
	// this prevents the application from hanging indefinitely if a request stalls
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Trigger the graceful shutdown
	// this stops accepting new connections and flushes current connection
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// final cleanup(closing db connections, flushing logs)
	log.Println("Cleanup complete, server exiting cleanly")
}
