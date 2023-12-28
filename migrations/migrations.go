package main

import (
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gorm.io/gorm"
)

// Migrate applies the migration files to the database in the correct order.
func Migrate(db *gorm.DB) error {
	migrationsDir := "./migrations" // Adjust the path to your migrations directory

	// Store all migration files in a slice
	var migrationFiles []string

	// Collect all .up.sql files
	subDirs, err := os.ReadDir(migrationsDir)
	if err != nil {
		return err
	}

	for _, dir := range subDirs {
		if dir.IsDir() {
			files, err := os.ReadDir(filepath.Join(migrationsDir, dir.Name()))
			if err != nil {
				return err
			}
			for _, file := range files {
				if strings.HasSuffix(file.Name(), ".up.sql") {
					migrationFiles = append(migrationFiles, filepath.Join(migrationsDir, dir.Name(), file.Name()))
				}
			}
		}
	}

	// Sort the migration files by their names (which include the sequence number)
	sort.Strings(migrationFiles)

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
