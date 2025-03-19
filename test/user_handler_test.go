package test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"server/internal/handler/user"
	"server/internal/models/orm"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// UserServiceMock은 UserService 인터페이스를 구현하는 모의 객체입니다.
type UserServiceMock struct {
	mock.Mock
}

func (m *UserServiceMock) Register(ctx context.Context, username, password, phoneNumber, countryCode string) (orm.User, error) {
	args := m.Called(ctx, username, password, phoneNumber, countryCode)
	return args.Get(0).(orm.User), args.Error(1)
}

func (m *UserServiceMock) Login(ctx context.Context, phoneNumber, password string) (string, error) {
	args := m.Called(ctx, phoneNumber, password)
	return args.String(0), args.Error(1)
}

func (m *UserServiceMock) GetUserByID(ctx context.Context, id uuid.UUID) (orm.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(orm.User), args.Error(1)
}

func (m *UserServiceMock) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (orm.User, error) {
	args := m.Called(ctx, phoneNumber)
	return args.Get(0).(orm.User), args.Error(1)
}

// AuthServiceMock은 AuthService 인터페이스를 구현하는 모의 객체입니다.
type AuthServiceMock struct {
	mock.Mock
}

func (m *AuthServiceMock) RequestAuthNumber(ctx context.Context, phoneNumber, countryCode, deviceID string) error {
	args := m.Called(ctx, phoneNumber, countryCode, deviceID)
	return args.Error(0)
}

func (m *AuthServiceMock) CheckAuthNumber(ctx context.Context, phoneNumber, countryCode, deviceID, authNumber string) (bool, error) {
	args := m.Called(ctx, phoneNumber, countryCode, deviceID, authNumber)
	return args.Bool(0), args.Error(1)
}

func (m *AuthServiceMock) CreateToken(ctx context.Context, userID string, expiryHours int) (string, error) {
	args := m.Called(ctx, userID, expiryHours)
	return args.String(0), args.Error(1)
}

func (m *AuthServiceMock) ValidateToken(ctx context.Context, token string) (string, error) {
	args := m.Called(ctx, token)
	return args.String(0), args.Error(1)
}

func TestSignUp(t *testing.T) {
	// 모의 서비스 생성
	userService := new(UserServiceMock)
	authService := new(AuthServiceMock)

	// 핸들러 생성
	handler := user.NewHandler(userService, authService)

	// 테스트 사용자 데이터
	testUser := orm.User{
		UserName:    "testuser",
		PhoneNumber: "1234567890",
		CountryCode: "82",
	}
	testUser.ID = uuid.New()

	// 모의 서비스 동작 설정
	userService.On("Register", mock.Anything, "testuser", "password123", "1234567890", "82").Return(testUser, nil)

	// 테스트 요청 생성
	reqBody := map[string]string{
		"username":    "testuser",
		"password":    "password123",
		"phoneNumber": "1234567890",
		"countryCode": "82",
	}
	reqJSON, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/signup", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")

	// 응답 레코더 생성
	rr := httptest.NewRecorder()

	// 핸들러 호출
	handler.SignUp(rr, req)

	// 응답 검증
	assert.Equal(t, http.StatusCreated, rr.Code)

	// 응답 본문 검증
	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, true, response["success"])
	assert.NotNil(t, response["user"])

	// 모의 서비스 호출 검증
	userService.AssertExpectations(t)
}

func TestLogin(t *testing.T) {
	// 모의 서비스 생성
	userService := new(UserServiceMock)
	authService := new(AuthServiceMock)

	// 핸들러 생성
	handler := user.NewHandler(userService, authService)

	// 모의 서비스 동작 설정
	userService.On("Login", mock.Anything, "1234567890", "password123").Return("test-jwt-token", nil)

	// 테스트 요청 생성
	reqBody := map[string]string{
		"phoneNumber": "1234567890",
		"password":    "password123",
	}
	reqJSON, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")

	// 응답 레코더 생성
	rr := httptest.NewRecorder()

	// 핸들러 호출
	handler.Login(rr, req)

	// 응답 검증
	assert.Equal(t, http.StatusOK, rr.Code)

	// 모의 서비스 호출 검증
	userService.AssertExpectations(t)
}

func TestRequestAuthNumber(t *testing.T) {
	// 모의 서비스 생성
	userService := new(UserServiceMock)
	authService := new(AuthServiceMock)

	// 핸들러 생성
	handler := user.NewHandler(userService, authService)

	// 모의 서비스 동작 설정
	authService.On("RequestAuthNumber", mock.Anything, "1234567890", "82", "device123").Return(nil)

	// 테스트 요청 생성
	reqBody := map[string]string{
		"phoneNumber": "1234567890",
		"countryCode": "82",
		"deviceId":    "device123",
	}
	reqJSON, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/authenticate", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")

	// 응답 레코더 생성
	rr := httptest.NewRecorder()

	// 핸들러 호출
	handler.RequestAuthNumber(rr, req)

	// 응답 검증
	assert.Equal(t, http.StatusOK, rr.Code)

	// 모의 서비스 호출 검증
	authService.AssertExpectations(t)
}

func TestCheckAuthNumber(t *testing.T) {
	// 모의 서비스 생성
	userService := new(UserServiceMock)
	authService := new(AuthServiceMock)

	// 핸들러 생성
	handler := user.NewHandler(userService, authService)

	// 모의 서비스 동작 설정
	authService.On("CheckAuthNumber", mock.Anything, "1234567890", "82", "device123", "123456").Return(true, nil)

	// 테스트 요청 생성
	reqBody := map[string]string{
		"phoneNumber": "1234567890",
		"countryCode": "82",
		"deviceId":    "device123",
		"authNumber":  "123456",
	}
	reqJSON, _ := json.Marshal(reqBody)
	req, _ := http.NewRequest("POST", "/checkauth", bytes.NewBuffer(reqJSON))
	req.Header.Set("Content-Type", "application/json")

	// 응답 레코더 생성
	rr := httptest.NewRecorder()

	// 핸들러 호출
	handler.CheckAuthNumber(rr, req)

	// 응답 검증
	assert.Equal(t, http.StatusOK, rr.Code)

	// 모의 서비스 호출 검증
	authService.AssertExpectations(t)
}
