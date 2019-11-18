package mw

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

// RecoverHandler recovers from panic, passes error to handler
func RecoverHandler(
	handler func(w http.ResponseWriter, r *http.Request, err error),
) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					err, ok := rvr.(error)
					if !ok {
						err = errors.New(fmt.Sprintf("%v", rvr))
					}
					handler(w, r, err)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(
	get func(path string, handle http.HandlerFunc),
	path, dirpath string,
) {
	if strings.ContainsAny(path, "{}*") {
		panic(errors.New("FileServer does not permit URL parameters"))
	}
	fs := http.StripPrefix(path, http.FileServer(http.Dir(dirpath)))
	if path != "/" && path[len(path)-1] != '/' {
		get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"
	get(path, fs.ServeHTTP)
}
