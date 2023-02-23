package rabbitmq

import (
	"testing"
)

func TestRabbitMQ(t *testing.T) {
	r := NewRabbitMQ()
	if r == nil {
		t.Fatalf("Failed to create RabbitMQ instance")
	}

	err := r.Start()
	if err != nil {
		t.Fatalf("Failed to start RabbitMQ: %s", err)
	}

	ch, err := r.GetChannel()
	if err != nil {
		t.Fatalf("Failed to get channel: %s", err)
	}

	err = ch.ExchangeDeclare(
		"test_exchange", // name
		"fanout",        // type
		true,            // durable
		false,           // auto-deleted
		false,           // internal
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		t.Fatalf("Failed to declare exchange: %s", err)
	}

	_, err = ch.QueueDeclare(
		"test_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		t.Fatalf("Failed to declare queue: %s", err)
	}

	err = ch.QueueBind(
		"test_queue",    // name
		"",              // routing key
		"test_exchange", // exchange
		false,           // no-wait
		nil,             // arguments
	)
	if err != nil {
		t.Fatalf("Failed to bind queue: %s", err)
	}

	msg := "Hello, World!"
	err = r.Publish("test_exchange", "", []byte(msg))
	if err != nil {
		t.Fatalf("Failed to publish message: %s", err)
	}

	received, err := r.Consume("test_queue")
	if err != nil {
		t.Fatalf("Failed to consume: %s", err)
	}

	if string(received.Body) != msg {
		t.Fatalf("Received message does not match published message")
	}

	err = r.Close()
	if err != nil {
		t.Fatalf("Failed to close RabbitMQ: %s", err)
	}
}
