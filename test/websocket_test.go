package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"server/internal/models/message"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// WebSocketChatServiceMock은 ChatService 인터페이스를 구현하는 모의 객체입니다.
type WebSocketChatServiceMock struct {
	mock.Mock
}

func (m *WebSocketChatServiceMock) SaveMessage(ctx context.Context, roomID string, msg message.Message) error {
	args := m.Called(ctx, roomID, msg)
	return args.Error(0)
}

func (m *WebSocketChatServiceMock) GetMessages(ctx context.Context, roomID string, lastMessageID int64) ([]message.Message, error) {
	args := m.Called(ctx, roomID, lastMessageID)
	return args.Get(0).([]message.Message), args.Error(1)
}

func (m *WebSocketChatServiceMock) GetMessagesByUUID(ctx context.Context, roomID string, lastMessageUUID uuid.UUID) ([]message.Message, error) {
	args := m.Called(ctx, roomID, lastMessageUUID)
	return args.Get(0).([]message.Message), args.Error(1)
}

func (m *WebSocketChatServiceMock) HandleWebSocketConnection(ctx context.Context, roomID, userID string, conn interface{}) error {
	args := m.Called(ctx, roomID, userID, conn)
	return args.Error(0)
}

// 간단한 WebSocket 핸들러 구현
func webSocketHandler(w http.ResponseWriter, r *http.Request) {
	// WebSocket 업그레이드
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "WebSocket 업그레이드 실패", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	// 메시지 처리 루프
	for {
		// 메시지 수신
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			return
		}

		// 메시지 파싱
		var msgData map[string]interface{}
		json.Unmarshal(p, &msgData)

		// 메시지 타입에 따라 처리
		switch msgData["type"] {
		case "message":
			// 메시지 에코
			conn.WriteMessage(messageType, p)
		case "typing":
			// 타이핑 상태 에코
			conn.WriteMessage(messageType, p)
		case "join":
			// 입장 메시지 에코
			conn.WriteMessage(messageType, p)
		case "leave":
			// 퇴장 메시지 에코
			conn.WriteMessage(messageType, p)
		}
	}
}

func TestWebSocketConnection(t *testing.T) {
	// 테스트 서버 설정
	server := httptest.NewServer(http.HandlerFunc(webSocketHandler))
	defer server.Close()

	// WebSocket URL 생성
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// WebSocket 연결
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer conn.Close()

	// 연결 성공 확인
	assert.NotNil(t, conn)
}

func TestWebSocketMessageSending(t *testing.T) {
	// 테스트 서버 설정
	server := httptest.NewServer(http.HandlerFunc(webSocketHandler))
	defer server.Close()

	// WebSocket URL 생성
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// WebSocket 연결
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer conn.Close()

	// 메시지 전송
	roomID := "room-123"
	messageData := map[string]interface{}{
		"type":    "message",
		"roomId":  roomID,
		"content": "Hello, WebSocket!",
	}

	err = conn.WriteJSON(messageData)
	assert.NoError(t, err)

	// 응답 수신 (타임아웃 설정)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	var response map[string]interface{}
	err = conn.ReadJSON(&response)
	assert.NoError(t, err)

	// 응답 검증
	assert.Equal(t, "message", response["type"])
	assert.Equal(t, roomID, response["roomId"])
	assert.Equal(t, "Hello, WebSocket!", response["content"])
}

func TestWebSocketTypingStatus(t *testing.T) {
	// 테스트 서버 설정
	server := httptest.NewServer(http.HandlerFunc(webSocketHandler))
	defer server.Close()

	// WebSocket URL 생성
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// WebSocket 연결
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer conn.Close()

	// 타이핑 상태 메시지 전송
	roomID := "room-123"
	typingData := map[string]interface{}{
		"type":   "typing",
		"roomId": roomID,
		"status": true,
	}

	err = conn.WriteJSON(typingData)
	assert.NoError(t, err)

	// 응답 수신 (타임아웃 설정)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	var response map[string]interface{}
	err = conn.ReadJSON(&response)
	assert.NoError(t, err)

	// 응답 검증
	assert.Equal(t, "typing", response["type"])
	assert.Equal(t, roomID, response["roomId"])
	assert.Equal(t, true, response["status"])
}

func TestWebSocketUserJoinLeave(t *testing.T) {
	// 테스트 서버 설정
	server := httptest.NewServer(http.HandlerFunc(webSocketHandler))
	defer server.Close()

	// WebSocket URL 생성
	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")

	// WebSocket 연결
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial(wsURL, nil)
	assert.NoError(t, err)
	defer conn.Close()

	// 사용자 입장 메시지 전송
	roomID := "room-123"
	joinData := map[string]interface{}{
		"type":   "join",
		"roomId": roomID,
	}

	err = conn.WriteJSON(joinData)
	assert.NoError(t, err)

	// 응답 수신 (타임아웃 설정)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	var joinResponse map[string]interface{}
	err = conn.ReadJSON(&joinResponse)
	assert.NoError(t, err)

	// 응답 검증
	assert.Equal(t, "join", joinResponse["type"])
	assert.Equal(t, roomID, joinResponse["roomId"])

	// 사용자 퇴장 메시지 전송
	leaveData := map[string]interface{}{
		"type":   "leave",
		"roomId": roomID,
	}

	err = conn.WriteJSON(leaveData)
	assert.NoError(t, err)

	// 응답 수신 (타임아웃 설정)
	conn.SetReadDeadline(time.Now().Add(2 * time.Second))

	var leaveResponse map[string]interface{}
	err = conn.ReadJSON(&leaveResponse)
	assert.NoError(t, err)

	// 응답 검증
	assert.Equal(t, "leave", leaveResponse["type"])
	assert.Equal(t, roomID, leaveResponse["roomId"])
}
