package myredis

import (
	"encoding/json"
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/go-redis/redis/v7"
	"platform-backend/server/sentry"
	"time"
)

type SortQueue struct {
	queueName    string // 队列名字
	redisCli     *redis.Client
	intervalTime time.Duration // 拉取数据间隔
}

func NewSortQueue(queueName string, intervalTime ...time.Duration) *SortQueue {
	t := time.Minute // 默认一分钟拉取一次
	if len(intervalTime) > 0 {
		t = intervalTime[0]
	}
	return &SortQueue{
		queueName:    queueName,
		redisCli:     client,
		intervalTime: t,
	}
}

func (s SortQueue) Publish(msg interface{}, score float64) error {
	z := redis.Z{
		Score:  score,
		Member: msg,
	}
	bt, _ := json.Marshal(msg)
	r, err := s.redisCli.ZAdd(s.queueName, &z).Result()
	if err != nil {
		sentry.LogAndSentry(err, &raven.Message{
			Message: string(bt),
			Params:  []interface{}{s.queueName},
		})
		return fmt.Errorf("%w", err)
	}
	if r != 1 {
		return fmt.Errorf("queueName:%s发布数据失败", s.queueName)
	}
	return nil
}

func (s SortQueue) Consume(fn func(msg string) error) error {
	tk := time.NewTicker(s.intervalTime)
	for range tk.C {
		zb := &redis.ZRangeBy{
			Min:    "0",
			Max:    fmt.Sprintf("%d", time.Now().Unix()),
			Offset: 0,
			Count:  0,
		}

		res, err := s.redisCli.ZRangeByScore(s.queueName, zb).Result()
		if err != nil {
			return err
		}
		for _, re := range res {
			if err = fn(re); err != nil {
				sentry.LogAndSentry(err, &raven.Message{
					Message: re,
					Params:  []interface{}{s.queueName},
				})
			}
			s.redisCli.ZRem(s.queueName, re)
		}
	}
	return nil
}
