package routes

import (
	"github.com/MarcelArt/refinery/internal/web/handlers"
	"github.com/gin-gonic/gin"
)

func setupExtractionResultRoutes(
	r *gin.Engine,
	authM *WebAuthMiddleware,
	h *handlers.ExtractionResultWebHandler,
) {
	// Results page (defaulting to first run)
	r.GET("/workflows/:id/results", authM.RequireAuth(), h.ShowResults)

	// Results page preselecting a specific run
	r.GET("/workflows/:id/results/:resultId", authM.RequireAuth(), h.ShowResults)

	// AJAX details endpoint to fetch the right side table fragment
	r.GET("/workflows/:id/results/details/:resultId", authM.RequireAuth(), h.ShowResultDetails)

	// Start workflow by uploading a file
	r.POST("/workflows/:id/results/upload", authM.RequireAuth(), h.Upload)
}
