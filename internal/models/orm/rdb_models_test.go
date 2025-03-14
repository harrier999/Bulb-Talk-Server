package orm_test

import (
	"server/internal/db/postgres_db"
	"server/internal/models/orm"
	"testing"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)



func TestMigration(t *testing.T) {
	err := godotenv.Load("../../../cmd/talk-server/.env")
	if err != nil {
		t.Errorf("Error loading .env file")
	}
	db := postgres_db.GetPostgresClient()
	if db == nil {
		t.Errorf("Expected db to be not nil")
	}
	db.AutoMigrate(&orm.Room{})
	room1 := orm.Room{}
	room1.RoomName = "test_room1"
	room2 := orm.Room{}
	room2.RoomName = "test_room2"
	db.Create(&room1)
	db.Create(&room2)

	room3 := orm.Room{}
	room3.ID = room1.ID
	room4 := orm.Room{}
	room4.ID = room2.ID

	db.First(&room3)
	db.First(&room4)

	defer db.Delete(&room1)
	defer db.Delete(&room2)

	assert.Equal(t, room1, room3)
	assert.Equal(t, room2, room4)

	// roomList := []orm.RoomList{}
	// db.Find(&roomList).Scan(&roomList)
	// for _, room := range roomList {
	// 	t.Logf("Room: %s", room.RoomID.String())
	// 	fmt.Printf("Room: %s\n", room.RoomID.String())
	// }

}
