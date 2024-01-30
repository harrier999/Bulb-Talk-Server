package room

import (
	"encoding/json"
	"net/http"
	"server/internal/db/postgres_db"
	"server/pkg/authenticator"
	"server/pkg/log"

	"github.com/google/uuid"
)

type GetRoomListResponse struct {
	RoomID   uuid.UUID `json:"room_id"`
	RoomName string    `json:"room_name"`
}

func GetRoomListHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.NewColorLog()
	user_id, err := authenticator.GetUserID(r)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	postgresClient := postgres_db.GetPostgresClient()
	var roomList []GetRoomListResponse
	postgresClient.Table("rooms").Select("rooms.room_id, rooms.room_name").
		Joins("left join room_users on rooms.room_id = room_users.room_id").
		Where("room_users.user_id = ?", user_id).Scan(&roomList)

	err = json.NewEncoder(w).Encode(roomList)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}