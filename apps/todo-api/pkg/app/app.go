package app

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	log "github.com/somprasongd/go-monorepo/common/logger"
	"github.com/somprasongd/go-monorepo/common/middleware"
	"github.com/somprasongd/go-monorepo/services/todo/pkg/app/database"
	"github.com/somprasongd/go-monorepo/services/todo/pkg/config"
	"github.com/somprasongd/go-monorepo/services/todo/pkg/util"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type Context struct {
	Config *config.Config
	Router *fiber.App
	DB     *database.GormDB
}

type app struct {
	*Context
}

func New(cfg *config.Config) *app {
	return &app{Context: &Context{
		Config: cfg,
	}}
}

func (a *app) CreateLivenessFile() {
	if a.Config.App.IsProdMode() {
		f, err := os.Create(a.Config.App.LivenessFile)
		if err != nil {
			panic(err)
		}
		f.Close()
	}
}

func (a *app) InitDS() {
	gorm, err := database.NewGormDB(a.Config)
	if err != nil {
		panic(err)
	}
	a.DB = gorm
}

func (a *app) Close() {
	// close database
	if a.DB != nil {
		log.Default.Info("Closing database")
		a.DB.CloseDB()
	}
	// remove liveness file
	log.Default.Info("Removing liveness file")
	os.Remove(a.Config.App.LivenessFile)
}

func (a *app) InitRouter() {
	cfg := fiber.Config{
		AppName:               fmt.Sprintf("%s v%s", a.Config.App.Name, a.Config.App.Version),
		ReadTimeout:           a.Config.Server.TimeoutRead,
		WriteTimeout:          a.Config.Server.TimeoutWrite,
		IdleTimeout:           a.Config.Server.TimeoutIdle,
		DisableStartupMessage: a.Config.App.IsProdMode(),
	}
	r := fiber.New(cfg)
	// Default middleware config
	r.Use(cors.New())
	r.Use(recover.New())
	r.Use(requestid.New())
	r.Use(util.WrapFiberHandler(middleware.LoggerMiddleware))
	r.Use(util.WrapFiberHandler(middleware.PublicRouteMiddleware()))
	// decode id token from gateway to user
	r.Use(util.WrapFiberHandler(middleware.DecodeUserMiddleware))

	a.Router = r
}

func (a *app) ServeHTTP() {
	serverShutdown := make(chan struct{})

	go func() {
		// Listen for syscall signals for process to interrupt/quit
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		<-sig
		err := a.Router.Shutdown()
		// err := a.Router.Shutdown(context.Background()) // for http server
		if err != nil {
			log.Default.Fatal(fmt.Sprintf("server shutdown failed: %+v", err))
		}
		serverShutdown <- struct{}{}
	}()

	// Run the server
	port := a.Config.Server.Port
	log.Default.Info(fmt.Sprintf("Starting server at port %v", port))

	err := a.Router.Listen(fmt.Sprintf("0.0.0.0:%v", port))
	if err != nil && err != http.ErrServerClosed {
		panic(err.Error())
	}

	<-serverShutdown
	log.Default.Info("Running cleanup tasks")
	// Your cleanup tasks go here
}
