package middleware

import (
	"net/http"

	"bm-staff/internal/domain/services"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// AuthMiddleware provides JWT authentication middleware
type AuthMiddleware struct {
	jwtService *services.JWTService
	logger     *zap.Logger
}

// NewAuthMiddleware creates a new authentication middleware
func NewAuthMiddleware(jwtService *services.JWTService, logger *zap.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
		logger:     logger,
	}
}

// RequireAuth middleware that requires valid JWT token
func (am *AuthMiddleware) RequireAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			am.logger.Warn("Missing authorization header",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authorization header is required",
			})
			c.Abort()
			return
		}

		// Extract token from header
		token, err := am.jwtService.ExtractTokenFromHeader(authHeader)
		if err != nil {
			am.logger.Warn("Invalid authorization header format",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.Error(err),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid authorization header format",
			})
			c.Abort()
			return
		}

		// Validate token
		claims, err := am.jwtService.ValidateToken(token)
		if err != nil {
			am.logger.Warn("Invalid token",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.Error(err),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid or expired token",
			})
			c.Abort()
			return
		}

		// Check if token is for API access (audience should be bm-staff-api)
		if len(claims.Audience) == 0 || claims.Audience[0] != "bm-staff-api" {
			am.logger.Warn("Invalid token audience",
				zap.String("path", c.Request.URL.Path),
				zap.String("method", c.Request.Method),
				zap.Strings("audience", claims.Audience),
			)
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid token type",
			})
			c.Abort()
			return
		}

		// Set user information in context
		c.Set("user_id", claims.UserID)
		c.Set("username", claims.Username)
		c.Set("email", claims.Email)
		c.Set("role_id", claims.RoleID)
		c.Set("claims", claims)

		am.logger.Debug("User authenticated successfully",
			zap.String("user_id", claims.UserID.String()),
			zap.String("username", claims.Username),
			zap.String("path", c.Request.URL.Path),
		)

		c.Next()
	}
}

// OptionalAuth middleware that validates JWT token if present
func (am *AuthMiddleware) OptionalAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			// No auth header, continue without authentication
			c.Next()
			return
		}

		// Extract token from header
		token, err := am.jwtService.ExtractTokenFromHeader(authHeader)
		if err != nil {
			// Invalid format, continue without authentication
			c.Next()
			return
		}

		// Validate token
		claims, err := am.jwtService.ValidateToken(token)
		if err != nil {
			// Invalid token, continue without authentication
			c.Next()
			return
		}

		// Check if token is for API access
		if len(claims.Audience) > 0 && claims.Audience[0] == "bm-staff-api" {
			// Set user information in context
			c.Set("user_id", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("email", claims.Email)
			c.Set("role_id", claims.RoleID)
			c.Set("claims", claims)

			am.logger.Debug("User authenticated successfully (optional)",
				zap.String("user_id", claims.UserID.String()),
				zap.String("username", claims.Username),
				zap.String("path", c.Request.URL.Path),
			)
		}

		c.Next()
	}
}

// RequireRole middleware that requires specific role
func (am *AuthMiddleware) RequireRole(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First check if user is authenticated
		userID, exists := c.Get("user_id")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Authentication required",
			})
			c.Abort()
			return
		}

		// Get role from context (this would need to be populated from database)
		// For now, we'll just check if user is authenticated
		// In a real implementation, you'd fetch the user's role from database
		_ = userID
		_ = requiredRole

		// TODO: Implement role-based access control
		// This would require:
		// 1. Fetching user's role from database
		// 2. Checking if user has required role
		// 3. Returning 403 Forbidden if not authorized

		c.Next()
	}
}

// GetCurrentUserID extracts current user ID from context
func GetCurrentUserID(c *gin.Context) (string, bool) {
	userID, exists := c.Get("user_id")
	if !exists {
		return "", false
	}
	return userID.(string), true
}

// GetCurrentUsername extracts current username from context
func GetCurrentUsername(c *gin.Context) (string, bool) {
	username, exists := c.Get("username")
	if !exists {
		return "", false
	}
	return username.(string), true
}

// GetCurrentClaims extracts current JWT claims from context
func GetCurrentClaims(c *gin.Context) (*services.JWTClaims, bool) {
	claims, exists := c.Get("claims")
	if !exists {
		return nil, false
	}
	return claims.(*services.JWTClaims), true
}
