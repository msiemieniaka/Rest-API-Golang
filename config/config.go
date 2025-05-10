package config

import (
	"fmt"
	"os"
	"strconv"
)

var (
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
)

func LoadConfig() {
	DBHost = os.Getenv("DB_HOST")
	DBUser = os.Getenv("DB_USER")
	DBPassword = os.Getenv("DB_PASSWORD")
	DBName = os.Getenv("DB_NAME")

	portStr := os.Getenv("DB_PORT")
	if portStr != "" {
		var err error
		DBPort, err = strconv.Atoi(portStr)
		if err != nil {
			panic(fmt.Sprintf("Invalid DB_PORT value: %s", portStr))
		}
	}

	// Sprawdzanie czy wszystkie wymagane zmienne są ustawione
	required := []string{"DB_HOST", "DB_PORT", "DB_USER", "DB_PASSWORD", "DB_NAME"}
	for _, env := range required {
		if os.Getenv(env) == "" {
			panic(fmt.Sprintf("Required environment variable %s is not set", env))
		}
	}
}
