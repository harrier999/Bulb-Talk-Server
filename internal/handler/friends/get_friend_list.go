package friends

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"server/internal/db/postgres_db"
	"server/internal/models/orm"
	"server/pkg/authenticator"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FriendListRequest struct {
	LastRequestTime time.Time `json:"last_request_time"`
}

type FriendListResponse struct {
	FriendList []orm.Friend `json:"friend_list"`
}

func GetFriendList(w http.ResponseWriter, r *http.Request) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	user_id, ok := r.Context().Value(authenticator.ContextKeyUserID).(string)
	if !ok {
		logger.Warn("Error while parsing user_id", "user_id", user_id)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !IsValidUUID(user_id) {
		logger.Warn("User id is invalid UUID", "user_id", user_id)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	friendList, err := getFriendListDB(user_id)
	if err != nil {
		logger.Warn("Error while getting friend list", "user_id", user_id)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	friendListResponse := FriendListResponse{FriendList: friendList}
	if err := json.NewEncoder(w).Encode(friendListResponse); err != nil {
		logger.Warn("Error while encoding friend list", "user_id", user_id)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func IsValidUUID(u string) bool {
	_, err := uuid.Parse(u)
	return err == nil
}

func getFriendListDB(user_id string) ([]orm.Friend, error) {
	postgresCleint := postgres_db.GetPostgresClient()
	var friends []orm.Friend
	if err := postgresCleint.Where("user_id = ?", user_id).Find(&friends).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return friends, nil
		}
		return nil, err
	}

	return friends, nil
}
