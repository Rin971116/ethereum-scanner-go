package kafka

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/segmentio/kafka-go"
)

var Writer *kafka.Writer

func InitKafkaWriter(brokers []string, topic string) {
	Writer = &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}
}

func Produce(entityTx interface{}) error {
	bytes, err := json.Marshal(entityTx)
	if err != nil {
		return err
	}
	msg := kafka.Message{Value: bytes}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return Writer.WriteMessages(ctx, msg)
}

func CloseKafkaWriter() {
	if Writer != nil {
		err := Writer.Close()
		if err != nil {
			log.Printf("關閉 kafka writer 錯誤: %v", err)
		} else {
			log.Println("kafka writer closed")
		}
	}
}
