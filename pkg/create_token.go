package authenticator

import "github.com/golang-jwt/jwt"

func CreateToken(user_id string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	token.Header["alg"] = "HS256"
	claims[string(contextKeyUserID)] = user_id
	tokenString, err := token.SignedString(hmacSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
