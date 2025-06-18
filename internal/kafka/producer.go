package producer

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

const (
	brokerAddress = "localhost:9092"
)

func Write(topic string, messageContent string) {
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{brokerAddress},
		Topic:   topic,
	})

	defer writer.Close()

	// Отправляем сообщение
	err := writer.WriteMessages(context.Background(),
		kafka.Message{
			Key:   []byte("event"),
			Value: []byte(messageContent),
		},
	)

	if err != nil {
		log.Fatalf("Failed to send message: %v", err)
	}

	fmt.Println("Message sent successfully!")
	time.Sleep(2 * time.Second) // Даем время на доставку
}
