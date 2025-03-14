package room

import (
	"encoding/json"
	"net/http"
	"server/internal/models/orm"
	"server/internal/service"
	"server/pkg/authenticator"

	"github.com/google/uuid"
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
	Rooms []orm.Room `json:"rooms"`
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
	json.NewEncoder(w).Encode(RoomListResponse{Rooms: rooms})
}

type CreateRoomRequest struct {
	RoomName string      `json:"room_name"`
	Users    []uuid.UUID `json:"room_user_list"`
}

type CreateRoomResponse struct {
	RoomID uuid.UUID `json:"room_id"`
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
	json.NewEncoder(w).Encode(CreateRoomResponse{RoomID: room.ID})
}

type AddUserRequest struct {
	RoomID string `json:"room_id"`
	UserID string `json:"user_id"`
}

func (h *Handler) AddUser(w http.ResponseWriter, r *http.Request) {

	if _, err := authenticator.GetUserID(r); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req AddUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.RoomID == "" || req.UserID == "" {
		http.Error(w, "Missing room ID or user ID", http.StatusBadRequest)
		return
	}

	roomID, err := uuid.Parse(req.RoomID)
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
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

	w.WriteHeader(http.StatusOK)
}

type RemoveUserRequest struct {
	RoomID string `json:"room_id"`
	UserID string `json:"user_id"`
}

func (h *Handler) RemoveUser(w http.ResponseWriter, r *http.Request) {
	if _, err := authenticator.GetUserID(r); err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req RemoveUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.RoomID == "" || req.UserID == "" {
		http.Error(w, "Missing room ID or user ID", http.StatusBadRequest)
		return
	}

	roomID, err := uuid.Parse(req.RoomID)
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}
	targetUserID, err := uuid.Parse(req.UserID)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.roomService.RemoveUserFromRoom(r.Context(), roomID, targetUserID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
