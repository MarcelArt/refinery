package routes

import (
	"github.com/MarcelArt/refinery/internal/v1/handlers"
	"github.com/MarcelArt/refinery/internal/v1/middlewares"
	"github.com/gin-gonic/gin"
)

func setupExtractionResultRoutes(v1 *gin.RouterGroup, authM *middlewares.AuthMiddleware, h *handlers.ExtractionResultHandler) {
	g := v1.Group("/extraction-results")
	g.POST("/:id/webhook", h.Webhook)

	extractionResults := v1.Group("/extraction-results", authM.Authn)

	extractionResults.POST("/", h.Create)

	extractionResults.GET("/", h.Read)
	extractionResults.GET("/:id", h.GetByID)

	extractionResults.PUT("/:id", h.Update)

	extractionResults.DELETE("/:id", h.Delete)
}
