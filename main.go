package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sangeeth200494/JWT-AUTH_Golang/login"
	userhandlers "github.com/sangeeth200494/JWT-AUTH_Golang/user-handlers"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/Home", login.Home).Methods("GET")
	router.HandleFunc("/login", login.LoginHandler).Methods("POST")
	router.HandleFunc("/RegisterUser", userhandlers.RegisterUser).Methods("POST")

	// privateRouter := router.PathPrefix("/").Subrouter()
	// privateRouter.Use(middleware.AuthMiddleware)
	router.HandleFunc("/protected", login.ProtectedHandler).Methods("GET")

	log.Println(":Listening on :4000")
	err := http.ListenAndServe("localhost:4000", router)
	if err != nil {
		fmt.Println("Could not start the server", err)
	}
}
