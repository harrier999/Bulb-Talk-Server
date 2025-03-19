package room

import (
	"encoding/json"
	"net/http"
	"server/internal/models/orm"
	"server/internal/service"
	"server/pkg/authenticator"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Handler struct {
	roomService service.RoomService
}

func NewHandler(roomService service.RoomService) *Handler {
	return &Handler{
		roomService: roomService,
	}
}

type RoomListResponse struct {
	Success bool       `json:"success"`
	Rooms   []orm.Room `json:"rooms"`
}

func (h *Handler) GetRoomList(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticator.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	rooms, err := h.roomService.GetUserRooms(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(RoomListResponse{Success: true, Rooms: rooms})
}

type CreateRoomRequest struct {
	RoomName string      `json:"roomName"`
	Users    []uuid.UUID `json:"roomUserList"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
}

type CreateRoomResponse struct {
	Success bool      `json:"success"`
	RoomID  uuid.UUID `json:"roomId"`
}

func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {

	userID, err := authenticator.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req CreateRoomRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.RoomName == "" {
		req.RoomName = "New Room"
	}
	if len(req.Users) < 1 {
		http.Error(w, "At least one user is required", http.StatusBadRequest)
		return
	}

	room, err := h.roomService.CreateRoom(r.Context(), req.RoomName, userID, req.Users)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(CreateRoomResponse{Success: true, RoomID: room.ID})
}

type AddUserRequest struct {
	UserID string `json:"userId"`
}

func (h *Handler) AddUser(w http.ResponseWriter, r *http.Request) {
	if _, err := authenticator.GetUserID(r); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// URL에서 roomId 파라미터 추출
	vars := mux.Vars(r)
	roomIDStr := vars["roomId"]

	if roomIDStr == "" {
		http.Error(w, "Missing room ID", http.StatusBadRequest)
		return
	}

	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}

	// 요청 본문에서 userId 추출
	var req AddUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserID == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	targetUserID, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.roomService.AddUserToRoom(r.Context(), roomID, targetUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SuccessResponse{Success: true})
}

func (h *Handler) RemoveUser(w http.ResponseWriter, r *http.Request) {
	if _, err := authenticator.GetUserID(r); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// URL에서 roomId와 userId 파라미터 추출
	vars := mux.Vars(r)
	roomIDStr := vars["roomId"]
	userIDStr := vars["userId"]

	if roomIDStr == "" || userIDStr == "" {
		http.Error(w, "Missing room ID or user ID", http.StatusBadRequest)
		return
	}

	roomID, err := uuid.Parse(roomIDStr)
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}

	targetUserID, err := uuid.Parse(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.roomService.RemoveUserFromRoom(r.Context(), roomID, targetUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SuccessResponse{Success: true})
}
