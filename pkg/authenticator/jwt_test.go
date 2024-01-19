package authenticator

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gorilla/mux"
)

func dummyHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func TestJWTMiddleware(t *testing.T) {
	m := mux.NewRouter()
	m.Use(JWTMiddleware)
	m.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		user_id := r.Context().Value(contextKeyUserID)
		assert.NotNil(t, user_id)
		assert.Equal(t, "test_user_id", user_id)

	})
	req, err := http.NewRequest("GET", "/test", nil)
	assert.Nil(t, err)
	token, _ := CreateToken("test_user_id", 100)
	req.Header.Set("Authorization", token)
	rr := httptest.NewRecorder()
	m.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestExpirationTime(t *testing.T) {
	token, _ := CreateToken("test_user_id", -1)
	m := mux.NewRouter()
	m.Use(JWTMiddleware)
	m.HandleFunc("/test", dummyHandler)
	req, err := http.NewRequest("GET", "/test", nil)
	assert.Nil(t, err)
	req.Header.Set("Authorization", token)
	rr := httptest.NewRecorder()
	m.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	token2, _ := CreateToken("test_user_id", 1)
	req2, err2 := http.NewRequest("GET", "/test", nil)
	assert.Nil(t, err2)
	req.Header.Set("Authorization", token2)
	rr2 := httptest.NewRecorder()
	m.ServeHTTP(rr2, req2)
	assert.Equal(t, http.StatusUnauthorized, rr2.Code)
}
