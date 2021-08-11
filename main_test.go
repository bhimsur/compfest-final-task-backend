package main_test

import (
	"os"
	"restgo/api/controllers"
	"testing"
)

var app controllers.App

func TestMain(m *testing.M) {
	app.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"))
	code := m.Run()
	os.Exit(code)
}
