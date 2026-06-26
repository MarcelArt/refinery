package models

import (
	"github.com/MarcelArt/refinery/internal/common"
	"github.com/MarcelArt/refinery/pkg/jsonb"
)

type ApiKeyInput struct {
	common.InputModel
	Name   string                `gorm:"not null" json:"name"`
	Key    string                `gorm:"not null;unique" json:"key"`
	Scopes jsonb.JSONB[[]string] `json:"scopes"`

	UserID uint `gorm:"not null" json:"userId"`
}

type ApiKeyPage struct {
	ID     uint                  `json:"ID"`
	Name   string                `gorm:"not null" json:"name"`
	Scopes jsonb.JSONB[[]string] `json:"scopes"`

	UserID uint `gorm:"not null" json:"userId"`
}
