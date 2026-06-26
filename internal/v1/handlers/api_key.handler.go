package handlers

import (
	"net/http"

	"github.com/MarcelArt/refinery/internal/common"
	"github.com/MarcelArt/refinery/internal/v1/models"
	"github.com/MarcelArt/refinery/internal/v1/services"
	"github.com/gin-gonic/gin"
	_ "github.com/morkid/paginate"
)

type ApiKeyHandler struct {
	service services.IApiKeyService
}

func NewApiKeyHandler(service services.IApiKeyService) *ApiKeyHandler {
	return &ApiKeyHandler{
		service: service,
	}
}

// Create godoc
// @Summary      Create a new API key
// @Description  Create a new API key with the provided details
// @Tags         api-keys
// @Accept       json
// @Produce      json
// @Param        apiKey  body      models.ApiKeyInput  true  "API key details"
// @Success      201   {object}  common.Result[uint]
// @Failure      400   {object}  common.Result[string]
// @Failure      401   {object}  common.Result[string]
// @Failure      500   {object}  common.Result[string]
// @Security     BearerAuth
// @Security     ApiKey
// @Router       /v1/api-keys [post]
func (h *ApiKeyHandler) Create(c *gin.Context) {
	userID, err := common.MustGet[float64](c, "userId")
	if err != nil {
		_, res := common.ResultErr(err, "token invalid")
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	var apiKey models.ApiKeyInput
	if err := c.ShouldBindJSON(&apiKey); err != nil {
		_, res := common.ResultErr(err, "failed parsing json")
		c.JSON(http.StatusBadRequest, res)
		return
	}
	apiKey.UserID = uint(userID)

	id, err := h.service.Generate(c, apiKey)
	if err != nil {
		c.JSON(common.ResultErr(err, "failed creating API key"))
		return
	}

	c.JSON(http.StatusCreated, common.ResultOk(id, "API key created"))
}

// Read godoc
// @Summary      List API keys
// @Description  Get a paginated list of API keys
// @Tags         api-keys
// @Accept       json
// @Produce      json
// @Param        page     query     int     false  "Page"
// @Param        size     query     int     false  "Size"
// @Param        sort     query     string  false  "Sort"
// @Param        filters  query     string  false  "Filter"
// @Success      200      {object}  paginate.Page{items=[]models.ApiKeyPage}
// @Failure      401      {object}  common.Result[string]
// @Failure      500      {object}  common.Result[string]
// @Security     BearerAuth
// @Security     ApiKey
// @Router       /v1/api-keys [get]
func (h *ApiKeyHandler) Read(c *gin.Context) {
	apiKeys, _ := h.service.Read(c)

	c.JSON(http.StatusOK, apiKeys)
}

// Update godoc
// @Summary      Update API key
// @Description  Update an existing API key's details
// @Tags         api-keys
// @Accept       json
// @Produce      json
// @Param        id      path      string              true  "API Key ID"
// @Param        apiKey  body      models.ApiKeyInput  true  "Updated API key details"
// @Success      200   {object}  common.Result[any]
// @Failure      400   {object}  common.Result[string]
// @Failure      401   {object}  common.Result[string]
// @Failure      500   {object}  common.Result[string]
// @Security     BearerAuth
// @Security     ApiKey
// @Router       /v1/api-keys/{id} [put]
func (h *ApiKeyHandler) Update(c *gin.Context) {
	id := c.Param("id")

	var apiKey models.ApiKeyInput
	if err := c.ShouldBindJSON(&apiKey); err != nil {
		_, res := common.ResultErr(err, "failed parsing json")
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if err := h.service.Update(c, id, apiKey); err != nil {
		c.JSON(common.ResultErr(err, "failed updating API key"))
		return
	}

	c.JSON(http.StatusOK, common.ResultOk[any](nil, "API key updated"))
}

// Delete godoc
// @Summary      Delete API key
// @Description  Delete an API key by ID
// @Tags         api-keys
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "API Key ID"
// @Success      200  {object}  common.Result[any]
// @Failure      401  {object}  common.Result[string]
// @Failure      500  {object}  common.Result[string]
// @Security     BearerAuth
// @Security     ApiKey
// @Router       /v1/api-keys/{id} [delete]
func (h *ApiKeyHandler) Delete(c *gin.Context) {
	if err := h.service.Delete(c, c.Param("id")); err != nil {
		c.JSON(common.ResultErr(err, "failed deleting API key"))
		return
	}

	c.JSON(http.StatusOK, common.ResultOk[any](nil, "API key deleted"))
}

// GetByID godoc
// @Summary      Get API key by ID
// @Description  Get detailed information about an API key by its ID
// @Tags         api-keys
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "API Key ID"
// @Success      200  {object}  common.Result[entities.ApiKey]
// @Failure      401  {object}  common.Result[string]
// @Failure      500  {object}  common.Result[string]
// @Security     BearerAuth
// @Security     ApiKey
// @Router       /v1/api-keys/{id} [get]
func (h *ApiKeyHandler) GetByID(c *gin.Context) {
	apiKey, err := h.service.GetByID(c, c.Param("id"))
	if err != nil {
		c.JSON(common.ResultErr(err, "failed getting API key"))
		return
	}

	c.JSON(http.StatusOK, common.ResultOk(apiKey, "API key found"))
}
