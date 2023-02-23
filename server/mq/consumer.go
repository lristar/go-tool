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
	*Channel
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

func NewConsumer(con Consumer, watchClose bool) (*Consumer, error) {
	reReceive := make(chan interface{}, 1)
	c, err := conn.newChannel()
	if err != nil {
		return nil, err
	}
	// 用于重新刷新接收数据的管道
	f := func() {
		reReceive <- struct{}{}
	}
	if watchClose {
		errM := make(chan *amqp.Error)
		c.ch.NotifyClose(errM)
		go watchChannel(c, errM, f)
	}
	return &Consumer{
		Channel:      c,
		QueueName:    con.QueueName,
		ExchangeName: con.ExchangeName,
		AutoAck:      con.AutoAck,
		Exclusive:    con.Exclusive,
		NoLocal:      con.NoLocal,
		NoWait:       con.NoWait,
		Args:         con.Args,
		QosCount:     con.QosCount,
		reReceive:    reReceive,
		Handle:       con.Handle,
	}, nil
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
		queue, err := c.ch.QueueDeclare(c.QueueName, durable, autoDelete, false, false, nil)
		if err != nil {
			return err
		}
		c.QueueName = queue.Name
		if err = c.ch.QueueBind(c.QueueName, c.Key, c.ExchangeName, false, nil); err != nil {
			return err
		}
	} else {
		queue, err := c.ch.QueueDeclare(c.QueueName, durable, autoDelete, false, false, nil)
		if err != nil {
			return err
		}
		c.QueueName = queue.Name
	}
	if err := c.ch.Qos(c.QosCount, 0, false); err != nil {
		return err
	}
	rev, err := c.ch.Consume(c.QueueName, "", false, false, false, false, nil)
	if err == nil {
		go func() {
			for {
				select {
				case r, ok := <-rev:
					if ok {
						if c.Handle == nil {
							if err = r.Ack(false); err != nil {
								logger.Errorf("Received a message ack false : %s\n", r.MessageId)
							}
						} else if err = c.Handle(&r); err != nil {
							for i := 0; i < 3; i++ {
								if r.Nack(false, true) == nil {
									break
								}
								time.Sleep(time.Millisecond * 200)
								logger.Errorf("Received a message nack false : %s\n", r.MessageId)
							}
						} else if err = r.Ack(false); err == nil {

						} else {
							logger.Errorf("Received a message ack false : %s\n", r.MessageId)
						}
					}
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
