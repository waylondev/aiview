// Package aiverr provides unified error types for multi-platform API clients.
package aiverr

import (
	"fmt"
)

// Error codes for standardized error handling.
const (
	CodeNotAuthenticated = "not_authenticated"
	CodeRateLimited      = "rate_limited"
	CodeAPIError         = "api_error"
	CodeNetworkError     = "network_error"
	CodeParseError       = "parse_error"
	CodeNotFound         = "not_found"
	CodeForbidden        = "forbidden"
	CodeInvalidInput     = "invalid_input"
)

// PlatformError represents a structured error with platform context.
type PlatformError struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Details *ErrorDetails `json:"details,omitempty"`
}

// ErrorDetails provides additional context about the error.
type ErrorDetails struct {
	Platform   string `json:"platform,omitempty"`
	Command    string `json:"command,omitempty"`
	HTTPStatus int    `json:"http_status,omitempty"`
	Suggestion string `json:"suggestion,omitempty"`
}

// Error implements the error interface.
func (e *PlatformError) Error() string {
	if e.Details != nil && e.Details.Suggestion != "" {
		return fmt.Sprintf("%s: %s (suggestion: %s)", e.Code, e.Message, e.Details.Suggestion)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// New creates a new PlatformError with the given code and message.
func New(code, message string) *PlatformError {
	return &PlatformError{
		Code:    code,
		Message: message,
	}
}

// WithDetails adds details to the error.
func (e *PlatformError) WithDetails(details *ErrorDetails) *PlatformError {
	e.Details = details
	return e
}

// WithPlatform sets the platform in error details.
func (e *PlatformError) WithPlatform(platform string) *PlatformError {
	if e.Details == nil {
		e.Details = &ErrorDetails{}
	}
	e.Details.Platform = platform
	return e
}

// WithCommand sets the command in error details.
func (e *PlatformError) WithCommand(command string) *PlatformError {
	if e.Details == nil {
		e.Details = &ErrorDetails{}
	}
	e.Details.Command = command
	return e
}

// WithHTTPStatus sets the HTTP status in error details.
func (e *PlatformError) WithHTTPStatus(status int) *PlatformError {
	if e.Details == nil {
		e.Details = &ErrorDetails{}
	}
	e.Details.HTTPStatus = status
	return e
}

// WithSuggestion sets the suggestion in error details.
func (e *PlatformError) WithSuggestion(suggestion string) *PlatformError {
	if e.Details == nil {
		e.Details = &ErrorDetails{}
	}
	e.Details.Suggestion = suggestion
	return e
}

// NotAuthenticated creates a not_authenticated error.
func NotAuthenticated(platform, message string) *PlatformError {
	return New(CodeNotAuthenticated, message).
		WithPlatform(platform).
		WithSuggestion(fmt.Sprintf("Please login with: aiview %s login --cookie", platform))
}

// RateLimited creates a rate_limited error.
func RateLimited(platform, message string) *PlatformError {
	return New(CodeRateLimited, message).
		WithPlatform(platform).
		WithSuggestion("Please wait a moment and try again")
}

// NotFound creates a not_found error.
func NotFound(platform, message string) *PlatformError {
	return New(CodeNotFound, message).WithPlatform(platform)
}

// APIError creates an api_error.
func APIError(platform, message string) *PlatformError {
	return New(CodeAPIError, message).WithPlatform(platform)
}

// NetworkError creates a network_error.
func NetworkError(platform, message string) *PlatformError {
	return New(CodeNetworkError, message).WithPlatform(platform)
}

// ParseError creates a parse_error.
func ParseError(platform, message string) *PlatformError {
	return New(CodeParseError, message).WithPlatform(platform)
}

// Forbidden creates a forbidden error.
func Forbidden(platform, message string) *PlatformError {
	return New(CodeForbidden, message).WithPlatform(platform)
}

// InvalidInput creates an invalid_input error.
func InvalidInput(platform, message string) *PlatformError {
	return New(CodeInvalidInput, message).WithPlatform(platform)
}

// IsPlatformError checks if an error is a PlatformError.
func IsPlatformError(err error) (*PlatformError, bool) {
	if err == nil {
		return nil, false
	}
	if pe, ok := err.(*PlatformError); ok {
		return pe, true
	}
	return nil, false
}

// Wrap wraps an existing error as a PlatformError.
func Wrap(err error, code, platform string) *PlatformError {
	if err == nil {
		return nil
	}
	if pe, ok := err.(*PlatformError); ok {
		return pe
	}
	return New(code, err.Error()).WithPlatform(platform)
}
