package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/unedtamps/gobackend/internal/config"
	"github.com/unedtamps/gobackend/pkg/utils"
)

func main() {
	db := flag.String("db", "", "database to migrate (e.g., primary, secondary, backup)")
	action := flag.String("action", "", "migration action (up, down, force, version, create)")
	steps := flag.Int(
		"steps",
		0,
		"number of migrations to apply for up and down. if 0, all migrations will be applied",
	)
	name := flag.String("name", "", "migration name for create action")
	list := flag.Bool("list", false, "list available databases")
	help := flag.Bool("help", false, "show help message")

	flag.Parse()

	if *help {
		printHelp()
		return
	}

	cfg, err := config.NewAppConfiguration(".")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	if *list {
		fmt.Println("Available databases:")
		for _, dbName := range cfg.Databases.ListNames() {
			if dbInstance, ok := cfg.Databases.GetByName(dbName); ok {
				fmt.Printf("  - %s (%s)\n", dbName, dbInstance.GetRDBMS())
			}
		}
		return
	}

	dbName := *db
	if dbName == "" {
		log.Println("Error: -db flag is required")
		printHelp()
		return
	}

	// Validate action
	if *action == "" {
		log.Println("Error: -action flag is required")
		printHelp()
		return
	}
	dbInstance, ok := cfg.Databases.GetByName(dbName)
	if !ok {
		log.Fatalf(
			"Error: database '%s' not found in configuration. Use -list to see available databases.",
			dbName,
		)
	}

	var databaseURL string
	switch dbInstance.GetRDBMS() {
	case "postgres":
		databaseURL = utils.LoadPostgresConnString(
			dbInstance.GetPort(),
			dbInstance.GetUser(),
			dbInstance.GetHost(),
			dbInstance.GetPassword(),
			dbInstance.GetDBName(),
		)
	case "mysql":
		databaseURL = utils.LoadMySQLConnStringWithDriver(
			dbInstance.GetPort(),
			dbInstance.GetUser(),
			dbInstance.GetHost(),
			dbInstance.GetPassword(),
			dbInstance.GetDBName(),
		)
	case "sqlite":
		databaseURL = fmt.Sprintf("sqlite3://%s", dbInstance.GetDBName())
	default:
		log.Fatalf("Error: unsupported RDBMS '%s'", dbInstance.GetDBName())
	}

	// Build source URL for migrations
	sourceURL := fmt.Sprintf(
		"file://internal/datastore/%s/migration",
		dbInstance.GetName(),
	)

	// Handle create action separately (doesn't need database connection)
	if *action == "create" {
		if *name == "" {
			log.Println("Error: -name flag is required for create action")
			printHelp()
			return
		}
		migrationDir := fmt.Sprintf("internal/datastore/%s/migration", dbInstance.GetName())
		if err := os.MkdirAll(migrationDir, os.ModePerm); err != nil {
			log.Fatalf("could not create migration directory: %v", err)
		}
		now := time.Now().UnixNano()
		versionStr := strconv.FormatInt(now, 10)
		migrationName := strings.ReplaceAll(strings.ToLower(*name), " ", "_")

		upFile := fmt.Sprintf("%s_%s.up.sql", versionStr, migrationName)
		downFile := fmt.Sprintf("%s_%s.down.sql", versionStr, migrationName)

		if err := os.WriteFile(filepath.Join(migrationDir, upFile), nil, 0o644); err != nil {
			log.Fatalf("could not create up migration file: %v", err)
		}
		log.Printf("Created migration file: %s", filepath.Join(migrationDir, upFile))

		if err := os.WriteFile(filepath.Join(migrationDir, downFile), nil, 0o644); err != nil {
			log.Fatalf("could not create down migration file: %v", err)
		}
		log.Printf("Created migration file: %s", filepath.Join(migrationDir, downFile))
		return
	}

	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		log.Fatalf("Error creating migrate instance: %v", err)
	}

	switch *action {
	case "up":
		if *steps > 0 {
			err = m.Steps(*steps)
		} else {
			err = m.Up()
		}
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Error applying migrations up: %v", err)
		}
		log.Println("Migrations applied successfully")
	case "down":
		if *steps > 0 {
			err = m.Steps(-*steps)
		} else {
			err = m.Down()
		}
		if err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Error rolling back migrations: %v", err)
		}
		log.Println("Migrations rolled back successfully")
	case "force":
		err = m.Force(*steps)
		if err != nil {
			log.Fatalf("Error forcing migration version: %v", err)
		}
		log.Printf("Forced migration version to %d", *steps)
	case "version":
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("Error getting migration version: %v", err)
		}
		log.Printf("Migration version: %d, dirty: %v", version, dirty)
	default:
		log.Printf("Error: unknown action '%s'", *action)
		printHelp()
	}
}

func printHelp() {
	fmt.Println("Usage: go run cmd/migration/main.go [flags]")
	fmt.Println("\nFlags:")
	flag.PrintDefaults()
	fmt.Println("\nExamples:")
	fmt.Println("  # List available databases")
	fmt.Println("  go run cmd/migration/main.go -list")
	fmt.Println("")
	fmt.Println("  # Interactive mode (asks for database)")
	fmt.Println("  go run cmd/migration/main.go -action=up")
	fmt.Println("")
	fmt.Println("  # Direct database selection")
	fmt.Println("  go run cmd/migration/main.go -db=primary -action=up")
	fmt.Println("  go run cmd/migration/main.go -db=secondary -action=up -steps=1")
	fmt.Println("  go run cmd/migration/main.go -db=backup -action=version")
	fmt.Println("")
	fmt.Println("  # Create new migration")
	fmt.Println("  go run cmd/migration/main.go -db=primary -action=create -name=my_new_migration")
}
