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

	c.Provide(services.NewFileService, dig.As(new(services.IFileService)))
	c.Provide(services.NewUserService, dig.As(new(services.IUserService)))

	c.Provide(middlewares.NewAuthMiddleware)

	c.Provide(handlers.NewUserHandler)

	c.Provide(app.New)

	return c
}
