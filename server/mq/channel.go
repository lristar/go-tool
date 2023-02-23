package mq

import (
	"github.com/streadway/amqp"
	"gitlab.gf.com.cn/hk-common/go-tool/server/logger"
	"time"
)

type Channel struct {
	ch *amqp.Channel // 发布/接收 频道
}

func watchChannel(c *Channel, errC chan *amqp.Error, f func()) {
	select {
	case reason, ok := <-errC:
		if ok {
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
		} else {
			logger.Error("错误管道被关闭")
		}
	}
}
