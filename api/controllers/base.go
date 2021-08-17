package controllers

import (
	"fmt"
	"log"
	"net/http"
	"restgo/api/middlewares"
	"restgo/api/models"
	"restgo/api/responses"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type App struct {
	Router *mux.Router
	DB     *gorm.DB
}

func (a *App) Initialize(DbHost, DbPort, DbUser, DbName, DbPassword string) {
	var err error

	DBURI := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)

	a.DB, err = gorm.Open("postgres", DBURI)
	if err != nil {
		fmt.Printf("\n Cannot connect to database %s", DbName)
		log.Fatal("This is the error:", err)
	} else {
		fmt.Printf("We are connected to the databse %s", DbName)
	}

	a.DB.Debug().AutoMigrate(&models.User{}, &models.DonationProgram{}, &models.Donation{}, &models.Wallet{}, &models.TopUp{}, &models.Withdrawal{})

	a.Router = mux.NewRouter().StrictSlash(true)
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.Use(middlewares.SetContentTypeMiddleware)
	a.Router.HandleFunc("/", home).Methods("GET")
	a.Router.HandleFunc("/register", a.UserSignUp).Methods("POST")
	a.Router.HandleFunc("/login", a.Login).Methods("POST")

	s := a.Router.PathPrefix("/api").Subrouter()
	s.Use(middlewares.AuthJwtVerify)
	s.HandleFunc("/donate", a.GetDonationPrograms).Methods("GET")
	s.HandleFunc("/donate", a.CreateDonationProgram).Methods("POST")
	s.HandleFunc("/donate/{id:[0-9]+}", a.GetDonationProgramById).Methods("GET")
	s.HandleFunc("/donate/{id:[0-9]+}", a.DonateToProgram).Methods("POST")
	s.HandleFunc("/donate/history", a.GetDonationProgramByFundraiser).Methods("GET")
	s.HandleFunc("/donate/verify/{id:[0-9]+}", a.VerifyDonationProgram).Methods("PUT")
	s.HandleFunc("/donate/unverified", a.GetUnverifiedDonationProgram).Methods("GET")

	s.HandleFunc("/user/verify/{id:[0-9]+}", a.VerifyFundraiser).Methods("PUT")
	s.HandleFunc("/user/donate/history", a.GetDonationHistoryFromUser).Methods("GET")
	s.HandleFunc("/user/unverified", a.GetUnverifiedUser).Methods("GET")

	s.HandleFunc("/withdraw/verify/{id:[0-9]+}", a.VerifyWithdrawal).Methods("PUT")
	s.HandleFunc("/withdraw/{id:[0-9]+}", a.CreateWithdrawal).Methods("POST")

	s.HandleFunc("/withdraw/unverified", a.GetUnverifiedWithdrawal).Methods("GET")

	//wallet
	s.HandleFunc("/wallet", a.GetWalletByUserId).Methods("GET")
	s.HandleFunc("/wallet", a.CreateTopUp).Methods("POST")
	s.HandleFunc("/wallet/history", a.TopupHistory).Methods("GET")
}

func (a *App) RunServer() {
	log.Printf("\nServer starting on port 5000")
	log.Fatal(http.ListenAndServe(":5000", handlers.CORS()(a.Router)))
}

func home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome to Fundraising")
}
