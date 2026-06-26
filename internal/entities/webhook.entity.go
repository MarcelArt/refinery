package entities

import (
	"gorm.io/gorm"
)

type Webhook struct {
	gorm.Model
	URL        string    `gorm:"not null" json:"url"`
	Method     string    `gorm:"not null" json:"method"`
	WorkflowID uint      `gorm:"not null" json:"workflowId"`
	Workflow   *Workflow `json:"workflow,omitempty"`
}
