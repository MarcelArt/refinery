package entities

import (
	"github.com/MarcelArt/refinery/pkg/jsonb"
	"gorm.io/gorm"
)

type ExtractionJSON []map[string]any
type ExtractionResult struct {
	gorm.Model
	Raw        string                      `gorm:"not null" json:"raw"`
	Json       jsonb.JSONB[ExtractionJSON] `json:"json"`
	WorkflowID uint                        `gorm:"not null" json:"workflow_id"`
	Workflow   *Workflow                   `json:"workflow,omitempty"`
}
