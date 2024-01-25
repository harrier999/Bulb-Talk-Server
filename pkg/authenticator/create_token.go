package authenticator

import (
	"time"

	"github.com/golang-jwt/jwt"
)

func CreateToken(user_id string, durationHours int64) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	token.Header["alg"] = "HS256"
	claims[string(ContextKeyUserID)] = user_id
	claims["exp"] = time.Now().Add(time.Hour * time.Duration(durationHours)).Unix()
	tokenString, err := token.SignedString(hmacSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
