package flow

import (
	"database/sql"
	"fmt"
	"slices"
	"strconv"
)

func (db DB) executeDownMigrationScript(step int) {
	query := func(tx *sql.Tx) error {
		var err error
		var message string
		migrations := readMigrationFiles()
		migratedVersions := fetchMigratedVersions(tx)

		slices.Reverse(migratedVersions)

		for index, version := range migratedVersions {
			migration := migrations[version]

			_, err = tx.Exec(migration.Down)
			if err != nil {
				logError(err, "Rollback error %v for "+migration.Name+" \n")
				return err
			}

			deleteMigratedVersion(migration, tx)

			message = "Reverting " + migration.Name + " successful: version - "
			deleteMigratedVersion(migration, tx)
			logInfo(message + strconv.Itoa(migration.Version))

			if index == step-1 {
				break
			}
		}

		return nil
	}

	db.WithTransaction("down", query)
}

func deleteMigratedVersion(migration Migration, tx *sql.Tx) error {
	_, err := tx.Exec(`
		DELETE FROM schema_migrations
		WHERE version = $1;
	`, migration.Version)

	if err != nil {
		logError(err, fmt.Sprintf("FATAL: could not delete migration %s (version: %d) because %v \n", migration.Name, migration.Version, err))
		return err
	}

	return nil
}
