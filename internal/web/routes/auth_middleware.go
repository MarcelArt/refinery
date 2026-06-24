package routes

import (
	"net/http"

	"github.com/MarcelArt/refinery/internal/common"
	"github.com/MarcelArt/refinery/internal/configs"
	"github.com/MarcelArt/refinery/internal/v1/services"
	"github.com/gin-gonic/gin"
)

type WebAuthMiddleware struct {
	userService services.IUserService
	jwtSecret   []byte
}

func NewWebAuthMiddleware(userService services.IUserService) *WebAuthMiddleware {
	return &WebAuthMiddleware{
		userService: userService,
		jwtSecret:   []byte(configs.Env.JwtSecret),
	}
}

// RequireAuth ensures the user is authenticated (via at cookie or refreshed rt cookie)
func (m *WebAuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		at, err := c.Cookie("at")
		if err == nil && at != "" {
			// Parse and validate at
			claims, err := common.ParseToken(at)
			if err == nil {
				// Set userId in context
				c.Set("userId", claims["userId"])
				c.Next()
				return
			}
		}

		// If access token is invalid or missing, check refresh token
		rt, err := c.Cookie("rt")
		if err == nil && rt != "" {
			claims, err := common.ParseToken(rt)
			if err == nil {
				userId := claims["userId"]
				isRemember := false
				if ir, ok := claims["isRemember"].(bool); ok {
					isRemember = ir
				}

				// Try to regenerate tokens
				_, err := m.userService.RegenerateTokenPair(c, userId, isRemember)
				if err == nil {
					c.Set("userId", userId)
					c.Next()
					return
				}
			}
		}

		// Not authenticated, clear cookies and redirect to /login
		c.SetCookie("at", "", -1, "/", "", false, true)
		c.SetCookie("rt", "", -1, "/", "", false, true)
		c.Redirect(http.StatusSeeOther, "/login")
		c.Abort()
	}
}

// RedirectIfAuthenticated redirects the user to /dashboard if they are already logged in
func (m *WebAuthMiddleware) RedirectIfAuthenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		at, err := c.Cookie("at")
		if err == nil && at != "" {
			_, err := common.ParseToken(at)
			if err == nil {
				c.Redirect(http.StatusSeeOther, "/dashboard")
				c.Abort()
				return
			}
		}

		rt, err := c.Cookie("rt")
		if err == nil && rt != "" {
			claims, err := common.ParseToken(rt)
			if err == nil {
				userId := claims["userId"]
				isRemember := false
				if ir, ok := claims["isRemember"].(bool); ok {
					isRemember = ir
				}
				_, err := m.userService.RegenerateTokenPair(c, userId, isRemember)
				if err == nil {
					c.Redirect(http.StatusSeeOther, "/dashboard")
					c.Abort()
					return
				}
			}
		}

		c.Next()
	}
}
