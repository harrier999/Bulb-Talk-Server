package test

import (
	"bytes"
	"encoding/json"

	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"server/internal/handler/friends"

	"server/pkg/authenticator"

	"github.com/stretchr/testify/assert"

	"github.com/gorilla/mux"
)

var (
	_WRONG_USER_ID_1 = "wrong_user_id"
	_WRONG_USER_ID_2 = "b5b06de3-5a19-4420-a6ae-c84b4439a61"

	_VALID_USER_ID_1 = "b5b06de3-5a19-4420-a6ae-cc84b4439a61"
)

func createRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(authenticator.JWTMiddleware)
	r.HandleFunc("/friends", friends.GetFriendList).Methods("POST")
	return r
}

func TestNormalCase(t *testing.T) {
	r := createRouter()

	requestData := friends.FriendListRequest{
		LastRequestTime: time.Now(),
	}
	reqBody, _ := json.Marshal(requestData)

	req := httptest.NewRequest("POST", "/friends", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()
	token, _ := authenticator.CreateToken(_VALID_USER_ID_1, 24)
	req.Header.Set("Authorization", token)

	r.ServeHTTP(rr, req)
	if rr.Result().StatusCode != http.StatusOK {
		t.Errorf("failed to get friend list")
	}

}

func TestInvalidUUID(t *testing.T) {
	r := createRouter()

	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/friends", nil)
	token, _ := authenticator.CreateToken(_WRONG_USER_ID_1, 24)
	req.Header.Set("Authorization", token)

	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode)

	res2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("POST", "/friends", nil)
	token2, _ := authenticator.CreateToken(_WRONG_USER_ID_2, 24)

	req2.Header.Set("Authorization", token2)
	r.ServeHTTP(res2, req2)

	assert.Equal(t, http.StatusUnauthorized, res2.Result().StatusCode)
}

func TestEmptyFriend(t *testing.T) {
	r := createRouter()

	res := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/friends", nil)
	token, _ := authenticator.CreateToken(_VALID_USER_ID_1, 24)
	req.Header.Set("Authorization", token)

	r.ServeHTTP(res, req)
	t.Logf("res: %v", res.Body)
	assert.Equal(t, http.StatusOK, res.Result().StatusCode)
}
