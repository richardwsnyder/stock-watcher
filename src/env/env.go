package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GoDotEnvVariable(key string) string {
	// load .env file
	err := godotenv.Load("../../.env")

	if err != nil {
		// If running heroku commands from Procfile, might require
		// root directory
		err = godotenv.Load(".env")
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
	}

	return os.Getenv(key)
}
