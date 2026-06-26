package entities

import (
	"github.com/MarcelArt/refinery/pkg/jsonb"
	"gorm.io/gorm"
)

type ApiKey struct {
	gorm.Model
	Name   string                `gorm:"not null" json:"name"`
	Key    string                `gorm:"not null;unique" json:"key"`
	Scopes jsonb.JSONB[[]string] `json:"scopes"`

	UserID uint `gorm:"not null" json:"userId"`

	User *User `json:"user,omitempty"`
}
