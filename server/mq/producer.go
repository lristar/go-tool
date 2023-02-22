package mq

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"reflect"
)

type Producer struct {
	*Connection
	Exchange     string
	Key          string // queueName or key
	ExchangeType string
	Mandatory    bool
	Immediate    bool
}

func (p *Producer) Send(bodys interface{}) error {
	if p.Channel == nil {
		return fmt.Errorf("channel为空")
	}
	defer p.Channel.Close()
	bType := reflect.TypeOf(bodys)
	realDatas := make([]interface{}, 0)
	if bType.Kind() != reflect.Slice {
		realDatas = append(realDatas, bodys)
	} else {
		value := reflect.ValueOf(bodys)
		for i := 0; i < value.Len(); i++ {
			realDatas = append(realDatas, value.Index(i).Interface())
		}
	}
	for _, body := range realDatas {
		msg := ""
		if data, ok := body.(string); ok {
			msg = data
		} else {
			bytes, _ := json.Marshal(body)
			msg = string(bytes)
		}
		err := p.Channel.Publish(p.Exchange, p.Key, p.Mandatory, p.Immediate,
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(msg),
				//Expiration:   "10000", // 设置过期时间
			})
		if err != nil {
			return err
		}
	}
	return nil
}
