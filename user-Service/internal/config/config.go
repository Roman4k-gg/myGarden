package config

import (
	"errors"
	"os"
)

type Config struct {
	DatabaseURL string
	JWTSecret   string
	GRPCPort    string
}

func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		GRPCPort:    getEnv("GRPC_PORT", ":50051"),
	}
	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL env var is required")
	}
	if cfg.JWTSecret == "" {
		return nil, errors.New("JWT_SECRET env var is required")
	}
	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
