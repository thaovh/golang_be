package oracle

import (
	"context"
	"database/sql"
	"fmt"

	"bm-staff/internal/domain/entities"
	"bm-staff/internal/domain/repositories"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// userRepository implements the UserRepository interface for Oracle
type userRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewUserRepository creates a new Oracle user repository
func NewUserRepository(db *sql.DB, logger *zap.Logger) repositories.UserRepository {
	return &userRepository{
		db:     db,
		logger: logger,
	}
}

// Create creates a new user
func (r *userRepository) Create(ctx context.Context, user *entities.User) error {
	query := `
		INSERT INTO BMSF_USER (
			ID, USERNAME, EMAIL, FIRST_NAME, LAST_NAME, PHONE, 
			STATUS, CREATED_AT, UPDATED_AT, CREATED_BY, UPDATED_BY, 
			DELETED_AT, VERSION, TENANT_ID
		) VALUES (
			:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11, :12, :13, :14
		)`

	_, err := r.db.ExecContext(ctx, query,
		user.ID.String(),
		user.Username,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Phone,
		string(user.Status),
		user.CreatedAt,
		user.UpdatedAt,
		user.CreatedBy,
		user.UpdatedBy,
		user.DeletedAt,
		user.Version,
		user.TenantID,
	)

	if err != nil {
		r.logger.Error("Failed to create user",
			zap.String("user_id", user.ID.String()),
			zap.String("username", user.Username),
			zap.Error(err),
		)
		return fmt.Errorf("failed to create user: %w", err)
	}

	r.logger.Info("User created successfully",
		zap.String("user_id", user.ID.String()),
		zap.String("username", user.Username),
	)

	return nil
}

// GetByID retrieves a user by ID
func (r *userRepository) GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	query := `
		SELECT ID, USERNAME, EMAIL, FIRST_NAME, LAST_NAME, PHONE, 
			   STATUS, CREATED_AT, UPDATED_AT, CREATED_BY, UPDATED_BY, 
			   DELETED_AT, VERSION, TENANT_ID
		FROM BMSF_USER 
		WHERE ID = :1 AND DELETED_AT IS NULL`

	var user entities.User
	var status string

	err := r.db.QueryRowContext(ctx, query, id.String()).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&status,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.CreatedBy,
		&user.UpdatedBy,
		&user.DeletedAt,
		&user.Version,
		&user.TenantID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.logger.Error("Failed to get user by ID",
			zap.String("user_id", id.String()),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get user by ID: %w", err)
	}

	user.Status = entities.UserStatus(status)
	return &user, nil
}

// GetByUsername retrieves a user by username
func (r *userRepository) GetByUsername(ctx context.Context, username string) (*entities.User, error) {
	query := `
		SELECT ID, USERNAME, EMAIL, FIRST_NAME, LAST_NAME, PHONE, 
			   STATUS, CREATED_AT, UPDATED_AT, CREATED_BY, UPDATED_BY, 
			   DELETED_AT, VERSION, TENANT_ID
		FROM BMSF_USER 
		WHERE USERNAME = :1 AND DELETED_AT IS NULL`

	var user entities.User
	var status string

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&status,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.CreatedBy,
		&user.UpdatedBy,
		&user.DeletedAt,
		&user.Version,
		&user.TenantID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.logger.Error("Failed to get user by username",
			zap.String("username", username),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get user by username: %w", err)
	}

	user.Status = entities.UserStatus(status)
	return &user, nil
}

// GetByEmail retrieves a user by email
func (r *userRepository) GetByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `
		SELECT ID, USERNAME, EMAIL, FIRST_NAME, LAST_NAME, PHONE, 
			   STATUS, CREATED_AT, UPDATED_AT, CREATED_BY, UPDATED_BY, 
			   DELETED_AT, VERSION, TENANT_ID
		FROM BMSF_USER 
		WHERE EMAIL = :1 AND DELETED_AT IS NULL`

	var user entities.User
	var status string

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.FirstName,
		&user.LastName,
		&user.Phone,
		&status,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.CreatedBy,
		&user.UpdatedBy,
		&user.DeletedAt,
		&user.Version,
		&user.TenantID,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		r.logger.Error("Failed to get user by email",
			zap.String("email", email),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	user.Status = entities.UserStatus(status)
	return &user, nil
}

// Update updates an existing user
func (r *userRepository) Update(ctx context.Context, user *entities.User) error {
	query := `
		UPDATE BMSF_USER 
		SET USERNAME = :1, EMAIL = :2, FIRST_NAME = :3, LAST_NAME = :4, 
			PHONE = :5, STATUS = :6, UPDATED_AT = :7, UPDATED_BY = :8, 
			VERSION = :9
		WHERE ID = :10 AND DELETED_AT IS NULL`

	result, err := r.db.ExecContext(ctx, query,
		user.Username,
		user.Email,
		user.FirstName,
		user.LastName,
		user.Phone,
		string(user.Status),
		user.UpdatedAt,
		user.UpdatedBy,
		user.Version,
		user.ID.String(),
	)

	if err != nil {
		r.logger.Error("Failed to update user",
			zap.String("user_id", user.ID.String()),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	r.logger.Info("User updated successfully",
		zap.String("user_id", user.ID.String()),
	)

	return nil
}

// Delete performs soft delete of a user by ID
func (r *userRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE BMSF_USER 
		SET DELETED_AT = CURRENT_TIMESTAMP, VERSION = VERSION + 1
		WHERE ID = :1 AND DELETED_AT IS NULL`

	result, err := r.db.ExecContext(ctx, query, id.String())
	if err != nil {
		r.logger.Error("Failed to delete user",
			zap.String("user_id", id.String()),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	r.logger.Info("User deleted successfully",
		zap.String("user_id", id.String()),
	)

	return nil
}

// List retrieves users with pagination
func (r *userRepository) List(ctx context.Context, limit, offset int) ([]*entities.User, error) {
	query := `
		SELECT ID, USERNAME, EMAIL, FIRST_NAME, LAST_NAME, PHONE, 
			   STATUS, CREATED_AT, UPDATED_AT, CREATED_BY, UPDATED_BY, 
			   DELETED_AT, VERSION, TENANT_ID
		FROM BMSF_USER 
		WHERE DELETED_AT IS NULL
		ORDER BY CREATED_AT DESC
		OFFSET :1 ROWS FETCH NEXT :2 ROWS ONLY`

	rows, err := r.db.QueryContext(ctx, query, offset, limit)
	if err != nil {
		r.logger.Error("Failed to list users",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		var user entities.User
		var status string

		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Phone,
			&status,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.CreatedBy,
			&user.UpdatedBy,
			&user.DeletedAt,
			&user.Version,
			&user.TenantID,
		)
		if err != nil {
			r.logger.Error("Failed to scan user row",
				zap.Error(err),
			)
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}

		user.Status = entities.UserStatus(status)
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return users, nil
}

// Count returns the total number of users
func (r *userRepository) Count(ctx context.Context) (int64, error) {
	query := `SELECT COUNT(*) FROM BMSF_USER WHERE DELETED_AT IS NULL`

	var count int64
	err := r.db.QueryRowContext(ctx, query).Scan(&count)
	if err != nil {
		r.logger.Error("Failed to count users",
			zap.Error(err),
		)
		return 0, fmt.Errorf("failed to count users: %w", err)
	}

	return count, nil
}

// GetByIDs retrieves multiple users by IDs (for DataLoader)
func (r *userRepository) GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*entities.User, error) {
	if len(ids) == 0 {
		return []*entities.User{}, nil
	}

	// For simplicity, we'll use Oracle's TABLE function for multiple IDs
	query := `
		SELECT ID, USERNAME, EMAIL, FIRST_NAME, LAST_NAME, PHONE, 
			   STATUS, CREATED_AT, UPDATED_AT, CREATED_BY, UPDATED_BY, 
			   DELETED_AT, VERSION, TENANT_ID
		FROM BMSF_USER 
		WHERE ID IN (SELECT COLUMN_VALUE FROM TABLE(SYS.ODCIVARCHAR2LIST(:1, :2, :3, :4, :5)))
		AND DELETED_AT IS NULL
		ORDER BY CREATED_AT DESC`

	// Convert UUIDs to strings and pad with empty strings if needed
	idStrings := make([]string, 5)
	for i, id := range ids {
		if i < 5 {
			idStrings[i] = id.String()
		}
	}
	// Pad remaining slots with empty strings
	for i := len(ids); i < 5; i++ {
		idStrings[i] = ""
	}

	rows, err := r.db.QueryContext(ctx, query, idStrings[0], idStrings[1], idStrings[2], idStrings[3], idStrings[4])
	if err != nil {
		r.logger.Error("Failed to get users by IDs",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get users by IDs: %w", err)
	}
	defer rows.Close()

	var users []*entities.User
	for rows.Next() {
		var user entities.User
		var status string

		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.FirstName,
			&user.LastName,
			&user.Phone,
			&status,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.CreatedBy,
			&user.UpdatedBy,
			&user.DeletedAt,
			&user.Version,
			&user.TenantID,
		)
		if err != nil {
			r.logger.Error("Failed to scan user row",
				zap.Error(err),
			)
			return nil, fmt.Errorf("failed to scan user row: %w", err)
		}

		user.Status = entities.UserStatus(status)
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating user rows: %w", err)
	}

	return users, nil
}
