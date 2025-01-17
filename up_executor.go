package flow

import (
	"database/sql"
	"sort"
	"strconv"
)

func (db DB) executeUpMigrationScript() {
	query := func(tx *sql.Tx) error {
		var err error
		var message string
		migrations := readMigrationFiles(db.Directory)
		migratedVersions := fetchMigratedVersions(tx)

		for _, version := range migratedVersions {
			delete(migrations, version)
		}

		// sort migration versions
		var migrationVersions []int
		for version := range migrations {
			migrationVersions = append(migrationVersions, version)
		}

		sort.Ints(migrationVersions)

		for _, version := range migrationVersions {
			migration := migrations[version]
			_, err = tx.Exec(migration.Up)
			if err != nil {
				logError(err, "Migration error %v for "+migration.Name+" \n")
				return err
			}

			insertMigratedVersion(migration, tx)

			message = "Migrating " + migration.Name + " successful: version - "
			logInfo(message + strconv.Itoa(migration.Version))
		}

		return nil
	}

	db.WithTransaction("up", query)
}

func insertMigratedVersion(migration Migration, tx *sql.Tx) error {
	_, err := tx.Exec(`
		INSERT INTO schema_migrations (version, name)
		VALUES ($1, $2);
	`, migration.Version, migration.Name)

	if err != nil {
		logError(err, "FATAL: could not insert migration "+migration.Name+" because %v \n")
		return err
	}

	return nil
}
