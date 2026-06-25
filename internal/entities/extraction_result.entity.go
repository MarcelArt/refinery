package entities

import (
	"time"

	"github.com/MarcelArt/refinery/pkg/jsonb"
	"gorm.io/gorm"
)

type ExtractionJSON []map[string]any
type ExtractionResult struct {
	gorm.Model
	Raw        string                      `gorm:"not null" json:"raw"`
	Json       jsonb.JSONB[ExtractionJSON] `json:"json"`
	Source     string                      `json:"source"`
	Status     string                      `gorm:"default:IN_PROGRESS" json:"string"` // IN_PROGRESS, DONE, FAILED
	FinishedAt *time.Time                  `json:"finishedAt"`
	WorkflowID uint                        `gorm:"not null" json:"workflowId"`
	Workflow   *Workflow                   `json:"workflow,omitempty"`
}
