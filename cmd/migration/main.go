package main

import (
	"server/internal/models/orm"
	"server/internal/db/postgres_db"
	"github.com/joho/godotenv"

)

func main() {

	godotenv.Load(".env")



	postgres_db.GetPostgresClient().AutoMigrate(&orm.User{}) 
	postgres_db.GetPostgresClient().AutoMigrate(&orm.Friend{})
	postgres_db.GetPostgresClient().AutoMigrate(&orm.Room{})
	postgres_db.GetPostgresClient().AutoMigrate(&orm.RoomUser{})
	postgres_db.GetPostgresClient().AutoMigrate(&orm.AuthenticateMessage{})
}
