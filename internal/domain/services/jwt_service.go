package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTService handles JWT token operations
type JWTService struct {
	secretKey     []byte
	accessExpiry  time.Duration
	refreshExpiry time.Duration
}

// JWTClaims represents the JWT claims
type JWTClaims struct {
	UserID   uuid.UUID  `json:"user_id"`
	Username string     `json:"username"`
	Email    string     `json:"email"`
	RoleID   *uuid.UUID `json:"role_id,omitempty"`
	jwt.RegisteredClaims
}

// TokenPair represents access and refresh token pair
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

// NewJWTService creates a new JWT service
func NewJWTService(secretKey string, accessExpiry, refreshExpiry time.Duration) *JWTService {
	return &JWTService{
		secretKey:     []byte(secretKey),
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
	}
}

// GenerateTokenPair generates both access and refresh tokens
func (js *JWTService) GenerateTokenPair(userID uuid.UUID, username, email string, roleID *uuid.UUID) (*TokenPair, error) {
	// Generate access token
	accessToken, _, err := js.generateAccessToken(userID, username, email, roleID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, _, err := js.generateRefreshToken(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(js.accessExpiry.Seconds()),
		TokenType:    "Bearer",
	}, nil
}

// generateAccessToken generates an access token
func (js *JWTService) generateAccessToken(userID uuid.UUID, username, email string, roleID *uuid.UUID) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(js.accessExpiry)

	claims := &JWTClaims{
		UserID:   userID,
		Username: username,
		Email:    email,
		RoleID:   roleID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "bm-staff",
			Subject:   userID.String(),
			Audience:  []string{"bm-staff-api"},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(js.secretKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// generateRefreshToken generates a refresh token
func (js *JWTService) generateRefreshToken(userID uuid.UUID) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(js.refreshExpiry)

	claims := &JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "bm-staff",
			Subject:   userID.String(),
			Audience:  []string{"bm-staff-refresh"},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			NotBefore: jwt.NewNumericDate(now),
			IssuedAt:  jwt.NewNumericDate(now),
			ID:        uuid.New().String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(js.secretKey)
	if err != nil {
		return "", time.Time{}, err
	}

	return tokenString, expiresAt, nil
}

// ValidateToken validates a JWT token and returns claims
func (js *JWTService) ValidateToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return js.secretKey, nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// RefreshToken generates a new access token from refresh token
func (js *JWTService) RefreshToken(refreshTokenString string, username, email string, roleID *uuid.UUID) (*TokenPair, error) {
	// Validate refresh token
	claims, err := js.ValidateToken(refreshTokenString)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Check if token is for refresh (audience should be bm-staff-refresh)
	if len(claims.Audience) == 0 || claims.Audience[0] != "bm-staff-refresh" {
		return nil, errors.New("invalid token type for refresh")
	}

	// Generate new token pair
	return js.GenerateTokenPair(claims.UserID, username, email, roleID)
}

// ExtractTokenFromHeader extracts token from Authorization header
func (js *JWTService) ExtractTokenFromHeader(authHeader string) (string, error) {
	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) < len(bearerPrefix) || authHeader[:len(bearerPrefix)] != bearerPrefix {
		return "", errors.New("authorization header must start with 'Bearer '")
	}

	return authHeader[len(bearerPrefix):], nil
}
