package handlers

import (
	"net/http"

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
