package app

import (
	"net/http"

	"github.com/WassimBannout/gatherops/internal/config"
	"github.com/WassimBannout/gatherops/internal/httpapi"
)

func NewHandler(cfg config.Config, db httpapi.Pinger) http.Handler {
	return httpapi.NewRouter(httpapi.Dependencies{
		DB:               db,
		ReadinessTimeout: cfg.ReadinessTimeout,
	})
}

func NewHTTPServer(cfg config.Config, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:         cfg.HTTPAddress(),
		Handler:      handler,
		ReadTimeout:  cfg.HTTPReadTimeout,
		WriteTimeout: cfg.HTTPWriteTimeout,
		IdleTimeout:  cfg.HTTPIdleTimeout,
	}
}
