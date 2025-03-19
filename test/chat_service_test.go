package test

import (
	"context"
	"server/internal/models/message"
	"server/internal/models/orm"
	"server/internal/service"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MessageRepositoryMock은 MessageRepository 인터페이스를 구현하는 모의 객체입니다.
type MessageRepositoryMock struct {
	mock.Mock
}

func (m *MessageRepositoryMock) SaveMessage(ctx context.Context, roomID string, msg message.Message) error {
	args := m.Called(ctx, roomID, msg)
	return args.Error(0)
}

func (m *MessageRepositoryMock) GetMessages(ctx context.Context, roomID string, lastMessageID int64) ([]message.Message, error) {
	args := m.Called(ctx, roomID, lastMessageID)
	return args.Get(0).([]message.Message), args.Error(1)
}

func (m *MessageRepositoryMock) GetMessagesByUUID(ctx context.Context, roomID string, lastMessageUUID uuid.UUID) ([]message.Message, error) {
	args := m.Called(ctx, roomID, lastMessageUUID)
	return args.Get(0).([]message.Message), args.Error(1)
}

// RoomRepositoryMock은 RoomRepository 인터페이스를 구현하는 모의 객체입니다.
type RoomRepositoryMock struct {
	mock.Mock
}

func (m *RoomRepositoryMock) Create(ctx context.Context, room orm.Room) (orm.Room, error) {
	args := m.Called(ctx, room)
	return args.Get(0).(orm.Room), args.Error(1)
}

func (m *RoomRepositoryMock) FindByID(ctx context.Context, id uuid.UUID) (orm.Room, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(orm.Room), args.Error(1)
}

func (m *RoomRepositoryMock) GetUserRooms(ctx context.Context, userID uuid.UUID) ([]orm.Room, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]orm.Room), args.Error(1)
}

func (m *RoomRepositoryMock) AddUserToRoom(ctx context.Context, roomID, userID uuid.UUID) error {
	args := m.Called(ctx, roomID, userID)
	return args.Error(0)
}

func (m *RoomRepositoryMock) RemoveUserFromRoom(ctx context.Context, roomID, userID uuid.UUID) error {
	args := m.Called(ctx, roomID, userID)
	return args.Error(0)
}

func (m *RoomRepositoryMock) CreateRoomWithUsers(ctx context.Context, roomName string, userIDs []uuid.UUID) (uuid.UUID, error) {
	args := m.Called(ctx, roomName, userIDs)
	return args.Get(0).(uuid.UUID), args.Error(1)
}

func TestSaveMessage(t *testing.T) {
	// 모의 리포지토리 생성
	msgRepo := new(MessageRepositoryMock)

	// 서비스 생성
	chatService := service.NewChatService(msgRepo)

	// 테스트 데이터
	roomID := "room-123"

	// 테스트 메시지 생성
	msg := &message.TextMessage{
		BaseMessage: message.BaseMessage{
			RoomId: roomID,
			Type:   "text",
			Author: message.User{Id: "user-123"},
		},
		Content: "Hello, world!",
	}

	// 모의 리포지토리 동작 설정
	msgRepo.On("SaveMessage", mock.Anything, roomID, mock.AnythingOfType("*message.TextMessage")).Return(nil)

	// 테스트 실행
	err := chatService.SaveMessage(context.Background(), roomID, msg)

	// 검증
	assert.NoError(t, err)

	// 모의 리포지토리 호출 검증
	msgRepo.AssertExpectations(t)
}

func TestGetMessagesWithID(t *testing.T) {
	// 모의 리포지토리 생성
	msgRepo := new(MessageRepositoryMock)

	// 서비스 생성
	chatService := service.NewChatService(msgRepo)

	// 테스트 데이터
	roomID := "room-123"
	lastMsgID := int64(100)

	msg1 := &message.BaseMessage{
		RoomId: roomID,
		Type:   "text",
		Author: message.User{Id: "user-123"},
	}

	msg2 := &message.BaseMessage{
		RoomId: roomID,
		Type:   "text",
		Author: message.User{Id: "user-456"},
	}

	messages := []message.Message{msg1, msg2}

	// 모의 리포지토리 동작 설정
	msgRepo.On("GetMessages", mock.Anything, roomID, lastMsgID).Return(messages, nil)

	// 테스트 실행
	result, err := chatService.GetMessages(context.Background(), roomID, lastMsgID)

	// 검증
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))

	// 모의 리포지토리 호출 검증
	msgRepo.AssertExpectations(t)
}

func TestGetMessagesWithUUID(t *testing.T) {
	// 모의 리포지토리 생성
	msgRepo := new(MessageRepositoryMock)

	// 서비스 생성
	chatService := service.NewChatService(msgRepo)

	// 테스트 데이터
	roomID := "room-123"
	lastMsgUUID := uuid.New()

	msg1 := &message.BaseMessage{
		RoomId: roomID,
		Type:   "text",
		Author: message.User{Id: "user-123"},
	}

	msg2 := &message.BaseMessage{
		RoomId: roomID,
		Type:   "text",
		Author: message.User{Id: "user-456"},
	}

	messages := []message.Message{msg1, msg2}

	// 모의 리포지토리 동작 설정
	msgRepo.On("GetMessagesByUUID", mock.Anything, roomID, lastMsgUUID).Return(messages, nil)

	// 테스트 실행
	result, err := chatService.GetMessagesByUUID(context.Background(), roomID, lastMsgUUID)

	// 검증
	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))

	// 모의 리포지토리 호출 검증
	msgRepo.AssertExpectations(t)
}
