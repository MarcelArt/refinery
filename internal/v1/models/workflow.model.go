package models

import (
	"github.com/MarcelArt/refinery/internal/common"
	"github.com/MarcelArt/refinery/internal/entities"
	"github.com/MarcelArt/refinery/pkg/jsonb"
)

type WorkflowInput struct {
	common.InputModel
	Title       string                                 `gorm:"not null" json:"title"`
	Description string                                 `gorm:"not null" json:"description"`
	Prompt      string                                 `json:"prompt"`
	Schemas     jsonb.JSONB[[]entities.WorkflowSchema] `json:"schemas"`
	Type        string                                 `gorm:"default:PDF_TEXT" json:"type"` // PDF_TEXT, PICTURE

	UserID uint `json:"userId"`
}

type WorkflowPage struct {
	ID          uint                                   `json:"ID"`
	Title       string                                 `json:"title"`
	Description string                                 `json:"description"`
	Prompt      string                                 `json:"prompt"`
	Schemas     jsonb.JSONB[[]entities.WorkflowSchema] `json:"schemas"`
	Type        string                                 `gorm:"default:PDF_TEXT" json:"type"` // PDF_TEXT, PICTURE
}

type WorkflowStartOption struct {
	AdditionalPrompt string
	Metadata         string
}
