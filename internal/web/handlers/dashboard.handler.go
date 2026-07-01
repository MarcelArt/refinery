package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"github.com/MarcelArt/refinery/internal/v1/services"
	"github.com/MarcelArt/refinery/internal/web/viewmodels"
	"github.com/gin-gonic/gin"
)

type DashboardWebHandler struct {
	dashboardService *services.DashboardService
	userService      services.IUserService
}

func NewDashboardWebHandler(
	dashboardService *services.DashboardService,
	userService services.IUserService,
) *DashboardWebHandler {
	return &DashboardWebHandler{
		dashboardService: dashboardService,
		userService:      userService,
	}
}

// ShowDashboard retrieves status counts, daily throughput, latency, and recent activities
func (h *DashboardWebHandler) ShowDashboard(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	user, err := h.userService.GetByID(c, userId)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	ctx := c.Request.Context()

	statusCounts, err := h.dashboardService.ExtractionStatusCounts(ctx, userId)
	if err != nil {
		log.Printf("Error fetching ExtractionStatusCounts: %v", err)
	}

	throughput, err := h.dashboardService.GetDailyThroughput(ctx, userId)
	if err != nil {
		log.Printf("Error fetching DailyThroughput: %v", err)
	}

	latency, err := h.dashboardService.GetLatencyStats(ctx, userId)
	if err != nil {
		log.Printf("Error fetching LatencyStats: %v", err)
	}

	breakdown, err := h.dashboardService.GetWorkflowBreakdown(ctx, userId)
	if err != nil {
		log.Printf("Error fetching WorkflowBreakdown: %v", err)
	}

	recent, err := h.dashboardService.GetLastNExtractions(ctx, userId, 10)
	if err != nil {
		log.Printf("Error fetching LastNExtractions: %v", err)
	}

	// Format Latency Stats to View Model
	var latencyVM viewmodels.LatencyStatsViewModel
	latencyVM = viewmodels.LatencyStatsViewModel{
		Completed:  latency.Completed,
		AvgSeconds: viewmodels.FormatLatency(latency.AvgSeconds),
		P50Seconds: viewmodels.FormatLatency(latency.P50Seconds),
		P95Seconds: viewmodels.FormatLatency(latency.P95Seconds),
	}

	// Format Workflow Breakdown to View Models
	breakdownVM := make([]viewmodels.WorkflowBreakdownViewModel, 0, len(breakdown))
	for _, b := range breakdown {
		successRateStr := "—"
		if b.SuccessRate != nil {
			successRateStr = fmt.Sprintf("%.1f%%", *b.SuccessRate)
		}
		breakdownVM = append(breakdownVM, viewmodels.WorkflowBreakdownViewModel{
			WorkflowID:        b.WorkflowID,
			WorkflowTitle:     b.WorkflowTitle,
			WorkflowType:      b.WorkflowType,
			TotalRuns:         b.TotalRuns,
			SuccessRate:       successRateStr,
			AvgLatencySeconds: viewmodels.FormatLatency(b.AvgLatencySeconds),
			P95LatencySeconds: viewmodels.FormatLatency(b.P95LatencySeconds),
			LastRunAt:         viewmodels.FormatRelativeTime(b.LastRunAt),
		})
	}

	// Format Recent Activity to View Models
	recentVM := make([]viewmodels.RecentActivityViewModel, 0, len(recent))
	for _, r := range recent {
		createdAtStr := viewmodels.FormatRelativeTime(&r.CreatedAt)
		
		durationStr := "—"
		if !r.FinishedAt.IsZero() {
			diffSec := r.FinishedAt.Sub(r.CreatedAt).Seconds()
			durationStr = viewmodels.FormatLatency(&diffSec)
		}

		attachmentName := filepath.Base(r.Attachment)

		recentVM = append(recentVM, viewmodels.RecentActivityViewModel{
			ExtractionID: r.ExtractionID,
			WorkflowID:   r.WorkflowID,
			Workflow:     r.Workflow,
			Attachment:   attachmentName,
			Status:       r.Status,
			StatusClass:  viewmodels.GetStatusClass(r.Status),
			CreatedAt:    createdAtStr,
			Duration:     durationStr,
			Route:        r.Route,
		})
	}

	// Marshal throughput data for Chart.js
	throughputBytes, err := json.Marshal(throughput)
	if err != nil {
		log.Printf("Error marshalling throughput data: %v", err)
		throughputBytes = []byte("[]")
	}

	renderTemplate(c, http.StatusOK, "dashboard.html", gin.H{
		"Title":          "Dashboard",
		"User":           user,
		"ActiveMenu":     "dashboard",
		"StatusCounts":    statusCounts,
		"Throughput":      throughput,
		"ThroughputJSON":  template.JS(throughputBytes),
		"Latency":         latencyVM,
		"Breakdown":       breakdownVM,
		"RecentActivity":  recentVM,
	})
}
