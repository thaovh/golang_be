package oracle

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"bm-staff/internal/domain/entities"
	"bm-staff/internal/domain/repositories"

	"go.uber.org/zap"
)

// RefreshTokenRepository implements the refresh token repository interface for Oracle
type RefreshTokenRepository struct {
	db     *sql.DB
	logger *zap.Logger
}

// NewRefreshTokenRepository creates a new Oracle refresh token repository
func NewRefreshTokenRepository(db *sql.DB, logger *zap.Logger) repositories.RefreshTokenRepository {
	return &RefreshTokenRepository{
		db:     db,
		logger: logger,
	}
}

// Create creates a new refresh token
func (r *RefreshTokenRepository) Create(ctx context.Context, refreshToken *entities.RefreshToken) error {
	query := `
		INSERT INTO BMSF_REFRESH_TOKEN (
			ID, CREATED_AT, UPDATED_AT, VERSION,
			USER_ID, TOKEN, EXPIRES_AT, IS_REVOKED, 
			REVOKED_AT, IP_ADDRESS, USER_AGENT
		) VALUES (
			:1, :2, :3, :4, :5, :6, :7, :8, :9, :10, :11
		)`

	_, err := r.db.ExecContext(ctx, query,
		refreshToken.ID,
		refreshToken.CreatedAt,
		refreshToken.UpdatedAt,
		refreshToken.Version,
		refreshToken.UserID,
		refreshToken.Token,
		refreshToken.ExpiresAt,
		refreshToken.IsRevoked,
		refreshToken.RevokedAt,
		refreshToken.IPAddress,
		refreshToken.UserAgent,
	)

	if err != nil {
		r.logger.Error("Failed to create refresh token",
			zap.String("user_id", refreshToken.UserID.String()),
			zap.Error(err),
		)
		return fmt.Errorf("failed to create refresh token: %w", err)
	}

	r.logger.Info("Refresh token created successfully",
		zap.String("user_id", refreshToken.UserID.String()),
		zap.String("token_id", refreshToken.ID.String()),
	)

	return nil
}

// GetByID gets a refresh token by ID
func (r *RefreshTokenRepository) GetByID(ctx context.Context, id string) (*entities.RefreshToken, error) {
	query := `
		SELECT ID, CREATED_AT, UPDATED_AT, CREATED_BY, UPDATED_BY, 
		       DELETED_AT, VERSION, TENANT_ID,
		       USER_ID, TOKEN, EXPIRES_AT, IS_REVOKED, 
		       REVOKED_AT, IP_ADDRESS, USER_AGENT
		FROM BMSF_REFRESH_TOKEN 
		WHERE ID = :1 AND DELETED_AT IS NULL`

	var refreshToken entities.RefreshToken
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&refreshToken.ID,
		&refreshToken.CreatedAt,
		&refreshToken.UpdatedAt,
		&refreshToken.CreatedBy,
		&refreshToken.UpdatedBy,
		&refreshToken.DeletedAt,
		&refreshToken.Version,
		&refreshToken.TenantID,
		&refreshToken.UserID,
		&refreshToken.Token,
		&refreshToken.ExpiresAt,
		&refreshToken.IsRevoked,
		&refreshToken.RevokedAt,
		&refreshToken.IPAddress,
		&refreshToken.UserAgent,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("refresh token not found")
		}
		r.logger.Error("Failed to get refresh token by ID",
			zap.String("id", id),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	return &refreshToken, nil
}

// GetByToken gets a refresh token by token string
func (r *RefreshTokenRepository) GetByToken(ctx context.Context, token string) (*entities.RefreshToken, error) {
	query := `
		SELECT ID, CREATED_AT, UPDATED_AT, CREATED_BY, UPDATED_BY, 
		       DELETED_AT, VERSION, TENANT_ID,
		       USER_ID, TOKEN, EXPIRES_AT, IS_REVOKED, 
		       REVOKED_AT, IP_ADDRESS, USER_AGENT
		FROM BMSF_REFRESH_TOKEN 
		WHERE TOKEN = :1 AND DELETED_AT IS NULL`

	var refreshToken entities.RefreshToken
	err := r.db.QueryRowContext(ctx, query, token).Scan(
		&refreshToken.ID,
		&refreshToken.CreatedAt,
		&refreshToken.UpdatedAt,
		&refreshToken.CreatedBy,
		&refreshToken.UpdatedBy,
		&refreshToken.DeletedAt,
		&refreshToken.Version,
		&refreshToken.TenantID,
		&refreshToken.UserID,
		&refreshToken.Token,
		&refreshToken.ExpiresAt,
		&refreshToken.IsRevoked,
		&refreshToken.RevokedAt,
		&refreshToken.IPAddress,
		&refreshToken.UserAgent,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("refresh token not found")
		}
		r.logger.Error("Failed to get refresh token by token",
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get refresh token: %w", err)
	}

	return &refreshToken, nil
}

// GetByUserID gets all refresh tokens for a user
func (r *RefreshTokenRepository) GetByUserID(ctx context.Context, userID string) ([]*entities.RefreshToken, error) {
	query := `
		SELECT ID, CREATED_AT, UPDATED_AT, CREATED_BY, UPDATED_BY, 
		       DELETED_AT, VERSION, TENANT_ID,
		       USER_ID, TOKEN, EXPIRES_AT, IS_REVOKED, 
		       REVOKED_AT, IP_ADDRESS, USER_AGENT
		FROM BMSF_REFRESH_TOKEN 
		WHERE USER_ID = :1 AND DELETED_AT IS NULL
		ORDER BY CREATED_AT DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		r.logger.Error("Failed to get refresh tokens by user ID",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to get refresh tokens: %w", err)
	}
	defer rows.Close()

	var refreshTokens []*entities.RefreshToken
	for rows.Next() {
		var refreshToken entities.RefreshToken
		err := rows.Scan(
			&refreshToken.ID,
			&refreshToken.CreatedAt,
			&refreshToken.UpdatedAt,
			&refreshToken.CreatedBy,
			&refreshToken.UpdatedBy,
			&refreshToken.DeletedAt,
			&refreshToken.Version,
			&refreshToken.TenantID,
			&refreshToken.UserID,
			&refreshToken.Token,
			&refreshToken.ExpiresAt,
			&refreshToken.IsRevoked,
			&refreshToken.RevokedAt,
			&refreshToken.IPAddress,
			&refreshToken.UserAgent,
		)
		if err != nil {
			r.logger.Error("Failed to scan refresh token",
				zap.Error(err),
			)
			return nil, fmt.Errorf("failed to scan refresh token: %w", err)
		}
		refreshTokens = append(refreshTokens, &refreshToken)
	}

	return refreshTokens, nil
}

// Update updates an existing refresh token
func (r *RefreshTokenRepository) Update(ctx context.Context, refreshToken *entities.RefreshToken) error {
	query := `
		UPDATE BMSF_REFRESH_TOKEN SET
			UPDATED_AT = :1,
			UPDATED_BY = :2,
			VERSION = :3,
			IS_REVOKED = :4,
			REVOKED_AT = :5
		WHERE ID = :6 AND VERSION = :7`

	result, err := r.db.ExecContext(ctx, query,
		refreshToken.UpdatedAt,
		refreshToken.UpdatedBy,
		refreshToken.Version,
		refreshToken.IsRevoked,
		refreshToken.RevokedAt,
		refreshToken.ID,
		refreshToken.Version-1, // Check against old version
	)

	if err != nil {
		r.logger.Error("Failed to update refresh token",
			zap.String("id", refreshToken.ID.String()),
			zap.Error(err),
		)
		return fmt.Errorf("failed to update refresh token: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("refresh token not found or version mismatch")
	}

	r.logger.Info("Refresh token updated successfully",
		zap.String("id", refreshToken.ID.String()),
	)

	return nil
}

// Delete deletes a refresh token
func (r *RefreshTokenRepository) Delete(ctx context.Context, id string) error {
	query := `UPDATE BMSF_REFRESH_TOKEN SET DELETED_AT = :1 WHERE ID = :2`

	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		r.logger.Error("Failed to delete refresh token",
			zap.String("id", id),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete refresh token: %w", err)
	}

	r.logger.Info("Refresh token deleted successfully",
		zap.String("id", id),
	)

	return nil
}

// RevokeAllForUser revokes all refresh tokens for a user
func (r *RefreshTokenRepository) RevokeAllForUser(ctx context.Context, userID string) error {
	query := `
		UPDATE BMSF_REFRESH_TOKEN SET
			IS_REVOKED = 1,
			REVOKED_AT = :1,
			UPDATED_AT = :2
		WHERE USER_ID = :3 AND IS_REVOKED = 0 AND DELETED_AT IS NULL`

	now := time.Now()
	_, err := r.db.ExecContext(ctx, query, now, now, userID)
	if err != nil {
		r.logger.Error("Failed to revoke all refresh tokens for user",
			zap.String("user_id", userID),
			zap.Error(err),
		)
		return fmt.Errorf("failed to revoke refresh tokens: %w", err)
	}

	r.logger.Info("All refresh tokens revoked for user",
		zap.String("user_id", userID),
	)

	return nil
}

// CleanupExpired removes expired refresh tokens
func (r *RefreshTokenRepository) CleanupExpired(ctx context.Context) error {
	query := `DELETE FROM BMSF_REFRESH_TOKEN WHERE EXPIRES_AT < :1`

	result, err := r.db.ExecContext(ctx, query, time.Now())
	if err != nil {
		r.logger.Error("Failed to cleanup expired refresh tokens",
			zap.Error(err),
		)
		return fmt.Errorf("failed to cleanup expired refresh tokens: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	r.logger.Info("Expired refresh tokens cleaned up",
		zap.Int64("count", rowsAffected),
	)

	return nil
}
