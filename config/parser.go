package config

import (
	"os"
)

// GetEnvWithFallback will retrieve a env var with a fallback string.
func GetEnvWithFallback(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}
