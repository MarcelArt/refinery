package models

import (
	"github.com/MarcelArt/refinery/internal/common"
)

type WebhookInput struct {
	common.InputModel
	DisplayName string `gorm:"not null" json:"displayName"`
	URL         string `json:"url"`
	Method      string `json:"method"`
	HmacKey     string `gorm:"not null" json:"hmacKey"`
	WorkflowID  uint   `json:"workflowId"`
}

type WebhookPage struct {
	ID          uint   `json:"ID"`
	DisplayName string `gorm:"not null" json:"displayName"`
	URL         string `json:"url"`
	Method      string `json:"method"`
	WorkflowID  uint   `json:"workflowId"`
}
