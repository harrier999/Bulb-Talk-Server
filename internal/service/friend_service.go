package service

import (
	"context"
	"errors"
	"server/internal/models/orm"
	"server/internal/repository"

	"github.com/google/uuid"
)

type FriendServiceImpl struct {
	friendRepo repository.FriendRepository
	userRepo   repository.UserRepository
}

func NewFriendService(friendRepo repository.FriendRepository, userRepo repository.UserRepository) FriendService {
	return &FriendServiceImpl{
		friendRepo: friendRepo,
		userRepo:   userRepo,
	}
}

func (s *FriendServiceImpl) GetFriendList(ctx context.Context, userID uuid.UUID) ([]orm.Friend, error) {
	return s.friendRepo.GetFriendList(ctx, userID)
}

func (s *FriendServiceImpl) AddFriend(ctx context.Context, userID uuid.UUID, phoneNumber string) error {

	friend, err := s.userRepo.FindByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return errors.New("user not found")
	}

	exists, err := s.friendRepo.CheckIfFriendExists(ctx, userID, friend.ID)
	if err != nil {
		return err
	}
	if exists {
		return errors.New("friend already exists")
	}

	return s.friendRepo.AddFriend(ctx, userID, friend.ID)
}

func (s *FriendServiceImpl) BlockFriend(ctx context.Context, userID, friendID uuid.UUID) error {

	exists, err := s.friendRepo.CheckIfFriendExists(ctx, userID, friendID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("friend not found")
	}

	return s.friendRepo.BlockFriend(ctx, userID, friendID)
}

func (s *FriendServiceImpl) UnblockFriend(ctx context.Context, userID, friendID uuid.UUID) error {

	exists, err := s.friendRepo.CheckIfFriendExists(ctx, userID, friendID)
	if err != nil {
		return err
	}
	if !exists {
		return errors.New("friend not found")
	}

	return s.friendRepo.UnblockFriend(ctx, userID, friendID)
}
