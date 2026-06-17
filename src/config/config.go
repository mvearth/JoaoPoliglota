// Package config centralizes reading runtime configuration from the environment.
package config

import "os"

// Getenv returns the value of the environment variable named by key, or
// fallback when the variable is unset or empty.
func Getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

// DB holds the connection settings for the PostgreSQL database.
type DB struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

// LoadDB reads the database configuration from the environment, applying
// sensible defaults for local development.
func LoadDB() DB {
	return DB{
		Host:     Getenv("DB_HOST", "localhost"),
		Port:     Getenv("DB_PORT", "5432"),
		User:     Getenv("DB_USER", "postgres"),
		Password: Getenv("DB_PASSWORD", ""),
		Name:     Getenv("DB_NAME", "joaopoliglota"),
		SSLMode:  Getenv("DB_SSLMODE", "disable"),
	}
}

// HTTPAddr returns the address the HTTP server should listen on.
func HTTPAddr() string {
	return Getenv("HTTP_ADDR", ":8000")
}
