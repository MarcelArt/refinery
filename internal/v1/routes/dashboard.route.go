package routes

import (
	"github.com/MarcelArt/refinery/internal/enums"
	"github.com/MarcelArt/refinery/internal/v1/handlers"
	"github.com/MarcelArt/refinery/internal/v1/middlewares"
	"github.com/gin-gonic/gin"
)

func setupDashboardRoutes(v1 *gin.RouterGroup, authM *middlewares.AuthMiddleware, h *handlers.DashboardHandler) {
	g := v1.Group("/dashboard", authM.Authn, authM.Authz(enums.PermDashboardRead))

	g.GET("/extraction-counts", h.ExtractionStatusCounts)
	g.GET("/daily-throughput", h.GetDailyThroughput)
	g.GET("/latency-stats", h.GetLatencyStats)
}
