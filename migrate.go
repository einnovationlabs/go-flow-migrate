package flow

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func Start(directory string) {
	logInfo(`Welcome to flow migrate.
1. Create a migration file
2. Run Migrations
3. Rollback Migrations`)

	choice := readInput("Select an option to proceed: ")
	switch choice {
	case "1":
		file_name := readInput("Enter Migration file name: ")
		Create(file_name, directory)
	case "2":
		MigrateUp(directory)
	case "3":
		step := readInput("Enter step count(int): ")
		step_int, err := strconv.Atoi(step)
		if err != nil {
			log.Fatalf("ERROR: Invalid input. try again")
		}
		MigrateDown(step_int, directory)
	default:
		logInfo("Invalid choice. Please select a valid option (1, 2, or 3).")
	}
}

func readInput(display string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(display)
	input, err := reader.ReadString('\n')
	if err != nil {
		log.Fatalf("ERROR: Failed to read input: %v", err)
	}

	input = input[:len(input)-1]

	return input
}

func MigrateUp(directory string) {
	db := ReadDatabaseConfiguration(directory)
	db.Connect()
	db.RunMigrations("up", 0)
}

func MigrateDown(step int, directory string) {
	db := ReadDatabaseConfiguration(directory)
	db.Connect()
	db.RunMigrations("down", step)
}

func Create(migration_name string, directory string) {
	createMigrationFile(migration_name, directory)
}
