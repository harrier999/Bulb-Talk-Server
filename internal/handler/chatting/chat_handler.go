package chatting

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"server/internal/models/message"
	"server/internal/service"
	"strings"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// 환경 변수에서 허용된 오리진 목록을 가져옵니다
func getAllowedOrigins() []string {
	allowedOriginsStr := os.Getenv("ALLOWED_ORIGINS")
	if allowedOriginsStr == "" {
		return []string{} // 기본값 없음
	}
	return strings.Split(allowedOriginsStr, ",")
}

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
		if origin == "" {
			return true // 개발 환경에서는 Origin이 없는 요청도 허용
		}

		parsedURL, err := url.Parse(origin)
		if err != nil {
			log.Println("Error parsing origin:", err)
			return false
		}

		allowedOrigins := getAllowedOrigins()

		// '*'가 허용된 오리진 목록에 있으면 모든 오리진 허용
		for _, allowedOrigin := range allowedOrigins {
			if strings.TrimSpace(allowedOrigin) == "*" {
				return true
			}
		}

		// 로컬호스트 확인 (디버깅 용도)
		host := parsedURL.Host
		if strings.HasPrefix(host, "localhost") ||
			strings.HasPrefix(host, "127.0.0.1") ||
			strings.HasSuffix(host, ".localhost") {
			return true
		}

		if len(allowedOrigins) == 0 {
			return true // 허용된 오리진이 없으면 모든 오리진 허용 (개발 환경용)
		}

		// 명시적으로 허용된 도메인 확인
		for _, allowedOrigin := range allowedOrigins {
			// 와일드카드 도메인 처리 (*.example.com)
			if strings.HasPrefix(allowedOrigin, "*.") {
				suffix := allowedOrigin[1:] // "*."를 제거
				if strings.HasSuffix(parsedURL.Host, suffix) {
					return true
				}
			} else if parsedURL.Host == allowedOrigin {
				return true
			}
		}

		log.Printf("Rejected WebSocket connection from origin: %s", origin)
		return false
	},
}

type initialMessage struct {
	RoomID string `json:"roomId"`
	UserID string `json:"userId"`
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

type MessageResponse struct {
	Success  bool              `json:"success"`
	Messages []message.Message `json:"messages"`
}

func (h *ChatHandler) GetMessages(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Query().Get("roomId")
	if roomID == "" {
		http.Error(w, "Missing room ID", http.StatusBadRequest)
		return
	}

	var messages []message.Message
	var err error

	lastMessageUUIDStr := r.URL.Query().Get("lastMessageId")
	if lastMessageUUIDStr != "" {
		// UUID 기반 조회
		lastMessageUUID, uuidErr := uuid.Parse(lastMessageUUIDStr)
		if uuidErr != nil {
			// UUID 파싱 실패 시 숫자 ID로 시도
			var lastMessageID int64 = 0
			lastMessageID, numErr := json.Number(lastMessageUUIDStr).Int64()
			if numErr != nil {
				http.Error(w, "Invalid last message ID", http.StatusBadRequest)
				return
			}
			messages, err = h.chatService.GetMessages(r.Context(), roomID, lastMessageID)
		} else {
			messages, err = h.chatService.GetMessagesByUUID(r.Context(), roomID, lastMessageUUID)
		}
	} else {
		// 모든 메시지 조회
		messages, err = h.chatService.GetMessages(r.Context(), roomID, 0)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(MessageResponse{Success: true, Messages: messages})
}
