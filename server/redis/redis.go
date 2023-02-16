package myredis

import (
	"context"
	"github.com/go-redis/redis/v8"
	utils "gitlab.gf.com.cn/hk-common/go-tool/lib"
	"strings"
	"sync"
	"time"
)

const (
	REDISKEYNIL         = "redis: nil"
	ZERO                = 0
	DEFAULTTIMEOUT      = 15 * time.Second
	DEFAULTDB           = 0
	DEFAULTPOOLSIZE     = 10
	DEFAULTLOCKTIME     = 10 * time.Second
	DEFAULTLOOPLOCKTIME = 100 * time.Millisecond
)

var (
	_cli *RedisClient
	once sync.Once
)

type RedisConfig struct {
	Host         string `json:"host"`
	Master       string `json:"master"`
	MPassword    string `json:"mpassword"`
	SPassword    string `json:"spassword"`
	Group        string `json:"group"` // 分组的key
	ReadTimeout  int    `json:"readtimeout"`
	WriteTimeout int    `json:"write_timeout"`
	DB           int    `json:"db"`
	PoolSize     int    `json:"poolsize"`
}

// RedisClient redis连接
type RedisClient struct {
	*redis.Client
	prefixKey string
	LockTime  time.Duration
	lockKey   string
	LockMap   map[string]string
}

//OnConnect OnConnect
func OnConnect(ctx context.Context, cn *redis.Conn) error {
	_, err := cn.Ping(ctx).Result()
	return err
}

func InitRedisClient(config RedisConfig) {
	once.Do(func() {
		op := &redis.FailoverOptions{
			MasterName:       config.Master,
			SentinelAddrs:    strings.Split(config.Host, ","),
			Password:         config.MPassword,
			SentinelPassword: config.SPassword,
			ReadTimeout:      DEFAULTTIMEOUT,
			WriteTimeout:     DEFAULTTIMEOUT,
			DB:               DEFAULTDB,
			PoolSize:         DEFAULTPOOLSIZE,
			OnConnect:        OnConnect,
		}
		if config.ReadTimeout != ZERO {
			op.ReadTimeout = time.Second * time.Duration(config.ReadTimeout)
		}
		if config.WriteTimeout != ZERO {
			op.WriteTimeout = time.Second * time.Duration(config.WriteTimeout)
		}
		if config.PoolSize != ZERO {
			op.PoolSize = config.PoolSize
		}
		if config.DB != ZERO {
			op.PoolSize = config.PoolSize
		}
		client := redis.NewFailoverClient(op)
		client.AddHook(NewHook(config.Group))
		_cli = &RedisClient{
			Client:    client,
			prefixKey: config.Group,
			LockTime:  0,
			lockKey:   "",
			LockMap:   make(map[string]string),
		}
	})
}

func NewClient() *RedisClient {
	return &RedisClient{Client: _cli.Client, prefixKey: _cli.prefixKey, LockMap: make(map[string]string)}
}

// Lock 通过redis加锁
func (r *RedisClient) Lock(ctx context.Context, id string) error {
	if r.LockTime == ZERO {
		r.LockTime = DEFAULTLOCKTIME
	}
	tc := time.NewTicker(DEFAULTLOOPLOCKTIME) // 100毫秒尝试一次加锁
	for range tc.C {
		lockValue := utils.Md5Sum(utils.GetRandomString(12))
		f, err := r.Client.SetNX(ctx, id, lockValue, r.LockTime).Result()
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
func (r *RedisClient) SingleLock(ctx context.Context, id string) (bool, error) {
	if r.LockTime == ZERO {
		r.LockTime = DEFAULTLOCKTIME
	}
	lockValue := utils.Md5Sum(utils.GetRandomString(12))
	f, err := r.Client.SetNX(ctx, id, lockValue, r.LockTime).Result()
	if err == nil && f {
		r.LockMap[id] = lockValue
	}
	return f, err
}

// Unlock 删除redis加锁（该方法不能判断是否解锁成功）
func (r *RedisClient) Unlock(ctx context.Context, id string) error {
	v, ok := r.LockMap[id]
	if !ok { // lockMap中找不到值，解锁失败，直接返回
		return nil
	}

	script := "if redis.call('get', KEYS[1]) == ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end"
	_, err := r.Eval(ctx, script, []string{id}, v).Result()
	return err
}

// Unlock 删除redis加锁（返回是否解锁成功）
func (r *RedisClient) UnlockWithFlag(ctx context.Context, id string) (bool, error) {
	flag := false
	v, ok := r.LockMap[id]
	if !ok { // lockMap中找不到值，解锁失败，直接返回
		return flag, nil
	}
	script := "if redis.call('get', KEYS[1]) == ARGV[1] then return redis.call('del', KEYS[1]) else return 0 end"
	res, err := r.Eval(ctx, script, []string{id}, v).Result()
	if err != nil {
		return flag, err
	}

	if res == int64(1) {
		flag = true
	}
	return flag, nil
}

// Nil redis查询空值错误
func Nil() error {
	return redis.Nil
}
