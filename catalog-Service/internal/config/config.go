package config

import (
	"errors"
	"os"
)

type Config struct {
	DatabaseURL string
	RedisAddr   string
	KafkaBroker string
	GRPCPort    string
}

func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		RedisAddr:   getEnv("REDIS_ADDR", "localhost:6379"),
		KafkaBroker: getEnv("KAFKA_BROKER", "localhost:9092"),
		GRPCPort:    getEnv("GRPC_PORT", ":50052"),
	}
	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL env var is required")
	}
	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
