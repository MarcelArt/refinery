package handlers

import (
	"net/http"

	"github.com/MarcelArt/refinery/internal/v1/services"
	"github.com/MarcelArt/refinery/internal/web/viewmodels"
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

// ShowDashboard renders the dashboard overview page
func (h *WorkflowWebHandler) ShowDashboard(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	user, err := h.userService.GetByID(c, userId)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	renderTemplate(c, http.StatusOK, "dashboard.html", gin.H{
		"Title":      "Dashboard",
		"User":       user,
		"ActiveMenu": "dashboard",
	})
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

	// Fetch workflows owned by the current logged-in user
	_, pages := h.workflowService.GetByUserID(c, userId)

	// Map backend model pages to frontend view models
	workflowsVM := make([]viewmodels.WorkflowRowViewModel, 0, len(pages))
	for _, p := range pages {
		schemas, _ := p.Schemas.Deserialize()
		workflowsVM = append(workflowsVM, viewmodels.WorkflowRowViewModel{
			ID:          p.ID,
			Title:       p.Title,
			Description: p.Description,
			Prompt:      p.Prompt,
			Schemas:     schemas,
		})
	}

	renderTemplate(c, http.StatusOK, "workflows.html", gin.H{
		"Title":      "Workflows",
		"User":       user,
		"Workflows":  workflowsVM,
		"ActiveMenu": "workflows",
	})
}
