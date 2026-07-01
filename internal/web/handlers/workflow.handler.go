package handlers

import (
	"encoding/json"
	"html/template"
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

	// Fetch workflows owned by the current logged-in user.
	// Pagination is carried by the browser URL query string parsed automatically from context.
	pageInfo, pages := h.workflowService.GetByUserID(c, userId)

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
			Type:        p.Type,
		})
	}

	// Compute item bounds (e.g., "Showing 1 to 10 of 25")
	start := pageInfo.Page*pageInfo.Size + 1
	end := start + int64(len(workflowsVM)) - 1
	if len(workflowsVM) == 0 {
		start = 0
		end = 0
	}

	paginationVM := viewmodels.PaginationViewModel{
		Total:      pageInfo.Total,
		Page:       pageInfo.Page,
		Size:       pageInfo.Size,
		TotalPages: pageInfo.TotalPages,
		Last:       pageInfo.Last,
		First:      pageInfo.First,
		PrevPage:   pageInfo.Page - 1,
		NextPage:   pageInfo.Page + 1,
		Start:      start,
		End:        end,
	}

	renderTemplate(c, http.StatusOK, "workflows.html", gin.H{
		"Title":      "Workflows",
		"User":       user,
		"Workflows":  workflowsVM,
		"Pagination": paginationVM,
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
	workflowType := c.PostForm("type")

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
		Type:        workflowType,
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

// ShowUpdateWorkflow renders the edit/update workflow page
func (h *WorkflowWebHandler) ShowUpdateWorkflow(c *gin.Context) {
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

	workflowIDStr := c.Param("id")
	workflow, err := h.workflowService.GetByID(c, workflowIDStr)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/workflows")
		return
	}

	if workflow.UserID != uint(userId.(float64)) {
		c.Redirect(http.StatusSeeOther, "/workflows")
		return
	}

	schemas, _ := workflow.Schemas.Deserialize()
	
	// Convert schemas to JSON string for Alpine.js initialization
	schemasJsonBytes, _ := json.Marshal(schemas)
	schemasJsonStr := string(schemasJsonBytes)

	renderWorkflowTemplate(c, http.StatusOK, "update_workflow.html", "edit", workflow.Title, workflow.Description, workflow.ID, user, gin.H{
		"Workflow":    workflow,
		"SchemasJson": template.JS(schemasJsonStr),
	})
}

// HandleUpdateWorkflow processes the edit form submission
func (h *WorkflowWebHandler) HandleUpdateWorkflow(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	workflowIDStr := c.Param("id")
	workflow, err := h.workflowService.GetByID(c, workflowIDStr)
	if err != nil {
		renderFragment(c, http.StatusOK, "error_alert.html", gin.H{
			"Error": "Workflow not found",
		})
		return
	}

	if workflow.UserID != uint(userId.(float64)) {
		renderFragment(c, http.StatusOK, "error_alert.html", gin.H{
			"Error": "Unauthorized to update this workflow",
		})
		return
	}

	title := c.PostForm("title")
	description := c.PostForm("description")
	prompt := c.PostForm("prompt")
	schemasJson := c.PostForm("schemasJson")
	workflowType := c.PostForm("type")

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
		Type:        workflowType,
		UserID:      workflow.UserID, // Maintain original owner ID
	}

	err = h.workflowService.Update(c, workflow.ID, input)
	if err != nil {
		renderFragment(c, http.StatusOK, "error_alert.html", gin.H{
			"Error": "Failed to update workflow: " + err.Error(),
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

