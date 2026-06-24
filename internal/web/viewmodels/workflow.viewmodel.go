package viewmodels

import (
	"github.com/MarcelArt/refinery/internal/entities"
)

type PaginationViewModel struct {
	Total      int64
	Page       int64
	Size       int64
	TotalPages int64
	Last       bool
	First      bool
	PrevPage   int64
	NextPage   int64
	Start      int64
	End        int64
}

type WorkflowRowViewModel struct {
	ID          uint
	Title       string
	Description string
	Prompt      string
	Schemas     []entities.WorkflowSchema
}

type WorkflowPageViewModel struct {
	Title      string
	User       entities.User
	Workflows  []WorkflowRowViewModel
	Pagination PaginationViewModel
}
