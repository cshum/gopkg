package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/cshum/gopkg/errw"
	"github.com/cshum/gopkg/paginator"
	"gopkg.in/go-playground/validator.v9"
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

func FailW(w http.ResponseWriter, err *errw.Error) {
	JSON(w, err.Status, Response{
		Status:  err.Status,
		Success: false,
		Error:   err,
	})
}

func Fail(w http.ResponseWriter, err error) {
	// errw
	if err, ok := err.(*errw.Error); ok {
		FailW(w, err)
		return
	}
	// validator errors
	if _, ok := err.(validator.ValidationErrors); ok {
		FailW(w, errw.Validate(err.Error()))
		return
	}
	// all others 500
	FailW(w, errw.InternalServer(err.Error()))
}

func FailNotFound(w http.ResponseWriter, message string) {
	FailW(w, errw.NotFound(message))
}

func FailUnauthorized(w http.ResponseWriter, message string) {
	FailW(w, errw.Unauthorized(message))
}

func FailValidate(w http.ResponseWriter, message string) {
	FailW(w, errw.Validate(message))
}

func FailValidateField(w http.ResponseWriter, field, reason string) {
	FailW(w, errw.ValidateField(field, reason))
}

func FailInternalServer(w http.ResponseWriter, message string) {
	FailW(w, errw.InternalServer(message))
}

func FailTimeout(w http.ResponseWriter, message string) {
	FailW(w, errw.Timeout(message))
}
