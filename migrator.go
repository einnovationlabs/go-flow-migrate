package flow

import (
	"database/sql"
)

type Migration struct {
	Version     int    `yaml:"version"`
	Name        string `yaml:"name"`
	Description string `yaml:"description"`
	Up          string `yaml:"up"`
	Down        string `yaml:"down"`
	ID          int
}

func (db DB) RunMigrations(action string, step int) {
	if action == "up" {
		db.executeUpMigrationScript()
	} else if action == "down" {
		db.executeDownMigrationScript(step)
	}
}

func fetchMigratedVersions(tx *sql.Tx) []int {
	_, err := tx.Exec(`
			CREATE TABLE IF NOT EXISTS schema_migrations (
			id SERIAL PRIMARY KEY,
			version BIGINT NOT NULL,
			name VARCHAR(255) NOT NULL
		);
	`)

	if err != nil {
		checkError(err, "FATAL: Error creating table: %v\n")
	}

	rows, err := tx.Query("SELECT version FROM schema_migrations ORDER BY version ASC;")
	if err != nil {
		checkError(err, "failed to execute query: %v\n")
	}
	defer rows.Close()

	var versions []int

	for rows.Next() {
		var version int

		if err := rows.Scan(&version); err != nil {
			checkError(err, "failed to scan rows: %v\n")
		}

		versions = append(versions, version)
	}

	return versions
}
