package user

import (
	"context"

	"bm-staff/internal/domain/entities"
	"bm-staff/internal/domain/repositories"
	"bm-staff/pkg/errors"

	"github.com/google/uuid"
)

// GetUserRequest represents the request to get a user
type GetUserRequest struct {
	ID string `json:"id" validate:"required,uuid"`
}

// GetUserResponse represents the response after getting a user
type GetUserResponse struct {
	User *entities.User `json:"user"`
}

// GetUserUseCase handles user retrieval business logic
type GetUserUseCase struct {
	userRepo repositories.UserRepository
}

// NewGetUserUseCase creates a new get user use case
func NewGetUserUseCase(userRepo repositories.UserRepository) *GetUserUseCase {
	return &GetUserUseCase{
		userRepo: userRepo,
	}
}

// Execute retrieves a user by ID
func (uc *GetUserUseCase) Execute(ctx context.Context, req *GetUserRequest) (*GetUserResponse, error) {
	// Parse UUID
	userID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, errors.NewValidationError("VAL_002", "Invalid user ID format", map[string]any{
			"id": req.ID,
		})
	}

	// Get user from repository
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.WrapError(err, "BIZ_001", "Failed to get user")
	}

	if user == nil {
		return nil, errors.NewBusinessError("BIZ_001", "User not found", map[string]any{
			"id": req.ID,
		})
	}

	return &GetUserResponse{
		User: user,
	}, nil
}
