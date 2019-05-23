package res

import (
	"encoding/json"
	"github.com/cshum/gopkg/paginator"
	"net/http"
	"strconv"
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
	Ok         bool                  `json:"ok"`
	Status     int                   `json:"status"`
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
	JSON(w, http.StatusOK, Response{
		Status:     http.StatusOK,
		Ok:         true,
		Data:       data,
		Pagination: p,
	})
}

// Fail response
func Fail(w http.ResponseWriter, status int, code string, message string) {
	JSON(w, status, Response{
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
