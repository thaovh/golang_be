package auth

import (
	"context"

	"bm-staff/internal/domain/entities"
	"bm-staff/internal/domain/repositories"
	"bm-staff/internal/domain/services"
	"bm-staff/pkg/errors"
)

// RefreshTokenRequest represents the request to refresh token
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// RefreshTokenResponse represents the response after refreshing token
type RefreshTokenResponse struct {
	Tokens    *services.TokenPair `json:"tokens"`
	ExpiresIn int64               `json:"expires_in"`
}

// RefreshTokenUseCase handles token refresh business logic
type RefreshTokenUseCase struct {
	userRepo         repositories.UserRepository
	refreshTokenRepo repositories.RefreshTokenRepository
	jwtService       *services.JWTService
}

// NewRefreshTokenUseCase creates a new refresh token use case
func NewRefreshTokenUseCase(
	userRepo repositories.UserRepository,
	refreshTokenRepo repositories.RefreshTokenRepository,
	jwtService *services.JWTService,
) *RefreshTokenUseCase {
	return &RefreshTokenUseCase{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtService:       jwtService,
	}
}

// Execute refreshes the access token
func (uc *RefreshTokenUseCase) Execute(ctx context.Context, req *RefreshTokenRequest, ipAddress, userAgent string) (*RefreshTokenResponse, error) {
	// Validate refresh token
	claims, err := uc.jwtService.ValidateToken(req.RefreshToken)
	if err != nil {
		return nil, errors.NewValidationError("AUTH_001", "Invalid refresh token", nil)
	}

	// Get refresh token from database
	refreshToken, err := uc.refreshTokenRepo.GetByToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, errors.NewValidationError("AUTH_001", "Invalid refresh token", nil)
	}

	// Check if token is valid
	if !refreshToken.IsValid() {
		return nil, errors.NewValidationError("AUTH_001", "Invalid refresh token", nil)
	}

	// Get user to get current information
	user, err := uc.userRepo.GetByID(ctx, claims.UserID)
	if err != nil {
		return nil, errors.NewValidationError("AUTH_001", "User not found", nil)
	}

	// Check if user is still active
	if !user.IsActive() {
		return nil, errors.NewValidationError("AUTH_003", "Account is not active", nil)
	}

	// Generate new token pair
	tokens, err := uc.jwtService.RefreshToken(req.RefreshToken, user.Username, user.Email, user.RoleID)
	if err != nil {
		return nil, errors.WrapError(err, "SYS_001", "Failed to refresh token")
	}

	// Revoke old refresh token
	refreshToken.Revoke(nil) // No updatedBy for token refresh
	if err := uc.refreshTokenRepo.Update(ctx, refreshToken); err != nil {
		// Log error but don't fail the refresh
	}

	// Save new refresh token to database
	newRefreshToken := entities.NewRefreshToken(
		user.ID,
		tokens.RefreshToken,
		refreshToken.ExpiresAt, // Keep same expiry
		ipAddress,
		userAgent,
	)

	if err := uc.refreshTokenRepo.Create(ctx, newRefreshToken); err != nil {
		// Log error but don't fail the refresh
	}

	return &RefreshTokenResponse{
		Tokens:    tokens,
		ExpiresIn: tokens.ExpiresIn,
	}, nil
}
