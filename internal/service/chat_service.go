package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"server/internal/models/message"
	"server/internal/repository"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type ChatServiceImpl struct {
	messageRepo repository.MessageRepository

	connections     map[string]map[string]*websocket.Conn
	connectionMutex sync.RWMutex
}

func NewChatService(messageRepo repository.MessageRepository) ChatService {
	return &ChatServiceImpl{
		messageRepo:     messageRepo,
		connections:     make(map[string]map[string]*websocket.Conn),
		connectionMutex: sync.RWMutex{},
	}
}

func (s *ChatServiceImpl) SaveMessage(ctx context.Context, roomID string, msg message.Message) error {
	err := s.messageRepo.SaveMessage(ctx, roomID, msg)
	if err != nil {
		return err
	}

	s.broadcastMessage(roomID, msg)

	return nil
}

func (s *ChatServiceImpl) GetMessages(ctx context.Context, roomID string, lastMessageID int64) ([]message.Message, error) {
	return s.messageRepo.GetMessages(ctx, roomID, lastMessageID)
}

func (s *ChatServiceImpl) GetMessagesByUUID(ctx context.Context, roomID string, lastMessageUUID uuid.UUID) ([]message.Message, error) {
	return s.messageRepo.GetMessagesByUUID(ctx, roomID, lastMessageUUID)
}

func (s *ChatServiceImpl) HandleWebSocketConnection(ctx context.Context, roomID, userID string, conn interface{}) error {
	wsConn, ok := conn.(*websocket.Conn)
	if !ok {
		return errors.New("invalid connection type")
	}

	s.addConnection(roomID, userID, wsConn)

	// 사용자 입장 이벤트 전송
	s.sendUserJoinedEvent(roomID, userID)

	go s.handleMessages(ctx, roomID, userID, wsConn)

	return nil
}

func (s *ChatServiceImpl) addConnection(roomID, userID string, conn *websocket.Conn) {
	s.connectionMutex.Lock()
	defer s.connectionMutex.Unlock()

	if _, ok := s.connections[roomID]; !ok {
		s.connections[roomID] = make(map[string]*websocket.Conn)
	}

	s.connections[roomID][userID] = conn
}

func (s *ChatServiceImpl) removeConnection(roomID, userID string) {
	s.connectionMutex.Lock()
	defer s.connectionMutex.Unlock()

	if room, ok := s.connections[roomID]; ok {
		delete(room, userID)

		if len(room) == 0 {
			delete(s.connections, roomID)
		}
	}
}

func (s *ChatServiceImpl) broadcastMessage(roomID string, msg message.Message) {
	s.connectionMutex.RLock()
	defer s.connectionMutex.RUnlock()

	room, ok := s.connections[roomID]
	if !ok {
		return
	}

	msgJSON := []byte(msg.ToJson())

	for _, conn := range room {
		err := conn.WriteMessage(websocket.TextMessage, msgJSON)
		if err != nil {
			log.Println("Error sending message:", err)
		}
	}
}

func (s *ChatServiceImpl) broadcastTypingStatus(roomID, userID string, isTyping bool) {
	s.connectionMutex.RLock()
	defer s.connectionMutex.RUnlock()

	room, ok := s.connections[roomID]
	if !ok {
		return
	}

	typingEvent := map[string]interface{}{
		"type":     "typing",
		"roomId":   roomID,
		"userId":   userID,
		"isTyping": isTyping,
	}

	msgJSON, _ := json.Marshal(typingEvent)

	for id, conn := range room {
		if id != userID { // 자신에게는 타이핑 상태를 보내지 않음
			err := conn.WriteMessage(websocket.TextMessage, msgJSON)
			if err != nil {
				log.Println("Error sending typing status:", err)
			}
		}
	}
}

func (s *ChatServiceImpl) sendUserJoinedEvent(roomID, userID string) {
	s.connectionMutex.RLock()
	defer s.connectionMutex.RUnlock()

	room, ok := s.connections[roomID]
	if !ok {
		return
	}

	joinEvent := map[string]interface{}{
		"type":      "userJoined",
		"roomId":    roomID,
		"userId":    userID,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	msgJSON, _ := json.Marshal(joinEvent)

	for id, conn := range room {
		if id != userID { // 자신에게는 입장 이벤트를 보내지 않음
			err := conn.WriteMessage(websocket.TextMessage, msgJSON)
			if err != nil {
				log.Println("Error sending user joined event:", err)
			}
		}
	}
}

func (s *ChatServiceImpl) sendUserLeftEvent(roomID, userID string) {
	s.connectionMutex.RLock()
	defer s.connectionMutex.RUnlock()

	room, ok := s.connections[roomID]
	if !ok {
		return
	}

	leftEvent := map[string]interface{}{
		"type":      "userLeft",
		"roomId":    roomID,
		"userId":    userID,
		"timestamp": time.Now().Format(time.RFC3339),
	}

	msgJSON, _ := json.Marshal(leftEvent)

	for _, conn := range room {
		err := conn.WriteMessage(websocket.TextMessage, msgJSON)
		if err != nil {
			log.Println("Error sending user left event:", err)
		}
	}
}

// WebSocketMessage는 WebSocket을 통해 주고받는 메시지의 구조를 정의합니다.
type WebSocketMessage struct {
	Type     string `json:"type"`
	RoomId   string `json:"roomId"`
	IsTyping bool   `json:"isTyping,omitempty"`
	Content  string `json:"content,omitempty"`
}

func (s *ChatServiceImpl) handleMessages(ctx context.Context, roomID, userID string, conn *websocket.Conn) {
	defer func() {
		conn.Close()
		s.removeConnection(roomID, userID)
		// 사용자 퇴장 이벤트 전송
		s.sendUserLeftEvent(roomID, userID)
	}()

	for {
		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Error reading message:", err)
			}
			break
		}

		var baseMsg WebSocketMessage
		err = json.Unmarshal(msgBytes, &baseMsg)
		if err != nil {
			log.Println("Error parsing message type:", err)
			continue
		}

		switch baseMsg.Type {
		case "message":
			textMsg := &message.TextMessage{}
			textMsg.Id, _ = uuid.NewV7()
			textMsg.Type = "message"
			textMsg.Author = message.User{Id: userID}
			textMsg.RoomId = roomID
			textMsg.Content = baseMsg.Content

			err = s.SaveMessage(ctx, roomID, textMsg)
			if err != nil {
				log.Println("Error saving message:", err)
			}
		case "typing":
			s.broadcastTypingStatus(roomID, userID, baseMsg.IsTyping)
		case "image":
			imageMsg := &message.ImageMessage{}
			imageMsg.FromJson(msgBytes)
			imageMsg.Author = message.User{Id: userID}
			imageMsg.RoomId = roomID
			imageMsg.Id, _ = uuid.NewV7()

			err = s.SaveMessage(ctx, roomID, imageMsg)
			if err != nil {
				log.Println("Error saving message:", err)
			}
		default:
			log.Println("Unknown message type:", baseMsg.Type)
			continue
		}
	}
}
