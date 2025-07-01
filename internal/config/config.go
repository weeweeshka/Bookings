package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseName     string `env:"DB_NAME"`
	DatabaseHost     string `env:"DB_HOST"`
	DatabasePort     string `env:"DB_PORT"`
	DatabaseUser     string `env:"DB_USER"`
	DatabasePassword string `env:"DB_PASSWORD"`
	DatabaseUrl      string `env:"DB_URL" env-default:"postgres://postgres:postgres@localhost:5432/bookings?sslmode=disable"`
}

func MustLoad() *Config {
	if err := godotenv.Load("D:/Bookings-1/.env"); err != nil {
		log.Printf("Warning: No .env file found: %v", err)
	}

	cfg := Config{
		DatabaseHost: "localhost",
		DatabasePort: "5432",
	}

	if dbName := os.Getenv("DB_NAME"); dbName != "" {
		cfg.DatabaseName = dbName
	}

	if dbHost := os.Getenv("DB_HOST"); dbHost != "" {
		cfg.DatabaseHost = dbHost
	}

	if dbPort := os.Getenv("DB_PORT"); dbPort != "" {
		cfg.DatabasePort = dbPort
	}

	if dbUser := os.Getenv("DB_USER"); dbUser != "" {
		cfg.DatabaseUser = dbUser
	}

	if dbPass := os.Getenv("DB_PASSWORD"); dbPass != "" {
		cfg.DatabasePassword = dbPass
	}

	if dbUrl := os.Getenv("DB_URL"); dbUrl != "" {
		cfg.DatabaseUrl = dbUrl
	}

	if cfg.DatabaseName == "" {
		log.Fatal("DB_NAME is required")
	}
	if cfg.DatabaseUser == "" {
		log.Fatal("DB_USER is required")
	}

	return &cfg
}
