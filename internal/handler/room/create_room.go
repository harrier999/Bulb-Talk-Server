package room

import (
	"encoding/json"
	"errors"
	"net/http"
	"server/internal/db/postgres_db"
	"server/internal/models/orm"
	"server/pkg/authenticator"
	"server/pkg/log"

	"github.com/google/uuid"
)

type createRoomRequest struct {
	Users []uuid.UUID `json:"room_user_list"`
}
type createRoomResponse struct {
	RoomID uuid.UUID `json:"room_id"`
}

func CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	logger := log.NewColorLog()
	var req createRoomRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := validCreateRoomRequest(req); err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	user_id, err := authenticator.GetUserID(r)
	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	req.Users = append(req.Users, user_id)
	
	room_id, err := createRoom(req.Users);

	if err != nil {
		logger.Error(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(createRoomResponse{RoomID: room_id})	
}

func createRoom(room_user_list []uuid.UUID) (uuid.UUID, error) {
	postgresClient := postgres_db.GetPostgresClient()
	room := orm.Room{RoomName: "group"} // TODO: get room name from users' name
	tx := postgresClient.Create(&room)
	if tx.Error != nil {
		return uuid.Nil, tx.Error
	}
	room_id := room.RoomID
	room_user_objects := make([]orm.RoomUser, len(room_user_list))
	for i, user_id := range room_user_list {
		room_user_objects[i] = orm.RoomUser{RoomID: room_id, UserID: user_id}
	}
	tx = postgresClient.Create(&room_user_objects)
	if tx.Error != nil {
		return uuid.Nil, tx.Error
	}
	return room_id, nil
}

func validCreateRoomRequest(req createRoomRequest) error {
	if len(req.Users) < 1 {
		return errors.New("not enough user")
	}
	for _, user := range req.Users {
		if _, err := uuid.Parse(user.String()); err != nil {
			return errors.New("invalid user uuid")
		}
	}
	return nil
}
