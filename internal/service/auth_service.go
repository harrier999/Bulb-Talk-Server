package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"server/internal/models/orm"
	"server/internal/repository"
	"server/pkg/authenticator"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthServiceImpl struct {
	authRepo repository.AuthRepository
}

func NewAuthService(authRepo repository.AuthRepository) AuthService {
	return &AuthServiceImpl{
		authRepo: authRepo,
	}
}

func (s *AuthServiceImpl) RequestAuthNumber(ctx context.Context, phoneNumber, countryCode, deviceID string) error {

	rand.Seed(time.Now().UnixNano())
	authNumber := rand.Intn(900000) + 100000
	authNumberStr := fmt.Sprintf("%d", authNumber)

	auth := orm.AuthenticateMessage{
		CountryCode:        countryCode,
		PhoneNumber:        phoneNumber,
		RequestTime:        time.Now(),
		ExpireTime:         time.Now().Add(3 * time.Minute),
		DeviceID:           deviceID,
		AuthenticateNumber: authNumberStr,
		Trial:              0,
	}

	return s.authRepo.SaveAuthMessage(ctx, auth)
}

func (s *AuthServiceImpl) CheckAuthNumber(ctx context.Context, phoneNumber, countryCode, deviceID, authNumber string) (bool, error) {

	auth, err := s.authRepo.GetAuthMessage(ctx, phoneNumber, countryCode, deviceID)
	if err != nil {
		return false, err
	}

	if time.Now().After(auth.ExpireTime) {
		return false, errors.New("authentication number expired")
	}

	if auth.Trial >= 5 {
		return false, errors.New("too many trials")
	}

	err = s.authRepo.UpdateAuthTrial(ctx, auth.ID, auth.Trial+1)
	if err != nil {
		return false, err
	}

	if auth.AuthenticateNumber != authNumber {
		return false, nil
	}

	return true, nil
}

func (s *AuthServiceImpl) CreateToken(ctx context.Context, userID string, expiryHours int) (string, error) {
	return authenticator.CreateToken(userID, int64(expiryHours))
}

func (s *AuthServiceImpl) ValidateToken(ctx context.Context, token string) (string, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(authenticator.GetHMACSecret()), nil
	})

	if err != nil || !parsedToken.Valid {
		return "", errors.New("invalid token")
	}

	userID, ok := parsedToken.Claims.(jwt.MapClaims)[string(authenticator.ContextKeyUserID)].(string)
	if !ok {
		return "", errors.New("invalid token")
	}

	return userID, nil
}
