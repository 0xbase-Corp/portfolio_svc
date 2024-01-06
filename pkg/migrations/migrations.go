package migrations

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// Migrate applies the migration files to the database in the correct order.
func Migrate(db *gorm.DB) error {
	migrationsDir := "/app/migrations" // Adjust the path to your migrations directory

	// Store all migration files in a slice
	var migrationFiles []string

	// Recursive function to collect all .up.sql files
	var collectMigrationFiles func(string) error
	collectMigrationFiles = func(path string) error {
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}

		for _, entry := range entries {
			fullPath := filepath.Join(path, entry.Name())
			if entry.IsDir() {
				err := collectMigrationFiles(fullPath)
				if err != nil {
					return err
				}
			} else if strings.HasSuffix(entry.Name(), ".up.sql") {
				migrationFiles = append(migrationFiles, fullPath)
			}
		}
		return nil
	}

	// Start the recursive file collection
	err := collectMigrationFiles(migrationsDir)
	if err != nil {
		return err
	}

	// Sort the migration files by their sequence number
	sort.Slice(migrationFiles, func(i, j int) bool {
		seqNumI, _ := strconv.Atoi(strings.Split(filepath.Base(migrationFiles[i]), "_")[0])
		seqNumJ, _ := strconv.Atoi(strings.Split(filepath.Base(migrationFiles[j]), "_")[0])
		return seqNumI < seqNumJ
	})

	// Create a table to track which migrations have been run
	db.Exec("CREATE TABLE IF NOT EXISTS migrations (name VARCHAR PRIMARY KEY)")

	// Execute each migration file in sorted order
	for _, file := range migrationFiles {
		log.Printf("Checking migration file: %s", file)
		var count int64
		db.Table("migrations").Where("name = ?", file).Count(&count)
		if count > 0 {
			log.Printf("Migration %s has already been applied, skipping", file)
			continue
		}
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		if err := db.Exec(string(content)).Error; err != nil {
			return err
		}
		db.Exec("INSERT INTO migrations (name) VALUES (?)", file)
	}

	return nil
}
