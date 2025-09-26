package di

import (
	"bm-staff/internal/domain/services"
	"bm-staff/internal/infrastructure/config"
	"bm-staff/internal/infrastructure/database"
	"bm-staff/internal/infrastructure/http"
	"bm-staff/internal/infrastructure/logging"
	"bm-staff/internal/interfaces/http/handlers"
	"bm-staff/internal/interfaces/repositories/oracle"
	"bm-staff/internal/usecases/user"

	"github.com/go-playground/validator/v10"
	"github.com/google/wire"
	"go.uber.org/zap"
)

// Container holds all dependencies
type Container struct {
	Config      *config.Config
	Logger      *zap.Logger
	Database    *database.OracleDB
	Migrator    *database.GORMMigrator
	UserHandler *handlers.UserHandler
	HTTPServer  *http.Server
}

// NewContainer creates a new dependency injection container
func NewContainer() (*Container, error) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	// Create logger
	logger, err := logging.NewLogger(cfg.Logging.Level, cfg.Logging.Format)
	if err != nil {
		return nil, err
	}

	// Create database connection
	dbConfig := &database.OracleConfig{
		Host:            cfg.Database.Host,
		Port:            cfg.Database.Port,
		Username:        cfg.Database.Username,
		Password:        cfg.Database.Password,
		ServiceName:     cfg.Database.ServiceName,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	}

	oracleDB, err := database.NewOracleDB(dbConfig, logger)
	if err != nil {
		return nil, err
	}

	// Create GORM migrator
	dsn := database.BuildOracleDSN(dbConfig)
	migrator, err := database.NewGORMMigrator(dsn, logger)
	if err != nil {
		return nil, err
	}

	// Create repository
	userRepo := oracle.NewUserRepository(oracleDB.DB(), logger)

	// Create domain service
	userService := services.NewUserService(userRepo)

	// Create use cases
	createUserUseCase := user.NewCreateUserUseCase(userRepo, userService)
	getUserUseCase := user.NewGetUserUseCase(userRepo)
	updateUserUseCase := user.NewUpdateUserUseCase(userRepo, userService)
	deleteUserUseCase := user.NewDeleteUserUseCase(userRepo, userService)

	// Create validator
	validator := validator.New()

	// Create handler
	userHandler := handlers.NewUserHandler(
		createUserUseCase,
		getUserUseCase,
		updateUserUseCase,
		deleteUserUseCase,
		validator,
		logger,
	)

	// Create HTTP server
	httpServer := http.NewServer(cfg, logger, userHandler)

	return &Container{
		Config:      cfg,
		Logger:      logger,
		Database:    oracleDB,
		Migrator:    migrator,
		UserHandler: userHandler,
		HTTPServer:  httpServer,
	}, nil
}

// WireSet is the Wire provider set
var WireSet = wire.NewSet(
	config.Load,
	logging.NewLogger,
	database.NewOracleDB,
	database.NewGORMMigrator,
	oracle.NewUserRepository,
	services.NewUserService,
	user.NewCreateUserUseCase,
	user.NewGetUserUseCase,
	user.NewUpdateUserUseCase,
	user.NewDeleteUserUseCase,
	handlers.NewUserHandler,
	http.NewServer,
	NewContainer,
)
