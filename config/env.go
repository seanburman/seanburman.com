package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Type Vars
//
// Environment variables
type Vars struct {
	PORT         string
	DATABASE_URL string
	ENVIRONMENT  string
}

// Env() returns Vars struct of environment variables
func Env() Vars {
	env := os.Getenv("ENVIRONMENT")
	if env == "development" {
		godotenv.Load(".env")
	}
	return Vars{
		PORT:         os.Getenv("PORT"),
		DATABASE_URL: os.Getenv("DATABASE_URL"),
		ENVIRONMENT:  env,
		// POSTGRES_URI: os.Getenv("POSTGRES_URI"),
		// MONGO_URI:    os.Getenv("MONGO_URI"),
		// JWT_SECRET:   os.Getenv("JWT_SECRET"),
		// VERSION:      os.Getenv("VERSION"),
	}
}
