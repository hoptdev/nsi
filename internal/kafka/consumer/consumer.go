package consumer

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/segmentio/kafka-go"
)

const (
	brokerAddress = "localhost:9092"
)

func Consume(ctx context.Context, topic string) {

	groupID := fmt.Sprintf("consumer-group-%d", os.Getpid())

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:     []string{brokerAddress},
		Topic:       topic,
		GroupID:     groupID,
		StartOffset: kafka.LastOffset,
	})
	defer reader.Close()

	fmt.Printf("Consumer %s started. Topic=%s, Waiting for messages...\n", groupID, topic)

	for {
		select {
		case <-ctx.Done():
			log.Println("Context canceled, stopping consumer")
			return
		default:
			msg, err := reader.ReadMessage(context.Background())
			if err != nil {
				log.Fatalf("Error reading message: %v", err)
			}

			reader.CommitMessages(context.Background(), msg)
			fmt.Printf("[%s]Received message: %s (partition=%d offset=%d)\n", topic,
				string(msg.Value), msg.Partition, msg.Offset)
		}
	}
}
