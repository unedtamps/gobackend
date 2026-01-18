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
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/unedtamps/gobackend/internal/config"
	"github.com/unedtamps/gobackend/pkg/utils"
)

func main() {
	db := flag.String("db", "", "database to migrate (e.g., postgres_primary, mysql_primary)")
	action := flag.String("action", "", "migration action (up, down, force, version, create)")
	steps := flag.Int("steps", 0, "number of migrations to apply for up and down. if 0, all migrations will be applied")
	name := flag.String("name", "", "migration name for create action")
	help := flag.Bool("help", false, "show help message")

	flag.Parse()

	if *help {
		printHelp()
		return
	}

	if *db == "" && *action != "help" {
		log.Println("Error: -db flag is required")
		printHelp()
		return
	}

	if *action == "" {
		log.Println("Error: -action flag is required")
		printHelp()
		return
	}

	cfg, err := config.NewAppConfiguration(".")
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	dbURL := make(map[string]string)
	for name, c := range cfg.Databases.Postgres {
		key := fmt.Sprintf("postgres_%s", name)
		dbURL[key] = utils.LoadPostgresConnString(c.Port, c.User, c.Host, c.Password, c.DB)
	}
	for name, c := range cfg.Databases.Mysql {
		key := fmt.Sprintf("mysql_%s", name)
		dbURL[key] = utils.LoadMySQLConnStringWithDriver(c.Port, c.User, c.Host, c.Password, c.DB)
	}
	dbParts := strings.SplitN(*db, "_", 2)
	if len(dbParts) != 2 {
		log.Fatalf("Error: invalid -db format. expected <type>_<name> (e.g., postgres_primary)")
	}
	dbType := dbParts[0]
	dbName := dbParts[1]

	sourceURL := fmt.Sprintf("file://internal/migration/%s/%s", dbType, dbName)

	if *action == "create" {
		if *name == "" {
			log.Println("Error: -name flag is required for create action")
			printHelp()
			return
		}
		migrationDir := fmt.Sprintf("internal/migration/%s/%s", dbType, dbName)
		if err := os.MkdirAll(migrationDir, os.ModePerm); err != nil {
			log.Fatalf("could not create migration directory: %v", err)
		}
		now := time.Now().UnixNano()
		versionStr := strconv.FormatInt(now, 10)
		migrationName := strings.ReplaceAll(strings.ToLower(*name), " ", "_")

		upFile := fmt.Sprintf("%s_%s.up.sql", versionStr, migrationName)
		downFile := fmt.Sprintf("%s_%s.down.sql", versionStr, migrationName)

		if err := os.WriteFile(filepath.Join(migrationDir, upFile), nil, 0644); err != nil {
			log.Fatalf("could not create up migration file: %v", err)
		}
		log.Printf("Created migration file: %s", filepath.Join(migrationDir, upFile))

		if err := os.WriteFile(filepath.Join(migrationDir, downFile), nil, 0644); err != nil {
			log.Fatalf("could not create down migration file: %v", err)
		}
		log.Printf("Created migration file: %s", filepath.Join(migrationDir, downFile))
		return
	}
	databaseURL, ok := dbURL[*db]
	if !ok {
		log.Fatalf("Error: database '%s' not found in configuration", *db)
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
	fmt.Println("\nExample:")
	fmt.Println("  go run cmd/migration/main.go -db=postgres_primary -action=up")
	fmt.Println("  go run cmd/migration/main.go -db=postgres_primary -action=create -name=my_new_migration")
}
