package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type User struct {
	ID        string `json:"id"`
	Firstname string `json:"firstname"`
	Email     string `json:"email"`
}

var users []User

func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Fatal(err)
		return
	}

	user.ID = strconv.Itoa(len(users) + 1)
	users = append(users, user)

	json.NewEncoder(w).Encode(user)
	w.WriteHeader(http.StatusOK)
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(users)
	w.WriteHeader(http.StatusOK)
}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	for index, user := range users {
		if user.ID == id {
			var updatedUser User
			err := json.NewDecoder(r.Body).Decode(&updatedUser)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				log.Fatal(err)
				return
			}

			user.Firstname = updatedUser.Firstname
			user.Email = updatedUser.Email
			users[index] = user
			json.NewEncoder(w).Encode(user)
			w.WriteHeader(http.StatusOK)
		}
	}
}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	for index, user := range users {
		if user.ID == id {
			users = append(users[:index], users[index+1:]...)
			json.NewEncoder(w).Encode(users)
			w.WriteHeader(http.StatusOK)
		}
	}
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/create", CreateUser).Methods("POST")
	r.HandleFunc("/users", GetUsers).Methods("GET")
	r.HandleFunc("/update/{id:[0-9]}", UpdateUser).Methods("PUT")
	r.HandleFunc("/delete/{id:[0-9]}", DeleteUser).Methods("DELETE")

	log.Fatal(http.ListenAndServe(":9090", r))
}
