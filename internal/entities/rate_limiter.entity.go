package entities

import "gorm.io/gorm"

type RateLimiter struct {
	gorm.Model
	Count uint `gorm:"not null;default:1" json:"count"`

	UserID uint `gorm:"not null" json:"userId"`

	User *User `json:"user,omitzero"`
}
