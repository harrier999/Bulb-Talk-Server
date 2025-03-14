package orm

import (
	"time"
)

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
