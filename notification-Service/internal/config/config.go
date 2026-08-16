package config

import (
	"errors"
	"os"
)

type Config struct {
	KafkaBroker string
	KafkaTopic  string
	GroupID     string
}

func Load() (*Config, error) {
	cfg := &Config{
		KafkaBroker: os.Getenv("KAFKA_BROKER"),
		KafkaTopic:  getEnv("KAFKA_TOPIC", "favorites_notifications"),
		GroupID:     getEnv("KAFKA_GROUP_ID", "notification-service-group"),
	}
	if cfg.KafkaBroker == "" {
		return nil, errors.New("KAFKA_BROKER env var is required")
	}
	return cfg, nil
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
