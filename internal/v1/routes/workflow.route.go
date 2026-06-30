package routes

import (
	"github.com/MarcelArt/refinery/internal/enums"
	"github.com/MarcelArt/refinery/internal/v1/handlers"
	"github.com/MarcelArt/refinery/internal/v1/middlewares"
	"github.com/gin-gonic/gin"
)

func setupWorkflowRoutes(v1 *gin.RouterGroup, authM *middlewares.AuthMiddleware, h *handlers.WorkflowHandler, rateLimitM *middlewares.RateLimiterMiddleware) {
	workflows := v1.Group("/workflows", authM.Authn)

	workflows.POST("/", authM.Authz(enums.PermWorkflowsCreate), h.Create)
	workflows.POST("/:id/upload", authM.Authz(enums.PermWorkflowsUpload), rateLimitM.Limit, h.Upload)

	workflows.GET("/", authM.Authz(enums.PermWorkflowsRead), h.Read)
	workflows.GET("/:id", authM.Authz(enums.PermWorkflowsRead), h.GetByID)

	workflows.PUT("/:id", authM.Authz(enums.PermWorkflowsUpdate), h.Update)

	workflows.DELETE("/:id", authM.Authz(enums.PermWorkflowsDelete), h.Delete)
}
