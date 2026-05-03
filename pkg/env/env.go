package env

import (
	"log"

	"github.com/joho/godotenv"
)

func LoadEnvWithExit(filenames ...string) {
	err := godotenv.Load(filenames...)

	if err != nil {
		log.Fatalf("Error loading environmental file: %v", err)
	}
}
