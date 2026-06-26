package handlers

import (
	"net/http"
	"strconv"

	"github.com/MarcelArt/refinery/internal/common"
	"github.com/MarcelArt/refinery/internal/v1/models"
	"github.com/MarcelArt/refinery/internal/v1/services"
	"github.com/MarcelArt/refinery/internal/web/viewmodels"
	"github.com/gin-gonic/gin"
)

type WebhookWebHandler struct {
	webhookService  services.IWebhookService
	workflowService services.IWorkflowService
	userService     services.IUserService
}

func NewWebhookWebHandler(
	webhookService services.IWebhookService,
	workflowService services.IWorkflowService,
	userService services.IUserService,
) *WebhookWebHandler {
	return &WebhookWebHandler{
		webhookService:  webhookService,
		workflowService: workflowService,
		userService:     userService,
	}
}

// ShowWorkflowWebhooks renders the webhooks settings tab for a workflow
func (h *WebhookWebHandler) ShowWorkflowWebhooks(c *gin.Context) {
	workflowIDStr := c.Param("id")
	h.renderListWithExtra(c, workflowIDStr, nil)
}

// HandleCreateWebhook handles the creation of a webhook
func (h *WebhookWebHandler) HandleCreateWebhook(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	workflowIDStr := c.Param("id")
	workflow, err := h.workflowService.GetByID(c, workflowIDStr)
	if err != nil {
		h.renderListWithExtra(c, workflowIDStr, gin.H{
			"Error": "Workflow not found",
		})
		return
	}

	if workflow.UserID != uint(userId.(float64)) {
		h.renderListWithExtra(c, workflowIDStr, gin.H{
			"Error": "Unauthorized access to workflow",
		})
		return
	}

	displayName := c.PostForm("displayName")
	url := c.PostForm("url")
	method := c.PostForm("method")
	hmacKey := c.PostForm("hmacKey")

	if displayName == "" || url == "" || method == "" || hmacKey == "" {
		h.renderListWithExtra(c, workflowIDStr, gin.H{
			"Error": "All fields (Display Name, URL, Method, and HMAC Key) are required",
		})
		return
	}

	input := models.WebhookInput{
		DisplayName: displayName,
		URL:         url,
		Method:      method,
		HmacKey:     hmacKey,
		WorkflowID:  workflow.ID,
	}

	_, err = h.webhookService.Create(c, input)
	if err != nil {
		h.renderListWithExtra(c, workflowIDStr, gin.H{
			"Error": "Failed to create webhook: " + err.Error(),
		})
		return
	}

	h.renderListWithExtra(c, workflowIDStr, nil)
}

// HandleUpdateWebhook handles the update/edit of a webhook
func (h *WebhookWebHandler) HandleUpdateWebhook(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	workflowIDStr := c.Param("id")
	workflow, err := h.workflowService.GetByID(c, workflowIDStr)
	if err != nil {
		h.renderListWithExtra(c, workflowIDStr, gin.H{
			"Error": "Workflow not found",
		})
		return
	}

	if workflow.UserID != uint(userId.(float64)) {
		h.renderListWithExtra(c, workflowIDStr, gin.H{
			"Error": "Unauthorized access to workflow",
		})
		return
	}

	webhookIDStr := c.PostForm("webhookId")
	webhookID, err := strconv.ParseUint(webhookIDStr, 10, 64)
	if err != nil {
		h.renderListWithExtra(c, workflowIDStr, gin.H{
			"Error": "Invalid Webhook ID",
		})
		return
	}

	webhook, err := h.webhookService.GetByID(c, uint(webhookID))
	if err != nil {
		h.renderListWithExtra(c, workflowIDStr, gin.H{
			"Error": "Webhook not found",
		})
		return
	}

	if webhook.WorkflowID != workflow.ID {
		h.renderListWithExtra(c, workflowIDStr, gin.H{
			"Error": "Webhook does not belong to this workflow",
		})
		return
	}

	displayName := c.PostForm("displayName")
	url := c.PostForm("url")
	method := c.PostForm("method")

	if displayName == "" || url == "" || method == "" {
		h.renderListWithExtra(c, workflowIDStr, gin.H{
			"Error": "Display Name, URL, and Method are required",
		})
		return
	}

	input := models.WebhookInput{
		DisplayName: displayName,
		URL:         url,
		Method:      method,
		HmacKey:     webhook.HmacKey, // preserve existing hmac key
		WorkflowID:  webhook.WorkflowID,
	}

	err = h.webhookService.Update(c, webhook.ID, input)
	if err != nil {
		h.renderListWithExtra(c, workflowIDStr, gin.H{
			"Error": "Failed to update webhook: " + err.Error(),
		})
		return
	}

	h.renderListWithExtra(c, workflowIDStr, nil)
}

// HandleDeleteWebhook handles revoking/deleting a webhook
func (h *WebhookWebHandler) HandleDeleteWebhook(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	workflowIDStr := c.Param("id")
	workflow, err := h.workflowService.GetByID(c, workflowIDStr)
	if err != nil {
		h.renderListWithExtra(c, workflowIDStr, gin.H{
			"Error": "Workflow not found",
		})
		return
	}

	if workflow.UserID != uint(userId.(float64)) {
		h.renderListWithExtra(c, workflowIDStr, gin.H{
			"Error": "Unauthorized access to workflow",
		})
		return
	}

	webhookIDStr := c.PostForm("webhookId")
	webhookID, err := strconv.ParseUint(webhookIDStr, 10, 64)
	if err != nil {
		h.renderListWithExtra(c, workflowIDStr, gin.H{
			"Error": "Invalid Webhook ID",
		})
		return
	}

	webhook, err := h.webhookService.GetByID(c, uint(webhookID))
	if err != nil {
		h.renderListWithExtra(c, workflowIDStr, gin.H{
			"Error": "Webhook not found",
		})
		return
	}

	if webhook.WorkflowID != workflow.ID {
		h.renderListWithExtra(c, workflowIDStr, gin.H{
			"Error": "Webhook does not belong to this workflow",
		})
		return
	}

	err = h.webhookService.Delete(c, webhook.ID)
	if err != nil {
		h.renderListWithExtra(c, workflowIDStr, gin.H{
			"Error": "Failed to delete webhook: " + err.Error(),
		})
		return
	}

	h.renderListWithExtra(c, workflowIDStr, nil)
}

func (h *WebhookWebHandler) renderListWithExtra(c *gin.Context, workflowIDStr string, extra gin.H) {
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

	workflow, err := h.workflowService.GetByID(c, workflowIDStr)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/workflows")
		return
	}

	if workflow.UserID != uint(userId.(float64)) {
		c.Redirect(http.StatusSeeOther, "/workflows")
		return
	}

	pageInfo, pages := h.webhookService.GetByWorkflowID(c, workflow.ID)

	start := pageInfo.Page*pageInfo.Size + 1
	end := start + int64(len(pages)) - 1
	if len(pages) == 0 {
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

	data := gin.H{
		"WorkflowID":       workflow.ID,
		"Webhooks":         pages,
		"GeneratedHmacKey": common.GenerateSecureKey("whsec_"),
		"Pagination":       paginationVM,
	}

	for k, v := range extra {
		data[k] = v
	}

	renderWorkflowTemplate(c, http.StatusOK, "webhooks_tab.html", "webhooks", workflow.Title, workflow.Description, workflow.ID, user, data)
}
