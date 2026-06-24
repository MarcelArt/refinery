package routes

import (
	"github.com/MarcelArt/refinery/internal/web/handlers"
	"github.com/gin-gonic/gin"
)

func setupAuthRoutes(r *gin.Engine, authM *WebAuthMiddleware, h *handlers.AuthWebHandler) {
	r.GET("/login", authM.RedirectIfAuthenticated(), h.ShowLogin)
	r.POST("/login", authM.RedirectIfAuthenticated(), h.HandleLogin)
	
	r.GET("/register", authM.RedirectIfAuthenticated(), h.ShowRegister)
	r.POST("/register", authM.RedirectIfAuthenticated(), h.HandleRegister)
	
	r.GET("/logout", h.HandleLogout)
}
