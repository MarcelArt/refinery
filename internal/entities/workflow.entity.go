package entities

import (
	"github.com/MarcelArt/refinery/pkg/jsonb"
	"gorm.io/gorm"
)

type Workflow struct {
	gorm.Model
	Title       string                        `gorm:"not null" json:"title"`
	Description string                        `gorm:"not null" json:"description"`
	Schemas     jsonb.JSONB[[]WorkflowSchema] `json:"schemas"`
}

type WorkflowSchema struct {
	Key         string `json:"key"`
	Type        string `json:"type"`
	Description string `json:"description"`
	Example     string `json:"example"`
}
