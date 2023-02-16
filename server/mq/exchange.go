package mq

import "github.com/streadway/amqp"

// ExchangeKind 交换机类型
type ExchangeKind string

const (
	// ExcKindDirect direct Kind Exchange
	ExcKindDirect ExchangeKind = "direct"
	// ExcKindFanout fanout Kind Exchange
	ExcKindFanout ExchangeKind = "fanout"
	// ExcKindTopic topic Kind Exchange
	ExcKindTopic ExchangeKind = "topic"
	// ExcKindHeaders headers Kind Exchange
	ExcKindHeaders ExchangeKind = "headers"
)

// Exchange 交换机
type Exchange struct {
	Name       string
	Kind       ExchangeKind
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqp.Table
}

// DefaultExchange 默认Exchange
func DefaultExchange() *Exchange {
	return &Exchange{"", ExcKindDirect, true, false, false, false, nil}
}

// EnsureExchange 保证交换机存在
func (mq *RabbitMQ) ensureExchange(exc *Exchange) error {
	if e := mq.getPublishChannel(); e != nil {
		return e
	}
	return mq.publishChan.ExchangeDeclare(exc.Name, string(exc.Kind), exc.Durable, exc.AutoDelete, exc.Internal, exc.NoWait, exc.Args)
}
