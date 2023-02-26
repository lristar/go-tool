package myredis

import (
	"context"
	"github.com/go-redis/redis/v8"
)

func NewHook(prefixKey string) *RedisHook {
	return &RedisHook{prefixKey: prefixKey}
}

type RedisHook struct {
	prefixKey string
}

func (r *RedisHook) BeforeProcess(ctx context.Context, cmd redis.Cmder) (context.Context, error) {
	prefix := ""
	if r.prefixKey != "" {
		prefix = r.prefixKey + ":"
	}
	ags := cmd.Args()
	if v, ok := ags[0].(string); ok {
		if v != "eval" {
			ags[1] = prefix + ags[1].(string)
		} else {
			agsLen := ags[2].(int)
			if agsLen > 0 {
				ags[3] = prefix + ags[3].(string)
			}
		}
	}
	return ctx, nil
}

func (r *RedisHook) AfterProcess(ctx context.Context, cmd redis.Cmder) error {
	return nil
}

func (r *RedisHook) BeforeProcessPipeline(ctx context.Context, cmds []redis.Cmder) (context.Context, error) {
	return ctx, nil
}

func (r *RedisHook) AfterProcessPipeline(ctx context.Context, cmds []redis.Cmder) error {
	return nil
}
