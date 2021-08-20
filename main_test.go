package main_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
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
	app.Initialize()
	code := m.Run()
	os.Exit(code)
}

func TestLoginSuccess(t *testing.T) {
	var jsonStr = []byte(`{"username":"doni","password":"doni"}`)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestLoginWithoutPassword(t *testing.T) {
	var jsonStr = []byte(`{"username":"doni"}`)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestLoginWithoutUsername(t *testing.T) {
	var jsonStr = []byte(`{"password":"doni"}`)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestLoginWrongPassword(t *testing.T) {
	var jsonStr = []byte(`{"username":"doni","password":"doni2"}`)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusForbidden, response.Code)
}

func TestLoginUserNotFound(t *testing.T) {
	var jsonStr = []byte(`{"username":"doni2","password":"doni"}`)
	req, _ := http.NewRequest("POST", "/auth/login", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)
	checkResponseCode(t, http.StatusInternalServerError, response.Code)
}

func TestGetAllDonationProgramSuccess(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/donate", nil)
	req.Header.Set("Authorization", os.Getenv("TOKEN_ADMIN"))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetAllDonationProgramFailed(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/donate", nil)
	req.Header.Set("Authorization", os.Getenv("SECRET"))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusForbidden, response.Code)
}

func TestGetDonationProgramByIdExists(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/donate/1", nil)
	req.Header.Set("Authorization", os.Getenv("TOKEN_DONOR"))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetDonationProgramByIdNotExists(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/donate/1000", nil)
	req.Header.Set("Authorization", os.Getenv("TOKEN_DONOR"))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusInternalServerError, response.Code)
}

func TestGetDonationProgramHistoryFromFundraiserSuccess(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/donate/history", nil)
	req.Header.Set("Authorization", os.Getenv("TOKEN_FUNDRAISER"))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetDonationProgramHistoryFromFundraiserFailed(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/donate/history", nil)
	req.Header.Set("Authorization", os.Getenv("TOKEN_DONOR"))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

func TestVerifyDonationProgramAlready(t *testing.T) {
	req, _ := http.NewRequest("PUT", "/api/donate/verify/7", nil)
	req.Header.Set("Authorization", os.Getenv("TOKEN_ADMIN"))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestCreateDonationProgramSuccess(t *testing.T) {
	var jsonStr = []byte(`{"title":"ini title", "detail":"ini detail", "amount":150000}`)
	req, _ := http.NewRequest("POST", "/api/donate", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", os.Getenv("TOKEN_FUNDRAISER"))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)
}

func TestCreateDonationProgramFailed(t *testing.T) {
	var jsonStr = []byte(`{"title":"ini title", "detail":"ini detail", "amount":150000}`)
	req, _ := http.NewRequest("POST", "/api/donate", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", os.Getenv("TOKEN_DONOR"))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestGetDonationHistorySuccess(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/user/donate/history", nil)
	req.Header.Set("Authorization", os.Getenv("TOKEN_DONOR"))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetDonationHistoryError(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/user/donate/history", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusForbidden, response.Code)
}

func TestGetWalletByUserIdSuccess(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/wallet", nil)
	req.Header.Set("Authorization", os.Getenv("TOKEN_DONOR"))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetWalletByUserIdInvalid(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/wallet", nil)
	req.Header.Set("Authorization", os.Getenv("TOKEN_ADMIN"))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusInternalServerError, response.Code)
}

func TestGetWalletByUserIdError(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/wallet", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusForbidden, response.Code)
}

func TestVerifyFundraiserAccountAlready(t *testing.T) {
	req, _ := http.NewRequest("PUT", "/api/user/verify/4", nil)
	req.Header.Set("Authorization", os.Getenv("TOKEN_ADMIN"))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestVerifyFundraiserAccountUnauthorized(t *testing.T) {
	req, _ := http.NewRequest("PUT", "/api/user/verify/4", nil)
	req.Header.Set("Authorization", os.Getenv("TOKEN_DONOR"))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusUnauthorized, response.Code)
}

func executeRequest(r *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	app.Router.ServeHTTP(rr, r)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
