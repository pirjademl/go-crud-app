package main

import (
	"fmt"
	"net/http"

	"github.com/go-rest-api/config"
	"github.com/go-rest-api/handlers"
	_ "github.com/go-rest-api/handlers"
	"github.com/gorilla/mux"
)

func main() {
	r := registerRoutes()
	http.ListenAndServe(":8080", r)
}

func registerRoutes() *mux.Router {
	r := mux.NewRouter()
	db := config.ConnectDB()
	defer db.Db.Close()
	userHandler := handlers.NewUserHandler(db.Db)
	r.HandleFunc("/users", userHandler.GETusers).Methods("GET")
	r.HandleFunc("/users", userHandler.POSTUser).Methods("POST")
	r.HandleFunc("/users/{userId}", userHandler.GETUser).Methods("GET")
	r.HandleFunc("/users/{userId}", userHandler.DELETEUser).Methods("DELETE")
	fmt.Print(http.ListenAndServe(":8080", r))
	return r
}
