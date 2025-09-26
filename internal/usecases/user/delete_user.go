package user

import (
	"context"

	"bm-staff/internal/domain/repositories"
	"bm-staff/internal/domain/services"
	"bm-staff/pkg/errors"

	"github.com/google/uuid"
)

// DeleteUserRequest represents the request to delete a user
type DeleteUserRequest struct {
	ID string `json:"id" validate:"required,uuid"`
}

// DeleteUserResponse represents the response after deleting a user
type DeleteUserResponse struct {
	Success bool `json:"success"`
}

// DeleteUserUseCase handles user deletion business logic
type DeleteUserUseCase struct {
	userRepo    repositories.UserRepository
	userService *services.UserService
}

// NewDeleteUserUseCase creates a new delete user use case
func NewDeleteUserUseCase(userRepo repositories.UserRepository, userService *services.UserService) *DeleteUserUseCase {
	return &DeleteUserUseCase{
		userRepo:    userRepo,
		userService: userService,
	}
}

// Execute deletes a user by ID
func (uc *DeleteUserUseCase) Execute(ctx context.Context, req *DeleteUserRequest) (*DeleteUserResponse, error) {
	// Parse UUID
	userID, err := uuid.Parse(req.ID)
	if err != nil {
		return nil, errors.NewValidationError("VAL_002", "Invalid user ID format", map[string]any{
			"id": req.ID,
		})
	}

	// Get existing user
	user, err := uc.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.WrapError(err, "BIZ_001", "Failed to get user")
	}

	if user == nil {
		return nil, errors.NewBusinessError("BIZ_001", "User not found", map[string]any{
			"id": req.ID,
		})
	}

	// Check if user can be deleted according to business rules
	if err := uc.userService.CanDelete(ctx, user); err != nil {
		return nil, errors.NewBusinessError("BIZ_002", "User cannot be deleted", map[string]any{
			"error": err.Error(),
		})
	}

	// Delete user from repository
	if err := uc.userRepo.Delete(ctx, userID); err != nil {
		return nil, errors.WrapError(err, "BIZ_001", "Failed to delete user")
	}

	return &DeleteUserResponse{
		Success: true,
	}, nil
}
