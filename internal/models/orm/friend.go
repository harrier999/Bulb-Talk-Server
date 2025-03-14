package orm

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Friend struct {
	gorm.Model
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	FriendID  uuid.UUID `gorm:"type:uuid;not null"`
	IsBlocked bool      `gorm:"type:boolean;not null;default:false"`
}
