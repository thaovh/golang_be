package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/sijms/go-ora/v2"
	"go.uber.org/zap"
)

// OracleConfig holds Oracle database configuration
type OracleConfig struct {
	Host            string
	Port            int
	Username        string
	Password        string
	ServiceName     string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// OracleDB wraps the Oracle database connection
type OracleDB struct {
	db     *sql.DB
	config *OracleConfig
	logger *zap.Logger
}

// NewOracleDB creates a new Oracle database connection
func NewOracleDB(config *OracleConfig, logger *zap.Logger) (*OracleDB, error) {
	dsn := fmt.Sprintf("oracle://%s:%s@%s:%d/%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.ServiceName,
	)

	db, err := sql.Open("oracle", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open Oracle database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping Oracle database: %w", err)
	}

	logger.Info("Successfully connected to Oracle database",
		zap.String("host", config.Host),
		zap.Int("port", config.Port),
		zap.String("service", config.ServiceName),
	)

	return &OracleDB{
		db:     db,
		config: config,
		logger: logger,
	}, nil
}

// DB returns the underlying sql.DB instance
func (o *OracleDB) DB() *sql.DB {
	return o.db
}

// Close closes the database connection
func (o *OracleDB) Close() error {
	if o.db != nil {
		o.logger.Info("Closing Oracle database connection")
		return o.db.Close()
	}
	return nil
}

// Health checks the database health
func (o *OracleDB) Health(ctx context.Context) error {
	return o.db.PingContext(ctx)
}

// Stats returns database connection statistics
func (o *OracleDB) Stats() sql.DBStats {
	return o.db.Stats()
}

// BuildOracleDSN builds Oracle DSN string from config
func BuildOracleDSN(config *OracleConfig) string {
	return fmt.Sprintf("oracle://%s:%s@%s:%d/%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.ServiceName,
	)
}
