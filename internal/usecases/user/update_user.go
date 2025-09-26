package user

import (
	"context"

	"bm-staff/internal/domain/entities"
	"bm-staff/internal/domain/repositories"
	"bm-staff/internal/domain/services"
	"bm-staff/pkg/errors"

	"github.com/google/uuid"
)

// UpdateUserRequest represents the request to update a user
type UpdateUserRequest struct {
	ID        string `json:"id" validate:"required,uuid"`
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"required,min=1,max=100"`
	LastName  string `json:"last_name" validate:"required,min=1,max=100"`
	Phone     string `json:"phone" validate:"omitempty,min=10,max=20"`
}

// UpdateUserResponse represents the response after updating a user
type UpdateUserResponse struct {
	User *entities.User `json:"user"`
}

// UpdateUserUseCase handles user update business logic
type UpdateUserUseCase struct {
	userRepo    repositories.UserRepository
	userService *services.UserService
}

// NewUpdateUserUseCase creates a new update user use case
func NewUpdateUserUseCase(userRepo repositories.UserRepository, userService *services.UserService) *UpdateUserUseCase {
	return &UpdateUserUseCase{
		userRepo:    userRepo,
		userService: userService,
	}
}

// Execute updates an existing user
func (uc *UpdateUserUseCase) Execute(ctx context.Context, req *UpdateUserRequest) (*UpdateUserResponse, error) {
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

	// Update user fields
	user.Username = req.Username
	user.Email = req.Email
	user.FirstName = req.FirstName
	user.LastName = req.LastName
	user.Phone = req.Phone
	user.UpdateVersion(nil) // TODO: Pass actual user ID from context

	// Validate user according to business rules
	if err := uc.userService.ValidateUser(ctx, user); err != nil {
		return nil, errors.NewValidationError("VAL_001", "User validation failed", map[string]any{
			"error": err.Error(),
		})
	}

	// Update user in repository
	if err := uc.userRepo.Update(ctx, user); err != nil {
		return nil, errors.WrapError(err, "BIZ_001", "Failed to update user")
	}

	return &UpdateUserResponse{
		User: user,
	}, nil
}
