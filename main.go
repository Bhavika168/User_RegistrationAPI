package main

import (
	"Project18/home"
	"Project18/user"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	user.InitialiseDb()
	user.InitialiseRedis()

	router := mux.NewRouter()

	signupRouter := router.PathPrefix("/signup").Subrouter()
	signupRouter.HandleFunc("", user.Signup).Methods(http.MethodGet)
	signupRouter.HandleFunc("/otp", user.SignupOTP).Methods(http.MethodGet)

	loginRouter := router.PathPrefix("/login").Subrouter()
	loginRouter.HandleFunc("", user.Login).Methods(http.MethodGet)
	loginRouter.HandleFunc("/otp", user.LoginOTP).Methods(http.MethodGet)

	homeRouter := router.PathPrefix("/home").Subrouter()
	homeRouter.Use(home.Auth)
	homeRouter.HandleFunc("", home.Home).Methods(http.MethodGet)

	log.Fatal(http.ListenAndServe(":8080", router))
}
