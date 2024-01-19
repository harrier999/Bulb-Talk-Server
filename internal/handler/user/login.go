package user

import (
	"encoding/json"
	"net/http"
	"server/internal/db/postgres_db"
	"server/internal/models/orm"
	"server/pkg/authenticator"
)

type loginRequest struct {
	PhoneNumber string `json:"phone_number"`
	PassWord    string `json:"password"`
}

func Login(w http.ResponseWriter, r *http.Request) {
	var loginData loginRequest
	json.NewDecoder(r.Body).Decode(&loginData)
	if !validLoginData(loginData) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var user orm.User
	user, err := verifyUser(loginData)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	token, err := authenticator.CreateToken(user.UserID.String(), 24*14)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(token)
}

func validLoginData(loginData loginRequest) bool {
	if loginData.PhoneNumber == "" || loginData.PassWord == "" {
		return false
	}
	return true
}

func verifyUser(loginData loginRequest) (user orm.User, e error) {
	postgresClient := postgres_db.GetPostgresClient()
	result := postgresClient.Where("phone_number = ?", loginData.PhoneNumber).First(&user)
	if result.Error != nil {
		return user, result.Error
	}
	return user, nil
}
