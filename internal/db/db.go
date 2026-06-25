package db

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

	_ "github.com/lib/pq"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

// Connect establishes a connection to the PostgreSQL database with retries
func Connect(host, port, user, password, dbname string) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var db *sql.DB
	var err error

	// Retry connection loop
	maxRetries := 10
	for i := 1; i <= maxRetries; i++ {
		log.Printf("Connecting to database (attempt %d/%d)...", i, maxRetries)
		db, err = sql.Open("postgres", connStr)
		if err == nil {
			err = db.Ping()
			if err == nil {
				log.Println("Successfully connected to the database!")
				return db, nil
			}
		}

		if db != nil {
			db.Close()
		}

		log.Printf("Database connection failed: %v. Retrying in 3 seconds...", err)
		time.Sleep(3 * time.Second)
	}

	return nil, fmt.Errorf("could not connect to database after %d attempts: %w", maxRetries, err)
}

// RunMigrations applies all pending migrations from the embedded file system
func RunMigrations(db *sql.DB) error {
	// Create schema_migrations table if not exists
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version VARCHAR(255) PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migration table: %w", err)
	}

	// Read migration files from embedded FS
	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	type Migration struct {
		Name string
		SQL  string
	}

	var migrations []Migration
	for _, entry := range entries {
		name := entry.Name()
		if !entry.IsDir() && strings.HasSuffix(name, ".up.sql") {
			content, err := migrationFiles.ReadFile("migrations/" + name)
			if err != nil {
				return fmt.Errorf("failed to read migration file %s: %w", name, err)
			}
			migrations = append(migrations, Migration{
				Name: name,
				SQL:  string(content),
			})
		}
	}

	// Sort migrations by name (e.g. 0001_..., 0002_...)
	sort.Slice(migrations, func(i, j int) bool {
		return migrations[i].Name < migrations[j].Name
	})

	// Get already applied migrations
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return fmt.Errorf("failed to fetch schema_migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[string]bool)
	for rows.Next() {
		var version string
		if err := rows.Scan(&version); err != nil {
			return err
		}
		applied[version] = true
	}

	// Apply pending migrations
	for _, m := range migrations {
		version := strings.TrimSuffix(m.Name, ".up.sql")
		if applied[version] {
			continue
		}

		log.Printf("Applying database migration: %s", m.Name)
		tx, err := db.BeginTx(context.Background(), nil)
		if err != nil {
			return fmt.Errorf("failed to begin transaction for %s: %w", m.Name, err)
		}

		if _, err := tx.Exec(m.SQL); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to execute migration %s: %w", m.Name, err)
		}

		if _, err := tx.Exec("INSERT INTO schema_migrations (version) VALUES ($1)", version); err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to record migration %s: %w", m.Name, err)
		}

		if err := tx.Commit(); err != nil {
			return fmt.Errorf("failed to commit migration %s: %w", m.Name, err)
		}
		log.Printf("Successfully applied database migration: %s", m.Name)
	}

	return nil
}
