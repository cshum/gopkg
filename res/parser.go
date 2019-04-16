package res

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

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

// ParseQuery decode req query and validate struct from string map
func ParseQuery(r *http.Request, dst interface{}) error {
	if err := decoder.Decode(dst, r.URL.Query()); err != nil {
		return err
	}
	if err := validate.Struct(dst); err != nil {
		return err.(validator.ValidationErrors)
	}
	return nil
}

// ParseBody decode req body and validate struct from string map
func ParseBody(r *http.Request, dst interface{}) error {
	if isBodyJSON(r) {
		reader := io.LimitReader(r.Body, int64(10<<20)) // 10MB
		body, err := ioutil.ReadAll(reader)
		if err != nil {
			return err
		}
		err = json.Unmarshal(body, &dst)
		if err != nil {
			return err
		}
	} else {
		err := r.ParseForm()
		if err != nil {
			return err
		}
		err = decoder.Decode(dst, r.PostForm)
		if err != nil {
			return err
		}
	}
	err := validate.Struct(dst)
	if err != nil {
		return err.(validator.ValidationErrors)
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
