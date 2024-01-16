package user

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	// "net/mail"
	"server/internal/db/postgres_db"
	"server/internal/models/orm"

	"github.com/ttacon/libphonenumber"
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
		w.Write([]byte(err.Error()))
		return
	}
	if err := checkIfPhoneNumberExists(userData); err != nil {
		w.WriteHeader(http.StatusConflict)
		log.Println(err.Error())
		return
	}
	_, err := createUser(userData)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func checkIfPhoneNumberExists(userData UserData) error {
	postgresClient := postgres_db.GetPostgresClient()
	var user orm.User
	result := postgresClient.Where("phone_number = ?", userData.PhoneNumber).First(&user)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil
		}
		log.Printf("Error: %s", result.Error)
		return errors.New("failed to check if user exists")
	}
	return errors.New("User already exists")
}

func validUserData(userData UserData) error {
	if userData.Password == "" || userData.Username == "" || userData.PhoneNumber == "" || userData.CountryCode == "" {
		return errors.New("fields are empty")
	}
	// _, err := mail.ParseAddress(userData.Email)
	// if err != nil {
	// 	return errors.New("invalid email address")
	// }
	if len(userData.Password) < 6 {
		return errors.New("password is too short")
	}
	if len(userData.Password) > 24 {
		return errors.New("password is too long")
	}
	_, err := libphonenumber.Parse(userData.PhoneNumber, "KR")
	if err != nil {
		return errors.New("invalid phone number")
	}
	if userData.CountryCode != "82" {
		return errors.New("korean phone number only")
	}
	return nil
}

func createUserObject(userData UserData) orm.User {
	var user orm.User
	user.UserName = userData.Username
	user.Email.String = ""
	user.PhoneNumber = userData.PhoneNumber
	user.CountryCode = userData.CountryCode
	user.PasswordHash = userData.Password
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
	result := postgresClient.Create(&user)
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
