package user

import (
	"encoding/json"
	"errors"
	"net/http"
	"server/internal/utils"
)

type authRequest struct {
	PhoneNumber        string `json:"phone_number"`
	CountryCode        string `json:"country_code"`
	DeviceID           string `json:"device_id"`
	AuthenticateNumber string `json:"authenticate_number"`
}

func ReqeustAuthNumber(w http.ResponseWriter, r *http.Request) {
	var authData authRequest
	json.NewDecoder(r.Body).Decode(&authData)
	if err := validAuthData(authData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	err := utils.SendSMS(authData.PhoneNumber, "인증번호는 1234입니다.")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func CheckAuthNumber(w http.ResponseWriter, r *http.Request) {
	var authData authRequest
	json.NewDecoder(r.Body).Decode(&authData)

	if err := validAuthData(authData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}

func validAuthData(authData authRequest) error {
	if authData.PhoneNumber == "" || authData.CountryCode == "" || authData.DeviceID == "" {
		return errors.New("invalid auth data")
	}
	return nil
}


