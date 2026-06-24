package routes

import (
	"github.com/MarcelArt/refinery/internal/v1/handlers"
	"github.com/MarcelArt/refinery/internal/v1/middlewares"
	"github.com/gin-gonic/gin"
)

func setupWorkflowRoutes(v1 *gin.RouterGroup, authM *middlewares.AuthMiddleware, h *handlers.WorkflowHandler) {
	workflows := v1.Group("/workflows", authM.Authn)

	workflows.POST("/", h.Create)
	workflows.GET("/", h.Read)
	workflows.GET("/:id", h.GetByID)
	workflows.PUT("/:id", h.Update)
	workflows.DELETE("/:id", h.Delete)
}
