package handlers

import (
	"net/http"
	"strings"

	"bm-staff/internal/usecases/auth"
	"bm-staff/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// AuthHandler handles authentication-related HTTP requests
type AuthHandler struct {
	loginUseCase        *auth.LoginUseCase
	logoutUseCase       *auth.LogoutUseCase
	refreshTokenUseCase *auth.RefreshTokenUseCase
	validator           *validator.Validate
	logger              *zap.Logger
}

// NewAuthHandler creates a new authentication handler
func NewAuthHandler(
	loginUseCase *auth.LoginUseCase,
	logoutUseCase *auth.LogoutUseCase,
	refreshTokenUseCase *auth.RefreshTokenUseCase,
	validator *validator.Validate,
	logger *zap.Logger,
) *AuthHandler {
	return &AuthHandler{
		loginUseCase:        loginUseCase,
		logoutUseCase:       logoutUseCase,
		refreshTokenUseCase: refreshTokenUseCase,
		validator:           validator,
		logger:              logger,
	}
}

// Login handles POST /api/v1/auth/login
// @Summary      User login
// @Description  Authenticate user with username and password
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        credentials body auth.LoginRequest true "Login credentials"
// @Success      200 {object} map[string]interface{} "Login successful"
// @Failure      400 {object} map[string]interface{} "Bad request - validation error"
// @Failure      401 {object} map[string]interface{} "Unauthorized - invalid credentials"
// @Failure      423 {object} map[string]interface{} "Account locked"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req auth.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		h.logger.Error("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	// Get client IP and User-Agent
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Execute login use case
	response, err := h.loginUseCase.Execute(c.Request.Context(), &req, ipAddress, userAgent)
	if err != nil {
		h.logger.Error("Login failed", zap.Error(err))

		if appErr, ok := err.(*errors.AppError); ok {
			switch appErr.Code {
			case "AUTH_001":
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": appErr.Message,
				})
			case "AUTH_002":
				c.JSON(http.StatusLocked, gin.H{
					"error":   appErr.Message,
					"details": appErr.Details,
				})
			case "AUTH_003":
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": appErr.Message,
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}
		return
	}

	// Set HTTP-only cookie for refresh token
	c.SetCookie(
		"refresh_token",
		response.Tokens.RefreshToken,
		int(response.Tokens.ExpiresIn),
		"/",
		"",
		true, // secure
		true, // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Login successful",
		"data": gin.H{
			"user": gin.H{
				"id":         response.User.ID,
				"username":   response.User.Username,
				"email":      response.User.Email,
				"first_name": response.User.FirstName,
				"last_name":  response.User.LastName,
				"status":     response.User.Status,
			},
			"access_token": response.Tokens.AccessToken,
			"token_type":   response.Tokens.TokenType,
			"expires_in":   response.ExpiresIn,
		},
	})
}

// Logout handles POST /api/v1/auth/logout
// @Summary      User logout
// @Description  Logout user and revoke refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        logout body auth.LogoutRequest true "Logout request"
// @Success      200 {object} map[string]interface{} "Logout successful"
// @Failure      400 {object} map[string]interface{} "Bad request - validation error"
// @Failure      401 {object} map[string]interface{} "Unauthorized - invalid token"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	var req auth.LogoutRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		h.logger.Error("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	// Execute logout use case
	response, err := h.logoutUseCase.Execute(c.Request.Context(), &req)
	if err != nil {
		h.logger.Error("Logout failed", zap.Error(err))

		if appErr, ok := err.(*errors.AppError); ok {
			switch appErr.Code {
			case "AUTH_001":
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": appErr.Message,
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}
		return
	}

	// Clear refresh token cookie
	c.SetCookie(
		"refresh_token",
		"",
		-1,
		"/",
		"",
		true, // secure
		true, // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"message": response.Message,
	})
}

// RefreshToken handles POST /api/v1/auth/refresh
// @Summary      Refresh access token
// @Description  Get new access token using refresh token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        refresh body auth.RefreshTokenRequest true "Refresh token request"
// @Success      200 {object} map[string]interface{} "Token refreshed successfully"
// @Failure      400 {object} map[string]interface{} "Bad request - validation error"
// @Failure      401 {object} map[string]interface{} "Unauthorized - invalid refresh token"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req auth.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request format",
		})
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		h.logger.Error("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Validation failed",
			"details": err.Error(),
		})
		return
	}

	// Get client IP and User-Agent
	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	// Execute refresh token use case
	response, err := h.refreshTokenUseCase.Execute(c.Request.Context(), &req, ipAddress, userAgent)
	if err != nil {
		h.logger.Error("Token refresh failed", zap.Error(err))

		if appErr, ok := err.(*errors.AppError); ok {
			switch appErr.Code {
			case "AUTH_001":
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": appErr.Message,
				})
			case "AUTH_003":
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": appErr.Message,
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Internal server error",
				})
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal server error",
			})
		}
		return
	}

	// Set new HTTP-only cookie for refresh token
	c.SetCookie(
		"refresh_token",
		response.Tokens.RefreshToken,
		int(response.ExpiresIn),
		"/",
		"",
		true, // secure
		true, // httpOnly
	)

	c.JSON(http.StatusOK, gin.H{
		"message": "Token refreshed successfully",
		"data": gin.H{
			"access_token": response.Tokens.AccessToken,
			"token_type":   response.Tokens.TokenType,
			"expires_in":   response.ExpiresIn,
		},
	})
}

// GetClientIP extracts client IP from request
func GetClientIP(c *gin.Context) string {
	// Check X-Forwarded-For header first
	if xff := c.GetHeader("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		if len(ips) > 0 {
			return strings.TrimSpace(ips[0])
		}
	}

	// Check X-Real-IP header
	if xri := c.GetHeader("X-Real-IP"); xri != "" {
		return xri
	}

	// Fallback to RemoteAddr
	return c.ClientIP()
}
