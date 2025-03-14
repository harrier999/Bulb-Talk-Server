package repository

import (
	"context"
	"server/internal/models/message"
	"server/internal/models/orm"

	"github.com/google/uuid"
)

type UserRepository interface {
	FindByID(ctx context.Context, id uuid.UUID) (orm.User, error)
	FindByPhoneNumber(ctx context.Context, phoneNumber string) (orm.User, error)
	Create(ctx context.Context, user orm.User) (orm.User, error)
	Update(ctx context.Context, user orm.User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type FriendRepository interface {
	GetFriendList(ctx context.Context, userID uuid.UUID) ([]orm.Friend, error)
	AddFriend(ctx context.Context, userID, friendID uuid.UUID) error
	BlockFriend(ctx context.Context, userID, friendID uuid.UUID) error
	UnblockFriend(ctx context.Context, userID, friendID uuid.UUID) error
	CheckIfFriendExists(ctx context.Context, userID, friendID uuid.UUID) (bool, error)
}

type RoomRepository interface {
	Create(ctx context.Context, room orm.Room) (orm.Room, error)
	FindByID(ctx context.Context, id uuid.UUID) (orm.Room, error)
	GetUserRooms(ctx context.Context, userID uuid.UUID) ([]orm.Room, error)
	AddUserToRoom(ctx context.Context, roomID, userID uuid.UUID) error
	RemoveUserFromRoom(ctx context.Context, roomID, userID uuid.UUID) error
	CreateRoomWithUsers(ctx context.Context, roomName string, userIDs []uuid.UUID) (uuid.UUID, error)
}

type MessageRepository interface {
	SaveMessage(ctx context.Context, roomID string, msg message.Message) error
	GetMessages(ctx context.Context, roomID string, lastMessageID int64) ([]message.Message, error)
}

type AuthRepository interface {
	SaveAuthMessage(ctx context.Context, auth orm.AuthenticateMessage) error
	GetAuthMessage(ctx context.Context, phoneNumber, countryCode, deviceID string) (orm.AuthenticateMessage, error)
	UpdateAuthTrial(ctx context.Context, id uuid.UUID, trial int) error
}
