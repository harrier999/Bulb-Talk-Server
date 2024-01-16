package user

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	"log"
	"math/big"
	"net/http"
	"server/internal/db/postgres_db"
	"server/internal/models/orm"
	"server/internal/utils"
	"time"

	"github.com/ttacon/libphonenumber"
)

type AuthRequest struct {
	PhoneNumber        string `json:"phone_number"`
	CountryCode        string `json:"country_code"`
	DeviceID           string `json:"device_id"`
	AuthenticateNumber string `json:"authenticate_number"`
}

func RequestAuthNumber(w http.ResponseWriter, r *http.Request) {
	var authData AuthRequest
	var err error
	json.NewDecoder(r.Body).Decode(&authData)
	if err := validAuthData(authData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err)
		return
	}
	authData.AuthenticateNumber, err = createRandomNumber()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	if err := storeAuthNumber(authData); err != nil {
		w.WriteHeader(http.StatusConflict)
		log.Println(err)
		return
	}
	err = utils.SendSMS(authData.PhoneNumber, "Bulb Talk \n인증번호는 ["+authData.AuthenticateNumber+"] 입니다.")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func CheckAuthNumber(w http.ResponseWriter, r *http.Request) {
	var authData AuthRequest
	json.NewDecoder(r.Body).Decode(&authData)

	if err := validAuthData(authData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	postgresClient := postgres_db.GetPostgresClient()
	var authenticateNumber orm.AuthenticateMessage
	result := postgresClient.Where("device_id = ? AND expire_time > ? AND phone_number = ?", authData.DeviceID, time.Now(), authData.PhoneNumber).First(&authenticateNumber)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		log.Println("Error checking auth number: ", result.Error.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	if authenticateNumber.AuthenticateNumber != authData.AuthenticateNumber {
		w.WriteHeader(http.StatusForbidden)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func validAuthData(authData AuthRequest) error {
	if authData.PhoneNumber == "" || authData.CountryCode == "" || authData.DeviceID == "" {
		return errors.New("invalid auth data")
	}
	if _, err := libphonenumber.Parse(authData.PhoneNumber, "KR"); err != nil {
		return errors.New("invalid phone number")
	}
	if authData.CountryCode != "82" {
		return errors.New("korean number only")
	}
	return nil
}

func storeAuthNumber(authData AuthRequest) error {
	postgresClient := postgres_db.GetPostgresClient()
	var authenticateNumber orm.AuthenticateMessage

	authenticateNumber.CountryCode = authData.CountryCode
	authenticateNumber.PhoneNumber = authData.PhoneNumber
	authenticateNumber.DeviceID = authData.DeviceID
	authenticateNumber.AuthenticateNumber = authData.AuthenticateNumber
	authenticateNumber.RequestTime = time.Now()
	authenticateNumber.ExpireTime = time.Now().Add(time.Minute * 3)

	err := checkIfAlreadyRequested(authData)
	if err != nil {
		return err
	}
	result := postgresClient.Create(&authenticateNumber)
	if result.Error != nil {
		log.Println("Error storing auth number: ", result.Error.Error())
		return errors.New("failed to store auth number")
	}

	return nil
}

func checkIfAlreadyRequested(authData AuthRequest) error {
	postgresClient := postgres_db.GetPostgresClient()
	var authenticateNumber orm.AuthenticateMessage
	result := postgresClient.Where("device_id = ? AND expire_time > ?", authData.DeviceID, time.Now()).First(&authenticateNumber)
	if result.Error != nil {
		if result.Error.Error() == "record not found" {
			return nil
		}
		log.Println("Error checking auth number: ", result.Error.Error())
		return result.Error
	}
	return errors.New("already requested auth number. please wait 3 minutes")
}

func createRandomNumber() (string, error) {
	var randomNumber *big.Int
	var err error
	for i := 0; i < 4; i++ {
		randomNumber, err = rand.Int(rand.Reader, big.NewInt(900000))
		if err != nil {
			log.Println("Error creating random number: ", err.Error())
		}
		if err == nil {
			s := randomNumber.Add(randomNumber, big.NewInt(100000)).String()
			return s, nil
		}
	}
	return "", errors.New("failed to create random number")
}
