package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"gitlab.gf.com.cn/hk-common/go-tool/lib/pool"
	"gitlab.gf.com.cn/hk-common/go-tool/server/logger"
	"gitlab.gf.com.cn/hk-common/go-tool/server/mq"
	mq2 "gitlab.gf.com.cn/hk-common/go-tool/server/mq_two"
	"time"
)

func temp1() {
	mq.InitConnect("amqp://guest:guest@10.68.41.31:5672")
	pool.NewChannelPool(pool.Config{
		InitialCap:  5,
		MaxCap:      20,
		Fac:         mq.Factory,
		IdleTimeout: 0,
	})
	go func() {
		for i := 1; i < 10000; i++ {
			for j := 0; j < 2; j++ {
				pub, err := mq.NewPublish("f_lzy", mq.FANOUT, "")
				if err != nil {
					panic(err)
				}
				if err = pub.Send(fmt.Sprintf("hahahahNo :%d-%d", i, j)); err != nil {
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

func temp2() {
	go func() {
		for i := 1; i < 10000; i++ {
			time.Sleep(time.Second * 10)
			for j := 0; j < 2; j++ {
				if err := mq2.Send(fmt.Sprintf("hahahahNo :%d-%d", i, j), "lzy_test"); err != nil {
					logger.Error(fmt.Errorf("发送失败 :%d-%d", i, j))
				}
			}
		}
	}()
	rb := mq2.NewRabbitMQ("amqp://guest:guest@10.68.41.31:5672")
	c := &mq2.Consumer{QosCount: 1, Queue: &mq2.Queue{Name: "lzy_test", Durable: true}, Handle: hello2}

	if err := rb.Consume(c, true); err != nil {
		panic(err)
	}
	for {
	}
}

func main() {
	temp1()
}
func hello2(r *amqp.Delivery, mq *mq2.RabbitMQ, queueName string) {
	logger.Infof("接收到数据%s", r.Body)
	mq.ResultHandle(r, queueName, false, false, nil)
}

func hello(r *amqp.Delivery) error {
	logger.Infof("接收到数据%s", r.Body)
	return nil
}
