package common

import (
	"net/http"
	"time"

	"github.com/MarcelArt/refinery/internal/configs"
	"github.com/MarcelArt/refinery/internal/enums"
	"github.com/gin-gonic/gin"
)

func GenerateCookies(c *gin.Context, at string, rt string, isRemember bool) {
	atExp := time.Minute * 5
	rtExp := enums.Day
	if isRemember {
		rtExp = enums.Month
	}

	isProd := configs.Env.ServerENV == "prod"

	c.SetCookieData(&http.Cookie{
		Name:     "at",
		Value:    at,
		HttpOnly: true,
		Secure:   isProd,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(atExp),
		Path:     "/",
	})
	c.SetCookieData(&http.Cookie{
		Name:     "rt",
		Value:    rt,
		HttpOnly: true,
		Secure:   isProd,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   int(rtExp),
		Path:     "/",
	})
}
