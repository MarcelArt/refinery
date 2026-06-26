package routes

import (
	"github.com/MarcelArt/refinery/internal/web/handlers"
	"github.com/gin-gonic/gin"
)

func setupApiKeyRoutes(
	r *gin.Engine,
	authM *WebAuthMiddleware,
	h *handlers.ApiKeyWebHandler,
) {
	r.GET("/api-keys", authM.RequireAuth(), h.ShowApiKeys)
	r.POST("/api-keys", authM.RequireAuth(), h.HandleCreateApiKey)
	r.POST("/api-keys/regenerate", authM.RequireAuth(), h.HandleRegenerateApiKey)
	r.POST("/api-keys/delete", authM.RequireAuth(), h.HandleDeleteApiKey)
}
