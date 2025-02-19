package flow

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type DB struct {
	Host       string
	Port       string
	User       string
	Password   string
	DBName     string
	Directory  string
	Connection *sql.DB
}

// reads the database credentials from config/database.yml
func ReadDatabaseConfiguration(directory string) *DB {
	godotenv.Load()

	config := DB{
		Host:      getEnv("DB_HOST"),
		Port:      getEnv("DB_PORT"),
		User:      getEnv("DB_USER"),
		Password:  getEnv("DB_PASSWORD"),
		DBName:    getEnv("DB_NAME"),
		Directory: directory,
	}

	return &config
}

func getEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	err := errors.New(key)
	checkError(err, "Fatal: Failed to fetch credentials - %v\n")
	return ""
}

// Connect establishes and returns a PostgreSQL DB instance
func (d *DB) Connect() {
	portInt, err := strconv.Atoi(d.Port)
	if err != nil {
		log.Fatalf("Invalid port: %v", err)
	}

	psqlconn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		d.Host, portInt, d.User, d.Password, d.DBName,
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
