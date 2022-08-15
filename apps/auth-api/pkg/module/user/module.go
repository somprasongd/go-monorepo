package user

import (
	"github.com/gofiber/fiber/v2"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/app"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/app/context"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/module/user/core/ports"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/module/user/core/service"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/module/user/handler"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/module/user/repository"
)

type RouteConfig struct {
	BaseURL     string
	Router      *fiber.App
	UserService ports.UserService
}

func Init(ctx *app.Context) {
	// สร้าง dependencies ทั้งหมด
	repo := repository.NewUserRepositoryDB(ctx.DB.DB)
	svc := service.NewUserService(repo)

	cfg := RouteConfig{
		BaseURL:     ctx.Config.App.BaseUrl,
		Router:      ctx.Router,
		UserService: svc,
	}

	SetupRoutes(cfg)
}

func SetupRoutes(cfg RouteConfig) {
	h := handler.NewUserHandler(cfg.UserService)

	userss := cfg.Router.Group(cfg.BaseURL + "/users")

	userss.Post("", context.WrapFiberHandler(h.CreateUser))
	userss.Get("", context.WrapFiberHandler(h.ListUser))
	userss.Get("/:id", context.WrapFiberHandler(h.GetUser))
	userss.Patch("/:id", context.WrapFiberHandler(h.UpdateUserPassword))
	userss.Delete("/:id", context.WrapFiberHandler(h.DeleteUser))
}
