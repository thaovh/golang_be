package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"bm-staff/internal/infrastructure/config"
	"bm-staff/internal/interfaces/http/handlers"
	"bm-staff/internal/interfaces/http/middleware"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

// Server represents the HTTP server
type Server struct {
	config  *config.Config
	logger  *zap.Logger
	handler *gin.Engine
	server  *http.Server
}

// NewServer creates a new HTTP server
func NewServer(config *config.Config, logger *zap.Logger, userHandler *handlers.UserHandler, authHandler *handlers.AuthHandler, authMiddleware *middleware.AuthMiddleware) *Server {
	// Set Gin mode
	if config.Logging.Level == "debug" {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	// Create Gin engine
	engine := gin.New()

	// Add middleware
	engine.Use(gin.Recovery())
	engine.Use(LoggerMiddleware(logger))

	// Setup routes
	setupRoutes(engine, userHandler, authHandler, authMiddleware)

	return &Server{
		config:  config,
		logger:  logger,
		handler: engine,
	}
}

// setupRoutes sets up all HTTP routes
func setupRoutes(engine *gin.Engine, userHandler *handlers.UserHandler, authHandler *handlers.AuthHandler, authMiddleware *middleware.AuthMiddleware) {
	// Swagger documentation
	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check
	engine.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ok",
			"timestamp": time.Now(),
		})
	})

	// API v1 routes
	v1 := engine.Group("/api/v1")
	{
		// Authentication routes (public)
		auth := v1.Group("/auth")
		{
			auth.POST("/login", authHandler.Login)
			auth.POST("/logout", authHandler.Logout)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// User routes (protected)
		users := v1.Group("/users")
		users.Use(authMiddleware.RequireAuth()) // Require authentication
		{
			users.POST("", userHandler.CreateUser)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
			users.GET("", userHandler.ListUsers)
		}
	}
}

// Start starts the HTTP server
func (s *Server) Start() error {
	s.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port),
		Handler:      s.handler,
		ReadTimeout:  s.config.Server.ReadTimeout,
		WriteTimeout: s.config.Server.WriteTimeout,
		IdleTimeout:  s.config.Server.IdleTimeout,
	}

	s.logger.Info("Starting HTTP server",
		zap.String("host", s.config.Server.Host),
		zap.Int("port", s.config.Server.Port),
	)

	return s.server.ListenAndServe()
}

// Stop stops the HTTP server gracefully
func (s *Server) Stop(ctx context.Context) error {
	s.logger.Info("Stopping HTTP server")
	return s.server.Shutdown(ctx)
}

// LoggerMiddleware creates a Gin middleware for logging
func LoggerMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
		logger.Info("HTTP Request",
			zap.String("method", param.Method),
			zap.String("path", param.Path),
			zap.Int("status", param.StatusCode),
			zap.Duration("latency", param.Latency),
			zap.String("client_ip", param.ClientIP),
			zap.String("user_agent", param.Request.UserAgent()),
		)
		return ""
	})
}
