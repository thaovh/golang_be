package handlers

import (
	"net/http"
	"strconv"

	"bm-staff/internal/usecases/user"
	"bm-staff/pkg/errors"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	createUserUseCase *user.CreateUserUseCase
	getUserUseCase    *user.GetUserUseCase
	updateUserUseCase *user.UpdateUserUseCase
	deleteUserUseCase *user.DeleteUserUseCase
	validator         *validator.Validate
	logger            *zap.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(
	createUserUseCase *user.CreateUserUseCase,
	getUserUseCase *user.GetUserUseCase,
	updateUserUseCase *user.UpdateUserUseCase,
	deleteUserUseCase *user.DeleteUserUseCase,
	validator *validator.Validate,
	logger *zap.Logger,
) *UserHandler {
	return &UserHandler{
		createUserUseCase: createUserUseCase,
		getUserUseCase:    getUserUseCase,
		updateUserUseCase: updateUserUseCase,
		deleteUserUseCase: deleteUserUseCase,
		validator:         validator,
		logger:            logger,
	}
}

// CreateUser handles POST /api/v1/users
// @Summary      Create a new user
// @Description  Create a new user with the provided information including password
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        user body user.CreateUserRequest true "User information"
// @Success      201 {object} map[string]interface{} "User created successfully"
// @Failure      400 {object} map[string]interface{} "Bad request - validation error"
// @Failure      409 {object} map[string]interface{} "Conflict - user already exists"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req user.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    errors.ErrValidationFormat,
				"message": "Invalid request format",
				"details": gin.H{"error": err.Error()},
			},
		})
		return
	}

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		h.logger.Error("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    errors.ErrValidationRequired,
				"message": "Validation failed",
				"details": gin.H{"error": err.Error()},
			},
		})
		return
	}

	// Execute use case
	resp, err := h.createUserUseCase.Execute(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data": resp.User,
	})
}

// GetUser handles GET /api/v1/users/:id
// @Summary      Get user by ID
// @Description  Retrieve a user by their unique identifier
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Success      200 {object} map[string]interface{} "User retrieved successfully"
// @Failure      400 {object} map[string]interface{} "Bad request - invalid user ID"
// @Failure      404 {object} map[string]interface{} "User not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userID := c.Param("id")
	req := &user.GetUserRequest{ID: userID}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		h.logger.Error("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    errors.ErrValidationFormat,
				"message": "Invalid user ID format",
				"details": gin.H{"error": err.Error()},
			},
		})
		return
	}

	// Execute use case
	resp, err := h.getUserUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp.User,
	})
}

// UpdateUser handles PUT /api/v1/users/:id
// @Summary      Update user
// @Description  Update an existing user's information
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Param        user body user.UpdateUserRequest true "Updated user information"
// @Success      200 {object} map[string]interface{} "User updated successfully"
// @Failure      400 {object} map[string]interface{} "Bad request - validation error"
// @Failure      404 {object} map[string]interface{} "User not found"
// @Failure      409 {object} map[string]interface{} "Conflict - username/email already exists"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")

	var req user.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		h.logger.Error("Failed to bind JSON", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    errors.ErrValidationFormat,
				"message": "Invalid request format",
				"details": gin.H{"error": err.Error()},
			},
		})
		return
	}

	// Set ID from URL parameter
	req.ID = userID

	// Validate request
	if err := h.validator.Struct(&req); err != nil {
		h.logger.Error("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    errors.ErrValidationRequired,
				"message": "Validation failed",
				"details": gin.H{"error": err.Error()},
			},
		})
		return
	}

	// Execute use case
	resp, err := h.updateUserUseCase.Execute(c.Request.Context(), &req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp.User,
	})
}

// DeleteUser handles DELETE /api/v1/users/:id
// @Summary      Delete user
// @Description  Soft delete a user by their ID
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        id path string true "User ID"
// @Success      200 {object} map[string]interface{} "User deleted successfully"
// @Failure      400 {object} map[string]interface{} "Bad request - invalid user ID"
// @Failure      404 {object} map[string]interface{} "User not found"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	req := &user.DeleteUserRequest{ID: userID}

	// Validate request
	if err := h.validator.Struct(req); err != nil {
		h.logger.Error("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"error": gin.H{
				"code":    errors.ErrValidationFormat,
				"message": "Invalid user ID format",
				"details": gin.H{"error": err.Error()},
			},
		})
		return
	}

	// Execute use case
	resp, err := h.deleteUserUseCase.Execute(c.Request.Context(), req)
	if err != nil {
		h.handleError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": resp,
	})
}

// ListUsers handles GET /api/v1/users
// @Summary      List users
// @Description  Retrieve a paginated list of users
// @Tags         users
// @Accept       json
// @Produce      json
// @Param        limit query int false "Number of users to return" default(10) minimum(1) maximum(100)
// @Param        offset query int false "Number of users to skip" default(0) minimum(0)
// @Success      200 {object} map[string]interface{} "Users retrieved successfully"
// @Failure      400 {object} map[string]interface{} "Bad request - invalid pagination parameters"
// @Failure      500 {object} map[string]interface{} "Internal server error"
// @Router       /users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Parse pagination parameters
	limitStr := c.DefaultQuery("limit", "10")
	offsetStr := c.DefaultQuery("offset", "0")

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit <= 0 {
		limit = 10
	}

	offset, err := strconv.Atoi(offsetStr)
	if err != nil || offset < 0 {
		offset = 0
	}

	// TODO: Implement ListUsers use case
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"users": []interface{}{},
			"pagination": gin.H{
				"limit":  limit,
				"offset": offset,
				"total":  0,
			},
		},
	})
}

// handleError handles application errors and returns appropriate HTTP responses
func (h *UserHandler) handleError(c *gin.Context, err error) {
	h.logger.Error("Handler error", zap.Error(err))

	if appErr, ok := err.(*errors.AppError); ok {
		statusCode := h.getStatusCodeFromErrorCode(appErr.Code)
		c.JSON(statusCode, gin.H{
			"error": gin.H{
				"code":      appErr.Code,
				"message":   appErr.Message,
				"details":   appErr.Details,
				"timestamp": appErr.Timestamp,
			},
		})
		return
	}

	// Generic error
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": gin.H{
			"code":    errors.ErrSystemInternal,
			"message": "Internal server error",
		},
	})
}

// getStatusCodeFromErrorCode maps error codes to HTTP status codes
func (h *UserHandler) getStatusCodeFromErrorCode(code string) int {
	switch code {
	case errors.ErrValidationRequired, errors.ErrValidationFormat, errors.ErrValidationRange:
		return http.StatusBadRequest
	case errors.ErrAuthInvalidToken, errors.ErrAuthExpiredToken, errors.ErrAuthInsufficient:
		return http.StatusUnauthorized
	case errors.ErrBusinessNotFound:
		return http.StatusNotFound
	case errors.ErrBusinessConflict:
		return http.StatusConflict
	case errors.ErrBusinessLimit:
		return http.StatusTooManyRequests
	case errors.ErrExternalTimeout, errors.ErrExternalUnavailable, errors.ErrExternalInvalid:
		return http.StatusBadGateway
	default:
		return http.StatusInternalServerError
	}
}
