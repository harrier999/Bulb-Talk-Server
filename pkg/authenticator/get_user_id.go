package authenticator

import (
	"errors"
	"net/http"
	"server/pkg/log"

	"github.com/google/uuid"
)

func GetUserID(r *http.Request) (uuid.UUID, error) {
	logger := log.NewColorLog()

	userID, ok := r.Context().Value(ContextKeyUserID).(string)
	if !ok {
		logger.Error("Could not get user_id from context")
		return uuid.Nil, errors.New("could not get user_id from context")
	}
	if userID == "" {
		logger.Error("User_id is empty")
		return uuid.Nil, errors.New("user_id is empty")
	}
	if !IsValidUUID(userID) {
		logger.Error("User_id is invalid UUID", "user_id", userID)
		return uuid.Nil, errors.New("user_id is invalid UUID")
	}
	return uuid.Parse(userID)
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}
