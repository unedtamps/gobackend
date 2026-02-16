package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"

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
	force := flag.Bool("force", false, "force action without confirmation (for down)")

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

		// Get next version number (auto-increment starting from 1)
		nextVersion := getNextMigrationVersion(migrationDir)
		versionStr := strconv.Itoa(nextVersion)
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
		// Handle down with confirmation
		if *steps == 0 {
			// Down all migrations
			if !*force {
				fmt.Println("WARNING: You are about to roll back ALL migrations.")
				fmt.Println("This will delete all data in the database!")
				fmt.Println("Tip: Use -steps=N to roll back specific number of migrations.")
				fmt.Print("Are you sure you want to continue? [y/N]: ")
				if !confirmAction() {
					log.Println("Operation cancelled.")
					return
				}
			}
			err = m.Down()
		} else {
			// Down specific number of steps
			if !*force {
				fmt.Printf("You are about to roll back %d migration(s).\n", *steps)
				fmt.Print("Are you sure you want to continue? [y/N]: ")
				if !confirmAction() {
					log.Println("Operation cancelled.")
					return
				}
			}
			err = m.Steps(-*steps)
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

// getNextMigrationVersion scans the migration directory and returns the next version number
func getNextMigrationVersion(migrationDir string) int {
	entries, err := os.ReadDir(migrationDir)
	if err != nil {
		// Directory doesn't exist or is empty, start from 1
		return 1
	}

	versionRegex := regexp.MustCompile(`^(\d+)_.*\.(up|down)\.sql$`)
	maxVersion := 0

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		matches := versionRegex.FindStringSubmatch(entry.Name())
		if len(matches) > 1 {
			version, err := strconv.Atoi(matches[1])
			if err == nil && version > maxVersion {
				maxVersion = version
			}
		}
	}

	return maxVersion + 1
}

// confirmAction reads user input and returns true if confirmed
func confirmAction() bool {
	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

// parseMigrationFiles parses migration files and returns sorted list
func parseMigrationFiles(migrationDir string) ([]migrationFile, error) {
	entries, err := os.ReadDir(migrationDir)
	if err != nil {
		return nil, err
	}

	versionRegex := regexp.MustCompile(`^(\d+)_(.+)\.(up|down)\.sql$`)
	migrations := make(map[int]*migrationFile)

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		matches := versionRegex.FindStringSubmatch(entry.Name())
		if len(matches) != 4 {
			continue
		}

		version, err := strconv.Atoi(matches[1])
		if err != nil {
			continue
		}

		name := matches[2]
		direction := matches[3]

		if migrations[version] == nil {
			migrations[version] = &migrationFile{
				Version: version,
				Name:    name,
			}
		}

		if direction == "up" {
			migrations[version].HasUp = true
		} else {
			migrations[version].HasDown = true
		}
	}

	result := make([]migrationFile, 0, len(migrations))
	for _, m := range migrations {
		result = append(result, *m)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].Version < result[j].Version
	})

	return result, nil
}

type migrationFile struct {
	Version int
	Name    string
	HasUp   bool
	HasDown bool
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
	fmt.Println("  # Create new migration (auto-increment version)")
	fmt.Println("  go run cmd/migration/main.go -db=primary -action=create -name=my_new_migration")
	fmt.Println("")
	fmt.Println("  # Down migrations with confirmation")
	fmt.Println("  go run cmd/migration/main.go -db=primary -action=down -steps=1")
	fmt.Println("  go run cmd/migration/main.go -db=primary -action=down -force  # Skip confirmation")
}
