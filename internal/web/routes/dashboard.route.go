package routes

import (
	"github.com/MarcelArt/refinery/internal/web/handlers"
	"github.com/gin-gonic/gin"
)

func setupDashboardRoutes(r *gin.Engine, authM *WebAuthMiddleware, h *handlers.DashboardWebHandler) {
	r.GET("/dashboard", authM.RequireAuth(), h.ShowDashboard)
}
