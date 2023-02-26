package main

import (
	"fmt"
	"github.com/lristar/go-tool/server/logger"
	"github.com/lristar/go-tool/server/mq"
	"github.com/streadway/amqp"
	"time"
)

func temp1() {
	if err := mq.InitConnect("amqp://guest:guest@10.68.41.31:5672"); err != nil {
		panic(err)
	}
	pub, err := mq.GetConn().NewPublish("f_lzy", mq.FANOUT, "", true)
	if err != nil {
		panic(err)
	}
	go func() {
		for i := 1; i < 10000; i++ {
			time.Sleep(time.Second * 10)
			for j := 0; j < 2; j++ {
				if err = pub.Send(fmt.Sprintf("hahahahNo :%d-%d", i, j)); err != nil {
					logger.Error(fmt.Errorf("发送失败 :%d-%d", i, j))
				}
			}
		}
	}()
	consumer, err := mq.GetConn().NewConsumer(mq.Consumer{ExchangeName: "f_lzy", QosCount: 0, Handle: hello}, true)
	if err != nil {
		panic(err)
	}
	if err := consumer.Receive(); err != nil {
		panic(err)
	}
	consumer1, err := mq.GetConn().NewConsumer(mq.Consumer{ExchangeName: "f_lzy", QosCount: 0, Handle: hello}, true)
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
