package user

import (
	"encoding/json"
	"net/http"
	"server/internal/service"
)

type Handler struct {
	userService service.UserService
	authService service.AuthService
}

func NewHandler(userService service.UserService, authService service.AuthService) *Handler {
	return &Handler{
		userService: userService,
		authService: authService,
	}
}

type SignUpRequest struct {
	UserName    string `json:"username"`
	Password    string `json:"password"`
	PhoneNumber string `json:"phone_number"`
	CountryCode string `json:"country_code"`
}

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var req SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserName == "" || req.Password == "" || req.PhoneNumber == "" || req.CountryCode == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	user, err := h.userService.Register(r.Context(), req.UserName, req.Password, req.PhoneNumber, req.CountryCode)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

type LoginRequest struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.PhoneNumber == "" || req.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	token, err := h.userService.Login(r.Context(), req.PhoneNumber, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(token)
}

type AuthNumberRequest struct {
	PhoneNumber string `json:"phone_number"`
	CountryCode string `json:"country_code"`
	DeviceID    string `json:"device_id"`
}

func (h *Handler) RequestAuthNumber(w http.ResponseWriter, r *http.Request) {
	var req AuthNumberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.PhoneNumber == "" || req.CountryCode == "" || req.DeviceID == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	err := h.authService.RequestAuthNumber(r.Context(), req.PhoneNumber, req.CountryCode, req.DeviceID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

type CheckAuthRequest struct {
	PhoneNumber string `json:"phone_number"`
	CountryCode string `json:"country_code"`
	DeviceID    string `json:"device_id"`
	AuthNumber  string `json:"auth_number"`
}

func (h *Handler) CheckAuthNumber(w http.ResponseWriter, r *http.Request) {
	var req CheckAuthRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.PhoneNumber == "" || req.CountryCode == "" || req.DeviceID == "" || req.AuthNumber == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	valid, err := h.authService.CheckAuthNumber(r.Context(), req.PhoneNumber, req.CountryCode, req.DeviceID, req.AuthNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if !valid {
		http.Error(w, "Invalid authentication number", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
}
