package orm

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UUIDv7BaseModel struct {
	ID uuid.UUID `gorm:"primaryKey; type:uuid; not null;"`
}

func (u *UUIDv7BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	uuidV7, err := uuid.NewV7();
	if err != nil {
		return err
	}
	u.ID = uuidV7
	return nil
}

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
	UserID   uuid.UUID `gorm:"type:uuid;not null"`
	Alarm    bool      `gorm:"type:boolean;not null;default:true"`
}

type Friend struct {
	gorm.Model
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	FriendID  uuid.UUID `gorm:"type:uuid;not null"`
	IsBlocked bool      `gorm:"type:boolean;not null;default:false"`
}

type Room struct {
	UUIDv7BaseModel
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
	RoomName string    `gorm:"type:varchar(40);not null"`
}

type RoomUser struct {
	gorm.Model
	RoomID     uuid.UUID `gorm:"type:uuid; not null;"`
	UserID     uuid.UUID `gorm:"type:uuid;not null"`
	StartIndex int       `gorm:"type:integer;not null;default:0"`
}

type AuthenticateMessage struct {
	UUIDv7BaseModel
	CountryCode        string    `gorm:"type:varchar(8);not null"`
	PhoneNumber        string    `gorm:"type:varchar(24);not null"`
	RequestTime        time.Time `gorm:"type:timestamp;not null;default:now()"`
	ExpireTime         time.Time `gorm:"type:timestamp;not null;default:now() + interval '3 minutes'"`
	DeviceID           string    `gorm:"type:varchar(24); not null"`
	AuthenticateNumber string    `gorm:"type:varchar(8);not null"`
	Trial              int       `gorm:"type:integer;not null;default:0"`
}
