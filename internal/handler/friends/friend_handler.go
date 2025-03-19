package friends

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
	friendService service.FriendService
}

func NewHandler(friendService service.FriendService) *Handler {
	return &Handler{
		friendService: friendService,
	}
}

type FriendListResponse struct {
	Success    bool         `json:"success"`
	FriendList []orm.Friend `json:"friendList"`
}

func (h *Handler) GetFriendList(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticator.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	friendList, err := h.friendService.GetFriendList(r.Context(), userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(FriendListResponse{Success: true, FriendList: friendList})
}

type AddFriendRequest struct {
	PhoneNumber string `json:"phoneNumber"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
}

func (h *Handler) AddFriend(w http.ResponseWriter, r *http.Request) {

	userID, err := authenticator.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req AddFriendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.PhoneNumber == "" {
		http.Error(w, "Missing phone number", http.StatusBadRequest)
		return
	}

	err = h.friendService.AddFriend(r.Context(), userID, req.PhoneNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SuccessResponse{Success: true})
}

type BlockFriendRequest struct {
	FriendID string `json:"friendId"`
}

func (h *Handler) BlockFriend(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticator.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// URL에서 friendId 파라미터 추출
	vars := mux.Vars(r)
	friendID := vars["friendId"]

	if friendID == "" {
		http.Error(w, "Missing friend ID", http.StatusBadRequest)
		return
	}

	friendUUID, err := uuid.Parse(friendID)
	if err != nil {
		http.Error(w, "Invalid friend ID", http.StatusBadRequest)
		return
	}

	err = h.friendService.BlockFriend(r.Context(), userID, friendUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SuccessResponse{Success: true})
}

func (h *Handler) UnblockFriend(w http.ResponseWriter, r *http.Request) {
	userID, err := authenticator.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// URL에서 friendId 파라미터 추출
	vars := mux.Vars(r)
	friendID := vars["friendId"]

	if friendID == "" {
		http.Error(w, "Missing friend ID", http.StatusBadRequest)
		return
	}

	friendUUID, err := uuid.Parse(friendID)
	if err != nil {
		http.Error(w, "Invalid friend ID", http.StatusBadRequest)
		return
	}

	err = h.friendService.UnblockFriend(r.Context(), userID, friendUUID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SuccessResponse{Success: true})
}
