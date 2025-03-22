package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-rest-api/dtos"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
)

type UserHandler struct {
	DB          *sql.DB
	redisClient *redis.Client
}

func NewUserHandler(db *sql.DB, reds *redis.Client) *UserHandler {
	return &UserHandler{db, reds}
}

func (h *UserHandler) GETusers(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	var users []dtos.UserResponse
	cntxt := context.Background()
	val, err := h.redisClient.Get(cntxt, "users").Result()
	if err != nil {
		log.Println(err.Error())
		fmt.Println("error retrieving users from redis")
	}
	if err == redis.Nil {
		log.Println("insider err==redis.Nil")
		rows, err := h.DB.Query(" SELECT userId,firstName,lastName,email from users")
		if err != nil {
			log.Println(err.Error())
			fmt.Println("errors selecting users from database ")
			return
		}
		for rows.Next() {
			var user dtos.UserResponse
			rows.Scan(&user.UserId, &user.FirstName, &user.LastName, &user.Email)
			users = append(users, user)

		}
		response, err := json.MarshalIndent(users, "", " ")
		_, err = h.redisClient.Set(cntxt, "users", response, 10*time.Minute).Result()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(response)
		return

	}
	response := []byte(val)
	w.Write(response)
}
func (h *UserHandler) POSTUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	var user dtos.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	r.Body.Close()
	_, err := h.DB.Exec(
		"INSERT INTO users values(?,?,?,?,?)",
		&user.UserId,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cntxt := context.Background()
	err = h.redisClient.Del(cntxt, "users").Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = h.redisClient.HSet(cntxt, user.UserId, map[string]interface{}{"userId": user.UserId, "email": user.Email, "firstName": user.FirstName, "lastName": user.LastName}).
		Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	response, err := json.Marshal(&user)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (h *UserHandler) GETUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	userId := vars["userId"]
	cntxt := context.Background()

	val, err := h.redisClient.HGetAll(cntxt, userId).Result()
	var user dtos.UserResponse
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(val) == 0 {
		result := h.DB.QueryRow(
			"SELECT userId firstName, lastName, email FROM users WHERE userId=?",
			userId,
		)
		err := result.Scan(&user.UserId, &user.Email, &user.FirstName, &user.LastName)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			http.Error(w, "user Not found", http.StatusNotFound)
			return
		}
		//cache the result in redis
		_, err = h.redisClient.HSet(cntxt, userId, map[string]interface{}{"userId": user.UserId, "email": user.Email, "firstName": user.FirstName, "lastName": user.LastName}).
			Result()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	} else {
		user = dtos.UserResponse{
			UserId:    val["userId"],
			FirstName: val["firstName"],
			LastName:  val["lastName"],
			Email:     val["email"],
		}

	}
	response, err := json.Marshal(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(response)
}

func (h *UserHandler) DELETEUser(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	vars := mux.Vars(r)
	userId := vars["userId"]

	result, err := h.DB.Exec("delete from users where userId=?", userId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := result.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	cntxt := context.Background()
	err = h.redisClient.Del(cntxt, "users").Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(201)
	response, err := json.Marshal(rows)
	w.Write(response)

}
