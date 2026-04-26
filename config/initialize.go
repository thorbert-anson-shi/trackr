package config

import (
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func InitializeEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Default().Fatalf("Failed to read .env file: %s\n", err.Error())
	}

	for _, keyVal := range os.Environ() {
		keyValArr := strings.Split(keyVal, "=")
		Config[keyValArr[0]] = keyValArr[1]
	}

	log.Default().Println("Successfully loaded .env file")
}
