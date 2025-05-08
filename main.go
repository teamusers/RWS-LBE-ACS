package main

import (
	"log"
	router "rlp-email-service/api"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router.Init()
}
