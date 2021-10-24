package routes

import (
	"net/http"

	controllers "example.com/gorestapi/controllers"
	"example.com/gorestapi/utils/auth"

	"github.com/gorilla/mux"
)

func Handlers() *mux.Router {

	r := mux.NewRouter().StrictSlash(true)
	r.Use(CommonMiddleware)
	r.HandleFunc("/vehicles", controllers.GetAllVehicles).Methods("GET")
	r.HandleFunc("/company", controllers.GetCompanyDetails).Methods("GET")
	r.HandleFunc("/feedback", controllers.GetAllFeedback).Methods("GET")
	r.HandleFunc("/feedbackCreate", controllers.PostFeedback).Methods("POST")
	r.HandleFunc("/register", controllers.CreateUser).Methods("POST")
	r.HandleFunc("/otp/authenticate", controllers.OtpAuthenticate).Methods("GET")
	r.HandleFunc("/login", controllers.Login).Methods("POST")

	s := r.PathPrefix("/auth").Subrouter()
	s.Use(auth.JwtVerify)
	s.HandleFunc("/user", controllers.FetchUsers).Methods("GET")

	return r

}

func CommonMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
		next.ServeHTTP(w, r)
	})

}
