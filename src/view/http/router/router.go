package router

import (
	"net/http"

	"github.com/fitant/xbin-api/config"
	"github.com/fitant/xbin-api/src/service"
	"github.com/fitant/xbin-api/src/view/http/handler/snippet"
	"github.com/fitant/xbin-api/src/view/http/middleware"
	"github.com/go-chi/chi/v5"
)

func New(svc service.Service, cfg *config.HTTPServerConfig) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.WithCors(cfg))

	r.Route(cfg.Enpoint, func(sr chi.Router) {
		sr.With(middleware.WithIngestion()).Put(
			"/{file}", snippet.Create(svc, cfg))
		sr.With(middleware.WithIngestion()).Put(
			"/", snippet.Create(svc, cfg))
		sr.With(middleware.WithIngestion()).Post(
			"/", snippet.Create(svc, cfg))
		sr.Post("/e2e/{snippetID}", snippet.CreateE2E(svc, cfg))
		sr.Get("/r/{snippetID}", snippet.Get(svc, "raw"))
		sr.Get("/{snippetID}", snippet.Get(svc, "json"))
		sr.Get("/", func(w http.ResponseWriter, r *http.Request) {
			return
		})
	})

	return r
}
