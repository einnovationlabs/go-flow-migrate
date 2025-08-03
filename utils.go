package flow

import (
	"log"
	"os"
)

func checkError(err error, message_wrapper string) {
	log.Fatalf(message_wrapper, err)
}

func logError(err error, messageWrapper string) {
	log.Printf("ERROR: "+messageWrapper, err)
}

func logInfo(message string) {
	log.Printf("%s", "INFO: "+message)
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}