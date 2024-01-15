package orm

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserName     string         `gorm:"type:varchar(40);not null"`
	PasswordHash string         `gorm:"type:varchar(255);not null"`
	ProfileImage sql.NullString `gorm:"type:varchar(255)"`
	PhoneNumber  string         `gorm:"type:varchar(24);not null"`
	CountryCode  string         `gorm:"type:varchar(8);not null"`
	Email        sql.NullString `gorm:"type:varchar(64)"`
	OauthID      sql.NullString `gorm:"type:varchar(255)"`
}

type DeviceList struct {
	gorm.Model
	DeviceID uuid.UUID `gorm:"primaryKey; type:uuid; not null; default:gen_random_uuid()"`
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
	gorm.Model
	RoomID   uuid.UUID `gorm:"primaryKey; type:uuid; not null; default:gen_random_uuid()"`
	RoomName string    `gorm:"type:varchar(40);not null"`
}

type RoomUser struct {
	gorm.Model
	RoomID     uuid.UUID `gorm:"type:uuid; not null; default:gen_random_uuid()"`
	UserID     uuid.UUID `gorm:"type:uuid;not null"`
	StartIndex int       `gorm:"type:integer;not null;default:0"`
}

type AuthenticateMessage struct {
	ID                 uint      `gorm:"primarykey"`
	CountryCode        string    `gorm:"type:varchar(8);not null"`
	PhoneNumber        string    `gorm:"type:varchar(24);not null"`
	RequestTime        time.Time `gorm:"type:timestamp;not null;default:now()"`
	ExpireTime         time.Time `gorm:"type:timestamp;not null;default:now() + interval '3 minutes'"`
	DeviceID           string    `gorm:"type:varchar(24); not null"`
	AuthenticateNumber string    `gorm:"type:varchar(8);not null"`
	Trial              int       `gorm:"type:integer;not null;default:0"`
}
