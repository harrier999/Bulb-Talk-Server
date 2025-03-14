package chatting

import (
	"context"
	"net/http"
	"net/http/httptest"
	"server/internal/models/message"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockChatService struct {
	mock.Mock
}

func (m *MockChatService) SaveMessage(ctx context.Context, roomID string, msg message.Message) error {
	args := m.Called(ctx, roomID, msg)
	return args.Error(0)
}

func (m *MockChatService) GetMessages(ctx context.Context, roomID string, lastMessageID int64) ([]message.Message, error) {
	args := m.Called(ctx, roomID, lastMessageID)
	return args.Get(0).([]message.Message), args.Error(1)
}

func (m *MockChatService) HandleWebSocketConnection(ctx context.Context, roomID, userID string, conn interface{}) error {
	args := m.Called(ctx, roomID, userID, conn)
	return args.Error(0)
}

func TestGetMessages(t *testing.T) {
	mockService := new(MockChatService)

	messages := []message.Message{}

	mockService.On("GetMessages", mock.Anything, "room1", int64(0)).Return(messages, nil)

	handler := NewChatHandler(mockService)

	req, err := http.NewRequest("GET", "/messages?room_id=room1", nil)
	assert.NoError(t, err)

	rr := httptest.NewRecorder()

	handler.GetMessages(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	mockService.AssertExpectations(t)
}
