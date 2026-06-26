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
	erH *handlers.ExtractionResultWebHandler,
	akH *handlers.ApiKeyWebHandler,
	whH *handlers.WebhookWebHandler,
) {
	// Serve static assets
	r.Static("/public", "internal/web/public")

	// Initialize modular routes
	setupAuthRoutes(r, authM, authH)
	setupWorkflowRoutes(r, authM, wfH)
	setupExtractionResultRoutes(r, authM, erH)
	setupApiKeyRoutes(r, authM, akH)
	setupWebhookRoutes(r, authM, whH)
}
