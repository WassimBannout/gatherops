package httpapi

import (
	"context"
	"net/http"
	"time"
)

type Pinger interface {
	Ping(context.Context) error
}

type HealthHandler struct {
	db               Pinger
	readinessTimeout time.Duration
}

func NewHealthHandler(db Pinger, readinessTimeout time.Duration) HealthHandler {
	if readinessTimeout <= 0 {
		readinessTimeout = 2 * time.Second
	}

	return HealthHandler{
		db:               db,
		readinessTimeout: readinessTimeout,
	}
}

func (h HealthHandler) Healthz(w http.ResponseWriter, _ *http.Request) {
	WriteJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h HealthHandler) Readyz(w http.ResponseWriter, r *http.Request) {
	if h.db == nil {
		WriteError(w, r, http.StatusServiceUnavailable, "database_unavailable", "Database is not reachable", nil)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), h.readinessTimeout)
	defer cancel()

	if err := h.db.Ping(ctx); err != nil {
		WriteError(w, r, http.StatusServiceUnavailable, "database_unavailable", "Database is not reachable", nil)
		return
	}

	WriteJSON(w, http.StatusOK, map[string]any{
		"status": "ready",
		"dependencies": map[string]string{
			"database": "ok",
		},
	})
}
