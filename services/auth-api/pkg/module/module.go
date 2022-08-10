package module

import (
	"fmt"

	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/app"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/docs"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/module/auth"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/module/user"

	fiberSwagger "github.com/swaggo/fiber-swagger"
)

func Init(ctx *app.Context) {
	auth.Init(ctx)
	user.Init(ctx)

	ctx.Router.Get("/healthz", healthCheckHandler)

	//Swagger Doc details
	host := ctx.Config.Gateway.Host
	basePath := ctx.Config.Gateway.BaseURL

	if len(host) == 0 {
		host = fmt.Sprintf("localhost:%v", ctx.Config.Server.Port)
	}

	if len(basePath) == 0 {
		basePath = ctx.Config.App.BaseUrl
	}

	docs.SwaggerInfo.Title = "Todo Service API Document"
	docs.SwaggerInfo.Description = "List of APIs for Todo Service."
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = host
	docs.SwaggerInfo.BasePath = basePath
	docs.SwaggerInfo.Schemes = []string{"https", "http"}

	//Init Swagger routes
	ctx.Router.Get("/swagger/*", fiberSwagger.WrapHandler)
}

func healthCheckHandler(c *fiber.Ctx) error {
	return c.SendStatus(http.StatusOK)
}
