package orm

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	UUIDv7BaseModel
	UserName     string         `gorm:"type:varchar(40);not null"`
	PasswordHash string         `gorm:"type:varchar(255);not null"`
	ProfileImage sql.NullString `gorm:"type:varchar(255)"`
	PhoneNumber  string         `gorm:"type:varchar(24);not null"`
	CountryCode  string         `gorm:"type:varchar(8);not null"`
	Email        sql.NullString `gorm:"type:varchar(64)"`
	OauthID      sql.NullString `gorm:"type:varchar(255)"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

type DeviceList struct {
	gorm.Model
	UUIDv7BaseModel
	UserID uuid.UUID `gorm:"type:uuid;not null"`
	Alarm  bool      `gorm:"type:boolean;not null;default:true"`
}
