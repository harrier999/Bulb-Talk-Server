package service

import (
	"context"
	"errors"
	"server/internal/models/orm"
	"server/internal/repository"
	"server/pkg/authenticator"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserServiceImpl struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &UserServiceImpl{
		userRepo: userRepo,
	}
}

func (s *UserServiceImpl) Register(ctx context.Context, username, password, phoneNumber, countryCode string) (orm.User, error) {

	_, err := s.userRepo.FindByPhoneNumber(ctx, phoneNumber)
	if err == nil {
		return orm.User{}, errors.New("user already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return orm.User{}, err
	}

	user := orm.User{
		UserName:     username,
		PasswordHash: string(hashedPassword),
		PhoneNumber:  phoneNumber,
		CountryCode:  countryCode,
	}

	return s.userRepo.Create(ctx, user)
}

func (s *UserServiceImpl) Login(ctx context.Context, phoneNumber, password string) (string, error) {

	user, err := s.userRepo.FindByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	token, err := authenticator.CreateToken(user.ID.String(), 24*14)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *UserServiceImpl) GetUserByID(ctx context.Context, id uuid.UUID) (orm.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *UserServiceImpl) GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (orm.User, error) {
	return s.userRepo.FindByPhoneNumber(ctx, phoneNumber)
}
