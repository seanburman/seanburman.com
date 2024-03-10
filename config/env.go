package config

import (
	"flag"
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Type Vars
//
// Environment variables
type Vars struct {
	PORT         string
	POSTGRES_URI string
	MONGO_URI    string
	JWT_SECRET   string
	VERSION      string
}

// Env() returns Vars struct of environment variables
func Env() Vars {
	// Load if not a test. This isn't required during testing.
	if flag.Lookup("test.v") == nil {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading environment variables")
		}
	}

	return Vars{
		PORT:         os.Getenv("PORT"),
		POSTGRES_URI: os.Getenv("POSTGRES_URI"),
		MONGO_URI:    os.Getenv("MONGO_URI"),
		JWT_SECRET:   os.Getenv("JWT_SECRET"),
		VERSION:      os.Getenv("VERSION"),
	}
}
