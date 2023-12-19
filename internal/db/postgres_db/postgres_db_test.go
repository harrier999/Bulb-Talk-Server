package postgres_db_test

import (
	"github.com/joho/godotenv"
	"server/internal/db/postgres_db"
	"testing"
	"os"
	"fmt"
)

func TestOne(t *testing.T) {
	fmt.Println(os.Getenv("FULLNODE_IMAGE"))
	err := godotenv.Load("../../../cmd/talk-server/.env")
	if err != nil {
		t.Errorf("Error loading .env file")
	}

	db := postgres_db.GetPostgresClient()
	if db == nil {
		t.Errorf("Expected db to be not nil")
	}
	db.Select("SELECT * FROM users")
}
