package main

import (
	"context"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

func main() {
	topic := "favorites_notifications"
	
	broker := os.Getenv("KAFKA_BROKER")
	if broker == "" {
		broker = "localhost:9092"
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{broker},
		Topic:     topic,
		Partition: 0,
		MinBytes:  10e3,
		MaxBytes:  10e6,
	})
	defer r.Close()

	log.Println("Notification Service started. Listening for Kafka messages...")

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error while reading message: %v\n", err)
			continue
		}
		
		log.Printf("========================================================\n")
		log.Printf("[EMAIL/PUSH SIMULATION] Message received from Kafka!\n")
		log.Printf("Topic: %s | Message: %s\n", m.Topic, string(m.Value))
		log.Printf("========================================================\n")
	}
}
