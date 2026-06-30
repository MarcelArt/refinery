package handlers

import (
	"fmt"
	"net/http"

	"github.com/MarcelArt/refinery/internal/common"
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

// ShowLanding renders the marketing landing page
func (h *AuthWebHandler) ShowLanding(c *gin.Context) {
	var isLoggedIn bool
	at, err := c.Cookie("at")
	if err == nil && at != "" {
		claims, err := common.ParseToken(at)
		if err == nil && claims["userId"] != nil {
			isLoggedIn = true
		}
	}
	if !isLoggedIn {
		rt, err := c.Cookie("rt")
		if err == nil && rt != "" {
			claims, err := common.ParseToken(rt)
			if err == nil && claims["userId"] != nil {
				isLoggedIn = true
			}
		}
	}

	renderTemplate(c, http.StatusOK, "landing.html", gin.H{
		"Title":      "AI-Powered Document Parsing",
		"HideLayout": true,
		"IsLoggedIn": isLoggedIn,
	})
}

// ShowLogin renders the login page
func (h *AuthWebHandler) ShowLogin(c *gin.Context) {
	var infoMsg string
	var errMsg string

	if c.Query("verified") == "true" {
		infoMsg = "Email verified successfully! You can now sign in."
	}
	if errParam := c.Query("error"); errParam != "" {
		errMsg = errParam
	}

	renderTemplate(c, http.StatusOK, "login.html", gin.H{
		"Title":      "Sign In",
		"HideLayout": true,
		"Info":       infoMsg,
		"Error":      errMsg,
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
			"Error": err.Error(),
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

	if err := validatePassword(password); err != nil {
		renderFragment(c, http.StatusOK, "error_alert.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	userInput := models.UserInput{
		Username: username,
		Email:    email,
		Password: password,
	}

	userID, err := h.userService.Create(c, userInput)
	if err != nil {
		renderFragment(c, http.StatusOK, "error_alert.html", gin.H{
			"Error": err.Error(),
		})
		return
	}

	// Send verification email
	_ = h.userService.SendEmailVerification(c, userID, userInput)

	// Display verification instructions page
	if c.GetHeader("HX-Request") == "true" {
		c.Header("HX-Retarget", ".auth-card")
		c.Header("HX-Reswap", "innerHTML")
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(`
			<div class="auth-header">
				<div style="display: flex; justify-content: center; margin-bottom: 16px;">
					<div class="logo">
						<span class="logo-icon">R</span>
						<span>Refinery</span>
					</div>
				</div>
				<h1 class="auth-title">Verify Your Email</h1>
				<p class="auth-subtitle">Check your inbox to activate your account</p>
			</div>

			<div style="background-color: rgba(16, 185, 129, 0.08); border: 1px solid rgba(16, 185, 129, 0.2); border-radius: 6px; padding: 20px; margin-bottom: 24px; color: #10b981; font-size: 13.5px; line-height: 1.6; display: flex; flex-direction: column; align-items: center; text-align: center; gap: 12px; width: 100%; box-sizing: border-box;">
				<div style="width: 44px; height: 44px; border-radius: 50%; background-color: rgba(16, 185, 129, 0.15); display: flex; align-items: center; justify-content: center; color: #10b981; flex-shrink: 0;">
					<svg xmlns="http://www.w3.org/2000/svg" width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round"><path d="M22 11.08V12a10 10 0 1 1-5.93-9.14"/><polyline points="22 4 12 14.01 9 11.01"/></svg>
				</div>
				<div>
					<strong style="font-size: 15px; display: block; margin-bottom: 6px;">Registration Successful!</strong>
					We have sent a verification link to your email address. Please click the link to activate your account.
				</div>
			</div>

			<div class="auth-footer" style="margin-top: 16px;">
				Already verified? <a href="/login" style="color: var(--color-amber); text-decoration: none;">Return to Sign In</a>
			</div>
		`))
	} else {
		c.Redirect(http.StatusSeeOther, "/login?error=Please+check+your+email+to+verify+your+account")
	}
}

// HandleVerifyEmail processes the email verification link
func (h *AuthWebHandler) HandleVerifyEmail(c *gin.Context) {
	tokenStr := c.Query("t")
	if tokenStr == "" {
		c.Redirect(http.StatusSeeOther, "/login?error=Verification+token+is+missing")
		return
	}

	claims, err := common.ParseToken(tokenStr)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login?error=Verification+link+is+invalid+or+expired")
		return
	}

	userIdVal, ok := claims["userId"]
	if !ok {
		c.Redirect(http.StatusSeeOther, "/login?error=Invalid+token+claims")
		return
	}

	var userId uint
	switch v := userIdVal.(type) {
	case float64:
		userId = uint(v)
	case int:
		userId = uint(v)
	case uint:
		userId = v
	default:
		c.Redirect(http.StatusSeeOther, "/login?error=Invalid+user+ID+format")
		return
	}

	err = h.userService.Verify(c, userId)
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login?error=Failed+to+verify+user")
		return
	}

	c.Redirect(http.StatusSeeOther, "/login?verified=true")
}

// HandleLogout clears authentication cookies and redirects to login
func (h *AuthWebHandler) HandleLogout(c *gin.Context) {
	// Clear access token cookie
	c.SetCookie("at", "", -1, "/", "", false, true)
	// Clear refresh token cookie
	c.SetCookie("rt", "", -1, "/", "", false, true)
	
	c.Redirect(http.StatusSeeOther, "/login")
}

func validatePassword(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}
	var hasUpper, hasDigit, hasSpecial bool
	for _, char := range password {
		if char >= 'A' && char <= 'Z' {
			hasUpper = true
		} else if char >= '0' && char <= '9' {
			hasDigit = true
		} else if char < 'a' || char > 'z' {
			// not lowercase a-z, not uppercase A-Z, not digit 0-9
			hasSpecial = true
		}
	}
	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasDigit {
		return fmt.Errorf("password must contain at least one numeric character")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}
	return nil
}
