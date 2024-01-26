package friends

import (
	"encoding/json"
	"net/http"
	"server/internal/db/postgres_db"
	"server/internal/models/orm"
	"server/pkg/authenticator"
	"server/pkg/log"
	"time"

	"gorm.io/gorm"
)

type FriendListRequest struct {
	LastRequestTime time.Time `json:"last_request_time"`
}

type FriendListResponse struct {
	FriendList []orm.Friend `json:"friend_list"`
}

func GetFriendList(w http.ResponseWriter, r *http.Request) {
	logger := log.NewColorLog()

	user_id, err := authenticator.GetUserID(r)
	if err != nil {
		logger.Warn("Error while getting user_id from context")
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
