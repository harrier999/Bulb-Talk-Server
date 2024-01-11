package orm

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	UserID       uuid.UUID      `gorm:"primaryKey; type:uuid; not null; default:gen_random_uuid()"`
	UserName     string         `gorm:"type:varchar(40);not null"`
	PasswordHash string         `gorm:"type:varchar(255);not null"`
	Salt         string         `gorm:"type:varchar(255);not null"`
	ProfileImage sql.NullString `gorm:"type:varchar(255)"`
	PhoneNumber  string         `gorm:"type:varchar(24);not null"`
	CountryCode  string         `gorm:"type:varchar(8);not null"`
	Email        sql.NullString `gorm:"type:varchar(64)"`
	OauthID      sql.NullString `gorm:"type:varchar(255)"`
	CreatedAt    sql.NullTime   `gorm:"type:timestamp;not null;default:now()"`
}

type DeviceList struct {
	DeviceID uuid.UUID `gorm:"primaryKey; type:uuid; not null; default:gen_random_uuid()"`
	UserID   uuid.UUID `gorm:"type:uuid;not null"`
	Alarm    bool      `gorm:"type:boolean;not null;default:true"`
}

type FriendList struct {
	UserID    uuid.UUID `gorm:"type:uuid;not null"`
	FriendID  uuid.UUID `gorm:"type:uuid;not null"`
	IsBlocked bool      `gorm:"type:boolean;not null;default:false"`
}

type RoomList struct {
	RoomID   uuid.UUID `gorm:"primaryKey; type:uuid; not null; default:gen_random_uuid()"`
	RoomName string    `gorm:"type:varchar(40);not null"`
}

type RoomUserList struct {
	ID     uint      `gorm:"primarykey"`
	RoomID uuid.UUID `gorm:"type:uuid; not null; default:gen_random_uuid()"`
	UserID uuid.UUID `gorm:"type:uuid;not null"`
}

type AuthenticateMessage struct {
	ID                 uint      `gorm:"primarykey"`
	PhoneNumber        string    `gorm:"type:varchar(24);not null"`
	RequestTime        time.Time `gorm:"type:timestamp;not null;default:now()"`
	DeviceID           string    `gorm:"type:varchar(24); not null"`
	AuthenticateNumber string    `gorm:"type:varchar(8);not null"`
}
