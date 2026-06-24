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
	Schemas     jsonb.JSONB[[]entities.WorkflowSchema] `json:"schemas"`
}

type WorkflowPage struct {
	ID          uint                                   `json:"ID"`
	Title       string                                 `json:"title"`
	Description string                                 `json:"description"`
	Schemas     jsonb.JSONB[[]entities.WorkflowSchema] `json:"schemas"`
}
