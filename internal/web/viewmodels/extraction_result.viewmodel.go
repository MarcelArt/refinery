package viewmodels

import (
	"time"
)

type ExtractionResultRowViewModel struct {
	ID        uint
	CreatedAt time.Time
	Status    string
}

type ExtractionResultDetailsViewModel struct {
	ID           uint
	CreatedAt    time.Time
	Status       string
	FinishedAt   *time.Time
	Columns      []string
	Rows         []map[string]any
	Attachment   string
	WorkflowType string
}

func (vm ExtractionResultDetailsViewModel) FormattedFinishedAt() string {
	if vm.FinishedAt == nil {
		return "-"
	}
	return vm.FinishedAt.Format("2006-01-02 15:04:05")
}

func (vm ExtractionResultDetailsViewModel) Duration() string {
	if vm.FinishedAt == nil {
		return ""
	}
	diff := vm.FinishedAt.Sub(vm.CreatedAt)
	// Round to nearest millisecond to avoid long decimals
	return diff.Round(time.Millisecond).String()
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
