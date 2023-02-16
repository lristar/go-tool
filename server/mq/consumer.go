package mq

import (
	"fmt"
	logger "gitlab.gf.com.cn/hk-common/go-tool/server/logger"

	"github.com/getsentry/raven-go"
	"github.com/streadway/amqp"
)

// Consumer rabbitMQ 消费者
type Consumer struct {
	Name      string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Args      amqp.Table
	QosCount  int
	Queue     *Queue
	Handle    func(*amqp.Delivery, *RabbitMQ, string)
}

// Consume 消费队列信息
func (mq *RabbitMQ) Consume(c *Consumer, keepAlive bool) error {
	var (
		err          error
		ch           *amqp.Channel
		deliveryChan <-chan amqp.Delivery
	)

	if err = mq.getConsumeChannel(c.Queue.Name); err == nil {
		ch = mq.consumeChan
		if _, err = ch.QueueDeclare(c.Queue.Name, c.Queue.Durable, c.Queue.AutoDelete, c.Queue.Exclusive, c.Queue.NoWait, c.Queue.Args); err != nil {
			return err
		}
		if c.Queue.Exchange != nil {
			if err = ch.QueueBind(c.Queue.Name, c.Queue.Key, c.Queue.Exchange.Name, c.Queue.NoWait, c.Queue.Args); err != nil {
				return err
			}
		}

		err = ch.Qos(c.QosCount, 0, false)
		if err != nil {
			return err
		}

		if deliveryChan, err = ch.Consume(c.Queue.Name, c.Name, c.AutoAck, c.Exclusive, c.NoLocal, c.NoWait, c.Args); err == nil {
			go func(m *RabbitMQ) {
				for delivery := range deliveryChan {
					logger.Infof("Received a message: %s\n", delivery.Body)
					c.Handle(&delivery, m, c.Queue.Name)
				}
			}(mq)

			if keepAlive {
				mq.consumer = append(mq.consumer, c)
			}
			logger.Infof("%s work start, [*] Waiting for messages......\n", c.Queue.Name)
		}
	}

	return err
}

// ConsumeCancel optimize: 关闭消费者
func (mq *RabbitMQ) ConsumeCancel(c Consumer) error {
	return nil
}

// ResultHandle 处理消费数据的结果 errorReport是否错误上报,resend 是否重发
func (mq *RabbitMQ) ResultHandle(d *amqp.Delivery, queueName string, errorReport, resend bool, err error) {
	msg := string(d.Body)
	if err != nil {
		err = fmt.Errorf("[队列:%s]%s", queueName, err.Error())
		logger.Error(err)
		if errorReport {
			raven.CaptureError(err, nil)
		}
	}
	if err = d.Ack(false); err != nil {
		err = fmt.Errorf("[队列:%s]消息Ack失败,[Msg]:%s;[Error]:%s", queueName, msg, err.Error())
		logger.Error(err)
		raven.CaptureError(err, nil)
	} else if resend {
		if err = mq.ReSend(queueName, d.Body); err != nil {
			err = fmt.Errorf("[队列:%s]重新推送消息失败,[Msg]:%s;[Error]:%s", queueName, msg, err.Error())
			logger.Error(err)
			raven.CaptureError(err, nil)
		}
	}
}
