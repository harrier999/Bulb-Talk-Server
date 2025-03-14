package chatting

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"server/internal/service"

	"github.com/gorilla/websocket"
)

var ALLOWED_ORIGIN_LIST = []string{"app.bulbtalk.com"}

type ChatHandler struct {
	chatService service.ChatService
}

func NewChatHandler(chatService service.ChatService) *ChatHandler {
	return &ChatHandler{
		chatService: chatService,
	}
}

var WebSocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		url, err := url.Parse(origin)
		if err != nil {
			log.Println(err)
			return false
		}
		for _, allowed_origin := range ALLOWED_ORIGIN_LIST {
			if url.Host == allowed_origin {
				return true
			}
		}
		return false
	},
}

type initialMessage struct {
	RoomID string `json:"room_id"`
	UserID string `json:"user_id"`
}

func (h *ChatHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := WebSocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
		return
	}

	_, msgBytes, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error reading initial message:", err)
		conn.Close()
		return
	}

	var initMsg initialMessage
	err = json.Unmarshal(msgBytes, &initMsg)
	if err != nil {
		log.Println("Error parsing initial message:", err)
		conn.Close()
		return
	}

	if initMsg.RoomID == "" || initMsg.UserID == "" {
		log.Println("Invalid room ID or user ID")
		conn.Close()
		return
	}

	err = h.chatService.HandleWebSocketConnection(r.Context(), initMsg.RoomID, initMsg.UserID, conn)
	if err != nil {
		log.Println("Error handling WebSocket connection:", err)
		conn.Close()
		return
	}
}

func (h *ChatHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("room_id")
	if roomID == "" {
		http.Error(w, "Missing room ID", http.StatusBadRequest)
		return
	}

	lastMessageIDStr := r.URL.Query().Get("last_message_id")
	var lastMessageID int64 = 0
	if lastMessageIDStr != "" {
		var err error
		lastMessageID, err = json.Number(lastMessageIDStr).Int64()
		if err != nil {
			http.Error(w, "Invalid last message ID", http.StatusBadRequest)
			return
		}
	}

	messages, err := h.chatService.GetMessages(r.Context(), roomID, lastMessageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
