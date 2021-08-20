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
	"github.com/rs/cors"
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
		fmt.Printf("We are connected to the database")
	}

	a.DB.Debug().AutoMigrate(&models.User{}, &models.DonationProgram{}, &models.Donation{}, &models.Wallet{}, &models.TopUp{}, &models.Withdrawal{})

	a.Router = mux.NewRouter().StrictSlash(true)
	a.initializeRoutes()
}

func (a *App) initializeRoutes() {
	a.Router.Methods("OPTIONS")

	a.Router.HandleFunc("/", home).Methods("GET")
	u := a.Router.PathPrefix("/auth").Subrouter()
	u.HandleFunc("/register", a.UserSignUp).Methods("POST", "OPTIONS")
	u.HandleFunc("/login", a.Login).Methods("POST", "OPTIONS")

	g := a.Router.PathPrefix("/api").Subrouter()
	g.HandleFunc("/donate", a.GetDonationPrograms).Methods("GET")
	g.HandleFunc("/donate/search", a.SearchDonationProgram).Methods("GET")
	g.HandleFunc("/donate/{id:[0-9]+}", a.GetDonationProgramById).Methods("GET")

	s := a.Router.PathPrefix("/api").Subrouter()
	s.Use(middlewares.AuthJwtVerify)
	s.HandleFunc("/donate", a.CreateDonationProgram).Methods("POST", "OPTIONS")
	s.HandleFunc("/donate/{id:[0-9]+}", a.DonateToProgram).Methods("POST", "OPTIONS")
	s.HandleFunc("/donate/program", a.GetDonationProgramByFundraiser).Methods("GET")
	s.HandleFunc("/donate/verify/{id:[0-9]+}", a.VerifyDonationProgram).Methods("PUT")
	s.HandleFunc("/donate/unverified", a.GetUnverifiedDonationProgram).Methods("GET")

	s.HandleFunc("/user/verify/{id:[0-9]+}", a.VerifyFundraiser).Methods("PUT")
	s.HandleFunc("/donation/history", a.GetDonationHistoryFromUser).Methods("GET")
	s.HandleFunc("/user/unverified", a.GetUnverifiedUser).Methods("GET")

	s.HandleFunc("/withdraw/verify/{id:[0-9]+}", a.VerifyWithdrawal).Methods("PUT")
	s.HandleFunc("/withdraw/{id:[0-9]+}", a.CreateWithdrawal).Methods("POST")

	s.HandleFunc("/withdraw/unverified", a.GetUnverifiedWithdrawal).Methods("GET")
	s.HandleFunc("/user", a.GetUserById).Methods("GET")
	s.HandleFunc("/user", a.UpdateUser).Methods("PUT")
	s.HandleFunc("/user/change-password", a.ChangePassword).Methods("PUT")

	s.HandleFunc("/admin/unverified", a.UnverifiedList).Methods("GET")

	//wallet
	s.HandleFunc("/wallet", a.GetWalletByUserId).Methods("GET")
	s.HandleFunc("/wallet", a.CreateTopUp).Methods("POST")
	s.HandleFunc("/wallet/history", a.WalletHistory).Methods("GET")
}

func (a *App) RunServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}

	log.Printf("\nServer starting on port " + port)

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://pentapeduli.hexalogi.cyou"},
		AllowedMethods:   []string{"OPTIONS", "GET", "POST", "PUT"},
		AllowedHeaders:   []string{"Content-Type", "X-Requested-With", "Authorization"},
		AllowCredentials: true,
		Debug:            true,
	})

	handler := corsMiddleware.Handler(a.Router)

	err := http.ListenAndServe(":"+port, handler)

	if err != nil {
		fmt.Print(err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome to Fundraising")
}
