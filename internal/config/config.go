package config

import (
	"os"

	"github.com/joho/godotenv"
)

// on cloud or container
func LoadEnv(path string) {
	_ = godotenv.Load(path)
}

// fall back
// local
func GetEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
