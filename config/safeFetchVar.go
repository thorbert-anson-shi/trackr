package config

import "log"

func SafeFetchVar(key string) string {
	value, exists := Config[key]
	if !exists {
		log.Default().Fatalf("Please provide a valid value for %s\n", key)
	}

	return value
}
