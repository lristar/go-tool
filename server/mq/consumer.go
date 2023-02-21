package mq

import (
	"fmt"
	"github.com/streadway/amqp"
	"gitlab.gf.com.cn/hk-common/go-tool/server/logger"
	"strings"
	"time"
)

// Consumer rabbitMQ 消费者
type Consumer struct {
	*Connection
	Key          string
	QueueName    string
	ExchangeName string
	AutoAck      bool
	Exclusive    bool
	NoLocal      bool
	NoWait       bool
	Args         amqp.Table
	QosCount     int
	reReceive    chan interface{}
	Handle       func(*amqp.Delivery) error
}

func (c *Consumer) Receive() error {
	durable := false
	autoDelete := true
	if c.QueueName != "" && !strings.HasPrefix(c.QueueName, "amq.gen-") {
		durable = true
		autoDelete = false
	}
	// 用于重新刷新接收数据的管道
	if c.ExchangeName != "" {
		if strings.HasPrefix(c.QueueName, "amq.gen-") {
			c.QueueName = ""
		}
		queue, err := c.Channel.QueueDeclare(c.QueueName, durable, autoDelete, false, false, nil)
		if err != nil {
			return err
		}
		c.QueueName = queue.Name
		if err = c.Channel.QueueBind(c.QueueName, c.Key, c.ExchangeName, false, nil); err != nil {
			return err
		}
	} else {
		queue, err := c.Channel.QueueDeclare(c.QueueName, durable, autoDelete, false, false, nil)
		if err != nil {
			return err
		}
		c.QueueName = queue.Name
	}
	if err := c.Channel.Qos(c.QosCount, 0, false); err != nil {
		return err
	}
	rev, err := c.Channel.Consume(c.QueueName, "", false, false, false, false, nil)
	if err == nil {
		go func() {
			for {
				select {
				case r, ok := <-rev:
					if ok {
						if err = c.Handle(&r); err != nil {
							for r.Nack(false, true) != nil {
								time.Sleep(time.Millisecond * 200)
								logger.Errorf("Received a message nack false : %s\n", r.MessageId)
							}
						} else if err = r.Ack(false); err == nil {

						} else {
							logger.Errorf("Received a message ack false : %s\n", r.MessageId)
						}
					}
					time.Sleep(time.Millisecond * 100)
				case <-c.reReceive:
					for {
						logger.Infof("Consumer Receive Restart%s\n", c.QueueName)
						err := c.Receive()
						if err == nil {
							return
						}
						time.Sleep(time.Second * 15)
					}
				}
			}
		}()
		logger.Infof("%s work start, [*] Waiting for messages......\n", c.QueueName)
	} else {
		err = fmt.Errorf("%s work start false %v\n", c.QueueName, err)
		logger.Error(err)
		return err
	}
	return nil
}
