package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
)

type OrderPlacer struct {
	producer  *kafka.Producer
	topic     string
	deliverch chan kafka.Event
}

func NewOrderPlacer(p *kafka.Producer, topic string) *OrderPlacer {
	return &OrderPlacer{
		producer:  p,
		topic:     topic,
		deliverch: make(chan kafka.Event, 10000),
	}
}

func (op *OrderPlacer) ingestData(payload []byte) error {
	err := op.producer.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &op.topic,
			Partition: kafka.PartitionAny,
		},
		Value: payload,
	},
		op.deliverch,
	)

	if err != nil {
		log.Fatal(err)
	}

	e := <-op.deliverch
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		fmt.Printf("Delivery failed: %v\n", m.TopicPartition.Error)
	} else {
		fmt.Printf("Delivered message to topic %s [%d] at offset %v\n",
			*m.TopicPartition.Topic, m.TopicPartition.Partition, m.TopicPartition.Offset)
	}

	return nil
}

func main() {
	topic := "unprocessed_data"
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "34.100.240.235",
		"client.id":         "http-to-kafka",
		"acks":              "all",
	})

	if err != nil {
		log.Fatalf("Failed to create producer: %s\n", err)
		return
	}

	op := NewOrderPlacer(p, topic)

	r := gin.Default()
	r.POST("/", func(c *gin.Context) {
		var myjson map[string]interface{}

		if err := c.BindJSON(&myjson); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		jsonBytes, err := json.Marshal(myjson)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error encoding JSON"})
			return
		}

		if err := op.ingestData(jsonBytes); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to place order in Kafka"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Sent to Kafka"})
	})

	r.Run(":8080")
}
