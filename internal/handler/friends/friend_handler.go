package friends

import (
	"encoding/json"
	"net/http"
	"server/internal/models/orm"
	"server/internal/service"
	"server/pkg/authenticator"

	"github.com/google/uuid"
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
	FriendList []orm.Friend `json:"friend_list"`
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
	json.NewEncoder(w).Encode(FriendListResponse{FriendList: friendList})
}

type AddFriendRequest struct {
	PhoneNumber string `json:"phone_number"`
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

	w.WriteHeader(http.StatusOK)
}

type BlockFriendRequest struct {
	FriendID string `json:"friend_id"`
}

func (h *Handler) BlockFriend(w http.ResponseWriter, r *http.Request) {

	userID, err := authenticator.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req BlockFriendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.FriendID == "" {
		http.Error(w, "Missing friend ID", http.StatusBadRequest)
		return
	}

	friendID, err := uuid.Parse(req.FriendID)
	if err != nil {
		http.Error(w, "Invalid friend ID", http.StatusBadRequest)
		return
	}

	err = h.friendService.BlockFriend(r.Context(), userID, friendID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handler) UnblockFriend(w http.ResponseWriter, r *http.Request) {

	userID, err := authenticator.GetUserID(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req BlockFriendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.FriendID == "" {
		http.Error(w, "Missing friend ID", http.StatusBadRequest)
		return
	}

	friendID, err := uuid.Parse(req.FriendID)
	if err != nil {
		http.Error(w, "Invalid friend ID", http.StatusBadRequest)
		return
	}

	err = h.friendService.UnblockFriend(r.Context(), userID, friendID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
