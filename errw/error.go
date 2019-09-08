package errw

import (
	"fmt"
	"net/http"
)

type Error struct {
	Code    string                 `json:"code,omitempty"`
	Message string                 `json:"message,omitempty"`
	Status  int                    `json:"status"`
	Extra   map[string]interface{} `json:"extra"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func NotFound(message string) *Error {
	return &Error{
		Status:  http.StatusNotFound,
		Code:    "NotFoundError",
		Message: message,
	}
}

func Unauthorized(message string) *Error {
	return &Error{
		Status:  http.StatusUnauthorized,
		Code:    "UnauthorizedError",
		Message: message,
	}
}
