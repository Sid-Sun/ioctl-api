package middleware

import (
	"net/http"
	"strings"

	"github.com/rs/cors"
	"github.com/sid-sun/ioctl-api/config"
)

func WithCors(cfg *config.HTTPServerConfig) func(h http.Handler) http.Handler {
	urls := strings.Split(cfg.CORS, ",")
	handler := cors.New(cors.Options{
		AllowedOrigins: urls,
		AllowedMethods: []string{"GET", "POST", "PUT", "OPTIONS"},
		AllowedHeaders: []string{"Origin", "X-Requested-With", "Content-Type", "Accept", "Ephemeral"},
		MaxAge:         30 * 60, // 30 mins of preflight caching
	}).Handler

	return handler
}

func WithMaxBodyReader(cfg *config.HTTPServerConfig) func(h http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			req.Body = http.MaxBytesReader(w, req.Body, cfg.MaxBodySize)
			next.ServeHTTP(w, req)
		})
	}
}
