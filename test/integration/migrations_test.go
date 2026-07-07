package integration_test

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/WassimBannout/gatherops/internal/config"
	"github.com/WassimBannout/gatherops/internal/database"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func TestCoreSchemaMigrationsApplyConstraintsAndRollback(t *testing.T) {
	if os.Getenv("GATHEROPS_INTEGRATION_TESTS") != "1" {
		t.Skip("set GATHEROPS_INTEGRATION_TESTS=1 to run migration integration tests")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	baseURL := os.Getenv("TEST_DATABASE_URL")
	if baseURL == "" {
		baseURL = os.Getenv("DATABASE_URL")
	}
	if baseURL == "" {
		baseURL = config.DefaultDatabaseURL
	}

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
	migrationsDir := migrationsDir(t)

	if err := database.RunMigrations(migrationURL, migrationsDir, database.MigrationUp, 0); err != nil {
		t.Fatalf("migrate up: %v", err)
	}

	pool, err := pgxpool.New(ctx, poolURL)
	if err != nil {
		t.Fatalf("connect to migrated schema: %v", err)
	}
	defer pool.Close()

	assertTablesExist(t, ctx, pool, schema, []string{
		"users",
		"refresh_tokens",
		"organizations",
		"organization_members",
		"events",
		"rsvps",
		"audit_logs",
	})
	assertCoreConstraints(t, ctx, pool)

	if err := database.RunMigrations(migrationURL, migrationsDir, database.MigrationDown, 0); err != nil {
		t.Fatalf("migrate down: %v", err)
	}

	assertNoTables(t, ctx, pool, schema)
}

func assertTablesExist(t *testing.T, ctx context.Context, pool *pgxpool.Pool, schema string, tables []string) {
	t.Helper()

	for _, table := range tables {
		var exists bool
		err := pool.QueryRow(ctx, `
            SELECT EXISTS (
                SELECT 1
                FROM information_schema.tables
                WHERE table_schema = $1 AND table_name = $2
            )
        `, schema, table).Scan(&exists)
		if err != nil {
			t.Fatalf("check table %s: %v", table, err)
		}
		if !exists {
			t.Fatalf("expected table %s.%s to exist", schema, table)
		}
	}
}

func assertNoTables(t *testing.T, ctx context.Context, pool *pgxpool.Pool, schema string) {
	t.Helper()

	var count int
	err := pool.QueryRow(ctx, `
        SELECT count(*)
        FROM information_schema.tables
        WHERE table_schema = $1
          AND table_name <> 'schema_migrations'
    `, schema).Scan(&count)
	if err != nil {
		t.Fatalf("count tables after down migration: %v", err)
	}
	if count != 0 {
		t.Fatalf("expected no tables after down migration, found %d", count)
	}
}

func assertCoreConstraints(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()

	var ownerID uuid.UUID
	err := pool.QueryRow(ctx, `
        INSERT INTO users (email, name, password_hash)
        VALUES ('owner@example.com', 'Owner', 'hashed-password')
        RETURNING id
    `).Scan(&ownerID)
	if err != nil {
		t.Fatalf("insert owner user: %v", err)
	}

	assertExecFails(t, ctx, pool, `
        INSERT INTO users (email, name, password_hash)
        VALUES ('OWNER@example.com', 'Owner', 'hashed-password')
    `)
	assertExecFails(t, ctx, pool, `
        INSERT INTO users (email, name, password_hash)
        VALUES ('owner@example.com', 'Other Owner', 'hashed-password')
    `)

	var organizationID uuid.UUID
	err = pool.QueryRow(ctx, `
        INSERT INTO organizations (name, slug, created_by)
        VALUES ('GatherOps Club', 'gatherops-club', $1)
        RETURNING id
    `, ownerID).Scan(&organizationID)
	if err != nil {
		t.Fatalf("insert organization: %v", err)
	}

	assertExecFails(t, ctx, pool, `
        INSERT INTO organizations (name, slug, created_by)
        VALUES ('Bad Slug', 'Bad Slug', $1)
    `, ownerID)

	_, err = pool.Exec(ctx, `
        INSERT INTO organization_members (organization_id, user_id, role)
        VALUES ($1, $2, 'owner')
    `, organizationID, ownerID)
	if err != nil {
		t.Fatalf("insert owner membership: %v", err)
	}
	assertExecFails(t, ctx, pool, `
        INSERT INTO organization_members (organization_id, user_id, role)
        VALUES ($1, $2, 'admin')
    `, organizationID, ownerID)

	var eventID uuid.UUID
	err = pool.QueryRow(ctx, `
        INSERT INTO events (organization_id, created_by, title, starts_at, ends_at, capacity, status)
        VALUES ($1, $2, 'Launch Night', now() + interval '1 hour', now() + interval '2 hours', 2, 'published')
        RETURNING id
    `, organizationID, ownerID).Scan(&eventID)
	if err != nil {
		t.Fatalf("insert event: %v", err)
	}
	assertExecFails(t, ctx, pool, `
        INSERT INTO events (organization_id, created_by, title, starts_at, ends_at, capacity)
        VALUES ($1, $2, 'Broken Event', now() + interval '2 hours', now() + interval '1 hour', 10)
    `, organizationID, ownerID)
	assertExecFails(t, ctx, pool, `
        INSERT INTO events (organization_id, created_by, title, starts_at, ends_at, capacity)
        VALUES ($1, $2, 'Broken Capacity', now() + interval '1 hour', now() + interval '2 hours', 0)
    `, organizationID, ownerID)

	_, err = pool.Exec(ctx, `
        INSERT INTO rsvps (event_id, user_id, status)
        VALUES ($1, $2, 'attending')
    `, eventID, ownerID)
	if err != nil {
		t.Fatalf("insert rsvp: %v", err)
	}
	assertExecFails(t, ctx, pool, `
        INSERT INTO rsvps (event_id, user_id, status)
        VALUES ($1, $2, 'attending')
    `, eventID, ownerID)
	assertExecFails(t, ctx, pool, `
        INSERT INTO rsvps (event_id, user_id, status)
        VALUES ($1, $2, 'maybe')
    `, eventID, ownerID)

	_, err = pool.Exec(ctx, `
        INSERT INTO audit_logs (organization_id, actor_user_id, action, entity_type, entity_id, metadata)
        VALUES ($1, $2, 'event.created', 'event', $3, '{"source":"test"}'::jsonb)
    `, organizationID, ownerID, eventID)
	if err != nil {
		t.Fatalf("insert audit log: %v", err)
	}
	assertExecFails(t, ctx, pool, `
        INSERT INTO audit_logs (organization_id, actor_user_id, action, entity_type, entity_id, metadata)
        VALUES ($1, $2, 'event.created', 'event', $3, '[]'::jsonb)
    `, organizationID, ownerID, eventID)
}

func assertExecFails(t *testing.T, ctx context.Context, pool *pgxpool.Pool, sql string, args ...any) {
	t.Helper()

	if _, err := pool.Exec(ctx, sql, args...); err == nil {
		t.Fatalf("expected SQL statement to fail: %s", sql)
	}
}

func databaseURLWithSearchPath(t *testing.T, rawURL, schema string, includeMigrateOptions bool) string {
	t.Helper()

	parsed, err := url.Parse(rawURL)
	if err != nil {
		t.Fatalf("parse database URL: %v", err)
	}

	query := parsed.Query()
	query.Set("options", "-c search_path="+schema+",public")
	if includeMigrateOptions {
		query.Set("x-migrations-table", "schema_migrations")
	}
	parsed.RawQuery = query.Encode()

	return parsed.String()
}

func migrationsDir(t *testing.T) string {
	t.Helper()

	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve test file path")
	}

	return filepath.Join(filepath.Dir(file), "..", "..", "migrations")
}
