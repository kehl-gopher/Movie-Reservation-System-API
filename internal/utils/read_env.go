package utils

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
)

var KeyDoesNotExists = errors.New("Environment key does not exists")

// read environment variable helper function
func ReadEnvVariable(key string) (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return "", err
	}
	val := os.Getenv(key)
	if val == "" {
		return "", KeyDoesNotExists
	}
	return val, nil
}
