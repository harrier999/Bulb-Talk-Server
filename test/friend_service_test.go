package test

import (
	"context"
	"server/internal/models/orm"
	"server/internal/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// FriendRepositoryMock은 FriendRepository 인터페이스를 구현하는 모의 객체입니다.
type FriendRepositoryMock struct {
	mock.Mock
}

func (m *FriendRepositoryMock) GetFriendList(ctx context.Context, userID uuid.UUID) ([]orm.Friend, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]orm.Friend), args.Error(1)
}

func (m *FriendRepositoryMock) AddFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	args := m.Called(ctx, userID, friendID)
	return args.Error(0)
}

func (m *FriendRepositoryMock) RemoveFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	args := m.Called(ctx, userID, friendID)
	return args.Error(0)
}

func (m *FriendRepositoryMock) BlockFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	args := m.Called(ctx, userID, friendID)
	return args.Error(0)
}

func (m *FriendRepositoryMock) UnblockFriend(ctx context.Context, userID, friendID uuid.UUID) error {
	args := m.Called(ctx, userID, friendID)
	return args.Error(0)
}

func (m *FriendRepositoryMock) GetBlockedList(ctx context.Context, userID uuid.UUID) ([]orm.Friend, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]orm.Friend), args.Error(1)
}

func (m *FriendRepositoryMock) CheckIfFriendExists(ctx context.Context, userID, friendID uuid.UUID) (bool, error) {
	args := m.Called(ctx, userID, friendID)
	return args.Bool(0), args.Error(1)
}

// UserRepositoryMock은 UserRepository 인터페이스를 구현하는 모의 객체입니다.
type UserRepositoryMock struct {
	mock.Mock
}

func (m *UserRepositoryMock) FindByID(ctx context.Context, id uuid.UUID) (orm.User, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(orm.User), args.Error(1)
}

func (m *UserRepositoryMock) FindByPhoneNumber(ctx context.Context, phoneNumber string) (orm.User, error) {
	args := m.Called(ctx, phoneNumber)
	return args.Get(0).(orm.User), args.Error(1)
}

func (m *UserRepositoryMock) Create(ctx context.Context, user orm.User) (orm.User, error) {
	args := m.Called(ctx, user)
	return args.Get(0).(orm.User), args.Error(1)
}

func (m *UserRepositoryMock) Update(ctx context.Context, user orm.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *UserRepositoryMock) Delete(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func TestGetFriendList(t *testing.T) {
	// 모의 리포지토리 생성
	friendRepo := new(FriendRepositoryMock)
	userRepo := new(UserRepositoryMock)

	// 서비스 생성
	friendService := service.NewFriendService(friendRepo, userRepo)

	// 테스트 데이터
	userID := uuid.New()
	friend1 := orm.Friend{
		UserID:    userID,
		FriendID:  uuid.New(),
		IsBlocked: false,
	}
	friend2 := orm.Friend{
		UserID:    userID,
		FriendID:  uuid.New(),
		IsBlocked: false,
	}
	friendList := []orm.Friend{friend1, friend2}

	// 모의 리포지토리 동작 설정
	friendRepo.On("GetFriendList", mock.Anything, userID).Return(friendList, nil)

	// 테스트 실행
	result, err := friendService.GetFriendList(context.Background(), userID)

	// 검증
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, friend1.FriendID, result[0].FriendID)
	assert.Equal(t, friend2.FriendID, result[1].FriendID)

	// 모의 리포지토리 호출 검증
	friendRepo.AssertExpectations(t)
}

func TestAddFriend(t *testing.T) {
	// 모의 리포지토리 생성
	friendRepo := new(FriendRepositoryMock)
	userRepo := new(UserRepositoryMock)

	// 서비스 생성
	friendService := service.NewFriendService(friendRepo, userRepo)

	// 테스트 데이터
	userID := uuid.New()
	phoneNumber := "1234567890"
	friendID := uuid.New()

	// 모의 리포지토리 동작 설정
	userRepo.On("FindByPhoneNumber", mock.Anything, phoneNumber).Return(orm.User{UUIDv7BaseModel: orm.UUIDv7BaseModel{ID: friendID}}, nil)
	friendRepo.On("CheckIfFriendExists", mock.Anything, userID, friendID).Return(false, nil)
	friendRepo.On("AddFriend", mock.Anything, userID, friendID).Return(nil)

	// 테스트 실행
	err := friendService.AddFriend(context.Background(), userID, phoneNumber)

	// 검증
	assert.NoError(t, err)

	// 모의 리포지토리 호출 검증
	friendRepo.AssertExpectations(t)
	userRepo.AssertExpectations(t)
}

func TestBlockFriend(t *testing.T) {
	// 모의 리포지토리 생성
	friendRepo := new(FriendRepositoryMock)
	userRepo := new(UserRepositoryMock)

	// 서비스 생성
	friendService := service.NewFriendService(friendRepo, userRepo)

	// 테스트 데이터
	userID := uuid.New()
	friendID := uuid.New()

	// 모의 리포지토리 동작 설정
	friendRepo.On("CheckIfFriendExists", mock.Anything, userID, friendID).Return(true, nil)
	friendRepo.On("BlockFriend", mock.Anything, userID, friendID).Return(nil)

	// 테스트 실행
	err := friendService.BlockFriend(context.Background(), userID, friendID)

	// 검증
	assert.NoError(t, err)

	// 모의 리포지토리 호출 검증
	friendRepo.AssertExpectations(t)
}

func TestUnblockFriend(t *testing.T) {
	// 모의 리포지토리 생성
	friendRepo := new(FriendRepositoryMock)
	userRepo := new(UserRepositoryMock)

	// 서비스 생성
	friendService := service.NewFriendService(friendRepo, userRepo)

	// 테스트 데이터
	userID := uuid.New()
	friendID := uuid.New()

	// 모의 리포지토리 동작 설정
	friendRepo.On("CheckIfFriendExists", mock.Anything, userID, friendID).Return(true, nil)
	friendRepo.On("UnblockFriend", mock.Anything, userID, friendID).Return(nil)

	// 테스트 실행
	err := friendService.UnblockFriend(context.Background(), userID, friendID)

	// 검증
	assert.NoError(t, err)

	// 모의 리포지토리 호출 검증
	friendRepo.AssertExpectations(t)
}
