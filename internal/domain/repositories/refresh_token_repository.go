package repositories

import (
	"bm-staff/internal/domain/entities"
	"context"
)

// RefreshTokenRepository defines the interface for refresh token repository operations
type RefreshTokenRepository interface {
	// Create creates a new refresh token
	Create(ctx context.Context, refreshToken *entities.RefreshToken) error

	// GetByID gets a refresh token by ID
	GetByID(ctx context.Context, id string) (*entities.RefreshToken, error)

	// GetByToken gets a refresh token by token string
	GetByToken(ctx context.Context, token string) (*entities.RefreshToken, error)

	// GetByUserID gets all refresh tokens for a user
	GetByUserID(ctx context.Context, userID string) ([]*entities.RefreshToken, error)

	// Update updates an existing refresh token
	Update(ctx context.Context, refreshToken *entities.RefreshToken) error

	// Delete deletes a refresh token
	Delete(ctx context.Context, id string) error

	// RevokeAllForUser revokes all refresh tokens for a user
	RevokeAllForUser(ctx context.Context, userID string) error

	// CleanupExpired removes expired refresh tokens
	CleanupExpired(ctx context.Context) error
}
