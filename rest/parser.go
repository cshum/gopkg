package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/gorilla/schema"
	"gopkg.in/go-playground/validator.v9"
)

var decoder *schema.Decoder
var Validate *validator.Validate

func init() {
	// schema decoder
	decoder = schema.NewDecoder()
	decoder.SetAliasTag("json")
	decoder.IgnoreUnknownKeys(true)
	// validator
	Validate = validator.New()
}

// ParamInt parse int from chi URL param
func ParamInt(r *http.Request, key string) int {
	id, _ := strconv.Atoi(chi.URLParam(r, key))
	return id
}

// ParseQuery decode req query and validate struct from string map
func ParseQuery(r *http.Request, dst interface{}) error {
	if err := decoder.Decode(dst, r.URL.Query()); err != nil {
		return err
	}
	if err := Validate.Struct(dst); err != nil {
		return err
	}
	return nil
}

// ParseBody decode req body and validate struct from string map
func ParseBody(r *http.Request, dst interface{}, maxMemory int64) error {
	if IsBodyJSON(r) {
		reader := io.LimitReader(r.Body, maxMemory)
		body, err := ioutil.ReadAll(reader)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(body, dst); err != nil {
			return err
		}
		// restore body state
		r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	} else {
		if err := r.ParseMultipartForm(maxMemory); err != nil && err != http.ErrNotMultipart {
			return err
		}
		if err := decoder.Decode(dst, r.PostForm); err != nil {
			return err
		}
	}
	if err := Validate.Struct(dst); err != nil {
		return err
	}
	return nil
}

func IsBodyJSON(r *http.Request) bool {
	if r.Header.Get("Content-Type") != "application/json" {
		return false
	}
	if r.Header.Get("Content-Length") == "0" {
		return false
	}
	return true
}
