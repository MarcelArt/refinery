package routes

import (
	"github.com/MarcelArt/refinery/internal/v1/handlers"
	"github.com/MarcelArt/refinery/internal/v1/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	api *gin.RouterGroup,
	uHandler *handlers.UserHandler,
	wHandler *handlers.WorkflowHandler,
	erHandler *handlers.ExtractionResultHandler,
	akHandler *handlers.ApiKeyHandler,
	authM *middlewares.AuthMiddleware,
	rateLimitM *middlewares.RateLimiterMiddleware,
	whHandler *handlers.WebhookHandler,
	dHandler *handlers.DashboardHandler,
	rlHandler *handlers.RateLimiterHandler,
) {
	v1 := api.Group("/v1")
	setupUserRoutes(v1, authM, uHandler)
	setupWorkflowRoutes(v1, authM, wHandler, rateLimitM)
	setupExtractionResultRoutes(v1, authM, erHandler)
	setupApiKeyRoutes(v1, authM, akHandler)
	setupWebhookRoutes(v1, authM, whHandler)
	setupDashboardRoutes(v1, authM, dHandler)
	setupRateLimiterRoutes(v1, authM, rlHandler)
}
