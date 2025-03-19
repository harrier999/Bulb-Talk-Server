package service

import (
	"context"
	"server/internal/models/message"
	"server/internal/models/orm"

	"github.com/google/uuid"
)

type UserService interface {
	Register(ctx context.Context, username, password, phoneNumber, countryCode string) (orm.User, error)
	Login(ctx context.Context, phoneNumber, password string) (string, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (orm.User, error)
	GetUserByPhoneNumber(ctx context.Context, phoneNumber string) (orm.User, error)
}

type AuthService interface {
	RequestAuthNumber(ctx context.Context, phoneNumber, countryCode, deviceID string) error
	CheckAuthNumber(ctx context.Context, phoneNumber, countryCode, deviceID, authNumber string) (bool, error)
	CreateToken(ctx context.Context, userID string, expiryHours int) (string, error)
	ValidateToken(ctx context.Context, token string) (string, error)
}

type FriendService interface {
	GetFriendList(ctx context.Context, userID uuid.UUID) ([]orm.Friend, error)
	AddFriend(ctx context.Context, userID uuid.UUID, phoneNumber string) error
	BlockFriend(ctx context.Context, userID, friendID uuid.UUID) error
	UnblockFriend(ctx context.Context, userID, friendID uuid.UUID) error
}

type RoomService interface {
	CreateRoom(ctx context.Context, name string, creatorID uuid.UUID, participantIDs []uuid.UUID) (orm.Room, error)
	GetRoomByID(ctx context.Context, id uuid.UUID) (orm.Room, error)
	GetUserRooms(ctx context.Context, userID uuid.UUID) ([]orm.Room, error)
	AddUserToRoom(ctx context.Context, roomID, userID uuid.UUID) error
	RemoveUserFromRoom(ctx context.Context, roomID, userID uuid.UUID) error
}

type ChatService interface {
	SaveMessage(ctx context.Context, roomID string, msg message.Message) error
	GetMessages(ctx context.Context, roomID string, lastMessageID int64) ([]message.Message, error)
	GetMessagesByUUID(ctx context.Context, roomID string, lastMessageUUID uuid.UUID) ([]message.Message, error)
	HandleWebSocketConnection(ctx context.Context, roomID, userID string, conn interface{}) error
}
