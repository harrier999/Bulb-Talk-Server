package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

func SendSMS(phone_number string, message string) error {
	log.Println("Sending SMS to " + phone_number)
	token := getTokenManager().getToken()
	err := sendSMS(token, phone_number, message)
	if err != nil {
		return err
	}
	return nil
}

var (
	SMS_API_KEY      string
	SMS_ID           string
	SMS_PHONE_NUMBER string
	once             sync.Once
	tm               *tokenManager
)

func init() {
	SMS_API_KEY = os.Getenv("SMS_API_KEY")
	SMS_ID = os.Getenv("SMS_ID")
	SMS_PHONE_NUMBER = os.Getenv("SMS_PHONE_NUMBER")
}

type tokenResponse struct {
	AccessToken string `json:"access_token"`
}

type tokenManager struct {
	tokenChannel chan string
}

func getTokenManager() *tokenManager {
	once.Do(func() {
		tm = &tokenManager{
			tokenChannel: make(chan string, 1),
		}
		initialToken, err := requestToken()
		if err != nil {
			log.Fatal("Failed to get initial token")
		}
		log.Println("Successfully got initial token")
		tm.tokenChannel <- initialToken
		tm.startRefreshRoutine(50 * time.Minute)
	})

	return tm
}

func (tm *tokenManager) getToken() string {
	token := <-tm.tokenChannel
	tm.tokenChannel <- token
	return token
}

func (tm *tokenManager) startRefreshRoutine(refreshInterval time.Duration) {
	go func() {
		n := 0
		for {
			time.Sleep(time.Duration(n*n) * time.Second)

			newToken, err := requestToken()
			if err != nil {
				if n > 9 {
					log.Fatal("Failed to refresh token")
				}
				log.Println("Error while refreshing token")
				log.Println(err)
				n++
				continue
			}
			n = 0
			var _ string
			_ = <-tm.tokenChannel // Throw away old token
			tm.tokenChannel <- newToken
			time.Sleep(refreshInterval)
		}
	}()
}

func requestToken() (string, error) {
	url := "https://sms.gabia.com/oauth/token"

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	payload := bytes.NewBuffer([]byte(`grant_type=client_credentials`))

	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(SMS_ID+":"+SMS_API_KEY)))

	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var a tokenResponse

	err = json.Unmarshal(body, &a)
	if err != nil {
		return "", err
	}

	return string(a.AccessToken), nil
}

func sendSMS(token string, phone_number string, message string) error {
	url := "https://sms.gabia.com/api/send/sms"

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	payload := bytes.NewBuffer([]byte(`phone=` + phone_number + `&callback=` + SMS_PHONE_NUMBER + `&message=` + message + `&refkey=` + `[[hello]`))
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", base64.StdEncoding.EncodeToString([]byte(SMS_ID+":"+token)))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if res.StatusCode != 200 {
		return err
	}
	defer res.Body.Close()
	log.Println("Successfully Sent SMS to " + phone_number)
	return nil
}
