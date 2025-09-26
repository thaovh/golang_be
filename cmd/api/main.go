package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "bm-staff/docs" // Import docs for Swagger
	"bm-staff/internal/di"

	"go.uber.org/zap"
)

// @title           BM Staff API
// @version         1.0
// @description     BM Staff Framework API with Clean Architecture
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	// Create dependency injection container
	container, err := di.NewContainer()
	if err != nil {
		log.Fatalf("Failed to create container: %v", err)
	}

	// Ensure database connection is closed
	defer func() {
		if err := container.Database.Close(); err != nil {
			container.Logger.Error("Failed to close database connection", zap.Error(err))
		}
	}()

	// Ensure logger is synced
	defer container.Logger.Sync()

	// Run auto-migration if enabled
	if container.Config.Database.AutoMigrate {
		container.Logger.Info("Auto-migration is enabled, running database migration...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		if err := container.Migrator.AutoMigrate(ctx); err != nil {
			container.Logger.Fatal("Failed to run auto-migration", zap.Error(err))
		}
		container.Logger.Info("Auto-migration completed successfully")
	} else {
		container.Logger.Info("Auto-migration is disabled")
	}

	// Start HTTP server in a goroutine
	go func() {
		container.Logger.Info("Starting application")
		if err := container.HTTPServer.Start(); err != nil {
			container.Logger.Fatal("Failed to start HTTP server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	container.Logger.Info("Shutting down server...")

	// Create a deadline for shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Shutdown HTTP server
	if err := container.HTTPServer.Stop(ctx); err != nil {
		container.Logger.Error("Server forced to shutdown", zap.Error(err))
	}

	container.Logger.Info("Server exited")
}
