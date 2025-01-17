package flow

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/lib/pq"
	"gopkg.in/yaml.v3"
)

type DB struct {
	Host       string `yaml:"host"`
	Port       int    `yaml:"port"`
	User       string `yaml:"user"`
	Password   string `yaml:"password"`
	DBName     string `yaml:"dbname"`
	Directory  string
	Connection *sql.DB
}

// reads the database credentials from config/database.yml
func ReadDatabaseConfiguration(directory string) *DB {
	dir := filepath.Join(directory, "database.yml")
	file, err := os.Open(dir)
	if err != nil {
		checkError(err, "FATAL: error opening YAML file: - %v\n")
	}

	defer file.Close()

	var config DB

	// Parse the YAML data into the struct
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		checkError(err, "FATAL: error decoding YAML file: - %v\n")
	}
	config.Directory = directory
	return &config
}

// Connect establishes and returns a PostgreSQL DB instance
func (d *DB) Connect() {
	psqlconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		d.Host, d.Port, d.User, d.Password, d.DBName,
	)

	db, err := sql.Open("postgres", psqlconn)
	if err != nil {
		checkError(err, "CRITICAL: Failed to open database - %v\n")
	}

	// Verify connection is alive
	if err := db.Ping(); err != nil {
		db.Close()
		checkError(err, "FATAL: to ping database - %v\n")
	}

	d.Connection = db
	logInfo("Database connected successfully.")
}

func (db DB) WithTransaction(action string, queryFunction func(tx *sql.Tx) error) {
	tx, err := db.Connection.Begin()
	if err != nil {
		checkError(err, "FATAL: - %v\n")
	}

	// Defer rollback in case of a panic or early return
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			checkError(nil, "Migration failed rollback complete.")
		}
	}()

	// Execute the provided query function
	var action_term string
	if action == "up" {
		action_term = "Migration"
	} else if action == "down" {
		action_term = "Rollback"
	}

	err = queryFunction(tx)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			checkError(rollbackErr, "FATAL: error during rollback: - %v\n")
		}

		logInfo(action_term + " failed database rolled back.")
		return
	}

	if commitErr := tx.Commit(); commitErr != nil {
		checkError(commitErr, "FATAL: error during commit -  %v\n")
	}

	logInfo(action_term + " completed.")
}
