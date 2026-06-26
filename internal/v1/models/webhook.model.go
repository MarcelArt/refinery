package models

import (
	"github.com/MarcelArt/refinery/internal/common"
)

type WebhookInput struct {
	common.InputModel
	URL        string `json:"url"`
	Method     string `json:"method"`
	WorkflowID uint   `json:"workflowId"`
}

type WebhookPage struct {
	ID         uint   `json:"ID"`
	URL        string `json:"url"`
	Method     string `json:"method"`
	WorkflowID uint   `json:"workflowId"`
}
