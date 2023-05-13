package main

import (
	"github.com/Shopify/sarama"
	"log"
	"sync"
	"time"
)

var (
	once     sync.Once
	producer sarama.SyncProducer
)

func initProducer(brokers []string) {
	var err error
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Compression = sarama.CompressionSnappy
	config.Producer.Flush.Frequency = 500 * time.Millisecond
	config.Producer.Return.Successes = true

	producer, err = sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
}

func GetKafkaProducer() sarama.SyncProducer {
	once.Do(func() {
		initProducer([]string{"localhost:9092"})
	})
	return producer
}

func SendMessage(topic string, message []byte) error {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}
	partition, offset, err := GetKafkaProducer().SendMessage(msg)
	if err != nil {
		return err
	}
	log.Printf("Produced message to topic %s: partition=%d, offset=%d", topic, partition, offset)

	err = Close()
	if err != nil {
		return err
	}
	return nil
}

func Close() error {
	err := GetKafkaProducer().Close()
	if err != nil {
		log.Printf("Failed to close Kafka producer: %v", err)
		return err
	}

	return nil
}
