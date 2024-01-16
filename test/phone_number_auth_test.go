package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"server/internal/db/postgres_db"
	"server/internal/handler/user"
	"server/internal/models/orm"

	"github.com/stretchr/testify/assert"
)

func TestReqeustAuthNumber(t *testing.T) {

	var requestData = user.AuthRequest{}
	db := postgres_db.GetPostgresClient()
	db.AutoMigrate(&orm.AuthenticateMessage{})
	requestData.PhoneNumber = os.Getenv("SMS_PHONE_NUMBER")
	requestData.CountryCode = "82"
	requestData.DeviceID = "test"
	reqBody, err := json.Marshal(requestData)
	if err != nil {
		t.Errorf("failed to marshal request data")
	}

	req := httptest.NewRequest("POST", "/auth", bytes.NewBuffer(reqBody))
	rr := httptest.NewRecorder()

	req2 := httptest.NewRequest("POST", "/auth", bytes.NewBuffer(reqBody))
	rr2 := httptest.NewRecorder()

	handler := http.HandlerFunc(user.ReqeustAuthNumber)
	handler.ServeHTTP(rr, req)
	defer db.Delete(&orm.AuthenticateMessage{}, "device_id = ?", requestData.DeviceID)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusOK)
	}
	handler.ServeHTTP(rr2, req2)
	if status := rr2.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v, want %v", status, http.StatusConflict)
	}
}

func TestCheckAuthNumber(t *testing.T) {
	var requestData user.AuthRequest
	requestData.PhoneNumber = os.Getenv("SMS_PHONE_NUMBER")
	requestData.CountryCode = "82"
	requestData.DeviceID = "test2"

	reqBody, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("failed to marshal request data")
	}

	req := httptest.NewRequest("POST", "/auth", bytes.NewBuffer(reqBody))
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(user.ReqeustAuthNumber)
	handler.ServeHTTP(res, req)
	if status := res.Code; status != http.StatusOK {
		t.Errorf("Failed to get auth message. Status code: %d", status)
		return
	}
	// var input string
	// fmt.Print("인증 번호를 입력해 주세요: ")
	// fmt.Scanln(&input)
	// fmt.Println(input)

	requestData.AuthenticateNumber = "190"
	reqBody, err = json.Marshal(requestData)
	if err != nil {
		t.Error("Failed to marshal JSON")
	}

	req2 := httptest.NewRequest("POST", "/auth", bytes.NewBuffer(reqBody))
	res2 := httptest.NewRecorder()
	checkAuthNumberHandler := http.HandlerFunc(user.CheckAuthNumber)
	checkAuthNumberHandler(res2, req2)
	if res2.Code == http.StatusConflict {
		log.Println("Wrong Authenticate Number")
		return
	}
	log.Println(res2.Code)
	assert.NotEqual(t, http.StatusConflict, res2.Code, "something wrong")
	assert.NotEqual(t, http.StatusOK, res2.Code, "Authenticate must be failed. But returned 200")

}

// func main() {
// 	var requestData user.AuthRequest
// 	requestData.PhoneNumber = os.Getenv("SMS_PHONE_NUMBER")
// 	requestData.CountryCode = "82"
// 	requestData.DeviceID = "test4"

// 	reqBody, err := json.Marshal(requestData)
// 	if err != nil {
// 		fmt.Println("failed to marshal request data")
// 	}

// 	req := httptest.NewRequest("POST", "/auth", bytes.NewBuffer(reqBody))
// 	res := httptest.NewRecorder()
// 	handler := http.HandlerFunc(user.ReqeustAuthNumber)
// 	handler.ServeHTTP(res, req)
// 	if status := res.Code; status != http.StatusOK {
// 		fmt.Println(status)
// 		fmt.Println("Failed to get auth message")
// 		log.Fatal()
// 	}
// 	var input string
// 	fmt.Print("인증 번호를 입력해 주세요: ")
// 	fmt.Scanln(&input)
// 	fmt.Println(input)

// 	requestData.AuthenticateNumber = input
// 	fmt.Println("input:", requestData.AuthenticateNumber)
// 	reqBody, err = json.Marshal(requestData)

// 	fmt.Println(reqBody)
// 	fmt.Println(requestData)

// 	req2 := httptest.NewRequest("POST", "/auth", bytes.NewBuffer(reqBody))
// 	res2 := httptest.NewRecorder()
// 	checkAuthNumberHandler := http.HandlerFunc(user.CheckAuthNumber)
// 	checkAuthNumberHandler(res2, req2)
// 	if res2.Code != http.StatusOK {
// 		fmt.Println("Wrong Authenticate Number")
// 	}

// }
