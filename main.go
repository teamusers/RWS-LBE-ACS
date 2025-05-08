package main

import (
	router "rlp-email-service/api"
)

func main() {
	//err := godotenv.Load()
	//if err != nil {
	//	log.Fatal("Error loading .env file")
	//}
	//topic.StartSubscription()

	router.Init()
}
