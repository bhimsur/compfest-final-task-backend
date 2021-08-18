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
	a.Router.Use(middlewares.SetContentTypeMiddleware)
	a.Router.HandleFunc("/", home).Methods("GET")
	u := a.Router.PathPrefix("/auth").Subrouter()
	u.HandleFunc("/register", a.UserSignUp).Methods("POST")
	u.HandleFunc("/login", a.Login).Methods("POST")

	s := a.Router.PathPrefix("/api").Subrouter()
	s.Use(middlewares.AuthJwtVerify)
	s.HandleFunc("/donate", a.GetDonationPrograms).Methods("GET")
	s.HandleFunc("/donate", a.CreateDonationProgram).Methods("POST")
	s.HandleFunc("/donate/{id:[0-9]+}", a.DonateNow).Methods("POST")
	s.HandleFunc("/donate/{id:[0-9]+}", a.GetDonationProgramById).Methods("GET")
	s.HandleFunc("/donate/history", a.GetDonationProgramByFundraiser).Methods("GET")
	s.HandleFunc("/donate/verify/{id:[0-9]+}", a.VerifyDonationProgram).Methods("PUT")

	s.HandleFunc("/user/verify/{id:[0-9]+}", a.VerifyFundraiser).Methods("PUT")
	s.HandleFunc("/user/donate/history", a.GetDonationHistoryFromUser).Methods("GET")
	s.HandleFunc("/user", a.GetUserById).Methods("GET")
	s.HandleFunc("/user", a.UpdateUser).Methods("PUT")

	//wallet
	s.HandleFunc("/wallet", a.GetWalletByUserId).Methods("GET")
	s.HandleFunc("/wallet", a.CreateTopUp).Methods("POST")
	s.HandleFunc("/wallet/history", a.GetTopUpHistoryByUserId).Methods("GET")
}

func corsHandler(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if (r.Method == "POST" || r.Method == "OPTIONS" ||  r.Method == "PUT" ||  r.Method == "DELETE") {
		log.Print("preflight detected: ", r.Header)
		w.Header().Add("Connection", "keep-alive")
		w.Header().Add("Access-Control-Allow-Origin", "https://pentapeduli.hexalogi.cyou")
		w.Header().Add("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Add("Access-Control-Allow-Headers", "Authorization, Content-Type")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
		} else {
			h.ServeHTTP(w, r)
		}
	}
}

func (a *App) RunServer() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	// c := cors.New(cors.Options{
	// 	AllowedOrigins: []string{"*"},
	// 	AllowedMethods: []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"},
	// })
	log.Printf("\nServer starting on port " + port)
	err := http.ListenAndServe(":"+port, corsHandler(a.Router))
	if err != nil {
		fmt.Print(err)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	responses.JSON(w, http.StatusOK, "Welcome to Fundraising")
}
