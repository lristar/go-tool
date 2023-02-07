package mq

import (
	"github.com/streadway/amqp"
)

// Queue 队列
type Queue struct {
	Name       string
	Durable    bool
	AutoDelete bool
	Exclusive  bool
	NoWait     bool
	Args       amqp.Table
	Key        string
	Exchange   *Exchange
}

// Producer 生产者
type Producer struct {
	Exchange Exchange
	Queue    Queue
	Msg      []byte
}

// Send 信息简易发送
func (mq *RabbitMQ) Send(queueName string, body []byte) error {
	// ensure queue
	queue := &Queue{
		Name:    queueName,
		Durable: true,
	}

	err := mq.ensureQueue(queue)
	if err != nil {
		return err
	}

	if err = mq.GetPublishChannel(); err == nil {
		msg := amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		}
		err = mq.publishChan.Publish("", queueName, false, false, msg)
	}

	return err
}

// ReSend 重新推送消息
func (mq *RabbitMQ) ReSend(queueName string, body []byte) error {
	err := mq.GetPublishChannel()
	if err == nil {
		msg := amqp.Publishing{
			ContentType: "text/plain",
			Body:        body,
		}
		err = mq.publishChan.Publish("", queueName, false, false, msg)
	}

	return err
}

// SendMsg 发送信息至队列
func (mq *RabbitMQ) SendMsg(queue *Queue, msg *amqp.Publishing, mandatory, immediate bool) error {
	var err error
	if err = mq.EnsureExchange(queue.Exchange); err != nil {
		return err
	}

	// ensure queue
	if err = mq.ensureQueue(queue); err != nil {
		return err
	}

	if err = mq.GetPublishChannel(); err == nil {
		err = mq.publishChan.Publish(queue.Exchange.Name, queue.Key, mandatory, immediate, *msg)
	}
	return err
}

// PublishMsg 发布消息至交换机
func (mq *RabbitMQ) PublishMsg(exc *Exchange, key string, msg *amqp.Publishing, mandatory, immediate bool) error {
	var err error
	if err = mq.EnsureExchange(exc); err != nil {
		return err
	}

	if err = mq.GetPublishChannel(); err == nil {
		err = mq.publishChan.Publish(exc.Name, key, mandatory, immediate, *msg)
	}
	return err
}

// EnsureQueue 保证队列存在
func (mq *RabbitMQ) ensureQueue(q *Queue) error {
	var err error

	if err = mq.GetPublishChannel(); err == nil {
		if _, err = mq.publishChan.QueueDeclare(q.Name, q.Durable, q.AutoDelete, q.Exclusive, q.NoWait, q.Args); err != nil {
			return err
		}
		if q.Exchange != nil {
			err = mq.publishChan.QueueBind(q.Name, q.Key, q.Exchange.Name, q.NoWait, q.Args)
		}
	}

	return err
}
