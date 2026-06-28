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
	r.GET("/account", authM.RequireAuth(), h.ShowAccount)
	r.POST("/account/api-keys", authM.RequireAuth(), h.HandleCreateApiKey)
	r.POST("/account/api-keys/regenerate", authM.RequireAuth(), h.HandleRegenerateApiKey)
	r.POST("/account/api-keys/delete", authM.RequireAuth(), h.HandleDeleteApiKey)
	r.POST("/account/password", authM.RequireAuth(), h.HandleChangePassword)
}
