package flow

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// Reading migration files
func readMigrationFiles(directory string) map[int]Migration {
	dir := filepath.Join(directory, "migrations")

	files, err := os.ReadDir(dir)
	if err != nil {
		checkError(err, "CRITICAL: %v\n")
	}

	migrations := make(map[int]Migration)

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".yml" {
			// Read the content of each .yml file
			filePath := filepath.Join(dir, file.Name())
			migration, err := migrationFromFile(filePath)
			if err != nil {
				checkError(err, "CRITICAL: %v\n")
			}

			migrations[migration.Version] = migration
		}
	}

	return migrations
}

func migrationFromFile(filePath string) (Migration, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return Migration{}, fmt.Errorf("error opening migration file %s: %w", filePath, err)
	}
	defer file.Close()

	var migration Migration
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&migration)
	if err != nil {
		return Migration{}, fmt.Errorf("error decoding migration file %s: %w", filePath, err)
	}

	return migration, nil
}

// Creating Migration files
func createMigrationFile(migration_name string, directory string) {
	currentTime := time.Now()
	version := strings.ReplaceAll(currentTime.Format("20060102150405.000"), ".", "")

	dir := filepath.Join(directory, "migrations")
	migration_name = strings.Join(strings.Fields(migration_name), "_")
	fileName := version + "_" + strings.ToLower(migration_name) + ".yml"

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		checkError(err, "FATAL: %v\n")
	}

	createFileInDir(dir, fileName, migration_name, version)
}

func createFileInDir(dir, fileName, migration_name, version string) {
	// Create the full path for the new file
	filePath := filepath.Join(dir, fileName)

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		checkError(err, "FATAL: error creating file - %v\n")
	}
	defer file.Close()

	migration_template := migrationFileTemplate(migration_name, version)

	// Write template to the file
	_, err = file.WriteString(migration_template)
	if err != nil {
		checkError(err, "FATAL: error writing file - %v\n")
	}

	logInfo("Migration File created " + filePath)
}

func migrationFileTemplate(migration_name, version string) string {
	migration_name = strings.Join(strings.Split(migration_name, "_"), " ")

	content := `version: %s
name: %s
description: #migration description goes here
up: | # up script goes under here
down: | # down sql script goes under here`

	return fmt.Sprintf(content, version, migration_name)
}
