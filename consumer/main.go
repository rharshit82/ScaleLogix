package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	totalMessagesReceived    int
	totalMessagesValidated   int
	totalMessagesInvalidated int
	totalDataSize            int64
)

type EventCommon struct {
	Event string `json:"event"`
}

type SimpleEvent struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age,omitempty"` // omitempty allows Age to be optional
}

func validateSimpleEvent(event SimpleEvent) bool {
	return event.ID != 0 && event.Name != ""
}

func handleMessage(msg *kafka.Message, producer *kafka.Producer, producerTopic string) {
    totalMessagesReceived++
    totalDataSize += int64(len(msg.Value))

    var event SimpleEvent
    err := json.Unmarshal(msg.Value, &event)
    if err != nil {
        log.Printf("Error unmarshalling message to SimpleEvent struct: %v", err)
        totalMessagesInvalidated++
        return
    }

    if validateSimpleEvent(event) {
        produceMessage(&event, producer, producerTopic)
    } else {
        log.Printf("Validation failed for SimpleEvent: %v", event)
        totalMessagesInvalidated++
    }
}


func produceMessage(event interface{}, producer *kafka.Producer, topic string) {
	jsonData, err := json.Marshal(event)
	if err != nil {
		log.Printf("Error marshalling event: %v", err)
		return
	}
	producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          jsonData,
	}, nil)
	totalMessagesValidated++
}

func main() {
	consumerGroup := "foo_data"
	topic := "unprocessed_data"
	bootstrapServers := "10.160.0.5"
	producerTopic := "validated_data"

	metricsTicker := time.NewTicker(10 * time.Second)
	defer metricsTicker.Stop()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			select {
			case <-metricsTicker.C:
				printMetrics()
			case <-sigChan:
				printMetrics()
				os.Exit(0)
			}
		}
	}()

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
		"group.id":          consumerGroup,
		"auto.offset.reset": "earliest",
	})
	if err != nil {
		log.Fatalf("Failed to create consumer: %s", err)
	}
	defer consumer.Close()

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServers,
	})
	if err != nil {
		log.Fatalf("Failed to create producer: %s", err)
	}
	defer producer.Close()

	err = consumer.Subscribe(topic, nil)
	if err != nil {
		log.Fatalf("Subscribe error: %s", err)
	}

	go func() {
		<-sigChan
		printMetrics()
		os.Exit(0)
	}()

	for {
		msg, err := consumer.ReadMessage(-1)
		if err != nil {
			log.Printf("Consumer error: %v (%v)\n", err, msg)
			continue
		}
		handleMessage(msg, producer, producerTopic)
	}
}

func printMetrics() {
	log.Printf("Total Messages Received: %d\n", totalMessagesReceived)
	log.Printf("Total Messages Validated: %d\n", totalMessagesValidated)
	log.Printf("Total Messages Invalidated: %d\n", totalMessagesInvalidated)
	log.Printf("Total Data Size: %d bytes\n", totalDataSize)
}
