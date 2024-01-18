package authenticator

import (
	"fmt"
	"net/http"
	"os"
	"context"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string
const contextKeyUserID contextKey = "user_id"

var hmacSecret = []byte(os.Getenv("JWT_SECRET"))

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return hmacSecret, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		user_id, ok := token.Claims.(jwt.MapClaims)["user_id"]
		if !ok {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), contextKeyUserID, user_id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
