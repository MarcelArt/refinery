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
