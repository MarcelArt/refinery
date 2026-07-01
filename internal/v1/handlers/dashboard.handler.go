package handlers

import (
	"net/http"
	"strconv"

	"github.com/MarcelArt/refinery/internal/common"
	"github.com/MarcelArt/refinery/internal/v1/services"
	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	service *services.DashboardService
}

func NewDashboardHandler(service *services.DashboardService) *DashboardHandler {
	return &DashboardHandler{
		service: service,
	}
}

// ExtractionStatusCounts godoc
// @Summary      Get extraction status counts
// @Description  Get extraction status counts (done, failed, success rate) for the authenticated user
// @Tags         dashboard
// @Accept       json
// @Produce      json
// @Success      200  {object}  common.Result[models.ExtractionStatusCounts]
// @Failure      401  {object}  common.Result[string]
// @Failure      403  {object}  common.Result[string]
// @Failure      500  {object}  common.Result[string]
// @Security     BearerAuth
// @Security     ApiKey
// @Router       /v1/dashboard/extraction-counts [get]
func (h *DashboardHandler) ExtractionStatusCounts(c *gin.Context) {
	userID := c.MustGet("userId")

	counts, err := h.service.ExtractionStatusCounts(c.Request.Context(), userID)
	if err != nil {
		c.JSON(common.ResultErr(err, "failed counting extraction statuses"))
		return
	}

	c.JSON(http.StatusOK, common.ResultOk(counts, "Extraction status counted"))
}

// GetDailyThroughput godoc
// @Summary      Get daily throughput
// @Description  Get daily throughput points (bucket, done, failed, in progress counts) for the authenticated user
// @Tags         dashboard
// @Accept       json
// @Produce      json
// @Success      200  {object}  common.Result[[]models.ThroughputPoint]
// @Failure      401  {object}  common.Result[string]
// @Failure      403  {object}  common.Result[string]
// @Failure      500  {object}  common.Result[string]
// @Security     BearerAuth
// @Security     ApiKey
// @Router       /v1/dashboard/daily-throughput [get]
func (h *DashboardHandler) GetDailyThroughput(c *gin.Context) {
	userID := c.MustGet("userId")
	throughputPoints, err := h.service.GetDailyThroughput(c.Request.Context(), userID)
	if err != nil {
		c.JSON(common.ResultErr(err, "failed getting daily throughput"))
		return
	}

	c.JSON(http.StatusOK, common.ResultOk(throughputPoints, "Daily throughput retrieved"))
}

// GetLatencyStats godoc
// @Summary      Get latency stats
// @Description  Get latency statistics (completed count, average seconds, P50, and P95 latency) for the authenticated user
// @Tags         dashboard
// @Accept       json
// @Produce      json
// @Success      200  {object}  common.Result[models.LatencyStats]
// @Failure      401  {object}  common.Result[string]
// @Failure      403  {object}  common.Result[string]
// @Failure      500  {object}  common.Result[string]
// @Security     BearerAuth
// @Security     ApiKey
// @Router       /v1/dashboard/latency-stats [get]
func (h *DashboardHandler) GetLatencyStats(c *gin.Context) {
	userID := c.MustGet("userId")
	latencyStats, err := h.service.GetLatencyStats(c.Request.Context(), userID)
	if err != nil {
		c.JSON(common.ResultErr(err, "failed getting latency stats"))
		return
	}

	c.JSON(http.StatusOK, common.ResultOk(latencyStats, "Latency stats retrieved"))
}

// GetWorkflowBreakdown godoc
// @Summary      Get workflow breakdown
// @Description  Get a breakdown of runs, statuses, latencies, and success rates grouped by workflow for the authenticated user
// @Tags         dashboard
// @Accept       json
// @Produce      json
// @Success      200  {object}  common.Result[[]models.WorkflowBreakdown]
// @Failure      401  {object}  common.Result[string]
// @Failure      403  {object}  common.Result[string]
// @Failure      500  {object}  common.Result[string]
// @Security     BearerAuth
// @Security     ApiKey
// @Router       /v1/dashboard/workflow-breakdown [get]
func (h *DashboardHandler) GetWorkflowBreakdown(c *gin.Context) {
	userID := c.MustGet("userId")
	breakdowns, err := h.service.GetWorkflowBreakdown(c.Request.Context(), userID)
	if err != nil {
		c.JSON(common.ResultErr(err, "failed getting workflow breakdown"))
		return
	}

	c.JSON(http.StatusOK, common.ResultOk(breakdowns, "Workflow breakdown retrieved"))
}

// GetLastNExtractions godoc
// @Summary      Get last N extractions
// @Description  Get last N extractions activity for the authenticated user
// @Tags         dashboard
// @Accept       json
// @Produce      json
// @Param        n    query     int  true  "Number of extractions to retrieve"
// @Success      200  {object}  common.Result[[]models.ExtractionActivity]
// @Failure      400  {object}  common.Result[string]
// @Failure      401  {object}  common.Result[string]
// @Failure      403  {object}  common.Result[string]
// @Failure      500  {object}  common.Result[string]
// @Security     BearerAuth
// @Security     ApiKey
// @Router       /v1/dashboard/last-extractions [get]
func (h *DashboardHandler) GetLastNExtractions(c *gin.Context) {
	userID := c.MustGet("userId")
	n := c.Query("n")
	limit, err := strconv.Atoi(n)
	if err != nil {
		_, res := common.ResultErr(err, "invalid limit parameter")
		c.JSON(http.StatusBadRequest, res)
		return
	}

	extractions, err := h.service.GetLastNExtractions(c.Request.Context(), userID, limit)
	if err != nil {
		c.JSON(common.ResultErr(err, "failed getting last N extractions"))
		return
	}

	c.JSON(http.StatusOK, common.ResultOk(extractions, "Last N extractions retrieved"))
}
