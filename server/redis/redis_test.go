package myredis

import (
	"context"
	"fmt"
	"testing"
	"time"
)

var rdClient *RedisClient

func init() {
	op := RedisConfig{
		Host:         "10.68.41.32:24377,10.68.41.33:24377,10.68.41.36:24377",
		Master:       "mymaster",
		MPassword:    "gfhs123456",
		SPassword:    "gfhs123456",
		Group:        "testdemo",
		ReadTimeout:  0,
		WriteTimeout: 0,
		DB:           0,
		PoolSize:     0,
	}
	rdClient = NewRedisClient(op)
	rdClient.LockTime = 15 * time.Second
}

func TestNewRedisClient(t *testing.T) {
	isLock, err := rdClient.SingleLock(context.Background(), "hahahah")
	if err != nil {
		t.Fatal(err)
	}
	if !isLock {
		t.Fatal("加锁失败")
	}
	fmt.Println("加锁成功")
	time.Sleep(5 * time.Second)
	if _, err := rdClient.UnlockWithFlag(context.Background(), "hahahah"); err != nil {
		t.Fatal(err)
	}
	fmt.Println("解锁成功")
}

// AddHook 添加中间件
func TestAddHook(t *testing.T) {
	rdClient.AddHook(NewHook("demo"))
}
