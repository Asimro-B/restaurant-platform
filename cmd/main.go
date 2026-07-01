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
