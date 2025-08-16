package main

import (
	"fmt"
	"os"

	"github.com/IBM/sarama"
)

func main() {
	brokers := []string{"localhost:9092"}
	topic := "orders"
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run producer.go <order.json>")
		os.Exit(1)
	}
	jsonFile := os.Args[1]
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		fmt.Printf("Ошибка чтения файла: %v\n", err)
		os.Exit(1)
	}
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		fmt.Printf("Ошибка создания продюсера: %v\n", err)
		os.Exit(1)
	}
	defer producer.Close()
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(data),
	}
	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		fmt.Printf("Ошибка отправки сообщения: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("Сообщение отправлено в partition %d, offset %d\n", partition, offset)
}
