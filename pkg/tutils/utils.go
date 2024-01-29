package tutils

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"server/pkg/authenticator"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func CreateRouterWithMiddleware(handler func(w http.ResponseWriter, r *http.Request)) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/test", handler).Methods("POST")
	router.Use(authenticator.JWTMiddleware)
	return router
}

func CreateRequestAndResponse(requestBody []byte) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest("POST", "/test", bytes.NewBuffer(requestBody))
	res := httptest.NewRecorder()
	return req, res
}

func CreateToken(user_id uuid.UUID) string {
	token, _ := authenticator.CreateToken(user_id.String(), 24)
	return token
}
func CreateTokenByString(user_id string) string {
	token, _ := authenticator.CreateToken(user_id, 24)
	return token
}
