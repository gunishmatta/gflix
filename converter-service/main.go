package main

import (
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"log"
	"os"
	"os/signal"
	_ "time"
)

const (
	topic          = "video-events"
	consumerGroup  = "video-converter"
	conversionPath = "/tmp/videos/conversion/"
)

type VideoCreatedEvent struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	URL       string `json:"url"`
	CreatedAt int64  `json:"createdAt"`
}
type Message struct {
	EventType string             `json:"eventType"`
	VideoID   primitive.ObjectID `json:"videoID"`
}

type Consumer struct{}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Received message: %v", string(message.Value))
		var event Message

		if err := json.Unmarshal(message.Value, &event); err != nil {
			log.Printf("Failed to unmarshal message: %v", err)
			session.MarkMessage(message, "")
			continue
		}
		if event.EventType != "VIDEO_CREATED" {
			log.Printf("Received event of type %s, ignoring...", event.EventType)
			session.MarkMessage(message, "")
			continue
		}
		// Get video by ID and convert it to other formats
		log.Printf("Video %s created, converting to other formats...", event.VideoID)
		// add your code here to convert video
		session.MarkMessage(message, "")

	}
	return nil
}

func main() {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Version = sarama.V2_6_0_0
	consumerGroup, err := sarama.NewConsumerGroup([]string{"localhost:9092"}, consumerGroup, config)
	if err != nil {
		log.Fatalf("Failed to create consumer group: %v", err)
	}
	defer func(consumerGroup sarama.ConsumerGroup) {
		err := consumerGroup.Close()
		if err != nil {
			log.Fatalf("Failed to close consumer group: %v", err)
		}
	}(consumerGroup)

	go func() {
		for err := range consumerGroup.Errors() {
			log.Fatalf("Consumer error: %v", err)
		}
	}()

	consumer := &Consumer{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start consuming messages
	if err := consumerGroup.Consume(ctx, []string{topic}, consumer); err != nil {
		log.Fatalf("Failed to consume messages: %v", err)
	}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, os.Interrupt)
	<-sigterm
	log.Println("Shutting down consumer...")
}
