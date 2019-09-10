package errw

import (
	"fmt"
	"net/http"
)

type Error struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Status  int                    `json:"status"`
	Extra   map[string]interface{} `json:"extra,omitempty"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NotFound(message string) error {
	return &Error{
		Code:    "NotFoundError",
		Status:  http.StatusNotFound,
		Message: message,
	}
}

func Unauthorized(message string) error {
	return &Error{
		Code:    "UnauthorizedError",
		Status:  http.StatusUnauthorized,
		Message: message,
	}
}

func Validate(message string) error {
	return &Error{
		Code:    "ValidateError",
		Status:  http.StatusBadRequest,
		Message: message,
	}
}

func ValidateField(field, reason string) error {
	return &Error{
		Code:    "ValidateError",
		Status:  http.StatusBadRequest,
		Message: fmt.Sprintf("%s %s", field, reason),
		Extra: map[string]interface{}{
			"field":  field,
			"reason": reason,
		},
	}
}
