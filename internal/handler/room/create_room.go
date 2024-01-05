package room

import (
	"gorm.io/gorm"
	"fmt"
	"log"
	"server/internal/db/postgres_db"
)

func CreateRoom(roomName string, user string) {
	client := postgres_db.GetPostgresClient()
	room := RoomList{}
	room.RoomName = "test_room"
	
}