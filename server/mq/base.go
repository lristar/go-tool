package mq

import (
	configs "platform-backend/config"

	"github.com/streadway/amqp"
)

// RabbitMQ rabbit mq
type RabbitMQ struct {
	uri         string
	conn        *amqp.Connection
	publishChan *amqp.Channel // 发布频道
	consumeChan *amqp.Channel // 消费频道
	consumer    []*Consumer   // 消费者
}

var (
	// DefaultMQURL 默认MQURL
	DefaultMQURL = configs.GetConf("MQUri")
	_rabbitMQ    = RabbitMQ{DefaultMQURL, nil, nil, nil, nil}
)

// 初始化单例RMQ
func init() {
	if err := _rabbitMQ.GetConn(); err != nil {
		panic(err)
	}
}

// GetDefaultRabbitMQ 获取RabbitMQ单例
func GetDefaultRabbitMQ() *RabbitMQ { return &_rabbitMQ }

// NewRabbitMQ 创建新的RabbitMQ
func NewRabbitMQ(uri string) *RabbitMQ {
	return &RabbitMQ{
		uri: uri,
	}
}
