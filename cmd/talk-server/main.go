package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
	"server/api/chatting"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `gorm:"type:varchar(100);unique_index" json:"email"`
}



func main() {


	r := mux.NewRouter()
	r.HandleFunc("/chat", chatting.Handler).Methods("GET")
	log.Fatal(http.ListenAndServe(":18000", r))
}

type UserError struct {
	Message string `json:"message"`
}

func signup(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		json.NewDecoder(r.Body).Decode(&user)
		if user.Email == "" || user.Password == "" || user.Username == "" {
			json.NewEncoder(w).Encode(UserError{Message: "Fields are empty"})
			return
		}
		hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
		if err != nil {
			log.Fatal(err)
		}
		user.Password = string(hash)
		result := db.Create(&user)
		if result.Error != nil{
			json.NewEncoder(w).Encode(UserError{Message: "User already exists"})
			return
		}
		user.Password = ""
		json.NewEncoder(w).Encode(user)
	}
}

func login(db *gorm.DB) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var user User
		var userFound User
		json.NewDecoder(r.Body).Decode(&user)
		if user.Email == "" || user.Password == "" {
			json.NewEncoder(w).Encode(UserError{Message: "Fields are empty"})
			return
		}
		db.Where("email = ?", user.Email).First(&userFound)
		if userFound.Email == "" {
			json.NewEncoder(w).Encode(UserError{Message: "User not found"})
			return
		}
		err := bcrypt.CompareHashAndPassword([]byte(userFound.Password), []byte(user.Password))
		if err != nil {
			json.NewEncoder(w).Encode(UserError{Message: "Invalid password"})
			return
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"email": userFound.Email,
			"exp":   time.Now().Add(time.Hour * 1).Unix(),
		})
		tokenString, error := token.SignedString([]byte("secret"))
		if error != nil {
			fmt.Println(error)
		}
		json.NewEncoder(w).Encode(UserError{Message: tokenString})
	}
}
