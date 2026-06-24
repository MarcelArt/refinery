package handlers

import (
	"net/http"

	"github.com/MarcelArt/refinery/internal/v1/services"
	"github.com/gin-gonic/gin"
)

type WorkflowWebHandler struct {
	workflowService services.IWorkflowService
	userService     services.IUserService
}

func NewWorkflowWebHandler(
	workflowService services.IWorkflowService,
	userService services.IUserService,
) *WorkflowWebHandler {
	return &WorkflowWebHandler{
		workflowService: workflowService,
		userService:     userService,
	}
}

// ShowWorkflows renders the main workflows dashboard page
func (h *WorkflowWebHandler) ShowWorkflows(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	user, err := h.userService.GetByID(c, userId)
	if err != nil {
		// If user cannot be fetched, redirect to login
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	renderTemplate(c, http.StatusOK, "workflows.html", gin.H{
		"Title": "Workflows",
		"User":  user,
	})
}
