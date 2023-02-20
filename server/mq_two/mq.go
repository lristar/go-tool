package mq

import (
	"encoding/json"
	"reflect"

	"github.com/streadway/amqp"
)

// Send  发送到队列支持string/map/struct或对应的数组
func Send(bodys interface{}, queueName string) error {
	mq := NewRabbitMQ("amqp://guest:guest@10.68.41.31:5672")
	var err error
	// ensure queue
	queue := &Queue{
		Name:    queueName,
		Durable: true,
	}
	defer mq.Close()
	if err = mq.ensureQueue(queue); err != nil {
		return err
	}

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

		err = mq.publishChan.Publish(
			"",        // exchange
			queueName, // routing key
			false,     // mandatory
			false,     // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg),
			})
		if err != nil {
			return err
		}
	}
	return nil
}

func SendExchange(bodys interface{}, exchangeName, queueName string) error {
	mq := NewRabbitMQ("amqp://guest:guest@10.68.41.31:5672")
	exc := Exchange{
		Name:       exchangeName,
		Kind:       ExcKindTopic,
		Durable:    true,
		AutoDelete: false,
		Internal:   false,
		NoWait:     false,
		Args:       nil,
	}
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
		err := mq.PublishMsg(&exc, queueName, &amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(msg),
		}, false, false)
		if err != nil {
			return err
		}
	}
	return nil
}
