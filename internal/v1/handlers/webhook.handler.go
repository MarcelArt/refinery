package handlers

import (
	"net/http"

	"github.com/MarcelArt/refinery/internal/common"
	"github.com/MarcelArt/refinery/internal/v1/models"
	"github.com/MarcelArt/refinery/internal/v1/services"
	"github.com/gin-gonic/gin"
	_ "github.com/morkid/paginate"
)

type WebhookHandler struct {
	service services.IWebhookService
}

func NewWebhookHandler(service services.IWebhookService) *WebhookHandler {
	return &WebhookHandler{
		service: service,
	}
}

// Create godoc
// @Summary      Create a new webhook
// @Description  Create a new webhook with the provided details
// @Tags         webhooks
// @Accept       json
// @Produce      json
// @Param        webhook  body      models.WebhookInput  true  "Webhook details"
// @Success      201   {object}  common.Result[uint]
// @Failure      400   {object}  common.Result[string]
// @Failure      401   {object}  common.Result[string]
// @Failure      500   {object}  common.Result[string]
// @Security     BearerAuth
// @Security     ApiKey
// @Router       /v1/webhooks [post]
func (h *WebhookHandler) Create(c *gin.Context) {
	var webhook models.WebhookInput
	if err := c.ShouldBindJSON(&webhook); err != nil {
		_, res := common.ResultErr(err, "failed parsing json")
		c.JSON(http.StatusBadRequest, res)
		return
	}

	id, err := h.service.CreateWithHMAC(c, webhook)
	if err != nil {
		c.JSON(common.ResultErr(err, "failed creating webhook"))
		return
	}

	c.JSON(http.StatusCreated, common.ResultOk(id, "Webhook created"))
}

// Read godoc
// @Summary      List webhooks
// @Description  Get a paginated list of webhooks
// @Tags         webhooks
// @Accept       json
// @Produce      json
// @Param        page     query     int     false  "Page"
// @Param        size     query     int     false  "Size"
// @Param        sort     query     string  false  "Sort"
// @Param        filters  query     string  false  "Filter"
// @Success      200      {object}  paginate.Page{items=[]models.WebhookPage}
// @Failure      401      {object}  common.Result[string]
// @Failure      500      {object}  common.Result[string]
// @Security     BearerAuth
// @Security     ApiKey
// @Router       /v1/webhooks [get]
func (h *WebhookHandler) Read(c *gin.Context) {
	webhooks, _ := h.service.Read(c)

	c.JSON(http.StatusOK, webhooks)
}

// Update godoc
// @Summary      Update webhook
// @Description  Update an existing webhook's details
// @Tags         webhooks
// @Accept       json
// @Produce      json
// @Param        id       path      string               true  "Webhook ID"
// @Param        webhook  body      models.WebhookInput  true  "Updated webhook details"
// @Success      200   {object}  common.Result[any]
// @Failure      400   {object}  common.Result[string]
// @Failure      401   {object}  common.Result[string]
// @Failure      500   {object}  common.Result[string]
// @Security     BearerAuth
// @Security     ApiKey
// @Router       /v1/webhooks/{id} [put]
func (h *WebhookHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var webhook models.WebhookInput
	if err := c.ShouldBindJSON(&webhook); err != nil {
		_, res := common.ResultErr(err, "failed parsing json")
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if err := h.service.Update(c, id, webhook); err != nil {
		c.JSON(common.ResultErr(err, "failed updating webhook"))
		return
	}

	c.JSON(http.StatusOK, common.ResultOk[any](nil, "Webhook updated"))
}

// Delete godoc
// @Summary      Delete webhook
// @Description  Delete a webhook by ID
// @Tags         webhooks
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Webhook ID"
// @Success      200  {object}  common.Result[any]
// @Failure      401  {object}  common.Result[string]
// @Failure      500  {object}  common.Result[string]
// @Security     BearerAuth
// @Security     ApiKey
// @Router       /v1/webhooks/{id} [delete]
func (h *WebhookHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c, c.Param("id")); err != nil {
		c.JSON(common.ResultErr(err, "failed deleting webhook"))
		return
	}

	c.JSON(http.StatusOK, common.ResultOk[any](nil, "Webhook deleted"))
}

// GetByID godoc
// @Summary      Get webhook by ID
// @Description  Get detailed information about a webhook by its ID
// @Tags         webhooks
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "Webhook ID"
// @Success      200  {object}  common.Result[entities.Webhook]
// @Failure      401  {object}  common.Result[string]
// @Failure      500  {object}  common.Result[string]
// @Security     BearerAuth
// @Security     ApiKey
// @Router       /v1/webhooks/{id} [get]
func (h *WebhookHandler) GetByID(c *gin.Context) {
	webhook, err := h.service.GetByID(c, c.Param("id"))
	if err != nil {
		c.JSON(common.ResultErr(err, "failed getting webhook"))
		return
	}

	c.JSON(http.StatusOK, common.ResultOk(webhook, "Webhook found"))
}
