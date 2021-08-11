package tests

import (
	"net/http"
	"net/http/httptest"
	"os"
	"restgo/api/controllers"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

var a controllers.App

func TestDonationProgram(m *testing.M) {
	a = controllers.App{}
	a.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_PASSWORD"),
	)
	code := m.Run()
	os.Exit(code)
}
func TestGetDonationProgram(t *testing.T) {
	addDonationPrograms(1)
	req, _ := http.NewRequest("GET", "/api/donate", nil)
	response := executeRequest(req)
	// checkResponseCode(t, http.StatusOK, response.Code)
	assert.Equal(t, http.StatusOK, response.Code)
}

func executeRequest(r *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, r)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	} else {
		assert.Equal(t, expected, actual)
	}
}

func addDonationPrograms(count int) {
	if count < 1 {
		count = 1
	}
	for i := 0; i < count; i++ {
		a.DB.Exec("INSERT INTO donation_programs(title, detail, amount, user_id) VALUES($1, $2, $3, $4)", "Donation "+strconv.Itoa(i), "Donation detail "+strconv.Itoa(i), (i+1)*10, i)
	}
}
