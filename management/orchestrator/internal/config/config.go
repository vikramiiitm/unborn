package config

import (
	"os"
	"strconv"
)

type Config struct {
	HTTPPort     int
	DatabaseURL  string
	RedisURL     string
	Environment  string
	MaxInstances int
}

func Load() *Config {
	port := 8080
	if p := os.Getenv("HTTP_PORT"); p != "" {
		if parsed, err := strconv.Atoi(p); err == nil {
			port = parsed
		}
	}

	maxInstances := 10
	if m := os.Getenv("MAX_INSTANCES"); m != "" {
		if parsed, err := strconv.Atoi(m); err == nil {
			maxInstances = parsed
		}
	}

	return &Config{
		HTTPPort:     port,
		DatabaseURL:  getEnv("DATABASE_URL", "postgres://unborn:unborn@localhost:5432/unborn?sslmode=disable"),
		RedisURL:     getEnv("REDIS_URL", "redis://localhost:6379"),
		Environment:  getEnv("UNBORN_ENV", "development"),
		MaxInstances: maxInstances,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
