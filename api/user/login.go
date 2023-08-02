package user

import (
	"encoding/json"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

type User struct {
	UserId       uint   `gorm:"primary_key;auto_increment" json:"user_id"`
	OAuthId      string `gorm:"type:varchar(100);unique_index" json:"oauth_id"`
	Username     string `json:"username"`
	Password     string `json:"password"`
	Email        string `gorm:"type:varchar(40);unique_index" json:"email"`
	PhoneNumber  string `gorm:"type:varchar(20);unique_index" json:"phone_number"`
	ProfileImage string `json:"profile_image"`
	Enabled      bool   `gorm:"default:true" json:"enabled"`
	Role         string `gorm:"default:ROLE_USER" json:"role"`

	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
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
		if result.Error != nil {
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
			log.Println(error)
		}
		json.NewEncoder(w).Encode(UserError{Message: tokenString})
	}
}
