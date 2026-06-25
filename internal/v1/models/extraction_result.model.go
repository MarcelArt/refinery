package models

import (
	"time"

	"github.com/MarcelArt/refinery/internal/common"
	"github.com/MarcelArt/refinery/internal/entities"
	"github.com/MarcelArt/refinery/pkg/jsonb"
)

type ExtractionResultInput struct {
	common.InputModel
	Raw        string                               `gorm:"not null" json:"raw"`
	Json       jsonb.JSONB[entities.ExtractionJSON] `json:"json"`
	Source     string                               `json:"source"`
	Status     string                               `gorm:"default:IN_PROGRESS" json:"string"` // IN_PROGRESS, DONE, FAILED
	FinishedAt *time.Time                           `json:"finishedAt"`
	WorkflowID uint                                 `gorm:"not null" json:"workflowId"`
}

type ExtractionResultPage struct {
	ID         uint                                 `json:"ID"`
	CreatedAt  time.Time                            `json:"CreatedAt"`
	Raw        string                               `json:"raw"`
	Json       jsonb.JSONB[entities.ExtractionJSON] `json:"json"`
	Source     string                               `json:"source"`
	Status     string                               `gorm:"default:IN_PROGRESS" json:"string"` // IN_PROGRESS, DONE, FAILED
	FinishedAt time.Time                            `json:"finishedAt"`
	WorkflowID uint                                 `json:"workflowId"`
}

type ContentLLM struct {
	Content string `json:"content"`
	Source  string `json:"source"`
}
