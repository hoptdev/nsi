package producer

import (
	"context"
	"fmt"
	"log"

	"github.com/segmentio/kafka-go"
)

const (
	brokerAddress = "localhost:9092"
)

func Write(topic string, messageContent string) {
	Init(topic)

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
}

func Init(topic string) {

	conn, err := kafka.DialLeader(context.Background(), "tcp", brokerAddress, topic, 0)
	if err != nil {
		panic(err)
	}

	defer conn.Close()
}
