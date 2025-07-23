package main

import (
	"context"
	"log"

	"github.com/felipeazsantos/concurrency-golang-fullcycle-lab03/configuration/database/mongodb"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("cmd/auction/.env"); err != nil {
		log.Fatal("Error trying to load env variables")
	}

	ctx := context.Background()

	_, err := mongodb.NewMongoDBConnection(ctx)
	if err != nil {
		log.Fatal(err.Error())
	}
}
