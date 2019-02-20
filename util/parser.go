package util

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
	validate = validator.New()
}

// QueryParse decode req query and validate struct from string map
func QueryParse(r *http.Request, dst interface{}) error {
	err := decoder.Decode(dst, r.URL.Query())
	if err != nil {
		return err
	}
	err = validate.Struct(dst)
	if err != nil {
		return err.(validator.ValidationErrors)
	}
	return nil
}

// BodyParse decode req body and validate struct from string map
func BodyParse(r *http.Request, dst interface{}) error {
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
