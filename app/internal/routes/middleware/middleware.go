package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/owicca/tr/internal/infra"

	"go.uber.org/zap"
)

func init() {
	LoadMd(infra.S)
}

// Load middlewares
func LoadMd(srv *infra.Server) {
	srv.Router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		LogRequest(w, r)
		template404Path := "front/404"
		if strings.HasPrefix(r.URL.Path, "/admin") {
			template404Path = "back/404"
		}
		srv.HTML(w, r, http.StatusNotFound, template404Path, nil)
		return
	})
	srv.Router.Use(setCSPHeader)

	srv.Router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			LogRequest(w, r)
			next.ServeHTTP(w, r)
		})
	})

	srv.Router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// middleware here
			next.ServeHTTP(w, r)
		})
	})
}

func setCSPHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		csp := []string{
			"default-src *",
			//"script-src 'self'",
			//"connect-src 'self'",
			//"img-src 'self'",
			//"style-src 'self'",
			//"base-uri 'self'",
			//"form-action 'self'",
		}

		header.Set("Content-Security-Policy", strings.Join(csp, ";"))
		// header.Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)
	})
}

func LogRequest(w http.ResponseWriter, r *http.Request) {
	timestamp := time.Now().Unix()
	url := r.RequestURI
	remote_addr := r.RemoteAddr
	method := r.Method

	logMsg := fmt.Sprintf("%s %s %s", remote_addr, method, url)
	zap.L().Info(logMsg, zap.Int64("timestamp", timestamp))
}
