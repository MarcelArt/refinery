package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/MarcelArt/refinery/internal/entities"
	"github.com/MarcelArt/refinery/internal/v1/models"
	"github.com/MarcelArt/refinery/internal/v1/services"
	"github.com/MarcelArt/refinery/internal/web/viewmodels"
	"github.com/MarcelArt/refinery/pkg/jsonb"
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

// ShowCreateWorkflow renders the create workflow page
func (h *WorkflowWebHandler) ShowCreateWorkflow(c *gin.Context) {
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

	renderTemplate(c, http.StatusOK, "create_workflow.html", gin.H{
		"Title":      "Create Workflow",
		"User":       user,
		"ActiveMenu": "workflows",
	})
}

// HandleCreateWorkflow processes the creation form submission
func (h *WorkflowWebHandler) HandleCreateWorkflow(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	title := c.PostForm("title")
	description := c.PostForm("description")
	prompt := c.PostForm("prompt")
	schemasJson := c.PostForm("schemasJson")

	if title == "" || description == "" {
		renderFragment(c, http.StatusOK, "error_alert.html", gin.H{
			"Error": "Title and description are required",
		})
		return
	}

	var schemas []entities.WorkflowSchema
	if schemasJson != "" {
		if err := json.Unmarshal([]byte(schemasJson), &schemas); err != nil {
			renderFragment(c, http.StatusOK, "error_alert.html", gin.H{
				"Error": "Invalid schema format: " + err.Error(),
			})
			return
		}
	}

	// Filter out any schema entries that have empty keys to avoid saving malformed objects
	filteredSchemas := make([]entities.WorkflowSchema, 0, len(schemas))
	for _, s := range schemas {
		if s.Key != "" {
			filteredSchemas = append(filteredSchemas, s)
		}
	}

	schemasJSONB, err := jsonb.New(filteredSchemas)
	if err != nil {
		renderFragment(c, http.StatusOK, "error_alert.html", gin.H{
			"Error": "Failed to create schema representation: " + err.Error(),
		})
		return
	}

	input := models.WorkflowInput{
		Title:       title,
		Description: description,
		Prompt:      prompt,
		Schemas:     schemasJSONB,
		UserID:      uint(userId.(float64)),
	}

	_, err = h.workflowService.Create(c, input)
	if err != nil {
		renderFragment(c, http.StatusOK, "error_alert.html", gin.H{
			"Error": "Failed to save workflow: " + err.Error(),
		})
		return
	}

	if c.GetHeader("HX-Request") == "true" {
		c.Header("HX-Redirect", "/workflows")
		c.Status(http.StatusOK)
	} else {
		c.Redirect(http.StatusSeeOther, "/workflows")
	}
}
