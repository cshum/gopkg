package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cshum/gopkg/errw"
	"github.com/cshum/gopkg/paginator"
)

// JSON write json to http response writer
func JSON(w http.ResponseWriter, status int, data interface{}) {
	bytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(bytes)))
	w.WriteHeader(status)
	_, _ = w.Write(bytes)
}

// Response standard response
type Response struct {
	Success    bool                  `json:"success"`
	Status     int                   `json:"status"`
	Data       interface{}           `json:"data,omitempty"`
	Error      *errw.Error           `json:"error,omitempty"`
	Pagination *paginator.Pagination `json:"pagination,omitempty"`
}

// Success response
func Success(w http.ResponseWriter, data interface{}) {
	SuccessPaginated(w, data, nil)
}

// SuccessPaginated response
func SuccessPaginated(w http.ResponseWriter, data interface{}, p *paginator.Pagination) {
	JSON(w, http.StatusOK, Response{
		Status:     http.StatusOK,
		Success:    true,
		Data:       data,
		Pagination: p,
	})
}

// Fail response
func Fail(w http.ResponseWriter, status int, code string, message string) {
	JSON(w, status, Response{
		Status:  status,
		Success: false,
		Error: &errw.Error{
			Code:    code,
			Message: message,
		},
	})
}

// FailOk fail response status ok
func FailOk(w http.ResponseWriter, code string, message string) {
	Fail(w, http.StatusOK, code, message)
}

// FailValidate 400 ValidateError
func FailValidate(w http.ResponseWriter, message string) {
	Fail(w, http.StatusBadRequest, "ValidateError", message)
}

// FailNotFound 400 NotFoundError
func FailNotFound(w http.ResponseWriter, message string) {
	Fail(w, http.StatusNotFound, "NotFoundError", message)
}

// FailUnauthorized 401 UnauthorizedError
func FailUnauthorized(w http.ResponseWriter, message string) {
	Fail(w, http.StatusUnauthorized, "UnauthorizedError", message)
}
