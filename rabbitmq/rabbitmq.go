package rabbitmq

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

const (
	maxConn     = 10
	maxCh       = 50
	maxRetries  = 3
	maxInterval = 60 // 最大重连间隔时间（秒）
)

type RabbitMQ struct {
	connPool    chan *amqp.Connection
	chPool      chan *amqp.Channel
	connRetries chan *amqp.Connection
}

func NewRabbitMQ() *RabbitMQ {
	return &RabbitMQ{
		connPool:    make(chan *amqp.Connection, maxConn),
		chPool:      make(chan *amqp.Channel, maxCh),
		connRetries: make(chan *amqp.Connection),
	}
}

func (r *RabbitMQ) Start() {
	var wg sync.WaitGroup
	for i := 0; i < maxConn; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				conn, err := amqp.Dial(os.Getenv("RABBITMQ_URL"))
				if err != nil {
					log.Printf("Failed to connect to RabbitMQ: %s", err)
					time.Sleep(time.Second * 5)
					continue
				}
				log.Printf("Connected to RabbitMQ")
				r.connPool <- conn

				// Reconnect if the connection is closed
				err = <-conn.NotifyClose(make(chan *amqp.Error))
				if err != nil {
					log.Printf("Lost connection to RabbitMQ: %s", err)
				}
				r.connRetries <- conn
			}
		}()
	}

	go func() {
		for {
			select {
			case conn := <-r.connPool:
				for i := 0; i < maxCh; i++ {
					ch, err := conn.Channel()
					if err != nil {
						log.Printf("Failed to create channel: %s", err)
						break
					}
					r.chPool <- ch
				}
			case conn := <-r.connRetries:
				for i := 0; i < maxRetries; i++ {
					interval := time.Duration(i+1) * time.Second
					if interval > maxInterval*time.Second {
						interval = maxInterval * time.Second
					}
					time.Sleep(interval)
					log.Printf("Trying to reconnect to RabbitMQ (%d/%d)", i+1, maxRetries)
					err := conn.Close()
					if err != nil {
						log.Printf("Failed to close connection: %s", err)
					}
					break
				}
			}
		}
	}()
}

func (r *RabbitMQ) GetChannel() (*amqp.Channel, error) {
	ch, ok := <-r.chPool
	if !ok {
		return nil, fmt.Errorf("channel pool closed")
	}
	return ch, nil
}

func (r *RabbitMQ) ReleaseChannel(ch *amqp.Channel) {
	select {
	case r.chPool <- ch:
	default:
		ch.Close()
	}
}

func (r *RabbitMQ) Close() {
	close(r.connPool)
	close(r.ch
