package models

import (
	"time"

	"github.com/MarcelArt/refinery/internal/common"
	"github.com/MarcelArt/refinery/internal/entities"
)

// UserInput
type UserInput struct {
	common.InputModel
	Username   string     `gorm:"not null;unique" json:"username"`
	Email      string     `gorm:"not null;unique" json:"email"`
	Password   string     `gorm:"not null" json:"password" mapper:"password"`
	VerifiedAt *time.Time `json:"-" mapper:"verifiedAt"`
	DailyLimit uint       `gorm:"default:10" json:"dailyLimit"`
}

// UserInput end

type UserPage struct {
	ID         uint       `json:"ID"`
	Username   string     `gorm:"not null;unique" json:"username"`
	Email      string     `gorm:"not null;unique" json:"email"`
	VerifiedAt *time.Time `json:"verifiedAt"`
	DailyLimit uint       `gorm:"default:10" json:"dailyLimit"`
	// Roles    jsonb.JSONB[[]string] `json:"roles"`
}

type LoginInput struct {
	Username   string `json:"username"`
	Password   string `json:"password"`
	IsRemember bool   `json:"isRemember"`
}

type LoginResponse struct {
	AccessToken  string        `json:"accessToken"`
	RefreshToken string        `json:"refreshToken"`
	User         entities.User `json:"user"`
}

type UserRole struct {
	ID          uint   `json:"ID"`
	Name        string `json:"name"`
	Description string `json:"description"`
}
