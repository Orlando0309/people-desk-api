package main

import (
	"database/sql"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"go-server/internal/config"

	_ "github.com/lib/pq"
)

const (
	migrationsDir   = "internal/migrations"
	migrationsTable = "schema_migrations"
)

type Migration struct {
	Version   string
	Name      string
	UpFile    string
	DownFile  string
	Applied   bool
	AppliedAt *string
}

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := connectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Ensure migrations table exists
	if err := ensureMigrationsTable(db); err != nil {
		log.Fatalf("Failed to create migrations table: %v", err)
	}

	switch command {
	case "up":
		if len(os.Args) > 2 {
			version := os.Args[2]
			if err := migrateUpTo(db, version); err != nil {
				log.Fatalf("Migration failed: %v", err)
			}
		} else {
			if err := migrateUp(db); err != nil {
				log.Fatalf("Migration failed: %v", err)
			}
		}
	case "down":
		if len(os.Args) > 2 {
			version := os.Args[2]
			if err := migrateDownTo(db, version); err != nil {
				log.Fatalf("Migration failed: %v", err)
			}
		} else {
			if err := migrateDown(db); err != nil {
				log.Fatalf("Migration failed: %v", err)
			}
		}
	case "status":
		printStatus(db)
	case "create":
		if len(os.Args) < 3 {
			log.Fatal("Please provide a migration name")
		}
		name := os.Args[2]
		if err := createMigration(name); err != nil {
			log.Fatalf("Failed to create migration: %v", err)
		}
	case "reset":
		if err := resetDatabase(db); err != nil {
			log.Fatalf("Reset failed: %v", err)
		}
	default:
		fmt.Printf("Unknown command: %s\n\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: go run cmd/migrate/main.go <command> [args]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  up [version]    Run all pending migrations or up to specific version")
	fmt.Println("  down [version]  Rollback last migration or down to specific version")
	fmt.Println("  status          Show migration status")
	fmt.Println("  create <name>   Create a new migration file")
	fmt.Println("  reset           Rollback all migrations and run them again")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  go run cmd/migrate/main.go up")
	fmt.Println("  go run cmd/migrate/main.go up 005")
	fmt.Println("  go run cmd/migrate/main.go down")
	fmt.Println("  go run cmd/migrate/main.go down 003")
	fmt.Println("  go run cmd/migrate/main.go status")
	fmt.Println("  go run cmd/migrate/main.go create add_new_column")
	fmt.Println("  go run cmd/migrate/main.go reset")
}

func connectDB(cfg config.Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBDatabase)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func ensureMigrationsTable(db *sql.DB) error {
	query := `
		CREATE TABLE IF NOT EXISTS ` + migrationsTable + ` (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMPTZ DEFAULT NOW()
		);
	`
	_, err := db.Exec(query)
	return err
}

func getMigrations() ([]Migration, error) {
	var migrations []Migration

	// Read migration files
	err := filepath.WalkDir(migrationsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		// Parse filename: 001_create_users_table.up.sql
		re := regexp.MustCompile(`^(\d+)_(.+)\.(up|down)\.sql$`)
		matches := re.FindStringSubmatch(d.Name())

		if len(matches) != 4 {
			return nil
		}

		version := matches[1]
		name := matches[2]
		direction := matches[3]

		// Find or create migration entry
		var migration *Migration
		for i := range migrations {
			if migrations[i].Version == version && migrations[i].Name == name {
				migration = &migrations[i]
				break
			}
		}

		if migration == nil {
			migrations = append(migrations, Migration{
				Version: version,
				Name:    name,
			})
			migration = &migrations[len(migrations)-1]
		}

		if direction == "up" {
			migration.UpFile = path
		} else {
			migration.DownFile = path
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort migrations by version
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Version < migrations[j].Version
	})

	return migrations, nil
}

func getAppliedMigrations(db *sql.DB) (map[string]string, error) {
	applied := make(map[string]string)

	rows, err := db.Query("SELECT version, applied_at FROM " + migrationsTable + " ORDER BY version")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var version, appliedAt string
		if err := rows.Scan(&version, &appliedAt); err != nil {
			return nil, err
		}
		applied[version] = appliedAt
	}

	return applied, nil
}

func loadMigrationsWithStatus(db *sql.DB) ([]Migration, error) {
	migrations, err := getMigrations()
	if err != nil {
		return nil, err
	}

	applied, err := getAppliedMigrations(db)
	if err != nil {
		return nil, err
	}

	for i := range migrations {
		if appliedAt, exists := applied[migrations[i].Version]; exists {
			migrations[i].Applied = true
			migrations[i].AppliedAt = &appliedAt
		}
	}

	return migrations, nil
}

func migrateUp(db *sql.DB) error {
	migrations, err := loadMigrationsWithStatus(db)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	count := 0
	for _, migration := range migrations {
		if !migration.Applied && migration.UpFile != "" {
			if err := applyMigration(tx, migration); err != nil {
				return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
			}
			fmt.Printf("Applied migration %s_%s\n", migration.Version, migration.Name)
			count++
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	if count == 0 {
		fmt.Println("No pending migrations to apply")
	} else {
		fmt.Printf("Successfully applied %d migration(s)\n", count)
	}

	return nil
}

func migrateUpTo(db *sql.DB, targetVersion string) error {
	migrations, err := loadMigrationsWithStatus(db)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	count := 0
	for _, migration := range migrations {
		if !migration.Applied && migration.UpFile != "" {
			if migration.Version > targetVersion {
				break
			}
			if err := applyMigration(tx, migration); err != nil {
				return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
			}
			fmt.Printf("Applied migration %s_%s\n", migration.Version, migration.Name)
			count++
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	if count == 0 {
		fmt.Println("No pending migrations to apply")
	} else {
		fmt.Printf("Successfully applied %d migration(s)\n", count)
	}

	return nil
}

func migrateDown(db *sql.DB) error {
	migrations, err := loadMigrationsWithStatus(db)
	if err != nil {
		return err
	}

	// Find the last applied migration
	var lastApplied *Migration
	for i := len(migrations) - 1; i >= 0; i-- {
		if migrations[i].Applied && migrations[i].DownFile != "" {
			lastApplied = &migrations[i]
			break
		}
	}

	if lastApplied == nil {
		fmt.Println("No migrations to rollback")
		return nil
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := rollbackMigration(tx, *lastApplied); err != nil {
		return fmt.Errorf("failed to rollback migration %s: %w", lastApplied.Version, err)
	}

	fmt.Printf("Rolled back migration %s_%s\n", lastApplied.Version, lastApplied.Name)

	if err := tx.Commit(); err != nil {
		return err
	}

	fmt.Println("Successfully rolled back 1 migration")
	return nil
}

func migrateDownTo(db *sql.DB, targetVersion string) error {
	migrations, err := loadMigrationsWithStatus(db)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	count := 0
	for i := len(migrations) - 1; i >= 0; i-- {
		if migrations[i].Applied && migrations[i].DownFile != "" {
			if migrations[i].Version <= targetVersion {
				break
			}
			if err := rollbackMigration(tx, migrations[i]); err != nil {
				return fmt.Errorf("failed to rollback migration %s: %w", migrations[i].Version, err)
			}
			fmt.Printf("Rolled back migration %s_%s\n", migrations[i].Version, migrations[i].Name)
			count++
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	if count == 0 {
		fmt.Println("No migrations to rollback")
	} else {
		fmt.Printf("Successfully rolled back %d migration(s)\n", count)
	}

	return nil
}

func applyMigration(tx *sql.Tx, migration Migration) error {
	// Read SQL file
	content, err := os.ReadFile(migration.UpFile)
	if err != nil {
		return err
	}

	// Execute SQL
	if _, err := tx.Exec(string(content)); err != nil {
		return err
	}

	// Record migration
	if _, err := tx.Exec("INSERT INTO "+migrationsTable+" (version) VALUES ($1)", migration.Version); err != nil {
		return err
	}

	return nil
}

func rollbackMigration(tx *sql.Tx, migration Migration) error {
	// Read SQL file
	content, err := os.ReadFile(migration.DownFile)
	if err != nil {
		return err
	}

	// Execute SQL
	if _, err := tx.Exec(string(content)); err != nil {
		return err
	}

	// Remove migration record
	if _, err := tx.Exec("DELETE FROM "+migrationsTable+" WHERE version = $1", migration.Version); err != nil {
		return err
	}

	return nil
}

func printStatus(db *sql.DB) {
	migrations, err := loadMigrationsWithStatus(db)
	if err != nil {
		log.Fatalf("Failed to load migrations: %v", err)
	}

	fmt.Println("\nMigration Status:")
	fmt.Println(strings.Repeat("-", 80))
	fmt.Printf("%-10s %-40s %-10s %s\n", "Version", "Name", "Status", "Applied At")
	fmt.Println(strings.Repeat("-", 80))

	for _, migration := range migrations {
		status := "Pending"
		appliedAt := "-"
		if migration.Applied {
			status = "Applied"
			if migration.AppliedAt != nil {
				appliedAt = *migration.AppliedAt
			}
		}
		fmt.Printf("%-10s %-40s %-10s %s\n", migration.Version, migration.Name, status, appliedAt)
	}

	fmt.Println(strings.Repeat("-", 80))

	appliedCount := 0
	pendingCount := 0
	for _, migration := range migrations {
		if migration.Applied {
			appliedCount++
		} else {
			pendingCount++
		}
	}

	fmt.Printf("\nTotal: %d migrations (%d applied, %d pending)\n", len(migrations), appliedCount, pendingCount)
}

func createMigration(name string) error {
	// Get the next version number
	migrations, err := getMigrations()
	if err != nil {
		return err
	}

	nextVersion := "001"
	if len(migrations) > 0 {
		lastVersion := migrations[len(migrations)-1].Version
		num := 0
		fmt.Sscanf(lastVersion, "%d", &num)
		nextVersion = fmt.Sprintf("%03d", num+1)
	}

	// Sanitize name
	name = strings.ToLower(strings.ReplaceAll(name, " ", "_"))

	// Create up file
	upFileName := fmt.Sprintf("%s/%s_%s.up.sql", migrationsDir, nextVersion, name)
	upContent := fmt.Sprintf("-- Migration: %s\n-- Description: %s\n\n", nextVersion, name)
	if err := os.WriteFile(upFileName, []byte(upContent), 0644); err != nil {
		return err
	}

	// Create down file
	downFileName := fmt.Sprintf("%s/%s_%s.down.sql", migrationsDir, nextVersion, name)
	downContent := fmt.Sprintf("-- Rollback: %s\n-- Description: %s\n\n", nextVersion, name)
	if err := os.WriteFile(downFileName, []byte(downContent), 0644); err != nil {
		return err
	}

	fmt.Printf("Created migration files:\n")
	fmt.Printf("  %s\n", upFileName)
	fmt.Printf("  %s\n", downFileName)

	return nil
}

func resetDatabase(db *sql.DB) error {
	fmt.Println("This will rollback all migrations and reapply them.")
	fmt.Print("Are you sure? (yes/no): ")

	var response string
	fmt.Scanln(&response)

	if strings.ToLower(response) != "yes" {
		fmt.Println("Reset cancelled")
		return nil
	}

	// Rollback all migrations
	migrations, err := loadMigrationsWithStatus(db)
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Rollback in reverse order
	rolledBack := 0
	for i := len(migrations) - 1; i >= 0; i-- {
		if migrations[i].Applied && migrations[i].DownFile != "" {
			if err := rollbackMigration(tx, migrations[i]); err != nil {
				return fmt.Errorf("failed to rollback migration %s: %w", migrations[i].Version, err)
			}
			fmt.Printf("Rolled back migration %s_%s\n", migrations[i].Version, migrations[i].Name)
			rolledBack++
		}
	}

	// Apply all migrations
	for _, migration := range migrations {
		if migration.UpFile != "" {
			if err := applyMigration(tx, migration); err != nil {
				return fmt.Errorf("failed to apply migration %s: %w", migration.Version, err)
			}
			fmt.Printf("Applied migration %s_%s\n", migration.Version, migration.Name)
		}
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	fmt.Printf("\nSuccessfully reset database (%d migrations rolled back and reapplied)\n", rolledBack)
	return nil
}
