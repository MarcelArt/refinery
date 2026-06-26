package middlewares

import (
	"errors"
	"net/http"

	"git.bangmarcel.art/marcel/arrays"
	"github.com/MarcelArt/refinery/internal/common"
	"github.com/MarcelArt/refinery/internal/configs"
	"github.com/MarcelArt/refinery/internal/enums"
	"github.com/MarcelArt/refinery/internal/v1/services"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type AuthMiddleware struct {
	jwtSecret []byte
	akService services.IApiKeyService
}

func NewAuthMiddleware(akService services.IApiKeyService) *AuthMiddleware {
	return &AuthMiddleware{
		jwtSecret: []byte(configs.Env.JwtSecret),
		akService: akService,
	}
}

func (m *AuthMiddleware) Authn(c *gin.Context) {
	// X-Api-Key
	apiKey := c.GetHeader("X-Api-Key")
	if apiKey != "" {
		apiKey, err := m.akService.GetByKey(c, apiKey)
		if err != nil {
			_, res := common.ResultErr(errors.New("invalid pat"), "")
			c.JSON(http.StatusUnauthorized, res)
			c.Abort()
			return
		}

		permissions, err := apiKey.Scopes.Deserialize()
		if err != nil {
			_, res := common.ResultErr(errors.New("broken pat"), "")
			c.JSON(http.StatusUnauthorized, res)
			c.Abort()
			return
		}

		c.Set("userId", float64(apiKey.UserID))
		c.Set("permissions", permissions)

		c.Next()
		return
	}

	// Bearer Auth
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" {
		_, res := common.ResultErr(errors.New("missing authorization token"), "")
		c.JSON(http.StatusUnauthorized, res)
		c.Abort()
		return
	}

	// Remove "Bearer " prefix if present
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		tokenString = tokenString[7:]
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		return m.jwtSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil || !token.Valid {
		_, res := common.ResultErr(errors.New("invalid token"), "")
		c.JSON(http.StatusUnauthorized, res)
		c.Abort()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		c.Set("userId", claims["userId"])
		c.Set("permissions", claims["permissions"])
		c.Next()
	}
}

func (m *AuthMiddleware) Refresh(c *gin.Context) {
	refreshToken := c.GetHeader("X-Refresh-Token")
	if refreshToken == "" {
		_, res := common.ResultErr(errors.New("missing refresh token"), "")
		c.JSON(http.StatusUnauthorized, res)
		c.Abort()
		return
	}

	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (any, error) {
		return m.jwtSecret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

	if err != nil || !token.Valid {
		_, res := common.ResultErr(errors.New("invalid refresh token"), "")
		c.JSON(http.StatusUnauthorized, res)
		c.Abort()
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		c.Set("userId", claims["userId"])
		c.Set("isRemember", claims["isRemember"])
		c.Next()
	}
}

func (m *AuthMiddleware) WebhookAuth(c *gin.Context) {
	webhookKey := c.GetHeader("X-Webhook-Key")
	if webhookKey != configs.Env.JwtSecret {
		_, res := common.ResultErr(errors.New("invalid webhook key"), "")
		c.JSON(http.StatusUnauthorized, res)
		c.Abort()
		return
	}

	c.Next()
}

func (m *AuthMiddleware) Authz(permissionKey string) gin.HandlerFunc {
	return func(c *gin.Context) {
		cPerms, ok := c.Get("permissions")
		if !ok {
			_, res := common.ResultErr(common.ErrNotStringSlice, "")
			c.JSON(http.StatusUnauthorized, res)
			c.Abort()
			return
		}

		permissions, err := common.ParseClaimsToStringSlice(cPerms)
		if err != nil {
			_, res := common.ResultErr(err, "failed parsing permission claims")
			c.JSON(http.StatusUnauthorized, res)
			c.Abort()
			return
		}

		permission := arrays.Find(permissions, func(p string) bool {
			return p == enums.PermFullAccess || p == permissionKey
		})

		if permission == nil {
			_, res := common.ResultErr(errors.New("user doesn't have required permission"), "unauthorized")
			c.JSON(http.StatusForbidden, res)
			c.Abort()
			return
		}

		c.Next()
	}
}
