package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/app"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/module/auth/core/ports"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/module/auth/core/service"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/module/auth/handler"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/module/auth/repository"
	"github.com/somprasongd/go-monorepo/services/auth/pkg/util"
)

type RouteConfig struct {
	BaseURL     string
	Router      *fiber.App
	AuthService ports.AuthService
}

func Init(ctx *app.Context) {
	// สร้าง dependencies ทั้งหมด
	tokenRepo := repository.NewTokenRepository(ctx.Cache)
	repo := repository.NewAuthRepositoryDB(ctx.DB.DB)
	svc := service.NewAuthService(ctx.Config, repo, tokenRepo)

	cfg := RouteConfig{
		BaseURL:     ctx.Config.App.BaseUrl,
		Router:      ctx.Router,
		AuthService: svc,
	}

	SetupRoutes(cfg)
}

func SetupRoutes(cfg RouteConfig) {
	h := handler.NewAuthHandler(cfg.AuthService)

	auth := cfg.Router.Group(cfg.BaseURL + "/auth")

	auth.Post("/register", util.WrapFiberHandler(h.Register))
	auth.Post("/login", util.WrapFiberHandler(h.Login))

	auth.Get("/profile", util.WrapFiberHandler(h.Profile))
	auth.Patch("/profile", util.WrapFiberHandler(h.UpdateProfile))

	auth.Post("/refresh", util.WrapFiberHandler(h.RefreshToken))
	auth.Post("/revoke", util.WrapFiberHandler(h.RevokeToken))
}
