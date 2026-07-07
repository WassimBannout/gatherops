package integration_test

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/WassimBannout/gatherops/internal/config"
	"github.com/WassimBannout/gatherops/internal/database"
	"github.com/WassimBannout/gatherops/internal/domain"
	"github.com/WassimBannout/gatherops/internal/repository"
	"github.com/WassimBannout/gatherops/internal/repository/postgres"
	"github.com/WassimBannout/gatherops/internal/security"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestAuthRepositoriesPersistUsersAndRefreshTokens(t *testing.T) {
	if os.Getenv("GATHEROPS_INTEGRATION_TESTS") != "1" {
		t.Skip("set GATHEROPS_INTEGRATION_TESTS=1 to run repository integration tests")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	baseURL := integrationDatabaseURL()
	adminPool, err := pgxpool.New(ctx, baseURL)
	if err != nil {
		t.Fatalf("connect to postgres: %v", err)
	}
	defer adminPool.Close()

	schema := "test_" + strings.ReplaceAll(uuid.NewString(), "-", "_")
	if _, err := adminPool.Exec(ctx, fmt.Sprintf("CREATE SCHEMA %s", schema)); err != nil {
		t.Fatalf("create schema: %v", err)
	}
	defer func() {
		_, _ = adminPool.Exec(context.Background(), fmt.Sprintf("DROP SCHEMA IF EXISTS %s CASCADE", schema))
	}()

	migrationURL := databaseURLWithSearchPath(t, baseURL, schema, true)
	poolURL := databaseURLWithSearchPath(t, baseURL, schema, false)
	if err := database.RunMigrations(migrationURL, migrationsDir(t), database.MigrationUp, 0); err != nil {
		t.Fatalf("migrate up: %v", err)
	}

	pool, err := pgxpool.New(ctx, poolURL)
	if err != nil {
		t.Fatalf("connect to migrated schema: %v", err)
	}
	defer pool.Close()

	users := postgres.NewUserRepository(pool)
	refreshTokens := postgres.NewRefreshTokenRepository(pool)

	created, err := users.Create(ctx, domain.User{
		Email:        "ada@example.com",
		Name:         "Ada Lovelace",
		PasswordHash: "hashed-password",
	})
	if err != nil {
		t.Fatalf("create user: %v", err)
	}
	if created.ID == uuid.Nil {
		t.Fatal("expected user id")
	}
	if created.PasswordHash != "hashed-password" {
		t.Fatal("repository should scan password hash for service comparison")
	}

	_, err = users.Create(ctx, domain.User{
		Email:        "ada@example.com",
		Name:         "Duplicate",
		PasswordHash: "hashed-password",
	})
	if !errors.Is(err, repository.ErrConflict) {
		t.Fatalf("duplicate create error = %v, want ErrConflict", err)
	}

	byEmail, err := users.FindByEmail(ctx, "ada@example.com")
	if err != nil {
		t.Fatalf("find by email: %v", err)
	}
	if byEmail.ID != created.ID {
		t.Fatalf("user id = %s, want %s", byEmail.ID, created.ID)
	}

	rawRefreshToken := "raw-refresh-token"
	createdToken, err := refreshTokens.Create(ctx, domain.RefreshToken{
		UserID:    created.ID,
		TokenHash: security.HashRefreshToken(rawRefreshToken),
		ExpiresAt: time.Now().Add(time.Hour),
	})
	if err != nil {
		t.Fatalf("create refresh token: %v", err)
	}
	if createdToken.TokenHash == rawRefreshToken {
		t.Fatal("raw refresh token must not be stored")
	}

	foundToken, err := refreshTokens.FindByHash(ctx, security.HashRefreshToken(rawRefreshToken))
	if err != nil {
		t.Fatalf("find refresh token by hash: %v", err)
	}
	if foundToken.UserID != created.ID {
		t.Fatalf("refresh token user id = %s, want %s", foundToken.UserID, created.ID)
	}
	if foundToken.RevokedAt != nil {
		t.Fatal("new refresh token should not be revoked")
	}

	if err := refreshTokens.Revoke(ctx, foundToken.ID); err != nil {
		t.Fatalf("revoke refresh token: %v", err)
	}
	revokedToken, err := refreshTokens.FindByHash(ctx, security.HashRefreshToken(rawRefreshToken))
	if err != nil {
		t.Fatalf("find revoked refresh token: %v", err)
	}
	if revokedToken.RevokedAt == nil {
		t.Fatal("expected revoked_at to be set")
	}
}

func integrationDatabaseURL() string {
	baseURL := os.Getenv("TEST_DATABASE_URL")
	if baseURL == "" {
		baseURL = os.Getenv("DATABASE_URL")
	}
	if baseURL == "" {
		baseURL = config.DefaultDatabaseURL
	}
	return baseURL
}
