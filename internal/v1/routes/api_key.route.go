package routes

import (
	"github.com/MarcelArt/refinery/internal/enums"
	"github.com/MarcelArt/refinery/internal/v1/handlers"
	"github.com/MarcelArt/refinery/internal/v1/middlewares"
	"github.com/gin-gonic/gin"
)

func setupApiKeyRoutes(v1 *gin.RouterGroup, authM *middlewares.AuthMiddleware, h *handlers.ApiKeyHandler) {
	apiKeys := v1.Group("/api-keys", authM.Authn)

	apiKeys.POST("/", authM.Authz(enums.PermApiKeysCreate), h.Create)

	apiKeys.GET("/", authM.Authz(enums.PermApiKeysRead), h.Read)
	apiKeys.GET("/:id", authM.Authz(enums.PermApiKeysRead), h.GetByID)

	apiKeys.PUT("/:id", authM.Authz(enums.PermApiKeysUpdate), h.Update)

	apiKeys.DELETE("/:id", authM.Authz(enums.PermApiKeysDelete), h.Delete)
}
