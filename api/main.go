package main

import (
	"context"
	"fmt"
	"log"
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
	con := &config.EnvConfig{}
	env := con.LoadEnv()
	log.Print(env)
	db := con.ConnectDB().Db

	defer db.Close()
    fmt.Println()
	redisClient := redis.NewClient(&redis.Options{
		Addr:     env.REDIS_URL,
		Password: env.REDIS_PASSWORD,
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

	userHandler := handlers.NewUserHandler(db, redisClient)
	r.HandleFunc("/users", userHandler.GETusers).Methods("GET")
	r.HandleFunc("/users", userHandler.POSTUser).Methods("POST")
	r.HandleFunc("/users/{userId}", userHandler.GETUser).Methods("GET")
	r.HandleFunc("/users/{userId}", userHandler.DELETEUser).Methods("DELETE")
	fmt.Print(http.ListenAndServe(":8080", r))
	return r
}
