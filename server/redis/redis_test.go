package myredis

import (
	"fmt"
	"gopkg.in/ffmt.v1"
	configs "platform-backend/config"
	"testing"
	"time"
)

func TestGetRedisClient(t *testing.T) {
	cli := GetRedisClient()
	result, err := cli.Set("b", "111", 1000*60*60*24).Result()
	fmt.Printf("%s %#v\n", result, err)
	result, err = cli.Get("b").Result()
	fmt.Printf("%s %#v\n", result, err)
}

func TestDel(t *testing.T) {
	cli := GetRedisClient()
	id := "testsuyifeng111"
	// v := "2341112"
	err := cli.Lock(id)
	if err != nil {
		t.Fatal(err)
	}

	cli.Lock(id)

	cli1 := GetRedisClient()
	err = cli1.Lock(id)
	if err != nil {
		t.Fatal(err)
	}

	err = cli1.Unlock(id)
	if err != nil {
		t.Fatal(err)
	}

	ffmt.Puts(err)

	f, err := cli1.UnlockWithFlag(id)
	if err != nil {
		t.Fatal(err)
	}

	ffmt.Puts(f)

	f, err = cli.UnlockWithFlag(id)
	if err != nil {
		t.Fatal(err)
	}

	ffmt.Puts(f)
}

func TestGetRedisClient2(t *testing.T) {
	cli := GetRedisClient()
	key := "token"
	result, err := cli.Set(key, "cabc033f-d91a-4450-a223-4b031e1092bc", 1000*60*60*24).Result()
	cli.Expire(key, 11*time.Second)
	fmt.Printf("%s %#v\n", result, err)
}

func TestTokenIsOver(t *testing.T) {
	cli := GetRedisClient()

	expire, err := cli.TTL("token").Result()
	over := expire.Milliseconds() / 1000
	if over < 1 {
		t.Fatal("timeout, please create")
	}
	fmt.Printf("expire=%d err =%v\n", over, err)

	result, err := cli.Get("token").Result()
	if err != nil {
		fmt.Println("获取失败 err =", err)
	}
	fmt.Println(result)
}

func TestRedisClient_LockOne(t *testing.T) {
	cli := GetRedisClient()
	err := cli.Lock("test1.aa.1111")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("加锁成功")
	defer func() {
		cli.Lock("test1.aa.1111")
		fmt.Println("加锁成功")
	}()
	err = cli.Unlock("test1.aa.1111")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("解锁成功")
}

func TestRedisClient_SingleLock(t *testing.T) {
	cli := GetRedisClient()
	f, err := cli.SingleLock("test1.aa.1111")
	if err != nil {
		t.Fatal()
	}
	if !f {
		fmt.Println("加锁失败")
		return
	}
	fmt.Println("加锁成功")
	f, err = cli.SingleLock("test1.aa.1111")
	if err != nil {
		t.Fatal()
	}
	if !f {
		fmt.Println("加锁失败")
		return
	}
}
func TestGetToken(t *testing.T) {
	cli := GetRedisClient()
	token, err := cli.Get(configs.HS_DICT_TOKEN_REDIS_KEY).Result()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(token)
}
