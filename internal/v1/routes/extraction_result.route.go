package routes

import (
	"github.com/MarcelArt/refinery/internal/enums"
	"github.com/MarcelArt/refinery/internal/v1/handlers"
	"github.com/MarcelArt/refinery/internal/v1/middlewares"
	"github.com/gin-gonic/gin"
)

func setupExtractionResultRoutes(v1 *gin.RouterGroup, authM *middlewares.AuthMiddleware, h *handlers.ExtractionResultHandler) {
	g := v1.Group("/extraction-results")
	g.POST("/:id/webhook", authM.WebhookAuth, h.Webhook)
	g.POST("/:id/webhook/error", authM.WebhookAuth, h.WebhookError)

	extractionResults := v1.Group("/extraction-results", authM.Authn)

	extractionResults.POST("/", authM.Authz(enums.PermExtractionResultsCreate), h.Create)

	extractionResults.GET("/", authM.Authz(enums.PermExtractionResultsRead), h.Read)
	extractionResults.GET("/:id", authM.Authz(enums.PermExtractionResultsRead), h.GetByID)

	extractionResults.PUT("/:id", authM.Authz(enums.PermExtractionResultsUpdate), h.Update)

	extractionResults.DELETE("/:id", authM.Authz(enums.PermExtractionResultsDelete), h.Delete)
}
