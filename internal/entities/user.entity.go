package entities

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username   string     `gorm:"not null;unique" json:"username"`
	Email      string     `gorm:"not null;unique" json:"email"`
	Password   string     `json:"-" mapper:"password"`
	VerifiedAt *time.Time `json:"verifiedAt" mapper:"verifiedAt"`
}
