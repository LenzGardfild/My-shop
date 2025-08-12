package kafka

import (
	"context"
	"encoding/json"
	"log"

	"my-shop/internal/cache"
	"my-shop/internal/db"
	"my-shop/internal/model"

	"github.com/IBM/sarama"
)

type Consumer struct {
	db    *db.DB
	cache *cache.OrderCache
}

func NewConsumer(db *db.DB, cache *cache.OrderCache) *Consumer {
	return &Consumer{db: db, cache: cache}
}

func (c *Consumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (c *Consumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (c *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var order model.Order
		err := json.Unmarshal(msg.Value, &order)
		if err != nil {
			log.Printf("Ошибка парсинга сообщения: %v", err)
			sess.MarkMessage(msg, "")
			continue
		}
		err = c.db.SaveOrder(context.Background(), &order)
		if err != nil {
			log.Printf("Ошибка сохранения заказа: %v", err)
			sess.MarkMessage(msg, "")
			continue
		}
		c.cache.Set(order.OrderUID, &order)
		sess.MarkMessage(msg, "")
	}
	return nil
}

func RunKafkaConsumer(brokers []string, group, topic string, db *db.DB, cache *cache.OrderCache) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumer := NewConsumer(db, cache)
	groupConsumer, err := sarama.NewConsumerGroup(brokers, group, config)
	if err != nil {
		log.Fatalf("Kafka consumer error: %v", err)
	}
	go func() {
		for {
			err := groupConsumer.Consume(context.Background(), []string{topic}, consumer)
			if err != nil {
				log.Printf("Kafka consume error: %v", err)
			}
		}
	}()
}
