package viewmodels

import (
	"github.com/MarcelArt/refinery/internal/entities"
)

type WorkflowRowViewModel struct {
	ID          uint
	Title       string
	Description string
	Prompt      string
	Schemas     []entities.WorkflowSchema
}

type WorkflowPageViewModel struct {
	Title     string
	User      entities.User
	Workflows []WorkflowRowViewModel
}
