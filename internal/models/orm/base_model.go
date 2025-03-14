package orm

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UUIDv7BaseModel struct {
	ID uuid.UUID `gorm:"primaryKey; type:uuid; not null;"`
}

func (u *UUIDv7BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	uuidV7, err := uuid.NewV7()
	if err != nil {
		return err
	}
	u.ID = uuidV7
	return nil
}
