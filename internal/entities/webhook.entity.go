package entities

import (
	"gorm.io/gorm"
)

type Webhook struct {
	gorm.Model
	DisplayName string    `gorm:"not null" json:"displayName"`
	URL         string    `gorm:"not null" json:"url"`
	Method      string    `gorm:"not null" json:"method"`
	WorkflowID  uint      `gorm:"not null" json:"workflowId"`
	HmacKey     string    `gorm:"not null" json:"hmacKey"`
	Workflow    *Workflow `json:"workflow,omitempty"`
}
