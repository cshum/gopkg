package errw

import (
	"fmt"
	"net/http"
)

type Error struct {
	Code    string                 `json:"code"`
	Status  int                    `json:"-"`
	Message string                 `json:"message"`
	Extra   map[string]interface{} `json:"extra,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NotFound(message string) *Error {
	return &Error{
		Code:    "NotFoundError",
		Status:  http.StatusNotFound,
		Message: message,
	}
}

func Unauthorized(message string) *Error {
	return &Error{
		Code:    "UnauthorizedError",
		Status:  http.StatusUnauthorized,
		Message: message,
	}
}

func Validate(message string) *Error {
	return &Error{
		Code:    "ValidateError",
		Status:  http.StatusBadRequest,
		Message: message,
	}
}

func Timeout(message string) *Error {
	return &Error{
		Code:    "TimeoutError",
		Status:  http.StatusRequestTimeout,
		Message: message,
	}
}

func InternalServer(message string) *Error {
	return &Error{
		Code:    "InternalServerError",
		Status:  http.StatusInternalServerError,
		Message: message,
	}
}

func ValidateField(field, reason string) *Error {
	return &Error{
		Code:    "ValidateError",
		Status:  http.StatusBadRequest,
		Message: field + " " + reason,
		Extra: map[string]interface{}{
			"field":  field,
			"reason": reason,
		},
	}
}
