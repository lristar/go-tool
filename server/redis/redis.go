package myredis

import (
	configs "platform-backend/config"
	utils "platform-backend/lib"
	logger "platform-backend/server/log"
	"strings"
	"time"

	// "github.com/go-redis/redis"
	"github.com/go-redis/redis/v7"
)

const (
	REDISKEYNIL = "redis: nil"
)

var (
	client *redis.Client
)

func init() {
	redisHost := configs.GetConf("Redishost")
	redisMaster := configs.GetConf("RedisMaster")
	password := configs.GetConf("RedisMasterPassword")
	sentinelPassword := configs.GetConf("RedisSentinelPassword")
	hosts := strings.Split(redisHost, ",")
	client = redis.NewFailoverClient(&redis.FailoverOptions{
		MasterName:       redisMaster,
		SentinelAddrs:    hosts,
		Password:         password,
		SentinelPassword: sentinelPassword,
		ReadTimeout:      15 * time.Second,
		WriteTimeout:     15 * time.Second,
		DB:               0,
		PoolSize:         10,
		OnConnect:        OnConnect,
	})
}

//OnConnect OnConnect
func OnConnect(conn *redis.Conn) error {
	_, err := conn.Ping().Result()
	if err != nil {
		logger.Error(err.Error())
	}
	return err
}

// RedisClient redis连接
type RedisClient struct {
	*redis.Client
	LockTime time.Duration
	lockKey  string
	LockMap  map[string]string
}

// GetRedisClient 获取redis连接
func GetRedisClient() *RedisClient {
	// _, err := client.Ping().Result()
	r := new(RedisClient)
	r.Client = client
	r.LockMap = make(map[string]string)
	return r
}

// Nil redis查询空值错误
func Nil() error {
	return redis.Nil
}

// Lock 通过redis加锁
func (r *RedisClient) Lock(id string) error {
	if r.LockTime == 0 {
		r.LockTime = 10 * time.Second
	}
	tc := time.NewTicker(100 * time.Millisecond) // 100毫秒尝试一次加锁
	for range tc.C {
		lockValue := utils.Md5Sum(utils.GetRandomString(12))
		f, err := r.Client.SetNX(id, lockValue, r.LockTime).Result()
		if err != nil {
			return err
		}
		if f {
			r.LockMap[id] = lockValue // 加锁成功，设置加锁lockValue
			return nil
		}
	}
	return nil
}

// SingleLock 通过redis尝试一次加锁,成功返回true,失败返回false
func (r *RedisClient) SingleLock(id string) (bool, error) {
	if r.LockTime == 0 {
		r.LockTime = 10 * time.Second
	}
	lockValue := utils.Md5Sum(utils.GetRandomString(12))
	f, err := r.Client.SetNX(id, lockValue, r.LockTime).Result()
	if err == nil && f {
		r.LockMap[id] = lockValue
	}
	return f, err
}

// Unlock 删除redis加锁（该方法不能判断是否解锁成功）
func (r *RedisClient) Unlock(id string) error {
	v, ok := r.LockMap[id]
	if !ok { // lockMap中找不到值，解锁失败，直接返回
		return nil
	}

	script := "if redis.call('get', KEYS[1]) == ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end"
	_, err := r.Eval(script, []string{id}, v).Result()
	return err
}

// Unlock 删除redis加锁（返回是否解锁成功）
func (r *RedisClient) UnlockWithFlag(id string) (bool, error) {
	flag := false
	v, ok := r.LockMap[id]
	if !ok { // lockMap中找不到值，解锁失败，直接返回
		return flag, nil
	}
	script := "if redis.call('get', KEYS[1]) == ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end"
	res, err := r.Eval(script, []string{id}, v).Result()
	if err != nil {
		return flag, err
	}

	if res == int64(1) {
		flag = true
	}
	return flag, nil
}

// Nil redis Nil error
func (r *RedisClient) Nil() error {
	return redis.Nil
}
