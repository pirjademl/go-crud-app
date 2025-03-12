package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-rest-api/config"
	"github.com/go-rest-api/handlers"
	_ "github.com/go-rest-api/handlers"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

func main() {
	r := registerRoutes()
	http.ListenAndServe(":8080", r)
}

func registerRoutes() *mux.Router {
	r := mux.NewRouter()
	db := config.ConnectDB()
	defer db.Db.Close()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "1234",
		DB:       0,
	})
	cntxt := context.Background()
	err := redisClient.Set(cntxt, "foo", "bar", 0).Err()
	if err != nil {
		panic(err)
	}
	val, err := redisClient.Get(cntxt, "foo").Result()
	if err != nil {
		panic(err)
	}

	fmt.Println("foo", val)

	userHandler := handlers.NewUserHandler(db.Db, redisClient)
	r.HandleFunc("/users", userHandler.GETusers).Methods("GET")
	r.HandleFunc("/users", userHandler.POSTUser).Methods("POST")
	r.HandleFunc("/users/{userId}", userHandler.GETUser).Methods("GET")
	r.HandleFunc("/users/{userId}", userHandler.DELETEUser).Methods("DELETE")
	fmt.Print(http.ListenAndServe(":8080", r))
	return r
}
