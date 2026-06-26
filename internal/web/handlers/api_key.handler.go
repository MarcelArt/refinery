package handlers

import (
	"net/http"
	"strconv"

	"github.com/MarcelArt/refinery/internal/enums"
	"github.com/MarcelArt/refinery/internal/v1/models"
	"github.com/MarcelArt/refinery/internal/v1/services"
	"github.com/MarcelArt/refinery/internal/web/viewmodels"
	"github.com/MarcelArt/refinery/pkg/jsonb"
	"github.com/gin-gonic/gin"
)

type ApiKeyWebHandler struct {
	apiKeyService services.IApiKeyService
	userService   services.IUserService
}

func NewApiKeyWebHandler(
	apiKeyService services.IApiKeyService,
	userService services.IUserService,
) *ApiKeyWebHandler {
	return &ApiKeyWebHandler{
		apiKeyService: apiKeyService,
		userService:   userService,
	}
}

// ShowApiKeys renders the main API keys management page
func (h *ApiKeyWebHandler) ShowApiKeys(c *gin.Context) {
	h.renderListWithExtra(c, nil)
}

// HandleCreateApiKey processes API key creation
func (h *ApiKeyWebHandler) HandleCreateApiKey(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	name := c.PostForm("name")
	scopes := c.PostFormArray("scopes")

	if name == "" {
		h.renderListWithExtra(c, gin.H{
			"Error": "API Key name is required",
		})
		return
	}

	userIdVal := uint(userId.(float64))
	scopesJSONB, err := jsonb.New(scopes)
	if err != nil {
		h.renderListWithExtra(c, gin.H{
			"Error": "Failed to serialize scopes: " + err.Error(),
		})
		return
	}

	input := models.ApiKeyInput{
		Name:   name,
		Scopes: scopesJSONB,
		UserID: userIdVal,
	}

	key, err := h.apiKeyService.Generate(c, input)
	if err != nil {
		h.renderListWithExtra(c, gin.H{
			"Error": "Failed to create API key: " + err.Error(),
		})
		return
	}

	h.renderListWithExtra(c, gin.H{
		"CreatedKey":     key,
		"CreatedKeyName": name,
	})
}

// HandleRegenerateApiKey processes API key regeneration
func (h *ApiKeyWebHandler) HandleRegenerateApiKey(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	apiKeyIDStr := c.PostForm("id")
	apiKeyID, err := strconv.ParseUint(apiKeyIDStr, 10, 64)
	if err != nil {
		h.renderListWithExtra(c, gin.H{
			"Error": "Invalid API Key ID",
		})
		return
	}

	apiKey, err := h.apiKeyService.GetByID(c, uint(apiKeyID))
	if err != nil {
		h.renderListWithExtra(c, gin.H{
			"Error": "API Key not found",
		})
		return
	}

	if apiKey.UserID != uint(userId.(float64)) {
		h.renderListWithExtra(c, gin.H{
			"Error": "Unauthorized to regenerate this API key",
		})
		return
	}

	newKey, err := h.apiKeyService.Regenerate(c, apiKey.ID)
	if err != nil {
		h.renderListWithExtra(c, gin.H{
			"Error": "Failed to regenerate API key: " + err.Error(),
		})
		return
	}

	h.renderListWithExtra(c, gin.H{
		"RegeneratedKey":     newKey,
		"RegeneratedKeyName": apiKey.Name,
	})
}

// HandleDeleteApiKey revokes/deletes an API key
func (h *ApiKeyWebHandler) HandleDeleteApiKey(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	apiKeyIDStr := c.PostForm("id")
	apiKeyID, err := strconv.ParseUint(apiKeyIDStr, 10, 64)
	if err != nil {
		h.renderListWithExtra(c, gin.H{
			"Error": "Invalid API Key ID",
		})
		return
	}

	apiKey, err := h.apiKeyService.GetByID(c, uint(apiKeyID))
	if err != nil {
		h.renderListWithExtra(c, gin.H{
			"Error": "API Key not found",
		})
		return
	}

	if apiKey.UserID != uint(userId.(float64)) {
		h.renderListWithExtra(c, gin.H{
			"Error": "Unauthorized to delete this API key",
		})
		return
	}

	err = h.apiKeyService.Delete(c, apiKey.ID)
	if err != nil {
		h.renderListWithExtra(c, gin.H{
			"Error": "Failed to delete API key: " + err.Error(),
		})
		return
	}

	h.renderListWithExtra(c, nil)
}

func (h *ApiKeyWebHandler) renderListWithExtra(c *gin.Context, extra gin.H) {
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

	pageInfo, pages := h.apiKeyService.GetByUserID(c, userId)

	apiKeysVM := make([]viewmodels.ApiKeyRowViewModel, 0, len(pages))
	for _, p := range pages {
		scopes, _ := p.Scopes.Deserialize()
		apiKeysVM = append(apiKeysVM, viewmodels.ApiKeyRowViewModel{
			ID:     p.ID,
			Name:   p.Name,
			Scopes: scopes,
			UserID: p.UserID,
		})
	}

	start := pageInfo.Page*pageInfo.Size + 1
	end := start + int64(len(apiKeysVM)) - 1
	if len(apiKeysVM) == 0 {
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
		"Title":          "API Keys",
		"User":           user,
		"ApiKeys":        apiKeysVM,
		"Pagination":     paginationVM,
		"AvailablePerms": enums.AvailablePerms,
		"ActiveMenu":     "api-keys",
	}

	for k, v := range extra {
		data[k] = v
	}

	renderTemplate(c, http.StatusOK, "api_keys.html", data)
}
