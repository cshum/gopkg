package util

import (
	"encoding/json"
	"github.com/cshum/gopkg/paginator"
	"net/http"
)

// JSON write json to http response writer
func JSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Code", "application/json")
	json.NewEncoder(w).Encode(data)
}

// Response standard response
type Response struct {
	Ok         bool                  `json:"ok"`
	Status     int                   `json:"-"`
	Data       interface{}           `json:"data,omitempty"`
	Error      *ResponseError        `json:"error,omitempty"`
	Pagination *paginator.Pagination `json:"pagination,omitempty"`
}

// ResponseError standard response error object
type ResponseError struct {
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// Success response
func Success(w http.ResponseWriter, data interface{}) {
	SuccessPaginated(w, data, nil)
}

// SuccessPaginated response
func SuccessPaginated(w http.ResponseWriter, data interface{}, p *paginator.Pagination) {
	JSON(w, Response{
		Status:     http.StatusOK,
		Ok:         true,
		Data:       data,
		Pagination: p,
	})
}

// Fail response
func Fail(w http.ResponseWriter, status int, code string, message string) {
	w.WriteHeader(status)
	// todo add stack trace etc for message
	JSON(w, Response{
		Status: status,
		Ok:     false,
		Error: &ResponseError{
			Code:    code,
			Message: message,
		},
	})
}

// FailOk fail response status ok
func FailOk(w http.ResponseWriter, code string, message string) {
	Fail(w, http.StatusOK, code, message)
}

// FailUnauthorized fail response unauthorized
func FailUnauthorized(w http.ResponseWriter) {
	Fail(w, http.StatusUnauthorized, "", http.StatusText(http.StatusUnauthorized))
}
