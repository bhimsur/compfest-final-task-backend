package main

import (
	"log"
	"restgo/api/controllers"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	app := controllers.App{}
	app.Initialize()
	app.RunServer()
}
