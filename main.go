package main

import (
	"restgo/api/controllers"

	_ "github.com/lib/pq"
)

func main() {
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	app := controllers.App{}
	app.Initialize()
	app.RunServer()
}
