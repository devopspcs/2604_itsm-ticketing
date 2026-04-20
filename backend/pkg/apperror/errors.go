package apperror

import "fmt"

type AppError struct {
	Code       string                 `json:"error_code"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
	HTTPStatus int                    `json:"-"`
}

func (e *AppError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func New(code, message string, httpStatus int) *AppError {
	return &AppError{Code: code, Message: message, HTTPStatus: httpStatus}
}

func (e *AppError) WithDetails(details map[string]interface{}) *AppError {
	return &AppError{
		Code:       e.Code,
		Message:    e.Message,
		Details:    details,
		HTTPStatus: e.HTTPStatus,
	}
}

var (
	ErrInvalidCredentials = New("INVALID_CREDENTIALS", "Invalid email or password", 401)
	ErrTokenExpired       = New("TOKEN_EXPIRED", "Token has expired", 401)
	ErrTokenInvalid       = New("TOKEN_INVALID", "Token is invalid", 401)
	ErrForbidden          = New("FORBIDDEN", "Access denied", 403)
	ErrNotFound           = New("NOT_FOUND", "Resource not found", 404)
	ErrConflict           = New("CONFLICT", "Resource already exists", 409)
	ErrValidation         = New("VALIDATION_ERROR", "Validation failed", 400)
	ErrRateLimit           = New("RATE_LIMIT_EXCEEDED", "Rate limit exceeded", 429)
	ErrRateLimitExceeded   = New("RATE_LIMIT_EXCEEDED", "Too many requests. Please try again later", 429)
	ErrInternal            = New("INTERNAL_ERROR", "Internal server error", 500)
	ErrUserInactive       = New("USER_INACTIVE", "User account is inactive", 401)
	ErrHasDependencies    = New("HAS_DEPENDENCIES", "Cannot delete: has dependent records", 400)
	ErrInvalidHierarchy   = New("INVALID_HIERARCHY", "Invalid organizational hierarchy", 400)
)
