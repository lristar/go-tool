package mq

import (
	"fmt"
	"github.com/streadway/amqp"
	"sync"
	"time"
)

type RabbitConfig struct {
	MQUri string `json:"mq_uri"`
}

// RabbitMQ rabbit mq
type RabbitMQ struct {
	uri         string
	conn        *amqp.Connection
	publishChan *amqp.Channel // 发布频道
	consumeChan *amqp.Channel // 消费频道
	consumer    []*Consumer   // 消费者
}

var (
	mx sync.Mutex
	// DefaultMQURL 默认MQURL
	_rabbitMQ *RabbitMQ
)

func NewDefaultRabbitMQ(url string) error {
	_rabbitMQ = NewRabbitMQ(url)
	return _rabbitMQ.getConn()
}

// GetMQ 获取默认
func getDefaultRabbitMQ() (*RabbitMQ, error) {
	if _rabbitMQ == nil {
		return nil, fmt.Errorf("_rabbitMQ is nil")
	}
	mq := &RabbitMQ{
		uri:         _rabbitMQ.uri,
		conn:        _rabbitMQ.conn,
		publishChan: nil,
		consumeChan: nil,
		consumer:    nil,
	}
	if _rabbitMQ.conn.IsClosed() {
		if ok := mx.TryLock(); ok {
			defer mx.Unlock()
			_, err := mq.getNewConn()
			if err != nil {
				return nil, err
			}
		} else {
			time.Sleep(time.Millisecond * 500)
			return getDefaultRabbitMQ()
		}
	}
	return mq, nil
}

// NewRabbitMQ 创建新的RabbitMQ
func NewRabbitMQ(uri string) *RabbitMQ {
	return &RabbitMQ{
		uri: uri,
	}
}
