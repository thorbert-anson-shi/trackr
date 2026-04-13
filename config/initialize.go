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

	for _, key_val := range os.Environ() {
		key_val_arr := strings.Split(key_val, "=")
		Config[key_val_arr[0]] = key_val_arr[1]
	}

	log.Default().Println("Successfully loaded .env file")
}
