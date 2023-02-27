package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"gitlab.gf.com.cn/hk-common/go-tool/lib/pool"
	"gitlab.gf.com.cn/hk-common/go-tool/server/logger"
	"gitlab.gf.com.cn/hk-common/go-tool/server/mq"
	"time"
)

func temp1() {
	mq.InitConnect("amqp://guest:guest@10.68.41.31:5672")
	if _, err := pool.InitPool(pool.Config{
		InitialCap:  5,
		MaxCap:      20,
		Fac:         mq.Factory,
		IdleTimeout: 20,
	}); err != nil {
		panic(err)
	}
	go func() {
		for {
			time.Sleep(time.Second)
			p, _ := pool.GetPool()
			logger.Infof("积极连接%d -- 管道存在%d", p.GetActive())
		}
	}()
	go func() {
		for i := 1; i < 10000; i++ {
			for j := 0; j < 2; j++ {
				if err := mq.Send("f_lzy", "", fmt.Sprintf("hahahahNo :%d-%d", i, j)); err != nil {
					logger.Error(fmt.Errorf("发送失败 :%d-%d", i, j))
				}
			}
			time.Sleep(time.Second * 10)
		}
	}()
	consumer, err := mq.NewConsumer(mq.Consumer{ExchangeName: "f_lzy", QosCount: 0, Handle: hello})
	if err != nil {
		panic(err)
	}
	if err := consumer.Receive(); err != nil {
		panic(err)
	}
	consumer1, err := mq.NewConsumer(mq.Consumer{ExchangeName: "f_lzy", QosCount: 0, Handle: hello})
	if err != nil {
		panic(err)
	}
	if err := consumer1.Receive(); err != nil {
		panic(err)
	}
	for {
	}
}

func main() {
	temp1()
}

func hello(r *amqp.Delivery) error {
	logger.Infof("接收到数据%s", r.Body)
	return nil
}
