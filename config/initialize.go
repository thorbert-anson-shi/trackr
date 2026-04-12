package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
	"tobtoby/trackr/logging"
)

func InitializeEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		logging.GlobalLogger.Fatalf("Failed to read .env file: %s", err.Error())
	}

	for _, key_val := range os.Environ() {
		key_val_arr := strings.Split(key_val, "=")
		Config[key_val_arr[0]] = key_val_arr[1]
	}

	logging.GlobalLogger.Println("Successfully loaded .env file")
}
