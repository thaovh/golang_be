package services

import (
	"context"
	"fmt"

	"bm-staff/internal/domain/entities"
	"bm-staff/internal/domain/repositories"
)

// UserService handles user-related business logic
type UserService struct {
	userRepo repositories.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repositories.UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

// ValidateUser validates user data according to business rules
func (s *UserService) ValidateUser(ctx context.Context, user *entities.User) error {
	if user.Username == "" {
		return fmt.Errorf("username is required")
	}

	if user.Email == "" {
		return fmt.Errorf("email is required")
	}

	if user.FirstName == "" {
		return fmt.Errorf("first name is required")
	}

	if user.LastName == "" {
		return fmt.Errorf("last name is required")
	}

	// Check if username already exists
	existingUser, err := s.userRepo.GetByUsername(ctx, user.Username)
	if err == nil && existingUser != nil && existingUser.ID != user.ID {
		return fmt.Errorf("username already exists")
	}

	// Check if email already exists
	existingUser, err = s.userRepo.GetByEmail(ctx, user.Email)
	if err == nil && existingUser != nil && existingUser.ID != user.ID {
		return fmt.Errorf("email already exists")
	}

	return nil
}

// CanActivate checks if a user can be activated
func (s *UserService) CanActivate(ctx context.Context, user *entities.User) error {
	if user.Status == entities.UserStatusBlocked {
		return fmt.Errorf("blocked user cannot be activated")
	}

	return nil
}

// CanDelete checks if a user can be deleted
func (s *UserService) CanDelete(ctx context.Context, user *entities.User) error {
	if user.Status == entities.UserStatusActive {
		return fmt.Errorf("active user cannot be deleted, please deactivate first")
	}

	return nil
}
