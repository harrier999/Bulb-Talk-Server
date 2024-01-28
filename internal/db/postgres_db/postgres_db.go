package postgres_db

import (
	"fmt"
	"os"
	"server/pkg/log"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	postgresClient *gorm.DB
	once           sync.Once
	logger = log.NewColorLog()
)

func GetPostgresClient() *gorm.DB {
	once.Do(func() {
		logger.Info("Connecting to postgres...")

		database_addr := os.Getenv("POSTGRES_ADDR")
		database_password := os.Getenv("POSTGRES_PASS")
		db_port := os.Getenv("POSTGRES_PORT")
		db_user := os.Getenv("POSTGRES_USER")

		connection_string := fmt.Sprintf("host=%s user=%s password=%s dbname=chat port=%s", database_addr, db_user, database_password, db_port)
		db, err := gorm.Open(postgres.Open(connection_string), &gorm.Config{})
		if err != nil {
			logger.Error("Error connecting to postgres. Error: ", "error", err.Error())
			os.Exit(1)
		}
		postgresClient = db
	})

	return postgresClient
}

func GetTestPostgresCleint() *gorm.DB {
	once.Do(func() {
		logger.Info("Connecting to test postgres ...", )

		database_addr := os.Getenv("POSTGRES_ADDR")
		database_password := os.Getenv("POSTGRES_PASS")
		db_port := os.Getenv("POSTGRES_PORT")
		db_user := os.Getenv("POSTGRES_USER")

		connection_string := fmt.Sprintf("host=%s user=%s password=%s dbname=chat_test port=%s", database_addr, db_user, database_password, db_port)

		db, err := gorm.Open(postgres.Open(connection_string), &gorm.Config{})
		if err != nil {
			logger.Error("Error connecting to postgres. Error: ", "error", err.Error())
			os.Exit(1)
		}
		postgresClient = db
	})

	return postgresClient
}
