package config

import (
	"errors"
	"os"
)

type Config struct {
	UserServiceAddr    string
	CatalogServiceAddr string
	JWTSecret          string
	HTTPPort           string
}

func Load() (*Config, error) {
	cfg := &Config{
		UserServiceAddr:    getEnv("USER_SERVICE_ADDR", "localhost:50051"),
		CatalogServiceAddr: getEnv("CATALOG_SERVICE_ADDR", "localhost:50052"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		HTTPPort:           getEnv("HTTP_PORT", ":3000"),
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
