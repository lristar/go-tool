package mq

import (
	"github.com/streadway/amqp"
	"gitlab.gf.com.cn/hk-common/go-tool/server/logger"
	"time"
)

const (
	DelayTime = 2 * time.Second

	// 交换机类型
	FANOUT = "fanout" //发布订阅模式
	TOPIC  = "topic"  // 主题模式
	DIRECT = "direct" // 路由模式
)

var (
	conn *Connection
	// 连接错误管道
	errConnChannel chan *amqp.Error
)

type Connection struct {
	url     string
	conn    *amqp.Connection
	Channel *amqp.Channel // 发布/接收 频道
}

// InitConnect 自带重连机制
func InitConnect(url string) error {
	var err error
	con, err := amqp.DialConfig(url, amqp.Config{
		Heartbeat: time.Second * 5,
		Locale:    "en_US",
	})
	if err != nil {
		return err
	}
	errConnChannel = make(chan *amqp.Error)
	con.NotifyClose(errConnChannel)
	// 添加重连机制
	go watchConn()
	conn = &Connection{
		url:     url,
		conn:    con,
		Channel: nil,
	}
	return nil
}

func (c *Connection) NewPublish(exchange, exchangeType, key string, watchClose bool) (*Producer, error) {
	if _, err := c.newChannel(); err != nil {
		return nil, err
	}
	if exchangeType != "" && exchange != "" {
		if err := c.Channel.ExchangeDeclare(exchange, exchangeType, true, false, false, false, nil); err != nil {
			return nil, err
		}
	}
	if watchClose {
		errM := make(chan *amqp.Error)
		c.Channel.NotifyClose(errM)
		go watchChannel(c, errM, nil)
	}
	return &Producer{
		Connection:   c,
		Exchange:     exchange,
		ExchangeType: exchangeType,
		Key:          key,
		Mandatory:    false,
		Immediate:    false,
	}, nil
}

func (c *Connection) NewConsumer(con Consumer, watchClose bool) (*Consumer, error) {
	reReceive := make(chan interface{}, 1)
	if _, err := c.newChannel(); err != nil {
		return nil, err
	}
	// 用于重新刷新接收数据的管道
	f := func() {
		reReceive <- struct{}{}
	}
	if watchClose {
		errM := make(chan *amqp.Error)
		c.Channel.NotifyClose(errM)
		go watchChannel(c, errM, f)
	}
	return &Consumer{
		Connection:   c,
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

func GetConn() *Connection {
	return &Connection{
		url:     conn.url,
		conn:    conn.conn,
		Channel: nil,
	}
}

func watchConn() {
	logger.Infof("开启连接监听")
	select {
	case reason, ok := <-errConnChannel:
		if !ok {
			logger.Error("连接连接关闭")
			return
		}
		logger.Error("err is:", reason)
		for ok {
			time.Sleep(DelayTime)
			con, err := amqp.Dial(conn.url)
			if err == nil {
				conn.conn = con
				conn.conn.ConnectionState()
				errConnChannel = make(chan *amqp.Error)
				con.NotifyClose(errConnChannel)
				go watchConn()
				logger.Info("连接重连成功")
				return
			}
			logger.Error("连接重连失败")
		}
	}
}

func (c *Connection) newChannel() (*Connection, error) {
	var err error
	c.Channel, err = c.conn.Channel()
	if err != nil {
		return c, err
	}
	return c, nil
}

func watchChannel(c *Connection, errC chan *amqp.Error, f func()) {
	select {
	case reason, ok := <-errC:
		if ok {
			logger.Error("管道断开", reason)
			logger.Info("开始重建........")
			time.Sleep(DelayTime)
			for {
				if !GetConn().conn.IsClosed() {
					c.conn = GetConn().conn
					var err error
					c.Channel, err = c.conn.Channel()
					if err == nil {
						errM := make(chan *amqp.Error)
						c.Channel.NotifyClose(errM)
						if f != nil {
							f()
						}
						go watchChannel(c, errM, f)
						logger.Info("管道重建成功")
						return
					}
					logger.Error("管道重建失败")
					time.Sleep(time.Second * 10)
				} else {
					time.Sleep(time.Second * 30)
				}
			}
		} else {
			logger.Error("错误管道被关闭")
		}
	}
}

type IRecover interface {
	handle()
	setNext()
}
