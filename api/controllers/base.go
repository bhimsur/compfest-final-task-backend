package controllers

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"restgo/api/middlewares"
	"restgo/api/models"
	"restgo/api/responses"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (a *App) Initialize() {
	var err error

	a.DB, err = gorm.Open("postgres", os.Getenv("HEROKU_DB_URI"))
	if err != nil {
		fmt.Printf("\n Cannot connect to database")
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the databse")
	}

	a.DB.Debug().AutoMigrate(&models.User{}, &models.DonationProgram{}, &models.Donation{}, &models.Wallet{}, &models.TopUp{})

	a.Router = mux.NewRouter().StrictSlash(true)
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.Use(middlewares.SetResponsesMiddleware)
	a.Router.HandleFunc("/", home).Methods("GET")
	u := a.Router.PathPrefix("/auth").Subrouter()
	u.HandleFunc("/register", a.UserSignUp).Methods("POST", "OPTIONS")
	u.HandleFunc("/login", a.Login).Methods("POST", "OPTIONS")

	s := a.Router.PathPrefix("/api").Subrouter()
	s.Use(middlewares.AuthJwtVerify)
	s.HandleFunc("/donate", a.GetDonationPrograms).Methods("GET")
	s.HandleFunc("/donate", a.CreateDonationProgram).Methods("POST", "OPTIONS")
	s.HandleFunc("/donate/{id:[0-9]+}", a.DonateNow).Methods("POST", "OPTIONS")
	s.HandleFunc("/donate/{id:[0-9]+}", a.GetDonationProgramById).Methods("GET")
	s.HandleFunc("/donate/history", a.GetDonationProgramByFundraiser).Methods("GET")
	s.HandleFunc("/donate/verify/{id:[0-9]+}", a.VerifyDonationProgram).Methods("PUT")

	s.HandleFunc("/user/verify/{id:[0-9]+}", a.VerifyFundraiser).Methods("PUT")
	s.HandleFunc("/user/donate/history", a.GetDonationHistoryFromUser).Methods("GET")
	s.HandleFunc("/user", a.GetUserById).Methods("GET")
	s.HandleFunc("/user", a.UpdateUser).Methods("PUT")

	//wallet
	s.HandleFunc("/wallet", a.GetWalletByUserId).Methods("GET")
	s.HandleFunc("/wallet", a.CreateTopUp).Methods("POST", "OPTIONS")
	s.HandleFunc("/wallet/history", a.GetTopUpHistoryByUserId).Methods("GET")
}

func (a *App) RunServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Printf("\nServer starting on port " + port)
	err := http.ListenAndServe(":"+port, a.Router)
	if err != nil {
		fmt.Print(err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome to Fundraising")
}
