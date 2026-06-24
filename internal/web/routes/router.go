package routes

import (
	"github.com/MarcelArt/refinery/internal/web/handlers"
	"github.com/gin-gonic/gin"
)

func SetupWebRoutes(
	r *gin.Engine,
	authM *WebAuthMiddleware,
	authH *handlers.AuthWebHandler,
	wfH *handlers.WorkflowWebHandler,
) {
	// Serve static assets
	r.Static("/public", "internal/web/public")

	// Initialize modular routes
	setupAuthRoutes(r, authM, authH)
	setupWorkflowRoutes(r, authM, wfH)
}
