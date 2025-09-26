package auth

import (
	"context"

	"bm-staff/internal/domain/repositories"
	"bm-staff/internal/domain/services"
	"bm-staff/pkg/errors"
)

// LogoutRequest represents the request to logout
type LogoutRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// LogoutResponse represents the response after logout
type LogoutResponse struct {
	Message string `json:"message"`
}

// LogoutUseCase handles user logout business logic
type LogoutUseCase struct {
	refreshTokenRepo repositories.RefreshTokenRepository
	jwtService       *services.JWTService
}

// NewLogoutUseCase creates a new logout use case
func NewLogoutUseCase(
	refreshTokenRepo repositories.RefreshTokenRepository,
	jwtService *services.JWTService,
) *LogoutUseCase {
	return &LogoutUseCase{
		refreshTokenRepo: refreshTokenRepo,
		jwtService:       jwtService,
	}
}

// Execute performs user logout
func (uc *LogoutUseCase) Execute(ctx context.Context, req *LogoutRequest) (*LogoutResponse, error) {
	// Validate refresh token
	_, err := uc.jwtService.ValidateToken(req.RefreshToken)
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

	// Revoke refresh token
	refreshToken.Revoke(nil) // No updatedBy for logout
	if err := uc.refreshTokenRepo.Update(ctx, refreshToken); err != nil {
		return nil, errors.WrapError(err, "SYS_001", "Failed to revoke refresh token")
	}

	return &LogoutResponse{
		Message: "Successfully logged out",
	}, nil
}
