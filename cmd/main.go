// @title           Restaurant Management Platform API
// @version         1.0
// @description     A production-grade multi-tenant restaurant management backend built with Go. Features durable order workflows with Temporal, real-time notifications with Pusher, Redis caching, and role-based access control.

// @contact.name    Asimro
// @contact.url     https://github.com/yourusername/restaurant-platform

// @host            localhost:8000
// @BasePath        /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"log"
	"restaurant-platform/cmd/initiator"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env.local"); err != nil {
		log.Println(".env.local not found")
	}
	initiator.Initiate()
}
