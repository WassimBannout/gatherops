package httpapi

import (
	"net/http"
	"time"

	"github.com/WassimBannout/gatherops/internal/security"
	"github.com/WassimBannout/gatherops/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Dependencies struct {
	DB               Pinger
	ReadinessTimeout time.Duration
	Auth             *service.AuthService
	Tokens           security.TokenManager
}

func NewRouter(deps Dependencies) http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)

	health := NewHealthHandler(deps.DB, deps.ReadinessTimeout)
	router.Get("/healthz", health.Healthz)
	router.Get("/readyz", health.Readyz)

	if deps.Auth != nil {
		authHandler := NewAuthHandler(deps.Auth)
		router.Route("/api/v1", func(api chi.Router) {
			api.Post("/auth/register", authHandler.Register)
			api.Post("/auth/login", authHandler.Login)
			api.Post("/auth/refresh", authHandler.Refresh)

			api.Group(func(protected chi.Router) {
				protected.Use(AuthMiddleware(deps.Tokens))
				protected.Post("/auth/logout", authHandler.Logout)
				protected.Get("/me", authHandler.Me)
			})
		})
	}

	return router
}
