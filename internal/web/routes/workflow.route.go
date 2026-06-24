package routes

import (
	"net/http"

	"github.com/MarcelArt/refinery/internal/web/handlers"
	"github.com/gin-gonic/gin"
)

func setupWorkflowRoutes(r *gin.Engine, authM *WebAuthMiddleware, h *handlers.WorkflowWebHandler) {
	r.GET("/dashboard", authM.RequireAuth(), h.ShowDashboard)
	r.GET("/workflows", authM.RequireAuth(), h.ShowWorkflows)
	
	// Workflow creation routes
	r.GET("/workflows/create", authM.RequireAuth(), h.ShowCreateWorkflow)
	r.POST("/workflows/create", authM.RequireAuth(), h.HandleCreateWorkflow)
	
	// Root redirects to dashboard page (which is protected by RequireAuth)
	r.GET("/", authM.RequireAuth(), func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/dashboard")
	})
}
