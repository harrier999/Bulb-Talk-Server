package service

import (
	"context"
	"errors"
	"server/internal/models/orm"
	"server/internal/repository"

	"github.com/google/uuid"
)

type RoomServiceImpl struct {
	roomRepo repository.RoomRepository
}

func NewRoomService(roomRepo repository.RoomRepository) RoomService {
	return &RoomServiceImpl{
		roomRepo: roomRepo,
	}
}

func (s *RoomServiceImpl) CreateRoom(ctx context.Context, name string, creatorID uuid.UUID, participantIDs []uuid.UUID) (orm.Room, error) {

	if name == "" {
		return orm.Room{}, errors.New("room name is required")
	}
	if len(participantIDs) < 1 {
		return orm.Room{}, errors.New("at least one participant is required")
	}

	creatorIncluded := false
	for _, id := range participantIDs {
		if id == creatorID {
			creatorIncluded = true
			break
		}
	}
	if !creatorIncluded {
		participantIDs = append(participantIDs, creatorID)
	}

	roomID, err := s.roomRepo.CreateRoomWithUsers(ctx, name, participantIDs)
	if err != nil {
		return orm.Room{}, err
	}

	return s.roomRepo.FindByID(ctx, roomID)
}

func (s *RoomServiceImpl) GetRoomByID(ctx context.Context, id uuid.UUID) (orm.Room, error) {
	return s.roomRepo.FindByID(ctx, id)
}

func (s *RoomServiceImpl) GetUserRooms(ctx context.Context, userID uuid.UUID) ([]orm.Room, error) {
	return s.roomRepo.GetUserRooms(ctx, userID)
}

func (s *RoomServiceImpl) AddUserToRoom(ctx context.Context, roomID, userID uuid.UUID) error {

	_, err := s.roomRepo.FindByID(ctx, roomID)
	if err != nil {
		return errors.New("room not found")
	}

	return s.roomRepo.AddUserToRoom(ctx, roomID, userID)
}

func (s *RoomServiceImpl) RemoveUserFromRoom(ctx context.Context, roomID, userID uuid.UUID) error {

	_, err := s.roomRepo.FindByID(ctx, roomID)
	if err != nil {
		return errors.New("room not found")
	}

	return s.roomRepo.RemoveUserFromRoom(ctx, roomID, userID)
}
