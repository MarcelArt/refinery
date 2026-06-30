package models

import (
	"github.com/MarcelArt/refinery/internal/common"
)

type RateLimiterInput struct {
	common.InputModel
	Count  uint `json:"count"`
	UserID uint `json:"userId"`
}

type RateLimiterPage struct {
	ID     uint `json:"ID"`
	Count  uint `json:"count"`
	UserID uint `json:"userId"`
}
