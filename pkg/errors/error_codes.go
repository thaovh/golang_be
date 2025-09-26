package errors

// Standard error codes
const (
	// 1xxx - System Errors
	ErrSystemInternal    = "SYS_001" // Internal system error
	ErrSystemTimeout     = "SYS_002" // Request timeout
	ErrSystemUnavailable = "SYS_003" // Service unavailable

	// 2xxx - Validation Errors
	ErrValidationRequired = "VAL_001" // Required field missing
	ErrValidationFormat   = "VAL_002" // Invalid format
	ErrValidationRange    = "VAL_003" // Value out of range

	// 3xxx - Authentication/Authorization
	ErrAuthInvalidToken = "AUTH_001" // Invalid token
	ErrAuthExpiredToken = "AUTH_002" // Token expired
	ErrAuthInsufficient = "AUTH_003" // Insufficient permissions

	// 4xxx - Business Logic
	ErrBusinessNotFound = "BIZ_001" // Resource not found
	ErrBusinessConflict = "BIZ_002" // Business rule conflict
	ErrBusinessLimit    = "BIZ_003" // Business limit exceeded

	// 5xxx - External Dependencies
	ErrExternalTimeout     = "EXT_001" // External service timeout
	ErrExternalUnavailable = "EXT_002" // External service unavailable
	ErrExternalInvalid     = "EXT_003" // External service error
)
