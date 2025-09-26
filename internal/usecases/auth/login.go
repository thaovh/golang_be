package auth

import (
	"context"
	"time"

	"bm-staff/internal/domain/entities"
	"bm-staff/internal/domain/repositories"
	"bm-staff/internal/domain/services"
	"bm-staff/pkg/errors"
)

// LoginRequest represents the request to login
type LoginRequest struct {
	Username string `json:"username" validate:"required,min=3,max=50"`
	Password string `json:"password" validate:"required,min=8,max=100"`
}

// LoginResponse represents the response after login
type LoginResponse struct {
	User      *entities.User      `json:"user"`
	Tokens    *services.TokenPair `json:"tokens"`
	ExpiresIn int64               `json:"expires_in"`
}

// LoginUseCase handles user login business logic
type LoginUseCase struct {
	userRepo         repositories.UserRepository
	refreshTokenRepo repositories.RefreshTokenRepository
	passwordService  *services.PasswordService
	jwtService       *services.JWTService
}

// NewLoginUseCase creates a new login use case
func NewLoginUseCase(
	userRepo repositories.UserRepository,
	refreshTokenRepo repositories.RefreshTokenRepository,
	passwordService *services.PasswordService,
	jwtService *services.JWTService,
) *LoginUseCase {
	return &LoginUseCase{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		passwordService:  passwordService,
		jwtService:       jwtService,
	}
}

// Execute performs user login
func (uc *LoginUseCase) Execute(ctx context.Context, req *LoginRequest, ipAddress, userAgent string) (*LoginResponse, error) {
	// Get user by username
	user, err := uc.userRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, errors.NewValidationError("AUTH_001", "Invalid credentials", nil)
	}

	// Check if user is locked
	if user.IsLocked() {
		return nil, errors.NewValidationError("AUTH_002", "Account is locked due to too many failed login attempts", map[string]any{
			"locked_until": user.LockedUntil,
		})
	}

	// Check if user is active
	if !user.IsActive() {
		return nil, errors.NewValidationError("AUTH_003", "Account is not active", nil)
	}

	// Verify password
	if !uc.passwordService.VerifyPassword(req.Password, user.PasswordHash, user.Salt) {
		// Record failed login attempt
		user.RecordFailedLogin(nil) // No updatedBy for failed login
		if err := uc.userRepo.Update(ctx, user); err != nil {
			// Log error but don't expose it
		}
		return nil, errors.NewValidationError("AUTH_001", "Invalid credentials", nil)
	}

	// Generate tokens
	tokens, err := uc.jwtService.GenerateTokenPair(user.ID, user.Username, user.Email, user.RoleID)
	if err != nil {
		return nil, errors.WrapError(err, "SYS_001", "Failed to generate tokens")
	}

	// Save refresh token to database
	refreshToken := entities.NewRefreshToken(
		user.ID,
		tokens.RefreshToken,
		time.Now().Add(7*24*time.Hour), // 7 days
		ipAddress,
		userAgent,
	)

	if err := uc.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		return nil, errors.WrapError(err, "SYS_001", "Failed to save refresh token")
	}

	// Record successful login
	user.RecordLogin(nil) // No updatedBy for login
	if err := uc.userRepo.Update(ctx, user); err != nil {
		// Log error but don't fail login
	}

	return &LoginResponse{
		User:      user,
		Tokens:    tokens,
		ExpiresIn: tokens.ExpiresIn,
	}, nil
}
