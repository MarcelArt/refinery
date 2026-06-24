package container

import (
	"github.com/MarcelArt/refinery/internal/app"
	"github.com/MarcelArt/refinery/internal/configs"
	"github.com/MarcelArt/refinery/internal/v1/handlers"
	"github.com/MarcelArt/refinery/internal/v1/middlewares"
	"github.com/MarcelArt/refinery/internal/v1/repositories"
	"github.com/MarcelArt/refinery/internal/v1/services"
	"go.uber.org/dig"
)

func New() *dig.Container {
	c := dig.New()

	c.Provide(configs.ConnectDB)

	c.Provide(repositories.NewN8NRepo, dig.As(new(repositories.IN8NRepo)))
	c.Provide(repositories.NewUserRepo, dig.As(new(repositories.IUserRepo)))
	c.Provide(repositories.NewWorkflowRepo, dig.As(new(repositories.IWorkflowRepo)))

	c.Provide(services.NewN8NService, dig.As(new(services.IN8NService)))
	c.Provide(services.NewUserService, dig.As(new(services.IUserService)))
	c.Provide(services.NewWorkflowService, dig.As(new(services.IWorkflowService)))

	c.Provide(middlewares.NewAuthMiddleware)

	c.Provide(handlers.NewUserHandler)
	c.Provide(handlers.NewWorkflowHandler)

	c.Provide(app.New)

	return c
}
