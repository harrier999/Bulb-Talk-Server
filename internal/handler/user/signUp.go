package user

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"server/internal/db/postgres_db"
	"server/internal/models/orm"

	"golang.org/x/crypto/bcrypt"
)

type UserData struct {
	Email              string `json:"email"`
	Password           string `json:"password"`
	Username           string `json:"username"`
	PhoneNumber        string `json:"phone_number"`
	CountryCode        string `json:"country_code"`
	AuthenticateNumber string `json:"authenticate_number"`
	DeviceID           string `json:"device_id"`
}

// Signup with Phone Number and Password
func SignUp(w http.ResponseWriter, r *http.Request) {
	var userData UserData
	json.NewDecoder(r.Body).Decode(&userData)

	if err := validUserData(userData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := checkDuplicatedUser(userData); err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	_, err := createUser(userData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func checkDuplicatedUser(userData UserData) error {
	postgresClient := postgres_db.GetPostgresClient()
	var user UserData
	result := postgresClient.Where("email = ?", userData.Email).First(&user)
	if result.Error != nil {
		log.Printf("Error: %s", result.Error)
		return errors.New("User already exists")
	}
	return nil
}

func validUserData(userData UserData) error {
	if userData.Email == "" || userData.Password == "" || userData.Username == "" || userData.PhoneNumber == "" || userData.CountryCode == "" {
		return errors.New("fields are empty")
	}
	return nil
}

func createUserObject(userData UserData) orm.User {
	var user orm.User
	user.UserName = userData.Username
	user.Email.String = userData.Email
	user.PhoneNumber = userData.PhoneNumber
	user.CountryCode = userData.CountryCode
	return user
}

func createUser(userData UserData) (orm.User, error) {
	postgresClient := postgres_db.GetPostgresClient()
	hash, err := encryptPassword(userData.Password)
	if err != nil {
		log.Printf("Error: %s", err)
		return orm.User{}, errors.New("failed to encrypt password")
	}
	userData.Password = hash
	user := createUserObject(userData)
	result := postgresClient.Create(&userData)
	if result.Error != nil {
		log.Printf("Error: %s", result.Error)
		return user, errors.New("failed to create user")
	}
	return user, nil
}

func encryptPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}
