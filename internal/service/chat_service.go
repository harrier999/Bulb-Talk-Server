package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"server/internal/models/message"
	"server/internal/repository"
	"sync"

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

func (s *ChatServiceImpl) HandleWebSocketConnection(ctx context.Context, roomID, userID string, conn interface{}) error {

	wsConn, ok := conn.(*websocket.Conn)
	if !ok {
		return errors.New("invalid connection type")
	}

	s.addConnection(roomID, userID, wsConn)

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

func (s *ChatServiceImpl) handleMessages(ctx context.Context, roomID, userID string, conn *websocket.Conn) {
	defer func() {
		conn.Close()
		s.removeConnection(roomID, userID)
	}()

	for {

		_, msgBytes, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Error reading message:", err)
			}
			break
		}

		var baseMsg struct {
			Type string `json:"type"`
		}
		err = json.Unmarshal(msgBytes, &baseMsg)
		if err != nil {
			log.Println("Error parsing message type:", err)
			continue
		}

		var msg message.Message
		switch baseMsg.Type {
		case "text":
			textMsg := &message.TextMessage{}
			textMsg.FromJson(msgBytes)
			textMsg.Author = message.User{Id: userID}
			textMsg.RoomId = roomID
			msg = textMsg
		case "image":
			imageMsg := &message.ImageMessage{}
			imageMsg.FromJson(msgBytes)
			imageMsg.Author = message.User{Id: userID}
			imageMsg.RoomId = roomID
			msg = imageMsg
		default:
			log.Println("Unknown message type:", baseMsg.Type)
			continue
		}

		err = s.SaveMessage(ctx, roomID, msg)
		if err != nil {
			log.Println("Error saving message:", err)
		}
	}
}
