package authenticator

import (
	"errors"
	"net/http"
	"server/pkg/log"

	"github.com/google/uuid"
)

func GetUserID(r *http.Request) (string, error) {
	logger := log.NewColorLog()

	userID, ok := r.Context().Value(ContextKeyUserID).(string)
	if !ok {
		logger.Error("Could not get user_id from context")
		return "", errors.New("could not get user_id from context")
	}
	if userID == "" {
		logger.Error("User_id is empty")
		return "", errors.New("user_id is empty")
	}
	if !IsValidUUID(userID) {
		logger.Error("User_id is invalid UUID", "user_id", userID)
		return "", errors.New("user_id is invalid UUID")
	}
	return userID, nil
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
