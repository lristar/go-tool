package rabbitmq

import (
	"log"
	"testing"
)

func TestRabbitMQ(t *testing.T) {
	rabbit, err := NewRabbit("amqp://guest:guest@localhost:5672/", 5)
	if err != nil {
		log.Fatalf("Failed to create rabbit instance: %s", err.Error())
	}

	defer rabbit.Close()

	msg := []byte("Hello, RabbitMQ!")
	err = rabbit.Publish("test_queue", msg)
	if err != nil {
		log.Fatalf("Failed to publish message: %s", err.Error())
	}
}
