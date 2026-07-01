package services

import (
	"context"
	"fmt"

	"github.com/MarcelArt/refinery/internal/enums"
	"github.com/MarcelArt/refinery/internal/v1/models"
	"github.com/MarcelArt/refinery/internal/v1/repositories"
)

type DashboardService struct {
	erRepo repositories.IExtractionResultRepo
}

func NewDashboardService(erRepo repositories.IExtractionResultRepo) *DashboardService {
	return &DashboardService{
		erRepo: erRepo,
	}
}

func (s *DashboardService) ExtractionStatusCounts(c context.Context, userID any) (models.ExtractionStatusCounts, error) {
	var counts models.ExtractionStatusCounts
	done, err := s.erRepo.GetStatusCount(c, enums.StatusDone, userID)
	if err != nil {
		return counts, fmt.Errorf("failed counting done status: %w", err)
	}
	counts.Done = done

	failed, err := s.erRepo.GetStatusCount(c, enums.StatusFailed, userID)
	if err != nil {
		return counts, fmt.Errorf("failed counting failed status: %w", err)
	}
	counts.Failed = failed

	counts.SuccessRate = (done / (done + failed)) * 100

	return counts, nil
}

func (s *DashboardService) GetDailyThroughput(c context.Context, userID any) ([]models.ThroughputPoint, error) {
	return s.erRepo.GetDailyThroughput(c, userID)
}

func (s *DashboardService) GetLatencyStats(c context.Context, userID any) (models.LatencyStats, error) {
	return s.erRepo.GetLatencyStats(c, userID)
}

func (s *DashboardService) GetWorkflowBreakdown(c context.Context, userID any) ([]models.WorkflowBreakdown, error) {
	rows, err := s.erRepo.GetWorkflowBreakdown(c, userID)
	if err != nil {
		return nil, fmt.Errorf("failed getting workflow breakdown: %w", err)
	}

	for i := range rows {
		row := &rows[i]
		terminal := row.Done + row.Failed
		if terminal > 0 {
			rate := float64(row.Done) / float64(terminal) * 100
			row.SuccessRate = &rate
		}
	}

	return rows, nil
}
