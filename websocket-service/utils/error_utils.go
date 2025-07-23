package utils

import (
	"fmt"
	"strconv"
)

type ServiceError struct {
	Err        error  // The underlying technical error
	Message    string // A user-friendly error message
	StatusCode int    // The HTTP status code to be returned to the client
}

func (se *ServiceError) Error() string {
	if se.Err != nil {
		return fmt.Sprintf("service error: %s (details: %v)", se.Message, se.Err)
	}
	return fmt.Sprintf("service error: %s", se.Message)
}

func NewServiceError(err error, message string, statusCode int) *ServiceError {
	return &ServiceError{
		Err:        err,
		Message:    message,
		StatusCode: statusCode,
	}
}

// ErrorResponseDTO defines the structure for generic error responses sent to the client.
type ErrorResponseDTO struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

func ParseIntOrDefault(s string, defaultValue int) (int, error) {
	if s == "" {
		return defaultValue, nil
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		// Return the error so the calling function can decide how to handle an invalid format
		return defaultValue, fmt.Errorf("invalid integer format for '%s': %w", s, err)
	}
	return val, nil
}

// ParseUintID parses a string to a uint ID.
// It's typically used for parsing IDs from URL paths or query parameters.
func ParseUintID(s string) (uint, error) {
	id, err := strconv.ParseUint(s, 10, 64) // Base 10, 64-bit unsigned integer
	if err != nil {
		// Return the error so the calling function can provide context
		return 0, fmt.Errorf("invalid ID format for '%s': %w", s, err)
	}
	return uint(id), nil
}
