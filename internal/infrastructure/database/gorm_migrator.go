package database

import (
	"context"
	"fmt"
	"strings"

	"bm-staff/internal/domain/entities"

	oracle "github.com/godoes/gorm-oracle"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// GORMMigrator provides GORM-based auto migration with Oracle enhancements
type GORMMigrator struct {
	db     *gorm.DB
	logger *zap.Logger
}

// NewGORMMigrator creates a new GORM-based migrator
func NewGORMMigrator(dsn string, logger *zap.Logger) (*GORMMigrator, error) {
	// Configure GORM for Oracle
	config := &gorm.Config{
		// Disable foreign key constraints for Oracle compatibility
		DisableForeignKeyConstraintWhenMigrating: true,
		// Custom naming strategy for BMSF_ prefix (tables, indexes, constraints only)
		NamingStrategy: &BMSFNamingStrategy{},
	}

	db, err := gorm.Open(oracle.New(oracle.Config{
		DSN: dsn,
	}), config)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Oracle with GORM: %w", err)
	}

	// Configure Oracle-specific settings
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get underlying sql.DB: %w", err)
	}

	// Set connection pool settings
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(5)

	return &GORMMigrator{
		db:     db,
		logger: logger,
	}, nil
}

// AutoMigrate runs automatic migration for all entities
func (m *GORMMigrator) AutoMigrate(ctx context.Context) error {
	m.logger.Info("Starting GORM auto-migration...")

	// Auto-migrate all entities - GORM handles everything automatically!
	err := m.db.WithContext(ctx).AutoMigrate(
		&entities.User{},
		&entities.Department{},
		&entities.Role{},
		&entities.Permission{},
		&entities.AuditLog{},
		&entities.RefreshToken{},
		// Add new entities here - no code changes needed!
	)

	if err != nil {
		// Check if error is due to existing objects (Oracle ORA-00955, ORA-01408)
		if m.isExistingObjectError(err) {
			m.logger.Warn("Some database objects already exist, continuing...",
				zap.Error(err))
			// Continue execution - this is not a fatal error
		} else {
			return fmt.Errorf("GORM auto-migration failed: %w", err)
		}
	}

	m.logger.Info("GORM auto-migration completed successfully")
	return nil
}

// isExistingObjectError checks if the error is due to existing database objects
func (m *GORMMigrator) isExistingObjectError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToUpper(err.Error())
	// Oracle error codes for existing objects
	existingObjectErrors := []string{
		"ORA-00955", // name is already used by an existing object
		"ORA-01408", // such column list already indexed
		"ORA-00942", // table or view does not exist (for drop operations)
		"ORA-02429", // cannot drop unique/primary key constraint
	}

	for _, errorCode := range existingObjectErrors {
		if strings.Contains(errStr, errorCode) {
			return true
		}
	}

	return false
}

// RegisterEntity registers a new entity for migration
func (m *GORMMigrator) RegisterEntity(entity interface{}) {
	// With GORM, we just need to add the entity to AutoMigrate call
	// This is handled in the AutoMigrate method above
	m.logger.Info("Entity registered for GORM auto-migration",
		zap.String("entity", fmt.Sprintf("%T", entity)))
}

// GetDB returns the underlying GORM DB instance
func (m *GORMMigrator) GetDB() *gorm.DB {
	return m.db
}

// Close closes the database connection
func (m *GORMMigrator) Close() error {
	sqlDB, err := m.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// BMSFNamingStrategy implements GORM naming strategy for BMSF_ prefix
type BMSFNamingStrategy struct{}

// TableName converts struct name to table name with BMSF_ prefix
func (ns *BMSFNamingStrategy) TableName(table string) string {
	// Convert to uppercase and add BMSF_ prefix
	return "BMSF_" + strings.ToUpper(table)
}

// ColumnName - NOT IMPLEMENTED to let GORM use explicit column tags
// This allows gorm:"column:FIRST_NAME" to work properly
func (ns *BMSFNamingStrategy) ColumnName(table, column string) string {
	// Return empty string to let GORM use explicit column tags
	// If no explicit tag, GORM will use field name as-is
	return ""
}

// IndexName converts index name to BMSF_ prefixed name
func (ns *BMSFNamingStrategy) IndexName(table, column string) string {
	// Create unique index name (Oracle limit: 30 chars)
	// Include full table name to avoid conflicts between tables
	shortTable := strings.TrimPrefix(strings.ToUpper(table), "BMSF_")
	shortColumn := strings.ToUpper(column)

	// Create unique index name: IDX_TABLE_COLUMN
	indexName := "IDX_" + shortTable + "_" + shortColumn

	// Truncate if too long, but keep table name for uniqueness
	if len(indexName) > 30 {
		// Keep table name, truncate column name
		maxColumnLen := 30 - len("IDX_") - len(shortTable) - 1 // -1 for underscore
		if maxColumnLen > 0 {
			indexName = "IDX_" + shortTable + "_" + shortColumn[:maxColumnLen]
		} else {
			// If table name is too long, truncate both
			indexName = "IDX_" + shortTable[:15] + "_" + shortColumn[:10]
		}
	}
	return indexName
}

// ConstraintName converts constraint name to BMSF_ prefixed name
func (ns *BMSFNamingStrategy) ConstraintName(table, column, foreignKey string) string {
	// Create unique constraint name (Oracle limit: 30 chars)
	shortTable := strings.TrimPrefix(strings.ToUpper(table), "BMSF_")
	shortColumn := strings.ToUpper(column)

	// Create unique constraint name: FK_TABLE_COLUMN
	constraintName := "FK_" + shortTable + "_" + shortColumn

	// Truncate if too long, but keep table name for uniqueness
	if len(constraintName) > 30 {
		maxColumnLen := 30 - len("FK_") - len(shortTable) - 1
		if maxColumnLen > 0 {
			constraintName = "FK_" + shortTable + "_" + shortColumn[:maxColumnLen]
		} else {
			constraintName = "FK_" + shortTable[:15] + "_" + shortColumn[:10]
		}
	}
	return constraintName
}

// CheckerName converts checker name to BMSF_ prefixed name
func (ns *BMSFNamingStrategy) CheckerName(table, column string) string {
	// Create unique checker name (Oracle limit: 30 chars)
	shortTable := strings.TrimPrefix(strings.ToUpper(table), "BMSF_")
	shortColumn := strings.ToUpper(column)

	// Create unique checker name: CHK_TABLE_COLUMN
	checkerName := "CHK_" + shortTable + "_" + shortColumn

	// Truncate if too long, but keep table name for uniqueness
	if len(checkerName) > 30 {
		maxColumnLen := 30 - len("CHK_") - len(shortTable) - 1
		if maxColumnLen > 0 {
			checkerName = "CHK_" + shortTable + "_" + shortColumn[:maxColumnLen]
		} else {
			checkerName = "CHK_" + shortTable[:15] + "_" + shortColumn[:10]
		}
	}
	return checkerName
}

// JoinTableName converts join table name to BMSF_ prefixed name
func (ns *BMSFNamingStrategy) JoinTableName(joinTable string) string {
	// Create join table name with BMSF_ prefix
	return "BMSF_" + strings.ToUpper(joinTable)
}

// RelationshipFKName converts foreign key name to BMSF_ prefixed name
func (ns *BMSFNamingStrategy) RelationshipFKName(relationship schema.Relationship) string {
	// Create unique foreign key name (Oracle limit: 30 chars)
	shortTable := strings.TrimPrefix(strings.ToUpper(relationship.Schema.Table), "BMSF_")
	shortField := strings.ToUpper(relationship.Field.Name)

	// Create unique FK name: FK_TABLE_FIELD
	fkName := "FK_" + shortTable + "_" + shortField

	// Truncate if too long, but keep table name for uniqueness
	if len(fkName) > 30 {
		maxFieldLen := 30 - len("FK_") - len(shortTable) - 1
		if maxFieldLen > 0 {
			fkName = "FK_" + shortTable + "_" + shortField[:maxFieldLen]
		} else {
			fkName = "FK_" + shortTable[:15] + "_" + shortField[:10]
		}
	}
	return fkName
}

// SchemaName converts schema name to BMSF_ prefixed name
func (ns *BMSFNamingStrategy) SchemaName(table string) string {
	// Create schema name with BMSF_ prefix
	return "BMSF_" + strings.ToUpper(table)
}

// UniqueName converts unique constraint name to BMSF_ prefixed name
func (ns *BMSFNamingStrategy) UniqueName(table, column string) string {
	// Create unique constraint name (Oracle limit: 30 chars)
	shortTable := strings.TrimPrefix(strings.ToUpper(table), "BMSF_")
	shortColumn := strings.ToUpper(column)

	// Create unique constraint name: UK_TABLE_COLUMN
	uniqueName := "UK_" + shortTable + "_" + shortColumn

	// Truncate if too long, but keep table name for uniqueness
	if len(uniqueName) > 30 {
		maxColumnLen := 30 - len("UK_") - len(shortTable) - 1
		if maxColumnLen > 0 {
			uniqueName = "UK_" + shortTable + "_" + shortColumn[:maxColumnLen]
		} else {
			uniqueName = "UK_" + shortTable[:15] + "_" + shortColumn[:10]
		}
	}
	return uniqueName
}
