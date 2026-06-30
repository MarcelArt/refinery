package routes

import (
	"github.com/MarcelArt/refinery/internal/enums"
	"github.com/MarcelArt/refinery/internal/v1/handlers"
	"github.com/MarcelArt/refinery/internal/v1/middlewares"
	"github.com/gin-gonic/gin"
)

func setupRateLimiterRoutes(v1 *gin.RouterGroup, authM *middlewares.AuthMiddleware, h *handlers.RateLimiterHandler) {
	rateLimiters := v1.Group("/rate-limiters", authM.Authn)

	rateLimiters.POST("/", authM.Authz(enums.PermRateLimitersCreate), h.Create)

	rateLimiters.GET("/", authM.Authz(enums.PermRateLimitersRead), h.Read)
	rateLimiters.GET("/:id", authM.Authz(enums.PermRateLimitersRead), h.GetByID)

	rateLimiters.PUT("/:id", authM.Authz(enums.PermRateLimitersUpdate), h.Update)

	rateLimiters.DELETE("/:id", authM.Authz(enums.PermRateLimitersDelete), h.Delete)
}
