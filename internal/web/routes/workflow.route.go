package routes

import (
	"net/http"

	"github.com/MarcelArt/refinery/internal/web/handlers"
	"github.com/gin-gonic/gin"
)

func setupWorkflowRoutes(r *gin.Engine, authM *WebAuthMiddleware, h *handlers.WorkflowWebHandler) {
	r.GET("/workflows", authM.RequireAuth(), h.ShowWorkflows)
	
	// Workflow creation routes
	r.GET("/workflows/create", authM.RequireAuth(), h.ShowCreateWorkflow)
	r.POST("/workflows/create", authM.RequireAuth(), h.HandleCreateWorkflow)
	
	// Workflow update routes
	r.GET("/workflows/:id/edit", authM.RequireAuth(), h.ShowUpdateWorkflow)
	r.POST("/workflows/:id/edit", authM.RequireAuth(), h.HandleUpdateWorkflow)
	

	// Redirect base workflow ID path to results tab
	r.GET("/workflows/:id", authM.RequireAuth(), func(c *gin.Context) {
		c.Redirect(http.StatusSeeOther, "/workflows/"+c.Param("id")+"/results")
	})
}
