package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	DefaultDatabaseURL     = "postgres://gatherops:gatherops@localhost:5433/gatherops?sslmode=disable"
	DefaultJWTAccessSecret = "development-only-change-me-gatherops-access-secret"
	minimumJWTSecretLength = 32
	defaultAccessTokenTTL  = 15 * time.Minute
	defaultRefreshTokenTTL = 30 * 24 * time.Hour
)

type Config struct {
	AppEnv           string
	HTTPPort         int
	DatabaseURL      string
	HTTPReadTimeout  time.Duration
	HTTPWriteTimeout time.Duration
	HTTPIdleTimeout  time.Duration
	ShutdownTimeout  time.Duration
	ReadinessTimeout time.Duration
	JWTAccessSecret  string
	AccessTokenTTL   time.Duration
	RefreshTokenTTL  time.Duration
}

func Load() (Config, error) {
	return load(os.LookupEnv)
}

func (c Config) HTTPAddress() string {
	return fmt.Sprintf(":%d", c.HTTPPort)
}

func (c Config) IsProduction() bool {
	return strings.EqualFold(c.AppEnv, "production")
}

type lookupFunc func(string) (string, bool)

func load(lookup lookupFunc) (Config, error) {
	appEnv := getString(lookup, "APP_ENV", "development")
	isProduction := strings.EqualFold(appEnv, "production")

	httpPort, err := getPort(lookup, "HTTP_PORT", 8080)
	if err != nil {
		return Config{}, err
	}

	databaseURL, ok := lookup("DATABASE_URL")
	databaseURL = strings.TrimSpace(databaseURL)
	if !ok || databaseURL == "" {
		if isProduction {
			return Config{}, errors.New("DATABASE_URL is required when APP_ENV=production")
		}
		databaseURL = DefaultDatabaseURL
	}

	jwtAccessSecret, ok := lookup("JWT_ACCESS_SECRET")
	jwtAccessSecret = strings.TrimSpace(jwtAccessSecret)
	if !ok || jwtAccessSecret == "" {
		if isProduction {
			return Config{}, errors.New("JWT_ACCESS_SECRET is required when APP_ENV=production")
		}
		jwtAccessSecret = DefaultJWTAccessSecret
	}
	if len(jwtAccessSecret) < minimumJWTSecretLength {
		return Config{}, fmt.Errorf("JWT_ACCESS_SECRET must be at least %d characters", minimumJWTSecretLength)
	}

	readTimeout, err := getDuration(lookup, "HTTP_READ_TIMEOUT", 5*time.Second)
	if err != nil {
		return Config{}, err
	}
	writeTimeout, err := getDuration(lookup, "HTTP_WRITE_TIMEOUT", 10*time.Second)
	if err != nil {
		return Config{}, err
	}
	idleTimeout, err := getDuration(lookup, "HTTP_IDLE_TIMEOUT", 60*time.Second)
	if err != nil {
		return Config{}, err
	}
	shutdownTimeout, err := getDuration(lookup, "SHUTDOWN_TIMEOUT", 10*time.Second)
	if err != nil {
		return Config{}, err
	}
	readinessTimeout, err := getDuration(lookup, "READINESS_TIMEOUT", 2*time.Second)
	if err != nil {
		return Config{}, err
	}
	accessTokenTTL, err := getDuration(lookup, "ACCESS_TOKEN_TTL", defaultAccessTokenTTL)
	if err != nil {
		return Config{}, err
	}
	refreshTokenTTL, err := getDuration(lookup, "REFRESH_TOKEN_TTL", defaultRefreshTokenTTL)
	if err != nil {
		return Config{}, err
	}

	return Config{
		AppEnv:           appEnv,
		HTTPPort:         httpPort,
		DatabaseURL:      databaseURL,
		HTTPReadTimeout:  readTimeout,
		HTTPWriteTimeout: writeTimeout,
		HTTPIdleTimeout:  idleTimeout,
		ShutdownTimeout:  shutdownTimeout,
		ReadinessTimeout: readinessTimeout,
		JWTAccessSecret:  jwtAccessSecret,
		AccessTokenTTL:   accessTokenTTL,
		RefreshTokenTTL:  refreshTokenTTL,
	}, nil
}

func getString(lookup lookupFunc, key, fallback string) string {
	value, ok := lookup(key)
	if !ok || strings.TrimSpace(value) == "" {
		return fallback
	}
	return strings.TrimSpace(value)
}

func getPort(lookup lookupFunc, key string, fallback int) (int, error) {
	raw := getString(lookup, key, strconv.Itoa(fallback))
	port, err := strconv.Atoi(raw)
	if err != nil || port < 1 || port > 65535 {
		return 0, fmt.Errorf("%s must be a valid TCP port", key)
	}
	return port, nil
}

func getDuration(lookup lookupFunc, key string, fallback time.Duration) (time.Duration, error) {
	raw := getString(lookup, key, fallback.String())
	duration, err := time.ParseDuration(raw)
	if err != nil || duration <= 0 {
		return 0, fmt.Errorf("%s must be a positive duration", key)
	}
	return duration, nil
}
