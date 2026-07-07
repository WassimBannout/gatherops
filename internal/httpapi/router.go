package httpapi

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Dependencies struct {
	DB               Pinger
	ReadinessTimeout time.Duration
}

func NewRouter(deps Dependencies) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)

	health := NewHealthHandler(deps.DB, deps.ReadinessTimeout)
	router.Get("/healthz", health.Healthz)
	router.Get("/readyz", health.Readyz)

	return router
}
