package config

import (
	"log"
	"os"
	"strings"
)

func InitializeEnv() {
	for _, keyVal := range os.Environ() {
		keyValArr := strings.Split(keyVal, "=")
		Config[keyValArr[0]] = keyValArr[1]
	}

	log.Default().Println("Successfully loaded environment variables")
}
