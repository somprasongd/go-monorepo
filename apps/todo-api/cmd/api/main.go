package main

import (
	"github.com/somprasongd/go-monorepo/services/todo/pkg/app"
	"github.com/somprasongd/go-monorepo/services/todo/pkg/config"
	"github.com/somprasongd/go-monorepo/services/todo/pkg/module"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	app := app.New(cfg)
	// Cleanup when server stopped
	defer app.Close()

	// For Liveness Probe
	app.CreateLivenessFile()

	// Initialize data sources
	app.InitDS()

	// Create router (mux/gin/fiber)
	app.InitRouter()

	// Initialize module with dependency injection
	module.Init(app.Context)
	// Start server
	app.ServeHTTP()
}
