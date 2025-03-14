package orm

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Room struct {
	UUIDv7BaseModel
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	RoomName  string         `gorm:"type:varchar(40);not null"`
}

type RoomUser struct {
	gorm.Model
	RoomID     uuid.UUID `gorm:"type:uuid; not null;"`
	UserID     uuid.UUID `gorm:"type:uuid;not null"`
	StartIndex int       `gorm:"type:integer;not null;default:0"`
}
