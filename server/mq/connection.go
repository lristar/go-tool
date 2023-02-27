package mq

import (
	"github.com/lristar/go-tool/lib/pool"
	"github.com/lristar/go-tool/server/logger"
	"github.com/streadway/amqp"
	"sync"
	"time"
)

const (
	DelayTime = 2 * time.Second
	Delay10S  = 10 * time.Second
	Delay30S  = 30 * time.Second
	// 交换机类型
	FANOUT = "fanout" //发布订阅模式
	TOPIC  = "topic"  // 主题模式
	DIRECT = "direct" // 路由模式
)

var (
	conn *Connection
	// 连接错误管道
	errConnChannel chan *amqp.Error
	once           sync.Once
)

type Connection struct {
	url  string
	conn *amqp.Connection
}

// InitConnect 自带重连机制
func InitConnect(url string) {
	once.Do(func() {
		var err error
		con, err := amqp.DialConfig(url, amqp.Config{
			Heartbeat: time.Second * 30,
			Locale:    "en_US",
		})
		if err != nil {
			panic(err)
		}
		errConnChannel = make(chan *amqp.Error)
		con.NotifyClose(errConnChannel)
		// 添加重连机制
		go watchConn()
		conn = &Connection{
			url:  url,
			conn: con,
		}
	})
}

func watchConn() {
	logger.Infof("开启连接监听")
	select {
	case reason := <-errConnChannel:
		for {
			logger.Error("err is:", reason)
			time.Sleep(DelayTime)
			con, err := amqp.Dial(conn.url)
			if err != nil {
				logger.Error("连接重连失败")
				continue
			}
			conn.conn = con
			errConnChannel = make(chan *amqp.Error)
			con.NotifyClose(errConnChannel)
			go watchConn()
			logger.Info("连接重连成功")
			return
		}
	}
}
func (c *Connection) newChannel() (*Channel, error) {
	var err error
	ch, err := conn.conn.Channel()
	if err != nil {
		return nil, err
	}
	return &Channel{ch: ch}, nil
}

func Factory() (pool.IConn, error) {
	var err error
	ch, err := conn.newChannel()
	if err != nil {
		return nil, err
	}
	return &Channel{ch: ch.ch}, nil
}

// IRecover 是想着用策略模式来recover,待实现
type IRecover interface {
	handle()
	setNext()
}
