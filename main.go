package main

import (
	"os"
	"restgo/api/controllers"

	_ "github.com/lib/pq"
)

func main() {
	// if err := godotenv.Load(); err != nil {
	// 	log.Fatal("Error loading .env file")
	// }
	app := controllers.App{}
	app.Initialize(
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PASSWORD"))
	app.RunServer()
}
