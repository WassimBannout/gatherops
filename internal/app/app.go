package app

import (
	"fmt"
	"net/http"

	"github.com/WassimBannout/gatherops/internal/config"
	"github.com/WassimBannout/gatherops/internal/httpapi"
	"github.com/WassimBannout/gatherops/internal/repository/postgres"
	"github.com/WassimBannout/gatherops/internal/security"
	"github.com/WassimBannout/gatherops/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NewHandler(cfg config.Config, db *pgxpool.Pool) (http.Handler, error) {
	tokenManager, err := security.NewTokenManager(cfg.JWTAccessSecret, cfg.AccessTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("create token manager: %w", err)
	}

	authService, err := service.NewAuthService(service.AuthServiceConfig{
		Users:         postgres.NewUserRepository(db),
		RefreshTokens: postgres.NewRefreshTokenRepository(db),
		Passwords:     security.NewPasswordHasher(0),
		Tokens:        tokenManager,
		RefreshTTL:    cfg.RefreshTokenTTL,
	})
	if err != nil {
		return nil, fmt.Errorf("create auth service: %w", err)
	}

	return httpapi.NewRouter(httpapi.Dependencies{
		DB:               db,
		ReadinessTimeout: cfg.ReadinessTimeout,
		Auth:             authService,
		Tokens:           tokenManager,
	}), nil
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
