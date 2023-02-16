package mq

import (
	"github.com/getsentry/raven-go"
	logger "gitlab.gf.com.cn/hk-common/go-tool/server/logger"
	"time"

	"github.com/streadway/amqp"
)

// getNewChannel 获取新的Channel
func (mq *RabbitMQ) getNewChannel() (*amqp.Channel, error) {
	if e := mq.getConn(); e != nil {
		return nil, e
	}
	return mq.conn.Channel()
}

// GetPublishChannel get publish channel
func (mq *RabbitMQ) getPublishChannel() error {
	var e error
	if mq.publishChan == nil {
		mq.publishChan, e = mq.getNewChannel()
	}

	return e
}

// GetConsumeChannel get consume channel
func (mq *RabbitMQ) getConsumeChannel(workName string) error {
	var (
		ch  *amqp.Channel
		err error
	)

	if mq.consumeChan == nil {
		if err = mq.getConn(); err == nil {
			ch, err = mq.conn.Channel()
			mq.consumeChan = ch
			closeErrorChan := ch.NotifyClose(make(chan *amqp.Error))
			go func(mq *RabbitMQ) {
				ce := <-closeErrorChan
				// mq.consumeChan = nil
				mq.Close()
				logger.Errorf("Channel Cancel,%s", ce)
				raven.CaptureError(ce, nil) // 错误上报
				// consume channel 尝试重新开启频道
				for i := 1; ; i++ {
					time.Sleep(30 * time.Second)
					if e := mq.getConsumeChannel(workName); e != nil {
						logger.Errorf("%s mq第%d次连接失败,30s后尝试重新启动", workName, i)
					} else if e = mq.channelRecover(); e != nil {
						logger.Error(e)
						raven.CaptureError(e, nil) // 频道重开成功，消费者重新监听失败
					} else {
						logger.Infof("%s mq重连成功", workName)
						break
					}
					mq.Close()
				}
			}(mq)
		}
	}
	return err
}

// channelRecover 频道恢复(恢复消费者监听)
// 返回最后一个消费者恢复失败的错误
func (mq *RabbitMQ) channelRecover() error {
	var err error
	for _, consumer := range mq.consumer {
		if e := mq.Consume(consumer, false); e != nil {
			err = e
		}
	}
	return err
}
