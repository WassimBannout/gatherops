package database

import (
	"path/filepath"
	"testing"
)

func TestMigrationSourceURLRejectsEmptyPath(t *testing.T) {
	_, err := MigrationSourceURL("")
	if err == nil {
		t.Fatal("expected empty migration path to fail")
	}
}

func TestMigrationSourceURLBuildsRelativeFileURL(t *testing.T) {
	got, err := MigrationSourceURL("migrations")
	if err != nil {
		t.Fatalf("build migration source URL: %v", err)
	}
	if got != "file://migrations" {
		t.Fatalf("source URL = %q, want file://migrations", got)
	}
}

func TestMigrationSourceURLBuildsAbsoluteFileURL(t *testing.T) {
	abs := filepath.Join(string(filepath.Separator), "tmp", "gatherops", "migrations")
	got, err := MigrationSourceURL(abs)
	if err != nil {
		t.Fatalf("build migration source URL: %v", err)
	}
	want := "file:///tmp/gatherops/migrations"
	if got != want {
		t.Fatalf("source URL = %q, want %q", got, want)
	}
}
