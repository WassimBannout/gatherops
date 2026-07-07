package config

import (
	"testing"
	"time"
)

func TestLoadUsesDevelopmentDefaults(t *testing.T) {
	cfg, err := load(mapLookup(nil))
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.AppEnv != "development" {
		t.Fatalf("AppEnv = %q, want development", cfg.AppEnv)
	}
	if cfg.HTTPAddress() != ":8080" {
		t.Fatalf("HTTPAddress() = %q, want :8080", cfg.HTTPAddress())
	}
	if cfg.DatabaseURL != DefaultDatabaseURL {
		t.Fatalf("DatabaseURL = %q, want default", cfg.DatabaseURL)
	}
	if cfg.JWTAccessSecret != DefaultJWTAccessSecret {
		t.Fatalf("JWTAccessSecret = %q, want development default", cfg.JWTAccessSecret)
	}
	if cfg.AccessTokenTTL != 15*time.Minute {
		t.Fatalf("AccessTokenTTL = %s, want 15m", cfg.AccessTokenTTL)
	}
	if cfg.RefreshTokenTTL != 30*24*time.Hour {
		t.Fatalf("RefreshTokenTTL = %s, want 720h", cfg.RefreshTokenTTL)
	}
	if cfg.ReadinessTimeout != 2*time.Second {
		t.Fatalf("ReadinessTimeout = %s, want 2s", cfg.ReadinessTimeout)
	}
}

func TestLoadRequiresProductionSecrets(t *testing.T) {
	_, err := load(mapLookup(map[string]string{
		"APP_ENV":      "production",
		"DATABASE_URL": "postgres://example",
	}))
	if err == nil {
		t.Fatal("expected production config without JWT_ACCESS_SECRET to fail")
	}
}

func TestLoadRequiresDatabaseURLInProduction(t *testing.T) {
	_, err := load(mapLookup(map[string]string{
		"APP_ENV":           "production",
		"JWT_ACCESS_SECRET": "abcdefghijklmnopqrstuvwxyz1234567890",
	}))
	if err == nil {
		t.Fatal("expected production config without DATABASE_URL to fail")
	}
}

func TestLoadRejectsShortJWTSecret(t *testing.T) {
	_, err := load(mapLookup(map[string]string{"JWT_ACCESS_SECRET": "too-short"}))
	if err == nil {
		t.Fatal("expected short JWT_ACCESS_SECRET to fail")
	}
}

func TestLoadParsesExplicitValues(t *testing.T) {
	cfg, err := load(mapLookup(map[string]string{
		"APP_ENV":           "test",
		"HTTP_PORT":         "9090",
		"DATABASE_URL":      "postgres://example",
		"JWT_ACCESS_SECRET": "abcdefghijklmnopqrstuvwxyz1234567890",
		"ACCESS_TOKEN_TTL":  "10m",
		"REFRESH_TOKEN_TTL": "168h",
		"HTTP_READ_TIMEOUT": "3s",
		"READINESS_TIMEOUT": "750ms",
	}))
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	if cfg.AppEnv != "test" {
		t.Fatalf("AppEnv = %q, want test", cfg.AppEnv)
	}
	if cfg.HTTPAddress() != ":9090" {
		t.Fatalf("HTTPAddress() = %q, want :9090", cfg.HTTPAddress())
	}
	if cfg.DatabaseURL != "postgres://example" {
		t.Fatalf("DatabaseURL = %q, want explicit value", cfg.DatabaseURL)
	}
	if cfg.AccessTokenTTL != 10*time.Minute {
		t.Fatalf("AccessTokenTTL = %s, want 10m", cfg.AccessTokenTTL)
	}
	if cfg.RefreshTokenTTL != 168*time.Hour {
		t.Fatalf("RefreshTokenTTL = %s, want 168h", cfg.RefreshTokenTTL)
	}
	if cfg.HTTPReadTimeout != 3*time.Second {
		t.Fatalf("HTTPReadTimeout = %s, want 3s", cfg.HTTPReadTimeout)
	}
	if cfg.ReadinessTimeout != 750*time.Millisecond {
		t.Fatalf("ReadinessTimeout = %s, want 750ms", cfg.ReadinessTimeout)
	}
}

func TestLoadRejectsInvalidPort(t *testing.T) {
	_, err := load(mapLookup(map[string]string{"HTTP_PORT": "70000"}))
	if err == nil {
		t.Fatal("expected invalid HTTP_PORT to fail")
	}
}

func TestLoadRejectsInvalidDuration(t *testing.T) {
	_, err := load(mapLookup(map[string]string{"HTTP_READ_TIMEOUT": "soon"}))
	if err == nil {
		t.Fatal("expected invalid duration to fail")
	}
}

func mapLookup(values map[string]string) lookupFunc {
	return func(key string) (string, bool) {
		value, ok := values[key]
		return value, ok
	}
}
