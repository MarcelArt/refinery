package middlewares

import (
	"errors"
	"net/http"

	"github.com/MarcelArt/refinery/internal/common"
	"github.com/MarcelArt/refinery/internal/v1/models"
	"github.com/MarcelArt/refinery/internal/v1/services"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RateLimiterMiddleware struct {
	service  services.IRateLimiterService
	uService services.IUserService
}

func NewRateLimiterMiddleware(service services.IRateLimiterService, uService services.IUserService) *RateLimiterMiddleware {
	return &RateLimiterMiddleware{
		service:  service,
		uService: uService,
	}
}

func (m *RateLimiterMiddleware) Limit(c *gin.Context) {
	userId := c.MustGet("userId")

	user, err := m.uService.GetByID(c.Request.Context(), userId)
	if err != nil {
		_, res := common.ResultErr(err, "unknown user")
		c.JSON(http.StatusUnauthorized, res)
		c.Abort()
		return
	}

	rateLimit, err := m.service.GetTodayByUserID(c, userId)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(common.ResultErr(err, "failed checking rate limit"))
		c.Abort()
		return
	}

	if rateLimit.Count > user.DailyLimit {
		_, res := common.ResultErr(errors.New("you have exceeded your daily quota"), "")
		c.JSON(http.StatusTooManyRequests, res)
		c.Abort()
		return
	}

	rateLimit.Count++
	input := models.RateLimiterInput{
		Count:  rateLimit.Count,
		UserID: user.ID,
	}

	if rateLimit.ID != 0 {
		if err := m.service.Update(c, rateLimit.ID, input); err != nil {
			c.JSON(common.ResultErr(err, "failed updating rate limit"))
			c.Abort()
			return
		}
	} else {
		if _, err := m.service.Create(c, input); err != nil {
			c.JSON(common.ResultErr(err, "failed creating rate limit"))
			c.Abort()
			return
		}
	}

	c.Next()
}
