package handlers

import (
	"net/http"
	"strconv"

	"github.com/MarcelArt/refinery/internal/enums"
	"github.com/MarcelArt/refinery/internal/v1/models"
	"github.com/MarcelArt/refinery/internal/v1/services"
	"github.com/MarcelArt/refinery/internal/web/viewmodels"
	"github.com/MarcelArt/refinery/pkg/jsonb"
	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
)

type ApiKeyWebHandler struct {
	apiKeyService      services.IApiKeyService
	userService        services.IUserService
	rateLimiterService services.IRateLimiterService
}

func NewApiKeyWebHandler(
	apiKeyService services.IApiKeyService,
	userService services.IUserService,
	rateLimiterService services.IRateLimiterService,
) *ApiKeyWebHandler {
	return &ApiKeyWebHandler{
		apiKeyService:      apiKeyService,
		userService:        userService,
		rateLimiterService: rateLimiterService,
	}
}

// ShowAccount renders the main account settings page
func (h *ApiKeyWebHandler) ShowAccount(c *gin.Context) {
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

	// Retrieve daily rate limit usage count
	var usageCount uint
	rateLimit, err := h.rateLimiterService.GetTodayByUserID(c, userId)
	if err == nil {
		usageCount = rateLimit.Count
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

	var usagePercent uint
	if user.DailyLimit > 0 {
		usagePercent = (usageCount * 100) / user.DailyLimit
		if usagePercent > 100 {
			usagePercent = 100
		}
	}

	data := gin.H{
		"Title":          "Account Settings",
		"User":           user,
		"ApiKeys":        apiKeysVM,
		"Pagination":     paginationVM,
		"AvailablePerms": enums.AvailablePerms,
		"ActiveMenu":     "account",
		"UsageCount":     usageCount,
		"DailyLimit":     user.DailyLimit,
		"UsagePercent":   usagePercent,
	}

	for k, v := range extra {
		data[k] = v
	}

	renderTemplate(c, http.StatusOK, "account.html", data)
}

// HandleChangePassword processes password change requests
func (h *ApiKeyWebHandler) HandleChangePassword(c *gin.Context) {
	userId, exists := c.Get("userId")
	if !exists {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	currentPassword := c.PostForm("currentPassword")
	newPassword := c.PostForm("password")
	confirmPassword := c.PostForm("confirmPassword")

	if currentPassword == "" || newPassword == "" || confirmPassword == "" {
		renderFragment(c, http.StatusOK, "password_error_alert.html", gin.H{
			"Error": "All fields are required",
		})
		return
	}

	if newPassword != confirmPassword {
		renderFragment(c, http.StatusOK, "password_error_alert.html", gin.H{
			"Error": "Passwords do not match",
		})
		return
	}

	if err := validatePassword(newPassword); err != nil {
		renderFragment(c, http.StatusOK, "password_error_alert.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	user, err := h.userService.GetByID(c, userId)
	if err != nil {
		renderFragment(c, http.StatusOK, "password_error_alert.html", gin.H{
			"Error": "User not found",
		})
		return
	}

	// Verify current password
	ok, _ := argon2id.ComparePasswordAndHash(currentPassword, user.Password)
	if !ok {
		renderFragment(c, http.StatusOK, "password_error_alert.html", gin.H{
			"Error": "Incorrect current password",
		})
		return
	}

	userInput := models.UserInput{
		Username: user.Username,
		Email:    user.Email,
		Password: newPassword,
	}

	err = h.userService.Update(c, user.ID, userInput)
	if err != nil {
		renderFragment(c, http.StatusOK, "password_error_alert.html", gin.H{
			"Error": "Failed to update password: " + err.Error(),
		})
		return
	}

	// Render success alert
	renderFragment(c, http.StatusOK, "password_success_alert.html", gin.H{
		"Message": "Password updated successfully",
	})
}
