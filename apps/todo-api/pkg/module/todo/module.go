package todo

import (
	_ "github.com/somprasongd/go-monorepo/common/swagdto"
	"github.com/somprasongd/go-monorepo/services/todo/pkg/app"
	"github.com/somprasongd/go-monorepo/services/todo/pkg/module/todo/core/ports"
	"github.com/somprasongd/go-monorepo/services/todo/pkg/module/todo/core/service"
	"github.com/somprasongd/go-monorepo/services/todo/pkg/module/todo/handler"
	"github.com/somprasongd/go-monorepo/services/todo/pkg/module/todo/repository"
	"github.com/somprasongd/go-monorepo/services/todo/pkg/util"

	"github.com/gofiber/fiber/v2"
)

type RouteConfig struct {
	BaseURL     string
	Router      *fiber.App
	TodoService ports.TodoService
}

func Init(ctx *app.Context) {
	// สร้าง dependencies ทั้งหมด
	repo := repository.NewTodoRepositoryDB(ctx.DB.DB)
	svc := service.NewTodoService(repo)

	cfg := RouteConfig{
		BaseURL:     ctx.Config.App.BaseUrl,
		Router:      ctx.Router,
		TodoService: svc,
	}

	SetupRoutes(cfg)
	// h := handler.NewTodoHandler(serv)

	// todos := ctx.Router.Group(ctx.Config.App.BaseUrl + "/todos")

	// todos.Post("", util.WrapFiberHandler(h.CreateTodo))
	// todos.Get("", util.WrapFiberHandler(h.ListTodo))
	// todos.Get("/:id", util.WrapFiberHandler(h.GetTodo))
	// todos.Patch("/:id", util.WrapFiberHandler(h.UpdateTodoStatus))
	// todos.Delete("/:id", util.WrapFiberHandler(h.DeleteTodo))
}

func SetupRoutes(cfg RouteConfig) {
	h := handler.NewTodoHandler(cfg.TodoService)

	todos := cfg.Router.Group(cfg.BaseURL + "/todos")

	todos.Post("", util.WrapFiberHandler(h.CreateTodo))
	todos.Get("", util.WrapFiberHandler(h.ListTodo))
	todos.Get("/:id", util.WrapFiberHandler(h.GetTodo))
	todos.Patch("/:id", util.WrapFiberHandler(h.UpdateTodoStatus))
	todos.Delete("/:id", util.WrapFiberHandler(h.DeleteTodo))
}
