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
	migrationsDir := "./pkg/migrations"

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
	err = db.Exec("CREATE TABLE IF NOT EXISTS migrations (name VARCHAR PRIMARY KEY)").Error
	if err != nil {
		return err
	}

	// Execute each migration file in sorted order
	for _, file := range migrationFiles {
		fileName := filepath.Base(file)
		log.Printf("Checking migration file: %s", fileName)
		var count int64
		db.Table("migrations").Where("name = ?", fileName).Count(&count)
		if count > 0 {
			log.Printf("Migration %s has already been applied, skipping", fileName)
			continue
		}
		content, err := os.ReadFile(file)
		if err != nil {
			return err
		}
		tx := db.Begin()
		if err := tx.Exec(string(content)).Error; err != nil {
			tx.Rollback()
			return err
		}
		if err := tx.Exec("INSERT INTO migrations (name) VALUES (?)", fileName).Error; err != nil {
			tx.Rollback()
			return err
		}
		tx.Commit()
	}

	return nil
}
