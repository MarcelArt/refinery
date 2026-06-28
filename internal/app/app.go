package app

import (
	"fmt"
	"log"
	"time"

	"github.com/MarcelArt/refinery/internal/configs"
	"github.com/MarcelArt/refinery/internal/v1/handlers"
	"github.com/MarcelArt/refinery/internal/v1/middlewares"
	"github.com/MarcelArt/refinery/internal/v1/routes"
	webhandlers "github.com/MarcelArt/refinery/internal/web/handlers"
	webroutes "github.com/MarcelArt/refinery/internal/web/routes"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type App struct {
	uHandler  *handlers.UserHandler
	wHandler  *handlers.WorkflowHandler
	erHandler *handlers.ExtractionResultHandler
	akHandler *handlers.ApiKeyHandler
	whHandler *handlers.WebhookHandler
	dHandler  *handlers.DashboardHandler
	authM     *middlewares.AuthMiddleware
	webAuthM  *webroutes.WebAuthMiddleware
	authWebH  *webhandlers.AuthWebHandler
	wfWebH    *webhandlers.WorkflowWebHandler
	erWebH    *webhandlers.ExtractionResultWebHandler
	akWebH    *webhandlers.ApiKeyWebHandler
	whWebH    *webhandlers.WebhookWebHandler
}

func New(
	uHandler *handlers.UserHandler,
	wHandler *handlers.WorkflowHandler,
	erHandler *handlers.ExtractionResultHandler,
	akHandler *handlers.ApiKeyHandler,
	whHandler *handlers.WebhookHandler,
	dHandler *handlers.DashboardHandler,
	authM *middlewares.AuthMiddleware,
	webAuthM *webroutes.WebAuthMiddleware,
	authWebH *webhandlers.AuthWebHandler,
	wfWebH *webhandlers.WorkflowWebHandler,
	erWebH *webhandlers.ExtractionResultWebHandler,
	akWebH *webhandlers.ApiKeyWebHandler,
	whWebH *webhandlers.WebhookWebHandler,
) *App {
	return &App{
		uHandler:  uHandler,
		wHandler:  wHandler,
		erHandler: erHandler,
		akHandler: akHandler,
		whHandler: whHandler,
		dHandler:  dHandler,
		authM:     authM,
		webAuthM:  webAuthM,
		authWebH:  authWebH,
		wfWebH:    wfWebH,
		erWebH:    erWebH,
		akWebH:    akWebH,
		whWebH:    whWebH,
	}
}

func (a *App) Run() error {
	if configs.Env.ServerENV == "prod" {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"POST, OPTIONS, GET, PUT, PATCH, DELETE"},
		AllowHeaders:     []string{"Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With", "X-Refresh-Token"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	api := r.Group("/api")
	routes.SetupRoutes(api, a.uHandler, a.wHandler, a.erHandler, a.akHandler, a.authM, a.whHandler, a.dHandler)

	webroutes.SetupWebRoutes(r, a.webAuthM, a.authWebH, a.wfWebH, a.erWebH, a.akWebH, a.whWebH)

	port := fmt.Sprintf(":%s", configs.Env.PORT)
	log.Printf("Listening on http://localhost%s", port)
	log.Printf("Open swagger doc on http://localhost%s/swagger/index.html", port)
	return r.Run(port)
}
