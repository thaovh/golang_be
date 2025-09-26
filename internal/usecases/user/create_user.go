package user

import (
	"context"

	"bm-staff/internal/domain/entities"
	"bm-staff/internal/domain/repositories"
	"bm-staff/internal/domain/services"
	"bm-staff/pkg/errors"
)

// CreateUserRequest represents the request to create a user
type CreateUserRequest struct {
	Username  string `json:"username" validate:"required,min=3,max=50"`
	Email     string `json:"email" validate:"required,email"`
	FirstName string `json:"first_name" validate:"required,min=1,max=100"`
	LastName  string `json:"last_name" validate:"required,min=1,max=100"`
	Phone     string `json:"phone" validate:"omitempty,min=10,max=20"`
	Password  string `json:"password" validate:"required,min=8,max=100"`
}

// CreateUserResponse represents the response after creating a user
type CreateUserResponse struct {
	User *entities.User `json:"user"`
}

// CreateUserUseCase handles user creation business logic
type CreateUserUseCase struct {
	userRepo        repositories.UserRepository
	userService     *services.UserService
	passwordService *services.PasswordService
}

// NewCreateUserUseCase creates a new create user use case
func NewCreateUserUseCase(userRepo repositories.UserRepository, userService *services.UserService, passwordService *services.PasswordService) *CreateUserUseCase {
	return &CreateUserUseCase{
		userRepo:        userRepo,
		userService:     userService,
		passwordService: passwordService,
	}
}

// Execute creates a new user
func (uc *CreateUserUseCase) Execute(ctx context.Context, req *CreateUserRequest) (*CreateUserResponse, error) {
	// Hash password
	passwordHash, salt, err := uc.passwordService.HashPassword(req.Password)
	if err != nil {
		return nil, errors.WrapError(err, "SYS_001", "Failed to hash password")
	}

	// Create user entity
	user := entities.NewUser(
		req.Username,
		req.Email,
		req.FirstName,
		req.LastName,
		req.Phone,
		passwordHash,
		salt,
	)

	// Validate user according to business rules
	if err := uc.userService.ValidateUser(ctx, user); err != nil {
		return nil, errors.NewValidationError("VAL_001", "User validation failed", map[string]any{
			"error": err.Error(),
		})
	}

	// Create user in repository
	if err := uc.userRepo.Create(ctx, user); err != nil {
		return nil, errors.WrapError(err, "BIZ_001", "Failed to create user")
	}

	return &CreateUserResponse{
		User: user,
	}, nil
}
