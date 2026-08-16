package main

import (
	"context"
	"log"

	"github.com/segmentio/kafka-go"
	"github.com/roman4k-gg/myGarden/notification-Service/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{cfg.KafkaBroker},
		Topic:    cfg.KafkaTopic,
		GroupID:  cfg.GroupID,
		MinBytes: 10e3,
		MaxBytes: 10e6,
	})
	defer r.Close()

	log.Printf("Notification Service started, consuming topic=%s group=%s", cfg.KafkaTopic, cfg.GroupID)

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error reading message: %v", err)
			continue
		}
		log.Printf("[NOTIFICATION] topic=%s key=%s value=%s", m.Topic, string(m.Key), string(m.Value))
	}
}
