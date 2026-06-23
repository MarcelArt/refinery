package routes

import (
	"github.com/MarcelArt/refinery/internal/v1/handlers"
	"github.com/MarcelArt/refinery/internal/v1/middlewares"
	"github.com/gin-gonic/gin"
)

func setupUserRoutes(v1 *gin.RouterGroup, authM *middlewares.AuthMiddleware, h *handlers.UserHandler) {
	users := v1.Group("/users")

	users.POST("/", h.Create)
	users.POST("/login", h.Login)
	users.POST("/refresh", authM.Refresh, h.Refresh)

	users.GET("/", h.Read)
	users.GET("/current", authM.Authn, h.GetCurrent)
	users.GET("/:id", h.GetByID)

	users.PUT("/:id", authM.Refresh, h.Update)

	users.DELETE("/:id", authM.Refresh, h.Delete)
}
