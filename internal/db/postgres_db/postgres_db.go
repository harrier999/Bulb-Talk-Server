package postgres_db

import (
	"fmt"
	"log/slog"
	"os"
	"server/pkg/log"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DBInstance struct {
	Client *gorm.DB
	logger *slog.Logger
}

var (
	instance          *DBInstance
	testInstance      *DBInstance
	instanceMutex     sync.Mutex
	testInstanceMutex sync.Mutex
)

func GetPostgresClient() *gorm.DB {
	if instance == nil {
		instanceMutex.Lock()
		defer instanceMutex.Unlock()

		if instance == nil {
			instance = &DBInstance{
				logger: log.NewColorLog(),
			}
			instance.connect("chat")
		}
	}

	return instance.Client
}

func GetTestPostgresClient() *gorm.DB {
	if testInstance == nil {
		testInstanceMutex.Lock()
		defer testInstanceMutex.Unlock()

		if testInstance == nil {
			testInstance = &DBInstance{
				logger: log.NewColorLog(),
			}
			testInstance.connect("chat_test")
		}
	}

	return testInstance.Client
}

func (db *DBInstance) connect(dbName string) {
	db.logger.Info(fmt.Sprintf("Connecting to PostgreSQL database: %s...", dbName))

	databaseAddr := os.Getenv("POSTGRES_ADDR")
	databasePassword := os.Getenv("POSTGRES_PASS")
	dbPort := os.Getenv("POSTGRES_PORT")
	dbUser := os.Getenv("POSTGRES_USER")

	connectionString := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s",
		databaseAddr, dbUser, databasePassword, dbName, dbPort,
	)

	client, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		db.logger.Error("Error connecting to PostgreSQL", "error", err.Error())
		os.Exit(1)
	}

	db.Client = client
	db.logger.Info(fmt.Sprintf("Successfully connected to PostgreSQL database: %s", dbName))
}

func CloseConnection() error {
	if instance != nil && instance.Client != nil {
		sqlDB, err := instance.Client.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}

func CloseTestConnection() error {
	if testInstance != nil && testInstance.Client != nil {
		sqlDB, err := testInstance.Client.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
