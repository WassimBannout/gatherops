package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"

	"github.com/WassimBannout/gatherops/internal/config"
	"github.com/WassimBannout/gatherops/internal/database"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	direction := flag.String("direction", string(database.MigrationUp), "migration direction: up or down")
	migrationsPath := flag.String("path", "migrations", "path to migration files")
	steps := flag.Int("steps", 0, "number of migration steps; 0 means all available for up/down")
	flag.Parse()

	if *steps < 0 {
		logger.Error("steps must be zero or positive")
		os.Exit(1)
	}

	cfg, err := config.Load()
	if err != nil {
		logger.Error("load config", "error", err)
		os.Exit(1)
	}

	migrationDirection := database.MigrationDirection(*direction)
	if err := database.RunMigrations(cfg.DatabaseURL, *migrationsPath, migrationDirection, *steps); err != nil {
		logger.Error("run migrations", "error", err)
		os.Exit(1)
	}

	fmt.Printf("migrations %s complete\n", migrationDirection)
}
