package main_test

import (
	"log"
	"os"
	"restgo/api/controllers"
	"testing"

	"github.com/joho/godotenv"
)

var app controllers.App

func TestMain(m *testing.M) {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	app.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"))
	code := m.Run()
	os.Exit(code)
}
