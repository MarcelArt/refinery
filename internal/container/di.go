package container

import (
	"github.com/MarcelArt/refinery/internal/app"
	"github.com/MarcelArt/refinery/internal/common"
	"github.com/MarcelArt/refinery/internal/configs"
	"github.com/MarcelArt/refinery/internal/v1/handlers"
	"github.com/MarcelArt/refinery/internal/v1/middlewares"
	"github.com/MarcelArt/refinery/internal/v1/repositories"
	"github.com/MarcelArt/refinery/internal/v1/services"
	webhandlers "github.com/MarcelArt/refinery/internal/web/handlers"
	webroutes "github.com/MarcelArt/refinery/internal/web/routes"
	"go.uber.org/dig"
)

func New() *dig.Container {
	c := dig.New()

	c.Provide(configs.ConnectDB)
	c.Provide(configs.ConnectR2)

	c.Provide(repositories.NewN8NRepo, dig.As(new(repositories.IN8NRepo)))
	c.Provide(repositories.NewUserRepo, dig.As(new(repositories.IUserRepo)))
	c.Provide(repositories.NewWorkflowRepo, dig.As(new(repositories.IWorkflowRepo)))
	c.Provide(repositories.NewExtractionResultRepo, dig.As(new(repositories.IExtractionResultRepo)))
	c.Provide(repositories.NewApiKeyRepo, dig.As(new(repositories.IApiKeyRepo)))
	c.Provide(repositories.NewWebhookRepo, dig.As(new(repositories.IWebhookRepo)))
	c.Provide(repositories.NewRateLimiterRepo, dig.As(new(repositories.IRateLimiterRepo)))
	c.Provide(repositories.NewR2Repo, dig.As(new(common.IS3Repo)))
	c.Provide(repositories.NewMailRepo)

	c.Provide(services.NewN8NService, dig.As(new(services.IN8NService)))
	c.Provide(services.NewUserService, dig.As(new(services.IUserService)))
	c.Provide(services.NewWorkflowService, dig.As(new(services.IWorkflowService)))
	c.Provide(services.NewExtractionResultService, dig.As(new(services.IExtractionResultService)))
	c.Provide(services.NewApiKeyService, dig.As(new(services.IApiKeyService)))
	c.Provide(services.NewWebhookService, dig.As(new(services.IWebhookService)))
	c.Provide(services.NewDashboardService)
	c.Provide(services.NewRateLimiterService, dig.As(new(services.IRateLimiterService)))

	c.Provide(middlewares.NewAuthMiddleware)
	c.Provide(middlewares.NewRateLimiterMiddleware)

	c.Provide(handlers.NewUserHandler)
	c.Provide(handlers.NewWorkflowHandler)
	c.Provide(handlers.NewExtractionResultHandler)
	c.Provide(handlers.NewApiKeyHandler)
	c.Provide(handlers.NewWebhookHandler)
	c.Provide(handlers.NewDashboardHandler)
	c.Provide(handlers.NewRateLimiterHandler)

	// Web components
	c.Provide(webroutes.NewWebAuthMiddleware)
	c.Provide(webhandlers.NewAuthWebHandler)
	c.Provide(webhandlers.NewWorkflowWebHandler)
	c.Provide(webhandlers.NewExtractionResultWebHandler)
	c.Provide(webhandlers.NewApiKeyWebHandler)
	c.Provide(webhandlers.NewWebhookWebHandler)
	c.Provide(webhandlers.NewDashboardWebHandler)

	c.Provide(app.New)

	return c
}
