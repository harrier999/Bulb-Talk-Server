package postgres_db

import (
	"log"
	"sync"
	"os"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var postgresClient *gorm.DB
var once sync.Once


func GetPostgresClient() *gorm.DB {
	once.Do(func() {
		log.Println("Connecting to postgres...")
		
		database_addr := os.Getenv("POSTGRES_ADDR")
		database_password := os.Getenv("POSTGRES_PASS")
		db_port := os.Getenv("POSTGRES_PORT")
		//db_ssl := os.Getenv("POSTGRES_SSL")
		db_user := os.Getenv("POSTGRES_USER")

		connection_string := fmt.Sprintf("host=%s user=%s password=%s dbname=chat port=%s", database_addr, db_user, database_password, db_port)
		fmt.Println(connection_string)
		db, err := gorm.Open(postgres.Open(connection_string), &gorm.Config{})
		if err != nil {
			log.Fatal("Error connecting to postgres. Error: ", err.Error())
		}
		postgresClient = db
	})

	return postgresClient
}
