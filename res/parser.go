package res

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"

	"github.com/gorilla/schema"
	validator "gopkg.in/go-playground/validator.v9"
)

var decoder *schema.Decoder
var validate *validator.Validate

func init() {
	// schema decoder
	decoder = schema.NewDecoder()
	decoder.SetAliasTag("json")
	decoder.IgnoreUnknownKeys(true)
	validate = validator.New()
}

// ParamInt parse int from chi URL param
func ParamInt(r *http.Request, key string) (int, error) {
	return strconv.Atoi(chi.URLParam(r, key))
}

// ParseQuery decode req query and validate struct from string map
func ParseQuery(r *http.Request, dst interface{}) error {
	if err := decoder.Decode(dst, r.URL.Query()); err != nil {
		return err
	}
	if err := validate.Struct(dst); err != nil {
		return err
	}
	return nil
}

const maxMemory = int64(10 << 20) // 10mb

// ParseBody decode req body and validate struct from string map
func ParseBody(r *http.Request, dst interface{}) error {
	if isBodyJSON(r) {
		reader := io.LimitReader(r.Body, maxMemory) // 10MB
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
	if err := validate.Struct(dst); err != nil {
		return err
	}
	return nil
}

func isBodyJSON(r *http.Request) bool {
	if r.Header.Get("Content-Type") != "application/json" {
		return false
	}
	if r.Header.Get("Content-Length") == "0" {
		return false
	}
	return true
}
