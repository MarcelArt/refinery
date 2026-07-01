package models

import "time"

type ExtractionStatusCounts struct {
	Done        float32 `json:"done"`
	Failed      float32 `json:"failed"`
	SuccessRate float32 `json:"successRate"`
}

type ThroughputPoint struct {
	Bucket     time.Time `json:"bucket"`
	Done       int64     `json:"done"`
	Failed     int64     `json:"failed"`
	InProgress int64     `json:"inProgress"`
}

type LatencyStats struct {
	Completed  int64    `json:"completed"`
	AvgSeconds *float64 `json:"avgSeconds"`
	P50Seconds *float64 `json:"p50Seconds"`
	P95Seconds *float64 `json:"p95Seconds"`
}

type WorkflowBreakdown struct {
	WorkflowID        uint       `json:"workflowId"`
	WorkflowTitle     string     `json:"workflowTitle"`
	WorkflowType      string     `json:"workflowType"`
	TotalRuns         int64      `json:"totalRuns"`
	Done              int64      `json:"done"`
	Failed            int64      `json:"failed"`
	InProgress        int64      `json:"inProgress"`
	LastRunAt         *time.Time `json:"lastRunAt"`
	AvgLatencySeconds *float64   `json:"avgLatencySeconds"`
	P95LatencySeconds *float64   `json:"p95LatencySeconds"`
	SuccessRate       *float64   `json:"successRate"` // computed in Go
}

type ExtractionActivity struct {
	ExtractionID uint      `json:"extractionId"`
	WorkflowID   uint      `json:"workflowId"`
	Workflow     string    `json:"workflow"`
	Attachment   string    `json:"attachment"`
	Status       string    `json:"status"`
	CreatedAt    time.Time `json:"createdAt"`
	FinishedAt   time.Time `json:"finishedAt"`
	Route        string    `json:"route"`
	UserID       uint      `json:"userId"`
}
