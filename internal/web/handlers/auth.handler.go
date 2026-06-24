package handlers

import (
	"net/http"

	"github.com/MarcelArt/refinery/internal/v1/models"
	"github.com/MarcelArt/refinery/internal/v1/services"
	"github.com/gin-gonic/gin"
)

type AuthWebHandler struct {
	userService services.IUserService
}

func NewAuthWebHandler(userService services.IUserService) *AuthWebHandler {
	return &AuthWebHandler{
		userService: userService,
	}
}

// ShowLogin renders the login page
func (h *AuthWebHandler) ShowLogin(c *gin.Context) {
	renderTemplate(c, http.StatusOK, "login.html", gin.H{
		"Title":      "Sign In",
		"HideLayout": true,
	})
}

// HandleLogin processes the login form submission
func (h *AuthWebHandler) HandleLogin(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	isRemember := c.PostForm("isRemember") == "on"

	if username == "" || password == "" {
		renderFragment(c, http.StatusOK, "error_alert.html", gin.H{
			"Error": "Username and password are required",
		})
		return
	}

	input := models.LoginInput{
		Username:   username,
		Password:   password,
		IsRemember: isRemember,
	}

	_, err := h.userService.Login(c, input)
	if err != nil {
		renderFragment(c, http.StatusOK, "error_alert.html", gin.H{
			"Error": "Invalid username/email or password",
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

// ShowRegister renders the registration page
func (h *AuthWebHandler) ShowRegister(c *gin.Context) {
	renderTemplate(c, http.StatusOK, "register.html", gin.H{
		"Title":      "Create Account",
		"HideLayout": true,
	})
}

// HandleRegister processes the registration form submission
func (h *AuthWebHandler) HandleRegister(c *gin.Context) {
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirmPassword")

	if username == "" || email == "" || password == "" {
		renderFragment(c, http.StatusOK, "error_alert.html", gin.H{
			"Error": "All fields are required",
		})
		return
	}

	if password != confirmPassword {
		renderFragment(c, http.StatusOK, "error_alert.html", gin.H{
			"Error": "Passwords do not match",
		})
		return
	}

	userInput := models.UserInput{
		Username: username,
		Email:    email,
		Password: password,
	}

	_, err := h.userService.Create(c, userInput)
	if err != nil {
		renderFragment(c, http.StatusOK, "error_alert.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	// Auto-login after successful registration
	loginInput := models.LoginInput{
		Username:   username,
		Password:   password,
		IsRemember: false,
	}
	_, err = h.userService.Login(c, loginInput)
	if err != nil {
		// If auto-login fails, redirect them to login page to sign in manually
		if c.GetHeader("HX-Request") == "true" {
			c.Header("HX-Redirect", "/login")
			c.Status(http.StatusOK)
		} else {
			c.Redirect(http.StatusSeeOther, "/login")
		}
		return
	}

	if c.GetHeader("HX-Request") == "true" {
		c.Header("HX-Redirect", "/workflows")
		c.Status(http.StatusOK)
	} else {
		c.Redirect(http.StatusSeeOther, "/workflows")
	}
}

// HandleLogout clears authentication cookies and redirects to login
func (h *AuthWebHandler) HandleLogout(c *gin.Context) {
	// Clear access token cookie
	c.SetCookie("at", "", -1, "/", "", false, true)
	// Clear refresh token cookie
	c.SetCookie("rt", "", -1, "/", "", false, true)
	
	c.Redirect(http.StatusSeeOther, "/login")
}
