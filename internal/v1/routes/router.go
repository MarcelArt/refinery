package routes

import (
	"github.com/MarcelArt/refinery/internal/v1/handlers"
	"github.com/MarcelArt/refinery/internal/v1/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(
	api *gin.RouterGroup,
	uHandler *handlers.UserHandler,
	wHandler *handlers.WorkflowHandler,
	erHandler *handlers.ExtractionResultHandler,
) {
	authM := middlewares.NewAuthMiddleware()

	v1 := api.Group("/v1")
	setupUserRoutes(v1, authM, uHandler)
	setupWorkflowRoutes(v1, authM, wHandler)
	setupExtractionResultRoutes(v1, authM, erHandler)
}
