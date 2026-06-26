package routes

import (
	"github.com/MarcelArt/refinery/internal/web/handlers"
	"github.com/gin-gonic/gin"
)

func setupWebhookRoutes(
	r *gin.Engine,
	authM *WebAuthMiddleware,
	h *handlers.WebhookWebHandler,
) {
	r.GET("/workflows/:id/webhooks", authM.RequireAuth(), h.ShowWorkflowWebhooks)
	r.POST("/workflows/:id/webhooks", authM.RequireAuth(), h.HandleCreateWebhook)
	r.POST("/workflows/:id/webhooks/edit", authM.RequireAuth(), h.HandleUpdateWebhook)
	r.POST("/workflows/:id/webhooks/delete", authM.RequireAuth(), h.HandleDeleteWebhook)
}
