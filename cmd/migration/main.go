package main

import (
	"github.com/joho/godotenv"
	"server/internal/db/postgres_db"
	"server/internal/models/orm"
)

func main() {

	godotenv.Load(".env")

	postgres_db.GetPostgresClient().AutoMigrate(&orm.User{})
	postgres_db.GetPostgresClient().AutoMigrate(&orm.Friend{})
	postgres_db.GetPostgresClient().AutoMigrate(&orm.Room{})
	postgres_db.GetPostgresClient().AutoMigrate(&orm.RoomUser{})
	postgres_db.GetPostgresClient().AutoMigrate(&orm.AuthenticateMessage{})
}
