package repositories

import (
	"context"

	"bm-staff/internal/domain/entities"

	"github.com/google/uuid"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	// Create creates a new user
	Create(ctx context.Context, user *entities.User) error

	// GetByID retrieves a user by ID
	GetByID(ctx context.Context, id uuid.UUID) (*entities.User, error)

	// GetByUsername retrieves a user by username
	GetByUsername(ctx context.Context, username string) (*entities.User, error)

	// GetByEmail retrieves a user by email
	GetByEmail(ctx context.Context, email string) (*entities.User, error)

	// Update updates an existing user
	Update(ctx context.Context, user *entities.User) error

	// Delete deletes a user by ID
	Delete(ctx context.Context, id uuid.UUID) error

	// List retrieves users with pagination
	List(ctx context.Context, limit, offset int) ([]*entities.User, error)

	// Count returns the total number of users
	Count(ctx context.Context) (int64, error)

	// GetByIDs retrieves multiple users by IDs (for DataLoader)
	GetByIDs(ctx context.Context, ids []uuid.UUID) ([]*entities.User, error)
}
