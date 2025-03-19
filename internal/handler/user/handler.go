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
	PhoneNumber string `json:"phoneNumber"`
	CountryCode string `json:"countryCode"`
}

type SuccessResponse struct {
	Success bool `json:"success"`
}

type UserResponse struct {
	Success bool        `json:"success"`
	User    interface{} `json:"user"`
}

type TokenResponse struct {
	Success bool   `json:"success"`
	Token   string `json:"token"`
}

type VerifiedResponse struct {
	Success  bool `json:"success"`
	Verified bool `json:"verified"`
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
	json.NewEncoder(w).Encode(UserResponse{Success: true, User: user})
}

type LoginRequest struct {
	PhoneNumber string `json:"phoneNumber"`
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
	json.NewEncoder(w).Encode(TokenResponse{Success: true, Token: token})
}

type AuthNumberRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	CountryCode string `json:"countryCode"`
	DeviceID    string `json:"deviceId"`
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(SuccessResponse{Success: true})
}

type CheckAuthRequest struct {
	PhoneNumber string `json:"phoneNumber"`
	CountryCode string `json:"countryCode"`
	DeviceID    string `json:"deviceId"`
	AuthNumber  string `json:"authNumber"`
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(VerifiedResponse{Success: true, Verified: true})
}
