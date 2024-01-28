package test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"server/internal/db/postgres_db"
	"server/internal/handler/friends"
	"server/pkg/authenticator"
	"testing"

	"server/internal/models/orm"

	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

var (
	_USER_1 = orm.User{
		UserName:    "test",
		PhoneNumber: "01012345678",
		CountryCode: "82",
	}
	_USER_2 = orm.User{
		UserName:    "test2",
		PhoneNumber: "01022222222",
		CountryCode: "82",
	}
)

func TestMain(m *testing.M) {
	client := postgres_db.GetTestPostgresCleint()
	client.Migrator().DropTable(&orm.User{}, &orm.Friend{}, &orm.Room{}, &orm.RoomUser{}, orm.AuthenticateMessage{})
	client.AutoMigrate(&orm.User{}, &orm.Friend{}, &orm.Room{}, &orm.RoomUser{}, orm.AuthenticateMessage{})
	
	client.Create(&_USER_1)
	client.Create(&_USER_2)
	
	m.Run()

}

func createAddFriendRouter() *mux.Router {
	r := mux.NewRouter()
	r.Use(authenticator.JWTMiddleware)
	r.HandleFunc("/friends", friends.AddFriend).Methods("POST")
	return r
}

func TestAddFriendNormalCase(t *testing.T) {
	r := createAddFriendRouter()

	requestData := friends.AddFriendRequest{
		PhoneNumber: "01012345678",
	}
	reqBody, _ := json.Marshal(requestData)
	req := httptest.NewRequest("POST", "/friends", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()
	token, _ := authenticator.CreateToken(_VALID_USER_ID_1, 24)
	req.Header.Set("Authorization", token)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
}

func TestAddFriendInvalidPhoneNumber(t *testing.T) {
	r := createAddFriendRouter()

	requestData := friends.AddFriendRequest{
		PhoneNumber: "010123456789",
	}
	reqBody, _ := json.Marshal(requestData)
	req := httptest.NewRequest("POST", "/friends", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()
	token, _ := authenticator.CreateToken(_VALID_USER_ID_1, 24)
	req.Header.Set("Authorization", token)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestAddFriendInvalidUUID(t *testing.T) {
	r := createAddFriendRouter()

	requestData := friends.AddFriendRequest{
		PhoneNumber: "01012345678",
	}
	reqBody, _ := json.Marshal(requestData)
	req := httptest.NewRequest("POST", "/friends", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()
	token, _ := authenticator.CreateToken(_WRONG_USER_ID_1, 24)
	req.Header.Set("Authorization", token)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)
}

func TestAddFriendEmptyPhoneNumber(t *testing.T) {
	r := createAddFriendRouter()

	requestData := friends.AddFriendRequest{
		PhoneNumber: "",
	}
	reqBody, _ := json.Marshal(requestData)
	req := httptest.NewRequest("POST", "/friends", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()
	token, _ := authenticator.CreateToken(_VALID_USER_ID_1, 24)
	req.Header.Set("Authorization", token)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)
}
// TODO: Fix this test case
// func TestAddFriendAlreadyFriend(t *testing.T) {
// 	r := createAddFriendRouter()

// 	requestData := friends.AddFriendRequest{
// 		PhoneNumber: "01012345678",
// 	}
// 	reqBody, _ := json.Marshal(requestData)
// 	req := httptest.NewRequest("POST", "/friends", bytes.NewBuffer(reqBody))
// 	rr := httptest.NewRecorder()
// 	token, _ := authenticator.CreateToken(_VALID_USER_ID_1, 24)
// 	req.Header.Set("Authorization", token)
// 	r.ServeHTTP(rr, req)
// 	assert.Equal(t, http.StatusBadRequest, rr.Code)
// }
