package routes

import (
	"github.com/MarcelArt/refinery/internal/enums"
	"github.com/MarcelArt/refinery/internal/v1/handlers"
	"github.com/MarcelArt/refinery/internal/v1/middlewares"
	"github.com/gin-gonic/gin"
)

func setupWebhookRoutes(v1 *gin.RouterGroup, authM *middlewares.AuthMiddleware, h *handlers.WebhookHandler) {
	webhooks := v1.Group("/webhooks", authM.Authn)

	webhooks.POST("/", authM.Authz(enums.PermWebhooksCreate), h.Create)

	webhooks.GET("/", authM.Authz(enums.PermWebhooksRead), h.Read)
	webhooks.GET("/:id", authM.Authz(enums.PermWebhooksRead), h.GetByID)

	webhooks.PUT("/:id", authM.Authz(enums.PermWebhooksUpdate), h.Update)

	webhooks.DELETE("/:id", authM.Authz(enums.PermWebhooksDelete), h.Delete)
}
