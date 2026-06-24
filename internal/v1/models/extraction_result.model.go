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
	WorkflowID uint                                 `gorm:"not null" json:"workflow_id"`
}

type ExtractionResultPage struct {
	ID         uint                                 `json:"ID"`
	CreatedAt  time.Time                            `json:"CreatedAt"`
	Raw        string                               `json:"raw"`
	Json       jsonb.JSONB[entities.ExtractionJSON] `json:"json"`
	WorkflowID uint                                 `json:"workflow_id"`
}

type ContentLLM struct {
	Content string `json:"content"`
}
