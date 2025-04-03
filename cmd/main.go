package main

import (
	"log"

	"github.com/TheDummyUser/goRss/cmd/api"
	"github.com/TheDummyUser/goRss/database"
)

func main() {
	db, err := database.ConnectionDb()

	if err != nil {
		log.Fatal("database connection failed")
	}

	cfg := api.NewServer(db)

	cfg.Listen(":8080")

	if err := cfg.Listen(":8080"); err != nil {
		panic(err)
	}
}
