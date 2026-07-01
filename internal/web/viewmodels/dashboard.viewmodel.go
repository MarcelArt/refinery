package viewmodels

import (
	"fmt"
	"math"
	"time"
)

type LatencyStatsViewModel struct {
	Completed  int64
	AvgSeconds string
	P50Seconds string
	P95Seconds string
}

type WorkflowBreakdownViewModel struct {
	WorkflowID        uint
	WorkflowTitle     string
	WorkflowType      string
	TotalRuns         int64
	SuccessRate       string
	AvgLatencySeconds string
	P95LatencySeconds string
	LastRunAt         string
}

type RecentActivityViewModel struct {
	ExtractionID uint
	WorkflowID   uint
	Workflow     string
	Attachment   string
	Status       string
	StatusClass  string
	CreatedAt    string
	Duration     string
	Route        string
}

// FormatLatency converts a latency pointer into a human-readable duration string (e.g. 4.2s, 1m 4s, or —).
func FormatLatency(seconds *float64) string {
	if seconds == nil {
		return "—"
	}
	sec := *seconds
	if sec < 0 {
		sec = 0
	}
	if sec < 60 {
		return fmt.Sprintf("%.1fs", sec)
	}
	mins := int(math.Floor(sec / 60))
	secs := math.Mod(sec, 60)
	return fmt.Sprintf("%dm %.0fs", mins, secs)
}

// FormatRelativeTime formats a time pointer into a relative string representation (e.g. "3 minutes ago", "yesterday").
func FormatRelativeTime(t *time.Time) string {
	if t == nil || t.IsZero() {
		return "—"
	}
	now := time.Now()
	diff := now.Sub(*t)
	if diff < 0 {
		return "just now"
	}
	if diff < time.Minute {
		return fmt.Sprintf("%ds ago", int(diff.Seconds()))
	}
	if diff < time.Hour {
		return fmt.Sprintf("%dm ago", int(diff.Minutes()))
	}
	if diff < 24*time.Hour {
		return fmt.Sprintf("%dh ago", int(diff.Hours()))
	}
	days := int(diff.Hours() / 24)
	if days == 1 {
		return "yesterday"
	}
	return fmt.Sprintf("%dd ago", days)
}

// GetStatusClass returns a CSS class name for a given status
func GetStatusClass(status string) string {
	switch status {
	case "DONE":
		return "status-done"
	case "FAILED":
		return "status-failed"
	case "IN_PROGRESS":
		return "status-in-progress"
	default:
		return "status-muted"
	}
}
