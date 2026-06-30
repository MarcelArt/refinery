package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/MarcelArt/refinery/internal/common"
	"github.com/MarcelArt/refinery/internal/configs"
	"github.com/MarcelArt/refinery/internal/entities"
	"github.com/MarcelArt/refinery/internal/enums"
	"github.com/MarcelArt/refinery/internal/v1/models"
	"github.com/MarcelArt/refinery/internal/v1/repositories"
	"github.com/alexedwards/argon2id"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/morkid/paginate"
)

type IUserService interface {
	common.IBaseCrudService[entities.User, models.UserInput, models.UserPage]
	Login(c *gin.Context, input models.LoginInput) (models.LoginResponse, error)
	RegenerateTokenPair(c *gin.Context, userID any, isRemember bool) (models.LoginResponse, error)
	// AssignRoles(c context.Context, userID uint, newRoleIDs []uint) error
	GetPermissions(userID any) ([]string, error)
	GetRoles(id any) ([]models.UserRole, error)
	SendEmailVerification(c *gin.Context, id uint, input models.UserInput) error
	Verify(c context.Context, id any) error
}

type UserService struct {
	repo  repositories.IUserRepo
	mRepo *repositories.MailRepo
	// urRepo repositories.IUserRoleRepo
}

var _ IUserService = &UserService{}

func NewUserService(
	repo repositories.IUserRepo,
	mRepo *repositories.MailRepo,
	// urRepo repositories.IUserRoleRepo,
) *UserService {
	return &UserService{
		repo:  repo,
		mRepo: mRepo,
		// urRepo: urRepo,
	}
}

func (s *UserService) Create(c context.Context, input models.UserInput) (uint, error) {
	password, err := argon2id.CreateHash(input.Password, argon2id.DefaultParams)
	if err != nil {
		return 0, fmt.Errorf("failed to hash password: %w", err)
	}
	input.Password = password
	input.VerifiedAt = nil

	return s.repo.Create(c, input)
}

func (s *UserService) Read(c *gin.Context) (paginate.Page, []models.UserPage) {
	return s.repo.Read(c)
}

func (s *UserService) Update(c context.Context, id any, input models.UserInput) error {
	if input.Password != "" {
		password, err := argon2id.CreateHash(input.Password, argon2id.DefaultParams)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}
		input.Password = password
	}

	return s.repo.Update(c, id, input)
}

func (s *UserService) Delete(c context.Context, id any) error {
	return s.repo.Delete(c, id)
}

func (s *UserService) GetByID(c context.Context, id any) (entities.User, error) {
	return s.repo.GetByID(c, id)
}

func (s *UserService) Login(c *gin.Context, input models.LoginInput) (models.LoginResponse, error) {
	var res models.LoginResponse
	user, err := s.repo.GetByUsernameOrEmail(c, input.Username)
	if err != nil {
		return res, errors.New("invalid username/email or password")
	}

	if ok, _ := argon2id.ComparePasswordAndHash(input.Password, user.Password); !ok {
		return res, errors.New("invalid username/email or password")
	}

	if user.VerifiedAt == nil {
		return res, errors.New("please verify your email address before signing in")
	}

	res, err = s.GenerateTokenPair(user, input.IsRemember, c.Request.Host)
	if err != nil {
		return res, fmt.Errorf("failed to generate token pair: %w", err)
	}
	common.GenerateCookies(c, res.AccessToken, res.RefreshToken, input.IsRemember)

	return res, nil
}

func (s *UserService) RegenerateTokenPair(c *gin.Context, userID any, isRemember bool) (models.LoginResponse, error) {
	var res models.LoginResponse
	user, err := s.repo.GetByID(c, userID)
	if err != nil {
		return res, err
	}

	res, err = s.GenerateTokenPair(user, isRemember, c.Request.Host)
	if err != nil {
		return res, fmt.Errorf("failed to generate token pair: %w", err)
	}
	common.GenerateCookies(c, res.AccessToken, res.RefreshToken, isRemember)

	return res, nil
}

// func (s *UserService) AssignRoles(c context.Context, userID uint, newRoleIDs []uint) error {
// 	tx := s.db.Begin()
// 	defer tx.Rollback()

// 	usecase := usecases.InitAssignRolesUsecase(tx)
// 	usecase.UserID = userID
// 	usecase.NewRoleIDs = newRoleIDs

// 	if err := usecase.Execute(c); err != nil {
// 		return err
// 	}

// 	return tx.Commit().Error
// }

func (s *UserService) GetPermissions(userID any) ([]string, error) {
	return s.repo.GetPermissions(userID)
}

func (s *UserService) GenerateTokenPair(user entities.User, isRemember bool, iss string) (models.LoginResponse, error) {
	var res models.LoginResponse
	claims := map[string]any{
		"sub":    user.Username,
		"userId": user.ID,
		"iss":    iss,
		"aud":    iss,
	}

	// permissions, err := s.repo.GetPermissions(user.ID)
	// if err != nil {
	// 	return res, fmt.Errorf("failed retrieving permissions: %w", err)
	// }

	permissions := []string{enums.PermFullAccess}

	at, rt, err := common.GenerateJWTPair(claims, permissions, isRemember)
	if err != nil {
		return res, fmt.Errorf("failed generating token pair: %w", err)
	}

	res.AccessToken = at
	res.RefreshToken = rt
	res.User = user

	return res, nil
}

func (s *UserService) GetRoles(id any) ([]models.UserRole, error) {
	return s.repo.GetRoles(id)
}

func (s *UserService) SendEmailVerification(c *gin.Context, id uint, input models.UserInput) error {
	body := `
		<!DOCTYPE html>
		<html lang="en">
		<head>
			<meta charset="UTF-8">
			<meta name="viewport" content="width=device-width, initial-scale=1.0">
			<title>Verify Your Email</title>
			<style>
				/* Embedded Web Fonts mimicking var(--font-sans) and var(--font-display) */
				@import url('https://fonts.googleapis.com/css2?family=Outfit:wght@600;700&family=Plus+Jakarta+Sans:wght@400;500;600&display=swap');

				body {
					font-family: 'Plus Jakarta Sans', -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif;
					background-color: #090d16; /* --bg-dark */
					margin: 0;
					padding: 0;
					-webkit-font-smoothing: antialiased;
				}
				.email-wrapper {
					background-color: #090d16;
					padding: 40px 20px;
				}
				.email-container {
					max-width: 440px;
					margin: 0 auto;
					background-color: #101626; /* --bg-card */
					border-radius: 6px;
					overflow: hidden;
					border: 1px solid #202b48; /* --border-color */
					box-shadow: 0 8px 24px -8px rgba(0, 0, 0, 0.5);
					position: relative;
				}
				/* Top border decorative gradient line matching your .auth-card::before style */
				.top-accent-line {
					height: 3px;
					background: linear-gradient(90deg, #10b981, #0ea5e9); /* --accent to --accent-hover */
				}
				.content {
					padding: 32px;
				}
				.logo-container {
					margin-bottom: 24px;
				}
				.logo-icon {
					display: inline-block;
					width: 24px;
					height: 24px;
					background: linear-gradient(135deg, #10b981, #0ea5e9);
					border-radius: 4px;
					text-align: center;
					line-height: 24px;
					font-family: 'Outfit', sans-serif;
					font-weight: 800;
					font-size: 13px;
					color: #070a13;
				}
				h1 {
					font-family: 'Outfit', sans-serif;
					color: #f8fafc; /* --text-primary */
					font-size: 26px;
					font-weight: 700;
					margin: 0 0 8px 0;
					letter-spacing: -0.02em;
				}
				.subtitle {
					color: #94a3b8; /* --text-secondary */
					font-size: 14px;
					margin: 0 0 24px 0;
				}
				p {
					color: #94a3b8; /* --text-secondary */
					font-size: 14px;
					line-height: 1.6;
					margin: 0 0 20px 0;
				}
				.button-container {
					margin: 28px 0;
					text-align: center;
				}
				/* Button styled precisely like your .btn-primary */
				.btn-primary {
					background-color: #10b981; /* --accent */
					color: #070a13 !important;
					text-decoration: none;
					padding: 10px 20px;
					border-radius: 4px;
					font-size: 13px;
					font-weight: 600;
					display: inline-block;
				}
				.security-notice {
					color: #64748b; /* --text-muted */
					font-size: 12px;
					margin-top: 24px;
				}
				.divider {
					border: 0;
					border-top: 1px solid #202b48; /* --border-color */
					margin: 24px 0;
				}
				.fallback-label {
					font-size: 11px;
					text-transform: uppercase;
					letter-spacing: 0.05em;
					color: #64748b; /* --text-muted */
					margin-bottom: 6px;
				}
				.fallback-box {
					font-family: monospace;
					font-size: 11px;
					color: #94a3b8; /* --text-secondary */
					background-color: #182035; /* --bg-input */
					padding: 8px 12px;
					border-radius: 4px;
					border: 1px solid #202b48; /* --border-color */
					word-break: break-all;
				}
				.fallback-box a {
					color: #10b981; /* --accent */
					text-decoration: none;
				}
				.footer {
					text-align: center;
					margin-top: 24px;
					font-size: 12px;
					color: #64748b; /* --text-muted */
				}
			</style>
		</head>
		<body>

			<div class="email-wrapper">
				<div class="email-container">
					<!-- Decorative Glow Line -->
					<div class="top-accent-line"></div>

					<div class="content">
						<!-- Brand Token Identifier -->
						<div class="logo-container">
							<div class="logo-icon">R</div>
						</div>

						<h1>Verify identity</h1>
						<div class="subtitle">Complete your account registration initialization.</div>

						<p>Hello, {{.USERNAME}}</p>
						<p>An account creation request linked to this email address was received. To secure your profile data structure and activate system permissions, please complete initialization below:</p>

						<!-- CTA Button -->
						<div class="button-container">
							<a href="{{.VERIFICATION_URL}}" class="btn-primary" target="_blank">VERIFY ACCOUNT</a>
						</div>

						<p class="security-notice">⚠️ This dynamic secure verification token is valid for 24 hours. If you did not initiate this system command, please ignore this communication.</p>

						<div class="divider"></div>

						<!-- Code/Link Fallback block styled to look like your platform's workflow prompts -->
						<div class="fallback-label">Fallback Token URL</div>
						<div class="fallback-box">
							<a href="{{.VERIFICATION_URL}}" target="_blank">{{.VERIFICATION_URL}}</a>
						</div>
					</div>
				</div>

				<!-- Footer -->
				<div class="footer">
					<p>&copy; 2026 Refinery.</p>
				</div>
			</div>

		</body>
		</html>
	`

	today := time.Now()
	tokenClaims := jwt.MapClaims{
		"iat":    today.Unix(),
		"exp":    today.Add(24 * time.Hour).Unix(),
		"sub":    input.Username,
		"userId": id,
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims).SignedString([]byte(configs.Env.JwtSecret))
	if err != nil {
		return fmt.Errorf("failed generating token: %w", err)
	}

	scheme := "https"
	if configs.Env.ServerENV == "local" {
		scheme = "http"
	}

	url := fmt.Sprintf("%s://%s/verify?t=%s", scheme, c.Request.Host, token)

	data := map[string]string{
		"USERNAME":         input.Username,
		"VERIFICATION_URL": url,
	}

	body, err = common.TextTemplating(body, data)
	if err != nil {
		return fmt.Errorf("failed composing email: %w", err)
	}

	mail := models.Mailer{
		To:      []string{input.Email},
		Subject: "🔒 Verify your account to get started!",
		Body:    body,
	}

	return s.mRepo.SendMail(mail)
}

func (s *UserService) Verify(c context.Context, id any) error {
	input := models.UserInput{
		VerifiedAt: new(time.Now()),
	}
	return s.repo.Update(c, id, input)
}
