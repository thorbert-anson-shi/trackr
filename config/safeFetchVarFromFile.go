package config

import (
	"log"
	"os"
	"strings"
)

func SafeFetchVarFromFile(key string) string {
	fileName, exists := Config[key]
	if !exists {
		log.Default().Fatalf("Please provide a valid value for %s\n", key)
	}

	valueBytes, err := os.ReadFile(fileName)
	if err != nil {
		log.Default().Fatalf("Failed to read file defined in %s: %s\n", key, err.Error())
	}

	value := strings.TrimSpace(string(valueBytes))

	return value
}
