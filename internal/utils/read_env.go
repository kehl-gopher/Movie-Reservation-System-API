package utils

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// read environment variable helper function
func ReadEnvVariable(key string) (string, error) {
	err := godotenv.Load(".env")
	if err != nil {
		return "", err
	}
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("environment key %s does not exist", key)
	}
	return val, nil
}
