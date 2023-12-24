package room

import (
	"log"
	"net/http"
	"server/internal/db/postgres_db"
	"server/internal/models/orm"
	"github.com/google/uuid"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	postgresClient := postgres_db.GetPostgresClient()

	user, err := GetUserId(r.Header)
	if err != nil {
		log.Println("Error getting user id. Error: ", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	roomUserList := []orm.RoomUserList{}
	roomUser := orm.RoomUserList{}
	roomUser.UserID = uuid.MustParse(user)

	postgresClient.Model(&orm.RoomUserList{}).Where(&roomUser).Find(&roomUserList)

	result := encoder(roomUserList)
	w.WriteHeader(http.StatusOK)
	w.Write(result)
	
}

func encoder(roomUserList []orm.RoomUserList) []byte {
	// TODO: implement encoder
	return []byte{}
}

func GetUserId(header http.Header) (string, error) {
	// TODO: implement jwt authorization
	return header.Get("user_id"), nil
}

