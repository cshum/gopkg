package util

import (
	"errors"
	"github.com/cshum/gopkg/res"
	"net/http"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// RecoverHandler recovers from panic, log a sentry and response 500
func RecoverHandler(handlers ...func(error)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					for _, handler := range handlers {
						handler(rvr.(error))
					}
					res.Fail(w, http.StatusInternalServerError, "InternalServerError", "internal server error")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// JWTVerifier verify jwt token
func JWTVerifier(ja *jwtauth.JWTAuth) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return jwtauth.Verify(
			ja,
			jwtauth.TokenFromHeader,
			func(r *http.Request) string {
				return r.URL.Query().Get("token")
			},
			func(r *http.Request) string {
				return r.Header.Get("x-access-token")
			},
		)(next)
	}
}

// JWTAuth middleware of jwtauth with custom response
func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())
		if err != nil {
			res.FailUnauthorized(w, "unauthorized")
			return
		}
		if token == nil || !token.Valid {
			res.FailUnauthorized(w, "unauthorized")
			return
		}
		// authorized
		next.ServeHTTP(w, r)
	})
}

// AccessLogHandler zap logger middleware
func AccessLogHandler(log func(string, ...zap.Field)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			next.ServeHTTP(ww, r)
			latency := time.Since(start)
			fields := []zapcore.Field{
				zap.Int("status", ww.Status()),
				zap.Duration("took", latency),
				zap.Int64("latency", latency.Nanoseconds()),
				zap.String("remote", r.RemoteAddr),
				zap.String("request", r.RequestURI),
				zap.String("method", r.Method),
			}
			log("access", fields...)
		})
	}
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, dirname string) {
	if strings.ContainsAny(path, "{}*") {
		panic(errors.New("FileServer does not permit URL parameters"))
	}
	_, callerFileName, _, _ := runtime.Caller(1)
	curr, err := filepath.Abs(filepath.Dir(callerFileName))
	if err != nil {
		panic(err)
	}
	filesDir := filepath.Join(curr, dirname)
	fs := http.StripPrefix(path, http.FileServer(http.Dir(filesDir)))
	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, http.HandlerFunc(fs.ServeHTTP))
}
