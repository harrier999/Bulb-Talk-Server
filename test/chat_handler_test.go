package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"server/internal/handler/chatting"
	"server/internal/models/message"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ChatService 인터페이스를 구현하는 mock 객체
type ChatServiceMock struct {
	mock.Mock
}

func (m *ChatServiceMock) SaveMessage(ctx context.Context, roomID string, msg message.Message) error {
	args := m.Called(ctx, roomID, msg)
	return args.Error(0)
}

func (m *ChatServiceMock) GetMessages(ctx context.Context, roomID string, lastMessageID int64) ([]message.Message, error) {
	args := m.Called(ctx, roomID, lastMessageID)
	return args.Get(0).([]message.Message), args.Error(1)
}

func (m *ChatServiceMock) GetMessagesByUUID(ctx context.Context, roomID string, lastMessageUUID uuid.UUID) ([]message.Message, error) {
	args := m.Called(ctx, roomID, lastMessageUUID)
	return args.Get(0).([]message.Message), args.Error(1)
}

func (m *ChatServiceMock) HandleWebSocketConnection(ctx context.Context, roomID, userID string, conn interface{}) error {
	args := m.Called(ctx, roomID, userID, conn)
	return args.Error(0)
}

func TestChatHandlerGetMessages(t *testing.T) {
	// mock 서비스 생성
	chatService := new(ChatServiceMock)

	// 핸들러 생성
	handler := chatting.NewChatHandler(chatService)

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

	// mock 서비스 동작 설정
	chatService.On("GetMessages", mock.Anything, roomID, lastMsgID).Return(messages, nil)

	// 테스트 요청 생성
	req, _ := http.NewRequest("GET", "/messages?roomId="+roomID+"&lastMessageId=100", nil)

	// 응답 레코더 생성
	rr := httptest.NewRecorder()

	// 핸들러 호출
	handler.GetMessages(rr, req)

	// 응답 검증
	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, true, response["success"])
	assert.NotNil(t, response["messages"])

	// mock 서비스 호출 확인
	chatService.AssertExpectations(t)
}

func TestChatHandlerGetMessagesByUUID(t *testing.T) {
	// mock 서비스 생성
	chatService := new(ChatServiceMock)

	// 핸들러 생성
	handler := chatting.NewChatHandler(chatService)

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

	// mock 서비스 동작 설정
	chatService.On("GetMessagesByUUID", mock.Anything, roomID, lastMsgUUID).Return(messages, nil)

	// 테스트 요청 생성
	req, _ := http.NewRequest("GET", "/messages?roomId="+roomID+"&lastMessageId="+lastMsgUUID.String(), nil)

	// 응답 레코더 생성
	rr := httptest.NewRecorder()

	// 핸들러 호출
	handler.GetMessages(rr, req)

	// 응답 검증
	assert.Equal(t, http.StatusOK, rr.Code)

	var response map[string]interface{}
	err := json.Unmarshal(rr.Body.Bytes(), &response)
	assert.NoError(t, err)

	assert.Equal(t, true, response["success"])
	assert.NotNil(t, response["messages"])

	// 모의 서비스 호출 확인
	chatService.AssertExpectations(t)
}
