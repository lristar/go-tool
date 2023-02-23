package rabbitmq

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type Rabbit struct {
	url            string             // RabbitMQ连接字符串
	conn           *amqp.Connection   // RabbitMQ连接
	channelPool    chan *amqp.Channel // Channel连接池
	channelNum     int                // Channel连接数
	reconnectCount int                // 重连次数
	reconnectMutex sync.RWMutex       // 重连锁
	connected      bool               // 是否已连接
}

// NewRabbit 创建一个新的Rabbit实例
func NewRabbit(url string, channelNum int) (*Rabbit, error) {
	rabbit := &Rabbit{
		url:         url,
		channelPool: make(chan *amqp.Channel, channelNum),
		channelNum:  channelNum,
		connected:   false,
	}

	err := rabbit.reconnect()
	if err != nil {
		return nil, err
	}

	go rabbit.keepAlive()

	return rabbit, nil
}

// keepAlive 维持RabbitMQ连接的心跳
func (r *Rabbit) keepAlive() {
	for {
		select {
		case <-time.After(5 * time.Second):
			if r.connected {
				if r.conn.IsClosed() {
					r.reconnect()
				}
			}
		}
	}
}

// reconnect 重连RabbitMQ
func (r *Rabbit) reconnect() error {
	r.reconnectMutex.Lock()
	defer r.reconnectMutex.Unlock()

	if r.connected {
		return nil
	}

	conn, err := amqp.Dial(r.url)
	if err != nil {
		return err
	}

	channelPool := make(chan *amqp.Channel, r.channelNum)

	for i := 0; i < r.channelNum; i++ {
		channel, err := conn.Channel()
		if err != nil {
			return err
		}

		channelPool <- channel
	}

	r.conn = conn
	r.channelPool = channelPool
	r.connected = true
	r.reconnectCount = 0

	return nil
}

// GetChannel 获取一个Channel连接
func (r *Rabbit) GetChannel() (*amqp.Channel, error) {
	if !r.connected {
		return nil, fmt.Errorf("connection is closed")
	}

	channel, ok := <-r.channelPool
	if !ok {
		return nil, fmt.Errorf("channel pool is closed")
	}

	return channel, nil
}

// ReturnChannel 归还一个Channel连接
func (r *Rabbit) ReturnChannel(channel *amqp.Channel) {
	if !r.connected {
		return
	}

	r.channelPool <- channel
}

// Close 关闭RabbitMQ连接
func (r *Rabbit) Close() error {
	r.reconnectMutex.Lock()
	defer r.reconnectMutex.Unlock()

	if !r.connected {
		return nil
	}

	err := r.conn.Close()
	if err != nil {
		return err
	}

	r.connected = false

	return nil
}

// Consume 消费RabbitMQ队列消息
func (r *Rabbit) Consume(queueName string, handler func(msg []byte) error) error {
	channel, err := r.GetChannel()
	if err != nil {
		return err
	}
	defer r.ReturnChannel(channel)

	err = channel.Qos(1, 0, false)
	if err != nil {
		return err
	}
	msgs, err := channel.Consume(
		queueName,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for {
		select {
		case msg := <-msgs:
			log.Printf("received message: %s", string(msg.Body))
			err := handler(msg.Body)
			if err != nil {
				log.Printf("failed to handle message: %s, error: %s", string(msg.Body), err.Error())
				return err
			}
			err = msg.Ack(false)
			if err != nil {
				log.Printf("failed to acknowledge message: %s, error: %s", string(msg.Body), err.Error())
				return err
			}
		}
	}
}

// Publish 发布消息到RabbitMQ队列
func (r *Rabbit) Publish(queueName string, msg []byte) error {
	channel, err := r.GetChannel()
	if err != nil {
		return err
	}
	defer r.ReturnChannel(channel)

	err = channel.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         msg,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
