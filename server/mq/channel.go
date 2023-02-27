package mq

import (
	"fmt"
	"github.com/lristar/go-tool/server/logger"
	"github.com/streadway/amqp"
	"time"
)

type Channel struct {
	ch *amqp.Channel // 发布/接收 频道
}

type PubBody struct {
	Exchange  string
	Key       string
	Mandatory bool
	Immediate bool
	Body      amqp.Publishing
}

// Close 关闭连接的方法
func (c *Channel) Close() error {
	return c.ch.Close()
}

// Ping 检查连接是否有效的方法
func (c *Channel) Ping() error {
	return nil
}

// Use 应用这个连接
func (c *Channel) Use(v interface{}) error {
	if v == nil {
		return fmt.Errorf("请求为空")
	}
	if body, ok := v.(PubBody); ok {
		return c.ch.Publish(body.Exchange, body.Key, body.Mandatory, body.Immediate, body.Body)
	}
	return fmt.Errorf("请求对象结构体错误")
}

func watchChannel(c *Channel, errC chan *amqp.Error, f func()) {
	select {
	case reason := <-errC:
		logger.Error("管道断开", reason)
		logger.Info("开始重建........")
		time.Sleep(DelayTime)
		for {
			if !conn.conn.IsClosed() {
				var err error
				c.ch, err = conn.conn.Channel()
				if err == nil {
					errM := make(chan *amqp.Error)
					c.ch.NotifyClose(errM)
					if f != nil {
						f()
					}
					go watchChannel(c, errM, f)
					logger.Info("管道重建成功")
					return
				}
				logger.Error("管道重建失败")
				time.Sleep(Delay10S)
			} else {
				time.Sleep(Delay30S)
			}
		}
	}
}
