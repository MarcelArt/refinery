package viewmodels

import (
	"time"
)

type ExtractionResultRowViewModel struct {
	ID        uint
	CreatedAt time.Time
}

type ExtractionResultDetailsViewModel struct {
	ID        uint
	CreatedAt time.Time
	Columns   []string
	Rows      []map[string]any
}

type ExtractionResultPageViewModel struct {
	Title          string
	User           interface{}
	WorkflowID     uint
	WorkflowTitle  string
	Results        []ExtractionResultRowViewModel
	SelectedResult *ExtractionResultDetailsViewModel
	Pagination     PaginationViewModel
}
